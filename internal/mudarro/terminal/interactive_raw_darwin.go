//go:build darwin

package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

func interactiveRaw(f *os.File) (func() error, error) {
	var old syscall.Termios
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCGETA, uintptr(unsafe.Pointer(&old))); e != 0 {
		return nil, e
	}
	raw := old
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	raw.Cc[syscall.VMIN] = 0
	raw.Cc[syscall.VTIME] = 1
	apply := func(t *syscall.Termios) error {
		_, _, e := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCSETA, uintptr(unsafe.Pointer(t)))
		if e != 0 {
			return e
		}
		return nil
	}
	if e := apply(&raw); e != nil {
		return nil, e
	}
	return func() error { return apply(&old) }, nil
}
