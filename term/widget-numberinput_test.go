package term

import (
	"testing"
	"time"
)

func TestNewNumberInput(t *testing.T) {
	ni := NewNumberInput()
	if v, err := ni.Value(); err == nil {
		t.Fatalf("Value() = %v, want error for empty input, got value %v", err, v)
	}
	if ni.ForceValue() != 0 {
		t.Fatalf("ForceValue() = %v, want 0", ni.ForceValue())
	}
	if ni.Focused() {
		t.Fatal("new number input should not be focused")
	}
	if ni.Width() != DefaultNumberInputWidth {
		t.Fatalf("Width() = %d, want %d", ni.Width(), DefaultNumberInputWidth)
	}
	if ni.Cursor() == nil {
		t.Fatal("Cursor() should not be nil")
	}
}

func TestNumberInputWithValue(t *testing.T) {
	ni := NewNumberInput()
	same := ni.WithValue(42)
	if same != ni {
		t.Fatal("WithValue should return the same *NumberInput for chaining")
	}
	if v, err := ni.Value(); err != nil || v != 42 {
		t.Fatalf("Value() = (%v, %v), want (42, nil)", v, err)
	}
}

func TestNumberInputWithValuePtrReadsInitial(t *testing.T) {
	v := 3.5
	ni := NewNumberInput().WithValuePtr(&v)
	if got, err := ni.Value(); err != nil || got != 3.5 {
		t.Fatalf("Value() = (%v, %v), want (3.5, nil)", got, err)
	}
	if ni.ValuePtr() != &v {
		t.Fatal("ValuePtr() should return the bound pointer")
	}
}

func TestNumberInputWithValuePtrSyncsOnEdit(t *testing.T) {
	var v float64
	ni := NewNumberInput().WithValuePtr(&v).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyRune, Rune: '1'})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '2'})
	if v != 12 {
		t.Fatalf("bound value = %v, want 12", v)
	}
	if ni.ForceValue() != 12 {
		t.Fatalf("ForceValue() = %v, want 12", ni.ForceValue())
	}
}

func TestNumberInputWithValuePtrDoesNotSyncOnInvalidState(t *testing.T) {
	v := 5.0
	ni := NewNumberInput().WithValuePtr(&v).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyBackspace})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '-'})
	if v != 5.0 {
		t.Fatalf("bound value = %v, want unchanged 5.0 while text is mid-edit", v)
	}
	if _, err := ni.Value(); err == nil {
		t.Fatal("Value() should error while text is just \"-\"")
	}
}

