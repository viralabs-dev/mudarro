//go:build !linux && !darwin && !windows

package terminal

import "os"

func terminalColumns(_ *os.File) (int, bool) { return 80, false }

func terminalSize(_ *os.File) (int, int, bool) { return 80, 24, false }

func terminalName() string { return os.Getenv("TERM") }
