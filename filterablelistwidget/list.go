package filterablelistwidget

import (
	"fmt"
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/inputwidget"
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
	_ bento.StatefulWidget[*State[StringItem]] = (*List[StringItem])(nil)
	_ bento.Widget                             = (*List[StringItem])(nil)
)

type List[I Item] struct {
	block                 *blockwidget.Block
	style                 bento.Style
	direction             Direction
	highlightStyle        bento.Style
	highlightSymbol       string
	repeatHighlightSymbol bool
	highlightSpacing      HighlightSpacing
	scrollPadding         int
}

func New[I Item]() List[I] {
	return List[I]{
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

func (l List[I]) WithScrollPadding(padding int) List[I] {
	l.scrollPadding = padding
	return l
}

func (l List[I]) RenderStateful(area bento.Rect, buffer *bento.Buffer, state *State[I]) {
	if state.filterState != FilterStateNoFilter {
		var inputArea bento.Rect

		bento.NewLayout(
			bento.ConstraintLength(3),
			bento.ConstraintFill(1),
		).Vertical().Split(area).Assign(&inputArea, &area)

		l.renderFilter(inputArea, buffer, state)
	}

	l.render(area, buffer, state)
}

func (l List[I]) renderFilter(area bento.Rect, buffer *bento.Buffer, state *State[I]) {
	title := "Filter"

	if state.filterState == FilterStateFiltered {
		title = fmt.Sprintf("%s %q", title, state.filterInput.String())
	}

	block := blockwidget.New().Bordered().WithTitleStr(title)

	inputwidget.
		New().
		WithBlock(block).
		WithPlaceholder("Search").
		RenderStateful(area, buffer, &state.filterInput)
}

func (l List[I]) render(area bento.Rect, buffer *bento.Buffer, state *State[I]) {
	buffer.SetStyle(area, l.style)

	listArea := area
	if l.block != nil {
		l.block.Render(area, buffer)
		listArea = l.block.Inner(area)
	}

	if listArea.IsEmpty() || len(state.filteredIndices) == 0 {
		return
	}

	listHeight := listArea.Height

	items := make([]textwidget.Text, 0, len(state.filteredIndices))
	for _, i := range state.filteredIndices {
		items = append(items, state.items[i].Title())
	}

	firstVisibleIndex, lastVisibleIndex := l.getItemsBounds(state.selected, items, state.offset, listHeight)

	// NOTE: this changes the state's offset to be the beginning of the now viewable items
	state.offset = firstVisibleIndex

	// Get our set highlighted symbol (if one was set)
	highlightSymbol := l.highlightSymbol
	blankSymbol := strings.Repeat(" ", uniseg.StringWidth(highlightSymbol))

	var currentHeight int

	selectionSpacing := l.highlightSpacing.shouldAdd(state.selected != nil)

	for i, item := range sliceutil.Take(sliceutil.Skip(state.filteredIndices, state.offset), lastVisibleIndex-firstVisibleIndex) {
		i += state.offset

		var x, y int

		text := state.items[item].Title()

		switch l.direction {
		case DirectionBottomToTop:
			currentHeight += text.Height()

			x = listArea.Left()
			y = listArea.Bottom() - currentHeight
		case DirectionTopToBottom:
			x = listArea.Left()
			y = listArea.Top() + currentHeight

			currentHeight += text.Height()
		}

		rowArea := bento.Rect{
			X:      x,
			Y:      y,
			Width:  listArea.Width,
			Height: text.Height(),
		}

		itemStyle := l.style.Patched(text.Style)
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

		text.Render(itemArea, buffer)

		if selectionSpacing {
			for j := 0; j < text.Height(); j++ {
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
func (l List[I]) Render(area bento.Rect, buffer *bento.Buffer) {
	l.RenderStateful(area, buffer, new(State[I]))
}

func (l List[I]) WithDirection(direction Direction) List[I] {
	l.direction = direction
	return l
}

func (l List[I]) WithHighlightSpacing(highlightSpacing HighlightSpacing) List[I] {
	l.highlightSpacing = highlightSpacing
	return l
}

func (l List[I]) WithRepeatHighlightSymbol(repeat bool) List[I] {
	l.repeatHighlightSymbol = repeat
	return l
}

func (l List[I]) WithHighlightStyle(style bento.Style) List[I] {
	l.highlightStyle = style
	return l
}

func (l List[I]) WithHighlightSymbol(symbol string) List[I] {
	l.highlightSymbol = symbol
	return l
}

func (l List[I]) WithStyle(style bento.Style) List[I] {
	l.style = style
	return l
}

func (l List[I]) WithBlock(block blockwidget.Block) List[I] {
	l.block = &block
	return l
}

// getItemsBounds given an offset, calculates which items can fit in a given area
func (l List[I]) getItemsBounds(selected *int, items []textwidget.Text, offset, maxHeight int) (int, int) {
	offset = min(offset, max(0, len(items)-1))

	// NOTE: visible here implies visible in the given area
	firstVisibleIndex := offset
	lastVisibleIndex := offset

	// Current height of all items in the list to render, beginning at the offset
	var heightFromOffset int

	// Calculate the last visible index and total height of the items
	// that will fit in the available space
	for _, item := range sliceutil.Skip(items, offset) {
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
			items,
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
		heightFromOffset += items[lastVisibleIndex].Height()

		lastVisibleIndex++

		// Now we need to hide previous items since we didn't have space
		// for the selected/offset item
		for heightFromOffset > maxHeight {
			heightFromOffset = max(0, heightFromOffset-items[firstVisibleIndex].Height())

			// Remove this item to view by starting at the next item index
			firstVisibleIndex++
		}
	}

	// Here we're doing something similar to what we just did above
	// If the selected item index is not in the viewable area, let's try to show the item
	for indexToDisplay < firstVisibleIndex {
		firstVisibleIndex--

		heightFromOffset += items[firstVisibleIndex].Height()

		// Don't show an item if it is beyond our viewable height
		for heightFromOffset > maxHeight {
			lastVisibleIndex--

			heightFromOffset = max(0, heightFromOffset-items[lastVisibleIndex].Height())
		}
	}

	return firstVisibleIndex, lastVisibleIndex
}

// applyScrollPaddingToSelectedIndex applies scroll padding to the selected index, reducing the padding value to keep the
// selected item on screen even with items of inconsistent sizes
//
// This function is sensitive to how the bounds checking function handles item height
func (l List[I]) applyScrollPaddingToSelectedIndex(selected int, items []textwidget.Text, maxHeight, firstVisibleIndex, lastVisibleIndex int) int {
	lastValidIndex := max(0, len(items)-1)

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
			heightAroundSelected += items[i].Height()
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
