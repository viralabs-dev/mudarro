package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"strings"
	"testing"
)

func TestPlainMenuHasNoANSIEscapes(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	root := t.TempDir()
	if err := saveConfig(root+"/mudarro.yaml", sample()); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Main([]string{"menu", "--root", root}, "test", strings.NewReader("abc\n1\n0\n0\n"), &out); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "\x1b") {
		t.Fatal("ANSI in piped menu")
	}
	for _, text := range []string{"Serviços", "qualidade", "Opção inválida", "Voltar / sair"} {
		if !strings.Contains(out.String(), text) {
			t.Fatal("missing", text)
		}
	}
}
