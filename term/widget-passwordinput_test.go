package term

import (
	"strings"
	"testing"
	"time"
)

func TestNewPasswordInput(t *testing.T) {
	pw := NewPasswordInput()
	if pw.Value() != "" {
		t.Fatalf("Value() = %q, want empty", pw.Value())
	}
	if pw.Focused() {
		t.Fatal("new PasswordInput should not be focused")
	}
	if pw.Width() != DefaultPasswordInputWidth {
		t.Fatalf("Width() = %d, want %d", pw.Width(), DefaultPasswordInputWidth)
	}
	if pw.Cursor() == nil {
		t.Fatal("Cursor() should not be nil")
	}
	if pw.MaskChar() != DefaultPasswordInputMaskChar {
		t.Fatalf("MaskChar() = %q, want %q", pw.MaskChar(), DefaultPasswordInputMaskChar)
	}
}

func TestPasswordInputWithValue(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.WithValue("hello")
	if same != pw {
		t.Fatal("WithValue should return the same *PasswordInput for chaining")
	}
	if pw.Value() != "hello" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "hello")
	}
}

func TestPasswordInputWithValuePtrReadsInitial(t *testing.T) {
	v := "seed"
	pw := NewPasswordInput().WithValuePtr(&v)
	if pw.Value() != "seed" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "seed")
	}
	if pw.ValuePtr() != &v {
		t.Fatal("ValuePtr() should return the bound pointer")
	}
}

func TestPasswordInputWithValuePtrSyncsOnEdit(t *testing.T) {
	var v string
	pw := NewPasswordInput().WithValuePtr(&v).AsFocused(true)
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'b'})
	if v != "ab" {
		t.Fatalf("bound value = %q, want %q", v, "ab")
	}
	if pw.Value() != "ab" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "ab")
	}
}

func TestPasswordInputWithoutStyle(t *testing.T) {
	pw := NewPasswordInput().WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if pw.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestPasswordInputWithMaskChar(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.WithMaskChar("*")
	if same != pw {
		t.Fatal("WithMaskChar should return the same *PasswordInput for chaining")
	}
	if pw.MaskChar() != "*" {
		t.Fatalf("MaskChar() = %q, want %q", pw.MaskChar(), "*")
	}
}

func TestPasswordInputWithPlaceholder(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.WithPlaceholder("type here")
	if same != pw {
		t.Fatal("WithPlaceholder should return the same *PasswordInput for chaining")
	}
	if pw.Placeholder() != "type here" {
		t.Fatalf("Placeholder() = %q, want %q", pw.Placeholder(), "type here")
	}
}

func TestPasswordInputWithPlaceholderStyle(t *testing.T) {
	pw := NewPasswordInput().WithPlaceholderStyle(NewStyle().WithForeground(ColorRed))
	if pw.PlaceholderStyle() == nil {
		t.Fatal("PlaceholderStyle() should not be nil after WithPlaceholderStyle")
	}

	pw.WithoutPlaceholderStyle()
	if pw.PlaceholderStyle() != nil {
		t.Fatal("WithoutPlaceholderStyle should clear the placeholder style")
	}
}

func TestPasswordInputWithPrefix(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.WithPrefix(">> ")
	if same != pw {
		t.Fatal("WithPrefix should return the same *PasswordInput for chaining")
	}
	if pw.Prefix() != ">> " {
		t.Fatalf("Prefix() = %q, want %q", pw.Prefix(), ">> ")
	}
}

func TestPasswordInputWithPrefixStyle(t *testing.T) {
	pw := NewPasswordInput().WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	if pw.PrefixStyle() == nil {
		t.Fatal("PrefixStyle() should not be nil after WithPrefixStyle")
	}

	pw.WithoutPrefixStyle()
	if pw.PrefixStyle() != nil {
		t.Fatal("WithoutPrefixStyle should clear the prefix style")
	}
}

func TestPasswordInputWithSuffix(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.WithSuffix(" <<")
	if same != pw {
		t.Fatal("WithSuffix should return the same *PasswordInput for chaining")
	}
	if pw.Suffix() != " <<" {
		t.Fatalf("Suffix() = %q, want %q", pw.Suffix(), " <<")
	}
}

func TestPasswordInputWithSuffixStyle(t *testing.T) {
	pw := NewPasswordInput().WithSuffixStyle(NewStyle().WithForeground(ColorRed))
	if pw.SuffixStyle() == nil {
		t.Fatal("SuffixStyle() should not be nil after WithSuffixStyle")
	}

	pw.WithoutSuffixStyle()
	if pw.SuffixStyle() != nil {
		t.Fatal("WithoutSuffixStyle should clear the suffix style")
	}
}

func TestPasswordInputPreferredWidthIncludesPrefixAndSuffix(t *testing.T) {
	pw := NewPasswordInput().WithWidth(10).WithPrefix(">> ").WithSuffix(" <<")
	if w := pw.PreferredWidth(); w != 16 {
		t.Fatalf("PreferredWidth() = %d, want 16", w)
	}
}

func TestPasswordInputWithWidth(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.WithWidth(10)
	if same != pw {
		t.Fatal("WithWidth should return the same *PasswordInput for chaining")
	}
	if pw.Width() != 10 {
		t.Fatalf("Width() = %d, want 10", pw.Width())
	}
	if pw.PreferredWidth() != 10 {
		t.Fatalf("PreferredWidth() = %d, want 10", pw.PreferredWidth())
	}
}

func TestPasswordInputPreferredHeight(t *testing.T) {
	pw := NewPasswordInput()
	if h := pw.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}

func TestPasswordInputFocusBlur(t *testing.T) {
	pw := NewPasswordInput()
	pw.Focus()
	if !pw.Focused() {
		t.Fatal("Focused() should be true after Focus")
	}

	pw.Blur()
	if pw.Focused() {
		t.Fatal("Focused() should be false after Blur")
	}
}

func TestPasswordInputAsFocused(t *testing.T) {
	pw := NewPasswordInput()
	same := pw.AsFocused(true)
	if same != pw {
		t.Fatal("AsFocused should return the same *PasswordInput for chaining")
	}
	if !pw.Focused() {
		t.Fatal("Focused() should be true after AsFocused(true)")
	}

	pw.AsFocused(false)
	if pw.Focused() {
		t.Fatal("Focused() should be false after AsFocused(false)")
	}
}

func TestPasswordInputIgnoresKeysWhenBlurred(t *testing.T) {
	pw := NewPasswordInput()
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	if pw.Value() != "" {
		t.Fatalf("Value() = %q, want empty (blurred input should ignore keys)", pw.Value())
	}
}

func TestPasswordInputInsertsRuneWhenFocused(t *testing.T) {
	pw := NewPasswordInput().AsFocused(true)
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'b'})
	if pw.Value() != "ab" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "ab")
	}
}

