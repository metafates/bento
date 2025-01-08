package main

import (
	"context"
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/textwidget"
)

var _ bento.Model = (*Model)(nil)

type Model struct{}

// Draw implements bento.Model.
func (m *Model) Draw(frame *bento.Frame) {
	var primaryArea, footnoteArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		WithConstraints(
			bento.ConstraintFill(1),
			bento.ConstraintLength(1),
		).
		Split(frame.Area()).
		Assign(&primaryArea, &footnoteArea)

	m.drawPrimary(frame, primaryArea)
	m.drawFootnote(frame, footnoteArea)
}

func (m *Model) drawFootnote(frame *bento.Frame, area bento.Rect) {
	help := textwidget.NewLineString("help").Left()
	version := textwidget.NewLineString("version").Right()

	frame.RenderWidget(help, area)
	frame.RenderWidget(version, area)
}

func (m *Model) drawPrimary(frame *bento.Frame, area bento.Rect) {
	var sidebarArea, rightArea bento.Rect

	bento.
		NewLayout().
		Horizontal().
		WithConstraints(
			bento.ConstraintFill(1),
			bento.ConstraintFill(2),
		).
		Split(area).
		Assign(&sidebarArea, &rightArea)

	m.drawSidebar(frame, sidebarArea)
	m.drawRight(frame, rightArea)
}

func (m *Model) drawSidebar(frame *bento.Frame, area bento.Rect) {
	var statusArea,
		filesArea,
		branchesArea,
		commitsArea,
		stashArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		WithConstraints(
			bento.ConstraintLength(3),
			bento.ConstraintFill(1),
			bento.ConstraintFill(1),
			bento.ConstraintFill(1),
			bento.ConstraintLength(3),
		).
		Split(area).
		Assign(
			&statusArea,
			&filesArea,
			&branchesArea,
			&commitsArea,
			&stashArea,
		)

	m.drawStatus(frame, statusArea)
	m.drawFiles(frame, filesArea)
	m.drawBranches(frame, branchesArea)
	m.drawCommits(frame, commitsArea)
	m.drawStash(frame, stashArea)
}

func (m *Model) drawStatus(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Status").Plain()

	innerArea := block.Inner(area)

	status := textwidget.NewLineString("some status here")

	frame.RenderWidget(block, area)
	frame.RenderWidget(status, innerArea)
}

func (m *Model) drawFiles(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Files").Plain()

	frame.RenderWidget(block, area)
}

func (m *Model) drawBranches(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Branches").Plain()

	frame.RenderWidget(block, area)
}

func (m *Model) drawCommits(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Commits").Plain()

	frame.RenderWidget(block, area)
}

func (m *Model) drawStash(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Stash").Plain()

	frame.RenderWidget(block, area)
}

func (m *Model) drawRight(frame *bento.Frame, area bento.Rect) {
	var infoArea, commandLogArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		WithConstraints(
			bento.ConstraintFill(1),
			bento.ConstraintLength(10),
		).
		Split(area).
		Assign(&infoArea, &commandLogArea)

	m.renderInfoArea(frame, infoArea)
	m.renderCommandLog(frame, commandLogArea)
}

func (m *Model) renderInfoArea(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Info").Plain()

	frame.RenderWidget(block, area)
}

func (m *Model) renderCommandLog(frame *bento.Frame, area bento.Rect) {
	block := blockwidget.NewBlock().WithTitleString("Command log").Plain()

	frame.RenderWidget(block, area)
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
		case "ctrl+c", "q":
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
