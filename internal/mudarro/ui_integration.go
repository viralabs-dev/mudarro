package mudarro

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	"github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"gopkg.in/yaml.v3"
)

func loadOptions(o options) (Config, error) {
	var c Config
	var err error
	if o.config != "" {
		c, err = LoadFile(o.root, o.config)
	} else {
		c, err = Load(o.root)
	}
	if err != nil {
		return c, err
	}
	if o.locale != "" || o.theme != "" {
		if c.UI == nil {
			c.UI = &UIConfig{}
		}
		if o.locale != "" {
			c.UI.Locale = o.locale
		}
		if o.theme != "" {
			c.UI.Theme = o.theme
		}
	}
	return c, c.Validate(o.root)
}
func saveSelectedConfig(root, path string, c Config) error {
	if err := c.Validate(root); err != nil {
		return err
	}
	var b []byte
	var err error
	if strings.ToLower(filepath.Ext(path)) == ".json" {
		b, err = json.MarshalIndent(c, "", "  ")
		b = append(b, '\n')
	} else {
		b, err = yaml.Marshal(c)
	}
	if err != nil {
		return err
	}
	return atomicWrite(path, b, 0644)
}
func uiText(c Config, en, pt string) string {
	if c.UIOptions().Locale == "pt-BR" {
		return pt
	}
	return en
}
func groupLabels(c Config, groups []string) []string {
	if c.UIOptions().Locale == "pt-BR" {
		return groups
	}
	names := map[string]string{"aplicacao": "application", "infraestrutura": "infrastructure", "banco": "database", "qualidade": "quality", "dependencias": "dependencies", "scripts": "scripts"}
	result := make([]string, len(groups))
	for i, g := range groups {
		result[i] = g
		if name, ok := names[g]; ok {
			result[i] = name
		}
	}
	return result
}
func chooseActions(v *terminal.Menu, r *bufio.Reader, in io.Reader, root string, c Config, s Service, group string, list []Action, labels []string) (int, error) {
	previews := []string{}
	if c.UIOptions().Preview.Enabled {
		for _, a := range list {
			previews = append(previews, actionPreview(root, s, a, c.UIOptions().Locale))
		}
	}
	v.SetPreviews(previews)
	defer v.SetPreviews(nil)
	title := groupLabels(c, []string{group})[0]
	if len(previews) > 0 && r.Buffered() == 0 {
		if f, ok := in.(*os.File); ok {
			n, supported, err := v.ChooseInteractive(f, title, labels, previews)
			if supported || err != nil {
				return n, err
			}
		}
	}
	return v.Choose(r, title, labels)
}

var sensitiveName = regexp.MustCompile(`(?i)(password|passwd|secret|token|api[_-]?key|credential)`)
var sensitiveFlag = regexp.MustCompile(`(?i)(--?(?:password|passwd|secret|token|api[_-]?key|credential)(?:[=\t ]+))('[^']*'|"[^"]*"|[^\s;]+)`)
var sensitiveAssignment = regexp.MustCompile(`(?i)((?:password|passwd|secret|token|api[_-]?key)\s*=\s*)('[^']*'|"[^"]*"|[^\s;]+)`)

var embeddedURI = regexp.MustCompile(`(?i)[a-z][a-z0-9+.-]*://[^\s"'<>]+`)
var sensitiveHeader = regexp.MustCompile(`(?i)(authorization\s*[:=]\s*(?:bearer|basic)\s+)[^\s"'<>]+`)

