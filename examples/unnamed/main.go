package main

import (
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/gaugewidget"
	"github.com/metafates/bento/helpwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	ratio float64

	footerState *helpwidget.State
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Render implements bento.Model.
func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	var mainArea, footerArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		Fill(1).
		Len(1).
		Split(area).
		Assign(&mainArea, &footerArea)

	gauge := gaugewidget.New().WithRatio(m.ratio).WithUnicode(true)

	gauge.Render(mainArea, buffer)

	helpwidget.
		New().
		RenderStateful(footerArea, buffer, m.footerState)
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	consumed, cmd := m.footerState.TryUpdate(msg)
	if consumed {
		return m, cmd
	}

	switch msg := msg.(type) {
	case ChangeMsg:
		m.ratio = max(0, min(1, m.ratio+float64(msg)))
	}

	return m, nil
}

type ChangeMsg float64

func Change(delta float64) bento.Cmd {
	return func() bento.Msg {
		return ChangeMsg(delta)
	}
}

func run() error {
	_, err := bento.
		NewApp(newModel()).
		Run()
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	return nil
}

func newModel() *Model {
	return &Model{
		ratio: 0,
		footerState: helpwidget.NewState(
			helpwidget.NewBinding("quit", "ctrl+c").
				WithAction(func() bento.Cmd { return bento.Quit }).
				Hidden(),

			helpwidget.NewBinding("increment", "up").
				WithAliases("right", "l", "+").
				WithDescription("Increment the gauge").
				WithAction(func() bento.Cmd {
					return Change(0.01)
				}),

			helpwidget.NewBinding("decrement", "down").
				WithAliases("left", "h", "-").
				WithDescription("Decrement the gauge").
				WithAction(func() bento.Cmd {
					return Change(-0.01)
				}),
		),
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
