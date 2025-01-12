package inputwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/internal/grapheme"
	"github.com/metafates/bento/internal/sliceutil"
	"github.com/metafates/bento/textwidget"
	"github.com/rivo/uniseg"
)

var _ bento.StatefulWidget[*State] = (*Input)(nil)

type Input struct {
	block            *blockwidget.Block
	style            bento.Style
	alignment        bento.Alignment
	placeholder      grapheme.Graphemes
	placeholderStyle bento.Style
	cursorStyle      bento.Style
	vertical         bento.Flex
	prompt           string
	promptStyle      bento.Style
	showCursor       bool
}

func New() Input {
	return Input{
		block:            nil,
		style:            bento.NewStyle(),
		alignment:        bento.AlignmentLeft,
		placeholder:      nil,
		placeholderStyle: bento.NewStyle().Dim().Italic(),
		cursorStyle:      bento.NewStyle().Reversed(),
		vertical:         bento.FlexLegacy,
		prompt:           "",
		promptStyle:      bento.NewStyle(),
		showCursor:       true,
	}
}

func (i Input) WithPrompt(prompt string) Input {
	i.prompt = prompt
	return i
}

func (i Input) WithPromptStyle(style bento.Style) Input {
	i.promptStyle = style
	return i
}

func (i Input) Focused() Input {
	i.showCursor = true
	return i
}

func (i Input) Blurred() Input {
	i.showCursor = false
	return i
}

func (i Input) Left() Input {
	return i.WithAlignment(bento.AlignmentLeft)
}

func (i Input) Center() Input {
	return i.WithAlignment(bento.AlignmentCenter)
}

func (i Input) Right() Input {
	return i.WithAlignment(bento.AlignmentRight)
}

func (i Input) Top() Input {
	return i.WithVertical(bento.FlexStart)
}

func (i Input) Middle() Input {
	return i.WithVertical(bento.FlexCenter)
}

func (i Input) Bottom() Input {
	return i.WithVertical(bento.FlexEnd)
}

func (i Input) WithVertical(flex bento.Flex) Input {
	i.vertical = flex
	return i
}

func (i Input) WithBlock(block blockwidget.Block) Input {
	i.block = &block
	return i
}

func (i Input) WithPlaceholder(placeholder string) Input {
	i.placeholder = grapheme.NewGraphemes(placeholder)
	return i
}

func (i Input) WithPlaceholderStyle(style bento.Style) Input {
	i.placeholderStyle = style
	return i
}

func (i Input) WithCursorStyle(style bento.Style) Input {
	i.cursorStyle = style
	return i
}

func (i Input) WithAlignment(alignment bento.Alignment) Input {
	i.alignment = alignment
	return i
}

func (i Input) RenderStateful(area bento.Rect, buffer *bento.Buffer, state *State) {
	buffer.SetStyle(area, i.style)

	if i.block != nil {
		i.block.Render(area, buffer)
		area = i.block.Inner(area)
	}

	area = bento.
		NewLayout(bento.ConstraintLength(1)).
		Vertical().
		WithFlex(i.vertical).
		Split(area).
		Unwrap()

	inputWidth := area.Width - uniseg.StringWidth(i.prompt) - 1 // 1 for cursor
	firstVisibleIndex, lastVisibleIndex := state.getBounds(inputWidth)

	state.offset = firstVisibleIndex

	tempState := State{
		cursor:    state.cursor - state.offset,
		graphemes: sliceutil.Take(sliceutil.Skip(state.graphemes, state.offset), lastVisibleIndex-firstVisibleIndex),
		offset:    state.offset,
	}

	before, cursor, after := i.split(tempState)

	if i.showCursor {
		cursor = cursor.WithStylePatch(i.cursorStyle)

		if cursor.Content == "" {
			cursor.Content = " "
		}
	}

	var prompt textwidget.Span

	if i.prompt != "" {
		prompt = textwidget.NewSpan(i.prompt).WithStyle(i.promptStyle)
	}

	line := textwidget.NewLine(
		prompt,
		before,
		cursor,
		after,
	).WithAlignment(i.alignment)

	line.Render(area, buffer)
}

func (i Input) split(state State) (beforeSpan, cursorSpan, afterSpan textwidget.Span) {
	if state.isEmpty() {
		placeholder := i.placeholder
		if len(placeholder) == 0 {
			return textwidget.NewSpan(""), textwidget.NewSpan(""), textwidget.NewSpan("")
		}

		cursor := textwidget.NewSpan(placeholder[0].String()).WithStyle(i.placeholderStyle)
		after := textwidget.NewSpan(placeholder[1:].String()).WithStyle(i.placeholderStyle)

		return textwidget.NewSpan(""), cursor, after
	}

	before, under, after := state.splitAtCursor()

	before = append(before, under)

	var cursor string
	if len(after) > 0 {
		cursor = after[0].String()
		after = after[1:]
	}

	beforeSpan = textwidget.NewSpan(before.String())
	cursorSpan = textwidget.NewSpan(cursor)
	afterSpan = textwidget.NewSpan(after.String())

	return beforeSpan, cursorSpan, afterSpan
}
