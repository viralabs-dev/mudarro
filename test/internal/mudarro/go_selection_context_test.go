package mudarro_test

import (
	"bytes"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func goSelectionFixture(t *testing.T) string {
	t.Helper()
	return fixture(t, map[string]string{"go.mod": "module example.test/selection\ngo 1.23\n", "cmd/z/main.go": "package main\nfunc main(){}\n", "cmd/a/main.go": "package main\nfunc main(){}\n"})
}
func TestGoSelectionSortedOfflineSuggestions(t *testing.T) {
	root := goSelectionFixture(t)
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 {
		t.Fatalf("services %+v", r.Config.Services)
	}
	if _, ok := r.Config.Services[0].Commands["start"]; ok {
		t.Fatal("ambiguous auto start")
	}
	if len(r.Suggestions) != 2 {
		t.Fatalf("suggestions %+v", r.Suggestions)
	}
	for i, target := range []string{"./cmd/a", "./cmd/z"} {
		su := r.Suggestions[i]
		if su.Name != []string{"go-entry-1", "go-entry-2"}[i] || !reflect.DeepEqual(su.Command.Args, []string{"go", "run", target}) || su.Command.Group != "aplicacao" || !strings.Contains(su.Evidence, strings.TrimPrefix(target, "./")) {
			t.Errorf("suggestion %+v", su)
		}
		b, _ := json.Marshal(su)
		var m map[string]any
		json.Unmarshal(b, &m)
		if m["purpose"] != "go-start" {
			t.Errorf("purpose %s", b)
		}
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("scan mutated fixture")
	}
}
func TestGoSelectionInitPromotesOnlyStart(t *testing.T) {
	root := goSelectionFixture(t)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	id := r.Config.Services[0].ID
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--select", id + ":go-entry-2"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	c, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	s := c.Services[0]
	if !reflect.DeepEqual(s.Commands["start"].Args, []string{"go", "run", "./cmd/z"}) {
		t.Fatalf("start %+v", s)
	}
	for k := range s.Commands {
		if strings.HasPrefix(k, "go-entry-") {
			t.Fatalf("synthetic action %s", k)
		}
	}
	for _, p := range s.Pending {
		if strings.Contains(p, "commands.start") {
			t.Fatalf("resolved start still pending %q", p)
		}
	}
	before := supportFileHashes(t, root)
	if e := Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("existing config accepted")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("existing config changed")
	}
}
func TestGoSelectionMultipleRejectedBeforeWrite(t *testing.T) {
	for _, mode := range []string{"select", "all", "interactive"} {
		t.Run(mode, func(t *testing.T) {
			root := goSelectionFixture(t)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			id := r.Config.Services[0].ID
			args := []string{"init", "--root", root}
			input := ""
			switch mode {
			case "select":
				args = append(args, "--select", id+":go-entry-1,"+id+":go-entry-2")
			case "all":
				args = append(args, "--all-scripts")
			case "interactive":
				args = append(args, "--interactive")
				input = "y\ny\n"
			}
			before := supportFileHashes(t, root)
			if e := Main(args, "test", strings.NewReader(input), &bytes.Buffer{}); e == nil {
				t.Fatal("ambiguous selection accepted")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("rejected init wrote files")
			}
		})
	}
}
func TestGoSelectionWorkspaceConfigRoundtrip(t *testing.T) {
	for _, format := range []string{"json", "yaml"} {
		for _, mode := range []string{"", "inherit", "off", "invalid"} {
			t.Run(format+"/"+mode, func(t *testing.T) {
				root := fixture(t, map[string]string{"go.mod": "module example.test/context\ngo 1.23\n"})
				m := map[string]any{"version": 1, "name": "context", "services": []any{map[string]any{"id": "go-app", "dir": ".", "language": "go", "go_workspace": mode, "commands": map[string]any{"test": map[string]any{"args": []string{"go", "test", "./..."}}}}}}
				var b []byte
				if format == "json" {
					b, _ = json.Marshal(m)
				} else {
					b, _ = yaml.Marshal(m)
				}
				p := filepath.Join(root, "mudarro."+format)
				if e := os.WriteFile(p, b, 0600); e != nil {
					t.Fatal(e)
				}
				c, e := Load(root)
				if mode == "invalid" {
					if e == nil {
						t.Fatal("invalid go_workspace accepted")
					}
					return
				}
				if e != nil {
					t.Fatal(e)
				}
				raw, _ := json.Marshal(c)
				var got map[string]any
				json.Unmarshal(raw, &got)
				s := got["services"].([]any)[0].(map[string]any)
				if mode == "" {
					if v, ok := s["go_workspace"]; ok && v != "" {
						t.Fatalf("default changed %v", v)
					}
				} else if s["go_workspace"] != mode {
					t.Fatalf("workspace lost %s", raw)
				}
				after, _ := os.ReadFile(p)
				if !bytes.Equal(b, after) {
					t.Fatal("load modified config")
				}
			})
		}
	}
}

