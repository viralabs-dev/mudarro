package mudarro_test

import (
	"bytes"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func jsWorkspaceService(t *testing.T, r Report, dir string) Service {
	t.Helper()
	for _, s := range r.Config.Services {
		if s.Dir == dir {
			return s
		}
	}
	t.Fatalf("service missing %q: %+v", dir, r.Config.Services)
	return Service{}
}
func jsWorkspaceRoot(t *testing.T, s Service) string {
	t.Helper()
	b, e := json.Marshal(s)
	if e != nil {
		t.Fatal(e)
	}
	var value struct {
		Root string `json:"workspace_root"`
	}
	if e = json.Unmarshal(b, &value); e != nil {
		t.Fatal(e)
	}
	return value.Root
}
func jsWorkspaceDescriptor(t *testing.T, r Report, root string) ([]string, string, bool) {
	t.Helper()
	b, e := json.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	var value struct {
		Workspaces []struct {
			Path    string   `json:"path"`
			Root    string   `json:"root"`
			Members []string `json:"members"`
			Manager string   `json:"manager"`
		} `json:"javascript_workspaces"`
	}
	if e = json.Unmarshal(b, &value); e != nil {
		t.Fatal(e)
	}
	for _, w := range value.Workspaces {
		if w.Root == root {
			return w.Members, w.Manager, true
		}
	}
	return nil, "", false
}
func jsWorkspaceDiagnostics(r Report, s Service) string {
	return strings.Join(r.Warnings, " ") + " " + strings.Join(s.Pending, " ")
}

func TestJavaScriptWorkspaceManagerOnlyInheritanceAndRootInstall(t *testing.T) {
	for _, decl := range []string{`["packages/*"]`, `{"packages":["packages/*"]}`} {
		t.Run(decl, func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": `{"name":"root","packageManager":"pnpm@10.0.0","workspaces":` + decl + `,"scripts":{"start":"touch ROOT_MUST_NOT_RUN","test":"touch ROOT_MUST_NOT_RUN","root-only":"touch ROOT_MUST_NOT_RUN"}}`, "compose.yaml": "services: {}\n", "prisma/schema.prisma": `datasource db { provider = "sqlite" }`, "packages/api/package.json": `{"scripts":{"test":"touch CHILD_MUST_NOT_RUN","hello":"touch CHILD_MUST_NOT_RUN"}}`, "outside/package.json": `{"scripts":{"test":"touch OUTSIDE_MUST_NOT_RUN"}}`})
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			child := jsWorkspaceService(t, r, "packages/api")
			owner := jsWorkspaceService(t, r, ".")
			outside := jsWorkspaceService(t, r, "outside")
			if child.Manager != "pnpm" || !reflect.DeepEqual(child.Commands["test"].Args, []string{"pnpm", "run", "test"}) {
				t.Errorf("manager inheritance missing: %+v", child)
			}
			if jsWorkspaceRoot(t, child) != "." {
				t.Error("child workspace owner absent")
			}
			if _, ok := child.Commands["install"]; ok {
				t.Error("member install planned in child cwd")
			}
			if !reflect.DeepEqual(owner.Commands["install"].Args, []string{"pnpm", "install"}) {
				t.Error("root install changed")
			}
			if _, ok := child.Commands["start"]; ok {
				t.Error("root startup inherited")
			}
			if child.Infrastructure.Kind != "local" || child.Database.Tool != "" {
				t.Error("root infrastructure/database inherited")
			}
			if outside.Manager != "npm" || jsWorkspaceRoot(t, outside) != "" {
				t.Errorf("nonmember contaminated: %+v", outside)
			}
			for _, su := range r.Suggestions {
				if su.Service == child.ID {
					if su.Name == "root-only" || su.Evidence != "packages/api/package.json" {
						t.Errorf("root script/evidence inherited: %+v", su)
					}
				}
			}
			members, manager, ok := jsWorkspaceDescriptor(t, r, ".")
			if !ok || manager != "pnpm" || !reflect.DeepEqual(members, []string{"packages/api"}) {
				t.Errorf("descriptor: %v %q %v", members, manager, ok)
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("scan changed files or executed scripts")
			}
		})
	}
}

