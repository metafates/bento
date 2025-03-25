package scrollwidget

type Orientation int

const (
	OrientationVerticalRight Orientation = iota
	OrientationVerticalLeft
	OrientationHorizontalBottom
	OrientationHorizontalTop
)

func (o Orientation) IsVertical() bool {
	switch o {
	case OrientationVerticalLeft, OrientationVerticalRight:
		return true
	default:
		return false
	}
}
