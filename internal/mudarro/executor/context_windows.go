//go:build windows

package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/viralabs-dev/mudarro/internal/mudarro/winapi"
)

// ContextOS cancels only the process tree it creates. Each command runs in a
// new process group and in its own Job Object: cancellation first sends
// CTRL_BREAK_EVENT to the group (the Windows counterpart of SIGTERM), then
// terminates the job, which ends every descendant that was assigned to it.
// Processes left running after a normal exit are not killed, as on Unix.
// It never changes console state; readers are never closed by this executor.
type ContextOS struct{ Context context.Context }

type windowsCommand struct {
	cmd       *exec.Cmd
	job       *winapi.Job
	cancelled chan struct{}
	killed    chan struct{}
	once      sync.Once
	mu        sync.Mutex
}

func (e ContextOS) command(directory string, args []string) (*windowsCommand, error) {
	if len(args) == 0 || args[0] == "" {
		return nil, fmt.Errorf("comando vazio")
	}
	ctx := e.Context
	if ctx == nil {
		ctx = context.Background()
	}
	args, env := prepare(args)
	c := exec.CommandContext(ctx, args[0], args[1:]...)
	c.Dir = directory
	c.Env = env
	c.SysProcAttr = &syscall.SysProcAttr{CreationFlags: winapi.CreateNewProcessGroup | winapi.CreateSuspended}
	c.WaitDelay = 750 * time.Millisecond
	w := &windowsCommand{cmd: c, cancelled: make(chan struct{}, 1), killed: make(chan struct{})}
	c.Cancel = w.cancel
	return w, nil
}

func (w *windowsCommand) terminate() {
	w.mu.Lock()
	job := w.job
	w.mu.Unlock()
	if job == nil || job.Terminate() != nil {
		_ = w.cmd.Process.Kill()
	}
}

func (w *windowsCommand) cancel() error {
	var err error
	w.once.Do(func() {
		if w.cmd.Process == nil {
			err = os.ErrProcessDone
			return
		}
		// Graceful first; a process without our console simply misses the event.
		if winapi.CtrlBreak(w.cmd.Process.Pid) != nil {
			w.terminate()
			close(w.killed)
			w.cancelled <- struct{}{}
			return
		}
		w.cancelled <- struct{}{}
		time.AfterFunc(300*time.Millisecond, func() { w.terminate(); close(w.killed) })
	})
	return err
}

func (w *windowsCommand) start() error {
	if err := w.cmd.Start(); err != nil {
		return err
	}
	// The process starts suspended, so it joins the job before it can create
	// descendants. Without a job (for example, a host that denies nested jobs),
	// cancellation still ends the direct child.
	job, err := winapi.NewJob(false)
	if err == nil {
		if err = job.Assign(w.cmd.Process.Pid); err != nil {
			job.Close()
			job = nil
		}
	}
	w.mu.Lock()
	w.job = job
	w.mu.Unlock()
	if err = winapi.ResumeProcess(w.cmd.Process.Pid); err != nil {
		w.terminate()
		_ = w.cmd.Wait()
		w.closeJob()
		return fmt.Errorf("resume suspended process: %w", err)
	}
	return nil
}

func (w *windowsCommand) wait() error {
	err := w.cmd.Wait()
	if errors.Is(err, exec.ErrWaitDelay) {
		w.cmd.Cancel()
	}
	select {
	case <-w.cancelled:
		<-w.killed
	default:
	}
	w.closeJob()
	return err
}

func (w *windowsCommand) closeJob() {
	w.mu.Lock()
	if w.job != nil {
		w.job.Close()
		w.job = nil
	}
	w.mu.Unlock()
}

func (e ContextOS) Run(directory string, args []string, in io.Reader, out io.Writer) error {
	w, err := e.command(directory, args)
	if err != nil {
		return err
	}
	w.cmd.Stdin = in
	w.cmd.Stdout = out
	w.cmd.Stderr = out
	if err = w.start(); err != nil {
		return err
	}
	return w.wait()
}

func (e ContextOS) Output(directory string, args []string) ([]byte, error) {
	w, err := e.command(directory, args)
	if err != nil {
		return nil, err
	}
	var stdout, stderr bytes.Buffer
	w.cmd.Stdout = &stdout
	w.cmd.Stderr = &stderr
	if err = w.start(); err != nil {
		return nil, err
	}
	err = w.wait()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitErr.Stderr = stderr.Bytes()
	}
	return stdout.Bytes(), err
}
