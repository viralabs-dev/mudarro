package executor

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// SplitEnv turns a leading `env NAME=value ...` prefix into environment
// entries when no env program exists (native Windows). Elsewhere, and when an
// env program is on PATH, args are returned unchanged and env is nil.
func SplitEnv(args []string) (command []string, env []string) {
	if runtime.GOOS != "windows" || len(args) < 2 || args[0] != "env" {
		return args, nil
	}
	if _, err := exec.LookPath("env"); err == nil {
		return args, nil
	}
	i := 1
	for ; i < len(args); i++ {
		name, _, ok := strings.Cut(args[i], "=")
		if !ok || name == "" || strings.HasPrefix(args[i], "-") {
			break
		}
		env = append(env, args[i])
	}
	if i == len(args) {
		return args, nil
	}
	return args[i:], env
}

// prepare applies SplitEnv to a command before it starts.
func prepare(args []string) ([]string, []string) {
	command, env := SplitEnv(args)
	if env != nil {
		env = append(os.Environ(), env...)
	}
	return command, env
}
