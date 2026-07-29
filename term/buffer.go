package term

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// cell is a single terminal column: its visible content, the number of
// columns it occupies (0 for the placeholder following a wide character, 1
// for a normal character, 2 for a wide one), and the raw SGR sequence
// active when it was written ("" for no styling).
type cell struct {
	text  string
	width int
	style string
}

// Buffer holds the last frame written to a terminal so that Flush can diff
// the next frame against it and write only the columns that changed.
type Buffer struct {
	width  int
	height int
	rows   [][]cell
}

// NewBuffer creates a Buffer for a terminal of the given size. The first
// call to Flush always performs a full repaint.
func NewBuffer(width, height int) *Buffer {
	return &Buffer{width: width, height: height}
}

// Resize changes the buffer's dimensions and discards the stored frame,
// forcing a full repaint on the next Flush.
func (b *Buffer) Resize(width, height int) {
	b.width = width
	b.height = height
	b.rows = nil
}

// Flush renders lines (one string per row, which may contain SGR escape
// sequences as produced by Style.Render) into the buffer, diffs it against
// the previously flushed frame, and writes only the changed columns to w as
// cursor moves and styled text.
func (b *Buffer) Flush(lines []string, w io.Writer) error {
	full := b.rows == nil

	rows := make([][]cell, b.height)
	for y := range rows {
		var line string
		if y < len(lines) {
			line = lines[y]
		}
		rows[y] = buildRow(line, b.width)
	}

	var out strings.Builder
	for y, row := range rows {
		var old []cell
		if !full {
			old = b.rows[y]
		}
		writeRowDiff(&out, y, old, row, full)
	}

	b.rows = rows
	if out.Len() == 0 {
		return nil
	}
	_, err := io.WriteString(w, out.String())
	return err
}

// buildRow scans line into exactly width columns, truncating content that
// overflows and padding short lines with blank cells.
func buildRow(line string, width int) []cell {
	row := make([]cell, 0, width)
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

		w := runeWidth(r)
		if w == 0 {
			if len(row) > 0 {
				row[len(row)-1].text += string(r)
			}
			i++
			continue
		}
		if x+w > width {
			break
		}

		row = append(row, cell{text: string(r), width: w, style: style})
		if w == 2 {
			row = append(row, cell{width: 0})
		}
		x += w
		i++
	}

	for x < width {
		row = append(row, cell{text: " ", width: 1})
		x++
	}
	return row
}

// writeRowDiff writes the changed columns of row y to out, comparing new
// against old. full forces every column to be treated as changed, which is
// used for the first Flush after NewBuffer/Resize.
func writeRowDiff(out *strings.Builder, y int, old, new []cell, full bool) {
	x := 0
	for x < len(new) {
		if !full && x < len(old) && old[x] == new[x] {
			x++
			continue
		}

		// A wide character's placeholder cell can't be rewritten on its
		// own, so a change starting on one pulls in the character before it.
		start := x
		if new[start].width == 0 && start > 0 {
			start--
		}

		end := start + 1
		for end < len(new) && (full || end >= len(old) || old[end] != new[end]) {
			end++
		}

		writeRun(out, y, start, new[start:end])
		x = end
	}
}

// writeRun positions the cursor at row y, column x and writes the run's
// cells, emitting an SGR sequence only where the active style changes.
func writeRun(out *strings.Builder, y, x int, run []cell) {
	fmt.Fprintf(out, "\x1b[%d;%dH", y+1, x+1)

	style := ""
	for _, c := range run {
		if c.width == 0 {
			continue
		}
		if c.style != style {
			if c.style == "" {
				out.WriteString("\x1b[0m")
			} else {
				out.WriteString(c.style)
			}
			style = c.style
		}
		out.WriteString(c.text)
	}
	if style != "" {
		out.WriteString("\x1b[0m")
	}
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

// runeWidth returns the number of terminal columns r occupies: 0 for
// combining marks and joiners that attach to the previous cell, 2 for wide
// East Asian characters, 1 otherwise.
func runeWidth(r rune) int {
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
