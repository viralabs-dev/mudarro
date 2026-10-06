from pathlib import Path
r=Path('/home/danielsouza/dev/viralabs/mudarro');p=r/'internal/mudarro/adapters/package_manager.go'
p.write_text('''package adapters

import (
 "fmt"
 "regexp"
 "strings"
 "unicode"
)

// Exact SemVer 2.0 syntax; build metadata (including Corepack hash text) is
// accepted syntactically, not verified or resolved by scanning.
var exactManagerVersion = regexp.MustCompile(`^(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$`)

func packageManagerName(value string) (string, error) {
 name, version, versioned := strings.Cut(value, "@")
 known := name == "npm" || name == "pnpm" || name == "yarn" || name == "bun"
 valid := known && strings.IndexFunc(value, unicode.IsSpace) < 0
 if versioned {
  valid = valid && exactManagerVersion.MatchString(version)
  core, _, _ := strings.Cut(version, "+")
  _, prerelease, hasPrerelease := strings.Cut(core, "-")
  if hasPrerelease {
   for _, part := range strings.Split(prerelease, ".") {
    numeric := part != "" && strings.IndexFunc(part, func(r rune) bool { return r < '0' || r > '9' }) < 0
    if numeric && len(part) > 1 && part[0] == '0' { valid = false }
   }
  }
 }
 if !valid {
  return "", fmt.Errorf("invalid packageManager %q: use npm, pnpm, yarn or bun, optionally with an exact SemVer version; custom tools require explicit commands", value)
 }
 return name, nil
}
''')
p=r/'internal/mudarro/model/types.go';t=p.read_text();t=t.replace('Manager        string             `yaml:"manager,omitempty" json:"manager,omitempty"`','Manager        string             `yaml:"manager,omitempty" json:"manager,omitempty"`\n GoWorkspace string `yaml:"go_workspace,omitempty" json:"go_workspace,omitempty"`');t=t.replace('type Suggestion struct {','type Suggestion struct {\n Purpose string `json:"purpose,omitempty"`');p.write_text(t)
p=r/'internal/mudarro/config.go';t=p.read_text();needle='\t\tseen[s.ID] = true\n';assert needle in t;t=t.replace(needle,needle+'''  if s.GoWorkspace != "" && (s.Language != "go" || (s.GoWorkspace != "inherit" && s.GoWorkspace != "off")) {
   return fmt.Errorf("%s: go_workspace requires a Go service and inherit or off", s.ID)
  }
''');p.write_text(t)
p=r/'internal/mudarro/adapters/language_adapters.go';t=p.read_text();old='''	if len(entries) == 1 {
		for p := range entries {
			s.Commands["start"] = cmd("aplicacao", "go", "run", "./"+filepath.ToSlash(p))
		}
	} else {
		s.Pending = append(s.Pending, "Declare commands.start: nenhuma entrada Go única")
	}
	return s, nil, nil
}''';new=''' targets := make([]string, 0, len(entries))
 for p := range entries { targets = append(targets, p) }
 sort.Strings(targets)
 suggestions := make([]model.Suggestion, 0, len(targets))
 for i, p := range targets {
  c := cmd("aplicacao", "go", "run", "./"+filepath.ToSlash(p))
  suggestions = append(suggestions, model.Suggestion{Name: fmt.Sprintf("go-entry-%d", i+1), Purpose: "go-start", Command: c, Evidence: filepath.Join(dir, p)})
  if len(targets) == 1 { s.Commands["start"] = c }
 }
 if len(targets) != 1 {
  s.Pending = append(s.Pending, "Declare commands.start: nenhuma entrada Go única")
 }
 return s, suggestions, nil
}''';assert old in t;t=t.replace(old,new);p.write_text(t)
p=r/'internal/mudarro/cli.go';t=p.read_text();t=t.replace('\t\tavailable := map[string]bool{}','\t\tavailable := map[string]bool{}\n  selectedGo := map[string]int{}');old='''						report.Config.Services[j].Commands[su.Name] = su.Command''';new='''      s := &report.Config.Services[j]
      if su.Purpose == "go-start" {
       selectedGo[s.ID]++
       if selectedGo[s.ID] > 1 { return fmt.Errorf("multiple Go entries selected for %s; select one start target", s.ID) }
       if _, exists := s.Commands["start"]; !exists { s.Commands["start"] = su.Command }
       pending := s.Pending[:0]
       for _, item := range s.Pending { if item != "Declare commands.start: nenhuma entrada Go única" { pending = append(pending, item) } }
       s.Pending = pending
      } else {
       s.Commands[su.Name] = su.Command
      }''';assert old in t;t=t.replace(old,new);p.write_text(t)
p=r/'internal/mudarro/runner.go';t=p.read_text();t=t.replace('\treturn r.execute(dir, args)','\treturn r.execute(dir, goWorkspaceArgs(s, args))');t+='''
// Apply an explicit workspace override only to this child's command. The
// default inherits Go's workspace discovery and any caller GOWORK value.
func goWorkspaceArgs(s Service, args []string) []string {
 if s.Language == "go" && s.GoWorkspace == "off" {
  return append([]string{"env", "GOWORK=off"}, args...)
 }
 return args
}
''';p.write_text(t)
p=r/'internal/mudarro/local_runner.go';t=p.read_text();needle='\tchild := exec.Command(v[0], v[1:]...)';assert needle in t;t=t.replace(needle,'\tv = goWorkspaceArgs(s, v)\n'+needle);p.write_text(t)
print('035 exact SemVer and038 explicit target/context integrated; no workspace parser added.')
