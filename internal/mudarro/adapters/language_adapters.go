package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (adapter JavaScript) Detect(root, dir string, f map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if !f["package.json"] {
		return nil, nil, nil
	}
	b, e := read(root, dir, "package.json")
	if e != nil {
		return nil, nil, e
	}
	var p struct {
		Name            string
		PackageManager  string
		Scripts         map[string]string
		Dependencies    map[string]string
		DevDependencies map[string]string
	}
	if e = json.Unmarshal(b, &p); e != nil {
		return nil, nil, fmt.Errorf("%s/package.json: %w", dir, e)
	}
	s := BaseService(dir, "javascript")
	if f["tsconfig.json"] || p.DevDependencies["typescript"] != "" || p.Dependencies["typescript"] != "" {
		s.Language = "typescript"
	}
	managers := []string{}
	seenManagers := map[string]bool{}
	for _, m := range []struct{ f, m string }{{"package-lock.json", "npm"}, {"pnpm-lock.yaml", "pnpm"}, {"yarn.lock", "yarn"}, {"bun.lock", "bun"}, {"bun.lockb", "bun"}} {
		if f[m.f] && !seenManagers[m.m] {
			managers = append(managers, m.m)
			seenManagers[m.m] = true
		}
	}
	if p.PackageManager != "" {
		manager, err := packageManagerName(p.PackageManager)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", filepath.Join(dir, "package.json"), err)
		}
		s.Manager = manager
	} else if len(managers) == 1 {
		s.Manager = managers[0]
	} else if len(managers) == 0 {
		s.Manager = "npm"
	}
	if len(managers) > 1 && p.PackageManager == "" {
		s.Pending = append(s.Pending, "Múltiplos lockfiles: declare manager")
	}
	if adapter.WorkspaceRoot != "" {
		s.WorkspaceRoot = adapter.WorkspaceRoot
		if adapter.WorkspaceConflict != "" {
			s.Manager = ""
			s.Pending = append(s.Pending, adapter.WorkspaceConflict)
		} else {
			s.Manager = adapter.WorkspaceManager
		}
	}
	for _, k := range []string{"next", "@nestjs/core", "vite", "express"} {
		if p.Dependencies[k] != "" || p.DevDependencies[k] != "" {
			s.Framework = k
			break
		}
	}
	suggestions := []model.Suggestion{}
	keys := []string{}
	for k := range p.Scripts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	names := map[string][]string{}
	for _, k := range keys {
		n := strings.NewReplacer(":", "-", ".", "-").Replace(k)
		names[n] = append(names[n], k)
	}
	reported := map[string]bool{}
	for _, k := range keys {
		name := strings.NewReplacer(":", "-", ".", "-").Replace(k)
		if !identifier.MatchString(name) {
			continue
		}
		if len(names[name]) > 1 {
			if !reported[name] {
				s.Pending = append(s.Pending, fmt.Sprintf("Scripts %s colidem no nome %s: declare comando explícito", strings.Join(names[name], ", "), name))
				reported[name] = true
			}
			continue
		}
		group := "scripts"
		if k == "test" || k == "lint" || strings.HasPrefix(k, "test:") || strings.HasPrefix(k, "lint:") {
			group = "qualidade"
		}
		if k == "dev" || k == "start" || k == "build" {
			group = "aplicacao"
		}
		c := cmd(group, s.Manager, "run", k)
		if adapter.WorkspaceConflict != "" {
			continue
		}
		suggestions = append(suggestions, model.Suggestion{Name: name, Command: c, Evidence: filepath.Join(dir, "package.json")})
		if s.Manager != "" && group != "scripts" {
			s.Commands[name] = c
		}
	}
	if s.Manager != "" && (adapter.WorkspaceRoot == "" || adapter.WorkspaceRoot == filepath.ToSlash(dir)) {
		s.Commands["install"] = cmd("dependencias", s.Manager, "install")
	}
	if _, ok := s.Commands["start"]; !ok {
		if c, ok := s.Commands["dev"]; ok {
			s.Commands["start"] = c
		} else {
			s.Pending = append(s.Pending, "Declare commands.start para execução local")
		}
	}
	return s, suggestions, nil
}
func (Python) Detect(root, dir string, f map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if f["Pipfile"] || f["Pipfile.lock"] {
		return detectPipenv(root, dir, f)
	}
	if !f["pyproject.toml"] && !f["requirements.txt"] && !f["manage.py"] && !f["Pipfile"] {
		return nil, nil, nil
	}
	s := BaseService(dir, "python")
	s.Manager = "pip"
	var suggestions []model.Suggestion
	if f["poetry.lock"] && f["uv.lock"] {
		s.Manager = ""
		s.Pending = append(s.Pending, "Múltiplos lockfiles: declare manager")
	} else if f["uv.lock"] {
		s.Manager = "uv"
	} else if f["poetry.lock"] {
		s.Manager = "poetry"
	}
	if f["pyproject.toml"] {
		b, e := read(root, dir, "pyproject.toml")
		if e != nil {
			return nil, nil, e
		}
		var p struct {
			Project struct {
				Scripts      map[string]string
				Dependencies []string
			}
			Tool struct {
				Poetry struct {
					Scripts      map[string]string
					Dependencies map[string]any
				}
			}
		}
		if e = toml.Unmarshal(b, &p); e != nil {
			return nil, nil, fmt.Errorf("%s/pyproject.toml: %w", dir, e)
		}
		for k := range p.Project.Scripts {
			suggestions = append(suggestions, model.Suggestion{Name: k, Command: cmd("scripts", ".venv/bin/"+k), Evidence: filepath.Join(dir, "pyproject.toml")})
		}
		for k := range p.Tool.Poetry.Scripts {
			suggestions = append(suggestions, model.Suggestion{Name: k, Command: cmd("scripts", "poetry", "run", k), Evidence: filepath.Join(dir, "pyproject.toml")})
		}
	}
	switch s.Manager {
	case "uv":
		s.Commands["install"] = cmd("dependencias", "uv", "sync")
	case "poetry":
		s.Commands["install"] = cmd("dependencias", "poetry", "install")
	case "pip":
		s.Commands["venv"] = cmd("dependencias", "python3", "-m", "venv", ".venv")
		if f["requirements.txt"] {
			s.Commands["install"] = cmd("dependencias", ".venv/bin/python", "-m", "pip", "install", "-r", "requirements.txt")
		} else {
			s.Commands["install"] = cmd("dependencias", ".venv/bin/python", "-m", "pip", "install", "-e", ".")
		}
	}
	if f["manage.py"] {
		s.Framework = "django"
		s.Commands["start"] = cmd("aplicacao", pythonArgs(s.Manager, "manage.py", "runserver")...)
		s.Commands["test"] = cmd("qualidade", pythonArgs(s.Manager, "manage.py", "test")...)
	} else {
		s.Pending = append(s.Pending, "Declare commands.start; módulo Python não é inferido")
	}
	return s, suggestions, nil
}
func pythonArgs(manager string, args ...string) []string {
	prefix := []string{".venv/bin/python"}
	if manager == "uv" {
		prefix = []string{"uv", "run", "python"}
	}
	if manager == "poetry" {
		prefix = []string{"poetry", "run", "python"}
	}
	if manager == "pipenv" {
		prefix = []string{"pipenv", "run", "python"}
	}
	return append(prefix, args...)
}
func (Go) Detect(root, dir string, f map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if !f["go.mod"] {
		return nil, nil, nil
	}
	s := BaseService(dir, "go")
	s.Manager = "go"
	s.Commands["build"] = cmd("aplicacao", "go", "build", "./...")
	s.Commands["test"] = cmd("qualidade", "go", "test", "./...")
	s.Commands["install"] = cmd("dependencias", "go", "mod", "download")
	entries := map[string]bool{}
	base := filepath.Join(root, dir)
	e := filepath.WalkDir(base, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		rel, _ := filepath.Rel(base, path)
		rootRel, _ := filepath.Rel(root, path)
		for _, pat := range excludes {
			match, _ := filepath.Match(pat, filepath.ToSlash(rootRel))
			if match || rootRel == pat || strings.HasPrefix(filepath.ToSlash(rootRel), strings.TrimSuffix(pat, "/")+"/") {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if d.IsDir() {
			if rel != "." && (ignored[d.Name()] || exists(filepath.Join(path, "go.mod"))) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("Go source must be a regular file: %s", rootRel)
		}
		source, err := read(root, dir, rel)
		if err != nil {
			return fmt.Errorf("Go source %s: %w", rootRel, err)
		}
		context := build.Default
		context.OpenFile = func(string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(source)), nil
		}
		active, err := context.MatchFile(filepath.Dir(path), filepath.Base(path))
		if err != nil || !active {
			return nil
		}
		p, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
		if err != nil || p.Name.Name != "main" {
			return nil
		}
		for _, decl := range p.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Name.Name == "main" && fn.Recv == nil && fn.Body != nil && fn.Type.Params.NumFields() == 0 && fn.Type.Results.NumFields() == 0 && fn.Type.TypeParams.NumFields() == 0 {
				entries[filepath.Dir(rel)] = true
			}
		}
		return nil
	})
	if e != nil {
		return nil, nil, e
	}
	targets := make([]string, 0, len(entries))
	for p := range entries {
		targets = append(targets, p)
	}
	sort.Strings(targets)
	suggestions := make([]model.Suggestion, 0, len(targets))
	for i, p := range targets {
		c := cmd("aplicacao", "go", "run", "./"+filepath.ToSlash(p))
		suggestions = append(suggestions, model.Suggestion{Name: fmt.Sprintf("go-entry-%d", i+1), Purpose: "go-start", Command: c, Evidence: filepath.Join(dir, p)})
		if len(targets) == 1 {
			s.Commands["start"] = c
		}
	}
	if len(targets) != 1 {
		s.Pending = append(s.Pending, "Declare commands.start: nenhuma entrada Go única")
	}
	return s, suggestions, nil
}

type JavaScript struct{ WorkspaceManager, WorkspaceRoot, WorkspaceConflict string }
type Python struct{}
type Go struct{}
