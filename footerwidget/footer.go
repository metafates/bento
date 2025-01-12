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

	left, right           *textwidget.Line
	leftStyle, rightStyle bento.Style
}

func New(bindings ...Binding) Footer {
	return Footer{
		bindings:    bindings,
		style:       bento.NewStyle(),
		keyStyle:    bento.NewStyle().Yellow().Bold(),
		actionStyle: bento.NewStyle(),
		left:        nil,
		right:       nil,
		leftStyle:   bento.NewStyle().OnRGB(41, 41, 41),
		rightStyle:  bento.NewStyle().OnRGB(41, 41, 41),
	}
}

func (f Footer) WithRightLine(line textwidget.Line) Footer {
	f.right = &line
	return f
}

func (f Footer) WithLeftLine(line textwidget.Line) Footer {
	f.left = &line
	return f
}

func (f Footer) Render(area bento.Rect, buffer *bento.Buffer) {
	buffer.SetStyle(area, f.style)

	if f.left != nil {
		var leftArea bento.Rect

		bento.
			NewLayout(
				bento.ConstraintLength(f.left.Width()+2),
				bento.ConstraintFill(1),
			).
			Horizontal().
			Split(area).
			Assign(&leftArea, &area)

		f.left.Center().WithStyle(f.leftStyle).Render(leftArea, buffer)
	}

	if f.right != nil {
		var rightArea bento.Rect

		bento.
			NewLayout(
				bento.ConstraintFill(1),
				bento.ConstraintLength(f.right.Width()+2),
			).
			Horizontal().
			Split(area).
			Assign(&area, &rightArea)

		f.right.Center().WithStyle(f.rightStyle).Render(rightArea, buffer)
	}

	footerLine := textwidget.NewLine()

	var width int

	for i, b := range f.bindings {
		key := b.Key
		action := b.Action

		var spans []textwidget.Span

		if i == 0 {
			spans = append(spans, textwidget.NewSpan("  "))
		}

		spans = append(spans, textwidget.NewSpan(key).WithStyle(f.keyStyle))
		spans = append(spans, textwidget.NewSpan(" "))
		spans = append(spans, textwidget.NewSpan(action).WithStyle(f.actionStyle))
		spans = append(spans, textwidget.NewSpan("  "))

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
