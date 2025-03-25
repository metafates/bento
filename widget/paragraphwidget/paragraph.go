package paragraphwidget

import (
	"fmt"

	"github.com/metafates/bento"
	"github.com/metafates/bento/internal/reflow"
	"github.com/metafates/bento/widget/blockwidget"
	"github.com/metafates/bento/widget/textwidget"
)

var _ bento.Widget = (*Paragraph)(nil)

type Paragraph struct {
	Block     *blockwidget.Block
	Style     bento.Style
	Wrap      *Wrap
	Text      textwidget.Text
	Scroll    bento.Position
	Alignment bento.Alignment
}

func NewStr(s string) Paragraph {
	return New(textwidget.NewTextStr(s))
}

func New(text textwidget.Text) Paragraph {
	return Paragraph{
		Block:     nil,
		Style:     bento.NewStyle(),
		Wrap:      nil,
		Text:      text,
		Scroll:    bento.Position{},
		Alignment: bento.AlignmentLeft,
	}
}

func (p Paragraph) WithStyle(style bento.Style) Paragraph {
	p.Style = style
	return p
}

// WithScroll sets the scroll offset for the given paragraph
//
// The scroll offset is a tuple of (y, x) offset. The y offset is the number of lines to
// scroll, and the x offset is the number of characters to scroll. The scroll offset is applied
// after the text is wrapped and aligned.
func (p Paragraph) WithScroll(vertical, horizontal int) Paragraph {
	p.Scroll = bento.NewPosition(horizontal, vertical)
	return p
}

func (p Paragraph) Right() Paragraph {
	p.Alignment = bento.AlignmentRight
	return p
}

func (p Paragraph) Left() Paragraph {
	p.Alignment = bento.AlignmentLeft
	return p
}

func (p Paragraph) Center() Paragraph {
	p.Alignment = bento.AlignmentCenter
	return p
}

func (p Paragraph) WithBlock(block blockwidget.Block) Paragraph {
	p.Block = &block
	return p
}

// WithWrap sets the wrapping configuration for the widget.
//
// See [Wrap] for more information on the different options.
func (p Paragraph) WithWrap(wrap Wrap) Paragraph {
	p.Wrap = &wrap
	return p
}

func (p Paragraph) Wrapped() Paragraph {
	return p.WithWrap(NewWrap().WithTrim(false))
}

func (p Paragraph) Render(area bento.Rect, buffer *bento.Buffer) {
	buffer.SetStyle(area, p.Style)

	if p.Block != nil {
		p.Block.Render(area, buffer)
		area = p.Block.Inner(area)
	}

	p.render(area, buffer)
}

func (p Paragraph) render(textArea bento.Rect, buffer *bento.Buffer) {
	if textArea.IsEmpty() {
		return
	}

	buffer.SetStyle(textArea, p.Style)

	styled := make([]reflow.InputLine, 0, len(p.Text.Lines))

	for _, line := range p.Text.Lines {
		graphemes := line.StyledGraphemes(p.Text.Style)

		alignment := line.Alignment
		if alignment == bento.AlignmentNone {
			alignment = p.Alignment
		}

		styled = append(styled, reflow.InputLine{
			Graphemes: graphemes,
			Alignment: alignment,
		})
	}

	if p.Wrap != nil {
		trim := p.Wrap.Trim

		lineComposer := reflow.NewWordWrapper(styled, textArea.Width, trim)

		p.renderText(&lineComposer, textArea, buffer)
	} else {
		lineComposer := reflow.NewLineTruncator(styled, textArea.Width)
		lineComposer.SetHorizontalOffset(p.Scroll.X)

		p.renderText(&lineComposer, textArea, buffer)
	}
}

func (p Paragraph) renderText(composer reflow.LineComposer, area bento.Rect, buffer *bento.Buffer) {
	var y int

	for {
		currentLine, ok := composer.NextLine()
		if !ok {
			break
		}

		if y >= p.Scroll.Y {
			x := getLineOffset(currentLine.Width, area.Width, currentLine.Alignment)

			for _, grapheme := range currentLine.Line {
				width := grapheme.Width()

				if width == 0 {
					continue
				}

				symbol := grapheme.String()
				if len(symbol) == 0 {
					symbol = " "
				}

				buffer.
					CellAt(bento.Position{
						X: area.Left() + x,
						Y: area.Top() + y - p.Scroll.Y,
					}).
					SetSymbol(symbol).
					SetStyle(grapheme.Style)

				x += width
			}
		}

		y++

		if y >= area.Height+p.Scroll.Y {
			break
		}
	}
}

func getLineOffset(lineWidth, textAreaWidth int, alignment bento.Alignment) int {
	switch alignment {
	case bento.AlignmentCenter:
		return max(0, textAreaWidth/2-lineWidth/2)
	case bento.AlignmentLeft:
		return 0
	case bento.AlignmentRight:
		return max(0, textAreaWidth-lineWidth)
	default:
		panic(fmt.Sprintf("unexpected bento.Alignment: %#v", alignment))
	}
}
