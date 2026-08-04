package term

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// DefaultIntCounterSpeed is the time an IntCounter takes to animate from its
// value at the moment WithValue is called to the new target, used until
// WithSpeed is called.
const DefaultIntCounterSpeed = 500 * time.Millisecond

// DefaultIntCounterFormat is the function used to render an IntCounter's
// value until WithFormat is called: the value as a plain decimal integer.
var DefaultIntCounterFormat = func(v int) string { return strconv.Itoa(v) }

// IntCounter renders an integer value that eases toward a target value over
// time instead of jumping to it immediately. WithValue sets a new target and
// lets it animate there, taking Speed to arrive; Jump sets a value
// immediately with no animation. It does not animate on its own; call Tick
// whenever a new frame should be shown (e.g. from a program-driven timer,
// see Program.WithTick).
type IntCounter struct {
	value       int
	target      int
	startValue  int
	speed       time.Duration
	elapsed     time.Duration
	lastTime    time.Time
	format      func(int) string
	prefix      string
	prefixStyle *Style
	suffix      string
	suffixStyle *Style
	style       *Style
}

func NewIntCounter(initial int) *IntCounter {
	return &IntCounter{
		value:      initial,
		target:     initial,
		startValue: initial,
		speed:      DefaultIntCounterSpeed,
		format:     DefaultIntCounterFormat,
	}
}

// Value returns the counter's current, possibly mid-animation, value.
func (c *IntCounter) Value() int {
	return c.value
}

// WithValue sets the value the counter animates toward, easing from its
// current value over the duration set by WithSpeed. Call Tick to advance the
// animation. Use Jump to change the value immediately instead.
func (c *IntCounter) WithValue(v int) *IntCounter {
	if v == c.target {
		return c
	}
	c.startValue = c.value
	c.target = v
	c.elapsed = 0
	c.lastTime = time.Time{}
	return c
}

// Target returns the value the counter is currently animating toward.
func (c *IntCounter) Target() int {
	return c.target
}

// Jump immediately sets both the current and target value, skipping the
// animation.
func (c *IntCounter) Jump(v int) *IntCounter {
	c.value = v
	c.target = v
	c.startValue = v
	c.elapsed = 0
	c.lastTime = time.Time{}
	return c
}

// Animating reports whether the counter's value has not yet reached its
// target.
func (c *IntCounter) Animating() bool {
	return c.value != c.target
}

func (c *IntCounter) Speed() time.Duration {
	return c.speed
}

// WithSpeed sets how long the counter takes to animate from its value at the
// time of the change to its current target, restarting the animation from
// the current value.
func (c *IntCounter) WithSpeed(d time.Duration) *IntCounter {
	c.speed = d
	c.startValue = c.value
	c.elapsed = 0
	c.lastTime = time.Time{}
	return c
}

func (c *IntCounter) Format() func(int) string {
	return c.format
}

// WithFormat sets the function used to render the counter's current value,
// e.g. to add thousands separators, a percent sign, or fixed-point scaling.
func (c *IntCounter) WithFormat(f func(int) string) *IntCounter {
	c.format = f
	return c
}

func (c *IntCounter) Prefix() string {
	return c.prefix
}

// WithPrefix sets the text rendered before the counter's value.
func (c *IntCounter) WithPrefix(s string) *IntCounter {
	c.prefix = s
	return c
}

func (c *IntCounter) PrefixStyle() *Style {
	return c.prefixStyle
}

// WithPrefixStyle sets the style applied to the counter's prefix when
// rendering.
func (c *IntCounter) WithPrefixStyle(s Style) *IntCounter {
	c.prefixStyle = &s
	return c
}

// WithoutPrefixStyle removes the counter's prefix style, rendering it plain.
func (c *IntCounter) WithoutPrefixStyle() *IntCounter {
	c.prefixStyle = nil
	return c
}

func (c *IntCounter) Suffix() string {
	return c.suffix
}

// WithSuffix sets the text rendered after the counter's value.
func (c *IntCounter) WithSuffix(s string) *IntCounter {
	c.suffix = s
	return c
}

func (c *IntCounter) SuffixStyle() *Style {
	return c.suffixStyle
}

// WithSuffixStyle sets the style applied to the counter's suffix when
// rendering.
func (c *IntCounter) WithSuffixStyle(s Style) *IntCounter {
	c.suffixStyle = &s
	return c
}

// WithoutSuffixStyle removes the counter's suffix style, rendering it plain.
func (c *IntCounter) WithoutSuffixStyle() *IntCounter {
	c.suffixStyle = nil
	return c
}

func (c *IntCounter) Style() *Style {
	return c.style
}

// WithStyle sets the style applied to the counter's value when rendering.
func (c *IntCounter) WithStyle(s Style) *IntCounter {
	c.style = &s
	return c
}

// WithoutStyle removes the counter's value style, rendering it plain.
func (c *IntCounter) WithoutStyle() *IntCounter {
	c.style = nil
	return c
}

// Tick advances the counter's value toward its target based on the time
// elapsed since the previous call, given the current time (typically
// TickEvent.Time). It has no effect once the value reaches the target, or
// while the speed is zero or negative.
func (c *IntCounter) Tick(now time.Time) *IntCounter {
	if c.value == c.target || c.speed <= 0 {
		return c
	}
	if c.lastTime.IsZero() {
		c.lastTime = now
		return c
	}

	c.elapsed += now.Sub(c.lastTime)
	c.lastTime = now

	if c.elapsed >= c.speed {
		c.value = c.target
		return c
	}

	frac := float64(c.elapsed) / float64(c.speed)
	c.value = c.startValue + int(math.Round(float64(c.target-c.startValue)*frac))
	return c
}

// PreferredWidth returns the widest of the counter's current and target
// values (its formatted text's two endpoints) plus its prefix and suffix,
// so its size doesn't jitter as it animates.
func (c *IntCounter) PreferredWidth() int {
	w := max(StringWidth(c.format(c.value)), StringWidth(c.format(c.target)))
	return StringWidth(c.prefix) + w + StringWidth(c.suffix)
}

func (c *IntCounter) PreferredHeight(width int) int {
	return 1
}

// Update advances the counter's animation on every TickEvent (see
// Program.WithTick).
func (c *IntCounter) Update(e Event) Event {
	if t, ok := e.(TickEvent); ok {
		c.Tick(t.Time)
	}
	return e
}

// Render draws the counter's prefix, current value, and suffix — each styled
// independently — clipped to width columns.
func (c *IntCounter) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	var b strings.Builder
	remaining := width
	remaining = renderCounterSegment(&b, c.prefix, c.prefixStyle, remaining)
	remaining = renderCounterSegment(&b, c.format(c.value), c.style, remaining)
	renderCounterSegment(&b, c.suffix, c.suffixStyle, remaining)

	return []string{b.String()}
}

// renderCounterSegment writes s onto b, styled with style if non-nil,
// clipped to at most remaining columns, and returns the columns left after
// it. Shared by IntCounter, FloatCounter, Confirm, and Select.
func renderCounterSegment(b *strings.Builder, s string, style *Style, remaining int) int {
	if remaining <= 0 || s == "" {
		return remaining
	}

	runes := []rune(s)
	n, cols := splitWidth(runes, remaining)
	text := string(runes[:n])
	if style != nil {
		text = style.Render(text)
	}
	b.WriteString(text)
	return remaining - cols
}
