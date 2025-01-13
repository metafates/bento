package footerwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Widget = (*Footer)(nil)

type Footer struct {
	bindings    []Binding
	style       bento.Style
	keyStyle    bento.Style
	actionStyle bento.Style
}

func New(bindings ...Binding) Footer {
	return Footer{
		bindings:    bindings,
		style:       bento.NewStyle(),
		keyStyle:    bento.NewStyle().Yellow().Bold(),
		actionStyle: bento.NewStyle(),
	}
}

func (f Footer) Render(area bento.Rect, buffer *bento.Buffer) {
	buffer.SetStyle(area, f.style)

	footerLine := textwidget.NewLine()

	var width int

	for i, b := range f.bindings {
		key := b.Key
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
