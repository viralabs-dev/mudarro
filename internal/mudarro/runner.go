package mudarro

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Runner struct {
	Executor   Executor
	In         io.Reader
	Out        io.Writer
	Dry        bool
	Name       string
	Locale     string
	ConfigPath string
}

func (r Runner) Run(root string, s Service, a Action) error {
	if a.Blocked != "" {
		return fmt.Errorf("%s: %s", a.Name, a.Blocked)
	}
	if a.Command.Destructive && !r.Dry {
		expected := "APAGAR " + s.ID
		fmt.Fprintf(r.Out, r.text("Destructive operation %s. Type %q: ", "Operação destrutiva %s. Digite %q: "), a.Name, expected)
		line, e := bufio.NewReader(r.In).ReadString('\n')
		if e != nil || strings.TrimSpace(line) != expected {
			return fmt.Errorf("operação cancelada")
		}
	}
	if strings.HasPrefix(a.Special, "local-") {
		return r.local(root, s, strings.TrimPrefix(a.Special, "local-"))
	}
	if strings.HasPrefix(a.Special, "k8s-") {
		return r.kubernetes(root, s, strings.TrimPrefix(a.Special, "k8s-"))
	}
	if a.Special == "sqlite-init" {
		return r.sqliteInit(root, s)
	}
	if r.Dry {
		fmt.Fprintf(r.Out, "%s: %s\n", s.ID, jsonString(a.Command))
		return nil
	}
	args, e := expandArgs(a.Command.Args, r.Name)
	if e != nil {
		return e
	}
	if a.Command.Shell != "" {
		args = []string{"bash", "-c", a.Command.Shell}
	}
	if len(args) == 0 {
		return fmt.Errorf("ação sem comando: %s", a.Name)
	}
	dir, e := serviceDirectory(root, s)
	if e != nil {
		return e
	}
	project := digest([]byte(root + "/" + s.Dir))[:12]
	for j := range args {
		args[j] = strings.ReplaceAll(args[j], "{project}", project)
	}
	// Dockerfile up is repeatable: resume a stopped owned container.
	if (a.Name == "up" || a.Name == "down" || a.Name == "restart") && s.Infrastructure.Mode == "dockerfile" && (s.Infrastructure.Kind == "docker" || s.Infrastructure.Kind == "podman") {
		engine := s.Infrastructure.Kind
		n := "mudarro-" + project + "-" + s.ID
		inspect := []string{engine, "inspect", "--format", `{{index .Config.Labels "dev.viralabs.mudarro"}}`, n}
		if out, e := r.executor().Output(dir, inspect); e == nil {
			if strings.TrimSpace(string(out)) != project+"/"+s.ID {
				return fmt.Errorf("container %s não pertence a este serviço", n)
			}
			if a.Name == "up" {
				args = []string{engine, "start", n}
			}
		}
	}
	return r.execute(dir, args)
}
func (r Runner) execute(dir string, args []string) error {
	return r.executor().Run(dir, args, r.In, r.Out)
}

func (r Runner) text(en, pt string) string {
	if r.Locale == "pt-BR" {
		return pt
	}
	return en
}
