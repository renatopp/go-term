package main

import (
	"fmt"
	"strings"

	"github.com/renatopp/go-term/term"
	"github.com/renatopp/go-term/term/ui"
)

type Page struct {
}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) Update(e ui.Event) ui.Event {
	switch e := e.(type) {
	case term.KeyEvent:
		if e.Rune == 'q' {
			return term.Quit
		}
	}
	return e
}

func (p *Page) Render(width, height int) []string {
	return ui.Column(
		ui.Item(ui.Row(
			ui.Item(term.NewLabel("Name")).WithBasisPercent(30),
			ui.Item(term.NewLabel("Renato")).WithBasisPercent(70),
		)),
		ui.Item(ui.Row(
			ui.Item(term.NewLabel("Email")).WithBasisPercent(30),
			ui.Item(term.NewLabel("renato@renato.com")).WithBasisPercent(70),
		)),
		ui.Spacer(),
		ui.Item(term.NewLabel("q quit").WithStyle(term.NewStyle().AsDim(true))),
	).Render(width, height)
}

func main() {
	l := term.NewLabel("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since 1966, when designers at Letraset and James Mosley, the librarian at St Bride Printing Library in London, took a 1914 Cicero translation and scrambled it to make dummy text for Letraset's Body Type sheets. It has survived not only many decades, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised thanks to these sheets and more recently with desktop publishing software like Aldus PageMaker and Microsoft Word including versions of Lorem Ipsum.")
	print(l, 100, 1)
	print(l, 20, 5)
	print(l, 10, 5)
	print(l, 5, 5)

	term.
		NewProgram(NewPage()).
		AsAlternateScreen().
		Run()
}

func print(label ui.Renderable, width, height int) {
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
// 	row := ui.NewContainer(ui.DirectionRow).
// 		WithJustify(ui.JustifySpaceBetween).
// 		WithGap(2).
// 		Add(
// 			ui.NewItem(term.NewLabel("Left")),
// 			ui.NewItem(term.NewLabel("Center")),
// 			ui.NewItem(term.NewLabel("Right")),
// 		)
// 	w, _, _ := term.GetScreenSize()
// 	for _, line := range row.Render(w, 1) {
// 		println("|" + line + "|")
// 	}

// 	println()
// 	println("Column layout (width 20, height 3):")
// 	col := ui.NewContainer(ui.DirectionColumn).
// 		Add(
// 			ui.NewItem(term.NewLabel("first")),
// 			ui.NewItem(term.NewLabel("second, longer")),
// 			ui.NewItem(term.NewLabel("third")),
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
