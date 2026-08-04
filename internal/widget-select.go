package term

import (
	"strings"
)

// DefaultSelectMarker is the text rendered before the highlighted option
// until WithMarker is called.
const DefaultSelectMarker = "> "

// Select renders a vertical list of options, one highlighted at a time and
// moved with Up/Down. The highlighted option is prefixed with Marker and
// styled with SelectedStyle; the rest are indented to align under it and
// styled with Style. When there are more options than fit height, the list
// scrolls to keep the highlighted option visible. It only reacts to key
// events while Focused.
type Select struct {
	options       []string
	cursor        int
	valuePtr      *string
	marker        string
	style         *Style
	selectedStyle *Style
	offset        int
	focused       bool
}

func NewSelect(options ...string) *Select {
	return &Select{
		options: options,
		marker:  DefaultSelectMarker,
	}
}

func (s *Select) Options() []string {
	return s.options
}

// WithOption appends the given options to the select.
func (s *Select) WithOption(options ...string) *Select {
	s.options = append(s.options, options...)
	return s
}

func (s *Select) Marker() string {
	return s.marker
}

// WithMarker sets the text rendered before the highlighted option, with the
// rest of the options indented to align under it.
func (s *Select) WithMarker(m string) *Select {
	s.marker = m
	return s
}

// Value returns the currently highlighted option's text, or "" if there are
// no options.
func (s *Select) Value() string {
	if s.cursor < 0 || s.cursor >= len(s.options) {
		return ""
	}
	return s.options[s.cursor]
}

// WithValue highlights the given option, if it's present among Options; it
// has no effect otherwise.
func (s *Select) WithValue(v string) *Select {
	for i, o := range s.options {
		if o == v {
			s.cursor = i
			s.sync()
			break
		}
	}
	return s
}

func (s *Select) ValuePtr() *string {
	return s.valuePtr
}

// WithValuePtr binds the select to p: its initial highlighted option is read
// from *p (see WithValue), and every subsequent move writes the current
// option's text back into *p.
func (s *Select) WithValuePtr(p *string) *Select {
	s.valuePtr = p
	if p != nil {
		s.WithValue(*p)
	}
	return s
}

func (s *Select) Style() *Style {
	return s.style
}

// WithStyle sets the style applied to the unhighlighted options.
func (s *Select) WithStyle(st Style) *Select {
	s.style = &st
	return s
}

// WithoutStyle removes the select's style, rendering unhighlighted options
// plain.
func (s *Select) WithoutStyle() *Select {
	s.style = nil
	return s
}

func (s *Select) SelectedStyle() *Style {
	return s.selectedStyle
}

// WithSelectedStyle sets the style applied to the highlighted option,
// overriding Style for it.
func (s *Select) WithSelectedStyle(st Style) *Select {
	s.selectedStyle = &st
	return s
}

// WithoutSelectedStyle removes the highlighted option's style, falling back
// to Style for it.
func (s *Select) WithoutSelectedStyle() *Select {
	s.selectedStyle = nil
	return s
}

func (s *Select) Focused() bool {
	return s.focused
}

// Focus gives the select keyboard focus.
func (s *Select) Focus() {
	s.focused = true
}

// Blur removes the select's keyboard focus, ignoring key events.
func (s *Select) Blur() {
	s.focused = false
}

// AsFocused sets the select's keyboard focus.
func (s *Select) AsFocused(v bool) *Select {
	if v {
		s.Focus()
	} else {
		s.Blur()
	}
	return s
}

// PreferredWidth is the marker's width plus the widest option's width.
func (s *Select) PreferredWidth() int {
	w := 0
	for _, o := range s.options {
		w = max(w, StringWidth(o))
	}
	return StringWidth(s.marker) + w
}

// PreferredHeight is the number of options, since each renders as a single
// line.
func (s *Select) PreferredHeight(width int) int {
	return len(s.options)
}

// Update moves the highlighted option, while Focused, in response to
// KeyEvents: Up and Down move it by one, clamped to the first and last
// option.
func (s *Select) Update(e Event) Event {
	if ev, ok := e.(KeyEvent); ok && s.focused {
		s.handleKey(ev)
	}
	return e
}

// Render draws the visible window of options, one per line, prefixing the
// highlighted option with Marker and indenting the rest to align under it,
// each clipped to width columns. If there are more options than height, the
// window scrolls to keep the highlighted option visible.
func (s *Select) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	s.scrollIntoView(height)

	indent := strings.Repeat(" ", StringWidth(s.marker))
	end := min(len(s.options), s.offset+height)

	var lines []string
	for i := s.offset; i < end; i++ {
		prefix := indent
		style := s.style
		if i == s.cursor {
			prefix = s.marker
			if s.selectedStyle != nil {
				style = s.selectedStyle
			}
		}

		var b strings.Builder
		renderCounterSegment(&b, prefix+s.options[i], style, width)
		lines = append(lines, b.String())
	}
	return lines
}

// handleKey applies a single key event's effect: moving the highlighted
// option up or down.
func (s *Select) handleKey(e KeyEvent) {
	switch e.Type {
	case KeyUp:
		s.moveCursor(-1)
	case KeyDown:
		s.moveCursor(1)
	}
}

// moveCursor shifts the highlighted option by delta, clamped to the first
// and last option, and syncs the bound pointer.
func (s *Select) moveCursor(delta int) {
	if len(s.options) == 0 {
		return
	}
	s.cursor = max(0, min(s.cursor+delta, len(s.options)-1))
	s.sync()
}

// sync writes the current selection back into the bound pointer, if any.
func (s *Select) sync() {
	if s.valuePtr != nil {
		*s.valuePtr = s.Value()
	}
}

// scrollIntoView adjusts the offset so the highlighted option's row stays
// within the visible window of height rows, scrolling as little as
// possible.
func (s *Select) scrollIntoView(height int) {
	if s.cursor < s.offset {
		s.offset = s.cursor
	}
	if s.cursor > s.offset+height-1 {
		s.offset = s.cursor - height + 1
	}
	s.offset = max(0, s.offset)
}
