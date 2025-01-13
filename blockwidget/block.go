package blockwidget

import (
	"slices"

	"github.com/metafates/bento"
	"github.com/metafates/bento/internal/bit"
)

var _ bento.Widget = (*Block)(nil)

type Block struct {
	titles          []Title
	titlesStyle     bento.Style
	titlesAlignment bento.Alignment
	titlesPosition  TitlePosition

	borders     Side
	borderStyle bento.Style
	borderSet   BorderSet

	style   bento.Style
	padding bento.Padding
}

func New() Block {
	return Block{
		titles:          nil,
		titlesStyle:     bento.Style{},
		titlesAlignment: bento.AlignmentLeft,
		borders:         SideNone,
		borderStyle:     bento.NewStyle(),
		borderSet:       BorderTypeSharp.Set(),
		style:           bento.NewStyle(),
		padding:         bento.NewPadding(),
	}
}

func (b Block) Rounded() Block {
	return b.WithBorderType(BorderTypeRounded)
}

func (b Block) Sharp() Block {
	return b.WithBorderType(BorderTypeSharp)
}

func (b Block) Thick() Block {
	return b.WithBorderType(BorderTypeThick)
}

func (b Block) Double() Block {
	return b.WithBorderType(BorderTypeDouble)
}

func (b Block) WithTitlesStyle(style bento.Style) Block {
	b.titlesStyle = style
	return b
}

func (b Block) WithTitlesAlignment(alignment bento.Alignment) Block {
	b.titlesAlignment = alignment
	return b
}

func (b Block) WithStyle(style bento.Style) Block {
	b.style = style
	return b
}

func (b Block) WithTitle(title Title) Block {
	if title.Alignment != bento.AlignmentNone {
		title.Line = title.Line.WithAlignment(title.Alignment)
	}

	b.titles = append(b.titles, title)
	return b
}

func (b Block) WithTitlePosition(position TitlePosition) Block {
	b.titlesPosition = position
	return b
}

func (b Block) WithTitleStr(title string) Block {
	return b.WithTitle(NewTitleStr(title))
}

func (b Block) Bordered() Block {
	return b.WithBorderSides()
}

func (b Block) WithBorderSides(borders ...Side) Block {
	if len(borders) == 0 {
		b.borders = SideAll
		return b
	}

	b.borders = SideNone

	for _, border := range borders {
		b.borders = bit.Union(b.borders, border)
	}

	return b
}

func (b Block) WithBorderType(borderType BorderType) Block {
	b.borderSet = borderType.Set()
	return b
}

func (b Block) WithPadding(padding bento.Padding) Block {
	b.padding = padding
	return b
}

func (b Block) WithBorderStyle(style bento.Style) Block {
	b.borderStyle = style
	return b
}

func (b Block) Inner(area bento.Rect) bento.Rect {
	inner := area

	if b.borders.intersects(SideLeft) {
		inner.X = min(inner.X+1, inner.Right())
		inner.Width = max(0, inner.Width-1)
	}

	if b.borders.intersects(SideTop) || b.hasTitleAtPosition(TitlePositionTop) {
		inner.Y = min(inner.Y+1, inner.Bottom())
		inner.Height = max(0, inner.Height-1)
	}

	if b.borders.intersects(SideRight) {
		inner.Width = max(0, inner.Width-1)
	}

	if b.borders.intersects(SideBottom) || b.hasTitleAtPosition(TitlePositionBottom) {
		inner.Height = max(0, inner.Height-1)
	}

	inner.X = inner.X + b.padding.Left
	inner.Y = inner.Y + b.padding.Top

	inner.Width = max(0, inner.Width-b.padding.Left-b.padding.Right)
	inner.Height = max(0, inner.Height-b.padding.Top-b.padding.Bottom)

	return inner
}

func (b Block) hasTitleAtPosition(position TitlePosition) bool {
	for _, t := range b.titles {
		p := b.titlesPosition
		if t.Position != nil {
			p = *t.Position
		}

		if p == position {
			return true
		}
	}

	return false
}

