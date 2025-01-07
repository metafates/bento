package bento

import "github.com/muesli/termenv"

type Color = termenv.Color

var _ Color = (*ResetColor)(nil)

type ResetColor struct{}

// Sequence implements termenv.Color.
func (ResetColor) Sequence(bg bool) string {
	if bg {
		return "49"
	}

	return "39"
}
