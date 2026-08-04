package term

import "testing"

func TestNewLabel(t *testing.T) {
	l := NewLabel("hello")
	if l.text != "hello" {
		t.Fatalf("text = %q, want %q", l.text, "hello")
	}
}

func TestLabelWithText(t *testing.T) {
	l := NewLabel("a")
	same := l.WithText("b")
	if same != l {
		t.Fatal("WithText should return the same *Label for chaining")
	}
	if l.text != "b" {
		t.Fatalf("text = %q, want %q", l.text, "b")
	}
}

func TestLabelRender(t *testing.T) {
	l := NewLabel("hello")
	lines := l.Render(10, 1)
	if len(lines) != 1 || lines[0] != "hello" {
		t.Fatalf("got %#v, want [\"hello\"]", lines)
	}
}

func TestLabelRenderTruncatesToWidth(t *testing.T) {
	l := NewLabel("hello world")
	lines := l.Render(5, 1)
	if len(lines) != 1 || lines[0] != "hell…" {
		t.Fatalf("got %#v, want [\"hell…\"]", lines)
	}
}

func TestLabelRenderTruncatesByRuneNotByte(t *testing.T) {
	l := NewLabel("héllo")
	lines := l.Render(3, 1)
	if len(lines) != 1 || lines[0] != "hé…" {
		t.Fatalf("got %#v, want [\"hé…\"]", lines)
	}
}

func TestLabelRenderZeroSize(t *testing.T) {
	l := NewLabel("hello")
	if lines := l.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := l.Render(5, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestLabelPreferredWidth(t *testing.T) {
	l := NewLabel("héllo")
	if w := l.PreferredWidth(); w != 5 {
		t.Fatalf("PreferredWidth = %d, want 5", w)
	}
}

func TestLabelPreferredHeight(t *testing.T) {
	l := NewLabel("hello world this is long")
	if h := l.PreferredHeight(30); h != 1 {
		t.Fatalf("PreferredHeight = %d, want 1", h)
	}
}
