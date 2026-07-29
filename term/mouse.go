package term

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type MouseButton uint8

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButtonNone // reported during a drag with no button held
)

type MouseAction uint8

const (
	MouseActionPress MouseAction = iota
	MouseActionRelease
	MouseActionMotion
)

// MouseMsg reports a mouse event: a button press/release, a drag while a
// button is held, or a wheel scroll. X and Y are 0-based screen columns and
// rows.
type MouseMsg struct {
	X, Y   int
	Button MouseButton
	Action MouseAction
	Shift  bool
	Alt    bool
	Ctrl   bool
}

var mouseSource io.Reader = os.Stdin

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

// OnMouse reads SGR mouse escape sequences from stdin and invokes fn for
// each one. Call the returned stop function to stop reading. EnableMouse
// must be called first, or the terminal won't send any events, and stdin
// must be in raw mode (see EnterRawMode) for the raw escape sequences to
// reach the process.
//
// If the terminal never sends another byte, the underlying read blocks
// forever and the goroutine outlives stop(); the process exiting is what
// ultimately reclaims it, the same tradeoff already made for other blocking
// stdin reads in this package.
func OnMouse(fn func(MouseMsg)) (stop func()) {
	done := make(chan struct{})
	r := bufio.NewReader(mouseSource)

	go func() {
		for {
			msg, err := readMouseEvent(r)
			if err != nil {
				return
			}
			select {
			case <-done:
				return
			default:
				fn(msg)
			}
		}
	}()

	return func() { close(done) }
}

// readMouseEvent scans r for the next SGR mouse escape sequence
// ("\x1b[<Cb;Cx;CyM" or "...m"), discarding any other bytes it sees along
// the way (including keystrokes, which this package doesn't otherwise
// consume).
func readMouseEvent(r *bufio.Reader) (MouseMsg, error) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return MouseMsg{}, err
		}
		if b != 0x1b {
			continue
		}

		if b, err = r.ReadByte(); err != nil {
			return MouseMsg{}, err
		} else if b != '[' {
			continue
		}

		if b, err = r.ReadByte(); err != nil {
			return MouseMsg{}, err
		} else if b != '<' {
			continue
		}

		if msg, ok, err := parseSGRMouse(r); err != nil {
			return MouseMsg{}, err
		} else if ok {
			return msg, nil
		}
	}
}

func parseSGRMouse(r *bufio.Reader) (msg MouseMsg, ok bool, err error) {
	var params strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			return MouseMsg{}, false, err
		}
		if b == 'M' || b == 'm' {
			msg, ok := decodeSGRMouse(params.String(), b == 'M')
			return msg, ok, nil
		}
		params.WriteByte(b)
	}
}

// decodeSGRMouse parses the "Cb;Cx;Cy" parameter string of an SGR mouse
// sequence. pressed is true when the sequence was terminated by 'M' (press,
// drag, or wheel), false for 'm' (release).
func decodeSGRMouse(params string, pressed bool) (MouseMsg, bool) {
	parts := strings.SplitN(params, ";", 3)
	if len(parts) != 3 {
		return MouseMsg{}, false
	}

	cb, err1 := strconv.Atoi(parts[0])
	cx, err2 := strconv.Atoi(parts[1])
	cy, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return MouseMsg{}, false
	}

	msg := MouseMsg{
		X:     cx - 1,
		Y:     cy - 1,
		Shift: cb&4 != 0,
		Alt:   cb&8 != 0,
		Ctrl:  cb&16 != 0,
	}

	switch {
	case cb&64 != 0:
		msg.Action = MouseActionPress
		if cb&1 != 0 {
			msg.Button = MouseButtonWheelDown
		} else {
			msg.Button = MouseButtonWheelUp
		}
	case cb&32 != 0:
		msg.Action = MouseActionMotion
		msg.Button = mouseButtonFromCb(cb)
	case pressed:
		msg.Action = MouseActionPress
		msg.Button = mouseButtonFromCb(cb)
	default:
		msg.Action = MouseActionRelease
		msg.Button = mouseButtonFromCb(cb)
	}
	return msg, true
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
