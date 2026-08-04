package term

import (
	"strings"
	"testing"
	"time"
)

var floatCounterTestBase = time.Unix(0, 0)

func TestNewFloatCounter(t *testing.T) {
	c := NewFloatCounter(5.5)
	if c.Value() != 5.5 {
		t.Fatalf("Value() = %v, want 5.5", c.Value())
	}
	if c.Target() != 5.5 {
		t.Fatalf("Target() = %v, want 5.5", c.Target())
	}
	if c.Speed() != DefaultFloatCounterSpeed {
		t.Fatalf("Speed() = %v, want %v", c.Speed(), DefaultFloatCounterSpeed)
	}
	if c.Animating() {
		t.Fatal("new counter should not be animating")
	}
}

func TestFloatCounterWithValue(t *testing.T) {
	c := NewFloatCounter(0)
	same := c.WithValue(10)
	if same != c {
		t.Fatal("WithValue should return the same *FloatCounter for chaining")
	}
	if c.Target() != 10 {
		t.Fatalf("Target() = %v, want 10", c.Target())
	}
	if c.Value() != 0 {
		t.Fatalf("Value() = %v, want 0 (should not jump)", c.Value())
	}
	if !c.Animating() {
		t.Fatal("counter should be animating after WithValue changes the target")
	}
}

