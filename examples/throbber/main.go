package main

import (
	"fmt"
	"log"
	"time"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/textwidget"
	"github.com/metafates/bento/throbberwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	throbber throbberwidget.State
}

func (m *Model) Init() bento.Cmd {
	return m.throbber.Tick
}

// Render implements bento.Model.
func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	var infoArea, throbbersArea, footerArea bento.Rect

	bento.
		NewLayout(
			bento.ConstraintLength(1),
			bento.ConstraintFill(1),
			bento.ConstraintLength(1),
		).
		Vertical().
		Split(area).
		Assign(&infoArea, &throbbersArea, &footerArea)

	m.renderInfo(infoArea, buffer)
	m.renderThrobbers(throbbersArea, buffer)
	m.renderFooter(footerArea, buffer)
}

func (m *Model) renderFooter(area bento.Rect, buffer *bento.Buffer) {
	text := textwidget.NewTextStr(`"+" - increment FPS | "-" - decrement FPS`)

	text.Render(area, buffer)
}

func (m *Model) renderInfo(area bento.Rect, buffer *bento.Buffer) {
	text := textwidget.NewTextStr(fmt.Sprintf("FPS %s", m.throbber.FPS))

	text.Render(area, buffer)
}

func (m *Model) renderThrobbers(area bento.Rect, buffer *bento.Buffer) {
	types := []throbberwidget.Type{
		throbberwidget.TypeLine,
		throbberwidget.TypeHorizontalBlock,
		throbberwidget.TypeVerticalBlock,
		throbberwidget.TypeBrailleEight,
		throbberwidget.TypeParenthesis,
		throbberwidget.TypeCanadian,
		throbberwidget.TypeEllipsis,
		throbberwidget.TypeMeter,
		throbberwidget.TypePulse,
	}

	grid := getGrid(area, len(types))

	for i, t := range types {
		block := blockwidget.NewBlock().Bordered().WithTitleStr(t.String())

		a := block.Inner(grid[i])

		block.Render(grid[i], buffer)

		throbberwidget.
			New().
			WithType(t).
			RenderStateful(a, buffer, m.throbber)
	}
}

func getGrid(area bento.Rect, itemsCount int) []bento.Rect {
	const columns = 3
	rows := itemsCount/columns + itemsCount%columns

	rowsConstraints := make([]bento.Constraint, 0, rows)

	for i := 0; i < rows; i++ {
		rowsConstraints = append(rowsConstraints, bento.ConstraintFill(1))
	}

	rowsAreas := bento.NewLayout(rowsConstraints...).Vertical().Split(area)

	areas := make([]bento.Rect, 0, rows+columns)

	for _, ra := range rowsAreas {
		constraints := make([]bento.Constraint, 0, columns)

		for i := 0; i < columns; i++ {
			constraints = append(constraints, bento.ConstraintFill(1))
		}

		columnsAreas := bento.NewLayout(constraints...).Horizontal().Split(ra)

		for _, ca := range columnsAreas {
			areas = append(areas, ca)
		}
	}

	return areas
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case throbberwidget.TickMsg:
		return m, m.throbber.Update(msg)
	case bento.KeyMsg:
		switch msg.String() {
		case "+":
			m.throbber.FPS += time.Millisecond
		case "-":
			m.throbber.FPS -= time.Millisecond
		case "ctrl+c":
			return m, bento.Quit
		}
	}

	return m, nil
}

func run() error {
	model := Model{
		throbber: throbberwidget.NewState(),
	}

	_, err := bento.NewApp(&model).Run()
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
