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

func TestMavenStaticContract(t *testing.T) {
	root := fixture(t, map[string]string{"pom.xml": `<project xmlns="http://maven.apache.org/POM/4.0.0"><modelVersion>4.0.0</modelVersion><artifactId>sample</artifactId><build><plugins><plugin><artifactId>exec-maven-plugin</artifactId></plugin></plugins></build></project>`})
	trap := t.TempDir()
	marker := filepath.Join(root, "EXECUTED")
	os.WriteFile(filepath.Join(trap, "mvn"), []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0755)
	t.Setenv("PATH", trap)
	before := supportFileHashes(t, root)
	s, su, e := (adapters.Java{}).Detect(root, ".", map[string]bool{"pom.xml": true}, nil)
	if e != nil || s == nil || s.Manager != "maven" || len(s.Commands) != 0 || len(su) != 2 {
		t.Fatalf("%+v %+v %v", s, su, e)
	}
	if !reflect.DeepEqual(su[0].Command.Args, []string{"mvn", "-f", "pom.xml", "compile"}) || su[1].Name != "test" {
		t.Fatal(su)
	}
	if exists(marker) || !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("scan executed or mutated source")
	}
}
func TestMavenInvalidEvidence(t *testing.T) {
	for _, src := range []string{"", `<project>`, `<other/>`, `<project/><project/>`, `<!DOCTYPE project [<!ENTITY remote SYSTEM "file:///etc/passwd">]><project>&remote;</project>`, `<project>&unknown;</project>`, strings.Repeat("x", 4*1024*1024+1)} {
		root := fixture(t, map[string]string{"pom.xml": src})
		if _, _, e := (adapters.Java{}).Detect(root, ".", map[string]bool{"pom.xml": true}, nil); e == nil {
			t.Fatal("invalid XML accepted")
		}
	}
	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "pom.xml")
	os.WriteFile(target, []byte("<project/>"), 0600)
	os.Symlink(target, filepath.Join(root, "pom.xml"))
	if _, _, e := (adapters.Java{}).Detect(root, ".", map[string]bool{"pom.xml": true}, nil); e == nil {
		t.Fatal("symlink accepted")
	}
	s, _, e := (adapters.Java{}).Detect(root, ".", map[string]bool{}, nil)
	if e != nil || s != nil {
		t.Fatal("absent detected")
	}
}
func TestMavenModulesAndExclusions(t *testing.T) {
	root := fixture(t, map[string]string{"pom.xml": `<project><modules><module>child</module><module>ignored</module></modules></project>`, "child/pom.xml": "<project/>", "ignored/pom.xml": strings.Repeat("x", 4*1024*1024+1)})
	s, su, e := (adapters.Java{}).Detect(root, ".", map[string]bool{"pom.xml": true}, []string{"ignored"})
	if e != nil || len(su) != 0 || len(s.Pending) < 2 {
		t.Fatalf("%+v %v", s, e)
	}
	for _, member := range []string{"../escape", "/absolute", "${dynamic}", "missing"} {
		os.WriteFile(filepath.Join(root, "pom.xml"), []byte("<project><modules><module>"+member+"</module></modules></project>"), 0600)
		if _, _, e := (adapters.Java{}).Detect(root, ".", map[string]bool{"pom.xml": true}, nil); e == nil {
			t.Fatalf("accepted %q", member)
		}
	}
}
func TestGradleUnevaluatedEvidence(t *testing.T) {
	for _, name := range []string{"build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"} {
		root := fixture(t, map[string]string{name: "throw new RuntimeException('must never execute')"})
		s, su, e := (adapters.Java{}).Detect(root, ".", map[string]bool{name: true}, nil)
		if e != nil || s.Manager != "gradle" || len(su) != 0 || len(s.Commands) != 0 || !bytes.Contains([]byte(s.Pending[0]), []byte("not evaluated")) {
			t.Fatalf("%+v %+v %v", s, su, e)
		}
	}
}

func TestMavenManifestExcludedBeforeRead(t *testing.T) {
	root := fixture(t, map[string]string{"pom.xml": "<project><modules><module>child</module></modules></project>", "child/pom.xml": strings.Repeat("x", 4*1024*1024+1)})
	s, su, e := (adapters.Java{}).Detect(root, ".", map[string]bool{"pom.xml": true}, []string{"child/pom.xml"})
	if e != nil || s == nil || len(su) != 0 {
		t.Fatal(e)
	}
}

func TestMavenScanSelectionGeneratePreservesManual(t *testing.T) {
	root := fixture(t, map[string]string{"pom.xml": `<project/>`, "excluded/pom.xml": strings.Repeat("x", 4*1024*1024+1)})
	report, e := Scan(root, []string{"excluded"})
	if e != nil || len(report.Config.Services) != 1 || len(report.Suggestions) != 2 || len(report.Warnings) != 0 {
		t.Fatalf("%+v %v", report, e)
	}
	var out bytes.Buffer
	if e = Main([]string{"init", "--root", root, "--exclude", "excluded", "--select", "app-java:compile,app-java:test"}, "test", strings.NewReader(""), &out); e != nil {
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
