package term

import (
	"strconv"
	"testing"
	"time"
)

var counterTestBase = time.Unix(0, 0)

func TestNewCounter(t *testing.T) {
	c := NewCounter(5)
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}
	if c.Target() != 5 {
		t.Fatalf("Target() = %d, want 5", c.Target())
	}
	if c.Speed() != DefaultCounterSpeed {
		t.Fatalf("Speed() = %v, want %v", c.Speed(), DefaultCounterSpeed)
	}
	if c.Animating() {
		t.Fatal("new counter should not be animating")
	}
}

func TestCounterWithValue(t *testing.T) {
	c := NewCounter(0)
	same := c.WithValue(10)
	if same != c {
		t.Fatal("WithValue should return the same *Counter for chaining")
	}
	if c.Target() != 10 {
		t.Fatalf("Target() = %d, want 10", c.Target())
	}
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0 (should not jump)", c.Value())
	}
	if !c.Animating() {
		t.Fatal("counter should be animating after WithValue changes the target")
	}
}

func TestCounterWithValueSameTargetDoesNotRestartAnimation(t *testing.T) {
	c := NewCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Tick(counterTestBase)
	c.Tick(counterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}

	c.WithValue(10) // same target again, should be a no-op
	c.Tick(counterTestBase.Add(100 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %d, want 10 (animation should not have restarted)", c.Value())
	}
}

func TestCounterJump(t *testing.T) {
	c := NewCounter(0).WithValue(10)
	same := c.Jump(7)
	if same != c {
		t.Fatal("Jump should return the same *Counter for chaining")
	}
	if c.Value() != 7 || c.Target() != 7 {
		t.Fatalf("Value()/Target() = %d/%d, want 7/7", c.Value(), c.Target())
	}
	if c.Animating() {
		t.Fatal("counter should not be animating right after Jump")
	}
}

func TestCounterWithSpeed(t *testing.T) {
	c := NewCounter(0)
	same := c.WithSpeed(50 * time.Millisecond)
	if same != c {
		t.Fatal("WithSpeed should return the same *Counter for chaining")
	}
	if c.Speed() != 50*time.Millisecond {
		t.Fatalf("Speed() = %v, want %v", c.Speed(), 50*time.Millisecond)
	}
}

func TestCounterWithSpeedRestartsFromCurrentValue(t *testing.T) {
	c := NewCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Tick(counterTestBase)
	c.Tick(counterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}

	c.WithSpeed(20 * time.Millisecond)
	c.Tick(counterTestBase.Add(60 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5 (this tick only re-baselines timing)", c.Value())
	}
	c.Tick(counterTestBase.Add(80 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %d, want 10 (new speed elapsed from value 5 to target 10)", c.Value())
	}
}

func TestCounterWithFormat(t *testing.T) {
	c := NewCounter(3).WithFormat(func(v int) string { return "#" + strconv.Itoa(v) })
	lines := c.Render(10, 1)
	if len(lines) != 1 || lines[0] != "#3" {
		t.Fatalf("got %#v, want [\"#3\"]", lines)
	}
}

func TestCounterWithPrefixAndSuffix(t *testing.T) {
	c := NewCounter(3).WithPrefix(">> ").WithSuffix(" <<")
	same := c.WithPrefix(">> ")
	if same != c {
		t.Fatal("WithPrefix should return the same *Counter for chaining")
	}
	if c.Prefix() != ">> " || c.Suffix() != " <<" {
		t.Fatalf("Prefix()/Suffix() = %q/%q", c.Prefix(), c.Suffix())
	}
	lines := c.Render(20, 1)
	if len(lines) != 1 || lines[0] != ">> 3 <<" {
		t.Fatalf("got %#v, want [\">> 3 <<\"]", lines)
	}
}

func TestCounterWithoutStyle(t *testing.T) {
	c := NewCounter(0).WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if c.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestCounterTickInterpolatesTowardTarget(t *testing.T) {
	c := NewCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)

	c.Tick(counterTestBase)
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0 (first tick only sets baseline)", c.Value())
	}

	c.Tick(counterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5 (halfway through the animation)", c.Value())
	}

	c.Tick(counterTestBase.Add(100 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %d, want 10 (animation complete)", c.Value())
	}
	if c.Animating() {
		t.Fatal("counter should stop animating once it reaches the target")
	}
}

func TestCounterTickCountsDown(t *testing.T) {
	c := NewCounter(10).WithSpeed(100 * time.Millisecond)
	c.WithValue(0)
	c.Tick(counterTestBase)
	c.Tick(counterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}
}

func TestCounterTickSnapsToTargetOnceSpeedElapses(t *testing.T) {
	c := NewCounter(0).WithSpeed(10 * time.Millisecond)
	c.WithValue(3)
	c.Tick(counterTestBase)
	c.Tick(counterTestBase.Add(time.Hour))
	if c.Value() != 3 {
		t.Fatalf("Value() = %d, want 3 (clamped to target)", c.Value())
	}
}

func TestCounterTickNoEffectAtTarget(t *testing.T) {
	c := NewCounter(5).WithSpeed(10 * time.Millisecond)
	c.Tick(counterTestBase)
	c.Tick(counterTestBase.Add(time.Hour))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5 (already at target)", c.Value())
	}
}

func TestCounterTickNoEffectWhenSpeedZero(t *testing.T) {
	c := NewCounter(0).WithValue(5).WithSpeed(0)
	c.Tick(counterTestBase)
	c.Tick(counterTestBase.Add(time.Hour))
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0 (speed disabled)", c.Value())
	}
}

func TestCounterUpdateAdvancesOnTickEvent(t *testing.T) {
	c := NewCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Update(TickEvent{Time: counterTestBase})
	c.Update(TickEvent{Time: counterTestBase.Add(50 * time.Millisecond)})
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}
}

func TestCounterUpdateIgnoresOtherEvents(t *testing.T) {
	c := NewCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Update(KeyEvent{Rune: 'x'})
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0", c.Value())
	}
}

func TestCounterUpdateReturnsEventUnchanged(t *testing.T) {
	c := NewCounter(0)
	e := TickEvent{Time: counterTestBase}
	if got := c.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestCounterRender(t *testing.T) {
	c := NewCounter(42)
	lines := c.Render(10, 1)
	if len(lines) != 1 || lines[0] != "42" {
		t.Fatalf("got %#v, want [\"42\"]", lines)
	}
}

func TestCounterRenderClipsToWidth(t *testing.T) {
	c := NewCounter(12345)
	lines := c.Render(3, 1)
	if len(lines) != 1 || lines[0] != "123" {
		t.Fatalf("got %#v, want [\"123\"]", lines)
	}
}

func TestCounterRenderZeroSize(t *testing.T) {
	c := NewCounter(0)
	if lines := c.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := c.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestCounterRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	c := NewCounter(1).WithStyle(NewStyle().WithForeground(ColorRed))
	lines := c.Render(5, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "1" {
		t.Fatalf("expected styled counter to contain an SGR sequence, got %q", lines[0])
	}
}

func TestCounterPreferredWidth(t *testing.T) {
	c := NewCounter(5).WithValue(100).WithPrefix(">").WithSuffix("!")
	// widest of "5" (1) and "100" (3) = 3, plus prefix (1) and suffix (1) = 5
	if w := c.PreferredWidth(); w != 5 {
		t.Fatalf("PreferredWidth() = %d, want 5", w)
	}
}

func TestCounterPreferredHeight(t *testing.T) {
	c := NewCounter(0)
	if h := c.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}
