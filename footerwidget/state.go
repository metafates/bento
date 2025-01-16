package footerwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/filterablelistwidget"
)

type State struct {
	BindingList filterablelistwidget.State[Binding]

	ShowPopup bool
}

func NewState(bindings ...Binding) State {
	return State{
		BindingList: filterablelistwidget.NewState(bindings...),
	}
}

func (s *State) Update(key bento.Key) bool {
	if s.ShowPopup {
		listUpdated := s.BindingList.Update(key)
		if listUpdated {
			return true
		}
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
