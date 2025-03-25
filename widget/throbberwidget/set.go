package throbberwidget

import (
	"fmt"

	"github.com/metafates/bento/widget/textwidget"
)

type Type int

const (
	TypeLine Type = iota
	TypeHorizontalBlock
	TypeVerticalBlock
	TypeBrailleEight
	TypeParenthesis
	TypeCanadian
	TypeEllipsis
	TypeMeter
	TypePulse
)

func (t Type) String() string {
	switch t {
	case TypeBrailleEight:
		return "BrailleEight"
	case TypeCanadian:
		return "Canadian"
	case TypeEllipsis:
		return "Ellipsis"
	case TypeHorizontalBlock:
		return "HorizontalBlock"
	case TypeLine:
		return "Line"
	case TypeMeter:
		return "Meter"
	case TypeParenthesis:
		return "Parenthesis"
	case TypePulse:
		return "Pulse"
	case TypeVerticalBlock:
		return "VerticalBlock"
	default:
		panic(fmt.Sprintf("unexpected throbberwidget.Type: %#v", t))
	}
}

func (t Type) Set() Set {
	switch t {
	case TypeLine:
		return _lineSet
	case TypeHorizontalBlock:
		return _horizontalBlockSet
	case TypeVerticalBlock:
		return _verticalBlockSet
	case TypeBrailleEight:
		return _brailleEightSet
	case TypeParenthesis:
		return _parenthesisSet
	case TypeCanadian:
		return _canadianSet
	case TypeEllipsis:
		return _ellipsisSet
	case TypeMeter:
		return _meterSet
	case TypePulse:
		return _pulseSet
	default:
		return _lineSet
	}
}

type Set []Frame

type Frame string

func (f Frame) String() string {
	return string(f)
}

func (f Frame) Dimensions() (width, height int) {
	lines := textwidget.NewLinesStr(f.String())

	return lines.Width(), lines.Height()
}

func (s Set) Frame(frame int) Frame {
	return s[frame%len(s)]
}

var (
	_lineSet            = Set{`|`, `/`, `-`, `\`}
	_horizontalBlockSet = Set{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
	_verticalBlockSet   = Set{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	_brailleEightSet    = Set{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}
	_parenthesisSet     = Set{"⎛", "⎜", "⎝", "⎞", "⎟", "⎠"}
	_canadianSet        = Set{"ᔐ", "ᯇ", "ᔑ", "ᯇ"}
	_ellipsisSet        = Set{"", ".", "..", "..."}
	_meterSet           = Set{"▱▱▱", "▰▱▱", "▰▰▱", "▰▰▰", "▰▰▱", "▰▱▱", "▱▱▱"}
	_pulseSet           = Set{"█", "▓", "▒", "░"}
)
