package ui

// Spacer is a Component with no content of its own, used to reserve empty
// space in a Container — typically flexible space via Item.WithGrow, or a
// fixed gap via WithWidth/WithHeight.
type Spacer struct {
	width  int
	height int
}

func NewSpacer() *Spacer {
	return &Spacer{}
}

var (
	_ Component = (*Spacer)(nil)
	_ Sized     = (*Spacer)(nil)
)

// WithWidth sets the spacer's fixed preferred width.
func (s *Spacer) WithWidth(n int) *Spacer {
	s.width = n
	return s
}

// WithHeight sets the spacer's fixed preferred height.
func (s *Spacer) WithHeight(n int) *Spacer {
	s.height = n
	return s
}

func (s *Spacer) Render(width, height int) []string {
	return nil
}

func (s *Spacer) PreferredWidth() int {
	return s.width
}

func (s *Spacer) PreferredHeight(width int) int {
	return s.height
}
