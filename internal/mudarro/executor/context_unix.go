//go:build linux || darwin

package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// ContextOS cancels only the process group it creates. It never changes terminal state.
// A blocking custom stdin Reader must support cancellation itself; nil and *os.File
// inputs avoid os/exec's stdin copy goroutine. Readers are never closed by this executor.
type ContextOS struct{ Context context.Context }

func (e ContextOS) command(directory string, args []string) (*exec.Cmd, func(error), error) {
	if len(args) == 0 || args[0] == "" {
		return nil, nil, fmt.Errorf("comando vazio")
	}
	ctx := e.Context
	if ctx == nil {
		ctx = context.Background()
	}
	c := exec.CommandContext(ctx, args[0], args[1:]...)
	c.Dir = directory
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.WaitDelay = 750 * time.Millisecond
	cancelled := make(chan struct{}, 1)
	killed := make(chan struct{})
	var cancelOnce sync.Once
	var cancelErr error
	c.Cancel = func() error {
		cancelOnce.Do(func() {
			cancelErr = syscall.Kill(-c.Process.Pid, syscall.SIGTERM)
			if errors.Is(cancelErr, syscall.ESRCH) {
				cancelErr = os.ErrProcessDone
				return
			}
			if cancelErr != nil {
				return
			}
			cancelled <- struct{}{}
			time.AfterFunc(300*time.Millisecond, func() { syscall.Kill(-c.Process.Pid, syscall.SIGKILL); close(killed) })
		})
		return cancelErr
	}
	finish := func(err error) {
		if errors.Is(err, exec.ErrWaitDelay) {
			c.Cancel()
		}
		select {
		case <-cancelled:
			<-killed
		default:
		}
	}
	return c, finish, nil
}
func (e ContextOS) Run(directory string, args []string, in io.Reader, out io.Writer) error {
	c, finish, err := e.command(directory, args)
	if err != nil {
		return err
	}
	c.Stdin = in
	c.Stdout = out
	c.Stderr = out
	err = c.Run()
	finish(err)
	return err
}
func (e ContextOS) Output(directory string, args []string) ([]byte, error) {
	c, finish, err := e.command(directory, args)
	if err != nil {
		return nil, err
	}
	body, err := c.Output()
	finish(err)
	return body, err
}
