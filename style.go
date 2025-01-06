package bento

import (
	"github.com/muesli/termenv"
)

type StyleModifier struct {
	enabled *bool
}

func (s *StyleModifier) IsSet() bool {
	return s.enabled != nil
}

func (s *StyleModifier) Bool() bool {
	if s.enabled == nil {
		return false
	}

	return *s.enabled
}

func (s *StyleModifier) Set(enabled bool) {
	s.enabled = &enabled
}

func (s *StyleModifier) Reset() {
	s.enabled = nil
}

type StyleColor struct {
	color termenv.Color
}

func (s *StyleColor) Color() termenv.Color {
	if s.color == nil {
		return termenv.NoColor{}
	}

	return s.color
}

func (s *StyleColor) IsSet() bool {
	return s.color == nil
}

func (s *StyleColor) Reset() {
	s.color = nil
}

func (s *StyleColor) Set(color termenv.Color) {
	s.color = color
}

type Style struct {
	Foreground, Background StyleColor

	Bold,
	Dim,
	Italic,
	Underlined,
	Reversed,
	CrossedOut,
	SlowBlink,
	RapidBlink,
	Hidden StyleModifier
}

func NewStyle() Style {
	// default values are ok
	return Style{}
}

func (s Style) WithBackground(color termenv.Color) Style {
	s.Background.Set(color)

	return s
}

func (s Style) WithForeground(color termenv.Color) Style {
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

	for prev, curr := range map[StyleModifier]StyleModifier{
		s.Bold:       other.Bold,
		s.Dim:        other.Dim,
		s.Italic:     other.Italic,
		s.Underlined: other.Underlined,
		s.Reversed:   other.Reversed,
		s.CrossedOut: other.CrossedOut,
		s.SlowBlink:  other.SlowBlink,
		s.RapidBlink: other.RapidBlink,
		s.Hidden:     other.Hidden,
	} {
		if prev.IsSet() && prev.Bool() == curr.Bool() {
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

	if patch.Bold.IsSet() {
		s.Bold = patch.Bold
	}

	if patch.Dim.IsSet() {
		s.Dim = patch.Dim
	}

	if patch.Italic.IsSet() {
		s.Italic = patch.Italic
	}

	if patch.Underlined.IsSet() {
		s.Underlined = patch.Underlined
	}

	if patch.Reversed.IsSet() {
		s.Reversed = patch.Reversed
	}

	if patch.CrossedOut.IsSet() {
		s.CrossedOut = patch.CrossedOut
	}

	if patch.SlowBlink.IsSet() {
		s.SlowBlink = patch.SlowBlink
	}

	if patch.RapidBlink.IsSet() {
		s.RapidBlink = patch.RapidBlink
	}

	if patch.Hidden.IsSet() {
		s.Hidden = patch.Hidden
	}

	return s
}
