package term

import "testing"

func TestNewSelect(t *testing.T) {
	s := NewSelect("a", "b", "c")
	opts := s.Options()
	if len(opts) != 3 || opts[0] != "a" || opts[1] != "b" || opts[2] != "c" {
		t.Fatalf("Options() = %#v, want [a b c]", opts)
	}
	if s.Marker() != DefaultSelectMarker {
		t.Fatalf("Marker() = %q, want %q", s.Marker(), DefaultSelectMarker)
	}
	if s.Value() != "a" {
		t.Fatalf("Value() = %q, want %q", s.Value(), "a")
	}
	if s.Focused() {
		t.Fatal("new select should not be focused")
	}
}

func TestSelectValueEmpty(t *testing.T) {
	s := NewSelect()
	if s.Value() != "" {
		t.Fatalf("Value() = %q, want empty", s.Value())
	}
}

func TestSelectWithOptionAppends(t *testing.T) {
	s := NewSelect("a")
	same := s.WithOption("b", "c")
	if same != s {
		t.Fatal("WithOption should return the same *Select for chaining")
	}
	opts := s.Options()
	if len(opts) != 3 || opts[0] != "a" || opts[1] != "b" || opts[2] != "c" {
		t.Fatalf("Options() = %#v, want [a b c]", opts)
	}
}

func TestSelectWithMarker(t *testing.T) {
	s := NewSelect("a").WithMarker("-> ")
	if s.Marker() != "-> " {
		t.Fatalf("Marker() = %q, want %q", s.Marker(), "-> ")
	}
}

func TestSelectWithValue(t *testing.T) {
	s := NewSelect("a", "b", "c")
	same := s.WithValue("b")
	if same != s {
		t.Fatal("WithValue should return the same *Select for chaining")
	}
	if s.Value() != "b" {
		t.Fatalf("Value() = %q, want %q", s.Value(), "b")
	}
}

func TestSelectWithValueIgnoresUnknown(t *testing.T) {
	s := NewSelect("a", "b").WithValue("z")
	if s.Value() != "a" {
		t.Fatalf("Value() = %q, want %q (unchanged)", s.Value(), "a")
	}
}

func TestSelectWithValuePtrReadsInitial(t *testing.T) {
	v := "b"
	s := NewSelect("a", "b", "c").WithValuePtr(&v)
	if s.Value() != "b" {
		t.Fatalf("Value() should reflect the bound pointer's initial value, got %q", s.Value())
	}
	if s.ValuePtr() != &v {
		t.Fatal("ValuePtr() should return the bound pointer")
	}
}

func TestSelectWithValuePtrSyncsOnMove(t *testing.T) {
	v := "a"
	s := NewSelect("a", "b", "c").WithValuePtr(&v).AsFocused(true)
	s.Update(KeyEvent{Type: KeyDown})
	if v != "b" {
		t.Fatalf("bound value = %q, want %q", v, "b")
	}
	if s.Value() != "b" {
		t.Fatalf("Value() = %q, want %q", s.Value(), "b")
	}
}

func TestSelectWithStyle(t *testing.T) {
	s := NewSelect("a").WithStyle(NewStyle().WithForeground(ColorRed))
	if s.Style() == nil {
		t.Fatal("Style() should not be nil after WithStyle")
	}
}

