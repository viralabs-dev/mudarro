package mudarro

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"gopkg.in/yaml.v3"
)

type JavaScriptWorkspace struct {
	Path    string   `json:"path"`
	Root    string   `json:"root"`
	Members []string `json:"members"`
	Manager string   `json:"manager,omitempty"`
}
type javascriptWorkspacePlan struct {
	JavaScriptWorkspace
	patterns []string
	version  string
}
type javascriptWorkspaceContext struct{ root, manager, conflict string }

// Membership is matched against the safe scan inventory, never an executing glob tool.
func javascriptWorkspacePlans(root string, paths []string, dirs map[string]map[string]bool) ([]javascriptWorkspacePlan, []string) {
	var plans []javascriptWorkspacePlan
	var warnings []string
	for _, dir := range paths {
		f := dirs[dir]
		if !f["package.json"] {
			continue
		}
		b, err := read(root, dir, "package.json")
		if err != nil {
			continue
		}
		var p struct {
			Workspaces     json.RawMessage
			PackageManager string
		}
		if json.Unmarshal(b, &p) != nil {
			continue
		}
		var patterns []string
		declaration := "package.json"
		declarationConflict := false
		declared := len(p.Workspaces) > 0 && string(p.Workspaces) != "null"
		if declared {
			if err := json.Unmarshal(p.Workspaces, &patterns); err != nil {
				var object struct {
					Packages []string `json:"packages"`
				}
				if json.Unmarshal(p.Workspaces, &object) != nil || object.Packages == nil {
					warnings = append(warnings, fmt.Sprintf("%s/package.json: invalid workspaces; use array or packages array", dir))
					continue
				}
				patterns = object.Packages
			}
		}
		if f["pnpm-workspace.yaml"] {
			b, err := read(root, dir, "pnpm-workspace.yaml")
			if err != nil {
				warnings = append(warnings, err.Error())
				continue
			}
			var document struct {
				Packages []string `yaml:"packages"`
			}
			if err = yaml.Unmarshal(b, &document); err != nil || document.Packages == nil {
				warnings = append(warnings, fmt.Sprintf("%s/pnpm-workspace.yaml: invalid packages array", dir))
				continue
			}
			if declared && !sameWorkspacePatterns(patterns, document.Packages) {
				warnings = append(warnings, fmt.Sprintf("%s: conflicting workspace declarations; configure ownership explicitly", dir))
				patterns = append(patterns, document.Packages...)
				declarationConflict = true
			} else {
				patterns = document.Packages
			}
			declaration = "pnpm-workspace.yaml"
			declared = true
		}
		if !declared {
			continue
		}
		invalid := false
		for _, pattern := range patterns {
			if err := validateWorkspacePattern(pattern); err != nil {
				warnings = append(warnings, fmt.Sprintf("%s/%s: %v", dir, declaration, err))
				invalid = true
			}
		}
		if invalid {
			continue
		}
		service, _, err := (adapters.JavaScript{}).Detect(root, dir, f, nil)
		manager := ""
		if err != nil {
			warnings = append(warnings, err.Error())
		} else if service != nil {
			manager = service.Manager
		}
		if err != nil || service == nil {
			declarationConflict = true
		}
		if declarationConflict {
			manager = ""
		} else if declaration == "pnpm-workspace.yaml" {
			if p.PackageManager == "" && !hasJavaScriptLock(f) {
				manager = "pnpm"
			} else if manager != "pnpm" {
				warnings = append(warnings, fmt.Sprintf("%s: pnpm workspace conflicts with manager %q", dir, manager))
				manager = ""
			}
		}
		plans = append(plans, javascriptWorkspacePlan{JavaScriptWorkspace: JavaScriptWorkspace{Path: filepath.ToSlash(filepath.Join(dir, declaration)), Root: filepath.ToSlash(dir), Members: []string{}, Manager: manager}, patterns: patterns, version: p.PackageManager})
	}
	return plans, warnings
}

func sameWorkspacePatterns(a, b []string) bool {
	setA, setB := map[string]bool{}, map[string]bool{}
	for _, p := range a {
		setA[p] = true
	}
	for _, p := range b {
		setB[p] = true
	}
	if len(setA) != len(setB) {
		return false
	}
	for p := range setA {
		if !setB[p] {
			return false
		}
	}
	return true
}

func hasJavaScriptLock(f map[string]bool) bool {
	for _, name := range []string{"package-lock.json", "pnpm-lock.yaml", "yarn.lock", "bun.lock", "bun.lockb"} {
		if f[name] {
			return true
		}
	}
	return false
}

