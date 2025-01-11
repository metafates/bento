package grapheme

import (
	"strings"

	"github.com/rivo/uniseg"
)

type Graphemes []Grapheme

func NewGraphemes(s string) Graphemes {
	g := uniseg.NewGraphemes(s)

	graphemes := make(Graphemes, 0, len([]rune(s)))

	for g.Next() {
		graphemes = append(graphemes, Grapheme{
			symbol: g.Str(),
			width:  g.Width(),
		})
	}

	return graphemes
}

func (gs Graphemes) Width() int {
	var width int

	for _, g := range gs {
		width += g.width
	}

	return width
}

func (gs Graphemes) String() string {
	var b strings.Builder

	for _, g := range gs {
		b.Grow(len(g.String()))
	}

	for _, g := range gs {
		b.WriteString(g.String())
	}

	return b.String()
}
