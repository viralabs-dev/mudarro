package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"path/filepath"
	"sort"
	"strings"
)

// Pipenv scripts remain opt-in suggestions; scan never evaluates their bodies.
func detectPipenv(root, dir string, files map[string]bool) (*model.Service, []model.Suggestion, error) {
	s := BaseService(dir, "python")
	if !files["Pipfile"] {
		s.Pending = append(s.Pending, "Pipfile.lock has no Pipfile: restore the manifest and declare commands")
		return s, nil, nil
	}
	data, err := read(root, dir, "Pipfile")
	if err != nil {
		return nil, nil, err
	}
	var manifest struct{ Scripts map[string]any }
	if err = toml.Unmarshal(data, &manifest); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", filepath.Join(dir, "Pipfile"), err)
	}
	if files["Pipfile.lock"] {
		data, err = read(root, dir, "Pipfile.lock")
		if err != nil {
			return nil, nil, err
		}
		var lock map[string]json.RawMessage
		if err = json.Unmarshal(data, &lock); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", filepath.Join(dir, "Pipfile.lock"), err)
		}
		if lock == nil {
			return nil, nil, fmt.Errorf("%s: expected lockfile object", filepath.Join(dir, "Pipfile.lock"))
		}
	}
	if files["uv.lock"] || files["poetry.lock"] {
		s.Pending = append(s.Pending, "Pipfile conflicts with other manager lockfiles: declare manager and commands explicitly")
		return s, nil, nil
	}
	s.Manager = "pipenv"
	if files["Pipfile.lock"] {
		s.Commands["install"] = cmd("dependencias", "pipenv", "sync")
	} else {
		s.Commands["install"] = cmd("dependencias", "pipenv", "install")
	}
	var suggestions []model.Suggestion
	names := map[string][]string{}
	keys := []string{}
	for k := range manifest.Scripts {
		keys = append(keys, k)
		normalized := strings.NewReplacer(":", "-", ".", "-").Replace(k)
		names[normalized] = append(names[normalized], k)
	}
	sort.Strings(keys)
	reported := map[string]bool{}
	for _, k := range keys {
		name := strings.NewReplacer(":", "-", ".", "-").Replace(k)
		if !identifier.MatchString(name) {
			s.Pending = append(s.Pending, fmt.Sprintf("Pipfile script %q needs an explicit command name", k))
			continue
		}
		if len(names[name]) > 1 {
			if !reported[name] {
				s.Pending = append(s.Pending, "Pipfile scripts collide at "+name+": declare commands explicitly")
				reported[name] = true
			}
			continue
		}
		body, ok := manifest.Scripts[k].(string)
		if !ok || strings.TrimSpace(body) == "" {
			s.Pending = append(s.Pending, "Pipfile script "+k+": only nonempty string scripts are supported as suggestions; declare command explicitly")
			continue
		}
		suggestions = append(suggestions, model.Suggestion{Name: name, Command: cmd("scripts", "pipenv", "run", k), Evidence: filepath.Join(dir, "Pipfile")})
	}
	if files["manage.py"] {
		s.Framework = "django"
		s.Commands["start"] = cmd("aplicacao", "pipenv", "run", "python", "manage.py", "runserver")
		s.Commands["test"] = cmd("qualidade", "pipenv", "run", "python", "manage.py", "test")
	} else {
		s.Pending = append(s.Pending, "Declare commands.start; módulo Python não é inferido")
	}
	return s, suggestions, nil
}
