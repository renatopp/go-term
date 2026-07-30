package term

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

// SetStdin sets the file to use for stdin. It is used by the package to read
// input events (keyboard, mouse, cursor position) and check terminal state.
func SetStdin(f *os.File) {
	mu.Lock()
	defer mu.Unlock()
	stdin = f
	stdinFile = f.Fd()
}

// SetStdout sets the file to use for stdout. It is used by the package to write
// terminal control sequences and messages and to check terminal size.
func SetStdout(f *os.File) {
	mu.Lock()
	defer mu.Unlock()
	stdout = f
	stdoutFile = f.Fd()
}

// EnterRawMode puts stdin into raw mode and starts reading it for mouse,
// keyboard, and other events (see startStdinReader), which are published on
// the package bus. Call ExitRawMode to restore the terminal and stop
// reading.
//
// Do not read from stdin while in raw mode, as it will interfere with the
// package's reading and event publishing. Use event callbacks instead.
func EnterRawMode() error {
	mu.Lock()
	defer mu.Unlock()
	if rawModeState != nil { // already in raw mode
		return nil
	}
	state, err := term.MakeRaw(int(stdinFile))
	if err != nil {
		return err
	}
	rawModeState = state
	stopStdinReader = startStdinReader()
	return nil
}

// ExitRawMode restores the terminal state to what it was before EnterRawMode
// and stops reading stdin for events. It is safe to call even if not in raw mode.
func ExitRawMode() error {
	mu.Lock()
	defer mu.Unlock()
	state := rawModeState
	rawModeState = nil
	if stopStdinReader != nil {
		stopStdinReader()
		stopStdinReader = nil
	}
	return RestoreTerminalState(state)
}

// WithinRawMode runs the given function while stdin is in raw mode. If stdin is
// not already in raw mode, it will enter raw mode before running the function and
// restore the terminal state afterwards.
//
// Just like EnterRawMode, do not read from stdin while in raw mode, as it will
// interfere with the package's event reading. Use event callbacks instead.
func WithinRawMode(f func()) error {
	if rawModeState != nil {
		f()
		return nil
	}
	if err := EnterRawMode(); err != nil {
		return err
	}
	defer ExitRawMode()
	f()
	return nil
}

// RestoreTerminalState restores the terminal state to the given state.
func RestoreTerminalState(state *State) error {
	return term.Restore(int(stdinFile), state)
}

// GetTerminalState returns the current terminal state, which can be used to
// restore the terminal later with RestoreTerminalState.
func GetTerminalState() (*State, error) {
	return term.GetState(int(stdinFile))
}

// SetWindowTitle sets the terminal window title to the given string. It may
// not work on all terminals.
func SetWindowTitle(title string) error {
	return write("\x1b]0;" + title + "\x07")
}

// EnterAlternateScreen switches the terminal to the alternate screen buffer,
// which is a separate screen that can be used for full-screen applications.
// The original screen is restored when ExitAlternateScreen is called.
func EnterAlternateScreen() error {
	return write("\x1b[?1049h")
}

// ExitAlternateScreen switches the terminal back to the original screen buffer,
// restoring the screen contents that were present before EnterAlternateScreen
// was called.
func ExitAlternateScreen() error {
	return write("\x1b[?1049l")
}

// GetScreenSize returns the current width and height of the terminal in columns
// and rows. It returns an error if the size cannot be determined.
func GetScreenSize() (width, height int, err error) {
	return term.GetSize(int(stdoutFile))
}

// ForceGetScreenSize returns the current width and height of the terminal in
// columns and rows, ignoring any errors. If the size cannot be determined, it
// returns zero for both width and height.
func ForceGetScreenSize() (width, height int) {
	width, height, _ = GetScreenSize()
	return
}

// ClearScreen clears the entire terminal screen and moves the cursor to the
// top-left corner.
func ClearScreen() error {
	return write("\x1b[2J\x1b[H")
}

// ClearScreenAroundCursor clears the entire screen after the cursor, leaving
// the cursor position unchanged.
func ClearScreenAfterCursor() error {
	return write("\x1b[0J")
}

