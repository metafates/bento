package blockwidget

import (
	"github.com/metafates/bento/internal/bit"
	"github.com/metafates/bento/symbol"
)

type BorderType int

const (
	BorderTypeSharp BorderType = iota
	BorderTypeRounded
	BorderTypeDouble
	BorderTypeThick
	// BorderTypeQuadrantInside
	// BorderTypeQuadrantOutside
)

func (b BorderType) Set() BorderSet {
	switch b {
	case BorderTypeSharp:
		return _sharpBorderSet
	case BorderTypeRounded:
		return _roundedBorderSet
	case BorderTypeDouble:
		return _doubleBorderSet
	case BorderTypeThick:
		return _thickBorderSet
	default:
		return _sharpBorderSet
	}
}

type Side uint8

const (
	SideNone   Side = 0b0000
	SideTop    Side = 0b0001
	SideRight  Side = 0b0010
	SideBottom Side = 0b0100
	SideLeft   Side = 0b1000
	SideAll         = SideTop | SideRight | SideBottom | SideLeft
)

func (b Side) intersects(other Side) bool {
	return bit.Intersects(b, other)
}

func (b Side) contains(other Side) bool {
	return bit.Contains(b, other)
}

var (
	_sharpBorderSet = BorderSet{
		TopLeft:          symbol.LineTopLeft,
		TopRight:         symbol.LineTopRight,
		BottomLeft:       symbol.LineBottomLeft,
		BottomRight:      symbol.LineBottomRight,
		VerticalLeft:     symbol.LineVertical,
		VerticalRight:    symbol.LineVertical,
		HorizontalTop:    symbol.LineHorizontal,
		HorizontalBottom: symbol.LineHorizontal,
	}

	_roundedBorderSet = _sharpBorderSet.
				WithTopLeft(symbol.LineRoundedTopLeft).
				WithTopRight(symbol.LineRoundedTopRight).
				WithBottomLeft(symbol.LineRoundedBottomLeft).
				WithBottomRight(symbol.LineRoundedBottomRight)

	_thickBorderSet = BorderSet{
		TopLeft:          symbol.LineThickTopLeft,
		TopRight:         symbol.LineThickTopRight,
		BottomLeft:       symbol.LineThickBottomLeft,
		BottomRight:      symbol.LineThickBottomRight,
		VerticalLeft:     symbol.LineThickVertical,
		VerticalRight:    symbol.LineThickVertical,
		HorizontalTop:    symbol.LineThickHorizontal,
		HorizontalBottom: symbol.LineThickHorizontal,
	}

	_doubleBorderSet = BorderSet{
		TopLeft:          symbol.LineDoubleTopLeft,
		TopRight:         symbol.LineDoubleTopRight,
		BottomLeft:       symbol.LineDoubleBottomLeft,
		BottomRight:      symbol.LineDoubleBottomRight,
		VerticalLeft:     symbol.LineDoubleVertical,
		VerticalRight:    symbol.LineDoubleVertical,
		HorizontalTop:    symbol.LineDoubleHorizontal,
		HorizontalBottom: symbol.LineDoubleHorizontal,
	}
)

type BorderSet struct {
	TopLeft,
	TopRight,
	BottomLeft,
	BottomRight,
	VerticalLeft,
	VerticalRight,
	HorizontalTop,
	HorizontalBottom string
}

func (b BorderSet) WithTopLeft(s string) BorderSet {
	b.TopLeft = s
	return b
}

func (b BorderSet) WithTopRight(s string) BorderSet {
	b.TopRight = s
	return b
}

func (b BorderSet) WithBottomLeft(s string) BorderSet {
	b.BottomLeft = s
	return b
}

func (b BorderSet) WithBottomRight(s string) BorderSet {
	b.BottomRight = s
	return b
}

func (b BorderSet) WithVerticalLeft(s string) BorderSet {
	b.VerticalLeft = s
	return b
}

func (b BorderSet) WithVerticalRight(s string) BorderSet {
	b.VerticalRight = s
	return b
}

func (b BorderSet) WithHorizontalTop(s string) BorderSet {
	b.HorizontalTop = s
	return b
}

func (b BorderSet) WithHorizontalBottom(s string) BorderSet {
	b.HorizontalBottom = s
	return b
}
