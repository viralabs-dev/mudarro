package adapters

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"path/filepath"
	"sort"
	"strings"
)

type Rust struct{}
type cargoManifest struct {
	Package *struct {
		Name      string
		Workspace string `toml:"workspace"`
	} `toml:"package"`
	Workspace *struct {
		Members        []string
		DefaultMembers []string `toml:"default-members"`
		Exclude        []string
	} `toml:"workspace"`
	Lib *cargoTarget  `toml:"lib"`
	Bin []cargoTarget `toml:"bin"`
}
type cargoTarget struct {
	Name string
	Path string
}

func readCargo(root, dir string) (cargoManifest, error) {
	var m cargoManifest
	b, e := read(root, dir, "Cargo.toml")
	if e == nil {
		e = toml.Unmarshal(b, &m)
	}
	if e != nil {
		return m, fmt.Errorf("%s: %w", filepath.ToSlash(filepath.Join(dir, "Cargo.toml")), e)
	}
	return m, nil
}
func cargoExcluded(rel string, excludes []string) bool {
	for _, part := range strings.Split(filepath.ToSlash(rel), "/") {
		if part == "target" || IgnoredDirectory(part) {
			return true
		}
	}
	for _, pattern := range excludes {
		matched, _ := filepath.Match(pattern, filepath.ToSlash(rel))
		if matched || rel == pattern || strings.HasPrefix(filepath.ToSlash(rel), strings.TrimSuffix(pattern, "/")+"/") {
			return true
		}
	}
	return false
}
func cargoMemberPath(root, dir, value string, excludes []string) (string, error) {
	if value == "" || filepath.IsAbs(value) || strings.ContainsAny(value, "*?[]{}") {
		return "", fmt.Errorf("workspace member %q requires an explicit relative path; glob expansion is not supported", value)
	}
	rel := filepath.Clean(filepath.Join(dir, value))
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("workspace member %q is outside project", value)
	}
	if cargoExcluded(rel, excludes) || cargoExcluded(filepath.Join(rel, "Cargo.toml"), excludes) {
		return "", fmt.Errorf("workspace member %q is excluded and not validated", value)
	}
	_, e := read(root, rel, "Cargo.toml")
	if e != nil {
		return "", fmt.Errorf("workspace member %q: %w", value, e)
	}
	return rel, nil
}
func cargoWorkspace(root, dir string, m cargoManifest, excludes []string) ([]string, []string) {
	if m.Workspace == nil {
		return nil, nil
	}
	members := []string{}
	issues := []string{}
	seen := map[string]bool{}
	names := map[string]string{}
	memberExcludes := append([]string{}, excludes...)
	for _, v := range m.Workspace.Exclude {
		if filepath.IsAbs(v) || strings.ContainsAny(v, "*?[]{}") {
			issues = append(issues, "workspace exclude requires explicit contained relative path: "+v)
			continue
		}
		rel := filepath.Clean(filepath.Join(dir, v))
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			issues = append(issues, "workspace exclude outside project: "+v)
			continue
		}
		memberExcludes = append(memberExcludes, filepath.ToSlash(rel))
	}
	for _, value := range m.Workspace.Members {
		rel, e := cargoMemberPath(root, dir, value, memberExcludes)
		if e != nil {
			issues = append(issues, e.Error())
			continue
		}
		if seen[rel] {
			issues = append(issues, "duplicate workspace member: "+value)
			continue
		}
		child, e := readCargo(root, rel)
		if e != nil {
			issues = append(issues, e.Error())
			continue
		}
		if child.Package == nil || strings.TrimSpace(child.Package.Name) == "" {
			issues = append(issues, "workspace member has no named package: "+value)
			continue
		}
		if prior, ok := names[child.Package.Name]; ok {
			issues = append(issues, fmt.Sprintf("duplicate workspace package name %q: %s and %s", child.Package.Name, prior, rel))
		}
		names[child.Package.Name] = rel
		seen[rel] = true
		members = append(members, rel)
	}
	if m.Package != nil && strings.TrimSpace(m.Package.Name) != "" {
		if !seen[dir] {
			if prior, ok := names[m.Package.Name]; ok {
				issues = append(issues, "duplicate workspace package name: "+m.Package.Name+" at "+prior+" and "+dir)
			}
			members = append(members, dir)
			seen[dir] = true
		}
	}
	if len(members) == 0 {
		issues = append(issues, "workspace has no validated explicit packages")
	}
	for _, value := range m.Workspace.DefaultMembers {
		rel, e := cargoMemberPath(root, dir, value, memberExcludes)
		if e != nil {
			issues = append(issues, "default-members: "+e.Error())
			continue
		}
		if !seen[rel] {
			issues = append(issues, "default-member is not a workspace member: "+value)
		}
	}
	sort.Strings(members)
	return members, issues
}
func (Rust) Detect(root, dir string, files map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if !files["Cargo.toml"] || cargoExcluded(filepath.Join(dir, "Cargo.toml"), excludes) {
		return nil, nil, nil
	}
	m, e := readCargo(root, dir)
	if e != nil {
		return nil, nil, e
	}
	if m.Package == nil && m.Workspace == nil {
		return nil, nil, fmt.Errorf("%s: Cargo manifest requires package or workspace", filepath.ToSlash(filepath.Join(dir, "Cargo.toml")))
	}
	if m.Package != nil && strings.TrimSpace(m.Package.Name) == "" {
		return nil, nil, fmt.Errorf("%s: package requires a name", filepath.ToSlash(filepath.Join(dir, "Cargo.toml")))
	}
	s := BaseService(dir, "rust")
	if m.Package == nil && m.Workspace != nil {
		s.ID += "-workspace"
	}
	s.Manager = "cargo"
	s.Pending = append(s.Pending, "Cargo build/test are opt-in; declare startup/bin command explicitly; dependency graphs/build scripts are not evaluated")
	issues := []string{}
	names := map[string]bool{}
	targets := append([]cargoTarget{}, m.Bin...)
	if m.Lib != nil {
		lib := *m.Lib
		lib.Name = ""
		targets = append(targets, lib)
	}
	for _, target := range targets {
		if target.Name != "" {
			if names[target.Name] {
				issues = append(issues, "duplicate explicit target name: "+target.Name)
			}
			names[target.Name] = true
		}
		if target.Path != "" {
			rel := filepath.Clean(filepath.Join(dir, target.Path))
			if filepath.IsAbs(target.Path) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || cargoExcluded(rel, excludes) {
				issues = append(issues, "target path outside project or excluded: "+target.Path)
				continue
			}
			if _, e := read(root, dir, target.Path); e != nil {
				issues = append(issues, "target path: "+e.Error())
			}
		}
	}
	if m.Package != nil && m.Package.Workspace != "" {
		issues = append(issues, "Explicit package.workspace is not resolved; declare service commands")
	} else if m.Workspace != nil {
		members, warnings := cargoWorkspace(root, dir, m, excludes)
		issues = append(issues, warnings...)
		s.Pending = append(s.Pending, "Cargo workspace explicit members: "+strings.Join(members, ", "))
	} else {
		for parent := filepath.Dir(dir); parent != ".." && parent != dir; parent = filepath.Dir(parent) {
			if cargoExcluded(filepath.Join(parent, "Cargo.toml"), excludes) {
				s.Pending = append(s.Pending, "Cargo ancestor manifest excluded; workspace ownership not validated")
				break
			}
			ancestor, err := readCargo(root, parent)
			if err == nil && ancestor.Workspace != nil {
				members, _ := cargoWorkspace(root, parent, ancestor, excludes)
				for _, member := range members {
					if member == dir {
						s.Pending = append(s.Pending, "Cargo explicit workspace owner: "+filepath.ToSlash(parent))
						break
					}
				}
				break
			}
			if parent == "." {
				break
			}
		}
	}
	s.Pending = append(s.Pending, issues...)
	if len(issues) > 0 {
		return s, nil, nil
	}
	evidence := filepath.ToSlash(filepath.Join(dir, "Cargo.toml"))
	suggestions := []model.Suggestion{{Name: "build", Purpose: "cargo", Command: cmd("aplicacao", "cargo", "build"), Evidence: evidence}, {Name: "test", Purpose: "cargo", Command: cmd("qualidade", "cargo", "test"), Evidence: evidence}}
	return s, suggestions, nil
}
