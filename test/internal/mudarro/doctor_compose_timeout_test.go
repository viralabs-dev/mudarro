//go:build linux || darwin

package mudarro_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestDoctorComposeVersionProbeHasDeadline(t *testing.T) {
	sleep, e := exec.LookPath("sleep")
	if e != nil {
		t.Fatal(e)
	}
	provider := t.TempDir()
	root := doctorFixture(t, mudarro.Infrastructure{Kind: "docker", Mode: "compose", File: "compose.yaml"}, nil)
	calls := filepath.Join(root, "probe.calls")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '" + calls + "'\nexec '" + sleep + "' 60\n"
	if e := os.WriteFile(filepath.Join(provider, "docker"), []byte(script), 0755); e != nil {
		t.Fatal(e)
	}
	config, e := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDoctorComposeTimeoutIsolatedProbe$")
	child.Env = append(os.Environ(), "MUDARRO_DOCTOR_TIMEOUT_PROBE="+root, "PATH="+provider)
	child.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var output bytes.Buffer
	child.Stdout = &output
	child.Stderr = &output
	start := time.Now()
	if e := child.Start(); e != nil {
		t.Fatal(e)
	}
	// This process group belongs exclusively to the helper/provider fixture.
	// Ensure provider cleanup even if the outer watchdog finds a regression.
	defer syscall.Kill(-child.Process.Pid, syscall.SIGKILL)
	e = child.Wait()
	elapsed := time.Since(start)
	if ctx.Err() != nil {
		t.Fatalf("provider exceeded bounded diagnostic deadline; owned group cleaned: %v", ctx.Err())
	}
	if e != nil {
		t.Fatalf("helper failed: %v %s", e, output.String())
	}
	if elapsed >= 15*time.Second {
		t.Fatalf("unbounded probe duration: %s", elapsed)
	}
	var result struct{ Output, Error string }
	if e := json.Unmarshal(output.Bytes(), &result); e != nil {
		t.Fatalf("invalid helper result: %v %s", e, output.String())
	}
	if result.Error == "" {
		t.Fatal("hanging provider reported success")
	}
	lower := strings.ToLower(result.Output)
	if !strings.Contains(lower, "timeout") && !strings.Contains(lower, "timed out") && !strings.Contains(lower, "deadline") {
		t.Fatalf("timeout not distinguished: %q", result.Output)
	}
	if strings.Contains(result.Output, "AUSENTE app: provedor Compose") {
		t.Fatal("timeout mislabeled provider missing")
	}
	recorded, e := os.ReadFile(calls)
	if e != nil || string(recorded) != "compose version\n" {
		t.Fatalf("unexpected provider calls: %q %v", recorded, e)
	}
	after, e := os.ReadFile(filepath.Join(root, "mudarro.yaml"))
	if e != nil || !bytes.Equal(config, after) {
		t.Fatal("doctor mutated config", e)
	}
}

func TestDoctorComposeTimeoutIsolatedProbe(t *testing.T) {
	root := os.Getenv("MUDARRO_DOCTOR_TIMEOUT_PROBE")
	if root == "" {
		return
	}
	output, e := doctorRun(root)
	result := struct{ Output, Error string }{Output: output}
	if e != nil {
		result.Error = e.Error()
	}
	if e := json.NewEncoder(os.Stdout).Encode(result); e != nil {
		os.Exit(2)
	}
	os.Exit(0)
}
