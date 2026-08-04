package term

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

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

// activeKeyboardProtocol is the keyboard-modifier protocol negotiated by the
// most recent EnterRawMode call, used by ExitRawMode to send the matching
// disable sequence.
var activeKeyboardProtocol = keyboardProtocolNone

// kittyProbeTimeout bounds how long negotiateKeyboardProtocol waits for a
// reply to the Kitty keyboard protocol query.
const kittyProbeTimeout = 100 * time.Millisecond

const (
	kittyQuerySeq             = "\x1b[?u"
	kittyEnableSeq            = "\x1b[>1u"
	kittyDisableSeq           = "\x1b[<1u"
	modifyOtherKeysEnableSeq  = "\x1b[>4;2m"
	modifyOtherKeysDisableSeq = "\x1b[>4;0m"
)

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

// negotiateKeyboardProtocol probes the terminal for Kitty keyboard protocol
// support and enables the richest keyboard-modifier protocol available, so
// Ctrl/Alt/Shift can be reported for keys that carry no room for modifiers
// in plain VT100 input (Enter, Tab, Backspace, Esc, Space, and printable
// keys). Falls back to xterm's modifyOtherKeys when Kitty isn't supported,
// or can't be probed because stdin isn't a file, or the file doesn't
// support read deadlines on this platform.
//
// The probe reads directly from stdin before startStdinReader begins; a
// keystroke racing the terminal's own query reply within the probe window
// is dropped, a known tradeoff for a synchronous capability check done once
// per EnterRawMode call.
//
// Must be called while holding mu (see EnterRawMode), and writes to stdout
// directly rather than through write to avoid deadlocking on it.
func negotiateKeyboardProtocol() keyboardProtocol {
	f, ok := stdin.(*os.File)
	if !ok {
		stdout.Write([]byte(modifyOtherKeysEnableSeq))
		return keyboardProtocolModifyOtherKeys
	}

	stdout.Write([]byte(kittyQuerySeq))
	if err := f.SetReadDeadline(time.Now().Add(kittyProbeTimeout)); err != nil {
		stdout.Write([]byte(modifyOtherKeysEnableSeq))
		return keyboardProtocolModifyOtherKeys
	}

	buf := make([]byte, 32)
	n, err := f.Read(buf)
	f.SetReadDeadline(time.Time{})

	if err == nil && isKittyQueryReply(buf[:n]) {
		stdout.Write([]byte(kittyEnableSeq))
		return keyboardProtocolKitty
	}

	stdout.Write([]byte(modifyOtherKeysEnableSeq))
	return keyboardProtocolModifyOtherKeys
}

// isKittyQueryReply reports whether b looks like a Kitty keyboard protocol
// query reply (CSI ? flags u).
func isKittyQueryReply(b []byte) bool {
	return len(b) >= 3 && b[0] == 0x1b && b[1] == '[' && b[2] == '?'
}

// disableKeyboardProtocol sends the disable sequence matching the protocol
// negotiated by negotiateKeyboardProtocol, restoring the terminal's default
// key reporting. Writes to stdout directly rather than through write to
// avoid deadlocking on it; see ExitRawMode.
func disableKeyboardProtocol(p keyboardProtocol) {
	switch p {
	case keyboardProtocolKitty:
		stdout.Write([]byte(kittyDisableSeq))
	case keyboardProtocolModifyOtherKeys:
		stdout.Write([]byte(modifyOtherKeysDisableSeq))
	}
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
		return EventKey, KeyEvent{Type: KeyRune, Rune: ch, Alt: true, Shift: isShiftedRune(ch)}, nil
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
	return KeyEvent{Type: KeyRune, Rune: ch, Shift: isShiftedRune(ch)}
}

// isShiftedRune reports whether ch is a character the terminal only
// produces when Shift (or Caps Lock) is held, letting a Shift signal be
// recovered for keys read without an escape sequence, where the byte itself
// carries no separate modifier bit.
func isShiftedRune(ch rune) bool {
	return unicode.IsUpper(ch)
}

// isCSIParamRune returns true if ch is a valid character in a CSI sequence's
// parameter string: a digit, semicolon, or colon (the sub-parameter
// separator used by the Kitty keyboard protocol).
func isCSIParamRune(ch rune) bool {
	return (ch >= '0' && ch <= '9') || ch == ';' || ch == ':'
}

// arrowKeys maps a CSI letter final byte to the named key it reports.
var arrowKeys = map[rune]KeyType{
	'A': KeyUp,
	'B': KeyDown,
	'C': KeyRight,
	'D': KeyLeft,
	'H': KeyHome,
	'F': KeyEnd,
}

// namedKeyCodes maps the numeric key codes used by both the Kitty keyboard
// protocol and xterm's modifyOtherKeys to this package's named keys.
var namedKeyCodes = map[int]KeyType{
	13:  KeyEnter,
	9:   KeyTab,
	127: KeyBackspace,
	27:  KeyEsc,
}

