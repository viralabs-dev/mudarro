//go:build darwin

package terminal_test

import (
	"os"
	"syscall"
	"unsafe"
)

func interactiveTermios(f *os.File) (syscall.Termios, error) {
	var t syscall.Termios
	e := ioctl(f, uintptr(syscall.TIOCGETA), unsafe.Pointer(&t))
	return t, e
}
