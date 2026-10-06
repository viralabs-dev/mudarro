package adapters

import (
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"github.com/viralabs-dev/mudarro/internal/mudarro/projectfs"
	"os"
	"regexp"
	"strings"
)

var identifier = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

func serviceID(dir, lang string) string {
	v := regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(dir, "-")
	v = strings.Trim(v, "-")
	if v == "" {
		v = "app"
	}
	return v + "-" + lang
}
func BaseService(dir, lang string) *model.Service {
	return &model.Service{ID: serviceID(dir, lang), Dir: dir, Language: lang, Commands: map[string]model.Command{}}
}
func cmd(group string, args ...string) model.Command { return model.Command{Args: args, Group: group} }

var ignored = map[string]bool{".git": true, "node_modules": true, "vendor": true, ".venv": true, "venv": true, "dist": true, "build": true, ".mudarro": true, "__pycache__": true, ".next": true, ".tox": true, ".cache": true}

func IgnoredDirectory(name string) bool           { return ignored[name] }
func exists(p string) bool                        { st, e := os.Lstat(p); return e == nil && st.Mode()&os.ModeSymlink == 0 }
func read(root, dir, file string) ([]byte, error) { return projectfs.ReadManifest(root, dir, file) }
func action(name, group string, args ...string) model.Action {
	return model.Action{Name: name, Group: group, Command: cmd(group, args...)}
}
func blocked(name, group, reason string) model.Action {
	return model.Action{Name: name, Group: group, Blocked: reason}
}
