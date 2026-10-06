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

func TestBunSupportLockDetection(t *testing.T) {
	for _, locks := range [][]string{{"bun.lock"}, {"bun.lockb"}, {"bun.lock", "bun.lockb"}} {
		t.Run(strings.Join(locks, "+"), func(t *testing.T) {
			files := map[string]string{"package.json": `{"name":"bun-probe","scripts":{"start":"bun app.js","test":"bun test","hello":"bun hello.js"}}`}
			for _, lock := range locks {
				files[lock] = "fixture lock detection only"
			}
			root := fixture(t, files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			if len(r.Config.Services) != 1 {
				t.Fatalf("services %+v", r.Config.Services)
			}
			s := r.Config.Services[0]
			if s.Manager != "bun" {
				t.Fatalf("Bun lock manager=%q", s.Manager)
			}
			for name, args := range map[string][]string{"start": {"bun", "run", "start"}, "test": {"bun", "run", "test"}, "install": {"bun", "install"}} {
				if !reflect.DeepEqual(s.Commands[name].Args, args) {
					t.Fatalf("%s argv %+v", name, s.Commands[name])
				}
			}
			if len(s.Pending) != 0 {
				t.Fatalf("same-manager locks became ambiguous: %v", s.Pending)
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("scan changed project")
			}
		})
	}
}
func TestBunSupportConflictingManagersRemainExplicit(t *testing.T) {
	for _, other := range []string{"package-lock.json", "pnpm-lock.yaml", "yarn.lock"} {
		t.Run(other, func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": `{"scripts":{"start":"bun app.js","test":"bun test"}}`, "bun.lock": "fixture", other: "fixture"})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			s := r.Config.Services[0]
			if s.Manager != "" || len(s.Pending) == 0 {
				t.Fatalf("conflict chose manager %+v", s)
			}
			for _, name := range []string{"start", "test", "install"} {
				if _, ok := s.Commands[name]; ok {
					t.Fatalf("conflict selected %s automatically", name)
				}
			}
		})
	}
}
func TestBunSupportExplicitPackageManagerOverridesLockConflict(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"packageManager":"bun@1.4.0","scripts":{"start":"bun app.js","test":"bun test"}}`, "bun.lock": "fixture", "package-lock.json": "fixture"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	s := r.Config.Services[0]
	if s.Manager != "bun" || !reflect.DeepEqual(s.Commands["test"].Args, []string{"bun", "run", "test"}) {
		t.Fatalf("explicit Bun override %+v", s)
	}
}
func TestBunSupportInitSelectedScriptAndGenerateIdempotent(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"name":"bun-probe","scripts":{"start":"bun app.js","test":"bun test","hello":"bun hello.js"}}`, "bun.lock": "fixture"})
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--select", "app-javascript:hello"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	c, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	s := c.Services[0]
	if s.Manager != "bun" || !reflect.DeepEqual(s.Commands["hello"].Args, []string{"bun", "run", "hello"}) {
		t.Fatalf("selected Bun script %+v", s)
	}
	body, e := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	if e := Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("repeated init accepted")
	}
	after, _ := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if !bytes.Equal(body, after) {
		t.Fatal("repeated init changed config")
	}
	if e := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, root)
	if e := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("repeated generation changed files")
	}
	if !exists(filepath.Join(root, ".mudarro/scripts/app-javascript/hello.sh")) {
		t.Fatal("selected Bun wrapper missing")
	}
}
