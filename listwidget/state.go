package listwidget

import "math"

type State struct {
	offset   int
	selected *int
}

func NewState() State {
	return State{}
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
// Returns ok = false if no item is selected
func (s *State) Selected() (index int, ok bool) {
	if s.selected == nil {
		return 0, false
	}

	return *s.selected, true
}
