package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func rustDetect(t *testing.T, root, dir string, excludes []string) (*Service, []Suggestion, error) {
	t.Helper()
	return (adapters.Rust{}).Detect(root, dir, map[string]bool{"Cargo.toml": true}, excludes)
}
func TestRustCargoPackageLibAndBinsStayOptIn(t *testing.T) {
	for _, manifest := range []string{"[package]\nname='owned'\nversion='0.1.0'\n", "[package]\nname='owned'\n[lib]\nname='owned'\npath='src/lib.rs'\n[[bin]]\nname='owned'\npath='src/main.rs'\n[[bin]]\nname='other'\npath='src/other.rs'\n"} {
		t.Run(manifest, func(t *testing.T) {
			root := fixture(t, map[string]string{"Cargo.toml": manifest, "src/lib.rs": "pub fn value()->i32{21}", "src/main.rs": "fn main(){}", "src/other.rs": "fn main(){}", "build.rs": "fn main(){panic!(\"NEVER EXECUTE\");}"})
			before := supportFileHashes(t, root)
			s, su, e := rustDetect(t, root, ".", nil)
			if e != nil {
				t.Fatal(e)
			}
			if s.Language != "rust" || s.Manager != "cargo" || len(s.Commands) != 0 || s.WorkspaceRoot != "" || len(su) != 2 {
				t.Fatalf("unexpected inference %+v %+v", s, su)
			}
			for i, name := range []string{"build", "test"} {
				if su[i].Name != name || su[i].Purpose != "cargo" || !reflect.DeepEqual(su[i].Command.Args, []string{"cargo", name}) {
					t.Errorf("opt-in command %+v", su[i])
				}
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Error("scan changed or executed source")
			}
		})
	}
}
func TestRustCargoAbsentMalformedAndOversized(t *testing.T) {
	s, su, e := (adapters.Rust{}).Detect(t.TempDir(), ".", nil, nil)
	if s != nil || su != nil || e != nil {
		t.Fatal("absent Cargo should not detect")
	}
	for _, tc := range []struct{ name, text string }{{"syntax", "[package"}, {"empty", ""}, {"missing-name", "[package]\nversion='1.0.0'"}, {"wrong-type", "[package]\nname=42"}, {"duplicate", "[package]\nname='a'\nname='b'"}, {"oversized", strings.Repeat("#", 4*1024*1024+1)}} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"Cargo.toml": tc.text})
			s, su, e := rustDetect(t, root, ".", nil)
			if e == nil || s != nil || len(su) != 0 {
				t.Errorf("invalid manifest accepted: %+v %+v %v", s, su, e)
			}
		})
	}
}
func TestRustCargoWorkspaceExplicitOwnershipAndDefaults(t *testing.T) {
	root := fixture(t, map[string]string{"Cargo.toml": "[workspace]\nmembers=['packages/lib','packages/app']\ndefault-members=['packages/app']\n", "packages/lib/Cargo.toml": "[package]\nname='owned-lib'\n", "packages/app/Cargo.toml": "[package]\nname='owned-app'\n", "outside/Cargo.toml": "[package]\nname='independent'\n"})
	s, su, e := rustDetect(t, root, ".", nil)
	if e != nil || len(su) != 2 || !strings.Contains(strings.Join(s.Pending, " "), "packages/app, packages/lib") {
		t.Fatalf("workspace validation %+v %+v %v", s, su, e)
	}
	for _, dir := range []string{"packages/app", "packages/lib"} {
		member, su, e := rustDetect(t, root, dir, nil)
		if e != nil || len(su) != 2 || !strings.Contains(strings.Join(member.Pending, " "), "workspace owner: .") || member.Dir != dir || member.WorkspaceRoot != "" {
			t.Errorf("member ownership %+v %+v %v", member, su, e)
		}
	}
	outside, _, e := rustDetect(t, root, "outside", nil)
	if e != nil || strings.Contains(strings.Join(outside.Pending, " "), "workspace owner") {
		t.Errorf("nonmember ownership inherited %+v %v", outside, e)
	}
}
func TestRustCargoWorkspaceInvalidMembersStayUnresolved(t *testing.T) {
	for _, tc := range []struct {
		name, extra string
		excludes    []string
	}{{"missing", "members=['missing']", nil}, {"outside", "members=['../outside']", nil}, {"glob", "members=['packages/*']", nil}, {"not-package", "members=['empty']", nil}, {"duplicate", "members=['valid','valid/../valid']", nil}, {"default-nonmember", "members=['valid']\ndefault-members=['other']", nil}, {"excluded-dir", "members=['valid']", []string{"valid"}}, {"excluded-file", "members=['valid']", []string{"valid/Cargo.toml"}}, {"target-artifacts", "members=['target/member']", nil}, {"same-package-name", "members=['valid','other']", nil}} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"Cargo.toml": "[workspace]\n" + tc.extra + "\n", "valid/Cargo.toml": "[package]\nname='same'\n", "other/Cargo.toml": "[package]\nname='same'\n", "empty/Cargo.toml": "[workspace]\n", "target/member/Cargo.toml": "[package]\nname='target'\n"})
			before := supportFileHashes(t, root)
			s, su, e := rustDetect(t, root, ".", tc.excludes)
			if e != nil || s == nil || len(su) != 0 || len(s.Pending) < 2 {
				t.Errorf("workspace ambiguity executable: %+v %+v %v", s, su, e)
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Error("member probing changed files")
			}
		})
	}
}
func TestRustCargoManifestAndTargetSymlinksAndExclusions(t *testing.T) {
	outside := fixture(t, map[string]string{"Cargo.toml": "[package]\nname='outside'\n", "main.rs": "outside"})
	root := t.TempDir()
	if e := os.Symlink(filepath.Join(outside, "Cargo.toml"), filepath.Join(root, "Cargo.toml")); e != nil {
		t.Fatal(e)
	}
	if _, _, e := rustDetect(t, root, ".", nil); e == nil {
		t.Error("symlink manifest accepted")
	}
	root = fixture(t, map[string]string{"Cargo.toml": "[package]\nname='owned'\n[[bin]]\nname='owned'\npath='main.rs'\n"})
	if e := os.Symlink(filepath.Join(outside, "main.rs"), filepath.Join(root, "main.rs")); e != nil {
		t.Fatal(e)
	}
	s, su, e := rustDetect(t, root, ".", nil)
	if e != nil || len(su) != 0 || len(s.Pending) < 2 {
		t.Error("symlink target executable")
	}
	s, su, e = rustDetect(t, root, ".", []string{"Cargo.toml"})
	if e != nil || s != nil || su != nil {
		t.Error("excluded manifest read")
	}
}
func TestRustCargoDuplicateAndOutsideExplicitTargets(t *testing.T) {
	for _, tail := range []string{"[[bin]]\nname='same'\n[[bin]]\nname='same'\n", "[[bin]]\nname='outside'\npath='../outside.rs'\n", "[lib]\npath='/tmp/outside.rs'\n"} {
		root := fixture(t, map[string]string{"Cargo.toml": "[package]\nname='owned'\n" + tail})
		s, su, e := rustDetect(t, root, ".", nil)
		if e != nil || len(su) != 0 || len(s.Pending) < 2 {
			t.Errorf("unresolved targets executable: %+v %+v %v", s, su, e)
		}
	}
}
func TestRustCargoScanOfflineSelectionGenerateAndManualPreservation(t *testing.T) {
	root := fixture(t, map[string]string{"Cargo.toml": "[package]\nname='owned'\n", "src/lib.rs": "pub fn value()->i32{21}", "build.rs": "NEVER EXECUTE"})
	tools := t.TempDir()
	marker := filepath.Join(tools, "EXECUTED")
	for _, name := range []string{"cargo", "rustc"} {
		if e := os.WriteFile(filepath.Join(tools, name), []byte("#!/bin/sh\nprintf executed > '"+marker+"'\nexit 88\n"), 0755); e != nil {
			t.Fatal(e)
		}
	}
	t.Setenv("PATH", tools)
	var out, errout bytes.Buffer
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	var rust Service
	for _, s := range r.Config.Services {
		if s.Language == "rust" {
			rust = s
		}
	}
	if rust.ID == "" {
		t.Fatal("Rust not registered")
	}
	run := func(args ...string) {
		t.Helper()
		out.Reset()
		errout.Reset()
		if Main(args, "test", strings.NewReader(""), &out) != nil {
			t.Fatalf("%v: %s", args, errout.String())
		}
	}
	run("init", "--root", root, "--select", rust.ID+":build,"+rust.ID+":test")
	run("generate", "--root", root)
	one := supportFileHashes(t, root)
	run("generate", "--root", root)
	if !reflect.DeepEqual(one, supportFileHashes(t, root)) {
		t.Error("generate not idempotent")
	}
	run("run", "--root", root, rust.ID+":test", "--dry-run")
	cfg, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	cfg.Services[0].Commands["start"] = Command{Args: []string{"cargo", "run", "--bin", "user-chosen"}, Group: "aplicacao"}
	data, e := yaml.Marshal(cfg)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(root, "mudarro.yaml"), data, 0644); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, root)
	if Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out) == nil {
		t.Error("init overwrote manual config")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Error("manual startup changed")
	}
	if _, e = os.Stat(marker); !os.IsNotExist(e) {
		t.Error("scan/init/generate/dry-run executed tool")
	}
}

