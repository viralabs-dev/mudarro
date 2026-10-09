package mudarro_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/viralabs-dev/mudarro/internal/mudarro"
)

func TestOwnedErrorsRespectLocale(t *testing.T) {
	for _, tc := range []struct{ locale, command, want string }{
		{"en", "unknown", "unknown command: unknown"},
		{"pt-BR", "unknown", "comando desconhecido: unknown"},
		{"en", "run", "action does not exist: app:missing"},
		{"pt-BR", "run", "ação inexistente: app:missing"},
	} {
		t.Run(tc.locale+"/"+tc.command, func(t *testing.T) {
			root := t.TempDir()
			config := "version: 1\nname: Aurora\nui:\n  locale: " + tc.locale + "\nservices:\n  - id: app\n    dir: .\n    infrastructure: {kind: custom}\n"
			if err := os.WriteFile(filepath.Join(root, "mudarro.yaml"), []byte(config), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{tc.command, "--root", root}
			if tc.command == "unknown" {
				args = append(args, "--locale", tc.locale)
			} else {
				args = append(args, "app:missing")
			}
			err := mudarro.Main(args, "test", strings.NewReader(""), &bytes.Buffer{})
			if err == nil || err.Error() != tc.want {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
}

func TestInvalidConfigErrorLocaleAndExternalOutput(t *testing.T) {
	skipOnWindows(t, "fixture runs sh")
	for _, tc := range []struct{ locale, want string }{{"en", "version must be 1"}, {"pt-BR", "version deve ser 1"}} {
		t.Run(tc.locale, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "mudarro.yaml"), []byte("version: 2\nname: Aurora\nservices: []\n"), 0600); err != nil {
				t.Fatal(err)
			}
			err := mudarro.Main([]string{"generate", "--root", root, "--locale", tc.locale}, "test", strings.NewReader(""), &bytes.Buffer{})
			if err == nil || err.Error() != tc.want {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
	root := t.TempDir()
	// The external output deliberately contains catalog phrases in both languages.
	config := `version: 1
name: Aurora
services:
  - id: app
    dir: .
    infrastructure: {kind: custom}
    commands:
      probe:
        args: [sh, -c, 'printf "%s\n" "comando vazio|empty command"; exit 7']
        group: scripts
`
	if err := os.WriteFile(filepath.Join(root, "mudarro.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err := mudarro.Main([]string{"run", "app:probe", "--root", root}, "test", strings.NewReader(""), &out)
	if out.String() != "comando vazio|empty command\n" {
		t.Fatalf("external output changed: %q", out.String())
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 7 {
		t.Fatalf("original execution error lost: %v", err)
	}
}
