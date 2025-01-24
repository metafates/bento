package tabs

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/clearwidget"
	"github.com/metafates/bento/examples/unnamed/gradient"
	"github.com/metafates/bento/examples/unnamed/theme"
	"github.com/metafates/bento/paragraphwidget"
)

var _ bento.Widget = (*About)(nil)

type About struct {
	rowIndex int
}

func NewAbout() About {
	return About{
		rowIndex: 0,
	}
}

func (a *About) PrevRow() {
	a.rowIndex = max(0, a.rowIndex-1)
}

func (a *About) NextRow() {
	a.rowIndex++
}

func (a *About) Render(area bento.Rect, buffer *bento.Buffer) {
	var logoArea, descriptionArea bento.Rect

	bento.
		NewLayout(
			bento.ConstraintLen(34),
			bento.ConstraintMin(0),
		).
		Horizontal().
		Split(area).
		Assign(&logoArea, &descriptionArea)

	a.renderDescription(area, buffer)
}

func (*About) renderDescription(area bento.Rect, buffer *bento.Buffer) {
	gradient.New().Render(area, buffer)

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
				WithTitleStr(" Bento ").
				WithTitlesAlignment(bento.AlignmentCenter).
				WithBorderSides(blockwidget.SideTop).
				WithBorderStyle(theme.Global.DescriptionTitle).
				WithPadding(bento.NewPadding()),
		).
		WithWrap(paragraphwidget.NewWrap().Trimmed()).
		WithScroll(0, 0).
		Render(area, buffer)
}
