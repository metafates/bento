package listwidget

import (
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/internal/sliceutil"
	"github.com/metafates/bento/textwidget"
	"github.com/rivo/uniseg"
)

type Direction int

const (
	DirectionTopToBottom Direction = iota
	DirectionBottomToTop
)

type HighlightSpacing int

const (
	HighlightSpacingWhenSelected HighlightSpacing = iota
	HighlightSpacingAlways
	HighlightSpacingNever
)

func (hs HighlightSpacing) shouldAdd(hasSelection bool) bool {
	switch hs {
	case HighlightSpacingWhenSelected:
		return hasSelection
	case HighlightSpacingAlways:
		return true
	case HighlightSpacingNever:
		return false
	default:
		return false
	}
}

var (
	_ bento.StatefulWidget[*State] = (*List)(nil)
	_ bento.Widget                 = (*List)(nil)
)

type List struct {
	items                 []textwidget.Text
	block                 *blockwidget.Block
	style                 bento.Style
	direction             Direction
	highlightStyle        bento.Style
	highlightSymbol       string
	repeatHighlightSymbol bool
	highlightSpacing      HighlightSpacing
	scrollPadding         int
}

func New(items ...textwidget.Text) List {
	return List{
		items:                 items,
		block:                 nil,
		style:                 bento.NewStyle(),
		direction:             DirectionTopToBottom,
		highlightStyle:        bento.NewStyle(),
		highlightSymbol:       "",
		repeatHighlightSymbol: false,
		highlightSpacing:      HighlightSpacingWhenSelected,
		scrollPadding:         0,
	}
}

func (l List) WithItems(items ...textwidget.Text) List {
	l.items = items
	return l
}

func (l List) WithScrollPadding(padding int) List {
	l.scrollPadding = padding
	return l
}

func (l List) RenderStateful(area bento.Rect, buffer *bento.Buffer, state *State) {
	buffer.SetStyle(area, l.style)

	listArea := area
	if l.block != nil {
		l.block.Render(area, buffer)
		listArea = l.block.Inner(area)
	}

	if listArea.IsEmpty() || len(l.items) == 0 {
		return
	}

	if state.selected != nil {
		selected := min(*state.selected, len(l.items)-1)

		state.selected = &selected
	}

	listHeight := listArea.Height

	firstVisibleIndex, lastVisibleIndex := l.getItemsBounds(state.selected, state.offset, listHeight)

	// NOTE: this changes the state's offset to be the beginning of the now viewable items
	state.offset = firstVisibleIndex

	// Get our set highlighted symbol (if one was set)
	highlightSymbol := l.highlightSymbol
	blankSymbol := strings.Repeat(" ", uniseg.StringWidth(highlightSymbol))

	var currentHeight int

	selectionSpacing := l.highlightSpacing.shouldAdd(state.selected != nil)

	for i, item := range sliceutil.Take(sliceutil.Skip(l.items, state.offset), lastVisibleIndex-firstVisibleIndex) {
		i += state.offset

		var x, y int

		switch l.direction {
		case DirectionBottomToTop:
			currentHeight += item.Height()

			x = listArea.Left()
			y = listArea.Bottom() - currentHeight
		case DirectionTopToBottom:
			x = listArea.Left()
			y = listArea.Top() + currentHeight

			currentHeight += item.Height()
		}

		rowArea := bento.Rect{
			X:      x,
			Y:      y,
			Width:  listArea.Width,
			Height: item.Height(),
		}

		itemStyle := l.style.Patched(item.Style)
		buffer.SetStyle(rowArea, itemStyle)

		var isSelected bool
		if state.selected != nil {
			isSelected = *state.selected == i
		}

		itemArea := rowArea
		if selectionSpacing {
			highlightSymbolWidth := uniseg.StringWidth(l.highlightSymbol)

			itemArea = bento.Rect{
				X:      rowArea.X + highlightSymbolWidth,
				Y:      rowArea.Y,
				Width:  max(0, rowArea.Width-highlightSymbolWidth),
				Height: rowArea.Height,
			}
		}

		item.Render(itemArea, buffer)

		if selectionSpacing {
			for j := 0; j < item.Height(); j++ {
				symbol := blankSymbol
				if isSelected && (j == 0 || l.repeatHighlightSymbol) {
					symbol = highlightSymbol
				}

				buffer.SetStringN(x, y+j, symbol, listArea.Width, itemStyle)
			}
		}

		if isSelected {
			buffer.SetStyle(rowArea, l.highlightStyle)
		}
	}
}

