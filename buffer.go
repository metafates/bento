package bento

import (
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/rivo/uniseg"
)

type Buffer struct {
	Area Rect

	// Content of the buffer. The length of this Vec should always be equal to [Area.Width] * [Area.Height]
	Content []Cell
}

func NewBufferEmpty(area Rect) *Buffer {
	return NewBufferFilled(area, Cell{Symbol: " "})
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

// SetString prints a string, starting at the position (x, y)
func (b *Buffer) SetString(x, y int, value string, style lipgloss.Style) {
	b.SetStringN(x, y, value, math.MaxInt, style)
}

// SetStringN prints at most the first n characters of a string if enough space is available
// until the end of the line. Skips zero-width graphemes and control characters.
//
// Use [Buffer.SetString] when the maximum amount of characters can be printed.
func (b *Buffer) SetStringN(x, y int, value string, maxWidth int, style lipgloss.Style) (int, int) {
	remainingWidth := min(maxWidth, b.Area.Right()-x)

	graphemes := uniseg.NewGraphemes(value)

	for remainingWidth > 0 && graphemes.Next() {
		symbol := graphemes.Str()
		width := graphemes.Width()

		remainingWidth -= width

		b.CellAt(Position{x, y}).Symbol = symbol
		b.CellAt(Position{x, y}).Style = style

		nextSymbol := x + width
		x++

		for x < nextSymbol {
			*b.CellAt(Position{x, y}) = Cell{}
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
	y := position.Y - b.Area.Y
	x := position.X - b.Area.X

	width := b.Area.Width

	return y*width + x
}

func (b *Buffer) Resize(area Rect) {
	length := area.Area()

	if len(b.Content) > length {
		// TODO: optimize with slices.Delete
		truncated := make([]Cell, length)
		copy(truncated, b.Content)

		b.Content = truncated
	} else {
		for i := 0; i < length-len(b.Content); i++ {
			b.Content = append(b.Content, Cell{})
		}
	}

	b.Area = area
}

func (b *Buffer) SetStyle(area Rect, style lipgloss.Style) {
	area = b.Area.Intersection(area)

	for y := area.Top(); y < area.Bottom(); y++ {
		for x := area.Left(); x < area.Right(); x++ {
			pos := Position{x, y}

			b.CellAt(pos).Style = style
		}
	}
}
