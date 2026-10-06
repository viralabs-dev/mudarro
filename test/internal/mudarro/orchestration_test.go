package mudarro_test

import (
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func orchestrationSuggestions(r Report) []Suggestion {
	var out []Suggestion
	for _, s := range r.Suggestions {
		if strings.HasPrefix(s.Name, "turbo-") || strings.HasPrefix(s.Name, "nx-") {
			out = append(out, s)
		}
	}
	return out
}
func orchestrationFind(t *testing.T, r Report, owner, name string, args []string, evidence string) {
	t.Helper()
	var found []Suggestion
	for _, s := range r.Suggestions {
		if s.Service == owner && s.Name == name {
			found = append(found, s)
		}
	}
	if len(found) != 1 {
		t.Fatalf("%s/%s suggestions=%+v", owner, name, found)
	}
	if !reflect.DeepEqual(found[0].Command.Args, args) || found[0].Command.Shell != "" || found[0].Evidence != evidence {
		t.Errorf("unexpected declarative command: %+v", found[0])
	}
}
func orchestrationNoCommands(t *testing.T, r Report) {
	t.Helper()
	for _, s := range r.Config.Services {
		for name, c := range s.Commands {
			if strings.HasPrefix(name, "turbo-") || strings.HasPrefix(name, "nx-") || (len(c.Args) > 0 && (c.Args[0] == "turbo" || c.Args[0] == "nx")) {
				t.Errorf("orchestration automatically adopted: %s %+v", name, c)
			}
		}
	}
}
func TestOrchestrationTurboExplicitTasksAreRootOptIn(t *testing.T) {
	tools := t.TempDir()
	marker := filepath.Join(tools, "RUNTIME_MUST_NOT_RUN")
	for _, name := range []string{"turbo", "nx"} {
		if err := os.WriteFile(filepath.Join(tools, name), []byte("#!/bin/sh\nprintf invoked > '"+marker+"'\nexit 91\n"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", tools)
	t.Cleanup(func() {
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Error("scan invoked orchestration runtime")
		}
	})
	root := fixture(t, map[string]string{"package.json": `{"scripts":{"test":"touch NEVER"},"workspaces":["apps/*"]}`, "turbo.json": `{"tasks":{"build":{},"test":{"dependsOn":["^build"]},"dev":{"persistent":true,"cache":false}}}`, "apps/web/package.json": `{"scripts":{"build":"touch NEVER","test":"touch NEVER"}}`})
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	owner := jsWorkspaceService(t, r, ".")
	for _, task := range []string{"build", "test", "dev"} {
		orchestrationFind(t, r, owner.ID, "turbo-"+task, []string{"turbo", "run", task}, "turbo.json")
	}
	if len(orchestrationSuggestions(r)) != 3 {
		t.Errorf("unexpected tasks: %+v", orchestrationSuggestions(r))
	}
	orchestrationNoCommands(t, r)
	if _, ok := owner.Commands["start"]; ok {
		t.Error("persistent task promoted to startup")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Error("scan mutated fixture or executed tasks")
	}
}
func TestOrchestrationNxExplicitProjectTargets(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{}`, "nx.json": `{"plugins":["PLUGIN_MUST_NOT_EXECUTE"],"targetDefaults":{"inferred":{"executor":"never"}}}`, "apps/web/project.json": `{"name":"web","targets":{"build":{"executor":"never:build"},"test":{"command":"touch NEVER"},"empty":{},"inferred":{"dependsOn":["build"]}}}`, "apps/web/package.json": `{}`, "apps/no-name/project.json": `{"targets":{"build":{"command":"touch NEVER"}}}`})
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	owner := jsWorkspaceService(t, r, ".")
	for _, task := range []string{"build", "test"} {
		orchestrationFind(t, r, owner.ID, "nx-web-"+task, []string{"nx", "run", "web:" + task}, "apps/web/project.json")
	}
	if len(orchestrationSuggestions(r)) != 2 {
		t.Errorf("inferred or unnamed targets exposed: %+v", orchestrationSuggestions(r))
	}
	orchestrationNoCommands(t, r)
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Error("scan executed plugin/command or changed files")
	}
}
func TestOrchestrationNestedOwnersAndSharedTaskNames(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{}`, "turbo.json": `{"tasks":{"build":{}}}`, "nx.json": `{}`, "apps/a/project.json": `{"name":"a","targets":{"build":{"command":"never"}}}`, "nested/package.json": `{}`, "nested/turbo.json": `{"tasks":{"build":{}}}`, "nested/nx.json": `{}`, "nested/apps/b/project.json": `{"name":"a","targets":{"build":{"executor":"never"}}}`})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	outer := jsWorkspaceService(t, r, ".")
	inner := jsWorkspaceService(t, r, "nested")
	orchestrationFind(t, r, outer.ID, "turbo-build", []string{"turbo", "run", "build"}, "turbo.json")
	orchestrationFind(t, r, inner.ID, "turbo-build", []string{"turbo", "run", "build"}, "nested/turbo.json")
	orchestrationFind(t, r, outer.ID, "nx-a-build", []string{"nx", "run", "a:build"}, "apps/a/project.json")
	orchestrationFind(t, r, inner.ID, "nx-a-build", []string{"nx", "run", "a:build"}, "nested/apps/b/project.json")
	if len(orchestrationSuggestions(r)) != 4 {
		t.Errorf("wrong nested ownership: %+v", orchestrationSuggestions(r))
	}
	orchestrationNoCommands(t, r)
}
func TestOrchestrationMalformedAndMissingOwnerDiagnostic(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files map[string]string
	}{
		{"turbo-malformed", map[string]string{"package.json": `{}`, "turbo.json": `{"tasks":`}},
		{"nx-malformed", map[string]string{"package.json": `{}`, "nx.json": `{`, "apps/a/project.json": `{"name":"a","targets":{"build":{"command":"never"}}}`}},
		{"project-malformed", map[string]string{"package.json": `{}`, "nx.json": `{}`, "apps/a/project.json": `{"name":`}},
		{"turbo-no-owner", map[string]string{"turbo.json": `{"tasks":{"build":{}}}`}},
		{"nx-no-owner", map[string]string{"nx.json": `{}`, "apps/a/project.json": `{"name":"a","targets":{"build":{"command":"never"}}}`}},
		{"json-comments", map[string]string{"package.json": `{}`, "turbo.json": "{ // JSONC unsupported\n\"tasks\":{\"build\":{}}}"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, tc.files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			if len(orchestrationSuggestions(r)) != 0 {
				t.Error("unsafe/unowned tasks suggested")
			}
			if len(r.Warnings) == 0 {
				t.Error("missing diagnostic")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Error("scan changed files")
			}
		})
	}
}
func TestOrchestrationExcludedAndSymlinkSourcesNeverFollowed(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{}`, "nx.json": `{}`, "private/project.json": `{"name":"secret","targets":{"build":{"command":"never"}}}`, "allowed/project.json": `{"name":"ok","targets":{"test":{"command":"never"}}}`})
	outside := fixture(t, map[string]string{"turbo.json": `{"tasks":{"leak":{}}}`, "project.json": `{"name":"outside","targets":{"build":{"command":"never"}}}`})
	if e := os.Symlink(filepath.Join(outside, "turbo.json"), filepath.Join(root, "turbo.json")); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(filepath.Join(outside, "project.json"), filepath.Join(root, "linked-project.json")); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(outside, filepath.Join(root, "outside")); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, outside)
	r, e := Scan(root, []string{"private"})
	if e != nil {
		t.Fatal(e)
	}
	owner := jsWorkspaceService(t, r, ".")
	orchestrationFind(t, r, owner.ID, "nx-ok-test", []string{"nx", "run", "ok:test"}, "allowed/project.json")
	if len(orchestrationSuggestions(r)) != 1 {
		t.Errorf("excluded/outside source followed: %+v", orchestrationSuggestions(r))
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, outside)) {
		t.Error("outside files changed")
	}
}
func TestOrchestrationDuplicateNxNamesAndUnsafeIdentifiers(t *testing.T) {
	for _, tc := range []struct{ name, one, two string }{
		{"duplicate-project", `{"name":"same","targets":{"build":{"command":"never"}}}`, `{"name":"same","targets":{"build":{"executor":"never"}}}`},
		{"unsafe-project", `{"name":"bad;touch NEVER","targets":{"build":{"command":"never"}}}`, `{"name":"-flag","targets":{"build":{"command":"never"}}}`},
		{"unsafe-target", `{"name":"a","targets":{"bad task":{"command":"never"},"-flag":{"command":"never"}}}`, `{"name":"b","targets":{"x;touch NEVER":{"command":"never"}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": `{}`, "nx.json": `{}`, "a/project.json": tc.one, "b/project.json": tc.two})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			if len(orchestrationSuggestions(r)) != 0 {
				t.Errorf("ambiguous/unsafe command suggested: %+v", orchestrationSuggestions(r))
			}
			if len(r.Warnings) == 0 {
				t.Error("missing collision/identifier diagnostic")
			}
		})
	}
}
func TestOrchestrationSuggestionCollisionIsDiagnosed(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"scripts":{"turbo-build":"never"}}`, "turbo.json": `{"tasks":{"build":{}}}`})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	owner := jsWorkspaceService(t, r, ".")
	count := 0
	for _, s := range r.Suggestions {
		if s.Service == owner.ID && s.Name == "turbo-build" {
			count++
		}
	}
	if count > 1 {
		t.Error("duplicate selectable suggestion IDs")
	}
	if len(r.Warnings) == 0 {
		t.Error("collision missing diagnostic")
	}
}

func TestOrchestrationNormalizedSuggestionCollisionsBlockBoth(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files map[string]string
	}{
		{"turbo-tasks", map[string]string{"package.json": `{}`, "turbo.json": `{"tasks":{"build/web":{},"build-web":{}}}`}},
		{"nx-projects", map[string]string{"package.json": `{}`, "nx.json": `{}`, "a/project.json": `{"name":"web/app","targets":{"build":{"command":"never"}}}`, "b/project.json": `{"name":"web-app","targets":{"build":{"executor":"never"}}}`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, tc.files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			if len(orchestrationSuggestions(r)) != 0 {
				t.Errorf("normalized collisions remain selectable: %+v", orchestrationSuggestions(r))
			}
			if !strings.Contains(strings.ToLower(strings.Join(r.Warnings, " ")), "collid") {
				t.Errorf("missing collision diagnostic: %v", r.Warnings)
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Error("scan changed declarations")
			}
		})
	}
}

func TestOrchestrationInvalidTurboDoesNotHideValidNxTasks(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{}`, "turbo.json": `{"tasks":`, "nx.json": `{}`, "apps/web/project.json": `{"name":"web","targets":{"test":{"command":"never"}}}`})
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	owner := jsWorkspaceService(t, r, ".")
	orchestrationFind(t, r, owner.ID, "nx-web-test", []string{"nx", "run", "web:test"}, "apps/web/project.json")
	if len(orchestrationSuggestions(r)) != 1 {
		t.Errorf("invalid Turbo contaminated Nx: %+v", orchestrationSuggestions(r))
	}
	if !strings.Contains(strings.Join(r.Warnings, " "), "turbo.json") {
		t.Errorf("missing invalid Turbo diagnostic: %v", r.Warnings)
	}
	orchestrationNoCommands(t, r)
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Error("scan changed declarations")
	}
}
