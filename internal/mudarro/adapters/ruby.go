package adapters

import (
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"sort"
	"strings"
)

// Ruby files are executable DSL evidence, never project-discovery programs.
type Ruby struct{}

func (Ruby) Detect(root, dir string, files map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	names := []string{}
	for name := range files {
		if name == "Gemfile" || strings.HasSuffix(name, ".gemspec") {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return nil, nil, nil
	}
	sort.Strings(names)
	for _, name := range names {
		if _, e := read(root, dir, name); e != nil {
			return nil, nil, fmt.Errorf("Ruby evidence %s: %w", name, e)
		}
	}
	s := BaseService(dir, "ruby")
	s.Manager = "bundler"
	s.Pending = append(s.Pending, "Ruby Gemfile/gemspec code was not evaluated: dependencies, aliases and project semantics remain unvalidated", "Declare build/test/start commands explicitly; Rake tasks and Rails startup were not inferred")
	if len(names) > 1 {
		s.Pending = append(s.Pending, "Multiple Ruby manifests: declare ownership explicitly; no DSL was evaluated")
	}
	if files["Gemfile.lock"] {
		if _, e := read(root, dir, "Gemfile.lock"); e != nil {
			return nil, nil, e
		}
		s.Pending = append(s.Pending, "Gemfile.lock is bounded evidence only; lock syntax and installed gems were not validated")
	}
	return s, nil, nil
}
