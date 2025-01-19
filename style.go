package bento

import (
	"fmt"

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

func (s Style) Reversed() Style {
	return s.WithModifier(ModifierReversed)
}

func (s Style) CrossedOut() Style {
	return s.WithModifier(ModifierCrossedOut)
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

func (s Style) RGB(r, g, b uint8) Style {
	return s.WithForeground(termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", r, g, b)))
}

func (s Style) OnRGB(r, g, b uint8) Style {
	return s.WithBackground(termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", r, g, b)))
}

func (s Style) Black() Style {
	return s.WithForeground(termenv.ANSIBlack)
}

func (s Style) OnBlack() Style {
	return s.WithBackground(termenv.ANSIBlack)
}

func (s Style) Red() Style {
	return s.WithForeground(termenv.ANSIRed)
}

func (s Style) OnRed() Style {
	return s.WithBackground(termenv.ANSIRed)
}

func (s Style) Green() Style {
	return s.WithForeground(termenv.ANSIGreen)
}

func (s Style) OnGreen() Style {
	return s.WithBackground(termenv.ANSIGreen)
}

func (s Style) Yellow() Style {
	return s.WithForeground(termenv.ANSIYellow)
}

func (s Style) OnYellow() Style {
	return s.WithBackground(termenv.ANSIYellow)
}

func (s Style) Blue() Style {
	return s.WithForeground(termenv.ANSIBlue)
}

func (s Style) OnBlue() Style {
	return s.WithBackground(termenv.ANSIBlue)
}

func (s Style) Magenta() Style {
	return s.WithForeground(termenv.ANSIMagenta)
}

func (s Style) OnMagenta() Style {
	return s.WithBackground(termenv.ANSIMagenta)
}

func (s Style) Cyan() Style {
	return s.WithForeground(termenv.ANSICyan)
}

func (s Style) OnCyan() Style {
	return s.WithBackground(termenv.ANSICyan)
}

func (s Style) White() Style {
	return s.WithForeground(termenv.ANSIWhite)
}

func (s Style) OnWhite() Style {
	return s.WithBackground(termenv.ANSIWhite)
}

func (s Style) BrightBlack() Style {
	return s.WithForeground(termenv.ANSIBrightBlack)
}

func (s Style) OnBrightBlack() Style {
	return s.WithBackground(termenv.ANSIBrightBlack)
}

func (s Style) BrightRed() Style {
	return s.WithForeground(termenv.ANSIBrightRed)
}

func (s Style) OnBrightRed() Style {
	return s.WithBackground(termenv.ANSIBrightRed)
}

func (s Style) BrightGreen() Style {
	return s.WithForeground(termenv.ANSIBrightGreen)
}

func (s Style) OnBrightGreen() Style {
	return s.WithBackground(termenv.ANSIBrightGreen)
}

func (s Style) BrightYellow() Style {
	return s.WithForeground(termenv.ANSIBrightYellow)
}

func (s Style) OnBrightYellow() Style {
	return s.WithBackground(termenv.ANSIBrightYellow)
}

func (s Style) BrightBlue() Style {
	return s.WithForeground(termenv.ANSIBrightBlue)
}

func (s Style) OnBrightBlue() Style {
	return s.WithBackground(termenv.ANSIBrightBlue)
}

func (s Style) BrightMagenta() Style {
	return s.WithForeground(termenv.ANSIBrightMagenta)
}

func (s Style) OnBrightMagenta() Style {
	return s.WithBackground(termenv.ANSIBrightMagenta)
}

func (s Style) BrightCyan() Style {
	return s.WithForeground(termenv.ANSIBrightCyan)
}

func (s Style) OnBrightCyan() Style {
	return s.WithBackground(termenv.ANSIBrightCyan)
}

func (s Style) BrightWhite() Style {
	return s.WithForeground(termenv.ANSIBrightWhite)
}

func (s Style) OnBrightWhite() Style {
	return s.WithBackground(termenv.ANSIBrightWhite)
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
