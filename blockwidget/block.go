package blockwidget

import "github.com/metafates/bento"

var _ bento.Widget = (*Block)(nil)

type Block struct {
	titles          []Title
	titlesStyle     bento.Style
	titlesAlignment bento.Alignment
	titlesPosition  TitlePosition

	borders     Borders
	borderStyle bento.Style
	borderSet   BorderSet

	style   bento.Style
	padding Padding
}

func NewBlock() Block {
	return Block{}
}

func (b Block) WithTitle(title Title) Block {
	if title.Alignment != bento.AlignmentNone {
		title.Line = title.Line.WithAlignment(title.Alignment)
	}

	b.titles = append(b.titles, title)
	return b
}

func (b Block) WithBorders(borders Borders) Block {
	b.borders = borders
	return b
}

func (b Block) WithBorderType(borderType BorderType) Block {
	b.borderSet = borderType.Set()
	return b
}

func (b Block) WithPadding(padding Padding) Block {
	b.padding = padding
	return b
}

func (b Block) Inner(area bento.Rect) bento.Rect {
	inner := area

	if b.borders.intersects(BordersLeft) {
		inner.X = min(inner.X, inner.Right())
		inner.Width = max(0, inner.Width-1)
	}

	if b.borders.intersects(BordersTop) || b.hasTitleAtPosition(TitlePositionTop) {
		inner.Y = min(inner.Y, inner.Bottom())
		inner.Height = max(0, inner.Height-1)
	}

	if b.borders.intersects(BordersRight) {
		inner.Width = max(0, inner.Width-1)
	}

	if b.borders.intersects(BordersBottom) || b.hasTitleAtPosition(TitlePositionBottom) {
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
	area = area.Intersection(buffer.Area)
	buffer.SetStyle(area, b.style)
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
	if b.borders.contains(BordersLeft) {
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
	if b.borders.contains(BordersTop) {
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

func (b Block) renderRightSide(area bento.Rect, buffer *bento.Buffer)  {}
func (b Block) renderBottomSide(area bento.Rect, buffer *bento.Buffer) {}

func (b Block) renderBottomRightCorner(area bento.Rect, buffer *bento.Buffer) {}
func (b Block) renderTopRightCorner(area bento.Rect, buffer *bento.Buffer)    {}
func (b Block) renderBottomLeftCorner(area bento.Rect, buffer *bento.Buffer)  {}
func (b Block) renderTopLeftCorner(area bento.Rect, buffer *bento.Buffer)     {}

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

func (b Block) renderRightTitles(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {}

func (b Block) renderCenterTitles(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {}

func (b Block) renderLeftTitles(position TitlePosition, area bento.Rect, buffer *bento.Buffer) {}
