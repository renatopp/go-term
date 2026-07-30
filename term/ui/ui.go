package ui

type Component interface {
	Render(width, height int) []string
}

// Width is context-free; height depends on the width it's given, since a
// component's content (e.g. wrapped text) can reflow to a different number
// of lines depending on the width it's rendered at.
type Sized interface {
	PreferredWidth() int
	PreferredHeight(width int) int
}

type Rect struct {
	X, Y, Width, Height int
}
