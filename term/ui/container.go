package ui

import "sort"

type Direction uint8

const (
	DirectionRow Direction = iota
	DirectionColumn
)

type Justify uint8

const (
	JustifyStart Justify = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
	JustifySpaceAround
)

type Align uint8

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
	AlignStretch
	AlignAuto // inherit the container's align
)

type Container struct {
	direction Direction
	justify   Justify
	align     Align
	gap       int
	padding   int
	items     []*Item
}

func NewContainer(direction Direction) *Container {
	return &Container{
		direction: direction,
		align:     AlignStretch,
	}
}

func Row(items ...*Item) *Container {
	return NewContainer(DirectionRow).WithItem(items...)
}

func Column(items ...*Item) *Container {
	return NewContainer(DirectionColumn).WithItem(items...)
}

// WithJustify sets the justification of the container's children along the main axis.
func (c *Container) WithJustify(j Justify) *Container {
	c.justify = j
	return c
}

// WithAlign sets the alignment of the container's children along the cross axis.
func (c *Container) WithAlign(a Align) *Container {
	c.align = a
	return c
}

// WithGap sets the gap between the container's children along the main axis.
func (c *Container) WithGap(n int) *Container {
	c.gap = n
	return c
}

// WithPadding sets the uniform space between the container's edges and its
// children, applied before children are laid out.
func (c *Container) WithPadding(n int) *Container {
	c.padding = n
	return c
}

// WithItem adds the given items as children of the container.
func (c *Container) WithItem(items ...*Item) *Container {
	c.items = append(c.items, items...)
	return c
}

// WithComponent adds the given components as children of the container. Each
// component is wrapped in an Item with default grow, shrink, and basis values.
func (c *Container) WithComponent(components ...Component) *Container {
	for _, comp := range components {
		c.WithItem(NewItem(comp))
	}
	return c
}

func (c *Container) WithSpacer() *Container {
	return c.WithItem(NewItem(NewSpacer()).WithGrow(1))
}

// RemoveItem removes the given items from the container's children, if
// present. Items not added via WithItem (e.g. those wrapped internally by
// WithComponent) cannot be matched and are ignored.
func (c *Container) RemoveItem(items ...*Item) *Container {
	remaining := c.items[:0]
outer:
	for _, existing := range c.items {
		for _, target := range items {
			if existing == target {
				continue outer
			}
		}
		remaining = append(remaining, existing)
	}
	c.items = remaining
	return c
}

// Clear removes all of the container's children.
func (c *Container) Clear() *Container {
	c.items = nil
	return c
}

// Layout computes the layout of the container's children given the available
// width and height. It returns a slice of Rects representing the position
// and size of each child.
func (c *Container) Layout(width, height int) []Rect {
	items := c.visibleItems()
	n := len(items)
	rects := make([]Rect, n)
	if n == 0 {
		return rects
	}

	width = max(0, width-2*c.padding)
	height = max(0, height-2*c.padding)

	grows := make([]int, n)
	shrinks := make([]int, n)
	for i, item := range items {
		grows[i] = item.grow
		shrinks[i] = item.shrink
	}

	if c.direction == DirectionRow {
		bases := make([]int, n)
		mins := make([]int, n)
		maxs := make([]int, n)
		for i, item := range items {
			bases[i] = item.mainWidth()
			mins[i] = item.minWidth
			maxs[i] = item.maxWidth
		}
		widths, xs := resolveAxis(bases, grows, shrinks, mins, maxs, width, c.gap, c.justify)

		for i, item := range items {
			align := item.resolveAlign(c.align)
			h := height
			if align != AlignStretch {
				h = item.crossHeight(widths[i])
			}
			h = item.clampHeight(h)
			rects[i] = Rect{X: xs[i], Y: crossOffset(align, height, h), Width: widths[i], Height: h}
		}
		offsetRects(rects, c.padding, c.padding)
		return rects
	}

	widths := make([]int, n)
	for i, item := range items {
		if item.resolveAlign(c.align) == AlignStretch {
			widths[i] = width
		} else {
			widths[i] = item.crossWidth()
		}
		widths[i] = item.clampWidth(widths[i])
	}

	bases := make([]int, n)
	mins := make([]int, n)
	maxs := make([]int, n)
	for i, item := range items {
		bases[i] = item.mainHeight(widths[i])
		mins[i] = item.minHeight
		maxs[i] = item.maxHeight
	}
	heights, ys := resolveAxis(bases, grows, shrinks, mins, maxs, height, c.gap, c.justify)

	for i, item := range items {
		rects[i] = Rect{X: crossOffset(item.resolveAlign(c.align), width, widths[i]), Y: ys[i], Width: widths[i], Height: heights[i]}
	}
	offsetRects(rects, c.padding, c.padding)
	return rects
}

// Render renders the container and its children into a slice of strings,
// each representing a line of text. The width and height parameters specify
// the available space for rendering.
func (c *Container) Render(width, height int) []string {
	rects := c.Layout(width, height)
	grid := make([][]rune, height)
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	for i, item := range c.visibleItems() {
		r := rects[i]
		lines := item.component.Render(r.Width, r.Height)
		for dy, line := range lines {
			y := r.Y + dy
			if dy >= r.Height || y < 0 || y >= height {
				continue
			}
			x := r.X
			for _, ch := range line {
				if x >= r.X+r.Width || x >= width {
					break
				}
				if x >= 0 {
					grid[y][x] = ch
				}
				x++
			}
		}
	}

	out := make([]string, height)
	for y, row := range grid {
		out[y] = string(row)
	}
	return out
}

