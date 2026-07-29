package ui

import (
	"strings"
	"testing"
)

type fakeComponent struct {
	width  int
	height func(width int) int
	render func(width, height int) []string
}

func (f *fakeComponent) Render(width, height int) []string {
	if f.render != nil {
		return f.render(width, height)
	}
	return nil
}

func (f *fakeComponent) PreferredWidth() int {
	return f.width
}

func (f *fakeComponent) PreferredHeight(width int) int {
	if f.height != nil {
		return f.height(width)
	}
	return 0
}

func TestContainerLayoutRowFixedWidths(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10),
		NewItem(&fakeComponent{}).WithBasis(20),
		NewItem(&fakeComponent{}).WithBasis(30),
	)

	rects := c.Layout(100, 5)

	want := []Rect{
		{X: 0, Y: 0, Width: 10, Height: 5},
		{X: 10, Y: 0, Width: 20, Height: 5},
		{X: 30, Y: 0, Width: 30, Height: 5},
	}
	for i, r := range rects {
		if r != want[i] {
			t.Errorf("item %d: got %+v, want %+v", i, r, want[i])
		}
	}
}

func TestContainerLayoutRowGrowFillsRemaining(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10),
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(1),
	)

	rects := c.Layout(100, 5)

	if rects[0].Width != 10 {
		t.Fatalf("fixed item width = %d, want 10", rects[0].Width)
	}
	if rects[1].Width != 90 {
		t.Fatalf("grow item width = %d, want 90", rects[1].Width)
	}
	if rects[1].X != 10 {
		t.Fatalf("grow item x = %d, want 10", rects[1].X)
	}
}

func TestContainerLayoutRowGrowWeighted(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(1),
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(2),
	)

	rects := c.Layout(90, 5)

	if rects[0].Width != 30 || rects[1].Width != 60 {
		t.Fatalf("got widths %d/%d, want 30/60", rects[0].Width, rects[1].Width)
	}
}

func TestContainerLayoutRowShrinkOverflow(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(60),
		NewItem(&fakeComponent{}).WithBasis(60),
	)

	rects := c.Layout(100, 5)

	total := rects[0].Width + rects[1].Width
	if total != 100 {
		t.Fatalf("shrunk total width = %d, want 100", total)
	}
}

func TestContainerLayoutRowGrowRespectsMaxWidth(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(1).WithMaxWidth(30),
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(1),
	)

	rects := c.Layout(100, 5)

	if rects[0].Width != 30 {
		t.Fatalf("capped item width = %d, want 30 (maxWidth)", rects[0].Width)
	}
}

func TestContainerLayoutRowShrinkRespectsMinWidth(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(60).WithMinWidth(55),
		NewItem(&fakeComponent{}).WithBasis(60),
	)

	rects := c.Layout(100, 5)

	if rects[0].Width < 55 {
		t.Fatalf("floored item width = %d, want >= 55 (minWidth)", rects[0].Width)
	}
}

func TestContainerLayoutRowStretchRespectsMaxHeight(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10).WithMaxHeight(3),
	)

	rects := c.Layout(20, 10)

	if rects[0].Height != 3 {
		t.Fatalf("stretched height = %d, want 3 (maxHeight)", rects[0].Height)
	}
}

func TestContainerLayoutColumnStretchRespectsMinWidth(t *testing.T) {
	c := NewContainer(DirectionColumn).WithItem(
		NewItem(&fakeComponent{}).WithBasis(5).WithMinWidth(15),
	)

	rects := c.Layout(10, 20)

	if rects[0].Width != 15 {
		t.Fatalf("stretched width = %d, want 15 (minWidth)", rects[0].Width)
	}
}

func TestContainerLayoutJustifyCenter(t *testing.T) {
	c := NewContainer(DirectionRow).WithJustify(JustifyCenter).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10),
	)

	rects := c.Layout(100, 5)

	if rects[0].X != 45 {
		t.Fatalf("centered x = %d, want 45", rects[0].X)
	}
}

func TestContainerLayoutJustifySpaceBetween(t *testing.T) {
	c := NewContainer(DirectionRow).WithJustify(JustifySpaceBetween).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10),
		NewItem(&fakeComponent{}).WithBasis(10),
		NewItem(&fakeComponent{}).WithBasis(10),
	)

	rects := c.Layout(100, 5)

	if rects[0].X != 0 || rects[2].X != 90 {
		t.Fatalf("got x positions %d/.../%d, want 0/.../90", rects[0].X, rects[2].X)
	}
}