func TestFloatCounterWithValueSameTargetDoesNotRestartAnimation(t *testing.T) {
	c := NewFloatCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Tick(floatCounterTestBase)
	c.Tick(floatCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %v, want 5", c.Value())
	}

	c.WithValue(10) // same target again, should be a no-op
	c.Tick(floatCounterTestBase.Add(100 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %v, want 10 (animation should not have restarted)", c.Value())
	}
}

func TestFloatCounterJump(t *testing.T) {
	c := NewFloatCounter(0).WithValue(10)
	same := c.Jump(7.5)
	if same != c {
		t.Fatal("Jump should return the same *FloatCounter for chaining")
	}
	if c.Value() != 7.5 || c.Target() != 7.5 {
		t.Fatalf("Value()/Target() = %v/%v, want 7.5/7.5", c.Value(), c.Target())
	}
	if c.Animating() {
		t.Fatal("counter should not be animating right after Jump")
	}
}

func TestFloatCounterWithSpeed(t *testing.T) {
	c := NewFloatCounter(0)
	same := c.WithSpeed(50 * time.Millisecond)
	if same != c {
		t.Fatal("WithSpeed should return the same *FloatCounter for chaining")
	}
	if c.Speed() != 50*time.Millisecond {
		t.Fatalf("Speed() = %v, want %v", c.Speed(), 50*time.Millisecond)
	}
}

func TestFloatCounterWithSpeedRestartsFromCurrentValue(t *testing.T) {
	c := NewFloatCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Tick(floatCounterTestBase)
	c.Tick(floatCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %v, want 5", c.Value())
	}

	c.WithSpeed(20 * time.Millisecond)
	c.Tick(floatCounterTestBase.Add(60 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %v, want 5 (this tick only re-baselines timing)", c.Value())
	}
	c.Tick(floatCounterTestBase.Add(80 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %v, want 10 (new speed elapsed from value 5 to target 10)", c.Value())
	}
}

func TestFloatCounterWithFormat(t *testing.T) {
	c := NewFloatCounter(3).WithFormat(func(v float64) string { return "#" + DefaultFloatCounterFormat(v) })
	lines := c.Render(20, 1)
	if len(lines) != 1 || lines[0] != "#3.00" {
		t.Fatalf("got %#v, want [\"#3.00\"]", lines)
	}
}

func TestFloatCounterWithPrefixAndSuffix(t *testing.T) {
	c := NewFloatCounter(3).WithPrefix(">> ").WithSuffix(" <<")
	same := c.WithPrefix(">> ")
	if same != c {
		t.Fatal("WithPrefix should return the same *FloatCounter for chaining")
	}
	if c.Prefix() != ">> " || c.Suffix() != " <<" {
		t.Fatalf("Prefix()/Suffix() = %q/%q", c.Prefix(), c.Suffix())
	}
	lines := c.Render(20, 1)
	if len(lines) != 1 || lines[0] != ">> 3.00 <<" {
		t.Fatalf("got %#v, want [\">> 3.00 <<\"]", lines)
	}
}

func TestFloatCounterWithoutStyle(t *testing.T) {
	c := NewFloatCounter(0).WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if c.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestFloatCounterWithPrefixStyle(t *testing.T) {
	c := NewFloatCounter(0)
	same := c.WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	if same != c {
		t.Fatal("WithPrefixStyle should return the same *FloatCounter for chaining")
	}
	if c.PrefixStyle() == nil {
		t.Fatal("PrefixStyle() should not be nil after WithPrefixStyle")
	}
}

func TestFloatCounterWithoutPrefixStyle(t *testing.T) {
	c := NewFloatCounter(0).WithPrefixStyle(NewStyle().WithForeground(ColorRed)).WithoutPrefixStyle()
	if c.PrefixStyle() != nil {
		t.Fatal("WithoutPrefixStyle should clear the prefix style")
	}
}

func TestFloatCounterWithSuffixStyle(t *testing.T) {
	c := NewFloatCounter(0)
	same := c.WithSuffixStyle(NewStyle().WithForeground(ColorBlue))
	if same != c {
		t.Fatal("WithSuffixStyle should return the same *FloatCounter for chaining")
	}
	if c.SuffixStyle() == nil {
		t.Fatal("SuffixStyle() should not be nil after WithSuffixStyle")
	}
}

func TestFloatCounterWithoutSuffixStyle(t *testing.T) {
	c := NewFloatCounter(0).WithSuffixStyle(NewStyle().WithForeground(ColorBlue)).WithoutSuffixStyle()
	if c.SuffixStyle() != nil {
		t.Fatal("WithoutSuffixStyle should clear the suffix style")
	}
}

func TestFloatCounterRenderStylesPrefixAndSuffixIndependently(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	c := NewFloatCounter(3).
		WithPrefix(">> ").WithPrefixStyle(NewStyle().WithForeground(ColorRed)).
		WithSuffix(" <<").WithSuffixStyle(NewStyle().WithForeground(ColorBlue))
	lines := c.Render(20, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	line := lines[0]
	if !strings.Contains(line, "3.00") {
		t.Fatalf("expected value to still be present, got %q", line)
	}
	if strings.Count(line, "\x1b[") < 2 {
		t.Fatalf("expected separate SGR sequences for prefix and suffix, got %q", line)
	}
}

func TestFloatCounterRenderClipsAcrossSegments(t *testing.T) {
	c := NewFloatCounter(1.5).WithPrefix("abc").WithSuffix("xyz")
	lines := c.Render(6, 1)
	if len(lines) != 1 || lines[0] != "abc1.5" {
		t.Fatalf("got %#v, want [\"abc1.5\"]", lines)
	}
}

func TestFloatCounterRenderNoOverflowWhenPrefixFillsWidth(t *testing.T) {
	c := NewFloatCounter(9.99).WithPrefix("abcde").WithSuffix("!")
	lines := c.Render(5, 1)
	if len(lines) != 1 || lines[0] != "abcde" {
		t.Fatalf("got %#v, want [\"abcde\"]", lines)
	}
}

func TestFloatCounterTickInterpolatesTowardTarget(t *testing.T) {
	c := NewFloatCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)

	c.Tick(floatCounterTestBase)
	if c.Value() != 0 {
		t.Fatalf("Value() = %v, want 0 (first tick only sets baseline)", c.Value())
	}

	c.Tick(floatCounterTestBase.Add(25 * time.Millisecond))
	if c.Value() != 2.5 {
		t.Fatalf("Value() = %v, want 2.5 (a quarter through the animation)", c.Value())
	}

	c.Tick(floatCounterTestBase.Add(100 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %v, want 10 (animation complete)", c.Value())
	}
	if c.Animating() {
		t.Fatal("counter should stop animating once it reaches the target")
	}
}

func TestFloatCounterTickCountsDown(t *testing.T) {
	c := NewFloatCounter(10).WithSpeed(100 * time.Millisecond)
	c.WithValue(0)
	c.Tick(floatCounterTestBase)
	c.Tick(floatCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %v, want 5", c.Value())
	}
}

func TestFloatCounterTickSnapsToTargetOnceSpeedElapses(t *testing.T) {
	c := NewFloatCounter(0).WithSpeed(10 * time.Millisecond)
	c.WithValue(3.25)
	c.Tick(floatCounterTestBase)
	c.Tick(floatCounterTestBase.Add(time.Hour))
	if c.Value() != 3.25 {
		t.Fatalf("Value() = %v, want 3.25 (clamped to target)", c.Value())
	}
}

func TestFloatCounterTickNoEffectAtTarget(t *testing.T) {
	c := NewFloatCounter(5).WithSpeed(10 * time.Millisecond)
	c.Tick(floatCounterTestBase)
	c.Tick(floatCounterTestBase.Add(time.Hour))
	if c.Value() != 5 {
		t.Fatalf("Value() = %v, want 5 (already at target)", c.Value())
	}
}

func TestFloatCounterTickNoEffectWhenSpeedZero(t *testing.T) {
	c := NewFloatCounter(0).WithValue(5).WithSpeed(0)
	c.Tick(floatCounterTestBase)
	c.Tick(floatCounterTestBase.Add(time.Hour))
	if c.Value() != 0 {
		t.Fatalf("Value() = %v, want 0 (speed disabled)", c.Value())
	}
}

func TestFloatCounterUpdateAdvancesOnTickEvent(t *testing.T) {
	c := NewFloatCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Update(TickEvent{Time: floatCounterTestBase})
	c.Update(TickEvent{Time: floatCounterTestBase.Add(50 * time.Millisecond)})
	if c.Value() != 5 {
		t.Fatalf("Value() = %v, want 5", c.Value())
	}
}

func TestFloatCounterUpdateIgnoresOtherEvents(t *testing.T) {
	c := NewFloatCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Update(KeyEvent{Rune: 'x'})
	if c.Value() != 0 {
		t.Fatalf("Value() = %v, want 0", c.Value())
	}
}

func TestFloatCounterUpdateReturnsEventUnchanged(t *testing.T) {
	c := NewFloatCounter(0)
	e := TickEvent{Time: floatCounterTestBase}
	if got := c.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestFloatCounterRenderDefaultFormat(t *testing.T) {
	c := NewFloatCounter(3.14159)
	lines := c.Render(10, 1)
	if len(lines) != 1 || lines[0] != "3.14" {
		t.Fatalf("got %#v, want [\"3.14\"]", lines)
	}
}

func TestFloatCounterRenderClipsToWidth(t *testing.T) {
	c := NewFloatCounter(12345.6)
	lines := c.Render(3, 1)
	if len(lines) != 1 || lines[0] != "123" {
		t.Fatalf("got %#v, want [\"123\"]", lines)
	}
}

func TestFloatCounterRenderZeroSize(t *testing.T) {
	c := NewFloatCounter(0)
	if lines := c.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := c.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestFloatCounterRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	c := NewFloatCounter(1).WithStyle(NewStyle().WithForeground(ColorRed))
	lines := c.Render(6, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "1.00" {
		t.Fatalf("expected styled counter to contain an SGR sequence, got %q", lines[0])
	}
}

func TestFloatCounterPreferredWidth(t *testing.T) {
	c := NewFloatCounter(5).WithValue(100).WithPrefix(">").WithSuffix("!")
	// widest of "5.00" (4) and "100.00" (6) = 6, plus prefix (1) and suffix (1) = 8
	if w := c.PreferredWidth(); w != 8 {
		t.Fatalf("PreferredWidth() = %d, want 8", w)
	}
}

func TestFloatCounterPreferredHeight(t *testing.T) {
	c := NewFloatCounter(0)
	if h := c.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}
