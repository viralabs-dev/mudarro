package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

// Composer reads static script declarations, never invoking hooks or plugins.
type Composer struct{}

func (Composer) Detect(root, dir string, files map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if !files["composer.json"] {
		return nil, nil, nil
	}
	data, e := read(root, dir, "composer.json")
	if e != nil {
		return nil, nil, e
	}
	if !utf8.Valid(data) {
		return nil, nil, fmt.Errorf("composer.json: invalid UTF-8")
	}
	if e = composerUniqueJSON(data); e != nil {
		return nil, nil, e
	}
	var manifest map[string]json.RawMessage
	if e = json.Unmarshal(data, &manifest); e != nil || manifest == nil {
		return nil, nil, fmt.Errorf("%s: expected Composer JSON object: %v", filepath.Join(dir, "composer.json"), e)
	}
	s := BaseService(dir, "php")
	s.Manager = "composer"
	s.Pending = append(s.Pending, "Declare commands.start explicitly; Composer metadata does not establish a framework or startup")
	var scripts map[string]json.RawMessage
	if raw, ok := manifest["scripts"]; ok {
		if string(raw) == "null" {
			return nil, nil, fmt.Errorf("%s: scripts must be an object", filepath.Join(dir, "composer.json"))
		}
		if e = json.Unmarshal(raw, &scripts); e != nil {
			return nil, nil, fmt.Errorf("%s: scripts must be an object: %w", filepath.Join(dir, "composer.json"), e)
		}
	}
	keys := []string{}
	names := map[string][]string{}
	for k := range scripts {
		keys = append(keys, k)
		n := strings.NewReplacer(":", "-", ".", "-").Replace(k)
		names[n] = append(names[n], k)
	}
	sort.Strings(keys)
	var suggestions []model.Suggestion
	for _, k := range keys {
		n := strings.NewReplacer(":", "-", ".", "-").Replace(k)
		if !identifier.MatchString(n) || strings.HasPrefix(k, "-") {
			s.Pending = append(s.Pending, fmt.Sprintf("Composer script %q requires explicit configuration: invalid action name", k))
			continue
		}
		if len(names[n]) != 1 {
			s.Pending = append(s.Pending, fmt.Sprintf("Composer script %q has a normalized-name collision: declare explicitly", k))
			continue
		}
		var body string
		if e = json.Unmarshal(scripts[k], &body); e != nil || strings.TrimSpace(body) == "" {
			s.Pending = append(s.Pending, fmt.Sprintf("Composer script %q has an unsupported or ambiguous body: declare explicitly", k))
			continue
		}
		suggestions = append(suggestions, model.Suggestion{Service: s.ID, Name: n, Purpose: "composer-script", Command: cmd("scripts", "composer", "run-script", "--", k), Evidence: filepath.ToSlash(filepath.Join(dir, "composer.json"))})
	}
	return s, suggestions, nil
}

// Reject duplicate declarations rather than silently choosing the last script.
func composerUniqueJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	var value func(int) error
	value = func(depth int) error {
		if depth > 128 {
			return fmt.Errorf("Composer JSON nesting exceeds 128")
		}
		token, e := d.Token()
		if e != nil {
			return e
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		if delim == '{' {
			seen := map[string]bool{}
			for d.More() {
				key, e := d.Token()
				if e != nil {
					return e
				}
				name, ok := key.(string)
				if !ok {
					return fmt.Errorf("invalid JSON object key")
				}
				if seen[name] {
					return fmt.Errorf("duplicate Composer JSON key %q", name)
				}
				seen[name] = true
				if e = value(depth + 1); e != nil {
					return e
				}
			}
		} else if delim == '[' {
			for d.More() {
				if e = value(depth + 1); e != nil {
					return e
				}
			}
		} else {
			return fmt.Errorf("unexpected JSON delimiter")
		}
		_, e = d.Token()
		return e
	}
	return value(0)
}
