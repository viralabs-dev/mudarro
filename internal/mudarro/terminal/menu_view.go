package terminal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Menu struct {
	previews                          []string
	out                               io.Writer
	style                             Style
	project, context                  string
	printed                           bool
	locale, density, lettering, mouse string
}

func (v *Menu) Choose(reader *bufio.Reader, title string, labels []string) (int, error) {
	s := v.style
	if s.color {
		fmt.Fprint(v.out, "\x1b[2J\x1b[H")
		for _, line := range v.mark() {
			fmt.Fprintln(v.out, line)
		}
		fmt.Fprintln(v.out)
		fmt.Fprintln(v.out, "   "+s.Hint(fitText(v.text("PROJECT / ", "PROJETO / ")+v.project, s.width-4)))
		fmt.Fprintln(v.out, "   "+s.Hint(fitText(v.text("LOCAL OPERATIONS · NO AI", "OPERAÇÕES LOCAIS · SEM IA"), s.width-4)))
		fmt.Fprintln(v.out, "   "+s.heading(fitText(v.context, s.width-4)))
		fmt.Fprintln(v.out, s.rule())
	} else if !v.printed {
		if v.density == "compact" || v.lettering == "text" {
			fmt.Fprintln(v.out, fitText(v.project, s.width))
		} else {
			fmt.Fprintln(v.out, asciiWidth(v.project, s.width))
		}
	}
	v.printed = true
	fmt.Fprintln(v.out, "\n  "+s.heading(fitText(v.text("Choose an area / ", "Escolha uma área / ")+title, s.width-4))+"\n")
	for i, label := range labels {
		fmt.Fprintf(v.out, "  %s  %s\n", s.key(fmt.Sprintf("%2d.", i+1)), fitText(label, s.width-8))
	}
	fmt.Fprintln(v.out, "\n  "+s.Hint(v.text(" 0.  Back / exit", " 0.  Voltar / sair")))
	if s.color {
		// Position the footer only when content fits. Small windows keep natural flow.
		used := len(v.mark()) + 9 + len(labels)
		if s.rows > used+3 {
			fmt.Fprintf(v.out, "\x1b[%d;1H", s.rows-2)
		}
		fmt.Fprintln(v.out, s.rule())
		if s.width < 70 {
			fmt.Fprintln(v.out, "  "+s.Hint(fitText(v.text("number / 0 / q back", "número / 0 / q voltar"), s.width-3)))
		} else {
			fmt.Fprintln(v.out, "  "+s.key("1-"+strconv.Itoa(len(labels)))+s.Hint(v.text(" select   ", " escolher   "))+s.key("0 / q")+s.Hint(v.text(" back   ", " voltar   "))+s.key("Enter")+s.Hint(v.text(" confirm", " confirmar")))
		}
	}
	if len(v.previews) > 0 {
		fmt.Fprintln(v.out, v.text("  p<number> preview only; number executes", "  p<número> só prévia; número executa"))
	}
	for {
		fmt.Fprint(v.out, "  "+s.key("> "))
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "p") && len(v.previews) > 0 {
			n, err := strconv.Atoi(line[1:])
			if err == nil && n > 0 && n <= len(v.previews) {
				for _, line := range strings.Split(v.previews[n-1], "\n") {
					fmt.Fprintln(v.out, cleanText(line))
				}
				continue
			}
		}
		if strings.EqualFold(line, "q") {
			return 0, nil
		}
		n, err := strconv.Atoi(line)
		if err == nil && n >= 0 && n <= len(labels) {
			return n, nil
		}
		fmt.Fprintln(v.out, s.Paint("31", fitText(v.text("  Invalid option. Choose a listed number.", "  Opção inválida. Escolha um número da lista."), s.width)))
	}
}

