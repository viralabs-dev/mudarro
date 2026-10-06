//go:build linux

package terminal_test

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func openPTY() (*os.File, *os.File, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	var unlock, number uint32
	if err = ioctl(master, uintptr(syscall.TIOCSPTLCK), unsafe.Pointer(&unlock)); err == nil {
		err = ioctl(master, uintptr(syscall.TIOCGPTN), unsafe.Pointer(&number))
	}
	if err != nil {
		master.Close()
		return nil, nil, err
	}
	slave, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		return nil, nil, err
	}
	return master, slave, nil
}
