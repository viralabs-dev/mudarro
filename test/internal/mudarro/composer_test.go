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

func TestComposerScriptsOptInAndExactArgv(t *testing.T) {
	root := fixture(t, map[string]string{"composer.json": `{"name":"sample/php","autoload":{"psr-4":{"App\\":"src/"}},"scripts":{"test:unit":"php tests.php","start":"EXECUTE_ONLY_IF_SELECTED","post-install-cmd":"NEVER_SCAN"}}`})
	before := supportFileHashes(t, root)
	s, su, e := (adapters.Composer{}).Detect(root, ".", map[string]bool{"composer.json": true}, nil)
	if e != nil {
		t.Fatal(e)
	}
	if s.Manager != "composer" || s.Language != "php" || s.Framework != "" || len(s.Commands) != 0 || len(su) != 3 {
		t.Fatalf("%+v %+v", s, su)
	}
	for _, x := range su {
		want := strings.ReplaceAll(x.Name, "test-unit", "test:unit")
		if !reflect.DeepEqual(x.Command.Args, []string{"composer", "run-script", "--", want}) {
			t.Fatal(x)
		}
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("modified source")
	}
}
func TestComposerMalformedAndAmbiguousScripts(t *testing.T) {
	for _, body := range []string{`[]`, `null`, `{"scripts":[]}`, `{"scripts":null}`, `{"scripts":"wrong"}`, `{`} {
		t.Run(body, func(t *testing.T) {
			root := fixture(t, map[string]string{"composer.json": body})
			if _, _, e := (adapters.Composer{}).Detect(root, ".", map[string]bool{"composer.json": true}, nil); e == nil {
				t.Fatal("invalid accepted")
			}
		})
	}
	root := fixture(t, map[string]string{"composer.json": `{"scripts":{"good":"php check.php","array":["php one.php","php two.php"],"number":7,"empty":"","a:b":"one","a-b":"two","--bad":"bad"}}`})
	s, su, e := (adapters.Composer{}).Detect(root, ".", map[string]bool{"composer.json": true}, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(su) != 1 || su[0].Name != "good" || len(s.Pending) < 6 {
		t.Fatalf("%+v %+v", s, su)
	}
}
func TestComposerEvidenceBoundaries(t *testing.T) {
	for _, kind := range []string{"oversized", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			p := filepath.Join(root, "composer.json")
			if kind == "oversized" {
				os.WriteFile(p, bytes.Repeat([]byte("x"), 4*1024*1024+1), 0600)
			} else {
				other := filepath.Join(t.TempDir(), "json")
				os.WriteFile(other, []byte(`{}`), 0600)
				if e := os.Symlink(other, p); e != nil {
					t.Skip(e)
				}
			}
			if _, _, e := (adapters.Composer{}).Detect(root, ".", map[string]bool{"composer.json": true}, nil); e == nil {
				t.Fatal("invalid evidence accepted")
			}
		})
	}
	s, _, e := (adapters.Composer{}).Detect(t.TempDir(), ".", nil, nil)
	if s != nil || e != nil {
		t.Fatal("missing evidence detected")
	}
}
func TestComposerScanSelectionPreservesConfiguration(t *testing.T) {
	root := fixture(t, map[string]string{"composer.json": `{"scripts":{"test:unit":"php tests.php"}}`, "excluded/composer.json": strings.Repeat("x", 4*1024*1024+1)})
	trap := t.TempDir()
	marker := filepath.Join(root, "NO_EXEC")
	os.WriteFile(filepath.Join(trap, "composer"), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0755)
	t.Setenv("PATH", trap)
	r, e := Scan(root, []string{"excluded"})
	if e != nil || len(r.Suggestions) != 1 {
		t.Fatalf("%+v %v", r, e)
	}
	if exists(marker) {
		t.Fatal("executed Composer")
	}
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--exclude", "excluded", "--select", "app-php:test-unit"}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, root)
	if e := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	h := supportFileHashes(t, root)
	if e := Main([]string{"generate", "--root", root}, "test", strings.NewReader(""), &out); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(h, supportFileHashes(t, root)) {
		t.Fatal("generate changed files")
	}
	if e := Main([]string{"init", "--root", root}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("existing config accepted")
	}
	now := supportFileHashes(t, root)
	if before["mudarro.yaml"] != now["mudarro.yaml"] {
		t.Fatal("manual config changed")
	}
}

func TestComposerDuplicateAndInvalidEncodingRejected(t *testing.T) {
	for _, body := range []string{`{"scripts":{"test":"one","test":"two"}}`, `{"scripts":{"test":"one"},"scripts":{"test":"two"}}`, string([]byte{'{', '"', 'x', '"', ':', '"', 0xff, '"', '}'})} {
		root := fixture(t, map[string]string{"composer.json": body})
		if _, _, e := (adapters.Composer{}).Detect(root, ".", map[string]bool{"composer.json": true}, nil); e == nil {
			t.Fatal("ambiguous/invalid JSON accepted")
		}
	}
}
