package term

import (
	"math"
	"strconv"
	"time"

	"github.com/renatopp/go-term/term/ui"
)

// DefaultCounterSpeed is the time a Counter takes to animate from its value
// at the moment WithValue is called to the new target, used until WithSpeed
// is called.
const DefaultCounterSpeed = 500 * time.Millisecond

// DefaultCounterFormat is the function used to render a Counter's value
// until WithFormat is called: the value as a plain decimal integer.
var DefaultCounterFormat = func(v int) string { return strconv.Itoa(v) }

// Counter renders an integer value that eases toward a target value over
// time instead of jumping to it immediately. WithValue sets a new target and
// lets it animate there, taking Speed to arrive; Jump sets a value
// immediately with no animation. It does not animate on its own; call Tick
// whenever a new frame should be shown (e.g. from a program-driven timer,
// see Program.WithTick).
type Counter struct {
	value      int
	target     int
	startValue int
	speed      time.Duration
	elapsed    time.Duration
	lastTime   time.Time
	format     func(int) string
	prefix     string
	suffix     string
	style      *Style
}

func NewCounter(initial int) *Counter {
	return &Counter{
		value:      initial,
		target:     initial,
		startValue: initial,
		speed:      DefaultCounterSpeed,
		format:     DefaultCounterFormat,
	}
}

// Value returns the counter's current, possibly mid-animation, value.
func (c *Counter) Value() int {
	return c.value
}

// WithValue sets the value the counter animates toward, easing from its
// current value over the duration set by WithSpeed. Call Tick to advance the
// animation. Use Jump to change the value immediately instead.
func (c *Counter) WithValue(v int) *Counter {
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
func (c *Counter) Target() int {
	return c.target
}

// Jump immediately sets both the current and target value, skipping the
// animation.
func (c *Counter) Jump(v int) *Counter {
	c.value = v
	c.target = v
	c.startValue = v
	c.elapsed = 0
	c.lastTime = time.Time{}
	return c
}

// Animating reports whether the counter's value has not yet reached its
// target.
func (c *Counter) Animating() bool {
	return c.value != c.target
}

func (c *Counter) Speed() time.Duration {
	return c.speed
}

// WithSpeed sets how long the counter takes to animate from its value at the
// time of the change to its current target, restarting the animation from
// the current value.
func (c *Counter) WithSpeed(d time.Duration) *Counter {
	c.speed = d
	c.startValue = c.value
	c.elapsed = 0
	c.lastTime = time.Time{}
	return c
}

func (c *Counter) Format() func(int) string {
	return c.format
}

// WithFormat sets the function used to render the counter's current value,
// e.g. to add thousands separators, a percent sign, or fixed-point scaling.
func (c *Counter) WithFormat(f func(int) string) *Counter {
	c.format = f
	return c
}

func (c *Counter) Prefix() string {
	return c.prefix
}

// WithPrefix sets the text rendered before the counter's value.
func (c *Counter) WithPrefix(s string) *Counter {
	c.prefix = s
	return c
}

func (c *Counter) Suffix() string {
	return c.suffix
}

// WithSuffix sets the text rendered after the counter's value.
func (c *Counter) WithSuffix(s string) *Counter {
	c.suffix = s
	return c
}

func (c *Counter) Style() *Style {
	return c.style
}

// WithStyle sets the style applied to the counter's text when rendering.
func (c *Counter) WithStyle(s Style) *Counter {
	c.style = &s
	return c
}

// WithoutStyle removes the counter's style, rendering plain text.
func (c *Counter) WithoutStyle() *Counter {
	c.style = nil
	return c
}

// Tick advances the counter's value toward its target based on the time
// elapsed since the previous call, given the current time (typically
// TickEvent.Time). It has no effect once the value reaches the target, or
// while the speed is zero or negative.
func (c *Counter) Tick(now time.Time) *Counter {
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
func (c *Counter) PreferredWidth() int {
	w := max(ui.StringWidth(c.format(c.value)), ui.StringWidth(c.format(c.target)))
	return ui.StringWidth(c.prefix) + w + ui.StringWidth(c.suffix)
}

func (c *Counter) PreferredHeight(width int) int {
	return 1
}

// Update advances the counter's animation on every TickEvent (see
// Program.WithTick).
func (c *Counter) Update(e Event) Event {
	if t, ok := e.(TickEvent); ok {
		c.Tick(t.Time)
	}
	return e
}

// Render draws the counter's prefix, current value, and suffix, clipped to
// width columns.
func (c *Counter) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	full := c.prefix + c.format(c.value) + c.suffix

	runes := []rune(full)
	n, _ := splitWidth(runes, width)
	line := string(runes[:n])

	if c.style != nil {
		line = c.style.Render(line)
	}
	return []string{line}
}
