//go:build !linux && !darwin

package terminal_test

import (
	"os"
	"testing"
)

func capturePTY(t *testing.T, width, rows uint16) (*os.File, func() string) {
	t.Helper()
	t.Skip("real PTY presentation tests require Linux or macOS, matching production terminal support")
	return nil, nil
}
