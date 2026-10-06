//go:build linux || darwin

package terminal_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func ioctl(f *os.File, request uintptr, arg unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), request, uintptr(arg))
	if errno != 0 {
		return errno
	}
	return nil
}

// Capture the destination PTY only; no user terminal or desktop is accessed.
func capturePTY(t *testing.T, width, rows uint16) (*os.File, func() string) {
	t.Helper()
	master, slave, err := openPTY()
	if err != nil {
		t.Fatalf("open isolated PTY: %v", err)
	}
	size := struct{ rows, cols, xpixel, ypixel uint16 }{rows: rows, cols: width}
	if err := ioctl(slave, uintptr(syscall.TIOCSWINSZ), unsafe.Pointer(&size)); err != nil {
		master.Close()
		slave.Close()
		t.Fatal(err)
	}
	type result struct {
		text string
		err  error
	}
	done := make(chan result, 1)
	go func() {
		var b bytes.Buffer
		_, err := io.Copy(&b, master)
		done <- result{b.String(), err}
	}()
	stopped := false
	var output string
	stop := func() string {
		t.Helper()
		if stopped {
			return output
		}
		stopped = true
		slave.Close()
		select {
		case r := <-done:
			output = r.text
			if r.err != nil && !errors.Is(r.err, syscall.EIO) {
				t.Errorf("PTY capture: %v", r.err)
			}
		case <-time.After(3 * time.Second):
			master.Close()
			t.Error("PTY capture did not finish after closing slave")
		}
		master.Close()
		return output
	}
	t.Cleanup(func() { stop() })
	return slave, stop
}
