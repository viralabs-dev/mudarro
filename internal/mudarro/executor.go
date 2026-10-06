package mudarro

import (
	"io"
	"os/exec"
)

// Executor is the operating-system boundary. Planners and adapters only describe commands.
type Executor interface {
	Run(directory string, args []string, in io.Reader, out io.Writer) error
	Output(directory string, args []string) ([]byte, error)
}
type OSExecutor struct{}

func (OSExecutor) Run(directory string, args []string, in io.Reader, out io.Writer) error {
	c := exec.Command(args[0], args[1:]...)
	c.Dir = directory
	c.Stdin = in
	c.Stdout = out
	c.Stderr = out
	return c.Run()
}
func (OSExecutor) Output(directory string, args []string) ([]byte, error) {
	c := exec.Command(args[0], args[1:]...)
	c.Dir = directory
	return c.Output()
}
func (r Runner) executor() Executor {
	if r.Executor != nil {
		return r.Executor
	}
	return OSExecutor{}
}
