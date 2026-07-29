package term

import "os"

// Msg is a command flowing through the event loop's queue, carrying events
// from OS signals, external callbacks, and timers to be processed by the
// loop.
type Msg any

// SignalMsg wraps an OS signal delivered to the process.
type SignalMsg struct {
	Signal os.Signal
}
