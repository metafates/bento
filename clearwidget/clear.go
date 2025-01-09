package clearwidget

import "github.com/metafates/bento"

var _ bento.Widget = (*Clear)(nil)

type Clear struct{}

func New() Clear {
	return Clear{}
}

func (Clear) Render(area bento.Rect, buffer *bento.Buffer) {
	for x := area.Left(); x < area.Right(); x++ {
		for y := area.Top(); y < area.Bottom(); y++ {
			buffer.CellAt(bento.NewPosition(x, y)).Reset()
		}
	}
}
