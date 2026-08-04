package term

import (
	"fmt"
	"io"
	"strings"
)

// buffer holds the last frame written to a terminal so that Flush can diff
// the next frame against it and write only the columns that changed.
type buffer struct {
	width  int
	height int
	rows   [][]Cell
}

// newBuffer creates a buffer for a terminal of the given size. The first
// call to Flush always performs a full repaint.
func newBuffer(width, height int) *buffer {
	return &buffer{width: width, height: height}
}

// Resize changes the buffer's dimensions and discards the stored frame,
// forcing a full repaint on the next Flush.
func (b *buffer) Resize(width, height int) {
	b.width = width
	b.height = height
	b.rows = nil
}

// Flush renders lines (one string per row, which may contain SGR escape
// sequences as produced by Style.Render) into the buffer, diffs it against
// the previously flushed frame, and writes only the changed columns to w as
// cursor moves and styled text.
func (b *buffer) Flush(lines []string, w io.Writer) error {
	full := b.rows == nil

	rows := make([][]Cell, b.height)
	for y := range rows {
		var line string
		if y < len(lines) {
			line = lines[y]
		}
		rows[y] = BuildRow(line, b.width)
	}

	var out strings.Builder
	for y, row := range rows {
		var old []Cell
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

// writeRowDiff writes the changed columns of row y to out, comparing new
// against old. full forces every column to be treated as changed, which is
// used for the first Flush after newBuffer/Resize.
func writeRowDiff(out *strings.Builder, y int, old, new []Cell, full bool) {
	x := 0
	for x < len(new) {
		if !full && x < len(old) && old[x] == new[x] {
			x++
			continue
		}

		// A wide character's placeholder cell can't be rewritten on its
		// own, so a change starting on one pulls in the character before it.
		start := x
		if new[start].Width == 0 && start > 0 {
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
func writeRun(out *strings.Builder, y, x int, run []Cell) {
	fmt.Fprintf(out, "\x1b[%d;%dH", y+1, x+1)

	style := ""
	for _, c := range run {
		if c.Width == 0 {
			continue
		}
		if c.Style != style {
			// Styles describe all attributes from a default state, so reset
			// before switching between two styles rather than stacking them.
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
}
