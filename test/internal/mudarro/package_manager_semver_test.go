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

func semverManifest(value string) string {
	body, _ := json.Marshal(map[string]any{"packageManager": value, "scripts": map[string]string{"start": "touch NEVER_EXECUTE", "test": "touch NEVER_EXECUTE", "hello": "touch NEVER_EXECUTE"}})
	return string(body)
}

func TestPackageManagerSemverAcceptsStrictVersions(t *testing.T) {
	for _, value := range []string{
		"npm@1.2.3", "pnpm@0.0.0", "yarn@4.18.1", "bun@1.4.0",
		"npm@1.2.3-alpha", "npm@1.2.3-alpha.0", "npm@1.2.3-0", "npm@1.2.3-0a", "npm@1.2.3--", "npm@1.2.3-a-b.10",
		"npm@1.2.3+001.02", "npm@1.2.3+build-meta.9", "npm@1.2.3-alpha.1+build.001",
		"yarn@4.18.1+sha224.abcdef0123456789", "npm@1.2.3+sha224.nothex",
		"npm@18446744073709551616.0.0",
	} {
		t.Run(value, func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": semverManifest(value)})
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			manager := strings.SplitN(value, "@", 2)[0]
			if len(r.Config.Services) != 1 || r.Config.Services[0].Manager != manager {
				t.Fatalf("valid SemVer rejected: %+v", r)
			}
			if !reflect.DeepEqual(r.Config.Services[0].Commands["test"].Args, []string{manager, "run", "test"}) {
				t.Fatal("version leaked into executable argv")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("scan modified manifest or ran script")
			}
		})
	}
}

func TestPackageManagerSemverRejectsOpaqueOrMalformedVersions(t *testing.T) {
	for _, version := range []string{
		"1", "1.2", "1.2.3.4", "01.2.3", "1.02.3", "1.2.03", "v1.2.3", "^1.2.3", "~1.2.3", "latest", "https://example.invalid/tool", "1.2.3@extra",
		"1.2.3-", "1.2.3-01", "1.2.3-alpha.01", "1.2.3-a..b", "1.2.3-alpha_1", "1.2.3-α", "1.2.3+", "1.2.3+a..b", "1.2.3+build_1", "1.2.3+é", "1.2.3++build", "1.2.3\t", "1.2.3\n", "1.2.3+build metadata",
	} {
		t.Run(strings.ReplaceAll(version, "/", "_"), func(t *testing.T) {
			root := fixture(t, map[string]string{"bad/package.json": semverManifest("npm@" + version), "good/package.json": semverManifest("bun@1.4.0")})
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			good := false
			diagnostic := strings.Join(r.Warnings, " ")
			for _, s := range r.Config.Services {
				if s.Dir == "good" {
					good = s.Manager == "bun"
				}
				if s.Dir == "bad" {
					diagnostic += " " + strings.Join(s.Pending, " ")
					if s.Manager != "" || len(s.Commands) > 0 {
						t.Errorf("invalid version planned execution: %+v", s)
					}
				}
			}
			if !good {
				t.Fatal("invalid version hid valid sibling")
			}
			if !strings.Contains(strings.ToLower(diagnostic), "packagemanager") {
				t.Errorf("missing field diagnostic: %q", diagnostic)
			}
			for _, su := range r.Suggestions {
				if su.Service == "bad-javascript" && len(su.Command.Args) > 0 && su.Command.Args[0] != "" {
					t.Errorf("invalid version executable suggestion: %+v", su)
				}
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("invalid scan changed files")
			}
		})
	}
}

func TestPackageManagerSemverInvalidInitAndScanPreserveArtifacts(t *testing.T) {
	t.Run("selected-script-rejected-before-config", func(t *testing.T) {
		root := fixture(t, map[string]string{"package.json": semverManifest("npm@latest"), "keep.txt": "user contents"})
		before := supportFileHashes(t, root)
		var out bytes.Buffer
		if e := Main([]string{"init", "--root", root, "--select", "app-javascript:hello"}, "test", strings.NewReader(""), &out); e == nil {
			t.Error("invalid SemVer selected at init")
		}
		if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
			t.Fatal("invalid init wrote config or changed fixture")
		}
	})
	t.Run("new-invalid-field-preserves-generated", func(t *testing.T) {
		root := fixture(t, map[string]string{"package.json": semverManifest("npm@1.2.3")})
		var out bytes.Buffer
		for _, command := range []string{"init", "generate"} {
			if e := Main([]string{command, "--root", root}, "test", strings.NewReader(""), &out); e != nil {
				t.Fatal(e)
			}
		}
		if e := os.WriteFile(filepath.Join(root, "package.json"), []byte(semverManifest("npm@1.2.3-alpha.01")), 0644); e != nil {
			t.Fatal(e)
		}
		before := supportFileHashes(t, root)
		r, e := Scan(root, nil)
		if e != nil {
			t.Fatal(e)
		}
		diagnostic := strings.Join(r.Warnings, " ")
		for _, s := range r.Config.Services {
			diagnostic += " " + strings.Join(s.Pending, " ")
		}
		if !strings.Contains(strings.ToLower(diagnostic), "packagemanager") {
			t.Error("new bad version silent")
		}
		if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
			t.Fatal("scan rewrote generated wrappers/config")
		}
	})
}
