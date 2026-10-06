//go:build linux || darwin

package terminal_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	. "github.com/viralabs-dev/mudarro/internal/mudarro/terminal"
)

// A deliberately small oracle for the test alphabet. Joined emoji/flags are bounded
// by the sum of their components, not a promise about any emulator's grapheme shaping.
func unicodeFixtureCells(text string) int {
	cells := 0
	for _, r := range text {
		switch {
		case unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r):
		case r == '界' || r == '語' || r == '日' || r == 'Ａ' || r == 'Ｂ' || r == '\u3000' || r == '🙂' || r == '👨' || r == '👩' || r == '👧' || r == '👦' || r == '🇧' || r == '🇷':
			cells += 2
		default:
			cells++
		}
	}
	return cells
}

var unicodeFixtureANSI = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

func TestUnicodeResizeKnownCellsStayInsideNarrowMenu(t *testing.T) {
	for _, tc := range []struct{ name, text string }{
		{"cjk", strings.Repeat("界語日", 12)}, {"fullwidth-letters", strings.Repeat("ＡＢ", 18)},
		{"combining", strings.Repeat("e\u0301", 30)}, {"emoji", strings.Repeat("🙂", 20)},
		{"joined-family-conservative", strings.Repeat("👨\u200d👩\u200d👧\u200d👦", 8)}, {"regional-pair-conservative", strings.Repeat("🇧🇷", 12)},
		{"ideographic-space", strings.Repeat("\u3000", 30)},
	} {
		for _, noColor := range []string{"", "1"} {
			t.Run(tc.name+"/NO_COLOR="+noColor, func(t *testing.T) {
				const width = 24
				t.Setenv("TERM", "xterm-256color")
				t.Setenv("NO_COLOR", noColor)
				t.Setenv("COLUMNS", "24")
				tty, finish := capturePTY(t, width, 24)
				view := NewMenu(tty, "Project "+tc.text, "context "+tc.text)
				view.Configure(Presentation{Locale: "en", Density: "compact", Lettering: "text"})
				n, err := view.Choose(bufio.NewReader(strings.NewReader("0\n")), "area "+tc.text, []string{"service " + tc.text})
				if err != nil || n != 0 {
					t.Fatal(n, err)
				}
				output := strings.ReplaceAll(finish(), "\r\n", "\n")
				if !utf8.ValidString(output) {
					t.Fatal("invalid UTF-8 emitted")
				}
				plain := unicodeFixtureANSI.ReplaceAllString(output, "")
				for _, line := range strings.Split(plain, "\n") {
					if cells := unicodeFixtureCells(line); cells > width {
						t.Errorf("%d cells exceed %d: %q", cells, width, line)
					}
					for _, r := range line {
						if unicode.Is(unicode.Cf, r) {
							t.Errorf("format control emitted: U+%04X", r)
						}
					}
				}
				if noColor != "" && regexp.MustCompile(`\x1b\[[0-9;]*m`).MatchString(strings.ReplaceAll(output, "\x1b[0m", "")) {
					t.Fatal("NO_COLOR emitted SGR")
				}
			})
		}
	}
}
func TestUnicodeResizePlainLetteringLongNameIsBounded(t *testing.T) {
	t.Setenv("TERM", "xterm")
	t.Setenv("NO_COLOR", "1")
	t.Setenv("COLUMNS", "24")
	tty, finish := capturePTY(t, 24, 24)
	view := NewMenu(tty, strings.Repeat("界", 40), "")
	if _, err := view.Choose(bufio.NewReader(strings.NewReader("0\n")), "Services", nil); err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.ReplaceAll(finish(), "\r\n", "\n"), "\n") {
		if unicodeFixtureCells(line) > 24 {
			t.Fatalf("plain name overflow: %q", line)
		}
	}
}
func TestUnicodeResizeExecutionWriterStripsFormatAndPreservesSplitCombining(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	tty, finish := capturePTY(t, 28, 12)
	view := NewMenu(tty, "Unicode", "")
	view.Configure(Presentation{Density: "compact"})
	active, err := view.BeginShell(tty)
	if err != nil || !active {
		t.Fatal(active, err)
	}
	writer := view.ShellOutput()
	for _, part := range []string{"e", "\xcc", "\x81", " + ", "界", "\u202e", "CONTROL_END", "\u2066", "\u200b", "\n"} {
		if n, err := writer.Write([]byte(part)); err != nil || n != len(part) {
			t.Fatal(n, err)
		}
	}
	view.CloseShell()
	output := finish()
	if !utf8.ValidString(output) || !strings.Contains(output, "e\u0301") {
		t.Fatal("split combining UTF-8 lost")
	}
	for _, r := range output {
		if unicode.Is(unicode.Cf, r) {
			t.Errorf("format/bidi control escaped writer: U+%04X", r)
		}
	}
	if !strings.Contains(output, "CONTROL_END") {
		t.Fatal("sanitization discarded legitimate text")
	}
}
func TestUnicodeResizeLivePreviewRestoresTermiosAfterShrinkAndGrow(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "1")
	master, slave, err := openPTY()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	defer slave.Close()
	resize := func(width, rows uint16) {
		t.Helper()
		s := struct{ rows, cols, x, y uint16 }{rows: rows, cols: width}
		if err := ioctl(slave, uintptr(syscall.TIOCSWINSZ), unsafe.Pointer(&s)); err != nil {
			t.Fatal(err)
		}
	}
	resize(80, 24)
	before, err := interactiveTermios(slave)
	if err != nil {
		t.Fatal(err)
	}
	type result struct {
		data string
		err  error
	}
	capture := make(chan result, 1)
	go func() { var b bytes.Buffer; _, err := io.Copy(&b, master); capture <- result{b.String(), err} }()
	view := NewMenu(slave, "Unicode 界", "")
	view.Configure(Presentation{Locale: "en", Density: "compact", Mouse: "off"})
	active, err := view.BeginShell(slave)
	if err != nil || !active {
		t.Fatal(active, err)
	}
	done := make(chan error, 1)
	go func() {
		n, supported, err := view.ChooseInteractive(slave, "Actions", []string{"inspect e\u0301 界 🙂"}, []string{strings.Repeat("界 e\u0301 🙂 preview\n", 30)})
		if err == nil && (!supported || n != 0) {
			err = fmt.Errorf("unexpected selection %d", n)
		}
		done <- err
	}()
	// Acknowledge raw mode, not merely a fixed scheduling delay.
	deadline := time.Now().Add(time.Second)
	for {
		term, e := interactiveTermios(slave)
		if e != nil {
			t.Fatal(e)
		}
		if term.Lflag&syscall.ICANON == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("raw mode not reached")
		}
		time.Sleep(time.Millisecond)
	}
	if _, err = master.Write([]byte("\x1b[B")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(130 * time.Millisecond)
	resize(12, 6)
	time.Sleep(130 * time.Millisecond)
	resize(8, 2)
	time.Sleep(130 * time.Millisecond)
	resize(32, 12)
	time.Sleep(130 * time.Millisecond)
	if _, err = master.Write([]byte("q")); err != nil {
		t.Fatal(err)
	}
	select {
	case err = <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		master.Write([]byte{3})
		select {
		case <-done:
		case <-time.After(time.Second):
		}
		t.Fatal("resized preview did not exit")
	}
	view.CloseShell()
	after, err := interactiveTermios(slave)
	if err != nil || after != before {
		t.Fatal("termios changed after Unicode/resize session", err)
	}
	slave.Close()
	var output string
	select {
	case result := <-capture:
		if result.err != nil && !errors.Is(result.err, syscall.EIO) {
			t.Fatal(result.err)
		}
		output = result.data
	case <-time.After(time.Second):
		t.Fatal("PTY capture did not end")
	}
	if !utf8.ValidString(output) {
		t.Fatal("resize emitted invalid UTF-8")
	}
	if regexp.MustCompile(`\x1b\[[0-9;]*m`).MatchString(strings.ReplaceAll(output, "\x1b[0m", "")) {
		t.Fatal("NO_COLOR live preview emitted SGR")
	}
	if strings.Count(output, "\x1b[?1049h") != 1 || strings.Count(output, "\x1b[?1049l") != 1 {
		t.Fatal("resize restarted/leaked alternate screen")
	}
	if !strings.Contains(output, "\x1b[6;1H") || !strings.Contains(output, "\x1b[2;1H") || !strings.Contains(output, "\x1b[12;1H") {
		t.Fatal("actual resized footer geometry not observed")
	}
}

