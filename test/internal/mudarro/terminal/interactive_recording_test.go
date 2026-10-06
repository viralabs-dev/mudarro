//go:build linux || darwin

package terminal_test

import (
	"encoding/json"
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// Optional real PTY capture. It records terminal output, never synthesizes frames.
func TestInteractiveRecording(t *testing.T) {
	path := os.Getenv("MUDARRO_INTERACTIVE_CAST")
	if path == "" {
		t.Skip("optional PTY recording")
	}
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "1")
	master, slave, e := openPTY()
	if e != nil {
		t.Fatal(e)
	}
	defer slave.Close()
	size := struct{ rows, cols, xpixel, ypixel uint16 }{24, 80, 0, 0}
	if e = ioctl(slave, uintptr(syscall.TIOCSWINSZ), unsafe.Pointer(&size)); e != nil {
		t.Fatal(e)
	}
	file, e := os.Create(path)
	if e != nil {
		t.Fatal(e)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(map[string]any{"version": 2, "width": 80, "height": 24, "title": "Mudarro genuine isolated PTY preview; injected keyboard and SGR mouse", "env": map[string]string{"TERM": "xterm-256color"}})
	start := time.Now()
	var mu sync.Mutex
	captureDone := make(chan struct{})
	go func() {
		defer close(captureDone)
		buf := make([]byte, 4096)
		for {
			n, e := master.Read(buf)
			if n > 0 {
				mu.Lock()
				encoder.Encode([]any{time.Since(start).Seconds(), "o", string(buf[:n])})
				mu.Unlock()
			}
			if e != nil {
				return
			}
		}
	}()
	done := make(chan error, 1)
	lines := []string{}
	for i := 0; i < 60; i++ {
		lines = append(lines, fmt.Sprintf("%02d READ ONLY PLAN: generated wrapper -> planned command; no action executed", i+1))
	}
	go func() {
		n, _, e := NewMenu(slave, "Mudarro", "").ChooseInteractive(slave, "Actions", []string{"test (preview only)", "up (preview only)"}, []string{strings.Join(lines, "\n"), "Second action preview"})
		if e == nil && n != 0 {
			e = fmt.Errorf("demo must exit without action")
		}
		done <- e
	}()
	for _, keys := range []string{"\x1b[B", "\x1b[6~", "\x1b[<65;20;8M", "\x1b[<0;80;8M", "\x1b[<32;80;19M", "\x1b[<0;80;19m", "\x1b[H", "\x1b[F", "\x1b[A", "\t", "\x1b[B", "q"} {
		time.Sleep(150 * time.Millisecond)
		if _, e = master.Write([]byte(keys)); e != nil {
			t.Fatal(e)
		}
	}
	select {
	case e = <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("recording timeout")
	}
	time.Sleep(100 * time.Millisecond)
	slave.Close()
	master.Close()
	<-captureDone
}
