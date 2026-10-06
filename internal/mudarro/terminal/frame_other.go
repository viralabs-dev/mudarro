//go:build !linux && !darwin

package terminal

import "os"

func frameTerminalSize(*os.File) (int, int, bool) { return 80, 24, false }
