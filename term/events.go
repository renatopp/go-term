package term

import (
	"os"
	"sync"
	"sync/atomic"

	"github.com/renatopp/go-term/term/ui"
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

type QuitEvent int

const Quit QuitEvent = 0

// KeyType identifies a key reported by KeyEvent. KeyRune covers ordinary
// character keys (including Ctrl/Alt/Shift combinations, reported via the
// Ctrl, Alt, and Shift fields); the rest name keys with no printable rune of
// their own.
type KeyType int

const (
	KeyRune KeyType = iota
	KeyEnter
	KeyEsc
	KeyTab
	KeyBackspace
	KeyDelete
	KeyInsert
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
)

// keyboardProtocol identifies which keyboard-modifier reporting protocol is
// active on the terminal, negotiated by EnterRawMode.
type keyboardProtocol uint8

const (
	keyboardProtocolNone keyboardProtocol = iota
	keyboardProtocolKitty
	keyboardProtocolModifyOtherKeys
)

type MouseButton uint8

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
	MouseButtonNone // reported during a drag with no button held
)

type MouseAction uint8

const (
	MouseActionPress MouseAction = iota
	MouseActionRelease
	MouseActionMotion
)

var listenerId = atomic.Uint64{}

// Event is a value flowing through the event loop's queue and the package
// bus, carrying events from OS signals, external callbacks, timers, and
// stdin input to be processed by the loop.
type Event = ui.Event

// SignalEvent wraps an OS signal delivered to the process.
type SignalEvent struct {
	Signal os.Signal
}

// KeyEvent reports a keyboard event read from stdin while in raw mode. Ctrl
// and Alt are always populated when the terminal reports them (Ctrl+letter,
// Alt+key). Shift on a plain or Alt-combined rune is inferred from the
// character's case (e.g. 'A' implies Shift, 'a' doesn't), so it can't be
// told apart from Caps Lock and never applies to Ctrl+letter, whose control
// byte carries no case. Shift on keys with no room for modifiers in plain
// VT100 input otherwise (Enter, Tab, Backspace, Esc, Space, arrows, ...),
// and Ctrl/Alt on those same keys, are only populated when the terminal
// supports the Kitty keyboard protocol or xterm's modifyOtherKeys — see
// EnterRawMode.
type KeyEvent struct {
	Type  KeyType
	Rune  rune
	Shift bool
	Alt   bool
	Ctrl  bool
}

// CursorPositionEvent reports the terminal's reply to a cursor position
// request (see GetCursorPosition). Row and Col are 1-based.
type CursorPositionEvent struct {
	Row, Col int
}

// MouseEvent reports a mouse event: a button press/release, a drag while a
// button is held, or a wheel scroll. X and Y are 1-based screen columns and
// rows, consistent with CursorPositionEvent.
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