func validateWorkspacePattern(pattern string) error {
	raw := strings.TrimPrefix(pattern, "!")
	if raw == "" || len(raw) > 4096 || strings.Count(raw, "/") > 127 || strings.ContainsAny(raw, "\\{}()") || filepath.IsAbs(raw) {
		return fmt.Errorf("unsupported workspace pattern %q", pattern)
	}
	for _, part := range strings.Split(raw, "/") {
		if part == ".." {
			return fmt.Errorf("workspace pattern outside project: %q", pattern)
		}
		if part != "**" {
			if _, err := filepath.Match(part, ""); err != nil {
				return fmt.Errorf("invalid workspace pattern %q", pattern)
			}
		}
	}
	return nil
}
func workspaceGlobMatch(pattern, path string) bool {
	a, b := strings.Split(strings.TrimPrefix(pattern, "./"), "/"), strings.Split(path, "/")
	memo := map[[2]int]bool{}
	visited := map[[2]int]bool{}
	var match func(int, int) bool
	match = func(i, j int) (result bool) {
		key := [2]int{i, j}
		if visited[key] {
			return memo[key]
		}
		defer func() { visited[key] = true; memo[key] = result }()
		if i == len(a) {
			return j == len(b)
		}
		if a[i] == "**" {
			return match(i+1, j) || (j < len(b) && match(i, j+1))
		}
		if j == len(b) {
			return false
		}
		ok, _ := filepath.Match(a[i], b[j])
		return ok && match(i+1, j+1)
	}
	return match(0, 0)
}
func workspacePatternsMatch(patterns []string, path string) bool {
	include := false
	for _, p := range patterns {
		if !strings.HasPrefix(p, "!") && workspaceGlobMatch(p, path) {
			include = true
		}
	}
	for _, p := range patterns {
		if strings.HasPrefix(p, "!") && workspaceGlobMatch(strings.TrimPrefix(p, "!"), path) {
			return false
		}
	}
	return include
}
func javascriptWorkspaceOwner(dir string, plans []javascriptWorkspacePlan) int {
	owner, best := -1, -1
	for i, p := range plans {
		rel, err := filepath.Rel(filepath.FromSlash(p.Root), dir)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			continue
		}
		if workspacePatternsMatch(p.patterns, filepath.ToSlash(rel)) && len(p.Root) > best {
			owner, best = i, len(p.Root)
		}
	}
	return owner
}
func javascriptWorkspaceContexts(root string, paths []string, dirs map[string]map[string]bool, plans []javascriptWorkspacePlan) (map[string]javascriptWorkspaceContext, []javascriptWorkspacePlan) {
	contexts := map[string]javascriptWorkspaceContext{}
	for _, p := range plans {
		c := javascriptWorkspaceContext{root: p.Root, manager: p.Manager}
		if p.Manager == "" {
			c.conflict = "Workspace manager unresolved: configure manager and commands explicitly"
		}
		contexts[filepath.FromSlash(p.Root)] = c
	}
	for _, dir := range paths {
		if !dirs[dir]["package.json"] {
			continue
		}
		owner := javascriptWorkspaceOwner(dir, plans)
		if owner < 0 {
			continue
		}
		// Inner declarations own their children; overlap conflicts remain explicit.
		p := &plans[owner]
		p.Members = append(p.Members, filepath.ToSlash(dir))
		c := javascriptWorkspaceContext{root: p.Root, manager: p.Manager}
		service, _, err := (adapters.JavaScript{}).Detect(root, dir, dirs[dir], nil)
		if p.Manager == "" {
			c.conflict = "Workspace manager unresolved: configure manager and commands explicitly"
		} else if err != nil || service == nil {
			continue
		} else {
			b, _ := read(root, dir, "package.json")
			var child struct{ PackageManager string }
			_ = json.Unmarshal(b, &child)
			if (child.PackageManager != "" || hasJavaScriptLock(dirs[dir])) && service.Manager != p.Manager {
				c.conflict = fmt.Sprintf("Workspace manager conflict: root %s uses %s, member uses %s; configure explicitly", p.Root, p.Manager, service.Manager)
			}
			if service.Manager == p.Manager && child.PackageManager != "" && p.version != "" && strings.Contains(child.PackageManager, "@") && strings.Contains(p.version, "@") && child.PackageManager != p.version {
				c.conflict = "Workspace manager version conflict: configure explicitly"
			}
			for _, pair := range []struct{ file, manager string }{{"package-lock.json", "npm"}, {"pnpm-lock.yaml", "pnpm"}, {"yarn.lock", "yarn"}, {"bun.lock", "bun"}, {"bun.lockb", "bun"}} {
				if dirs[dir][pair.file] && pair.manager != p.Manager {
					c.conflict = "Workspace lockfile manager conflict: configure explicitly"
				}
			}
		}
		for _, inner := range plans {
			if inner.Root == filepath.ToSlash(dir) {
				c.root = inner.Root
				if inner.Manager == "" {
					c.conflict = "Workspace manager unresolved: configure manager and commands explicitly"
				}
				if c.conflict == "" {
					c.manager = inner.Manager
				}
			}
		}
		contexts[dir] = c
	}
	for i := range plans {
		sort.Strings(plans[i].Members)
	}
	return contexts, plans
}
