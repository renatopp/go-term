package term

import "fmt"

type Color interface {
	Mode() ColorMode
	Fallback() Color
	Sequence(bg bool) string
}

type ColorAnsi struct{ code uint8 }

func NewColorAnsi(code uint8) ColorAnsi {
	return ColorAnsi{code: code}
}
func (c ColorAnsi) Mode() ColorMode {
	return ColorModeAnsi
}
func (c ColorAnsi) Fallback() Color {
	return c
}
func (c ColorAnsi) Sequence(bg bool) string {
	return fmt.Sprintf("%d", c.code)
}

type Color256 struct{ code uint8 }

func NewColor256(code uint8) Color256 {
	return Color256{code: code}
}
func (c Color256) Mode() ColorMode {
	return ColorMode256
}
func (c Color256) Fallback() Color {
	return c
}
func (c Color256) Sequence(bg bool) string {
	if bg {
		return fmt.Sprintf("48;5;%d", c.code)
	}
	return fmt.Sprintf("38;5;%d", c.code)
}

type ColorTrue struct{ r, g, b uint8 }

func NewColor(r, g, b uint8) ColorTrue {
	return ColorTrue{r: r, g: g, b: b}
}
func NewColorHex(hex string) ColorTrue {
	var r, g, b uint8
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) == 6 {
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	} else if len(hex) == 3 {
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r = r * 17
		g = g * 17
		b = b * 17
	}
	return ColorTrue{r: r, g: g, b: b}
}
func (c ColorTrue) Mode() ColorMode {
	return ColorModeTrue
}
func (c ColorTrue) Fallback() Color {
	return Color256{code: rgbTo256(c.r, c.g, c.b)}
}
func (c ColorTrue) Sequence(bg bool) string {
	if bg {
		return fmt.Sprintf("48;2;%d;%d;%d", c.r, c.g, c.b)
	}
	return fmt.Sprintf("38;2;%d;%d;%d", c.r, c.g, c.b)
}

func rgbTo256(r, g, b uint8) uint8 {
	// Convert RGB to 256-color code.
	if r == g && g == b {
		if r < 8 {
			return 16
		}
		if r > 248 {
			return 231
		}
		return uint8((int(r)-8)*24/247 + 232)
	}
	r6 := int(r) * 5 / 255
	g6 := int(g) * 5 / 255
	b6 := int(b) * 5 / 255
	return uint8(16 + (36 * r6) + (6 * g6) + b6)
}

// ------

var (
	ColorSuccess = NewColor(34, 197, 94)
	ColorInfo    = NewColor(59, 130, 246)
	ColorWarning = NewColor(245, 158, 11)
	ColorError   = NewColor(239, 68, 68)

	ColorRed       = NewColor(239, 68, 68)
	ColorOrange    = NewColor(249, 115, 22)
	ColorAmber     = NewColor(245, 158, 11)
	ColorYellow    = NewColor(234, 179, 8)
	ColorLime      = NewColor(132, 204, 22)
	ColorGreen     = NewColor(34, 197, 94)
	ColorEmerald   = NewColor(16, 185, 129)
	ColorTeal      = NewColor(20, 184, 166)
	ColorCyan      = NewColor(6, 182, 212)
	ColorSky       = NewColor(14, 165, 233)
	ColorBlue      = NewColor(59, 130, 246)
	ColorIndigo    = NewColor(99, 102, 241)
	ColorViolet    = NewColor(139, 92, 246)
	ColorPurple    = NewColor(168, 85, 247)
	ColorFuchsia   = NewColor(217, 70, 239)
	ColorPink      = NewColor(236, 72, 153)
	ColorRose      = NewColor(244, 63, 94)
	ColorWhite     = NewColor(255, 255, 255)
	ColorBlack     = NewColor(0, 0, 0)
	ColorGray      = NewColor(107, 114, 128)
	ColorLightGray = NewColor(209, 213, 219)
	ColorDarkGray  = NewColor(55, 65, 8)
)
