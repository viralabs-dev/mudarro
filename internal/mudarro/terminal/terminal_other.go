//go:build !linux && !darwin

package terminal

import "os"

func terminalColumns(_ *os.File) (int, bool) { return 80, false }

func terminalSize(_ *os.File) (int, int, bool) { return 80, 24, false }
