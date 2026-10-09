package adapters

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/viralabs-dev/mudarro/internal/mudarro/model"
	"github.com/viralabs-dev/mudarro/internal/mudarro/projectfs"
	"io"
	"path"
	"path/filepath"
	"strings"
)

// Java identifies Maven XML; Gradle source is evidence only, never evaluated.
type Java struct{}

// declarativeXML rejects directives (including DTDs) and malformed or multiple
// roots. encoding/xml has no external entity resolver; custom entities fail.
func declarativeXML(data []byte, expected string) error {
	d := xml.NewDecoder(bytes.NewReader(data))
	depth, roots := 0, 0
	for {
		t, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		switch v := t.(type) {
		case xml.Directive:
			return fmt.Errorf("XML directives/DTD are not supported")
		case xml.StartElement:
			if depth == 0 {
				roots++
				if v.Name.Local != expected {
					return fmt.Errorf("expected XML root %s", expected)
				}
			}
			depth++
		case xml.EndElement:
			depth--
		case xml.CharData:
			if depth == 0 && strings.TrimSpace(string(v)) != "" {
				return fmt.Errorf("text outside XML root")
			}
		}
	}
	if roots != 1 || depth != 0 {
		return fmt.Errorf("expected one complete XML root")
	}
	return nil
}
func declarativeExcluded(rel string, excludes []string) bool {
	for _, part := range strings.Split(filepath.ToSlash(rel), "/") {
		if IgnoredDirectory(part) {
			return true
		}
	}
	for _, pat := range excludes {
		match, _ := path.Match(pat, filepath.ToSlash(rel))
		if match || rel == pat || strings.HasPrefix(filepath.ToSlash(rel), strings.TrimSuffix(pat, "/")+"/") {
			return true
		}
	}
	return false
}
func (Java) Detect(root, dir string, files map[string]bool, excludes []string) (*model.Service, []model.Suggestion, error) {
	if !files["pom.xml"] {
		for _, name := range []string{"build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"} {
			if files[name] {
				if _, e := read(root, dir, name); e != nil {
					return nil, nil, e
				}
				s := BaseService(dir, "jvm")
				s.Manager = "gradle"
				s.Pending = []string{"Gradle DSL is executable source: modules, plugins, language and tasks were not evaluated; declare commands explicitly"}
				return s, nil, nil
			}
		}
		return nil, nil, nil
	}
	b, e := read(root, dir, "pom.xml")
	if e != nil {
		return nil, nil, e
	}
	if e = declarativeXML(b, "project"); e != nil {
		return nil, nil, e
	}
	var pom struct {
		Modules   []string `xml:"modules>module"`
		Packaging string   `xml:"packaging"`
	}
	if e = xml.Unmarshal(b, &pom); e != nil {
		return nil, nil, e
	}
	s := BaseService(dir, "java")
	s.Manager = "maven"
	s.Pending = []string{"Maven profiles, plugins, dependencies and properties were not evaluated; declare main/start explicitly"}
	if len(pom.Modules) > 0 || strings.TrimSpace(pom.Packaging) == "pom" {
		for _, member := range pom.Modules {
			member = strings.TrimSpace(member)
			rel := filepath.Join(dir, member)
			if member == "" || filepath.IsAbs(member) || strings.Contains(member, "$") || strings.Contains(member, "\\") {
				return nil, nil, fmt.Errorf("unsupported Maven module %q", member)
			}
			if _, e = projectfs.SafePath(root, rel); e != nil {
				return nil, nil, e
			}
			if declarativeExcluded(rel, excludes) || declarativeExcluded(filepath.Join(rel, "pom.xml"), excludes) {
				s.Pending = append(s.Pending, "Excluded Maven module: "+member)
				continue
			}
			mb, err := read(root, rel, "pom.xml")
			if err != nil {
				return nil, nil, err
			}
			if err = declarativeXML(mb, "project"); err != nil {
				return nil, nil, err
			}
		}
		s.Pending = append(s.Pending, "Maven aggregator: select independently detected contained modules explicitly; no aggregate command inferred")
		return s, nil, nil
	}
	evidence := filepath.ToSlash(filepath.Join(dir, "pom.xml"))
	return s, []model.Suggestion{{Service: s.ID, Name: "compile", Purpose: "maven", Command: cmd("aplicacao", "mvn", "-f", "pom.xml", "compile"), Evidence: evidence}, {Service: s.ID, Name: "test", Purpose: "maven", Command: cmd("qualidade", "mvn", "-f", "pom.xml", "test"), Evidence: evidence}}, nil
}
