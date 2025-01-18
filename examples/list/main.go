package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/filterablelistwidget"
	"github.com/metafates/bento/listwidget"
	"github.com/metafates/bento/paragraphwidget"
	"github.com/metafates/bento/popupwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	listState filterablelistwidget.State[Item]
	showPopup bool

	currentItem *int
}

type Item struct {
	id int
}

func (i Item) Text() textwidget.Text {
	return textwidget.NewText(
		textwidget.NewLineStr("Item #"+strconv.Itoa(i.id+1)).WithStyle(bento.NewStyle().Bold()),
		textwidget.NewLineStr("Description").WithStyle(bento.NewStyle().Italic().Dim()),
	)
}

func newItems(count int) []Item {
	var items []Item

	for i := 0; i < count; i++ {
		items = append(items, Item{id: i})
	}

	return items
}

func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	var bottomTitle string
	if m.currentItem != nil {
		bottomTitle = strconv.Itoa(*m.currentItem+1) + " of "
	}

	bottomTitle += strconv.Itoa(m.listState.LenFiltered()) + " items"

	block := blockwidget.
		New().
		WithBorderSides().
		Rounded().
		WithTitle(blockwidget.NewTitleStr("List")).
		WithTitle(blockwidget.NewTitleStr(bottomTitle).Bottom().Right())

	list := listwidget.
		New().
		WithHighlightSymbol("> ").
		WithHighlightSpacing(listwidget.HighlightSpacingAlways).
		WithBlock(block).WithHighlightStyle(bento.NewStyle().Black().OnBlue())

	filterablelistwidget.New[Item](list).RenderStateful(area, buffer, &m.listState)

	if m.showPopup {
		paragraph := paragraphwidget.NewStr("Hello, world!").Center()

		popup := popupwidget.
			New(paragraph).
			WithBlock(blockwidget.New().WithBorderSides().Thick().WithTitleStr("Popup")).
			Middle().
			Center().
			WithHeight(bento.ConstraintLen(3)).
			WithWidth(bento.ConstraintLen(30))

		popup.Render(area, buffer)
	}
}

func (m *Model) Init() bento.Cmd {
	return nil
}

func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	consumed, cmd := m.listState.TryUpdate(msg)
	if consumed {
		return m, cmd
	}

	switch msg := msg.(type) {
	case bento.KeyMsg:
		switch msg.String() {
		case " ":
			m.showPopup = !m.showPopup

		case "q", "ctrl+c":
			return m, bento.Quit

		}
	}

	return m, nil
}

func run() error {
	model := Model{listState: filterablelistwidget.NewState(newItems(100)...)}

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
