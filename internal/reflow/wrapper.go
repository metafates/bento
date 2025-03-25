package reflow

import (
	"slices"

	"github.com/edwingeng/deque/v2"
	"github.com/metafates/bento"
	"github.com/metafates/bento/widget/textwidget"
)

var _ LineComposer = (*WordWrapper)(nil)

type WordWrapper struct {
	inputLines       []InputLine
	maxLineWidth     int
	wrappedLines     *deque.Deque[[]textwidget.StyledGrapheme]
	currentAlignment bento.Alignment
	currentLine      []textwidget.StyledGrapheme
	trim             bool

	pendingWord       []textwidget.StyledGrapheme
	pendingWhitespace *deque.Deque[textwidget.StyledGrapheme]
	pendingLinePool   [][]textwidget.StyledGrapheme
}

func NewWordWrapper(lines []InputLine, maxLineWidth int, trim bool) WordWrapper {
	return WordWrapper{
		inputLines:        lines,
		maxLineWidth:      maxLineWidth,
		wrappedLines:      deque.NewDeque[[]textwidget.StyledGrapheme](),
		currentAlignment:  bento.AlignmentLeft,
		currentLine:       nil,
		trim:              trim,
		pendingWord:       nil,
		pendingWhitespace: deque.NewDeque[textwidget.StyledGrapheme](),
		pendingLinePool:   nil,
	}
}

func (ww *WordWrapper) processInput(lineSymbols []textwidget.StyledGrapheme) {
	var pendingLine []textwidget.StyledGrapheme

	if len(ww.pendingLinePool) > 0 {
		pendingLine = ww.pendingLinePool[len(ww.pendingLinePool)-1]
		ww.pendingLinePool = slices.Delete(ww.pendingLinePool, len(ww.pendingLinePool)-1, len(ww.pendingLinePool))
	}

	var (
		lineWidth             int
		wordWidth             int
		whitespaceWidth       int
		nonWhitespacePrevious bool
	)

	clear(ww.pendingWord)
	ww.pendingWhitespace.Clear()
	clear(pendingLine)

	for _, g := range lineSymbols {
		isWhitespace := g.IsWhitespace()
		symbolWidth := g.Width()

		if symbolWidth > ww.maxLineWidth {
			continue
		}

		wordFound := nonWhitespacePrevious && isWhitespace

		trimmedOverflow := len(pendingLine) == 0 && ww.trim && wordWidth+symbolWidth > ww.maxLineWidth
		whitespaceOverflow := len(pendingLine) == 0 && ww.trim && whitespaceWidth+symbolWidth > ww.maxLineWidth
		untrimmedOverflow := len(pendingLine) == 0 && !ww.trim && wordWidth+whitespaceWidth+symbolWidth > ww.maxLineWidth

		if wordFound || trimmedOverflow || whitespaceOverflow || untrimmedOverflow {
			if len(pendingLine) > 0 || !ww.trim {
				pendingLine = append(pendingLine, ww.pendingWhitespace.Dump()...)
				lineWidth += whitespaceWidth
				ww.pendingWhitespace.Clear()
			}

			pendingLine = append(pendingLine, ww.pendingWord...)
			clear(ww.pendingWord)

			lineWidth += wordWidth

			ww.pendingWhitespace.Clear()
			whitespaceWidth = 0
			wordWidth = 0
		}

		lineFull := lineWidth >= ww.maxLineWidth
		pendingWordOverflow := symbolWidth > 0 && lineWidth+whitespaceWidth+wordWidth >= ww.maxLineWidth

		if lineFull || pendingWordOverflow {
			remainingWidth := max(0, ww.maxLineWidth-lineWidth)

			ww.wrappedLines.PushBack(slices.Clone(pendingLine))
			clear(pendingLine)

			lineWidth = 0

			for {
				grapheme, ok := ww.pendingWhitespace.Front()
				if !ok {
					break
				}

				width := grapheme.Width()

				if width > remainingWidth {
					break
				}

				whitespaceWidth -= width
				remainingWidth -= width

				ww.pendingWhitespace.PopFront()
			}

			if isWhitespace && ww.pendingWhitespace.IsEmpty() {
				continue
			}
		}

		if isWhitespace {
			whitespaceWidth += symbolWidth
			ww.pendingWhitespace.PushBack(g)
		} else {
			wordWidth += symbolWidth
			ww.pendingWord = append(ww.pendingWord, g)
		}

		nonWhitespacePrevious = !isWhitespace
	}

	if len(pendingLine) == 0 && len(ww.pendingWord) == 0 && !ww.pendingWhitespace.IsEmpty() {
		ww.wrappedLines.PushBack(nil)
	}

	if len(pendingLine) > 0 || !ww.trim {
		pendingLine = append(pendingLine, ww.pendingWhitespace.Dump()...)
		ww.pendingWhitespace.Clear()
	}

	pendingLine = append(pendingLine, ww.pendingWord...)
	clear(ww.pendingWord)

	if len(pendingLine) > 0 {
		ww.wrappedLines.PushBack(pendingLine)
	} else if cap(pendingLine) > 0 {
		ww.pendingLinePool = append(ww.pendingLinePool, pendingLine)
	}

	if ww.wrappedLines.IsEmpty() {
		ww.wrappedLines.PushBack(nil)
	}
}

func (ww *WordWrapper) replaceCurrentLine(line []textwidget.StyledGrapheme) {
	cache := ww.currentLine
	ww.currentLine = line

	if cap(cache) > 0 {
		ww.pendingLinePool = append(ww.pendingLinePool, cache)
	}
}

// NextLine implements LineComposer.
func (ww *WordWrapper) NextLine() (WrappedLine, bool) {
	if ww.maxLineWidth == 0 {
		return WrappedLine{}, false
	}

	for {
		line, ok := ww.wrappedLines.TryPopFront()
		if ok {
			var lineWidth int

			for _, g := range line {
				lineWidth += g.Width()
			}

			ww.replaceCurrentLine(line)

			return WrappedLine{
				Line:      ww.currentLine,
				Width:     lineWidth,
				Alignment: ww.currentAlignment,
			}, true
		}

		if len(ww.inputLines) == 0 {
			return WrappedLine{}, false
		}

		// TODO: use iterators
		inputLine := ww.inputLines[0]
		ww.inputLines = ww.inputLines[1:]

		ww.currentAlignment = inputLine.Alignment
		ww.processInput(inputLine.Graphemes)
	}
}
