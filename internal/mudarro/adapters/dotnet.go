package adapters

import (
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// DotNet offers explicitly named project targets without evaluating MSBuild or
// guessing a solution graph, startup project, framework, or server.
type DotNet struct{}

func (DotNet) Detect(root, dir string, files map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	var projects, solutions []string
	for name := range files {
		switch strings.ToLower(filepath.Ext(name)) {
		case ".csproj":
			projects = append(projects, name)
		case ".sln", ".slnx":
			solutions = append(solutions, name)
		}
	}
	sort.Strings(projects)
	sort.Strings(solutions)
	if len(projects)+len(solutions) == 0 {
		return nil, nil, nil
	}
	s := BaseService(dir, "csharp")
	s.Manager = "dotnet"
	s.Pending = []string{"MSBuild properties, imports and targets were not evaluated; declare run/start explicitly", "Solution project graphs were not resolved; select contained independently detected projects explicitly"}
	for _, name := range solutions {
		b, e := read(root, dir, name)
		if e != nil {
			return nil, nil, e
		}
		if strings.HasSuffix(strings.ToLower(name), ".slnx") {
			if e = declarativeXML(b, "Solution"); e != nil {
				return nil, nil, e
			}
		}
		s.Pending = append(s.Pending, "Solution evidence only: "+name)
	}
	var suggestions []model.Suggestion
	for index, name := range projects {
		b, e := read(root, dir, name)
		if e != nil {
			return nil, nil, e
		}
		if e = declarativeXML(b, "Project"); e != nil {
			return nil, nil, e
		}
		for _, task := range []string{"build", "test"} {
			choice := task
			if len(projects) > 1 {
				choice = task + "-project-" + strconv.Itoa(index+1)
			}
			group := "aplicacao"
			if task == "test" {
				group = "qualidade"
			}
			suggestions = append(suggestions, model.Suggestion{Service: s.ID, Name: choice, Purpose: "dotnet", Command: cmd(group, "dotnet", task, "./"+name), Evidence: filepath.ToSlash(filepath.Join(dir, name))})
		}
	}
	return s, suggestions, nil
}
