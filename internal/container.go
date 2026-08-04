package term

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

type container struct {
	direction Direction
	justify   Justify
	align     Align
	gap       int
	paddingX  int
	paddingY  int
	items     []*ContainerItem
}

func Container(direction Direction, items ...*ContainerItem) *container {
	return (&container{
		direction: direction,
		align:     AlignStretch,
	}).WithItem(items...)
}

func Row(items ...*ContainerItem) *container {
	return Container(DirectionRow).WithItem(items...)
}

func Column(items ...*ContainerItem) *container {
	return Container(DirectionColumn).WithItem(items...)
}

// WithJustify sets the justification of the container's children along the main axis.
func (c *container) WithJustify(j Justify) *container {
	c.justify = j
	return c
}

// WithAlign sets the alignment of the container's children along the cross axis.
func (c *container) WithAlign(a Align) *container {
	c.align = a
	return c
}

// WithGap sets the gap between the container's children along the main axis.
func (c *container) WithGap(n int) *container {
	c.gap = n
	return c
}

// WithPadding sets the uniform space between the container's edges and its
// children, applied before children are laid out.
func (c *container) WithPadding(n int) *container {
	return c.WithPaddingXY(n, n)
}

// WithPaddingXY sets the space between the container's edges and its children
// with independent horizontal (x) and vertical (y) values.
func (c *container) WithPaddingXY(x, y int) *container {
	c.paddingX = x
	c.paddingY = y
	return c
}

// WithItem adds the given items as children of the container.
func (c *container) WithItem(items ...*ContainerItem) *container {
	c.items = append(c.items, items...)
	return c
}

// RemoveItem removes the given items from the container's children, if
// present. Items not added via WithItem (e.g. those wrapped internally by
// WithComponent) cannot be matched and are ignored.
func (c *container) RemoveItem(items ...*ContainerItem) *container {
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
func (c *container) Clear() *container {
	c.items = nil
	return c
}

// Layout computes the layout of the container's children given the available
// width and height. It returns a slice of Rects representing the position
// and size of each child.
func (c *container) Layout(width, height int) []Rect {
	items := c.visibleItems()
	n := len(items)
	rects := make([]Rect, n)
	if n == 0 {
		return rects
	}

	width = max(0, width-2*c.paddingX)
	height = max(0, height-2*c.paddingY)

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
			bases[i] = item.mainWidth(width)
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
		offsetRects(rects, c.paddingX, c.paddingY)
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
		bases[i] = item.mainHeight(widths[i], height)
		mins[i] = item.minHeight
		maxs[i] = item.maxHeight
	}
	heights, ys := resolveAxis(bases, grows, shrinks, mins, maxs, height, c.gap, c.justify)

	for i, item := range items {
		rects[i] = Rect{X: crossOffset(item.resolveAlign(c.align), width, widths[i]), Y: ys[i], Width: widths[i], Height: heights[i]}
	}
	offsetRects(rects, c.paddingX, c.paddingY)
	return rects
}

// Render renders the container and its children into a slice of strings,
// each representing a line of text. Children are composited as cells, so
// their lines may contain SGR escape sequences and wide characters. The
// width and height parameters specify the available space for rendering.
func (c *container) Render(width, height int) []string {
	rects := c.Layout(width, height)
	grid := make([][]Cell, height)
	for y := range grid {
		grid[y] = make([]Cell, width)
		for x := range grid[y] {
			grid[y][x] = Cell{Text: " ", Width: 1}
		}
	}

	for i, it := range c.visibleItems() {
		r := rects[i]
		lines := it.component.Render(r.Width, r.Height)
		for dy, line := range lines {
			y := r.Y + dy
			if dy >= r.Height || y < 0 || y >= height {
				continue
			}
			for dx, cl := range BuildRow(line, r.Width) {
				x := r.X + dx
				if x < 0 || x >= width {
					continue
				}
				// A wide character whose placeholder would be clipped by the
				// container's edge can't be drawn whole; blank it instead.
				if cl.Width == 2 && x+1 >= width {
					cl = Cell{Text: " ", Width: 1, Style: cl.Style}
				}
				grid[y][x] = cl
			}
		}
	}

	out := make([]string, height)
	for y, row := range grid {
		out[y] = RenderRow(row)
	}
	return out
}

// PreferredWidth returns the preferred width of the container, which is the
// sum of the preferred widths of its children plus the gaps between them.
func (c *container) PreferredWidth() int {
	items := c.visibleItems()
	n := len(items)
	if n == 0 {
		return 0
	}

	if c.direction == DirectionRow {
		sum := c.gap * (n - 1)
		for _, item := range items {
			sum += item.mainWidth(-1)
		}
		return sum + 2*c.paddingX
	}

	w := 0
	for _, item := range items {
		w = max(w, item.crossWidth())
	}
	return w + 2*c.paddingX
}

// PreferredHeight returns the preferred height of the container given a
// specific width.
func (c *container) PreferredHeight(width int) int {
	items := c.visibleItems()
	n := len(items)
	if n == 0 {
		return 0
	}

	width = max(0, width-2*c.paddingX)

	if c.direction == DirectionColumn {
		sum := c.gap * (n - 1)
		for _, item := range items {
			w := width
			if item.resolveAlign(c.align) != AlignStretch {
				w = item.crossWidth()
			}
			sum += item.mainHeight(w, -1)
		}
		return sum + 2*c.paddingY
	}

	bases := make([]int, n)
	grows := make([]int, n)
	shrinks := make([]int, n)
	mins := make([]int, n)
	maxs := make([]int, n)
	for i, item := range items {
		bases[i] = item.mainWidth(width)
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
	return h + 2*c.paddingY
}

func (c *container) visibleItems() []*ContainerItem {
	items := make([]*ContainerItem, 0, len(c.items))
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
