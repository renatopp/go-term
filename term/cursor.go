package term

import (
	"fmt"
	"os"
	"strconv"
)

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
		// Request cursor position
		if err := write("\x1b[6n"); err != nil {
			row, col, oerr = 0, 0, err
			return
		}

		// Read response from terminal
		var buf [32]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			row, col, oerr = 0, 0, err
			return
		}

		// Parse response
		var r, c int
		if _, err := fmt.Sscanf(string(buf[:n]), "\x1b[%d;%dR", &r, &c); err != nil {
			row, col, oerr = 0, 0, err
			return
		}

		row, col = r, c
	})
	return row, col, nil
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
