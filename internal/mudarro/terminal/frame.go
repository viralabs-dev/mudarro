package terminal

import (
	"fmt"
	"io"
	"os"
	"sync"
	"unicode"
	"unicode/utf8"
)

// shellFrame owns terminal chrome. All application output is rendered as text in its central region.
type shellFrame struct {
	mu                       sync.Mutex
	menu                     *Menu
	out                      *os.File
	closed                   bool
	width, rows, top, height int
	signature                string
	generation               uint64
	in                       *os.File
	headerContext            string
}

func (v *Menu) BeginShell(in *os.File) (bool, error) {
	if v.ShellActive() {
		return true, nil
	}
	out, ok := v.out.(*os.File)
	if !ok || in == nil || out == nil || os.Getenv("TERM") == "" || os.Getenv("TERM") == "dumb" {
		return false, nil
	}
	_, _, inputTTY := frameTerminalSize(in)
	width, rows, outputTTY := frameTerminalSize(out)
	if !inputTTY || !outputTTY || width < 4 || rows < 3 {
		return false, nil
	}
	f := &shellFrame{menu: v, out: out, in: in, headerContext: v.context}
	v.frame = f
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, err := fmt.Fprint(out, "\x1b[?1049h\x1b[?7l\x1b[?25l\x1b[2J\x1b[H"); err != nil {
		f.closed = true
		fmt.Fprint(out, "\x1b[r\x1b[?7h\x1b[?25h\x1b[?1049l")
		return false, err
	}
	f.ensure()
	return true, nil
}
func (v *Menu) CloseShell() {
	f := v.frame
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return
	}
	f.closed = true
	fmt.Fprint(f.out, "\x1b[?1006l\x1b[?1002l\x1b[r\x1b[0m\x1b[?7h\x1b[?25h\x1b[?1049l")
}
func (v *Menu) ShellActive() bool {
	f := v.frame
	if f == nil {
		return false
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return !f.closed
}
func (v *Menu) ShellGeometry() (width, top, height int) {
	f := v.frame
	if f == nil {
		return v.style.width, 1, v.style.rows
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.closed {
		f.ensure()
	}
	return f.width, f.top, f.height
}
func (v *Menu) ShellRender(title string, lines []string, hint string) {
	f := v.frame
	if f == nil {
		fmt.Fprintln(v.out, cleanText(title))
		for _, line := range lines {
			fmt.Fprintln(v.out, cleanText(line))
		}
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return
	}
	f.ensure()
	f.render(title, lines, hint)
}
func (v *Menu) ShellOutput() io.Writer {
	f := v.frame
	if f == nil {
		return v.out
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.generation++
	if !f.closed {
		f.ensure()
		f.render("", []string{v.text("Execution", "Execução")}, "")
	}
	return &frameOutput{frame: f, generation: f.generation}
}

func (f *shellFrame) ensure() {
	width, rows, ok := frameTerminalSize(f.out)
	if !ok || width < 1 || rows < 1 {
		width, rows = f.width, f.rows
		if width == 0 {
			width, rows = 80, 24
		}
	}
	v := f.menu
	v.style.width = width
	v.style.rows = rows
	mark := v.mark()
	header := append(append([]string{}, mark...), " "+v.style.Hint(fitText(v.text("PROJECT / ", "PROJETO / ")+v.project, width-2)), " "+v.style.Hint(fitText(f.headerContext, width-2)))
	footerRows := 2
	if len(header)+footerRows+4 > rows {
		header = []string{" " + v.style.heading(fitText(v.project, width-2))}
		if rows >= 7 {
			header = append(header, " "+v.style.Hint(fitText(f.headerContext, width-2)))
		}
	}
	if rows < 5 {
		footerRows = 1
	}
	if rows < 3 {
		header = nil
		if rows == 2 {
			header = []string{v.style.heading(fitText(v.project, width))}
		}
	}
	top := len(header) + 1
	height := rows - len(header) - footerRows
	signature := fmt.Sprintf("%d/%d/%s/%s/%s/%s", width, rows, v.project, f.headerContext, v.locale, v.style.theme)
	changed := signature != f.signature
	f.width, f.rows, f.top, f.height = width, rows, top, height
	if !changed {
		return
	}
	f.signature = signature
	// Absolute positioning does not scroll the header/footer. Only a multi-row body gets DECSTBM.
	fmt.Fprint(f.out, "\x1b[r")
	for row := 1; row <= rows; row++ {
		fmt.Fprintf(f.out, "\x1b[%d;1H\x1b[2K", row)
	}
	for index, line := range header {
		fmt.Fprintf(f.out, "\x1b[%d;1H%s", index+1, line)
	}
	if height > 1 {
		fmt.Fprintf(f.out, "\x1b[%d;%dr", top, top+height-1)
	}
	f.footer("")
}
func (f *shellFrame) footer(hint string) {
	v := f.menu
	start := f.top + f.height
	if start < f.rows {
		fmt.Fprintf(f.out, "\x1b[%d;1H\x1b[2K%s", start, v.style.rule())
	}
	hint = v.text("Mudarro · Enter confirm · q back · Ctrl+C cancel", "Mudarro · Enter confirma · q volta · Ctrl+C cancela")
	fmt.Fprintf(f.out, "\x1b[%d;1H\x1b[2K%s", f.rows, fitText(hint, f.width-1))
}
func (f *shellFrame) render(title string, lines []string, hint string) {
	if f.height == 0 {
		return
	}
	for row := f.top; row < f.top+f.height; row++ {
		fmt.Fprintf(f.out, "\x1b[%d;1H\x1b[2K", row)
	}
	body := lines
	for index, line := range body {
		if index >= f.height {
			break
		}
		fmt.Fprintf(f.out, "\x1b[%d;1H%s", f.top+index, fitText(line, f.width))
	}
	row := f.top + len(body)
	if row >= f.top+f.height {
		row = f.top + f.height - 1
	}
	if row < 1 {
		row = 1
	}
	fmt.Fprintf(f.out, "\x1b[%d;1H", row)
}

// frameOutput has a streaming control parser: escape payloads are never buffered or emitted.
// A line is bounded to 4 KiB and history to 256 lines; old writers cannot overwrite a newer action.
type frameOutput struct {
	frame      *shellFrame
	generation uint64
	state      byte
	tail       []byte
	line       string
	history    []string
	carriage   bool
}

func (w *frameOutput) Write(p []byte) (int, error) {
	f := w.frame
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed || w.generation != f.generation {
		return len(p), nil
	}
	for _, b := range p {
		w.consume(b)
	}
	f.ensure()
	lines := append(append([]string{}, w.history...), w.line)
	available := f.height - 1
	if available < 1 {
		available = 1
	}
	if len(lines) > available {
		lines = lines[len(lines)-available:]
	}
	f.render("", append([]string{f.menu.text("Execution", "Execução")}, lines...), "")
	return len(p), nil
}
func (w *frameOutput) consume(b byte) {
	if w.carriage {
		if b != '\n' {
			w.line = ""
		}
		w.carriage = false
	}
	switch w.state {
	case 1:
		switch b {
		case '[':
			w.state = 2
		case ']':
			w.state = 3
		case 'P', 'X', '^', '_':
			w.state = 4
		case 27:
			w.state = 1
		default:
			w.state = 0
		}
		return
	case 2:
		if b >= 0x40 && b <= 0x7e {
			w.state = 0
		} else if b == 27 {
			w.state = 1
		}
		return
	case 3:
		if b == 7 {
			w.state = 0
		} else if b == 27 {
			w.state = 5
		}
		return
	case 4:
		if b == 27 {
			w.state = 6
		}
		return
	case 5:
		if b == '\\' {
			w.state = 0
		} else if b == 7 {
			w.state = 0
		} else if b != 27 {
			w.state = 3
		}
		return
	case 6:
		if b == '\\' {
			w.state = 0
		} else if b != 27 {
			w.state = 4
		}
		return
	}
	if b == 27 {
		w.tail = nil
		w.state = 1
		return
	}
	if len(w.tail) == 0 {
		switch b {
		case '\n':
			w.history = append(w.history, w.line)
			w.line = ""
			if len(w.history) > 256 {
				w.history = w.history[len(w.history)-256:]
			}
			return
		case '\r':
			w.carriage = true
			return
		case '\b':
			if w.line != "" {
				_, size := utf8.DecodeLastRuneInString(w.line)
				w.line = w.line[:len(w.line)-size]
			}
			return
		case '\t':
			w.appendText("    ")
			return
		}
		if b < 32 || b == 127 {
			return
		}
	}
	w.tail = append(w.tail, b)
	if !utf8.FullRune(w.tail) {
		return
	}
	r, size := utf8.DecodeRune(w.tail)
	w.tail = w.tail[size:]
	if r == utf8.RuneError && size == 1 {
		return
	}
	if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
		return
	}
	w.appendText(string(r))
}
func (w *frameOutput) appendText(s string) {
	if len(w.line)+len(s) <= 4096 {
		w.line += s
	}
}

// Shell output is intentionally plain; arbitrary children opening /dev/tty bypass this presentation boundary.
var _ io.Writer = (*frameOutput)(nil)
