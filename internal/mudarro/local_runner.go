package mudarro

import (
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
			fmt.Fprintf(r.Out, r.text("%s: running (PID %d)\n", "%s: em execução (PID %d)\n"), s.ID, st.PID)
		} else {
			fmt.Fprintf(r.Out, r.text("%s: stopped\n", "%s: parado\n"), s.ID)
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
			fmt.Fprintln(r.Out, s.ID+r.text(": stopped", ": parado"))
			return nil
		}
	case "up":
		if alive {
			fmt.Fprintln(r.Out, s.ID+r.text(": already running", ": já está em execução"))
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
	supervisorArgs := []string{"__supervise", root, s.ID}
	if r.ConfigPath != "" {
		supervisorArgs = append(supervisorArgs, r.ConfigPath)
	}
	supervisorArgs = append(supervisorArgs, t)
	child := exec.Command(executable, supervisorArgs...)
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
	fmt.Fprintf(r.Out, r.text("%s: started (PID %d)\n", "%s: iniciado (PID %d)\n"), s.ID, st.PID)
	return nil
}
func supervise(args []string) error {
	if len(args) != 3 && len(args) != 4 {
		return fmt.Errorf("supervisor inválido")
	}
	root, id := args[0], args[1]
	var c Config
	var e error
	if len(args) == 4 {
		c, e = LoadFile(root, args[2])
	} else {
		c, e = Load(root)
	}
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
