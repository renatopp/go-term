package term

import (
	"context"
	"time"
)

// DefaultFPS is used when no FPS is configured via WithFPS.
const DefaultFPS = 30

type Program struct {
	ctx       context.Context
	cancel    context.CancelFunc
	queue     chan Event
	alternate bool
	mouse     bool
	dirty     bool
	fps       int
}

func NewProgram() *Program {
	return &Program{
		ctx:    context.Background(),
		cancel: func() {},
		queue:  make(chan Event, 64),
	}
}

func (p *Program) FPS() int {
	return p.fps
}

func (p *Program) AsAlternateScreen() *Program {
	p.alternate = true
	return p
}

func (p *Program) WithMouse() *Program {
	p.mouse = true
	return p
}

func (p *Program) WithFPS(fps int) *Program {
	p.fps = fps
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

	OnEvent(p.Send)

	ClearScreen()

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
	fps := p.fps
	if fps <= 0 {
		fps = DefaultFPS
	}
	ticker := time.NewTicker(time.Second / time.Duration(fps))
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case event := <-p.queue:
			p.update(event)
		case <-ticker.C:
			if p.dirty {
				p.draw()
				p.dirty = false
			}
		}
	}
}

func (p *Program) update(event Event) {
	p.dirty = true

	switch event.(type) {
	case SignalEvent:
		p.Stop()
	}
}

func (p *Program) draw() {

}