// runeColumns is a conservative terminal-cell width policy, not a grapheme engine.
// Combining marks occupy zero cells; CJK/fullwidth and emoji occupy two.
func runeColumns(r rune) int {
	if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) {
		return 0
	}
	if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hangul, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) ||
		r >= 0xff01 && r <= 0xff60 || r >= 0xffe0 && r <= 0xffe6 || r >= 0x1f000 && r <= 0x1faff || r >= 0x2600 && r <= 0x27bf {
		return 2
	}
	return 1
}
func cleanText(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, text)
}
func fitText(text string, width int) string {
	text = cleanText(text)
	if width < 4 {
		width = 4
	}
	cells := 0
	for _, r := range text {
		cells += runeColumns(r)
	}
	if cells <= width {
		return text
	}
	var b strings.Builder
	cells = 0
	for _, r := range text {
		n := runeColumns(r)
		if cells+n > width-3 {
			break
		}
		b.WriteRune(r)
		cells += n
	}
	return b.String() + "..."
}

// Original bitmap lettering, independent of FIGlet fonts and third-party logos.
var largeGlyphs = map[rune]string{
	'M': "#   #|## ##|# # #|#   #|#   #|#   #|#   #",
	'U': "#   #|#   #|#   #|#   #|#   #|#   #| ### ",
	'D': "#### |#   #|#   #|#   #|#   #|#   #|#### ",
	'A': " ### |#   #|#   #|#####|#   #|#   #|#   #",
	'R': "#### |#   #|#   #|#### |# #  |#  # |#   #",
	'O': " ### |#   #|#   #|#   #|#   #|#   #| ### ",
}

func wordmark(name string, s Style) []string {
	name = fitText(strings.ToUpper(name), 64)
	letters := []rune(name)
	// Names outside the embedded alphabet remain readable as text, never question glyphs.
	fallback := func() []string { return []string{" " + s.key("▌") + " " + s.heading(fitText(name, s.width-4))} }
	for _, r := range letters {
		if _, ok := glyphs[r]; !ok {
			return fallback()
		}
	}
	scale := 1
	if len(letters)*12+5 <= s.width {
		scale = 2
	}
	max := (s.width - 5) / (6 * scale)
	if max < 1 {
		max = 1
	}
	if len(letters) > max {
		return fallback()
	}
	width := len(letters) * 6
	pixels := make([][]bool, 7)
	for y := range pixels {
		pixels[y] = make([]bool, width)
	}
	for i, r := range letters {
		raw, ok := largeGlyphs[r]
		if !ok {
			small, exists := glyphs[r]
			if !exists {
				small = "###|  #| # |   | # "
			}
			rows := strings.Split(small, "|")
			raw = "     | " + rows[0] + " | " + rows[1] + " | " + rows[2] + " | " + rows[3] + " | " + rows[4] + " |     "
		}
		for y, row := range strings.Split(raw, "|") {
			for x, r := range row {
				if r == '#' {
					pixels[y][i*6+x] = true
				}
			}
		}
	}
	lines := []string{}
	for y := 0; y < 8; y++ {
		var line strings.Builder
		line.WriteString(" " + s.key("▌") + " ")
		for x := 0; x < width+1; x++ {
			if y < 7 && x < width && pixels[y][x] {
				line.WriteString(s.Paint("1;39", strings.Repeat("█", scale)))
			} else if y > 0 && x > 0 && x-1 < width && pixels[y-1][x-1] {
				line.WriteString(s.Hint(strings.Repeat("░", scale)))
			} else {
				line.WriteString(strings.Repeat(" ", scale))
			}
		}
		lines = append(lines, strings.TrimRight(line.String(), " "))
	}
	return lines
}

// NewMenu binds presentation to its output stream; it does not run application actions.
func NewMenu(out io.Writer, project, context string) *Menu {
	return &Menu{out: out, style: styleFor(out), project: project, context: context}
}
func (v *Menu) SetContext(context string) { v.context = context }
func (v *Menu) Style() Style              { return v.style }
