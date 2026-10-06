//go:build linux || darwin

package mudarro_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/viralabs-dev/mudarro/internal/mudarro/executor"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestContextExecutorOutputAndStdin(t *testing.T) {
	var out bytes.Buffer
	e := executor.ContextOS{Context: context.Background()}
	if err := e.Run(t.TempDir(), []string{"sh", "-c", "cat; printf '\\033[31mresult\\033[0m'"}, strings.NewReader("input\n"), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "input\n\x1b[31mresult\x1b[0m") {
		t.Fatalf("output %q", out.String())
	}
	body, err := e.Output(t.TempDir(), []string{"sh", "-c", "printf output"})
	if err != nil || string(body) != "output" {
		t.Fatalf("Output %q %v", body, err)
	}
}
func TestContextExecutorRejectsEmptyAndCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	e := executor.ContextOS{Context: ctx}
	for _, args := range [][]string{nil, {""}, {"sh", "-c", "touch SHOULD_NOT_EXIST"}} {
		if err := e.Run(t.TempDir(), args, nil, io.Discard); err == nil {
			t.Fatal("invalid/cancelled command accepted")
		}
	}
}
func TestContextExecutorCancelsOwnedProcessGroup(t *testing.T) {
	root := t.TempDir()
	pidFile := filepath.Join(root, "child.pid")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	script := `trap '' TERM; sh -c 'trap "" TERM; exec sleep 30' & child=$!; printf '%s' "$child" > child.pid; wait`
	go func() {
		done <- (executor.ContextOS{Context: ctx}).Run(root, []string{"sh", "-c", script}, nil, io.Discard)
	}()
	var pid int
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		body, e := os.ReadFile(pidFile)
		if e == nil {
			pid, _ = strconv.Atoi(string(body))
			if pid > 0 {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	if pid == 0 {
		cancel()
		<-done
		t.Fatal("child readiness")
	}
	start := time.Now()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancel no error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("group cancellation unbounded")
	}
	if time.Since(start) > time.Second {
		t.Fatal("group cancellation too slow")
	}
	// Linux can briefly retain an orphan zombie until init reaps it. Kill(0) alone
	// must not mistake that non-running process for a surviving child.
	deadline = time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		err := syscall.Kill(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return
		}
		body, _ := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
		if strings.Contains(string(body), ") Z ") {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("owned descendant still running after cancellation")
}
func TestContextExecutorBoundsInheritedOutputPipe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := (executor.ContextOS{Context: ctx}).Output(t.TempDir(), []string{"sh", "-c", `trap '' TERM; sleep 30 & wait`})
	if err == nil || time.Since(start) > 2*time.Second {
		t.Fatalf("output cancellation %v after %s", err, time.Since(start))
	}
}

func TestContextExecutorBoundsPipeAfterParentExits(t *testing.T) {
	root := t.TempDir()
	start := time.Now()
	_, err := (executor.ContextOS{}).Output(root, []string{"sh", "-c", `sleep 30 & printf '%s' "$!" > orphan.pid; exit 0`})
	if !errors.Is(err, exec.ErrWaitDelay) || time.Since(start) > 2*time.Second {
		t.Fatalf("parent exit pipe wait %v after %s", err, time.Since(start))
	}
	body, e := os.ReadFile(filepath.Join(root, "orphan.pid"))
	if e != nil {
		t.Fatal(e)
	}
	pid, e := strconv.Atoi(string(body))
	if e != nil {
		t.Fatal(e)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if errors.Is(syscall.Kill(pid, 0), syscall.ESRCH) {
			return
		}
		stat, _ := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
		if strings.Contains(string(stat), ") Z ") {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("inherited pipe child survived cleanup")
}
