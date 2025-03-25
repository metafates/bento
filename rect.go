package bento

type Rect struct {
	Position
	// X, Y int

	Width, Height int
}

func NewRect(width, height int) Rect {
	return Rect{Width: width, Height: height}
}

func (r Rect) Positioned(x, y int) Rect {
	r.X = x
	r.Y = y
	return r
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

func (r Rect) Inner(padding Padding) Rect {
	horizontal := padding.Right + padding.Left
	vertical := padding.Top + padding.Bottom

	if r.Width < horizontal || r.Height < vertical {
		return Rect{}
	}

	return Rect{
		Position: Position{
			X: r.X + padding.Left,
			Y: r.Y + padding.Top,
		},
		Width:  max(0, r.Width-horizontal),
		Height: max(0, r.Height-vertical),
	}
}

func (r Rect) Intersection(other Rect) Rect {
	x1 := max(r.X, other.X)
	y1 := max(r.Y, other.Y)

	x2 := min(r.Right(), other.Right())
	y2 := min(r.Bottom(), other.Bottom())

	return Rect{
		Position: Position{
			X: x1,
			Y: y1,
		},
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
			Position: Position{
				X: r.X,
				Y: currentRowFwd,
			},
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
		Position: Position{
			X: r.X + offset,
			Y: r.Y,
		},
		Width:  max(0, r.Width-offset),
		Height: r.Height,
	}
}

func (r Rect) Columns() []Rect {
	currentColumnFwd := r.X
	currentColumnBack := r.Right()

	var columns []Rect

	for currentColumnFwd < currentColumnBack {
		column := NewRect(1, r.Height).Positioned(currentColumnFwd, r.Y)
		currentColumnFwd++

		columns = append(columns, column)
	}

	return columns
}
