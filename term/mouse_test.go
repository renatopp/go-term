package term

import (
	"io"
	"testing"
	"time"
)

func withMouseSource(t *testing.T, r io.Reader) {
	t.Helper()
	original := mouseSource
	mouseSource = r
	t.Cleanup(func() { mouseSource = original })
}

func TestDecodeSGRMouseLeftPress(t *testing.T) {
	msg, ok := decodeSGRMouse("0;10;5", true)
	if !ok {
		t.Fatal("expected ok")
	}
	want := MouseMsg{X: 9, Y: 4, Button: MouseButtonLeft, Action: MouseActionPress}
	if msg != want {
		t.Fatalf("got %+v, want %+v", msg, want)
	}
}

func TestDecodeSGRMouseRelease(t *testing.T) {
	msg, ok := decodeSGRMouse("2;1;1", false)
	if !ok {
		t.Fatal("expected ok")
	}
	if msg.Button != MouseButtonRight || msg.Action != MouseActionRelease {
		t.Fatalf("got %+v, want right release", msg)
	}
}

func TestDecodeSGRMouseWheel(t *testing.T) {
	up, ok := decodeSGRMouse("64;1;1", true)
	if !ok || up.Button != MouseButtonWheelUp || up.Action != MouseActionPress {
		t.Fatalf("got %+v, ok=%v, want wheel up press", up, ok)
	}

	down, ok := decodeSGRMouse("65;1;1", true)
	if !ok || down.Button != MouseButtonWheelDown {
		t.Fatalf("got %+v, ok=%v, want wheel down", down, ok)
	}
}

func TestDecodeSGRMouseModifiers(t *testing.T) {
	msg, ok := decodeSGRMouse("28;1;1", true) // 4 shift + 8 alt + 16 ctrl
	if !ok || !msg.Shift || !msg.Alt || !msg.Ctrl {
		t.Fatalf("got %+v, ok=%v, want shift+alt+ctrl", msg, ok)
	}
}

func TestDecodeSGRMouseDrag(t *testing.T) {
	msg, ok := decodeSGRMouse("32;1;1", true) // motion bit + left button held
	if !ok || msg.Action != MouseActionMotion || msg.Button != MouseButtonLeft {
		t.Fatalf("got %+v, ok=%v, want left drag", msg, ok)
	}
}

func TestDecodeSGRMouseMalformed(t *testing.T) {
	if _, ok := decodeSGRMouse("x;1;1", true); ok {
		t.Fatal("expected not ok for non-numeric param")
	}
	if _, ok := decodeSGRMouse("1;2", true); ok {
		t.Fatal("expected not ok for too few params")
	}
}

func TestOnMouseParsesEventAndIgnoresOtherEscapes(t *testing.T) {
	pr, pw := io.Pipe()
	withMouseSource(t, pr)
	t.Cleanup(func() { pw.Close() })

	got := make(chan MouseMsg, 1)
	stop := OnMouse(func(m MouseMsg) { got <- m })
	defer stop()

	go pw.Write([]byte("garbage\x1b[A\x1b[<0;3;4M"))

	select {
	case msg := <-got:
		want := MouseMsg{X: 2, Y: 3, Button: MouseButtonLeft, Action: MouseActionPress}
		if msg != want {
			t.Fatalf("got %+v, want %+v", msg, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for mouse event")
	}
}

func TestOnMouseStopStopsDelivery(t *testing.T) {
	pr, pw := io.Pipe()
	withMouseSource(t, pr)
	t.Cleanup(func() { pw.Close() })

	calls := 0
	stop := OnMouse(func(m MouseMsg) { calls++ })
	stop()

	pw.Write([]byte("\x1b[<0;1;1M"))
	time.Sleep(20 * time.Millisecond)

	if calls != 0 {
		t.Fatalf("callback invoked after stop: %d calls", calls)
	}
}