func TestPasswordInputInsertAtCursorPosition(t *testing.T) {
	pw := NewPasswordInput().WithValue("ac").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyLeft})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'b'})
	if pw.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "abc")
	}
}

func TestPasswordInputIgnoresCtrlAndAltRunes(t *testing.T) {
	pw := NewPasswordInput().AsFocused(true)
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'c', Ctrl: true})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'x', Alt: true})
	if pw.Value() != "" {
		t.Fatalf("Value() = %q, want empty", pw.Value())
	}
}

func TestPasswordInputBackspace(t *testing.T) {
	pw := NewPasswordInput().WithValue("abc").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyBackspace})
	if pw.Value() != "ab" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "ab")
	}
}

func TestPasswordInputBackspaceAtStartIsNoop(t *testing.T) {
	pw := NewPasswordInput().AsFocused(true)
	pw.Update(KeyEvent{Type: KeyBackspace})
	if pw.Value() != "" {
		t.Fatalf("Value() = %q, want empty", pw.Value())
	}
}

func TestPasswordInputDelete(t *testing.T) {
	pw := NewPasswordInput().WithValue("abc").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyHome})
	pw.Update(KeyEvent{Type: KeyDelete})
	if pw.Value() != "bc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "bc")
	}
}

func TestPasswordInputDeleteAtEndIsNoop(t *testing.T) {
	pw := NewPasswordInput().WithValue("abc").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyDelete})
	if pw.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "abc")
	}
}

func TestPasswordInputHomeThenInsert(t *testing.T) {
	pw := NewPasswordInput().WithValue("bc").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyHome})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	if pw.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "abc")
	}
}

func TestPasswordInputEndThenInsert(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyHome})
	pw.Update(KeyEvent{Type: KeyEnd})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'c'})
	if pw.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "abc")
	}
}

func TestPasswordInputLeftAtStartIsNoop(t *testing.T) {
	pw := NewPasswordInput().WithValue("bc").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyLeft})
	pw.Update(KeyEvent{Type: KeyLeft})
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	if pw.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "abc")
	}
}

func TestPasswordInputRightAtEndIsNoop(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").AsFocused(true)
	for range 5 {
		pw.Update(KeyEvent{Type: KeyRight})
	}
	pw.Update(KeyEvent{Type: KeyRune, Rune: 'c'})
	if pw.Value() != "abc" {
		t.Fatalf("Value() = %q, want %q", pw.Value(), "abc")
	}
}

