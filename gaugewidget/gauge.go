package gaugewidget

import (
	"math"
	"strconv"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/symbol"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Widget = (*Gauge)(nil)

type Gauge struct {
	block      *blockwidget.Block
	ratio      float64
	label      *textwidget.Span
	useUnicode bool
	style      bento.Style
	gaugeStyle bento.Style
}

func New() Gauge {
	return Gauge{
		block:      nil,
		ratio:      0,
		label:      nil,
		useUnicode: false,
		style:      bento.NewStyle(),
		gaugeStyle: bento.NewStyle(),
	}
}

// WithRatio sets the bar progression from a ratio (float).
//
// `ratio` is the ratio between filled bar over empty bar (i.e. `3/4` completion is `0.75`).
// This is more easily seen as a floating point percentage (e.g. 42% = `0.42`).
//
// This method panics if `ratio` is not between 0 and 1 inclusively.
//
// See [Gauge.WithPercent] to set from a percentage.
func (g Gauge) WithRatio(ratio float64) Gauge {
	if ratio < 0 && ratio > 1 {
		panic("ratio should be between 0 and 1 inclusively")
	}

	g.ratio = ratio
	return g
}

// WithPercent sets the bar progression from a percentage.
//
// This method panics if `percent` is **not** between 0 and 100 inclusively.
//
// # See also
//
// See [Gauge.WithRatio] to set from a float.
func (g Gauge) WithPercent(percent int) Gauge {
	if percent < 0 || percent > 100 {
		panic("percentage should be between 0 and 100 inclusively")
	}

	g.ratio = float64(percent) / 100.0
	return g
}

// WithLabel sets the label to display in the center of the bar.
//
// If the label is not defined, it is the percentage filled.
func (g Gauge) WithLabel(label textwidget.Span) Gauge {
	g.label = &label
	return g
}

// WithLabelStr is the same as [Gauge.WithLabel] but accepts string
func (g Gauge) WithLabelStr(label string) Gauge {
	span := textwidget.NewSpan(label)
	g.label = &span
	return g
}

// WithStyle sets the widget style.
func (g Gauge) WithStyle(style bento.Style) Gauge {
	g.style = style
	return g
}

// WithGaugeStyle sets the style of the bar.
func (g Gauge) WithGaugeStyle(style bento.Style) Gauge {
	g.gaugeStyle = style
	return g
}

// WithUnicode sets whether to use unicode characters to display the progress bar.
//
// This enables the use of
// [unicode block characters](https://en.wikipedia.org/wiki/Block_Elements).
// This is useful to display a higher precision bar (8 extra fractional parts per cell).
func (g Gauge) WithUnicode(unicode bool) Gauge {
	g.useUnicode = unicode
	return g
}

func (g Gauge) Render(area bento.Rect, buffer *bento.Buffer) {
	buffer.SetStyle(area, g.style)

	if g.block != nil {
		g.block.Render(area, buffer)
		area = g.block.Inner(area)
	}

	g.render(area, buffer)
}

func (g Gauge) render(area bento.Rect, buffer *bento.Buffer) {
	if area.IsEmpty() {
		return
	}

	buffer.SetStyle(area, g.gaugeStyle)

	label := textwidget.NewSpan(strconv.Itoa(int(math.Round(g.ratio*100))) + "%")
	if g.label != nil {
		label = *g.label
	}

	clampedLabelWidth := min(area.Width, label.Width())
	labelCol := area.Left() + (area.Width-clampedLabelWidth)/2
	labelRow := area.Top() + area.Height/2

	filledWidth := float64(area.Width) * g.ratio

	end := area.Left() + int(math.Round(filledWidth))
	if g.useUnicode {
		end = area.Left() + int(math.Floor(filledWidth))
	}

	var fg, bg bento.Color = bento.ResetColor{}, bento.ResetColor{}
	if g.gaugeStyle.Foreground.IsSet() {
		fg = g.gaugeStyle.Foreground.Color()
	}

	if g.gaugeStyle.Background.IsSet() {
		bg = g.gaugeStyle.Background.Color()
	}

	frac := math.Mod(filledWidth, 1)

	for y := area.Top(); y < area.Bottom(); y++ {
		for x := area.Left(); x < end; x++ {
			pos := bento.NewPosition(x, y)

			if x < labelCol || x >= labelCol+clampedLabelWidth || y != labelRow {
				buffer.CellAt(pos).SetSymbol(symbol.BlockFull).SetFg(fg).SetBg(bg)
			} else {
				buffer.CellAt(pos).SetSymbol(" ").SetFg(fg).SetBg(bg)
			}
		}

		if g.useUnicode && g.ratio < 1 {
			buffer.CellAt(bento.NewPosition(end, y)).SetSymbol(getUnicodeBlock(frac))
		}
	}

	buffer.SetStringN(labelCol, labelRow, label.Content, clampedLabelWidth, label.Style)
}

func getUnicodeBlock(frac float64) string {
	switch int(math.Round(frac * 8.0)) {
	case 1:
		return symbol.BlockOneEighth
	case 2:
		return symbol.BlockOneQuarter
	case 3:
		return symbol.BlockThreeEighths
	case 4:
		return symbol.BlockHalf
	case 5:
		return symbol.BlockFiveEighths
	case 6:
		return symbol.BlockThreeQuarters
	case 7:
		return symbol.BlockSevenEighths
	case 8:
		return symbol.BlockFull
	default:
		return " "
	}
}
