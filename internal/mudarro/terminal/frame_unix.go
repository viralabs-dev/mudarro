//go:build linux || darwin

package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

// Keep real small-window dimensions rather than pretending a tiny terminal has 24 rows.
func frameTerminalSize(f *os.File) (int, int, bool) {
	var s struct{ rows, cols, x, y uint16 }
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&s)))
	if e != 0 {
		return 80, 24, false
	}
	width, height := int(s.cols), int(s.rows)
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 24
	}
	return width, height, true
}
