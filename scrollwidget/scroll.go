package scrollwidget

import (
	"math"

	"github.com/metafates/bento"
	"github.com/metafates/bento/internal/sliceutil"
	"github.com/rivo/uniseg"
)

var _ bento.StatefulWidget[State] = (*Scroll)(nil)

type Scroll struct {
	orientation Orientation

	thumbStyle  bento.Style
	thumbSymbol string

	trackStyle  bento.Style
	trackSymbol *string

	beginStyle  bento.Style
	beginSymbol *string

	endStyle  bento.Style
	endSymbol *string
}

func New(orientation Orientation) Scroll {
	symbols := SymbolsDoubleHorizontal
	if orientation.IsVertical() {
		symbols = SymbolsDoubleVertical
	}

	set := symbols.Set()

	return Scroll{
		orientation: orientation,
		thumbStyle:  bento.NewStyle(),
		thumbSymbol: set.Thumb,
		trackStyle:  bento.NewStyle(),
		trackSymbol: &set.Track,
		beginStyle:  bento.NewStyle(),
		beginSymbol: &set.Begin,
		endStyle:    bento.NewStyle(),
		endSymbol:   &set.End,
	}
}

func (s Scroll) RenderStateful(area bento.Rect, buffer *bento.Buffer, state State) {
	if state.contentLen == 0 || s.trackLenExcludingArrowHeads(area) == 0 {
		return
	}

	area, ok := s.scrollbarArea(area)
	if !ok {
		return
	}

	areas := sliceutil.FlatMap(area.Columns(), func(rect bento.Rect) []bento.Rect {
		return rect.Rows()
	})

	barSymbols := s.barSymbols(area, state)

	for i := 0; i < min(len(areas), len(barSymbols)); i++ {
		bar := barSymbols[i]

		if bar == nil {
			continue
		}

		symbol, style := bar.Symbol, bar.Style
		x, y := areas[i].X, areas[i].Y

		buffer.SetString(x, y, symbol, style)
	}
}

type _Symbol struct {
	Symbol string
	Style  bento.Style
}

func (s Scroll) barSymbols(area bento.Rect, state State) []*_Symbol {
	trackStartLen, thumbLen, trackEndLen := s.partLens(area, state)

	var begin *_Symbol

	if s.beginSymbol != nil {
		begin = &_Symbol{
			Symbol: *s.beginSymbol,
			Style:  s.beginStyle,
		}
	}

	var track *_Symbol

	if s.trackSymbol != nil {
		track = &_Symbol{
			Symbol: *s.trackSymbol,
			Style:  s.trackStyle,
		}
	}

	thumb := _Symbol{
		Symbol: s.thumbSymbol,
		Style:  s.thumbStyle,
	}

	var end *_Symbol

	if s.endSymbol != nil {
		end = &_Symbol{
			Symbol: *s.endSymbol,
			Style:  s.endStyle,
		}
	}

	// <
	symbols := []*_Symbol{begin}

	// <═══
	symbols = append(symbols, sliceutil.Repeat(track, trackStartLen)...)

	// <═══█████
	symbols = append(symbols, sliceutil.Repeat(&thumb, thumbLen)...)

	// <═══█████═══════
	symbols = append(symbols, sliceutil.Repeat(track, trackEndLen)...)

	// <═══█████═══════>
	symbols = append(symbols, end)

	return symbols
}

// Calculates length of the track excluding the arrow heads
//
//	       ┌────────── track_length
//	 vvvvvvvvvvvvvvv
//	<═══█████═══════>
func (s Scroll) trackLenExcludingArrowHeads(area bento.Rect) int {
	var startLen int

	if s.beginSymbol != nil {
		startLen = uniseg.StringWidth(*s.beginSymbol)
	}

	var endLen int

	if s.endSymbol != nil {
		endLen = uniseg.StringWidth(*s.endSymbol)
	}

	arrowsLen := startLen + endLen

	if s.orientation.IsVertical() {
		return max(0, area.Height-arrowsLen)
	}

	return max(0, area.Width-arrowsLen)
}

// Returns the lengths of the parts of a scrollbar
//
// The scrollbar has 3 parts of note:
//
//	<═══█████═══════> full scrollbar
//	 ═══              track start
//	    █████         thumb
//	         ═══════  track end
//
// This method returns the length of the start, thumb, and end as a tuple.
func (s Scroll) partLens(
	area bento.Rect,
	state State,
) (thumbStart, thumbEnd, trackEndLen int) {
	trackLen := s.trackLenExcludingArrowHeads(area)
	viewportLen := s.viewportLen(area, state)

	// Ensure that the position of the thumb is within the bounds of the content taking into
	// account the content and viewport length. When the last line of the content is at the top
	// of the viewport, the thumb should be at the bottom of the track.
	maxPosition := max(0, state.contentLen-1)
	startPosition := max(0, min(maxPosition, state.position))
	maxViewportPosition := maxPosition + viewportLen
	endPosition := startPosition + viewportLen

	// Calculate the start and end positions of the thumb. The size will be proportional to the
	// viewport length compared to the total amount of possible visible rows.
	thumbStartF := float64(startPosition*trackLen) / float64(maxViewportPosition)
	thumbEndF := float64(endPosition*trackLen) / float64(maxViewportPosition)

	// Make sure that the thumb is at least 1 cell long by ensuring that the start of the thumb
	// is less than the track_len. We use the positions instead of the sizes and use nearest
	// integer instead of floor / ceil to avoid problems caused by rounding errors.
	thumbStart = max(0, min(trackLen-1, int(math.Round(thumbStartF))))
	thumbEnd = max(0, min(trackLen, int(math.Round(thumbEndF))))

	thumbLen := max(1, max(0, thumbEnd-thumbStart))
	trackEndLen = max(0, trackLen-thumbStart-thumbLen)

	return thumbStart, thumbLen, trackEndLen
}

func (s Scroll) viewportLen(area bento.Rect, state State) int {
	if state.viewportContentLen != 0 {
		return state.viewportContentLen
	}

	if s.orientation.IsVertical() {
		return area.Height
	}

	return area.Width
}

func (s Scroll) scrollbarArea(area bento.Rect) (bento.Rect, bool) {
	switch s.orientation {
	case OrientationVerticalLeft:
		columns := area.Columns()
		if len(columns) == 0 {
			return bento.Rect{}, false
		}

		return columns[0], true

	case OrientationVerticalRight:
		columns := area.Columns()
		if len(columns) == 0 {
			return bento.Rect{}, false
		}

		return columns[len(columns)-1], true

	case OrientationHorizontalBottom:
		rows := area.Rows()
		if len(rows) == 0 {
			return bento.Rect{}, false
		}

		return rows[0], true

	case OrientationHorizontalTop:
		rows := area.Rows()
		if len(rows) == 0 {
			return bento.Rect{}, false
		}

		return rows[len(rows)-1], true

	default:
		return bento.Rect{}, false
	}
}