func TestJavaScriptWorkspaceGlobsNegationsAndExcludes(t *testing.T) {
	for _, tc := range []struct {
		name, patterns string
		excludes       []string
		members        []string
	}{
		{"single-star", `["packages/*"]`, nil, []string{"packages/a", "packages/b", "packages/drop", "packages/long1"}},
		{"question-and-class", `["packages/[ab]","packages/long?"]`, nil, []string{"packages/a", "packages/b", "packages/long1"}},
		{"double-star-negation", `["packages/**","!packages/drop","!packages/drop/**"]`, nil, []string{"packages/a", "packages/b", "packages/long1", "packages/nested/deep"}},
		{"scan-exclusion", `["packages/**"]`, []string{"packages/drop", "packages/nested"}, []string{"packages/a", "packages/b", "packages/long1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"package.json": `{"packageManager":"bun@1.4.0","workspaces":` + tc.patterns + `}`}
			for _, dir := range []string{"packages/a", "packages/b", "packages/drop", "packages/drop/deep", "packages/long1", "packages/nested/deep"} {
				files[dir+"/package.json"] = `{"scripts":{"test":"NEVER_EXECUTE"}}`
			}
			root := fixture(t, files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, tc.excludes)
			if e != nil {
				t.Fatal(e)
			}
			members, _, ok := jsWorkspaceDescriptor(t, r, ".")
			if !ok || !reflect.DeepEqual(members, tc.members) {
				t.Errorf("members=%v want=%v present=%v", members, tc.members, ok)
			}
			for _, dir := range tc.members {
				s := jsWorkspaceService(t, r, dir)
				if s.Manager != "bun" || jsWorkspaceRoot(t, s) != "." {
					t.Errorf("matched member manager/owner: %+v", s)
				}
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("workspace scan changed source")
			}
		})
	}
}

func TestJavaScriptWorkspaceManagerAndVersionConflicts(t *testing.T) {
	for _, tc := range []struct {
		name, child, lock string
		conflict          bool
	}{
		{"same-manager", `{"packageManager":"pnpm@10.0.0","scripts":{"test":"NEVER_EXECUTE"}}`, "", false},
		{"different-manager", `{"packageManager":"yarn@4.18.1","scripts":{"test":"NEVER_EXECUTE"}}`, "", true},
		{"different-lock", `{"scripts":{"test":"NEVER_EXECUTE"}}`, "yarn.lock", true},
		{"different-exact-version", `{"packageManager":"pnpm@9.0.0","scripts":{"test":"NEVER_EXECUTE"}}`, "", true},
		{"same-manager-lock", `{"scripts":{"test":"NEVER_EXECUTE"}}`, "pnpm-lock.yaml", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{"package.json": `{"packageManager":"pnpm@10.0.0","workspaces":["packages/*"]}`, "packages/api/package.json": tc.child}
			if tc.lock != "" {
				files["packages/api/"+tc.lock] = "fixture"
			}
			r, e := Scan(fixture(t, files), nil)
			if e != nil {
				t.Fatal(e)
			}
			s := jsWorkspaceService(t, r, "packages/api")
			if tc.conflict {
				if s.Manager != "" || len(s.Commands) > 0 {
					t.Errorf("conflicting member executable plan: %+v", s)
				}
				if len(s.Pending) == 0 {
					t.Error("conflict not pending")
				}
				for _, su := range r.Suggestions {
					if su.Service == s.ID && len(su.Command.Args) > 0 && su.Command.Args[0] != "" {
						t.Error("conflicting member executable suggestion")
					}
				}
			} else if s.Manager != "pnpm" || len(s.Commands["test"].Args) == 0 {
				t.Errorf("compatible member unresolved: %+v", s)
			}
		})
	}
}

