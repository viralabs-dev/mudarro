package mudarro

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/viralabs-dev/mudarro/internal/mudarro/projectfs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// scanFrameworkMetadata annotates declared dependencies only. It never imports
// modules, resolves packages, invents an entrypoint, or changes commands.
func scanFrameworkMetadata(root string, r *Report) {
	for i := range r.Config.Services {
		s := &r.Config.Services[i]
		var files []string
		switch s.Language {
		case "javascript", "typescript":
			files = []string{"package.json"}
		case "python":
			files = []string{"pyproject.toml", "requirements.txt"}
		default:
			continue
		}
		found := map[string]bool{}
		for _, file := range files {
			rel := filepath.ToSlash(filepath.Join(s.Dir, file))
			if frameworkExcluded(rel, r.Config.Exclude) {
				continue
			}
			b, e := projectfs.ReadManifest(root, s.Dir, file)
			if os.IsNotExist(e) {
				continue
			}
			if e != nil {
				r.Warnings = append(r.Warnings, e.Error())
				continue
			}
			deps, e := frameworkDependencies(file, b)
			if e != nil {
				r.Warnings = append(r.Warnings, fmt.Sprintf("%s: framework metadata: %v", rel, e))
				continue
			}
			for _, dep := range deps {
				if s.Language == "python" {
					dep = strings.ToLower(dep)
				}
				accepted := false
				if s.Language == "python" {
					switch dep {
					case "flask", "fastapi", "django", "pytest":
						accepted = true
					}
				} else {
					switch dep {
					case "next", "@nestjs/core", "vite", "express":
						accepted = true
					}
				}
				if !accepted {
					continue
				}
				r.Evidence = append(r.Evidence, Evidence{rel, "declared-dependency:" + dep})
				if dep != "pytest" {
					found[dep] = true
				}
			}
		}
		names := make([]string, 0, len(found))
		for name := range found {
			names = append(names, name)
		}
		sort.Strings(names)
		if len(names) == 1 {
			s.Framework = names[0]
			s.Pending = append(s.Pending, "Framework declared in metadata only; installation and entrypoint are not validated")
		}
		if len(names) > 1 {
			s.Framework = ""
			s.Pending = append(s.Pending, "Multiple declared frameworks: "+strings.Join(names, ", ")+"; declare framework and commands explicitly")
		}
	}
}
func frameworkExcluded(rel string, excludes []string) bool {
	for _, pat := range excludes {
		match, _ := path.Match(pat, rel)
		if match || rel == pat || strings.HasPrefix(rel, strings.TrimSuffix(pat, "/")+"/") {
			return true
		}
	}
	return false
}

var frameworkRequirement = regexp.MustCompile(`^([A-Za-z0-9][A-Za-z0-9._-]*)(?:\[[A-Za-z0-9_, .-]+\])?(?:\s*(?:[<>=!~;@]|$))`)

func frameworkDependencies(file string, b []byte) ([]string, error) {
	var deps []string
	switch file {
	case "package.json":
		var p struct {
			Dependencies         map[string]string `json:"dependencies"`
			DevDependencies      map[string]string `json:"devDependencies"`
			OptionalDependencies map[string]string `json:"optionalDependencies"`
		}
		if e := json.Unmarshal(b, &p); e != nil {
			return nil, e
		}
		for _, m := range []map[string]string{p.Dependencies, p.DevDependencies, p.OptionalDependencies} {
			for k, v := range m {
				if strings.TrimSpace(v) != "" {
					deps = append(deps, k)
				}
			}
		}
	case "pyproject.toml":
		var p struct {
			Project struct {
				Dependencies         []string
				OptionalDependencies map[string][]string `toml:"optional-dependencies"`
			}
			Tool struct {
				Poetry struct {
					Dependencies map[string]any
					Group        map[string]struct{ Dependencies map[string]any }
				}
			}
		}
		if e := toml.Unmarshal(b, &p); e != nil {
			return nil, e
		}
		reqs := append([]string{}, p.Project.Dependencies...)
		for _, v := range p.Project.OptionalDependencies {
			reqs = append(reqs, v...)
		}
		for _, v := range reqs {
			if m := frameworkRequirement.FindStringSubmatch(strings.TrimSpace(v)); len(m) > 1 {
				deps = append(deps, strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(m[1]), "_", "-"), ".", "-"))
			}
		}
		for k, v := range p.Tool.Poetry.Dependencies {
			if frameworkPoetryVersion(v) {
				deps = append(deps, k)
			}
		}
		for _, g := range p.Tool.Poetry.Group {
			for k, v := range g.Dependencies {
				if frameworkPoetryVersion(v) {
					deps = append(deps, k)
				}
			}
		}
	case "requirements.txt":
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
				continue
			}
			if m := frameworkRequirement.FindStringSubmatch(line); len(m) > 1 {
				deps = append(deps, strings.ToLower(m[1]))
			}
		}
	}
	sort.Strings(deps)
	out := deps[:0]
	for _, d := range deps {
		if len(out) == 0 || out[len(out)-1] != d {
			out = append(out, d)
		}
	}
	return out, nil
}
func frameworkPoetryVersion(v any) bool {
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x) != ""
	case map[string]any:
		for _, k := range []string{"version", "git", "path", "url"} {
			if s, ok := x[k].(string); ok && strings.TrimSpace(s) != "" {
				return true
			}
		}
	}
	return false
}
