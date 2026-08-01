package term

import "strings"

// DefaultTextInputWidth is used until WithWidth is called. It's large enough
// to be effectively unbounded, so by default the input fills whatever width
// it's rendered with rather than a fixed size.
const DefaultTextInputWidth = 9999

// TextInput renders a single line of editable text with a blinking cursor,
// scrolling its content horizontally so the cursor always stays visible
// within its width. It only reacts to key events while Focused.
type TextInput struct {
	value    []rune
	valuePtr *string
	style    *Style
	cursor   *Cursor
	pos      int
	offset   int
	width    int
	focused  bool
}

func NewTextInput() *TextInput {
	return &TextInput{
		cursor: NewCursor(),
		width:  DefaultTextInputWidth,
	}
}

func (t *TextInput) Value() string {
	return string(t.value)
}

// WithValue sets the input's text and moves the cursor to its end.
func (t *TextInput) WithValue(text string) *TextInput {
	t.value = []rune(text)
	t.pos = len(t.value)
	t.sync()
	return t
}

func (t *TextInput) ValuePtr() *string {
	return t.valuePtr
}

// WithValuePtr binds the input to p: its initial text is read from *p, and
// every subsequent edit writes the current value back into *p.
func (t *TextInput) WithValuePtr(p *string) *TextInput {
	t.valuePtr = p
	if p != nil {
		t.value = []rune(*p)
		t.pos = len(t.value)
	}
	return t
}

func (t *TextInput) Style() *Style {
	return t.style
}

// WithStyle sets the style applied to the input's text when rendering.
func (t *TextInput) WithStyle(s Style) *TextInput {
	t.style = &s
	return t
}

// WithoutStyle removes the input's style, rendering plain text.
func (t *TextInput) WithoutStyle() *TextInput {
	t.style = nil
	return t
}

// Cursor returns the input's internal cursor, so its character, style,
// blinking, visibility, and blink speed can be customized.
func (t *TextInput) Cursor() *Cursor {
	return t.cursor
}

func (t *TextInput) Width() int {
	return t.width
}

// WithWidth sets the input's preferred display width, in columns.
func (t *TextInput) WithWidth(n int) *TextInput {
	t.width = n
	return t
}

func (t *TextInput) Focused() bool {
	return t.focused
}

// Focus gives the input keyboard focus and restarts its cursor's blink cycle
// so it starts out fully shown.
func (t *TextInput) Focus() *TextInput {
	t.focused = true
	t.cursor.AsBlinking(t.cursor.Blinking())
	return t
}

// Blur removes the input's keyboard focus, hiding its cursor and ignoring
// key events.
func (t *TextInput) Blur() *TextInput {
	t.focused = false
	return t
}

func (t *TextInput) PreferredWidth() int {
	return t.width
}

func (t *TextInput) PreferredHeight(width int) int {
	return 1
}

// Update advances the cursor's blink state on every TickEvent, and, while
// Focused, edits the value and moves the cursor in response to KeyEvents.
func (t *TextInput) Update(e Event) Event {
	switch ev := e.(type) {
	case TickEvent:
		t.cursor.Update(ev)
	case KeyEvent:
		if t.focused {
			t.handleKey(ev)
		}
	}
	return e
}

// Render draws the visible window of the input's text, scrolled so the
// cursor stays within the smaller of its configured width and the available
// render width, overlaying the cursor's character at its position while
// Focused and showing. The returned line is always padded to width columns.
func (t *TextInput) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	visibleWidth := min(t.width, width)

	t.scrollIntoView(visibleWidth)

	end := min(len(t.value), t.offset+visibleWidth)
	visible := t.value[t.offset:end]
	cursorCol := t.pos - t.offset
	showCursor := t.focused && t.cursor.Showing()

	var b strings.Builder
	written := 0
	for i, r := range visible {
		if showCursor && i == cursorCol {
			b.WriteString(t.renderCursor(string(r)))
		} else {
			b.WriteString(t.renderText(string(r)))
		}
		written++
	}
	if showCursor && cursorCol == len(visible) {
		b.WriteString(t.renderCursor(" "))
		written++
	}
	for ; written < width; written++ {
		b.WriteByte(' ')
	}

	return []string{b.String()}
}

// handleKey applies a single key event's editing effect: inserting a rune,
// deleting one, or moving the cursor.
func (t *TextInput) handleKey(e KeyEvent) {
	switch e.Type {
	case KeyRune:
		if !e.Ctrl && !e.Alt {
			t.insert(e.Rune)
		}
	case KeyBackspace:
		t.deleteBefore()
	case KeyDelete:
		t.deleteAt()
	case KeyLeft:
		t.pos = max(0, t.pos-1)
	case KeyRight:
		t.pos = min(len(t.value), t.pos+1)
	case KeyHome:
		t.pos = 0
	case KeyEnd:
		t.pos = len(t.value)
	}
}

// insert adds r at the cursor's position and advances past it.
func (t *TextInput) insert(r rune) {
	t.value = append(t.value[:t.pos:t.pos], append([]rune{r}, t.value[t.pos:]...)...)
	t.pos++
	t.sync()
}

// deleteBefore removes the rune immediately before the cursor, if any.
func (t *TextInput) deleteBefore() {
	if t.pos == 0 {
		return
	}
	t.value = append(t.value[:t.pos-1], t.value[t.pos:]...)
	t.pos--
	t.sync()
}

// deleteAt removes the rune at the cursor, if any.
func (t *TextInput) deleteAt() {
	if t.pos >= len(t.value) {
		return
	}
	t.value = append(t.value[:t.pos], t.value[t.pos+1:]...)
	t.sync()
}

// sync writes the current value back into the bound pointer, if any.
func (t *TextInput) sync() {
	if t.valuePtr != nil {
		*t.valuePtr = string(t.value)
	}
}

// scrollIntoView adjusts the offset so the cursor's column stays within the
// visible window of width columns, scrolling as little as possible.
func (t *TextInput) scrollIntoView(width int) {
	if t.pos < t.offset {
		t.offset = t.pos
	}
	if t.pos > t.offset+width-1 {
		t.offset = t.pos - width + 1
	}
	t.offset = max(0, t.offset)
}

// renderText styles a single rune of plain text using the input's style.
func (t *TextInput) renderText(s string) string {
	if t.style != nil {
		return t.style.Render(s)
	}
	return s
}

// renderCursor styles the cursor's character (falling back to the character
// it's overlaying, under, when the cursor has none) using the cursor's
// style. Only the cursor's first rune is used, so it always occupies a
// single column.
func (t *TextInput) renderCursor(under string) string {
	ch := under
	if c := []rune(t.cursor.Char()); len(c) > 0 {
		ch = string(c[0])
	}
	if s := t.cursor.Style(); s != nil {
		return s.Render(ch)
	}
	return ch
}
