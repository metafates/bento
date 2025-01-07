package bento

type Rect struct {
	X, Y int

	Width, Height int
}

func (r Rect) IsEmpty() bool {
	return r.Width == 0 || r.Height == 0
}

func (r Rect) Area() int {
	return r.Width * r.Height
}

func (r Rect) Left() int {
	return r.X
}

func (r Rect) Right() int {
	return r.X + r.Width
}

func (r Rect) Top() int {
	return r.Y
}

func (r Rect) Bottom() int {
	return r.Y + r.Height
}

func (r Rect) Position() Position {
	return Position{
		X: r.X,
		Y: r.Y,
	}
}

func (r Rect) Inner(margin Margin) Rect {
	doubleHorizontal := margin.Horizontal * 2
	doubleVertical := margin.Vertical * 2

	if r.Width < doubleHorizontal || r.Height < doubleVertical {
		return Rect{}
	}

	return Rect{
		X:      r.X + margin.Horizontal,
		Y:      r.Y + margin.Vertical,
		Width:  max(0, r.Width-doubleHorizontal),
		Height: max(0, r.Height-doubleVertical),
	}
}

func (r Rect) Intersection(other Rect) Rect {
	x1 := max(r.X, other.X)
	y1 := max(r.Y, other.Y)

	x2 := min(r.Right(), other.Right())
	y2 := min(r.Bottom(), other.Bottom())

	return Rect{
		X:      x1,
		Y:      y1,
		Width:  max(0, x2-x1),
		Height: max(0, y2-y1),
	}
}

func (r Rect) Contains(position Position) bool {
	return position.X >= r.X && position.X < r.Right() && position.Y >= r.Y && position.Y < r.Bottom()
}

func (r Rect) Rows() []Rect {
	currentRowFwd := r.Y
	currentRowBack := r.Bottom()

	var rows []Rect

	for currentRowFwd < currentRowBack {
		row := Rect{
			X:      r.X,
			Y:      currentRowFwd,
			Width:  r.Width,
			Height: 1,
		}

		currentRowFwd++

		rows = append(rows, row)
	}

	return rows
}

func (r Rect) IndentX(offset int) Rect {
	return Rect{
		X:      r.X + offset,
		Y:      r.Y,
		Width:  max(0, r.Width-offset),
		Height: r.Height,
	}
}
