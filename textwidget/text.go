package textwidget

import (
	"github.com/metafates/bento"
)

var _ bento.Widget = (*Text)(nil)

type Text struct {
	Style     bento.Style
	Lines     Lines
	Alignment bento.Alignment
}

func NewText(lines Lines) *Text {
	return &Text{
		Style:     bento.NewStyle(),
		Lines:     lines,
		Alignment: bento.AlignmentNone,
	}
}

// Render implements bento.Widget.
func (t *Text) Render(area bento.Rect, buffer *bento.Buffer) {
	area = area.Intersection(buffer.Area)

	buffer.SetStyle(area, t.Style)

	rows := area.Rows()

	for i := 0; i < min(len(t.Lines), len(rows)); i++ {
		line := t.Lines[i]
		lineArea := rows[i]

		line.renderWithAlignment(lineArea, buffer, t.Alignment)
	}
}
