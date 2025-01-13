package listwidget

import (
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
	"github.com/stretchr/testify/require"
)

func TestSelectedItemEnsuresVisibleOffsetBeforeRange(t *testing.T) {
	items := []Item{
		NewItemStr("Item 0"),
		NewItemStr("Item 1"),
		NewItemStr("Item 2"),
		NewItemStr("Item 3"),
		NewItemStr("Item 4"),
		NewItemStr("Item 5"),
		NewItemStr("Item 6"),
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
