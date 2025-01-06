package ansi

import "strconv"

type Attr int

var _sgr = [...]int{0, 1, 2, 3, 4, 2, 3, 4, 5, 5, 6, 7, 8, 9, 20, 21, 22, 23, 24, 25, 27, 28, 29, 51, 52, 53, 54, 55}

const (
	AttrReset = iota
	AttrBold
	AttrDim
	AttrItalic
	AttrUnderlined
	AttrDoubleUnderlined
	AttrUndercurled
	AttrUnderdotted
	AttrUnderdashed
	AttrSlowBlink
	AttrRapidBlink
	AttrReverse
	AttrHidden
	AttrCrossedOut
	AttrFraktur
	AttrNoBold
	AttrNormalIntensity
	AttrNoItalic
	AttrNoUnderline
	AttrNoBlink
	AttrNoReverse
	AttrNoHidden
	AttrNotCrossedOut
	AttrFramed
	AttrEncircled
	AttrOverLined
	AttrNotFramedOrEncircled
	AttrNotOverLined
)

func (a Attr) Bytes() uint32 {
	return 1 << (uint32(a) + 1)
}

func (a Attr) SGR() string {
	if a > 4 && a < 9 {
		return "4:" + strconv.Itoa(_sgr[a])
	}

	return strconv.Itoa(_sgr[a])
}