func TestRustCargoExclusionsPrecedeWorkspaceReads(t *testing.T) {
	t.Run("excluded-invalid-ancestor", func(t *testing.T) {
		root := fixture(t, map[string]string{"Cargo.toml": "[workspace]\nmembers=['member']\nBROKEN TOML", "member/Cargo.toml": "[package]\nname='owned'\n"})
		s, su, e := rustDetect(t, root, "member", []string{"Cargo.toml"})
		if e != nil || len(su) != 2 || !strings.Contains(strings.Join(s.Pending, " "), "ancestor manifest excluded") || strings.Contains(strings.Join(s.Pending, " "), "workspace owner:") {
			t.Errorf("excluded ancestor inspected/inherited: %+v %+v %v", s, su, e)
		}
	})
	t.Run("declared-excluded-oversize", func(t *testing.T) {
		root := fixture(t, map[string]string{"Cargo.toml": "[workspace]\nmembers=['private']\nexclude=['private']\n", "private/Cargo.toml": strings.Repeat("#", 4*1024*1024+1)})
		s, su, e := rustDetect(t, root, ".", nil)
		pending := strings.Join(s.Pending, " ")
		if e != nil || len(su) != 0 || !strings.Contains(pending, "excluded") || strings.Contains(pending, "grande demais") {
			t.Errorf("excluded member was read: %+v %+v %v", s, su, e)
		}
	})
	t.Run("member-symlink", func(t *testing.T) {
		outside := fixture(t, map[string]string{"Cargo.toml": "[package]\nname='outside'\n"})
		root := fixture(t, map[string]string{"Cargo.toml": "[workspace]\nmembers=['linked']\n"})
		if e := os.Symlink(outside, filepath.Join(root, "linked")); e != nil {
			t.Fatal(e)
		}
		s, su, e := rustDetect(t, root, ".", nil)
		if e != nil || len(su) != 0 || !strings.Contains(strings.Join(s.Pending, " "), "symlink") {
			t.Errorf("outside member symlink followed: %+v %+v %v", s, su, e)
		}
	})
}

