package mudarro

import (
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

// Workspace describes parsed membership, not a module dependency-resolution plan.
type Workspace struct {
	Path         string            `json:"path"`
	Members      []string          `json:"members"`
	Excluded     []string          `json:"excluded,omitempty"`
	Modules      map[string]string `json:"modules,omitempty"`
	GoVersion    string            `json:"go_version,omitempty"`
	Toolchain    string            `json:"toolchain,omitempty"`
	Replacements int               `json:"replacements,omitempty"`
}

func scanWorkspace(root, dir string, excludes []string) (Workspace, []string) {
	w := Workspace{Path: filepath.ToSlash(filepath.Join(dir, "go.work")), Members: []string{}, Modules: map[string]string{}}
	warnings := []string{}
	warn := func(err error) { warnings = append(warnings, fmt.Sprintf("%s: %v", w.Path, err)) }
	b, err := read(root, dir, "go.work")
	if err != nil {
		warn(err)
		return w, warnings
	}
	parsed, err := modfile.ParseWork(w.Path, b, nil)
	if err != nil {
		warn(err)
		return w, warnings
	}
	if parsed.Go != nil {
		w.GoVersion = parsed.Go.Version
	}
	if parsed.Toolchain != nil {
		w.Toolchain = parsed.Toolchain.Name
	}
	w.Replacements = len(parsed.Replace)
	seen := map[string]bool{}
	moduleOwners := map[string]string{}
	excludedSeen := map[string]bool{}
	for _, use := range parsed.Use {
		path := use.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, dir, path)
		}
		rel, err := filepath.Rel(root, filepath.Clean(path))
		if err != nil {
			warn(fmt.Errorf("member %q: %w", use.Path, err))
			continue
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			warn(fmt.Errorf("member %q is outside the project", use.Path))
			continue
		}
		if workspaceMemberExcluded(filepath.ToSlash(rel), excludes) || workspaceMemberExcluded(filepath.ToSlash(filepath.Join(rel, "go.mod")), excludes) {
			normalized := filepath.ToSlash(rel)
			if !excludedSeen[normalized] {
				w.Excluded = append(w.Excluded, normalized)
				excludedSeen[normalized] = true
			}
			warn(fmt.Errorf("member %q excluded from scan; membership not validated for this entry", normalized))
			continue
		}
		member, err := safePath(root, rel)
		if err != nil {
			warn(fmt.Errorf("member %q: %w", use.Path, err))
			continue
		}
		st, err := os.Stat(member)
		if err != nil {
			warn(fmt.Errorf("member %q: %w", use.Path, err))
			continue
		}
		if !st.IsDir() {
			warn(fmt.Errorf("member %q is not a directory", use.Path))
			continue
		}
		rel = filepath.ToSlash(rel)
		if seen[rel] {
			warn(fmt.Errorf("duplicate workspace member %q", rel))
			continue
		}
		seen[rel] = true
		source, err := read(root, rel, "go.mod")
		if err != nil {
			warn(fmt.Errorf("member %q go.mod: %w", rel, err))
			continue
		}
		module, err := modfile.Parse(filepath.ToSlash(filepath.Join(rel, "go.mod")), source, nil)
		if err != nil {
			warn(fmt.Errorf("member %q go.mod: %w", rel, err))
			continue
		}
		if module.Module == nil || module.Module.Mod.Path == "" {
			warn(fmt.Errorf("member %q go.mod: missing module declaration", rel))
			continue
		}
		if previous, exists := moduleOwners[module.Module.Mod.Path]; exists {
			warn(fmt.Errorf("module %q appears in multiple workspace members: %s and %s", module.Module.Mod.Path, previous, rel))
		} else {
			moduleOwners[module.Module.Mod.Path] = rel
		}
		w.Members = append(w.Members, rel)
		w.Modules[rel] = module.Module.Mod.Path
	}
	sort.Strings(w.Members)
	sort.Strings(w.Excluded)
	// Replacement targets and toolchains remain declarative; Go resolves them
	// only at explicitly requested execution time under inherited/off policy.
	return w, warnings
}

func workspaceMemberExcluded(rel string, excludes []string) bool {
	for _, part := range strings.Split(rel, "/") {
		if adapters.IgnoredDirectory(part) {
			return true
		}
	}
	for _, pattern := range excludes {
		match, _ := filepath.Match(pattern, rel)
		if match || rel == pattern || strings.HasPrefix(rel, strings.TrimSuffix(pattern, "/")+"/") {
			return true
		}
	}
	return false
}