func (b Block) Render(area bento.Rect, buffer *bento.Buffer) {
	area = area.Intersection(area)
	if area.IsEmpty() {
		return
	}

	buffer.SetStyle(area, b.style)

	b.renderBorders(area, buffer)
	b.renderTitles(area, buffer)
}

func (b Block) renderBorders(area bento.Rect, buffer *bento.Buffer) {
	b.renderLeftSide(area, buffer)
	b.renderTopSide(area, buffer)
	b.renderRightSide(area, buffer)
	b.renderBottomSide(area, buffer)

	b.renderBottomRightCorner(area, buffer)
	b.renderTopRightCorner(area, buffer)
	b.renderBottomLeftCorner(area, buffer)
	b.renderTopLeftCorner(area, buffer)
}

func (b Block) renderLeftSide(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideLeft) {
		for y := area.Top(); y < area.Bottom(); y++ {
			buffer.
				CellAt(bento.Position{
					X: area.Left(),
					Y: y,
				}).
				SetSymbol(b.borderSet.VerticalLeft).
				SetStyle(b.borderStyle)
		}
	}
}

func (b Block) renderTopSide(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideTop) {
		for x := area.Left(); x < area.Right(); x++ {
			buffer.
				CellAt(bento.Position{
					X: x,
					Y: area.Top(),
				}).
				SetSymbol(b.borderSet.HorizontalTop).
				SetStyle(b.borderStyle)
		}
	}
}

func (b Block) renderRightSide(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideRight) {
		x := area.Right() - 1

		for y := area.Top(); y < area.Bottom(); y++ {
			buffer.
				CellAt(bento.Position{
					X: x,
					Y: y,
				}).
				SetSymbol(b.borderSet.VerticalRight).
				SetStyle(b.borderStyle)
		}
	}
}

func (b Block) renderBottomSide(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideBottom) {
		y := area.Bottom() - 1

		for x := area.Left(); x < area.Right(); x++ {
			buffer.
				CellAt(bento.Position{
					X: x,
					Y: y,
				}).
				SetSymbol(b.borderSet.HorizontalBottom).
				SetStyle(b.borderStyle)
		}
	}
}

func (b Block) renderBottomRightCorner(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideRight | SideBottom) {
		buffer.
			CellAt(bento.Position{
				X: area.Right() - 1,
				Y: area.Bottom() - 1,
			}).
			SetSymbol(b.borderSet.BottomRight).
			SetStyle(b.borderStyle)
	}
}

func (b Block) renderTopRightCorner(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideRight | SideTop) {
		buffer.
			CellAt(bento.Position{
				X: area.Right() - 1,
				Y: area.Top(),
			}).
			SetSymbol(b.borderSet.TopRight).
			SetStyle(b.borderStyle)
	}
}

func (b Block) renderBottomLeftCorner(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideLeft | SideBottom) {
		buffer.
			CellAt(bento.Position{
				X: area.Left(),
				Y: area.Bottom() - 1,
			}).
			SetSymbol(b.borderSet.BottomLeft).
			SetStyle(b.borderStyle)
	}
}

func (b Block) renderTopLeftCorner(area bento.Rect, buffer *bento.Buffer) {
	if b.borders.contains(SideLeft | SideTop) {
		buffer.
			CellAt(bento.Position{
				X: area.Left(),
				Y: area.Top(),
			}).
			SetSymbol(b.borderSet.TopLeft).
			SetStyle(b.borderStyle)
	}
}

func (b Block) renderTitles(area bento.Rect, buffer *bento.Buffer) {
	b.renderTitlePosition(TitlePositionTop, area, buffer)
	b.renderTitlePosition(TitlePositionBottom, area, buffer)
}

