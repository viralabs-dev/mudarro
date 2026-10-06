//go:build linux || darwin

package mudarro_test

import (
	"context"
	"encoding/json"
	. "github.com/viralabs-dev/mudarro/internal/mudarro"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestGoSourceIOBoundedRegularSources(t *testing.T) {
	for _, tc := range []struct {
		name      string
		size      int
		wantStart bool
	}{
		{"ordinary", 64, true}, {"exact-limit", 4 << 20, true}, {"over-limit", (4 << 20) + 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prefix, suffix := "package main\n/*", "*/\nfunc main(){}\n"
			source := prefix + strings.Repeat("x", tc.size-len(prefix)-len(suffix)) + suffix
			root := fixture(t, map[string]string{"go.mod": "module example.test/source\ngo 1.23\n", "main.go": source})
			r, e := Scan(root, nil)
			if e != nil {
				t.Fatal(e)
			}
			hasStart, unresolved := false, len(r.Warnings) > 0
			for _, s := range r.Config.Services {
				hasStart = hasStart || len(s.Commands["start"].Args) > 0
				unresolved = unresolved || len(s.Pending) > 0
			}
			if hasStart != tc.wantStart {
				t.Errorf("start=%v want=%v; report=%+v", hasStart, tc.wantStart, r)
			}
			if !tc.wantStart && !unresolved {
				t.Error("rejected source must remain visibly unresolved")
			}
			after, e := os.ReadFile(filepath.Join(root, "main.go"))
			if e != nil || string(after) != source {
				t.Fatal("source mutated", e)
			}
		})
	}
}

// Isolate a potentially blocking scan in this owned test process. CommandContext
// kills and Wait joins it on timeout; no terminal or unrelated process is used.
func TestGoSourceIOFIFOIsNotOpened(t *testing.T) {
	root := fixture(t, map[string]string{"go.mod": "module example.test/fifo\ngo 1.23\n"})
	if e := syscall.Mkfifo(filepath.Join(root, "main.go"), 0600); e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	child := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGoSourceIOIsolatedProbe$")
	child.Env = append(os.Environ(), "MUDARRO_GO_SOURCE_PROBE="+root)
	output, e := child.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("scan opened FIFO and blocked; owned helper killed/joined: %v", ctx.Err())
	}
	if e != nil {
		t.Fatalf("helper: %v %s", e, output)
	}
	var r Report
	if e = json.Unmarshal(output, &r); e != nil {
		t.Fatalf("helper report: %v %s", e, output)
	}
	for _, s := range r.Config.Services {
		if len(s.Commands["start"].Args) > 0 {
			t.Fatalf("FIFO inferred entry: %+v", s)
		}
	}
	st, e := os.Lstat(filepath.Join(root, "main.go"))
	if e != nil || st.Mode()&os.ModeNamedPipe == 0 {
		t.Fatal("FIFO fixture changed", e)
	}
}

func TestGoSourceIOIsolatedProbe(t *testing.T) {
	root := os.Getenv("MUDARRO_GO_SOURCE_PROBE")
	if root == "" {
		return
	}
	r, e := Scan(root, nil)
	if e != nil {
		os.Exit(2)
	}
	if e = json.NewEncoder(os.Stdout).Encode(r); e != nil {
		os.Exit(3)
	}
	os.Exit(0)
}
