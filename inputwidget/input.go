package inputwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.StatefulWidget[State] = (*Input)(nil)

type Input struct {
	block            *blockwidget.Block
	style            bento.Style
	alignment        bento.Alignment
	placeholder      string
	placeholderStyle bento.Style
	cursorStyle      bento.Style
}

func NewInput() Input {
	return Input{
		block:            nil,
		style:            bento.NewStyle(),
		alignment:        bento.AlignmentLeft,
		placeholder:      "",
		placeholderStyle: bento.NewStyle(),
		cursorStyle:      bento.NewStyle().Reversed(),
	}
}

func (i Input) RenderStateful(area bento.Rect, buffer *bento.Buffer, state State) {
	buffer.SetStyle(area, i.style)

	if i.block != nil {
		i.block.Render(area, buffer)
		area = i.block.Inner(area)
	}

	before, under, after := state.splitAtCursor()
	before = append(before, under)

	cursor := " "
	if len(after) > 0 {
		if !after[0].IsEmpty() {
			cursor = after[0].String()
		}
		after = after[1:]
	}

	line := textwidget.NewLine(
		textwidget.NewSpan(before.String()),
		textwidget.NewSpan(cursor).WithStyle(i.cursorStyle),
		textwidget.NewSpan(after.String()),
	)

	line.Render(area, buffer)
}
