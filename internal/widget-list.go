package term

import (
	"strings"
)

// DefaultBullet is the bullet character used by a List until WithBullet is
// called.
const DefaultBullet = "•"

// List renders a vertical list of text items, each prefixed with a bullet.
// Items that don't fit width are wrapped at word boundaries, with wrapped
// lines indented to align under the item's text.
type List struct {
	items       []string
	bullet      string
	bulletStyle *Style
	style       *Style
	paddingLeft int
}

// listLine is a single rendered line of a List: unstyled indent, an optional
// styled bullet with its separating space, and styled text.
type listLine struct {
	indent string
	bullet string
	sep    string
	text   string
}

func (ln listLine) plain() string {
	return ln.indent + ln.bullet + ln.sep + ln.text
}

func NewList(items ...string) *List {
	return (&List{bullet: DefaultBullet}).WithItem(items...)
}

func (l *List) Items() []string {
	return l.items
}

// WithItem appends the given items to the list.
func (l *List) WithItem(items ...string) *List {
	l.items = append(l.items, items...)
	return l
}

// Clear removes all of the list's items.
func (l *List) Clear() *List {
	l.items = nil
	return l
}

func (l *List) Bullet() string {
	return l.bullet
}

// WithBullet sets the character (or string) rendered before each item.
func (l *List) WithBullet(bullet string) *List {
	l.bullet = bullet
	return l
}

func (l *List) BulletStyle() *Style {
	return l.bulletStyle
}

// WithBulletStyle sets the style applied to each item's bullet when
// rendering.
func (l *List) WithBulletStyle(s Style) *List {
	l.bulletStyle = &s
	return l
}

// WithoutBulletStyle removes the list's bullet style, rendering bullets
// plain.
func (l *List) WithoutBulletStyle() *List {
	l.bulletStyle = nil
	return l
}

func (l *List) Style() *Style {
	return l.style
}

// WithStyle sets the style applied to each item's text when rendering. The
// bullet is styled independently; see WithBulletStyle.
func (l *List) WithStyle(s Style) *List {
	l.style = &s
	return l
}

// WithoutStyle removes the list's text style, rendering item text plain.
func (l *List) WithoutStyle() *List {
	l.style = nil
	return l
}

func (l *List) PaddingLeft() int {
	return l.paddingLeft
}

// WithPaddingLeft sets the number of blank columns rendered before each
// item's bullet, pushing the whole list to the right.
func (l *List) WithPaddingLeft(n int) *List {
	l.paddingLeft = n
	return l
}

func (l *List) PreferredWidth() int {
	prefix := l.prefixWidth()
	w := 0
	for _, item := range l.items {
		w = max(w, prefix+StringWidth(item))
	}
	return w + l.paddingLeft
}

func (l *List) PreferredHeight(width int) int {
	h := 0
	for _, item := range l.items {
		h += len(l.wrapItem(item, width))
	}
	return h
}

func (l *List) Update(e Event) Event {
	return e
}

// Render wraps each item's text to fit width, prefixing the first line with
// the bullet and indenting continuation lines to align under the text. If
// the combined lines don't fit within height, they are truncated and the
// last visible line ends with an ellipsis. The bullet and text are styled
// independently.
func (l *List) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	var items []listLine
	for _, item := range l.items {
		items = append(items, l.wrapItem(item, width)...)
	}

	if len(items) > height {
		items = items[:height]
		ellipsizeListLine(&items[height-1], width)
	}

	lines := make([]string, len(items))
	for i, ln := range items {
		bullet := ln.bullet
		if l.bulletStyle != nil && bullet != "" {
			bullet = l.bulletStyle.Render(bullet)
		}
		text := ln.text
		if l.style != nil && text != "" {
			text = l.style.Render(text)
		}
		lines[i] = ln.indent + bullet + ln.sep + text
	}
	return lines
}

// ellipsizeListLine truncates ln's text so its plain rendering fits within
// width, ending with an ellipsis, while preserving its indent and bullet.
func ellipsizeListLine(ln *listLine, width int) {
	prefix := []rune(ln.indent + ln.bullet + ln.sep)
	ellipsized := []rune(ellipsize(ln.plain(), width))
	if len(ellipsized) < len(prefix) {
		*ln = listLine{bullet: string(ellipsized)}
		return
	}
	ln.text = string(ellipsized[len(prefix):])
}

// prefixWidth returns the number of columns occupied by the bullet plus the
// single space separating it from an item's text.
func (l *List) prefixWidth() int {
	return StringWidth(l.bullet) + 1
}

// wrapItem wraps a single item's text into lines of at most width columns,
// prefixing the first with the left padding and bullet, and indenting the
// rest to align under the text.
func (l *List) wrapItem(item string, width int) []listLine {
	pad := strings.Repeat(" ", l.paddingLeft)
	prefix := l.prefixWidth()
	textWidth := width - l.paddingLeft - prefix
	if textWidth <= 0 {
		return []listLine{{bullet: ellipsize(pad+l.bullet, width)}}
	}

	wrapped := wrapText(item, textWidth)
	indent := pad + strings.Repeat(" ", prefix)
	lines := make([]listLine, len(wrapped))
	for i, w := range wrapped {
		if i == 0 {
			lines[i] = listLine{indent: pad, bullet: l.bullet, sep: " ", text: w}
		} else {
			lines[i] = listLine{indent: indent, text: w}
		}
	}
	return lines
}
