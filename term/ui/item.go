package ui

type Item struct {
	component Component
	grow      int
	shrink    int
	basis     int
	minWidth  int
	maxWidth  int
	minHeight int
	maxHeight int
	alignSelf Align
	hidden    bool
}

func NewItem(c Component) *Item {
	return &Item{
		component: c,
		shrink:    1,
		basis:     -1,
		maxWidth:  -1,
		maxHeight: -1,
		alignSelf: AlignAuto,
	}
}

func (i *Item) WithGrow(weight int) *Item {
	i.grow = weight
	return i
}

func (i *Item) WithShrink(weight int) *Item {
	i.shrink = weight
	return i
}

func (i *Item) WithBasis(n int) *Item {
	i.basis = n
	return i
}

// WithMinWidth sets the smallest width the item may be shrunk to.
func (i *Item) WithMinWidth(n int) *Item {
	i.minWidth = n
	return i
}

// WithMaxWidth sets the largest width the item may be grown to.
func (i *Item) WithMaxWidth(n int) *Item {
	i.maxWidth = n
	return i
}

// WithMinHeight sets the smallest height the item may be shrunk to.
func (i *Item) WithMinHeight(n int) *Item {
	i.minHeight = n
	return i
}

// WithMaxHeight sets the largest height the item may be grown to.
func (i *Item) WithMaxHeight(n int) *Item {
	i.maxHeight = n
	return i
}

// WithAlignSelf overrides the container's align for this item along the
// cross axis. AlignAuto (the default) inherits the container's align.
func (i *Item) WithAlignSelf(a Align) *Item {
	i.alignSelf = a
	return i
}

// AsHidden excludes the item from layout and rendering entirely, as if it
// were never added to the container.
func (i *Item) AsHidden(hidden bool) *Item {
	i.hidden = hidden
	return i
}

func (i *Item) mainWidth() int {
	if i.basis >= 0 {
		return i.clampWidth(i.basis)
	}
	return i.clampWidth(i.preferredWidth())
}

func (i *Item) mainHeight(width int) int {
	if i.basis >= 0 {
		return i.clampHeight(i.basis)
	}
	return i.clampHeight(i.preferredHeight(width))
}

func (i *Item) crossWidth() int {
	return i.clampWidth(i.preferredWidth())
}

func (i *Item) crossHeight(width int) int {
	return i.clampHeight(i.preferredHeight(width))
}

func (i *Item) clampWidth(w int) int {
	return clampSize(w, i.minWidth, i.maxWidth)
}

func (i *Item) clampHeight(h int) int {
	return clampSize(h, i.minHeight, i.maxHeight)
}

func (i *Item) resolveAlign(containerAlign Align) Align {
	if i.alignSelf == AlignAuto {
		return containerAlign
	}
	return i.alignSelf
}

func (i *Item) preferredWidth() int {
	if s, ok := i.component.(Sizeable); ok {
		return s.PreferredWidth()
	}
	return 0
}

func (i *Item) preferredHeight(width int) int {
	if s, ok := i.component.(Sizeable); ok {
		return s.PreferredHeight(width)
	}
	return 0
}

// clampSize bounds v to [min, max]. A negative max means unbounded.
func clampSize(v, min, max int) int {
	if v < min {
		v = min
	}
	if max >= 0 && v > max {
		v = max
	}
	return v
}
