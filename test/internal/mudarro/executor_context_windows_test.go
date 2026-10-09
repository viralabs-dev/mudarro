//go:build windows

package mudarro_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/viralabs-dev/mudarro/internal/mudarro/executor"
	"github.com/viralabs-dev/mudarro/internal/mudarro/winapi"
)

func TestContextExecutorWindowsOutputAndStdin(t *testing.T) {
	t.Setenv(windowsHelperMode, "echo")
	var out bytes.Buffer
	e := executor.ContextOS{Context: context.Background()}
	if err := e.Run(t.TempDir(), windowsHelperArgs(t), strings.NewReader("input\n"), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "input\n\x1b[31mresult\x1b[0m") {
		t.Fatalf("output %q", out.String())
	}
	t.Setenv(windowsHelperMode, "print")
	body, err := e.Output(t.TempDir(), windowsHelperArgs(t))
	if err != nil || string(body) != "output" {
		t.Fatalf("Output %q %v", body, err)
	}
}

func TestContextExecutorWindowsExitCodeAndStderr(t *testing.T) {
	t.Setenv(windowsHelperMode, "exit3")
	_, err := (executor.ContextOS{Context: context.Background()}).Output(t.TempDir(), windowsHelperArgs(t))
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 3 || !strings.Contains(string(exitErr.Stderr), "failure detail") {
		t.Fatalf("exit status not propagated: %v", err)
	}
}

func TestContextExecutorWindowsRejectsEmptyAndCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	root := t.TempDir()
	e := executor.ContextOS{Context: ctx}
	for _, args := range [][]string{nil, {""}, {"cmd", "/c", "type nul > SHOULD_NOT_EXIST"}} {
		if err := e.Run(root, args, nil, io.Discard); err == nil {
			t.Fatal("invalid/cancelled command accepted")
		}
	}
	if _, err := os.Stat(filepath.Join(root, "SHOULD_NOT_EXIST")); err == nil {
		t.Fatal("cancelled command ran")
	}
}

// The child ignores CTRL_BREAK and its own child is created immediately after
// start: the suspended start plus the Job Object must still end both.
func TestContextExecutorWindowsCancelsOwnedProcessTree(t *testing.T) {
	root := t.TempDir()
	pidFile := filepath.Join(root, "grandchild.pid")
	t.Setenv(windowsHelperMode, "tree")
	t.Setenv(windowsHelperPIDFile, pidFile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- (executor.ContextOS{Context: ctx}).Run(root, windowsHelperArgs(t), nil, io.Discard)
	}()
	grandchild := waitWindowsPIDFile(t, pidFile)
	if _, alive := winapi.ProcessAlive(grandchild); !alive {
		t.Fatal("grandchild not running before cancellation")
	}
	started := time.Now()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancelled tree reported success")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancellation did not finish")
	}
	if elapsed := time.Since(started); elapsed > 3*time.Second {
		t.Fatalf("cancellation took %s", elapsed)
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, alive := winapi.ProcessAlive(grandchild); !alive {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("grandchild %d survived cancellation", grandchild)
}

func TestExecutorSplitEnvWithoutEnvProgram(t *testing.T) {
	if _, err := exec.LookPath("env"); err == nil {
		t.Skip("an env program is on PATH; the prefix is passed through unchanged")
	}
	args, env := executor.SplitEnv([]string{"env", "GOWORK=off", "go", "test", "./..."})
	if strings.Join(args, " ") != "go test ./..." || len(env) != 1 || env[0] != "GOWORK=off" {
		t.Fatalf("SplitEnv = %q %q", args, env)
	}
	args, env = executor.SplitEnv([]string{"go", "version"})
	if strings.Join(args, " ") != "go version" || env != nil {
		t.Fatalf("plain command changed: %q %q", args, env)
	}
}