func TestNumberInputWithoutStyle(t *testing.T) {
	ni := NewNumberInput().WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if ni.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestNumberInputWithPlaceholder(t *testing.T) {
	ni := NewNumberInput()
	same := ni.WithPlaceholder("0")
	if same != ni {
		t.Fatal("WithPlaceholder should return the same *NumberInput for chaining")
	}
	if ni.Placeholder() != "0" {
		t.Fatalf("Placeholder() = %q, want %q", ni.Placeholder(), "0")
	}
}

func TestNumberInputWithPlaceholderStyle(t *testing.T) {
	ni := NewNumberInput().WithPlaceholderStyle(NewStyle().WithForeground(ColorRed))
	if ni.PlaceholderStyle() == nil {
		t.Fatal("PlaceholderStyle() should not be nil after WithPlaceholderStyle")
	}

	ni.WithoutPlaceholderStyle()
	if ni.PlaceholderStyle() != nil {
		t.Fatal("WithoutPlaceholderStyle should clear the placeholder style")
	}
}

func TestNumberInputWithPrefix(t *testing.T) {
	ni := NewNumberInput()
	same := ni.WithPrefix("$")
	if same != ni {
		t.Fatal("WithPrefix should return the same *NumberInput for chaining")
	}
	if ni.Prefix() != "$" {
		t.Fatalf("Prefix() = %q, want %q", ni.Prefix(), "$")
	}
}

func TestNumberInputWithPrefixStyle(t *testing.T) {
	ni := NewNumberInput().WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	if ni.PrefixStyle() == nil {
		t.Fatal("PrefixStyle() should not be nil after WithPrefixStyle")
	}

	ni.WithoutPrefixStyle()
	if ni.PrefixStyle() != nil {
		t.Fatal("WithoutPrefixStyle should clear the prefix style")
	}
}

func TestNumberInputWithSuffix(t *testing.T) {
	ni := NewNumberInput()
	same := ni.WithSuffix("%")
	if same != ni {
		t.Fatal("WithSuffix should return the same *NumberInput for chaining")
	}
	if ni.Suffix() != "%" {
		t.Fatalf("Suffix() = %q, want %q", ni.Suffix(), "%")
	}
}

func TestNumberInputWithSuffixStyle(t *testing.T) {
	ni := NewNumberInput().WithSuffixStyle(NewStyle().WithForeground(ColorRed))
	if ni.SuffixStyle() == nil {
		t.Fatal("SuffixStyle() should not be nil after WithSuffixStyle")
	}

	ni.WithoutSuffixStyle()
	if ni.SuffixStyle() != nil {
		t.Fatal("WithoutSuffixStyle should clear the suffix style")
	}
}

func TestNumberInputPreferredWidthIncludesPrefixAndSuffix(t *testing.T) {
	ni := NewNumberInput().WithWidth(10).WithPrefix("$").WithSuffix("%")
	if w := ni.PreferredWidth(); w != 12 {
		t.Fatalf("PreferredWidth() = %d, want 12", w)
	}
}

func TestNumberInputWithWidth(t *testing.T) {
	ni := NewNumberInput()
	same := ni.WithWidth(10)
	if same != ni {
		t.Fatal("WithWidth should return the same *NumberInput for chaining")
	}
	if ni.Width() != 10 {
		t.Fatalf("Width() = %d, want 10", ni.Width())
	}
	if ni.PreferredWidth() != 10 {
		t.Fatalf("PreferredWidth() = %d, want 10", ni.PreferredWidth())
	}
}

func TestNumberInputPreferredHeight(t *testing.T) {
	ni := NewNumberInput()
	if h := ni.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}

func TestNumberInputFocusBlur(t *testing.T) {
	ni := NewNumberInput()
	ni.Focus()
	if !ni.Focused() {
		t.Fatal("Focused() should be true after Focus")
	}

	ni.Blur()
	if ni.Focused() {
		t.Fatal("Focused() should be false after Blur")
	}
}

func TestNumberInputAsFocused(t *testing.T) {
	ni := NewNumberInput()
	same := ni.AsFocused(true)
	if same != ni {
		t.Fatal("AsFocused should return the same *NumberInput for chaining")
	}
	if !ni.Focused() {
		t.Fatal("Focused() should be true after AsFocused(true)")
	}

	ni.AsFocused(false)
	if ni.Focused() {
		t.Fatal("Focused() should be false after AsFocused(false)")
	}
}

func TestNumberInputIgnoresKeysWhenBlurred(t *testing.T) {
	ni := NewNumberInput()
	ni.Update(KeyEvent{Type: KeyRune, Rune: '1'})
	if ni.ForceValue() != 0 {
		t.Fatalf("ForceValue() = %v, want 0 (blurred input should ignore keys)", ni.ForceValue())
	}
}

func TestNumberInputInsertsDigitsWhenFocused(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	ni.Update(KeyEvent{Type: KeyRune, Rune: '1'})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '2'})
	if v, err := ni.Value(); err != nil || v != 12 {
		t.Fatalf("Value() = (%v, %v), want (12, nil)", v, err)
	}
}

func TestNumberInputInsertsDecimalPoint(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	for _, r := range "3.14" {
		ni.Update(KeyEvent{Type: KeyRune, Rune: r})
	}
	if v, err := ni.Value(); err != nil || v != 3.14 {
		t.Fatalf("Value() = (%v, %v), want (3.14, nil)", v, err)
	}
}

