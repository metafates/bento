package helpwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/filterablelistwidget"
)

type State struct {
	bindingList filterablelistwidget.State[Binding]
	showPopup   bool
}

func NewState(bindings ...Binding) *State {
	state := State{
		bindingList: filterablelistwidget.NewState(bindings...),
	}

	state.bindingList.OnSelect(state.closePopup)

	return &state
}

func (s *State) closePopup() {
	s.showPopup = false
	s.bindingList.Reset()
}

func (s *State) TryUpdate(msg bento.Msg) (bool, bento.Cmd) {
	if s.showPopup {
		consumed, cmd := s.bindingList.TryUpdate(msg)
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
		return s.callBinding(bento.Key(keyMsg))
	}
}

func (s *State) callBinding(key bento.Key) (bool, bento.Cmd) {
	// TODO: cache
	for _, b := range s.bindingList.AllItems() {
		if b.Matches(key) {
			return true, b.Call()
		}
	}

	return false, nil
}

func (s *State) TogglePopup() {
	s.showPopup = !s.showPopup
}
