package term

import (
	"fmt"
	"math"
	"strings"

	"github.com/renatopp/go-term/term/ui"
)

// DefaultProgressBarWidth is used until WithWidth is called. It's large
// enough to be effectively unbounded, so by default the bar's track fills
// whatever width it's rendered with rather than a fixed size.
const DefaultProgressBarWidth = 9999

// DefaultProgressBarFilledChar is the character used to draw the filled
// portion of a ProgressBar until WithFilledChar is called.
const DefaultProgressBarFilledChar = "█"

// DefaultProgressBarEmptyChar is the character used to draw the empty
// portion of a ProgressBar until WithEmptyChar is called.
const DefaultProgressBarEmptyChar = "░"

// ProgressBar renders a horizontal bar filled in proportion to a value
// between 0 and 1, optionally preceded by a prefix (touching the bar) and
// followed by a suffix (also touching the bar) and a percentage. It does not
// animate on its own; call WithValue whenever progress changes.
type ProgressBar struct {
	value        float64
	width        int
	filledChar   string
	emptyChar    string
	prefix       string
	suffix       string
	showPercent  bool
	style        *Style
	emptyStyle   *Style
	prefixStyle  *Style
	suffixStyle  *Style
	percentStyle *Style
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		width:      DefaultProgressBarWidth,
		filledChar: DefaultProgressBarFilledChar,
		emptyChar:  DefaultProgressBarEmptyChar,
	}
}

func (p *ProgressBar) Value() float64 {
	return p.value
}

// WithValue sets the bar's progress, clamped to [0, 1].
func (p *ProgressBar) WithValue(v float64) *ProgressBar {
	switch {
	case v < 0:
		v = 0
	case v > 1:
		v = 1
	}
	p.value = v
	return p
}

func (p *ProgressBar) Width() int {
	return p.width
}

// WithWidth sets the number of cells the bar track occupies, clamped to a
// minimum of 0.
func (p *ProgressBar) WithWidth(n int) *ProgressBar {
	if n < 0 {
		n = 0
	}
	p.width = n
	return p
}

func (p *ProgressBar) FilledChar() string {
	return p.filledChar
}

// WithFilledChar sets the character rendered for each filled cell of the
// bar.
func (p *ProgressBar) WithFilledChar(ch string) *ProgressBar {
	p.filledChar = ch
	return p
}

func (p *ProgressBar) EmptyChar() string {
	return p.emptyChar
}

// WithEmptyChar sets the character rendered for each empty cell of the bar.
func (p *ProgressBar) WithEmptyChar(ch string) *ProgressBar {
	p.emptyChar = ch
	return p
}

// WithChars sets the characters rendered for each filled and empty cell of
// the bar in one call.
func (p *ProgressBar) WithChars(filled, empty string) *ProgressBar {
	p.filledChar = filled
	p.emptyChar = empty
	return p
}

func (p *ProgressBar) Prefix() string {
	return p.prefix
}

// WithPrefix sets the text rendered immediately before the bar, with no
// separating space.
func (p *ProgressBar) WithPrefix(text string) *ProgressBar {
	p.prefix = text
	return p
}

func (p *ProgressBar) PrefixStyle() *Style {
	return p.prefixStyle
}

// WithPrefixStyle sets the style applied to the bar's prefix when rendering.
func (p *ProgressBar) WithPrefixStyle(s Style) *ProgressBar {
	p.prefixStyle = &s
	return p
}

// WithoutPrefixStyle removes the bar's prefix style, rendering it plain.
func (p *ProgressBar) WithoutPrefixStyle() *ProgressBar {
	p.prefixStyle = nil
	return p
}

func (p *ProgressBar) Suffix() string {
	return p.suffix
}

// WithSuffix sets the text rendered immediately after the bar (and before
// its percentage, if shown), with no separating space.
func (p *ProgressBar) WithSuffix(text string) *ProgressBar {
	p.suffix = text
	return p
}

func (p *ProgressBar) SuffixStyle() *Style {
	return p.suffixStyle
}

// WithSuffixStyle sets the style applied to the bar's suffix when rendering.
func (p *ProgressBar) WithSuffixStyle(s Style) *ProgressBar {
	p.suffixStyle = &s
	return p
}

// WithoutSuffixStyle removes the bar's suffix style, rendering it plain.
func (p *ProgressBar) WithoutSuffixStyle() *ProgressBar {
	p.suffixStyle = nil
	return p
}

