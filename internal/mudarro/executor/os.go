package executor

import (
	"fmt"
	"io"
	"os/exec"
)

type OS struct{}

func (OS) Run(directory string, args []string, in io.Reader, out io.Writer) error {
	if len(args) == 0 || args[0] == "" {
		return fmt.Errorf("comando vazio")
	}
	args, env := prepare(args)
	c := exec.Command(args[0], args[1:]...)
	c.Dir = directory
	c.Env = env
	c.Stdin = in
	c.Stdout = out
	c.Stderr = out
	return c.Run()
}
func (OS) Output(directory string, args []string) ([]byte, error) {
	if len(args) == 0 || args[0] == "" {
		return nil, fmt.Errorf("comando vazio")
	}
	args, env := prepare(args)
	c := exec.Command(args[0], args[1:]...)
	c.Dir = directory
	c.Env = env
	return c.Output()
}