func TestPasswordInputUpdateReturnsEventUnchanged(t *testing.T) {
	pw := NewPasswordInput().AsFocused(true)
	e := KeyEvent{Type: KeyRune, Rune: 'a'}
	if got := pw.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestPasswordInputRenderMasksValue(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab")
	lines := pw.Render(5, 1)
	want := DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + "   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderWithCustomMaskChar(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").WithMaskChar("*")
	lines := pw.Render(5, 1)
	want := "**   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderShrinksToConfiguredWidth(t *testing.T) {
	// A configured width smaller than the render width caps the visible
	// window, padding the remainder with blanks up to the full render width.
	pw := NewPasswordInput().WithValue("abcdef").WithWidth(3).AsFocused(true)
	pw.Update(KeyEvent{Type: KeyHome})
	pw.Blur()
	lines := pw.Render(10, 1)
	want := DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + "       "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderDefaultWidthFillsRenderWidth(t *testing.T) {
	// An unconfigured input's width (DefaultPasswordInputWidth) is effectively
	// unbounded, so it fills whatever width it's rendered with.
	pw := NewPasswordInput().WithValue("abcdef").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyHome})
	pw.Blur()
	lines := pw.Render(4, 1)
	want := strings.Repeat(DefaultPasswordInputMaskChar, 4)
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderZeroSize(t *testing.T) {
	pw := NewPasswordInput()
	if lines := pw.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := pw.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestPasswordInputRenderHidesCursorWhenBlurred(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab")
	lines := pw.Render(5, 1)
	want := DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + "   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderShowsCursorAtEndWhenFocusedAndEmpty(t *testing.T) {
	pw := NewPasswordInput().AsFocused(true)
	lines := pw.Render(3, 1)
	want := DefaultCursorChar + "  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderCursorOverlaysMaskedCharacterAtPosition(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyLeft})
	lines := pw.Render(5, 1)
	want := DefaultPasswordInputMaskChar + DefaultCursorChar + "   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderRevealsMaskedCharacterWhenCursorBlinkedOff(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").AsFocused(true)
	pw.Update(KeyEvent{Type: KeyLeft})

	base := time.Unix(0, 0)
	pw.Cursor().WithBlinkSpeed(10 * time.Millisecond)
	pw.Cursor().Tick(base)
	pw.Cursor().Tick(base.Add(15 * time.Millisecond))

	lines := pw.Render(5, 1)
	want := DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + "   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputCursorCustomChar(t *testing.T) {
	pw := NewPasswordInput().AsFocused(true)
	pw.Cursor().WithChar("_")
	lines := pw.Render(3, 1)
	want := "_  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderScrollsRightToKeepCursorVisible(t *testing.T) {
	pw := NewPasswordInput().WithValue("abcdef").WithWidth(3).AsFocused(true)
	lines := pw.Render(3, 1)
	want := DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + DefaultCursorChar
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderScrollsLeftAfterMovingCursorBack(t *testing.T) {
	pw := NewPasswordInput().WithValue("abcdef").WithWidth(3).AsFocused(true)
	pw.Render(3, 1)
	for range 6 {
		pw.Update(KeyEvent{Type: KeyLeft})
	}
	lines := pw.Render(3, 1)
	want := DefaultCursorChar + DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	pw := NewPasswordInput().WithValue("ab").WithStyle(NewStyle().WithForeground(ColorRed))
	lines := pw.Render(2, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == DefaultPasswordInputMaskChar+DefaultPasswordInputMaskChar {
		t.Fatalf("expected styled text to contain an SGR sequence, got %q", lines[0])
	}
}

func TestPasswordInputRenderWithPrefixAndSuffix(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").WithPrefix("[").WithSuffix("]")
	lines := pw.Render(6, 1)
	want := "[" + DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + "  ]"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderClipsPrefixAndSuffixToWidth(t *testing.T) {
	// The field shrinks to make room for the prefix and suffix, and the
	// suffix is dropped entirely once the prefix alone fills the width.
	pw := NewPasswordInput().WithValue("x").WithPrefix("abcde").WithSuffix("fghij")
	lines := pw.Render(3, 1)
	want := "abc"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderWithPrefixAndSuffixStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	pw := NewPasswordInput().
		WithValue("ab").
		WithPrefix("[").WithPrefixStyle(NewStyle().WithForeground(ColorRed)).
		WithSuffix("]").WithSuffixStyle(NewStyle().WithForeground(ColorBlue))
	lines := pw.Render(6, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "["+DefaultPasswordInputMaskChar+DefaultPasswordInputMaskChar+"  ]" {
		t.Fatalf("expected styled prefix/suffix to contain SGR sequences, got %q", lines[0])
	}
}

func TestPasswordInputRenderShowsPlaceholderWhenEmpty(t *testing.T) {
	pw := NewPasswordInput().WithPlaceholder("search")
	lines := pw.Render(10, 1)
	want := "search    "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderHidesPlaceholderWhenValueSet(t *testing.T) {
	pw := NewPasswordInput().WithValue("ab").WithPlaceholder("search")
	lines := pw.Render(10, 1)
	want := DefaultPasswordInputMaskChar + DefaultPasswordInputMaskChar + "        "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderPlaceholderCursorAtStartWhenFocused(t *testing.T) {
	pw := NewPasswordInput().WithPlaceholder("abc").AsFocused(true)
	lines := pw.Render(5, 1)
	want := DefaultCursorChar + "bc  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestPasswordInputRenderWithPlaceholderStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	pw := NewPasswordInput().WithPlaceholder("ab").WithPlaceholderStyle(NewStyle().WithForeground(ColorRed))
	lines := pw.Render(2, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "ab" {
		t.Fatalf("expected styled placeholder to contain an SGR sequence, got %q", lines[0])
	}
}
