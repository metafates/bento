package bento

import (
	"bufio"
	"fmt"
	"io"

	"github.com/charmbracelet/x/term"
	"github.com/metafates/bento/internal/ansi"
	"github.com/metafates/bento/internal/bit"
	"github.com/muesli/termenv"
)

type TerminalBackend interface {
	io.Reader

	Draw(cells []PositionedCell) error
	HideCursor() error
	ShowCursor() error
	GetCursorPosition() (Position, error)
	SetCursorPosition(position Position) error
	GetSize() (Size, error)
	Flush() error

	ClearAll() error
	ClearAfterCursor() error
	ClearBeforeCursor() error
	ClearCurrentLine() error
	ClearUntilNewLine() error

	EnableRawMode() error
	DisableRawMode() error

	EnableAlternateScreen() error
	LeaveAlternateScreen() error

	EnableBracketedPaste() error
	DisableBracketedPaste() error

	Input() io.Reader
	Output() io.Writer
}

var _ TerminalBackend = (*DefaultBackend)(nil)

type DefaultBackend struct {
	colorProfile termenv.Profile
	input        io.Reader
	output       io.Writer
	outputBuf    *bufio.Writer

	termState ansi.State
}

func NewDefaultBackend(input io.Reader, output io.Writer) DefaultBackend {
	return DefaultBackend{
		colorProfile: termenv.NewOutput(output).ColorProfile(),
		input:        input,
		output:       output,
		outputBuf:    bufio.NewWriter(output),
	}
}

func (d *DefaultBackend) EnableBracketedPaste() error {
	return d.execute(ansi.EnableBracketedPaste{})
}

func (d *DefaultBackend) DisableBracketedPaste() error {
	return d.execute(ansi.DisableBracketedPaste{})
}

func (d *DefaultBackend) Output() io.Writer {
	return d.output
}

func (d *DefaultBackend) Input() io.Reader {
	return d.input
}

// Read implements TerminalBackend.
func (d *DefaultBackend) Read(p []byte) (n int, err error) {
	return d.input.Read(p)
}

func (d *DefaultBackend) EnableRawMode() error {
	if d.isRawMode() {
		return nil
	}

	fd, ok := d.inputFd()
	if !ok {
		return nil
	}

	if !term.IsTerminal(fd) {
		return nil
	}

	state, err := ansi.EnableRawMode(fd)
	if err != nil {
		return fmt.Errorf("enable raw mode: %w", err)
	}

	d.termState = state

	return nil
}

func (d *DefaultBackend) DisableRawMode() error {
	if !d.isRawMode() {
		return nil
	}

	fd, ok := d.inputFd()
	if !ok {
		return nil
	}

	if !term.IsTerminal(fd) {
		return nil
	}

	if err := ansi.Restore(fd, d.termState); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	d.termState = nil

	return nil
}

func (d *DefaultBackend) isRawMode() bool {
	return d.termState != nil
}

// ClearAfterCursor implements TerminalBackend.
func (d *DefaultBackend) ClearAfterCursor() error {
	return d.execute(ansi.ClearAfterCursor{})
}

// ClearAll implements TerminalBackend.
func (d *DefaultBackend) ClearAll() error {
	return d.execute(ansi.ClearAll{})
}

// ClearBeforeCursor implements TerminalBackend.
func (d *DefaultBackend) ClearBeforeCursor() error {
	return d.execute(ansi.ClearBeforeCursor{})
}

// ClearCurrentLine implements TerminalBackend.
func (d *DefaultBackend) ClearCurrentLine() error {
	return d.execute(ansi.ClearCurrentLine{})
}

// ClearUntilNewLine implements TerminalBackend.
func (d *DefaultBackend) ClearUntilNewLine() error {
	return d.execute(ansi.ClearUntilNewLine{})
}

// Draw implements TerminalBackend.
func (d *DefaultBackend) Draw(cells []PositionedCell) error {
	var (
		lastPos  *Position
		modifier Modifier
	)

	var (
		fg Color = ResetColor{}
		bg Color = ResetColor{}
	)

	for _, pc := range cells {
		x, y := pc.Position.X, pc.Position.Y
		cell := pc.Cell

		if lastPos == nil || lastPos.X+1 != x || lastPos.Y != y {
			if err := d.queue(ansi.MoveTo{
				Column: x,
				Row:    y,
			}); err != nil {
				return fmt.Errorf("set cursor position: %w", err)
			}
		}

		lastPos = &Position{X: x, Y: y}

		if cell.Modifier != modifier {
			diff := _StyleModifierDiff{
				From: modifier,
				To:   cell.Modifier,
			}

			if err := diff.queue(d.outputBuf); err != nil {
				return fmt.Errorf("queue: %w", err)
			}

			modifier = cell.Modifier
		}

		if cell.Fg != fg || cell.Bg != bg {
			if err := d.queue(ansi.SetColors(ansi.Colors{
				Foreground: d.colorProfile.Convert(cell.Fg),
				Background: d.colorProfile.Convert(cell.Bg),
			})); err != nil {
				return fmt.Errorf("queue: %w", err)
			}

			fg = cell.Fg
			bg = cell.Bg
		}

		if err := d.queue(ansi.Print(cell.Symbol)); err != nil {
			return fmt.Errorf("queue: %w", err)
		}
	}

	if err := d.queue(
		ansi.SetColors(ansi.Colors{
			Foreground: ResetColor{},
			Background: ResetColor{},
		}),
		ansi.SetAttribute(ansi.AttrReset),
	); err != nil {
		return fmt.Errorf("queue: %w", err)
	}

	return nil
}

