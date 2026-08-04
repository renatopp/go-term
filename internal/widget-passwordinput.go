package term

import "strings"

// DefaultPasswordInputWidth is used until WithWidth is called. It's large enough
// to be effectively unbounded, so by default the input fills whatever width
// it's rendered with rather than a fixed size.
const DefaultPasswordInputWidth = 9999

// DefaultPasswordInputMaskChar is the character rendered in place of each entered
// rune until WithMaskChar is called.
const DefaultPasswordInputMaskChar = "•"

// PasswordInput renders a single line of editable text like TextInput, except
// each entered rune is rendered as MaskChar instead of its actual character,
// so the value stays hidden on screen while still scrolling horizontally so
// the cursor always stays visible within its width. It only reacts to key
// events while Focused.
type PasswordInput struct {
	value            []rune
	valuePtr         *string
	style            *Style
	maskChar         string
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

func NewPasswordInput() *PasswordInput {
	return &PasswordInput{
		cursor:   NewCursor(),
		width:    DefaultPasswordInputWidth,
		maskChar: DefaultPasswordInputMaskChar,
	}
}

func (p *PasswordInput) Value() string {
	return string(p.value)
}

// WithValue sets the input's text and moves the cursor to its end.
func (p *PasswordInput) WithValue(text string) *PasswordInput {
	p.value = []rune(text)
	p.pos = len(p.value)
	p.sync()
	return p
}

func (p *PasswordInput) ValuePtr() *string {
	return p.valuePtr
}

// WithValuePtr binds the input to v: its initial text is read from *v, and
// every subsequent edit writes the current value back into *v.
func (p *PasswordInput) WithValuePtr(v *string) *PasswordInput {
	p.valuePtr = v
	if v != nil {
		p.value = []rune(*v)
		p.pos = len(p.value)
	}
	return p
}

func (p *PasswordInput) Style() *Style {
	return p.style
}

// WithStyle sets the style applied to the input's masked text when
// rendering.
func (p *PasswordInput) WithStyle(s Style) *PasswordInput {
	p.style = &s
	return p
}

// WithoutStyle removes the input's style, rendering plain text.
func (p *PasswordInput) WithoutStyle() *PasswordInput {
	p.style = nil
	return p
}

func (p *PasswordInput) MaskChar() string {
	return p.maskChar
}

// WithMaskChar sets the character rendered in place of each entered rune.
// Only its first rune is used, so it always occupies a single column.
func (p *PasswordInput) WithMaskChar(s string) *PasswordInput {
	p.maskChar = s
	return p
}

func (p *PasswordInput) Placeholder() string {
	return p.placeholder
}

// WithPlaceholder sets the text shown in place of the field when the input's
// value is empty.
func (p *PasswordInput) WithPlaceholder(s string) *PasswordInput {
	p.placeholder = s
	return p
}

func (p *PasswordInput) PlaceholderStyle() *Style {
	return p.placeholderStyle
}

// WithPlaceholderStyle sets the style applied to the input's placeholder when rendering.
func (p *PasswordInput) WithPlaceholderStyle(s Style) *PasswordInput {
	p.placeholderStyle = &s
	return p
}

// WithoutPlaceholderStyle removes the input's placeholder style, rendering it plain.
func (p *PasswordInput) WithoutPlaceholderStyle() *PasswordInput {
	p.placeholderStyle = nil
	return p
}

func (p *PasswordInput) Prefix() string {
	return p.prefix
}

// WithPrefix sets the text rendered before the input's editable field.
func (p *PasswordInput) WithPrefix(s string) *PasswordInput {
	p.prefix = s
	return p
}

func (p *PasswordInput) PrefixStyle() *Style {
	return p.prefixStyle
}

// WithPrefixStyle sets the style applied to the input's prefix when rendering.
func (p *PasswordInput) WithPrefixStyle(s Style) *PasswordInput {
	p.prefixStyle = &s
	return p
}

// WithoutPrefixStyle removes the input's prefix style, rendering it plain.
func (p *PasswordInput) WithoutPrefixStyle() *PasswordInput {
	p.prefixStyle = nil
	return p
}

func (p *PasswordInput) Suffix() string {
	return p.suffix
}

// WithSuffix sets the text rendered after the input's editable field.
func (p *PasswordInput) WithSuffix(s string) *PasswordInput {
	p.suffix = s
	return p
}

func (p *PasswordInput) SuffixStyle() *Style {
	return p.suffixStyle
}

// WithSuffixStyle sets the style applied to the input's suffix when rendering.
func (p *PasswordInput) WithSuffixStyle(s Style) *PasswordInput {
	p.suffixStyle = &s
	return p
}

// WithoutSuffixStyle removes the input's suffix style, rendering it plain.
func (p *PasswordInput) WithoutSuffixStyle() *PasswordInput {
	p.suffixStyle = nil
	return p
}

// Cursor returns the input's internal cursor, so its character, style,
// blinking, visibility, and blink speed can be customized.
func (p *PasswordInput) Cursor() *Cursor {
	return p.cursor
}

func (p *PasswordInput) Width() int {
	return p.width
}

// WithWidth sets the input's preferred display width, in columns.
func (p *PasswordInput) WithWidth(n int) *PasswordInput {
	p.width = n
	return p
}

func (p *PasswordInput) Focused() bool {
	return p.focused
}

// Focus gives the input keyboard focus and restarts its cursor's blink cycle
// so it starts out fully shown.
func (p *PasswordInput) Focus() {
	p.focused = true
	p.cursor.AsBlinking(p.cursor.Blinking())
}

// Blur removes the input's keyboard focus, hiding its cursor and ignoring
// key events.
func (p *PasswordInput) Blur() {
	p.focused = false
}

// AsFocused sets the input's keyboard focus, restarting its cursor's blink
// cycle when focused (see Focus).
func (p *PasswordInput) AsFocused(v bool) *PasswordInput {
	if v {
		p.Focus()
	} else {
		p.Blur()
	}
	return p
}

// PreferredWidth is the input's preferred field width (see WithWidth) plus
// its prefix and suffix.
func (p *PasswordInput) PreferredWidth() int {
	return len([]rune(p.prefix)) + p.width + len([]rune(p.suffix))
}

func (p *PasswordInput) PreferredHeight(width int) int {
	return 1
}

// Update advances the cursor's blink state on every TickEvent, and, while
// Focused, edits the value and moves the cursor in response to KeyEvents.
func (p *PasswordInput) Update(e Event) Event {
	switch ev := e.(type) {
	case TickEvent:
		p.cursor.Update(ev)
	case KeyEvent:
		if p.focused {
			p.handleKey(ev)
		}
	}
	return e
}

// Render draws the input's prefix, its editable field, and its suffix, in
// that order. The field is scrolled so the cursor stays within the space
// left over once the prefix and suffix (clipped first, since they don't
// scroll) are accounted for, overlaying the cursor's character at its
// position while Focused and showing. Every entered rune is rendered as
// MaskChar. The returned line is always padded to width columns.
func (p *PasswordInput) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	prefix := []rune(p.prefix)
	suffix := []rune(p.suffix)
	prefixWidth := min(len(prefix), width)
	suffixWidth := min(len(suffix), width-prefixWidth)
	fieldWidth := width - prefixWidth - suffixWidth

	var b strings.Builder
	written := 0
	if prefixWidth > 0 {
		b.WriteString(p.renderPrefix(string(prefix[:prefixWidth])))
		written += prefixWidth
	}

	if fieldWidth > 0 {
		visibleWidth := min(p.width, fieldWidth)
		showCursor := p.focused && p.cursor.Showing()

		var fieldWritten int
		if len(p.value) == 0 && p.placeholder != "" {
			fieldWritten = p.renderPlaceholderInto(&b, visibleWidth, showCursor)
		} else {
			p.scrollIntoView(visibleWidth)

			end := min(len(p.value), p.offset+visibleWidth)
			visible := p.value[p.offset:end]
			cursorCol := p.pos - p.offset

			for i := range visible {
				if showCursor && i == cursorCol {
					b.WriteString(p.renderCursor(p.mask()))
				} else {
					b.WriteString(p.renderText(p.mask()))
				}
				fieldWritten++
			}
			if showCursor && cursorCol == len(visible) {
				b.WriteString(p.renderCursor(" "))
				fieldWritten++
			}
		}
		for ; fieldWritten < fieldWidth; fieldWritten++ {
			b.WriteByte(' ')
		}
		written += fieldWidth
	}

	if suffixWidth > 0 {
		b.WriteString(p.renderSuffix(string(suffix[:suffixWidth])))
		written += suffixWidth
	}

	for ; written < width; written++ {
		b.WriteByte(' ')
	}

	return []string{b.String()}
}

