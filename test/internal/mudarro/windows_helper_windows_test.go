//go:build windows

package mudarro_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"testing"
	"time"
)

// The Windows fixtures re-run this test binary as a child process instead of
// depending on sh, bash or sleep, which a stock Windows runner does not have.
const (
	windowsHelperMode    = "MUDARRO_WINDOWS_HELPER"
	windowsHelperPIDFile = "MUDARRO_WINDOWS_HELPER_PIDFILE"
)

func windowsHelperArgs(t *testing.T) []string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return []string{exe, "-test.run=^TestWindowsHelperProcess$"}
}

// TestWindowsHelperProcess is not a test: it only acts when a fixture sets the mode.
func TestWindowsHelperProcess(t *testing.T) {
	mode := os.Getenv(windowsHelperMode)
	if mode == "" {
		return
	}
	switch mode {
	case "echo":
		io.Copy(os.Stdout, os.Stdin)
		fmt.Print("\x1b[31mresult\x1b[0m")
	case "print":
		fmt.Print("output")
	case "exit3":
		fmt.Fprint(os.Stderr, "failure detail")
		os.Exit(3)
	case "tree":
		// Ignore CTRL_BREAK so that only Job Object termination can end the tree.
		signal.Ignore(os.Interrupt)
		exe, _ := os.Executable()
		child := exec.Command(exe, "-test.run=^TestWindowsHelperProcess$")
		child.Env = append(os.Environ(), windowsHelperMode+"=sleep")
		if err := child.Start(); err != nil {
			os.Exit(2)
		}
		time.Sleep(60 * time.Second)
	case "sleep":
		signal.Ignore(os.Interrupt)
		if path := os.Getenv(windowsHelperPIDFile); path != "" {
			tmp := path + ".tmp"
			if os.WriteFile(tmp, []byte(strconv.Itoa(os.Getpid())), 0600) == nil {
				os.Rename(tmp, path)
			}
		}
		time.Sleep(60 * time.Second)
	default:
		os.Exit(2)
	}
	os.Exit(0)
}

func waitWindowsPIDFile(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if body, err := os.ReadFile(path); err == nil {
			if pid, err := strconv.Atoi(string(body)); err == nil && pid > 0 {
				return pid
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("helper process did not report its PID")
	return 0
}
