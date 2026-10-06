//go:build linux || darwin

package terminal_test

import (
	"context"
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"io"
	"os"
	"testing"
	"time"
)

func TestShellInputCancelRestoreAndRepeat(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	m, s, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer m.Close()
	defer s.Close()
	done := make(chan struct{})
	go func() { io.Copy(io.Discard, m); close(done) }()
	before, e := interactiveTermios(s)
	if e != nil {
		t.Fatal(e)
	}
	v := NewMenu(s, "Aurora", "")
	ok, e := v.BeginShell(s)
	if !ok || e != nil {
		t.Fatalf("shell %v %v", ok, e)
	}
	defer v.CloseShell()
	for n := 0; n < 2; n++ {
		ctx, cancel := context.WithCancel(context.Background())
		input, stop, e := v.ExecutionInput(ctx, cancel)
		if e != nil {
			t.Fatal(e)
		}
		result := make(chan error, 1)
		go func() {
			line, e := v.ReadShellLine(input, io.Discard)
			if e == nil && line != "APAGAR aurora" {
				e = fmt.Errorf("line %q", line)
			}
			result <- e
		}()
		m.Write([]byte("APAGAR aurora\r"))
		select {
		case e := <-result:
			if e != nil {
				t.Fatal(e)
			}
		case <-time.After(2 * time.Second):
			cancel()
			stop()
			t.Fatal("input stalled")
		}
		m.Write([]byte{3})
		select {
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
			cancel()
			stop()
			t.Fatal("Ctrl+C stalled")
		}
		stop()
		stop()
		cancel()
		after, e := interactiveTermios(s)
		if e != nil || before != after {
			t.Fatalf("termios not restored %v", e)
		}
	}
	v.CloseShell()
	s.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		m.Close()
		<-done
	}
}

var _ *os.File
