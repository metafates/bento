package bento

type Margin struct {
	Top, Right, Bottom, Left int
}

func NewMargin(sides ...int) Margin {
	switch len(sides) {
	case 0:
		return Margin{}
	case 1:
		side := sides[0]
		return Margin{Top: side, Right: side, Bottom: side, Left: side}
	case 2:
		horizontal := sides[0]
		vertical := sides[1]

		return Margin{
			Top:    horizontal,
			Right:  vertical,
			Bottom: horizontal,
			Left:   vertical,
		}
	case 4:
		return Margin{
			Top:    sides[0],
			Right:  sides[1],
			Bottom: sides[2],
			Left:   sides[3],
		}
	default:
		panic("unexpected sides count")
	}
}

func (m Margin) WithTop(s int) Margin {
	m.Top = s
	return m
}

func (m Margin) WithRight(s int) Margin {
	m.Right = s
	return m
}

func (m Margin) WithBottom(s int) Margin {
	m.Bottom = s
	return m
}

func (m Margin) WithLeft(s int) Margin {
	m.Bottom = s
	return m
}
