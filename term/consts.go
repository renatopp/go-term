package term

// KeyType identifies a key reported by KeyEvent. KeyRune covers ordinary
// character keys (including Ctrl/Alt/Shift combinations, reported via the
// Ctrl, Alt, and Shift fields); the rest name keys with no printable rune of
// their own.
type KeyType int

const (
	KeyRune KeyType = iota
	KeyEnter
	KeyEsc
	KeyTab
	KeyBackspace
	KeyDelete
	KeyInsert
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
)

// keyboardProtocol identifies which keyboard-modifier reporting protocol is
// active on the terminal, negotiated by EnterRawMode.
type keyboardProtocol uint8

const (
	keyboardProtocolNone keyboardProtocol = iota
	keyboardProtocolKitty
	keyboardProtocolModifyOtherKeys
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

type ColorMode uint8

const (
	ColorModeNone ColorMode = iota
	ColorModeAnsi
	ColorMode256
	ColorModeTrue
)

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

var (
	Success = NewStyle().WithBackground(ColorSuccess)
	Info    = NewStyle().WithBackground(ColorInfo)
	Warning = NewStyle().WithBackground(ColorWarning)
	Error   = NewStyle().WithBackground(ColorError)
)
