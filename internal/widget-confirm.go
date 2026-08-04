package term

import (
	"strings"
)

// DefaultConfirmYesLabel is the label rendered for the affirmative option
// until WithYesLabel is called.
const DefaultConfirmYesLabel = "Yes"

// DefaultConfirmNoLabel is the label rendered for the negative option until
// WithNoLabel is called.
const DefaultConfirmNoLabel = "No"

// Confirm renders a message followed by a Yes/No choice, toggled with
// Left/Right/Tab or the y/n keys. The selected option is bracketed so it
// reads even without color, and additionally styled with SelectedStyle. It
// only reacts to key events while Focused.
type Confirm struct {
	message       string
	yesLabel      string
	noLabel       string
	value         bool
	valuePtr      *bool
	style         *Style
	selectedStyle *Style
	focused       bool
}

func NewConfirm(message string) *Confirm {
	return &Confirm{
		message:  message,
		yesLabel: DefaultConfirmYesLabel,
		noLabel:  DefaultConfirmNoLabel,
		value:    true,
	}
}

func (c *Confirm) Message() string {
	return c.message
}

// WithMessage sets the text rendered before the Yes/No options.
func (c *Confirm) WithMessage(s string) *Confirm {
	c.message = s
	return c
}

func (c *Confirm) YesLabel() string {
	return c.yesLabel
}

// WithYesLabel sets the label rendered for the affirmative option.
func (c *Confirm) WithYesLabel(s string) *Confirm {
	c.yesLabel = s
	return c
}

func (c *Confirm) NoLabel() string {
	return c.noLabel
}

// WithNoLabel sets the label rendered for the negative option.
func (c *Confirm) WithNoLabel(s string) *Confirm {
	c.noLabel = s
	return c
}

func (c *Confirm) Value() bool {
	return c.value
}

// WithValue sets whether Yes (true) or No (false) is currently selected.
func (c *Confirm) WithValue(v bool) *Confirm {
	c.value = v
	c.sync()
	return c
}

func (c *Confirm) ValuePtr() *bool {
	return c.valuePtr
}

// WithValuePtr binds the confirm to p: its initial selection is read from
// *p, and every subsequent toggle writes the current selection back into *p.
func (c *Confirm) WithValuePtr(p *bool) *Confirm {
	c.valuePtr = p
	if p != nil {
		c.value = *p
	}
	return c
}

func (c *Confirm) Style() *Style {
	return c.style
}

// WithStyle sets the style applied to the message and the unselected option.
func (c *Confirm) WithStyle(s Style) *Confirm {
	c.style = &s
	return c
}

// WithoutStyle removes the confirm's style, rendering plain text.
func (c *Confirm) WithoutStyle() *Confirm {
	c.style = nil
	return c
}

func (c *Confirm) SelectedStyle() *Style {
	return c.selectedStyle
}

// WithSelectedStyle sets the style applied to the currently selected option,
// overriding Style for that option.
func (c *Confirm) WithSelectedStyle(s Style) *Confirm {
	c.selectedStyle = &s
	return c
}

// WithoutSelectedStyle removes the selected option's style, falling back to
// Style for it.
func (c *Confirm) WithoutSelectedStyle() *Confirm {
	c.selectedStyle = nil
	return c
}

func (c *Confirm) Focused() bool {
	return c.focused
}

// Focus gives the confirm keyboard focus.
func (c *Confirm) Focus() {
	c.focused = true
}

// Blur removes the confirm's keyboard focus, ignoring key events.
func (c *Confirm) Blur() {
	c.focused = false
}

// AsFocused sets the confirm's keyboard focus.
func (c *Confirm) AsFocused(v bool) *Confirm {
	if v {
		c.Focus()
	} else {
		c.Blur()
	}
	return c
}

// PreferredWidth is the combined width of the message and both options, as
// currently rendered.
func (c *Confirm) PreferredWidth() int {
	return StringWidth(c.message) + 1 + StringWidth(c.yesText()) + 2 + StringWidth(c.noText())
}

func (c *Confirm) PreferredHeight(width int) int {
	return 1
}

// Update toggles the selected option, while Focused, in response to
// KeyEvents: Left, Right, and Tab flip it, and y/n (either case) select it
// directly.
func (c *Confirm) Update(e Event) Event {
	if ev, ok := e.(KeyEvent); ok && c.focused {
		c.handleKey(ev)
	}
	return e
}

// Render draws the message followed by the Yes and No options, bracketing
// whichever is currently selected.
func (c *Confirm) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	var b strings.Builder
	remaining := width
	remaining = renderCounterSegment(&b, c.message, c.style, remaining)
	remaining = renderCounterSegment(&b, " ", c.style, remaining)
	remaining = renderCounterSegment(&b, c.yesText(), c.yesStyle(), remaining)
	remaining = renderCounterSegment(&b, "  ", c.style, remaining)
	renderCounterSegment(&b, c.noText(), c.noStyle(), remaining)

	return []string{b.String()}
}

// handleKey applies a single key event's effect: flipping the selected
// option or setting it directly.
func (c *Confirm) handleKey(e KeyEvent) {
	switch e.Type {
	case KeyLeft, KeyRight, KeyTab:
		c.WithValue(!c.value)
	case KeyRune:
		switch e.Rune {
		case 'y', 'Y':
			c.WithValue(true)
		case 'n', 'N':
			c.WithValue(false)
		}
	}
}

// sync writes the current selection back into the bound pointer, if any.
func (c *Confirm) sync() {
	if c.valuePtr != nil {
		*c.valuePtr = c.value
	}
}

// yesText renders the Yes label, bracketed when it's the selected option.
func (c *Confirm) yesText() string {
	if c.value {
		return "[" + c.yesLabel + "]"
	}
	return c.yesLabel
}

// noText renders the No label, bracketed when it's the selected option.
func (c *Confirm) noText() string {
	if !c.value {
		return "[" + c.noLabel + "]"
	}
	return c.noLabel
}

// yesStyle returns the style used for the Yes option: SelectedStyle when
// it's selected and set, Style otherwise.
func (c *Confirm) yesStyle() *Style {
	if c.value && c.selectedStyle != nil {
		return c.selectedStyle
	}
	return c.style
}

// noStyle returns the style used for the No option: SelectedStyle when it's
// selected and set, Style otherwise.
func (c *Confirm) noStyle() *Style {
	if !c.value && c.selectedStyle != nil {
		return c.selectedStyle
	}
	return c.style
}
