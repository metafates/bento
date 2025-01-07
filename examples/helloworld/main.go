package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	count int
}

// Draw implements bento.Model.
func (m *Model) Draw(frame *bento.Frame) {
	mainBlock := blockwidget.
		NewBlock().
		WithTitle(blockwidget.NewTitleString("Example")).
		WithPadding(blockwidget.NewPadding(3))

	statusBlock := blockwidget.NewBlock().WithTitle(blockwidget.NewTitleString("Status"))

	mainArea, statusArea := bento.Layout{
		Direction: bento.DirectionVertical,
		Constraints: []bento.Constraint{
			bento.ConstraintMin(3),
			bento.ConstraintLength(3),
		},
	}.Split2(frame.Area())

	mainInnerArea := mainBlock.Inner(mainArea)

	textChunks := bento.Layout{
		Direction: bento.DirectionVertical,
		Constraints: []bento.Constraint{
			bento.ConstraintPercentage(25),
			bento.ConstraintPercentage(25),
			bento.ConstraintPercentage(25),
			bento.ConstraintPercentage(25),
		},
	}.Split(mainInnerArea)

	style := bento.
		NewStyle().
		WithModifier(bento.StyleModifierItalic)
		// WithForeground(termenv.ANSIBlue)

	span := textwidget.NewSpan("Hello, World! " + strconv.Itoa(m.count)).WithStyle(style)

	w := textwidget.NewText(textwidget.NewLine(span))

	frame.RenderWidget(mainBlock, mainArea)
	frame.RenderWidget(w.WithAlignment(bento.AlignmentLeft), textChunks[0])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentCenter), textChunks[1])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentRight), textChunks[2])
	frame.RenderWidget(w.WithAlignment(bento.AlignmentLeft), textChunks[3])
	frame.RenderWidget(statusBlock, statusArea)
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
		case "a":
			m.count++
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
