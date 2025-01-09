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

type Panel int

const (
	PanelStatus Panel = iota
	PanelFiles
	PanelBranches
	PanelCommits
	PanelStash
	PanelInfo
	PanelCommandLog
)

func (p Panel) Next() Panel {
	return (p + 1) % PanelInfo
}

func (p Panel) Prev() Panel {
	if p-1 < 0 {
		return PanelStash
	}

	return p - 1
}

func (p Panel) Title() string {
	switch p {
	case PanelBranches:
		return "Branches"
	case PanelCommits:
		return "Commits"
	case PanelFiles:
		return "Files"
	case PanelStash:
		return "Stash"
	case PanelStatus:
		return "Status"
	case PanelCommandLog:
		return "Command log"
	case PanelInfo:
		return "Info"
	default:
		panic(fmt.Sprintf("unexpected Panel: %#v", p))
	}
}

type Model struct {
	size        bento.Size
	activePanel Panel
}

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
	left := textwidget.NewLine(
		textwidget.NewSpan("Quit: q / ctrl+c").WithStyle(bento.NewStyle().Blue()),
	).Left()

	right := textwidget.NewLine(
		textwidget.NewSpan("Foo").WithStyle(bento.NewStyle().Underlined().Magenta()),
		textwidget.NewSpan(" "),
		textwidget.NewSpan("Bar").WithStyle(bento.NewStyle().Underlined().Yellow()),
		textwidget.NewSpan(" "),
		textwidget.NewSpan("0.0.1"),
	).Right()

	frame.RenderWidget(left, area)
	frame.RenderWidget(right, area)
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

const (
	HeightMedium = 30
	HeightSmall  = 22
)

func (m *Model) filesConstraint() bento.Constraint {
	if m.size.Height > HeightMedium || m.activePanel == PanelFiles {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLength(3)
	}

	return bento.ConstraintLength(1)
}

func (m *Model) branchesConstraint() bento.Constraint {
	if m.size.Height > HeightMedium || m.activePanel == PanelBranches {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLength(3)
	}

	return bento.ConstraintLength(1)
}

func (m *Model) commitsConstraint() bento.Constraint {
	if m.size.Height > HeightMedium || m.activePanel == PanelCommits {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLength(3)
	}

	return bento.ConstraintLength(1)
}

func (m *Model) statusConstraint() bento.Constraint {
	if m.size.Height > HeightMedium {
		return bento.ConstraintLength(3)
	}

	if m.activePanel == PanelStatus {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLength(3)
	}

	return bento.ConstraintLength(1)
}

func (m *Model) stashConstraint() bento.Constraint {
	isActive := m.activePanel == PanelStash

	if m.size.Height > HeightMedium || (m.size.Height > HeightSmall && !isActive) {
		return bento.ConstraintLength(3)
	}

	if isActive {
		return bento.ConstraintFill(1)
	}

	return bento.ConstraintLength(1)
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
			m.statusConstraint(),
			m.filesConstraint(),
			m.branchesConstraint(),
			m.commitsConstraint(),
			m.stashConstraint(),
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
	block := m.newBlock(PanelStatus, "")

	innerArea := block.Inner(area)

	status := textwidget.NewLineString(fmt.Sprintf("%dx%d", m.size.Width, m.size.Height))

	frame.RenderWidget(block, area)
	frame.RenderWidget(status, innerArea)
}

func (m *Model) drawFiles(frame *bento.Frame, area bento.Rect) {
	block := m.newBlock(PanelFiles, "1 of 10")

	frame.RenderWidget(block, area)
}

func (m *Model) drawBranches(frame *bento.Frame, area bento.Rect) {
	block := m.newBlock(PanelBranches, "1 of 2")

	frame.RenderWidget(block, area)
}

func (m *Model) drawCommits(frame *bento.Frame, area bento.Rect) {
	block := m.newBlock(PanelCommits, "1 of 42")

	frame.RenderWidget(block, area)
}

func (m *Model) drawStash(frame *bento.Frame, area bento.Rect) {
	block := m.newBlock(PanelStash, "1 of 1")

	frame.RenderWidget(block, area)
}

func (m *Model) drawRight(frame *bento.Frame, area bento.Rect) {
	var infoArea, commandLogArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		WithConstraints(
			bento.ConstraintFill(1),
			bento.ConstraintLength(m.commandLogHeight()),
		).
		Split(area).
		Assign(&infoArea, &commandLogArea)

	m.renderInfoArea(frame, infoArea)
	m.renderCommandLog(frame, commandLogArea)
}

func (m *Model) commandLogHeight() int {
	if m.size.Height <= 40 {
		return 3
	}

	return 10
}

func (m *Model) renderInfoArea(frame *bento.Frame, area bento.Rect) {
	block := m.newBlock(PanelInfo, "")

	frame.RenderWidget(block, area)
}

func (m *Model) renderCommandLog(frame *bento.Frame, area bento.Rect) {
	block := m.newBlock(PanelCommandLog, "")

	frame.RenderWidget(block, area)
}

func (m *Model) newBlock(panel Panel, footer string) blockwidget.Block {
	block := blockwidget.
		NewBlock().
		WithBorders().
		Rounded().
		WithTitle(blockwidget.NewTitleString(panel.Title()).Top().Left())

	if m.activePanel == panel {
		block = block.WithBorderStyle(bento.NewStyle().Bold().Green())
	}

	if footer == "" {
		return block
	}

	return block.WithTitle(blockwidget.NewTitleString(footer).Bottom().Right())
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case bento.WindowSizeMsg:
		m.size = bento.Size(msg)
	case bento.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, bento.Quit
		case "shift+tab":
			m.activePanel = m.activePanel.Prev()
		case "tab":
			m.activePanel = m.activePanel.Next()
		}
	}

	return m, nil
}

func run() error {
	app, err := bento.NewApp(context.Background(), &Model{
		activePanel: PanelFiles,
	})
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
