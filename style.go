package bento

import (
	"github.com/metafates/bento/internal/bit"
	"github.com/muesli/termenv"
)

type Modifier uint16

const (
	ModifierNone       Modifier = 0b0000_0000_0000
	ModifierBold       Modifier = 0b0000_0000_0001
	ModifierDim        Modifier = 0b0000_0000_0010
	ModifierItalic     Modifier = 0b0000_0000_0100
	ModifierUnderlined Modifier = 0b0000_0000_1000
	ModifierSlowBlink  Modifier = 0b0000_0001_0000
	ModifierRapidBlink Modifier = 0b0000_0010_0000
	ModifierReversed   Modifier = 0b0000_0100_0000
	ModifierHidden     Modifier = 0b0000_1000_0000
	ModifierCrossedOut Modifier = 0b0001_0000_0000
	ModifierAll        Modifier = ModifierBold |
		ModifierDim |
		ModifierItalic |
		ModifierUnderlined |
		ModifierSlowBlink |
		ModifierRapidBlink |
		ModifierReversed |
		ModifierHidden |
		ModifierCrossedOut
)

func (m Modifier) Contains(other Modifier) bool {
	return bit.Contains(m, other)
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

	addModifier, subModifier Modifier
}

func NewStyle() Style {
	return Style{
		Foreground: StyleColor{},
		Background: StyleColor{},

		addModifier: ModifierNone,
		subModifier: ModifierNone,
	}
}

func (s Style) Italic() Style {
	return s.WithModifier(ModifierItalic)
}

func (s Style) Bold() Style {
	return s.WithModifier(ModifierBold)
}

func (s Style) Underlined() Style {
	return s.WithModifier(ModifierUnderlined)
}

func (s Style) Dim() Style {
	return s.WithModifier(ModifierDim)
}

func (s Style) WithModifier(modifier Modifier) Style {
	s.subModifier = bit.Difference(s.subModifier, modifier)
	s.addModifier = bit.Union(s.addModifier, modifier)

	return s
}

func (s Style) WithoutModifier(modifier Modifier) Style {
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
