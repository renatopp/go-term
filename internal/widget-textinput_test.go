package term

import (
	"testing"
	"time"
)

func TestNewTextInput(t *testing.T) {
	ti := NewTextInput()
	if ti.Value() != "" {
		t.Fatalf("Value() = %q, want empty", ti.Value())
	}
	if ti.Focused() {
		t.Fatal("new text input should not be focused")
	}
	if ti.Width() != DefaultTextInputWidth {
		t.Fatalf("Width() = %d, want %d", ti.Width(), DefaultTextInputWidth)
	}
	if ti.Cursor() == nil {
		t.Fatal("Cursor() should not be nil")
	}
}

func TestTextInputWithValue(t *testing.T) {
	ti := NewTextInput()
	same := ti.WithValue("hello")
	if same != ti {
		t.Fatal("WithValue should return the same *TextInput for chaining")
	}
	if ti.Value() != "hello" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "hello")
	}
}

func TestTextInputWithValuePtrReadsInitial(t *testing.T) {
	v := "seed"
	ti := NewTextInput().WithValuePtr(&v)
	if ti.Value() != "seed" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "seed")
	}
	if ti.ValuePtr() != &v {
		t.Fatal("ValuePtr() should return the bound pointer")
	}
}

func TestTextInputWithValuePtrSyncsOnEdit(t *testing.T) {
	var v string
	ti := NewTextInput().WithValuePtr(&v).AsFocused(true)
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'b'})
	if v != "ab" {
		t.Fatalf("bound value = %q, want %q", v, "ab")
	}
	if ti.Value() != "ab" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "ab")
	}
}

func TestTextInputWithoutStyle(t *testing.T) {
	ti := NewTextInput().WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if ti.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestTextInputWithPlaceholder(t *testing.T) {
	ti := NewTextInput()
	same := ti.WithPlaceholder("type here")
	if same != ti {
		t.Fatal("WithPlaceholder should return the same *TextInput for chaining")
	}
	if ti.Placeholder() != "type here" {
		t.Fatalf("Placeholder() = %q, want %q", ti.Placeholder(), "type here")
	}
}

func TestTextInputWithPlaceholderStyle(t *testing.T) {
	ti := NewTextInput().WithPlaceholderStyle(NewStyle().WithForeground(ColorRed))
	if ti.PlaceholderStyle() == nil {
		t.Fatal("PlaceholderStyle() should not be nil after WithPlaceholderStyle")
	}

	ti.WithoutPlaceholderStyle()
	if ti.PlaceholderStyle() != nil {
		t.Fatal("WithoutPlaceholderStyle should clear the placeholder style")
	}
}

func TestTextInputWithPrefix(t *testing.T) {
	ti := NewTextInput()
	same := ti.WithPrefix(">> ")
	if same != ti {
		t.Fatal("WithPrefix should return the same *TextInput for chaining")
	}
	if ti.Prefix() != ">> " {
		t.Fatalf("Prefix() = %q, want %q", ti.Prefix(), ">> ")
	}
}

func TestTextInputWithPrefixStyle(t *testing.T) {
	ti := NewTextInput().WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	if ti.PrefixStyle() == nil {
		t.Fatal("PrefixStyle() should not be nil after WithPrefixStyle")
	}

	ti.WithoutPrefixStyle()
	if ti.PrefixStyle() != nil {
		t.Fatal("WithoutPrefixStyle should clear the prefix style")
	}
}

func TestTextInputWithSuffix(t *testing.T) {
	ti := NewTextInput()
	same := ti.WithSuffix(" <<")
	if same != ti {
		t.Fatal("WithSuffix should return the same *TextInput for chaining")
	}
	if ti.Suffix() != " <<" {
		t.Fatalf("Suffix() = %q, want %q", ti.Suffix(), " <<")
	}
}

func TestTextInputWithSuffixStyle(t *testing.T) {
	ti := NewTextInput().WithSuffixStyle(NewStyle().WithForeground(ColorRed))
	if ti.SuffixStyle() == nil {
		t.Fatal("SuffixStyle() should not be nil after WithSuffixStyle")
	}

	ti.WithoutSuffixStyle()
	if ti.SuffixStyle() != nil {
		t.Fatal("WithoutSuffixStyle should clear the suffix style")
	}
}

func TestTextInputPreferredWidthIncludesPrefixAndSuffix(t *testing.T) {
	ti := NewTextInput().WithWidth(10).WithPrefix(">> ").WithSuffix(" <<")
	if w := ti.PreferredWidth(); w != 16 {
		t.Fatalf("PreferredWidth() = %d, want 16", w)
	}
}

func TestTextInputWithWidth(t *testing.T) {
	ti := NewTextInput()
	same := ti.WithWidth(10)
	if same != ti {
		t.Fatal("WithWidth should return the same *TextInput for chaining")
	}
	if ti.Width() != 10 {
		t.Fatalf("Width() = %d, want 10", ti.Width())
	}
	if ti.PreferredWidth() != 10 {
		t.Fatalf("PreferredWidth() = %d, want 10", ti.PreferredWidth())
	}
}

func TestTextInputPreferredHeight(t *testing.T) {
	ti := NewTextInput()
	if h := ti.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}

func TestTextInputFocusBlur(t *testing.T) {
	ti := NewTextInput()
	ti.Focus()
	if !ti.Focused() {
		t.Fatal("Focused() should be true after Focus")
	}

	ti.Blur()
	if ti.Focused() {
		t.Fatal("Focused() should be false after Blur")
	}
}

func TestTextInputAsFocused(t *testing.T) {
	ti := NewTextInput()
	same := ti.AsFocused(true)
	if same != ti {
		t.Fatal("AsFocused should return the same *TextInput for chaining")
	}
	if !ti.Focused() {
		t.Fatal("Focused() should be true after AsFocused(true)")
	}

	ti.AsFocused(false)
	if ti.Focused() {
		t.Fatal("Focused() should be false after AsFocused(false)")
	}
}

func TestTextInputIgnoresKeysWhenBlurred(t *testing.T) {
	ti := NewTextInput()
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	if ti.Value() != "" {
		t.Fatalf("Value() = %q, want empty (blurred input should ignore keys)", ti.Value())
	}
}

func TestTextInputInsertsRuneWhenFocused(t *testing.T) {
	ti := NewTextInput().AsFocused(true)
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'b'})
	if ti.Value() != "ab" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "ab")
	}
}

