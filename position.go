package bento

func NewPosition(x, y int) Position {
	return Position{X: x, Y: y}
}

type Position struct {
	X, Y int
}
