//go:build linux || darwin

package terminal_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
	"io"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestPresentationInteractiveReadablePaletteAndLocale(t *testing.T) {
	for _, theme := range []string{"dark", "light"} {
		for _, noColor := range []string{"", "1"} {
			for _, locale := range []string{"en", "pt-BR"} {
				t.Run(fmt.Sprintf("%s/no-color=%s/%s", theme, noColor, locale), func(t *testing.T) {
					t.Setenv("TERM", "xterm-256color")
					t.Setenv("NO_COLOR", noColor)
					master, slave, e := openPTY()
					if e != nil {
						t.Fatal(e)
					}
					defer master.Close()
					defer slave.Close()
					type captured struct {
						data string
						err  error
					}
					capture := make(chan captured, 1)
					go func() { var b bytes.Buffer; _, e := io.Copy(&b, master); capture <- captured{b.String(), e} }()
					done := make(chan error, 1)
					view := NewMenu(slave, "Mudarro", "")
					view.Configure(Presentation{Locale: locale, Theme: theme, Mouse: "off"})
					title, label, hint := "Actions", "test action", "Down preview"
					if locale == "pt-BR" {
						title, label, hint = "Ações", "ação de teste", "↓ prévia"
					}
					go func() {
						n, supported, e := view.ChooseInteractive(slave, title, []string{label}, []string{"read-only plan"})
						if e == nil && (!supported || n != 0) {
							e = fmt.Errorf("unexpected selection %d", n)
						}
						done <- e
					}()
					time.Sleep(50 * time.Millisecond)
					if _, e = master.Write([]byte("q")); e != nil {
						t.Fatal(e)
					}
					select {
					case e = <-done:
						if e != nil {
							t.Fatal(e)
						}
					case <-time.After(time.Second):
						t.Fatal("interactive presentation timeout")
					}
					slave.Close()
					var output string
					select {
					case c := <-capture:
						if c.err != nil && !errors.Is(c.err, syscall.EIO) {
							t.Fatal(c.err)
						}
						output = c.data
					case <-time.After(time.Second):
						t.Fatal("capture timeout")
					}
					ansi := regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
					plain := ansi.ReplaceAllString(output, "")
					for _, want := range []string{"Mudarro / " + title, label, hint} {
						if !strings.Contains(plain, want) {
							t.Fatalf("readable %q missing: %q", want, plain)
						}
					}
					if strings.Contains(plain, "\x1b") {
						t.Fatalf("malformed terminal escape: %q", plain)
					}
					if noColor != "" && regexp.MustCompile(`\x1b\[[0-9;]*m`).MatchString(output) {
						t.Fatal("NO_COLOR emitted color SGR")
					}
					if strings.Contains(output, "\x1b[?1002h") || strings.Contains(output, "\x1b[?1006h") {
						t.Fatal("mouse off enabled reporting")
					}
				})
			}
		}
	}
}
func TestPresentationCompactTextBaseline(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	tty, finish := capturePTY(t, 80, 24)
	view := NewMenu(tty, "Mudarro", "")
	view.Configure(Presentation{Locale: "en", Theme: "dark", Density: "compact", Lettering: "text"})
	if n, e := view.Choose(bufio.NewReader(strings.NewReader("0\n")), "Services", []string{"app"}); e != nil || n != 0 {
		t.Fatalf("choose %d %v", n, e)
	}
	output := finish()
	plain := regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`).ReplaceAllString(output, "")
	if !strings.Contains(plain, "Mudarro") || !strings.Contains(plain, "Services") || !strings.Contains(plain, "app") {
		t.Fatalf("text baseline unreadable: %q", plain)
	}
	if strings.Contains(output, "█") || strings.Contains(output, "░") {
		t.Fatal("compact text emitted block lettering")
	}
}
