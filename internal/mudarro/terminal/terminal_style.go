package terminal

import (
	"io"
	"os"
	"strconv"
	"strings"
)

// Style is presentation only; plans, JSON and generated files stay plain.
type Style struct {
	color bool
	theme string
	width int
	rows  int
	tty   bool
	rich  bool
}

func styleFor(out io.Writer) Style {
	interactive := false
	terminalWidth, rows := 80, 24
	if f, ok := out.(*os.File); ok {
		terminalWidth, rows, interactive = terminalSize(f)
	}
	noColor := os.Getenv("NO_COLOR")
	term := os.Getenv("TERM")
	width, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || width < 20 || width > 240 {
		width = terminalWidth
	}
	return Style{color: colorAllowed(interactive, term, noColor), width: width, rows: rows, tty: interactive, rich: strings.Contains(term, "256color") || os.Getenv("COLORTERM") == "truecolor"}
}
func (s Style) Paint(code, text string) string {
	if !s.color {
		return text
	}
	if s.theme == "light" {
		switch code {
		case "1;36":
			code = "1;34"
		case "36":
			code = "34"
		case "90":
			code = "30"
		case "33", "38;5;208":
			code = "35"
		case "1;39":
			code = "1;30"
		}
	} else if s.theme == "dark" {
		if code == "1;39" {
			code = "1;37"
		}
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}
func (s Style) heading(text string) string { return s.Paint("1;36", text) }
func (s Style) Hint(text string) string    { return s.Paint("90", text) }
func (s Style) key(text string) string {
	code := "33"
	if s.rich {
		code = "38;5;208"
	}
	return s.Paint(code, text)
}
func (s Style) rule() string {
	width := s.width
	if width > 64 {
		width = 64
	}
	return s.Paint("36", strings.Repeat("-", width))
}

func colorAllowed(tty bool, term, noColor string) bool {
	return tty && noColor == "" && term != "" && term != "dumb"
}

func (s Style) HasColor() bool { return s.color }
