package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePreflightProtectsAllDestinations(t *testing.T) {
	for _, bad := range []string{"edited-launcher", "invalid-manifest"} {
		t.Run(bad, func(t *testing.T) {
			root := t.TempDir()
			c := sample()
			var out bytes.Buffer
			if err := Generate(root, c, false, &out); err != nil {
				t.Fatal(err)
			}
			target := "menu.sh"
			if bad == "invalid-manifest" {
				target = ".mudarro/generated.json"
			}
			full := filepath.Join(root, target)
			if err := os.WriteFile(full, []byte("invalid user content"), 0644); err != nil {
				t.Fatal(err)
			}
			c.Services[0].Commands["new-action"] = cmd("scripts", "echo", "new")
			if err := Generate(root, c, false, &out); err == nil {
				t.Fatal("accepted protected content")
			}
			if exists(filepath.Join(root, ".mudarro/scripts/app/new-action.sh")) {
				t.Fatal("partial generation before preflight error")
			}
			got, err := os.ReadFile(full)
			if err != nil || string(got) != "invalid user content" {
				t.Fatalf("modified user content: %q %v", got, err)
			}
		})
	}
}

func TestInitSelectionFailureDoesNotSave(t *testing.T) {
	for _, tc := range []struct {
		name  string
		args  []string
		input string
	}{
		{"unknown-selection", []string{"--select", "app-javascript:missing"}, ""},
		{"interactive-eof", []string{"--interactive"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": `{"scripts":{"hello":"echo hello"}}`})
			args := append([]string{"init", "--root", root}, tc.args...)
			var out bytes.Buffer
			if err := Main(args, "test", strings.NewReader(tc.input), &out); err == nil {
				t.Fatal("accepted incomplete selection")
			}
			if exists(filepath.Join(root, "mudarro.yaml")) {
				t.Fatal("saved configuration on error")
			}
		})
	}
}