// handleKey applies a single key event's editing effect: inserting a rune,
// deleting one, or moving the cursor.
func (p *PasswordInput) handleKey(e KeyEvent) {
	switch e.Type {
	case KeyRune:
		if !e.Ctrl && !e.Alt {
			p.insert(e.Rune)
		}
	case KeyBackspace:
		p.deleteBefore()
	case KeyDelete:
		p.deleteAt()
	case KeyLeft:
		p.pos = max(0, p.pos-1)
	case KeyRight:
		p.pos = min(len(p.value), p.pos+1)
	case KeyHome:
		p.pos = 0
	case KeyEnd:
		p.pos = len(p.value)
	}
}

// insert adds r at the cursor's position and advances past it.
func (p *PasswordInput) insert(r rune) {
	p.value = append(p.value[:p.pos:p.pos], append([]rune{r}, p.value[p.pos:]...)...)
	p.pos++
	p.sync()
}

// deleteBefore removes the rune immediately before the cursor, if any.
func (p *PasswordInput) deleteBefore() {
	if p.pos == 0 {
		return
	}
	p.value = append(p.value[:p.pos-1], p.value[p.pos:]...)
	p.pos--
	p.sync()
}

// deleteAt removes the rune at the cursor, if any.
func (p *PasswordInput) deleteAt() {
	if p.pos >= len(p.value) {
		return
	}
	p.value = append(p.value[:p.pos], p.value[p.pos+1:]...)
	p.sync()
}

