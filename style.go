package bento

import (
	"github.com/metafates/bento/internal/bit"
	"github.com/muesli/termenv"
)

type StyleModifier uint16

const (
	StyleModifierNone       StyleModifier = 0b0000_0000_0000
	StyleModifierBold       StyleModifier = 0b0000_0000_0001
	StyleModifierDim        StyleModifier = 0b0000_0000_0010
	StyleModifierItalic     StyleModifier = 0b0000_0000_0100
	StyleModifierUnderlined StyleModifier = 0b0000_0000_1000
	StyleModifierSlowBlink  StyleModifier = 0b0000_0001_0000
	StyleModifierRapidBlink StyleModifier = 0b0000_0010_0000
	StyleModifierReversed   StyleModifier = 0b0000_0100_0000
	StyleModifierHidden     StyleModifier = 0b0000_1000_0000
	StyleModifierCrossedOut StyleModifier = 0b0001_0000_0000
	StyleModifierAll        StyleModifier = StyleModifierBold |
		StyleModifierDim |
		StyleModifierItalic |
		StyleModifierUnderlined |
		StyleModifierSlowBlink |
		StyleModifierRapidBlink |
		StyleModifierReversed |
		StyleModifierHidden |
		StyleModifierCrossedOut
)

func (m *StyleModifier) Insert(other StyleModifier) {
	*m = bit.Union(*m, other)
}

func (m *StyleModifier) Remove(other StyleModifier) {
	*m = bit.Difference(*m, other)
}

func (m *StyleModifier) Contains(other StyleModifier) bool {
	return bit.Contains(*m, other)
}

type StyleColor struct {
	color Color
}

func (s *StyleColor) Color() Color {
	if s.color == nil {
		return termenv.NoColor{}
	}

	return s.color
}

func (s *StyleColor) IsSet() bool {
	return s.color != nil
}

func (s *StyleColor) Reset() {
	s.color = nil
}

func (s *StyleColor) Set(color Color) {
	s.color = color
}

type Style struct {
	Foreground, Background StyleColor

	addModifier, subModifier StyleModifier
}

func NewStyle() Style {
	return Style{
		Foreground: StyleColor{},
		Background: StyleColor{},

		addModifier: StyleModifierNone,
		subModifier: StyleModifierNone,
	}
}

func (s Style) WithModifier(modifier StyleModifier) Style {
	s.subModifier = bit.Difference(s.subModifier, modifier)
	s.addModifier = bit.Union(s.addModifier, modifier)

	return s
}

func (s Style) WithoutModifier(modifier StyleModifier) Style {
	s.addModifier = bit.Difference(s.addModifier, modifier)
	s.subModifier = bit.Union(s.subModifier, modifier)

	return s
}

func (s Style) WithBackground(color Color) Style {
	s.Background.Set(color)

	return s
}

func (s Style) WithForeground(color Color) Style {
	s.Foreground.Set(color)

	return s
}

func (s Style) Sub(other Style) Style {
	for prev, curr := range map[StyleColor]StyleColor{
		s.Background: other.Background,
		s.Foreground: other.Foreground,
	} {
		if prev.IsSet() && prev.Color() == curr.Color() {
			prev.Reset()
		}
	}

	return s
}

func (s Style) Patched(patch Style) Style {
	if patch.Foreground.IsSet() {
		s.Foreground = patch.Foreground
	}

	if patch.Background.IsSet() {
		s.Background = patch.Background
	}

	s.addModifier = bit.Difference(s.addModifier, patch.subModifier)
	s.addModifier = bit.Union(s.addModifier, patch.addModifier)

	s.subModifier = bit.Difference(s.subModifier, patch.addModifier)
	s.subModifier = bit.Union(s.subModifier, patch.subModifier)

	return s
}
