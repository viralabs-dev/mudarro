package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func matrixAction(t *testing.T, s Service, name string) Action {
	t.Helper()
	for _, a := range Actions(s) {
		if a.Name == name {
			return a
		}
	}
	t.Fatalf("missing %s", name)
	return Action{}
}
func TestInfrastructureComposeProviders(t *testing.T) {
	for _, engine := range []string{"docker", "podman"} {
		t.Run(engine, func(t *testing.T) {
			s := sample().Services[0]
			s.Infrastructure = Infrastructure{Kind: engine, Mode: "compose", File: "compose with space.yaml"}
			prefix := []string{"docker", "compose", "-f", s.Infrastructure.File}
			if engine == "podman" {
				prefix = []string{"podman-compose", "-f", s.Infrastructure.File}
			}
			for name, tail := range map[string][]string{"up": {"up", "-d"}, "down": {"down"}, "restart": {"restart"}, "status": {"ps"}, "logs": {"logs", "--tail", "100"}, "image-build": {"build"}} {
				a := matrixAction(t, s, name)
				want := append(append([]string{}, prefix...), tail...)
				if a.Blocked != "" || !reflect.DeepEqual(a.Command.Args, want) {
					t.Fatalf("%s: %+v", name, a)
				}
			}
			s.Infrastructure.File = ""
			if matrixAction(t, s, "up").Blocked == "" {
				t.Fatal("missing compose file accepted")
			}
		})
	}
}
func TestInfrastructureDockerfileContracts(t *testing.T) {
	for _, engine := range []string{"docker", "podman"} {
		t.Run(engine, func(t *testing.T) {
			s := sample().Services[0]
			s.Infrastructure = Infrastructure{Kind: engine, Mode: "dockerfile"}
			if matrixAction(t, s, "up").Blocked == "" {
				t.Fatal("image inferred")
			}
			s.Infrastructure.Image = "fixture:local"
			s.Infrastructure.Port = 8080
			up := matrixAction(t, s, "up").Command.Args
			want := []string{engine, "run", "-d", "--name", "mudarro-{project}-app", "--label", "dev.viralabs.mudarro={project}/app", "-p", "127.0.0.1:8080:8080", "fixture:local"}
			if !reflect.DeepEqual(up, want) {
				t.Fatal(up)
			}
			if !reflect.DeepEqual(matrixAction(t, s, "image-build").Command.Args, []string{engine, "build", "-f", "Dockerfile", "-t", "fixture:local", "."}) {
				t.Fatal("default build contract")
			}
			s.Infrastructure.Port = 0
			s.Infrastructure.File = "Dockerfile.custom"
			if strings.Contains(strings.Join(matrixAction(t, s, "up").Command.Args, " "), "-p") {
				t.Fatal("port published implicitly")
			}
			if matrixAction(t, s, "image-build").Command.Args[3] != "Dockerfile.custom" {
				t.Fatal("custom file ignored")
			}
			if !reflect.DeepEqual(matrixAction(t, s, "down").Command.Args, []string{engine, "stop", "mudarro-{project}-app"}) {
				t.Fatal("down removes container")
			}
		})
	}
}
func TestInfrastructureKubernetesCompleteness(t *testing.T) {
	for _, mode := range []string{"manifests", "kustomize", "helm"} {
		t.Run(mode, func(t *testing.T) {
			s := sample().Services[0]
			ready := Infrastructure{Kind: "kubernetes", Mode: mode, Context: "isolated", Namespace: "fixture", File: "resources"}
			for _, missing := range []string{"context", "namespace", "file"} {
				s.Infrastructure = ready
				switch missing {
				case "context":
					s.Infrastructure.Context = ""
				case "namespace":
					s.Infrastructure.Namespace = ""
				case "file":
					s.Infrastructure.File = ""
				}
				if matrixAction(t, s, "up").Blocked == "" {
					t.Fatalf("missing %s accepted", missing)
				}
			}
			s.Infrastructure = ready
			up := matrixAction(t, s, "up")
			if up.Blocked != "" {
				t.Fatal(up)
			}
			if mode == "helm" {
				if !reflect.DeepEqual(up.Command.Args, []string{"helm", "upgrade", "--install", "app", "resources", "--kube-context", "isolated", "--namespace", "fixture"}) {
					t.Fatal(up)
				}
			} else {
				flag := "-f"
				if mode == "kustomize" {
					flag = "-k"
				}
				if !reflect.DeepEqual(up.Command.Args, []string{"kubectl", "--context", "isolated", "--namespace", "fixture", "apply", flag, "resources"}) {
					t.Fatal(up)
				}
			}
			down := matrixAction(t, s, "down")
			joined := strings.Join(down.Command.Args, " ")
			if strings.Contains(joined, "delete") || strings.Contains(joined, "uninstall") {
				t.Fatal("destructive down", down)
			}
			if mode == "helm" && !strings.Contains(joined, "scale deployment,statefulset") {
				t.Fatal(down)
			}
		})
	}
}
func TestInfrastructureKubernetesProbeErrorsAndKinds(t *testing.T) {
	for _, tc := range []struct {
		name, body, special string
		wantError           bool
		calls               int
	}{
		{"stateful-list", `{"kind":"List","items":[{"kind":"StatefulSet","metadata":{"name":"db"}},{"kind":"PersistentVolumeClaim","metadata":{"name":"data"}}]}`, "k8s-down", false, 2},
		{"single-deployment", `{"kind":"Deployment","metadata":{"name":"app"}}`, "k8s-restart", false, 2},
		{"daemonset-logs", `{"kind":"List","items":[{"kind":"DaemonSet","metadata":{"name":"agent"}}]}`, "k8s-logs", false, 2},
		{"only-pvc", `{"kind":"List","items":[{"kind":"PersistentVolumeClaim","metadata":{"name":"data"}}]}`, "k8s-down", true, 1},
		{"daemonset-down", `{"kind":"DaemonSet","metadata":{"name":"agent"}}`, "k8s-down", true, 1},
		{"invalid-json", `invalid`, "k8s-down", true, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			executor := &recordingExecutor{response: []byte(tc.body)}
			s := sample().Services[0]
			s.Infrastructure = Infrastructure{Kind: "kubernetes", Mode: "kustomize", File: "k8s", Context: "test", Namespace: "sandbox"}
			var out bytes.Buffer
			err := (Runner{Out: &out, Executor: executor}).Run(t.TempDir(), s, Action{Name: strings.TrimPrefix(tc.special, "k8s-"), Special: tc.special})
			if (err != nil) != tc.wantError || len(executor.calls) != tc.calls {
				t.Fatalf("err=%v calls=%v", err, executor.calls)
			}
			if !strings.Contains(strings.Join(executor.calls[0], " "), "get -k k8s -o json") {
				t.Fatal(executor.calls)
			}
			if len(executor.calls) > 1 {
				joined := strings.Join(executor.calls[1], " ")
				if strings.Contains(joined, "data") || strings.Contains(joined, "delete") {
					t.Fatal(joined)
				}
				if tc.name == "stateful-list" && !strings.Contains(joined, "scale statefulset/db --replicas=0") {
					t.Fatal(joined)
				}
				if tc.name == "single-deployment" && !strings.Contains(joined, "rollout restart deployment/app") {
					t.Fatal(joined)
				}
				if tc.name == "daemonset-logs" && !strings.Contains(joined, "logs daemonset/agent --all-containers=true --tail=100") {
					t.Fatal(joined)
				}
			}
		})
	}
}
func TestInfrastructureSQLiteInitPreservesAndRejectsSymlink(t *testing.T) {
	root := fixture(t, map[string]string{"data/app.db": "existing database bytes"})
	s := sample().Services[0]
	s.Database = Database{Kind: "sqlite", Path: "data/app.db"}
	var out bytes.Buffer
	a := Action{Name: "db-init", Special: "sqlite-init"}
	if err := (Runner{Out: &out}).Run(root, s, a); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "data/app.db"))
	if err != nil || string(b) != "existing database bytes" {
		t.Fatalf("existing changed %q %v", b, err)
	}
	s.Database.Path = "absent/deep/app.db"
	if err := (Runner{Out: &out, Dry: true}).Run(root, s, a); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "absent")); !os.IsNotExist(err) {
		t.Fatal("dry wrote directory", err)
	}
	outside := fixture(t, map[string]string{"secret.db": "outside untouched"})
	if err := os.Symlink(filepath.Join(outside, "secret.db"), filepath.Join(root, "linked.db")); err != nil {
		t.Fatal(err)
	}
	s.Database.Path = "linked.db"
	if err := (Runner{Out: &out}).Run(root, s, a); err == nil {
		t.Fatal("symlink accepted")
	}
	b, err = os.ReadFile(filepath.Join(outside, "secret.db"))
	if err != nil || string(b) != "outside untouched" {
		t.Fatal("outside changed")
	}
}

