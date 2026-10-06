package mudarro

import (
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type InfrastructureAdapter interface{ Actions(Service) []Action }
type DatabaseAdapter interface{ Actions(Service) []Action }

var infrastructureAdapters = map[string]InfrastructureAdapter{
	"local": adapters.Local{}, "docker": adapters.Container{}, "podman": adapters.Container{}, "kubernetes": adapters.Kubernetes{},
}
var databaseAdapters = map[string]DatabaseAdapter{
	"prisma": adapters.Prisma{}, "django": adapters.Django{}, "alembic": adapters.Alembic{}, "goose": adapters.Goose{},
}

func blocked(name, group, reason string) Action {
	return Action{Name: name, Group: group, Blocked: reason}
}
func Actions(s Service) []Action {
	var a []Action
	if adapter, ok := infrastructureAdapters[s.Infrastructure.Kind]; ok {
		a = adapter.Actions(s)
	} else if s.Infrastructure.Kind != "custom" {
		a = []Action{blocked("up", "infraestrutura", "declare infrastructure.kind (detecção ambígua)")}
	}
	if s.Database.Tool != "" {
		if s.Database.Kind != "postgresql" && s.Database.Kind != "sqlite" {
			a = append(a, blocked("db-migrate", "banco", "declare database.kind: postgresql ou sqlite"))
		} else if adapter, ok := databaseAdapters[s.Database.Tool]; ok {
			a = append(a, adapter.Actions(s)...)
		} else {
			a = append(a, blocked("db-migrate", "banco", "ferramenta de migration não suportada"))
		}
	}
	a = append(a, adapters.SetupActions(s)...)
	byName := map[string]Action{}
	for _, v := range a {
		byName[v.Name] = v
	}
	for n, c := range s.Commands {
		g := c.Group
		if g == "" {
			g = "scripts"
		}
		if n == "db-reset" {
			c.Destructive = true
		}
		byName[n] = Action{Name: n, Group: g, Command: c}
	}
	a = nil
	for _, v := range byName {
		a = append(a, v)
	}
	sort.Slice(a, func(i, j int) bool {
		if a[i].Group == a[j].Group {
			return a[i].Name < a[j].Name
		}
		return a[i].Group < a[j].Group
	})
	return a
}
func expandArgs(args []string, name string) ([]string, error) {
	out := make([]string, len(args))
	for i, v := range args {
		if strings.Contains(v, "{name}") {
			if !identifier.MatchString(name) {
				return nil, fmt.Errorf("informe --name com letras, números, _ ou -")
			}
			v = strings.ReplaceAll(v, "{name}", name)
		}
		var missing string
		v = os.Expand(v, func(k string) string {
			value, ok := os.LookupEnv(k)
			if !ok || value == "" {
				missing = k
			}
			return value
		})
		if missing != "" {
			return nil, fmt.Errorf("variável de ambiente ausente: %s", missing)
		}
		out[i] = v
	}
	return out, nil
}
func serviceDirectory(root string, s Service) (string, error) {
	return safePath(root, filepath.Clean(s.Dir))
}

func action(name, group string, args ...string) Action {
	return Action{Name: name, Group: group, Command: cmd(group, args...)}
}
