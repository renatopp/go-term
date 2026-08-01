package term

import (
	"testing"
	"time"
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

func TestSpinnerWithText(t *testing.T) {
	s := NewSpinner()
	same := s.WithText("Loading")
	if same != s {
		t.Fatal("WithText should return the same *Spinner for chaining")
	}
	if s.Text() != "Loading" {
		t.Fatalf("Text() = %q, want %q", s.Text(), "Loading")
	}
}

func TestSpinnerRenderWithText(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithText("Loading")
	lines := s.Render(20, 1)
	if len(lines) != 1 || lines[0] != "* Loading" {
		t.Fatalf("got %#v, want [\"* Loading\"]", lines)
	}
}

func TestSpinnerRenderWithTextClipsToWidth(t *testing.T) {
	s := NewSpinner().WithFrames("*").WithText("Loading")
	lines := s.Render(5, 1)
	if len(lines) != 1 || lines[0] != "* Loa" {
		t.Fatalf("got %#v, want [\"* Loa\"]", lines)
	}
}

func TestSpinnerPreferredWidthWithText(t *testing.T) {
	s := NewSpinner().WithFrames("a", "abc").WithText("Loading")
	// widest frame (3) + " " (1) + "Loading" (7) = 11
	if w := s.PreferredWidth(); w != 11 {
		t.Fatalf("PreferredWidth = %d, want 11", w)
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
