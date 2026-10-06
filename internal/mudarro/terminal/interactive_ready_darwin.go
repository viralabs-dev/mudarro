//go:build darwin

package terminal

import (
	"fmt"
	"syscall"
	"unsafe"
)

// interactiveReady distinguishes bounded input polling from a readable EOF.
// A zero read with VMIN=0/VTIME>0 alone cannot distinguish them on BSD PTYs.
func interactiveReady(fd int) (bool, error) {
	var set syscall.FdSet
	bits := int(unsafe.Sizeof(set.Bits[0])) * 8
	if fd < 0 || fd >= len(set.Bits)*bits {
		return false, fmt.Errorf("terminal descriptor outside select range: %d", fd)
	}
	set.Bits[fd/bits] |= 1 << uint(fd%bits)
	timeout := syscall.Timeval{Usec: 100000}
	err := syscall.Select(fd+1, &set, nil, nil, &timeout)
	if err != nil {
		return false, err
	}
	return set.Bits[fd/bits]&(1<<uint(fd%bits)) != 0, nil
}