func TestNumberInputRejectsSecondDecimalPoint(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	for _, r := range "1.2.3" {
		ni.Update(KeyEvent{Type: KeyRune, Rune: r})
	}
	if v, err := ni.Value(); err != nil || v != 1.23 {
		t.Fatalf("Value() = (%v, %v), want (1.23, nil)", v, err)
	}
}

func TestNumberInputInsertsLeadingMinus(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	for _, r := range "-5" {
		ni.Update(KeyEvent{Type: KeyRune, Rune: r})
	}
	if v, err := ni.Value(); err != nil || v != -5 {
		t.Fatalf("Value() = (%v, %v), want (-5, nil)", v, err)
	}
}

func TestNumberInputRejectsMinusNotAtStart(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	ni.Update(KeyEvent{Type: KeyRune, Rune: '5'})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '-'})
	if v, err := ni.Value(); err != nil || v != 5 {
		t.Fatalf("Value() = (%v, %v), want (5, nil)", v, err)
	}
}

func TestNumberInputRejectsSecondMinus(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	ni.Update(KeyEvent{Type: KeyRune, Rune: '-'})
	ni.Update(KeyEvent{Type: KeyHome})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '-'})
	ni.Update(KeyEvent{Type: KeyEnd})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '5'})
	if v, err := ni.Value(); err != nil || v != -5 {
		t.Fatalf("Value() = (%v, %v), want (-5, nil)", v, err)
	}
}

func TestNumberInputRejectsNonNumericRunes(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	ni.Update(KeyEvent{Type: KeyRune, Rune: 'a'})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '5'})
	if v, err := ni.Value(); err != nil || v != 5 {
		t.Fatalf("Value() = (%v, %v), want (5, nil)", v, err)
	}
}

func TestNumberInputIgnoresCtrlAndAltRunes(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	ni.Update(KeyEvent{Type: KeyRune, Rune: '1', Ctrl: true})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '2', Alt: true})
	if _, err := ni.Value(); err == nil {
		t.Fatal("Value() should error, want empty input")
	}
}

func TestNumberInputBackspace(t *testing.T) {
	ni := NewNumberInput().WithValue(123).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyBackspace})
	if v, err := ni.Value(); err != nil || v != 12 {
		t.Fatalf("Value() = (%v, %v), want (12, nil)", v, err)
	}
}

func TestNumberInputBackspaceAtStartIsNoop(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	ni.Update(KeyEvent{Type: KeyBackspace})
	if _, err := ni.Value(); err == nil {
		t.Fatal("Value() should error, want empty input")
	}
}

func TestNumberInputDelete(t *testing.T) {
	ni := NewNumberInput().WithValue(123).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyHome})
	ni.Update(KeyEvent{Type: KeyDelete})
	if v, err := ni.Value(); err != nil || v != 23 {
		t.Fatalf("Value() = (%v, %v), want (23, nil)", v, err)
	}
}

func TestNumberInputDeleteAtEndIsNoop(t *testing.T) {
	ni := NewNumberInput().WithValue(123).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyDelete})
	if v, err := ni.Value(); err != nil || v != 123 {
		t.Fatalf("Value() = (%v, %v), want (123, nil)", v, err)
	}
}

func TestNumberInputHomeThenInsert(t *testing.T) {
	ni := NewNumberInput().WithValue(23).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyHome})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '1'})
	if v, err := ni.Value(); err != nil || v != 123 {
		t.Fatalf("Value() = (%v, %v), want (123, nil)", v, err)
	}
}

func TestNumberInputEndThenInsert(t *testing.T) {
	ni := NewNumberInput().WithValue(12).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyHome})
	ni.Update(KeyEvent{Type: KeyEnd})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '3'})
	if v, err := ni.Value(); err != nil || v != 123 {
		t.Fatalf("Value() = (%v, %v), want (123, nil)", v, err)
	}
}

func TestNumberInputLeftAtStartIsNoop(t *testing.T) {
	ni := NewNumberInput().WithValue(23).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyLeft})
	ni.Update(KeyEvent{Type: KeyLeft})
	ni.Update(KeyEvent{Type: KeyRune, Rune: '1'})
	if v, err := ni.Value(); err != nil || v != 123 {
		t.Fatalf("Value() = (%v, %v), want (123, nil)", v, err)
	}
}

