package mudarro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// UIConfig is optional; omitted fields are interpreted at read time, never written back.
type UIConfig struct {
	Locale    string         `yaml:"locale,omitempty" json:"locale,omitempty"`
	Theme     string         `yaml:"theme,omitempty" json:"theme,omitempty"`
	Density   string         `yaml:"density,omitempty" json:"density,omitempty"`
	Lettering string         `yaml:"lettering,omitempty" json:"lettering,omitempty"`
	Preview   *PreviewConfig `yaml:"preview,omitempty" json:"preview,omitempty"`
}
type PreviewConfig struct {
	Enabled *bool  `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Mouse   string `yaml:"mouse,omitempty" json:"mouse,omitempty"`
}
type UISettings struct {
	Locale, Theme, Density, Lettering string
	Preview                           PreviewSettings
}
type PreviewSettings struct {
	Enabled bool
	Mouse   string
}

func (c Config) UIOptions() UISettings {
	result := UISettings{Locale: "en", Theme: "auto", Density: "comfortable", Lettering: "auto", Preview: PreviewSettings{Enabled: true, Mouse: "auto"}}
	if c.UI == nil {
		return result
	}
	u := c.UI
	if u.Locale != "" {
		result.Locale = u.Locale
	}
	if u.Theme != "" {
		result.Theme = u.Theme
	}
	if u.Density != "" {
		result.Density = u.Density
	}
	if u.Lettering != "" {
		result.Lettering = u.Lettering
	}
	if u.Preview != nil {
		if u.Preview.Enabled != nil {
			result.Preview.Enabled = *u.Preview.Enabled
		}
		if u.Preview.Mouse != "" {
			result.Preview.Mouse = u.Preview.Mouse
		}
	}
	return result
}
func (c Config) validateUI() error {
	u := c.UIOptions()
	for _, field := range []struct {
		name, value string
		allowed     []string
	}{
		{"locale", u.Locale, []string{"en", "pt-BR"}}, {"theme", u.Theme, []string{"auto", "light", "dark"}}, {"density", u.Density, []string{"comfortable", "compact"}}, {"lettering", u.Lettering, []string{"auto", "ascii", "text"}}, {"preview.mouse", u.Preview.Mouse, []string{"auto", "on", "off"}},
	} {
		valid := false
		for _, allowed := range field.allowed {
			if field.value == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid ui.%s %q", field.name, field.value)
		}
	}
	return nil
}

// Load chooses one conventional configuration. Ambiguity requires an explicit LoadFile selection.
func Load(root string) (Config, error) {
	present := []string{}
	for _, name := range []string{"mudarro.yaml", "mudarro.json"} {
		p, err := safePath(root, name)
		if err != nil {
			return Config{}, err
		}
		_, err = os.Stat(p)
		if err == nil {
			present = append(present, name)
		} else if !os.IsNotExist(err) {
			return Config{}, err
		}
	}
	if len(present) > 1 {
		return Config{}, fmt.Errorf("both mudarro.yaml and mudarro.json exist; select a configuration explicitly")
	}
	if len(present) == 0 {
		return LoadFile(root, "mudarro.yaml")
	}
	return LoadFile(root, present[0])
}

// LoadFile reads only a guarded project path, accepting .yaml/.yml or .json.
func LoadFile(root, configPath string) (Config, error) {
	var c Config
	if filepath.IsAbs(configPath) {
		absoluteRoot, err := filepath.Abs(root)
		if err != nil {
			return c, err
		}
		relative, err := filepath.Rel(absoluteRoot, configPath)
		if err != nil {
			return c, err
		}
		configPath = relative
	}
	path, err := safePath(root, configPath)
	if err != nil {
		return c, err
	}
	extension := strings.ToLower(filepath.Ext(path))
	if extension != ".yaml" && extension != ".yml" && extension != ".json" {
		return c, fmt.Errorf("unsupported configuration format: %s", extension)
	}
	// Check special files before opening: opening a FIFO could otherwise block indefinitely.
	initial, err := os.Stat(path)
	if err != nil {
		return c, err
	}
	if !initial.Mode().IsRegular() {
		return c, fmt.Errorf("configuration must be a regular file: %s", configPath)
	}
	f, err := os.Open(path)
	if err != nil {
		return c, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return c, err
	}
	if !info.Mode().IsRegular() {
		return c, fmt.Errorf("configuration must be a regular file: %s", configPath)
	}
	const limit = 4 << 20
	if info.Size() > limit {
		return c, fmt.Errorf("configuration exceeds 4 MiB: %s", configPath)
	}
	body, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return c, err
	}
	if len(body) > limit {
		return c, fmt.Errorf("configuration exceeds 4 MiB: %s", configPath)
	}
	if extension == ".json" {
		if err = checkJSONKeys(body); err != nil {
			return c, fmt.Errorf("%s: %w", configPath, err)
		}
		d := json.NewDecoder(bytes.NewReader(body))
		d.DisallowUnknownFields()
		if err = d.Decode(&c); err != nil {
			return c, fmt.Errorf("%s: %w", configPath, err)
		}
		var extra any
		if err = d.Decode(&extra); err != io.EOF {
			return c, fmt.Errorf("%s: expected one JSON document", configPath)
		}
	} else {
		d := yaml.NewDecoder(bytes.NewReader(body))
		d.KnownFields(true)
		if err = d.Decode(&c); err != nil {
			return c, fmt.Errorf("%s: %w", configPath, err)
		}
		var extra any
		if err = d.Decode(&extra); err != io.EOF {
			return c, fmt.Errorf("%s: expected one YAML document", configPath)
		}
	}
	if err := c.Validate(root); err != nil {
		return c, err
	}
	c.configPath = filepath.Clean(configPath)
	return c, nil
}

// The standard JSON decoder accepts duplicate object keys; reject them before decoding the schema.
func checkJSONKeys(body []byte) error {
	d := json.NewDecoder(bytes.NewReader(body))
	d.UseNumber()
	var value func(int) error
	value = func(depth int) error {
		if depth > 256 {
			return fmt.Errorf("JSON nesting exceeds 256 levels")
		}
		token, err := d.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			keys := map[string]bool{}
			for d.More() {
				key, err := d.Token()
				if err != nil {
					return err
				}
				name, ok := key.(string)
				if !ok {
					return fmt.Errorf("invalid JSON object key")
				}
				if keys[name] {
					return fmt.Errorf("duplicate JSON key %q", name)
				}
				keys[name] = true
				if err = value(depth + 1); err != nil {
					return err
				}
			}
		case '[':
			for d.More() {
				if err = value(depth + 1); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
		}
		_, err = d.Token()
		return err
	}
	if err := value(0); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return fmt.Errorf("expected one JSON document")
	}
	return nil
}
