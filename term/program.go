package term

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Program struct {
	ctx    context.Context
	cancel context.CancelFunc
	queue  chan Event
	mouse  bool
}

func NewProgram() *Program {
	return &Program{
		ctx:    context.Background(),
		cancel: func() {},
		queue:  make(chan Event, 64),
	}
}

// WithMouse enables mouse reporting (press, release, drag, and wheel) for
// the program's lifetime; events arrive on the queue as MouseEvent. Mouse
// reporting is off by default, since capturing it also disables the
// terminal's own text selection for callers that don't need it.
func (p *Program) WithMouse() *Program {
	p.mouse = true
	return p
}

// Run starts the program's event loop. It terminates when ctx is done or
// when Stop is called.
func (p *Program) Run() error {
	ctx, cancel := context.WithCancel(p.ctx)
	p.ctx = ctx
	p.cancel = cancel
	defer cancel()

	EnterAlternateScreen()
	defer ExitAlternateScreen()

	EnterRawMode()
	defer ExitRawMode()

	ClearScreen()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigs)

	go p.forwardSignals(sigs)

	if p.mouse {
		EnableMouse()
		defer DisableMouse()

		stopMouse := OnMouse(func(m MouseEvent) { p.Send(m) })
		defer stopMouse()
	}

	p.eventLoop()
	return nil
}

// Stop requests the event loop to terminate.
func (p *Program) Stop() {
	p.cancel()
}

// Send enqueues an event onto the event loop's queue. Signal handlers,
// timers, and external callbacks all use this to feed events into the loop.
func (p *Program) Send(event Event) {
	select {
	case p.queue <- event:
	case <-p.ctx.Done():
	}
}

func (p *Program) eventLoop() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case event := <-p.queue:
			p.dispatch(event)
		}
	}
}

func (p *Program) forwardSignals(sigs chan os.Signal) {
	for {
		select {
		case sig := <-sigs:
			p.Send(SignalEvent{Signal: sig})
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Program) dispatch(event Event) {
	switch event.(type) {
	case SignalEvent:
		p.Stop()
	}
}
