//go:build linux || darwin

package mudarro_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// Every process signalled by these fixtures was created by this test in its own session.
// CLI subprocesses are necessary because local up relaunches os.Executable as the supervisor.
type supervisorIdentityState struct {
	PID     int    `json:"pid"`
	Token   string `json:"token"`
	Root    string `json:"root"`
	Service string `json:"service"`
}
type supervisorIdentityFixture struct {
	t            *testing.T
	binary, root string
	configPath   string
	owned        *supervisorIdentityState
}

func supervisorIdentityCLI(t *testing.T) *supervisorIdentityFixture {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "mudarro")
	_, file, _, _ := runtime.Caller(0)
	repo := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "./cmd/mudarro")
	build.Dir = repo
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v %s", err, out)
	}
	f := &supervisorIdentityFixture{t: t, binary: binary, root: canonicalTempDir(t)}
	f.configure("exec sleep 30")
	t.Cleanup(func() {
		if f.owned != nil {
			f.write(*f.owned)
			if out, err := f.run("down"); err != nil {
				t.Errorf("exclusive supervisor cleanup: %v %s", err, out)
			}
		}
	})
	return f
}
func (f *supervisorIdentityFixture) configure(shell string) {
	f.t.Helper()
	body := map[string]any{"version": 1, "name": "exclusive supervisor fixture", "services": []any{map[string]any{"id": "app", "dir": ".", "infrastructure": map[string]any{"kind": "local"}, "commands": map[string]any{"start": map[string]any{"shell": shell, "group": "application"}}}}}
	b, err := json.Marshal(body)
	if err != nil {
		f.t.Fatal(err)
	}
	name := f.configPath
	if name == "" {
		name = "mudarro.json"
	}
	target := filepath.Join(f.root, name)
	if err = os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		f.t.Fatal(err)
	}
	if err = os.WriteFile(target, b, 0600); err != nil {
		f.t.Fatal(err)
	}
}
func (f *supervisorIdentityFixture) path() string {
	return filepath.Join(f.root, ".mudarro/run/app/state.json")
}
func (f *supervisorIdentityFixture) run(op string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	args := []string{"run", "app:" + op, "--root", f.root}
	if f.configPath != "" {
		args = append(args, "--config", f.configPath)
	}
	out, err := exec.CommandContext(ctx, f.binary, args...).CombinedOutput()
	return string(out), err
}
func (f *supervisorIdentityFixture) state() supervisorIdentityState {
	f.t.Helper()
	b, err := os.ReadFile(f.path())
	if err != nil {
		f.t.Fatal(err)
	}
	var st supervisorIdentityState
	if err = json.Unmarshal(b, &st); err != nil {
		f.t.Fatal(err)
	}
	return st
}
func (f *supervisorIdentityFixture) write(st supervisorIdentityState) {
	f.t.Helper()
	b, err := json.Marshal(st)
	if err != nil {
		f.t.Fatal(err)
	}
	if err = os.MkdirAll(filepath.Dir(f.path()), 0700); err != nil {
		f.t.Fatal(err)
	}
	if err = os.WriteFile(f.path(), b, 0600); err != nil {
		f.t.Fatal(err)
	}
}
func (f *supervisorIdentityFixture) up() supervisorIdentityState {
	f.t.Helper()
	out, err := f.run("up")
	if err != nil {
		f.t.Fatalf("up: %v %s", err, out)
	}
	st := f.state()
	if st.PID < 2 || len(st.Token) != 32 || st.Root != f.root || st.Service != "app" {
		f.t.Fatalf("invalid ownership: %+v", st)
	}
	f.owned = &st
	return st
}
func supervisorIdentityAlive(st supervisorIdentityState) bool {
	out, err := exec.Command("ps", "-ww", "-p", fmt.Sprint(st.PID), "-o", "args=").Output()
	return err == nil && strings.Contains(string(out), "__supervise") && strings.HasSuffix(strings.TrimSpace(string(out)), st.Token)
}

