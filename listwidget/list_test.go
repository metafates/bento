package listwidget

import (
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
	"github.com/stretchr/testify/require"
)

func TestSelectedItemEnsuresVisibleOffsetBeforeRange(t *testing.T) {
	items := []textwidget.Text{
		textwidget.NewTextStr("Item 0"),
		textwidget.NewTextStr("Item 1"),
		textwidget.NewTextStr("Item 2"),
		textwidget.NewTextStr("Item 3"),
		textwidget.NewTextStr("Item 4"),
		textwidget.NewTextStr("Item 5"),
		textwidget.NewTextStr("Item 6"),
	}

	list := New(items...).WithHighlightSymbol(">>")
	state := NewState()
	state.Select(1)
	state.SetOffset(3)

	buffer := statefulWidget(list, &state, 10, 3)

	want := textwidget.NewLinesStr(
		">>Item 1  ",
		"  Item 2  ",
		"  Item 3  ",
	).NewBuffer()

	require.Equal(t, want, buffer)
}

func statefulWidget(widget List, state *State, width, height int) bento.Buffer {
	buffer := bento.NewBufferEmpty(bento.Rect{Width: width, Height: height})

	widget.RenderStateful(buffer.Area(), &buffer, state)

	return buffer
}
