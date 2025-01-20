package scrollwidget

// State is a struct representing the state of a [Scroll] widget.
//
// It's essential to set the `contentLen` field when using this struct. This field
// represents the total length of the scrollable content. The default value is zero
// which will result in the Scrollbar not rendering.
//
// For example, in the following list, assume there are 4 bullet points:
//
// The `contentLen` is 4, the `position` is 0, the `viewportContentLen` is 2
//
//	┌───────────────┐
//	│1. this is a   █
//	│   single item █
//	│2. this is a   ║
//	│   second item ║
//	└───────────────┘
//
// If you don't have multi-line content, you can leave the `viewportContentLength` set to the
// default and it'll use the track size as a `viewportContentLength`.
type State struct {
	contentLen         int
	position           int
	viewportContentLen int
}

// NewState constructs a new [State] with the specified content length.
// `contentLen` is the total number of element, that can be scrolled.
//
// See [State] for more details.
func NewState(contentLen int) State {
	return State{
		contentLen:         contentLen,
		position:           0,
		viewportContentLen: 0,
	}
}

// Decrements the scroll position by one, ensuring it doesn't go below zero.
func (s *State) Prev() {
	s.position = max(0, s.position-1)
}

// Increments the scroll position by one, ensuring it doesn't exceed the length of the content.
func (s *State) Next() {
	s.position = min(s.position+1, max(0, s.contentLen-1))
}

func (s *State) First() {
	s.position = 0
}

func (s *State) Last() {
	s.position = max(0, s.contentLen-1)
}

func (s *State) Scroll(direction Direction) {
	switch direction {
	case DirectionForward:
		s.Next()
	case DirectionBackward:
		s.Prev()
	}
}

func (s *State) Ratio() float64 {
	return float64(s.position) / float64(s.contentLen)
}

func (s *State) Position() int {
	return s.position
}
