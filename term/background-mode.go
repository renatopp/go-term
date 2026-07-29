package term

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// var backgroundMode BackgroundMode = BackgroundModeDark

// type BackgroundMode uint8

// const (
// 	BackgroundModeDark BackgroundMode = iota
// 	BackgroundModeLight
// )

// func init() {
// 	backgroundMode = detectBackgroundMode()
// }

// func SetBackgroundMode(mode BackgroundMode) {
// 	backgroundMode = mode
// }

// func GetBackgroundMode() BackgroundMode {
// 	return backgroundMode
// }

// // QueryBackgroundColor asks the terminal for its actual background color via
// // OSC 11. Unlike detectBackgroundMode, this requires the terminal to support
// // and respond to the query, so it can time out on terminals that don't.
// func QueryBackgroundColor() (color ColorTrue, oerr error) {
// 	err := WithinRawMode(func() {
// 		if err := write("\x1b]11;?\x07"); err != nil {
// 			oerr = err
// 			return
// 		}

// 		type readResult struct {
// 			n   int
// 			buf [64]byte
// 			err error
// 		}
// 		ch := make(chan readResult, 1)
// 		// If the terminal never replies, this goroutine leaks for the life
// 		// of the process; there's no portable way to cancel a blocking read
// 		// on stdin, so we just stop waiting on it via the timeout below.
// 		go func() {
// 			var res readResult
// 			res.n, res.err = os.Stdin.Read(res.buf[:])
// 			ch <- res
// 		}()

// 		select {
// 		case res := <-ch:
// 			if res.err != nil {
// 				oerr = res.err
// 				return
// 			}
// 			r, g, b, err := parseOSC11Response(string(res.buf[:res.n]))
// 			if err != nil {
// 				oerr = err
// 				return
// 			}
// 			color = NewColor(r, g, b)
// 		case <-time.After(200 * time.Millisecond):
// 			oerr = fmt.Errorf("timeout waiting for OSC 11 response")
// 		}
// 	})
// 	if err != nil {
// 		return color, err
// 	}
// 	return color, oerr
// }

// func ForceQueryBackgroundColor() ColorTrue {
// 	color, _ := QueryBackgroundColor()
// 	return color
// }

// // RefreshBackgroundMode queries the terminal's real background color via
// // OSC 11 and updates the mode from its luminance. It leaves the current
// // mode untouched if the terminal doesn't respond in time.
// func RefreshBackgroundMode() error {
// 	color, err := QueryBackgroundColor()
// 	if err != nil {
// 		return err
// 	}

// 	if luminance(color.r, color.g, color.b) > 127.5 {
// 		backgroundMode = BackgroundModeLight
// 	} else {
// 		backgroundMode = BackgroundModeDark
// 	}
// 	return nil
// }

// func parseOSC11Response(resp string) (r, g, b uint8, err error) {
// 	idx := strings.Index(resp, "rgb:")
// 	if idx == -1 {
// 		return 0, 0, 0, fmt.Errorf("unexpected OSC 11 response: %q", resp)
// 	}

// 	var r16, g16, b16 uint32
// 	if _, err := fmt.Sscanf(resp[idx+len("rgb:"):], "%4x/%4x/%4x", &r16, &g16, &b16); err != nil {
// 		return 0, 0, 0, err
// 	}
// 	return uint8(r16 >> 8), uint8(g16 >> 8), uint8(b16 >> 8), nil
// }

// func luminance(r, g, b uint8) float64 {
// 	return 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
// }

// func detectBackgroundMode() BackgroundMode {
// 	// COLORFGBG is set by some terminals/multiplexers (rxvt, tmux) as
// 	// "fg;bg" or "fg;default;bg", using the standard 16-color ANSI palette.
// 	v, ok := os.LookupEnv("COLORFGBG")
// 	if !ok {
// 		return BackgroundModeDark
// 	}

// 	parts := strings.Split(v, ";")
// 	bg, err := strconv.Atoi(parts[len(parts)-1])
// 	if err != nil {
// 		return BackgroundModeDark
// 	}

// 	// 7 (white) and 15 (bright white) are the only palette entries
// 	// terminals commonly use as a light background.
// 	if bg == 7 || bg == 15 {
// 		return BackgroundModeLight
// 	}
// 	return BackgroundModeDark
// }
