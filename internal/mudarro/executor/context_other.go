//go:build !linux && !darwin && !windows

package executor

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type ContextOS struct{ Context context.Context }

func (e ContextOS) command(directory string, args []string) (*exec.Cmd, error) {
	if len(args) == 0 || args[0] == "" {
		return nil, fmt.Errorf("comando vazio")
	}
	ctx := e.Context
	if ctx == nil {
		ctx = context.Background()
	}
	c := exec.CommandContext(ctx, args[0], args[1:]...)
	c.Dir = directory
	c.WaitDelay = 750 * time.Millisecond
	return c, nil
}
func (e ContextOS) Run(directory string, args []string, in io.Reader, out io.Writer) error {
	c, err := e.command(directory, args)
	if err != nil {
		return err
	}
	c.Stdin = in
	c.Stdout = out
	c.Stderr = out
	return c.Run()
}
func (e ContextOS) Output(directory string, args []string) ([]byte, error) {
	c, err := e.command(directory, args)
	if err != nil {
		return nil, err
	}
	return c.Output()
}
