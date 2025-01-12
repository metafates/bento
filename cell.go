package bento

import (
	"github.com/metafates/bento/internal/bit"
	"github.com/muesli/termenv"
)

type Cell struct {
	Symbol   string
	Fg, Bg   Color
	Modifier Modifier
	Skip     bool
}

func NewEmptyCell() Cell {
	return NewCell(" ")
}

func NewCell(symbol string) Cell {
	return Cell{
		Symbol: symbol,
		Fg:     ResetColor{},
		Bg:     ResetColor{},
		Skip:   false,
	}
}

func (c *Cell) SetSymbol(symbol string) *Cell {
	c.Symbol = symbol

	return c
}

func (c *Cell) AppendSymbol(symbol string) *Cell {
	c.Symbol += symbol

	return c
}

func (c *Cell) SetFg(color termenv.Color) *Cell {
	c.Fg = color
	return c
}

func (c *Cell) SetBg(color termenv.Color) *Cell {
	c.Bg = color
	return c
}

func (c *Cell) SetStyle(style Style) *Cell {
	if style.Foreground.IsSet() {
		c.Fg = style.Foreground.Color()
	}

	if style.Background.IsSet() {
		c.Bg = style.Background.Color()
	}

	c.Modifier = bit.Union(c.Modifier, style.addModifier)
	c.Modifier = bit.Difference(c.Modifier, style.subModifier)

	return c
}

func (c *Cell) Reset() {
	c.Symbol = " "
	c.Skip = false
	c.Fg = ResetColor{}
	c.Bg = ResetColor{}
	c.Modifier = 0
}
