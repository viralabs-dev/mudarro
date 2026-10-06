package mudarro

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Command struct {
	Args        []string `yaml:"args,omitempty" json:"args,omitempty"`
	Shell       string   `yaml:"shell,omitempty" json:"shell,omitempty"`
	Group       string   `yaml:"group,omitempty" json:"group,omitempty"`
	Destructive bool     `yaml:"destructive,omitempty" json:"destructive,omitempty"`
	Requires    []string `yaml:"requires,omitempty" json:"requires,omitempty"`
}
type Infrastructure struct {
	Kind      string `yaml:"kind,omitempty" json:"kind,omitempty"`
	File      string `yaml:"file,omitempty" json:"file,omitempty"`
	Mode      string `yaml:"mode,omitempty" json:"mode,omitempty"`
	Context   string `yaml:"context,omitempty" json:"context,omitempty"`
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Image     string `yaml:"image,omitempty" json:"image,omitempty"`
	Port      int    `yaml:"port,omitempty" json:"port,omitempty"`
	Generate  bool   `yaml:"generate,omitempty" json:"generate,omitempty"`
}
type Database struct {
	Kind     string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Tool     string `yaml:"tool,omitempty" json:"tool,omitempty"`
	Generate bool   `yaml:"generate,omitempty" json:"generate,omitempty"`
	URLenv   string `yaml:"url_env,omitempty" json:"url_env,omitempty"`
	Path     string `yaml:"path,omitempty" json:"path,omitempty"`
}
type Service struct {
	ID             string             `yaml:"id" json:"id"`
	Dir            string             `yaml:"dir" json:"dir"`
	Language       string             `yaml:"language,omitempty" json:"language,omitempty"`
	Manager        string             `yaml:"manager,omitempty" json:"manager,omitempty"`
	Framework      string             `yaml:"framework,omitempty" json:"framework,omitempty"`
	Infrastructure Infrastructure     `yaml:"infrastructure" json:"infrastructure"`
	Database       Database           `yaml:"database,omitempty" json:"database,omitempty"`
	Commands       map[string]Command `yaml:"commands,omitempty" json:"commands,omitempty"`
	Pending        []string           `yaml:"pending,omitempty" json:"pending,omitempty"`
}
type Config struct {
	Version  int       `yaml:"version" json:"version"`
	Name     string    `yaml:"name" json:"name"`
	Exclude  []string  `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	Services []Service `yaml:"services" json:"services"`
}

var identifier = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// safePath checks existing ancestors too: a symlink cannot redirect a generated file.
func safePath(root, rel string) (string, error) {
	if filepath.IsAbs(rel) || rel == "" {
		return "", fmt.Errorf("caminho relativo inválido: %q", rel)
	}
	clean := filepath.Clean(rel)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("caminho fora do projeto: %q", rel)
	}
	current := root
	for _, part := range strings.Split(clean, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		st, err := os.Lstat(current)
		if err == nil && st.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symlink não permitido: %s", current)
		}
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}
	return filepath.Join(root, clean), nil
}

func Load(root string) (Config, error) {
	var c Config
	path, err := safePath(root, "mudarro.yaml")
	if err != nil {
		return c, err
	}
	f, err := os.Open(path)
	if err != nil {
		return c, err
	}
	defer f.Close()
	d := yaml.NewDecoder(io.LimitReader(f, 4<<20))
	d.KnownFields(true)
	if err = d.Decode(&c); err != nil {
		return c, fmt.Errorf("mudarro.yaml: %w", err)
	}
	var extra any
	if err = d.Decode(&extra); err != io.EOF {
		return c, fmt.Errorf("use um único documento YAML")
	}
	return c, c.Validate(root)
}
func (c Config) Validate(root string) error {
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
