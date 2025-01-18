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

func (s *State) TryUpdate(msg bento.Msg) (bool, bento.Cmd) {
	if s.ShowPopup {
		consumed, cmd := s.BindingList.TryUpdate(msg)
		if consumed {
			return true, cmd
		}

		keyMsg, ok := msg.(bento.KeyMsg)
		if !ok {
			return false, nil
		}

		switch keyMsg.String() {
		case "esc", "q":
			s.closePopup()
			return true, nil
		}
	}

	keyMsg, ok := msg.(bento.KeyMsg)
	if !ok {
		return false, nil
	}

	switch keyMsg.String() {
	case "?":
		s.TogglePopup()
		return true, nil

	default:
		return s.callBinding(bento.Key(keyMsg)), nil
	}
}

func (s *State) callBinding(key bento.Key) bool {
	// TODO: cache
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
