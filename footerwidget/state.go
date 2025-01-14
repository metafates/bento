package footerwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/listwidget"
)

type State struct {
	BindingList listwidget.State

	ShowPopup bool
}

func NewState() State {
	return State{
		BindingList: listwidget.NewState(),
	}
}

func (s *State) Update(key bento.Key) bool {
	listUpdated := s.BindingList.Update(key)
	if listUpdated {
		return true
	}

	switch key.String() {
	case "?":
		s.TogglePopup()
		return true

	case "esc", "q":
		s.ShowPopup = false
		return true

	default:
		return false
	}
}

func (s *State) TogglePopup() {
	s.ShowPopup = !s.ShowPopup
}
