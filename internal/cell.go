package term

import (
	"strings"
	"unicode"
)

// Cell is a single terminal column: its visible content, the number of
// columns it occupies (0 for the placeholder following a wide character, 1
// for a normal character, 2 for a wide one), and the raw SGR sequence
// active when it was written ("" for no styling).
type Cell struct {
	Text  string
	Width int
	Style string
}

// BuildRow scans line into exactly width columns, truncating content that
// overflows and padding short lines with blank cells. The line may contain
// SGR escape sequences, which are captured as each cell's Style.
func BuildRow(line string, width int) []Cell {
	row := make([]Cell, 0, width)
	runes := []rune(line)
	style := ""
	x := 0

	for i := 0; i < len(runes) && x < width; {
		r := runes[i]
		if r == '\x1b' {
			if seq, n := scanSGR(runes[i:]); n > 0 {
				if seq == "\x1b[0m" || seq == "\x1b[m" {
					style = ""
				} else {
					style = seq
				}
				i += n
			} else {
				i++
			}
			continue
		}

		w := RuneWidth(r)
		if w == 0 {
			if len(row) > 0 {
				row[len(row)-1].Text += string(r)
			}
			i++
			continue
		}
		if x+w > width {
			break
		}

		row = append(row, Cell{Text: string(r), Width: w, Style: style})
		if w == 2 {
			row = append(row, Cell{Width: 0})
		}
		x += w
		i++
	}

	for x < width {
		row = append(row, Cell{Text: " ", Width: 1})
		x++
	}
	return row
}

// RenderRow serializes cells back into a single line, emitting an SGR
// sequence only where the active style changes and a reset at the end.
func RenderRow(cells []Cell) string {
	var out strings.Builder
	style := ""
	for _, c := range cells {
		if c.Width == 0 {
			continue
		}
		if c.Style != style {
			if style != "" {
				out.WriteString("\x1b[0m")
			}
			out.WriteString(c.Style)
			style = c.Style
		}
		out.WriteString(c.Text)
	}
	if style != "" {
		out.WriteString("\x1b[0m")
	}
	return out.String()
}

// RuneWidth returns the number of terminal columns r occupies: 0 for
// combining marks and joiners that attach to the previous cell, 2 for wide
// East Asian characters, 1 otherwise.
func RuneWidth(r rune) int {
	switch {
	case r == 0x200D, r == 0xFEFF, r >= 0xFE00 && r <= 0xFE0F, r >= 0xE0100 && r <= 0xE01EF:
		return 0
	case unicode.Is(unicode.Mn, r), unicode.Is(unicode.Me, r):
		return 0
	case isWideRune(r):
		return 2
	default:
		return 1
	}
}

// StringWidth returns the number of terminal columns s occupies, ignoring
// SGR escape sequences.
func StringWidth(s string) int {
	runes := []rune(s)
	w := 0
	for i := 0; i < len(runes); {
		if runes[i] == '\x1b' {
			if _, n := scanSGR(runes[i:]); n > 0 {
				i += n
				continue
			}
		}
		w += RuneWidth(runes[i])
		i++
	}
	return w
}

// scanSGR reads a CSI SGR escape sequence ("\x1b[<params>m") from the start
// of runes. It returns the raw sequence and the number of runes consumed,
// or ("", 0) if runes doesn't start with one.
func scanSGR(runes []rune) (string, int) {
	if len(runes) < 3 || runes[0] != '\x1b' || runes[1] != '[' {
		return "", 0
	}
	for i := 2; i < len(runes); i++ {
		switch r := runes[i]; {
		case r >= '0' && r <= '9', r == ';':
			continue
		case r == 'm':
			return string(runes[:i+1]), i + 1
		default:
			return "", 0
		}
	}
	return "", 0
}

// isWideRune reports whether r is an East Asian Wide or Fullwidth
// character that occupies two terminal columns.
func isWideRune(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x115F,
		r >= 0x2329 && r <= 0x232A,
		r >= 0x2E80 && r <= 0xA4CF,
		r >= 0xAC00 && r <= 0xD7A3,
		r >= 0xF900 && r <= 0xFAFF,
		r >= 0xFE10 && r <= 0xFE6F,
		r >= 0xFF00 && r <= 0xFF60,
		r >= 0xFFE0 && r <= 0xFFE6,
		r >= 0x1F300 && r <= 0x1FAFF,
		r >= 0x20000 && r <= 0x3FFFD:
		return true
	default:
		return false
	}
}
