package mudarro

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/viralabs-dev/mudarro/internal/mudarro/executor"
	"github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
)

func runShellAction(view *terminal.Menu, root string, c Config, s Service, a Action) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	input, stop, err := view.ExecutionInput(ctx, cancel)
	if err != nil {
		return err
	}
	output := view.ShellOutput()
	name := ""
	if strings.Contains(strings.Join(a.Command.Args, " "), "{name}") {
		fmt.Fprint(output, uiText(c, "Migration name: ", "Nome da migration: "))
		name, err = view.ReadShellLine(input, output)
	}
	if err == nil {
		err = (Runner{Executor: executor.ContextOS{Context: ctx}, In: input, Out: output, Name: strings.TrimSpace(name), Locale: c.UIOptions().Locale, ConfigPath: c.configPath}).Run(root, s, a)
	}
	stop()
	if err != nil {
		fmt.Fprintln(output, uiText(c, "Error:", "Erro:"), err)
	}
	if ctx.Err() != nil {
		fmt.Fprintln(output, uiText(c, "Execution cancelled.", "Execução cancelada."))
		return nil
	}
	// A fresh input lifetime preserves the next menu's ability to execute after cancellation.
	waitCtx, waitCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer waitCancel()
	input, stop, err = view.ExecutionInput(waitCtx, waitCancel)
	if err != nil {
		return err
	}
	defer stop()
	fmt.Fprint(output, uiText(c, "\nEnter to return to the menu...", "\nEnter para voltar ao menu..."))
	_, err = view.ReadShellLine(input, output)
	if err == io.EOF || waitCtx.Err() != nil {
		return nil
	}
	return err
}
