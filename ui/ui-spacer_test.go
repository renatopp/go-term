package ui

import "testing"

func TestSpacerDefaults(t *testing.T) {
	s := NewSpacer()

	if s.PreferredWidth() != 0 {
		t.Fatalf("default PreferredWidth = %d, want 0", s.PreferredWidth())
	}
	if s.PreferredHeight(100) != 0 {
		t.Fatalf("default PreferredHeight = %d, want 0", s.PreferredHeight(100))
	}
	if s.Render(10, 10) != nil {
		t.Fatalf("Render = %v, want nil", s.Render(10, 10))
	}
}

func TestSpacerWithSize(t *testing.T) {
	s := NewSpacer().WithWidth(5).WithHeight(3)

	if s.PreferredWidth() != 5 {
		t.Fatalf("PreferredWidth = %d, want 5", s.PreferredWidth())
	}
	if s.PreferredHeight(100) != 3 {
		t.Fatalf("PreferredHeight = %d, want 3", s.PreferredHeight(100))
	}
}

func TestSpacerWithersChain(t *testing.T) {
	s := NewSpacer()
	same := s.WithWidth(1).WithHeight(2)

	if same != s {
		t.Fatal("With* methods should return the same *Spacer for chaining")
	}
}
