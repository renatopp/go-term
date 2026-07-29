package term

import (
	"time"

	"golang.org/x/term"
)

var resizePollInterval = 100 * time.Millisecond
var screenSizeGetter = GetScreenSize

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
