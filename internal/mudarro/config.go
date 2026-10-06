package mudarro

import (
	"encoding/json"
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/projectfs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	configPath string    // Relative selected file; runtime metadata, excluded from serialization.
	Version    int       `yaml:"version" json:"version"`
	Name       string    `yaml:"name" json:"name"`
	Exclude    []string  `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	Services   []Service `yaml:"services" json:"services"`
	UI         *UIConfig `yaml:"ui,omitempty" json:"ui,omitempty"`
}

var identifier = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// safePath checks existing ancestors too: a symlink cannot redirect a generated file.
func safePath(root, rel string) (string, error) { return projectfs.SafePath(root, rel) }

func (c Config) Validate(root string) error {
	if err := c.validateUI(); err != nil {
		return err
	}
	if c.Version != 1 {
		return fmt.Errorf("version deve ser 1")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("name obrigatório")
	}
	seen := map[string]bool{}
	for _, s := range c.Services {
		if !identifier.MatchString(s.ID) || seen[s.ID] {
			return fmt.Errorf("id inválido ou duplicado: %q", s.ID)
		}
		seen[s.ID] = true
		if s.GoWorkspace != "" && (s.Language != "go" || (s.GoWorkspace != "inherit" && s.GoWorkspace != "off")) {
			return fmt.Errorf("%s: go_workspace requires a Go service and inherit or off", s.ID)
		}
		if s.WorkspaceRoot != "" {
			if s.Language != "javascript" && s.Language != "typescript" {
				return fmt.Errorf("%s: workspace_root requires JavaScript/TypeScript", s.ID)
			}
			if _, err := safePath(root, s.WorkspaceRoot); err != nil {
				return err
			}
			rel, err := filepath.Rel(s.WorkspaceRoot, s.Dir)
			if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return fmt.Errorf("%s: workspace_root must contain service directory", s.ID)
			}
		}
		p, e := safePath(root, s.Dir)
		if e != nil {
			return e
		}
		st, e := os.Stat(p)
		if e != nil || !st.IsDir() {
			return fmt.Errorf("diretório inexistente: %s", s.Dir)
		}
		for _, p := range []string{s.Infrastructure.File, s.Database.Path} {
			if p != "" {
				if _, e := safePath(root, filepath.Join(s.Dir, p)); e != nil {
					return e
				}
			}
		}
		if s.Infrastructure.Port < 0 || s.Infrastructure.Port > 65535 {
			return fmt.Errorf("porta inválida: %s", s.ID)
		}
		for n, cmd := range s.Commands {
			if !identifier.MatchString(n) {
				return fmt.Errorf("ação inválida: %s", n)
			}
			if (len(cmd.Args) > 0) == (cmd.Shell != "") {
				return fmt.Errorf("%s/%s: defina args OU shell", s.ID, n)
			}
			if len(cmd.Args) > 0 && strings.TrimSpace(cmd.Args[0]) == "" {
				return fmt.Errorf("executável vazio: %s", n)
			}
		}
		if s.Database.URLenv != "" && !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(s.Database.URLenv) {
			return fmt.Errorf("url_env inválida")
		}
	}
	return nil
}
func saveConfig(path string, c Config) error {
	b, e := yaml.Marshal(c)
	if e != nil {
		return e
	}
	f, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = f.Write(b)
	return e
}
func jsonString(v any) string    { b, _ := json.Marshal(v); return string(b) }
func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
