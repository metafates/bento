package throbberwidget

import (
	"sync/atomic"
	"time"

	"github.com/metafates/bento"
)

var lastID atomic.Int64

func nextID() int {
	return int(lastID.Add(1))
}

type TickMsg struct {
	time time.Time
	tag  int

	id int
}

// ID of the spinner that this message belongs to. This can be
// helpful when routing messages, however bear in mind that spinners
// will ignore messages that don't contain id by default.
func (t TickMsg) ID() int {
	return t.id
}

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

func (s *State) Update(msg TickMsg) bento.Cmd {
	if msg.id > 0 && msg.id != s.id {
		return nil
	}

	if msg.tag > 0 && msg.tag != s.tag {
		return nil
	}

	s.frame++
	s.tag++

	return s.tick()
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
