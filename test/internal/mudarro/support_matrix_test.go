package mudarro_test

import (
	"bytes"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSupportJavascriptManagers(t *testing.T) {
	for _, tc := range []struct{ name, lock, explicit, manager, language string }{
		{"default", "", "", "npm", "javascript"},
		{"npm", "package-lock.json", "", "npm", "javascript"},
		{"yarn", "yarn.lock", "", "yarn", "javascript"},
		{"pnpm-ts", "pnpm-lock.yaml", "", "pnpm", "typescript"},
		{"explicit-wins", "yarn.lock", "pnpm@9.0.0", "pnpm", "javascript"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			manifest := map[string]any{"name": "matrix", "scripts": map[string]string{"dev": "NEVER_EXECUTE", "test:unit": "NEVER_EXECUTE", "custom": "NEVER_EXECUTE"}}
			if tc.explicit != "" {
				manifest["packageManager"] = tc.explicit
			}
			b, _ := json.Marshal(manifest)
			files := map[string]string{"package.json": string(b)}
			if tc.lock != "" {
				files[tc.lock] = ""
			}
			if tc.language == "typescript" {
				files["tsconfig.json"] = "{}"
			}
			r, err := Scan(fixture(t, files), nil)
			if err != nil {
				t.Fatal(err)
			}
			s := r.Config.Services[0]
			if s.Manager != tc.manager || s.Language != tc.language {
				t.Fatalf("%+v", s)
			}
			for name, want := range map[string][]string{"start": {tc.manager, "run", "dev"}, "test-unit": {tc.manager, "run", "test:unit"}, "install": {tc.manager, "install"}} {
				if !reflect.DeepEqual(s.Commands[name].Args, want) {
					t.Fatalf("%s: %+v", name, s.Commands[name])
				}
			}
			if _, ok := s.Commands["custom"]; ok {
				t.Fatal("custom script selected implicitly")
			}
		})
	}
}

func TestSupportPythonManagers(t *testing.T) {
	for _, tc := range []struct {
		manager, lock  string
		install, start []string
	}{
		{"pip", "", []string{".venv/bin/python", "-m", "pip", "install", "-r", "requirements.txt"}, []string{".venv/bin/python", "manage.py", "runserver"}},
		{"uv", "uv.lock", []string{"uv", "sync"}, []string{"uv", "run", "python", "manage.py", "runserver"}},
		{"poetry", "poetry.lock", []string{"poetry", "install"}, []string{"poetry", "run", "python", "manage.py", "runserver"}},
	} {
		t.Run(tc.manager, func(t *testing.T) {
			files := map[string]string{"requirements.txt": "django", "manage.py": "raise RuntimeError('NEVER_EXECUTE')"}
			if tc.lock != "" {
				files[tc.lock] = ""
			}
			r, err := Scan(fixture(t, files), nil)
			if err != nil {
				t.Fatal(err)
			}
			s := r.Config.Services[0]
			if s.Manager != tc.manager || s.Framework != "django" || !reflect.DeepEqual(s.Commands["install"].Args, tc.install) || !reflect.DeepEqual(s.Commands["start"].Args, tc.start) {
				t.Fatalf("%+v", s)
			}
		})
	}
	t.Run("ambiguous-locks", func(t *testing.T) {
		r, err := Scan(fixture(t, map[string]string{"requirements.txt": "", "uv.lock": "", "poetry.lock": ""}), nil)
		if err != nil {
			t.Fatal(err)
		}
		s := r.Config.Services[0]
		if s.Manager != "" || len(s.Pending) == 0 {
			t.Fatalf("%+v", s)
		}
		if _, ok := s.Commands["install"]; ok {
			t.Fatal("installation inferred with ambiguous manager")
		}
	})
}

