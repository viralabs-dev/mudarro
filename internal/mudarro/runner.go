package mudarro

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Runner struct {
	Executor Executor
	In       io.Reader
	Out      io.Writer
	Dry      bool
	Name     string
}

func (r Runner) Run(root string, s Service, a Action) error {
	if a.Blocked != "" {
		return fmt.Errorf("%s: %s", a.Name, a.Blocked)
	}
	if a.Command.Destructive && !r.Dry {
		expected := "APAGAR " + s.ID
		fmt.Fprintf(r.Out, "Operação destrutiva %s. Digite %q: ", a.Name, expected)
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
		path := s.Database.Path
		if path == "" {
			path = "app.db"
		}
		full, err := safePath(root, filepath.Join(s.Dir, path))
		if err != nil {
			return err
		}
		if r.Dry {
			fmt.Fprintln(r.Out, "criar SQLite se ausente:", full)
			return nil
		}
		if err = os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return err
		}
		f, err := os.OpenFile(full, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if os.IsExist(err) {
			fmt.Fprintln(r.Out, "SQLite existente preservado")
			return nil
		}
		if err != nil {
			return err
		}
		return f.Close()
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
	if r.Dry {
		fmt.Fprintf(r.Out, "%s: %s\n", s.ID, jsonString(args))
		return nil
	}
	// Dockerfile up is repeatable: resume a stopped owned container.
	if (a.Name == "up" || a.Name == "down" || a.Name == "restart") && s.Infrastructure.Mode == "dockerfile" && (s.Infrastructure.Kind == "docker" || s.Infrastructure.Kind == "podman") {
		engine := s.Infrastructure.Kind
		n := "mudarro-" + project + "-" + s.ID
		check := exec.Command(engine, "inspect", "--format", `{{index .Config.Labels "dev.viralabs.mudarro"}}`, n)
		check.Dir = dir
		if out, e := check.Output(); e == nil {
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

type processState struct {
	PID     int    `json:"pid"`
	Token   string `json:"token"`
	Root    string `json:"root"`
	Service string `json:"service"`
}

func running(st processState) bool {
	if st.PID < 2 || len(st.Token) != 32 {
		return false
	}
	b, e := exec.Command("ps", "-ww", "-p", strconv.Itoa(st.PID), "-o", "args=").Output()
	if e != nil {
		return false
	}
	fields := strings.Fields(string(b))
	for i, f := range fields {
		if f == "__supervise" && len(fields) > i+1 && fields[len(fields)-1] == st.Token {
			return true
		}
	}
	return false
}
func readState(path string) processState {
	var st processState
	b, _ := os.ReadFile(path)
	_ = json.Unmarshal(b, &st)
	return st
}
func (r Runner) local(root string, s Service, op string) error {
	base, e := safePath(root, filepath.Join(".mudarro", "run", s.ID))
	if e != nil {
		return e
	}
	if r.Dry {
		fmt.Fprintf(r.Out, "%s: local %s\n", s.ID, op)
		return nil
	}
	if e = os.MkdirAll(base, 0700); e != nil {
		return e
	}
	lockPath, e := safePath(root, filepath.Join(".mudarro", "run", s.ID, "lock"))
	if e != nil {
		return e
	}
	lock, e := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return e
	}
	defer lock.Close()
	if e = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); e != nil {
		return e
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	for _, name := range []string{"state.json", "output.log"} {
		if _, e = safePath(root, filepath.Join(".mudarro", "run", s.ID, name)); e != nil {
			return e
		}
	}
	path := filepath.Join(base, "state.json")
	st := readState(path)
	alive := running(st) && st.Root == root && st.Service == s.ID
	switch op {
	case "status":
		if alive {
			fmt.Fprintf(r.Out, "%s: em execução (PID %d)\n", s.ID, st.PID)
		} else {
			fmt.Fprintf(r.Out, "%s: parado\n", s.ID)
		}
		return nil
	case "logs":
		f, e := os.Open(filepath.Join(base, "output.log"))
		if e != nil {
			return e
		}
		defer f.Close()
		_, e = io.Copy(r.Out, f)
		return e
	case "down", "restart":
		if alive {
			if e = syscall.Kill(-st.PID, syscall.SIGTERM); e != nil {
				return e
			}
			for j := 0; j < 100 && running(st); j++ {
				time.Sleep(50 * time.Millisecond)
			}
			if running(st) {
				return fmt.Errorf("processo não encerrou; verifique os logs antes de tentar novamente")
			}
		}
		if op == "down" {
			fmt.Fprintln(r.Out, s.ID+": parado")
			return nil
		}
	case "up":
		if alive {
			fmt.Fprintln(r.Out, s.ID+": já está em execução")
			return nil
		}
	default:
		return fmt.Errorf("operação local desconhecida: %s", op)
	}
	executable, e := os.Executable()
	if e != nil {
		return e
	}
	token := make([]byte, 16)
	if _, e = rand.Read(token); e != nil {
		return e
	}
	t := hex.EncodeToString(token)
	log, e := os.OpenFile(filepath.Join(base, "output.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	defer log.Close()
	child := exec.Command(executable, "__supervise", root, s.ID, t)
	child.Stdout = log
	child.Stderr = log
	child.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if e = child.Start(); e != nil {
		return e
	}
	st = processState{child.Process.Pid, t, root, s.ID}
	b, _ := json.Marshal(st)
	if e = atomicWrite(path, b, 0600); e != nil {
		_ = syscall.Kill(-st.PID, syscall.SIGTERM)
		return e
	}
	_ = child.Process.Release()
	time.Sleep(150 * time.Millisecond)
	if !running(st) {
		return fmt.Errorf("aplicação encerrou durante a subida; consulte logs")
	}
	fmt.Fprintf(r.Out, "%s: iniciado (PID %d)\n", s.ID, st.PID)
	return nil
}
func supervise(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("supervisor inválido")
	}
	root, id := args[0], args[1]
	c, e := Load(root)
	if e != nil {
		return e
	}
	var s Service
	found := false
	for _, v := range c.Services {
		if v.ID == id {
			s = v
			found = true
		}
	}
	if !found {
		return fmt.Errorf("serviço inexistente")
	}
	a, ok := s.Commands["start"]
	if !ok {
		return fmt.Errorf("commands.start ausente")
	}
	v, e := expandArgs(a.Args, "")
	if e != nil {
		return e
	}
	if a.Shell != "" {
		v = []string{"bash", "-c", a.Shell}
	}
	if len(v) == 0 {
		return fmt.Errorf("comando vazio")
	}
	dir, e := serviceDirectory(root, s)
	if e != nil {
		return e
	}
	child := exec.Command(v[0], v[1:]...)
	child.Dir = dir
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sig)
	if e = child.Start(); e != nil {
		return e
	}
	done := make(chan error, 1)
	go func() { done <- child.Wait() }()
	select {
	case e = <-done:
		return e
	case <-sig:
		_ = child.Process.Signal(syscall.SIGTERM)
		select {
		case e = <-done:
			return e
		case <-time.After(4 * time.Second):
			_ = child.Process.Kill()
			return <-done
		}
	}
}
func (r Runner) kubernetes(root string, s Service, op string) error {
	i := s.Infrastructure
	base := []string{"kubectl", "--context", i.Context, "--namespace", i.Namespace}
	dir, e := serviceDirectory(root, s)
	if e != nil {
		return e
	}
	flag := "-f"
	if i.Mode == "kustomize" {
		flag = "-k"
	}
	query := append(append([]string{}, base...), "get", flag, i.File, "-o", "json")
	if i.Mode == "helm" {
		query = append(append([]string{}, base...), "get", "deployments,statefulsets,daemonsets", "-l", "app.kubernetes.io/instance="+s.ID, "-o", "json")
	}
	if r.Dry {
		fmt.Fprintf(r.Out, "%s: selecionar workloads via %s e executar %s; preservar PVCs\n", s.ID, jsonString(query), op)
		return nil
	}
	b, e := r.executor().Output(dir, query)
	if e != nil {
		return e
	}
	var result struct {
		Kind     string
		Metadata struct{ Name string }
		Items    []struct {
			Kind     string
			Metadata struct{ Name string }
		}
	}
	if e = json.Unmarshal(b, &result); e != nil {
		return e
	}
	if result.Kind != "List" && result.Kind != "" {
		result.Items = append(result.Items, struct {
			Kind     string
			Metadata struct{ Name string }
		}{result.Kind, result.Metadata})
	}
	count := 0
	for _, item := range result.Items {
		kind := strings.ToLower(item.Kind)
		if kind != "deployment" && kind != "statefulset" && kind != "daemonset" {
			continue
		}
		resource := kind + "/" + item.Metadata.Name
		var args []string
		switch op {
		case "down":
			if kind == "daemonset" {
				return fmt.Errorf("DaemonSet não escala a zero; declare ação down personalizada")
			}
			args = append(append([]string{}, base...), "scale", resource, "--replicas=0")
		case "restart":
			args = append(append([]string{}, base...), "rollout", "restart", resource)
		case "logs":
			args = append(append([]string{}, base...), "logs", resource, "--all-containers=true", "--tail=100")
		}
		if e = r.execute(dir, args); e != nil {
			return e
		}
		count++
	}
	if count == 0 {
		return fmt.Errorf("nenhum workload encontrado nos recursos declarados")
	}
	return nil
}
