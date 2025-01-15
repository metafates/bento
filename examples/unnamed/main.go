package main

import (
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/footerwidget"
	"github.com/metafates/bento/gaugewidget"
	"github.com/metafates/bento/popupwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	ratio float64

	footerState footerwidget.State
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Render implements bento.Model.
func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	var mainArea, footerArea bento.Rect

	bento.
		NewLayout(
			bento.ConstraintFill(1),
			bento.ConstraintLength(1),
		).
		Vertical().
		Split(area).
		Assign(&mainArea, &footerArea)

	message := textwidget.NewTextStr("Try scrolling")
	popup := popupwidget.
		New(message).
		Top().
		Right().
		WithWidth(bento.ConstraintLength(message.Width() + 2)).
		WithHeight(bento.ConstraintLength(message.Height() + 2)).
		WithBlock(blockwidget.New().Bordered().WithTitleStr("Hint"))

	gauge := gaugewidget.New().WithRatio(m.ratio).WithUnicode(true)

	gauge.Render(mainArea, buffer)
	popup.Render(area, buffer)

	footerwidget.
		New().
		RenderStateful(footerArea, buffer, &m.footerState)
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case bento.KeyMsg:
		footerUpdated := m.footerState.Update(bento.Key(msg))
		if footerUpdated {
			return m, nil
		}

		switch msg.String() {
		case "shift+up", "shift+right", "L", "+":
			m.ratio = min(1, m.ratio+0.02)
		case "up", "right", "l":
			m.ratio = min(1, m.ratio+0.001)
		case "shift+down", "shift+left", "H", "-":
			m.ratio = max(0, m.ratio-0.02)
		case "down", "left", "h":
			m.ratio = max(0, m.ratio-0.001)
		case "ctrl+c":
			return m, bento.Quit
		}
	}

	return m, nil
}

func run() error {
	model := Model{
		ratio: 0,
		footerState: footerwidget.NewState(
			footerwidget.NewBinding("^c", "quit"),
			footerwidget.NewBinding("up", "increment").WithDescription("Increment the gauge"),
			footerwidget.NewBinding("down", "decrement").WithDescription("Decrement the gauge"),
		),
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
