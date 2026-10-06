//go:build linux || darwin

package terminal_test

import (
	"bytes"
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"io"
	"strings"
	"testing"
	"time"
)

func TestInteractivePreservesPartialEscapeAcrossIdle(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	master, slave, err := openPTY()
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	drained := make(chan struct{})
	go func() { defer close(drained); _, _ = io.Copy(&output, master) }()
	done := make(chan error, 1)
	go func() {
		n, supported, err := NewMenu(slave, "Aurora", "").ChooseInteractive(slave, "Actions", []string{"test"}, []string{"PARTIAL_ESCAPE_PREVIEW"})
		if err == nil && (!supported || n != 0) {
			err = fmt.Errorf("unexpected selection %d", n)
		}
		done <- err
	}()
	// Preserve ESC[ across several poll deadlines, then complete Down.
	time.Sleep(150 * time.Millisecond)
	_, _ = master.Write([]byte("\x1b["))
	time.Sleep(350 * time.Millisecond)
	_, _ = master.Write([]byte("B"))
	time.Sleep(150 * time.Millisecond)
	_, _ = master.Write([]byte("q"))
	select {
	case err = <-done:
	case <-time.After(2 * time.Second):
		_, _ = master.Write([]byte{3})
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("menu did not cancel")
		}
		err = fmt.Errorf("menu exit timed out")
	}
	_ = slave.Close()
	_ = master.Close()
	<-drained
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "PARTIAL_ESCAPE_PREVIEW") {
		t.Fatal("partial Down sequence was lost across idle polling")
	}
}
