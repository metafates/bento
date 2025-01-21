package treewidget

import (
	"slices"
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/internal/sliceutil"
	"github.com/metafates/bento/scrollwidget"
	"github.com/metafates/bento/symbol"
	"github.com/rivo/uniseg"
)

var _ bento.StatefulWidget[State[int]] = (*Tree[int])(nil)

type Tree[T comparable] struct {
	items     []Item[T]
	block     *blockwidget.Block
	scrollbar *scrollwidget.Scroll

	style bento.Style

	highlightStyle  bento.Style
	highlightSymbol string

	nodeClosedSymbol     string
	nodeOpenSymbol       string
	nodeNoChildrenSymbol string
}

func New[T comparable](items ...Item[T]) Tree[T] {
	return Tree[T]{
		items:                items,
		block:                nil,
		scrollbar:            nil,
		style:                bento.NewStyle(),
		highlightStyle:       bento.NewStyle(),
		highlightSymbol:      "",
		nodeClosedSymbol:     symbol.ArrowRight,
		nodeOpenSymbol:       symbol.ArrowDown,
		nodeNoChildrenSymbol: "  ",
	}
}

func (t Tree[T]) RenderStateful(area bento.Rect, buffer *bento.Buffer, state State[T]) {
	buffer.SetStyle(area, t.style)

	fullArea := area

	if t.block != nil {
		t.block.Render(area, buffer)
		area = t.block.Inner(area)
	}

	state.lastArea = area
	clear(state.lastRenderedIDs)

	if area.Width < 1 || area.Height < 1 {
		return
	}

	visible := state.flatten(t.items)

	state.lastBiggestIndex = max(0, len(visible)-1)

	if len(visible) == 0 {
		return
	}

	availableHeight := area.Height

	var ensureIndexInView *int
	if state.ensureSelectedInView && len(state.selected) != 0 {
		index := slices.IndexFunc(visible, func(f _Flattened[T]) bool {
			return slices.Equal(f.id, state.selected)
		})

		ensureIndexInView = &index
	}

	start := min(state.offset, state.lastBiggestIndex)

	if ensureIndexInView != nil {
		start = min(start, *ensureIndexInView)
	}

	end := start

	var height int

	for _, flattenedItem := range sliceutil.Skip(visible, start) {
		itemHeight := flattenedItem.item.Height()

		if height+itemHeight > availableHeight {
			break
		}

		height += itemHeight
		end++
	}

	if ensureIndexInView != nil {
		for *ensureIndexInView >= end {
			height += visible[end].item.Height()
			end++

			for height > availableHeight {
				height = max(0, height-visible[start].item.Height())
				start++
			}
		}
	}

	state.offset = start
	state.ensureSelectedInView = false

	if t.scrollbar != nil {
		scrollState := scrollwidget.NewState(max(0, len(visible)-height))
		scrollState.SetPosition(start)
		scrollState.SetViewportnContentLen(height)

		scrollArea := bento.Rect{
			X:      fullArea.X,
			Y:      area.Y,
			Width:  fullArea.Width,
			Height: area.Height,
		}

		t.scrollbar.RenderStateful(scrollArea, buffer, scrollState)
	}

	blankSymbol := strings.Repeat(" ", uniseg.StringWidth(t.highlightSymbol))

	var currentHeight int

	hasSelection := len(state.selected) != 0

	for _, flattened := range sliceutil.Take(sliceutil.Skip(visible, state.offset), end-start) {
		id, item := flattened.id, flattened.item

		x := area.X
		y := area.Y
		height := item.Height()

		currentHeight += height

		area := bento.Rect{
			X:      x,
			Y:      y,
			Width:  area.Width,
			Height: height,
		}

		text := item.text
		itemStyle := text.Style

		isSelected := slices.Equal(state.selected, id)

		afterHighlightSymbolX := x
		if hasSelection {
			symbol := blankSymbol
			if isSelected {
				symbol = t.highlightSymbol
			}

			afterHighlightSymbolX, _ = buffer.SetStringN(x, y, symbol, area.Width, itemStyle)
		}

		var afterDepthX int
		{
			indentWidth := flattened.depth() * 2

			afterIndentX, _ := buffer.SetStringN(
				afterHighlightSymbolX,
				y,
				strings.Repeat(" ", indentWidth),
				indentWidth,
				itemStyle,
			)

			symbol := t.nodeClosedSymbol
			if len(item.children) == 0 {
				symbol = t.nodeNoChildrenSymbol
			} else {
				for _, s := range state.opened {
					if slices.Equal(s, id) {
						symbol = t.nodeOpenSymbol
						break
					}
				}
			}

			maxWidth := max(0, area.Width-afterIndentX+x)

			afterDepthX, _ = buffer.SetStringN(afterIndentX, y, symbol, maxWidth, itemStyle)
		}

		textArea := bento.Rect{
			X:      afterDepthX,
			Y:      area.Y,
			Width:  max(0, area.Width-afterDepthX+x),
			Height: area.Height,
		}

		text.Render(textArea, buffer)

		if isSelected {
			buffer.SetStyle(area, t.highlightStyle)
		}

		state.lastRenderedIDs = append(state.lastRenderedIDs, _LastRenderedIDs[T]{
			Y:   area.Y,
			IDs: id,
		})

		var lastIDs [][]T

		for _, f := range visible {
			lastIDs = append(lastIDs, f.id)
		}

		state.lastIDs = lastIDs
	}
}
