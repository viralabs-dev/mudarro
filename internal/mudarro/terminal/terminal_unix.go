//go:build linux || darwin

package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

// Query only the destination stream. Character devices such as /dev/null are not TTYs.
func terminalSize(f *os.File) (int, int, bool) {
	var size struct{ rows, cols, xpixel, ypixel uint16 }
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)))
	if errno != 0 {
		return 80, 24, false
	}
	width := int(size.cols)
	if width < 20 || width > 240 {
		width = 80
	}
	height := int(size.rows)
	if height < 12 || height > 200 {
		height = 24
	}
	return width, height, true
}

func terminalColumns(f *os.File) (int, bool) { width, _, tty := terminalSize(f); return width, tty }

func terminalName() string { return os.Getenv("TERM") }
