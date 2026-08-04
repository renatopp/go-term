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
	value            []rune
	valuePtr         *string
	style            *Style
	placeholder      string
	placeholderStyle *Style
	prefix           string
	prefixStyle      *Style
	suffix           string
	suffixStyle      *Style
	cursor           *Cursor
	pos              int
	offset           int
	width            int
	focused          bool
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

func (t *TextInput) Placeholder() string {
	return t.placeholder
}

// WithPlaceholder sets the text shown in place of the field when the input's
// value is empty.
func (t *TextInput) WithPlaceholder(s string) *TextInput {
	t.placeholder = s
	return t
}

func (t *TextInput) PlaceholderStyle() *Style {
	return t.placeholderStyle
}

// WithPlaceholderStyle sets the style applied to the input's placeholder when rendering.
func (t *TextInput) WithPlaceholderStyle(s Style) *TextInput {
	t.placeholderStyle = &s
	return t
}

// WithoutPlaceholderStyle removes the input's placeholder style, rendering it plain.
func (t *TextInput) WithoutPlaceholderStyle() *TextInput {
	t.placeholderStyle = nil
	return t
}

func (t *TextInput) Prefix() string {
	return t.prefix
}

// WithPrefix sets the text rendered before the input's editable field.
func (t *TextInput) WithPrefix(s string) *TextInput {
	t.prefix = s
	return t
}

func (t *TextInput) PrefixStyle() *Style {
	return t.prefixStyle
}

// WithPrefixStyle sets the style applied to the input's prefix when rendering.
func (t *TextInput) WithPrefixStyle(s Style) *TextInput {
	t.prefixStyle = &s
	return t
}

// WithoutPrefixStyle removes the input's prefix style, rendering it plain.
func (t *TextInput) WithoutPrefixStyle() *TextInput {
	t.prefixStyle = nil
	return t
}

func (t *TextInput) Suffix() string {
	return t.suffix
}

// WithSuffix sets the text rendered after the input's editable field.
func (t *TextInput) WithSuffix(s string) *TextInput {
	t.suffix = s
	return t
}

func (t *TextInput) SuffixStyle() *Style {
	return t.suffixStyle
}

// WithSuffixStyle sets the style applied to the input's suffix when rendering.
func (t *TextInput) WithSuffixStyle(s Style) *TextInput {
	t.suffixStyle = &s
	return t
}

// WithoutSuffixStyle removes the input's suffix style, rendering it plain.
func (t *TextInput) WithoutSuffixStyle() *TextInput {
	t.suffixStyle = nil
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
func (t *TextInput) Focus() {
	t.focused = true
	t.cursor.AsBlinking(t.cursor.Blinking())
}

// Blur removes the input's keyboard focus, hiding its cursor and ignoring
// key events.
func (t *TextInput) Blur() {
	t.focused = false
}

// AsFocused sets the input's keyboard focus, restarting its cursor's blink
// cycle when focused (see Focus).
func (t *TextInput) AsFocused(v bool) *TextInput {
	if v {
		t.Focus()
	} else {
		t.Blur()
	}
	return t
}

// PreferredWidth is the input's preferred field width (see WithWidth) plus
// its prefix and suffix.
func (t *TextInput) PreferredWidth() int {
	return len([]rune(t.prefix)) + t.width + len([]rune(t.suffix))
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

// Render draws the input's prefix, its editable field, and its suffix, in
// that order. The field is scrolled so the cursor stays within the space
// left over once the prefix and suffix (clipped first, since they don't
// scroll) are accounted for, overlaying the cursor's character at its
// position while Focused and showing. The returned line is always padded to
// width columns.
func (t *TextInput) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	prefix := []rune(t.prefix)
	suffix := []rune(t.suffix)
	prefixWidth := min(len(prefix), width)
	suffixWidth := min(len(suffix), width-prefixWidth)
	fieldWidth := width - prefixWidth - suffixWidth

	var b strings.Builder
	written := 0
	if prefixWidth > 0 {
		b.WriteString(t.renderPrefix(string(prefix[:prefixWidth])))
		written += prefixWidth
	}

	if fieldWidth > 0 {
		visibleWidth := min(t.width, fieldWidth)
		showCursor := t.focused && t.cursor.Showing()

		var fieldWritten int
		if len(t.value) == 0 && t.placeholder != "" {
			fieldWritten = t.renderPlaceholderInto(&b, visibleWidth, showCursor)
		} else {
			t.scrollIntoView(visibleWidth)

			end := min(len(t.value), t.offset+visibleWidth)
			visible := t.value[t.offset:end]
			cursorCol := t.pos - t.offset

			for i, r := range visible {
				if showCursor && i == cursorCol {
					b.WriteString(t.renderCursor(string(r)))
				} else {
					b.WriteString(t.renderText(string(r)))
				}
				fieldWritten++
			}
			if showCursor && cursorCol == len(visible) {
				b.WriteString(t.renderCursor(" "))
				fieldWritten++
			}
		}
		for ; fieldWritten < fieldWidth; fieldWritten++ {
			b.WriteByte(' ')
		}
		written += fieldWidth
	}

	if suffixWidth > 0 {
		b.WriteString(t.renderSuffix(string(suffix[:suffixWidth])))
		written += suffixWidth
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

// renderPlaceholderInto writes up to width runes of the input's placeholder
// to b, styled with the placeholder style, overlaying the cursor's character
// at the first column while showCursor. It returns the number of columns
// written.
func (t *TextInput) renderPlaceholderInto(b *strings.Builder, width int, showCursor bool) int {
	placeholder := []rune(t.placeholder)
	end := min(len(placeholder), width)
	visible := placeholder[:end]

	written := 0
	for i, r := range visible {
		if showCursor && i == 0 {
			b.WriteString(t.renderCursor(string(r)))
		} else {
			b.WriteString(t.renderPlaceholder(string(r)))
		}
		written++
	}
	if showCursor && len(visible) == 0 {
		b.WriteString(t.renderCursor(" "))
		written++
	}
	return written
}

// renderPlaceholder styles a single rune of placeholder text using the
// input's placeholder style.
func (t *TextInput) renderPlaceholder(s string) string {
	if t.placeholderStyle != nil {
		return t.placeholderStyle.Render(s)
	}
	return s
}

// renderPrefix styles the input's prefix using its prefix style.
func (t *TextInput) renderPrefix(s string) string {
	if t.prefixStyle != nil {
		return t.prefixStyle.Render(s)
	}
	return s
}

// renderSuffix styles the input's suffix using its suffix style.
func (t *TextInput) renderSuffix(s string) string {
	if t.suffixStyle != nil {
		return t.suffixStyle.Render(s)
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
