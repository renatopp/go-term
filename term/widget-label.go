package term

import "github.com/renatopp/go-term/term/ui"

type Label struct {
	text string
}

func NewLabel(text string) *Label {
	return &Label{text: text}
}

var (
	_ ui.Component = (*Label)(nil)
	_ ui.Sizeable  = (*Label)(nil)
)

func (l *Label) WithText(text string) *Label {
	l.text = text
	return l
}

func (l *Label) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}
	runes := []rune(l.text)
	if len(runes) > width {
		runes = runes[:width]
	}
	return []string{string(runes)}
}

func (l *Label) PreferredWidth() int {
	return len([]rune(l.text))
}

func (l *Label) PreferredHeight(width int) int {
	return 1
}
