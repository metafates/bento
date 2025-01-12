package bento

type Padding struct {
	Top, Right, Bottom, Left int
}

func NewPadding(sides ...int) Padding {
	switch len(sides) {
	case 0:
		return Padding{}
	case 1:
		side := sides[0]
		return Padding{Top: side, Right: side, Bottom: side, Left: side}
	case 2:
		vertical := sides[0]
		horizontal := sides[1]

		return Padding{
			Top:    vertical,
			Right:  horizontal,
			Bottom: vertical,
			Left:   horizontal,
		}
	case 4:
		return Padding{
			Top:    sides[0],
			Right:  sides[1],
			Bottom: sides[2],
			Left:   sides[3],
		}
	default:
		panic("unexpected sides count")
	}
}

func (p Padding) WithTop(s int) Padding {
	p.Top = s
	return p
}

func (p Padding) WithRight(s int) Padding {
	p.Right = s
	return p
}

func (p Padding) WithBottom(s int) Padding {
	p.Bottom = s
	return p
}

func (p Padding) WithLeft(s int) Padding {
	p.Bottom = s
	return p
}
