package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const pipenvFixture = "[packages]\ndjango = '*'\n[scripts]\nhello = 'python hello.py'\n"
const pipenvLockFixture = `{"_meta":{"hash":{"sha256":"0000000000000000000000000000000000000000000000000000000000000000"},"pipfile-spec":6,"requires":{},"sources":[]},"default":{},"develop":{}}`

func TestPipenvSupportDetectionAndDjangoArgv(t *testing.T) {
	for _, locked := range []bool{false, true} {
		name := "Pipfile"
		if locked {
			name += "+lock"
		}
		t.Run(name, func(t *testing.T) {
			files := map[string]string{"Pipfile": pipenvFixture, "manage.py": "raise RuntimeError('NEVER_EXECUTE')"}
			install := []string{"pipenv", "install"}
			if locked {
				files["Pipfile.lock"] = pipenvLockFixture
				install = []string{"pipenv", "sync"}
			}
			root := fixture(t, files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			s := r.Config.Services[0]
			if s.Manager != "pipenv" || s.Language != "python" {
				t.Fatalf("Pipenv detection %+v", s)
			}
			for key, args := range map[string][]string{"install": install, "start": {"pipenv", "run", "python", "manage.py", "runserver"}, "test": {"pipenv", "run", "python", "manage.py", "test"}} {
				if !reflect.DeepEqual(s.Commands[key].Args, args) {
					t.Fatalf("%s argv %v want %v", key, s.Commands[key].Args, args)
				}
			}
			if _, ok := s.Commands["venv"]; ok {
				t.Fatal("Pipenv inferred manual .venv creation")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("scan mutated manifests")
			}
		})
	}
}
func TestPipenvSupportMalformedManifestsAreDiagnosed(t *testing.T) {
	for _, tc := range []struct{ name, pipfile, lock string }{{"toml", "[packages\nbad='*'", ""}, {"lock-json", pipenvFixture, "{invalid"}, {"script-value", "[scripts]\nhello=123\n", ""}} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"Pipfile": tc.pipfile}
			if tc.lock != "" {
				files["Pipfile.lock"] = tc.lock
			}
			root := fixture(t, files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			messages := strings.Join(r.Warnings, " ")
			for _, s := range r.Config.Services {
				messages += " " + strings.Join(s.Pending, " ")
				for _, su := range r.Suggestions {
					if su.Service == s.ID && su.Name == "hello" {
						t.Fatal("malformed script selected")
					}
				}
			}
			if !strings.Contains(messages, "Pipfile") && !(tc.name == "script-value" && strings.Contains(messages, "hello")) {
				t.Fatal("malformed Pipenv files silently accepted")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("diagnostic scan changed files")
			}
		})
	}
}
func TestPipenvSupportManagerConflictsRemainPending(t *testing.T) {
	for _, other := range []string{"uv.lock", "poetry.lock"} {
		t.Run(other, func(t *testing.T) {
			root := fixture(t, map[string]string{"Pipfile": pipenvFixture, "Pipfile.lock": pipenvLockFixture, other: "fixture", "manage.py": "raise RuntimeError('NEVER_EXECUTE')"})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			s := r.Config.Services[0]
			if s.Manager != "" || len(s.Pending) == 0 {
				t.Fatalf("ambiguous manager selected %+v", s)
			}
			for _, key := range []string{"install", "start", "test", "venv"} {
				if _, ok := s.Commands[key]; ok {
					t.Fatalf("manager conflict inferred %s: %+v", key, s.Commands[key])
				}
			}
			for _, su := range r.Suggestions {
				if len(su.Command.Args) > 0 && su.Command.Args[0] != "" {
					t.Fatalf("manager conflict offered executable suggestion %+v", su)
				}
			}
		})
	}
}
func TestPipenvSupportOrphanLockIsActionable(t *testing.T) {
	root := fixture(t, map[string]string{"Pipfile.lock": pipenvLockFixture})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	messages := strings.Join(r.Warnings, " ")
	for _, s := range r.Config.Services {
		messages += " " + strings.Join(s.Pending, " ")
		if _, ok := s.Commands["install"]; ok {
			t.Fatal("orphan lock inferred install")
		}
	}
	if !strings.Contains(messages, "Pipfile") {
		t.Fatalf("orphan lock not diagnosed: %s", messages)
	}
}
func TestPipenvSupportScriptSuggestionAndSelectedInit(t *testing.T) {
	root := fixture(t, map[string]string{"Pipfile": pipenvFixture, "hello.py": "raise RuntimeError('NEVER_EXECUTE')"})
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 {
		t.Fatal(r.Config.Services)
	}
	s := r.Config.Services[0]
	if _, ok := s.Commands["hello"]; ok {
		t.Fatal("custom Pipenv script selected automatically")
	}
	found := false
	for _, su := range r.Suggestions {
		if su.Name == "hello" {
			found = true
			if !reflect.DeepEqual(su.Command.Args, []string{"pipenv", "run", "hello"}) || su.Evidence != "Pipfile" {
				t.Fatalf("suggestion %+v", su)
			}
		}
	}
	if !found {
		t.Fatal("Pipfile script not suggested")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("scan executed or changed files")
	}
	var out bytes.Buffer
	if e = Main([]string{"init", "--root", root, "--select", "app-python:hello"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	c, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(c.Services[0].Commands["hello"].Args, []string{"pipenv", "run", "hello"}) {
		t.Fatal("selected script argv")
	}
	configBefore, e := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	if e = Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("init overwrote config")
	}
	configAfter, _ := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if !bytes.Equal(configBefore, configAfter) {
		t.Fatal("rejected init changed config")
	}
	if e = Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	hashes := supportFileHashes(t, root)
	if e = Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(hashes, supportFileHashes(t, root)) {
		t.Fatal("repeated generate changed files")
	}
}
func TestPipenvSupportScriptCollisionsStayExplicit(t *testing.T) {
	root := fixture(t, map[string]string{"Pipfile": "[packages]\n[scripts]\n'hello:world'='python one.py'\n'hello.world'='python two.py'\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, su := range r.Suggestions {
		if su.Name == "hello-world" {
			t.Fatal("script collision resolved arbitrarily")
		}
	}
	s := r.Config.Services[0]
	if !strings.Contains(strings.Join(s.Pending, " "), "hello-world") {
		t.Fatalf("collision not actionable: %v", s.Pending)
	}
}

func TestPipenvSupportDatabaseActionsUseOwnedEnvironment(t *testing.T) {
	for _, tool := range []string{"django", "alembic"} {
		t.Run(tool, func(t *testing.T) {
			files := map[string]string{"Pipfile": pipenvFixture}
			if tool == "django" {
				files["manage.py"] = "raise RuntimeError('NEVER_EXECUTE')"
			} else {
				files["alembic.ini"] = "[alembic]\n"
			}
			r, err := Scan(fixture(t, files), nil)
			if err != nil {
				t.Fatal(err)
			}
			s := r.Config.Services[0]
			s.Database.Kind = "sqlite"
			s.Database.Generate = true
			foundMigrate, foundInstall := false, false
			for _, a := range Actions(s) {
				if a.Name == "db-migrate" {
					foundMigrate = true
					prefix := []string{"pipenv", "run", "python"}
					if len(a.Command.Args) < len(prefix) || !reflect.DeepEqual(a.Command.Args[:len(prefix)], prefix) {
						t.Errorf("Pipenv db-migrate uses foreign environment: %v", a.Command.Args)
					}
				}
				if a.Name == "db-install" {
					foundInstall = true
					prefix := []string{"pipenv", "install"}
					if len(a.Command.Args) < len(prefix) || !reflect.DeepEqual(a.Command.Args[:len(prefix)], prefix) {
						t.Errorf("Pipenv db-install uses foreign environment: %v", a.Command.Args)
					}
				}
			}
			if !foundMigrate || !foundInstall {
				t.Fatalf("database actions missing migrate=%v install=%v", foundMigrate, foundInstall)
			}
		})
	}
}
