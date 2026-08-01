package term

import (
	"testing"
	"time"
)

var cursorTestBase = time.Unix(0, 0)

func TestNewCursor(t *testing.T) {
	c := NewCursor()
	if c.Char() != DefaultCursorChar {
		t.Fatalf("Char() = %q, want %q", c.Char(), DefaultCursorChar)
	}
	if !c.Blinking() {
		t.Fatal("new cursor should blink by default")
	}
	if !c.Visible() {
		t.Fatal("new cursor should be visible by default")
	}
	if !c.Showing() {
		t.Fatal("new cursor should be showing by default")
	}
	if c.BlinkSpeed() != DefaultCursorBlinkSpeed {
		t.Fatalf("BlinkSpeed() = %v, want %v", c.BlinkSpeed(), DefaultCursorBlinkSpeed)
	}
}

func TestCursorWithChar(t *testing.T) {
	c := NewCursor()
	same := c.WithChar("<>")
	if same != c {
		t.Fatal("WithChar should return the same *Cursor for chaining")
	}
	if c.Char() != "<>" {
		t.Fatalf("Char() = %q, want %q", c.Char(), "<>")
	}
}

func TestCursorWithoutStyle(t *testing.T) {
	c := NewCursor().WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if c.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestCursorAsBlinkingFalseKeepsShowing(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(10 * time.Millisecond)
	c.Tick(cursorTestBase)
	c.Tick(cursorTestBase.Add(15 * time.Millisecond))
	if c.Showing() {
		t.Fatal("cursor should be blinked off before disabling blinking")
	}

	same := c.AsBlinking(false)
	if same != c {
		t.Fatal("AsBlinking should return the same *Cursor for chaining")
	}
	if !c.Showing() {
		t.Fatal("disabling blinking should show the cursor")
	}

	c.Tick(cursorTestBase.Add(time.Hour))
	if !c.Showing() {
		t.Fatal("cursor should keep showing while blinking is disabled")
	}
}

func TestCursorAsBlinkingTrueResetsCycle(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(10 * time.Millisecond).AsBlinking(false)
	c.AsBlinking(true)
	if !c.Showing() {
		t.Fatal("re-enabling blinking should restart from showing")
	}
	c.Tick(cursorTestBase.Add(time.Hour))
	if !c.Showing() {
		t.Fatal("first tick after re-enabling should only set the baseline, not toggle")
	}
}

func TestCursorWithBlinkSpeed(t *testing.T) {
	c := NewCursor()
	same := c.WithBlinkSpeed(50 * time.Millisecond)
	if same != c {
		t.Fatal("WithBlinkSpeed should return the same *Cursor for chaining")
	}
	if c.BlinkSpeed() != 50*time.Millisecond {
		t.Fatalf("BlinkSpeed() = %v, want %v", c.BlinkSpeed(), 50*time.Millisecond)
	}
}

func TestCursorAsVisible(t *testing.T) {
	c := NewCursor()
	same := c.AsVisible(false)
	if same != c {
		t.Fatal("AsVisible should return the same *Cursor for chaining")
	}
	if c.Visible() {
		t.Fatal("Visible() should be false after AsVisible(false)")
	}
	if c.Showing() {
		t.Fatal("Showing() should be false when not visible")
	}
}

func TestCursorTickTogglesAfterBlinkSpeedElapses(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(100 * time.Millisecond)

	c.Tick(cursorTestBase)
	if !c.Showing() {
		t.Fatal("first Tick should only set the baseline, not toggle")
	}

	c.Tick(cursorTestBase.Add(150 * time.Millisecond))
	if c.Showing() {
		t.Fatal("cursor should be blinked off after one blink speed has elapsed")
	}

	c.Tick(cursorTestBase.Add(210 * time.Millisecond))
	if !c.Showing() {
		t.Fatal("cursor should be blinked back on after a second blink speed has elapsed")
	}
}

func TestCursorTickNoEffectWhenNotBlinking(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(10 * time.Millisecond).AsBlinking(false)
	c.Tick(cursorTestBase)
	c.Tick(cursorTestBase.Add(time.Hour))
	if !c.Showing() {
		t.Fatal("Tick should have no effect while blinking is disabled")
	}
}

func TestCursorUpdateAdvancesOnTickEvent(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(10 * time.Millisecond)
	c.Update(TickEvent{Time: cursorTestBase})
	c.Update(TickEvent{Time: cursorTestBase.Add(15 * time.Millisecond)})
	if c.Showing() {
		t.Fatal("Update should advance the blink state on TickEvent")
	}
}

func TestCursorUpdateIgnoresOtherEvents(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(10 * time.Millisecond)
	c.Update(KeyEvent{Rune: 'x'})
	if !c.Showing() {
		t.Fatal("Update should ignore non-TickEvent events")
	}
}

func TestCursorUpdateReturnsEventUnchanged(t *testing.T) {
	c := NewCursor()
	e := TickEvent{Time: cursorTestBase}
	if got := c.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestCursorRender(t *testing.T) {
	c := NewCursor()
	lines := c.Render(1, 1)
	if len(lines) != 1 || lines[0] != DefaultCursorChar {
		t.Fatalf("got %#v, want [%q]", lines, DefaultCursorChar)
	}
}

func TestCursorRenderHiddenWhenNotVisible(t *testing.T) {
	c := NewCursor().AsVisible(false)
	lines := c.Render(1, 1)
	if len(lines) != 1 || lines[0] != " " {
		t.Fatalf("got %#v, want [\" \"]", lines)
	}
}

func TestCursorRenderBlankDuringBlinkOff(t *testing.T) {
	c := NewCursor().WithBlinkSpeed(10 * time.Millisecond)
	c.Tick(cursorTestBase)
	c.Tick(cursorTestBase.Add(15 * time.Millisecond))
	lines := c.Render(1, 1)
	if len(lines) != 1 || lines[0] != " " {
		t.Fatalf("got %#v, want [\" \"]", lines)
	}
}

func TestCursorRenderWithChar(t *testing.T) {
	c := NewCursor().WithChar("<>")
	lines := c.Render(10, 1)
	if len(lines) != 1 || lines[0] != "<>" {
		t.Fatalf("got %#v, want [\"<>\"]", lines)
	}
}

func TestCursorRenderClipsToWidth(t *testing.T) {
	c := NewCursor().WithChar("abc")
	lines := c.Render(2, 1)
	if len(lines) != 1 || lines[0] != "ab" {
		t.Fatalf("got %#v, want [\"ab\"]", lines)
	}
}

func TestCursorRenderZeroSize(t *testing.T) {
	c := NewCursor()
	if lines := c.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := c.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestCursorRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	c := NewCursor().WithStyle(NewStyle().WithForeground(ColorRed))
	lines := c.Render(1, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == DefaultCursorChar {
		t.Fatalf("expected styled cursor to contain an SGR sequence, got %q", lines[0])
	}
}

func TestCursorPreferredWidth(t *testing.T) {
	c := NewCursor().WithChar("abc")
	if w := c.PreferredWidth(); w != 3 {
		t.Fatalf("PreferredWidth() = %d, want 3", w)
	}
}

func TestCursorPreferredHeight(t *testing.T) {
	c := NewCursor()
	if h := c.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}
