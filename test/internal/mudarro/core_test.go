package mudarro_test

import (
	"bytes"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for p, b := range files {
		full := filepath.Join(root, p)
		if e := os.MkdirAll(filepath.Dir(full), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(full, []byte(b), 0644); e != nil {
			t.Fatal(e)
		}
	}
	return root
}
func sample() Config {
	return Config{Version: 1, Name: "Minha app", UI: &UIConfig{Locale: "pt-BR"}, Services: []Service{{ID: "app", Dir: ".", Language: "go", Manager: "go", Infrastructure: Infrastructure{Kind: "local"}, Commands: map[string]Command{"start": cmd("aplicacao", "sleep", "30"), "test": cmd("qualidade", "go", "test", "./...")}}}}
}
func TestScanLanguagesAndSuggestions(t *testing.T) {
	root := fixture(t, map[string]string{"web/package.json": `{"name":"web","scripts":{"dev":"vite","test":"vitest","custom":"echo yes"},"devDependencies":{"typescript":"5"}}`, "web/pnpm-lock.yaml": "", "api/pyproject.toml": "[project]\nname='api'\n[project.scripts]\nhello='app:main'\n", "api/uv.lock": "", "worker/go.mod": "module worker\n\ngo 1.23\n", "worker/main.go": "package main\nfunc main(){}\n", "scripts/check.sh": "touch NEVER_EXECUTE", "Makefile": "check: deps\n\techo ok\n", "node_modules/bad/package.json": "bad", "dist/package.json": "bad"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 3 {
		t.Fatalf("%+v", r.Config.Services)
	}
	by := map[string]Service{}
	for _, s := range r.Config.Services {
		by[s.Language] = s
	}
	if by["typescript"].Manager != "pnpm" || by["python"].Manager != "uv" || len(by["go"].Commands["start"].Args) == 0 {
		t.Fatalf("%+v", by)
	}
	if len(r.Suggestions) < 5 {
		t.Fatalf("%+v", r.Suggestions)
	}
	if exists(filepath.Join(root, "NEVER_EXECUTE")) {
		t.Fatal("scan executed code")
	}
}
func TestAmbiguityExclusionsAndMalformed(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"scripts":{"start":"node app.js"}}`, "pnpm-lock.yaml": "", "package-lock.json": "{}", "compose.yaml": "services: {}", "k8s/a.yaml": "", "ignored/go.mod": "module nope", "bad/pyproject.toml": "[oops"})
	r, e := Scan(root, []string{"ignored"})
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 || r.Config.Services[0].Manager != "" || r.Config.Services[0].Infrastructure.Kind != "" || len(r.Warnings) != 1 {
		t.Fatalf("%+v", r)
	}
}
func TestScanDoesNotFollowSymlinks(t *testing.T) {
	outside := fixture(t, map[string]string{"package.json": `{"name":"outside"}`})
	root := t.TempDir()
	if e := os.Symlink(outside, filepath.Join(root, "linked")); e != nil {
		t.Fatal(e)
	}
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if r.Config.Services[0].Language != "custom" {
		t.Fatal(r)
	}
}
func TestConfigValidation(t *testing.T) {
	for _, tt := range []struct{ name, body string }{{"unknown", "version: 1\nname: x\nunknown: true\n"}, {"duplicate", "version: 1\nversion: 1\nname: x\n"}, {"multiple", "version: 1\nname: x\n---\nname: y\n"}, {"escape", "version: 1\nname: x\nservices:\n- id: app\n  dir: ../escape\n"}, {"commands", "version: 1\nname: x\nservices:\n- id: app\n  dir: .\n  commands:\n    test:\n      args: [echo, hi]\n      shell: echo hi\n"}} {
		t.Run(tt.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"mudarro.yaml": tt.body})
			if _, e := Load(root); e == nil {
				t.Fatal("accepted invalid config")
			}
		})
	}
}
func TestGenerateIdempotenceAndProtection(t *testing.T) {
	root := t.TempDir()
	c := sample()
	var out bytes.Buffer
	if e := Generate(root, c, true, &out); e != nil {
		t.Fatal(e)
	}
	if exists(filepath.Join(root, "menu.sh")) {
		t.Fatal("dry-run mutated files")
	}
	if e := Generate(root, c, false, &out); e != nil {
		t.Fatal(e)
	}
	before, _ := os.ReadFile(filepath.Join(root, ".mudarro/generated.json"))
	out.Reset()
	if e := Generate(root, c, false, &out); e != nil {
		t.Fatal(e)
	}
	after, _ := os.ReadFile(filepath.Join(root, ".mudarro/generated.json"))
	if !bytes.Equal(before, after) || out.Len() != 0 {
		t.Fatal("not idempotent")
	}
	if e := os.WriteFile(filepath.Join(root, "menu.sh"), []byte("# edited"), 0755); e != nil {
		t.Fatal(e)
	}
	if e := Generate(root, c, false, &out); e == nil {
		t.Fatal("overwrote user edit")
	}
}
func TestGeneratePreservesExistingAndRejectsSymlink(t *testing.T) {
	root := fixture(t, map[string]string{"menu.sh": "# user-owned"})
	var out bytes.Buffer
	if e := Generate(root, sample(), false, &out); e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(filepath.Join(root, "menu.sh"))
	if string(b) != "# user-owned" {
		t.Fatal("overwritten")
	}
	root = t.TempDir()
	if e := os.Symlink(t.TempDir(), filepath.Join(root, ".mudarro")); e != nil {
		t.Fatal(e)
	}
	if e := Generate(root, sample(), false, &out); e == nil {
		t.Fatal("followed symlink")
	}
}
func TestExplicitCommandsOverrideAdapters(t *testing.T) {
	s := sample().Services[0]
	s.Commands["up"] = cmd("custom", "echo", "explicit")
	a := Actions(s)
	for _, v := range a {
		if v.Name == "up" {
			if v.Special != "" || v.Command.Args[1] != "explicit" {
				t.Fatal(v)
			}
			return
		}
	}
	t.Fatal("missing up")
}
func TestStructuredArgumentsAndExitCode(t *testing.T) {
	skipOnWindows(t, "fixture runs POSIX printf")
	root := t.TempDir()
	s := sample().Services[0]
	var out bytes.Buffer
	a := action("test", "qualidade", "printf", "%s", "$(touch SHOULD_NOT_EXIST); spaces ' quotes")
	if e := (Runner{In: strings.NewReader(""), Out: &out}).Run(root, s, a); e != nil {
		t.Fatal(e)
	}
	if exists(filepath.Join(root, "SHOULD_NOT_EXIST")) {
		t.Fatal("shell injection")
	}
	if !strings.Contains(out.String(), "spaces ' quotes") {
		t.Fatal(out.String())
	}
	e := (Runner{Out: &out}).Run(root, s, action("fail", "scripts", "sh", "-c", "exit 7"))
	if v, ok := e.(interface{ ExitCode() int }); !ok || v.ExitCode() != 7 {
		t.Fatalf("%v", e)
	}
}
func TestDestructiveConfirmation(t *testing.T) {
	skipOnWindows(t, "fixture runs POSIX touch")
	root := t.TempDir()
	s := sample().Services[0]
	a := action("reset", "banco", "touch", "deleted")
	a.Command.Destructive = true
	var out bytes.Buffer
	if e := (Runner{In: strings.NewReader("yes\n"), Out: &out}).Run(root, s, a); e == nil {
		t.Fatal("accepted generic yes")
	}
	if exists(filepath.Join(root, "deleted")) {
		t.Fatal("executed")
	}
	if e := (Runner{In: strings.NewReader("APAGAR app\n"), Out: &out}).Run(root, s, a); e != nil {
		t.Fatal(e)
	}
}
func TestMenuAndCLI(t *testing.T) {
	root := t.TempDir()
	if e := saveConfig(filepath.Join(root, "mudarro.yaml"), sample()); e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e := Main([]string{"menu", "--root", root}, "test", strings.NewReader("abc\n1\n0\n0\n"), &out); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(out.String(), "Opção inválida") || !strings.Contains(out.String(), "infraestrutura") {
		t.Fatal(out.String())
	}
	out.Reset()
	if e := Main([]string{"scan", "--root", root, "--json"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	var r Report
	if e := json.Unmarshal(out.Bytes(), &r); e != nil {
		t.Fatal(e)
	}
}
func TestDatabaseAndInfrastructurePlans(t *testing.T) {
	for _, tool := range []string{"prisma", "django", "alembic", "goose"} {
		s := sample().Services[0]
		s.Database = Database{Kind: "sqlite", Tool: tool}
		a := Actions(s)
		found := map[string]bool{}
		for _, v := range a {
			found[v.Name] = true
			if v.Name == "db-reset" && !v.Command.Destructive {
				t.Fatal("unguarded reset")
			}
		}
		for _, n := range []string{"db-migrate", "db-migration-new", "db-seed"} {
			if !found[n] {
				t.Fatalf("%s missing %s", tool, n)
			}
		}
	}
	s := sample().Services[0]
	s.Infrastructure = Infrastructure{Kind: "docker", Mode: "compose", File: "compose.yaml"}
	for _, a := range Actions(s) {
		if a.Name == "down" && strings.Contains(strings.Join(a.Command.Args, " "), "--volumes") {
			t.Fatal("deletes volumes")
		}
	}
	s.Infrastructure = Infrastructure{Kind: "kubernetes", Mode: "manifests", File: "k8s"}
	if Actions(s)[0].Name == "up" && Actions(s)[0].Blocked == "" {
		t.Fatal("context missing")
	}
}
func TestScaffoldsParseAndDoNotInventModels(t *testing.T) {
	for _, tool := range []string{"prisma", "django", "alembic", "goose"} {
		s := sample().Services[0]
		s.Database = Database{Kind: "postgresql", Tool: tool, Generate: true}
		s.Infrastructure = Infrastructure{Kind: "docker", Mode: "compose", File: "compose.yaml", Generate: true}
		f, e := generatedScaffolds(t, s)
		if e != nil {
			t.Fatal(e)
		}
		if len(f) < 4 {
			t.Fatal(tool)
		}
		if bytes.Contains(f["prisma/schema.prisma"].Data, []byte("model User")) {
			t.Fatal("invented models")
		}
	}
	s := sample().Services[0]
	s.Infrastructure = Infrastructure{Kind: "kubernetes", Mode: "kustomize", File: "k8s", Generate: true, Image: "example:v1", Port: 8000, Context: "test", Namespace: "example"}
	s.Database = Database{Kind: "postgresql", Tool: "goose", Generate: true}
	f, e := generatedScaffolds(t, s)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(f["k8s/database.yaml"].Data), "NAMESPACE") || strings.Contains(string(f["k8s/database.yaml"].Data), "appSPACE") {
		t.Fatal("replacement collision")
	}
}
func TestInvalidMigrationNameAndEnvironment(t *testing.T) {
	root := t.TempDir()
	s := sample().Services[0]
	var out bytes.Buffer
	if err := (Runner{Out: &out, Name: "a; touch x"}).Run(root, s, action("migration", "banco", "echo", "{name}")); err == nil {
		t.Fatal("invalid name")
	}
	t.Setenv("MUDARRO_TEST_ABSENT", "")
	if err := (Runner{Out: &out}).Run(root, s, action("migration", "banco", "echo", "${MUDARRO_TEST_ABSENT}")); err == nil {
		t.Fatal("missing env")
	}
}
