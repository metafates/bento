package bento

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/metafates/bento/internal/ansi"
	"github.com/muesli/termenv"
)

var _ TerminalBackend = (*DefaultBackend)(nil)

type DefaultBackend struct {
	input  io.Reader
	tty    *os.File
	ttyBuf *bufio.Writer

	termState ansi.State
}

func NewDefaultBackend() DefaultBackend {
	return DefaultBackend{
		input:  os.Stdin,
		tty:    os.Stdout,
		ttyBuf: bufio.NewWriter(os.Stdout),
	}
}

// Read implements TerminalBackend.
func (d *DefaultBackend) Read(p []byte) (n int, err error) {
	return d.input.Read(p)
}

func (d *DefaultBackend) EnableRawMode() error {
	if d.isRawMode() {
		return nil
	}

	state, err := ansi.EnableRawMode(int(d.fd()))
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

	if err := ansi.Restore(int(d.fd()), d.termState); err != nil {
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
		lastPos *Position
		style   Style
	)

	var (
		fg termenv.Color = ansi.ResetColor{}
		bg termenv.Color = ansi.ResetColor{}
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

		if cell.Style != style {
			diff := _StyleDiff{
				From: style,
				To:   cell.Style,
			}

			if err := diff.queue(d.ttyBuf); err != nil {
				return fmt.Errorf("queue: %w", err)
			}

			style = cell.Style
		}

		if cell.Style.Foreground.Color() != fg || cell.Style.Background.Color() != bg {
			if err := d.queue(ansi.SetColors(ansi.Colors{
				Foreground: fg,
				Background: bg,
			})); err != nil {
				return fmt.Errorf("queue: %w", err)
			}

			fg = cell.Style.Foreground.Color()
			bg = cell.Style.Background.Color()
		}

		if err := d.queue(ansi.Print(cell.Symbol)); err != nil {
			return fmt.Errorf("queue: %w", err)
		}
	}

	if err := d.queue(
		ansi.SetColors(ansi.Colors{
			Foreground: ansi.ResetColor{},
			Background: ansi.ResetColor{},
		}),
		ansi.SetAttribute(ansi.AttrReset),
	); err != nil {
		return fmt.Errorf("queue: %w", err)
	}

	return nil
}

// Flush implements TerminalBackend.
func (d *DefaultBackend) Flush() error {
	return d.ttyBuf.Flush()
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
	fd := d.tty.Fd()

	width, height, err := ansi.GetSize(int(fd))
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
	return queue(d.ttyBuf, commands...)
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

func (d *DefaultBackend) fd() uintptr {
	return d.tty.Fd()
}

type _StyleDiff struct {
	From, To Style
}

func (d _StyleDiff) queue(w io.Writer) error {
	var cmds []ansi.Command

	removed := d.From.Sub(d.To)

	if removed.Reversed.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNormalIntensity))

		if d.To.Dim.IsSet() {
			cmds = append(cmds, ansi.SetAttribute(ansi.AttrDim))
		}

		if d.To.Bold.IsSet() {
			cmds = append(cmds, ansi.SetAttribute(ansi.AttrBold))
		}
	}

	if removed.Italic.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoItalic))
	}

	if removed.Underlined.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoUnderline))
	}

	if removed.CrossedOut.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNotCrossedOut))
	}

	if removed.SlowBlink.IsSet() || removed.RapidBlink.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrNoBlink))
	}

	added := d.To.Sub(d.From)

	if added.Reversed.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrReverse))
	}

	if added.Bold.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrBold))
	}

	if added.Italic.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrItalic))
	}

	if added.Underlined.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrUnderlined))
	}

	if added.Dim.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrDim))
	}

	if added.CrossedOut.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrCrossedOut))
	}

	if added.SlowBlink.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrSlowBlink))
	}

	if added.RapidBlink.IsSet() {
		cmds = append(cmds, ansi.SetAttribute(ansi.AttrRapidBlink))
	}

	for _, c := range cmds {
		if err := queue(w, c); err != nil {
			return fmt.Errorf("queue: %w", err)
		}
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
