package reflow

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/internal/grapheme"
	"github.com/metafates/bento/textwidget"
	"github.com/rivo/uniseg"
)

var _ LineComposer = (*LineTruncator)(nil)

type LineTruncator struct {
	inputLines       []InputLine
	maxLineWidth     int
	currentLine      []textwidget.StyledGrapheme
	horizontalOffset int
}

func NewLineTruncator(lines []InputLine, maxLineWidth int) LineTruncator {
	return LineTruncator{
		inputLines:       lines,
		maxLineWidth:     maxLineWidth,
		currentLine:      nil,
		horizontalOffset: 0,
	}
}

func (lt *LineTruncator) SetHorizontalOffset(horizontalOffset int) {
	lt.horizontalOffset = horizontalOffset
}

// NextLine implements LineComposer.
func (lt *LineTruncator) NextLine() (WrappedLine, bool) {
	if lt.maxLineWidth == 0 {
		return WrappedLine{}, false
	}

	clear(lt.currentLine)

	var currentLineWidth int

	linesExhausted := true
	horizontalOffset := lt.horizontalOffset
	currentAlignment := bento.AlignmentLeft

	var lastIndex int
	for i, line := range lt.inputLines {
		lastIndex = i
		currentLine := line.Graphemes
		currentAlignment = line.Alignment
		alignment := line.Alignment
		linesExhausted = false

		for _, g := range currentLine {
			if g.Width() > lt.maxLineWidth {
				continue
			}

			if currentLineWidth+g.Width() > lt.maxLineWidth {
				break
			}

			var symbol string

			if horizontalOffset == 0 || alignment != bento.AlignmentLeft {
				symbol = g.String()
			} else {
				width := g.Width()

				if width > horizontalOffset {
					symbol = trimOffset(g.String(), horizontalOffset)
					horizontalOffset = 0
				} else {
					horizontalOffset -= width
				}
			}

			currentLineWidth += uniseg.StringWidth(symbol)
			lt.currentLine = append(lt.currentLine, textwidget.StyledGrapheme{
				Style:    g.Style,
				Grapheme: grapheme.New(symbol),
			})
		}
	}

	lt.inputLines = lt.inputLines[min(len(lt.inputLines), lastIndex+1):]

	if linesExhausted {
		return WrappedLine{}, false
	}

	return WrappedLine{
		Line:      lt.currentLine,
		Width:     currentLineWidth,
		Alignment: currentAlignment,
	}, true
}

func trimOffset(src string, offset int) string {
	var start int

	graphemes := uniseg.NewGraphemes(src)

	for graphemes.Next() {
		width := graphemes.Width()

		if width > offset {
			break
		}

		offset -= width
		start += len(graphemes.Str())
	}

	return src[start:]
}
