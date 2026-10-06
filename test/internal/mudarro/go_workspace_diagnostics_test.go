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

func goDiagnosticFixture(t *testing.T, mode string) string {
	t.Helper()
	root := t.TempDir()
	c := Config{Version: 1, Name: "Go diagnostic fixture", Services: []Service{{ID: "app", Dir: ".", Language: "go", GoWorkspace: mode, Infrastructure: Infrastructure{Kind: "custom"}, Commands: map[string]Command{"test": {Args: []string{"go", "test", "./..."}, Group: "qualidade"}}}}}
	if e := saveConfig(filepath.Join(root, "mudarro.yaml"), c); e != nil {
		t.Fatal(e)
	}
	return root
}

func TestGoWorkspaceDiagnosticsPreviewShowsEffectiveIsolation(t *testing.T) {
	root := goDiagnosticFixture(t, "off")
	before := supportFileHashes(t, root)
	var out bytes.Buffer
	if e := Main([]string{"preview", "app:test", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(out.String(), "GOWORK=off") || !strings.Contains(out.String(), `"env"`) {
		t.Errorf("preview hides effective environment argv: %q", out.String())
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("preview changed config/project")
	}
}

func TestGoWorkspaceDiagnosticsDryRunShowsIsolationWithoutExecution(t *testing.T) {
	root := goDiagnosticFixture(t, "off")
	tools := t.TempDir()
	t.Setenv("PATH", tools)
	for _, tool := range []string{"go", "env"} {
		if e := os.WriteFile(filepath.Join(tools, tool), []byte("#!/bin/sh\nprintf executed > '"+filepath.Join(root, "MUST_NOT_EXECUTE")+"'\nexit 91\n"), 0755); e != nil {
			t.Fatal(e)
		}
	}
	before := supportFileHashes(t, root)
	var out bytes.Buffer
	if e := Main([]string{"run", "app:test", "--root", root, "--dry-run"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(out.String(), "GOWORK=off") || !strings.Contains(out.String(), `"env"`) {
		t.Errorf("dry run hides effective environment argv: %q", out.String())
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("dry run executed tool or changed config")
	}
}

func TestGoWorkspaceDiagnosticsDoctorRequiresEnvOnlyForIsolation(t *testing.T) {
	tools := t.TempDir()
	t.Setenv("PATH", tools)
	for _, mode := range []string{"off", "", "inherit"} {
		t.Run(mode, func(t *testing.T) {
			root := goDiagnosticFixture(t, mode)
			if e := os.WriteFile(filepath.Join(tools, "go"), []byte("#!/bin/sh\nprintf executed > '"+filepath.Join(root, "MUST_NOT_EXECUTE")+"'\n"), 0755); e != nil {
				t.Fatal(e)
			}
			before := supportFileHashes(t, root)
			var out bytes.Buffer
			e := Main([]string{"doctor", "--root", root}, "test", strings.NewReader(""), &out)
			if mode == "off" {
				if e == nil || !strings.Contains(out.String(), "AUSENTE app: env") {
					t.Errorf("env absent was not diagnosed: %v %q", e, out.String())
				}
			} else if e != nil || strings.Contains(out.String(), "app: env") {
				t.Errorf("inherited context gained env dependency: %v %q", e, out.String())
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("doctor executed application or changed config")
			}
		})
	}
}
