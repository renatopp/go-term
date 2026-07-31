package ui

type Event any

// MultiEvent groups several events into one, so an Update method can react
// to a single event by producing more than one in return. Use Events to
// build one.
type MultiEvent []Event

// Events wraps the given events into a MultiEvent for returning from Update.
func Events(events ...Event) Event {
	return MultiEvent(events)
}

type Component interface {
	Renderable
	Updatable
}

type Renderable interface {
	Render(width, height int) []string
}

// Updatable reacts to an event and optionally returns one for further
// processing. A nil return stops that event from being processed any
// further, which lets a component suppress default handling (e.g. the
// program's built-in Ctrl+C behavior) for events it has fully handled
// itself. Returning the event unchanged, a different event, or a MultiEvent
// lets processing continue.
type Updatable interface {
	Update(Event) Event
}

type Sizeable interface {
	PreferredWidth() int
	PreferredHeight(width int) int
}

type Rect struct {
	X, Y, Width, Height int
}
