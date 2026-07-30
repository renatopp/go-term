package term

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

func SetWriter(w io.Writer) {
	writer = w
}

func SetStdin(fd uintptr) {
	stdin = fd
}

func SetStdout(fd uintptr) {
	stdout = fd
}

// EnterRawMode puts stdin into raw mode and starts reading it for mouse,
// keyboard, and other events (see startStdinReader), which are published on
// the package bus. Call ExitRawMode to restore the terminal and stop
// reading.
func EnterRawMode() error {
	if rawModeState != nil {
		return nil
	}
	state, err := term.MakeRaw(int(stdin))
	if err != nil {
		return err
	}
	rawModeState = state
	stopStdinReader = startStdinReader()
	return nil
}

func ExitRawMode() error {
	state := rawModeState
	rawModeState = nil
	if stopStdinReader != nil {
		stopStdinReader()
		stopStdinReader = nil
	}
	return RestoreTerminalState(state)
}

func GetTerminalState() (*State, error) {
	return term.GetState(int(stdin))
}

func RestoreTerminalState(state *State) error {
	return term.Restore(int(stdin), state)
}

func SetWindowTitle(title string) error {
	return write("\x1b]0;" + title + "\x07")
}

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

func EnterAlternateScreen() error {
	return write("\x1b[?1049h")
}

func ExitAlternateScreen() error {
	return write("\x1b[?1049l")
}

func GetScreenSize() (width, height int, err error) {
	return term.GetSize(int(stdout))
}

func ForceGetScreenSize() (width, height int) {
	width, height, _ = GetScreenSize()
	return
}

func ClearScreen() error {
	return write("\x1b[2J\x1b[H")
}

func ClearLine() error {
	return write("\x1b[2K")
}

func ClearLineAfterCursor() error {
	return write("\x1b[0K")
}

func ClearLineBeforeCursor() error {
	return write("\x1b[1K")
}

func MoveCursorUp(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "A")
}

func MoveCursorDown(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "B")
}

func MoveCursorForward(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "C")
}

func MoveCursorBackward(n int) error {
	return write("\x1b[" + strconv.Itoa(n) + "D")
}

func MoveCursorTo(row, col int) error {
	return write("\x1b[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H")
}

func MoveCursorToRow(row int) error {
	return write("\x1b[" + strconv.Itoa(row) + "d")
}

func MoveCursorToColumn(col int) error {
	return write("\x1b[" + strconv.Itoa(col) + "G")
}

func MoveCursorToHome() error {
	return write("\x1b[H")
}

func MoveCursorToStartOfLine() error {
	return write("\x1b[1G")
}

func MoveCursorToEndOfLine() error {
	return write("\x1b[999C")
}

func MoveCursorToBottom() error {
	return write("\x1b[999B")
}

func SaveCursorPosition() error {
	return write("\x1b[s")
}

func RestoreCursorPosition() error {
	return write("\x1b[u")
}

func HideCursor() error {
	return write("\x1b[?25l")
}

func ShowCursor() error {
	return write("\x1b[?25h")
}

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

func ForceGetCursorPosition() (row, col int) {
	row, col, _ = GetCursorPosition()
	return
}

func EnableCursorWrap() error {
	return write("\x1b[?7h")
}

func DisableCursorWrap() error {
	return write("\x1b[?7l")
}

func EnableCursorOriginMode() error {
	return write("\x1b[?6h")
}

func DisableCursorOriginMode() error {
	return write("\x1b[?6l")
}

// EnableMouse turns on mouse reporting (press, release, drag, and wheel)
// using SGR extended coordinates, which are needed for terminals wider or
// taller than 223 cells.
func EnableMouse() error {
	return write("\x1b[?1002h\x1b[?1006h")
}

// DisableMouse turns off mouse reporting enabled by EnableMouse.
func DisableMouse() error {
	return write("\x1b[?1006l\x1b[?1002l")
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
	for _, file := range f {
		if !term.IsTerminal(int(file.Fd())) {
			return false
		}
	}
	if len(f) == 0 {
		return term.IsTerminal(int(stdin))
	}
	return true
}

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

func IsDumb() bool {
	return os.Getenv("TERM") == "dumb"
}

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

func IsSSH() bool {
	return os.Getenv("SSH_CONNECTION") != "" ||
		os.Getenv("SSH_CLIENT") != "" ||
		os.Getenv("SSH_TTY") != ""
}

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

func IsRaw() bool {
	return rawModeState != nil
}

func SetColorLevel(level ColorMode) {
	colorLevel = level
}

func GetColorLevel() ColorMode {
	return colorLevel
}

func SupportsColor() bool {
	return colorLevel != ColorModeNone
}

func SupportsColorMode(mode ColorMode) bool {
	return colorLevel >= mode
}

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

// OnResize polls the screen size and invokes fn whenever it changes. Call
// the returned stop function to stop polling.
func OnResize(fn func(width, height int)) (stop func()) {
	done := make(chan struct{})
	lastWidth, lastHeight, _ := screenSizeGetter()

	go func() {
		ticker := time.NewTicker(resizePollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				width, height, err := screenSizeGetter()
				if err != nil || (width == lastWidth && height == lastHeight) {
					continue
				}
				lastWidth, lastHeight = width, height
				fn(width, height)
			}
		}
	}()

	return func() {
		close(done)
	}
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
