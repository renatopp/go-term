package ui

type Event any

type Component interface {
	Render(width, height int) []string
}

type Updatable interface {
	Update(Event)
}

type Sizeable interface {
	PreferredWidth() int
	PreferredHeight(width int) int
}

type Rect struct {
	X, Y, Width, Height int
}
