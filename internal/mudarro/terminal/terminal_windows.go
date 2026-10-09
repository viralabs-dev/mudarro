//go:build windows

package terminal

import (
	"os"

	"github.com/viralabs-dev/mudarro/internal/mudarro/winapi"
)

// consoleSize reports the visible console window. An output console counts as
// a terminal only when it accepts VT sequences (enabled here, Windows 10+);
// an input console has no window size and reports the defaults.
func consoleSize(f *os.File) (int, int, bool) {
	if width, height, ok := winapi.WindowSize(f); ok {
		if !winapi.EnableVirtualTerminalOutput(f) {
			return 80, 24, false
		}
		return width, height, true
	}
	if _, ok := winapi.ConsoleMode(f); ok {
		return 80, 24, true
	}
	return 80, 24, false
}

func terminalSize(f *os.File) (int, int, bool) {
	width, height, tty := consoleSize(f)
	if width < 20 || width > 240 {
		width = 80
	}
	if height < 12 || height > 200 {
		height = 24
	}
	return width, height, tty
}

func terminalColumns(f *os.File) (int, bool) { width, _, tty := terminalSize(f); return width, tty }

// Keep real small-window dimensions rather than pretending a tiny console has 24 rows.
func frameTerminalSize(f *os.File) (int, int, bool) {
	width, height, tty := consoleSize(f)
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}
	return width, height, tty
}

// terminalName honours TERM (so TERM=dumb still disables colour and the
// interactive menu). Windows consoles usually leave TERM unset: Windows
// Terminal is xterm-compatible; the classic console gets a conservative name.
func terminalName() string {
	if term := os.Getenv("TERM"); term != "" {
		return term
	}
	if os.Getenv("WT_SESSION") != "" {
		return "xterm-256color"
	}
	return "windows-console"
}