func TestContainerLayoutAlignCenterCrossAxis(t *testing.T) {
	c := NewContainer(DirectionRow).WithAlign(AlignCenter).WithItem(
		NewItem(&fakeComponent{width: 5, height: func(int) int { return 2 }}),
	)

	rects := c.Layout(20, 10)

	if rects[0].Height != 2 {
		t.Fatalf("height = %d, want 2", rects[0].Height)
	}
	if rects[0].Y != 4 {
		t.Fatalf("y = %d, want 4", rects[0].Y)
	}
}

func TestContainerLayoutAlignSelfOverridesContainerAlignRow(t *testing.T) {
	c := NewContainer(DirectionRow).WithAlign(AlignStart).WithItem(
		NewItem(&fakeComponent{width: 5, height: func(int) int { return 2 }}).WithAlignSelf(AlignCenter),
	)

	rects := c.Layout(20, 10)

	if rects[0].Y != 4 {
		t.Fatalf("y = %d, want 4 (centered via alignSelf)", rects[0].Y)
	}
}

func TestContainerLayoutAlignSelfOverridesContainerAlignColumn(t *testing.T) {
	c := NewContainer(DirectionColumn).WithAlign(AlignStart).WithItem(
		NewItem(&fakeComponent{width: 5}).WithAlignSelf(AlignStretch),
	)

	rects := c.Layout(20, 10)

	if rects[0].Width != 20 {
		t.Fatalf("width = %d, want 20 (stretched via alignSelf)", rects[0].Width)
	}
}

func TestContainerLayoutHiddenItemExcludedFromFlow(t *testing.T) {
	c := NewContainer(DirectionRow).WithGap(5).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10),
		NewItem(&fakeComponent{}).WithBasis(10).AsHidden(true),
		NewItem(&fakeComponent{}).WithBasis(10),
	)

	rects := c.Layout(100, 5)

	if len(rects) != 2 {
		t.Fatalf("got %d rects, want 2 (hidden item excluded)", len(rects))
	}
	if rects[1].X != 15 {
		t.Fatalf("second visible item x = %d, want 15 (no gap reserved for the hidden item)", rects[1].X)
	}
}

func TestContainerPreferredWidthExcludesHidden(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{width: 5}).WithBasis(5),
		NewItem(&fakeComponent{width: 100}).WithBasis(100).AsHidden(true),
	)

	if w := c.PreferredWidth(); w != 5 {
		t.Fatalf("preferred width = %d, want 5 (hidden item excluded)", w)
	}
}

func TestContainerRenderSkipsHiddenItem(t *testing.T) {
	visible := &fakeComponent{render: func(w, h int) []string { return []string{"AA"} }}
	hidden := &fakeComponent{render: func(w, h int) []string { return []string{"XX"} }}

	c := NewContainer(DirectionRow).WithItem(
		NewItem(visible).WithBasis(2),
		NewItem(hidden).WithBasis(2).AsHidden(true),
	)

	out := c.Render(4, 1)
	if out[0] != "AA  " {
		t.Fatalf("got %q, want %q (hidden item not rendered)", out[0], "AA  ")
	}
}

func TestContainerRemoveItemRemovesMatchingPointer(t *testing.T) {
	first := NewItem(&fakeComponent{}).WithBasis(10)
	middle := NewItem(&fakeComponent{}).WithBasis(10)
	last := NewItem(&fakeComponent{}).WithBasis(10)

	c := NewContainer(DirectionRow).WithGap(5).WithItem(first, middle, last)
	c.RemoveItem(middle)

	rects := c.Layout(100, 5)
	if len(rects) != 2 {
		t.Fatalf("got %d rects, want 2 (middle item removed)", len(rects))
	}
	if rects[1].X != 15 {
		t.Fatalf("second item x = %d, want 15 (no gap reserved for the removed item)", rects[1].X)
	}
}

func TestContainerRemoveItemVariadicRemovesMultiple(t *testing.T) {
	a := NewItem(&fakeComponent{}).WithBasis(10)
	b := NewItem(&fakeComponent{}).WithBasis(10)
	keep := NewItem(&fakeComponent{}).WithBasis(10)

	c := NewContainer(DirectionRow).WithItem(a, b, keep)
	c.RemoveItem(a, b)

	if len(c.items) != 1 || c.items[0] != keep {
		t.Fatalf("got %d items, want 1 matching keep", len(c.items))
	}
}

func TestContainerRemoveItemIgnoresUnknownItem(t *testing.T) {
	kept := NewItem(&fakeComponent{}).WithBasis(10)
	unrelated := NewItem(&fakeComponent{}).WithBasis(10)

	c := NewContainer(DirectionRow).WithItem(kept)
	c.RemoveItem(unrelated)

	if len(c.items) != 1 || c.items[0] != kept {
		t.Fatalf("got %d items, want 1 unchanged item", len(c.items))
	}
}

