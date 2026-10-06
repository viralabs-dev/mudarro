package mudarro

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type options struct {
	root, name, selectScripts   string
	json, dry, all, interactive bool
	exclude                     []string
	pos                         []string
}

func parseOptions(args []string) (options, error) {
	o := options{root: "."}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--json":
			o.json = true
		case "--dry-run":
			o.dry = true
		case "--all-scripts":
			o.all = true
		case "--interactive":
			o.interactive = true
		case "--root", "--name", "--select", "--exclude":
			if i+1 >= len(args) {
				return o, fmt.Errorf("valor ausente: %s", a)
			}
			i++
			switch a {
			case "--root":
				o.root = args[i]
			case "--name":
				o.name = args[i]
			case "--select":
				o.selectScripts = args[i]
			case "--exclude":
				o.exclude = append(o.exclude, args[i])
			}
		default:
			if strings.HasPrefix(a, "-") {
				return o, fmt.Errorf("opção desconhecida: %s", a)
			}
			o.pos = append(o.pos, a)
		}
	}
	root, e := filepath.Abs(o.root)
	if e != nil {
		return o, e
	}
	o.root, e = filepath.EvalSymlinks(root)
	return o, e
}
func Main(args []string, version string, in io.Reader, out io.Writer) error {
	if len(args) == 0 {
		args = []string{"help"}
	}
	if args[0] == "__supervise" {
		return supervise(args[1:])
	}
	if args[0] == "version" || args[0] == "--version" {
		fmt.Fprintln(out, version)
		return nil
	}
	if args[0] == "help" || args[0] == "--help" {
		fmt.Fprint(out, "Mudarro — menus operacionais sem IA\n\n  scan [--json] [--exclude caminho]\n  init [--interactive | --select serviço:script,... | --all-scripts]\n  generate [--dry-run]\n  menu\n  run serviço:ação [--name nome] [--dry-run]\n  doctor\n  version\n\nTodos os comandos aceitam --root diretório. Configure mudarro.yaml antes de gerar.\n")
		return nil
	}
	o, e := parseOptions(args[1:])
	if e != nil {
		return e
	}
	switch args[0] {
	case "scan", "init":
		if args[0] == "scan" && exists(filepath.Join(o.root, "mudarro.yaml")) {
			c, e := Load(o.root)
			if e != nil {
				return e
			}
			o.exclude = append(o.exclude, c.Exclude...)
		}
		report, e := Scan(o.root, o.exclude)
		if e != nil {
			return e
		}
		if o.name != "" {
			report.Config.Name = o.name
		}
		if args[0] == "scan" {
			if o.json {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}
			printReport(out, report)
			return nil
		}
		if exists(filepath.Join(o.root, "mudarro.yaml")) {
			return fmt.Errorf("mudarro.yaml já existe; edite a configuração preservada")
		}
		selected := map[string]bool{}
		for _, v := range strings.Split(o.selectScripts, ",") {
			if v != "" {
				selected[v] = true
			}
		}
		reader := bufio.NewReader(in)
		available := map[string]bool{}
		for _, su := range report.Suggestions {
			key := su.Service + ":" + su.Name
			available[key] = true
			include := o.all || selected[key]
			if o.interactive {
				fmt.Fprintf(out, "Incluir %s (%s)? [s/N] ", key, su.Evidence)
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("entrada encerrada antes de salvar configuração")
				}
				include = strings.EqualFold(strings.TrimSpace(line), "s")
			}
			if include {
				for j := range report.Config.Services {
					if report.Config.Services[j].ID == su.Service {
						if len(su.Command.Args) > 0 && su.Command.Args[0] == "" {
							return fmt.Errorf("defina o gerenciador antes de incluir %s", key)
						}
						report.Config.Services[j].Commands[su.Name] = su.Command
					}
				}
			}
		}
		for key := range selected {
			if !available[key] {
				return fmt.Errorf("script sugerido inexistente: %s", key)
			}
		}
		if e = saveConfig(filepath.Join(o.root, "mudarro.yaml"), report.Config); e != nil {
			return e
		}
		printReport(out, report)
		fmt.Fprintln(out, "Criado mudarro.yaml. Revise as pendências e execute mudarro generate.")
		return nil
	case "generate", "menu", "run", "doctor":
		c, e := Load(o.root)
		if e != nil {
			return e
		}
		switch args[0] {
		case "generate":
			return Generate(o.root, c, o.dry, out)
		case "menu":
			return menu(o.root, c, in, out)
		case "doctor":
			return doctor(o.root, c, out)
		case "run":
			if len(o.pos) != 1 {
				return fmt.Errorf("uso: mudarro run serviço:ação")
			}
			s, a, e := findAction(c, o.pos[0])
			if e != nil {
				return e
			}
			return (Runner{In: in, Out: out, Dry: o.dry, Name: o.name}).Run(o.root, s, a)
		}
	}
	return fmt.Errorf("comando desconhecido: %s", args[0])
}
func printReport(w io.Writer, r Report) {
	fmt.Fprintln(w, "Projeto:", r.Config.Name)
	for _, s := range r.Config.Services {
		fmt.Fprintf(w, "- %s (%s): %s / %s / %s\n", s.ID, s.Dir, s.Language, s.Manager, s.Infrastructure.Kind)
		for _, p := range s.Pending {
			fmt.Fprintln(w, "  pendência:", p)
		}
	}
	for _, su := range r.Suggestions {
		fmt.Fprintf(w, "  sugestão: %s:%s [%s]\n", su.Service, su.Name, su.Evidence)
	}
	for _, warning := range r.Warnings {
		fmt.Fprintln(w, "aviso:", warning)
	}
}
func findAction(c Config, key string) (Service, Action, error) {
	parts := strings.SplitN(key, ":", 2)
	var service, name string
	if len(parts) == 2 {
		service, name = parts[0], parts[1]
	} else {
		if len(c.Services) != 1 {
			return Service{}, Action{}, fmt.Errorf("informe serviço:ação em projetos com vários serviços")
		}
		service, name = c.Services[0].ID, key
	}
	for _, s := range c.Services {
		if s.ID == service {
			for _, a := range Actions(s) {
				if a.Name == name {
					return s, a, nil
				}
			}
		}
	}
	return Service{}, Action{}, fmt.Errorf("ação inexistente: %s", key)
}
func choose(r *bufio.Reader, w io.Writer, title string, labels []string) (int, error) {
	fmt.Fprintln(w, "\n"+title)
	for i, l := range labels {
		fmt.Fprintf(w, "%d) %s\n", i+1, l)
	}
	fmt.Fprintln(w, "0) Voltar / sair")
	for {
		fmt.Fprint(w, "> ")
		line, e := r.ReadString('\n')
		if e != nil {
			return 0, e
		}
		n, e := strconv.Atoi(strings.TrimSpace(line))
		if e == nil && n >= 0 && n <= len(labels) {
			return n, nil
		}
		fmt.Fprintln(w, "Opção inválida.")
	}
}
func menu(root string, c Config, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	fmt.Fprintln(out, ascii(c.Name))
	for {
		labels := []string{}
		for _, s := range c.Services {
			labels = append(labels, s.ID+" — "+s.Dir)
		}
		n, e := choose(reader, out, "Serviços", labels)
		if e == io.EOF || n == 0 {
			return nil
		}
		if e != nil {
			return e
		}
		s := c.Services[n-1]
		for {
			actions := Actions(s)
			groups := []string{}
			byGroup := map[string][]Action{}
			for _, a := range actions {
				if _, ok := byGroup[a.Group]; !ok {
					groups = append(groups, a.Group)
				}
				byGroup[a.Group] = append(byGroup[a.Group], a)
			}
			sort.Strings(groups)
			n, e = choose(reader, out, s.ID, groups)
			if e == io.EOF {
				return nil
			}
			if e != nil {
				return e
			}
			if n == 0 {
				break
			}
			group := groups[n-1]
			for {
				list := byGroup[group]
				labels = nil
				for _, a := range list {
					label := a.Name
					if a.Blocked != "" {
						label += " [pendente: " + a.Blocked + "]"
					}
					labels = append(labels, label)
				}
				n, e = choose(reader, out, group, labels)
				if e == io.EOF {
					return nil
				}
				if e != nil {
					return e
				}
				if n == 0 {
					break
				}
				a := list[n-1]
				name := ""
				if strings.Contains(strings.Join(a.Command.Args, " "), "{name}") {
					fmt.Fprint(out, "Nome da migration: ")
					line, e := reader.ReadString('\n')
					if e != nil {
						return nil
					}
					name = strings.TrimSpace(line)
				}
				if e = (Runner{In: reader, Out: out, Name: name}).Run(root, s, a); e != nil {
					fmt.Fprintln(out, "Erro:", e)
				}
			}
		}
	}
}
func doctor(root string, c Config, w io.Writer) error {
	failures := 0
	for _, s := range c.Services {
		dir, _ := serviceDirectory(root, s)
		seen := map[string]bool{}
		for _, a := range Actions(s) {
			if a.Blocked != "" {
				fmt.Fprintf(w, "PENDENTE %s:%s — %s\n", s.ID, a.Name, a.Blocked)
				failures++
				continue
			}
			tools := append([]string{}, a.Command.Requires...)
			if a.Command.Shell != "" {
				tools = append(tools, "bash")
			}
			if len(a.Command.Args) > 0 {
				tools = append(tools, a.Command.Args[0])
			}
			if strings.HasPrefix(a.Special, "local-") {
				tools = append(tools, "ps")
			}
			if strings.HasPrefix(a.Special, "k8s-") {
				tools = append(tools, "kubectl")
			}
			for _, tool := range tools {
				if seen[tool] {
					continue
				}
				seen[tool] = true
				var e error
				if strings.Contains(tool, "/") {
					var p string
					p, e = safePath(root, filepath.Join(s.Dir, tool))
					if e == nil {
						var st os.FileInfo
						st, e = os.Stat(p)
						if e == nil && (st.IsDir() || st.Mode()&0111 == 0) {
							e = fmt.Errorf("não executável")
						}
					}
				} else {
					_, e = exec.LookPath(tool)
				}
				if e != nil {
					fmt.Fprintf(w, "AUSENTE %s: %s — configure ou instale explicitamente\n", s.ID, tool)
					failures++
				} else {
					fmt.Fprintf(w, "OK %s: %s\n", s.ID, tool)
				}
			}
		}
		if s.Infrastructure.Mode == "compose" && (s.Infrastructure.Kind == "docker" || s.Infrastructure.Kind == "podman") {
			provider := s.Infrastructure.Kind
			args := []string{"compose", "version"}
			if provider == "podman" {
				provider = "podman-compose"
				args = []string{"--version"}
			}
			check := exec.Command(provider, args...)
			check.Dir = dir
			if e := check.Run(); e != nil {
				fmt.Fprintf(w, "AUSENTE %s: provedor Compose\n", s.ID)
				failures++
			}
		}
	}
	if failures > 0 {
		return fmt.Errorf("%d pendência(s); as ações independentes continuam disponíveis", failures)
	}
	fmt.Fprintln(w, "Nenhuma dependência executável ausente. Isso não valida conexão, credenciais ou saúde da aplicação.")
	return nil
}
