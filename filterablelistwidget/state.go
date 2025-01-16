package filterablelistwidget

import (
	"math"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/metafates/bento"
	"github.com/metafates/bento/inputwidget"
	"github.com/metafates/bento/listwidget"
)

type FilterState int

const (
	FilterStateNoFilter FilterState = iota
	FilterStateFiltering
	FilterStateFiltered
)

type State[I Item] struct {
	list listwidget.State

	filterState FilterState
	filterInput inputwidget.State

	items           []I
	filteredIndices []int

	cache map[string][]int

	onSelect func()
}

func NewState[I Item](items ...I) State[I] {
	input := inputwidget.NewState()
	input.ShowCursor(true)

	state := State[I]{
		list: listwidget.NewState(),

		filterState: FilterStateNoFilter,

		filterInput: input,

		items:           items,
		filteredIndices: make([]int, 0, len(items)),
		cache:           make(map[string][]int),
		onSelect:        nil,
	}

	state.applyFilter()

	return state
}

func (s *State[I]) OnSelect(f func()) {
	s.onSelect = f
}

func (s *State[I]) reselect() {
	selected, ok := s.list.Selected()
	if !ok {
		return
	}

	s.Select(selected)
}

func (s *State[I]) applyFilter() {
	filter := s.filterInput.String()

	if cached, ok := s.cache[filter]; ok {
		s.filteredIndices = cached
		return
	}

	s.filteredIndices = make([]int, 0, cap(s.filteredIndices))

	defer func() {
		s.reselect()
		s.cache[filter] = s.filteredIndices
	}()

	if filter == "" {
		for i := range s.items {
			s.filteredIndices = append(s.filteredIndices, i)
		}

		return
	}

	for i, item := range s.items {
		var value string

		if f, ok := Item(item).(Filterable); ok {
			value = f.FilterValue()
		} else {
			value = item.Text().String()
		}

		if fuzzy.MatchNormalizedFold(filter, value) {
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

// AllItems returns all items
func (s *State[I]) AllItems() []I {
	return s.items
}

// Items returns currently filtered items
//
// See [State.AllItems] to get all items
func (s *State[I]) Items() []I {
	items := make([]I, 0, len(s.filteredIndices))

	for _, i := range s.filteredIndices {
		items = append(items, s.items[i])
	}

	return items
}

func (s *State[I]) Update(key bento.Key) bool {
	if s.filterState == FilterStateFiltering {
		if s.filterInput.Update(key) {
			s.applyFilter()
			return true
		}
	} else if s.list.Update(key) {
		s.reselect()

		return true
	}

	switch key.String() {
	case "enter":
		if s.filterState != FilterStateFiltering {
			selected, ok := s.Selected()
			if !ok {
				s.SelectFirst()
				return true
			}

			if callable, ok := Item(selected).(Callable); ok {
				if s.onSelect != nil {
					s.onSelect()
				}

				callable.Call()

				return true
			}

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
			if _, ok := s.list.Selected(); ok {
				s.Unselect()
				return true
			}

			return false
		}

		s.setFilteringState(FilterStateNoFilter)

		return true

	case "/":
		s.setFilteringState(FilterStateFiltering)

		return true

	default:
		return false
	}
}

func (s *State[I]) SetOffset(offset int) {
	s.list.SetOffset(offset)
}

func (s *State[I]) Reset() {
	s.Unselect()
	s.filterInput.DeleteLine()
	s.applyFilter()
	s.filterState = FilterStateNoFilter
}

func (s *State[I]) Select(index int) {
	index = s.clampIndex(index)
	s.list.Select(index)
}

func (s *State[I]) clampIndex(index int) int {
	return max(0, min(len(s.filteredIndices)-1, index))
}

func (s *State[I]) SelectNext() {
	s.list.SelectNext()
	s.reselect()
}

func (s *State[I]) SelectPrevious() {
	s.list.SelectPrevious()
	s.reselect()
}

func (s *State[I]) SelectFirst() {
	s.Select(0)
}

func (s *State[I]) SelectLast() {
	s.Select(math.MaxInt)
}

func (s *State[I]) ScrollDownBy(amount int) {
	s.list.ScrollDownBy(amount)
	s.reselect()
}

func (s *State[I]) ScrollUpBy(amount int) {
	s.list.ScrollUpBy(amount)
	s.reselect()
}

func (s *State[I]) Unselect() {
	s.list.Unselect()
}

// Selected returns selected item.
// Returns ok = false if no item is selected.
func (s *State[I]) Selected() (item I, ok bool) {
	index, ok := s.list.Selected()
	if !ok {
		var empty I
		return empty, false
	}

	return s.items[s.filteredIndices[index]], true
}
