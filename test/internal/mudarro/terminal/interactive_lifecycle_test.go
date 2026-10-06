//go:build linux || darwin

package terminal_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestInteractiveSignalCleanup(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run(sig.String(), func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestInteractiveSignalHelper$")
			cmd.Env = append(os.Environ(), "MUDARRO_SIGNAL_HELPER=1", "TERM=xterm-256color")
			stdout, e := cmd.StdoutPipe()
			if e != nil {
				t.Fatal(e)
			}
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			if e = cmd.Start(); e != nil {
				t.Fatal(e)
			}
			ready := make(chan bool, 1)
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "RAW_READY") {
						ready <- true
					}
				}
			}()
			select {
			case <-ready:
			case <-time.After(3 * time.Second):
				cmd.Process.Kill()
				t.Fatal("helper raw readiness timeout")
			}
			if e = cmd.Process.Signal(sig); e != nil {
				t.Fatal(e)
			}
			if e = cmd.Wait(); e != nil {
				t.Fatalf("signal cleanup: %v %s", e, stderr.String())
			}
		})
	}
}
func TestInteractiveSignalHelper(t *testing.T) {
	if os.Getenv("MUDARRO_SIGNAL_HELPER") != "1" {
		return
	}
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
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, e := master.Read(buf); e != nil {
				return
			}
		}
	}()
	probeStop := make(chan struct{})
	probeDone := make(chan struct{})
	// Join the termios observer before deferred slave.Close on every return path.
	defer func() { close(probeStop); <-probeDone }()
	go func() {
		defer close(probeDone)
		for i := 0; i < 300; i++ {
			select {
			case <-probeStop:
				return
			default:
			}
			now, err := interactiveTermios(slave)
			if err == nil && now.Lflag&syscall.ICANON == 0 {
				fmt.Println("RAW_READY")
				return
			}
			select {
			case <-probeStop:
				return
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
	_, _, e = NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", []string{"one"}, nil)
	if e == nil {
		t.Fatal("signal no error")
	}
	after, e := interactiveTermios(slave)
	if e != nil || after != before {
		t.Fatal("signal termios restore failed")
	}
}
func TestInteractiveMasterEOFReturns(t *testing.T) {
	if os.Getenv("MUDARRO_EOF_HELPER") != "1" {
		// A broken EOF loop is killed and joined as a separate process. Never close a
		// slave concurrently with a still-running ChooseInteractive goroutine.
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestInteractiveMasterEOFReturns$")
		cmd.Env = append(os.Environ(), "MUDARRO_EOF_HELPER=1", "TERM=xterm-256color")
		output, e := cmd.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("EOF loop: %v; helper killed and joined: %s", ctx.Err(), output)
		}
		if e != nil {
			t.Fatalf("EOF helper: %v %s", e, output)
		}
		return
	}
	t.Setenv("TERM", "xterm-256color")
	master, slave, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer slave.Close()
	closed := make(chan struct{})
	go func() { time.Sleep(50 * time.Millisecond); master.Close(); close(closed) }()
	_, _, e = NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", []string{"one"}, nil)
	<-closed
	if e == nil {
		t.Fatal("EOF no error")
	}
}