func TestGoSelectionWorkspaceNonGoRejected(t *testing.T) {
	root := fixture(t, map[string]string{"mudarro.json": `{"version":1,"name":"context","services":[{"id":"node","dir":".","language":"javascript","go_workspace":"off"}]}`})
	if _, e := Load(root); e == nil {
		t.Fatal("non-Go workspace option accepted")
	}
}
func TestGoSelectionSingleKeepsAutomaticStartAndDefaultOmitted(t *testing.T) {
	root := fixture(t, map[string]string{"go.mod": "module example.test/unique\ngo 1.23\n", "cmd/app/main.go": "package main\nfunc main(){}\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 || !reflect.DeepEqual(r.Config.Services[0].Commands["start"].Args, []string{"go", "run", "./cmd/app"}) {
		t.Fatalf("single %+v", r.Config.Services)
	}
	b, _ := json.Marshal(r.Config)
	if strings.Contains(string(b), `"go_workspace"`) {
		t.Fatalf("default serialized %s", b)
	}
	id := r.Config.Services[0].ID
	if e := Main([]string{"init", "--root", root, "--select", id + ":go-entry-1"}, "test", strings.NewReader(""), &bytes.Buffer{}); e != nil {
		t.Fatal(e)
	}
	c, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(c.Services[0].Commands["start"].Args, []string{"go", "run", "./cmd/app"}) {
		t.Fatal("explicit unique changed inferred start")
	}
}

func TestGoSelectionWorkspaceRunnerChildEnvironmentScoped(t *testing.T) {
	t.Setenv("GOWORK", filepath.Join(t.TempDir(), "parent-workspace"))
	parent := os.Getenv("GOWORK")
	for _, mode := range []string{"", "inherit", "off"} {
		for _, shell := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "/argv", true: "/shell"}[shell], func(t *testing.T) {
				root := t.TempDir()
				data := map[string]any{"id": "go-app", "dir": ".", "language": "go", "go_workspace": mode}
				body, _ := json.Marshal(data)
				var s Service
				if e := json.Unmarshal(body, &s); e != nil {
					t.Fatal(e)
				}
				command := Command{Args: []string{"sh", "-c", "printenv GOWORK | tr -d '\\n'"}}
				if shell {
					command = Command{Shell: `printf '%s' "$GOWORK"`}
				}
				var out bytes.Buffer
				if e := (Runner{Out: &out}).Run(root, s, Action{Name: "probe", Command: command}); e != nil {
					t.Fatal(e)
				}
				want := parent
				if mode == "off" {
					want = "off"
				}
				if out.String() != want {
					t.Fatalf("child GOWORK=%q want=%q", out.String(), want)
				}
				if os.Getenv("GOWORK") != parent {
					t.Fatal("parent environment changed")
				}
			})
		}
	}
}