func TestRustCargoVirtualWorkspaceAndAppMemberIDsAllowSelection(t *testing.T) {
	root := fixture(t, map[string]string{"Cargo.toml": "[workspace]\nmembers=['app']\n", "app/Cargo.toml": "[package]\nname='owned-app'\n"})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 2 {
		t.Fatalf("workspace/member lost: %+v", r.Config.Services)
	}
	ids := map[string]bool{}
	for _, s := range r.Config.Services {
		if ids[s.ID] {
			t.Fatalf("duplicate IDs %+v", r.Config.Services)
		}
		ids[s.ID] = true
	}
	selectKeys := []string{}
	for _, su := range r.Suggestions {
		if su.Purpose == "cargo" {
			selectKeys = append(selectKeys, su.Service+":"+su.Name)
		}
	}
	if len(selectKeys) != 4 {
		t.Fatalf("root/member suggestions %+v", r.Suggestions)
	}
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--select", strings.Join(selectKeys, ",")}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	cfg, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	if len(cfg.Services) != 2 {
		t.Fatal("selected workspace/member config lost")
	}
	for _, s := range cfg.Services {
		if !reflect.DeepEqual(s.Commands["test"].Args, []string{"cargo", "test"}) {
			t.Errorf("commands owner lost %+v", s)
		}
	}
}

func TestRustCargoExplicitPackageWorkspaceRemainsUnresolvedWithoutRead(t *testing.T) {
	for _, tc := range []struct {
		name, owner string
		excludes    []string
	}{{"outside", "../../outside", nil}, {"contained", "owner", nil}, {"excluded", "owner", []string{"owner"}}, {"symlink", "linked", nil}} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"Cargo.toml": "[package]\nname='owned'\nworkspace='" + tc.owner + "'\n", "owner/Cargo.toml": strings.Repeat("#", 4*1024*1024+1)})
			outside := fixture(t, map[string]string{"Cargo.toml": "MUST_NOT_BE_READ"})
			if tc.name == "symlink" {
				if e := os.Symlink(outside, filepath.Join(root, "linked")); e != nil {
					t.Fatal(e)
				}
			}
			hashFiles := func() map[string]string {
				result := map[string]string{}
				for _, name := range []string{"Cargo.toml", "owner/Cargo.toml"} {
					data, err := os.ReadFile(filepath.Join(root, name))
					if err != nil {
						t.Fatal(err)
					}
					result[name] = string(data)
				}
				return result
			}
			before := hashFiles()
			s, su, e := rustDetect(t, root, ".", tc.excludes)
			if e != nil || s == nil || len(su) != 0 {
				t.Fatalf("explicit owner executable %+v %+v %v", s, su, e)
			}
			pending := strings.Join(s.Pending, " ")
			if !strings.Contains(pending, "package.workspace is not resolved") || strings.Contains(pending, "workspace owner:") || strings.Contains(pending, "grande demais") {
				t.Errorf("explicit owner ignored or read: %s", pending)
			}
			if !reflect.DeepEqual(before, hashFiles()) {
				t.Error("explicit owner probing changed files")
			}
		})
	}
}