func TestSelectWithoutStyle(t *testing.T) {
	s := NewSelect("a").WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if s.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestSelectWithSelectedStyle(t *testing.T) {
	s := NewSelect("a").WithSelectedStyle(NewStyle().WithForeground(ColorRed))
	if s.SelectedStyle() == nil {
		t.Fatal("SelectedStyle() should not be nil after WithSelectedStyle")
	}
}

func TestSelectWithoutSelectedStyle(t *testing.T) {
	s := NewSelect("a").WithSelectedStyle(NewStyle().WithForeground(ColorRed)).WithoutSelectedStyle()
	if s.SelectedStyle() != nil {
		t.Fatal("WithoutSelectedStyle should clear the selected style")
	}
}

func TestSelectFocusBlurAsFocused(t *testing.T) {
	s := NewSelect("a")
	s.Focus()
	if !s.Focused() {
		t.Fatal("Focus() should set Focused() to true")
	}
	s.Blur()
	if s.Focused() {
		t.Fatal("Blur() should set Focused() to false")
	}
	if same := s.AsFocused(true); same != s {
		t.Fatal("AsFocused should return the same *Select for chaining")
	}
	if !s.Focused() {
		t.Fatal("AsFocused(true) should set Focused() to true")
	}
}

func TestSelectUpdateIgnoresKeysWhenBlurred(t *testing.T) {
	s := NewSelect("a", "b")
	s.Update(KeyEvent{Type: KeyDown})
	if s.Value() != "a" {
		t.Fatal("Update should ignore key events while blurred")
	}
}

func TestSelectUpdateMovesWithArrowKeys(t *testing.T) {
	s := NewSelect("a", "b", "c").AsFocused(true)
	s.Update(KeyEvent{Type: KeyDown})
	if s.Value() != "b" {
		t.Fatalf("Down should move to %q, got %q", "b", s.Value())
	}
	s.Update(KeyEvent{Type: KeyDown})
	if s.Value() != "c" {
		t.Fatalf("Down should move to %q, got %q", "c", s.Value())
	}
	s.Update(KeyEvent{Type: KeyUp})
	if s.Value() != "b" {
		t.Fatalf("Up should move to %q, got %q", "b", s.Value())
	}
}

func TestSelectUpdateClampsAtBounds(t *testing.T) {
	s := NewSelect("a", "b").AsFocused(true)
	s.Update(KeyEvent{Type: KeyUp})
	if s.Value() != "a" {
		t.Fatalf("Up at first option should stay at %q, got %q", "a", s.Value())
	}
	s.Update(KeyEvent{Type: KeyDown})
	s.Update(KeyEvent{Type: KeyDown})
	if s.Value() != "b" {
		t.Fatalf("Down at last option should stay at %q, got %q", "b", s.Value())
	}
}

func TestSelectRenderHighlightsCursor(t *testing.T) {
	s := NewSelect("one", "two", "three")
	lines := s.Render(20, 3)
	want := []string{"> one", "  two", "  three"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestSelectRenderMovesHighlight(t *testing.T) {
	s := NewSelect("one", "two").AsFocused(true)
	s.Update(KeyEvent{Type: KeyDown})
	lines := s.Render(20, 2)
	want := []string{"  one", "> two"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestSelectRenderScrollsToKeepCursorVisible(t *testing.T) {
	s := NewSelect("one", "two", "three", "four").AsFocused(true)
	s.Update(KeyEvent{Type: KeyDown})
	s.Update(KeyEvent{Type: KeyDown})
	s.Update(KeyEvent{Type: KeyDown})
	lines := s.Render(20, 2)
	want := []string{"  three", "> four"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestSelectRenderZeroSize(t *testing.T) {
	s := NewSelect("a")
	if lines := s.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := s.Render(5, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestSelectRenderNoOptions(t *testing.T) {
	s := NewSelect()
	if lines := s.Render(10, 3); lines != nil {
		t.Fatalf("got %#v, want nil for no options", lines)
	}
}

func TestSelectPreferredWidth(t *testing.T) {
	s := NewSelect("a", "longer")
	// DefaultSelectMarker "> " (2) + "longer" (6) = 8
	if w := s.PreferredWidth(); w != 8 {
		t.Fatalf("PreferredWidth() = %d, want 8", w)
	}
}

func TestSelectPreferredHeight(t *testing.T) {
	s := NewSelect("a", "b", "c")
	if h := s.PreferredHeight(30); h != 3 {
		t.Fatalf("PreferredHeight() = %d, want 3", h)
	}
}
