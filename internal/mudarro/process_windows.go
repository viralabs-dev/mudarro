//go:build windows

package mudarro

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/viralabs-dev/mudarro/internal/mudarro/winapi"
)

func lockRunFile(f *os.File) error   { return winapi.LockFile(f) }
func unlockRunFile(f *os.File) error { return winapi.UnlockFile(f) }

// The supervisor runs without a console and in its own process group, so
// closing the terminal or pressing Ctrl+C in it does not stop the service.
func supervisorProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: winapi.CreateNewProcessGroup | winapi.DetachedProcess, HideWindow: true}
}

// stopSupervisor ends the supervisor. Its Job Object is kill-on-close, so the
// whole service tree ends with it. Windows has no SIGTERM for a process
// without a console; services must tolerate a forced stop.
func stopSupervisor(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	defer p.Release()
	return p.Kill()
}

// supervisorMutexName binds the supervisor argv (root, service, config and
// token) to a named mutex the supervisor holds for its lifetime.
func supervisorMutexName(args []string) string {
	sum := sha256.Sum256([]byte(strings.Join(args, "\x00")))
	return `Local\mudarro-supervisor-` + hex.EncodeToString(sum[:])
}

// supervisorIdentity cannot read another process's argv without undocumented
// APIs. It requires the live PID to run this executable and the argv-bound
// mutex to exist; a reused PID or a forged process fails one of them.
func supervisorIdentity(st processState, expected []string) bool {
	image, alive := winapi.ProcessAlive(st.PID)
	if !alive || !winapi.SameExecutable(image, expected[0]) {
		return false
	}
	return winapi.MutexExists(supervisorMutexName(expected[1:]))
}

func holdSupervisorIdentity(args []string) (func(), error) {
	return winapi.HoldMutex(supervisorMutexName(append([]string{"__supervise"}, args...)))
}

type supervisedChild struct {
	cmd *exec.Cmd
	job *winapi.Job
}

// startSupervised assigns the service to a kill-on-close Job Object before it
// runs (it starts suspended), so its descendants end when the supervisor ends,
// however it ends.
func startSupervised(c *exec.Cmd) (*supervisedChild, error) {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.CreationFlags |= winapi.CreateSuspended
	job, err := winapi.NewJob(true)
	if err != nil {
		return nil, err
	}
	if err = c.Start(); err != nil {
		job.Close()
		return nil, err
	}
	if err = job.Assign(c.Process.Pid); err == nil {
		err = winapi.ResumeProcess(c.Process.Pid)
	}
	if err != nil {
		_ = c.Process.Kill()
		job.Close()
		_ = c.Wait()
		return nil, err
	}
	return &supervisedChild{cmd: c, job: job}, nil
}
func (s *supervisedChild) terminate() error { return s.job.Terminate() }
func (s *supervisedChild) kill() error      { return s.job.Terminate() }
func (s *supervisedChild) release()         { s.job.Close() }

// supervisorStarted waits for the supervisor to publish its identity mutex
// (process start is slower on Windows), then still requires it to survive
// the same fail-fast window as on Unix.
func supervisorStarted(st processState, configPath string) bool {
	deadline := time.Now().Add(5 * time.Second)
	for !running(st, configPath) {
		if _, alive := winapi.ProcessAlive(st.PID); !alive || time.Now().After(deadline) {
			return false
		}
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(150 * time.Millisecond)
	return running(st, configPath)
}
