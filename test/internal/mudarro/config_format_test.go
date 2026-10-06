package mudarro_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"gopkg.in/yaml.v3"
)

func TestConfigFormatsLegacyDefaultsWithoutWrites(t *testing.T) {
	for _, format := range []string{"yaml", "json"} {
		t.Run(format, func(t *testing.T) {
			c := sample()
			c.UI = nil // This fixture must remain legacy regardless of shared menu locale fixtures.
			var body []byte
			var err error
			if format == "json" {
				body, err = json.Marshal(c)
			} else {
				body, err = yaml.Marshal(c)
			}
			if err != nil {
				t.Fatal(err)
			}
			root := fixture(t, map[string]string{"mudarro." + format: string(body)})
			before := supportFileHashes(t, root)
			loaded, err := Load(root)
			if err != nil {
				t.Fatal(err)
			}
			expectedJSON, _ := json.Marshal(c)
			actualJSON, _ := json.Marshal(loaded)
			if string(expectedJSON) != string(actualJSON) || loaded.UI != nil {
				t.Fatalf("legacy changed: %+v", loaded)
			}
			want := UISettings{Locale: "en", Theme: "auto", Density: "comfortable", Lettering: "auto", Preview: PreviewSettings{Enabled: true, Mouse: "auto"}}
			if loaded.UIOptions() != want {
				t.Fatalf("defaults: %+v", loaded.UIOptions())
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("loader wrote configuration")
			}
		})
	}
}
func TestConfigUIExplicitAndPartialDefaults(t *testing.T) {
	for _, format := range []string{"yaml", "json"} {
		t.Run(format, func(t *testing.T) {
			disabled := false
			c := sample()
			c.UI = &UIConfig{Locale: "pt-BR", Theme: "dark", Density: "compact", Lettering: "text", Preview: &PreviewConfig{Enabled: &disabled, Mouse: "off"}}
			var body []byte
			if format == "json" {
				body, _ = json.Marshal(c)
			} else {
				body, _ = yaml.Marshal(c)
			}
			root := fixture(t, map[string]string{"custom." + format: string(body)})
			loaded, err := LoadFile(root, "custom."+format)
			if err != nil {
				t.Fatal(err)
			}
			want := UISettings{Locale: "pt-BR", Theme: "dark", Density: "compact", Lettering: "text", Preview: PreviewSettings{Enabled: false, Mouse: "off"}}
			if loaded.UIOptions() != want {
				t.Fatal(loaded.UIOptions())
			}
		})
	}
	c := sample()
	c.UI = &UIConfig{Theme: "light", Lettering: "ascii", Preview: &PreviewConfig{Mouse: "on"}}
	got := c.UIOptions()
	if got.Locale != "en" || got.Density != "comfortable" || !got.Preview.Enabled || got.Preview.Mouse != "on" {
		t.Fatal(got)
	}
	if err := c.Validate(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}
func TestConfigUIInvalidEnums(t *testing.T) {
	for _, field := range []string{"locale", "theme", "density", "lettering", "mouse"} {
		t.Run(field, func(t *testing.T) {
			c := sample()
			c.UI = &UIConfig{}
			switch field {
			case "locale":
				c.UI.Locale = "fr"
			case "theme":
				c.UI.Theme = "neon"
			case "density":
				c.UI.Density = "dense"
			case "lettering":
				c.UI.Lettering = "figlet"
			case "mouse":
				c.UI.Preview = &PreviewConfig{Mouse: "yes"}
			}
			b, _ := json.Marshal(c)
			root := fixture(t, map[string]string{"mudarro.json": string(b)})
			if _, err := Load(root); err == nil {
				t.Fatal("invalid enum accepted")
			}
		})
	}
}
func TestConfigStrictDocumentsAndKeys(t *testing.T) {
	for _, tc := range []struct{ name, format, body string }{
		{"json-unknown", "json", `{"version":1,"name":"x","services":[],"unknown":true}`},
		{"json-ui-unknown", "json", `{"version":1,"name":"x","services":[],"ui":{"preview":{"unknown":true}}}`},
		{"json-duplicate", "json", `{"version":1,"version":1,"name":"x","services":[]}`},
		{"json-nested-duplicate", "json", `{"version":1,"name":"x","services":[],"ui":{"preview":{"mouse":"auto","mouse":"off"}}}`},
		{"json-multiple", "json", `{"version":1,"name":"x","services":[]} {}`},
		{"json-trailing", "json", `{"version":1,"name":"x","services":[]} garbage`},
		{"json-comment", "json", "{\"version\":1,\"name\":\"x\",\"services\":[]} // comment"},
		{"yaml-unknown", "yaml", "version: 1\nname: x\nservices: []\nui:\n  preview:\n    unknown: true\n"},
		{"yaml-duplicate", "yaml", "version: 1\nname: x\nservices: []\nui:\n  theme: dark\n  theme: light\n"},
		{"yaml-multiple", "yaml", "version: 1\nname: x\nservices: []\n---\n{}\n"},
		{"yaml-trailing-empty", "yaml", "version: 1\nname: x\nservices: []\n---\n"},
		{"json-wrong-boolean", "json", `{"version":1,"name":"x","services":[],"ui":{"preview":{"enabled":"false"}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"mudarro." + tc.format: tc.body})
			if _, err := Load(root); err == nil {
				t.Fatal("invalid document accepted")
			}
		})
	}
}
func TestConfigDiscoveryAndExplicitSelection(t *testing.T) {
	root := fixture(t, map[string]string{"mudarro.yaml": "version: 1\nname: yaml\nservices: []\n", "mudarro.json": `{"version":1,"name":"json","services":[]}`})
	if _, err := Load(root); err == nil {
		t.Fatal("ambiguous discovery accepted")
	}
	for _, name := range []string{"mudarro.yaml", "mudarro.json"} {
		c, err := LoadFile(root, name)
		if err != nil {
			t.Fatal(err)
		}
		if c.Name != strings.TrimPrefix(filepath.Ext(name), ".") {
			t.Fatal(c.Name)
		}
	}
	c, err := LoadFile(root, filepath.Join(root, "mudarro.json"))
	if err != nil || c.Name != "json" {
		t.Fatalf("absolute inside root %v %v", c, err)
	}
	if _, err := Load(t.TempDir()); !os.IsNotExist(err) {
		t.Fatal("missing config should expose absence", err)
	}
}
func TestConfigBoundedReadsAndGuardedPaths(t *testing.T) {
	body := `{"version":1,"name":"x","services":[]}`
	for _, extra := range []int{0, 1} {
		t.Run(string(rune('0'+extra)), func(t *testing.T) {
			root := fixture(t, map[string]string{"mudarro.json": body + strings.Repeat(" ", (4<<20)-len(body)+extra)})
			_, err := Load(root)
			if (err != nil) != (extra == 1) {
				t.Fatalf("boundary err=%v", err)
			}
		})
	}
	root := t.TempDir()
	outside := fixture(t, map[string]string{"other.json": body})
	if err := os.Symlink(filepath.Join(outside, "other.json"), filepath.Join(root, "linked.json")); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"linked.json", "../other.json", filepath.Join(outside, "other.json")} {
		if _, err := LoadFile(root, name); err == nil {
			t.Fatal("unguarded path accepted", name)
		}
	}
	if err := os.Mkdir(filepath.Join(root, "dir.json"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(root, "dir.json"); err == nil {
		t.Fatal("directory accepted")
	}
	if _, err := LoadFile(root, "config.toml"); err == nil {
		t.Fatal("unsupported extension accepted")
	}
}
