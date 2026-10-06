package mudarro_test

import (
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"strings"
	"testing"
)

func TestLanguageServiceIDCollisionDiagnosed(t *testing.T) {
	root := fixture(t, map[string]string{"a/b/composer.json": `{"scripts":{"check":"php check.php"}}`, "a-b/composer.json": `{"scripts":{"other":"php other.php"}}`})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Config.Services) != 1 || len(r.Suggestions) != 1 {
		t.Fatalf("duplicate services/suggestions: %+v", r)
	}
	if !strings.Contains(strings.Join(r.Warnings, "\n"), "service ID collision") {
		t.Fatal(r.Warnings)
	}
	if e := r.Config.Validate(root); e != nil {
		t.Fatal(e)
	}
}
