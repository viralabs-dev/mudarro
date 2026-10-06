package mudarro_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/viralabs-dev/mudarro/internal/mudarro"
)

func TestMenuUsesExplicitConsumerName(t *testing.T) {
	root := t.TempDir()
	c := sample()
	c.Name = "Atlas API"
	if err := saveConfig(filepath.Join(root, "mudarro.yaml"), c); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Main([]string{"menu", "--root", root}, "test", strings.NewReader("0\n"), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Atlas API") {
		t.Fatal("explicit configuration name missing")
	}
	if strings.Contains(out.String(), filepath.Base(root)) {
		t.Fatal("directory replaced explicit project name")
	}
}
