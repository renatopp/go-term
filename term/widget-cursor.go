package term

import (
	"strings"
	"time"

	"github.com/renatopp/go-term/term/ui"
)

// DefaultCursorChar is the character used by a Cursor until WithChar is
// called.
const DefaultCursorChar = "█"

// DefaultCursorBlinkSpeed is the interval between blink toggles used by a
// Cursor until WithBlinkSpeed is called.
const DefaultCursorBlinkSpeed = 530 * time.Millisecond

// Cursor renders a single character that can blink on and off, useful for
// marking the insertion point of a text input. Its blink state advances by
// calling Tick, typically forwarded from TickEvent (see Program.WithTick).
type Cursor struct {
	char       string
	style      *Style
	blinking   bool
	blinkSpeed time.Duration
	visible    bool
	on         bool
	elapsed    time.Duration
	lastTime   time.Time
}

func NewCursor() *Cursor {
	return &Cursor{
		char:       DefaultCursorChar,
		blinking:   true,
		blinkSpeed: DefaultCursorBlinkSpeed,
		visible:    true,
		on:         true,
	}
}

func (c *Cursor) Char() string {
	return c.char
}

// WithChar sets the character (or string) rendered by the cursor.
func (c *Cursor) WithChar(char string) *Cursor {
	c.char = char
	return c
}

func (c *Cursor) Style() *Style {
	return c.style
}

// WithStyle sets the style applied to the cursor when rendering.
func (c *Cursor) WithStyle(s Style) *Cursor {
	c.style = &s
	return c
}

// WithoutStyle removes the cursor's style, rendering plain text.
func (c *Cursor) WithoutStyle() *Cursor {
	c.style = nil
	return c
}

func (c *Cursor) Blinking() bool {
	return c.blinking
}

// AsBlinking toggles whether the cursor blinks over time. Disabling it keeps
// the cursor solidly shown (subject to Visible); enabling it restarts the
// blink cycle from fully shown.
func (c *Cursor) AsBlinking(v bool) *Cursor {
	c.blinking = v
	c.on = true
	c.elapsed = 0
	c.lastTime = time.Time{}
	return c
}

func (c *Cursor) BlinkSpeed() time.Duration {
	return c.blinkSpeed
}

// WithBlinkSpeed sets the interval between blink toggles.
func (c *Cursor) WithBlinkSpeed(d time.Duration) *Cursor {
	c.blinkSpeed = d
	c.elapsed = 0
	return c
}

func (c *Cursor) Visible() bool {
	return c.visible
}

// AsVisible toggles whether the cursor is shown at all, independent of
// blinking.
func (c *Cursor) AsVisible(v bool) *Cursor {
	c.visible = v
	return c
}

// Showing reports whether the cursor is currently rendering its character,
// accounting for both Visible and, while Blinking, the current blink phase.
func (c *Cursor) Showing() bool {
	return c.visible && (!c.blinking || c.on)
}

// Tick advances the cursor's blink state based on the time elapsed since the
// previous call, given the current time (typically TickEvent.Time). It has
// no effect while blinking is disabled.
func (c *Cursor) Tick(now time.Time) *Cursor {
	if !c.blinking || c.blinkSpeed <= 0 {
		return c
	}
	if c.lastTime.IsZero() {
		c.lastTime = now
		return c
	}

	c.elapsed += now.Sub(c.lastTime)
	c.lastTime = now

	if toggles := int(c.elapsed / c.blinkSpeed); toggles > 0 {
		c.elapsed -= time.Duration(toggles) * c.blinkSpeed
		if toggles%2 == 1 {
			c.on = !c.on
		}
	}
	return c
}

func (c *Cursor) PreferredWidth() int {
	return ui.StringWidth(c.char)
}

func (c *Cursor) PreferredHeight(width int) int {
	return 1
}

// Update advances the cursor's blink state on every TickEvent (see
// Program.WithTick).
func (c *Cursor) Update(e Event) Event {
	if t, ok := e.(TickEvent); ok {
		c.Tick(t.Time)
	}
	return e
}

// Render draws the cursor's character when showing, or blank space of the
// same width otherwise, clipped to width columns.
func (c *Cursor) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	text := c.char
	if !c.Showing() {
		text = strings.Repeat(" ", ui.StringWidth(c.char))
	}

	runes := []rune(text)
	n, _ := splitWidth(runes, width)
	line := string(runes[:n])

	if c.style != nil {
		line = c.style.Render(line)
	}
	return []string{line}
}
