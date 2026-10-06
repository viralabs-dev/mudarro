package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"io"
	"strings"
	"testing"
)

type recordingExecutor struct {
	calls    [][]string
	response []byte
}

func (e *recordingExecutor) Run(_ string, args []string, _ io.Reader, _ io.Writer) error {
	e.calls = append(e.calls, append([]string{}, args...))
	return nil
}
func (e *recordingExecutor) Output(_ string, args []string) ([]byte, error) {
	e.calls = append(e.calls, append([]string{}, args...))
	return e.response, nil
}
func TestKubernetesStopOnlyScalesWorkloads(t *testing.T) {
	executor := &recordingExecutor{response: []byte(`{"kind":"List","items":[{"kind":"Deployment","metadata":{"name":"app"}},{"kind":"PersistentVolumeClaim","metadata":{"name":"data"}}]}`)}
	s := sample().Services[0]
	s.Infrastructure = Infrastructure{Kind: "kubernetes", Mode: "manifests", File: "k8s", Context: "test", Namespace: "sandbox"}
	var out bytes.Buffer
	if e := (Runner{Out: &out, Executor: executor}).Run(t.TempDir(), s, Action{Name: "down", Special: "k8s-down"}); e != nil {
		t.Fatal(e)
	}
	if len(executor.calls) != 2 {
		t.Fatal(executor.calls)
	}
	command := strings.Join(executor.calls[1], " ")
	if strings.Contains(command, "delete") || strings.Contains(command, "data") || !strings.Contains(command, "scale deployment/app --replicas=0") {
		t.Fatal(command)
	}
}
func TestDryRunDoesNotExposeEnvironment(t *testing.T) {
	t.Setenv("MUDARRO_TEST_SECRET", "do-not-print")
	var out bytes.Buffer
	if e := (Runner{Out: &out, Dry: true}).Run(t.TempDir(), sample().Services[0], action("probe", "scripts", "echo", "${MUDARRO_TEST_SECRET}")); e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out.String(), "do-not-print") {
		t.Fatal("secret exposed")
	}
}
