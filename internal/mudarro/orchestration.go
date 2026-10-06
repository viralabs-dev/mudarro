package mudarro

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var orchestrationToken = regexp.MustCompile(`^[A-Za-z0-9@][A-Za-z0-9@/_.-]*$`)

// scanOrchestration reads declarations only; task graphs/plugins are evaluated
// by the selected runtime at explicitly requested execution time.
func scanOrchestration(root string, paths []string, dirs map[string]map[string]bool, r *Report) {
	owners := map[string]int{}
	ids := map[string]int{}
	for _, s := range r.Config.Services {
		ids[s.ID]++
	}
	for i, s := range r.Config.Services {
		if s.Language == "javascript" || s.Language == "typescript" {
			owners[s.Dir] = i
		}
	}
	var candidates []Suggestion
	warn := func(path string, err error) {
		r.Warnings = append(r.Warnings, fmt.Sprintf("%s: %v", filepath.ToSlash(path), err))
	}
	owner := func(dir, file string) (int, bool) {
		i, ok := owners[dir]
		if !ok {
			warn(filepath.Join(dir, file), fmt.Errorf("orchestration requires a JavaScript/TypeScript service in its root"))
			return 0, false
		}
		if ids[r.Config.Services[i].ID] != 1 {
			warn(filepath.Join(dir, file), fmt.Errorf("ambiguous service ID; configure orchestration commands explicitly"))
			return 0, false
		}
		return i, true
	}
	readObject := func(dir, file string) (map[string]json.RawMessage, bool) {
		b, err := read(root, dir, file)
		if err != nil {
			warn(filepath.Join(dir, file), err)
			return nil, false
		}
		var object map[string]json.RawMessage
		if err = json.Unmarshal(b, &object); err != nil || object == nil {
			warn(filepath.Join(dir, file), fmt.Errorf("invalid strict JSON object (JSONC/comments are not inferred)"))
			return nil, false
		}
		return object, true
	}
	appendTask := func(s Service, name, evidence string, args []string) {
		normalized := strings.NewReplacer("@", "", "/", "-", ".", "-", ":", "-").Replace(name)
		if !identifier.MatchString(normalized) {
			warn(evidence, fmt.Errorf("invalid orchestration suggestion name %q", name))
			return
		}
		group := "scripts"
		task := args[len(args)-1]
		if task == "test" || strings.HasSuffix(task, ":test") {
			group = "qualidade"
		} else if task == "build" || strings.HasSuffix(task, ":build") {
			group = "aplicacao"
		}
		candidates = append(candidates, Suggestion{Service: s.ID, Name: normalized, Purpose: "orchestration", Command: cmd(group, args...), Evidence: filepath.ToSlash(evidence)})
	}
	nxRoots := map[string]int{}
	for _, dir := range paths {
		f := dirs[dir]
		if f["turbo.json"] {
			r.Evidence = append(r.Evidence, Evidence{filepath.ToSlash(filepath.Join(dir, "turbo.json")), "turborepo"})
			object, ok := readObject(dir, "turbo.json")
			if !ok {
				continue
			}
			i, ok := owner(dir, "turbo.json")
			if !ok {
				continue
			}
			var tasks map[string]json.RawMessage
			if err := json.Unmarshal(object["tasks"], &tasks); err != nil || tasks == nil {
				warn(filepath.Join(dir, "turbo.json"), fmt.Errorf("declare explicit Turbo tasks; legacy pipeline/extended task inference is unsupported"))
				continue
			}
			if len(object["extends"]) > 0 {
				warn(filepath.Join(dir, "turbo.json"), fmt.Errorf("extends/task graph is declarative only; inherited tasks are not inferred"))
			}
			names := []string{}
			for name := range tasks {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				var config map[string]json.RawMessage
				if !orchestrationToken.MatchString(name) || json.Unmarshal(tasks[name], &config) != nil || config == nil {
					warn(filepath.Join(dir, "turbo.json"), fmt.Errorf("invalid task %q; configure explicitly", name))
					continue
				}
				appendTask(r.Config.Services[i], "turbo-"+name, filepath.Join(dir, "turbo.json"), []string{"turbo", "run", name})
			}
		}
	}
	for _, dir := range paths {
		f := dirs[dir]
		if f["nx.json"] {
			r.Evidence = append(r.Evidence, Evidence{filepath.ToSlash(filepath.Join(dir, "nx.json")), "nx"})
			// Invalid inner declarations remain a boundary, never borrowed from outer roots.
			nxRoots[dir] = -1
			object, ok := readObject(dir, "nx.json")
			if !ok {
				continue
			}
			i, ok := owner(dir, "nx.json")
			if !ok {
				continue
			}
			nxRoots[dir] = i
			if plugins := object["plugins"]; len(plugins) > 0 && string(plugins) != "[]" && string(plugins) != "null" {
				warn(filepath.Join(dir, "nx.json"), fmt.Errorf("plugins/inferred targets are not executed or resolved during scan"))
			}
		}
	}
	type project struct {
		dir, name string
		owner     int
		targets   map[string]json.RawMessage
	}
	projects := []project{}
	projectNames := map[string]int{}
	for _, dir := range paths {
		if !dirs[dir]["project.json"] {
			continue
		}
		best := -1
		rootDir := ""
		i := -1
		for nxDir, nxOwner := range nxRoots {
			rel, err := filepath.Rel(nxDir, dir)
			if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				continue
			}
			if len(nxDir) > best {
				best = len(nxDir)
				rootDir = nxDir
				i = nxOwner
			}
		}
		if best < 0 {
			continue
		}
		if i < 0 {
			warn(filepath.Join(dir, "project.json"), fmt.Errorf("Nx root declaration/owner unresolved; configure explicitly"))
			continue
		}
		object, ok := readObject(dir, "project.json")
		if !ok {
			continue
		}
		var name string
		_ = json.Unmarshal(object["name"], &name)
		if !orchestrationToken.MatchString(name) {
			warn(filepath.Join(dir, "project.json"), fmt.Errorf("explicit valid Nx project name required"))
			continue
		}
		var targets map[string]json.RawMessage
		if json.Unmarshal(object["targets"], &targets) != nil || targets == nil {
			warn(filepath.Join(dir, "project.json"), fmt.Errorf("explicit Nx targets object required"))
			continue
		}
		projects = append(projects, project{dir: dir, name: name, owner: i, targets: targets})
		projectNames[rootDir+"/"+name]++
	}
	for _, p := range projects {
		s := r.Config.Services[p.owner]
		if projectNames[s.Dir+"/"+p.name] > 1 {
			warn(filepath.Join(p.dir, "project.json"), fmt.Errorf("duplicate Nx project name %q; configure explicitly", p.name))
			continue
		}
		names := []string{}
		for name := range p.targets {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			var target map[string]json.RawMessage
			valid := json.Unmarshal(p.targets[name], &target) == nil && target != nil
			var executor, command string
			if value, exists := target["executor"]; exists {
				valid = valid && json.Unmarshal(value, &executor) == nil
			}
			if value, exists := target["command"]; exists {
				valid = valid && json.Unmarshal(value, &command) == nil
			}
			if !orchestrationToken.MatchString(name) || !valid || (strings.TrimSpace(executor) == "" && strings.TrimSpace(command) == "") {
				warn(filepath.Join(p.dir, "project.json"), fmt.Errorf("target %q needs explicit executor or command; inferred/default targets are not resolved", name))
				continue
			}
			appendTask(s, "nx-"+p.name+"-"+name, filepath.Join(p.dir, "project.json"), []string{"nx", "run", p.name + ":" + name})
		}
	}
	counts := map[string]int{}
	for _, su := range candidates {
		counts[su.Service+"/"+su.Name]++
	}
	for _, su := range r.Suggestions {
		counts[su.Service+"/"+su.Name]++
	}
	for _, su := range candidates {
		service := r.Config.Services[ownersByID(r.Config.Services, su.Service)]
		_, commandCollision := service.Commands[su.Name]
		if counts[su.Service+"/"+su.Name] > 1 || commandCollision {
			warn(su.Evidence, fmt.Errorf("orchestration suggestion %s collides; configure an explicit command", su.Name))
			continue
		}
		r.Suggestions = append(r.Suggestions, su)
	}
}
func ownersByID(services []Service, id string) int {
	for i, s := range services {
		if s.ID == id {
			return i
		}
	}
	return 0
}