// Resize below startup's supported width after opening a real shell: every
// positioned text run must stay inside the actual terminal, even with no body.
func TestUnicodeResizeTinyFrameBounds(t *testing.T) {
	for _, width := range []uint16{1, 3, 4} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			t.Setenv("TERM", "xterm-256color")
			t.Setenv("NO_COLOR", "1")
			tty, finish := capturePTY(t, 80, 24)
			view := NewMenu(tty, "界界界 Project", "context")
			view.Configure(Presentation{Density: "compact", Lettering: "text"})
			active, err := view.BeginShell(tty)
			if err != nil || !active {
				t.Fatal(active, err)
			}
			size := struct{ rows, cols, x, y uint16 }{rows: 6, cols: width}
			if err := ioctl(tty, uintptr(syscall.TIOCSWINSZ), unsafe.Pointer(&size)); err != nil {
				t.Fatal(err)
			}
			view.ShellRender("", []string{"界界 e\u0301 🙂 long"}, "")
			actual, _, _ := view.ShellGeometry()
			if actual != int(width) {
				t.Fatalf("actual width %d, want %d", actual, width)
			}
			view.CloseShell()
			output := finish()
			// The resize begins with resetting the scroll region; ignore startup at 80.
			at := strings.LastIndex(output, "\x1b[r\x1b[1;1H\x1b[2K")
			if at < 0 {
				t.Fatal("resize redraw missing")
			}
			output = output[at:]
			if !utf8.ValidString(output) {
				t.Fatal("invalid UTF-8")
			}
			cursor := regexp.MustCompile(`\x1b\[([0-9]+);([0-9]+)H`)
			for _, run := range cursor.Split(output, -1) {
				plain := unicodeFixtureANSI.ReplaceAllString(run, "")
				if cells := unicodeFixtureCells(plain); cells > int(width) {
					t.Errorf("text run has %d cells at width %d: %q", cells, width, plain)
				}
			}
			for _, address := range cursor.FindAllStringSubmatch(output, -1) {
				var row, col int
				fmt.Sscan(address[1], &row)
				fmt.Sscan(address[2], &col)
				if row < 1 || row > 6 || col < 1 || col > int(width) {
					t.Errorf("invalid actual-size address: %q", address[0])
				}
			}
		})
	}
}
