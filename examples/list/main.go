package main

import (
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

	currentItem *int
	itemsCount  int
}

func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	var bottomTitle string
	if m.currentItem != nil {
		bottomTitle = strconv.Itoa(*m.currentItem+1) + " of "
	}

	bottomTitle += strconv.Itoa(m.itemsCount) + " items"

	block := blockwidget.
		New().
		WithBorderSides().
		Rounded().
		WithTitle(blockwidget.NewTitleStr("List")).
		WithTitle(blockwidget.NewTitleStr(bottomTitle).Bottom().Right())

	var items []listwidget.Item

	for i := 0; i < m.itemsCount; i++ {
		items = append(items, listwidget.NewItem(textwidget.NewText(
			textwidget.NewLineStr("Item #"+strconv.Itoa(i+1)).WithStyle(bento.NewStyle().Bold()),
			textwidget.NewLineStr("Description").WithStyle(bento.NewStyle().Italic().Dim()),
		)))
	}

	list := listwidget.
		NewList(items...).
		WithHighlightSymbol("> ").
		WithHighlightSpacing(listwidget.HighlightSpacingAlways).
		WithBlock(block).WithHighlightStyle(bento.NewStyle().Black().OnBlue())

	list.RenderStateful(area, buffer, &m.listState)

	if m.showPopup {
		paragraph := paragraphwidget.NewParagraphStr("Hello, world!").Center()

		popup := popupwidget.
			New(paragraph).
			WithBlock(blockwidget.New().WithBorderSides().Thick().WithTitleStr("Popup")).
			Middle().
			Center().
			WithHeight(bento.ConstraintLength(3)).
			WithWidth(bento.ConstraintLength(30))

		popup.Render(area, buffer)
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

	index, ok := m.listState.SelectedWithLimit(m.itemsCount - 1)
	if ok {
		m.currentItem = &index
	} else {
		m.currentItem = nil
	}

	return m, nil
}

func run() error {
	model := Model{listState: listwidget.NewState(), itemsCount: 100}

	_, err := bento.NewApp(&model).Run()
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
