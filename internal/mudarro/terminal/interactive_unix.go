//go:build linux || darwin

package terminal

import (
	"io"
	"os"
	"syscall"
)

// interactiveNext returns the next input bytes, or none when nothing is ready
// within the bounded poll. A readable EOF is io.EOF.
func interactiveNext(in *os.File) ([]byte, error) {
	fd := int(in.Fd())
	ready, e := interactiveReady(fd)
	if e == syscall.EINTR {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	if !ready {
		return nil, nil
	}
	buf := make([]byte, 1)
	n, e := syscall.Read(fd, buf)
	if e != nil {
		if e == syscall.EINTR || e == syscall.EAGAIN {
			return nil, nil
		}
		return nil, e
	}
	if n == 0 {
		return nil, io.EOF
	}
	return buf[:n], nil
}