func (p *ProgressBar) ShowPercent() bool {
	return p.showPercent
}

// AsShowPercent sets whether the bar's rounded percentage is rendered after
// it, separated by a single space.
func (p *ProgressBar) AsShowPercent(v bool) *ProgressBar {
	p.showPercent = v
	return p
}

func (p *ProgressBar) PercentStyle() *Style {
	return p.percentStyle
}

// WithPercentStyle sets the style applied to the bar's percentage when
// rendering.
func (p *ProgressBar) WithPercentStyle(s Style) *ProgressBar {
	p.percentStyle = &s
	return p
}

// WithoutPercentStyle removes the bar's percentage style, rendering it
// plain.
func (p *ProgressBar) WithoutPercentStyle() *ProgressBar {
	p.percentStyle = nil
	return p
}

func (p *ProgressBar) Style() *Style {
	return p.style
}

// WithStyle sets the style applied to the bar's filled portion when
// rendering.
func (p *ProgressBar) WithStyle(s Style) *ProgressBar {
	p.style = &s
	return p
}

// WithoutStyle removes the bar's filled-portion style, rendering it plain.
func (p *ProgressBar) WithoutStyle() *ProgressBar {
	p.style = nil
	return p
}

func (p *ProgressBar) EmptyStyle() *Style {
	return p.emptyStyle
}

// WithEmptyStyle sets the style applied to the bar's empty portion when
// rendering.
func (p *ProgressBar) WithEmptyStyle(s Style) *ProgressBar {
	p.emptyStyle = &s
	return p
}

// WithoutEmptyStyle removes the bar's empty-portion style, rendering it
// plain.
func (p *ProgressBar) WithoutEmptyStyle() *ProgressBar {
	p.emptyStyle = nil
	return p
}

// PreferredWidth returns the bar's track width plus its prefix and suffix
// (touching the bar, so no separating space) and its percentage (whichever
// are set), which keeps a separating space before it.
func (p *ProgressBar) PreferredWidth() int {
	w := p.width
	if p.prefix != "" {
		w += ui.StringWidth(p.prefix)
	}
	if p.suffix != "" {
		w += ui.StringWidth(p.suffix)
	}
	if p.showPercent {
		w += 5 // " " + "100%"
	}
	return w
}

func (p *ProgressBar) PreferredHeight(width int) int {
	return 1
}

func (p *ProgressBar) Update(e Event) Event {
	return e
}

// Render draws the bar's prefix (if set, touching the bar), its filled and
// empty portions, its suffix (if set, also touching the bar), and its
// percentage (if shown), clipped or padded to fit exactly width columns. The
// bar's track is sized to the smaller of its configured width and the render
// width left over once its prefix, suffix, and percentage (whichever are
// set) reserve their own space.
func (p *ProgressBar) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	reserved := 0
	if p.prefix != "" {
		reserved += ui.StringWidth(p.prefix)
	}
	if p.suffix != "" {
		reserved += ui.StringWidth(p.suffix)
	}
	if p.showPercent {
		reserved += 5 // " " + "100%"
	}

	barWidth := min(p.width, max(0, width-reserved))
	filled := max(0, min(barWidth, int(math.Round(p.value*float64(barWidth)))))
	empty := barWidth - filled

	bar := strings.Repeat(p.filledChar, filled)
	if p.style != nil {
		bar = p.style.Render(bar)
	}
	track := strings.Repeat(p.emptyChar, empty)
	if p.emptyStyle != nil {
		track = p.emptyStyle.Render(track)
	}

	var line strings.Builder
	if p.prefix != "" {
		prefix := p.prefix
		if p.prefixStyle != nil {
			prefix = p.prefixStyle.Render(prefix)
		}
		line.WriteString(prefix)
	}
	line.WriteString(bar)
	line.WriteString(track)
	if p.suffix != "" {
		suffix := p.suffix
		if p.suffixStyle != nil {
			suffix = p.suffixStyle.Render(suffix)
		}
		line.WriteString(suffix)
	}
	if p.showPercent {
		line.WriteString(" ")
		percent := fmt.Sprintf("%3d%%", int(math.Round(p.value*100)))
		if p.percentStyle != nil {
			percent = p.percentStyle.Render(percent)
		}
		line.WriteString(percent)
	}

	return []string{ui.RenderRow(ui.BuildRow(line.String(), width))}
}