// ClearScreenBeforeCursor clears the entire screen before the cursor, leaving
// the cursor position unchanged.
func ClearScreenBeforeCursor() error {
	return write("\x1b[1J")
}

// ClearLine clears the entire current line and moves the cursor to the
// beginning of the line. It does not move the cursor vertically.
func ClearLine() error {
	return write("\x1b[2K\x1b[1G")
}

// ClearLineAfterCursor clears the line from the current cursor position to the
// end of the line, leaving the cursor position unchanged.
func ClearLineAfterCursor() error {
	return write("\x1b[0K")
}

// ClearLineBeforeCursor clears the line from the beginning of the line to the
// current cursor position, leaving the cursor position unchanged.
func ClearLineBeforeCursor() error {
	return write("\x1b[1K")
}

// ClearLineAroundCursor clears the line from the beginning of the line to the
// end of the line, leaving the cursor position unchanged.
func ClearLineAroundCursor() error {
	return write("\x1b[2K")
}

// MoveCursorUp moves the cursor up by n rows. If n is zero, it does nothing.
func MoveCursorUp(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "A")
}

// MoveCursorDown moves the cursor down by n rows. If n is zero, it does nothing.
func MoveCursorDown(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "B")
}

// MoveCursorForward moves the cursor forward (right) by n columns. If n is
// zero, it does nothing.
func MoveCursorForward(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "C")
}

// MoveCursorBackward moves the cursor backward (left) by n columns. If n is
// zero, it does nothing.
func MoveCursorBackward(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "D")
}

// MoveCursorTo moves the cursor to the specified row and column (1-based). If
// either row or col is zero, it does nothing.
func MoveCursorTo(row, col int) error {
	return write("\x1b[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H")
}

// MoveCursorToRow moves the cursor to the specified row (1-based) without
// changing the column. If row is zero, it does nothing.
func MoveCursorToRow(row int) error {
	return write("\x1b[" + strconv.Itoa(row) + "d")
}

// MoveCursorToColumn moves the cursor to the specified column (1-based) without
// changing the row. If col is zero, it does nothing.
func MoveCursorToColumn(col int) error {
	return write("\x1b[" + strconv.Itoa(col) + "G")
}

// MoveCursorToHome moves the cursor to the top-left corner of the screen
// (row 1, column 1).
func MoveCursorToHome() error {
	return write("\x1b[H")
}

// MoveCursorToStartOfLine moves the cursor to the beginning of the current line
// (column 1) without changing the row.
func MoveCursorToStartOfLine() error {
	return write("\x1b[1G")
}

// MoveCursorToEndOfLine moves the cursor to the end of the current line
// (column 9999) without changing the row.
func MoveCursorToEndOfLine() error {
	return write("\x1b[9999C")
}

// MoveCursorToBottom moves the cursor to the bottom of the screen (row 9999)
// without changing the column.
func MoveCursorToBottom() error {
	return write("\x1b[9999B")
}

// SaveCursorPosition saves the current cursor position, which can be restored
// later with RestoreCursorPosition.
func SaveCursorPosition() error {
	return write("\x1b[s")
}

// RestoreCursorPosition restores the cursor position to the last saved
// position using SaveCursorPosition. If no position was saved, it does nothing.
func RestoreCursorPosition() error {
	return write("\x1b[u")
}

// HideCursor hides the cursor, making it invisible. It does not move the cursor
// or change its position. Call ShowCursor to make the cursor visible again.
func HideCursor() error {
	return write("\x1b[?25l")
}

// ShowCursor shows the cursor, making it visible. It does not move the cursor
// or change its position. Call HideCursor to make the cursor invisible again.
func ShowCursor() error {
	return write("\x1b[?25h")
}

// GetCursorPosition queries the terminal for the current cursor position and
// returns the row and column (1-based). It returns an error if the position
// cannot be determined or if the terminal does not respond in time.
func GetCursorPosition() (row, col int, oerr error) {
	WithinRawMode(func() {
		got := make(chan CursorPositionEvent, 1)
		id := bus.Subscribe(EventCursorPosition, func(event Event) {
			if m, ok := event.(CursorPositionEvent); ok {
				select {
				case got <- m:
				default:
				}
			}
		})
		defer bus.Unsubscribe(EventCursorPosition, id)

		if err := write("\x1b[6n"); err != nil {
			oerr = err
			return
		}

		select {
		case m := <-got:
			row, col = m.Row, m.Col
		case <-time.After(200 * time.Millisecond):
			oerr = fmt.Errorf("timeout waiting for cursor position response")
		}
	})
	return row, col, oerr
}

