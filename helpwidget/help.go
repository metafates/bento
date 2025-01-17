package helpwidget

import (
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/filterablelistwidget"
	"github.com/metafates/bento/listwidget"
	"github.com/metafates/bento/popupwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.StatefulWidget[*State] = (*Footer)(nil)

type Footer struct {
	style       bento.Style
	keyStyle    bento.Style
	actionStyle bento.Style
	keyPadding  int
}

func New() Footer {
	return Footer{
		style:       bento.NewStyle(),
		keyStyle:    bento.NewStyle().Reversed(),
		actionStyle: bento.NewStyle(),
		keyPadding:  1,
	}
}

// WithKeyPadding sets key padding to use on both sides.
func (f Footer) WithKeyPadding(padding int) Footer {
	f.keyPadding = padding
	return f
}

func (f Footer) RenderStateful(area bento.Rect, buffer *bento.Buffer, state *State) {
	if state.ShowPopup {
		f.renderPopup(buffer.Area(), buffer, state)
	}

	helpBinding := f.bindingLine(_helpBinding)

	var otherKeysArea, helpKeyArea bento.Rect

	bento.NewLayout(
		bento.ConstraintFill(1),
		bento.ConstraintLength(1),
		bento.ConstraintLength(helpBinding.Width()),
	).Horizontal().Split(area).Assign(&otherKeysArea, nil, &helpKeyArea)

	f.renderFooter(otherKeysArea, buffer, state)

	helpBinding.Render(helpKeyArea, buffer)
}

func (f Footer) bindingLine(b Binding) textwidget.Line {
	padding := strings.Repeat(" ", f.keyPadding)
	key := padding + b.String() + padding

	return textwidget.NewLine(
		textwidget.NewSpan(key).WithStyle(f.keyStyle),
		textwidget.NewSpan(" "),
		textwidget.NewSpan(b.Name).WithStyle(f.actionStyle),
	)
}

func (f Footer) renderFooter(area bento.Rect, buffer *bento.Buffer, state *State) {
	buffer.SetStyle(area, f.style)

	footerLine := textwidget.NewLine()

	var width, shownCount int

	for _, b := range state.BindingList.AllItems() {
		if b.IsHidden {
			continue
		}

		var spans []textwidget.Span

		if shownCount > 0 {
			spans = append(spans, textwidget.NewSpan("  "))
		}

		spans = append(spans, f.bindingLine(b).Spans...)

		line := textwidget.NewLine(spans...)
		lineWidth := line.Width()

		if width+lineWidth > area.Width {
			break
		}

		footerLine.Spans = append(footerLine.Spans, line.Spans...)
		width += lineWidth
		shownCount++
	}

	footerLine.Render(area, buffer)
}

func (f Footer) renderPopup(area bento.Rect, buffer *bento.Buffer, state *State) {
	block := blockwidget.New().Bordered().WithTitleStr("Help")
	list := listwidget.
		New().
		WithHighlightStyle(bento.NewStyle().Reversed()).
		WithBlock(block)

	popup := popupwidget.
		NewStateful(filterablelistwidget.New[Binding](list)).
		Center().
		Middle().
		WithHeight(bento.ConstraintPercentage(60))

	popup.RenderStateful(area, buffer, &state.BindingList)
}
