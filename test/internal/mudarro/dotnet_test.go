package mudarro_test

import (
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"github.com/viralabs-dev/mudarro/internal/mudarro/adapters"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDotNetExplicitProjectsOffline(t *testing.T) {
	root := fixture(t, map[string]string{"A.csproj": `<Project Sdk="Microsoft.NET.Sdk.Web"><PropertyGroup><OutputType>Exe</OutputType></PropertyGroup><Target Name="Trap"><Exec Command="touch EXECUTED"/></Target></Project>`, "B.csproj": "<Project/>", "sample.sln": "unresolved external graph ../outside.csproj"})
	trap := t.TempDir()
	marker := filepath.Join(root, "EXECUTED")
	os.WriteFile(filepath.Join(trap, "dotnet"), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0755)
	t.Setenv("PATH", trap)
	before := supportFileHashes(t, root)
	s, su, e := (adapters.DotNet{}).Detect(root, ".", map[string]bool{"A.csproj": true, "B.csproj": true, "sample.sln": true}, nil)
	if e != nil || s == nil || s.Manager != "dotnet" || s.Framework != "" || len(s.Commands) != 0 || len(su) != 4 {
		t.Fatalf("%+v %+v %v", s, su, e)
	}
	if su[0].Name != "build-project-1" || !reflect.DeepEqual(su[0].Command.Args, []string{"dotnet", "build", "./A.csproj"}) || su[2].Name != "build-project-2" {
		t.Fatal(su)
	}
	if exists(marker) || !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("MSBuild executed or source mutated")
	}
}
func TestDotNetInvalidAndBoundedXML(t *testing.T) {
	for _, src := range []string{"", `<Project>`, `<project/>`, `<Project/><Project/>`, `<!DOCTYPE Project SYSTEM "https://example.invalid/x"><Project/>`, `<Project>&external;</Project>`, strings.Repeat("x", 4*1024*1024+1)} {
		root := fixture(t, map[string]string{"app.csproj": src})
		if _, _, e := (adapters.DotNet{}).Detect(root, ".", map[string]bool{"app.csproj": true}, nil); e == nil {
			t.Fatal("invalid accepted")
		}
	}
	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "app.csproj")
	os.WriteFile(target, []byte("<Project/>"), 0600)
	os.Symlink(target, filepath.Join(root, "app.csproj"))
	if _, _, e := (adapters.DotNet{}).Detect(root, ".", map[string]bool{"app.csproj": true}, nil); e == nil {
		t.Fatal("symlink accepted")
	}
}
func TestDotNetSolutionEvidenceOnly(t *testing.T) {
	root := fixture(t, map[string]string{"app.sln": "invalid executable-looking metadata", "modern.slnx": `<Solution><Project Path="../../foreign.csproj"/></Solution>`})
	s, su, e := (adapters.DotNet{}).Detect(root, ".", map[string]bool{"app.sln": true, "modern.slnx": true}, nil)
	if e != nil || s == nil || len(su) != 0 {
		t.Fatalf("%+v %v", s, e)
	}
	s, _, e = (adapters.DotNet{}).Detect(root, ".", map[string]bool{}, nil)
	if e != nil || s != nil {
		t.Fatal("absent detected")
	}
}

func TestDotNetUnusualNamesSafeArgv(t *testing.T) {
	root := fixture(t, map[string]string{"-flag.csproj": "<Project/>", "space +.csproj": "<Project/>"})
	_, su, e := (adapters.DotNet{}).Detect(root, ".", map[string]bool{"-flag.csproj": true, "space +.csproj": true}, nil)
	if e != nil || len(su) != 4 {
		t.Fatal(e)
	}
	for _, s := range su {
		if !strings.HasPrefix(s.Command.Args[2], "./") || !strings.Contains(s.Name, "project-") {
			t.Fatal(s)
		}
	}
}

func TestDotNetScanSelectionGeneratePreservesManual(t *testing.T) {
	root := fixture(t, map[string]string{"app.csproj": `<Project/>`, "excluded/app.csproj": strings.Repeat("x", 4*1024*1024+1)})
	report, e := Scan(root, []string{"excluded"})
	if e != nil || len(report.Config.Services) != 1 || len(report.Suggestions) != 2 || len(report.Warnings) != 0 {
		t.Fatalf("%+v %v", report, e)
	}
	var out bytes.Buffer
	if e = Main([]string{"init", "--root", root, "--exclude", "excluded", "--select", "app-csharp:build,app-csharp:test"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	c, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	if _, ok := c.Services[0].Commands["start"]; ok {
		t.Fatal("invented start")
	}
	if e = Generate(root, c, false, &out); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, root)
	if e = Generate(root, c, false, &out); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("generation not idempotent")
	}
	if e = Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("existing config overwritten")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("manual config mutated")
	}
}
