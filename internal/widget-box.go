package term

// BorderStyle defines the characters used to draw a Box's border. The zero
// value draws no border at all.
type BorderStyle struct {
	Top         string
	Bottom      string
	Left        string
	Right       string
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
}

var (
	BorderNone = BorderStyle{}

	BorderSingle = BorderStyle{
		Top: "─", Bottom: "─", Left: "│", Right: "│",
		TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘",
	}

	BorderDouble = BorderStyle{
		Top: "═", Bottom: "═", Left: "║", Right: "║",
		TopLeft: "╔", TopRight: "╗", BottomLeft: "╚", BottomRight: "╝",
	}

	BorderRound = BorderStyle{
		Top: "─", Bottom: "─", Left: "│", Right: "│",
		TopLeft: "╭", TopRight: "╮", BottomLeft: "╰", BottomRight: "╯",
	}

	BorderThick = BorderStyle{
		Top: "━", Bottom: "━", Left: "┃", Right: "┃",
		TopLeft: "┏", TopRight: "┓", BottomLeft: "┗", BottomRight: "┛",
	}
)

// hasBorder reports whether b draws any border at all.
func (b BorderStyle) hasBorder() bool {
	return b != BorderStyle{}
}

// Box draws an optional border around a single child, with padding between
// the border and the child and margin outside the border.
type Box struct {
	child    Renderable
	border   BorderStyle
	style    *Style
	paddingX int
	paddingY int
	marginX  int
	marginY  int
}

func NewBox(child Renderable) *Box {
	return &Box{child: child}
}

func (b *Box) Child() Renderable {
	return b.child
}

func (b *Box) WithChild(child Renderable) *Box {
	b.child = child
	return b
}

func (b *Box) Border() BorderStyle {
	return b.border
}

// WithBorder sets the characters used to draw the box's border.
func (b *Box) WithBorder(border BorderStyle) *Box {
	b.border = border
	return b
}

// WithoutBorder removes the box's border.
func (b *Box) WithoutBorder() *Box {
	b.border = BorderStyle{}
	return b
}

func (b *Box) Style() *Style {
	return b.style
}

// WithStyle sets the style applied to the box's border characters when rendering.
func (b *Box) WithStyle(s Style) *Box {
	b.style = &s
	return b
}

// WithoutStyle removes the box's style, rendering the border unstyled.
func (b *Box) WithoutStyle() *Box {
	b.style = nil
	return b
}

// WithPadding sets the uniform space between the box's border and its child.
func (b *Box) WithPadding(n int) *Box {
	return b.WithPaddingXY(n, n)
}

// WithPaddingXY sets the space between the box's border and its child with
// independent horizontal (x) and vertical (y) values.
func (b *Box) WithPaddingXY(x, y int) *Box {
	b.paddingX = x
	b.paddingY = y
	return b
}

// WithMargin sets the uniform space between the box's border and its
// surroundings.
func (b *Box) WithMargin(n int) *Box {
	return b.WithMarginXY(n, n)
}

// WithMarginXY sets the space between the box's border and its surroundings
// with independent horizontal (x) and vertical (y) values.
func (b *Box) WithMarginXY(x, y int) *Box {
	b.marginX = x
	b.marginY = y
	return b
}

func (b *Box) PreferredWidth() int {
	w := 0
	if s, ok := b.child.(Sizeable); ok {
		w = s.PreferredWidth()
	}
	return w + 2*b.paddingX + 2*b.marginX + 2*b.borderWidth()
}

func (b *Box) PreferredHeight(width int) int {
	inner := max(0, width-2*b.paddingX-2*b.marginX-2*b.borderWidth())
	h := 0
	if s, ok := b.child.(Sizeable); ok {
		h = s.PreferredHeight(inner)
	}
	return h + 2*b.paddingY + 2*b.marginY + 2*b.borderWidth()
}

func (b *Box) Update(e Event) Event {
	if u, ok := b.child.(Updatable); ok {
		return u.Update(e)
	}
	return e
}

// Render draws the box's border (if any) and child within width and height,
// with margin reserved outside the border and padding reserved between the
// border and the child.
func (b *Box) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	bw := b.borderWidth()
	left := b.marginX + bw + b.paddingX
	right := left
	top := b.marginY + bw + b.paddingY
	bottom := top

	contentWidth := max(0, width-left-right)
	contentHeight := max(0, height-top-bottom)

	var childLines []string
	if b.child != nil && contentWidth > 0 && contentHeight > 0 {
		childLines = b.child.Render(contentWidth, contentHeight)
	}

	grid := make([][]Cell, height)
	for y := range grid {
		grid[y] = make([]Cell, width)
		for x := range grid[y] {
			grid[y][x] = Cell{Text: " ", Width: 1}
		}
	}

	if b.border.hasBorder() {
		b.renderBorder(grid, width, height)
	}

	for dy, line := range childLines {
		y := top + dy
		if dy >= contentHeight || y < 0 || y >= height {
			continue
		}
		for dx, cl := range BuildRow(line, contentWidth) {
			x := left + dx
			if x < 0 || x >= width {
				continue
			}
			// A wide character whose placeholder would be clipped by the
			// box's edge can't be drawn whole; blank it instead.
			if cl.Width == 2 && x+1 >= width {
				cl = Cell{Text: " ", Width: 1, Style: cl.Style}
			}
			grid[y][x] = cl
		}
	}

	out := make([]string, height)
	for y, row := range grid {
		out[y] = RenderRow(row)
	}
	return out
}

// borderWidth returns the thickness, in columns or rows, that the border
// occupies on a single side.
func (b *Box) borderWidth() int {
	if b.border.hasBorder() {
		return 1
	}
	return 0
}

// renderBorder draws the box's border edges and corners onto grid at the
// position determined by the box's margin.
func (b *Box) renderBorder(grid [][]Cell, width, height int) {
	x0 := b.marginX
	x1 := width - 1 - b.marginX
	y0 := b.marginY
	y1 := height - 1 - b.marginY
	// A border needs at least 2 distinct columns and rows so its corners
	// don't collapse into the same cell.
	if x0 < 0 || y0 < 0 || x1 <= x0 || y1 <= y0 || x1 >= width || y1 >= height {
		return
	}

	b.setCell(grid, x0, y0, b.border.TopLeft)
	b.setCell(grid, x1, y0, b.border.TopRight)
	b.setCell(grid, x0, y1, b.border.BottomLeft)
	b.setCell(grid, x1, y1, b.border.BottomRight)
	for x := x0 + 1; x < x1; x++ {
		b.setCell(grid, x, y0, b.border.Top)
		b.setCell(grid, x, y1, b.border.Bottom)
	}
	for y := y0 + 1; y < y1; y++ {
		b.setCell(grid, x0, y, b.border.Left)
		b.setCell(grid, x1, y, b.border.Right)
	}
}

// setCell writes a single-width cell at (x, y) in grid, if in bounds. An
// empty ch (an unset border character) renders as a blank space so column
// alignment isn't disturbed. If the box has a style, it's applied to the
// character.
func (b *Box) setCell(grid [][]Cell, x, y int, ch string) {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return
	}
	if ch == "" {
		ch = " "
	}
	cell := Cell{Text: ch, Width: 1}
	if b.style != nil {
		cell = BuildRow(b.style.Render(ch), 1)[0]
	}
	grid[y][x] = cell
}
