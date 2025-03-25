package throbberwidget

import (
	"time"

	"github.com/metafates/bento"
)

var _ bento.TryUpdater = (*State)(nil)

type State struct {
	frame int
	id    int
	tag   int
	FPS   time.Duration
}

func NewState() State {
	return State{
		frame: 0,
		id:    nextID(),
		tag:   0,
		FPS:   time.Second / 9,
	}
}

// ID of the spinner. This can be
// helpful when routing messages, however bear in mind that spinners
// will ignore messages that don't contain id by default.
func (s *State) ID() int {
	return s.id
}

func (s *State) TryUpdate(msg bento.Msg) (bento.Cmd, bool) {
	tickMsg, ok := msg.(TickMsg)
	if !ok {
		return nil, false
	}

	if tickMsg.id > 0 && tickMsg.id != s.id {
		return nil, false
	}

	if tickMsg.tag > 0 && tickMsg.tag != s.tag {
		return nil, false
	}

	s.frame++
	s.tag++

	return s.tick(), true
}

// Tick is the command used to advance the spinner one frame. Use this command
// to effectively start the throbber.
func (s *State) Tick() bento.Msg {
	return TickMsg{
		id:   s.id,
		time: time.Now(),
		tag:  s.tag,
	}
}

func (s *State) WithFPS(fps time.Duration) *State {
	s.FPS = fps
	return s
}

func (s *State) tick() bento.Cmd {
	id := s.id
	tag := s.tag

	return bento.Tick(s.FPS, func(t time.Time) bento.Msg {
		return TickMsg{
			id:   id,
			time: t,
			tag:  tag,
		}
	})
}
