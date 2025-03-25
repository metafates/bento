package bento

import (
	"math"
	"slices"

	"github.com/rivo/uniseg"
)

type Buffer struct {
	area Rect

	// content of the buffer. The length of this Vec should always be equal to [Area.Width] * [Area.Height]
	content []Cell
}

func NewBufferEmpty(area Rect) Buffer {
	return NewBufferFilled(area, NewEmptyCell())
}

func NewBufferFilled(area Rect, cell Cell) Buffer {
	size := area.Area()

	content := make([]Cell, 0, size)
	for i := 0; i < size; i++ {
		content = append(content, cell)
	}

	return Buffer{
		area:    area,
		content: content,
	}
}

// Area of the buffer
func (b *Buffer) Area() Rect {
	return b.area
}

// Diff builds a minimal sequence of coordinates and Cells necessary to update the UI from
// self to other.
//
// We're assuming that buffers are well-formed, that is no double-width cell is followed by
// a non-blank cell.
func (b *Buffer) Diff(other *Buffer) []PositionedCell {
	prevBuffer := b.content
	nextBuffer := other.content

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
	if index >= len(b.content) {
		panic("trying to get coords of a cell outside the buffer")
	}

	x := index%b.area.Width + b.area.X
	y := index/b.area.Width + b.area.Y

	return NewPosition(x, y)
}

// SetString prints a string, starting at the position (x, y)
func (b *Buffer) SetString(x, y int, value string, style Style) {
	b.SetStringN(x, y, value, math.MaxInt, style)
}

func (b *Buffer) Reset() {
	for i := range b.content {
		b.content[i].Reset()
	}
}

// SetStringN prints at most the first n characters of a string if enough space is available
// until the end of the line. Skips zero-width graphemes and control characters.
//
// Use [Buffer.SetString] when the maximum amount of characters can be printed.
func (b *Buffer) SetStringN(x, y int, value string, maxWidth int, style Style) (int, int) {
	remainingWidth := min(maxWidth, max(0, b.area.Right()-x))

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
	return &b.content[b.indexOf(position)]
}

// indexOf returns the (global) coordinates of a cell given its index
//
// Global coordinates are offset by the Buffer's area offset (`x`/`y`).
//
// Panics when given an index that is outside the Buffer's content.
func (b *Buffer) indexOf(position Position) int {
	if !b.area.Contains(position) {
		panic("position out of bounds")
	}

	// remove offset
	y := max(0, position.Y-b.area.Y)
	x := max(0, position.X-b.area.X)

	width := b.area.Width

	return y*width + x
}

func (b *Buffer) Resize(area Rect) {
	length := area.Area()

	if len(b.content) > length {
		b.content = slices.Delete(b.content, length, len(b.content))
	} else {
		toAdd := length - len(b.content)

		for i := 0; i < toAdd; i++ {
			b.content = append(b.content, NewEmptyCell())
		}
	}

	b.area = area
}

func (b *Buffer) SetStyle(area Rect, style Style) {
	area = b.area.Intersection(area)

	for y := area.Top(); y < area.Bottom(); y++ {
		for x := area.Left(); x < area.Right(); x++ {
			pos := NewPosition(x, y)

			b.CellAt(pos).SetStyle(style)
		}
	}
}
