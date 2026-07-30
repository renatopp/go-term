package term

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

var writer io.Writer = os.Stdout
var stdin uintptr = os.Stdin.Fd()
var stdout uintptr = os.Stdout.Fd()
var bus = newEventBus()
var rawModeState *State
var stopStdinReader func()
var stdinSource io.Reader = os.Stdin
var resizePollInterval = 100 * time.Millisecond
var screenSizeGetter = GetScreenSize
var colorLevel ColorMode = ColorModeNone

type State = term.State

func init() {
	colorLevel = detectColorMode()
}

// startStdinReader reads stdin and publishes decoded events (MouseEvent,
// KeyEvent, CursorPositionEvent) on the package bus until stop is called.
// Call the returned stop function to stop reading.
//
// If the terminal never sends another byte, the underlying read blocks
// forever and the goroutine outlives stop(); the process exiting is what
// ultimately reclaims it, the same tradeoff already made for other blocking
// stdin reads in this package.
func startStdinReader() (stop func()) {
	done := make(chan struct{})
	r := bufio.NewReader(stdinSource)

	go func() {
		for {
			eventType, event, err := readInputEvent(r)
			if err != nil {
				return
			}
			select {
			case <-done:
				return
			default:
				bus.Publish(eventType, event)
			}
		}
	}()

	return func() { close(done) }
}

func write(s string) error {
	_, err := writer.Write([]byte(s))
	return err
}

// readInputEvent reads and decodes the next event from r: a key press, a
// mouse event, or a cursor position report.
func readInputEvent(r *bufio.Reader) (eventType string, event Event, err error) {
	ch, _, err := r.ReadRune()
	if err != nil {
		return "", nil, err
	}

	if ch == 0x1b {
		return readEscapeSequence(r)
	}

	return EventKey, decodeControlKey(ch), nil
}

// decodeControlKey decodes a single non-escape rune read from stdin into a
// KeyEvent: named keys (enter, tab, backspace), Ctrl+letter combinations
// (bytes 0x01-0x1a), or a plain rune.
func decodeControlKey(ch rune) KeyEvent {
	switch ch {
	case '\r', '\n':
		return KeyEvent{Type: KeyEnter}
	case '\t':
		return KeyEvent{Type: KeyTab}
	case 0x7f, 0x08:
		return KeyEvent{Type: KeyBackspace}
	}
	if ch >= 1 && ch <= 26 {
		return KeyEvent{Type: KeyRune, Rune: 'a' + ch - 1, Ctrl: true}
	}
	return KeyEvent{Type: KeyRune, Rune: ch}
}

// readEscapeSequence decodes the byte(s) following an ESC (0x1b) already
// read from r: a lone Escape key, an Alt+key combination, an SGR mouse
// event, or a CSI key sequence (arrows, Home/End/Delete/Insert/PgUp/PgDown,
// or a cursor position report).
//
// A lone Escape key and the start of a multi-byte escape sequence both
// begin with the same 0x1b byte; distinguishing them relies on terminals
// writing every byte of a sequence in a single burst, so if nothing is
// buffered yet, we treat it as a standalone Escape rather than blocking
// until the next keypress arrives.
func readEscapeSequence(r *bufio.Reader) (eventType string, event Event, err error) {
	if r.Buffered() == 0 {
		return EventKey, KeyEvent{Type: KeyEsc}, nil
	}

	ch, _, err := r.ReadRune()
	if err != nil {
		return "", nil, err
	}

	if ch != '[' {
		return EventKey, KeyEvent{Type: KeyRune, Rune: ch, Alt: true}, nil
	}

	final, _, err := r.ReadRune()
	if err != nil {
		return "", nil, err
	}

	if final == '<' {
		mouse, ok, err := parseSGRMouse(r)
		if err != nil {
			return "", nil, err
		}
		if !ok {
			return EventKey, KeyEvent{Type: KeyEsc}, nil
		}
		return EventMouse, mouse, nil
	}

	var params []rune
	for isCSIParamRune(final) {
		params = append(params, final)
		if final, _, err = r.ReadRune(); err != nil {
			return "", nil, err
		}
	}

	return decodeCSI(string(params), final)
}

func isCSIParamRune(ch rune) bool {
	return (ch >= '0' && ch <= '9') || ch == ';'
}

