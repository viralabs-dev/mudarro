package mudarro_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/viralabs-dev/mudarro/internal/mudarro"
)

func doctorFixture(t *testing.T, infrastructure mudarro.Infrastructure, commands map[string]mudarro.Command) string {
	t.Helper()
	root := t.TempDir()
	c := mudarro.Config{Version: 1, Name: "Doctor fixture", Services: []mudarro.Service{{ID: "app", Dir: ".", Infrastructure: infrastructure, Commands: commands}}}
	if err := saveConfig(filepath.Join(root, "mudarro.yaml"), c); err != nil {
		t.Fatal(err)
	}
	return root
}
func doctorRun(root string) (string, error) {
	var out bytes.Buffer
	err := mudarro.Main([]string{"doctor", "--root", root}, "test", strings.NewReader(""), &out)
	return out.String(), err
}

func TestDoctorMissingAndNonExecutableTools(t *testing.T) {
	for _, tc := range []struct {
		name, args            string
		directory, executable bool
	}{{"missing", "unavailable-tool", false, false}, {"non-executable-path", "./tool", false, false}, {"directory", "./tool", true, true}, {"executable", "./tool", false, true}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PATH", t.TempDir())
			root := doctorFixture(t, mudarro.Infrastructure{Kind: "custom"}, map[string]mudarro.Command{"probe": {Args: []string{tc.args}, Group: "scripts"}})
			if tc.args == "./tool" {
				if tc.directory {
					if err := os.Mkdir(filepath.Join(root, "tool"), 0755); err != nil {
						t.Fatal(err)
					}
				} else {
					mode := os.FileMode(0644)
					if tc.executable {
						mode = 0755
					}
					if err := os.WriteFile(filepath.Join(root, "tool"), []byte("#!/bin/sh\ntouch SHOULD_NOT_EXECUTE\n"), mode); err != nil {
						t.Fatal(err)
					}
				}
			}
			output, err := doctorRun(root)
			success := tc.executable && !tc.directory
			if (err == nil) != success {
				t.Fatalf("err=%v output=%q", err, output)
			}
			if !success && !strings.Contains(output, tc.args) {
				t.Fatalf("missing tool identity: %q", output)
			}
			if _, e := os.Stat(filepath.Join(root, "SHOULD_NOT_EXECUTE")); !os.IsNotExist(e) {
				t.Fatal("doctor executed application")
			}
		})
	}
}

func TestDoctorSymlinkAndBlockedConfiguration(t *testing.T) {
	root := doctorFixture(t, mudarro.Infrastructure{Kind: "custom"}, map[string]mudarro.Command{"probe": {Args: []string{"./linked"}}})
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	output, err := doctorRun(root)
	if err == nil || !strings.Contains(output, "linked") {
		t.Fatalf("symlink accepted: %v %q", err, output)
	}
	root = doctorFixture(t, mudarro.Infrastructure{Kind: "kubernetes", Mode: "manifests", File: "k8s"}, nil)
	output, err = doctorRun(root)
	if err == nil || !strings.Contains(output, "app:up") {
		t.Fatalf("blocked configuration accepted: %v %q", err, output)
	}
}

func TestDoctorComposeProviderFailureAndScope(t *testing.T) {
	for _, engine := range []string{"docker", "podman"} {
		for _, status := range []string{"0", "9"} {
			t.Run(engine+"/"+status, func(t *testing.T) {
				bin := t.TempDir()
				t.Setenv("PATH", bin)
				root := doctorFixture(t, mudarro.Infrastructure{Kind: engine, Mode: "compose", File: "compose.yaml"}, nil)
				executable := "docker"
				expected := "compose version"
				if engine == "podman" {
					executable = "podman-compose"
					expected = "--version"
				}
				script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '" + filepath.Join(root, "provider.calls") + "'\nexit " + status + "\n"
				if err := os.WriteFile(filepath.Join(bin, executable), []byte(script), 0755); err != nil {
					t.Fatal(err)
				}
				output, err := doctorRun(root)
				if (err == nil) != (status == "0") {
					t.Fatalf("provider result=%v output=%q", err, output)
				}
				calls, e := os.ReadFile(filepath.Join(root, "provider.calls"))
				if e != nil {
					t.Fatal(e)
				}
				if string(calls) != expected+"\n" {
					t.Fatalf("doctor invoked operational command: %q", calls)
				}
				if status == "0" && !strings.Contains(output, "não valida conexão") {
					t.Fatalf("missing health-check limitation: %q", output)
				}
			})
		}
	}
}

func TestDoctorRequiresAndConfigSymlink(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	root := doctorFixture(t, mudarro.Infrastructure{Kind: "custom"}, map[string]mudarro.Command{"probe": {Args: []string{"./tool"}, Requires: []string{"missing-required-module"}}})
	if err := os.WriteFile(filepath.Join(root, "tool"), []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	output, err := doctorRun(root)
	if err == nil || !strings.Contains(output, "missing-required-module") {
		t.Fatalf("requires failure lost: %v %q", err, output)
	}
	outside := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.Rename(filepath.Join(root, "mudarro.yaml"), outside); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "mudarro.yaml")); err != nil {
		t.Fatal(err)
	}
	output, err = doctorRun(root)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("configuration symlink accepted: %v %q", err, output)
	}
}
