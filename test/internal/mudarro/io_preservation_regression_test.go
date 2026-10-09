package mudarro_test

import (
	"bytes"
	"context"
	"errors"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"github.com/viralabs-dev/mudarro/internal/mudarro/executor"
	"github.com/viralabs-dev/mudarro/internal/mudarro/projectfs"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type preservationFailWriter struct{ err error }

func (w preservationFailWriter) Write(p []byte) (int, error) { return 0, w.err }

type preservationHookWriter struct{ hook func(string) }

func (w preservationHookWriter) Write(p []byte) (int, error) { w.hook(string(p)); return len(p), nil }
func assertNoAtomicTemps(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if strings.HasPrefix(d.Name(), ".mudarro-tmp-") {
			t.Errorf("atomic temp leaked: %s", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestIOPreservationManifestReadFailurePreventsWrites(t *testing.T) {
	root := fixture(t, map[string]string{"keep.txt": "user data"})
	if e := os.MkdirAll(filepath.Join(root, ".mudarro/generated.json"), 0755); e != nil {
		t.Fatal(e)
	}
	before := supportFileHashes(t, root)
	if e := Generate(root, sample(), false, io.Discard); e == nil {
		t.Fatal("manifest read failure accepted")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("read failure modified user files")
	}
	if exists(filepath.Join(root, "menu.sh")) {
		t.Fatal("preflight wrote launcher")
	}
	assertNoAtomicTemps(t, root)
}
func TestIOPreservationDestinationCreationFailure(t *testing.T) {
	root := fixture(t, map[string]string{"keep.txt": "user data"})
	injected := false
	writer := preservationHookWriter{hook: func(line string) {
		if !injected && strings.HasPrefix(line, "gerar:") {
			injected = true
			if e := os.MkdirAll(filepath.Join(root, ".mudarro"), 0755); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(root, ".mudarro/scripts"), []byte("owned fixture obstacle"), 0644); e != nil {
				t.Fatal(e)
			}
		}
	}}
	if e := Generate(root, sample(), false, writer); e == nil {
		t.Fatal("destination creation failure accepted")
	}
	body, e := os.ReadFile(filepath.Join(root, ".mudarro/scripts"))
	if e != nil || string(body) != "owned fixture obstacle" {
		t.Fatalf("obstacle changed %q %v", body, e)
	}
	if exists(filepath.Join(root, "menu.sh")) || exists(filepath.Join(root, ".mudarro/generated.json")) {
		t.Fatal("failed first write produced launcher/manifest")
	}
	assertNoAtomicTemps(t, root)
}
func TestIOPreservationAtomicRenameFailureKeepsDestination(t *testing.T) {
	root := fixture(t, map[string]string{})
	injected := false
	writer := preservationHookWriter{hook: func(line string) {
		if !injected && strings.Contains(line, "gerar: menu.sh") {
			injected = true
			if e := os.Mkdir(filepath.Join(root, "menu.sh"), 0755); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(root, "menu.sh/user.txt"), []byte("keep"), 0644); e != nil {
				t.Fatal(e)
			}
		}
	}}
	if e := Generate(root, sample(), false, writer); e == nil {
		t.Fatal("rename into directory accepted")
	}
	if !injected {
		t.Fatal("rename fixture was not exercised")
	}
	body, e := os.ReadFile(filepath.Join(root, "menu.sh/user.txt"))
	if e != nil || string(body) != "keep" {
		t.Fatal("rename failure damaged destination")
	}
	if exists(filepath.Join(root, ".mudarro/generated.json")) {
		t.Fatal("failed rename committed ownership manifest")
	}
	assertNoAtomicTemps(t, root)
	// Generation is per-file atomic, not a multi-file transaction: earlier wrappers
	// can already exist. This test asserts destination/temp/manifest preservation.
}
func TestIOPreservationConfigReadAndInitFailure(t *testing.T) {
	root := fixture(t, map[string]string{"config.json": "{invalid", "mudarro.json": "{invalid"})
	before := supportFileHashes(t, root)
	if _, e := LoadFile(root, "config.json"); e == nil {
		t.Fatal("invalid config read accepted")
	}
	var out bytes.Buffer
	if e := Main([]string{"init", "--root", root, "--format", "json"}, "test", strings.NewReader(""), &out); e == nil {
		t.Fatal("existing invalid config overwritten")
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("failed read/init changed files")
	}
}
func TestIOPreservationExecutorOutputExitStatusAndMissingCommand(t *testing.T) {
	skipOnWindows(t, "fixture runs sh")
	e := executor.OS{}
	root := t.TempDir()
	var output bytes.Buffer
	err := e.Run(root, []string{"sh", "-c", "printf 'stdout-marker'; printf 'stderr-marker' >&2; exit 23"}, nil, &output)
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 23 {
		t.Fatalf("exit status lost: %v", err)
	}
	for _, marker := range []string{"stdout-marker", "stderr-marker"} {
		if !strings.Contains(output.String(), marker) {
			t.Fatalf("%s missing: %q", marker, output.String())
		}
	}
	body, err := e.Output(root, []string{"sh", "-c", "printf stdout-marker; printf stderr-marker >&2; exit 19"})
	if !errors.As(err, &exit) || exit.ExitCode() != 19 || string(body) != "stdout-marker" || !strings.Contains(string(exit.Stderr), "stderr-marker") {
		t.Fatalf("Output result/stderr/status %q %v", body, err)
	}
	err = e.Run(root, []string{filepath.Join(root, "missing-command")}, nil, io.Discard)
	if err == nil {
		t.Fatal("missing executable accepted")
	}
	if _, err = e.Output(filepath.Join(root, "missing-directory"), []string{"sh", "-c", "printf unreachable"}); err == nil {
		t.Fatal("missing working directory accepted")
	}
}
func TestIOPreservationGenerateWriterFailurePropagatesBeforeMutation(t *testing.T) {
	root := fixture(t, map[string]string{"keep.txt": "user data"})
	before := supportFileHashes(t, root)
	sentinel := errors.New("fixture output unavailable")
	err := Generate(root, sample(), false, preservationFailWriter{sentinel})
	if !errors.Is(err, sentinel) {
		t.Errorf("writer error not propagated: %v", err)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Error("output failure still mutated project")
	}
}
func TestIOPreservationManifestSpecialFileDoesNotBlock(t *testing.T) {
	if os.Getenv("MUDARRO_FIFO_READ_HELPER") == "1" {
		_, err := projectfs.ReadManifest(os.Getenv("MUDARRO_FIFO_ROOT"), "", "package.json")
		if err == nil {
			t.Fatal("special manifest accepted")
		}
		return
	}
	if os.PathSeparator != '/' {
		t.Skip("FIFO regression requires Unix")
	}
	root := t.TempDir()
	cmd := exec.Command("mkfifo", filepath.Join(root, "package.json"))
	if e := cmd.Run(); e != nil {
		t.Skipf("mkfifo unavailable: %v", e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, os.Args[0], "-test.run=^TestIOPreservationManifestSpecialFileDoesNotBlock$")
	cmd.Env = append(os.Environ(), "MUDARRO_FIFO_READ_HELPER=1", "MUDARRO_FIFO_ROOT="+root)
	out, e := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("manifest FIFO read blocked; helper killed and joined: %v", ctx.Err())
	}
	if e != nil {
		t.Fatalf("helper read error handling: %v %s", e, out)
	}
}
func TestIOPreservationPermissionFailures(t *testing.T) {
	skipOnWindows(t, "Unix permission bits (chmod 000) do not deny reads on Windows")
	if os.Getuid() == 0 {
		t.Skip("permission failure cannot be demonstrated as root")
	}
	root := fixture(t, map[string]string{"config.json": `{"version":1,"name":"test","services":[]}`})
	path := filepath.Join(root, "config.json")
	if e := os.Chmod(path, 0); e != nil {
		t.Fatal(e)
	}
	defer os.Chmod(path, 0644)
	if _, e := LoadFile(root, "config.json"); e == nil {
		t.Fatal("unreadable configuration accepted")
	}
	// Restrict only our fixture directory and always restore it for cleanup.
	if e := os.Chmod(root, 0500); e != nil {
		t.Fatal(e)
	}
	defer os.Chmod(root, 0700)
	if e := Main([]string{"init", "--root", root}, "test", strings.NewReader(""), io.Discard); e == nil {
		t.Fatal("init creation failure accepted")
	}
	assertNoAtomicTemps(t, root)
}

type preservationNthFailWriter struct {
	calls, failAt int
	err           error
}

func (w *preservationNthFailWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls == w.failAt {
		return 0, w.err
	}
	return len(p), nil
}
func TestIOPreservationSecondAnnouncementFailurePreventsAllWrites(t *testing.T) {
	for _, dry := range []bool{false, true} {
		t.Run(map[bool]string{false: "write", true: "dry-run"}[dry], func(t *testing.T) {
			root := fixture(t, map[string]string{"keep.txt": "user data"})
			before := supportFileHashes(t, root)
			sentinel := errors.New("fixture second announcement unavailable")
			writer := &preservationNthFailWriter{failAt: 2, err: sentinel}
			if err := Generate(root, sample(), dry, writer); !errors.Is(err, sentinel) {
				t.Fatalf("second announcement failure: %v", err)
			}
			if writer.calls != 2 {
				t.Fatalf("failed point not exercised: %d", writer.calls)
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("later output failure mutated project")
			}
			assertNoAtomicTemps(t, root)
		})
	}
}
func TestIOPreservationPreservedMessageFailurePreventsWrites(t *testing.T) {
	root := fixture(t, map[string]string{"menu.sh": "#!/bin/sh\nprintf user-owned\n", "keep.txt": "user data"})
	before := supportFileHashes(t, root)
	sentinel := errors.New("fixture preservation report unavailable")
	if err := Generate(root, sample(), false, preservationFailWriter{sentinel}); !errors.Is(err, sentinel) {
		t.Fatalf("preservation message failure: %v", err)
	}
	if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
		t.Fatal("preservation message failure changed files")
	}
	assertNoAtomicTemps(t, root)
}
func TestIOPreservationReadManifestBoundedRegularFiles(t *testing.T) {
	for _, size := range []int{4 << 20, (4 << 20) + 1} {
		t.Run(map[bool]string{false: "limit", true: "over-limit"}[size > 4<<20], func(t *testing.T) {
			root := fixture(t, map[string]string{"package.json": strings.Repeat("x", size)})
			body, err := projectfs.ReadManifest(root, ".", "package.json")
			if size > 4<<20 {
				if err == nil {
					t.Fatal("oversized manifest accepted")
				}
				if len(body) != 0 {
					t.Fatal("oversized manifest returned data")
				}
			} else if err != nil || len(body) != size {
				t.Fatalf("bounded regular file len=%d err=%v", len(body), err)
			}
		})
	}
	root := t.TempDir()
	if _, err := projectfs.ReadManifest(root, ".", "missing.json"); !os.IsNotExist(err) {
		t.Fatalf("missing file error lost: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "directory.json"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := projectfs.ReadManifest(root, ".", "directory.json"); err == nil {
		t.Fatal("directory manifest accepted")
	}
	outside := fixture(t, map[string]string{"outside.json": "personal-free fixture"})
	if err := os.Symlink(filepath.Join(outside, "outside.json"), filepath.Join(root, "linked.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := projectfs.ReadManifest(root, ".", "linked.json"); err == nil {
		t.Fatal("symlink manifest accepted")
	}
}

func TestIOPreservationOwnershipManifestFIFORejectsWithoutBlocking(t *testing.T) {
	if os.Getenv("MUDARRO_OWNERSHIP_FIFO_HELPER") == "1" {
		if err := Generate(os.Getenv("MUDARRO_FIFO_ROOT"), sample(), false, io.Discard); err == nil {
			t.Fatal("ownership FIFO accepted")
		}
		return
	}
	if os.PathSeparator != '/' {
		t.Skip("FIFO regression requires Unix")
	}
	root := fixture(t, map[string]string{"keep.txt": "owned fixture user data"})
	if err := os.Mkdir(filepath.Join(root, ".mudarro"), 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("mkfifo", filepath.Join(root, ".mudarro/generated.json"))
	if err := cmd.Run(); err != nil {
		t.Skipf("mkfifo unavailable: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, os.Args[0], "-test.run=^TestIOPreservationOwnershipManifestFIFORejectsWithoutBlocking$")
	cmd.Env = append(os.Environ(), "MUDARRO_OWNERSHIP_FIFO_HELPER=1", "MUDARRO_FIFO_ROOT="+root)
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("ownership FIFO blocked; helper killed and joined: %v", ctx.Err())
	}
	if err != nil {
		t.Fatalf("ownership FIFO helper: %v %s", err, out)
	}
	body, err := os.ReadFile(filepath.Join(root, "keep.txt"))
	if err != nil || string(body) != "owned fixture user data" {
		t.Fatal("ownership read error damaged user file")
	}
	if exists(filepath.Join(root, "menu.sh")) || exists(filepath.Join(root, ".mudarro/scripts")) {
		t.Fatal("ownership read error generated files")
	}
	assertNoAtomicTemps(t, root)
}
func TestIOPreservationMalformedAndOversizedOwnershipManifest(t *testing.T) {
	for _, tc := range []struct{ name, body string }{{"malformed", "{invalid ownership manifest"}, {"oversized", `{"version":1,"files":{}}` + strings.Repeat(" ", 4<<20)}} {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t, map[string]string{".mudarro/generated.json": tc.body, "keep.txt": "user data"})
			before := supportFileHashes(t, root)
			if err := Generate(root, sample(), false, io.Discard); err == nil {
				t.Fatal("bad ownership manifest accepted")
			}
			if !reflect.DeepEqual(before, supportFileHashes(t, root)) {
				t.Fatal("ownership manifest failure changed project")
			}
			if exists(filepath.Join(root, "menu.sh")) {
				t.Fatal("ownership manifest failure wrote launcher")
			}
			assertNoAtomicTemps(t, root)
		})
	}
}