func TestTextInputInsertAtCursorPosition(t *testing.T) {
	ti := NewTextInput().WithValue("ac").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyLeft})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'b'})
	if ti.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "abc")
	}
}

func TestTextInputIgnoresCtrlAndAltRunes(t *testing.T) {
	ti := NewTextInput().AsFocused(true)
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'c', Ctrl: true})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'x', Alt: true})
	if ti.Value() != "" {
		t.Fatalf("Value() = %q, want empty", ti.Value())
	}
}

func TestTextInputBackspace(t *testing.T) {
	ti := NewTextInput().WithValue("abc").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyBackspace})
	if ti.Value() != "ab" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "ab")
	}
}

func TestTextInputBackspaceAtStartIsNoop(t *testing.T) {
	ti := NewTextInput().AsFocused(true)
	ti.Update(KeyEvent{Type: KeyBackspace})
	if ti.Value() != "" {
		t.Fatalf("Value() = %q, want empty", ti.Value())
	}
}

func TestTextInputDelete(t *testing.T) {
	ti := NewTextInput().WithValue("abc").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyHome})
	ti.Update(KeyEvent{Type: KeyDelete})
	if ti.Value() != "bc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "bc")
	}
}

func TestTextInputDeleteAtEndIsNoop(t *testing.T) {
	ti := NewTextInput().WithValue("abc").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyDelete})
	if ti.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "abc")
	}
}

func TestTextInputHomeThenInsert(t *testing.T) {
	ti := NewTextInput().WithValue("bc").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyHome})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	if ti.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "abc")
	}
}

func TestTextInputEndThenInsert(t *testing.T) {
	ti := NewTextInput().WithValue("ab").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyHome})
	ti.Update(KeyEvent{Type: KeyEnd})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'c'})
	if ti.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "abc")
	}
}

func TestTextInputLeftAtStartIsNoop(t *testing.T) {
	ti := NewTextInput().WithValue("bc").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyLeft})
	ti.Update(KeyEvent{Type: KeyLeft})
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	if ti.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "abc")
	}
}

func TestTextInputRightAtEndIsNoop(t *testing.T) {
	ti := NewTextInput().WithValue("ab").AsFocused(true)
	for range 5 {
		ti.Update(KeyEvent{Type: KeyRight})
	}
	ti.Update(KeyEvent{Type: KeyRune, Rune: 'c'})
	if ti.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", ti.Value(), "abc")
	}
}

func TestTextInputUpdateReturnsEventUnchanged(t *testing.T) {
	ti := NewTextInput().AsFocused(true)
	e := KeyEvent{Type: KeyRune, Rune: 'a'}
	if got := ti.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestTextInputRenderPadsToWidth(t *testing.T) {
	ti := NewTextInput().WithValue("ab")
	lines := ti.Render(5, 1)
	if len(lines) != 1 || lines[0] != "ab   " {
		t.Fatalf("got %#v, want [\"ab   \"]", lines)
	}
}

func TestTextInputRenderShrinksToConfiguredWidth(t *testing.T) {
	// A configured width smaller than the render width caps the visible
	// window, padding the remainder with blanks up to the full render width.
	ti := NewTextInput().WithValue("abcdef").WithWidth(3).AsFocused(true)
	ti.Update(KeyEvent{Type: KeyHome})
	ti.Blur()
	lines := ti.Render(10, 1)
	if len(lines) != 1 || lines[0] != "abc       " {
		t.Fatalf("got %#v, want [\"abc       \"]", lines)
	}
}

func TestTextInputRenderDefaultWidthFillsRenderWidth(t *testing.T) {
	// An unconfigured input's width (DefaultTextInputWidth) is effectively
	// unbounded, so it fills whatever width it's rendered with.
	ti := NewTextInput().WithValue("abcdef").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyHome})
	ti.Blur()
	lines := ti.Render(4, 1)
	if len(lines) != 1 || lines[0] != "abcd" {
		t.Fatalf("got %#v, want [\"abcd\"]", lines)
	}
}

