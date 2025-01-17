package helpwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/filterablelistwidget"
)

type State struct {
	BindingList filterablelistwidget.State[Binding]

	ShowPopup bool
}

func NewState(bindings ...Binding) *State {
	state := State{
		BindingList: filterablelistwidget.NewState(bindings...),
	}

	state.BindingList.OnSelect(state.closePopup)

	return &state
}

func (s *State) closePopup() {
	s.ShowPopup = false
	s.BindingList.Reset()
}

func (s *State) Update(key bento.Key) bool {
	if s.ShowPopup {
		listUpdated := s.BindingList.Update(key)
		if listUpdated {
			return true
		}

		switch key.String() {
		case "esc", "q":
			s.closePopup()
			return true
		}
	}

	switch key.String() {
	case "?":
		s.TogglePopup()
		return true

	default:
		return s.callBinding(key)
	}
}

func (s *State) callBinding(key bento.Key) bool {
	for _, b := range s.BindingList.AllItems() {
		if b.Matches(key) {
			b.Call()
			return true
		}
	}

	return false
}

func (s *State) TogglePopup() {
	s.ShowPopup = !s.ShowPopup
}
