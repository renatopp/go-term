package term

import "testing"

func TestNewList(t *testing.T) {
	l := NewList("a", "b")
	items := l.Items()
	if len(items) != 2 || items[0] != "a" || items[1] != "b" {
		t.Fatalf("Items() = %#v, want [\"a\" \"b\"]", items)
	}
}

func TestListDefaultBullet(t *testing.T) {
	l := NewList()
	if l.Bullet() != DefaultBullet {
		t.Fatalf("Bullet() = %q, want %q", l.Bullet(), DefaultBullet)
	}
}

func TestListWithItemAppends(t *testing.T) {
	l := NewList("a")
	same := l.WithItem("b", "c")
	if same != l {
		t.Fatal("WithItem should return the same *List for chaining")
	}
	items := l.Items()
	if len(items) != 3 || items[0] != "a" || items[1] != "b" || items[2] != "c" {
		t.Fatalf("Items() = %#v, want [\"a\" \"b\" \"c\"]", items)
	}
}

func TestListClear(t *testing.T) {
	l := NewList("a", "b").Clear()
	if len(l.Items()) != 0 {
		t.Fatalf("Items() = %#v, want empty", l.Items())
	}
}

func TestListWithBullet(t *testing.T) {
	l := NewList("a").WithBullet("-")
	if l.Bullet() != "-" {
		t.Fatalf("Bullet() = %q, want %q", l.Bullet(), "-")
	}
}

func TestListWithPaddingLeft(t *testing.T) {
	l := NewList("a").WithPaddingLeft(3)
	if l.PaddingLeft() != 3 {
		t.Fatalf("PaddingLeft() = %d, want 3", l.PaddingLeft())
	}
}

func TestListRenderWithPaddingLeft(t *testing.T) {
	l := NewList("one", "two").WithBullet("-").WithPaddingLeft(2)
	lines := l.Render(10, 2)
	want := []string{"  - one", "  - two"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestListRenderWithPaddingLeftWrapsAndIndents(t *testing.T) {
	l := NewList("hello world").WithBullet("-").WithPaddingLeft(2)
	lines := l.Render(9, 2)
	want := []string{"  - hello", "    world"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestListPreferredWidthWithPaddingLeft(t *testing.T) {
	l := NewList("a").WithBullet("-").WithPaddingLeft(2)
	// padding 2 + "- " (2) + "a" (1) = 5
	if w := l.PreferredWidth(); w != 5 {
		t.Fatalf("PreferredWidth = %d, want 5", w)
	}
}

func TestListRender(t *testing.T) {
	l := NewList("one", "two").WithBullet("-")
	lines := l.Render(10, 2)
	want := []string{"- one", "- two"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestListRenderWrapsAndIndents(t *testing.T) {
	l := NewList("hello world").WithBullet("-")
	lines := l.Render(7, 2)
	want := []string{"- hello", "  world"}
	if len(lines) != len(want) {
		t.Fatalf("got %#v, want %#v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestListRenderTruncatesToHeight(t *testing.T) {
	l := NewList("one", "two", "three").WithBullet("-")
	lines := l.Render(10, 2)
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %#v", len(lines), lines)
	}
	if lines[0] != "- one" {
		t.Fatalf("line 0 = %q, want %q", lines[0], "- one")
	}
	if lines[1] != "- two…" {
		t.Fatalf("line 1 = %q, want %q", lines[1], "- two…")
	}
}

func TestListRenderZeroSize(t *testing.T) {
	l := NewList("a")
	if lines := l.Render(0, 1); lines != nil {
		t.Fatalf("got %#v, want nil for zero width", lines)
	}
	if lines := l.Render(5, 0); lines != nil {
		t.Fatalf("got %#v, want nil for zero height", lines)
	}
}

func TestListPreferredWidth(t *testing.T) {
	l := NewList("a", "longer item").WithBullet("-")
	// "- " (2) + "longer item" (11) = 13
	if w := l.PreferredWidth(); w != 13 {
		t.Fatalf("PreferredWidth = %d, want 13", w)
	}
}

func TestListPreferredHeight(t *testing.T) {
	l := NewList("one", "hello world").WithBullet("-")
	// "one" -> 1 line, "hello world" wrapped at width 7 -> 2 lines
	if h := l.PreferredHeight(7); h != 3 {
		t.Fatalf("PreferredHeight = %d, want 3", h)
	}
}
