package mudarro_test

import (
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func frameworkFixture(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, b := range files {
		p = filepath.Join(root, p)
		if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(p, []byte(b), 0644); e != nil {
			t.Fatal(e)
		}
	}
	return root
}
func TestFrameworkDeclaredMetadataDoesNotInferEntrypoint(t *testing.T) {
	for _, tc := range []struct{ name, file, text, want string }{
		{"pep621", "pyproject.toml", "[project]\nname='sample'\ndependencies=['Flask>=3; python_version >= \"3\"']\n", "flask"},
		{"requirements", "requirements.txt", "# fastapi elsewhere\nfastapi[standard]>=0.100\n-r hidden.txt\n", "fastapi"},
		{"poetry", "pyproject.toml", "[tool.poetry.dependencies]\nFlask={version='*', optional=true}\n", "flask"},
		{"node", "package.json", `{"dependencies":{"express":"^5"},"scripts":{"check":"node owned.cjs"}}`, "express"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := frameworkFixture(t, map[string]string{tc.file: tc.text, "hidden.txt": "django\n", "app.py": "raise RuntimeError('NEVER IMPORT')"})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			s := r.Config.Services[0]
			if s.Framework != tc.want {
				t.Fatalf("framework=%q report=%+v", s.Framework, r)
			}
			if _, ok := s.Commands["start"]; ok {
				t.Fatal("invented start")
			}
			for _, su := range r.Suggestions {
				if su.Name == "start" {
					t.Fatal("invented suggestion")
				}
			}
			if len(r.Evidence) < 2 {
				t.Fatal("missing metadata evidence")
			}
		})
	}
}
func TestFrameworkAmbiguityAndUnknownRemainExplicit(t *testing.T) {
	for _, tc := range []struct{ name, text string }{{"multiple", `{"dependencies":{"express":"5","next":"15"}}`}, {"unknown", `{"dependencies":{"myframework":"1"}}`}, {"empty", `{"dependencies":{"express":""}}`}} {
		t.Run(tc.name, func(t *testing.T) {
			r, e := Scan(frameworkFixture(t, map[string]string{"package.json": tc.text}), nil)
			if e != nil {
				t.Fatal(e)
			}
			if r.Config.Services[0].Framework != "" {
				t.Fatalf("silently selected %+v", r.Config.Services[0])
			}
		})
	}
}
func TestFrameworkPytestIsToolOnlyAndOptionalMetadata(t *testing.T) {
	root := frameworkFixture(t, map[string]string{"pyproject.toml": "[project]\nname='sample'\n[project.optional-dependencies]\nquality=['pytest>=8']\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if r.Config.Services[0].Framework != "" {
		t.Fatal("pytest is not application framework")
	}
	found := false
	for _, ev := range r.Evidence {
		if ev.Kind == "declared-dependency:pytest" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing pytest evidence")
	}
	if _, ok := r.Config.Services[0].Commands["test"]; ok {
		t.Fatal("invented pytest command")
	}
}
func TestFrameworkMetadataWrongTypesExcludedSymlinkAndBounded(t *testing.T) {
	t.Run("wrong-types", func(t *testing.T) {
		r, e := Scan(frameworkFixture(t, map[string]string{"requirements.txt": "flask>=3\n", "pyproject.toml": "[project]\ndependencies=[42]\n"}), nil)
		if e != nil {
			t.Fatal(e)
		}
		if len(r.Warnings) == 0 {
			t.Fatal("wrong typed metadata accepted")
		}
	})
	t.Run("excluded", func(t *testing.T) {
		r, e := Scan(frameworkFixture(t, map[string]string{"requirements.txt": "flask\n", "pyproject.toml": "[project]\ndependencies=['fastapi']\n"}), []string{"requirements.txt"})
		if e != nil {
			t.Fatal(e)
		}
		if r.Config.Services[0].Framework != "fastapi" {
			t.Fatalf("excluded dependency read %+v", r)
		}
	})
	t.Run("symlink", func(t *testing.T) {
		root := frameworkFixture(t, map[string]string{"pyproject.toml": "[project]\nname='sample'\n"})
		outside := frameworkFixture(t, map[string]string{"outside": "flask\n"})
		if e := os.Symlink(filepath.Join(outside, "outside"), filepath.Join(root, "requirements.txt")); e != nil {
			t.Fatal(e)
		}
		r, e := Scan(root, nil)
		if e != nil {
			t.Fatal(e)
		}
		if r.Config.Services[0].Framework != "" {
			t.Fatal("symlink metadata read")
		}
		if len(r.Warnings) == 0 {
			t.Fatal("missing symlink diagnostic")
		}
	})
	t.Run("oversized", func(t *testing.T) {
		r, e := Scan(frameworkFixture(t, map[string]string{"pyproject.toml": "[project]\nname='sample'\n", "requirements.txt": strings.Repeat("#", (4<<20)+1)}), nil)
		if e != nil {
			t.Fatal(e)
		}
		if len(r.Warnings) == 0 {
			t.Fatal("oversized metadata read")
		}
	})
}
func TestFrameworkScanNoExecutablesOrImports(t *testing.T) {
	root := frameworkFixture(t, map[string]string{"requirements.txt": "Flask>=3\n", "app.py": "open('marker','w').write('imported')"})
	tools := frameworkFixture(t, map[string]string{"python3": "#!/bin/sh\ntouch '" + filepath.Join(root, "marker") + "'\n", "flask": "#!/bin/sh\nexit 89\n"})
	if e := os.Chmod(filepath.Join(tools, "python3"), 0755); e != nil {
		t.Fatal(e)
	}
	t.Setenv("PATH", tools)
	if _, e := Scan(root, nil); e != nil {
		t.Fatal(e)
	}
	if _, e := os.Stat(filepath.Join(root, "marker")); !os.IsNotExist(e) {
		t.Fatal("executed framework discovery")
	}
}

func TestFrameworkDeclaredEntrypointSuggestionPreservedWithoutImport(t *testing.T) {
	root := frameworkFixture(t, map[string]string{"pyproject.toml": "[project]\nname='sample'\ndependencies=['FastAPI>=0.100']\n[project.scripts]\nowned='owned.module:main'\n", "owned/module.py": "raise RuntimeError('NEVER IMPORT')"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, su := range r.Suggestions {
		if su.Name == "owned" {
			found = true
			if len(su.Command.Args) != 1 || su.Command.Args[0] != ".venv/bin/owned" {
				t.Fatalf("entrypoint changed: %+v", su)
			}
		}
	}
	if !found {
		t.Fatal("explicit entrypoint removed")
	}
	if _, ok := r.Config.Services[0].Commands["start"]; ok {
		t.Fatal("explicit script treated as implicit start")
	}
}
