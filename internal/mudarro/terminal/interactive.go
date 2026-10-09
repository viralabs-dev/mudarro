//go:build linux || darwin || windows

package terminal

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// InteractiveState contains presentation state only. Navigation never runs actions.
type InteractiveState struct {
	Selected, Offset int
	Expanded         bool
}

func (s *InteractiveState) Navigate(key string, count, page, total int) (int, bool) {
	max := total - page
	if max < 0 {
		max = 0
	}
	switch key {
	case "q", "0":
		return 0, true
	case "enter":
		if count > 0 {
			return s.Selected + 1, true
		}
	case "down":
		s.Expanded = true
	case "up":
		s.Expanded = false
		s.Offset = 0
	case "shift-tab":
		if count > 0 {
			s.Selected = (s.Selected + count - 1) % count
			s.Offset = 0
		}
	case "tab":
		if count > 0 {
			s.Selected = (s.Selected + 1) % count
			s.Offset = 0
		}
	case "page-down":
		s.Offset += page
	case "page-up":
		s.Offset -= page
	case "home":
		s.Offset = 0
	case "end":
		s.Offset = max
	case "wheel-down":
		s.Offset += 3
	case "wheel-up":
		s.Offset -= 3
	default:
		if n, e := strconv.Atoi(key); e == nil && n > 0 && n <= count {
			s.Selected = n - 1
			s.Offset = 0
		}
	}
	if s.Offset < 0 {
		s.Offset = 0
	}
	if s.Offset > max {
		s.Offset = max
	}
	return 0, false
}

// ChooseInteractive returns supported=false without reading when either stream is not a TTY.
// Only Enter returns a selected action. Both raw mode and mouse reporting are restored first.
func (v *Menu) ChooseInteractive(in *os.File, title string, labels, previews []string) (index int, supported bool, err error) {
	out, ok := v.out.(*os.File)
	if !ok || in == nil || out == nil {
		return 0, false, nil
	}
	_, _, inputTTY := terminalSize(in)
	width, rows, outputTTY := terminalSize(out)
	if !inputTTY || !outputTTY || (terminalName() == "dumb" || terminalName() == "") {
		return 0, false, nil
	}
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupts)
	restore, e := interactiveRaw(in)
	if e != nil {
		return 0, false, e
	}
	supported = true
	defer func() {
		if e := restore(); err == nil {
			err = e
		}
	}()
	mouse := v.mouse == "on" || (v.mouse != "off" && interactiveMouseCapable(terminalName()))
	if !v.ShellActive() {
		fmt.Fprint(out, "\x1b[?1049h")
	}
	fmt.Fprint(out, "\x1b[?25l")
	if mouse {
		fmt.Fprint(out, "\x1b[?1002h\x1b[?1006h")
	}
	defer func() {
		if mouse {
			fmt.Fprint(out, "\x1b[?1006l\x1b[?1002l")
		}
		fmt.Fprint(out, "\x1b[?25h")
		if !v.ShellActive() {
			fmt.Fprint(out, "\x1b[?1049l")
		}
	}()
	state := InteractiveState{}
	pending := []byte{}
	digits := ""
	invalidDigits := false
	dragging := false
	for {
		select {
		case sig := <-interrupts:
			return 0, true, fmt.Errorf("menu interrompido: %s", sig)
		default:
		}
		width, rows, _ = terminalSize(out)
		centerTop := 1
		if v.ShellActive() {
			width, centerTop, rows = v.ShellGeometry()
		}
		if rows == 0 {
			key, more, err := interactiveRead(in, pending)
			pending = more
			if err != nil {
				return 0, true, err
			}
			if key == "q" || key == "0" {
				return 0, true, nil
			}
			continue
		}
		visible := len(labels)
		limit := (rows - 6) / 2
		if limit < 1 {
			limit = 1
		}
		if visible > limit {
			visible = limit
		}
		first := 0
		if state.Selected >= visible {
			first = state.Selected - visible + 1
		}
		page := rows - visible - 5
		if page < 1 {
			page = 1
		}
		preview := []string{}
		if state.Expanded && state.Selected < len(previews) {
			for _, line := range strings.SplitN(previews[state.Selected], "\n", 4096) {
				preview = append(preview, fitText(line, width-4))
			}
		}
		state.Navigate("", len(labels), page, len(preview))
		lines := []string{fitText(v.project+" / "+title, width-1)}
		for i := first; i < first+visible; i++ {
			marker := " "
			if i == state.Selected {
				marker = ">"
			}
			lines = append(lines, fmt.Sprintf("%s %d. %s", marker, i+1, fitText(labels[i], width-8)))
		}
		hint := v.text("Down preview · Up close · Tab/number focus · Enter execute · q back", "↓ prévia · ↑ recolher · Tab/número foco · Enter executar · q voltar")
		lines = append(lines, fitText(hint, width-1), "")
		top := centerTop + len(lines)
		if state.Expanded {
			for row := 0; row < page; row++ {
				line := ""
				n := state.Offset + row
				if n < len(preview) {
					line = preview[n]
				}
				if len(preview) > page {
					thumb := state.Offset * (page - 1) / (len(preview) - page)
					ch := "│"
					if row == thumb {
						ch = "█"
					}
					// Reserve the final cell for the drag target, without wrapping.
					line = fitText(line, width-2)
					cells := 0
					for _, r := range line {
						cells += runeColumns(r)
					}
					line += strings.Repeat(" ", max(0, width-1-cells)) + ch
				}
				lines = append(lines, line)
			}
		}
		if v.ShellActive() {
			v.ShellRender(title, lines, hint)
		} else {
			fmt.Fprint(out, "\x1b[2J\x1b[H")
			for i, line := range lines {
				fmt.Fprintf(out, "\x1b[%d;1H%s", i+1, line)
			}
		}

		key, more, readErr := interactiveRead(in, pending)
		pending = more
		if readErr != nil {
			return 0, true, readErr
		}
		if key == "" {
			continue
		}
		if strings.HasPrefix(key, "mouse:") {
			if !mouse {
				continue
			}
			var button, x, y int
			var final rune
			if _, e := fmt.Sscanf(key, "mouse:%d;%d;%d%c", &button, &x, &y, &final); e == nil {
				if final == 'm' {
					dragging = false
				}
				if y < top || y >= top+page {
					continue
				}
				if button&64 != 0 {
					if button&1 == 0 {
						key = "wheel-up"
					} else {
						key = "wheel-down"
					}
				} else {
					if final == 'm' {
						dragging = false
					}
					if x == width && y >= top && y < top+page && final == 'M' && button&3 == 0 {
						dragging = true
					}
					if dragging && final == 'M' && len(preview) > page {
						relative := y - top
						if relative < 0 {
							relative = 0
						}
						if relative >= page {
							relative = page - 1
						}
						if page > 1 {
							state.Offset = relative * (len(preview) - page) / (page - 1)
						}
					}
					continue
				}
			} else {
				continue
			}
		}
		if len(key) == 1 && key[0] >= '0' && key[0] <= '9' && (key != "0" || digits != "") {
			digits += key
			if len(digits) > 3 {
				digits = key
			}
			if n, e := strconv.Atoi(digits); e == nil && n > 0 && n <= len(labels) {
				state.Navigate(digits, len(labels), page, len(preview))
				invalidDigits = false
			} else {
				invalidDigits = true
			}
			continue
		}
		if key == "enter" && invalidDigits {
			digits = ""
			continue
		}
		invalidDigits = false
		digits = ""
		if selected, done := state.Navigate(key, len(labels), page, len(preview)); done {
			return selected, true, nil
		}
	}
}

