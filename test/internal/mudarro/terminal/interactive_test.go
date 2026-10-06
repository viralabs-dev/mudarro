//go:build linux || darwin

package terminal_test

import (
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestInteractiveNavigationDoesNotExecute(t *testing.T) {
	s := InteractiveState{}
	for _, key := range []string{"down", "tab", "2", "page-down", "end", "wheel-up", "up"} {
		if _, done := s.Navigate(key, 4, 5, 40); done {
			t.Fatalf("navigation %s executed", key)
		}
	}
	s.Navigate("3", 4, 5, 40)
	n, done := s.Navigate("enter", 4, 5, 40)
	if !done || n != 3 {
		t.Fatalf("explicit enter: %d %v", n, done)
	}
	s.Navigate("end", 4, 5, 40)
	if s.Offset != 35 {
		t.Fatal(s.Offset)
	}
	s.Navigate("home", 4, 5, 40)
	if s.Offset != 0 {
		t.Fatal(s.Offset)
	}
}
func TestInteractiveRequiresBothTTY(t *testing.T) {
	f, e := os.Open(os.DevNull)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	if _, supported, e := NewMenu(f, "Mudarro", "").ChooseInteractive(f, "Actions", []string{"test"}, nil); supported || e != nil {
		t.Fatalf("non tty: %v %v", supported, e)
	}
}
func TestInteractiveRealPTYNavigationAndRestore(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "1")
	master, slave, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer master.Close()
	defer slave.Close()
	before, e := interactiveTermios(slave)
	if e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() {
		n, supported, e := NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", []string{"one", "two"}, []string{strings.Repeat("read only plan\n", 40), "second"})
		if e == nil && (!supported || n != 2) {
			e = fmt.Errorf("selection %d supported %v", n, supported)
		}
		done <- e
	}()
	// Drain our disposable PTY to avoid blocking a full terminal output buffer.
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, e := master.Read(buf); e != nil {
				return
			}
		}
	}()
	time.Sleep(50 * time.Millisecond)
	master.Write([]byte("\x1b[B\x1b[6~\x1b[<65;20;6M\x1b[<0;80;8M\x1b[<32;80;10M\x1b[<0;80;10m\t\r"))
	select {
	case e := <-done:
		after, te := interactiveTermios(slave)
		if te != nil || after != before {
			t.Fatalf("termios not restored: %v", te)
		}
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("PTY interactive timeout")
	}
}

func TestInteractiveKeepsChildInputQueued(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	master, slave, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer master.Close()
	defer slave.Close()
	done := make(chan error, 1)
	labels := make([]string, 10)
	for i := range labels {
		labels[i] = fmt.Sprint("action ", i+1)
	}
	go func() {
		n, _, e := NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", labels, nil)
		if e == nil && n != 10 {
			e = fmt.Errorf("focus10 selected %d", n)
		}
		done <- e
	}()
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, e := master.Read(buf); e != nil {
				return
			}
		}
	}()
	time.Sleep(50 * time.Millisecond)
	master.Write([]byte("10\rchild-input\n"))
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}
	fd := int(slave.Fd())
	syscall.SetNonblock(fd, true)
	defer syscall.SetNonblock(fd, false)
	buf := make([]byte, 128)
	n, e := syscall.Read(fd, buf)
	if e != nil || string(buf[:n]) != "child-input\n" {
		t.Fatalf("child stdin lost: %q %v", buf[:n], e)
	}
}
func TestInteractiveCtrlCRestoresTerminal(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	master, slave, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer master.Close()
	defer slave.Close()
	before, e := interactiveTermios(slave)
	if e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() {
		_, _, e := NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", []string{"one"}, nil)
		done <- e
	}()
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, e := master.Read(buf); e != nil {
				return
			}
		}
	}()
	time.Sleep(50 * time.Millisecond)
	master.Write([]byte{3})
	select {
	case e := <-done:
		if e == nil {
			t.Fatal("CtrlC no cancellation")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}
	after, e := interactiveTermios(slave)
	if e != nil || after != before {
		t.Fatal("CtrlC terminal not restored")
	}
}

func TestInteractiveInvalidNumberDoesNotExecutePreviousFocus(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	master, slave, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer master.Close()
	defer slave.Close()
	done := make(chan error, 1)
	go func() {
		n, _, e := NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", []string{"one"}, nil)
		if e == nil && n != 0 {
			e = fmt.Errorf("invalid number executed priorfocus: %d", n)
		}
		done <- e
	}()
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, e := master.Read(buf); e != nil {
				return
			}
		}
	}()
	time.Sleep(50 * time.Millisecond)
	master.Write([]byte("99\rq"))
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("invalid number timeout")
	}
}
