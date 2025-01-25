package mascotwidget

import (
	"strings"

	"github.com/metafates/bento"
	"github.com/muesli/termenv"
	"github.com/rivo/uniseg"
)

var _ bento.Widget = (*Mascot)(nil)

type Mascot struct {
	riceColor   bento.Color
	eyeColor    bento.Color
	noriColor   bento.Color
	borderColor bento.Color

	horizontal, vertical bento.Flex
}

func New() Mascot {
	return Mascot{
		riceColor:   termenv.RGBColor("#FFFFFF"),
		eyeColor:    termenv.ANSIBlue,
		noriColor:   termenv.RGBColor("#000000"),
		borderColor: termenv.RGBColor("#000000"),

		horizontal: bento.FlexCenter,
		vertical:   bento.FlexCenter,
	}
}

func (m Mascot) Render(area bento.Rect, buffer *bento.Buffer) {
	lines := strings.Split(strings.Trim(mascot, "\n"), "\n")

	var width int
	for _, l := range lines {
		width = max(width, uniseg.StringWidth(l))
	}

	height := len(lines)

	if !area.Contains(bento.NewPosition(width, height)) {
		return
	}

	area = bento.NewLayout(bento.ConstraintLen(width)).Horizontal().WithFlex(m.horizontal).Split(area).Unwrap()
	area = bento.NewLayout(bento.ConstraintLen(height)).Vertical().WithFlex(m.vertical).Split(area).Unwrap()

	for y, line := range lines {
		for x := 0; x < len([]rune(line)); x++ {
			posX := area.Left() + x
			posY := area.Top() + y

			position := bento.NewPosition(posX, posY)

			if !area.Contains(position) {
				continue
			}

			cell := buffer.CellAt(position)

			r := []rune(line)[x]

			color, hasColor := m.colorFor(r)

			if hasColor {
				cell.SetFg(color)
				cell.SetBg(color)
			}

			cell.SetSymbol(string(r))
		}
	}
}

func (m Mascot) colorFor(r rune) (bento.Color, bool) {
	switch r {
	case rice:
		return m.riceColor, true
	case eye:
		return m.eyeColor, true
	case nori:
		return m.noriColor, true
	case border:
		return m.borderColor, true
	default:
		return nil, false
	}
}
