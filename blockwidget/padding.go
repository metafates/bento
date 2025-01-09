package blockwidget

import "github.com/metafates/bento"

type Padding bento.Margin

func NewPadding(sides ...int) Padding {
	return Padding(bento.NewMargin(sides...))
}

func (p Padding) WithTop(s int) Padding {
	return Padding(bento.Margin(p).WithTop(s))
}

func (p Padding) WithRight(s int) Padding {
	return Padding(bento.Margin(p).WithRight(s))
}

func (p Padding) WithBottom(s int) Padding {
	return Padding(bento.Margin(p).WithBottom(s))
}

func (p Padding) WithLeft(s int) Padding {
	return Padding(bento.Margin(p).WithLeft(s))
}
