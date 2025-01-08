package textwidget

import (
	"bufio"
	"strings"

	"github.com/metafates/bento"
)

var _ bento.Widget = (*Line)(nil)

type Lines []Line

func NewLines(lines ...Line) Lines {
	return lines
}

func NewLinesString(s ...string) Lines {
	var lines []Line

	joined := strings.Join(s, "\n")

	if joined == "" {
		lines = []Line{NewLineString("")}
	} else {
		scanner := bufio.NewScanner(strings.NewReader(joined))

		for scanner.Scan() {
			line := NewLineString(scanner.Text())

			lines = append(lines, line)
		}
	}

	return Lines(lines)
}

func (l Lines) NewBuffer() *bento.Buffer {
	height := len(l)
	width := l.Width()

	buffer := bento.NewBufferEmpty(bento.Rect{
		Width:  width,
		Height: height,
	})

	for y, line := range l {
		setLine(0, y, buffer, &line, width)
	}

	return buffer
}

func (l Lines) Width() int {
	var width int

	for _, line := range l {
		width = max(width, line.Width())
	}

	return width
}

type Line struct {
	Style     bento.Style
	Spans     []Span
	Alignment bento.Alignment
}

func NewLine(spans ...Span) Line {
	return Line{
		Style:     bento.NewStyle(),
		Spans:     spans,
		Alignment: bento.AlignmentNone,
	}
}

func NewLineString(s string) Line {
	lines := bufio.NewScanner(strings.NewReader(s))

	var spans []Span

	for lines.Scan() {
		line := lines.Text()

		spans = append(spans, NewSpan(line))
	}

	return NewLine(spans...)
}

func (l Line) StyledGraphemes(baseStyle bento.Style) []StyledGrapheme {
	style := baseStyle.Patched(l.Style)

	graphemes := make([]StyledGrapheme, 0, len(l.Spans))

	for _, s := range l.Spans {
		graphemes = append(graphemes, s.StyledGraphemes(style)...)
	}

	return graphemes
}

func (l Line) WithStyle(style bento.Style) Line {
	l.Style = style
	return l
}

func (l Line) WithAlignment(alignment bento.Alignment) Line {
	l.Alignment = alignment
	return l
}

func (l Line) Render(area bento.Rect, buffer *bento.Buffer) {
	l.renderWithAlignment(area, buffer, bento.AlignmentNone)
}

func (l Line) renderWithAlignment(
	area bento.Rect,
	buffer *bento.Buffer,
	parentAlignment bento.Alignment,
) {
	area = area.Intersection(buffer.Area)
	if area.IsEmpty() {
		return
	}

	area = bento.Rect{
		X:      area.X,
		Y:      area.Y,
		Width:  area.Width,
		Height: 1,
	}

	lineWidth := l.Width()
	if lineWidth == 0 {
		return
	}

	buffer.SetStyle(area, l.Style)

	alignment := l.Alignment
	if alignment == bento.AlignmentNone {
		alignment = parentAlignment
	}

	areaWidth := area.Width

	canRenderCompleteLine := lineWidth <= areaWidth

	if canRenderCompleteLine {
		var indentWidth int

		switch alignment {
		case bento.AlignmentLeft, bento.AlignmentNone:
			indentWidth = 0
		case bento.AlignmentCenter:
			indentWidth = max(0, areaWidth-lineWidth) / 2
		case bento.AlignmentRight:
			indentWidth = max(0, areaWidth-lineWidth)
		}

		area = area.IndentX(indentWidth)

		renderSpans(l.Spans, area, buffer, 0)
	} else {
		// There is not enough space to render the whole line. As the right side is truncated by
		// the area width, only truncate the left.
		var skipWidth int

		switch alignment {
		case bento.AlignmentLeft, bento.AlignmentNone:
			skipWidth = 0
		case bento.AlignmentCenter:
			skipWidth = max(0, lineWidth-areaWidth) / 2
		case bento.AlignmentRight:
			skipWidth = max(0, lineWidth-areaWidth)
		}

		renderSpans(l.Spans, area, buffer, skipWidth)
	}
}

func (l Line) Width() int {
	var width int

	for _, s := range l.Spans {
		width += s.Width()
	}

	return width
}

func setLine(x, y int, buffer *bento.Buffer, line *Line, maxWidth int) (int, int) {
	remainingWidth := maxWidth

	for _, s := range line.Spans {
		if remainingWidth == 0 {
			break
		}

		style := line.Style.Patched(s.Style)
		newX, _ := buffer.SetStringN(x, y, s.Content, remainingWidth, style)

		w := max(0, newX-x)
		x = newX

		remainingWidth = max(0, remainingWidth, w)
	}

	return x, y
}
