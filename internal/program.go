package term

import (
	"context"
	"time"
)

// DefaultFPS is used when no FPS is configured via WithFPS.
const DefaultFPS = 30

// DefaultTick is used when no tick interval is configured via WithTick.
const DefaultTick = 100 * time.Millisecond

type Program struct {
	ctx       context.Context
	cancel    context.CancelFunc
	queue     chan Event
	root      Component
	buffer    *buffer
	width     int
	height    int
	alternate bool
	mouse     bool
	dirty     bool
	fps       int
	tick      time.Duration
}

func NewProgram(root Component) *Program {
	return &Program{
		ctx:    context.Background(),
		cancel: func() {},
		queue:  make(chan Event, 64),
		root:   root,
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

func (p *Program) Tick() time.Duration {
	return p.tick
}

// WithTick sets the interval at which a TickEvent is dispatched to the root
// component, driving time-based components like Spinner.
func (p *Program) WithTick(d time.Duration) *Program {
	p.tick = d
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

	HideCursor()
	defer ShowCursor()

	OnEvent(p.Send)

	p.width, p.height = ForceGetScreenSize()

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

	tick := p.tick
	if tick <= 0 {
		tick = DefaultTick
	}
	tickTicker := time.NewTicker(tick)
	defer tickTicker.Stop()

	lastTick := time.Now()

	p.draw()
	for {
		select {
		case <-p.ctx.Done():
			return
		case event := <-p.queue:
			p.update(event)
		case now := <-tickTicker.C:
			p.update(TickEvent{Time: now, Duration: now.Sub(lastTick)})
			lastTick = now
		case <-ticker.C:
			if p.dirty {
				p.draw()
				p.dirty = false
			}
		}
	}
}

func (p *Program) update(event Event) {
	switch e := event.(type) {
	case ResizeEvent:
		p.width, p.height = e.Width, e.Height
	}

	p.dirty = true
	p.dispatch(p.root.Update(event))
}

// dispatch stops the program on a SignalEvent or Ctrl+C. A MultiEvent is
// unpacked and each of its events is dispatched in turn. A nil event (as
// returned by Component.Update to suppress default handling) is ignored.
func (p *Program) dispatch(event Event) {
	if multi, ok := event.(MultiEvent); ok {
		for _, e := range multi {
			p.dispatch(e)
		}
		return
	}

	switch e := event.(type) {
	case SignalEvent, QuitEvent:
		p.Stop()
	case KeyEvent:
		if e.Rune == 'c' && e.Ctrl {
			p.Stop()
		}
	}
}

func (p *Program) draw() {
	if p.buffer == nil {
		p.buffer = newBuffer(p.width, p.height)
	} else if p.buffer.width != p.width || p.buffer.height != p.height {
		p.buffer.Resize(p.width, p.height)
	}

	p.buffer.Flush(p.root.Render(p.width, p.height), stdout)
}