// PreferredWidth returns the preferred width of the container, which is the
// sum of the preferred widths of its children plus the gaps between them.
func (c *Container) PreferredWidth() int {
	items := c.visibleItems()
	n := len(items)
	if n == 0 {
		return 0
	}

	if c.direction == DirectionRow {
		sum := c.gap * (n - 1)
		for _, item := range items {
			sum += item.mainWidth()
		}
		return sum + 2*c.padding
	}

	w := 0
	for _, item := range items {
		w = max(w, item.crossWidth())
	}
	return w + 2*c.padding
}

// PreferredHeight returns the preferred height of the container given a
// specific width.
func (c *Container) PreferredHeight(width int) int {
	items := c.visibleItems()
	n := len(items)
	if n == 0 {
		return 0
	}

	width = max(0, width-2*c.padding)

	if c.direction == DirectionColumn {
		sum := c.gap * (n - 1)
		for _, item := range items {
			w := width
			if item.resolveAlign(c.align) != AlignStretch {
				w = item.crossWidth()
			}
			sum += item.mainHeight(w)
		}
		return sum + 2*c.padding
	}

	bases := make([]int, n)
	grows := make([]int, n)
	shrinks := make([]int, n)
	mins := make([]int, n)
	maxs := make([]int, n)
	for i, item := range items {
		bases[i] = item.mainWidth()
		grows[i] = item.grow
		shrinks[i] = item.shrink
		mins[i] = item.minWidth
		maxs[i] = item.maxWidth
	}
	widths, _ := resolveAxis(bases, grows, shrinks, mins, maxs, width, c.gap, c.justify)

	h := 0
	for i, item := range items {
		h = max(h, item.crossHeight(widths[i]))
	}
	return h + 2*c.padding
}

func (c *Container) visibleItems() []*Item {
	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		if !item.hidden {
			items = append(items, item)
		}
	}
	return items
}

func resolveAxis(bases, grows, shrinks, mins, maxs []int, container, gap int, justify Justify) (sizes, offsets []int) {
	n := len(bases)
	sizes = make([]int, n)
	copy(sizes, bases)
	offsets = make([]int, n)
	if n == 0 {
		return
	}

	totalGap := gap * (n - 1)
	used := totalGap
	for _, s := range sizes {
		used += s
	}
	free := container - used

	switch {
	case free > 0:
		distributeGrow(sizes, grows, free)
	case free < 0:
		distributeShrink(sizes, shrinks, bases, -free)
	}

	for i := range sizes {
		sizes[i] = clampSize(sizes[i], mins[i], maxs[i])
	}

	leftover := container - totalGap
	for _, s := range sizes {
		leftover -= s
	}
	if leftover < 0 {
		leftover = 0
	}

	totalGrow := 0
	for _, g := range grows {
		totalGrow += g
	}

	var start, between int
	if totalGrow == 0 {
		switch justify {
		case JustifyCenter:
			start = leftover / 2
		case JustifyEnd:
			start = leftover
		case JustifySpaceBetween:
			if n > 1 {
				between = leftover / (n - 1)
			}
		case JustifySpaceAround:
			between = leftover / n
			start = between / 2
		}
	}

	pos := start
	for i, s := range sizes {
		offsets[i] = pos
		pos += s + gap + between
	}

	return sizes, offsets
}

func distributeGrow(sizes, weights []int, extra int) {
	total := 0
	for _, w := range weights {
		total += w
	}
	if total <= 0 {
		return
	}

	type share struct {
		idx int
		rem int
	}

	remaining := extra
	shares := make([]share, 0, len(weights))
	for i, w := range weights {
		if w <= 0 {
			continue
		}
		add := extra * w / total
		sizes[i] += add
		remaining -= add
		shares = append(shares, share{idx: i, rem: extra * w % total})
	}

	sort.Slice(shares, func(a, b int) bool { return shares[a].rem > shares[b].rem })
	for k := 0; k < remaining && k < len(shares); k++ {
		sizes[shares[k].idx]++
	}
}

func distributeShrink(sizes, shrinks, bases []int, deficit int) {
	weights := make([]int, len(shrinks))
	total := 0
	for i, s := range shrinks {
		weights[i] = s * bases[i]
		total += weights[i]
	}
	if total <= 0 {
		return
	}

	type share struct {
		idx int
		rem int
	}

	remaining := deficit
	shares := make([]share, 0, len(weights))
	for i, w := range weights {
		if w <= 0 {
			continue
		}
		sub := deficit * w / total
		if sub > sizes[i] {
			sub = sizes[i]
		}
		sizes[i] -= sub
		remaining -= sub
		shares = append(shares, share{idx: i, rem: deficit * w % total})
	}

	sort.Slice(shares, func(a, b int) bool { return shares[a].rem > shares[b].rem })
	extra := remaining
	for k := 0; k < extra && k < len(shares); k++ {
		if sizes[shares[k].idx] > 0 {
			sizes[shares[k].idx]--
		}
	}
}

func crossOffset(align Align, container, item int) int {
	switch align {
	case AlignCenter:
		return (container - item) / 2
	case AlignEnd:
		return container - item
	default:
		return 0
	}
}

func offsetRects(rects []Rect, dx, dy int) {
	for i := range rects {
		rects[i].X += dx
		rects[i].Y += dy
	}
}
