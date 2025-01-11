package inputwidget

import (
	"slices"

	"github.com/metafates/bento"
	"github.com/metafates/bento/internal/grapheme"
	"github.com/metafates/bento/internal/sliceutil"
	"github.com/rivo/uniseg"
)

type State struct {
	cursor int

	graphemes grapheme.Graphemes

	offset int
}

func NewState() State {
	return State{}
}

func (s *State) String() string {
	return s.graphemes.String()
}

func (s *State) setCursor(cursor int) {
	cursor = s.clampCursor(cursor)
	s.cursor = cursor
}

func (s *State) MoveCursorLeft() {
	s.setCursor(s.cursor - 1)
}

func (s *State) MoveCursorRight() {
	s.setCursor(s.cursor + 1)
}

func (s *State) MoveCursorBegin() {
	s.setCursor(0)
	s.offset = 0
}

func (s *State) MoveCursorEnd() {
	s.setCursor(len(s.graphemes))
}

func (s *State) Append(content string) {
	graphemes := uniseg.NewGraphemes(content)

	for graphemes.Next() {
		s.graphemes = slices.Insert(s.graphemes, s.cursor, grapheme.New(graphemes.Str()))
		s.MoveCursorRight()
	}
}

func (s *State) DeleteLine() {
	s.graphemes = nil
	s.cursor = 0
}

func (s *State) DeleteWordUnderCursor() {
	current := s.underCursor()

	// exit early
	if current.IsEmpty() {
		return
	}

	if current.IsWhitespace() {
		s.deleteWhile(func(g grapheme.Grapheme) bool { return g.IsWhitespace() })
	}

	s.deleteWhile(func(g grapheme.Grapheme) bool { return !g.IsEmpty() && !g.IsWhitespace() })
}

func (s *State) deleteWhile(cond func(g grapheme.Grapheme) bool) {
	before, under, after := s.splitAtCursor()

	current := under

	for cond(current) {
		if len(before) == 0 {
			current = grapheme.Grapheme{}
			break
		}
		s.MoveCursorLeft()

		current = before[len(before)-1]
		before = before[:len(before)-1]
	}

	s.graphemes = append(before, current)
	s.graphemes = append(s.graphemes, after...)
}

func (s *State) DeleteUnderCursor() {
	before, _, after := s.splitAtCursor()

	s.graphemes = append(before, after...)
	s.MoveCursorLeft()
}

func (s *State) MoveWordRight() {
	// TODO
}

func (s *State) MoveWordLeft() {
	// TODO
}

func (s *State) HandleKey(key bento.Key) {
	switch key.Type {
	case bento.KeyLeft:
		s.MoveCursorLeft()

	case bento.KeyShiftLeft:
		s.MoveWordLeft()

	case bento.KeyRight:
		s.MoveCursorRight()

	case bento.KeyShiftRight:
		s.MoveWordRight()

	case bento.KeyBackspace, bento.KeyDelete:
		s.DeleteUnderCursor()

	case bento.KeyCtrlA:
		s.MoveCursorBegin()

	case bento.KeyCtrlE:
		s.MoveCursorEnd()

	case bento.KeyCtrlW:
		s.DeleteWordUnderCursor()

	case bento.KeyRunes, bento.KeySpace:
		s.Append(string(key.Runes))
	}
}

func (s *State) clampCursor(cursor int) int {
	cursor = max(0, cursor)
	cursor = min(cursor, len(s.graphemes))

	return cursor
}

func (s *State) underCursor() grapheme.Grapheme {
	if len(s.graphemes) == 0 || s.cursor == 0 {
		return grapheme.Empty()
	}

	return s.graphemes[s.cursor-1]
}

func (s *State) isEmpty() bool {
	if s.graphemes == nil {
		return true
	}

	for _, g := range s.graphemes {
		if !g.IsEmpty() {
			return false
		}
	}

	return true
}

func (s *State) splitAtCursor() (before grapheme.Graphemes, under grapheme.Grapheme, after grapheme.Graphemes) {
	if len(s.graphemes) == 0 || s.cursor == 0 {
		return nil, grapheme.Empty(), s.graphemes
	}

	before = s.graphemes[:s.cursor-1]
	under = s.graphemes[s.cursor-1]

	if len(s.graphemes) > s.cursor {
		after = s.graphemes[s.cursor:]
	}

	return before, under, after
}

func (s *State) getBounds(maxWidth int) (int, int) {
	if len(s.graphemes) == 0 {
		return 0, 0
	}

	offset := min(s.offset, max(0, len(s.graphemes)-1))

	firstVisibleIndex := offset
	lastVisibleIndex := offset

	var widthFromOffset int

	for _, g := range sliceutil.Skip(s.graphemes, offset) {
		if widthFromOffset+g.Width() > maxWidth {
			break
		}

		widthFromOffset += g.Width()

		lastVisibleIndex++
	}

	indexToDisplay := offset
	if s.cursor != 0 {
		indexToDisplay = min(s.cursor-1, lastVisibleIndex)
	}

	for indexToDisplay >= lastVisibleIndex {
		widthFromOffset += s.graphemes[lastVisibleIndex].Width()

		lastVisibleIndex++

		for widthFromOffset > maxWidth {
			widthFromOffset = max(0, widthFromOffset-s.graphemes[firstVisibleIndex].Width())

			firstVisibleIndex++
		}
	}

	for indexToDisplay < firstVisibleIndex {
		firstVisibleIndex--

		widthFromOffset += s.graphemes[firstVisibleIndex].Width()

		for widthFromOffset > maxWidth {
			lastVisibleIndex--

			widthFromOffset = max(0, widthFromOffset-s.graphemes[lastVisibleIndex].Width())
		}
	}

	if s.cursor-firstVisibleIndex >= lastVisibleIndex {
		diff := s.cursor - lastVisibleIndex

		firstVisibleIndex += diff
		lastVisibleIndex += diff
	}

	return firstVisibleIndex, lastVisibleIndex
}
