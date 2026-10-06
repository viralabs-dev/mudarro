package adapters

import (
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"path/filepath"
)

// Elixir recognizes project evidence without evaluating executable Mix source.
type Elixir struct{}

func (Elixir) Detect(root, dir string, files map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if !files["mix.exs"] {
		return nil, nil, nil
	}
	// Reading through the shared boundary verifies containment, regular-file status
	// and the manifest size limit. Source syntax and project functions stay unevaluated.
	if _, err := read(root, dir, "mix.exs"); err != nil {
		return nil, nil, err
	}
	s := BaseService(dir, "elixir")
	s.Manager = "mix"
	s.Pending = append(s.Pending, "Mix source is executable code: project syntax, aliases and dependencies were not evaluated", "Declare commands.start explicitly; Mix startup was not inferred", "Declare umbrella ownership explicitly when applicable; apps_path was not evaluated")
	evidence := filepath.ToSlash(filepath.Join(dir, "mix.exs"))
	return s, []model.Suggestion{
		{Service: s.ID, Name: "compile", Purpose: "mix", Command: cmd("aplicacao", "mix", "compile"), Evidence: evidence},
		{Service: s.ID, Name: "test", Purpose: "mix", Command: cmd("qualidade", "mix", "test"), Evidence: evidence},
	}, nil
}
