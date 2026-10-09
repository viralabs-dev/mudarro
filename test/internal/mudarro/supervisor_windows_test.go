//go:build windows

package mudarro_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/viralabs-dev/mudarro/internal/mudarro/winapi"
)

type windowsSupervisor struct {
	t                     *testing.T
	binary, root, pidFile string
}

// The CLI is built because local up relaunches os.Executable as the supervisor.
func windowsSupervisorCLI(t *testing.T) *windowsSupervisor {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "mudarro.exe")
	_, file, _, _ := runtime.Caller(0)
	repo := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go.exe"), "build", "-o", binary, "./cmd/mudarro")
	build.Dir = repo
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v %s", err, out)
	}
	f := &windowsSupervisor{t: t, binary: binary, root: t.TempDir()}
	f.pidFile = filepath.Join(f.root, "service.pid")
	body := map[string]any{"version": 1, "name": "windows supervisor fixture", "services": []any{map[string]any{"id": "app", "dir": ".", "infrastructure": map[string]any{"kind": "local"}, "commands": map[string]any{"start": map[string]any{"args": windowsHelperArgs(t), "group": "application"}}}}}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(f.root, "mudarro.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.run("down") })
	return f
}

func (f *windowsSupervisor) run(op string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, f.binary, "run", "app:"+op, "--root", f.root)
	// The supervisor inherits this environment and passes it to the service.
	c.Env = append(os.Environ(), windowsHelperMode+"=tree", windowsHelperPIDFile+"="+f.pidFile)
	out, err := c.CombinedOutput()
	return string(out), err
}

func (f *windowsSupervisor) state() supervisorWindowsState {
	f.t.Helper()
	b, err := os.ReadFile(filepath.Join(f.root, ".mudarro", "run", "app", "state.json"))
	if err != nil {
		f.t.Fatal(err)
	}
	var st supervisorWindowsState
	if err = json.Unmarshal(b, &st); err != nil {
		f.t.Fatal(err)
	}
	return st
}

type supervisorWindowsState struct {
	PID     int    `json:"pid"`
	Token   string `json:"token"`
	Root    string `json:"root"`
	Service string `json:"service"`
}

func TestWindowsLocalSupervisorUpStatusDownEndsTree(t *testing.T) {
	f := windowsSupervisorCLI(t)
	if out, err := f.run("up"); err != nil || !strings.Contains(out, "started") {
		t.Fatalf("up: %v %s", err, out)
	}
	st := f.state()
	if st.PID < 2 || len(st.Token) != 32 || st.Root != f.root || st.Service != "app" {
		t.Fatalf("invalid ownership: %+v", st)
	}
	grandchild := waitWindowsPIDFile(t, f.pidFile)
	if out, err := f.run("status"); err != nil || !strings.Contains(out, "running") {
		t.Fatalf("status: %v %s", err, out)
	}
	if out, err := f.run("up"); err != nil || !strings.Contains(out, "already running") {
		t.Fatalf("second up: %v %s", err, out)
	}
	if out, err := f.run("down"); err != nil || !strings.Contains(out, "stopped") {
		t.Fatalf("down: %v %s", err, out)
	}
	if out, err := f.run("status"); err != nil || !strings.Contains(out, "stopped") {
		t.Fatalf("status after down: %v %s", err, out)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		_, supervisorAlive := winapi.ProcessAlive(st.PID)
		_, serviceAlive := winapi.ProcessAlive(grandchild)
		if !supervisorAlive && !serviceAlive {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("supervisor or service descendant survived down")
}

func TestWindowsLocalSupervisorRejectsForgedState(t *testing.T) {
	f := windowsSupervisorCLI(t)
	if out, err := f.run("up"); err != nil {
		t.Fatalf("up: %v %s", err, out)
	}
	original := f.state()
	statePath := filepath.Join(f.root, ".mudarro", "run", "app", "state.json")
	for name, change := range map[string]func(*supervisorWindowsState){
		"wrong-token":  func(s *supervisorWindowsState) { s.Token = strings.Repeat("0", 32) },
		"foreign-root": func(s *supervisorWindowsState) { s.Root = filepath.Join(f.root, "foreign") },
		"foreign-pid":  func(s *supervisorWindowsState) { s.PID = os.Getpid() },
	} {
		st := original
		change(&st)
		b, _ := json.Marshal(st)
		if err := os.WriteFile(statePath, b, 0600); err != nil {
			t.Fatal(err)
		}
		if out, err := f.run("status"); err != nil || !strings.Contains(out, "stopped") {
			t.Fatalf("%s: forged state accepted: %v %s", name, err, out)
		}
		if out, err := f.run("down"); err != nil {
			t.Fatalf("%s: down: %v %s", name, err, out)
		}
		if _, alive := winapi.ProcessAlive(original.PID); !alive {
			t.Fatalf("%s: forged state stopped the owned supervisor", name)
		}
	}
	b, _ := json.Marshal(original)
	if err := os.WriteFile(statePath, b, 0600); err != nil {
		t.Fatal(err)
	}
}
