package term

import (
	"strings"

	"github.com/renatopp/go-term/term/ui"
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
	style       *Style
	paddingLeft int
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

func (l *List) Style() *Style {
	return l.style
}

// WithStyle sets the style applied to each rendered line, including its
// bullet.
func (l *List) WithStyle(s Style) *List {
	l.style = &s
	return l
}

// WithoutStyle removes the list's style, rendering plain text.
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
		w = max(w, prefix+ui.StringWidth(item))
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
// last visible line ends with an ellipsis.
func (l *List) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	var lines []string
	for _, item := range l.items {
		lines = append(lines, l.wrapItem(item, width)...)
	}

	if len(lines) > height {
		lines = lines[:height]
		lines[height-1] = ellipsize(lines[height-1], width)
	}

	if l.style != nil {
		for i, line := range lines {
			lines[i] = l.style.Render(line)
		}
	}
	return lines
}

// prefixWidth returns the number of columns occupied by the bullet plus the
// single space separating it from an item's text.
func (l *List) prefixWidth() int {
	return ui.StringWidth(l.bullet) + 1
}

// wrapItem wraps a single item's text into lines of at most width columns,
// prefixing the first with the left padding and bullet, and indenting the
// rest to align under the text.
func (l *List) wrapItem(item string, width int) []string {
	pad := strings.Repeat(" ", l.paddingLeft)
	prefix := l.prefixWidth()
	textWidth := width - l.paddingLeft - prefix
	if textWidth <= 0 {
		return []string{ellipsize(pad+l.bullet, width)}
	}

	wrapped := wrapText(item, textWidth)
	indent := pad + strings.Repeat(" ", prefix)
	lines := make([]string, len(wrapped))
	for i, w := range wrapped {
		if i == 0 {
			lines[i] = pad + l.bullet + " " + w
		} else {
			lines[i] = indent + w
		}
	}
	return lines
}
