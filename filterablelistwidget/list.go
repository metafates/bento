package filterablelistwidget

import (
	"fmt"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/inputwidget"
	"github.com/metafates/bento/listwidget"
	"github.com/metafates/bento/textwidget"
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
	list listwidget.List
}

func New[I Item](list listwidget.List) List[I] {
	return List[I]{
		list: list,
	}
}

func (l List[I]) WithList(list listwidget.List) List[I] {
	l.list = list
	return l
}

func (l List[I]) RenderStateful(area bento.Rect, buffer *bento.Buffer, state *State[I]) {
	if state.filterState != FilterStateNoFilter {
		var inputArea bento.Rect

		bento.NewLayout(
			bento.ConstraintLen(3),
			bento.ConstraintFill(1),
		).Vertical().Split(area).Assign(&inputArea, &area)

		l.renderFilter(inputArea, buffer, state)
	}

	items := make([]textwidget.Text, 0, len(state.filteredIndices))
	for _, i := range state.filteredIndices {
		items = append(items, state.items[i].Text())
	}

	l.list.WithItems(items...).RenderStateful(area, buffer, &state.list)
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

// Render implements bento.Widget.
func (l List[I]) Render(area bento.Rect, buffer *bento.Buffer) {
	l.RenderStateful(area, buffer, new(State[I]))
}
