package term

import (
	"strings"
	"testing"

	"github.com/renatopp/go-term/term/ui"
)

func TestNewProgressBar(t *testing.T) {
	p := NewProgressBar()
	if p.Value() != 0 {
		t.Fatalf("Value() = %v, want 0", p.Value())
	}
	if p.Width() != DefaultProgressBarWidth {
		t.Fatalf("Width() = %d, want %d", p.Width(), DefaultProgressBarWidth)
	}
	if p.FilledChar() != DefaultProgressBarFilledChar {
		t.Fatalf("FilledChar() = %q, want %q", p.FilledChar(), DefaultProgressBarFilledChar)
	}
	if p.EmptyChar() != DefaultProgressBarEmptyChar {
		t.Fatalf("EmptyChar() = %q, want %q", p.EmptyChar(), DefaultProgressBarEmptyChar)
	}
}

func TestProgressBarWithValue(t *testing.T) {
	p := NewProgressBar()
	same := p.WithValue(0.5)
	if same != p {
		t.Fatal("WithValue should return the same *ProgressBar for chaining")
	}
	if p.Value() != 0.5 {
		t.Fatalf("Value() = %v, want 0.5", p.Value())
	}
}

func TestProgressBarWithValueClampsBelowZero(t *testing.T) {
	p := NewProgressBar().WithValue(-1)
	if p.Value() != 0 {
		t.Fatalf("Value() = %v, want 0", p.Value())
	}
}

func TestProgressBarWithValueClampsAboveOne(t *testing.T) {
	p := NewProgressBar().WithValue(2)
	if p.Value() != 1 {
		t.Fatalf("Value() = %v, want 1", p.Value())
	}
}

func TestProgressBarWithWidth(t *testing.T) {
	p := NewProgressBar().WithWidth(10)
	if p.Width() != 10 {
		t.Fatalf("Width() = %d, want 10", p.Width())
	}
}

func TestProgressBarWithWidthClampsNegative(t *testing.T) {
	p := NewProgressBar().WithWidth(-5)
	if p.Width() != 0 {
		t.Fatalf("Width() = %d, want 0", p.Width())
	}
}

func TestProgressBarWithFilledChar(t *testing.T) {
	p := NewProgressBar().WithFilledChar("#")
	if p.FilledChar() != "#" {
		t.Fatalf("FilledChar() = %q, want %q", p.FilledChar(), "#")
	}
}

func TestProgressBarWithEmptyChar(t *testing.T) {
	p := NewProgressBar().WithEmptyChar("-")
	if p.EmptyChar() != "-" {
		t.Fatalf("EmptyChar() = %q, want %q", p.EmptyChar(), "-")
	}
}

func TestProgressBarWithChars(t *testing.T) {
	p := NewProgressBar()
	same := p.WithChars("#", "-")
	if same != p {
		t.Fatal("WithChars should return the same *ProgressBar for chaining")
	}
	if p.FilledChar() != "#" {
		t.Fatalf("FilledChar() = %q, want %q", p.FilledChar(), "#")
	}
	if p.EmptyChar() != "-" {
		t.Fatalf("EmptyChar() = %q, want %q", p.EmptyChar(), "-")
	}
}

func TestProgressBarWithPrefix(t *testing.T) {
	p := NewProgressBar().WithPrefix("Loading")
	if p.Prefix() != "Loading" {
		t.Fatalf("Prefix() = %q, want %q", p.Prefix(), "Loading")
	}
}

func TestProgressBarWithoutPrefixStyle(t *testing.T) {
	p := NewProgressBar().WithPrefixStyle(NewStyle().WithForeground(ColorRed)).WithoutPrefixStyle()
	if p.PrefixStyle() != nil {
		t.Fatal("WithoutPrefixStyle should clear the style")
	}
}

func TestProgressBarWithSuffix(t *testing.T) {
	p := NewProgressBar().WithSuffix("Done")
	if p.Suffix() != "Done" {
		t.Fatalf("Suffix() = %q, want %q", p.Suffix(), "Done")
	}
}

