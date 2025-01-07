package blockwidget

import (
	"github.com/metafates/bento/internal/bit"
	"github.com/metafates/bento/symbol"
)

type BorderType int

const (
	BorderTypePlain BorderType = iota
	BorderTypeRounded
	// BorderDouble
	// BorderThick
	// BorderQuadrantInside
	// BorderQuadrantOutside
)

func (b BorderType) Set() BorderSet {
	switch b {
	case BorderTypePlain:
		return _plainBorderSet
	case BorderTypeRounded:
		return _roundedBorderSet
	default:
		return _plainBorderSet
	}
}

type Borders uint8

const (
	BordersNone   Borders = 0b0000
	BordersTop    Borders = 0b0001
	BordersRight  Borders = 0b0010
	BordersBottom Borders = 0b0100
	BordersLeft   Borders = 0b1000
	BordersAll            = BordersTop | BordersRight | BordersBottom | BordersLeft
)

func (b Borders) intersects(other Borders) bool {
	return bit.Intersects(b, other)
}

func (b Borders) contains(other Borders) bool {
	return bit.Contains(b, other)
}

type BorderSet struct {
	TopLeft, TopRight,
	BottomLeft, BottomRight,
	VerticalLeft, VerticalRight,
	HorizontalTop, HorizontalBottom string
}

var (
	_plainBorderSet = BorderSet{
		TopLeft:          symbol.LineTopLeft,
		TopRight:         symbol.LineTopRight,
		BottomLeft:       symbol.LineBottomLeft,
		BottomRight:      symbol.LineBottomRight,
		VerticalLeft:     symbol.LineVertical,
		VerticalRight:    symbol.LineVertical,
		HorizontalTop:    symbol.LineHorizontal,
		HorizontalBottom: symbol.LineHorizontal,
	}

	_roundedBorderSet = BorderSet{
		TopLeft:          symbol.LineRoundedTopLeft,
		TopRight:         symbol.LineRoundedTopRight,
		BottomLeft:       symbol.LineRoundedBottomLeft,
		BottomRight:      symbol.LineRoundedBottomRight,
		VerticalLeft:     symbol.LineVertical,
		VerticalRight:    symbol.LineVertical,
		HorizontalTop:    symbol.LineHorizontal,
		HorizontalBottom: symbol.LineHorizontal,
	}
)
