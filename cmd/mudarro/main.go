package main

import (
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
)

var version = "dev"

func main() {
	if err := mudarro.Main(os.Args[1:], version, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "mudarro:", err)
		if e, ok := err.(interface{ ExitCode() int }); ok && e.ExitCode() > 0 {
			os.Exit(e.ExitCode())
		}
		os.Exit(1)
	}
}
