package term

import (
	"strings"
	"testing"

	"github.com/renatopp/go-term/term/ui"
)

func TestNewBox(t *testing.T) {
	l := NewLabel("hi")
	b := NewBox(l)
	if b.Child() != l {
		t.Fatalf("Child() = %v, want %v", b.Child(), l)
	}
}

func TestBoxWithChild(t *testing.T) {
	b := NewBox(NewLabel("a"))
	l := NewLabel("b")
	same := b.WithChild(l)
	if same != b {
		t.Fatal("WithChild should return the same *Box for chaining")
	}
	if b.Child() != l {
		t.Fatalf("Child() = %v, want %v", b.Child(), l)
	}
}

func TestBoxRenderNoBorderNoPaddingNoMargin(t *testing.T) {
	b := NewBox(NewLabel("hi"))
	lines := b.Render(4, 2)
	want := []string{"hi  ", "    "}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestBoxRenderWithBorder(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithBorder(BorderSingle)
	lines := b.Render(4, 3)
	want := []string{
		"┌──┐",
		"│hi│",
		"└──┘",
	}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestBoxRenderWithPadding(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithBorder(BorderSingle).WithPadding(1)
	lines := b.Render(6, 5)
	if len(lines) != 5 {
		t.Fatalf("got %d lines, want 5: %#v", len(lines), lines)
	}
	for _, line := range lines {
		if w := len([]rune(line)); w != 6 {
			t.Fatalf("line %q has width %d, want 6", line, w)
		}
	}
	if lines[0] != "┌────┐" || lines[4] != "└────┘" {
		t.Fatalf("got %#v", lines)
	}
	if lines[2] != "│ hi │" {
		t.Fatalf("line 2 = %q, want %q", lines[2], "│ hi │")
	}
}

func TestBoxRenderWithMargin(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithBorder(BorderSingle).WithMargin(1)
	lines := b.Render(6, 5)
	want := []string{
		"      ",
		" ┌──┐ ",
		" │hi│ ",
		" └──┘ ",
		"      ",
	}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestBoxWithoutBorder(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithBorder(BorderSingle).WithoutBorder()
	if b.Border().hasBorder() {
		t.Fatal("WithoutBorder should clear the border")
	}
}

func TestBoxRenderWithStyle(t *testing.T) {
	prev := GetColorLevel()
	SetColorLevel(ColorModeTrue)
	defer SetColorLevel(prev)

	b := NewBox(NewLabel("hi")).
		WithBorder(BorderSingle).
		WithStyle(NewStyle().WithForeground(ColorRed))
	lines := b.Render(4, 3)

	if !strings.Contains(lines[0], "\x1b[") {
		t.Fatalf("expected styled border to contain an SGR sequence, got %q", lines[0])
	}
	if w := ui.StringWidth(lines[0]); w != 4 {
		t.Fatalf("StringWidth(lines[0]) = %d, want 4", w)
	}
	if w := ui.StringWidth(lines[1]); w != 4 {
		t.Fatalf("StringWidth(lines[1]) = %d, want 4", w)
	}
}

func TestBoxWithoutStyle(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithStyle(NewStyle().WithForeground(ColorRed)).WithoutStyle()
	if b.Style() != nil {
		t.Fatal("WithoutStyle should clear the style")
	}
}

func TestBoxPreferredWidth(t *testing.T) {
	b := NewBox(NewLabel("hello")).WithBorder(BorderSingle).WithPadding(1).WithMargin(2)
	// label width 5 + padding 2 + border 2 + margin 4 = 13
	if w := b.PreferredWidth(); w != 13 {
		t.Fatalf("PreferredWidth = %d, want 13", w)
	}
}

func TestBoxPreferredHeight(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithBorder(BorderSingle).WithPadding(1)
	// label height 1 + padding 2 + border 2 = 5
	if h := b.PreferredHeight(10); h != 5 {
		t.Fatalf("PreferredHeight = %d, want 5", h)
	}
}

func TestBoxRenderZeroSize(t *testing.T) {
	b := NewBox(NewLabel("hi"))
	if lines := b.Render(0, 2); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := b.Render(2, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestBoxRenderTooSmallForBorder(t *testing.T) {
	b := NewBox(NewLabel("hi")).WithBorder(BorderSingle)
	lines := b.Render(1, 1)
	if len(lines) != 1 || lines[0] != " " {
		t.Fatalf("got %#v, want [\" \"]", lines)
	}
}
