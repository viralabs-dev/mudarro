package terminal_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
)

func TestConsumerProjectNamesAndWidths(t *testing.T) {
	for _, tc := range []struct{ name, visible string }{
		{"Atlas API", "Atlas API"}, {"Aurora", "Aurora"}, {"mudarro", "mudarro"},
		{"Café Equipe", "Café Equipe"}, {"漢字😀", "漢字😀"},
		{"Equipe com um nome de projeto muito longo", "Equipe com um nom"},
		{"Atlas\x1b]52;secret\a\n", "Atlas]52;secret"},
	} {
		for _, mode := range []string{"rich", "narrow", "plain", "no-color"} {
			t.Run(fmt.Sprintf("%s/%s", tc.name, mode), func(t *testing.T) {
				width := 100
				if mode == "narrow" {
					width = 24
				}
				t.Setenv("TERM", "xterm-256color")
				t.Setenv("COLUMNS", fmt.Sprint(width))
				t.Setenv("NO_COLOR", "")
				var out bytes.Buffer
				var writer io.Writer = &out
				finish := func() string { return out.String() }
				if mode != "plain" {
					writer, finish = capturePTY(t, uint16(width), 32)
				}
				if mode == "no-color" {
					t.Setenv("NO_COLOR", "1")
				}
				_, err := NewMenu(writer, tc.name, "offline").Choose(bufio.NewReader(strings.NewReader("0\n")), "Área", nil)
				if err != nil {
					t.Fatal(err)
				}
				output := strings.ReplaceAll(finish(), "\r\n", "\n")
				if (mode == "plain" || mode == "no-color") && strings.Contains(output, "\x1b") {
					t.Fatal("ANSI in plain/NO_COLOR")
				}
				ansi := regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)
				text := ansi.ReplaceAllString(output, "")
				if strings.ContainsAny(text, "\x1b\a") {
					t.Fatal("project name injected control sequence")
				}
				// Narrow headings reserve prefix space, so assert the identifiable project prefix.
				want := tc.visible
				if width == 24 && len(want) > 16 {
					want = want[:16]
				}
				if !strings.Contains(text, want) && !strings.Contains(text, strings.ToUpper(want)) {
					t.Fatalf("name not readable: %q in %q", want, text)
				}
				if !strings.EqualFold(tc.name, "mudarro") && strings.Contains(text, "MUDARRO") {
					t.Fatal("consumer menu replaced project with tool branding")
				}
				for _, line := range strings.Split(text, "\n") {
					line = strings.TrimPrefix(line, "  > ")
					// These fixtures contain three known double-cell characters; all others use one.
					cells := utf8.RuneCountInString(line) + strings.Count(line, "漢") + strings.Count(line, "字") + strings.Count(line, "😀")
					if cells > width {
						t.Fatalf("width %d>%d: %q", cells, width, line)
					}
				}
			})
		}
	}
}
