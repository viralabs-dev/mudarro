package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"io"
	"os"
	"testing"
)

type confirmationExecutor struct{ got string }

func (e *confirmationExecutor) Run(_ string, _ []string, in io.Reader, _ io.Writer) error {
	b, err := io.ReadAll(in)
	e.got = string(b)
	return err
}
func (e *confirmationExecutor) Output(string, []string) ([]byte, error) { return nil, nil }
func TestDestructiveConfirmationPreservesQueuedChildInput(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	w.WriteString("APAGAR app\nqueued-child-input\n")
	w.Close()
	e := &confirmationExecutor{}
	var out bytes.Buffer
	err = (Runner{In: r, Out: &out, Executor: e}).Run(t.TempDir(), Service{ID: "app", Dir: "."}, Action{Name: "confirm", Command: Command{Args: []string{"true"}, Destructive: true}})
	if err != nil || e.got != "queued-child-input\n" {
		t.Fatalf("queued input %q error %v", e.got, err)
	}
}
