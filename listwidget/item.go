package listwidget

import (
	"github.com/metafates/bento/textwidget"
)

type Item interface {
	Title() textwidget.Text
}

type FilterableItem interface {
	Item

	FilterValue() string
}

type StringItem string

func (i StringItem) Title() textwidget.Text {
	return textwidget.NewTextStr(string(i))
}

func (i StringItem) FilterValue() string {
	return string(i)
}

func NewStringItem(s string) StringItem {
	return StringItem(s)
}
