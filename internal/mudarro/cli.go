package mudarro

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type options struct {
	root, name, selectScripts, config, format, locale, theme string
	json, dry, all, interactive                              bool
	exclude                                                  []string
	pos                                                      []string
}

func parseOptions(args []string) (options, error) {
	o := options{root: "."}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--json":
			o.json = true
		case "--dry-run":
			o.dry = true
		case "--all-scripts":
			o.all = true
		case "--interactive":
			o.interactive = true
		case "--root", "--name", "--select", "--exclude", "--config", "--format", "--locale", "--theme":
			if i+1 >= len(args) {
				return o, fmt.Errorf("valor ausente: %s", a)
			}
			i++
			switch a {
			case "--root":
				o.root = args[i]
			case "--name":
				o.name = args[i]
			case "--select":
				o.selectScripts = args[i]
			case "--config":
				o.config = args[i]
			case "--format":
				o.format = args[i]
			case "--locale":
				o.locale = args[i]
			case "--theme":
				o.theme = args[i]
			case "--exclude":
				o.exclude = append(o.exclude, args[i])
			}
		default:
			if strings.HasPrefix(a, "-") {
				return o, fmt.Errorf("opção desconhecida: %s", a)
			}
			o.pos = append(o.pos, a)
		}
	}
	root, e := filepath.Abs(o.root)
	if e != nil {
		return o, e
	}
	o.root, e = filepath.EvalSymlinks(root)
	return o, e
}
func Main(args []string, version string, in io.Reader, out io.Writer) (result error) {
	locale := "en"
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--locale" {
			locale = args[i+1]
		}
	}
	defer func() {
		if result != nil {
			result = localizedError{result, localizedMessage(result.Error(), locale)}
		}
	}()
	if len(args) == 0 {
		args = []string{"help"}
	}
	if args[0] == "__supervise" {
		return supervise(args[1:])
	}
	if args[0] == "version" || args[0] == "--version" {
		fmt.Fprintln(out, version)
		return nil
	}
	if args[0] == "help" || args[0] == "--help" {
		text := "Mudarro — operational menus without AI\n\n  scan [--json] [--exclude path]\n  init [--interactive | --select service:script,... | --all-scripts] [--format yaml|json]\n  generate [--dry-run]\n  menu\n  preview service:action\n  run service:action [--name name] [--dry-run]\n  doctor\n  version\n\nAll commands accept --root directory and --config file. Configure mudarro.yaml or mudarro.json before generating. UI flags: --locale en|pt-BR --theme auto|light|dark.\n"
		if locale == "pt-BR" {
			text = "Mudarro — menus operacionais sem IA\n\n  scan [--json] [--exclude caminho]\n  init [--interactive | --select serviço:script,... | --all-scripts] [--format yaml|json]\n  generate [--dry-run]\n  menu\n  preview serviço:ação\n  run serviço:ação [--name nome] [--dry-run]\n  doctor\n  version\n\nTodos os comandos aceitam --root diretório e --config arquivo. Configure mudarro.yaml ou mudarro.json antes de gerar. UI: --locale en|pt-BR --theme auto|light|dark.\n"
		}
		fmt.Fprint(out, text)
		return nil
	}

	o, e := parseOptions(args[1:])
	if e != nil {
		return e
	}
	switch args[0] {
	case "scan", "init":
		var scanUI *UIConfig
		if args[0] == "scan" && (o.config != "" || exists(filepath.Join(o.root, "mudarro.yaml")) || exists(filepath.Join(o.root, "mudarro.json"))) {
			c, e := loadOptions(o)
			if e != nil {
				return e
			}
			o.exclude = append(o.exclude, c.Exclude...)
			locale = c.UIOptions().Locale
			scanUI = c.UI
		}
		report, e := Scan(o.root, o.exclude)
		if e != nil {
			return e
		}
		report.Config.UI = scanUI
		if o.locale != "" || o.theme != "" {
			if report.Config.UI == nil {
				report.Config.UI = &UIConfig{}
			}
			if o.locale != "" {
				report.Config.UI.Locale = o.locale
			}
			if o.theme != "" {
				report.Config.UI.Theme = o.theme
			}
			if err := report.Config.validateUI(); err != nil {
				return err
			}
		}
		if o.name != "" {
			report.Config.Name = o.name
		}
		if args[0] == "scan" {
			if o.json {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}
			printReport(out, report)
			return nil
		}
		configName := o.config
		if configName == "" {
			configName = "mudarro.yaml"
			if o.format == "json" {
				configName = "mudarro.json"
			}
		}
		if o.format != "" && o.format != "yaml" && o.format != "json" {
			return fmt.Errorf("invalid --format: %s", o.format)
		}
		ext := strings.ToLower(filepath.Ext(configName))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return fmt.Errorf("unsupported configuration format: %s", ext)
		}
		if o.format == "json" && ext != ".json" || o.format == "yaml" && ext == ".json" {
			return fmt.Errorf("--format conflicts with --config extension")
		}
		target, err := safePath(o.root, configName)
		if err != nil {
			return err
		}
		if exists(target) || exists(filepath.Join(o.root, "mudarro.yaml")) || exists(filepath.Join(o.root, "mudarro.json")) {
			return fmt.Errorf("configuration already exists; edit the preserved file")
		}
		if o.locale != "" || o.theme != "" {
			report.Config.UI = &UIConfig{Locale: o.locale, Theme: o.theme}
		}
		selected := map[string]bool{}
		for _, v := range strings.Split(o.selectScripts, ",") {
			if v != "" {
				selected[v] = true
			}
		}
		reader := bufio.NewReader(in)
		available := map[string]bool{}
		for _, su := range report.Suggestions {
			key := su.Service + ":" + su.Name
			available[key] = true
			include := o.all || selected[key]
			if o.interactive {
				fmt.Fprintf(out, uiText(report.Config, "Include %s (%s)? [y/N] ", "Incluir %s (%s)? [s/N] "), key, su.Evidence)
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("entrada encerrada antes de salvar configuração")
				}
				answer := strings.TrimSpace(line)
				include = strings.EqualFold(answer, uiText(report.Config, "y", "s"))
			}
			if include {
				for j := range report.Config.Services {
					if report.Config.Services[j].ID == su.Service {
						if len(su.Command.Args) > 0 && su.Command.Args[0] == "" {
							return fmt.Errorf("defina o gerenciador antes de incluir %s", key)
						}
						report.Config.Services[j].Commands[su.Name] = su.Command
					}
				}
			}
		}
		for key := range selected {
			if !available[key] {
				return fmt.Errorf("script sugerido inexistente: %s", key)
			}
		}
		if e = saveSelectedConfig(o.root, target, report.Config); e != nil {
			return e
		}
		printReport(out, report)
		fmt.Fprintf(out, "Created %s. Review pending choices, then run mudarro generate.\n", configName)
		return nil
	case "generate", "menu", "run", "doctor", "preview":
		c, e := loadOptions(o)
		if e != nil {
			return e
		}
		locale = c.UIOptions().Locale
		switch args[0] {
		case "generate":
			return Generate(o.root, c, o.dry, out)
		case "menu":
			return menu(o.root, c, in, out)
		case "doctor":
			return doctor(o.root, c, out)
		case "run", "preview":
			if len(o.pos) != 1 {
				return fmt.Errorf("uso: mudarro run serviço:ação")
			}
			s, a, e := findAction(c, o.pos[0])
			if e != nil {
				return e
			}
			if args[0] == "preview" {
				fmt.Fprint(out, actionPreview(o.root, s, a, c.UIOptions().Locale))
				return nil
			}
			return (Runner{In: in, Out: out, Dry: o.dry, Name: o.name, Locale: c.UIOptions().Locale, ConfigPath: c.configPath}).Run(o.root, s, a)
		}
	}
	return fmt.Errorf("comando desconhecido: %s", args[0])
}
func printReport(w io.Writer, r Report) {
	fmt.Fprintln(w, uiText(r.Config, "Project:", "Projeto:"), r.Config.Name)
	for _, s := range r.Config.Services {
		fmt.Fprintf(w, "- %s (%s): %s / %s / %s\n", s.ID, s.Dir, s.Language, s.Manager, s.Infrastructure.Kind)
		for _, p := range s.Pending {
			fmt.Fprintln(w, uiText(r.Config, "  pending:", "  pendente:"), localizedMessage(p, r.Config.UIOptions().Locale))
		}
	}
	for _, su := range r.Suggestions {
		fmt.Fprintf(w, uiText(r.Config, "  suggestion: %s:%s [%s]\n", "  sugestão: %s:%s [%s]\n"), su.Service, su.Name, su.Evidence)
	}
	for _, warning := range r.Warnings {
		fmt.Fprintln(w, uiText(r.Config, "warning:", "aviso:"), localizedMessage(warning, r.Config.UIOptions().Locale))
	}
}
func findAction(c Config, key string) (Service, Action, error) {
	parts := strings.SplitN(key, ":", 2)
	var service, name string
	if len(parts) == 2 {
		service, name = parts[0], parts[1]
	} else {
		if len(c.Services) != 1 {
			return Service{}, Action{}, fmt.Errorf("informe serviço:ação em projetos com vários serviços")
		}
		service, name = c.Services[0].ID, key
	}
	for _, s := range c.Services {
		if s.ID == service {
			for _, a := range Actions(s) {
				if a.Name == name {
					return s, a, nil
				}
			}
		}
	}
	return Service{}, Action{}, fmt.Errorf("ação inexistente: %s", key)
}
func menu(root string, c Config, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	u := c.UIOptions()
	view := terminal.NewMenu(out, c.Name, fmt.Sprintf(uiText(c, "%d service(s) · offline detection · no AI", "%d serviço(s) · detecção offline · sem IA"), len(c.Services)))
	view.Configure(terminal.Presentation{Locale: u.Locale, Theme: u.Theme, Density: u.Density, Lettering: u.Lettering, Mouse: u.Preview.Mouse})
	if f, ok := in.(*os.File); ok {
		if _, err := view.BeginShell(f); err != nil {
			return err
		}
		defer view.CloseShell()
	}
	style := view.Style()
	for {
		labels := []string{}
		for _, s := range c.Services {
			labels = append(labels, s.ID+" — "+s.Dir)
		}
		view.SetContext(fmt.Sprintf(uiText(c, "%d service(s) / offline detection", "%d serviço(s) / detecção offline"), len(c.Services)))
		n, e := view.Choose(reader, uiText(c, "Services", "Serviços"), labels)
		if e == io.EOF || n == 0 {
			return nil
		}
		if e != nil {
			return e
		}
		s := c.Services[n-1]
		view.SetContext(s.ID + " / " + s.Language + " / " + s.Infrastructure.Kind)
		for {
			actions := Actions(s)
			groups := []string{}
			byGroup := map[string][]Action{}
			for _, a := range actions {
				if _, ok := byGroup[a.Group]; !ok {
					groups = append(groups, a.Group)
				}
				byGroup[a.Group] = append(byGroup[a.Group], a)
			}
			sort.Strings(groups)
			n, e = view.Choose(reader, s.ID, groupLabels(c, groups))
			if e == io.EOF {
				return nil
			}
			if e != nil {
				return e
			}
			if n == 0 {
				break
			}
			group := groups[n-1]
			for {
				list := byGroup[group]
				labels = nil
				for _, a := range list {
					label := a.Name
					if a.Blocked != "" {
						label += " [" + uiText(c, "pending: ", "pendente: ") + a.Blocked + "]"
					}
					labels = append(labels, label)
				}
				n, e = chooseActions(view, reader, in, root, c, s, group, list, labels)
				if e == io.EOF {
					return nil
				}
				if e != nil {
					return e
				}
				if n == 0 {
					break
				}
				a := list[n-1]
				if view.ShellActive() {
					if err := runShellAction(view, root, c, s, a); err != nil {
						return err
					}
					continue
				}
				name := ""
				if strings.Contains(strings.Join(a.Command.Args, " "), "{name}") {
					fmt.Fprint(out, uiText(c, "Migration name: ", "Nome da migration: "))
					line, e := reader.ReadString('\n')
					if e != nil {
						return nil
					}
					name = strings.TrimSpace(line)
				}
				if e = (Runner{In: reader, Out: out, Name: name, Locale: c.UIOptions().Locale, ConfigPath: c.configPath}).Run(root, s, a); e != nil {
					fmt.Fprintln(out, style.Paint("1;31", uiText(c, "Error:", "Erro:")), e)
				}
				if style.HasColor() {
					fmt.Fprint(out, style.Hint(uiText(c, "\n  Enter to return to the menu...", "\n  Enter para voltar ao menu...")))
					if _, err := reader.ReadString('\n'); err != nil {
						return nil
					}
				}
			}
		}
	}
}
func doctor(root string, c Config, w io.Writer) error {
	failures := 0
	for _, s := range c.Services {
		dir, _ := serviceDirectory(root, s)
		seen := map[string]bool{}
		for _, a := range Actions(s) {
			if a.Blocked != "" {
				fmt.Fprintf(w, "PENDENTE %s:%s — %s\n", s.ID, a.Name, a.Blocked)
				failures++
				continue
			}
			tools := append([]string{}, a.Command.Requires...)
			if a.Command.Shell != "" {
				tools = append(tools, "bash")
			}
			if len(a.Command.Args) > 0 {
				tools = append(tools, a.Command.Args[0])
			}
			if strings.HasPrefix(a.Special, "local-") {
				tools = append(tools, "ps")
			}
			if strings.HasPrefix(a.Special, "k8s-") {
				tools = append(tools, "kubectl")
			}
			for _, tool := range tools {
				if seen[tool] {
					continue
				}
				seen[tool] = true
				var e error
				if strings.Contains(tool, "/") {
					var p string
					p, e = safePath(root, filepath.Join(s.Dir, tool))
					if e == nil {
						var st os.FileInfo
						st, e = os.Stat(p)
						if e == nil && (st.IsDir() || st.Mode()&0111 == 0) {
							e = fmt.Errorf("não executável")
						}
					}
				} else {
					_, e = exec.LookPath(tool)
				}
				if e != nil {
					fmt.Fprintf(w, "AUSENTE %s: %s — configure ou instale explicitamente\n", s.ID, tool)
					failures++
				} else {
					fmt.Fprintf(w, "OK %s: %s\n", s.ID, tool)
				}
			}
		}
		if s.Infrastructure.Mode == "compose" && (s.Infrastructure.Kind == "docker" || s.Infrastructure.Kind == "podman") {
			provider := s.Infrastructure.Kind
			args := []string{"compose", "version"}
			if provider == "podman" {
				provider = "podman-compose"
				args = []string{"--version"}
			}
			check := exec.Command(provider, args...)
			check.Dir = dir
			if e := check.Run(); e != nil {
				fmt.Fprintf(w, "AUSENTE %s: provedor Compose\n", s.ID)
				failures++
			}
		}
	}
	if failures > 0 {
		return fmt.Errorf("%d pendência(s); as ações independentes continuam disponíveis", failures)
	}
	fmt.Fprintln(w, "Nenhuma dependência executável ausente. Isso não valida conexão, credenciais ou saúde da aplicação.")
	return nil
}
