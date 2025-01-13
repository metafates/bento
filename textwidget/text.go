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

func NewTextStr(s string) Text {
	lines := NewLinesStr(s)

	return NewText(lines...)
}

func NewText(lines ...Line) Text {
	return Text{
		Style:     bento.NewStyle(),
		Lines:     lines,
		Alignment: bento.AlignmentNone,
	}
}

func (t Text) Height() int {
	return len(t.Lines)
}

func (t Text) Width() int {
	return t.Lines.Width()
}

func (t Text) WithStyle(style bento.Style) Text {
	t.Style = style
	return t
}

func (t Text) WithAlignment(alignment bento.Alignment) Text {
	t.Alignment = alignment
	return t
}

func (t Text) Left() Text {
	return t.WithAlignment(bento.AlignmentLeft)
}

func (t Text) Right() Text {
	return t.WithAlignment(bento.AlignmentRight)
}

func (t Text) Center() Text {
	return t.WithAlignment(bento.AlignmentCenter)
}

// Render implements bento.Widget.
func (t Text) Render(area bento.Rect, buffer *bento.Buffer) {
	area = area.Intersection(buffer.Area())

	buffer.SetStyle(area, t.Style)

	rows := area.Rows()

	for i := 0; i < min(len(t.Lines), len(rows)); i++ {
		line := t.Lines[i]
		lineArea := rows[i]

		line.renderWithAlignment(lineArea, buffer, t.Alignment)
	}
}

func AppendTextSpans(text *Text, spans ...Span) {
	if len(text.Lines) > 0 {
		last := len(text.Lines) - 1
		text.Lines[last].Spans = append(text.Lines[last].Spans, spans...)
	} else {
		text.Lines = append(text.Lines, NewLine(spans...))
	}
}
