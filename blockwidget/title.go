package blockwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
)

type TitlePosition int

const (
	TitlePositionTop TitlePosition = iota
	TitlePositionBottom
)

type Title struct {
	Alignment bento.Alignment
	Position  *TitlePosition
	Line      textwidget.Line
}