func TestNumberInputRightAtEndIsNoop(t *testing.T) {
	ni := NewNumberInput().WithValue(12).AsFocused(true)
	for range 5 {
		ni.Update(KeyEvent{Type: KeyRight})
	}
	ni.Update(KeyEvent{Type: KeyRune, Rune: '3'})
	if v, err := ni.Value(); err != nil || v != 123 {
		t.Fatalf("Value() = (%v, %v), want (123, nil)", v, err)
	}
}

func TestNumberInputUpdateReturnsEventUnchanged(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	e := KeyEvent{Type: KeyRune, Rune: '1'}
	if got := ni.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestNumberInputRenderPadsToWidth(t *testing.T) {
	ni := NewNumberInput().WithValue(12)
	lines := ni.Render(5, 1)
	if len(lines) != 1 || lines[0] != "12   " {
		t.Fatalf("got %#v, want [\"12   \"]", lines)
	}
}

func TestNumberInputRenderShrinksToConfiguredWidth(t *testing.T) {
	ni := NewNumberInput().WithValue(123456).WithWidth(3).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyHome})
	ni.Blur()
	lines := ni.Render(10, 1)
	if len(lines) != 1 || lines[0] != "123       " {
		t.Fatalf("got %#v, want [\"123       \"]", lines)
	}
}

func TestNumberInputRenderZeroSize(t *testing.T) {
	ni := NewNumberInput()
	if lines := ni.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := ni.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestNumberInputRenderHidesCursorWhenBlurred(t *testing.T) {
	ni := NewNumberInput().WithValue(12)
	lines := ni.Render(5, 1)
	if len(lines) != 1 || lines[0] != "12   " {
		t.Fatalf("got %#v, want [\"12   \"]", lines)
	}
}

func TestNumberInputRenderShowsCursorAtEndWhenFocusedAndEmpty(t *testing.T) {
	ni := NewNumberInput().AsFocused(true)
	lines := ni.Render(3, 1)
	want := DefaultCursorChar + "  "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestNumberInputRenderCursorOverlaysCharacterAtPosition(t *testing.T) {
	ni := NewNumberInput().WithValue(12).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyLeft})
	lines := ni.Render(5, 1)
	want := "1" + DefaultCursorChar + "   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestNumberInputRenderRevealsCharacterWhenCursorBlinkedOff(t *testing.T) {
	ni := NewNumberInput().WithValue(12).AsFocused(true)
	ni.Update(KeyEvent{Type: KeyLeft})

	base := time.Unix(0, 0)
	ni.Cursor().WithBlinkSpeed(10 * time.Millisecond)
	ni.Cursor().Tick(base)
	ni.Cursor().Tick(base.Add(15 * time.Millisecond))

	lines := ni.Render(5, 1)
	want := "12   "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestNumberInputRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	ni := NewNumberInput().WithValue(12).WithStyle(NewStyle().WithForeground(ColorRed))
	lines := ni.Render(2, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "12" {
		t.Fatalf("expected styled text to contain an SGR sequence, got %q", lines[0])
	}
}

func TestNumberInputRenderWithPrefixAndSuffix(t *testing.T) {
	ni := NewNumberInput().WithValue(12).WithPrefix("[").WithSuffix("]")
	lines := ni.Render(6, 1)
	want := "[12  ]"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestNumberInputRenderShowsPlaceholderWhenEmpty(t *testing.T) {
	ni := NewNumberInput().WithPlaceholder("0")
	lines := ni.Render(10, 1)
	want := "0         "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}

func TestNumberInputRenderHidesPlaceholderWhenValueSet(t *testing.T) {
	ni := NewNumberInput().WithValue(12).WithPlaceholder("0")
	lines := ni.Render(10, 1)
	want := "12        "
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("got %#v, want [%q]", lines, want)
	}
}
