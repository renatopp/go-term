package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/renatopp/go-term"
)

type Page struct {
	components []term.Component
}

func NewPage() *Page {
	return &Page{
		components: []term.Component{
			term.NewSpinner().WithSuffix("Loading..."),
			term.NewCursor(),
			term.NewProgressBar().
				AsShowPercent(true).
				WithPrefix("renato").
				WithSuffix("pereira"),
			term.NewIntCounter(0).
				WithSpeed(100 * time.Millisecond).
				WithPrefix("Int counter: "),
			term.NewFloatCounter(0).
				WithPrefix("Float counter: "),
			term.NewTextInput().
				WithPrefix("> ").
				WithPlaceholder("username").
				WithPlaceholderStyle(term.NewStyle().AsDim(true)).
				WithWidth(10).
				AsFocused(true),
			term.NewPasswordInput().
				WithPrefix("* ").
				WithPlaceholder("password").
				AsFocused(true),
			term.NewNumberInput().
				WithPrefix("$ ").
				WithPlaceholder("number").
				AsFocused(true),
			term.NewConfirm("Confirm??").AsFocused(true),
			term.NewSelect("a", "b", "c").AsFocused(true),
		},
	}
}

func (p *Page) Update(e term.Event) term.Event {
	switch e := e.(type) {
	case term.KeyEvent:
		if e.Rune == 'q' {
			return term.Quit
		}
		// p.intCounter.WithValue(rand.IntN(10000))
		// p.floatCounter.WithValue(rand.Float64() * 10)
		// p.progress.WithValue(rand.Float64())

	case term.TickEvent:
	}

	for _, p := range p.components {
		p.Update(e)
	}

	return e
}

func (p *Page) Render(width, height int) []string {
	items := []*term.ContainerItem{
		term.Item(term.Row(
			term.Item(term.NewLabel("Name")).WithBasisPercent(30),
			term.Item(term.NewLabel("Renato")).WithBasisPercent(70),
		)),
		term.Item(
			term.NewList("a", "b", "c"),
		),
	}
	for _, c := range p.components {
		items = append(items, term.Item(c))
	}

	return term.Column(
		items...,
	).Render(width, height)
}

func main() {
	term.
		NewProgram(NewPage()).
		AsAlternateScreen().
		WithTick(time.Second / 60).
		Run()
}

func print(label term.Renderable, width, height int) {
	lines := label.Render(width, height)
	fmt.Println("|" + strings.Join(lines, "\n|"))
	fmt.Println("")
}

// 	term.EnterAlternateScreen()
// 	defer term.ExitAlternateScreen()

// 	term.ClearScreen()
// 	println("Hello, World!")
// 	checkers()

// 	time.Sleep(2 * time.Second)
// 	labels()

// 	term.SetWindowTitle("Done...")
// 	time.Sleep(1 * time.Second)
// }

// func labels() {
// 	term.ClearScreen()
// 	term.MoveCursorToHome()

// 	println("Row layout (justify: space-between, gap: 2):")
// 	row := term.NewContainer(term.DirectionRow).
// 		WithJustify(term.JustifySpaceBetween).
// 		WithGap(2).
// 		Add(
// 			term.NewItem(term.NewLabel("Left")),
// 			term.NewItem(term.NewLabel("Center")),
// 			term.NewItem(term.NewLabel("Right")),
// 		)
// 	w, _, _ := term.GetScreenSize()
// 	for _, line := range row.Render(w, 1) {
// 		println("|" + line + "|")
// 	}

// 	println()
// 	println("Column layout (width 20, height 3):")
// 	col := term.NewContainer(term.DirectionColumn).
// 		Add(
// 			term.NewItem(term.NewLabel("first")),
// 			term.NewItem(term.NewLabel("second, longer")),
// 			term.NewItem(term.NewLabel("third")),
// 		)
// 	for _, line := range col.Render(20, 3) {
// 		println("|" + line + "|")
// 	}

// 	time.Sleep(3 * time.Second)
// }

// func checkers() {
// 	term.SetWindowTitle("Checking...")
// 	w, h, e1 := term.GetScreenSize()
// 	r, c, e2 := term.GetCursorPosition()
// 	println("- screen size:       ", w, "x", h, "--", e1)
// 	println("- cursor pos :       ", r, "x", c, "--", e2)
// 	println("- is stdout terminal?", term.IsTerminal(os.Stdout))
// 	println("- is stdin terminal? ", term.IsTerminal(os.Stdin))
// 	println("- is stderr terminal?", term.IsTerminal(os.Stderr))
// 	println("- is stdout pipe?    ", term.IsPipe(os.Stdout))
// 	println("- is stdin pipe?     ", term.IsPipe(os.Stdin))
// 	println("- is stderr pipe?    ", term.IsPipe(os.Stderr))
// 	println("- is stdout file?    ", term.IsFile(os.Stdout))
// 	println("- is stdin file?     ", term.IsFile(os.Stdin))
// 	println("- is stderr file?    ", term.IsFile(os.Stderr))
// 	println("-----")
// 	println("- is dumb?           ", term.IsDumb())
// 	println("- is wsl?            ", term.IsWSL())
// 	println("- is ssh?            ", term.IsSSH())
// 	println("- is Docker?         ", term.IsDocker())
// 	println("- supports colors?   ", term.SupportsColor())
// 	println("- supports true?     ", term.SupportsColorMode(term.ColorModeTrue))
// }

// // os.Getenv("TERM") != "dumb"