func TestSupervisorIdentityRejectsInvalidStateAndPreservesOwnedProcess(t *testing.T) {
	f := supervisorIdentityCLI(t)
	original := f.up()
	cases := []struct {
		name   string
		change func(*supervisorIdentityState)
	}{
		{"pid-zero", func(s *supervisorIdentityState) { s.PID = 0 }}, {"pid-one", func(s *supervisorIdentityState) { s.PID = 1 }}, {"pid-negative", func(s *supervisorIdentityState) { s.PID = -1 }},
		{"empty-token", func(s *supervisorIdentityState) { s.Token = "" }}, {"short-token", func(s *supervisorIdentityState) { s.Token = "abc" }}, {"wrong-token", func(s *supervisorIdentityState) { s.Token = strings.Repeat("z", 32) }},
		{"foreign-root-metadata", func(s *supervisorIdentityState) { s.Root = filepath.Join(f.root, "foreign") }}, {"foreign-service-metadata", func(s *supervisorIdentityState) { s.Service = "other" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			st := original
			tc.change(&st)
			f.write(st)
			out, err := f.run("status")
			if err != nil || !strings.Contains(out, "stopped") {
				t.Fatalf("invalid state accepted: %v %s", err, out)
			}
			if out, err = f.run("down"); err != nil {
				t.Fatalf("down: %v %s", err, out)
			}
			if !supervisorIdentityAlive(original) {
				t.Fatal("invalid state signalled exclusively owned original supervisor")
			}
		})
	}
}
func TestSupervisorIdentityCorruptJSONMustNotAuthorizeSignals(t *testing.T) {
	f := supervisorIdentityCLI(t)
	original := f.up()
	b, _ := json.Marshal(original)
	for _, body := range []string{"not JSON", strings.TrimSuffix(string(b), "}") + `,"pid":"invalid type"}`} {
		if err := os.WriteFile(f.path(), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		out, err := f.run("status")
		if err != nil || !strings.Contains(out, "stopped") {
			t.Errorf("corrupt state accepted as running: %v %s", err, out)
			continue
		}
		if out, err = f.run("down"); err != nil {
			t.Errorf("down: %v %s", err, out)
		}
		if !supervisorIdentityAlive(original) {
			t.Fatal("corrupt state authorized signalling owned supervisor")
		}
	}
}
func TestSupervisorIdentityDoesNotAdoptAnotherProjectSupervisor(t *testing.T) {
	owner := supervisorIdentityCLI(t)
	foreign := owner.up()
	claimant := supervisorIdentityCLI(t)
	claimant.binary = owner.binary // Same executable: the real project arguments must still be checked.
	forged := foreign
	forged.Root = claimant.root
	forged.Service = "app"
	claimant.write(forged)
	out, err := claimant.run("status")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Errorf("another project adopted: %v %s", err, out)
	}
	// This exercises a real refusal. The only possible target is the exclusive owner fixture above.
	out, err = claimant.run("down")
	if err != nil {
		t.Errorf("down: %v %s", err, out)
	}
	if !supervisorIdentityAlive(foreign) {
		t.Fatal("forged state stopped another exclusively owned project supervisor")
	}
}
func TestSupervisorIdentityStalePIDAndUnrelatedOwnedPID(t *testing.T) {
	f := supervisorIdentityCLI(t)
	probe := exec.Command("sleep", "30")
	probe.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := probe.Start(); err != nil {
		t.Fatal(err)
	}
	waited := false
	t.Cleanup(func() {
		if !waited {
			probe.Process.Kill()
			probe.Wait()
		}
	})
	st := supervisorIdentityState{PID: probe.Process.Pid, Token: strings.Repeat("a", 32), Root: f.root, Service: "app"}
	f.write(st)
	out, err := f.run("status")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Fatal("unrelated process adopted", out, err)
	}
	if _, err = f.run("down"); err != nil {
		t.Fatal(err)
	}
	if err = probe.Process.Signal(syscall.Signal(0)); err != nil {
		t.Fatal("unrelated exclusively owned probe was signalled", err)
	}
	probe.Process.Kill()
	probe.Wait()
	waited = true
	f.write(st)
	out, err = f.run("status")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Fatal("stale PID adopted", out, err)
	}
	if _, err = f.run("down"); err != nil {
		t.Fatal(err)
	}
}
func TestSupervisorIdentityConcurrentUpAndRestart(t *testing.T) {
	f := supervisorIdentityCLI(t)
	f.configure("printf 'started\\n' >> launches.txt; exec sleep 30")
	var wg sync.WaitGroup
	results := make(chan error, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := f.run("up")
			if err != nil {
				err = fmt.Errorf("up %v %s", err, out)
			}
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	// Register cleanup before assertions, even if a concurrent up reports an error.
	st := f.state()
	f.owned = &st
	for err := range results {
		if err != nil {
			t.Error(err)
		}
	}
	b, err := os.ReadFile(filepath.Join(f.root, "launches.txt"))
	if err != nil || strings.Count(string(b), "started\n") != 1 {
		t.Fatalf("concurrent up started duplicate applications: %q %v", b, err)
	}
	if out, err := f.run("restart"); err != nil {
		t.Fatalf("restart: %v %s", err, out)
	}
	next := f.state()
	f.owned = &next
	if next.PID == st.PID || next.Token == st.Token || supervisorIdentityAlive(st) {
		t.Fatal("restart reused/leaked old identity", st, next)
	}
	b, err = os.ReadFile(filepath.Join(f.root, "launches.txt"))
	if err != nil || strings.Count(string(b), "started\n") != 2 {
		t.Fatal("restart did not launch exactly once", string(b), err)
	}
}
func TestSupervisorIdentityStartupFailureAndSupervisorSignal(t *testing.T) {
	f := supervisorIdentityCLI(t)
	f.configure("exit 7")
	if out, err := f.run("up"); err == nil {
		t.Fatal("early exit accepted", out)
	}
	st := f.state()
	if supervisorIdentityAlive(st) {
		f.owned = &st
		t.Fatal("failed startup leaked supervisor")
	}
	f.configure("printf '%s\\n' \"$$\" > child.pid; exec sleep 30")
	st = f.up()
	childBytes, err := os.ReadFile(filepath.Join(f.root, "child.pid"))
	if err != nil {
		t.Fatal(err)
	}
	childPID, err := strconv.Atoi(strings.TrimSpace(string(childBytes)))
	if err != nil || childPID < 2 {
		t.Fatal("invalid exclusive child PID", err)
	}
	// Signal only the exact supervisor PID created by this fixture, never a process selected by a broad name.
	if err := syscall.Kill(st.PID, syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for supervisorIdentityAlive(st) && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if supervisorIdentityAlive(st) {
		t.Fatal("supervisor SIGINT did not stop its child")
	}
	if err := syscall.Kill(childPID, 0); err == nil {
		t.Fatal("supervisor exited but its exclusive child survived")
	}
	out, err := f.run("status")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Fatal(out, err)
	}
}

func TestSupervisorIdentityRejectsForeignExecutableAndEmbeddedMarker(t *testing.T) {
	f := supervisorIdentityCLI(t)
	token := strings.Repeat("b", 32)
	probe := exec.Command("bash", "-c", "while :; do sleep 1; done", "unrelated-prefix", "__supervise", f.root, "app", token)
	probe.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := probe.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() { probe.Wait(); close(done) }()
	t.Cleanup(func() {
		select {
		case <-done:
			return
		default:
		}
		syscall.Kill(-probe.Process.Pid, syscall.SIGTERM)
		select {
		case <-done:
		case <-time.After(time.Second):
			syscall.Kill(-probe.Process.Pid, syscall.SIGKILL)
			<-done
		}
	})
	f.write(supervisorIdentityState{PID: probe.Process.Pid, Token: token, Root: f.root, Service: "app"})
	out, err := f.run("status")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Errorf("foreign executable/embedded marker adopted: %v %s", err, out)
	}
	if out, err = f.run("down"); err != nil {
		t.Errorf("down: %v %s", err, out)
	}
	select {
	case <-done:
		t.Fatal("embedded marker authorized signalling foreign executable fixture")
	case <-time.After(100 * time.Millisecond):
	}
}
func TestSupervisorIdentityRootAndSelectedConfigWithSpaces(t *testing.T) {
	f := supervisorIdentityCLI(t)
	f.root = filepath.Join(canonicalTempDir(t), "project with spaces")
	f.configPath = "operations/selected config.json"
	f.configure("exec sleep 30")
	st := f.up()
	out, err := f.run("status")
	if err != nil || !strings.Contains(out, "running") {
		t.Fatalf("spaces broke identity: %v %s", err, out)
	}
	if out, err = f.run("up"); err != nil {
		t.Fatal(err, out)
	}
	if next := f.state(); next != st {
		t.Fatal("space-delimited arguments created duplicate supervisor", st, next)
	}
	if out, err = f.run("down"); err != nil {
		t.Fatal(err, out)
	}
	if supervisorIdentityAlive(st) {
		t.Fatal("selected configuration supervisor remained running")
	}
}

