package term

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

var intCounterTestBase = time.Unix(0, 0)

func TestNewIntCounter(t *testing.T) {
	c := NewIntCounter(5)
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}
	if c.Target() != 5 {
		t.Fatalf("Target() = %d, want 5", c.Target())
	}
	if c.Speed() != DefaultIntCounterSpeed {
		t.Fatalf("Speed() = %v, want %v", c.Speed(), DefaultIntCounterSpeed)
	}
	if c.Animating() {
		t.Fatal("new counter should not be animating")
	}
}

func TestIntCounterWithValue(t *testing.T) {
	c := NewIntCounter(0)
	same := c.WithValue(10)
	if same != c {
		t.Fatal("WithValue should return the same *IntCounter for chaining")
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

func TestIntCounterWithValueSameTargetDoesNotRestartAnimation(t *testing.T) {
	c := NewIntCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Tick(intCounterTestBase)
	c.Tick(intCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}

	c.WithValue(10) // same target again, should be a no-op
	c.Tick(intCounterTestBase.Add(100 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %d, want 10 (animation should not have restarted)", c.Value())
	}
}

func TestIntCounterJump(t *testing.T) {
	c := NewIntCounter(0).WithValue(10)
	same := c.Jump(7)
	if same != c {
		t.Fatal("Jump should return the same *IntCounter for chaining")
	}
	if c.Value() != 7 || c.Target() != 7 {
		t.Fatalf("Value()/Target() = %d/%d, want 7/7", c.Value(), c.Target())
	}
	if c.Animating() {
		t.Fatal("counter should not be animating right after Jump")
	}
}

func TestIntCounterWithSpeed(t *testing.T) {
	c := NewIntCounter(0)
	same := c.WithSpeed(50 * time.Millisecond)
	if same != c {
		t.Fatal("WithSpeed should return the same *IntCounter for chaining")
	}
	if c.Speed() != 50*time.Millisecond {
		t.Fatalf("Speed() = %v, want %v", c.Speed(), 50*time.Millisecond)
	}
}

func TestIntCounterWithSpeedRestartsFromCurrentValue(t *testing.T) {
	c := NewIntCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Tick(intCounterTestBase)
	c.Tick(intCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}

	c.WithSpeed(20 * time.Millisecond)
	c.Tick(intCounterTestBase.Add(60 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5 (this tick only re-baselines timing)", c.Value())
	}
	c.Tick(intCounterTestBase.Add(80 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %d, want 10 (new speed elapsed from value 5 to target 10)", c.Value())
	}
}

func TestIntCounterWithFormat(t *testing.T) {
	c := NewIntCounter(3).WithFormat(func(v int) string { return "#" + strconv.Itoa(v) })
	lines := c.Render(10, 1)
	if len(lines) != 1 || lines[0] != "#3" {
		t.Fatalf("got %#v, want [\"#3\"]", lines)
	}
}

func TestIntCounterWithPrefixAndSuffix(t *testing.T) {
	c := NewIntCounter(3).WithPrefix(">> ").WithSuffix(" <<")
	same := c.WithPrefix(">> ")
	if same != c {
		t.Fatal("WithPrefix should return the same *IntCounter for chaining")
	}
	if c.Prefix() != ">> " || c.Suffix() != " <<" {
		t.Fatalf("Prefix()/Suffix() = %q/%q", c.Prefix(), c.Suffix())
	}
	lines := c.Render(20, 1)
	if len(lines) != 1 || lines[0] != ">> 3 <<" {
		t.Fatalf("got %#v, want [\">> 3 <<\"]", lines)
	}
}

func TestIntCounterWithoutStyle(t *testing.T) {
	c := NewIntCounter(0).WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if c.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestIntCounterWithPrefixStyle(t *testing.T) {
	c := NewIntCounter(0)
	same := c.WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	if same != c {
		t.Fatal("WithPrefixStyle should return the same *IntCounter for chaining")
	}
	if c.PrefixStyle() == nil {
		t.Fatal("PrefixStyle() should not be nil after WithPrefixStyle")
	}
}

func TestIntCounterWithoutPrefixStyle(t *testing.T) {
	c := NewIntCounter(0).WithPrefixStyle(NewStyle().WithForeground(ColorRed)).WithoutPrefixStyle()
	if c.PrefixStyle() != nil {
		t.Fatal("WithoutPrefixStyle should clear the prefix style")
	}
}

func TestIntCounterWithSuffixStyle(t *testing.T) {
	c := NewIntCounter(0)
	same := c.WithSuffixStyle(NewStyle().WithForeground(ColorBlue))
	if same != c {
		t.Fatal("WithSuffixStyle should return the same *IntCounter for chaining")
	}
	if c.SuffixStyle() == nil {
		t.Fatal("SuffixStyle() should not be nil after WithSuffixStyle")
	}
}

func TestIntCounterWithoutSuffixStyle(t *testing.T) {
	c := NewIntCounter(0).WithSuffixStyle(NewStyle().WithForeground(ColorBlue)).WithoutSuffixStyle()
	if c.SuffixStyle() != nil {
		t.Fatal("WithoutSuffixStyle should clear the suffix style")
	}
}

func TestIntCounterRenderStylesPrefixAndSuffixIndependently(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	c := NewIntCounter(3).
		WithPrefix(">> ").WithPrefixStyle(NewStyle().WithForeground(ColorRed)).
		WithSuffix(" <<").WithSuffixStyle(NewStyle().WithForeground(ColorBlue))
	lines := c.Render(20, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	line := lines[0]
	if !strings.Contains(line, "3") {
		t.Fatalf("expected value to still be present, got %q", line)
	}
	if strings.Count(line, "\x1b[") < 2 {
		t.Fatalf("expected separate SGR sequences for prefix and suffix, got %q", line)
	}
}

func TestIntCounterRenderClipsAcrossSegments(t *testing.T) {
	c := NewIntCounter(123).WithPrefix("abc").WithSuffix("xyz")
	lines := c.Render(5, 1)
	if len(lines) != 1 || lines[0] != "abc12" {
		t.Fatalf("got %#v, want [\"abc12\"]", lines)
	}
}

func TestIntCounterRenderNoOverflowWhenPrefixFillsWidth(t *testing.T) {
	c := NewIntCounter(999).WithPrefix("abcde").WithSuffix("!")
	lines := c.Render(5, 1)
	if len(lines) != 1 || lines[0] != "abcde" {
		t.Fatalf("got %#v, want [\"abcde\"]", lines)
	}
}

func TestIntCounterTickInterpolatesTowardTarget(t *testing.T) {
	c := NewIntCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)

	c.Tick(intCounterTestBase)
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0 (first tick only sets baseline)", c.Value())
	}

	c.Tick(intCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5 (halfway through the animation)", c.Value())
	}

	c.Tick(intCounterTestBase.Add(100 * time.Millisecond))
	if c.Value() != 10 {
		t.Fatalf("Value() = %d, want 10 (animation complete)", c.Value())
	}
	if c.Animating() {
		t.Fatal("counter should stop animating once it reaches the target")
	}
}

func TestIntCounterTickCountsDown(t *testing.T) {
	c := NewIntCounter(10).WithSpeed(100 * time.Millisecond)
	c.WithValue(0)
	c.Tick(intCounterTestBase)
	c.Tick(intCounterTestBase.Add(50 * time.Millisecond))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}
}

func TestIntCounterTickSnapsToTargetOnceSpeedElapses(t *testing.T) {
	c := NewIntCounter(0).WithSpeed(10 * time.Millisecond)
	c.WithValue(3)
	c.Tick(intCounterTestBase)
	c.Tick(intCounterTestBase.Add(time.Hour))
	if c.Value() != 3 {
		t.Fatalf("Value() = %d, want 3 (clamped to target)", c.Value())
	}
}

func TestIntCounterTickNoEffectAtTarget(t *testing.T) {
	c := NewIntCounter(5).WithSpeed(10 * time.Millisecond)
	c.Tick(intCounterTestBase)
	c.Tick(intCounterTestBase.Add(time.Hour))
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5 (already at target)", c.Value())
	}
}

func TestIntCounterTickNoEffectWhenSpeedZero(t *testing.T) {
	c := NewIntCounter(0).WithValue(5).WithSpeed(0)
	c.Tick(intCounterTestBase)
	c.Tick(intCounterTestBase.Add(time.Hour))
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0 (speed disabled)", c.Value())
	}
}

func TestIntCounterUpdateAdvancesOnTickEvent(t *testing.T) {
	c := NewIntCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Update(TickEvent{Time: intCounterTestBase})
	c.Update(TickEvent{Time: intCounterTestBase.Add(50 * time.Millisecond)})
	if c.Value() != 5 {
		t.Fatalf("Value() = %d, want 5", c.Value())
	}
}

func TestIntCounterUpdateIgnoresOtherEvents(t *testing.T) {
	c := NewIntCounter(0).WithSpeed(100 * time.Millisecond)
	c.WithValue(10)
	c.Update(KeyEvent{Rune: 'x'})
	if c.Value() != 0 {
		t.Fatalf("Value() = %d, want 0", c.Value())
	}
}

func TestIntCounterUpdateReturnsEventUnchanged(t *testing.T) {
	c := NewIntCounter(0)
	e := TickEvent{Time: intCounterTestBase}
	if got := c.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestIntCounterRender(t *testing.T) {
	c := NewIntCounter(42)
	lines := c.Render(10, 1)
	if len(lines) != 1 || lines[0] != "42" {
		t.Fatalf("got %#v, want [\"42\"]", lines)
	}
}

func TestIntCounterRenderClipsToWidth(t *testing.T) {
	c := NewIntCounter(12345)
	lines := c.Render(3, 1)
	if len(lines) != 1 || lines[0] != "123" {
		t.Fatalf("got %#v, want [\"123\"]", lines)
	}
}

func TestIntCounterRenderZeroSize(t *testing.T) {
	c := NewIntCounter(0)
	if lines := c.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := c.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestIntCounterRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	c := NewIntCounter(1).WithStyle(NewStyle().WithForeground(ColorRed))
	lines := c.Render(5, 1)
	if len(lines) != 1 {
		t.Fatalf("got %#v, want 1 line", lines)
	}
	if lines[0] == "1" {
		t.Fatalf("expected styled counter to contain an SGR sequence, got %q", lines[0])
	}
}

func TestIntCounterPreferredWidth(t *testing.T) {
	c := NewIntCounter(5).WithValue(100).WithPrefix(">").WithSuffix("!")
	// widest of "5" (1) and "100" (3) = 3, plus prefix (1) and suffix (1) = 5
	if w := c.PreferredWidth(); w != 5 {
		t.Fatalf("PreferredWidth() = %d, want 5", w)
	}
}

func TestIntCounterPreferredHeight(t *testing.T) {
	c := NewIntCounter(0)
	if h := c.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}
