package mudarro_test

import (
	"bytes"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestUICLIInitJSONPreservesExistingAndGenerateIdempotent(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"name":"probe","scripts":{"start":"node app.js","test":"node --test"}}`, "app.js": "console.log('probe')"})
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--format", "json"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	body, e := os.ReadFile(filepath.Join(root, "mudarro.json"))
	if e != nil {
		t.Fatal(e)
	}
	var cfg Config
	if e = json.Unmarshal(body, &cfg); e != nil {
		t.Fatal(e)
	}
	if exists(filepath.Join(root, "mudarro.yaml")) {
		t.Fatal("JSON init created YAML")
	}
	if e = Main([]string{"init", "--root", root, "--format", "json"}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("init overwrote config")
	}
	after, _ := os.ReadFile(filepath.Join(root, "mudarro.json"))
	if !bytes.Equal(body, after) {
		t.Fatal("config changed on rejected init")
	}
	for i := 0; i < 2; i++ {
		if e = Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
			t.Fatal(e)
		}
		if i == 0 {
			body, _ = os.ReadFile(filepath.Join(root, "menu.sh"))
		}
	}
	after, _ = os.ReadFile(filepath.Join(root, "menu.sh"))
	if !bytes.Equal(body, after) {
		t.Fatal("repeated generation changed launcher")
	}
}
func TestUICLISelectedConfigWrapperWithSpaces(t *testing.T) {
	root := fixture(t, map[string]string{})
	name := "config selected.json"
	c := sample()
	body, _ := json.Marshal(c)
	if e := os.WriteFile(filepath.Join(root, name), body, 0644); e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e := Main([]string{"generate", "--root", root, "--config", name}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	bin := t.TempDir()
	capture := filepath.Join(bin, "arguments")
	script := "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\" > \"$CAPTURE_ARGS\"\n"
	if e := os.WriteFile(filepath.Join(bin, "mudarro"), []byte(script), 0755); e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		path string
		want []string
	}{{"menu.sh", []string{"menu", "--config", name, "--root", root}}, {".mudarro/scripts/app/test.sh", []string{"run", "app:test", "--config", name, "--root", root}}} {
		cmd := exec.Command("bash", filepath.Join(root, tc.path))
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "PATH="+bin+":"+os.Getenv("PATH"), "CAPTURE_ARGS="+capture)
		if b, e := cmd.CombinedOutput(); e != nil {
			t.Fatalf("wrapper %v %s", e, b)
		}
		got, e := os.ReadFile(capture)
		if e != nil {
			t.Fatal(e)
		}
		if !reflect.DeepEqual(strings.Split(strings.TrimSuffix(string(got), "\n"), "\n"), tc.want) {
			t.Fatalf("wrapper argv %q", got)
		}
	}
}
func TestUICLIAmbiguousConfigNeedsSelection(t *testing.T) {
	c := sample()
	body, _ := json.Marshal(c)
	root := fixture(t, map[string]string{"mudarro.json": string(body)})
	if e := saveConfig(filepath.Join(root, "mudarro.yaml"), c); e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e := Main([]string{"menu", "--root", root}, "test", strings.NewReader("0\n"), &out); e == nil {
		t.Fatal("ambiguous config accepted")
	}
	if e := Main([]string{"menu", "--root", root, "--config", "mudarro.json"}, "test", strings.NewReader("0\n"), &out); e != nil {
		t.Fatal(e)
	}
}
func TestUICLILocaleOverrideDoesNotWriteConfig(t *testing.T) {
	c := sample()
	root := fixture(t, map[string]string{})
	path := filepath.Join(root, "mudarro.yaml")
	if e := saveConfig(path, c); e != nil {
		t.Fatal(e)
	}
	before, _ := os.ReadFile(path)
	for _, tc := range []struct{ locale, want string }{{"en", "Services"}, {"pt-BR", "Serviços"}} {
		var out bytes.Buffer
		if e := Main([]string{"menu", "--root", root, "--locale", tc.locale}, "test", strings.NewReader("0\n"), &out); e != nil {
			t.Fatal(e)
		}
		if !strings.Contains(out.String(), tc.want) {
			t.Fatalf("locale %s: %s", tc.locale, out.String())
		}
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("UI override wrote config")
	}
	var out bytes.Buffer
	if e := Main([]string{"menu", "--root", root, "--theme", "neon"}, "test", strings.NewReader("0\n"), &out); e == nil {
		t.Fatal("invalid theme accepted")
	}
}
func TestUICLIPreviewDisabledHasNoPreviewHint(t *testing.T) {
	c := sample()
	disabled := false
	c.UI.Preview = &PreviewConfig{Enabled: &disabled, Mouse: "off"}
	root := fixture(t, map[string]string{})
	if e := saveConfig(filepath.Join(root, "mudarro.yaml"), c); e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e := Main([]string{"menu", "--root", root}, "test", strings.NewReader("1\n1\n0\n0\n0\n"), &out); e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out.String(), "pN") || strings.Contains(out.String(), "p1") {
		t.Fatalf("disabled preview advertised: %s", out.String())
	}
}
