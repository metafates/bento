package filterablelistwidget

import (
	"math"
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/inputwidget"
)

type FilterState int

const (
	FilterStateNoFilter FilterState = iota
	FilterStateFiltering
	FilterStateFiltered
)

type State[I Item] struct {
	offset   int
	selected *int

	filterState FilterState
	filterInput inputwidget.State

	items           []I
	filteredIndices []int
}

func NewState[I Item](items ...I) State[I] {
	input := inputwidget.NewState()
	input.ShowCursor(true)

	state := State[I]{
		offset:   0,
		selected: nil,

		filterState: FilterStateNoFilter,
		filterInput: input,

		items:           items,
		filteredIndices: make([]int, 0, len(items)),
	}

	state.applyFilter()

	return state
}

func (s *State[I]) reselect() {
	if s.selected == nil {
		return
	}

	s.Select(*s.selected)
}

func (s *State[I]) applyFilter() {
	// TODO: optimize
	s.filteredIndices = make([]int, 0, cap(s.filteredIndices))

	filter := s.filterInput.String()

	defer s.reselect()

	if filter == "" {
		for i := range s.items {
			s.filteredIndices = append(s.filteredIndices, i)
		}

		return
	}

	for i, item := range s.items {
		var value string

		if f, ok := Item(item).(FilterableItem); ok {
			value = f.FilterValue()
		} else {
			value = item.Title().String()
		}

		if strings.Contains(value, filter) {
			s.filteredIndices = append(s.filteredIndices, i)
		}
	}
}

func (s *State[I]) Len() int {
	return len(s.items)
}

func (s *State[I]) LenFiltered() int {
	return len(s.filteredIndices)
}

func (s *State[I]) setFilteringState(state FilterState) {
	s.filterState = state
	s.filterInput.ShowCursor(s.filterState == FilterStateFiltering)

	if state == FilterStateNoFilter {
		s.filterInput.DeleteLine()
		s.Unselect()
		s.applyFilter()
	}
}

func (s *State[I]) SetItems(items ...I) {
	s.items = items
	s.applyFilter()
}

func (s *State[I]) Items() []I {
	return s.items
}

func (s *State[I]) Update(key bento.Key) bool {
	if s.filterState == FilterStateFiltering {
		if s.filterInput.Update(key) {
			s.applyFilter()
			return true
		}
	}

	switch key.String() {
	case "enter":
		if s.filterState != FilterStateFiltering {
			return false
		}

		if s.filterInput.String() == "" {
			s.setFilteringState(FilterStateNoFilter)
		} else {
			s.setFilteringState(FilterStateFiltered)
		}

		s.applyFilter()

		return true
	case "esc":
		if s.filterState == FilterStateNoFilter {
			if s.selected != nil {
				s.Unselect()
				return true
			}

			return false
		}

		s.setFilteringState(FilterStateNoFilter)

		return true

	case "ctrl+u":
		if s.filterState == FilterStateFiltering {
			return false
		}

		s.ScrollUpBy(8)
		return true

	case "ctrl+d":
		if s.filterState == FilterStateFiltering {
			return false
		}

		s.ScrollDownBy(8)
		return true

	case "G":
		if s.filterState == FilterStateFiltering {
			return false
		}

		s.SelectLast()
		return true

	case "g":
		if s.filterState == FilterStateFiltering {
			return false
		}

		s.SelectFirst()
		return true

	case "j", "down":
		if s.filterState == FilterStateFiltering {
			return false
		}

		s.SelectNext()
		return true

	case "k", "up":
		if s.filterState == FilterStateFiltering {
			return false
		}

		s.SelectPrevious()
		return true

	case "/":
		s.setFilteringState(FilterStateFiltering)

		return true

	default:
		return false
	}
}

func (s *State[I]) SetOffset(offset int) {
	s.offset = offset
}

func (s *State[I]) Select(index int) {
	index = s.clampIndex(index)
	s.selected = &index
}

func (s *State[I]) clampIndex(index int) int {
	return max(0, min(len(s.filteredIndices)-1, index))
}

func (s *State[I]) SelectNext() {
	var next int

	if s.selected != nil {
		next = *s.selected + 1
	}

	s.Select(next)
}

func (s *State[I]) SelectPrevious() {
	previous := math.MaxInt

	if s.selected != nil {
		previous = max(0, *s.selected-1)
	}

	s.Select(previous)
}

func (s *State[I]) SelectFirst() {
	s.Select(0)
}

func (s *State[I]) SelectLast() {
	s.Select(math.MaxInt)
}

func (s *State[I]) ScrollDownBy(amount int) {
	var selected int

	if s.selected != nil {
		selected = *s.selected
	}

	s.Select(selected + amount)
}

func (s *State[I]) ScrollUpBy(amount int) {
	var selected int

	if s.selected != nil {
		selected = *s.selected
	}

	s.Select(selected - amount)
}

func (s *State[I]) Unselect() {
	s.selected = nil
	s.offset = 0
}

// Selected returns index of the selected item.
// Returns ok = false if no item is selected.
//
// NOTE: Returned index may be greater, than you would expect, since it is trimmed only after [List.RenderStateful] call
//
// Use [State.SelectedWithLimit] if you know the limit or [GetSelectedItem] to select an item from the slice
func (s *State[I]) Selected() (item I, ok bool) {
	if s.selected == nil {
		var empty I
		return empty, false
	}

	index := *s.selected
	return s.items[s.filteredIndices[index]], true
}
