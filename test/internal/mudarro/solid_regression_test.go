package mudarro_test

import (
	"bytes"
	"errors"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"github.com/viralabs-dev/mudarro/internal/mudarro/executor"
	"io"
	"strings"
	"testing"
)

type containerProbe struct {
	recordingExecutor
	probeError error
}

func (e *containerProbe) Output(dir string, args []string) ([]byte, error) {
	e.calls = append(e.calls, append([]string{}, args...))
	return e.response, e.probeError
}
func TestDockerfileOwnershipUsesExecutor(t *testing.T) {
	for _, engine := range []string{"docker", "podman"} {
		for _, state := range []string{"owned", "foreign", "absent"} {
			t.Run(engine+"/"+state, func(t *testing.T) {
				root := t.TempDir()
				s := sample().Services[0]
				s.Infrastructure = Infrastructure{Kind: engine, Mode: "dockerfile", Image: "example:test"}
				project := digest([]byte(root + "/" + s.Dir))[:12]
				probe := &containerProbe{}
				switch state {
				case "owned":
					probe.response = []byte(project + "/" + s.ID + "\n")
				case "foreign":
					probe.response = []byte("other/app")
				case "absent":
					probe.probeError = errors.New("not found")
				}
				var out bytes.Buffer
				_, a, err := findAction(Config{Services: []Service{s}}, "app:up")
				if err != nil {
					t.Fatal(err)
				}
				err = (Runner{Executor: probe, Out: &out}).Run(root, s, a)
				if state == "foreign" {
					if err == nil || !strings.Contains(err.Error(), "não pertence") || len(probe.calls) != 1 {
						t.Fatalf("foreign container accepted: %v %v", err, probe.calls)
					}
					return
				}
				if err != nil || len(probe.calls) != 2 {
					t.Fatalf("missing injected inspect: %v %v", err, probe.calls)
				}
				if probe.calls[0][1] != "inspect" {
					t.Fatal(probe.calls)
				}
				verb := "start"
				if state == "absent" {
					verb = "run"
				}
				if probe.calls[1][1] != verb {
					t.Fatal(probe.calls)
				}
			})
		}
	}
}
func TestOSExecutorRejectsEmptyCommand(t *testing.T) {
	for _, args := range [][]string{nil, {}, {""}} {
		func() {
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("command validation panicked: %v", p)
				}
			}()
			if err := (executor.OS{}).Run(t.TempDir(), args, nil, io.Discard); err == nil {
				t.Error("accepted empty run")
			}
			if _, err := (executor.OS{}).Output(t.TempDir(), args); err == nil {
				t.Error("accepted empty output")
			}
		}()
	}
}
