//go:build linux || darwin || windows

package terminal

import "io"

// ReadShellLine displays input only through the central writer; it never enables echo.
func (v *Menu) ReadShellLine(in io.Reader, out io.Writer) (string, error) {
	var line []byte
	b := make([]byte, 1)
	for {
		n, e := in.Read(b)
		if e != nil {
			return "", e
		}
		if n == 0 {
			continue
		}
		if b[0] == '\n' {
			out.Write([]byte{'\n'})
			return string(line), nil
		}
		if b[0] == 127 || b[0] == 8 {
			if len(line) > 0 {
				line = line[:len(line)-1]
			}
			continue
		}
		if b[0] >= 32 && len(line) < 4096 {
			line = append(line, b[0])
			out.Write(b)
		}
	}
}
