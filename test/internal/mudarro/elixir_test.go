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

func TestElixirEvidenceOnlyOffline(t *testing.T) {
	for _, source := range []string{"defmodule Project do\n use Mix.Project\nend\n", "INVALID executable syntax", string([]byte{0xff, 0xfe})} {
		root := fixture(t, map[string]string{"mix.exs": source})
		trapdir := t.TempDir()
		marker := filepath.Join(root, "EXECUTED")
		os.WriteFile(filepath.Join(trapdir, "mix"), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0755)
		t.Setenv("PATH", trapdir)
		before := supportFileHashes(t, root)
		s, su, e := (adapters.Elixir{}).Detect(root, ".", map[string]bool{"mix.exs": true}, nil)
		if e != nil {
			t.Fatal(e)
		}
		if s.Language != "elixir" || s.Manager != "mix" || len(s.Commands) != 0 || len(s.Pending) < 3 || len(su) != 2 {
			t.Fatalf("service %+v suggestions %+v", s, su)
		}
		for i, name := range []string{"compile", "test"} {
			if su[i].Name != name || su[i].Purpose != "mix" || su[i].Evidence != "mix.exs" || !reflect.DeepEqual(su[i].Command.Args, []string{"mix", name}) {
				t.Fatal(su[i])
			}
		}
		if exists(marker) || !reflect.DeepEqual(before, supportFileHashes(t, root)) {
			t.Fatal("detector executed or modified source")
		}
	}
}
func TestElixirContainedRegularBoundedEvidence(t *testing.T) {
	for _, kind := range []string{"oversized", "symlink", "directory"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			p := filepath.Join(root, "mix.exs")
			switch kind {
			case "oversized":
				os.WriteFile(p, bytes.Repeat([]byte("x"), 4*1024*1024+1), 0600)
			case "symlink":
				target := filepath.Join(t.TempDir(), "source")
				os.WriteFile(target, []byte("fixture"), 0600)
				if e := os.Symlink(target, p); e != nil {
					t.Skip(e)
				}
			case "directory":
				os.Mkdir(p, 0700)
			}
			if _, _, e := (adapters.Elixir{}).Detect(root, ".", map[string]bool{"mix.exs": true}, nil); e == nil {
				t.Fatal("invalid evidence accepted")
			}
		})
	}
	s, su, e := (adapters.Elixir{}).Detect(t.TempDir(), ".", map[string]bool{}, nil)
	if e != nil || s != nil || su != nil {
		t.Fatal("manifest absent was detected")
	}
}
func TestElixirScanNestedExcludedAndOptIn(t *testing.T) {
	root := fixture(t, map[string]string{"mix.exs": "root executable evidence", "apps/member/mix.exs": "member evidence", "excluded/mix.exs": strings.Repeat("x", 4*1024*1024+1)})
	report, e := Scan(root, []string{"excluded"})
	if e != nil {
		t.Fatal(e)
	}
	if len(report.Config.Services) != 2 || len(report.Suggestions) != 4 {
		t.Fatalf("report %+v", report)
	}
	for _, s := range report.Config.Services {
		if s.Language != "elixir" || len(s.Commands) != 0 {
			t.Fatalf("automatic command %+v", s)
		}
	}
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--exclude", "excluded", "--select", "app-elixir:compile,apps-member-elixir:test"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	c, e := Load(root)
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range c.Services {
		if _, ok := s.Commands["start"]; ok {
			t.Fatal("invented start")
		}
		if s.Dir == "." && !reflect.DeepEqual(s.Commands["compile"].Args, []string{"mix", "compile"}) {
			t.Fatal(s)
		}
	}
	before := supportFileHashes(t, root)
	if e := Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("existing config accepted")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("manual config changed")
	}
}

func TestElixirDeclaredUmbrellaConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name, language, manager string
		manifest                bool
		ok                      bool
	}{
		{"declared", "elixir", "mix", true, true}, {"wrong-language", "javascript", "mix", true, false}, {"wrong-manager", "elixir", "npm", true, false}, {"missing-manifest", "elixir", "mix", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if tc.manifest {
				os.WriteFile(filepath.Join(root, "mix.exs"), []byte("not evaluated source"), 0600)
			}
			body := `{"version":1,"name":"declared","services":[{"id":"umbrella","dir":".","language":"` + tc.language + `","manager":"` + tc.manager + `","mix_umbrella":true,"commands":{"compile":{"args":["mix","compile"]}}}]}`
			os.WriteFile(filepath.Join(root, "mudarro.json"), []byte(body), 0600)
			before := supportFileHashes(t, root)
			c, e := Load(root)
			if (e == nil) != tc.ok {
				t.Fatalf("load err %v", e)
			}
			if tc.ok {
				if !reflect.DeepEqual(c.Services[0].Commands["compile"].Args, []string{"mix", "compile"}) {
					t.Fatal("explicit command changed")
				}
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("validation modified fixtures")
			}
		})
	}
}

func TestElixirBuildAndDependencySourcesIgnored(t *testing.T) {
	for _, dir := range []string{"_build", "deps"} {
		t.Run(dir, func(t *testing.T) {
			root := fixture(t, map[string]string{"mix.exs": "project evidence", dir + "/foreign/mix.exs": strings.Repeat("x", 4*1024*1024+1)})
			report, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			if len(report.Config.Services) != 1 || report.Config.Services[0].Dir != "." || len(report.Suggestions) != 2 {
				t.Fatalf("generated/dependency project leaked %+v", report)
			}
			if len(report.Warnings) != 0 {
				t.Fatalf("ignored oversized source read %v", report.Warnings)
			}
		})
	}
}
