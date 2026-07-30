package term

import (
	"os"
	"sync"
	"sync/atomic"
)

// Event types published on the package-level bus by the stdin reader
// started in EnterRawMode.
const (
	EventKey            = "key"
	EventMouse          = "mouse"
	EventCursorPosition = "cursor-position"
	EventResize         = "resize"
	EventSignal         = "signal"
)

var listenerId = atomic.Uint64{}

// Event is a value flowing through the event loop's queue and the package
// bus, carrying events from OS signals, external callbacks, timers, and
// stdin input to be processed by the loop.
type Event any

// SignalEvent wraps an OS signal delivered to the process.
type SignalEvent struct {
	Signal os.Signal
}

// KeyEvent reports a keyboard event read from stdin while in raw mode.
type KeyEvent struct {
	Type KeyType
	Rune rune
	Ctrl bool
	Alt  bool
}

// CursorPositionEvent reports the terminal's reply to a cursor position
// request (see GetCursorPosition). Row and Col are 1-based.
type CursorPositionEvent struct {
	Row, Col int
}

// MouseEvent reports a mouse event: a button press/release, a drag while a
// button is held, or a wheel scroll. X and Y are 0-based screen columns and
// rows.
type MouseEvent struct {
	X, Y   int
	Button MouseButton
	Action MouseAction
	Shift  bool
	Alt    bool
	Ctrl   bool
}

type ResizeEvent struct {
	Width, Height int
}

type eventBus struct {
	mu        sync.Mutex
	listeners map[string]map[uint64]func(Event)
}

func newEventBus() *eventBus {
	return &eventBus{
		listeners: make(map[string]map[uint64]func(Event)),
	}
}

func (eb *eventBus) CountListeners(eventType string) int {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if listeners, exists := eb.listeners[eventType]; exists {
		return len(listeners)
	}
	return 0
}

func (eb *eventBus) Subscribe(eventType string, listener func(Event)) uint64 {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, exists := eb.listeners[eventType]; !exists {
		eb.listeners[eventType] = make(map[uint64]func(Event))
	}
	id := listenerId.Add(1)
	eb.listeners[eventType][id] = listener
	return id
}

func (eb *eventBus) Unsubscribe(eventType string, id uint64) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if listeners, exists := eb.listeners[eventType]; exists {
		delete(listeners, id)
		if len(listeners) == 0 {
			delete(eb.listeners, eventType)
		}
	}
}

func (eb *eventBus) UnsubscribeAll(eventType string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	delete(eb.listeners, eventType)
}

func (eb *eventBus) Publish(eventType string, event Event) {
	eb.mu.Lock()
	listeners1, exists1 := eb.listeners[eventType]
	listeners2, exists2 := eb.listeners["*"]
	eb.mu.Unlock()

	if exists1 {
		for _, listener := range listeners1 {
			listener(event)
		}
	}

	if exists2 {
		for _, listener := range listeners2 {
			listener(event)
		}
	}
}
