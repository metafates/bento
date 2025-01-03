package textwidget

import (
	"bufio"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/bento"
)

var _ bento.Widget = (*Text)(nil)

type Text struct {
	Style     lipgloss.Style
	Lines     Lines
	Alignment bento.Alignment
}

func NewText(s string) *Text {
	var lines []Line

	if s == "" {
		lines = []Line{*NewLine("")}
	} else {
		scanner := bufio.NewScanner(strings.NewReader(s))

		for scanner.Scan() {
			line := NewLine(scanner.Text())

			lines = append(lines, *line)
		}
	}

	return &Text{
		Style:     lipgloss.NewStyle(),
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
