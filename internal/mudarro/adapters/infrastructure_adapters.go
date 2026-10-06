package adapters

import "github.com/viralabs-dev/mudarro/internal/mudarro/model"

import "fmt"

func (Local) Actions(s model.Service) []model.Action {
	a := []model.Action{}
	for _, n := range []string{"up", "down", "restart", "status", "logs"} {
		v := model.Action{Name: n, Group: "infraestrutura", Special: "local-" + n}
		if (n == "up" || n == "restart") && len(s.Commands["start"].Args) == 0 && s.Commands["start"].Shell == "" {
			v.Blocked = "declare commands.start"
		}
		a = append(a, v)
	}
	return a
}
func (Container) Actions(s model.Service) []model.Action {
	i := s.Infrastructure
	engine := i.Kind
	var a []model.Action
	if i.Mode == "compose" {
		if i.File == "" {
			return []model.Action{blocked("up", "infraestrutura", "declare infrastructure.file")}
		}
		prefix := []string{engine, "compose", "-f", i.File}
		if engine == "podman" {
			// Podman's bundled `compose` delegates to any provider it finds. Using
			// podman-compose directly avoids silently selecting Docker Compose.
			prefix = []string{"podman-compose", "-f", i.File}
		}
		for _, v := range []struct {
			n    string
			args []string
		}{{"up", []string{"up", "-d"}}, {"down", []string{"down"}}, {"restart", []string{"restart"}}, {"status", []string{"ps"}}, {"logs", []string{"logs", "--tail", "100"}}, {"image-build", []string{"build"}}} {
			args := append(append([]string{}, prefix...), v.args...)
			a = append(a, action(v.n, "infraestrutura", args...))
		}
	} else if i.Mode == "dockerfile" {
		if i.Image == "" {
			return []model.Action{blocked("up", "infraestrutura", "declare infrastructure.image para Dockerfile")}
		}
		file := i.File
		if file == "" {
			file = "Dockerfile"
		}
		a = append(a, action("image-build", "infraestrutura", engine, "build", "-f", file, "-t", i.Image, "."))
		args := []string{engine, "run", "-d", "--name", "mudarro-{project}-" + s.ID, "--label", "dev.viralabs.mudarro={project}/" + s.ID}
		if i.Port > 0 {
			args = append(args, "-p", fmt.Sprintf("127.0.0.1:%d:%d", i.Port, i.Port))
		}
		args = append(args, i.Image)
		a = append(a, action("up", "infraestrutura", args...))
		for _, v := range []struct{ n, c string }{{"down", "stop"}, {"restart", "restart"}, {"status", "inspect"}, {"logs", "logs"}} {
			a = append(a, action(v.n, "infraestrutura", engine, v.c, "mudarro-{project}-"+s.ID))
		}
	} else {
		a = append(a, blocked("up", "infraestrutura", "declare infrastructure.mode: compose ou dockerfile"))
	}
	return a
}
func (Kubernetes) Actions(s model.Service) []model.Action {
	i := s.Infrastructure
	if i.Context == "" || i.Namespace == "" || i.File == "" {
		return []model.Action{blocked("up", "infraestrutura", "declare context, namespace e file Kubernetes")}
	}
	base := []string{"kubectl", "--context", i.Context, "--namespace", i.Namespace}
	var a []model.Action
	if i.Mode == "helm" {
		a = append(a, action("up", "infraestrutura", "helm", "upgrade", "--install", s.ID, i.File, "--kube-context", i.Context, "--namespace", i.Namespace))
		// Helm uninstall can delete PVCs; stopping is deliberately a workload scale operation.
		a = append(a, action("down", "infraestrutura", append(base, "scale", "deployment,statefulset", "-l", "app.kubernetes.io/instance="+s.ID, "--replicas=0")...))
		a = append(a, action("status", "infraestrutura", "helm", "status", s.ID, "--kube-context", i.Context, "--namespace", i.Namespace))
	} else if i.Mode == "manifests" || i.Mode == "kustomize" {
		flag := "-f"
		if i.Mode == "kustomize" {
			flag = "-k"
		}
		a = append(a, action("up", "infraestrutura", append(base, "apply", flag, i.File)...))
		// Only scale named workloads selected from the declared resources. Never delete PVCs.
		a = append(a, model.Action{Name: "down", Group: "infraestrutura", Special: "k8s-down"})
		a = append(a, action("status", "infraestrutura", append(base, "get", flag, i.File)...))
	} else {
		return []model.Action{blocked("up", "infraestrutura", "mode Kubernetes deve ser manifests, kustomize ou helm")}
	}
	a = append(a, model.Action{Name: "restart", Group: "infraestrutura", Special: "k8s-restart"}, model.Action{Name: "logs", Group: "infraestrutura", Special: "k8s-logs"})
	return a
}

type Local struct{}
type Container struct{}
type Kubernetes struct{}
