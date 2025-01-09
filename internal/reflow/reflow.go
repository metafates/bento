package reflow

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
)

type LineComposer interface {
	NextLine() (WrappedLine, bool)
}

type WrappedLine struct {
	Line      []textwidget.StyledGrapheme
	Width     int
	Alignment bento.Alignment
}

type InputLine struct {
	Graphemes []textwidget.StyledGrapheme
	Alignment bento.Alignment
}
