package term

import (
	"strconv"
	"strings"
)

// DefaultNumberInputWidth is used until WithWidth is called. It's large
// enough to be effectively unbounded, so by default the input fills whatever
// width it's rendered with rather than a fixed size.
const DefaultNumberInputWidth = 9999

// NumberInput renders a single line of editable numeric text with a
// blinking cursor, scrolling its content horizontally so the cursor always
// stays visible within its width. It only accepts digits, a single leading
// minus sign, and a single decimal point, and only reacts to key events
// while Focused.
type NumberInput struct {
	value            []rune
	valuePtr         *float64
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

func NewNumberInput() *NumberInput {
	return &NumberInput{
		cursor: NewCursor(),
		width:  DefaultNumberInputWidth,
	}
}

// Value parses the input's current text as a float64. It returns an error
// while the text is empty or in a mid-edit state that isn't a valid number
// yet (e.g. "-" or "").
func (n *NumberInput) Value() (float64, error) {
	return strconv.ParseFloat(string(n.value), 64)
}

// ForceValue returns the input's value, ignoring any parse error and
// returning 0 in its place.
func (n *NumberInput) ForceValue() float64 {
	v, _ := n.Value()
	return v
}

// WithValue sets the input's numeric value and moves the cursor to its end.
func (n *NumberInput) WithValue(v float64) *NumberInput {
	n.value = []rune(strconv.FormatFloat(v, 'f', -1, 64))
	n.pos = len(n.value)
	n.sync()
	return n
}

func (n *NumberInput) ValuePtr() *float64 {
	return n.valuePtr
}

// WithValuePtr binds the input to p: its initial text is read from *p, and
// every subsequent edit that results in a valid number writes it back into
// *p.
func (n *NumberInput) WithValuePtr(p *float64) *NumberInput {
	n.valuePtr = p
	if p != nil {
		n.value = []rune(strconv.FormatFloat(*p, 'f', -1, 64))
		n.pos = len(n.value)
	}
	return n
}

func (n *NumberInput) Style() *Style {
	return n.style
}

// WithStyle sets the style applied to the input's text when rendering.
func (n *NumberInput) WithStyle(s Style) *NumberInput {
	n.style = &s
	return n
}

// WithoutStyle removes the input's style, rendering plain text.
func (n *NumberInput) WithoutStyle() *NumberInput {
	n.style = nil
	return n
}

func (n *NumberInput) Placeholder() string {
	return n.placeholder
}

// WithPlaceholder sets the text shown in place of the field when the
// input's value is empty.
func (n *NumberInput) WithPlaceholder(s string) *NumberInput {
	n.placeholder = s
	return n
}

func (n *NumberInput) PlaceholderStyle() *Style {
	return n.placeholderStyle
}

// WithPlaceholderStyle sets the style applied to the input's placeholder
// when rendering.
func (n *NumberInput) WithPlaceholderStyle(s Style) *NumberInput {
	n.placeholderStyle = &s
	return n
}

// WithoutPlaceholderStyle removes the input's placeholder style, rendering
// it plain.
func (n *NumberInput) WithoutPlaceholderStyle() *NumberInput {
	n.placeholderStyle = nil
	return n
}

func (n *NumberInput) Prefix() string {
	return n.prefix
}

// WithPrefix sets the text rendered before the input's editable field.
func (n *NumberInput) WithPrefix(s string) *NumberInput {
	n.prefix = s
	return n
}

func (n *NumberInput) PrefixStyle() *Style {
	return n.prefixStyle
}

// WithPrefixStyle sets the style applied to the input's prefix when
// rendering.
func (n *NumberInput) WithPrefixStyle(s Style) *NumberInput {
	n.prefixStyle = &s
	return n
}

// WithoutPrefixStyle removes the input's prefix style, rendering it plain.
func (n *NumberInput) WithoutPrefixStyle() *NumberInput {
	n.prefixStyle = nil
	return n
}

func (n *NumberInput) Suffix() string {
	return n.suffix
}

// WithSuffix sets the text rendered after the input's editable field.
func (n *NumberInput) WithSuffix(s string) *NumberInput {
	n.suffix = s
	return n
}

func (n *NumberInput) SuffixStyle() *Style {
	return n.suffixStyle
}

// WithSuffixStyle sets the style applied to the input's suffix when
// rendering.
func (n *NumberInput) WithSuffixStyle(s Style) *NumberInput {
	n.suffixStyle = &s
	return n
}

// WithoutSuffixStyle removes the input's suffix style, rendering it plain.
func (n *NumberInput) WithoutSuffixStyle() *NumberInput {
	n.suffixStyle = nil
	return n
}

// Cursor returns the input's internal cursor, so its character, style,
// blinking, visibility, and blink speed can be customized.
func (n *NumberInput) Cursor() *Cursor {
	return n.cursor
}

func (n *NumberInput) Width() int {
	return n.width
}

// WithWidth sets the input's preferred display width, in columns.
func (n *NumberInput) WithWidth(w int) *NumberInput {
	n.width = w
	return n
}

func (n *NumberInput) Focused() bool {
	return n.focused
}

// Focus gives the input keyboard focus and restarts its cursor's blink
// cycle so it starts out fully shown.
func (n *NumberInput) Focus() {
	n.focused = true
	n.cursor.AsBlinking(n.cursor.Blinking())
}

// Blur removes the input's keyboard focus, hiding its cursor and ignoring
// key events.
func (n *NumberInput) Blur() {
	n.focused = false
}

// AsFocused sets the input's keyboard focus, restarting its cursor's blink
// cycle when focused (see Focus).
func (n *NumberInput) AsFocused(v bool) *NumberInput {
	if v {
		n.Focus()
	} else {
		n.Blur()
	}
	return n
}

// PreferredWidth is the input's preferred field width (see WithWidth) plus
// its prefix and suffix.
func (n *NumberInput) PreferredWidth() int {
	return len([]rune(n.prefix)) + n.width + len([]rune(n.suffix))
}

func (n *NumberInput) PreferredHeight(width int) int {
	return 1
}

// Update advances the cursor's blink state on every TickEvent, and, while
// Focused, edits the value and moves the cursor in response to KeyEvents.
func (n *NumberInput) Update(e Event) Event {
	switch ev := e.(type) {
	case TickEvent:
		n.cursor.Update(ev)
	case KeyEvent:
		if n.focused {
			n.handleKey(ev)
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
func (n *NumberInput) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	prefix := []rune(n.prefix)
	suffix := []rune(n.suffix)
	prefixWidth := min(len(prefix), width)
	suffixWidth := min(len(suffix), width-prefixWidth)
	fieldWidth := width - prefixWidth - suffixWidth

	var b strings.Builder
	written := 0
	if prefixWidth > 0 {
		b.WriteString(n.renderPrefix(string(prefix[:prefixWidth])))
		written += prefixWidth
	}

	if fieldWidth > 0 {
		visibleWidth := min(n.width, fieldWidth)
		showCursor := n.focused && n.cursor.Showing()

		var fieldWritten int
		if len(n.value) == 0 && n.placeholder != "" {
			fieldWritten = n.renderPlaceholderInto(&b, visibleWidth, showCursor)
		} else {
			n.scrollIntoView(visibleWidth)

			end := min(len(n.value), n.offset+visibleWidth)
			visible := n.value[n.offset:end]
			cursorCol := n.pos - n.offset

			for i, r := range visible {
				if showCursor && i == cursorCol {
					b.WriteString(n.renderCursor(string(r)))
				} else {
					b.WriteString(n.renderText(string(r)))
				}
				fieldWritten++
			}
			if showCursor && cursorCol == len(visible) {
				b.WriteString(n.renderCursor(" "))
				fieldWritten++
			}
		}
		for ; fieldWritten < fieldWidth; fieldWritten++ {
			b.WriteByte(' ')
		}
		written += fieldWidth
	}

	if suffixWidth > 0 {
		b.WriteString(n.renderSuffix(string(suffix[:suffixWidth])))
		written += suffixWidth
	}

	for ; written < width; written++ {
		b.WriteByte(' ')
	}

	return []string{b.String()}
}

// handleKey applies a single key event's editing effect: inserting a rune,
// deleting one, or moving the cursor.
func (n *NumberInput) handleKey(e KeyEvent) {
	switch e.Type {
	case KeyRune:
		if !e.Ctrl && !e.Alt && n.acceptsRune(e.Rune) {
			n.insert(e.Rune)
		}
	case KeyBackspace:
		n.deleteBefore()
	case KeyDelete:
		n.deleteAt()
	case KeyLeft:
		n.pos = max(0, n.pos-1)
	case KeyRight:
		n.pos = min(len(n.value), n.pos+1)
	case KeyHome:
		n.pos = 0
	case KeyEnd:
		n.pos = len(n.value)
	}
}

// acceptsRune reports whether r can be inserted into the input's value:
// digits are always allowed, '-' only as the very first character, and '.'
// only once.
func (n *NumberInput) acceptsRune(r rune) bool {
	switch {
	case r >= '0' && r <= '9':
		return true
	case r == '-':
		return n.pos == 0 && (len(n.value) == 0 || n.value[0] != '-')
	case r == '.':
		return !strings.ContainsRune(string(n.value), '.')
	default:
		return false
	}
}

// insert adds r at the cursor's position and advances past it.
func (n *NumberInput) insert(r rune) {
	n.value = append(n.value[:n.pos:n.pos], append([]rune{r}, n.value[n.pos:]...)...)
	n.pos++
	n.sync()
}

// deleteBefore removes the rune immediately before the cursor, if any.
func (n *NumberInput) deleteBefore() {
	if n.pos == 0 {
		return
	}
	n.value = append(n.value[:n.pos-1], n.value[n.pos:]...)
	n.pos--
	n.sync()
}

// deleteAt removes the rune at the cursor, if any.
func (n *NumberInput) deleteAt() {
	if n.pos >= len(n.value) {
		return
	}
	n.value = append(n.value[:n.pos], n.value[n.pos+1:]...)
	n.sync()
}

// sync writes the current value back into the bound pointer, if it parses
// as a valid number.
func (n *NumberInput) sync() {
	if n.valuePtr != nil {
		if v, err := n.Value(); err == nil {
			*n.valuePtr = v
		}
	}
}

// scrollIntoView adjusts the offset so the cursor's column stays within the
// visible window of width columns, scrolling as little as possible.
func (n *NumberInput) scrollIntoView(width int) {
	if n.pos < n.offset {
		n.offset = n.pos
	}
	if n.pos > n.offset+width-1 {
		n.offset = n.pos - width + 1
	}
	n.offset = max(0, n.offset)
}

// renderText styles a single rune of plain text using the input's style.
func (n *NumberInput) renderText(s string) string {
	if n.style != nil {
		return n.style.Render(s)
	}
	return s
}

// renderPlaceholderInto writes up to width runes of the input's placeholder
// to b, styled with the placeholder style, overlaying the cursor's
// character at the first column while showCursor. It returns the number of
// columns written.
func (n *NumberInput) renderPlaceholderInto(b *strings.Builder, width int, showCursor bool) int {
	placeholder := []rune(n.placeholder)
	end := min(len(placeholder), width)
	visible := placeholder[:end]

	written := 0
	for i, r := range visible {
		if showCursor && i == 0 {
			b.WriteString(n.renderCursor(string(r)))
		} else {
			b.WriteString(n.renderPlaceholder(string(r)))
		}
		written++
	}
	if showCursor && len(visible) == 0 {
		b.WriteString(n.renderCursor(" "))
		written++
	}
	return written
}

// renderPlaceholder styles a single rune of placeholder text using the
// input's placeholder style.
func (n *NumberInput) renderPlaceholder(s string) string {
	if n.placeholderStyle != nil {
		return n.placeholderStyle.Render(s)
	}
	return s
}

// renderPrefix styles the input's prefix using its prefix style.
func (n *NumberInput) renderPrefix(s string) string {
	if n.prefixStyle != nil {
		return n.prefixStyle.Render(s)
	}
	return s
}

// renderSuffix styles the input's suffix using its suffix style.
func (n *NumberInput) renderSuffix(s string) string {
	if n.suffixStyle != nil {
		return n.suffixStyle.Render(s)
	}
	return s
}

// renderCursor styles the cursor's character (falling back to the character
// it's overlaying, under, when the cursor has none) using the cursor's
// style. Only the cursor's first rune is used, so it always occupies a
// single column.
func (n *NumberInput) renderCursor(under string) string {
	ch := under
	if c := []rune(n.cursor.Char()); len(c) > 0 {
		ch = string(c[0])
	}
	if s := n.cursor.Style(); s != nil {
		return s.Render(ch)
	}
	return ch
}
