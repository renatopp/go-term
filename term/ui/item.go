package ui

type ContainerItem struct {
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

func Item(c Renderable) *ContainerItem {
	return &ContainerItem{
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
func Spacer() *ContainerItem {
	return Item(spacer{}).WithGrow(1)
}

func (i *ContainerItem) WithGrow(weight int) *ContainerItem {
	i.grow = weight
	return i
}

func (i *ContainerItem) WithShrink(weight int) *ContainerItem {
	i.shrink = weight
	return i
}

func (i *ContainerItem) WithBasis(n int) *ContainerItem {
	i.basis = n
	return i
}

// WithBasisPercent sets the item's basis as a percentage of the container's
// main axis size. It takes precedence over WithBasis when the container size
// is known; otherwise the item falls back to its basis or preferred size.
func (i *ContainerItem) WithBasisPercent(pct int) *ContainerItem {
	i.basisPercent = pct
	return i
}

// WithMinWidth sets the smallest width the item may be shrunk to.
func (i *ContainerItem) WithMinWidth(n int) *ContainerItem {
	i.minWidth = n
	return i
}

// WithMaxWidth sets the largest width the item may be grown to.
func (i *ContainerItem) WithMaxWidth(n int) *ContainerItem {
	i.maxWidth = n
	return i
}

// WithMinHeight sets the smallest height the item may be shrunk to.
func (i *ContainerItem) WithMinHeight(n int) *ContainerItem {
	i.minHeight = n
	return i
}

// WithMaxHeight sets the largest height the item may be grown to.
func (i *ContainerItem) WithMaxHeight(n int) *ContainerItem {
	i.maxHeight = n
	return i
}

// WithMinSize sets the smallest width and height the item may be shrunk to.
func (i *ContainerItem) WithMinSize(w, h int) *ContainerItem {
	i.minWidth = w
	i.minHeight = h
	return i
}

// WithMaxSize sets the largest width and height the item may be grown to.
func (i *ContainerItem) WithMaxSize(w, h int) *ContainerItem {
	i.maxWidth = w
	i.maxHeight = h
	return i
}

// WithAlignSelf overrides the container's align for this item along the
// cross axis. AlignAuto (the default) inherits the container's align.
func (i *ContainerItem) WithAlignSelf(a Align) *ContainerItem {
	i.alignSelf = a
	return i
}

// AsHidden excludes the item from layout and rendering entirely, as if it
// were never added to the container.
func (i *ContainerItem) AsHidden(hidden bool) *ContainerItem {
	i.hidden = hidden
	return i
}

// mainWidth returns the item's base width along a row's main axis. A negative
// available width means the container size is unknown, disabling percent basis.
func (i *ContainerItem) mainWidth(available int) int {
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
func (i *ContainerItem) mainHeight(width, available int) int {
	if i.basisPercent >= 0 && available >= 0 {
		return i.clampHeight(available * i.basisPercent / 100)
	}
	if i.basis >= 0 {
		return i.clampHeight(i.basis)
	}
	return i.clampHeight(i.preferredHeight(width))
}

func (i *ContainerItem) crossWidth() int {
	return i.clampWidth(i.preferredWidth())
}

func (i *ContainerItem) crossHeight(width int) int {
	return i.clampHeight(i.preferredHeight(width))
}

func (i *ContainerItem) clampWidth(w int) int {
	return clampSize(w, i.minWidth, i.maxWidth)
}

func (i *ContainerItem) clampHeight(h int) int {
	return clampSize(h, i.minHeight, i.maxHeight)
}

func (i *ContainerItem) resolveAlign(containerAlign Align) Align {
	if i.alignSelf == AlignAuto {
		return containerAlign
	}
	return i.alignSelf
}

func (i *ContainerItem) preferredWidth() int {
	if s, ok := i.component.(Sizeable); ok {
		return s.PreferredWidth()
	}
	return 0
}

func (i *ContainerItem) preferredHeight(width int) int {
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