// tildeKeys maps a CSI ~ sequence's leading numeric parameter to the named
// key it reports.
var tildeKeys = map[int]KeyType{
	1: KeyHome,
	7: KeyHome,
	2: KeyInsert,
	3: KeyDelete,
	4: KeyEnd,
	8: KeyEnd,
	5: KeyPgUp,
	6: KeyPgDown,
}

// decodeCSI decodes a CSI sequence's parameter string and final byte into an
// event. params is the accumulated digits/semicolons/colons; final is the
// byte that terminated the sequence.
func decodeCSI(params string, final rune) (eventType string, event Event, err error) {
	switch final {
	case 'A', 'B', 'C', 'D', 'H', 'F':
		return EventKey, decodeArrowKey(params, final), nil
	case 'R':
		if row, col, ok := parseCursorPosition(params); ok {
			return EventCursorPosition, CursorPositionEvent{Row: row, Col: col}, nil
		}
	case '~':
		if key, ok := decodeTildeKey(params); ok {
			return EventKey, key, nil
		}
	case 'u':
		if key, ok := decodeKittyKey(params); ok {
			return EventKey, key, nil
		}
	}
	return EventKey, KeyEvent{Type: KeyEsc}, nil
}

// splitParams splits a CSI parameter string on ';' into its numeric fields.
// Each field may itself carry ':'-separated sub-parameters (as used by the
// Kitty keyboard protocol); only the first sub-parameter of each field is
// used. Unparsable or empty fields decode to 0.
func splitParams(params string) []int {
	if params == "" {
		return nil
	}
	fields := strings.Split(params, ";")
	nums := make([]int, len(fields))
	for i, f := range fields {
		if colon := strings.IndexByte(f, ':'); colon >= 0 {
			f = f[:colon]
		}
		nums[i], _ = strconv.Atoi(f)
	}
	return nums
}

// decodeModifier decodes an xterm-style modifier parameter (1 + bitmask)
// into Shift, Alt, and Ctrl. A mod of 0 or less (absent) reports no
// modifiers.
func decodeModifier(mod int) (shift, alt, ctrl bool) {
	if mod <= 0 {
		return false, false, false
	}
	bits := mod - 1
	return bits&1 != 0, bits&2 != 0, bits&4 != 0
}

// decodeArrowKey decodes an arrow/Home/End CSI sequence's parameter string
// (empty, or "1;mod" carrying an xterm modifier) into a KeyEvent.
func decodeArrowKey(params string, final rune) KeyEvent {
	event := KeyEvent{Type: arrowKeys[final]}
	if nums := splitParams(params); len(nums) >= 2 {
		event.Shift, event.Alt, event.Ctrl = decodeModifier(nums[1])
	}
	return event
}

// keyEventFromCode decodes a Kitty/modifyOtherKeys numeric key code into a
// KeyEvent: a named key for control keys with no printable rune, or the rune
// itself otherwise.
func keyEventFromCode(code int) KeyEvent {
	if keyType, ok := namedKeyCodes[code]; ok {
		return KeyEvent{Type: keyType}
	}
	return KeyEvent{Type: KeyRune, Rune: rune(code)}
}

// decodeTildeKey decodes the parameter string of a CSI ~ sequence into a
// KeyEvent: a legacy nav-key id (optionally followed by an xterm modifier
// field), or an xterm modifyOtherKeys sequence ("27;mod;code"). Returns
// false if the parameter string is unrecognized.
func decodeTildeKey(params string) (KeyEvent, bool) {
	nums := splitParams(params)
	if len(nums) == 0 {
		return KeyEvent{}, false
	}

	if len(nums) == 3 && nums[0] == 27 {
		event := keyEventFromCode(nums[2])
		event.Shift, event.Alt, event.Ctrl = decodeModifier(nums[1])
		return event, true
	}

	keyType, ok := tildeKeys[nums[0]]
	if !ok {
		return KeyEvent{}, false
	}
	event := KeyEvent{Type: keyType}
	if len(nums) >= 2 {
		event.Shift, event.Alt, event.Ctrl = decodeModifier(nums[1])
	}
	return event, true
}

// decodeKittyKey decodes the parameter string of a CSI u sequence (Kitty
// keyboard protocol) into a KeyEvent. Returns false if the parameter string
// is unrecognized.
func decodeKittyKey(params string) (KeyEvent, bool) {
	nums := splitParams(params)
	if len(nums) == 0 {
		return KeyEvent{}, false
	}

	event := keyEventFromCode(nums[0])
	if len(nums) >= 2 {
		event.Shift, event.Alt, event.Ctrl = decodeModifier(nums[1])
	}
	return event, true
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
		X:     cx,
		Y:     cy,
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
