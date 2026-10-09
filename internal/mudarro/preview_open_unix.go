//go:build !windows

package mudarro

import (
	"os"
	"syscall"
)

// openPreviewSource never follows a final symlink and never blocks on a FIFO.
func openPreviewSource(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
}
