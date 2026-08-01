package term

import (
	"strings"
	"testing"
	"time"

	"github.com/renatopp/go-term/term/ui"
)

func TestNewSpinner(t *testing.T) {
	s := NewSpinner()
	if s.Frame() != DefaultSpinnerFrames[0] {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), DefaultSpinnerFrames[0])
	}
}

func TestSpinnerWithFrames(t *testing.T) {
	s := NewSpinner()
	same := s.WithFrames("a", "b", "c")
	if same != s {
		t.Fatal("WithFrames should return the same *Spinner for chaining")
	}
	if s.Frame() != "a" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "a")
	}
}

func TestSpinnerWithFramesResetsToFirstFrame(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b")
	s.Tick()
	if s.Frame() != "b" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "b")
	}
	s.WithFrames("x", "y")
	if s.Frame() != "x" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "x")
	}
}

func TestSpinnerUpdateAdvancesOnTickEvent(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c")
	s.Update(TickEvent{Duration: s.FrameRate()})
	if s.Frame() != "b" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "b")
	}
}

func TestSpinnerUpdateDoesNotAdvanceBelowFrameRate(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c")
	s.Update(TickEvent{Duration: s.FrameRate() - 1})
	if s.Frame() != "a" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "a")
	}
}

func TestSpinnerUpdateAccumulatesElapsedTime(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c").WithFrameRate(100 * time.Millisecond)
	s.Update(TickEvent{Duration: 60 * time.Millisecond})
	if s.Frame() != "a" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "a")
	}
	s.Update(TickEvent{Duration: 60 * time.Millisecond})
	if s.Frame() != "b" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "b")
	}
}

func TestSpinnerUpdateAdvancesMultipleFramesWhenLagging(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c").WithFrameRate(100 * time.Millisecond)
	s.Update(TickEvent{Duration: 250 * time.Millisecond})
	if s.Frame() != "c" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "c")
	}
}

func TestSpinnerUpdateIgnoresTickEventWhenFrameRateIsZero(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c").WithFrameRate(0)
	s.Update(TickEvent{Duration: time.Hour})
	if s.Frame() != "a" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "a")
	}
}

func TestSpinnerFrameRateDefault(t *testing.T) {
	s := NewSpinner()
	if s.FrameRate() != DefaultSpinnerFrameRate {
		t.Fatalf("FrameRate() = %v, want %v", s.FrameRate(), DefaultSpinnerFrameRate)
	}
}

func TestSpinnerWithFrameRate(t *testing.T) {
	s := NewSpinner()
	same := s.WithFrameRate(50 * time.Millisecond)
	if same != s {
		t.Fatal("WithFrameRate should return the same *Spinner for chaining")
	}
	if s.FrameRate() != 50*time.Millisecond {
		t.Fatalf("FrameRate() = %v, want %v", s.FrameRate(), 50*time.Millisecond)
	}
}

func TestSpinnerUpdateIgnoresOtherEvents(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c")
	s.Update(KeyEvent{Rune: 'x'})
	if s.Frame() != "a" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "a")
	}
}

func TestSpinnerUpdateReturnsEventUnchanged(t *testing.T) {
	s := NewSpinner()
	e := TickEvent{}
	if got := s.Update(e); got != Event(e) {
		t.Fatalf("Update should return the event unchanged, got %#v", got)
	}
}

func TestSpinnerTick(t *testing.T) {
	s := NewSpinner().WithFrames("a", "b", "c")
	s.Tick()
	if s.Frame() != "b" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "b")
	}
	s.Tick()
	if s.Frame() != "c" {
		t.Fatalf("Frame() = %q, want %q", s.Frame(), "c")
	}
	s.Tick()
	if s.Frame() != "a" {
		t.Fatalf("Frame() = %q, want %q (wrap around)", s.Frame(), "a")
	}
}

func TestSpinnerTickNoFrames(t *testing.T) {
	s := NewSpinner().WithFrames()
	s.Tick()
	if s.Frame() != "" {
		t.Fatalf("Frame() = %q, want empty", s.Frame())
	}
}

func TestSpinnerRender(t *testing.T) {
	s := NewSpinner().WithFrames("*")
	lines := s.Render(5, 1)
	if len(lines) != 1 || lines[0] != "*" {
		t.Fatalf("got %#v, want [\"*\"]", lines)
	}
}

func TestSpinnerWithSuffix(t *testing.T) {
	s := NewSpinner()
	same := s.WithSuffix("Loading")
	if same != s {
		t.Fatal("WithSuffix should return the same *Spinner for chaining")
	}
	if s.Suffix() != "Loading" {
		t.Fatalf("Suffix() = %q, want %q", s.Suffix(), "Loading")
	}
}

func TestSpinnerRenderWithSuffix(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithSuffix("Loading")
	lines := s.Render(20, 1)
	if len(lines) != 1 || lines[0] != "* Loading" {
		t.Fatalf("got %#v, want [\"* Loading\"]", lines)
	}
}

