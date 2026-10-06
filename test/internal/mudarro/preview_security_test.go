package mudarro_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	. "github.com/viralabs-dev/mudarro/internal/mudarro"
)

func previewSecurityFixture(t *testing.T, command Command, files map[string]string) string {
	t.Helper()
	if files == nil {
		files = map[string]string{}
	}
	c := Config{Version: 1, Name: "preview security fixture", Services: []Service{{ID: "app", Dir: ".", Infrastructure: Infrastructure{Kind: "custom"}, Commands: map[string]Command{"probe": command}}}}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	files["mudarro.json"] = string(b)
	return fixture(t, files)
}
func publicSecurityPreview(root string) (string, error) {
	var out bytes.Buffer
	err := Main([]string{"preview", "app:probe", "--root", root}, "test", strings.NewReader(""), &out)
	return out.String(), err
}

func TestPreviewSecurityFIFOIsNonblocking(t *testing.T) {
	root := previewSecurityFixture(t, Command{Args: []string{"bash", "probe.sh"}}, nil)
	fifo := filepath.Join(root, "probe.sh")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, err := publicSecurityPreview(root); done <- err }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		// Unblock a regressed reader without invoking a shell or leaving a blocked goroutine.
		if writer, err := os.OpenFile(fifo, os.O_WRONLY|syscall.O_NONBLOCK, 0600); err == nil {
			writer.Close()
		}
		t.Fatal("preview blocked while opening a FIFO")
	}
}
func TestPreviewSecurityRedactsEmbeddedCredentials(t *testing.T) {
	for _, tc := range []struct {
		name    string
		command Command
		files   map[string]string
		secrets []string
	}{
		{"embedded-uri", Command{Shell: "echo postgresql://user:dummy-uri-password@localhost/db"}, nil, []string{"dummy-uri-password"}},
		{"uri-query", Command{Args: []string{"echo", "https://localhost/resource?token=dummy-query-token&password=dummy-query-password"}}, nil, []string{"dummy-query-token", "dummy-query-password"}},
		{"authorization", Command{Shell: "curl --header 'Authorization: Bearer dummy-bearer-token' --header 'Authorization: Basic dummy-basic-value' https://localhost"}, nil, []string{"dummy-bearer-token", "dummy-basic-value"}},
		{"flags-assignments", Command{Args: []string{"echo", "--token", "dummy-flag-token", "PASSWORD=dummy-assignment-password"}}, nil, []string{"dummy-flag-token", "dummy-assignment-password"}},
		{"script-source", Command{Args: []string{"bash", "probe.sh"}}, map[string]string{"probe.sh": "#!/bin/sh\necho https://user:dummy-source-password@localhost/?token=dummy-source-token\necho 'Authorization: Bearer dummy-source-bearer'\n"}, []string{"dummy-source-password", "dummy-source-token", "dummy-source-bearer"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := publicSecurityPreview(previewSecurityFixture(t, tc.command, tc.files))
			if err != nil {
				t.Fatal(err)
			}
			for _, secret := range tc.secrets {
				if strings.Contains(out, secret) {
					t.Fatalf("preview exposed synthetic credential %q", secret)
				}
			}
			if !strings.Contains(out, "REDACTED") {
				t.Fatal("redaction missing")
			}
		})
	}
}
func TestPreviewSecurityNavigationDoesNotExecute(t *testing.T) {
	root := previewSecurityFixture(t, Command{Args: []string{"touch", "execution-marker"}, Group: "scripts"}, nil)
	var out bytes.Buffer
	err := Main([]string{"menu", "--root", root}, "test", strings.NewReader("1\n1\np1\n0\n0\n0\n"), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "app:probe") {
		t.Fatal("preview not displayed")
	}
	if _, err := os.Stat(filepath.Join(root, "execution-marker")); !os.IsNotExist(err) {
		t.Fatal("preview executed action", err)
	}
}
func TestPreviewSecurityTotalSourceBudget(t *testing.T) {
	args := []string{"bash"}
	for i := 0; i < 100; i++ {
		args = append(args, "large.sh")
	}
	args = append(args, "missing.sh", "../outside.sh")
	root := previewSecurityFixture(t, Command{Args: args}, map[string]string{"large.sh": strings.Repeat("bounded source line\n", 1000)})
	out, err := publicSecurityPreview(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) > 20<<10 {
		t.Fatalf("aggregate preview exceeded budget: %d bytes", len(out))
	}
	if !strings.Contains(out, "truncated") && !strings.Contains(out, "limit") && !strings.Contains(out, "budget") {
		t.Fatal("source truncation missing")
	}
}

func TestPreviewSecurityRejectsSourceSymlinks(t *testing.T) {
	root := previewSecurityFixture(t, Command{Args: []string{"bash", "linked.sh", "../outside.sh"}}, nil)
	outside := fixture(t, map[string]string{"outside.sh": "OUTSIDE_SOURCE_MUST_NOT_APPEAR"})
	if err := os.Symlink(filepath.Join(outside, "outside.sh"), filepath.Join(root, "linked.sh")); err != nil {
		t.Fatal(err)
	}
	out, err := publicSecurityPreview(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "OUTSIDE_SOURCE_MUST_NOT_APPEAR") {
		t.Fatal("preview followed a source symlink")
	}
}

func TestPreviewSecuritySourceCannotEmitTerminalControls(t *testing.T) {
	root := previewSecurityFixture(t, Command{Args: []string{"bash", "probe.sh"}}, map[string]string{"probe.sh": "#!/bin/sh\necho visible-preview-text\n\x1b]52;c;synthetic-clipboard\a\x1b[2J"})
	out, err := publicSecurityPreview(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(out, "\x1b\a") {
		t.Fatal("source emitted live terminal controls")
	}
	if !strings.Contains(out, "visible-preview-text") {
		t.Fatal("source disappeared instead of being sanitized")
	}
}