// ForceGetCursorPosition queries the terminal for the current cursor position and returns the row and column (1-based), ignoring any errors. If the position
// cannot be determined, it returns zero for both row and column.
func ForceGetCursorPosition() (row, col int) {
	row, col, _ = GetCursorPosition()
	return
}

// EnableCursorWrap enables line wrapping at the right edge of the terminal.
// When enabled, the cursor moves to the beginning of the next line when it
// reaches the right edge of the terminal. When disabled, the cursor stays at
// the right edge and overwrites existing characters.
func EnableCursorWrap() error {
	return write("\x1b[?7h")
}

// DisableCursorWrap disables line wrapping at the right edge of the terminal.
func DisableCursorWrap() error {
	return write("\x1b[?7l")
}

// ScrollUp scrolls the terminal content up by n lines, moving the cursor down
// by n lines. If n is zero, it does nothing. New lines are added at the bottom
// of the screen.
func ScrollUp(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "S")
}

// ScrollDown scrolls the terminal content down by n lines, moving the cursor up
// by n lines. If n is zero, it does nothing. New lines are added at the top of
// the screen.
func ScrollDown(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "T")
}

// EnableMouse turns on mouse reporting (press, release, drag, and wheel)
// using SGR extended coordinates, which are needed for terminals wider or
// taller than 223 cells.
func EnableMouse() error {
	return write("\x1b[?1003h\x1b[?1006h")
}

// DisableMouse turns off mouse reporting enabled by EnableMouse.
func DisableMouse() error {
	return write("\x1b[?1006l\x1b[?1003l")
}

// OnMouse invokes fn for every mouse event decoded by the stdin reader
// started in EnterRawMode. Call the returned stop function to stop
// receiving events. EnableMouse must be called first, or the terminal won't
// send any events, and stdin must be in raw mode (see EnterRawMode) for
// events to be read and published at all.
func OnMouse(fn func(MouseEvent)) (stop func()) {
	id := bus.Subscribe(EventMouse, func(event Event) {
		if m, ok := event.(MouseEvent); ok {
			fn(m)
		}
	})
	return func() { bus.Unsubscribe(EventMouse, id) }
}

// IsTerminal reports whether the given files (or stdin, if none are given)
// are connected to a terminal.
func IsTerminal(f ...*os.File) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, file := range f {
		if !term.IsTerminal(int(file.Fd())) {
			return false
		}
	}
	if len(f) == 0 {
		return term.IsTerminal(int(stdinFile))
	}
	return true
}

// IsPipe reports whether the given files (or stdin, if none are given) are
// connected to a named pipe (FIFO). It returns false for regular files and
// terminals.
func IsPipe(f ...*os.File) bool {
	for _, file := range f {
		if !isPipe(file) {
			return false
		}
	}
	if len(f) == 0 {
		return isPipe(os.Stdin)
	}
	return true
}

// IsFile reports whether the given files (or stdin, if noare
// regular files. It returns false for directories, pipes, and terminals.
func IsFile(f ...*os.File) bool {
	for _, file := range f {
		if !isFile(file) {
			return false
		}
	}
	if len(f) == 0 {
		return isFile(os.Stdin)
	}
	return true
}

// IsDumb reports whether the terminal is a "dumb" terminal, which has very
// limited capabilities and does not support advanced features like cursor
// movement or color. It is determined by checking the TERM environment variable.
func IsDumb() bool {
	return os.Getenv("TERM") == "dumb"
}

// IsWSL reports whether the program is running inside Windows Subsystem for
// Linux (WSL). It checks for specific environment variables and the presence
// of "microsoft" in /proc/version.
func IsWSL() bool {
	if os.Getenv("WSL_INTEROP") != "" ||
		os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}

	v := strings.ToLower(string(data))
	return strings.Contains(v, "microsoft")
}

