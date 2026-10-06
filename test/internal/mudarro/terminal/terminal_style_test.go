package terminal_test

import (
	"bufio"
	"bytes"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"os"
	"strings"
	"testing"
)

func TestTerminalStyleResetsAndPlainFallback(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	tty, _ := capturePTY(t, 80, 24)
	styled := NewMenu(tty, "Mudarro", "").Style().Paint("1;36", "Menu")
	if styled != "\x1b[1;36mMenu\x1b[0m" {
		t.Fatal(styled)
	}
	if (Style{}).Paint("1;36", "Menu") != "Menu" {
		t.Fatal("plain fallback")
	}
	t.Setenv("NO_COLOR", "1")
	if NewMenu(tty, "Mudarro", "").Style().HasColor() {
		t.Fatal("NO_COLOR ignored")
	}
}
func TestBannerNarrowWidthAndControlSanitization(t *testing.T) {
	t.Setenv("COLUMNS", "24")
	var out bytes.Buffer
	_, err := NewMenu(&out, "long sample name\x1b\n", "").Choose(bufio.NewReader(strings.NewReader("0\n")), "Area", nil)
	if err != nil {
		t.Fatal(err)
	}
	// The public plain menu starts with the complete ASCII banner.
	b, _, ok := strings.Cut(out.String(), "\n\n  ")
	if !ok {
		t.Fatalf("banner delimiter missing: %q", out.String())
	}
	b += "\n"
	lines := strings.Split(b, "\n")
	for _, line := range lines[:len(lines)-2] {
		if len(line) > 24 {
			t.Fatalf("banner overflow: %q", line)
		}
	}
	if strings.ContainsAny(b, "\x1b") {
		t.Fatal("control characters in banner")
	}
}

func TestColorPolicy(t *testing.T) {
	for _, tc := range []struct {
		name          string
		tty           bool
		term, noColor string
		want          bool
	}{
		{"tty", true, "xterm-256color", "", true},
		{"pipe", false, "xterm-256color", "", false},
		{"dumb", true, "dumb", "", false},
		{"unset-term", true, "", "", false},
		{"no-color", true, "xterm", "1", false},
		{"zero-no-color", true, "xterm", "0", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TERM", tc.term)
			t.Setenv("NO_COLOR", tc.noColor)
			var got bool
			if tc.tty {
				tty, _ := capturePTY(t, 80, 24)
				got = NewMenu(tty, "Mudarro", "").Style().HasColor()
			} else {
				got = NewMenu(&bytes.Buffer{}, "Mudarro", "").Style().HasColor()
			}
			if got != tc.want {
				t.Fatalf("color=%v want %v", got, tc.want)
			}
		})
	}
}
func TestDevNullIsNotTerminal(t *testing.T) {
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	if NewMenu(f, "Mudarro", "").Style().HasColor() {
		t.Fatal("character device treated as terminal")
	}
}
