package main

import (
	"context"
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/listwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	listState listwidget.State
}

func (m *Model) Draw(frame *bento.Frame) {
	block := blockwidget.NewBlock().WithBorders().Rounded().WithTitle(blockwidget.NewTitleString("List"))

	var items []listwidget.Item

	for i := 0; i < 100; i++ {
		items = append(items, listwidget.NewItemString(fmt.Sprintf("Item #%d", i)))
	}

	list := listwidget.
		NewList(items...).
		WithHighlightSymbol("> ").
		WithHighlightSpacing(listwidget.HighlightSpacingAlways).
		WithBlock(block).WithHighlightStyle(bento.NewStyle().Black().OnBlue())

	bento.RenderStatefulWidget(frame, list, frame.Area(), &m.listState)
}

func (m *Model) Init() bento.Cmd {
	return nil
}

func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case bento.KeyMsg:
		switch msg.String() {
		case "ctrl+u":
			m.listState.ScrollUpBy(8)
		case "ctrl+d":
			m.listState.ScrollDownBy(8)
		case "j", "down":
			m.listState.SelectNext()
		case "k", "up":
			m.listState.SelectPrevious()
		case "esc":
			m.listState.Unselect()
		case "q", "ctrl+c":
			return m, bento.Quit
		}
	}

	return m, nil
}

func run() error {
	model := Model{listState: listwidget.NewState()}

	app, err := bento.NewApp(context.Background(), &model)
	if err != nil {
		return fmt.Errorf("new app")
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
