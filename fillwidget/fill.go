package fillwidget

import "github.com/metafates/bento"

var _ bento.Widget = (*Fill)(nil)

type Fill struct {
	Symbol string
	Style  bento.Style
}

func New(symbol string) Fill {
	return Fill{
		Symbol: symbol,
		Style:  bento.NewStyle(),
	}
}

func (f Fill) WithStyle(style bento.Style) Fill {
	f.Style = style
	return f
}

func (f Fill) Render(area bento.Rect, buffer *bento.Buffer) {
	for x := area.Left(); x < area.Right(); x++ {
		for y := area.Top(); y < area.Bottom(); y++ {
			pos := bento.NewPosition(x, y)

			cell := buffer.CellAt(pos)

			cell.Reset()
			cell.SetStyle(f.Style)
			cell.SetSymbol(f.Symbol)
		}
	}
}
