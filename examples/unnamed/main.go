package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/gaugewidget"
	"github.com/metafates/bento/paragraphwidget"
	"github.com/metafates/bento/popupwidget"
	"github.com/metafates/bento/scrollwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct {
	text           string
	verticalScroll scrollwidget.State
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Render implements bento.Model.
func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	content := paragraphwidget.
		NewStr(m.text).
		Wrapped(paragraphwidget.NewWrap()).
		WithScroll(m.verticalScroll.Position(), 0)

	scroll := scrollwidget.New(scrollwidget.OrientationVerticalRight)

	innerArea := scroll.Inner(area)

	gauge := gaugewidget.New().WithRatio(m.verticalScroll.Ratio()).WithUnicode(true)
	popup := popupwidget.New(gauge).Bottom().Left().WithHeight(bento.ConstraintLen(3)).WithWidth(bento.ConstraintPercentage(30))

	content.Render(innerArea, buffer)
	scroll.RenderStateful(area, buffer, m.verticalScroll)
	popup.Render(area, buffer)
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case bento.KeyMsg:
		switch msg.String() {
		case "up":
			m.verticalScroll.Prev()
		case "down":
			m.verticalScroll.Next()
		case "ctrl+c":
			return m, bento.Quit
		}
	}

	return m, nil
}

func run() error {
	model := Model{
		text:           strings.Repeat("Lorem ipsum dolor sit amet. ", 40),
		verticalScroll: scrollwidget.NewState(100),
	}

	_, err := bento.NewApp(&model).Run()
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
