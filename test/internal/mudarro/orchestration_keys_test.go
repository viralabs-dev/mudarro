package mudarro_test

import (
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"testing"
)

func TestOrchestrationNxTargetKeysAreCaseSensitive(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{}`, "nx.json": `{}`, "project.json": `{"name":"app","targets":{"upper":{"COMMAND":"never"},"mixed":{"Executor":"never"},"wrong-type":{"executor":3},"valid":{"command":"never"}}}`})
	r, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	suggestions := orchestrationSuggestions(r)
	if len(suggestions) != 1 || suggestions[0].Name != "nx-app-valid" {
		t.Fatalf("unexpected targets: %+v", suggestions)
	}
	if len(r.Warnings) < 3 {
		t.Fatalf("invalid targets lack diagnostic: %+v", r.Warnings)
	}
}
