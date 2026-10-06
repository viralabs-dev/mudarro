package mudarro

import (
	"github.com/viralabs-dev/mudarro/internal/mudarro/executor"
	"io"
)

// Executor is the operating-system boundary. Planners and adapters only describe commands.
type Executor interface {
	Run(directory string, args []string, in io.Reader, out io.Writer) error
	Output(directory string, args []string) ([]byte, error)
}

func (r Runner) executor() Executor {
	if r.Executor != nil {
		return r.Executor
	}
	return executor.OS{}
}