// IsSSH reports whether the program is running in an SSH session. It checks for
// the presence of SSH-related environment variables.
func IsSSH() bool {
	return os.Getenv("SSH_CONNECTION") != "" ||
		os.Getenv("SSH_CLIENT") != "" ||
		os.Getenv("SSH_TTY") != ""
}

// IsDocker reports whether the program is running inside a Docker container.
func IsDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}

	s := string(data)

	return strings.Contains(s, "docker") ||
		strings.Contains(s, "containerd") ||
		strings.Contains(s, "kubepods")
}

// IsRaw reports whether stdin is currently in raw mode, which means that input
// is being read directly from the terminal without line buffering or
// processing. It returns true if stdin is in raw mode, false otherwise.
func IsRaw() bool {
	mu.Lock()
	defer mu.Unlock()
	return rawModeState != nil
}

// SetColorLevel sets the color mode for the terminal. It can be used to force
// a specific color mode (None, Basic, 256, TrueColor) regardless of the
// terminal's capabilities or environment variables. This can be useful for
// testing or when the terminal detection is incorrect.
func SetColorLevel(level ColorMode) {
	mu.Lock()
	defer mu.Unlock()
	colorLevel = level
}

// GetColorLevel returns the current color mode of the terminal, which can be
// None, Basic, 256, or TrueColor. It reflects the terminal's capabilities and
// any overrides set by SetColorLevel.
func GetColorLevel() ColorMode {
	mu.Lock()
	defer mu.Unlock()
	return colorLevel
}

// SupportsColor reports whether the terminal supports any color output. It
// returns true if the color mode is Basic, 256, or TrueColor, and false if it
// is None. This can be used to conditionally enable or disable color output in
// applications. Use SetColorLevel to force a specific color mode if needed.
func SupportsColor() bool {
	mu.Lock()
	defer mu.Unlock()
	return colorLevel != ColorModeNone
}

// SupportsColorMode reports whether the terminal supports the specified color
// mode. It returns true if the current color mode is equal to or greater than
// the requested mode, and false otherwise. This can be used to check for
// specific color capabilities before using advanced color features.
func SupportsColorMode(mode ColorMode) bool {
	mu.Lock()
	defer mu.Unlock()
	return colorLevel >= mode
}

// OnResize polls the screen size and invokes fn whenever it changes. Call
// the returned stop function to stop polling.
func OnResize(fn func(e ResizeEvent)) (stop func()) {
	mu.Lock()
	defer mu.Unlock()
	id := bus.Subscribe(EventResize, func(event Event) {
		if e, ok := event.(ResizeEvent); ok {
			fn(e)
		}
	})

	stopResizeWatcher = startResizeWatcher()

	return func() {
		bus.Unsubscribe(EventResize, id)
		mu.Lock()
		defer mu.Unlock()
		if bus.CountListeners(EventResize) == 0 {
			if stopResizeWatcher != nil {
				stopResizeWatcher()
			}
			stopResizeWatcher = nil
		}
	}
}

// OnKey invokes fn for every key event decoded by the stdin reader started in
// EnterRawMode. Call the returned stop function to stop receiving events. Stdin
// must be in raw mode (see EnterRawMode) for events to be read and published at
// all.
func OnKey(fn func(e KeyEvent)) (stop func()) {
	id := bus.Subscribe(EventKey, func(event Event) {
		if e, ok := event.(KeyEvent); ok {
			fn(e)
		}
	})
	return func() { bus.Unsubscribe(EventKey, id) }
}

// OnEvent invokes fn for every event (key, mouse, resize, cursor position, etc.)
// decoded by the stdin reader started in EnterRawMode. Call the returned stop
// function to stop receiving events. Stdin must be in raw mode (see
// EnterRawMode) for events to be read and published at all.
//
// Some events may come from other sources, such as OS signals or external
// callbacks, so fn may receive events that are not directly related to stdin
// input.
func OnEvent(fn func(e Event)) (stop func()) {
	id := bus.Subscribe("*", func(event Event) {
		if e, ok := event.(KeyEvent); ok {
			fn(e)
		}
	})
	return func() { bus.Unsubscribe("*", id) }
}

func isPipe(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeNamedPipe != 0
}

func isFile(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}