// Flush implements TerminalBackend.
func (d *DefaultBackend) Flush() error {
	return d.outputBuf.Flush()
}

// GetCursorPosition implements TerminalBackend.
func (d *DefaultBackend) GetCursorPosition() (Position, error) {
	column, row, err := ansi.GetCursorPosition()
	if err != nil {
		return Position{}, fmt.Errorf("get cursor position: %w", err)
	}

	return Position{
		X: column,
		Y: row,
	}, nil
}

// GetSize implements TerminalBackend.
func (d *DefaultBackend) GetSize() (Size, error) {
	file, ok := d.output.(term.File)
	if !ok {
		return Size{}, nil
	}

	fd := file.Fd()

	width, height, err := ansi.GetSize(fd)
	if err != nil {
		return Size{}, fmt.Errorf("get size: %w", err)
	}

	return Size{
		Width:  width,
		Height: height,
	}, nil
}

// HideCursor implements TerminalBackend.
func (d *DefaultBackend) HideCursor() error {
	return d.execute(ansi.HideCursor{})
}

func (d *DefaultBackend) EnableAlternateScreen() error {
	return d.execute(ansi.EnterAlternateScreen{})
}

func (d *DefaultBackend) LeaveAlternateScreen() error {
	return d.execute(ansi.LeaveAlternateScreen{})
}

// SetCursorPosition implements TerminalBackend.
func (d *DefaultBackend) SetCursorPosition(position Position) error {
	return d.execute(ansi.MoveTo{
		Column: position.X,
		Row:    position.Y,
	})
}

// ShowCursor implements TerminalBackend.
func (d *DefaultBackend) ShowCursor() error {
	return d.execute(ansi.ShowCursor{})
}

func (d *DefaultBackend) queue(commands ...ansi.Command) error {
	return queue(d.outputBuf, commands...)
}

func (d *DefaultBackend) execute(commands ...ansi.Command) error {
	if err := d.queue(commands...); err != nil {
		return fmt.Errorf("queue: %w", err)
	}

	if err := d.Flush(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	return nil
}

func (d *DefaultBackend) outputFd() (uintptr, bool) {
	file, ok := d.output.(term.File)
	if !ok {
		return 0, false
	}

	return file.Fd(), true
}

func (d *DefaultBackend) inputFd() (uintptr, bool) {
	file, ok := d.input.(term.File)
	if !ok {
		return 0, false
	}

	return file.Fd(), true
}

type _StyleModifierDiff struct {
	From, To Modifier
}

func (d _StyleModifierDiff) queue(w io.Writer) error {
	var cmds []ansi.Command

	removed := bit.Difference(d.From, d.To)

	if removed.Contains(ModifierReversed) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoReverse))
	}

	if removed.Contains(ModifierBold) || removed.Contains(ModifierDim) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNormalIntensity))

		if d.To.Contains(ModifierDim) {
			cmds = append(cmds, ansi.SetAttribute(ansi.AttrDim))
		}

		if d.To.Contains(ModifierBold) {
			cmds = append(cmds, ansi.SetAttribute(ansi.AttrBold))
		}
	}

	if removed.Contains(ModifierItalic) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoItalic))
	}

	if removed.Contains(ModifierUnderlined) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoUnderline))
	}

	if removed.Contains(ModifierCrossedOut) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNotCrossedOut))
	}

	if removed.Contains(ModifierSlowBlink) || removed.Contains(ModifierRapidBlink) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoBlink))
	}

	added := bit.Difference(d.To, d.From)

	if added.Contains(ModifierReversed) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrReverse))
	}

	if added.Contains(ModifierBold) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrBold))
	}

	if added.Contains(ModifierItalic) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrItalic))
	}

	if added.Contains(ModifierUnderlined) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrUnderlined))
	}

	if added.Contains(ModifierDim) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrDim))
	}

	if added.Contains(ModifierCrossedOut) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrCrossedOut))
	}

	if added.Contains(ModifierSlowBlink) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrSlowBlink))
	}

	if added.Contains(ModifierRapidBlink) {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrRapidBlink))
	}

	if err := queue(w, cmds...); err != nil {
		return fmt.Errorf("queue: %w", err)
	}

	return nil
}

func queue(w io.Writer, commands ...ansi.Command) error {
	for _, c := range commands {
		if err := c.WriteANSI(w); err != nil {
			return fmt.Errorf("write ansi: %w", err)
		}
	}

	return nil
}
