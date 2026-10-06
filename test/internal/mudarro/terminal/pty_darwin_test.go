//go:build darwin

package terminal_test

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

func openPTY() (*os.File, *os.File, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	var name [128]byte
	for _, request := range []uintptr{syscall.TIOCPTYGRANT, syscall.TIOCPTYUNLK} {
		if err = ioctl(master, request, nil); err != nil {
			master.Close()
			return nil, nil, err
		}
	}
	if err = ioctl(master, syscall.TIOCPTYGNAME, unsafe.Pointer(&name[0])); err != nil {
		master.Close()
		return nil, nil, err
	}
	end := bytes.IndexByte(name[:], 0)
	if end < 0 {
		end = len(name)
	}
	slave, err := os.OpenFile(string(name[:end]), os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		return nil, nil, err
	}
	return master, slave, nil
}
