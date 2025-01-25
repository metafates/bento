package gradient

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/metafates/bento"
	"github.com/muesli/termenv"
)

var _ bento.Widget = (*Gradient)(nil)

type Gradient struct{}

func New() Gradient {
	return Gradient{}
}

func (Gradient) Render(area bento.Rect, buffer *bento.Buffer) {
	var i int

	for y := area.Top(); y < area.Bottom(); y++ {

		value := area.Height - i
		valueFg := float64(value) / float64(area.Height)
		valueBg := (float64(value) - 0.5) / float64(area.Height)

		var j int
		for x := area.Left(); x < area.Right(); x++ {
			hue := float64(j*360) / float64(area.Width)

			fg := colorFromOklab(hue, 1, valueFg)
			bg := colorFromOklab(hue, 1, valueBg)

			buffer.CellAt(bento.NewPosition(x, y)).SetSymbol("▀").SetFg(fg).SetBg(bg)

			j++
		}

		i++
	}
}

func colorFromOklab(hue, saturation, value float64) bento.Color {
	color := colorful.Hsv(hue, saturation, value)

	return termenv.RGBColor(color.Hex())
}