func TestProgressBarWithoutSuffixStyle(t *testing.T) {
	p := NewProgressBar().WithSuffixStyle(NewStyle().WithForeground(ColorRed)).WithoutSuffixStyle()
	if p.SuffixStyle() != nil {
		t.Fatal("WithoutSuffixStyle should clear the style")
	}
}

func TestProgressBarWithoutPercentStyle(t *testing.T) {
	p := NewProgressBar().WithPercentStyle(NewStyle().WithForeground(ColorRed)).WithoutPercentStyle()
	if p.PercentStyle() != nil {
		t.Fatal("WithoutPercentStyle should clear the style")
	}
}

func TestProgressBarAsShowPercent(t *testing.T) {
	p := NewProgressBar()
	same := p.AsShowPercent(true)
	if same != p {
		t.Fatal("AsShowPercent should return the same *ProgressBar for chaining")
	}
	if !p.ShowPercent() {
		t.Fatal("ShowPercent() = false, want true")
	}
	p.AsShowPercent(false)
	if p.ShowPercent() {
		t.Fatal("ShowPercent() = true, want false")
	}
}

func TestProgressBarWithoutStyle(t *testing.T) {
	p := NewProgressBar().WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if p.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestProgressBarWithoutEmptyStyle(t *testing.T) {
	p := NewProgressBar().WithEmptyStyle(NewStyle().WithForeground(ColorRed)).WithoutEmptyStyle()
	if p.EmptyStyle() != nil {
		t.Fatal("WithoutEmptyStyle should clear the style")
	}
}

func TestProgressBarPreferredWidth(t *testing.T) {
	p := NewProgressBar().WithWidth(10)
	if w := p.PreferredWidth(); w != 10 {
		t.Fatalf("PreferredWidth() = %d, want 10", w)
	}
}

func TestProgressBarPreferredWidthWithPrefix(t *testing.T) {
	p := NewProgressBar().WithWidth(10).WithPrefix("Loading")
	// "Loading" (7, touching the bar) + bar 10 = 17
	if w := p.PreferredWidth(); w != 17 {
		t.Fatalf("PreferredWidth() = %d, want 17", w)
	}
}

func TestProgressBarPreferredWidthWithSuffix(t *testing.T) {
	p := NewProgressBar().WithWidth(10).WithSuffix("Done")
	// bar 10 + "Done" (4, touching the bar) = 14
	if w := p.PreferredWidth(); w != 14 {
		t.Fatalf("PreferredWidth() = %d, want 14", w)
	}
}

func TestProgressBarPreferredWidthWithPercent(t *testing.T) {
	p := NewProgressBar().WithWidth(10).AsShowPercent(true)
	// bar 10 + " " (1) + "100%" (4) = 15
	if w := p.PreferredWidth(); w != 15 {
		t.Fatalf("PreferredWidth() = %d, want 15", w)
	}
}

func TestProgressBarPreferredWidthWithPrefixSuffixAndPercent(t *testing.T) {
	p := NewProgressBar().WithWidth(10).WithPrefix("Loading").WithSuffix("Done").AsShowPercent(true)
	// "Loading" (7) + bar 10 + "Done" (4) + " " (1) + "100%" (4) = 26
	if w := p.PreferredWidth(); w != 26 {
		t.Fatalf("PreferredWidth() = %d, want 26", w)
	}
}

func TestProgressBarPreferredHeight(t *testing.T) {
	p := NewProgressBar()
	if h := p.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight() = %d, want 1", h)
	}
}