func TestSupportGoWorkspaceBoundaries(t *testing.T) {
	root := fixture(t, map[string]string{"go.work": "go 1.23\nuse (\n ./api\n ./worker\n)\n", "api/go.mod": "module api\ngo 1.23\n", "api/cmd/server/main.go": "package main\nfunc main(){}", "api/nested/go.mod": "module nested\ngo 1.23\n", "api/nested/main.go": "package main", "worker/go.mod": "module worker\ngo 1.23\n", "worker/one/main.go": "package main", "worker/two/main.go": "package main", "worker/example_test.go": "package main"})
	r, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Config.Services) != 3 {
		t.Fatalf("services: %+v", r.Config.Services)
	}
	workspace := false
	for _, e := range r.Evidence {
		if e.Kind == "go-workspace" {
			workspace = true
		}
	}
	if !workspace {
		t.Fatal("workspace evidence missing")
	}
	for _, s := range r.Config.Services {
		if s.Dir == "api" && !reflect.DeepEqual(s.Commands["start"].Args, []string{"go", "run", "./cmd/server"}) {
			t.Fatalf("nested module contaminated parent: %+v", s)
		}
		if s.Dir == "worker" {
			if _, ok := s.Commands["start"]; ok {
				t.Fatal("ambiguous entry selected")
			}
			if len(s.Pending) == 0 {
				t.Fatal("ambiguity missing")
			}
		}
	}
}

func TestSupportAmbiguousDatabaseAndInfrastructure(t *testing.T) {
	r, err := Scan(fixture(t, map[string]string{"package.json": "{}", "Chart.yaml": "apiVersion: v2\nname: x\nversion: 1.0.0", "kustomization.yaml": "resources: []", "prisma/schema.prisma": "datasource db { provider = \"sqlite\" }", "alembic.ini": "[alembic]"}), nil)
	if err != nil {
		t.Fatal(err)
	}
	s := r.Config.Services[0]
	if s.Infrastructure.Kind != "" || s.Database.Tool != "" {
		t.Fatalf("ambiguity resolved silently: %+v", s)
	}
	pending := strings.Join(s.Pending, "\n")
	if !strings.Contains(pending, "Múltiplas infraestruturas") || !strings.Contains(pending, "Múltiplas ferramentas") {
		t.Fatal(pending)
	}
}

func TestSupportCLIRunDryMatrix(t *testing.T) {
	for _, tool := range []string{"prisma", "django", "alembic", "goose"} {
		t.Run(tool, func(t *testing.T) {
			root := t.TempDir()
			c := sample()
			s := &c.Services[0]
			s.Database = Database{Kind: "postgresql", Tool: tool, URLenv: "MUDARRO_MATRIX_URL"}
			if tool == "django" || tool == "alembic" {
				s.Language = "python"
				s.Manager = "uv"
			}
			if err := saveConfig(filepath.Join(root, "mudarro.yaml"), c); err != nil {
				t.Fatal(err)
			}
			t.Setenv("MUDARRO_MATRIX_URL", "SECRET_NEVER_PRINT")
			for _, action := range Actions(*s) {
				if action.Group != "banco" || action.Blocked != "" {
					continue
				}
				var out bytes.Buffer
				err := Main([]string{"run", s.ID + ":" + action.Name, "--root", root, "--dry-run", "--name", "matrix_migration"}, "test", strings.NewReader(""), &out)
				if err != nil {
					t.Fatalf("%s: %v", action.Name, err)
				}
				if strings.Contains(out.String(), "SECRET_NEVER_PRINT") {
					t.Fatal("dry-run exposed secret")
				}
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 1 {
				t.Fatalf("dry-run wrote files: %v %v", entries, err)
			}
		})
	}
}

func TestSupportScanJSONDeterministicAndReadOnly(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"name":"json-matrix","scripts":{"z-last":"echo z","a-first":"echo a","test":"echo test"}}`, "Makefile": "check: deps\n\techo check\n"})
	var before, after bytes.Buffer
	args := []string{"scan", "--root", root, "--json"}
	if err := Main(args, "test", strings.NewReader(""), &before); err != nil {
		t.Fatal(err)
	}
	if err := Main(args, "test", strings.NewReader(""), &after); err != nil {
		t.Fatal(err)
	}
	if before.String() != after.String() {
		t.Fatal("JSON nondeterministic")
	}
	var report Report
	if err := json.Unmarshal(before.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Config.Name != "json-matrix" || len(report.Suggestions) != 4 {
		t.Fatalf("%+v", report)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 2 {
		t.Fatalf("scan wrote files: %v %v", entries, err)
	}
}

func TestSupportOversizedManifestWarning(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": strings.Repeat(" ", (4<<20)+1), "worker/go.mod": "module worker\ngo 1.23\n"})
	r, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Warnings) != 1 || !strings.Contains(r.Warnings[0], "grande demais") || len(r.Config.Services) != 1 || r.Config.Services[0].Language != "go" {
		t.Fatalf("oversized manifest should not hide valid sibling: %+v", r)
	}
}
