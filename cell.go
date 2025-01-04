package bento

type Cell struct {
	Symbol string
	Style  Style
	Skip   bool
}

func NewCell(symbol string) *Cell {
	return &Cell{
		Symbol: symbol,
		Style:  NewStyle(),
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

func (c *Cell) PatchStyle(style Style) *Cell {
	c.Style = c.Style.Patched(style)

	return c
}

func (c *Cell) Reset() {
	c.Symbol = " "
	c.Style = NewStyle()
	c.Skip = false
}