func TestProgressBarUpdateReturnsEventUnchanged(t *testing.T) {
	p := NewProgressBar()
	e := KeyEvent{Rune: 'x'}
	if got := p.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestProgressBarRenderEmpty(t *testing.T) {
	p := NewProgressBar().WithWidth(4)
	lines := p.Render(4, 1)
	want := []string{"░░░░"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderHalfFilled(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(0.5)
	lines := p.Render(4, 1)
	want := []string{"██░░"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderFull(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(1)
	lines := p.Render(4, 1)
	want := []string{"████"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderWithCustomChars(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(0.5).WithFilledChar("#").WithEmptyChar("-")
	lines := p.Render(4, 1)
	want := []string{"##--"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderWithPrefix(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(1).WithPrefix("Load")
	lines := p.Render(9, 1)
	want := []string{"Load████ "}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderWithSuffix(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(1).WithSuffix("Done")
	lines := p.Render(9, 1)
	want := []string{"████Done "}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderWithPercent(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(0.5).AsShowPercent(true)
	lines := p.Render(9, 1)
	want := []string{"██░░  50%"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderWithPrefixSuffixAndPercent(t *testing.T) {
	p := NewProgressBar().WithWidth(4).WithValue(0.5).WithPrefix("Load").WithSuffix("Done").AsShowPercent(true)
	lines := p.Render(19, 1)
	want := []string{"Load██░░Done  50%  "}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderDefaultWidthReservesSpaceForTextAndPercent(t *testing.T) {
	// With the default (effectively unbounded) width, the bar must still
	// leave room for the prefix, percentage, and suffix instead of consuming
	// the whole render width and truncating them off; the leftover space
	// (after reserving for text and percent) all goes to the track.
	p := NewProgressBar().WithValue(0.5).WithPrefix("Load").WithSuffix("Done").AsShowPercent(true)
	lines := p.Render(19, 1)
	want := []string{"Load███░░░Done  50%"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderClipsToWidth(t *testing.T) {
	p := NewProgressBar().WithWidth(10).WithValue(1)
	lines := p.Render(4, 1)
	want := []string{"████"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderTrackShrinksToRenderWidthPreservingProportion(t *testing.T) {
	// Configured width (10) exceeds the render width (4), so the track
	// itself shrinks to 4 cells and 50% still renders as half-filled,
	// rather than the first 4 of 10 cells (which would show as full).
	p := NewProgressBar().WithWidth(10).WithValue(0.5)
	lines := p.Render(4, 1)
	want := []string{"██░░"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarDefaultWidthFillsRenderWidth(t *testing.T) {
	// An unconfigured bar's width (DefaultProgressBarWidth) is effectively
	// unbounded, so it fills whatever width it's rendered with.
	p := NewProgressBar().WithValue(0.5)
	lines := p.Render(4, 1)
	want := []string{"██░░"}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderPadsToWidth(t *testing.T) {
	p := NewProgressBar().WithWidth(2).WithValue(1)
	lines := p.Render(5, 1)
	want := []string{"██   "}
	if len(lines) != 1 || lines[0] != want[0] {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
}

func TestProgressBarRenderZeroSize(t *testing.T) {
	p := NewProgressBar()
	if lines := p.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := p.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestProgressBarRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	p := NewProgressBar().WithWidth(4).WithValue(0.5).WithStyle(NewStyle().WithForeground(ColorRed))
	lines := p.Render(4, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled fill to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 4 {
		t.Fatalf("StringWidth(lines[0]) = %d, want 4", w)
	}
}

func TestProgressBarRenderWithEmptyStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	p := NewProgressBar().WithWidth(4).WithValue(0.5).WithEmptyStyle(NewStyle().WithForeground(ColorRed))
	lines := p.Render(4, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled track to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 4 {
		t.Fatalf("StringWidth(lines[0]) = %d, want 4", w)
	}
}

func TestProgressBarRenderWithPrefixStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	p := NewProgressBar().WithWidth(4).WithPrefix("Load").WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	lines := p.Render(9, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled prefix to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 9 {
		t.Fatalf("StringWidth(lines[0]) = %d, want 9", w)
	}
}

func TestProgressBarRenderWithSuffixStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	p := NewProgressBar().WithWidth(4).WithSuffix("Done").WithSuffixStyle(NewStyle().WithForeground(ColorRed))
	lines := p.Render(9, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled suffix to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 9 {
		t.Fatalf("StringWidth(lines[0]) = %d, want 9", w)
	}
}

func TestProgressBarRenderWithPercentStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	p := NewProgressBar().WithWidth(4).WithValue(0.5).AsShowPercent(true).WithPercentStyle(NewStyle().WithForeground(ColorRed))
	lines := p.Render(9, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled percent to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 9 {
		t.Fatalf("StringWidth(lines[0]) = %d, want 9", w)
	}
}
