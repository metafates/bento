package tabs

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/clearwidget"
	"github.com/metafates/bento/examples/demo/theme"
	"github.com/metafates/bento/mascotwidget"
	"github.com/metafates/bento/paragraphwidget"
)

var _ bento.Widget = (*AboutTab)(nil)

type AboutTab struct {
	rowIndex int
}

func NewAboutTab() AboutTab {
	return AboutTab{
		rowIndex: 0,
	}
}

func (a *AboutTab) PrevRow() {
	a.rowIndex = max(0, a.rowIndex-1)
}

func (a *AboutTab) NextRow() {
	a.rowIndex++
}

func (a *AboutTab) Render(area bento.Rect, buffer *bento.Buffer) {
	var logoArea, descriptionArea bento.Rect

	bento.
		NewLayout(
			bento.ConstraintLen(44),
			bento.ConstraintMin(0),
		).
		Horizontal().
		Split(area).
		Assign(&logoArea, &descriptionArea)

	a.renderDescription(descriptionArea, buffer)

	mascotwidget.New().Render(logoArea.Inner(bento.NewPadding(0, 2)), buffer)
}

func (*AboutTab) renderDescription(area bento.Rect, buffer *bento.Buffer) {
	area = area.Inner(bento.NewPadding(4, 2))

	clearwidget.New().Render(area, buffer)

	blockwidget.New().WithStyle(theme.Global.Content).Render(area, buffer)

	area = area.Inner(bento.NewPadding(1, 2))

	const description = `- cooking up terminal user interfaces -

Bento is a Go framework that provides widgets (e.g. Paragraph, Table) and draws them to the screen efficiently every frame.`

	paragraphwidget.
		NewStr(description).
		WithStyle(theme.Global.Description).
		WithBlock(
			blockwidget.
				New().
				WithTitleStr(" Bento 🍱 ").
				WithTitlesAlignment(bento.AlignmentCenter).
				WithBorderSides(blockwidget.SideTop).
				WithBorderStyle(theme.Global.DescriptionTitle).
				WithPadding(bento.NewPadding()),
		).
		WithWrap(paragraphwidget.NewWrap().WithTrim(true)).
		WithScroll(0, 0).
		Center().
		Render(area, buffer)
}
