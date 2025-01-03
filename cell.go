package bento

import "github.com/charmbracelet/lipgloss"

type Cell struct {
	Symbol string
	Style  lipgloss.Style
	Skip   bool
}

func NewCell(symbol string) *Cell {
	return &Cell{
		Symbol: symbol,
		Style:  lipgloss.NewStyle(),
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

func (c *Cell) SetStyle(style lipgloss.Style) *Cell {
	c.Style = style

	return c
}

func (c *Cell) Reset() {
	c.Symbol = " "
	c.Style = lipgloss.NewStyle()
	c.Skip = false
}