func interactiveRead(in *os.File, pending []byte) (string, []byte, error) {
	for {
		if len(pending) > 0 && pending[0] != 27 {
			b := pending[0]
			rest := pending[1:]
			switch b {
			case 3:
				return "q", rest, fmt.Errorf("menu interrompido")
			case 13, 10:
				return "enter", rest, nil
			case 9:
				return "tab", rest, nil
			}
			return string(b), rest, nil
		}
		if len(pending) >= 3 && pending[0] == 27 {
			for i := 2; i < len(pending); i++ {
				c := pending[i]
				if c >= 0x40 && c <= 0x7e {
					seq := string(pending[:i+1])
					rest := pending[i+1:]
					switch seq {
					case "\x1b[A":
						return "up", rest, nil
					case "\x1b[Z":
						return "shift-tab", rest, nil
					case "\x1b[B":
						return "down", rest, nil
					case "\x1b[5~":
						return "page-up", rest, nil
					case "\x1b[6~":
						return "page-down", rest, nil
					case "\x1b[H", "\x1b[1~":
						return "home", rest, nil
					case "\x1b[F", "\x1b[4~":
						return "end", rest, nil
					}
					if strings.HasPrefix(seq, "\x1b[<") {
						return "mouse:" + seq[3:], rest, nil
					}
					return "", rest, nil
				}
			}
		}
		if len(pending) > 64 {
			return "", nil, nil
		}
		chunk, e := interactiveNext(in)
		if e != nil {
			return "", nil, e
		}
		if len(chunk) == 0 {
			return "", pending, nil
		}
		pending = append(pending, chunk...)
	}
}

func interactiveMouseCapable(term string) bool {
	for _, prefix := range []string{"xterm", "screen", "tmux", "kitty", "wezterm"} {
		if strings.HasPrefix(term, prefix) {
			return true
		}
	}
	return false
}
