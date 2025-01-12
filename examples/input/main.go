package main

import (
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

func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	block := blockwidget.New().Bordered().Thick().WithTitleStr("Input")
	input := inputwidget.NewInput().WithPlaceholder("Placeholder...").WithPrompt("> ")

	popup := popupwidget.NewStateful(input).WithBlock(block).WithHeight(bento.ConstraintLength(3))

	fill := fillwidget.New("╲").WithStyle(bento.NewStyle().Dim())

	fill.Render(area, buffer)
	popup.RenderStateful(area, buffer, &m.input)
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

	_, err := bento.NewApp(&model).Run()
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
