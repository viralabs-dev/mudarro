//go:build windows

package winapi

import (
	"os"
	"syscall"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

type coord struct{ X, Y int16 }
type smallRect struct{ Left, Top, Right, Bottom int16 }
type screenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

// inputRecord mirrors INPUT_RECORD: a 2-byte type, 2 bytes of padding and a
// 16-byte union. For KEY_EVENT the union starts with KEY_EVENT_RECORD.
type inputRecord struct {
	EventType uint16
	_         uint16
	Event     [16]byte
}

// ConsoleMode returns the console mode of f, or ok=false when f is not a console.
func ConsoleMode(f *os.File) (uint32, bool) {
	var mode uint32
	if f == nil || syscall.GetConsoleMode(syscall.Handle(f.Fd()), &mode) != nil {
		return 0, false
	}
	return mode, true
}

// SetConsoleMode applies mode to the console handle behind f.
func SetConsoleMode(f *os.File, mode uint32) error {
	r, _, err := procSetConsoleMode.Call(f.Fd(), uintptr(mode))
	return callErr(r, err)
}

// WindowSize returns the visible window (not the scrollback buffer) of an output console.
func WindowSize(f *os.File) (int, int, bool) {
	var info screenBufferInfo
	if f == nil {
		return 0, 0, false
	}
	r, _, _ := procGetConsoleScreenBufferInfo.Call(f.Fd(), uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 0, 0, false
	}
	return int(info.Window.Right-info.Window.Left) + 1, int(info.Window.Bottom-info.Window.Top) + 1, true
}

// EnableVirtualTerminalOutput turns on ANSI/VT processing (Windows 10 1511+).
// It reports false when the console cannot interpret escape sequences.
func EnableVirtualTerminalOutput(f *os.File) bool {
	mode, ok := ConsoleMode(f)
	if !ok {
		return false
	}
	if mode&EnableVTProcessing != 0 {
		return true
	}
	return SetConsoleMode(f, mode|EnableVTProcessing) == nil
}

// InputReady waits up to timeoutMS for a key press that produces characters.
// Focus, mouse, resize, key-up and modifier-only records are consumed, because
// ReadConsole would otherwise block on them while the handle stays signalled.
func InputReady(f *os.File, timeoutMS uint32) (bool, error) {
	h := syscall.Handle(f.Fd())
	for {
		event, err := syscall.WaitForSingleObject(h, timeoutMS)
		if err != nil {
			return false, err
		}
		if event == waitTimeout {
			return false, nil
		}
		var rec inputRecord
		var n uint32
		r, _, callErrno := procPeekConsoleInputW.Call(uintptr(h), uintptr(unsafe.Pointer(&rec)), 1, uintptr(unsafe.Pointer(&n)))
		if r == 0 {
			return false, callErr(r, callErrno)
		}
		if n == 0 {
			return false, nil
		}
		keyDown := *(*int32)(unsafe.Pointer(&rec.Event[0])) != 0
		char := *(*uint16)(unsafe.Pointer(&rec.Event[10]))
		if rec.EventType == keyEvent && keyDown && char != 0 {
			return true, nil
		}
		r, _, callErrno = procReadConsoleInputW.Call(uintptr(h), uintptr(unsafe.Pointer(&rec)), 1, uintptr(unsafe.Pointer(&n)))
		if r == 0 {
			return false, callErr(r, callErrno)
		}
		// Keep polling within the same bounded wait; a later wait will time out.
		timeoutMS = 0
	}
}

// ReadInput reads one console character (two UTF-16 units for a surrogate
// pair) and returns its UTF-8 bytes. Call it only after InputReady.
func ReadInput(f *os.File) ([]byte, error) {
	h := syscall.Handle(f.Fd())
	unit := make([]uint16, 1)
	var n uint32
	if err := syscall.ReadConsole(h, &unit[0], 1, &n, nil); err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	units := unit[:1]
	if utf16.IsSurrogate(rune(unit[0])) {
		low := make([]uint16, 1)
		if err := syscall.ReadConsole(h, &low[0], 1, &n, nil); err == nil && n == 1 {
			units = append(units, low[0])
		}
	}
	out := []byte{}
	for _, r := range utf16.Decode(units) {
		out = utf8.AppendRune(out, r)
	}
	return out, nil
}
