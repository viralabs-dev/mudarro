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

func TestPackageManagerValidationKnownNamesPreserveContracts(t *testing.T) {
	for _, manager := range []string{"npm", "pnpm", "yarn", "bun"} {
		t.Run(manager, func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{"packageManager": manager, "scripts": map[string]string{"start": "NEVER_EXECUTE", "test": "NEVER_EXECUTE", "hello": "NEVER_EXECUTE"}})
			root := fixture(t, map[string]string{"package.json": string(body), "package-lock.json": "fixture", "yarn.lock": "fixture"})
			before := supportFileHashes(t, root)
			r, err := Scan(root, nil)
			if err != nil {
				t.Fatal(err)
			}
			s := r.Config.Services[0]
			if s.Manager != manager {
				t.Fatalf("known manager %+v", s)
			}
			if !reflect.DeepEqual(s.Commands["test"].Args, []string{manager, "run", "test"}) {
				t.Fatal("explicit override changed argv")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("scan mutated project")
			}
		})
	}
}
func TestPackageManagerValidationBadNamesDoNotPlanExecution(t *testing.T) {
	for _, value := range []string{"unknown", "../manager", "/tmp/manager", "https://example.invalid/tool", "npm ", " npm", "npm@", "@1.2.3", "npm\t", "npm\n"} {
		t.Run(strings.ReplaceAll(value, "/", "_"), func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{"packageManager": value, "scripts": map[string]string{"start": "NEVER_EXECUTE", "test": "NEVER_EXECUTE", "hello": "NEVER_EXECUTE"}})
			root := fixture(t, map[string]string{"bad/package.json": string(body), "good/package.json": `{"packageManager":"npm","scripts":{"test":"NEVER_EXECUTE"}}`})
			before := supportFileHashes(t, root)
			r, err := Scan(root, nil)
			if err != nil {
				t.Fatal(err)
			}
			good := false
			diagnostic := strings.Join(r.Warnings, " ")
			for _, s := range r.Config.Services {
				if s.Dir == "good" {
					good = true
					if s.Manager != "npm" {
						t.Fatal("valid sibling corrupted")
					}
				}
				if s.Dir == "bad" {
					diagnostic += " " + strings.Join(s.Pending, " ")
					if s.Manager != "" {
						t.Errorf("invalid field chosen as executable: %q", s.Manager)
					}
					if len(s.Commands) > 0 {
						t.Errorf("invalid field generated commands: %+v", s.Commands)
					}
				}
			}
			if !good {
				t.Fatal("invalid manager dropped valid sibling")
			}
			if !strings.Contains(strings.ToLower(diagnostic), "packagemanager") {
				t.Errorf("missing actionable packageManager diagnostic: %s", diagnostic)
			}
			for _, su := range r.Suggestions {
				if su.Service == "bad-javascript" && len(su.Command.Args) > 0 && su.Command.Args[0] != "" {
					t.Errorf("invalid manager executable suggestion: %+v", su)
				}
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("invalid scan modified project")
			}
		})
	}
}
func TestPackageManagerValidationWrongJSONTypesDiagnosed(t *testing.T) {
	for _, value := range []string{"123", "true", "[]", "{}"} {
		t.Run(value, func(t *testing.T) {
			root := fixture(t, map[string]string{"bad/package.json": `{"packageManager":` + value + `,"scripts":{"start":"NEVER_EXECUTE"}}`, "good/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
			r, err := Scan(root, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Warnings) == 0 {
				t.Fatal("invalid JSON field type silent")
			}
			if len(r.Config.Services) != 1 || r.Config.Services[0].Dir != "good" {
				t.Fatalf("valid sibling lost %+v", r.Config.Services)
			}
		})
	}
}
func TestPackageManagerValidationInitInvalidSelectionPreservesProject(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"packageManager":"unknown","scripts":{"hello":"touch NEVER_EXECUTE"}}`, "keep.txt": "user data"})
	before := supportFileHashes(t, root)
	var out bytes.Buffer
	if err := Main([]string{"init", "--root", root, "--select", "app-javascript:hello"}, "test", strings.NewReader(""), &out); err == nil {
		t.Error("invalid manager selected at init")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Error("invalid init changed files")
	}
}
func TestPackageManagerValidationScanPreservesExistingGeneratedFiles(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"scripts":{"start":"NEVER_EXECUTE","test":"NEVER_EXECUTE"}}`})
	var out bytes.Buffer
	if err := Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); err != nil {
		t.Fatal(err)
	}
	if err := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"packageManager":"unknown","scripts":{"start":"NEVER_EXECUTE"}}`), 0644); err != nil {
		t.Fatal(err)
	}
	before := supportFileHashes(t, root)
	r, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	diagnostic := strings.Join(r.Warnings, " ")
	for _, s := range r.Config.Services {
		diagnostic += " " + strings.Join(s.Pending, " ")
	}
	if !strings.Contains(strings.ToLower(diagnostic), "packagemanager") {
		t.Error("invalid field was not diagnosed")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("invalid scan changed generated files/config")
	}
}

func TestPackageManagerValidationExplicitConfigurationKeepsCustomCommands(t *testing.T) {
	root := t.TempDir()
	c := sample()
	c.Services[0].Language = "javascript"
	c.Services[0].Manager = "custom-tool"
	c.Services[0].Commands["test"] = cmd("qualidade", "sh", "-c", "printf custom-manager-safe")
	if err := saveConfig(filepath.Join(root, "mudarro.yaml"), c); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err = Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err = Main([]string{"run", "app:test", "--root", root}, "test", strings.NewReader(""), &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "custom-manager-safe" {
		t.Fatalf("custom configured command altered: %q", out.String())
	}
	after, _ := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if !bytes.Equal(before, after) {
		t.Fatal("configured custom manager was rewritten")
	}
}
