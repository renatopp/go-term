package term

import (
	"strings"
	"time"
)

// DefaultSpinnerFrames are the frames used by a Spinner until WithFrames is
// called: a Braille-based spinner animation.
var DefaultSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// SpinnerFramesLine is a classic ASCII line spinner.
var SpinnerFramesLine = []string{"-", "\\", "|", "/"}

// SpinnerFramesCircle cycles through quarter-filled circles.
var SpinnerFramesCircle = []string{"◐", "◓", "◑", "◒"}

// SpinnerFramesClock cycles through quarter-turn clock faces.
var SpinnerFramesClock = []string{"◴", "◷", "◶", "◵"}

// SpinnerFramesSquare cycles through quarter-filled squares.
var SpinnerFramesSquare = []string{"◰", "◳", "◲", "◱"}

// SpinnerFramesArc cycles through arc segments.
var SpinnerFramesArc = []string{"◜", "◠", "◝", "◞", "◡", "◟"}

// SpinnerFramesDots2 is an alternate Braille-based spinner animation.
var SpinnerFramesDots2 = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

// SpinnerFramesDotsOrbit cycles a single Braille dot around its cell.
var SpinnerFramesDotsOrbit = []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}

// SpinnerFramesEllipsis shows a growing then resetting sequence of dots.
var SpinnerFramesEllipsis = []string{".  ", ".. ", "...", "   "}

// SpinnerFramesBounceBar bounces a block up and down through its height levels.
var SpinnerFramesBounceBar = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▁"}

// SpinnerFramesQuadrant cycles through single quadrant blocks.
var SpinnerFramesQuadrant = []string{"▖", "▘", "▝", "▗"}

// SpinnerFramesSquareToggle toggles between filled and empty squares.
var SpinnerFramesSquareToggle = []string{"■", "□", "▪", "▫"}

// SpinnerFramesArrow cycles through the 8 compass-direction arrows.
var SpinnerFramesArrow = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

// SpinnerFramesEarth cycles through the rotating earth emoji.
var SpinnerFramesEarth = []string{"🌍", "🌎", "🌏"}

// SpinnerFramesMoon cycles through the moon phase emoji.
var SpinnerFramesMoon = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}

// DefaultSpinnerFrameRate is the interval between frame advances used by a
// Spinner until WithFrameRate is called.
const DefaultSpinnerFrameRate = 100 * time.Millisecond

// Spinner renders one of a sequence of frames, optionally preceded by a
// prefix and followed by a suffix, advanced by calling Tick. It does not
// animate on its own; call Tick whenever a new frame should be shown (e.g.
// from a program-driven timer). Update, in contrast, paces advancement to
// FrameRate regardless of how often TickEvent fires.
type Spinner struct {
	frames      []string
	frame       int
	style       *Style
	prefix      string
	prefixStyle *Style
	suffix      string
	suffixStyle *Style
	frameRate   time.Duration
	elapsed     time.Duration
}

func NewSpinner() *Spinner {
	return &Spinner{frames: DefaultSpinnerFrames, frameRate: DefaultSpinnerFrameRate}
}

func (s *Spinner) Frames() []string {
	return s.frames
}

// WithFrames sets the sequence of frames the spinner cycles through and
// resets it to the first one.
func (s *Spinner) WithFrames(frames ...string) *Spinner {
	s.frames = frames
	s.frame = 0
	return s
}

// Frame returns the text of the spinner's current frame.
func (s *Spinner) Frame() string {
	if len(s.frames) == 0 {
		return ""
	}
	return s.frames[s.frame]
}

// Tick advances the spinner to its next frame, wrapping around at the end.
func (s *Spinner) Tick() *Spinner {
	if len(s.frames) == 0 {
		return s
	}
	s.frame = (s.frame + 1) % len(s.frames)
	return s
}

func (s *Spinner) Style() *Style {
	return s.style
}

// WithStyle sets the style applied to the spinner's current frame when
// rendering.
func (s *Spinner) WithStyle(style Style) *Spinner {
	s.style = &style
	return s
}

// WithoutStyle removes the spinner's frame style, rendering it plain.
func (s *Spinner) WithoutStyle() *Spinner {
	s.style = nil
	return s
}

