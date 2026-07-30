package term

import "strings"

type Style struct {
	bold            bool
	dim             bool
	italic          bool
	underline       bool
	slowBlink       bool
	rapidBlink      bool
	inverse         bool
	hidden          bool
	strikeThrough   bool
	fraktur         bool
	doubleUnderline bool
	framed          bool
	encircled       bool
	overline        bool
	foreground      *Color
	background      *Color
}

func NewStyle() *Style {
	return &Style{}
}

func (s Style) Bold() bool                   { return s.bold }
func (s Style) Dim() bool                    { return s.dim }
func (s Style) Italic() bool                 { return s.italic }
func (s Style) Underline() bool              { return s.underline }
func (s Style) SlowBlink() bool              { return s.slowBlink }
func (s Style) RapidBlink() bool             { return s.rapidBlink }
func (s Style) Inverse() bool                { return s.inverse }
func (s Style) Hidden() bool                 { return s.hidden }
func (s Style) StrikeThrough() bool          { return s.strikeThrough }
func (s Style) Fraktur() bool                { return s.fraktur }
func (s Style) DoubleUnderline() bool        { return s.doubleUnderline }
func (s Style) Framed() bool                 { return s.framed }
func (s Style) Encircled() bool              { return s.encircled }
func (s Style) Overline() bool               { return s.overline }
func (s Style) Foreground() *Color           { return s.foreground }
func (s Style) Background() *Color           { return s.background }
func (s Style) AsBold() Style                { s.bold = true; return s }
func (s Style) AsDim() Style                 { s.dim = true; return s }
func (s Style) AsItalic() Style              { s.italic = true; return s }
func (s Style) AsUnderline() Style           { s.underline = true; return s }
func (s Style) AsSlowBlink() Style           { s.slowBlink = true; return s }
func (s Style) AsRapidBlink() Style          { s.rapidBlink = true; return s }
func (s Style) AsInverse() Style             { s.inverse = true; return s }
func (s Style) AsHidden() Style              { s.hidden = true; return s }
func (s Style) AsStrikeThrough() Style       { s.strikeThrough = true; return s }
func (s Style) AsFraktur() Style             { s.fraktur = true; return s }
func (s Style) AsDoubleUnderline() Style     { s.doubleUnderline = true; return s }
func (s Style) AsFramed() Style              { s.framed = true; return s }
func (s Style) AsEncircled() Style           { s.encircled = true; return s }
func (s Style) AsOverline() Style            { s.overline = true; return s }
func (s Style) WithForeground(c Color) Style { s.foreground = &c; return s }
func (s Style) WithBackground(c Color) Style { s.background = &c; return s }
func (s Style) WithoutForeground() Style     { s.foreground = nil; return s }
func (s Style) WithoutBackground() Style     { s.background = nil; return s }
func (s Style) WithoutColors() Style         { s.foreground = nil; s.background = nil; return s }

func (s Style) Render(v string) string {
	var codes []string

	if s.bold {
		codes = append(codes, "1")
	}
	if s.dim {
		codes = append(codes, "2")
	}
	if s.italic {
		codes = append(codes, "3")
	}
	if s.underline {
		codes = append(codes, "4")
	}
	if s.slowBlink {
		codes = append(codes, "5")
	}
	if s.rapidBlink {
		codes = append(codes, "6")
	}
	if s.inverse {
		codes = append(codes, "7")
	}
	if s.hidden {
		codes = append(codes, "8")
	}
	if s.strikeThrough {
		codes = append(codes, "9")
	}
	if s.fraktur {
		codes = append(codes, "20")
	}
	if s.doubleUnderline {
		codes = append(codes, "21")
	}
	if s.framed {
		codes = append(codes, "51")
	}
	if s.encircled {
		codes = append(codes, "52")
	}
	if s.overline {
		codes = append(codes, "53")
	}

	if level := GetColorLevel(); level != ColorModeNone {
		if s.foreground != nil {
			codes = append(codes, colorSequence(*s.foreground, level, false))
		}
		if s.background != nil {
			codes = append(codes, colorSequence(*s.background, level, true))
		}
	}

	if len(codes) == 0 {
		return v
	}
	return "\x1b[" + strings.Join(codes, ";") + "m" + v + "\x1b[0m"
}

func colorSequence(c Color, level ColorMode, bg bool) string {
	for c.Mode() > level {
		c = c.Fallback()
	}
	return c.Sequence(bg)
}
