package ui

type item struct {
	component    Renderable
	grow         int
	shrink       int
	basis        int
	basisPercent int
	minWidth     int
	maxWidth     int
	minHeight    int
	maxHeight    int
	alignSelf    Align
	hidden       bool
}

func Item(c Renderable) *item {
	return &item{
		component:    c,
		shrink:       1,
		basis:        -1,
		basisPercent: -1,
		maxWidth:     -1,
		maxHeight:    -1,
		alignSelf:    AlignAuto,
	}
}

// spacer is a blank component that renders nothing.
type spacer struct{}

func (spacer) Render(width, height int) []string {
	return nil
}

// Spacer returns an empty item that grows to fill the available space along
// the container's main axis, pushing the surrounding items apart.
func Spacer() *item {
	return Item(spacer{}).WithGrow(1)
}

func (i *item) WithGrow(weight int) *item {
	i.grow = weight
	return i
}

func (i *item) WithShrink(weight int) *item {
	i.shrink = weight
	return i
}

func (i *item) WithBasis(n int) *item {
	i.basis = n
	return i
}

// WithBasisPercent sets the item's basis as a percentage of the container's
// main axis size. It takes precedence over WithBasis when the container size
// is known; otherwise the item falls back to its basis or preferred size.
func (i *item) WithBasisPercent(pct int) *item {
	i.basisPercent = pct
	return i
}

// WithMinWidth sets the smallest width the item may be shrunk to.
func (i *item) WithMinWidth(n int) *item {
	i.minWidth = n
	return i
}

// WithMaxWidth sets the largest width the item may be grown to.
func (i *item) WithMaxWidth(n int) *item {
	i.maxWidth = n
	return i
}

// WithMinHeight sets the smallest height the item may be shrunk to.
func (i *item) WithMinHeight(n int) *item {
	i.minHeight = n
	return i
}

// WithMaxHeight sets the largest height the item may be grown to.
func (i *item) WithMaxHeight(n int) *item {
	i.maxHeight = n
	return i
}

// WithMinSize sets the smallest width and height the item may be shrunk to.
func (i *item) WithMinSize(w, h int) *item {
	i.minWidth = w
	i.minHeight = h
	return i
}

// WithMaxSize sets the largest width and height the item may be grown to.
func (i *item) WithMaxSize(w, h int) *item {
	i.maxWidth = w
	i.maxHeight = h
	return i
}

// WithAlignSelf overrides the container's align for this item along the
// cross axis. AlignAuto (the default) inherits the container's align.
func (i *item) WithAlignSelf(a Align) *item {
	i.alignSelf = a
	return i
}

// AsHidden excludes the item from layout and rendering entirely, as if it
// were never added to the container.
func (i *item) AsHidden(hidden bool) *item {
	i.hidden = hidden
	return i
}

// mainWidth returns the item's base width along a row's main axis. A negative
// available width means the container size is unknown, disabling percent basis.
func (i *item) mainWidth(available int) int {
	if i.basisPercent >= 0 && available >= 0 {
		return i.clampWidth(available * i.basisPercent / 100)
	}
	if i.basis >= 0 {
		return i.clampWidth(i.basis)
	}
	return i.clampWidth(i.preferredWidth())
}

// mainHeight returns the item's base height along a column's main axis. A
// negative available height means the container size is unknown, disabling
// percent basis.
func (i *item) mainHeight(width, available int) int {
	if i.basisPercent >= 0 && available >= 0 {
		return i.clampHeight(available * i.basisPercent / 100)
	}
	if i.basis >= 0 {
		return i.clampHeight(i.basis)
	}
	return i.clampHeight(i.preferredHeight(width))
}

func (i *item) crossWidth() int {
	return i.clampWidth(i.preferredWidth())
}

func (i *item) crossHeight(width int) int {
	return i.clampHeight(i.preferredHeight(width))
}

func (i *item) clampWidth(w int) int {
	return clampSize(w, i.minWidth, i.maxWidth)
}

func (i *item) clampHeight(h int) int {
	return clampSize(h, i.minHeight, i.maxHeight)
}

func (i *item) resolveAlign(containerAlign Align) Align {
	if i.alignSelf == AlignAuto {
		return containerAlign
	}
	return i.alignSelf
}

func (i *item) preferredWidth() int {
	if s, ok := i.component.(Sizeable); ok {
		return s.PreferredWidth()
	}
	return 0
}

func (i *item) preferredHeight(width int) int {
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