func (b Block) renderTitlePosition(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {
	// NOTE: the order in which these functions are called defines the overlapping behavior
	b.renderRightTitles(position, area, buffer)
	b.renderCenterTitles(position, area, buffer)
	b.renderLeftTitles(position, area, buffer)
}

func (b Block) renderRightTitles(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {
	titles := b.filterTitles(position, bento.AlignmentRight)
	titlesArea := b.titlesArea(area, position)

	slices.Reverse(titles)

	for _, t := range titles {
		if titlesArea.IsEmpty() {
			break
		}

		titleWidth := t.Line.Width()
		titleArea := bento.Rect{
			X:      max(0, max(titlesArea.Left(), titlesArea.Right()-titleWidth)),
			Y:      titlesArea.Y,
			Width:  min(titlesArea.Width, titleWidth),
			Height: titlesArea.Height,
		}

		buffer.SetStyle(titleArea, b.titlesStyle)
		t.Line.Render(titleArea, buffer)

		titlesArea.Width = max(0, titlesArea.Width-titleWidth-1)
	}
}

func (b Block) renderCenterTitles(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {
	titles := b.filterTitles(position, bento.AlignmentCenter)

	var totalWidth int

	for _, t := range titles {
		width := t.Line.Width() + 1

		totalWidth += width
	}

	totalWidth = max(0, totalWidth-1)

	titlesArea := b.titlesArea(area, position)

	titlesArea = bento.Rect{
		X:      titlesArea.Left() + max(0, titlesArea.Width-totalWidth)/2,
		Y:      titlesArea.Y,
		Width:  titlesArea.Width,
		Height: titlesArea.Height,
	}

	for _, t := range titles {
		if titlesArea.IsEmpty() {
			break
		}

		titleWidth := t.Line.Width()

		titleArea := bento.Rect{
			X:      titlesArea.X,
			Y:      titlesArea.Y,
			Width:  min(titleWidth, titlesArea.Width),
			Height: titlesArea.Height,
		}

		buffer.SetStyle(titleArea, b.titlesStyle)
		t.Line.Render(titleArea, buffer)

		titlesArea.X += titleWidth + 1
		titlesArea.Width = max(0, titlesArea.Width-titleWidth-1)
	}
}

func (b Block) renderLeftTitles(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {
	titles := b.filterTitles(position, bento.AlignmentLeft)
	titlesArea := b.titlesArea(area, position)

	for _, t := range titles {
		if titlesArea.IsEmpty() {
			break
		}

		titleWidth := t.Line.Width()
		titleArea := bento.Rect{
			X:      titlesArea.X,
			Y:      titlesArea.Y,
			Width:  min(titleWidth, titlesArea.Width),
			Height: titlesArea.Height,
		}

		buffer.SetStyle(titleArea, b.titlesStyle)
		t.Line.Render(titleArea, buffer)

		titlesArea.X += titleWidth + 1
		titlesArea.Width = max(0, titlesArea.Width-titleWidth-1)
	}
}

func (b Block) filterTitles(position TitlePosition, alignment bento.Alignment) []Title {
	titles := make([]Title, 0, len(b.titles))

	for _, t := range b.titles {
		titlePosition := b.titlesPosition
		if t.Position != nil {
			titlePosition = *t.Position
		}

		if titlePosition != position {
			continue
		}

		titleAlignment := b.titlesAlignment
		if t.Alignment != bento.AlignmentNone {
			titleAlignment = t.Alignment
		}

		if titleAlignment != alignment {
			continue
		}

		titles = append(titles, t)
	}

	return titles
}

func (b Block) titlesArea(area bento.Rect, position TitlePosition) bento.Rect {
	var leftBorder, rightBorder int

	if b.borders.contains(SideLeft) {
		leftBorder = 1
	}

	if b.borders.contains(SideRight) {
		rightBorder = 1
	}

	var y int

	switch position {
	case TitlePositionTop:
		y = area.Top()
	case TitlePositionBottom:
		y = area.Bottom() - 1
	}

	return bento.Rect{
		X:      area.Left() + leftBorder,
		Y:      y,
		Width:  max(0, area.Width-leftBorder-rightBorder),
		Height: 1,
	}
}
