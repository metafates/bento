package throbberwidget

import (
	"sync/atomic"
	"time"
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
