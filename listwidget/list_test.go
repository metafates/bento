package listwidget

import (
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
	"github.com/stretchr/testify/require"
)

func TestSelectedItemEnsuresVisibleOffsetBeforeRange(t *testing.T) {
	items := []Item{
		NewItemString("Item 0"),
		NewItemString("Item 1"),
		NewItemString("Item 2"),
		NewItemString("Item 3"),
		NewItemString("Item 4"),
		NewItemString("Item 5"),
		NewItemString("Item 6"),
	}

	list := NewList(items...).WithHighlightSymbol(">>")
	state := NewState()
	state.Select(1)
	state.SetOffset(3)

	buffer := statefulWidget(list, &state, 10, 3)

	want := textwidget.NewLinesString(
		">>Item 1  ",
		"  Item 2  ",
		"  Item 3  ",
	).NewBuffer()

	require.Equal(t, want, buffer)
}

func statefulWidget(widget List, state *State, width, height int) *bento.Buffer {
	buffer := bento.NewBufferEmpty(bento.Rect{Width: width, Height: height})

	widget.RenderStateful(buffer.Area, buffer, state)

	return buffer
}
