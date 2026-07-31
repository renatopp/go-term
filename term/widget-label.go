package term

import "strings"

type Label struct {
	text  string
	width int
	style *Style
}

func NewLabel(text string) *Label {
	return (&Label{}).WithText(text)
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) WithText(text string) *Label {
	l.text = text
	l.width = stringWidth(l.text)
	return l
}

func (l *Label) Style() *Style {
	return l.style
}

// WithStyle sets the style applied to each line of the label when rendering.
func (l *Label) WithStyle(s Style) *Label {
	l.style = &s
	return l
}

// WithoutStyle removes the label's style, rendering plain text.
func (l *Label) WithoutStyle() *Label {
	l.style = nil
	return l
}

func (l *Label) PreferredWidth() int {
	return l.width
}

// PreferredHeight returns the number of lines needed to render the label's
// text wrapped greedily at word boundaries within width columns. Words wider
// than width are broken across multiple lines.
func (l *Label) PreferredHeight(width int) int {
	return len(wrapText(l.text, width))
}

func (l *Label) Update(e Event) Event {
	return e
}

// Render wraps the label's text to fit width, greedily breaking at word
// boundaries (and splitting words longer than width). If the wrapped text
// doesn't fit within height, it is truncated and the last visible line ends
// with an ellipsis.
func (l *Label) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	lines := wrapText(l.text, width)
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

// wrapText greedily wraps text into lines of at most width columns, breaking
// at word boundaries and splitting words wider than width across lines.
func wrapText(text string, width int) []string {
	if width <= 0 {
		return nil
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var line []rune
	lineWidth := 0
	for _, word := range words {
		w := []rune(word)
		wordWidth := runesWidth(w)
		for wordWidth > width {
			if len(line) > 0 {
				lines = append(lines, string(line))
				line = nil
				lineWidth = 0
			}
			n, cols := splitWidth(w, width)
			lines = append(lines, string(w[:n]))
			w = w[n:]
			wordWidth -= cols
		}

		switch {
		case len(line) == 0:
			line = w
			lineWidth = wordWidth
		case lineWidth+1+wordWidth <= width:
			line = append(append(line, ' '), w...)
			lineWidth += 1 + wordWidth
		default:
			lines = append(lines, string(line))
			line = w
			lineWidth = wordWidth
		}
	}

	return append(lines, string(line))
}

// ellipsize truncates s to fit within width columns, replacing the trailing
// content with an ellipsis when s doesn't already fit.
func ellipsize(s string, width int) string {
	if width <= 0 {
		return ""
	}

	runes := []rune(s)
	w := runesWidth(runes)
	if w < width {
		return string(runes) + "…"
	}
	if width == 1 {
		return "…"
	}
	for len(runes) > 0 && w > width-1 {
		w -= runeWidth(runes[len(runes)-1])
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}

// splitWidth returns the number of leading runes of w that fit within width
// columns (at least one, so oversized runes still make progress), along with
// their total column width.
func splitWidth(w []rune, width int) (n, cols int) {
	for n < len(w) {
		rw := runeWidth(w[n])
		if n > 0 && cols+rw > width {
			break
		}
		cols += rw
		n++
	}
	return n, cols
}

// runesWidth returns the number of terminal columns runes occupy.
func runesWidth(runes []rune) int {
	w := 0
	for _, r := range runes {
		w += runeWidth(r)
	}
	return w
}
