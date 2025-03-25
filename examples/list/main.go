package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/metafates/bento"
	"github.com/metafates/bento/widget/blockwidget"
	"github.com/metafates/bento/widget/listwidget"
	"github.com/metafates/bento/widget/paragraphwidget"
	"github.com/metafates/bento/widget/popupwidget"
	"github.com/metafates/bento/widget/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	listState listwidget.State
	showPopup bool

	items       []Item
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

	bottomTitle += strconv.Itoa(len(m.items)) + " items"

	block := blockwidget.
		New().
		WithBorderSides().
		Rounded().
		WithTitle(blockwidget.NewTitleStr("List")).
		WithTitle(blockwidget.NewTitleStr(bottomTitle).Bottom().Right())

	items := make([]textwidget.Text, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item.Text())
	}

	list := listwidget.
		New(items...).
		WithHighlightSymbol("> ").
		WithHighlightSpacing(listwidget.HighlightSpacingAlways).
		WithBlock(block).WithHighlightStyle(bento.NewStyle().Black().OnBlue())

	list.RenderStateful(area, buffer, &m.listState)

	if m.showPopup {
		paragraph := paragraphwidget.NewStr("Hello, world!").Center()

		popup := popupwidget.
			New().
			WithBlock(blockwidget.New().WithBorderSides().Thick().WithTitleStr("Popup")).
			Middle().
			Center().
			WithHeight(bento.ConstraintLen(3)).
			WithWidth(bento.ConstraintLen(30))

		popup.Render(area, buffer)
		paragraph.Render(popup.Inner(area), buffer)
	}
}

func (m *Model) Init() bento.Cmd {
	return nil
}

func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	cmd, ok := m.listState.TryUpdate(msg)
	if ok {
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
	model := Model{
		listState: listwidget.NewState(),
		items:     newItems(100),
	}

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
