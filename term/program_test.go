package term

import (
	"context"
	"testing"
	"time"

	"github.com/renatopp/go-term/term/ui"
)

func newTestProgram() *Program {
	p := NewProgram(nil)
	p.ctx, p.cancel = context.WithCancel(context.Background())
	return p
}

func (p *Program) stopped() bool {
	select {
	case <-p.ctx.Done():
		return true
	default:
		return false
	}
}

func TestDispatchStopsOnSignalEvent(t *testing.T) {
	p := newTestProgram()
	p.dispatch(SignalEvent{})
	if !p.stopped() {
		t.Fatal("expected SignalEvent to stop the program")
	}
}

func TestDispatchStopsOnCtrlC(t *testing.T) {
	p := newTestProgram()
	p.dispatch(KeyEvent{Rune: 'c', Ctrl: true})
	if !p.stopped() {
		t.Fatal("expected Ctrl+C to stop the program")
	}
}

func TestDispatchIgnoresNilEvent(t *testing.T) {
	p := newTestProgram()
	p.dispatch(nil)
	if p.stopped() {
		t.Fatal("expected nil event not to stop the program")
	}
}

func TestDispatchUnwrapsMultiEvent(t *testing.T) {
	p := newTestProgram()
	p.dispatch(ui.Events(KeyEvent{Rune: 'x'}, SignalEvent{}))
	if !p.stopped() {
		t.Fatal("expected a SignalEvent nested in a MultiEvent to stop the program")
	}
}

func TestProgramWithTick(t *testing.T) {
	p := NewProgram(nil).WithTick(50 * time.Millisecond)
	if p.Tick() != 50*time.Millisecond {
		t.Fatalf("Tick() = %v, want %v", p.Tick(), 50*time.Millisecond)
	}
}
