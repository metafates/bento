package footerwidget

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

	f.renderFooter(area, buffer, state)
}

func (f Footer) renderFooter(area bento.Rect, buffer *bento.Buffer, state *State) {
	buffer.SetStyle(area, f.style)

	footerLine := textwidget.NewLine()

	var width int

	padding := strings.Repeat(" ", f.keyPadding)

	for i, b := range state.BindingList.Items() {
		key := padding + b.Key + padding
		action := b.Action

		var spans []textwidget.Span

		if i > 0 {
			spans = append(spans, textwidget.NewSpan("  "))
		}

		spans = append(spans, textwidget.NewSpan(key).WithStyle(f.keyStyle))
		spans = append(spans, textwidget.NewSpan(" "))
		spans = append(spans, textwidget.NewSpan(action).WithStyle(f.actionStyle))

		line := textwidget.NewLine(spans...)
		lineWidth := line.Width()

		if width+lineWidth > area.Width {
			break
		}

		footerLine.Spans = append(footerLine.Spans, line.Spans...)
		width += lineWidth
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
