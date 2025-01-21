package tabswidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/symbol"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Widget = (*Tabs)(nil)

type Tabs struct {
	block          *blockwidget.Block
	titles         []textwidget.Line
	selected       *int
	style          bento.Style
	highlightStyle bento.Style
	divider        textwidget.Span
	paddingLeft    textwidget.Line
	paddingRight   textwidget.Line
}

func New(titles ...textwidget.Line) Tabs {
	return Tabs{
		block:          nil,
		titles:         titles,
		selected:       nil,
		style:          bento.NewStyle(),
		highlightStyle: bento.NewStyle().Reversed(),
		divider:        textwidget.NewSpan(symbol.LineVertical),
		paddingLeft:    textwidget.NewLineStr(" "),
		paddingRight:   textwidget.NewLineStr(" "),
	}
}

func (t Tabs) WithStyle(style bento.Style) Tabs {
	t.style = style
	return t
}

func (t Tabs) WithPaddingLeft(padding textwidget.Line) Tabs {
	t.paddingLeft = padding
	return t
}

func (t Tabs) WithPaddingRight(padding textwidget.Line) Tabs {
	t.paddingRight = padding
	return t
}

func (t Tabs) WithDivider(divider textwidget.Span) Tabs {
	t.divider = divider
	return t
}

func (t Tabs) WithHighlightStyle(style bento.Style) Tabs {
	t.highlightStyle = style
	return t
}

func (t Tabs) WithBlock(block blockwidget.Block) Tabs {
	t.block = &block
	return t
}

func (t Tabs) WithTitles(titles ...textwidget.Line) Tabs {
	t.titles = titles
	return t
}

func (t Tabs) Select(index int) Tabs {
	t.selected = &index
	return t
}

func (t Tabs) Unselect() Tabs {
	t.selected = nil
	return t
}

func (t Tabs) Render(area bento.Rect, buffer *bento.Buffer) {
	buffer.SetStyle(area, t.style)

	if t.block != nil {
		t.block.Render(area, buffer)
		area = t.block.Inner(area)
	}

	t.render(area, buffer)
}

func (t Tabs) render(area bento.Rect, buffer *bento.Buffer) {
	if area.IsEmpty() {
		return
	}

	x := area.Left()
	titlesLen := len(t.titles)

	for i, title := range t.titles {
		isLast := i == titlesLen-1

		remainingWidth := max(0, area.Right()-x)
		if remainingWidth == 0 {
			break
		}

		// Left padding
		x, _ = t.paddingLeft.Print(buffer, x, area.Top(), remainingWidth)

		remainingWidth = max(0, area.Right()-x)
		if remainingWidth == 0 {
			break
		}

		newX, _ := title.Print(buffer, x, area.Top(), remainingWidth)

		if t.selected != nil {
			buffer.SetStyle(bento.Rect{
				X:      x,
				Y:      area.Top(),
				Width:  max(0, newX-x),
				Height: 1,
			}, t.highlightStyle)
		}

		x = newX

		remainingWidth = max(0, area.Right()-x)
		if remainingWidth == 0 {
			break
		}

		// Right Padding
		x, _ = t.paddingRight.Print(buffer, x, area.Top(), remainingWidth)

		remainingWidth = max(0, area.Right()-x)
		if remainingWidth == 0 || isLast {
			break
		}

		x, _ = t.divider.Print(buffer, x, area.Top(), remainingWidth)
	}
}
