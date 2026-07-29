package main

import (
	"time"

	"github.com/renatopp/go-term"
	"github.com/renatopp/go-term/ui"
)

func main() {
	term.EnterAlternateScreen()
	defer term.ExitAlternateScreen()

	term.ClearScreen()

	row := ui.Row().
		WithJustify(ui.JustifySpaceBetween).
		WithItem(
			ui.NewItem(term.NewLabel("Left")).WithGrow(0),
			ui.NewItem(term.NewLabel("Center")).WithGrow(1),
			ui.NewItem(term.NewLabel("Right")).WithGrow(0),
		)

	w, h := term.ForceGetScreenSize()
	for _, line := range row.Render(w, h) {
		print(line)
	}

	time.Sleep(2 * time.Second)
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
