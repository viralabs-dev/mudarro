//go:build linux || darwin

package mudarro_test

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
)

type frameCapture struct {
	mu            sync.Mutex
	data          bytes.Buffer
	master, slave *os.File
	stop, done    chan struct{}
}

func framePTY(t *testing.T, cols, rows uint16) *frameCapture {
	t.Helper()
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Fatal(err)
	}
	call := func(request uintptr, arg unsafe.Pointer) error {
		_, _, e := syscall.Syscall(syscall.SYS_IOCTL, master.Fd(), request, uintptr(arg))
		if e != 0 {
			return e
		}
		return nil
	}
	var path string
	if runtime.GOOS == "linux" {
		var unlock, number uint32
		if err = call(0x40045431, unsafe.Pointer(&unlock)); err == nil {
			err = call(0x80045430, unsafe.Pointer(&number))
		}
		path = fmt.Sprintf("/dev/pts/%d", number)
	} else {
		var name [128]byte
		for _, r := range []uintptr{0x20007454, 0x20007452} {
			if err = call(r, nil); err != nil {
				break
			}
		}
		if err == nil {
			err = call(0x40807453, unsafe.Pointer(&name[0]))
		}
		end := bytes.IndexByte(name[:], 0)
		if end < 0 {
			end = len(name)
		}
		path = string(name[:end])
	}
	if err != nil {
		master.Close()
		t.Fatal(err)
	}
	slave, err := os.OpenFile(path, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		t.Fatal(err)
	}
	capture := &frameCapture{master: master, slave: slave, stop: make(chan struct{}), done: make(chan struct{})}
	capture.resize(t, cols, rows)
	fd := int(master.Fd())
	if err = syscall.SetNonblock(fd, true); err != nil {
		master.Close()
		slave.Close()
		t.Fatal(err)
	}
	go func() {
		defer close(capture.done)
		buf := make([]byte, 32768)
		for {
			select {
			case <-capture.stop:
				return
			default:
			}
			n, err := syscall.Read(fd, buf)
			if n > 0 {
				capture.mu.Lock()
				capture.data.Write(buf[:n])
				capture.mu.Unlock()
			}
			if err != nil && err != syscall.EAGAIN && err != syscall.EINTR {
				return
			}
			if n <= 0 {
				time.Sleep(time.Millisecond)
			}
		}
	}()
	t.Cleanup(func() { close(capture.stop); <-capture.done; slave.Close(); master.Close() })
	return capture
}
func (c *frameCapture) resize(t *testing.T, cols, rows uint16) {
	t.Helper()
	size := struct{ rows, cols, x, y uint16 }{rows: rows, cols: cols}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, c.slave.Fd(), uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&size)))
	if err != 0 {
		t.Fatal(err)
	}
}
func (c *frameCapture) text() string {
	time.Sleep(20 * time.Millisecond)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data.String()
}
func (c *frameCapture) reset() { c.text(); c.mu.Lock(); c.data.Reset(); c.mu.Unlock() }