// Render implements bento.Widget.
func (l List) Render(area bento.Rect, buffer *bento.Buffer) {
	l.RenderStateful(area, buffer, new(State))
}

func (l List) WithDirection(direction Direction) List {
	l.direction = direction
	return l
}

func (l List) WithHighlightSpacing(highlightSpacing HighlightSpacing) List {
	l.highlightSpacing = highlightSpacing
	return l
}

func (l List) WithRepeatHighlightSymbol(repeat bool) List {
	l.repeatHighlightSymbol = repeat
	return l
}

func (l List) WithHighlightStyle(style bento.Style) List {
	l.highlightStyle = style
	return l
}

func (l List) WithHighlightSymbol(symbol string) List {
	l.highlightSymbol = symbol
	return l
}

func (l List) WithStyle(style bento.Style) List {
	l.style = style
	return l
}

func (l List) WithBlock(block blockwidget.Block) List {
	l.block = &block
	return l
}

// getItemsBounds given an offset, calculates which items can fit in a given area
func (l List) getItemsBounds(selected *int, offset, maxHeight int) (int, int) {
	offset = min(offset, max(0, len(l.items)-1))

	// NOTE: visible here implies visible in the given area
	firstVisibleIndex := offset
	lastVisibleIndex := offset

	// Current height of all items in the list to render, beginning at the offset
	var heightFromOffset int

	// Calculate the last visible index and total height of the items
	// that will fit in the available space
	for _, item := range sliceutil.Skip(l.items, offset) {
		if heightFromOffset+item.Height() > maxHeight {
			break
		}

		heightFromOffset += item.Height()

		lastVisibleIndex++
	}

	// Get the selected index and apply scroll_padding to it, but still honor the offset if
	// nothing is selected. This allows for the list to stay at a position after selecting
	// None.
	indexToDisplay := offset
	if selected != nil {
		indexToDisplay = l.applyScrollPaddingToSelectedIndex(
			*selected,
			maxHeight,
			firstVisibleIndex,
			lastVisibleIndex,
		)
	}

	// Recall that last_visible_index is the index of what we
	// can render up to in the given space after the offset
	// If we have an item selected that is out of the viewable area (or
	// the offset is still set), we still need to show this item
	for indexToDisplay >= lastVisibleIndex {
		heightFromOffset += l.items[lastVisibleIndex].Height()

		lastVisibleIndex++

		// Now we need to hide previous items since we didn't have space
		// for the selected/offset item
		for heightFromOffset > maxHeight {
			heightFromOffset = max(0, heightFromOffset-l.items[firstVisibleIndex].Height())

			// Remove this item to view by starting at the next item index
			firstVisibleIndex++
		}
	}

	// Here we're doing something similar to what we just did above
	// If the selected item index is not in the viewable area, let's try to show the item
	for indexToDisplay < firstVisibleIndex {
		firstVisibleIndex--

		heightFromOffset += l.items[firstVisibleIndex].Height()

		// Don't show an item if it is beyond our viewable height
		for heightFromOffset > maxHeight {
			lastVisibleIndex--

			heightFromOffset = max(0, heightFromOffset-l.items[lastVisibleIndex].Height())
		}
	}

	return firstVisibleIndex, lastVisibleIndex
}

// applyScrollPaddingToSelectedIndex applies scroll padding to the selected index, reducing the padding value to keep the
// selected item on screen even with items of inconsistent sizes
//
// This function is sensitive to how the bounds checking function handles item height
func (l List) applyScrollPaddingToSelectedIndex(selected int, maxHeight, firstVisibleIndex, lastVisibleIndex int) int {
	lastValidIndex := max(0, len(l.items)-1)

	selected = min(selected, lastValidIndex)

	// The below loop handles situations where the list item sizes may not be consistent,
	// where the offset would have excluded some items that we want to include, or could
	// cause the offset value to be set to an inconsistent value each time we render.
	// The padding value will be reduced in case any of these issues would occur
	scrollPadding := l.scrollPadding

	for scrollPadding > 0 {
		var heightAroundSelected int

		from := max(0, selected-scrollPadding)
		to := min(lastValidIndex, selected+scrollPadding)

		for i := from; i <= to; i++ {
			heightAroundSelected += l.items[i].Height()
		}

		if heightAroundSelected <= maxHeight {
			break
		}

		scrollPadding--
	}

	res := selected

	if min(lastValidIndex, selected+scrollPadding) >= lastVisibleIndex {
		res = selected + scrollPadding
	} else if max(0, selected-scrollPadding) < firstVisibleIndex {
		res = max(0, selected-scrollPadding)
	}

	return min(res, lastValidIndex)
}