func redactArg(s string) string {
	s = embeddedURI.ReplaceAllStringFunc(s, func(raw string) string {
		u, err := url.Parse(raw)
		if err != nil {
			return "[REDACTED URI]"
		}
		if u.User != nil {
			if _, ok := u.User.Password(); ok {
				u.User = url.UserPassword(u.User.Username(), "[REDACTED]")
			}
		}
		q := u.Query()
		for key := range q {
			if sensitiveName.MatchString(key) {
				q.Set(key, "[REDACTED]")
			}
		}
		if u.RawQuery != "" {
			u.RawQuery = q.Encode()
		}
		return u.String()
	})
	if key, _, ok := strings.Cut(s, "="); ok && sensitiveName.MatchString(key) && !strings.ContainsAny(key, " \n") {
		return key + "=[REDACTED]"
	}
	s = sensitiveHeader.ReplaceAllString(s, "${1}[REDACTED]")
	s = sensitiveFlag.ReplaceAllString(s, "${1}[REDACTED]")
	return sensitiveAssignment.ReplaceAllString(s, "${1}[REDACTED]")
}
func actionPreview(root string, s Service, a Action, locale string) string {
	text := func(en, pt string) string {
		if locale == "pt-BR" {
			return pt
		}
		return en
	}
	var b strings.Builder
	fmt.Fprintf(&b, text("Action: %s:%s\nDirectory: %s\n", "Ação: %s:%s\nDiretório: %s\n"), s.ID, a.Name, s.Dir)
	if a.Blocked != "" {
		fmt.Fprintf(&b, text("Pending: %s\n", "Pendente: %s\n"), a.Blocked)
	}
	if a.Command.Destructive {
		b.WriteString(text("Destructive: explicit confirmation is required\n", "Destrutiva: confirmação explícita obrigatória\n"))
	}
	if a.Special != "" {
		fmt.Fprintf(&b, text("Planned operation: %s\nRuntime process/container/workload values are resolved only on execution.\n", "Operação planejada: %s\nValores de processo/container/workload são resolvidos somente na execução.\n"), a.Special)
	}
	if a.Command.Shell != "" {
		b.WriteString(text("Explicit shell (not evaluated):\n", "Shell explícito (sem avaliar):\n") + redactArg(a.Command.Shell) + "\n")
	}
	if len(a.Command.Args) > 0 {
		b.WriteString(text("argv (not shell evaluation):\n", "argv (sem avaliação de shell):\n"))
		hide := false
		for i, arg := range a.Command.Args {
			if b.Len() >= 16384 {
				b.WriteString("[preview budget reached]\n")
				break
			}
			shown := redactArg(arg)
			if hide {
				shown = "[REDACTED]"
			}
			hide = strings.HasPrefix(arg, "-") && sensitiveName.MatchString(arg) && !strings.Contains(arg, "=")
			fmt.Fprintf(&b, "  [%d] %s\n", i, strconv.Quote(shown))
		}
	}
	// Show only bounded files explicitly named in argv; never follow symlinks/source/eval.
	seenSources := map[string]bool{}
	for _, arg := range a.Command.Args {
		if b.Len() >= 16384 {
			b.WriteString("[preview budget reached]\n")
			break
		}
		ext := strings.ToLower(filepath.Ext(arg))
		if ext != ".sh" && ext != ".py" && ext != ".js" && ext != ".go" {
			continue
		}
		path, err := safePath(root, filepath.Join(s.Dir, arg))
		if err != nil {
			continue
		}
		if seenSources[path] {
			continue
		}
		seenSources[path] = true
		initial, err := os.Stat(path)
		if err != nil || !initial.Mode().IsRegular() {
			continue
		}
		f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
		if err != nil {
			continue
		}
		st, err := f.Stat()
		if err != nil || !st.Mode().IsRegular() {
			f.Close()
			continue
		}
		bytes, err := io.ReadAll(io.LimitReader(f, 8193))
		f.Close()
		if err != nil {
			continue
		}
		capped := len(bytes) > 8192
		if capped {
			bytes = bytes[:8192]
		}
		fmt.Fprintf(&b, text("Source %s (read-only):\n%s\n", "Fonte %s (somente leitura):\n%s\n"), arg, redactArg(string(bytes)))
		if capped {
			b.WriteString("[preview truncated at 8 KiB]\n")
		}
	}
	b.WriteString(text("Environment references are not expanded. Embedded arbitrary shell secrets may require manual redaction.\n", "Referências ao ambiente não são expandidas. Segredos arbitrários em shell podem exigir ocultação manual.\n"))
	result := strings.Map(func(r rune) rune {
		if (unicode.IsControl(r) && r != '\n' && r != '\t') || unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, b.String())
	if len(result) > 16384 {
		result = strings.ToValidUTF8(result[:16384], "�") + "\n[preview truncated at 16 KiB]\n"
	}
	return result
}