func TestFramePersistentChromeAndGeometry(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	c := framePTY(t, 80, 24)
	v := terminal.NewMenu(c.slave, "FRAME_HEADER_PROJECT", "initial context")
	active, err := v.BeginShell(c.slave)
	if err != nil || !active {
		t.Fatalf("begin %v %v", active, err)
	}
	defer v.CloseShell()
	width, top, height := v.ShellGeometry()
	if width != 80 || top < 2 || height < 1 || top+height > 24 {
		t.Fatalf("geometry %d %d %d", width, top, height)
	}
	initial := c.text()
	if strings.Count(initial, "\x1b[?1049h") != 1 || !strings.Contains(initial, "FRAME_HEADER_PROJECT") {
		t.Fatal("missing persistent header")
	}
	c.reset()
	v.SetContext("changed central context")
	v.ShellRender("title", []string{"central first", "central second"}, "changing hint")
	changed := c.text()
	if strings.Contains(changed, "FRAME_HEADER_PROJECT") || strings.Contains(changed, "changing hint") || strings.Contains(changed, "\x1b[2J") || strings.Contains(changed, "\x1b[?1049h") {
		t.Fatal("central render changed chrome")
	}
	if !strings.Contains(changed, fmt.Sprintf("\x1b[%d;1Hcentral first", top)) {
		t.Fatal("body origin mismatch")
	}
	v.CloseShell()
	v.CloseShell()
	closed := c.text()
	if strings.Count(closed, "\x1b[?1049l") != 1 || !strings.Contains(closed, "\x1b[r") || !strings.Contains(closed, "\x1b[?7h") {
		t.Fatal("frame cleanup missing/duplicated")
	}
}
func TestFrameOutputFiltersFragmentedControlsAndBoundsHistory(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	c := framePTY(t, 80, 24)
	v := terminal.NewMenu(c.slave, "frame", "")
	active, err := v.BeginShell(c.slave)
	if err != nil || !active {
		t.Fatal(active, err)
	}
	defer v.CloseShell()
	w := v.ShellOutput()
	c.reset()
	for _, chunk := range []string{"safe text\r\n", "\x1b[", "2J", "\x1b]52;c;DO_NOT_SHOW_OSC", "\x1b", "\\", "\x1bPDO_NOT_SHOW_DCS", "\x1b\\", "unicode ", "\xe2", "\x82", "\xac", "\n"} {
		n, err := w.Write([]byte(chunk))
		if err != nil || n != len(chunk) {
			t.Fatal(n, err)
		}
	}
	output := c.text()
	if strings.Contains(output, "DO_NOT_SHOW_OSC") || strings.Contains(output, "DO_NOT_SHOW_DCS") || strings.Contains(output, "\x1b[2J") {
		t.Fatal("child control string escaped frame")
	}
	if !strings.Contains(output, "safe text") || !strings.Contains(output, "unicode €") {
		t.Fatal("valid output was lost")
	}
	c.reset()
	payload := strings.Repeat("x", 100000) + "\n"
	for i := 0; i < 300; i++ {
		payload += fmt.Sprintf("line-%03d\n", i)
	}
	w.Write([]byte(payload))
	final := c.text()
	if !strings.Contains(final, "line-299") || strings.Contains(final, "line-000") || len(final) > 8192 {
		t.Fatalf("history not bounded: %d bytes", len(final))
	}
	old := w
	fresh := v.ShellOutput()
	fresh.Write([]byte("fresh output"))
	c.reset()
	old.Write([]byte("obsolete output"))
	if strings.Contains(c.text(), "obsolete output") {
		t.Fatal("old writer replaced current action")
	}
}
func TestFrameResizeUsesActualSmallWindow(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	c := framePTY(t, 80, 24)
	v := terminal.NewMenu(c.slave, "frame", "")
	active, err := v.BeginShell(c.slave)
	if err != nil || !active {
		t.Fatal(active, err)
	}
	defer v.CloseShell()
	c.resize(t, 12, 6)
	width, top, height := v.ShellGeometry()
	if width != 12 || top+height > 6 || height < 1 {
		t.Fatal(width, top, height)
	}
	c.resize(t, 8, 2)
	width, top, height = v.ShellGeometry()
	if width != 8 || height != 0 || top+height > 2 {
		t.Fatal(width, top, height)
	}
	c.resize(t, 80, 24)
	width, top, height = v.ShellGeometry()
	if width != 80 || height < 1 || top+height > 24 {
		t.Fatal(width, top, height)
	}
}
func TestFramePipeFallback(t *testing.T) {
	var out bytes.Buffer
	v := terminal.NewMenu(&out, "frame", "")
	active, err := v.BeginShell(os.Stdin)
	if active || err != nil || v.ShellActive() {
		t.Fatal(active, err)
	}
	v.ShellRender("title", []string{"plain line"}, "")
	if strings.Contains(out.String(), "\x1b") || !strings.Contains(out.String(), "plain line") {
		t.Fatal(out.String())
	}
}
