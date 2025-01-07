package bento

import (
	"math"
	"slices"

	"github.com/rivo/uniseg"
)

type Buffer struct {
	Area Rect

	// Content of the buffer. The length of this Vec should always be equal to [Area.Width] * [Area.Height]
	Content []Cell
}

func NewBufferEmpty(area Rect) *Buffer {
	return NewBufferFilled(area, NewEmptyCell())
}

func NewBufferFilled(area Rect, cell Cell) *Buffer {
	size := area.Area()

	content := make([]Cell, 0, size)
	for i := 0; i < size; i++ {
		content = append(content, cell)
	}

	return &Buffer{
		Area:    area,
		Content: content,
	}
}

// Diff builds a minimal sequence of coordinates and Cells necessary to update the UI from
// self to other.
//
// We're assuming that buffers are well-formed, that is no double-width cell is followed by
// a non-blank cell.
func (b *Buffer) Diff(other *Buffer) []PositionedCell {
	prevBuffer := b.Content
	nextBuffer := other.Content

	var (
		updates     []PositionedCell
		invalidated int
		toSkip      int
	)

	for i := 0; i < min(len(nextBuffer), len(prevBuffer)); i++ {
		current := nextBuffer[i]
		previous := prevBuffer[i]

		if !current.Skip && (current != previous || invalidated > 0) && toSkip == 0 {
			pos := b.PosOf(i)

			updates = append(updates, PositionedCell{
				Cell:     nextBuffer[i],
				Position: pos,
			})
		}

		previousWidth := uniseg.StringWidth(previous.Symbol)
		currentWidth := uniseg.StringWidth(current.Symbol)

		toSkip = max(0, currentWidth-1)

		affectedWidth := max(previousWidth, currentWidth)

		invalidated = max(0, max(invalidated, affectedWidth)-1)
	}

	return updates
}

func (b *Buffer) PosOf(index int) Position {
	if index >= len(b.Content) {
		panic("trying to get coords of a cell outside the buffer")
	}

	x := index%b.Area.Width + b.Area.X
	y := index/b.Area.Width + b.Area.Y

	return Position{X: x, Y: y}
}

// SetString prints a string, starting at the position (x, y)
func (b *Buffer) SetString(x, y int, value string, style Style) {
	b.SetStringN(x, y, value, math.MaxInt, style)
}

func (b *Buffer) Reset() {
	for i := range b.Content {
		b.Content[i].Reset()
	}
}

// SetStringN prints at most the first n characters of a string if enough space is available
// until the end of the line. Skips zero-width graphemes and control characters.
//
// Use [Buffer.SetString] when the maximum amount of characters can be printed.
func (b *Buffer) SetStringN(x, y int, value string, maxWidth int, style Style) (int, int) {
	remainingWidth := min(maxWidth, max(0, b.Area.Right()-x))

	graphemes := uniseg.NewGraphemes(value)

	for remainingWidth > 0 && graphemes.Next() {
		symbol := graphemes.Str()
		width := graphemes.Width()

		remainingWidth -= width

		b.CellAt(Position{x, y}).SetSymbol(symbol).SetStyle(style)

		nextSymbol := x + width
		x++

		for x < nextSymbol {
			b.CellAt(Position{x, y}).Reset()
			x++
		}
	}

	return x, y
}

func (b *Buffer) CellAt(position Position) *Cell {
	return &b.Content[b.indexOf(position)]
}

func (b *Buffer) indexOf(position Position) int {
	if !b.Area.Contains(position) {
		panic("position out of bounds")
	}

	// remove offset
	y := max(0, position.Y-b.Area.Y)
	x := max(0, position.X-b.Area.X)

	width := b.Area.Width

	return y*width + x
}

func (b *Buffer) Resize(area Rect) {
	length := area.Area()

	if len(b.Content) > length {
		b.Content = slices.Delete(b.Content, length, len(b.Content))
	} else {
		toAdd := length - len(b.Content)

		for i := 0; i < toAdd; i++ {
			b.Content = append(b.Content, NewEmptyCell())
		}
	}

	b.Area = area
}

func (b *Buffer) SetStyle(area Rect, style Style) {
	area = b.Area.Intersection(area)

	for y := area.Top(); y < area.Bottom(); y++ {
		for x := area.Left(); x < area.Right(); x++ {
			pos := Position{X: x, Y: y}

			b.CellAt(pos).SetStyle(style)
		}
	}
}