func TestJavaScriptWorkspaceNestedBoundaryUsesNearestDeclaration(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"packageManager":"pnpm@10.0.0","workspaces":["packages/**"]}`, "packages/nested/package.json": `{"packageManager":"bun@1.4.0","workspaces":["children/*"],"scripts":{"test":"NEVER_EXECUTE"}}`, "packages/nested/children/app/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`, "packages/plain/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	nested := jsWorkspaceService(t, r, "packages/nested")
	child := jsWorkspaceService(t, r, "packages/nested/children/app")
	if nested.Manager != "" || len(nested.Pending) == 0 {
		t.Errorf("overlapping nested root conflict hidden: %+v", nested)
	}
	if jsWorkspaceRoot(t, child) != "packages/nested" || child.Manager != "bun" {
		t.Errorf("outer workspace took nested child: %+v", child)
	}
	members, manager, ok := jsWorkspaceDescriptor(t, r, "packages/nested")
	if !ok || manager != "bun" || !reflect.DeepEqual(members, []string{"packages/nested/children/app"}) {
		t.Errorf("nested descriptor=%v %q %v", members, manager, ok)
	}
	t.Run("unresolved-nested-owner-is-not-overwritten", func(t *testing.T) {
		root := fixture(t, map[string]string{"package.json": `{"packageManager":"pnpm@10.0.0","workspaces":["packages/*"]}`, "packages/inner/package.json": `{"packageManager":"pnpm@10.0.0","workspaces":["children/*"],"scripts":{"hello":"NEVER_EXECUTE"}}`, "packages/inner/pnpm-workspace.yaml": "packages:\n - 'other/*'\n", "packages/inner/children/app/package.json": `{"scripts":{"test":"NEVER_EXECUTE","hello":"NEVER_EXECUTE"}}`, "packages/inner/other/app/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
		r, e := Scan(root, nil)
		if e != nil {
			t.Fatal(e)
		}
		for _, dir := range []string{"packages/inner", "packages/inner/children/app", "packages/inner/other/app"} {
			svc := jsWorkspaceService(t, r, dir)
			if svc.Manager != "" || len(svc.Commands) > 0 || len(svc.Pending) == 0 {
				t.Errorf("outer context overwrote unresolved nested owner: %+v", svc)
			}
			for _, su := range r.Suggestions {
				if su.Service == svc.ID && len(su.Command.Args) > 0 && su.Command.Args[0] != "" {
					t.Errorf("unresolved nested executable suggestion: %+v", su)
				}
			}
		}
	})

}

func TestJavaScriptWorkspaceSymlinksAndScanStayOffline(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"packageManager":"bun@1.4.0","workspaces":["packages/*"]}`, "packages/app/package.json": `{"scripts":{"test":"touch NEVER_EXECUTE"}}`})
	outside := fixture(t, map[string]string{"package.json": `{"scripts":{"start":"NEVER_EXECUTE"}}`})
	outsideBefore := supportFileHashes(t, outside)
	if e := os.Symlink(outside, filepath.Join(root, "packages", "linked")); e != nil {
		t.Fatal(e)
	}
	tools := t.TempDir()
	for _, tool := range []string{"node", "npm", "pnpm", "yarn", "bun"} {
		if e := os.WriteFile(filepath.Join(tools, tool), []byte("#!/bin/sh\nprintf executed > '"+filepath.Join(root, "MUST_NOT_RUN")+"'\n"), 0755); e != nil {
			t.Fatal(e)
		}
	}
	t.Setenv("PATH", tools)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, _, ok := jsWorkspaceDescriptor(t, r, ".")
	if !ok || !reflect.DeepEqual(members, []string{"packages/app"}) {
		t.Errorf("symlink admitted: %v %v", members, ok)
	}
	if len(r.Config.Services) != 2 {
		t.Error("scanner followed linked workspace")
	}
	if _, e := os.Stat(filepath.Join(root, "MUST_NOT_RUN")); !os.IsNotExist(e) {
		t.Fatal("scanner executed runtime")
	}
	if !reflect.DeepEqual(outsideBefore, supportFileHashes(t, outside)) {
		t.Fatal("outside fixture changed")
	}
}

func TestJavaScriptWorkspaceUnsupportedPatternsAreDiagnosed(t *testing.T) {
	for _, pattern := range []string{"packages/{a,b}", "packages/@(a|b)", "../outside", "/outside"} {
		t.Run(pattern, func(t *testing.T) {
			patterns, _ := json.Marshal([]string{pattern})
			root := fixture(t, map[string]string{"package.json": `{"packageManager":"bun@1.4.0","workspaces":` + string(patterns) + `}`, "packages/a/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			owner := jsWorkspaceService(t, r, ".")
			if !strings.Contains(strings.ToLower(jsWorkspaceDiagnostics(r, owner)), "workspace") {
				t.Error("unsupported workspace pattern silent")
			}
			child := jsWorkspaceService(t, r, "packages/a")
			if jsWorkspaceRoot(t, child) != "" || child.Manager != "npm" {
				t.Error("unsupported pattern silently applied")
			}
		})
	}
}

func TestJavaScriptWorkspacePnpmManifestInheritance(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"packageManager":"pnpm@10.0.0"}`, "pnpm-workspace.yaml": "packages:\n - 'packages/*'\n - '!packages/drop'\n", "packages/api/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`, "packages/drop/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
	before := supportFileHashes(t, root)
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	members, manager, ok := jsWorkspaceDescriptor(t, r, ".")
	if !ok || manager != "pnpm" || !reflect.DeepEqual(members, []string{"packages/api"}) {
		t.Errorf("pnpm descriptor=%v %q %v", members, manager, ok)
	}
	s := jsWorkspaceService(t, r, "packages/api")
	if s.Manager != "pnpm" || jsWorkspaceRoot(t, s) != "." {
		t.Errorf("pnpm inherited manager absent: %+v", s)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("pnpm scan changed fixture")
	}
}

