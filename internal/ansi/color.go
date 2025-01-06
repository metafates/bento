package ansi

type Color interface {
	Sequence(bg bool) string
}

var _ Color = (*ResetColor)(nil)

type ResetColor struct{}

// Sequence implements termenv.Color.
func (r ResetColor) Sequence(bg bool) string {
	if bg {
		return "49"
	}

	return "39"
}
