package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"strings"
	"testing"
)

func TestScriptNormalizationCollisionRequiresExplicitCommand(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"scripts":{"start":"node app.js","check:fast":"echo first","check.fast":"echo second"}}`})
	r, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range r.Suggestions {
		if s.Name == "check-fast" {
			t.Fatal("ambiguous normalized script offered")
		}
	}
	if !strings.Contains(strings.Join(r.Config.Services[0].Pending, " "), "check-fast") {
		t.Fatal("collision missing from pending")
	}
	var out bytes.Buffer
	if Main([]string{"init", "--root", root, "--select", "app-javascript:check-fast"}, "test", strings.NewReader(""), &out) == nil {
		t.Fatal("ambiguous selection accepted")
	}
	if _, err := os.Stat(root + "/mudarro.yaml"); !os.IsNotExist(err) {
		t.Fatal("failed selection wrote config")
	}
}
