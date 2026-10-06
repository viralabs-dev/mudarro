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

func TestRubyExecutableEvidenceNeverEvaluated(t *testing.T) {
	root := fixture(t, map[string]string{"Gemfile": "File.write('EXECUTED', 'bad')\nINVALID Ruby source", "sample.gemspec": "raise 'never evaluate'", "Gemfile.lock": "invalid lock evidence"})
	before := supportFileHashes(t, root)
	s, su, e := (adapters.Ruby{}).Detect(root, ".", map[string]bool{"Gemfile": true, "sample.gemspec": true, "Gemfile.lock": true}, nil)
	if e != nil {
		t.Fatal(e)
	}
	if s.Language != "ruby" || s.Manager != "bundler" || len(s.Commands) != 0 || len(su) != 0 || len(s.Pending) < 4 {
		t.Fatalf("%+v %+v", s, su)
	}
	if exists(filepath.Join(root, "EXECUTED")) || !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("Ruby code executed or modified")
	}
}
func TestRubyEvidenceBoundaries(t *testing.T) {
	for _, name := range []string{"Gemfile", "sample.gemspec", "Gemfile.lock"} {
		t.Run(name, func(t *testing.T) {
			root := fixture(t, map[string]string{"Gemfile": "evidence"})
			os.WriteFile(filepath.Join(root, name), bytes.Repeat([]byte("x"), 4*1024*1024+1), 0600)
			if _, _, e := (adapters.Ruby{}).Detect(root, ".", map[string]bool{"Gemfile": true, name: true}, nil); e == nil {
				t.Fatal("oversized accepted")
			}
		})
	}
	root := t.TempDir()
	other := filepath.Join(t.TempDir(), "Gemfile")
	os.WriteFile(other, []byte("evidence"), 0600)
	if e := os.Symlink(other, filepath.Join(root, "Gemfile")); e != nil {
		t.Skip(e)
	}
	if _, _, e := (adapters.Ruby{}).Detect(root, ".", map[string]bool{"Gemfile": true}, nil); e == nil {
		t.Fatal("symlink accepted")
	}
	s, _, e := (adapters.Ruby{}).Detect(t.TempDir(), ".", map[string]bool{"Gemfile.lock": true}, nil)
	if e != nil || s != nil {
		t.Fatal("lock alone invented project")
	}
}
func TestRubyScanNestedExcludedNoCommands(t *testing.T) {
	root := fixture(t, map[string]string{"Gemfile": "evidence", "member/example.gemspec": "raise 'never evaluate'", "excluded/Gemfile": strings.Repeat("x", 4*1024*1024+1)})
	trap := t.TempDir()
	marker := filepath.Join(root, "EXECUTED")
	for _, tool := range []string{"ruby", "bundle"} {
		os.WriteFile(filepath.Join(trap, tool), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0755)
	}
	t.Setenv("PATH", trap)
	r, e := Scan(root, []string{"excluded"})
	if e != nil || len(r.Config.Services) != 2 {
		t.Fatalf("%+v %v", r, e)
	}
	for _, s := range r.Config.Services {
		if len(s.Commands) != 0 || s.Framework != "" {
			t.Fatal("invented Ruby command/framework")
		}
	}
	if len(r.Suggestions) != 0 || exists(marker) {
		t.Fatal("unexpected suggestion or execution")
	}
}
