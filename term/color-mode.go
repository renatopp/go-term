package term

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var colorLevel ColorMode = ColorModeNone

type ColorMode uint8

const (
	ColorModeNone ColorMode = iota
	ColorModeAnsi
	ColorMode256
	ColorModeTrue
)

func init() {
	colorLevel = detectColorMode()
}

func SetColorLevel(level ColorMode) {
	colorLevel = level
}

func GetColorLevel() ColorMode {
	return colorLevel
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
