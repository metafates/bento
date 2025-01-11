package grapheme

import (
	"unicode"

	"github.com/rivo/uniseg"
)

func New(symbol string) Grapheme {
	return Grapheme{
		symbol: symbol,
		width:  uniseg.StringWidth(symbol),
	}
}

type Grapheme struct {
	symbol string
	width  int
}

func (g Grapheme) String() string {
	return g.symbol
}

func (g Grapheme) Width() int {
	return g.width
}

func (g Grapheme) IsEmpty() bool {
	return g.symbol == ""
}

func (g Grapheme) IsWhitespace() bool {
	symbol := g.symbol

	const (
		zwsp = "\u200b"
		nbsp = "\u00a0"
	)

	if symbol == zwsp {
		return true
	}

	for _, r := range symbol {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return symbol != nbsp
}
