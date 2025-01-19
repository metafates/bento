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

var _ bento.StatefulWidget[*State] = (*Help)(nil)

type Help struct {
	style       bento.Style
	keyStyle    bento.Style
	actionStyle bento.Style
	keyPadding  int
}

func New() Help {
	return Help{
		style:       bento.NewStyle(),
		keyStyle:    bento.NewStyle().Reversed(),
		actionStyle: bento.NewStyle(),
		keyPadding:  1,
	}
}

// WithKeyPadding sets key padding to use on both sides.
func (h Help) WithKeyPadding(padding int) Help {
	h.keyPadding = padding
	return h
}

func (h Help) RenderStateful(area bento.Rect, buffer *bento.Buffer, state *State) {
	if state.showPopup {
		h.renderPopup(buffer.Area(), buffer, state)
	}

	helpBinding := h.bindingLine(_helpBinding)

	var otherKeysArea, helpKeyArea bento.Rect

	bento.NewLayout(
		bento.ConstraintFill(1),
		bento.ConstraintLen(1),
		bento.ConstraintLen(helpBinding.Width()),
	).Horizontal().Split(area).Assign(&otherKeysArea, nil, &helpKeyArea)

	h.renderFooter(otherKeysArea, buffer, state)

	helpBinding.Render(helpKeyArea, buffer)
}

func (h Help) bindingLine(b Binding) textwidget.Line {
	padding := strings.Repeat(" ", h.keyPadding)
	key := padding + b.String() + padding

	return textwidget.NewLine(
		textwidget.NewSpan(key).WithStyle(h.keyStyle),
		textwidget.NewSpan(" "),
		textwidget.NewSpan(b.Name).WithStyle(h.actionStyle),
	)
}

func (h Help) renderFooter(area bento.Rect, buffer *bento.Buffer, state *State) {
	buffer.SetStyle(area, h.style)

	footerLine := textwidget.NewLine()

	var width, shownCount int

	for _, b := range state.bindingList.AllItems() {
		if b.IsHidden || !b.IsActive() {
			continue
		}

		var spans []textwidget.Span

		if shownCount > 0 {
			spans = append(spans, textwidget.NewSpan("  "))
		}

		spans = append(spans, h.bindingLine(b).Spans...)

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

func (h Help) renderPopup(area bento.Rect, buffer *bento.Buffer, state *State) {
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

	popup.RenderStateful(area, buffer, &state.bindingList)
}
