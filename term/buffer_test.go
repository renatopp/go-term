package term

import (
	"errors"
	"strings"
	"testing"
)

func flush(t *testing.T, b *Buffer, lines []string) string {
	t.Helper()
	var out strings.Builder
	if err := b.Flush(lines, &out); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}
	return out.String()
}

func TestBufferFlushFullRepaintOnFirstCall(t *testing.T) {
	b := NewBuffer(3, 2)
	got := flush(t, b, []string{"ab", "c"})
	want := "\x1b[1;1Hab \x1b[2;1Hc  "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushNoChangesWritesNothing(t *testing.T) {
	b := NewBuffer(3, 2)
	flush(t, b, []string{"ab", "c"})
	got := flush(t, b, []string{"ab", "c"})
	if got != "" {
		t.Fatalf("got %q, want empty string for an unchanged frame", got)
	}
}

func TestBufferFlushOnlyWritesChangedColumn(t *testing.T) {
	b := NewBuffer(3, 1)
	flush(t, b, []string{"abc"})
	got := flush(t, b, []string{"abx"})
	want := "\x1b[1;3Hx"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferResizeForcesFullRepaint(t *testing.T) {
	b := NewBuffer(3, 1)
	flush(t, b, []string{"abc"})
	b.Resize(3, 1)
	got := flush(t, b, []string{"abc"})
	want := "\x1b[1;1Habc"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushWideCharacter(t *testing.T) {
	b := NewBuffer(3, 1)
	got := flush(t, b, []string{"中"})
	want := "\x1b[1;1H中 "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushDropsWideCharThatDoesNotFit(t *testing.T) {
	b := NewBuffer(2, 1)
	got := flush(t, b, []string{"a中"})
	want := "\x1b[1;1Ha "
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushCombiningMarkStaysInOneCell(t *testing.T) {
	b := NewBuffer(2, 1)
	got := flush(t, b, []string{"éx"}) // 'e' + combining acute accent, then 'x'
	want := "\x1b[1;1Héx"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushTruncatesToWidth(t *testing.T) {
	b := NewBuffer(2, 1)
	got := flush(t, b, []string{"abcdef"})
	want := "\x1b[1;1Hab"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushRewritesOnStyleChangeAlone(t *testing.T) {
	b := NewBuffer(2, 1)
	flush(t, b, []string{"hi"})
	got := flush(t, b, []string{NewStyle().AsBold(true).Render("hi")})
	want := "\x1b[1;1H\x1b[1mhi\x1b[0m"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBufferFlushPropagatesWriteError(t *testing.T) {
	b := NewBuffer(2, 1)
	wantErr := errors.New("boom")
	err := b.Flush([]string{"ab"}, errWriter{wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got err %v, want %v", err, wantErr)
	}
}

type errWriter struct{ err error }

func (w errWriter) Write(p []byte) (int, error) { return 0, w.err }

func TestWriteRowDiffPullsInPrimaryForChangedPlaceholder(t *testing.T) {
	old := []cell{{text: "A", width: 1}, {text: "B", width: 1}}
	new := []cell{{text: "中", width: 2}, {width: 0}}

	var out strings.Builder
	writeRowDiff(&out, 0, old, new, false)

	want := "\x1b[1;1H中"
	if got := out.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRuneWidth(t *testing.T) {
	cases := []struct {
		name string
		r    rune
		want int
	}{
		{"ascii", 'a', 1},
		{"wide CJK", '中', 2},
		{"combining acute", '́', 0},
		{"zero-width joiner", '‍', 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := runeWidth(c.r); got != c.want {
				t.Fatalf("runeWidth(%q) = %d, want %d", c.r, got, c.want)
			}
		})
	}
}

func TestScanSGR(t *testing.T) {
	seq := "\x1b[1;31m"
	runes := []rune(seq + "x")
	gotSeq, gotN := scanSGR(runes)
	if gotSeq != seq || gotN != len([]rune(seq)) {
		t.Fatalf("got (%q, %d), want (%q, %d)", gotSeq, gotN, seq, len([]rune(seq)))
	}

	if seq, n := scanSGR([]rune("abc")); n != 0 || seq != "" {
		t.Fatalf("got (%q, %d), want (\"\", 0) for non-escape input", seq, n)
	}

	if seq, n := scanSGR([]rune("\x1b[1")); n != 0 || seq != "" {
		t.Fatalf("got (%q, %d), want (\"\", 0) for unterminated sequence", seq, n)
	}
}