// decodeCSI decodes a CSI sequence's parameter string and final byte into an
// event. params is the accumulated digits/semicolons; final is the byte
// that terminated the sequence.
func decodeCSI(params string, final rune) (eventType string, event Event, err error) {
	switch final {
	case 'A':
		return EventKey, KeyEvent{Type: KeyUp}, nil
	case 'B':
		return EventKey, KeyEvent{Type: KeyDown}, nil
	case 'C':
		return EventKey, KeyEvent{Type: KeyRight}, nil
	case 'D':
		return EventKey, KeyEvent{Type: KeyLeft}, nil
	case 'H':
		return EventKey, KeyEvent{Type: KeyHome}, nil
	case 'F':
		return EventKey, KeyEvent{Type: KeyEnd}, nil
	case 'R':
		if row, col, ok := parseCursorPosition(params); ok {
			return EventCursorPosition, CursorPositionEvent{Row: row, Col: col}, nil
		}
	case '~':
		if key, ok := decodeTildeKey(params); ok {
			return EventKey, key, nil
		}
	}
	return EventKey, KeyEvent{Type: KeyEsc}, nil
}

func decodeTildeKey(params string) (KeyEvent, bool) {
	switch params {
	case "1", "7":
		return KeyEvent{Type: KeyHome}, true
	case "2":
		return KeyEvent{Type: KeyInsert}, true
	case "3":
		return KeyEvent{Type: KeyDelete}, true
	case "4", "8":
		return KeyEvent{Type: KeyEnd}, true
	case "5":
		return KeyEvent{Type: KeyPgUp}, true
	case "6":
		return KeyEvent{Type: KeyPgDown}, true
	}
	return KeyEvent{}, false
}

// parseCursorPosition parses the "row;col" parameter string of a cursor
// position report (CSI row ; col R).
func parseCursorPosition(params string) (row, col int, ok bool) {
	sep := -1
	for i, c := range params {
		if c == ';' {
			sep = i
			break
		}
	}
	if sep < 0 {
		return 0, 0, false
	}

	row, err1 := strconv.Atoi(params[:sep])
	col, err2 := strconv.Atoi(params[sep+1:])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return row, col, true
}

func parseSGRMouse(r *bufio.Reader) (event MouseEvent, ok bool, err error) {
	var params strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			return MouseEvent{}, false, err
		}
		if b == 'M' || b == 'm' {
			event, ok := decodeSGRMouse(params.String(), b == 'M')
			return event, ok, nil
		}
		params.WriteByte(b)
	}
}

// decodeSGRMouse parses the "Cb;Cx;Cy" parameter string of an SGR mouse
// sequence. pressed is true when the sequence was terminated by 'M' (press,
// drag, or wheel), false for 'm' (release).
func decodeSGRMouse(params string, pressed bool) (MouseEvent, bool) {
	parts := strings.SplitN(params, ";", 3)
	if len(parts) != 3 {
		return MouseEvent{}, false
	}

	cb, err1 := strconv.Atoi(parts[0])
	cx, err2 := strconv.Atoi(parts[1])
	cy, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return MouseEvent{}, false
	}

	event := MouseEvent{
		X:     cx - 1,
		Y:     cy - 1,
		Shift: cb&4 != 0,
		Alt:   cb&8 != 0,
		Ctrl:  cb&16 != 0,
	}

	switch {
	case cb&64 != 0:
		event.Action = MouseActionPress
		if cb&1 != 0 {
			event.Button = MouseButtonWheelDown
		} else {
			event.Button = MouseButtonWheelUp
		}
	case cb&32 != 0:
		event.Action = MouseActionMotion
		event.Button = mouseButtonFromCb(cb)
	case pressed:
		event.Action = MouseActionPress
		event.Button = mouseButtonFromCb(cb)
	default:
		event.Action = MouseActionRelease
		event.Button = mouseButtonFromCb(cb)
	}
	return event, true
}

func mouseButtonFromCb(cb int) MouseButton {
	switch cb & 3 {
	case 0:
		return MouseButtonLeft
	case 1:
		return MouseButtonMiddle
	case 2:
		return MouseButtonRight
	default:
		return MouseButtonNone
	}
}
