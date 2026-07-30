package term

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// bus Event bus for publishing and subscribing to terminal events (keyboard, mouse, cursor position).
var bus = newEventBus()

// mu is a mutex to protect concurrent access to terminal state and settings.
var mu sync.Mutex

// stdout is the output stdout for terminal control sequences and messages.
var stdout io.Writer = os.Stdout

// stdin is the input stdin for terminal events (keyboard, mouse, cursor position).
var stdin io.Reader = os.Stdin

// stdoutFile are the file descriptors to check terminal size.
var stdoutFile uintptr = os.Stdout.Fd()

// stdinFile are the file desriptors to watch for input events (keyboard, mouse,
// cursor position) and check terminal.
var stdinFile uintptr = os.Stdin.Fd()

// rawModeState is the stored state before entering raw mode, used to restore
// terminal settings on exit.
var rawModeState *State = nil

// colorLevel is the detected color mode of the terminal (None, Basic, 256, TrueColor).
var colorLevel ColorMode = ColorModeNone

// stopStdinReader is the function to stop reading stdin events, set when
// EnterRawMode is called and cleared on ExitRawMode.
var stopStdinReader func() = nil

// stopResizeWatcher is the function to stop watching for terminal resize
// events, set when StartResizeWatcher is called and cleared on StopResizeWatcher.
var stopResizeWatcher func() = nil

// resizePollInterval is the interval at which terminal size is polled for changes.
const resizePollInterval = 100 * time.Millisecond

// State is an alias for term.State, representing the terminal state for raw
// mode and restoration.
type State = term.State

func init() {
	colorLevel = detectColorMode()
	startResizeWatcher()
}

// write writes the given string to the configured stdout and returns any error
// encountered.
func write(s string) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := stdout.Write([]byte(s))
	return err
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
	if stopStdinReader != nil {
		return stopStdinReader
	}

	done := make(chan struct{})
	r := bufio.NewReader(stdin)

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

// startResizeWatcher polls the terminal size and publishes ResizeEvent on the
// package bus whenever it changes. Call the returned stop function to stop
// polling.
func startResizeWatcher() (stop func()) {
	if stopResizeWatcher != nil {
		return stopResizeWatcher
	}

	done := make(chan struct{})
	lastWidth, lastHeight, _ := GetScreenSize()

	go func() {
		ticker := time.NewTicker(resizePollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				width, height, err := GetScreenSize()
				if err != nil || (width == lastWidth && height == lastHeight) {
					continue
				}
				lastWidth, lastHeight = width, height
				bus.Publish(EventResize, ResizeEvent{Width: width, Height: height})
			}
		}
	}()

	return func() { close(done) }
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

// isCSIParamRune returns true if ch is a valid character in a CSI sequence's
// parameter string: a digit or semicolon.
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

// decodeTildeKey decodes the parameter string of a CSI ~ sequence into a
// KeyEvent. Returns false if the parameter string is unrecognized.
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

// parseSGRMouse reads the parameter string of an SGR mouse sequence from r
// until the terminating 'M' or 'm' is found. It returns the decoded MouseEvent
// and whether the sequence was valid.
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

// mouseButtonFromCb returns the MouseButton corresponding to the lower two bits of cb.
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

// detectColorMode detects the terminal's color mode based on environment
// variables and terminal capabilities.
func detectColorMode() ColorMode {
	// FORCE_COLOR overrides everything.
	if v, ok := os.LookupEnv("FORCE_COLOR"); ok {
		switch strings.TrimSpace(strings.ToLower(v)) {
		case "0", "false":
			return ColorModeNone
		case "3":
			return ColorModeTrue
		case "2":
			return ColorMode256
		case "1", "", "true":
			return ColorModeAnsi
		default:
			if n, err := strconv.Atoi(v); err == nil {
				switch {
				case n >= 3:
					return ColorModeTrue
				case n == 2:
					return ColorMode256
				case n >= 1:
					return ColorModeAnsi
				}
			}
			return ColorModeAnsi
		}
	}

	// NO_COLOR always disables colors.
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return ColorModeNone
	}

	// BSD/macOS convention.
	if os.Getenv("CLICOLOR_FORCE") != "" &&
		os.Getenv("CLICOLOR_FORCE") != "0" {
		return ColorModeAnsi
	}

	if os.Getenv("CLICOLOR") == "0" {
		return ColorModeNone
	}

	// Not writing to a terminal.
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return ColorModeNone
	}

	// TrueColor hint.
	switch strings.ToLower(os.Getenv("COLORTERM")) {
	case "truecolor", "24bit":
		return ColorModeTrue
	}

	termEnv := strings.ToLower(os.Getenv("TERM"))
	windowsTerm := strings.ToLower(os.Getenv("WT_SESSION"))

	switch {
	case termEnv == "", termEnv == "dumb":
		return ColorModeNone

	case strings.Contains(termEnv, "direct"):
		return ColorModeTrue

	case windowsTerm != "":
		return ColorModeTrue

	case strings.Contains(termEnv, "kitty"),
		strings.Contains(termEnv, "wezterm"),
		strings.Contains(termEnv, "ghostty"):
		return ColorModeTrue

	case strings.Contains(termEnv, "256color"):
		return ColorMode256

	default:
		return ColorModeAnsi
	}
}
