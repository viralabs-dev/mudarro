//go:build linux || darwin

package mudarro_test

import (
	"context"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Read descriptors through public JSON so the before-fix test compiles without
// adding a production API solely to accommodate the regression.
func parsedWorkspaceMembers(t *testing.T, r Report, path string) ([]string, bool) {
	t.Helper()
	b, e := json.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	var envelope struct {
		Workspaces []struct {
			Path    string   `json:"path"`
			Members []string `json:"members"`
		} `json:"workspaces"`
	}
	if e = json.Unmarshal(b, &envelope); e != nil {
		t.Fatal(e)
	}
	for _, w := range envelope.Workspaces {
		if filepath.ToSlash(w.Path) == path {
			return w.Members, true
		}
	}
	return nil, false
}

func TestGoWorkParserValidQuotedBlockAndNormalizedMembership(t *testing.T) {
	root := fixture(t, map[string]string{
		"go.work":          "// actual work syntax, no subprocess\ngo 1.23\ntoolchain go1.27.1\nuse (\n \"./space app\" // quoted member\n ./z/../z\n)\n",
		"space app/go.mod": "module example.test/space\ngo 1.23\n", "space app/main.go": "package main\nfunc main(){}\n",
		"z/go.mod": "module example.test/z\ngo 1.23\n", "z/main.go": "package main\nfunc main(){}\n",
		"nonmember/go.mod": "module example.test/nonmember\ngo 1.23\n", "nonmember/main.go": "package main\nfunc main(){}\n",
	})
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, ok := parsedWorkspaceMembers(t, r, "go.work")
	if !ok || !reflect.DeepEqual(members, []string{"space app", "z"}) {
		t.Errorf("members=%v present=%v", members, ok)
	}
	if len(r.Warnings) > 0 {
		t.Errorf("valid membership received blanket/error warning: %v", r.Warnings)
	}
	if len(r.Config.Services) != 3 {
		t.Errorf("workspace hid independent module: %+v", r.Config.Services)
	}
	evidence := false
	for _, v := range r.Evidence {
		if v.Kind == "go-workspace" && v.Path == "go.work" {
			evidence = true
		}
	}
	if !evidence {
		t.Error("workspace evidence lost")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("parser changed files")
	}
}

func TestGoWorkParserOfficialReplacementSyntaxKeepsMembershipScope(t *testing.T) {
	root := fixture(t, map[string]string{"go.work": "go 1.23\nuse ./app\nreplace example.test/remote v1.0.0 => ./replacement-not-member\n", "app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, ok := parsedWorkspaceMembers(t, r, "go.work")
	if !ok || !reflect.DeepEqual(members, []string{"app"}) {
		t.Errorf("replace incorrectly treated as membership: %v %v", members, ok)
	}
}

func TestGoWorkParserInvalidMembershipWarnsAndPreservesModules(t *testing.T) {
	for _, tc := range []struct {
		name, work  string
		extra       map[string]string
		wantMembers []string
	}{
		{"syntax", "go 1.23\nuse (\n", nil, nil},
		{"missing", "go 1.23\nuse (\n ./app\n ./missing\n)\n", nil, []string{"app"}},
		{"duplicate-canonical-member", "go 1.23\nuse (\n ./app\n ./app/../app\n)\n", nil, []string{"app"}},
		{"not-module", "go 1.23\nuse (\n ./app\n ./plain\n)\n", map[string]string{"plain/readme.txt": "not module"}, []string{"app"}},
		{"invalid-member-module", "go 1.23\nuse (\n ./app\n ./invalid\n)\n", map[string]string{"invalid/go.mod": "not a module directive\n"}, []string{"app"}},
		{"missing-module-declaration", "go 1.23\nuse (\n ./app\n ./invalid\n)\n", map[string]string{"invalid/go.mod": "go 1.23\n"}, []string{"app"}},
		{"outside-relative", "go 1.23\nuse (\n ./app\n ../outside\n)\n", nil, []string{"app"}},
		{"outside-absolute", "go 1.23\nuse (\n ./app\n /does-not-belong-to-project\n)\n", nil, []string{"app"}},
		{"work-over-limit", "//" + strings.Repeat("x", 4<<20) + "\ngo 1.23\nuse ./app\n", nil, nil},
		{"module-over-limit", "go 1.23\nuse (\n ./app\n ./invalid\n)\n", map[string]string{"invalid/go.mod": "module example.test/invalid\n//" + strings.Repeat("x", 4<<20)}, []string{"app"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"go.work": tc.work, "app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n"}
			for p, b := range tc.extra {
				files[p] = b
			}
			root := fixture(t, files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			members, ok := parsedWorkspaceMembers(t, r, "go.work")
			if !ok {
				t.Error("diagnosed workspace descriptor absent")
			}
			if len(members) != len(tc.wantMembers) || !reflect.DeepEqual(append([]string{}, members...), append([]string{}, tc.wantMembers...)) {
				t.Errorf("members=%v want=%v", members, tc.wantMembers)
			}
			diagnostic := strings.Join(r.Warnings, " ")
			if !strings.Contains(diagnostic, "go.work") {
				t.Errorf("missing path-specific warning: %q", diagnostic)
			}
			app := false
			for _, s := range r.Config.Services {
				if s.Dir == "app" {
					app = true
				}
			}
			if !app {
				t.Error("invalid workspace hid valid independent app")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("invalid scan changed fixture")
			}
		})
	}
}

func TestGoWorkParserSymlinkMemberNeverTraversed(t *testing.T) {
	root := fixture(t, map[string]string{"go.work": "go 1.23\nuse ./linked\n", "app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n"})
	outside := fixture(t, map[string]string{"go.mod": "module example.test/outside\ngo 1.23\n", "main.go": "package main\nfunc main(){}\n"})
	before := supportFileHashes(t, outside)
	if e := os.Symlink(outside, filepath.Join(root, "linked")); e != nil {
		t.Fatal(e)
	}
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, ok := parsedWorkspaceMembers(t, r, "go.work")
	if !ok || len(members) != 0 {
		t.Errorf("symlink admitted: %v %v", members, ok)
	}
	if len(r.Warnings) == 0 {
		t.Error("symlink member silent")
	}
	if len(r.Config.Services) != 1 || r.Config.Services[0].Dir != "app" {
		t.Error("scanner traversed symlink")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, outside)) {
		t.Fatal("outside fixture modified")
	}
}

func TestGoWorkParserOfflineNeverRunsGo(t *testing.T) {
	root := fixture(t, map[string]string{"go.work": "go 1.23\nuse ./app\n", "app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n"})
	tools := t.TempDir()
	if e := os.WriteFile(filepath.Join(tools, "go"), []byte("#!/bin/sh\nprintf executed > '"+filepath.Join(root, "SCAN_MUST_NOT_EXECUTE")+"'\nexit 92\n"), 0755); e != nil {
		t.Fatal(e)
	}
	t.Setenv("PATH", tools)
	t.Setenv("GOPROXY", "off")
	t.Setenv("GOTOOLCHAIN", "local")
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, ok := parsedWorkspaceMembers(t, r, "go.work")
	if !ok || !reflect.DeepEqual(members, []string{"app"}) {
		t.Errorf("offline parser descriptors absent: %v %v", members, ok)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("scan executed helper or changed source")
	}
}

// Workspace/source FIFOs are never opened. The child scan is exclusive and
// killed/joined by its watchdog if a future IO regression blocks opening.
func TestGoWorkParserFIFOIsBounded(t *testing.T) {
	for _, kind := range []string{"workspace", "member-module"} {
		t.Run(kind, func(t *testing.T) {
			files := map[string]string{"app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n"}
			if kind == "member-module" {
				files["go.work"] = "go 1.23\nuse (\n ./app\n ./invalid\n)\n"
				files["invalid/keep.txt"] = "fixture"
			}
			root := fixture(t, files)
			fifo := "go.work"
			if kind == "member-module" {
				fifo = "invalid/go.mod"
			}
			if e := syscall.Mkfifo(filepath.Join(root, fifo), 0600); e != nil {
				t.Fatal(e)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGoWorkParserIsolatedProbe$")
			child.Env = append(os.Environ(), "MUDARRO_WORK_PARSER_PROBE="+root)
			output, e := child.CombinedOutput()
			if ctx.Err() != nil {
				t.Fatal("scan opened FIFO; owned helper killed/joined", ctx.Err())
			}
			if e != nil {
				t.Fatalf("helper: %v %s", e, output)
			}
			var r Report
			if e := json.Unmarshal(output, &r); e != nil {
				t.Fatalf("report: %v %s", e, output)
			}
			members, ok := parsedWorkspaceMembers(t, r, "go.work")
			want := []string{}
			if kind == "member-module" {
				want = []string{"app"}
			}
			if !ok || !reflect.DeepEqual(append([]string{}, members...), want) {
				t.Errorf("FIFO descriptor members=%v present=%v", members, ok)
			}
			if !strings.Contains(strings.Join(r.Warnings, " "), "go.work") {
				t.Error("FIFO workspace issue not diagnosed")
			}
			st, e := os.Lstat(filepath.Join(root, fifo))
			if e != nil || st.Mode()&os.ModeNamedPipe == 0 {
				t.Fatal("FIFO changed", e)
			}
		})
	}
}

func TestGoWorkParserIsolatedProbe(t *testing.T) {
	root := os.Getenv("MUDARRO_WORK_PARSER_PROBE")
	if root == "" {
		return
	}
	var excludes []string
	if raw := os.Getenv("MUDARRO_WORK_PARSER_EXCLUDES"); raw != "" {
		if e := json.Unmarshal([]byte(raw), &excludes); e != nil {
			os.Exit(4)
		}
	}
	r, e := Scan(root, excludes)
	if e != nil {
		os.Exit(2)
	}
	if e := json.NewEncoder(os.Stdout).Encode(r); e != nil {
		os.Exit(3)
	}
	os.Exit(0)
}

func TestGoWorkParserDuplicateModulePathIsDiagnosedWithoutHidingMembers(t *testing.T) {
	root := fixture(t, map[string]string{"go.work": "go 1.23\nuse (\n ./one\n ./two\n)\n", "one/go.mod": "module example.test/identical\ngo 1.23\n", "two/go.mod": "module example.test/identical\ngo 1.23\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, ok := parsedWorkspaceMembers(t, r, "go.work")
	if !ok || !reflect.DeepEqual(members, []string{"one", "two"}) {
		t.Errorf("conflict hid actual members: %v %v", members, ok)
	}
	warnings := strings.Join(r.Warnings, " ")
	if !strings.Contains(warnings, "go.work") || !strings.Contains(warnings, "example.test/identical") {
		t.Errorf("duplicate module path silent: %q", warnings)
	}
}

func TestGoWorkParserExcludedMembershipIsNotRead(t *testing.T) {
	for _, exclusion := range []string{"private", "private/go.mod"} {
		t.Run(exclusion, func(t *testing.T) {
			root := fixture(t, map[string]string{"go.work": "go 1.23\nuse (\n ./app\n ./private\n ./.mudarro\n)\n", "app/go.mod": "module example.test/app\ngo 1.23\n", "app/main.go": "package main\nfunc main(){}\n", "private/keep.txt": "fixture", ".mudarro/keep.txt": "fixture"})
			for _, path := range []string{"private/go.mod", ".mudarro/go.mod"} {
				if e := syscall.Mkfifo(filepath.Join(root, path), 0600); e != nil {
					t.Fatal(e)
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGoWorkParserIsolatedProbe$")
			excludedJSON, _ := json.Marshal([]string{exclusion})
			child.Env = append(os.Environ(), "MUDARRO_WORK_PARSER_PROBE="+root, "MUDARRO_WORK_PARSER_EXCLUDES="+string(excludedJSON))
			output, e := child.CombinedOutput()
			if ctx.Err() != nil {
				t.Fatal("parser read excluded FIFO; exclusive helper killed/joined", ctx.Err())
			}
			if e != nil {
				t.Fatalf("helper: %v %s", e, output)
			}
			var r Report
			if e := json.Unmarshal(output, &r); e != nil {
				t.Fatalf("report: %v %s", e, output)
			}
			members, ok := parsedWorkspaceMembers(t, r, "go.work")
			if !ok || !reflect.DeepEqual(members, []string{"app"}) {
				t.Errorf("excluded member admitted: %v %v", members, ok)
			}
			var envelope struct {
				Workspaces []struct {
					Path     string            `json:"path"`
					Excluded []string          `json:"excluded"`
					Modules  map[string]string `json:"modules"`
				} `json:"workspaces"`
			}
			if e := json.Unmarshal(output, &envelope); e != nil {
				t.Fatal(e)
			}
			if len(envelope.Workspaces) != 1 || !reflect.DeepEqual(envelope.Workspaces[0].Excluded, []string{".mudarro", "private"}) {
				t.Errorf("excluded declared metadata missing: %+v", envelope)
			}
			if len(envelope.Workspaces) > 0 {
				for _, name := range []string{"private", ".mudarro"} {
					if _, present := envelope.Workspaces[0].Modules[name]; present {
						t.Errorf("excluded module read: %s", name)
					}
				}
			}
			diagnostic := strings.ToLower(strings.Join(r.Warnings, " "))
			if !strings.Contains(diagnostic, "excluded") || !strings.Contains(diagnostic, "not validated") {
				t.Errorf("exclusion not disclosed: %q", diagnostic)
			}
			if strings.Contains(diagnostic, "regular") {
				t.Errorf("excluded FIFO content/type was validated: %q", diagnostic)
			}
			if len(r.Config.Services) != 1 || r.Config.Services[0].Dir != "app" {
				t.Errorf("excluded services returned: %+v", r.Config.Services)
			}
			for _, path := range []string{"private/go.mod", ".mudarro/go.mod"} {
				st, e := os.Lstat(filepath.Join(root, path))
				if e != nil || st.Mode()&os.ModeNamedPipe == 0 {
					t.Fatal("excluded FIFO changed", e)
				}
			}
		})
	}
}
