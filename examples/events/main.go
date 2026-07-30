package main

import (
	"fmt"

	"github.com/renatopp/go-term/term"
)

func main() {
	term.EnterRawMode()
	defer term.ExitRawMode()

	close := make(chan bool)

	term.OnEvent(func(e term.Event) {
		fmt.Printf("[%T] %#v\n", e, e)
		switch e := e.(type) {
		case term.KeyEvent:
			if e.Rune == 'c' && e.Ctrl {
				close <- true
				return
			}
		case term.MouseEvent:
			fmt.Printf("Mouse event: %+v\n", e)
		}
	})

	<-close
}
