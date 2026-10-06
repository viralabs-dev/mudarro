//go:build !linux && !darwin

package terminal

import (
	"context"
	"fmt"
	"io"
	"os"
)

func (v *Menu) ExecutionInput(context.Context, context.CancelFunc) (*os.File, func(), error) {
	return nil, nil, fmt.Errorf("TTY unavailable")
}
func (v *Menu) ReadShellLine(io.Reader, io.Writer) (string, error) {
	return "", fmt.Errorf("TTY unavailable")
}
