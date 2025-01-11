package main

import (
	"context"
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/fillwidget"
	"github.com/metafates/bento/inputwidget"
	"github.com/metafates/bento/popupwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	input inputwidget.State
}

// Draw implements bento.Model.
func (m *Model) Draw(frame *bento.Frame) {
	block := blockwidget.NewBlock().Bordered().Thick().WithTitleStr("Input")
	input := inputwidget.NewInput().WithPlaceholder("Placeholder...").WithPrompt("> ")

	popup := popupwidget.NewStateful(input).WithBlock(block).WithHeight(bento.ConstraintLength(3))

	fill := fillwidget.New("⣿").WithStyle(bento.NewStyle().Dim())
	frame.RenderWidget(fill, frame.Area())
	bento.RenderStatefulWidget(frame, popup, frame.Area(), m.input)
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
		default:
			m.input.HandleKey(bento.Key(msg))
		}
	}

	return m, nil
}

func run() error {
	model := Model{input: inputwidget.NewState()}

	app, err := bento.NewApp(context.Background(), &model)
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	_, err = app.Run()
	if err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
