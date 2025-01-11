package inputwidget

import (
	"slices"
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/internal/grapheme"
	"github.com/metafates/bento/internal/sliceutil"
	"github.com/rivo/uniseg"
)

type _Graphemes []grapheme.Grapheme

func (g _Graphemes) String() string {
	var b strings.Builder

	for _, gr := range g {
		b.Grow(len(gr.String()))
	}

	for _, gr := range g {
		b.WriteString(gr.String())
	}

	return b.String()
}

type State struct {
	Cursor int

	graphemes _Graphemes
}

func NewState() State {
	return State{}
}

func (s *State) String() string {
	return s.graphemes.String()
}

func (s *State) MoveCursorLeft() {
	s.Cursor = s.clampCursor(s.Cursor - 1)
}

func (s *State) MoveCursorRight() {
	s.Cursor = s.clampCursor(s.Cursor + 1)
}

func (s *State) MoveCursorBegin() {
	s.Cursor = 0
}

func (s *State) MoveCursorEnd() {
	s.Cursor = len(s.graphemes)
}

func (s *State) Append(content string) {
	graphemes := uniseg.NewGraphemes(content)

	for graphemes.Next() {
		s.graphemes = slices.Insert(s.graphemes, s.Cursor, grapheme.New(graphemes.Str()))
		s.MoveCursorRight()
	}
}

func (s *State) DeleteWordUnderCursor() {
	before, under, after := s.splitAtCursor()

	current := under

	for !current.IsEmpty() && current.IsWhitespace() {
		s.MoveCursorLeft()

		if len(before) == 0 {
			current = grapheme.Grapheme{}
			before = nil
			break
		}

		current = before[len(before)-1]
		before = before[:len(before)-1]
	}

	s.graphemes = append(before, current)
	s.graphemes = append(s.graphemes, after...)
}

func (s *State) DeleteUnderCursor() {
	if s.Cursor == 0 {
		return
	}

	before, _, after := s.splitAtCursor()

	s.graphemes = append(before, after...)
	s.MoveCursorLeft()
}

func (s *State) HandleKey(key bento.Key) {
	switch key.Type {
	case bento.KeyLeft:
		s.MoveCursorLeft()

	case bento.KeyRight:
		s.MoveCursorRight()

	case bento.KeyBackspace:
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

func (s *State) splitAtCursor() (before _Graphemes, under grapheme.Grapheme, after _Graphemes) {
	if len(s.graphemes) == 0 || s.Cursor == 0 {
		return nil, grapheme.Grapheme{}, s.graphemes
	}

	before = sliceutil.Take(s.graphemes, s.Cursor-1)

	if len(s.graphemes) > s.Cursor-1 {
		under = s.graphemes[s.Cursor-1]
		after = sliceutil.Skip(s.graphemes, s.Cursor)
	}

	return before, under, after
}