func TestTextInputRenderZeroSize(t *testing.T) {
	ti := NewTextInput()
	if lines := ti.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := ti.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestTextInputRenderHidesCursorWhenBlurred(t *testing.T) {
	ti := NewTextInput().WithValue("ab")
	lines := ti.Render(5, 1)
	if len(lines) != 1 || lines[0] != "ab   " {
		t.Fatalf("got %#v, want [\"ab   \"]", lines)
	}
}

func TestTextInputRenderShowsCursorAtEndWhenFocusedAndEmpty(t *testing.T) {
	ti := NewTextInput().AsFocused(true)
	lines := ti.Render(3, 1)
	want := DefaultCursorChar + "  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderCursorOverlaysCharacterAtPosition(t *testing.T) {
	ti := NewTextInput().WithValue("ab").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyLeft})
	lines := ti.Render(5, 1)
	want := "a" + DefaultCursorChar + "   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderRevealsCharacterWhenCursorBlinkedOff(t *testing.T) {
	ti := NewTextInput().WithValue("ab").AsFocused(true)
	ti.Update(KeyEvent{Type: KeyLeft})

	base := time.Unix(0, 0)
	ti.Cursor().WithBlinkSpeed(10 * time.Millisecond)
	ti.Cursor().Tick(base)
	ti.Cursor().Tick(base.Add(15 * time.Millisecond))

	lines := ti.Render(5, 1)
	want := "ab   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputCursorCustomChar(t *testing.T) {
	ti := NewTextInput().AsFocused(true)
	ti.Cursor().WithChar("_")
	lines := ti.Render(3, 1)
	want := "_  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderScrollsRightToKeepCursorVisible(t *testing.T) {
	ti := NewTextInput().WithValue("abcdef").WithWidth(3).AsFocused(true)
	lines := ti.Render(3, 1)
	want := "ef" + DefaultCursorChar
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderScrollsLeftAfterMovingCursorBack(t *testing.T) {
	ti := NewTextInput().WithValue("abcdef").WithWidth(3).AsFocused(true)
	ti.Render(3, 1)
	for range 6 {
		ti.Update(KeyEvent{Type: KeyLeft})
	}
	lines := ti.Render(3, 1)
	want := DefaultCursorChar + "bc"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	ti := NewTextInput().WithValue("ab").WithStyle(NewStyle().WithForeground(ColorRed))
	lines := ti.Render(2, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "ab" {
		t.Fatalf("expected styled text to contain an SGR sequence, got %q", lines[0])
	}
}

func TestTextInputRenderWithPrefixAndSuffix(t *testing.T) {
	ti := NewTextInput().WithValue("ab").WithPrefix("[").WithSuffix("]")
	lines := ti.Render(6, 1)
	want := "[ab  ]"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderClipsPrefixAndSuffixToWidth(t *testing.T) {
	// The field shrinks to make room for the prefix and suffix, and the
	// suffix is dropped entirely once the prefix alone fills the width.
	ti := NewTextInput().WithValue("x").WithPrefix("abcde").WithSuffix("fghij")
	lines := ti.Render(3, 1)
	want := "abc"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderWithPrefixAndSuffixStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	ti := NewTextInput().
		WithValue("ab").
		WithPrefix("[").WithPrefixStyle(NewStyle().WithForeground(ColorRed)).
		WithSuffix("]").WithSuffixStyle(NewStyle().WithForeground(ColorBlue))
	lines := ti.Render(6, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "[ab  ]" {
		t.Fatalf("expected styled prefix/suffix to contain SGR sequences, got %q", lines[0])
	}
}

func TestTextInputRenderShowsPlaceholderWhenEmpty(t *testing.T) {
	ti := NewTextInput().WithPlaceholder("search")
	lines := ti.Render(10, 1)
	want := "search    "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderHidesPlaceholderWhenValueSet(t *testing.T) {
	ti := NewTextInput().WithValue("ab").WithPlaceholder("search")
	lines := ti.Render(10, 1)
	want := "ab        "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderPlaceholderCursorAtStartWhenFocused(t *testing.T) {
	ti := NewTextInput().WithPlaceholder("abc").AsFocused(true)
	lines := ti.Render(5, 1)
	want := DefaultCursorChar + "bc  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestTextInputRenderWithPlaceholderStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	ti := NewTextInput().WithPlaceholder("ab").WithPlaceholderStyle(NewStyle().WithForeground(ColorRed))
	lines := ti.Render(2, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "ab" {
		t.Fatalf("expected styled placeholder to contain an SGR sequence, got %q", lines[0])
	}
}