func TestContainerRemoveItemChains(t *testing.T) {
	c := NewContainer(DirectionRow)
	same := c.RemoveItem()

	if same != c {
		t.Fatal("RemoveItem should return the same *Container for chaining")
	}
}

func TestContainerClearRemovesAllItems(t *testing.T) {
	c := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{}).WithBasis(10),
		NewItem(&fakeComponent{}).WithBasis(10),
	)
	c.Clear()

	if len(c.items) != 0 {
		t.Fatalf("got %d items, want 0 after Clear", len(c.items))
	}
	if w := c.PreferredWidth(); w != 0 {
		t.Fatalf("preferred width = %d, want 0 after Clear", w)
	}
	if rects := c.Layout(100, 5); len(rects) != 0 {
		t.Fatalf("got %d rects, want 0 after Clear", len(rects))
	}
}

func TestContainerClearChains(t *testing.T) {
	c := NewContainer(DirectionRow)
	same := c.Clear()

	if same != c {
		t.Fatal("Clear should return the same *Container for chaining")
	}
}

func TestContainerLayoutColumnHeightDependsOnResolvedWidth(t *testing.T) {
	wrap := &fakeComponent{
		height: func(width int) int {
			return (40 + width - 1) / width // simulate wrapping 40 chars of text
		},
	}

	c := NewContainer(DirectionColumn).WithItem(NewItem(wrap))

	rects := c.Layout(10, 20)
	if rects[0].Width != 10 {
		t.Fatalf("width = %d, want 10", rects[0].Width)
	}
	if rects[0].Height != 4 {
		t.Fatalf("height = %d, want 4 (40 chars wrapped at width 10)", rects[0].Height)
	}
}

func TestContainerRenderComposites(t *testing.T) {
	left := &fakeComponent{render: func(w, h int) []string { return []string{"AA"} }}
	right := &fakeComponent{render: func(w, h int) []string { return []string{"BB"} }}

	c := NewContainer(DirectionRow).WithItem(
		NewItem(left).WithBasis(2),
		NewItem(right).WithBasis(2),
	)

	out := c.Render(6, 1)
	if len(out) != 1 {
		t.Fatalf("got %d lines, want 1", len(out))
	}
	if !strings.HasPrefix(out[0], "AABB") {
		t.Fatalf("got %q, want prefix AABB", out[0])
	}
}

func TestContainerLayoutPaddingOffsetsChildrenAndShrinksSpace(t *testing.T) {
	c := NewContainer(DirectionRow).WithPadding(2).WithItem(
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(1),
	)

	rects := c.Layout(20, 10)

	if rects[0].X != 2 || rects[0].Y != 2 {
		t.Fatalf("got X=%d Y=%d, want X=2 Y=2", rects[0].X, rects[0].Y)
	}
	if rects[0].Width != 16 {
		t.Fatalf("width = %d, want 16 (20 - 2*2)", rects[0].Width)
	}
	if rects[0].Height != 6 {
		t.Fatalf("height = %d, want 6 (10 - 2*2)", rects[0].Height)
	}
}

func TestContainerLayoutPaddingLargerThanSpaceClampsToZero(t *testing.T) {
	c := NewContainer(DirectionRow).WithPadding(20).WithItem(
		NewItem(&fakeComponent{}).WithBasis(0).WithGrow(1),
	)

	rects := c.Layout(10, 10)

	if rects[0].Width != 0 || rects[0].Height != 0 {
		t.Fatalf("got Width=%d Height=%d, want 0/0", rects[0].Width, rects[0].Height)
	}
}

func TestContainerPreferredSizeIncludesPadding(t *testing.T) {
	c := NewContainer(DirectionRow).WithPadding(3).WithItem(
		NewItem(&fakeComponent{width: 5, height: func(int) int { return 2 }}).WithBasis(5),
	)

	if w := c.PreferredWidth(); w != 11 {
		t.Fatalf("preferred width = %d, want 11 (5 + 2*3)", w)
	}
	if h := c.PreferredHeight(11); h != 8 {
		t.Fatalf("preferred height = %d, want 8 (2 + 2*3)", h)
	}
}

func TestContainerPreferredSizeNested(t *testing.T) {
	inner := NewContainer(DirectionRow).WithItem(
		NewItem(&fakeComponent{width: 5}).WithBasis(5),
		NewItem(&fakeComponent{width: 7}).WithBasis(7),
	)

	if w := inner.PreferredWidth(); w != 12 {
		t.Fatalf("preferred width = %d, want 12", w)
	}
}
