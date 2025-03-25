package scrollwidget

import "github.com/metafates/bento/symbol"

type Symbols int

const (
	SymbolsDoubleVertical Symbols = iota
	SymbolsVertical
	SymbolsDoubleHorizontal
	SymbolsHorizontal
)

func (s Symbols) Set() Set {
	switch s {
	case SymbolsDoubleVertical:
		return _doubleVertical
	case SymbolsVertical:
		return _vertical
	case SymbolsDoubleHorizontal:
		return _doubleHorizontal
	case SymbolsHorizontal:
		return _horizontal
	default:
		return Set{}
	}
}

// Set for scrollbar
//
//	<--▮------->
//	^  ^   ^   ^
//	│  │   │   └ end
//	│  │   └──── track
//	│  └──────── thumb
//	└─────────── begin
type Set struct {
	Track, Thumb, Begin, End string
}

var _doubleVertical = Set{
	Track: symbol.LineDoubleVertical,
	Thumb: symbol.BlockFull,
	Begin: symbol.ArrowUp,
	End:   symbol.ArrowDown,
}

var _vertical = Set{
	Track: symbol.LineVertical,
	Thumb: symbol.BlockFull,
	Begin: symbol.ArrowLeft,
	End:   symbol.ArrowRight,
}

var _doubleHorizontal = Set{
	Track: symbol.LineDoubleHorizontal,
	Thumb: symbol.BlockFull,
	Begin: "",
	End:   "",
}

var _horizontal = Set{
	Track: symbol.LineHorizontal,
	Thumb: symbol.BlockFull,
	Begin: "",
	End:   "",
}
