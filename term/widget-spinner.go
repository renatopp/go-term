package term

import (
	"time"

	"github.com/renatopp/go-term/term/ui"
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

// Spinner renders one of a sequence of frames, advanced by calling Tick.
// It does not animate on its own; call Tick whenever a new frame should be
// shown (e.g. from a program-driven timer). Update, in contrast, paces
// advancement to FrameRate regardless of how often TickEvent fires.
type Spinner struct {
	frames    []string
	frame     int
	text      string
	style     *Style
	frameRate time.Duration
	elapsed   time.Duration
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

func (s *Spinner) Text() string {
	return s.text
}

// WithText sets the text rendered after the spinner's frame, separated by a
// single space.
func (s *Spinner) WithText(text string) *Spinner {
	s.text = text
	return s
}

func (s *Spinner) Style() *Style {
	return s.style
}

// WithStyle sets the style applied to the spinner's current frame and text
// when rendering.
func (s *Spinner) WithStyle(style Style) *Spinner {
	s.style = &style
	return s
}

// WithoutStyle removes the spinner's style, rendering plain text.
func (s *Spinner) WithoutStyle() *Spinner {
	s.style = nil
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

// PreferredWidth returns the widest of the spinner's frames plus its text,
// so its size doesn't jitter as it animates.
func (s *Spinner) PreferredWidth() int {
	w := 0
	for _, f := range s.frames {
		w = max(w, ui.StringWidth(f))
	}
	if s.text != "" {
		w += 1 + ui.StringWidth(s.text)
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

// Render draws the spinner's current frame followed by its text (if set),
// clipped to width columns.
func (s *Spinner) Render(width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	full := s.Frame()
	if s.text != "" {
		full += " " + s.text
	}

	runes := []rune(full)
	n, _ := splitWidth(runes, width)
	line := string(runes[:n])

	if s.style != nil {
		line = s.style.Render(line)
	}
	return []string{line}
}
