package main

import (
	"context"
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
		Direction: bento.DirectionVertical,
		Constraints: []bento.Constraint{
			bento.ConstraintPercentage(25),
			bento.ConstraintPercentage(25),
			bento.ConstraintPercentage(25),
			bento.ConstraintPercentage(25),
		},
	}.Split(frame.Area())

	w := textwidget.NewTextString("Hello, world!")

	frame.RenderWidget(w.WithAlignment(bento.AlignmentLeft), chunks[0])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentCenter), chunks[1])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentRight), chunks[2])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentLeft), chunks[3])
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case bento.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, bento.Quit
		}
	}

	return m, nil
}

func run() error {
	app, err := bento.NewApp(context.Background(), &Model{})
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
