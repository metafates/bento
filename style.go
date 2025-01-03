package bento

import (
	"github.com/muesli/termenv"
)

type StyleModifier struct {
	enabled *bool
}

func (s *StyleModifier) isSet() bool {
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

type StyleColor struct {
	color termenv.Color
}

func (s *StyleColor) Color() termenv.Color {
	if s.color == nil {
		return termenv.NoColor{}
	}

	return s.color
}

func (s *StyleColor) isSet() bool {
	return s.color == nil
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

// TODO: modifiers

func (s Style) Patched(patch Style) Style {
	if patch.Foreground.isSet() {
		s.Foreground = patch.Foreground
	}

	if patch.Background.isSet() {
		s.Background = patch.Background
	}

	if patch.Bold.isSet() {
		s.Bold = patch.Bold
	}

	if patch.Dim.isSet() {
		s.Dim = patch.Dim
	}

	if patch.Italic.isSet() {
		s.Italic = patch.Italic
	}

	if patch.Underlined.isSet() {
		s.Underlined = patch.Underlined
	}

	if patch.Reversed.isSet() {
		s.Reversed = patch.Reversed
	}

	if patch.CrossedOut.isSet() {
		s.CrossedOut = patch.CrossedOut
	}

	if patch.Hidden.isSet() {
		s.Hidden = patch.Hidden
	}

	return s
}
