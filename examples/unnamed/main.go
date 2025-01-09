package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/listwidget"
	"github.com/metafates/bento/paragraphwidget"
	"github.com/metafates/bento/popupwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	listState listwidget.State
	showPopup bool
}

func (m *Model) Draw(frame *bento.Frame) {
	block := blockwidget.
		NewBlock().
		WithBorders().
		Rounded().
		WithTitle(blockwidget.NewTitleString("List"))

	var items []listwidget.Item

	for i := 0; i < 100; i++ {
		items = append(items, listwidget.NewItem(textwidget.NewText(
			textwidget.NewLineString("Item #"+strconv.Itoa(i+1)).WithStyle(bento.NewStyle().Bold()),
			textwidget.NewLineString("Description").WithStyle(bento.NewStyle().Italic().Dim()),
		)))
	}

	list := listwidget.
		NewList(items...).
		WithHighlightSymbol("> ").
		WithHighlightSpacing(listwidget.HighlightSpacingAlways).
		WithBlock(block).WithHighlightStyle(bento.NewStyle().Black().OnBlue())

	bento.RenderStatefulWidget(frame, list, frame.Area(), &m.listState)

	if m.showPopup {
		popup := popupwidget.New(paragraphwidget.NewParagraphString("Hello, world!").Center()).WithBlock(blockwidget.NewBlock().WithBorders().WithTitleString("Popup"))
		popup = popup.Top().Right().WithHeight(bento.ConstraintLength(3)).WithWidth(bento.ConstraintPercentage(30))

		frame.RenderWidget(popup, frame.Area())
	}
}

func (m *Model) Init() bento.Cmd {
	return nil
}

func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case bento.KeyMsg:
		switch msg.String() {
		case " ":
			m.showPopup = !m.showPopup
		case "ctrl+u":
			m.listState.ScrollUpBy(8)
		case "ctrl+d":
			m.listState.ScrollDownBy(8)
		case "G":
			m.listState.SelectLast()
		case "g":
			m.listState.SelectFirst()
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
