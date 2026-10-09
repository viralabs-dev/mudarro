//go:build linux || darwin

package mudarro_test

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	. "github.com/viralabs-dev/mudarro/internal/mudarro"
)

// FIFOs are a Unix special file; Windows named pipes do not live in the project tree.
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