func TestInfrastructureHelmProbeSelectsReleaseOnly(t *testing.T) {
	executor := &recordingExecutor{response: []byte(`{"kind":"List","items":[{"kind":"Deployment","metadata":{"name":"web"}},{"kind":"StatefulSet","metadata":{"name":"db"}},{"kind":"PersistentVolumeClaim","metadata":{"name":"data"}}]}`)}
	s := sample().Services[0]
	s.Infrastructure = Infrastructure{Kind: "kubernetes", Mode: "helm", File: "chart", Context: "isolated", Namespace: "fixture"}
	var out bytes.Buffer
	if err := (Runner{Out: &out, Executor: executor}).Run(t.TempDir(), s, Action{Name: "restart", Special: "k8s-restart"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"kubectl", "--context", "isolated", "--namespace", "fixture", "get", "deployments,statefulsets,daemonsets", "-l", "app.kubernetes.io/instance=app", "-o", "json"}
	if len(executor.calls) != 3 || !reflect.DeepEqual(executor.calls[0], want) {
		t.Fatal(executor.calls)
	}
	if !strings.Contains(strings.Join(executor.calls[1], " "), "rollout restart deployment/web") || !strings.Contains(strings.Join(executor.calls[2], " "), "rollout restart statefulset/db") {
		t.Fatal(executor.calls)
	}
}
