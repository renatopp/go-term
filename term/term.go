package term

import (
	"io"
	"os"

	"golang.org/x/term"
)

var writer io.Writer = os.Stdout
var stdin uintptr = os.Stdin.Fd()
var stdout uintptr = os.Stdout.Fd()

func SetWriter(w io.Writer) {
	writer = w
}

func SetStdin(fd uintptr) {
	stdin = fd
}

func SetStdout(fd uintptr) {
	stdout = fd
}

type State = term.State

func write(s string) error {
	_, err := writer.Write([]byte(s))
	return err
}
