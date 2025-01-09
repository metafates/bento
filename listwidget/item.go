package listwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
)

type Item struct {
	style   bento.Style
	content textwidget.Text
}

func NewItemString(text string) Item {
	return NewItem(textwidget.NewTextString(text))
}

func NewItem(text textwidget.Text) Item {
	return Item{
		style:   bento.NewStyle(),
		content: text,
	}
}

func (i Item) WithStyle(style bento.Style) Item {
	i.style = style
	return i
}

func (i Item) Height() int {
	return i.content.Height()
}

func (i Item) Width() int {
	return i.content.Width()
}
