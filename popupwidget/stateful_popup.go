package popupwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/clearwidget"
)

var _ bento.StatefulWidget[any] = (*StatefulPopup[any])(nil)

type StatefulPopup[S any] struct {
	Block                *blockwidget.Block
	Horizontal, Vertical bento.Flex
	Style                bento.Style
	Width, Height        bento.Constraint
	Content              bento.StatefulWidget[S]
	Padding              bento.Padding
}

func NewStateful[S any](content bento.StatefulWidget[S]) StatefulPopup[S] {
	return StatefulPopup[S]{
		Block:      nil,
		Horizontal: bento.FlexCenter,
		Vertical:   bento.FlexCenter,
		Style:      bento.NewStyle(),
		Width:      bento.ConstraintPercentage(60),
		Height:     bento.ConstraintPercentage(30),
		Padding:    bento.NewPadding(),
		Content:    content,
	}
}

func (p StatefulPopup[S]) WithBlock(block blockwidget.Block) StatefulPopup[S] {
	p.Block = &block
	return p
}

func (p StatefulPopup[S]) WithStyle(style bento.Style) StatefulPopup[S] {
	p.Style = style
	return p
}

func (p StatefulPopup[S]) WithHeight(height bento.Constraint) StatefulPopup[S] {
	p.Height = height
	return p
}

func (p StatefulPopup[S]) WithWidth(width bento.Constraint) StatefulPopup[S] {
	p.Width = width
	return p
}

func (p StatefulPopup[S]) WithPadding(padding bento.Padding) StatefulPopup[S] {
	p.Padding = padding
	return p
}

func (p StatefulPopup[S]) Top() StatefulPopup[S] {
	p.Vertical = bento.FlexStart
	return p
}

func (p StatefulPopup[S]) Middle() StatefulPopup[S] {
	p.Vertical = bento.FlexCenter
	return p
}

func (p StatefulPopup[S]) Bottom() StatefulPopup[S] {
	p.Vertical = bento.FlexEnd
	return p
}

func (p StatefulPopup[S]) Left() StatefulPopup[S] {
	p.Horizontal = bento.FlexStart
	return p
}

func (p StatefulPopup[S]) Center() StatefulPopup[S] {
	p.Horizontal = bento.FlexCenter
	return p
}

func (p StatefulPopup[S]) Right() StatefulPopup[S] {
	p.Horizontal = bento.FlexEnd
	return p
}

func (p StatefulPopup[S]) RenderStateful(area bento.Rect, buffer *bento.Buffer, state S) {
	vertical := bento.NewLayout(p.Height).Vertical().WithPadding(p.Padding).WithFlex(p.Vertical)
	horizontal := bento.NewLayout(p.Width).Horizontal().WithPadding(p.Padding).WithFlex(p.Horizontal)

	area = vertical.Split(area).Unwrap()
	area = horizontal.Split(area).Unwrap()

	clearwidget.New().Render(area, buffer)
	buffer.SetStyle(area, p.Style)

	if p.Block != nil {
		p.Block.Render(area, buffer)
		area = p.Block.Inner(area)
	}

	if area.IsEmpty() {
		return
	}

	p.Content.RenderStateful(area, buffer, state)
}
