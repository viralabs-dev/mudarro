//go:build !windows

package mudarro

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func lockRunFile(f *os.File) error   { return syscall.Flock(int(f.Fd()), syscall.LOCK_EX) }
func unlockRunFile(f *os.File) error { return syscall.Flock(int(f.Fd()), syscall.LOCK_UN) }

// The supervisor leads its own session, so its process group is the service tree.
func supervisorProcAttr() *syscall.SysProcAttr { return &syscall.SysProcAttr{Setsid: true} }

func stopSupervisor(pid int) error { return syscall.Kill(-pid, syscall.SIGTERM) }

// supervisorIdentity compares the live process with the exact supervisor argv.
func supervisorIdentity(st processState, expected []string) bool {
	if runtime.GOOS == "linux" {
		// argv[0] can be forged: also require the actual executable inode.
		self, err := os.Stat("/proc/self/exe")
		if err != nil {
			return false
		}
		other, err := os.Stat(filepath.Join("/proc", strconv.Itoa(st.PID), "exe"))
		if err != nil || !os.SameFile(self, other) {
			return false
		}
		// NUL separators preserve argument boundaries, including roots/configs with spaces.
		b, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(st.PID), "cmdline"))
		return err == nil && bytes.Equal(b, []byte(strings.Join(expected, "\x00")+"\x00"))
	}
	// Darwin ps does not expose NUL-separated argv. Compare the entire expected
	// command, never search for a supervisor/token substring in an unrelated process.
	b, err := exec.Command("ps", "-ww", "-p", strconv.Itoa(st.PID), "-o", "args=").Output()
	return err == nil && strings.TrimSuffix(string(b), "\n") == strings.Join(expected, " ")
}

// holdSupervisorIdentity is a no-op on Unix: the kernel already exposes argv.
func holdSupervisorIdentity([]string) (func(), error) { return func() {}, nil }

type supervisedChild struct{ cmd *exec.Cmd }

func startSupervised(c *exec.Cmd) (*supervisedChild, error) {
	if err := c.Start(); err != nil {
		return nil, err
	}
	return &supervisedChild{cmd: c}, nil
}
func (s *supervisedChild) terminate() error { return s.cmd.Process.Signal(syscall.SIGTERM) }
func (s *supervisedChild) kill() error      { return s.cmd.Process.Kill() }
func (s *supervisedChild) release()         {}

// supervisorStarted gives the supervisor a moment to fail fast; argv is
// visible from the first instruction, so one check suffices.
func supervisorStarted(st processState, configPath string) bool {
	time.Sleep(150 * time.Millisecond)
	return running(st, configPath)
}
