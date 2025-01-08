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

func NewTitle(line textwidget.Line) Title {
	return Title{
		Line: line,
	}
}

func NewTitleString(title string) Title {
	return NewTitle(textwidget.NewLineString(title))
}

func (t Title) WithPosition(position TitlePosition) Title {
	t.Position = &position
	return t
}

func (t Title) WithAlignment(alignment bento.Alignment) Title {
	t.Alignment = alignment
	return t
}

func (t Title) Top() Title {
	return t.WithPosition(TitlePositionTop)
}

func (t Title) Bottom() Title {
	return t.WithPosition(TitlePositionBottom)
}

func (t Title) Right() Title {
	return t.WithAlignment(bento.AlignmentRight)
}

func (t Title) Left() Title {
	return t.WithAlignment(bento.AlignmentLeft)
}

func (t Title) Center() Title {
	return t.WithAlignment(bento.AlignmentCenter)
}
