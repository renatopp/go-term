package ui

import "testing"

type plainComponent struct{}

func (plainComponent) Render(width, height int) []string { return nil }

func TestItemDefaults(t *testing.T) {
	i := NewItem(plainComponent{})

	if i.grow != 0 {
		t.Fatalf("default grow = %d, want 0", i.grow)
	}
	if i.shrink != 1 {
		t.Fatalf("default shrink = %d, want 1", i.shrink)
	}
	if i.mainWidth() != 0 {
		t.Fatalf("mainWidth without Sized = %d, want 0", i.mainWidth())
	}
}

func TestItemBasisOverridesPreferredSize(t *testing.T) {
	c := &fakeComponent{width: 5, height: func(int) int { return 3 }}
	i := NewItem(c).WithBasis(20)

	if i.mainWidth() != 20 {
		t.Fatalf("mainWidth with basis = %d, want 20", i.mainWidth())
	}
	if i.mainHeight(100) != 20 {
		t.Fatalf("mainHeight with basis = %d, want 20", i.mainHeight(100))
	}
	if i.crossWidth() != 5 {
		t.Fatalf("crossWidth should ignore basis, got %d, want 5", i.crossWidth())
	}
	if i.crossHeight(100) != 3 {
		t.Fatalf("crossHeight should ignore basis, got %d, want 3", i.crossHeight(100))
	}
}

func TestItemWithersChain(t *testing.T) {
	i := NewItem(plainComponent{})
	same := i.WithGrow(2).WithShrink(3).WithBasis(4)

	if same != i {
		t.Fatal("With* methods should return the same *Item for chaining")
	}
	if i.grow != 2 || i.shrink != 3 || i.basis != 4 {
		t.Fatalf("got grow=%d shrink=%d basis=%d, want 2/3/4", i.grow, i.shrink, i.basis)
	}
}

func TestItemMinMaxClampPreferredSize(t *testing.T) {
	c := &fakeComponent{width: 5, height: func(int) int { return 5 }}
	i := NewItem(c).WithMinWidth(10).WithMaxHeight(2)

	if w := i.mainWidth(); w != 10 {
		t.Fatalf("mainWidth below min = %d, want 10", w)
	}
	if h := i.mainHeight(100); h != 2 {
		t.Fatalf("mainHeight above max = %d, want 2", h)
	}
	if w := i.crossWidth(); w != 10 {
		t.Fatalf("crossWidth below min = %d, want 10", w)
	}
	if h := i.crossHeight(100); h != 2 {
		t.Fatalf("crossHeight above max = %d, want 2", h)
	}
}

func TestItemMaxUnboundedByDefault(t *testing.T) {
	i := NewItem(&fakeComponent{width: 1000, height: func(int) int { return 1000 }})

	if w := i.mainWidth(); w != 1000 {
		t.Fatalf("mainWidth = %d, want 1000 (unbounded)", w)
	}
	if h := i.mainHeight(0); h != 1000 {
		t.Fatalf("mainHeight = %d, want 1000 (unbounded)", h)
	}
}

func TestItemResolveAlignDefaultsToAuto(t *testing.T) {
	i := NewItem(plainComponent{})

	if i.alignSelf != AlignAuto {
		t.Fatalf("default alignSelf = %v, want AlignAuto", i.alignSelf)
	}
	if got := i.resolveAlign(AlignCenter); got != AlignCenter {
		t.Fatalf("resolveAlign with AlignAuto = %v, want inherited AlignCenter", got)
	}
}

func TestItemResolveAlignOverridesContainer(t *testing.T) {
	i := NewItem(plainComponent{}).WithAlignSelf(AlignEnd)

	if got := i.resolveAlign(AlignCenter); got != AlignEnd {
		t.Fatalf("resolveAlign with explicit alignSelf = %v, want AlignEnd", got)
	}
}

func TestItemAsHiddenDefaultsToVisible(t *testing.T) {
	i := NewItem(plainComponent{})

	if i.hidden {
		t.Fatal("default hidden = true, want false")
	}
}

func TestItemAsHiddenChains(t *testing.T) {
	i := NewItem(plainComponent{})
	same := i.AsHidden(true)

	if same != i {
		t.Fatal("AsHidden should return the same *Item for chaining")
	}
	if !i.hidden {
		t.Fatal("hidden = false after AsHidden(true), want true")
	}
}
