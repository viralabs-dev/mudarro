//go:build !linux && !darwin

package terminal

import (
	"fmt"
	"os"
)

func interactiveRaw(f *os.File) (func() error, error) {
	return nil, fmt.Errorf("interactive terminal unavailable")
}
