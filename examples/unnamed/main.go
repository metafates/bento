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

	footerState *footerwidget.State
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
		RenderStateful(footerArea, buffer, m.footerState)
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case ChangeMsg:
		m.ratio = max(0, min(1, m.ratio+float64(msg)))
	case bento.KeyMsg:
		m.footerState.Update(bento.Key(msg))
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
	_, err := bento.NewAppWithProxy(func(proxy bento.AppProxy) bento.Model {
		return &Model{
			ratio: 0,
			footerState: footerwidget.NewState(
				footerwidget.NewBinding("quit", "ctrl+c").
					WithAction(proxy.Quit).
					Hidden(),

				footerwidget.NewBinding("increment", "up", "right", "l", "+").
					WithDescription("Increment the gauge").
					WithAction(func() {
						proxy.Send(Change(0.01))
					}),

				footerwidget.NewBinding("decrement", "down", "left", "h", "-").
					WithDescription("Decrement the gauge").
					WithAction(func() {
						proxy.Send(Change(-0.01))
					}),
			),
		}
	}).Run()
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
