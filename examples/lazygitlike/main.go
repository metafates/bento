package main

import (
	"fmt"
	"log"

	"github.com/metafates/bento"
	"github.com/metafates/bento/widget/blockwidget"
	"github.com/metafates/bento/widget/textwidget"
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

func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	var primaryArea, footnoteArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		WithConstraints(
			bento.ConstraintFill(1),
			bento.ConstraintLen(1),
		).
		Split(area).
		Assign(&primaryArea, &footnoteArea)

	m.drawPrimary(primaryArea, buffer)
	m.drawFootnote(footnoteArea, buffer)
}

func (m *Model) drawFootnote(area bento.Rect, buffer *bento.Buffer) {
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

	left.Render(area, buffer)
	right.Render(area, buffer)
}

func (m *Model) drawPrimary(area bento.Rect, buffer *bento.Buffer) {
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

	m.drawSidebar(sidebarArea, buffer)
	m.drawRight(rightArea, buffer)
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
		return bento.ConstraintLen(3)
	}

	return bento.ConstraintLen(1)
}

func (m *Model) branchesConstraint() bento.Constraint {
	if m.size.Height > HeightMedium || m.activePanel == PanelBranches {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLen(3)
	}

	return bento.ConstraintLen(1)
}

func (m *Model) commitsConstraint() bento.Constraint {
	if m.size.Height > HeightMedium || m.activePanel == PanelCommits {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLen(3)
	}

	return bento.ConstraintLen(1)
}

func (m *Model) statusConstraint() bento.Constraint {
	if m.size.Height > HeightMedium {
		return bento.ConstraintLen(3)
	}

	if m.activePanel == PanelStatus {
		return bento.ConstraintFill(1)
	}

	if m.size.Height > HeightSmall {
		return bento.ConstraintLen(3)
	}

	return bento.ConstraintLen(1)
}

func (m *Model) stashConstraint() bento.Constraint {
	isActive := m.activePanel == PanelStash

	if m.size.Height > HeightMedium || (m.size.Height > HeightSmall && !isActive) {
		return bento.ConstraintLen(3)
	}

	if isActive {
		return bento.ConstraintFill(1)
	}

	return bento.ConstraintLen(1)
}

func (m *Model) drawSidebar(area bento.Rect, buffer *bento.Buffer) {
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

	m.drawStatus(statusArea, buffer)
	m.drawFiles(filesArea, buffer)
	m.drawBranches(branchesArea, buffer)
	m.drawCommits(commitsArea, buffer)
	m.drawStash(stashArea, buffer)
}

func (m *Model) drawStatus(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelStatus, "")

	innerArea := block.Inner(area)

	status := textwidget.NewLineStr(fmt.Sprintf("%dx%d", m.size.Width, m.size.Height))

	block.Render(area, buffer)
	status.Render(innerArea, buffer)
}

func (m *Model) drawFiles(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelFiles, "1 of 10")

	block.Render(area, buffer)
}

func (m *Model) drawBranches(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelBranches, "1 of 2")

	block.Render(area, buffer)
}

func (m *Model) drawCommits(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelCommits, "1 of 42")

	block.Render(area, buffer)
}

func (m *Model) drawStash(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelStash, "1 of 1")

	block.Render(area, buffer)
}

func (m *Model) drawRight(area bento.Rect, buffer *bento.Buffer) {
	var infoArea, commandLogArea bento.Rect

	bento.
		NewLayout().
		Vertical().
		WithConstraints(
			bento.ConstraintFill(1),
			bento.ConstraintLen(m.commandLogHeight()),
		).
		Split(area).
		Assign(&infoArea, &commandLogArea)

	m.renderInfoArea(infoArea, buffer)
	m.renderCommandLog(commandLogArea, buffer)
}

func (m *Model) commandLogHeight() int {
	if m.size.Height <= 40 {
		return 3
	}

	return 10
}

func (m *Model) renderInfoArea(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelInfo, "")

	block.Render(area, buffer)
}

func (m *Model) renderCommandLog(area bento.Rect, buffer *bento.Buffer) {
	block := m.newBlock(PanelCommandLog, "")

	block.Render(area, buffer)
}

func (m *Model) newBlock(panel Panel, footer string) blockwidget.Block {
	block := blockwidget.
		New().
		WithBorderSides().
		Rounded().
		WithTitle(blockwidget.NewTitleStr(panel.Title()).Top().Left())

	if m.activePanel == panel {
		block = block.WithBorderStyle(bento.NewStyle().Bold().Green())
	}

	if footer == "" {
		return block
	}

	return block.WithTitle(blockwidget.NewTitleStr(footer).Bottom().Right())
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
	_, err := bento.NewApp(&Model{activePanel: PanelFiles}).Run()
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
