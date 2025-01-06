package main

import (
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct{}

// Draw implements bento.Model.
func (m *Model) Draw(frame *bento.Frame) {
	chunks := bento.Layout{
		Direction: bento.DirectionHorizontal,
		Constraints: []bento.Constraint{
			bento.ConstraintPercentage(50),
			bento.ConstraintPercentage(50),
		},
	}.Split(frame.Area())

	w := textwidget.NewTextString("Hello, world!")

	frame.RenderWidget(w.WithAlignment(bento.AlignmentLeft), chunks[0])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentRight), chunks[1])
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	return m, nil
}

func run() error {
	app, err := bento.NewApp(&Model{})
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	_, err = app.Run()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