func TestSpinnerRenderWithSuffixClipsToWidth(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithSuffix("Loading")
	lines := s.Render(5, 1)
	if len(lines) != 1 || lines[0] != "* Loa" {
		t.Fatalf("got %#v, want [\"* Loa\"]", lines)
	}
}

func TestSpinnerPreferredWidthWithSuffix(t *testing.T) {
	s := NewSpinner().WithFrames("a", "abc").WithSuffix("Loading")
	// widest frame (3) + " " (1) + "Loading" (7) = 11
	if w := s.PreferredWidth(); w != 11 {
		t.Fatalf("PreferredWidth = %d, want 11", w)
	}
}

func TestSpinnerWithPrefix(t *testing.T) {
	s := NewSpinner()
	same := s.WithPrefix("Status:")
	if same != s {
		t.Fatal("WithPrefix should return the same *Spinner for chaining")
	}
	if s.Prefix() != "Status:" {
		t.Fatalf("Prefix() = %q, want %q", s.Prefix(), "Status:")
	}
}

func TestSpinnerRenderWithPrefix(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithPrefix("Status:")
	lines := s.Render(20, 1)
	if len(lines) != 1 || lines[0] != "Status: *" {
		t.Fatalf("got %#v, want [\"Status: *\"]", lines)
	}
}

func TestSpinnerRenderWithPrefixClipsToWidth(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithPrefix("Status:")
	lines := s.Render(5, 1)
	if len(lines) != 1 || lines[0] != "Statu" {
		t.Fatalf("got %#v, want [\"Statu\"]", lines)
	}
}

func TestSpinnerRenderWithPrefixAndSuffix(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithPrefix("Status:").WithSuffix("Loading")
	lines := s.Render(30, 1)
	if len(lines) != 1 || lines[0] != "Status: * Loading" {
		t.Fatalf("got %#v, want [\"Status: * Loading\"]", lines)
	}
}

func TestSpinnerPreferredWidthWithPrefixAndSuffix(t *testing.T) {
	s := NewSpinner().WithFrames("a", "abc").WithPrefix("Status:").WithSuffix("Loading")
	// "Status:" (7) + " " (1) + widest frame (3) + " " (1) + "Loading" (7) = 19
	if w := s.PreferredWidth(); w != 19 {
		t.Fatalf("PreferredWidth = %d, want 19", w)
	}
}

func TestSpinnerRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	s := NewSpinner().WithFrames("*").WithStyle(NewStyle().WithForeground(ColorRed))
	lines := s.Render(5, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled frame to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 1 {
		t.Fatalf("StringWidth = %d, want 1", w)
	}
}

func TestSpinnerRenderWithPrefixStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	s := NewSpinner().WithFrames("*").WithPrefix("Status:").WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	lines := s.Render(20, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled prefix to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 9 {
		t.Fatalf("StringWidth = %d, want 9", w)
	}
}

func TestSpinnerRenderWithSuffixStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	s := NewSpinner().WithFrames("*").WithSuffix("Loading").WithSuffixStyle(NewStyle().WithForeground(ColorRed))
	lines := s.Render(20, 1)
	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled suffix to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 9 {
		t.Fatalf("StringWidth = %d, want 9", w)
	}
}

func TestSpinnerWithoutStyleRemovesFrameStyleOnly(t *testing.T) {
	s := NewSpinner().WithStyle(NewStyle().WithForeground(ColorRed))
	s.WithoutStyle()
	if s.Style() != nil {
		t.Fatalf("Style() = %v, want nil", s.Style())
	}
}

func TestSpinnerWithoutPrefixStyle(t *testing.T) {
	s := NewSpinner().WithPrefixStyle(NewStyle().WithForeground(ColorRed))
	s.WithoutPrefixStyle()
	if s.PrefixStyle() != nil {
		t.Fatalf("PrefixStyle() = %v, want nil", s.PrefixStyle())
	}
}

func TestSpinnerWithoutSuffixStyle(t *testing.T) {
	s := NewSpinner().WithSuffixStyle(NewStyle().WithForeground(ColorRed))
	s.WithoutSuffixStyle()
	if s.SuffixStyle() != nil {
		t.Fatalf("SuffixStyle() = %v, want nil", s.SuffixStyle())
	}
}

func TestSpinnerRenderClipsToWidth(t *testing.T) {
	s := NewSpinner().WithFrames("abc")
	lines := s.Render(2, 1)
	if len(lines) != 1 || lines[0] != "ab" {
		t.Fatalf("got %#v, want [\"ab\"]", lines)
	}
}

func TestSpinnerRenderZeroSize(t *testing.T) {
	s := NewSpinner()
	if lines := s.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := s.Render(1, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestSpinnerPreferredWidth(t *testing.T) {
	s := NewSpinner().WithFrames("a", "abc", "ab")
	if w := s.PreferredWidth(); w != 3 {
		t.Fatalf("PreferredWidth = %d, want 3", w)
	}
}

func TestSpinnerPreferredHeight(t *testing.T) {
	s := NewSpinner()
	if h := s.PreferredHeight(10); h != 1 {
		t.Fatalf("PreferredHeight = %d, want 1", h)
	}
}
