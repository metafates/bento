package listwidget

import (
	"math"

	"github.com/metafates/bento"
	"github.com/metafates/bento/inputwidget"
)

type State struct {
	offset   int
	selected *int
}

func NewState() State {
	input := inputwidget.NewState()
	input.ShowCursor(true)

	return State{
		offset:   0,
		selected: nil,
	}
}

func (s *State) TryUpdate(msg bento.Msg) (bool, bento.Cmd) {
	keyMsg, ok := msg.(bento.KeyMsg)
	if !ok {
		return false, nil
	}

	return s.update(bento.Key(keyMsg)), nil
}

func (s *State) update(key bento.Key) bool {
	switch key.String() {
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

	s.Select(next)
}

func (s *State) SelectPrevious() {
	previous := math.MaxInt

	if s.selected != nil {
		previous = max(0, *s.selected-1)
	}

	s.Select(previous)
}

func (s *State) SelectFirst() {
	s.Select(0)
}

func (s *State) SelectLast() {
	s.Select(math.MaxInt)
}

func (s *State) ScrollDownBy(amount int) {
	var selected int

	if s.selected != nil {
		selected = *s.selected
	}

	s.Select(selected + amount)
}

func (s *State) ScrollUpBy(amount int) {
	var selected int

	if s.selected != nil {
		selected = *s.selected
	}

	s.Select(selected - amount)
}

func (s *State) Unselect() {
	s.selected = nil
	s.offset = 0
}

// Selected returns index of the selected item.
// Returns ok = false if no item is selected.
//
// NOTE: Returned index may be greater, than you would expect, since the state does not the items count.
// Therefore, you must supply a max value as an argument
func (s *State) Selected(max_ int) (index int, ok bool) {
	if s.selected == nil {
		return 0, false
	}

	return min(max_, max(0, *s.selected)), true
}
