//go:build linux || darwin

package mudarro_test

import (
	"github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestDoctorSpecialFileExecutableIsRejected(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fifo    bool
		mode    os.FileMode
		success bool
	}{
		{"fifo-executable", true, 0755, false}, {"regular-executable", false, 0755, true}, {"regular-not-executable", false, 0644, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PATH", t.TempDir())
			root := doctorFixture(t, mudarro.Infrastructure{Kind: "custom"}, map[string]mudarro.Command{"probe": {Args: []string{"./probe-tool"}, Group: "scripts"}})
			target := filepath.Join(root, "probe-tool")
			body := []byte("#!/bin/sh\nprintf executed > ACTION_MUST_NOT_RUN\n")
			if tc.fifo {
				if e := syscall.Mkfifo(target, uint32(tc.mode)); e != nil {
					t.Fatal(e)
				}
			} else if e := os.WriteFile(target, body, tc.mode); e != nil {
				t.Fatal(e)
			}
			if e := os.Chmod(target, tc.mode); e != nil {
				t.Fatal(e)
			}
			before, e := os.Lstat(target)
			if e != nil {
				t.Fatal(e)
			}
			config, e := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
			if e != nil {
				t.Fatal(e)
			}
			output, e := doctorRun(root)
			if (e == nil) != tc.success {
				t.Errorf("doctor success=%v want=%v output=%q err=%v", e == nil, tc.success, output, e)
			}
			expected := "AUSENTE app: ./probe-tool"
			if tc.success {
				expected = "OK app: ./probe-tool"
			}
			if !strings.Contains(output, expected) {
				t.Errorf("missing diagnostic %q in %q", expected, output)
			}
			if _, e := os.Stat(filepath.Join(root, "ACTION_MUST_NOT_RUN")); !os.IsNotExist(e) {
				t.Fatal("doctor executed fixture action", e)
			}
			after, e := os.Lstat(target)
			if e != nil || !os.SameFile(before, after) || before.Mode() != after.Mode() || before.Size() != after.Size() || before.ModTime() != after.ModTime() {
				t.Fatal("doctor changed tool fixture", e)
			}
			if !tc.fifo {
				after, e := os.ReadFile(target)
				if e != nil || string(after) != string(body) {
					t.Fatal("doctor changed regular source", e)
				}
			}
			afterConfig, e := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
			if e != nil || string(afterConfig) != string(config) {
				t.Fatal("doctor changed configuration", e)
			}
		})
	}
}
