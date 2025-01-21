package treewidget

import (
	"slices"

	"github.com/metafates/bento"
)

type _LastRenderedIDs[T comparable] struct {
	Y   int
	IDs []T
}

type State[T comparable] struct {
	offset               int
	opened               [][]T
	selected             []T
	ensureSelectedInView bool

	lastArea         bento.Rect
	lastBiggestIndex int
	lastIDs          [][]T
	lastRenderedIDs  []_LastRenderedIDs[T]
}

func NewState[T comparable]() State[T] {
	return State[T]{}
}

func (s *State[T]) flatten(items []Item[T]) []_Flattened[T] {
	return flatten(s.opened, items, nil)
}

func (s *State[T]) Select(id []T) bool {
	s.ensureSelectedInView = true

	changed := !slices.Equal(s.selected, id)

	s.selected = id

	return changed
}

func (s *State[T]) Open(id []T) bool {
	if len(id) == 0 {
		return false
	}

	for _, o := range s.opened {
		if slices.Equal(o, id) {
			return false
		}
	}

	s.opened = append(s.opened, id)

	return true
}

type _Flattened[T comparable] struct {
	id   []T
	item Item[T]
}

func (f _Flattened[T]) depth() int {
	return len(f.id) - 1
}

func flatten[T comparable](
	openIDs [][]T,
	items []Item[T],
	current []T,
) []_Flattened[T] {
	var result []_Flattened[T]

	for _, item := range items {
		childID := current

		childID = append(childID, item.id)

		result = append(result, _Flattened[T]{
			id:   childID,
			item: item,
		})

		for _, s := range openIDs {
			if slices.Equal(s, childID) {
				childResult := flatten(openIDs, item.children, childID)

				result = append(result, childResult...)
				break
			}
		}
	}

	return result
}
