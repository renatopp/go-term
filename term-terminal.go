package term

import (
	"golang.org/x/term"
)

var rawModeState *State

func EnterRawMode() error {
	if rawModeState != nil {
		return nil
	}
	state, err := term.MakeRaw(int(stdin))
	if err != nil {
		return err
	}
	rawModeState = state
	return nil
}

func ExitRawMode() error {
	state := rawModeState
	rawModeState = nil
	return RestoreTerminalState(state)
}

func GetTerminalState() (*State, error) {
	return term.GetState(int(stdin))
}

func RestoreTerminalState(state *State) error {
	return term.Restore(int(stdin), state)
}

func SetWindowTitle(title string) error {
	return write("\x1b]0;" + title + "\x07")
}

func WithinRawMode(f func()) error {
	if rawModeState != nil {
		f()
		return nil
	}
	if err := EnterRawMode(); err != nil {
		return err
	}
	defer ExitRawMode()
	f()
	return nil
}
