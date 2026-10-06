package mudarro_test

import (
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func nativeFixture(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for path, text := range files {
		p := filepath.Join(root, path)
		if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(p, []byte(text), 0644); e != nil {
			t.Fatal(e)
		}
	}
	return root
}
func nativeService(t *testing.T, r Report, dir string) Service {
	t.Helper()
	for _, s := range r.Config.Services {
		if s.Dir == dir {
			return s
		}
	}
	t.Fatalf("no service %s: %+v", dir, r.Config.Services)
	return Service{}
}
func nativeAssertUninferred(t *testing.T, s Service) {
	t.Helper()
	if s.Language != "custom" || len(s.Commands) != 0 || s.Manager != "" || s.Database.Tool != "" || s.Infrastructure.Kind != "local" || len(s.Pending) == 0 {
		t.Errorf("native commands/infra guessed: %+v", s)
	}
}
func TestNativeSourcesMixedMakeTargetsRemainOptInAndScanOffline(t *testing.T) {
	root := nativeFixture(t, map[string]string{"Makefile": "$(shell printf SHOULD_NOT_RUN > marker)\nuser-check:\n\tprintf NEVER\ncustom-build: user-check\n\tprintf NEVER\n", "src/main.c": "not parsed or compiled", "src/util.cpp": "not parsed or compiled", "include/api.h": "header"})
	tool := t.TempDir()
	marker := filepath.Join(tool, "executed")
	for _, name := range []string{"make", "gcc", "g++", "cmake", "meson"} {
		if e := os.WriteFile(filepath.Join(tool, name), []byte("#!/bin/sh\nprintf executed > '"+marker+"'\nexit 89\n"), 0755); e != nil {
			t.Fatal(e)
		}
	}
	t.Setenv("PATH", tool)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 {
		t.Fatalf("mixed sources duplicated: %+v", r.Config.Services)
	}
	s := nativeService(t, r, ".")
	nativeAssertUninferred(t, s)
	if !strings.Contains(strings.Join(s.Pending, " "), "Make") {
		t.Error("explicit Make contract absent")
	}
	names := map[string]bool{}
	for _, su := range r.Suggestions {
		names[su.Name] = true
		if su.Service != s.ID {
			t.Errorf("wrong Make owner: %+v", su)
		}
	}
	if !names["make-user-check"] || !names["make-custom-build"] {
		t.Errorf("user-defined targets missing: %+v", r.Suggestions)
	}
	for _, p := range []string{marker, filepath.Join(root, "marker")} {
		if _, e := os.Stat(p); !os.IsNotExist(e) {
			t.Error("scan executed tools or Make shell expression")
		}
	}
	kinds := map[string]bool{}
	for _, ev := range r.Evidence {
		kinds[ev.Kind] = true
	}
	if !kinds["c-source"] || !kinds["cpp-source"] || !kinds["c-header"] {
		t.Errorf("evidence missing %+v", r.Evidence)
	}
	again, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(r, again) {
		t.Error("native scan not deterministic")
	}
}
func TestNativeSourcesNearestMakeAndUnrelatedDirectories(t *testing.T) {
	root := nativeFixture(t, map[string]string{"a/Makefile": "a-check:\n\ttrue\n", "a/src/main.c": "", "a/deeper/Makefile": "deep-check:\n\ttrue\n", "a/deeper/src/main.cc": "", "b/main.cxx": "", "c/main.cpp": ""})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 4 {
		t.Fatalf("anchors: %+v", r.Config.Services)
	}
	for _, dir := range []string{"a", "a/deeper", "b", "c"} {
		nativeAssertUninferred(t, nativeService(t, r, dir))
	}
	for _, su := range r.Suggestions {
		want := "a"
		if su.Name == "make-deep-check" {
			want = "a/deeper"
		}
		if su.Service != nativeService(t, r, want).ID {
			t.Errorf("target owner %+v", su)
		}
	}
}
func TestNativeSourcesExistingLanguageOwnerAndNestedMakeBoundary(t *testing.T) {
	root := nativeFixture(t, map[string]string{"package.json": `{"scripts":{"test":"NEVER"}}`, "native/main.c": "", "native/sub/Makefile": "native-check:\n\ttrue\n", "native/sub/src/main.cpp": ""})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 2 {
		t.Fatalf("existing service duplicated: %+v", r.Config.Services)
	}
	js := nativeService(t, r, ".")
	if js.Language != "javascript" || !reflect.DeepEqual(js.Commands["test"].Args, []string{"npm", "run", "test"}) {
		t.Errorf("JS changed: %+v", js)
	}
	native := nativeService(t, r, "native/sub")
	nativeAssertUninferred(t, native)
	for _, su := range r.Suggestions {
		if su.Name == "make-native-check" && su.Service != native.ID {
			t.Errorf("nested boundary lost: %+v", su)
		}
	}
}
func TestNativeSourcesExcludeSymlinkAndHeadersOnly(t *testing.T) {
	root := nativeFixture(t, map[string]string{"kept/main.c": "", "private/main.cpp": "", "headers/api.hpp": "", "headers/api.h": ""})
	outside := nativeFixture(t, map[string]string{"main.c": "outside"})
	if e := os.Symlink(filepath.Join(outside, "main.c"), filepath.Join(root, "linked.cpp")); e != nil {
		t.Fatal(e)
	}
	r, e := Scan(root, []string{"private"})
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 || r.Config.Services[0].Dir != "kept" {
		t.Fatalf("excluded/header/symlink creates service: %+v", r.Config.Services)
	}
	for _, ev := range r.Evidence {
		if strings.Contains(ev.Path, "private") || strings.Contains(ev.Path, "linked") {
			t.Errorf("excluded source evidence: %+v", ev)
		}
	}
}
func TestNativeSourcesIDCollisionDoesNotOverwriteService(t *testing.T) {
	root := nativeFixture(t, map[string]string{"a/b/main.c": "", "a-b/main.cpp": ""})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 {
		t.Errorf("colliding IDs retained: %+v", r.Config.Services)
	}
	if !strings.Contains(strings.Join(r.Warnings, " "), "ambiguous") {
		t.Error("ID ambiguity missing")
	}
	kinds := 0
	for _, ev := range r.Evidence {
		if ev.Kind == "c-source" || ev.Kind == "cpp-source" {
			kinds++
		}
	}
	if kinds != 2 {
		t.Error("conflict discarded source evidence")
	}
}

func TestNativeSourcesSpecialFilesAreRejectedWithoutOpening(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix socket fixture")
	}
	root := nativeFixture(t, map[string]string{"Makefile": "custom-check:\n\ttrue\n"})
	listener, e := net.Listen("unix", filepath.Join(root, "main.c"))
	if e != nil {
		t.Fatal(e)
	}
	defer listener.Close()
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, ev := range r.Evidence {
		if ev.Kind == "c-source" || ev.Kind == "cpp-source" {
			t.Errorf("special file accepted: %+v", ev)
		}
	}
	if !strings.Contains(strings.Join(r.Warnings, " "), "regular file") {
		t.Error("special source lacks diagnostic")
	}
}
func TestNativeSourcesUnownedMakefileDoesNotUseSiblingService(t *testing.T) {
	root := nativeFixture(t, map[string]string{"a/src/main.c": "", "b/Makefile": "other-check:\n\ttrue\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	nativeAssertUninferred(t, nativeService(t, r, "a/src"))
	for _, su := range r.Suggestions {
		if su.Name == "make-other-check" {
			t.Errorf("unrelated Make target assigned to sibling: %+v", su)
		}
	}
	if !strings.Contains(strings.Join(r.Warnings, " "), "unowned") {
		t.Error("unowned Makefile not disclosed")
	}
}