func TestSupervisorIdentityLinuxRejectsSpoofedArgvZero(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("actual executable inode identity requires Linux /proc; Darwin runtime/identity follow-up remains MUD-031")
	}
	f := supervisorIdentityCLI(t)
	dir := t.TempDir()
	source := filepath.Join(dir, "foreign.go")
	foreign := filepath.Join(dir, "foreign-binary")
	if err := os.WriteFile(source, []byte("package main\nimport \"time\"\nfunc main(){time.Sleep(30*time.Second)}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", foreign, source)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatal(err, string(out))
	}
	token := strings.Repeat("c", 32)
	probe := exec.Command(foreign)
	probe.Args = []string{f.binary, "__supervise", f.root, "app", "mudarro.json", token}
	probe.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := probe.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() { probe.Wait(); close(done) }()
	t.Cleanup(func() {
		select {
		case <-done:
			return
		default:
		}
		syscall.Kill(-probe.Process.Pid, syscall.SIGTERM)
		<-done
	})
	f.write(supervisorIdentityState{PID: probe.Process.Pid, Token: token, Root: f.root, Service: "app"})
	out, err := f.run("status")
	if err != nil || !strings.Contains(out, "stopped") {
		t.Errorf("spoofed argv[0] adopted foreign executable: %v %s", err, out)
	}
	if out, err = f.run("down"); err != nil {
		t.Errorf("down: %v %s", err, out)
	}
	select {
	case <-done:
		t.Fatal("spoofed argv[0] authorized signalling foreign executable fixture")
	case <-time.After(100 * time.Millisecond):
	}
}
