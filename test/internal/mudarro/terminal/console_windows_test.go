//go:build windows

package terminal_test

import (
	"os"
	"testing"

	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
)

// Redirected handles are not consoles: the Windows backend must fall back to
// the plain, colourless, non-interactive path without touching console modes.
func TestWindowsRedirectedHandlesAreNotInteractive(t *testing.T) {
	t.Setenv("TERM", "")
	t.Setenv("NO_COLOR", "")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	v := NewMenu(w, "project", "context")
	if v.Style().HasColor() {
		t.Fatal("pipe output reported colour")
	}
	if n, supported, err := v.ChooseInteractive(r, "title", []string{"a"}, []string{"preview"}); supported || err != nil || n != 0 {
		t.Fatalf("pipe input entered interactive mode: %d %v %v", n, supported, err)
	}
	if active, err := v.BeginShell(r); active || err != nil {
		t.Fatalf("pipe input opened the shell frame: %v %v", active, err)
	}
}
