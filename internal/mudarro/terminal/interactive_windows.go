//go:build windows

package terminal

import (
	"fmt"
	"os"

	"github.com/viralabs-dev/mudarro/internal/mudarro/winapi"
)

// interactiveRaw puts the console input in raw VT mode: no echo, no line
// editing, Ctrl+C delivered as byte 3, and keys reported as VT sequences.
func interactiveRaw(f *os.File) (func() error, error) {
	old, ok := winapi.ConsoleMode(f)
	if !ok {
		return nil, fmt.Errorf("interactive terminal unavailable")
	}
	raw := old&^(winapi.EnableProcessedInput|winapi.EnableLineInput|winapi.EnableEchoInput) | winapi.EnableVirtualTerminal
	if err := winapi.SetConsoleMode(f, raw); err != nil {
		return nil, fmt.Errorf("console without virtual terminal input (Windows 10 or later required): %w", err)
	}
	return func() error { return winapi.SetConsoleMode(f, old) }, nil
}

// interactiveNext waits up to 100 ms for a key, matching the Unix select poll.
func interactiveNext(in *os.File) ([]byte, error) {
	ready, err := winapi.InputReady(in, 100)
	if err != nil || !ready {
		return nil, err
	}
	return winapi.ReadInput(in)
}
