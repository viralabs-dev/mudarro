package mudarro

import (
	"encoding/json"
	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"github.com/viralabs-dev/mudarro/internal/mudarro/projectfs"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Evidence struct {
	Path string `json:"path"`
	Kind string `json:"kind"`
}
type Report struct {
	JavaScriptWorkspaces []JavaScriptWorkspace `json:"javascript_workspaces,omitempty"`
	Workspaces           []Workspace           `json:"workspaces,omitempty"`
	Config               Config                `json:"config"`
	Evidence             []Evidence            `json:"evidence"`
	Suggestions          []Suggestion          `json:"suggestions"`
	Warnings             []string              `json:"warnings"`
}
type LanguageAdapter interface {
	Detect(root, dir string, files map[string]bool, excludes []string) (*Service, []Suggestion, error)
}

var languageAdapters = []LanguageAdapter{adapters.JavaScript{}, adapters.Python{}, adapters.Go{}, adapters.Elixir{}, adapters.Rust{}, adapters.Java{}, adapters.Composer{}, adapters.DotNet{}, adapters.Ruby{}}

func read(root, dir, file string) ([]byte, error) { return projectfs.ReadManifest(root, dir, file) }
func cmd(group string, args ...string) Command    { return Command{Args: args, Group: group} }
func exists(p string) bool                        { st, e := os.Lstat(p); return e == nil && st.Mode()&os.ModeSymlink == 0 }

func Scan(root string, excludes []string) (Report, error) {
	r := Report{Config: Config{Version: 1, Name: filepath.Base(root), Exclude: excludes}}
	dirs := map[string]map[string]bool{}
	e := filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		rel, _ := filepath.Rel(root, path)
		for _, pat := range excludes {
			match, _ := filepath.Match(pat, filepath.ToSlash(rel))
			if match || rel == pat || strings.HasPrefix(filepath.ToSlash(rel), strings.TrimSuffix(pat, "/")+"/") {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if d.IsDir() {
			// Mix artifacts are scoped to an adjacent declared project, never
			// globally hidden from unrelated language/workspace scans.
			if d.Name() == "deps" || d.Name() == "_build" {
				manifest := filepath.Join(filepath.Dir(rel), "mix.exs")
				excluded := false
				for _, pat := range excludes {
					match, _ := filepath.Match(pat, filepath.ToSlash(manifest))
					if match || manifest == pat || strings.HasPrefix(filepath.ToSlash(manifest), strings.TrimSuffix(pat, "/")+"/") {
						excluded = true
						break
					}
				}
				if !excluded {
					if st, err := os.Lstat(filepath.Join(root, manifest)); err == nil && st.Mode().IsRegular() {
						return filepath.SkipDir
					}
				}
			}
			if rel != "." && adapters.IgnoredDirectory(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		dir := filepath.Dir(rel)
		if dirs[dir] == nil {
			dirs[dir] = map[string]bool{}
		}
		dirs[dir][d.Name()] = true
		return nil
	})
	if e != nil {
		return r, e
	}
	paths := []string{}
	for p := range dirs {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	jsPlans, jsWarnings := javascriptWorkspacePlans(root, paths, dirs)
	r.Warnings = append(r.Warnings, jsWarnings...)
	jsContexts, jsPlans := javascriptWorkspaceContexts(root, paths, dirs, jsPlans)
	for _, plan := range jsPlans {
		r.JavaScriptWorkspaces = append(r.JavaScriptWorkspaces, plan.JavaScriptWorkspace)
		r.Evidence = append(r.Evidence, Evidence{plan.Path, "javascript-workspace"})
	}
	for _, dir := range paths {
		f := dirs[dir]
		if f["go.work"] {
			workspace := filepath.Join(dir, "go.work")
			r.Evidence = append(r.Evidence, Evidence{workspace, "go-workspace"})
			parsed, warnings := scanWorkspace(root, dir, excludes)
			r.Workspaces = append(r.Workspaces, parsed)
			r.Warnings = append(r.Warnings, warnings...)
		}
		for _, a := range languageAdapters {
			if _, ok := a.(adapters.JavaScript); ok {
				if context, found := jsContexts[dir]; found {
					a = adapters.JavaScript{WorkspaceManager: context.manager, WorkspaceRoot: context.root, WorkspaceConflict: context.conflict}
				}
			}
			s, suggestions, e := a.Detect(root, dir, f, excludes)
			if e != nil {
				r.Warnings = append(r.Warnings, e.Error())
				continue
			}
			if s == nil {
				continue
			}
			collision := false
			for _, prior := range r.Config.Services {
				if prior.ID == s.ID {
					r.Warnings = append(r.Warnings, "service ID collision: "+s.ID+" at "+dir+" conflicts with "+prior.Dir+"; configure services explicitly")
					collision = true
					break
				}
			}
			if collision {
				continue
			}
			detectInfrastructure(root, s, &r)
			detectDatabase(root, s)
			r.Config.Services = append(r.Config.Services, *s)
			r.Evidence = append(r.Evidence, Evidence{dir, s.Language})
			for _, su := range suggestions {
				su.Service = s.ID
				r.Suggestions = append(r.Suggestions, su)
			}
			if dir == "." && f["package.json"] {
				b, _ := read(root, dir, "package.json")
				var p struct{ Name string }
				if json.Unmarshal(b, &p) == nil && p.Name != "" {
					r.Config.Name = p.Name
				}
			}
		}
	}
	scanFrameworkMetadata(root, &r)
	scanNativeSources(root, paths, dirs, &r)
	if len(r.Config.Services) == 0 {
		s := adapters.BaseService(".", "custom")
		detectInfrastructure(root, s, &r)
		s.Pending = append(s.Pending, "Declare comandos da aplicação")
		r.Config.Services = append(r.Config.Services, *s)
	}
	scanOrchestration(root, paths, dirs, &r)
	for _, dir := range paths {
		for file := range dirs[dir] {
			if file != "Makefile" && !strings.HasSuffix(file, ".sh") {
				continue
			}
			if file == "menu.sh" {
				continue
			}
			owner := -1
			best := -1
			for i, s := range r.Config.Services {
				if s.Dir == "." || dir == s.Dir || strings.HasPrefix(dir, s.Dir+string(filepath.Separator)) {
					if len(s.Dir) > best {
						owner = i
						best = len(s.Dir)
					}
				}
			}
			if owner < 0 {
				r.Warnings = append(r.Warnings, "unowned script/Makefile: "+filepath.ToSlash(filepath.Join(dir, file))+"; configure a service explicitly")
				continue
			}
			s := r.Config.Services[owner]
			rel, _ := filepath.Rel(s.Dir, filepath.Join(dir, file))
			if file == "Makefile" {
				b, e := read(root, dir, file)
				if e != nil {
					r.Warnings = append(r.Warnings, e.Error())
					continue
				}
				re := regexp.MustCompile(`(?m)^([A-Za-z0-9][A-Za-z0-9_-]*):[^=]`)
				for _, m := range re.FindAllStringSubmatch(string(b), -1) {
					r.Suggestions = append(r.Suggestions, Suggestion{Service: s.ID, Name: "make-" + m[1], Command: cmd("scripts", "make", "-C", filepath.Dir(rel), m[1]), Evidence: filepath.Join(dir, file)})
				}
			} else {
				n := strings.TrimSuffix(strings.ReplaceAll(filepath.ToSlash(rel), "/", "-"), ".sh")
				if identifier.MatchString(n) {
					r.Suggestions = append(r.Suggestions, Suggestion{Service: s.ID, Name: "script-" + n, Command: cmd("scripts", "bash", rel), Evidence: filepath.Join(dir, file)})
				}
			}
		}
	}
	sort.Slice(r.Suggestions, func(i, j int) bool {
		return r.Suggestions[i].Service+"/"+r.Suggestions[i].Name < r.Suggestions[j].Service+"/"+r.Suggestions[j].Name
	})
	return r, nil
}

func detectInfrastructure(root string, s *Service, r *Report) {
	candidates := []Infrastructure{}
	for _, f := range []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"} {
		if exists(filepath.Join(root, s.Dir, f)) {
			candidates = append(candidates, Infrastructure{Kind: "docker", Mode: "compose", File: f})
			r.Evidence = append(r.Evidence, Evidence{filepath.Join(s.Dir, f), "compose (Docker ou Podman)"})
		}
	}
	if len(candidates) == 0 && exists(filepath.Join(root, s.Dir, "Dockerfile")) {
		candidates = append(candidates, Infrastructure{Kind: "docker", Mode: "dockerfile", File: "Dockerfile"})
	}
	for _, f := range []string{"k8s", "kubernetes"} {
		if exists(filepath.Join(root, s.Dir, f)) {
			candidates = append(candidates, Infrastructure{Kind: "kubernetes", Mode: "manifests", File: f})
		}
	}
	if exists(filepath.Join(root, s.Dir, "Chart.yaml")) {
		candidates = append(candidates, Infrastructure{Kind: "kubernetes", Mode: "helm", File: "."})
	}
	for _, f := range []string{"kustomization.yaml", "kustomization.yml"} {
		if exists(filepath.Join(root, s.Dir, f)) {
			candidates = append(candidates, Infrastructure{Kind: "kubernetes", Mode: "kustomize", File: "."})
			break
		}
	}
	if len(candidates) == 0 {
		s.Infrastructure.Kind = "local"
	} else if len(candidates) == 1 {
		s.Infrastructure = candidates[0]
		if s.Infrastructure.Mode == "compose" {
			s.Infrastructure.Kind = ""
			s.Pending = append(s.Pending, "Compose detectado: declare infrastructure.kind como docker ou podman")
		}
	} else {
		s.Pending = append(s.Pending, "Múltiplas infraestruturas: declare infrastructure.kind, mode e file")
		for _, c := range candidates {
			r.Evidence = append(r.Evidence, Evidence{filepath.Join(s.Dir, c.File), c.Kind + "/" + c.Mode})
		}
	}
	if s.Infrastructure.Kind == "kubernetes" {
		s.Pending = append(s.Pending, "Declare context e namespace Kubernetes")
	}
}
func detectDatabase(root string, s *Service) {
	candidates := []string{}
	if exists(filepath.Join(root, s.Dir, "prisma/schema.prisma")) {
		candidates = append(candidates, "prisma")
		b, _ := read(root, s.Dir, "prisma/schema.prisma")
		if strings.Contains(string(b), `"postgresql"`) {
			s.Database.Kind = "postgresql"
		}
		if strings.Contains(string(b), `"sqlite"`) {
			s.Database.Kind = "sqlite"
		}
	}
	if s.Framework == "django" {
		candidates = append(candidates, "django")
	}
	if exists(filepath.Join(root, s.Dir, "alembic.ini")) {
		candidates = append(candidates, "alembic")
	}
	if paths, e := filepath.Glob(filepath.Join(root, s.Dir, "migrations", "*.sql")); e == nil {
		for _, p := range paths {
			rel, _ := filepath.Rel(filepath.Join(root, s.Dir), p)
			b, e := read(root, s.Dir, rel)
			if e == nil && strings.Contains(string(b), "-- +goose") {
				candidates = append(candidates, "goose")
				break
			}
		}
	}
	if len(candidates) > 1 {
		s.Pending = append(s.Pending, "Múltiplas ferramentas de banco: declare database.tool")
		return
	}
	if len(candidates) == 1 {
		s.Database.Tool = candidates[0]
		s.Database.URLenv = "DATABASE_URL"
		if s.Database.Kind == "" {
			s.Pending = append(s.Pending, "Declare database.kind; conexão não será executada durante detecção")
		}
	}
}
