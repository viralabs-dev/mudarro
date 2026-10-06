package mudarro

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