func (s *Spinner) Prefix() string {
	return s.prefix
}

// WithPrefix sets the text rendered before the spinner's frame, separated by
// a single space.
func (s *Spinner) WithPrefix(text string) *Spinner {
	s.prefix = text
	return s
}

func (s *Spinner) PrefixStyle() *Style {
	return s.prefixStyle
}

// WithPrefixStyle sets the style applied to the spinner's prefix when
// rendering.
func (s *Spinner) WithPrefixStyle(style Style) *Spinner {
	s.prefixStyle = &style
	return s
}

// WithoutPrefixStyle removes the spinner's prefix style, rendering it plain.
func (s *Spinner) WithoutPrefixStyle() *Spinner {
	s.prefixStyle = nil
	return s
}

func (s *Spinner) Suffix() string {
	return s.suffix
}

// WithSuffix sets the text rendered after the spinner's frame, separated by
// a single space.
func (s *Spinner) WithSuffix(text string) *Spinner {
	s.suffix = text
	return s
}

func (s *Spinner) SuffixStyle() *Style {
	return s.suffixStyle
}

// WithSuffixStyle sets the style applied to the spinner's suffix when
// rendering.
func (s *Spinner) WithSuffixStyle(style Style) *Spinner {
	s.suffixStyle = &style
	return s
}

// WithoutSuffixStyle removes the spinner's suffix style, rendering it plain.
func (s *Spinner) WithoutSuffixStyle() *Spinner {
	s.suffixStyle = nil
	return s
}

func (s *Spinner) FrameRate() time.Duration {
	return s.frameRate
}

// WithFrameRate sets the interval between frame advances applied by Update.
// A rate <= 0 stops Update from advancing the spinner; Tick is unaffected.
func (s *Spinner) WithFrameRate(d time.Duration) *Spinner {
	s.frameRate = d
	s.elapsed = 0
	return s
}

// PreferredWidth returns the widest of the spinner's frames plus its prefix
// and suffix (whichever are set), each with their separating space, so its
// size doesn't jitter as it animates.
func (s *Spinner) PreferredWidth() int {
	w := 0
	for _, f := range s.frames {
		w = max(w, StringWidth(f))
	}
	if s.prefix != "" {
		w += StringWidth(s.prefix) + 1
	}
	if s.suffix != "" {
		w += StringWidth(s.suffix) + 1
	}
	return w
}

func (s *Spinner) PreferredHeight(width int) int {
	return 1
}

// Update advances the spinner on TickEvent (see Program.WithTick), pacing
// frame changes to FrameRate rather than the raw tick rate, accumulating
// leftover time so drift doesn't accumulate.
func (s *Spinner) Update(e Event) Event {
	if t, ok := e.(TickEvent); ok && s.frameRate > 0 {
		s.elapsed += t.Duration
		if ticks := int(s.elapsed / s.frameRate); ticks > 0 {
			s.elapsed -= time.Duration(ticks) * s.frameRate
			for range ticks {
				s.Tick()
			}
		}
	}
	return e
}

// Render draws the spinner's prefix (if set), its current frame, and its
// suffix (if set), clipped to width columns. Each part keeps its own style;
// the spaces separating them are left unstyled.
func (s *Spinner) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	type segment struct {
		runes []rune
		style *Style
	}

	var segs []segment
	if s.prefix != "" {
		segs = append(segs, segment{[]rune(s.prefix), s.prefixStyle}, segment{[]rune(" "), nil})
	}
	segs = append(segs, segment{[]rune(s.Frame()), s.style})
	if s.suffix != "" {
		segs = append(segs, segment{[]rune(" "), nil}, segment{[]rune(s.suffix), s.suffixStyle})
	}

	var full []rune
	for _, seg := range segs {
		full = append(full, seg.runes...)
	}
	n, _ := splitWidth(full, width)

	var line strings.Builder
	remaining := n
	for _, seg := range segs {
		if remaining <= 0 {
			break
		}
		take := min(len(seg.runes), remaining)
		text := string(seg.runes[:take])
		if seg.style != nil && take > 0 {
			text = seg.style.Render(text)
		}
		line.WriteString(text)
		remaining -= take
	}

	return []string{line.String()}
}