func TestJavaScriptWorkspaceExistingManualConfigurationIsPreserved(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{"packageManager":"pnpm@10.0.0","workspaces":["packages/*"]}`, "packages/api/package.json": `{"packageManager":"yarn@4.18.1","scripts":{"test":"NEVER_EXECUTE"}}`})
	c := Config{Version: 1, Name: "manual consumer", Services: []Service{{ID: "manual", Dir: "packages/api", Language: "javascript", Manager: "custom-tool", Infrastructure: Infrastructure{Kind: "custom"}, Commands: map[string]Command{"test": {Args: []string{"custom-tool", "explicit-test"}, Group: "qualidade"}}}}}
	if e := saveConfig(filepath.Join(root, "mudarro.yaml"), c); e != nil {
		t.Fatal(e)
	}
	var out bytes.Buffer
	if e := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, root)
	if _, e := Scan(root, nil); e != nil {
		t.Fatal(e)
	}
	if e := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("workspace scan/regeneration changed manual configuration or wrappers")
	}
	loaded, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	if loaded.Services[0].Manager != "custom-tool" || !reflect.DeepEqual(loaded.Services[0].Commands["test"].Args, []string{"custom-tool", "explicit-test"}) {
		t.Fatal("explicit configured command was globally blocked/rewritten")
	}
}

func TestJavaScriptWorkspaceDualDeclarationsEquivalentOrAmbiguous(t *testing.T) {
	for _, same := range []bool{true, false} {
		t.Run(map[bool]string{true: "equivalent-reordered", false: "divergent-union-blocked"}[same], func(t *testing.T) {
			yaml := "packages:\n - '!packages/drop'\n - 'packages/*'\n"
			jsonPatterns := `["packages/*","!packages/drop"]`
			if !same {
				yaml = "packages:\n - 'packages/other'\n"
				jsonPatterns = `["packages/api"]`
			}
			root := fixture(t, map[string]string{"package.json": `{"packageManager":"pnpm@10.0.0","workspaces":` + jsonPatterns + `}`, "pnpm-workspace.yaml": yaml, "packages/api/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`, "packages/other/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`, "packages/drop/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			for _, dir := range []string{"packages/api", "packages/other"} {
				s := jsWorkspaceService(t, r, dir)
				if same {
					if s.Manager != "pnpm" || jsWorkspaceRoot(t, s) != "." {
						t.Errorf("equivalent declaration failed ownership: %+v", s)
					}
				} else {
					if s.Manager != "" || len(s.Commands) > 0 || len(s.Pending) == 0 {
						t.Errorf("ambiguous workspace planned fallback execution: %+v", s)
					}
				}
			}
			members, manager, ok := jsWorkspaceDescriptor(t, r, ".")
			if !ok || !reflect.DeepEqual(members, []string{"packages/api", "packages/other"}) {
				t.Errorf("dual membership=%v present=%v", members, ok)
			}
			if same && manager != "pnpm" {
				t.Error("equivalent declarations lost manager")
			}
			if !same && manager != "" {
				t.Error("ambiguous declaration selected manager")
			}
		})
	}
}

func TestJavaScriptWorkspacePnpmDeclarationWithoutManagerInfersRootManager(t *testing.T) {
	root := fixture(t, map[string]string{"package.json": `{}`, "pnpm-workspace.yaml": "packages:\n - 'packages/*'\n", "packages/api/package.json": `{"scripts":{"test":"NEVER_EXECUTE"}}`})
	r, e := Scan(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	for _, dir := range []string{".", "packages/api"} {
		s := jsWorkspaceService(t, r, dir)
		if s.Manager != "pnpm" {
			t.Errorf("pnpm workspace without PM fell back npm: %+v", s)
		}
	}
	if !reflect.DeepEqual(jsWorkspaceService(t, r, ".").Commands["install"].Args, []string{"pnpm", "install"}) {
		t.Error("root install did not use pnpm")
	}
}

func TestJavaScriptWorkspaceUnresolvedRootDoesNotFallbackMemberManager(t *testing.T) {
	for _, kind := range []string{"invalid-packageManager", "ambiguous-locks"} {
		t.Run(kind, func(t *testing.T) {
			rootBody := `{"packageManager":"pnpm@latest","workspaces":["packages/*"]}`
			files := map[string]string{"package.json": rootBody, "packages/api/package.json": `{"scripts":{"test":"NEVER_EXECUTE","hello":"NEVER_EXECUTE"}}`}
			if kind == "ambiguous-locks" {
				files["package.json"] = `{"workspaces":["packages/*"]}`
				files["pnpm-lock.yaml"] = "fixture"
				files["yarn.lock"] = "fixture"
			}
			root := fixture(t, files)
			before := supportFileHashes(t, root)
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			s := jsWorkspaceService(t, r, "packages/api")
			if s.Manager != "" || len(s.Commands) > 0 || len(s.Pending) == 0 || jsWorkspaceRoot(t, s) != "." {
				t.Errorf("unresolved root lost ownership or planned fallback: %+v", s)
			}
			for _, su := range r.Suggestions {
				if su.Service == s.ID && len(su.Command.Args) > 0 && su.Command.Args[0] != "" {
					t.Error("unresolved root executable suggestion")
				}
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("conflict scan changed source")
			}
		})
	}
}
