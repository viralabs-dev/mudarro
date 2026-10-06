package mudarro_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func supportFileHashes(t *testing.T, root string) map[string][32]byte {
	t.Helper()
	result := map[string][32]byte{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[rel] = sha256.Sum256(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func TestSupportGenerationPipelineByManager(t *testing.T) {
	cases := []struct {
		name     string
		files    map[string]string
		testArgs []string
	}{
		{"npm", map[string]string{"package.json": `{"scripts":{"start":"NEVER_EXECUTE","test":"NEVER_EXECUTE"}}`, "package-lock.json": "{}"}, []string{"npm", "run", "test"}},
		{"pnpm", map[string]string{"package.json": `{"scripts":{"start":"NEVER_EXECUTE","test":"NEVER_EXECUTE"}}`, "pnpm-lock.yaml": ""}, []string{"pnpm", "run", "test"}},
		{"yarn", map[string]string{"package.json": `{"scripts":{"start":"NEVER_EXECUTE","test":"NEVER_EXECUTE"}}`, "yarn.lock": ""}, []string{"yarn", "run", "test"}},
		{"pip", map[string]string{"requirements.txt": "django", "manage.py": "raise Exception('NEVER_EXECUTE')"}, []string{".venv/bin/python", "manage.py", "test"}},
		{"uv", map[string]string{"requirements.txt": "django", "manage.py": "raise Exception('NEVER_EXECUTE')", "uv.lock": ""}, []string{"uv", "run", "python", "manage.py", "test"}},
		{"poetry", map[string]string{"requirements.txt": "django", "manage.py": "raise Exception('NEVER_EXECUTE')", "poetry.lock": ""}, []string{"poetry", "run", "python", "manage.py", "test"}},
		{"go", map[string]string{"go.mod": "module fixture\ngo 1.23\n", "main.go": "package main\nfunc main(){}"}, []string{"go", "test", "./..."}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, tc.files)
			before := supportFileHashes(t, root)
			var out bytes.Buffer
			invoke := func(command string, extra ...string) error {
				out.Reset()
				args := append([]string{command, "--root", root}, extra...)
				return Main(args, "test", strings.NewReader(""), &out)
			}
			if err := invoke("scan", "--json"); err != nil {
				t.Fatal(err)
			}
			var report Report
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if len(report.Config.Services) != 1 {
				t.Fatalf("%+v", report)
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("scan changed source")
			}
			if err := invoke("init"); err != nil {
				t.Fatal(err)
			}
			initialized := supportFileHashes(t, root)
			if err := invoke("init"); err == nil {
				t.Fatal("second init overwrote config")
			}
			if !reflect.DeepEqual(initialized, supportFileHashes(t, root)) {
				t.Fatal("second init changed files")
			}
			if err := invoke("generate"); err != nil {
				t.Fatal(err)
			}
			generated := supportFileHashes(t, root)
			if err := invoke("generate"); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(generated, supportFileHashes(t, root)) {
				t.Fatal("generation not idempotent")
			}
			id := report.Config.Services[0].ID
			for _, name := range []string{"menu.sh", filepath.Join(".mudarro", "scripts", id, "test.sh")} {
				info, err := os.Stat(filepath.Join(root, name))
				if err != nil || info.Mode()&0111 == 0 {
					t.Fatalf("launcher missing/nonexecutable %s %v", name, err)
				}
			}
			// Choose service and first group, return to services, then quit. Never choose an action.
			out.Reset()
			if err := Main([]string{"menu", "--root", root}, "test", strings.NewReader("1\n1\n0\n0\n0\n"), &out); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), "Services") || !strings.Contains(out.String(), id) || !strings.Contains(out.String(), "Back / exit") {
				t.Fatal("menu hierarchy missing")
			}
			if err := invoke("run", id+":test", "--dry-run"); err != nil {
				t.Fatal(err)
			}
			_, commandJSON, ok := strings.Cut(strings.TrimSpace(out.String()), ": ")
			var planned Command
			if !ok {
				t.Fatal("missing plan JSON", out.String())
			}
			if err := json.Unmarshal([]byte(commandJSON), &planned); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(planned.Args, tc.testArgs) {
				t.Fatalf("argv got %v want %v", planned.Args, tc.testArgs)
			}
			if !reflect.DeepEqual(generated, supportFileHashes(t, root)) {
				t.Fatal("menu/dry-run changed files")
			}
		})
	}
}
