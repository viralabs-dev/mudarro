package mudarro_test

import (
	"crypto/sha256"
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Fixtures use the public data model; these helpers never reimplement production parsing or planning.
func cmd(group string, args ...string) Command { return Command{Args: args, Group: group} }
func action(name, group string, args ...string) Action {
	return Action{Name: name, Group: group, Command: cmd(group, args...)}
}
func exists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink == 0
}
func saveConfig(path string, c Config) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
func digest(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func findAction(c Config, key string) (Service, Action, error) {
	id, name, ok := strings.Cut(key, ":")
	if ok {
		for _, s := range c.Services {
			if s.ID == id {
				for _, a := range Actions(s) {
					if a.Name == name {
						return s, a, nil
					}
				}
			}
		}
	}
	return Service{}, Action{}, fmt.Errorf("action missing: %s", key)
}

// Inspect actual generated scaffolds through Generate, excluding operational launchers/ownership manifest.
func generatedScaffolds(t *testing.T, s Service) (map[string]GeneratedFile, error) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, s.Dir), 0755); err != nil {
		return nil, err
	}
	c := Config{Version: 1, Name: "scaffold fixture", Services: []Service{s}}
	if err := Generate(root, c, false, &strings.Builder{}); err != nil {
		return nil, err
	}
	files := map[string]GeneratedFile{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			if rel == ".mudarro" {
				return filepath.SkipDir
			}
			return nil
		}
		if rel == "menu.sh" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = GeneratedFile{Data: b, Mode: info.Mode()}
		return nil
	})
	return files, err
}

// skipOnWindows marks fixtures that depend on POSIX tools or Unix file
// semantics. Native Windows behaviour is covered by *_windows_test.go.
func skipOnWindows(t *testing.T, reason string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("Unix-only fixture: " + reason)
	}
}

// canonicalTempDir returns t.TempDir() resolved the way the runner records a
// project root (macOS: /var -> /private/var; Windows: 8.3 names -> long names).
func canonicalTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}
