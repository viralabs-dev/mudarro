package terminal_test

import (
	"bufio"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRichMenuLayoutAndQuit(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	t.Setenv("COLUMNS", "100")
	tty, finish := capturePTY(t, 100, 32)
	v := NewMenu(tty, "Mudarro", "1 serviço / offline")
	v.Configure(Presentation{Locale: "pt-BR"})
	n, err := v.Choose(bufio.NewReader(strings.NewReader("q\n")), "Serviços", []string{"app"})
	if err != nil || n != 0 {
		t.Fatalf("quit %d %v", n, err)
	}
	output := finish()
	for _, want := range []string{"\x1b[2J", "PROJETO / Mudarro", "█", "░", "Serviços", "0 / q", "Enter"} {
		if !strings.Contains(output, want) {
			t.Fatal("missing", want)
		}
	}
}
func TestRichMenuFitsNarrowWindow(t *testing.T) {
	t.Setenv("TERM", "xterm")
	t.Setenv("NO_COLOR", "")
	t.Setenv("COLUMNS", "24")
	tty, finish := capturePTY(t, 24, 24)
	v := NewMenu(tty, "Mudarro with a very long name", "many services and a long context")
	_, err := v.Choose(bufio.NewReader(strings.NewReader("invalid\n0\n")), "Very long title needing compact layout", []string{"a very long path with spaces and a long command"})
	if err != nil {
		t.Fatal(err)
	}
	output := strings.ReplaceAll(finish(), "\r\n", "\n")
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)
	for _, line := range strings.Split(ansi.ReplaceAllString(output, ""), "\n") {
		// Input readers do not echo Enter as a real terminal does.
		line = strings.TrimPrefix(line, "  > ")
		if utf8.RuneCountInString(line) > 24 {
			t.Fatalf("width overflow: %q", line)
		}
	}
}
