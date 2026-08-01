package term

import "testing"

func TestNewConfirm(t *testing.T) {
	c := NewConfirm("Are you sure?")
	if c.Message() != "Are you sure?" {
		t.Fatalf("Message() = %q, want %q", c.Message(), "Are you sure?")
	}
	if c.YesLabel() != DefaultConfirmYesLabel {
		t.Fatalf("YesLabel() = %q, want %q", c.YesLabel(), DefaultConfirmYesLabel)
	}
	if c.NoLabel() != DefaultConfirmNoLabel {
		t.Fatalf("NoLabel() = %q, want %q", c.NoLabel(), DefaultConfirmNoLabel)
	}
	if !c.Value() {
		t.Fatal("new confirm should default to Yes selected")
	}
	if c.Focused() {
		t.Fatal("new confirm should not be focused")
	}
}

func TestConfirmWithMessage(t *testing.T) {
	c := NewConfirm("old")
	same := c.WithMessage("new")
	if same != c {
		t.Fatal("WithMessage should return the same *Confirm for chaining")
	}
	if c.Message() != "new" {
		t.Fatalf("Message() = %q, want %q", c.Message(), "new")
	}
}

func TestConfirmWithYesLabel(t *testing.T) {
	c := NewConfirm("").WithYesLabel("Sure")
	if c.YesLabel() != "Sure" {
		t.Fatalf("YesLabel() = %q, want %q", c.YesLabel(), "Sure")
	}
}

func TestConfirmWithNoLabel(t *testing.T) {
	c := NewConfirm("").WithNoLabel("Nah")
	if c.NoLabel() != "Nah" {
		t.Fatalf("NoLabel() = %q, want %q", c.NoLabel(), "Nah")
	}
}

func TestConfirmWithValue(t *testing.T) {
	c := NewConfirm("")
	same := c.WithValue(false)
	if same != c {
		t.Fatal("WithValue should return the same *Confirm for chaining")
	}
	if c.Value() {
		t.Fatal("Value() should be false after WithValue(false)")
	}
}

func TestConfirmWithValuePtrReadsInitial(t *testing.T) {
	v := false
	c := NewConfirm("").WithValuePtr(&v)
	if c.Value() {
		t.Fatal("Value() should reflect the bound pointer's initial value")
	}
	if c.ValuePtr() != &v {
		t.Fatal("ValuePtr() should return the bound pointer")
	}
}

func TestConfirmWithValuePtrSyncsOnToggle(t *testing.T) {
	v := true
	c := NewConfirm("").WithValuePtr(&v).AsFocused(true)
	c.Update(KeyEvent{Type: KeyTab})
	if v {
		t.Fatal("bound value should be false after toggling")
	}
	if c.Value() {
		t.Fatal("Value() should be false after toggling")
	}
}

func TestConfirmWithStyle(t *testing.T) {
	c := NewConfirm("").WithStyle(NewStyle().WithForeground(ColorRed))
	if c.Style() == nil {
		t.Fatal("Style() should not be nil after WithStyle")
	}
}

func TestConfirmWithoutStyle(t *testing.T) {
	c := NewConfirm("").WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if c.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestConfirmWithSelectedStyle(t *testing.T) {
	c := NewConfirm("").WithSelectedStyle(NewStyle().WithForeground(ColorRed))
	if c.SelectedStyle() == nil {
		t.Fatal("SelectedStyle() should not be nil after WithSelectedStyle")
	}
}

func TestConfirmWithoutSelectedStyle(t *testing.T) {
	c := NewConfirm("").WithSelectedStyle(NewStyle().WithForeground(ColorRed)).WithoutSelectedStyle()
	if c.SelectedStyle() != nil {
		t.Fatal("WithoutSelectedStyle should clear the selected style")
	}
}

func TestConfirmFocusBlurAsFocused(t *testing.T) {
	c := NewConfirm("")
	c.Focus()
	if !c.Focused() {
		t.Fatal("Focus() should set Focused() to true")
	}
	c.Blur()
	if c.Focused() {
		t.Fatal("Blur() should set Focused() to false")
	}
	if same := c.AsFocused(true); same != c {
		t.Fatal("AsFocused should return the same *Confirm for chaining")
	}
	if !c.Focused() {
		t.Fatal("AsFocused(true) should set Focused() to true")
	}
}

func TestConfirmUpdateIgnoresKeysWhenBlurred(t *testing.T) {
	c := NewConfirm("")
	c.Update(KeyEvent{Type: KeyTab})
	if !c.Value() {
		t.Fatal("Update should ignore key events while blurred")
	}
}

func TestConfirmUpdateTogglesWithArrowKeys(t *testing.T) {
	c := NewConfirm("").AsFocused(true)
	c.Update(KeyEvent{Type: KeyLeft})
	if c.Value() {
		t.Fatal("Left should toggle Value() to false")
	}
	c.Update(KeyEvent{Type: KeyRight})
	if !c.Value() {
		t.Fatal("Right should toggle Value() back to true")
	}
	c.Update(KeyEvent{Type: KeyTab})
	if c.Value() {
		t.Fatal("Tab should toggle Value() to false")
	}
}

func TestConfirmUpdateSetsWithYN(t *testing.T) {
	c := NewConfirm("").AsFocused(true)
	c.Update(KeyEvent{Type: KeyRune, Rune: 'n'})
	if c.Value() {
		t.Fatal("'n' should select No")
	}
	c.Update(KeyEvent{Type: KeyRune, Rune: 'Y'})
	if !c.Value() {
		t.Fatal("'Y' should select Yes")
	}
}

func TestConfirmRenderBracketsSelected(t *testing.T) {
	c := NewConfirm("Continue?")
	lines := c.Render(30, 1)
	want := "Continue? [Yes]  No"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("Render() = %#v, want [%q]", lines, want)
	}
}

func TestConfirmRenderBracketsNoWhenSelected(t *testing.T) {
	c := NewConfirm("Continue?").WithValue(false)
	lines := c.Render(30, 1)
	want := "Continue? Yes  [No]"
	if len(lines) != 1 || lines[0] != want {
		t.Fatalf("Render() = %#v, want [%q]", lines, want)
	}
}

func TestConfirmRenderZeroSize(t *testing.T) {
	c := NewConfirm("Continue?")
	if lines := c.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := c.Render(10, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestConfirmPreferredWidth(t *testing.T) {
	c := NewConfirm("Continue?")
	// "Continue?" (9) + " " (1) + "[Yes]" (5) + "  " (2) + "No" (2) = 19
	if w := c.PreferredWidth(); w != 19 {
		t.Fatalf("PreferredWidth() = %d, want 19", w)
	}
}

func TestConfirmPreferredHeight(t *testing.T) {
	c := NewConfirm("Continue?")
	if h := c.PreferredHeight(30); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}
