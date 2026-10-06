package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestGoWorkspaceEntrypointConstraints(t *testing.T) {
	opposite := "windows"
	if runtime.GOOS == "windows" {
		opposite = "linux"
	}
	for _, tc := range []struct {
		name  string
		files map[string]string
		want  []string
	}{
		{"library", map[string]string{"lib.go": "package example\nfunc Work(){}\n"}, nil},
		{"tests-only", map[string]string{"main_test.go": "package main\nfunc main(){}\n"}, nil},
		{"main-package-without-entry", map[string]string{"main.go": "package main\nfunc Helper(){}\n"}, nil},
		{"ignored-main", map[string]string{"main.go": "//go:build ignore\n\npackage main\nfunc main(){}\n"}, nil},
		{"inactive-host-tag", map[string]string{"main.go": "//go:build " + opposite + "\n\npackage main\nfunc main(){}\n"}, nil},
		{"inactive-host-suffix", map[string]string{"main_" + opposite + ".go": "package main\nfunc main(){}\n"}, nil},
		{"active-plus-ignored-other-package", map[string]string{"cmd/live/main.go": "package main\nfunc main(){}\n", "cmd/tool/main.go": "//go:build ignore\n\npackage main\nfunc main(){}\n"}, []string{"go", "run", "./cmd/live"}},
		{"same-package-main-files", map[string]string{"cmd/live/main.go": "package main\nfunc main(){}\n", "cmd/live/helper.go": "package main\nfunc Helper(){}\n"}, []string{"go", "run", "./cmd/live"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.files["go.mod"] = "module example.test/local\ngo 1.23\n"
			root := fixture(t, tc.files)
			r, err := Scan(root, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Config.Services) != 1 {
				t.Fatalf("services: %+v", r.Config.Services)
			}
			s := r.Config.Services[0]
			if !reflect.DeepEqual(s.Commands["start"].Args, tc.want) {
				t.Errorf("start=%v want=%v; pending=%v", s.Commands["start"].Args, tc.want, s.Pending)
			}
			if tc.want == nil && len(s.Pending) == 0 {
				t.Error("unresolved start must stay pending")
			}
			for _, name := range []string{"build", "test"} {
				if !reflect.DeepEqual(s.Commands[name].Args, []string{"go", name, "./..."}) {
					t.Errorf("%s crossed module scope: %v", name, s.Commands[name].Args)
				}
			}
		})
	}
}

func TestGoWorkspaceInvalidMembershipIsVisible(t *testing.T) {
	for _, tc := range []struct {
		name, work string
		extras     map[string]string
	}{
		{"syntax", "go 1.23\nuse (\n", nil},
		{"missing", "go 1.23\nuse ./missing\n", nil},
		{"member-without-module", "go 1.23\nuse ./plain\n", map[string]string{"plain/lib.go": "package plain\n"}},
		{"outside-root", "go 1.23\nuse ../outside\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"go.work": tc.work, "app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n"}
			for k, v := range tc.extras {
				files[k] = v
			}
			root := fixture(t, files)
			r, err := Scan(root, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Warnings) == 0 {
				t.Errorf("invalid go.work membership silently accepted: %q", tc.work)
			}
			if len(r.Config.Services) != 1 || r.Config.Services[0].Dir != "app" {
				t.Errorf("workspace issue hid independent module: %+v", r.Config.Services)
			}
		})
	}
}

func TestGoWorkspaceNonmembersRemainIndependentAndOffline(t *testing.T) {
	root := fixture(t, map[string]string{"go.work": "go 1.23\nuse ./member\n", "member/go.mod": "module example.test/member\ngo 1.23\n", "member/main.go": "package main\nfunc main(){}\n", "other/go.mod": "module example.test/other\ngo 1.23\n", "other/cmd/cli/main.go": "package main\nfunc main(){}\n"})
	trap := filepath.Join(t.TempDir(), "go")
	marker := filepath.Join(root, "SCAN_MUST_NOT_EXECUTE")
	if err := os.WriteFile(trap, []byte("#!/bin/sh\nprintf executed > '"+marker+"'\nexit 99\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", filepath.Dir(trap))
	t.Setenv("GOPROXY", "off")
	t.Setenv("GOTOOLCHAIN", "local")
	before := map[string][]byte{}
	for _, p := range []string{"go.work", "member/go.mod", "other/go.mod"} {
		b, e := os.ReadFile(filepath.Join(root, p))
		if e != nil {
			t.Fatal(e)
		}
		before[p] = b
	}
	for i := 0; i < 2; i++ {
		r, e := Scan(root, nil)
		if e != nil {
			t.Fatal(e)
		}
		if len(r.Warnings) != 0 {
			t.Errorf("valid parsed membership warned: %v", r.Warnings)
		}
		if len(r.Workspaces) != 1 || !reflect.DeepEqual(r.Workspaces[0].Members, []string{"member"}) {
			t.Fatalf("workspace membership not parsed: %+v", r.Workspaces)
		}
		if len(r.Config.Services) != 2 {
			t.Fatalf("nonmember hidden: %+v", r.Config.Services)
		}
		for _, s := range r.Config.Services {
			if s.Language != "go" || !strings.HasPrefix(strings.Join(s.Commands["start"].Args, " "), "go run ./") {
				t.Errorf("unexpected module: %+v", s)
			}
		}
	}
	if _, e := os.Stat(marker); !os.IsNotExist(e) {
		t.Fatalf("scan invoked executable: %v", e)
	}
	for p, b := range before {
		after, e := os.ReadFile(filepath.Join(root, p))
		if e != nil || !bytes.Equal(b, after) {
			t.Errorf("scan mutated %s", p)
		}
	}
}
