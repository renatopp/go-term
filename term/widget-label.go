package term

import (
	"strings"
	"unicode/utf8"
)

type Label struct {
	text  string
	width int
}

func NewLabel(text string) *Label {
	return (&Label{}).WithText(text)
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) WithText(text string) *Label {
	l.text = text
	l.width = utf8.RuneCountInString(l.text)
	return l
}

func (l *Label) PreferredWidth() int {
	return l.width
}

// PreferredHeight returns the number of lines needed to render the label's
// text wrapped greedily at word boundaries within width runes. Words longer
// than width are broken across multiple lines.
func (l *Label) PreferredHeight(width int) int {
	return len(wrapText(l.text, width))
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
	if len(lines) <= height {
		return lines
	}

	lines = lines[:height]
	lines[height-1] = ellipsize(lines[height-1], width)
	return lines
}

// wrapText greedily wraps text into lines of at most width runes, breaking
// at word boundaries and splitting words longer than width across lines.
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
	for _, word := range words {
		w := []rune(word)
		for len(w) > width {
			if len(line) > 0 {
				lines = append(lines, string(line))
				line = nil
			}
			lines = append(lines, string(w[:width]))
			w = w[width:]
		}

		switch {
		case len(line) == 0:
			line = w
		case len(line)+1+len(w) <= width:
			line = append(append(line, ' '), w...)
		default:
			lines = append(lines, string(line))
			line = w
		}
	}

	return append(lines, string(line))
}

// ellipsize truncates s to fit within width runes, replacing the trailing
// rune with an ellipsis when s doesn't already fit.
func ellipsize(s string, width int) string {
	if width <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) < width {
		return string(runes) + "…"
	}
	if width == 1 {
		return "…"
	}
	return string(runes[:width-1]) + "…"
}
