package listwidget

import (
	"math"

	"github.com/metafates/bento"
	"github.com/metafates/bento/inputwidget"
)

type FilterState int

const (
	FilterStateNoFilter FilterState = iota
	FilterStateFiltering
	FilterStateFiltered
)

type State struct {
	offset   int
	selected *int

	FilterState FilterState
	input       inputwidget.State
}

func NewState() State {
	input := inputwidget.NewState()
	input.ShowCursor(true)

	return State{
		offset:      0,
		selected:    nil,
		FilterState: FilterStateNoFilter,
		input:       input,
	}
}

func (s *State) setFilteringState(state FilterState) {
	s.FilterState = state
	s.input.ShowCursor(s.FilterState == FilterStateFiltering)

	if state == FilterStateNoFilter {
		s.input.DeleteLine()
		s.Unselect()
	}
}

func (s *State) Update(key bento.Key) bool {
	if s.FilterState == FilterStateFiltering {
		if s.input.Update(key) {
			return true
		}
	}

	switch key.String() {
	case "enter":
		if s.FilterState == FilterStateFiltering {
			s.setFilteringState(FilterStateFiltered)
			return true
		}

		return false
	case "esc":
		if s.FilterState == FilterStateNoFilter {
			return false
		}

		s.setFilteringState(FilterStateNoFilter)
		return true

	case "ctrl+u":
		s.ScrollUpBy(8)
		return true

	case "ctrl+d":
		s.ScrollDownBy(8)
		return true

	case "G":
		s.SelectLast()
		return true

	case "g":
		s.SelectFirst()
		return true

	case "j", "down":
		s.SelectNext()
		return true

	case "k", "up":
		s.SelectPrevious()
		return true

	case "/":
		s.setFilteringState(FilterStateFiltering)

		return true

	default:
		return false
	}
}

func (s *State) SetOffset(offset int) {
	s.offset = offset
}

func (s *State) Select(index int) {
	s.selected = &index
}

func (s *State) SelectNext() {
	var next int

	if s.selected != nil {
		next = *s.selected + 1
	}

	s.selected = &next
}

func (s *State) SelectPrevious() {
	previous := math.MaxInt

	if s.selected != nil {
		previous = max(0, *s.selected-1)
	}

	s.selected = &previous
}

func (s *State) SelectFirst() {
	first := 0

	s.selected = &first
}

func (s *State) SelectLast() {
	last := math.MaxInt

	s.selected = &last
}

func (s *State) ScrollDownBy(amount int) {
	var selected int

	if s.selected != nil {
		selected = *s.selected
	}

	s.Select(max(0, selected+amount))
}

func (s *State) ScrollUpBy(amount int) {
	var selected int

	if s.selected != nil {
		selected = *s.selected
	}

	s.Select(max(0, selected-amount))
}

func (s *State) Unselect() {
	s.selected = nil
	s.offset = 0
}

// Selected returns index of the selected item.
// Returns ok = false if no item is selected.
//
// NOTE: Returned index may be greater, than you would expect, since it is trimmed only after [List.RenderStateful] call
//
// Use [State.SelectedWithLimit] if you know the limit or [GetSelectedItem] to select an item from the slice
func (s *State) Selected() (index int, ok bool) {
	if s.selected == nil {
		return 0, false
	}

	return *s.selected, true
}

// SelectedWithLimit returns index of the selected item wrapped at the limit
//
// See: [State.Selected]
func (s *State) SelectedWithLimit(limit int) (index int, ok bool) {
	index, ok = s.Selected()
	if !ok {
		return 0, false
	}

	return min(limit, index), true
}

func GetSelectedItem[S ~[]E, E any](state State, slice S) (element E, ok bool) {
	index, ok := state.SelectedWithLimit(len(slice) - 1)
	if !ok {
		var e E
		return e, false
	}

	return slice[index], true
}