// sync writes the current value back into the bound pointer, if any.
func (p *PasswordInput) sync() {
	if p.valuePtr != nil {
		*p.valuePtr = string(p.value)
	}
}

// scrollIntoView adjusts the offset so the cursor's column stays within the
// visible window of width columns, scrolling as little as possible.
func (p *PasswordInput) scrollIntoView(width int) {
	if p.pos < p.offset {
		p.offset = p.pos
	}
	if p.pos > p.offset+width-1 {
		p.offset = p.pos - width + 1
	}
	p.offset = max(0, p.offset)
}

// mask returns the single-column character rendered in place of each entered
// rune, falling back to a space when MaskChar is empty.
func (p *PasswordInput) mask() string {
	if c := []rune(p.maskChar); len(c) > 0 {
		return string(c[0])
	}
	return " "
}

// renderText styles a single rune of masked text using the input's style.
func (p *PasswordInput) renderText(s string) string {
	if p.style != nil {
		return p.style.Render(s)
	}
	return s
}

// renderPlaceholderInto writes up to width runes of the input's placeholder
// to b, styled with the placeholder style, overlaying the cursor's character
// at the first column while showCursor. It returns the number of columns
// written.
func (p *PasswordInput) renderPlaceholderInto(b *strings.Builder, width int, showCursor bool) int {
	placeholder := []rune(p.placeholder)
	end := min(len(placeholder), width)
	visible := placeholder[:end]

	written := 0
	for i, r := range visible {
		if showCursor && i == 0 {
			b.WriteString(p.renderCursor(string(r)))
		} else {
			b.WriteString(p.renderPlaceholder(string(r)))
		}
		written++
	}
	if showCursor && len(visible) == 0 {
		b.WriteString(p.renderCursor(" "))
		written++
	}
	return written
}

// renderPlaceholder styles a single rune of placeholder text using the
// input's placeholder style.
func (p *PasswordInput) renderPlaceholder(s string) string {
	if p.placeholderStyle != nil {
		return p.placeholderStyle.Render(s)
	}
	return s
}

// renderPrefix styles the input's prefix using its prefix style.
func (p *PasswordInput) renderPrefix(s string) string {
	if p.prefixStyle != nil {
		return p.prefixStyle.Render(s)
	}
	return s
}

// renderSuffix styles the input's suffix using its suffix style.
func (p *PasswordInput) renderSuffix(s string) string {
	if p.suffixStyle != nil {
		return p.suffixStyle.Render(s)
	}
	return s
}

// renderCursor styles the cursor's character (falling back to the character
// it's overlaying, under, when the cursor has none) using the cursor's
// style. Only the cursor's first rune is used, so it always occupies a
// single column.
func (p *PasswordInput) renderCursor(under string) string {
	ch := under
	if c := []rune(p.cursor.Char()); len(c) > 0 {
		ch = string(c[0])
	}
	if s := p.cursor.Style(); s != nil {
		return s.Render(ch)
	}
	return ch
}

