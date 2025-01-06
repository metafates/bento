package bento

import (
	"bufio"
	"fmt"
	"io"

	"github.com/metafates/bento/internal/ansi"
	"github.com/muesli/termenv"
)

var _ TerminalBackend = (*DefaultBackend)(nil)

type DefaultBackend struct {
	writer *bufio.Writer
	output *termenv.Output
}

// ClearAfterCursor implements TerminalBackend.
func (d *DefaultBackend) ClearAfterCursor() error {
	d.output.ClearLineRight()

	return nil
}

// ClearAll implements TerminalBackend.
func (d *DefaultBackend) ClearAll() error {
	d.output.ClearScreen()

	return nil
}

// ClearBeforeCursor implements TerminalBackend.
func (d *DefaultBackend) ClearBeforeCursor() error {
	d.output.ClearLineLeft()

	return nil
}

// ClearCurrentLine implements TerminalBackend.
func (d *DefaultBackend) ClearCurrentLine() error {
	d.output.ClearLine()

	return nil
}

// ClearUntilNewLine implements TerminalBackend.
func (d *DefaultBackend) ClearUntilNewLine() error {
	panic("unimplemented")
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
			d.output.MoveCursor(y, x)
		}

		lastPos = &Position{X: x, Y: y}

		if cell.Style != style {
			diff := _StyleDiff{
				From: style,
				To:   cell.Style,
			}

			if err := diff.queue(d.output); err != nil {
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
	panic("unimplemented")
}

// GetCursorPosition implements TerminalBackend.
func (d *DefaultBackend) GetCursorPosition() (Position, error) {
	panic("unimplemented")
}

// GetSize implements TerminalBackend.
func (d *DefaultBackend) GetSize() (Size, error) {
	panic("unimplemented")
}

// HideCursor implements TerminalBackend.
func (d *DefaultBackend) HideCursor() error {
	d.output.HideCursor()

	return nil
}

// SetCursorPosition implements TerminalBackend.
func (d *DefaultBackend) SetCursorPosition(position Position) error {
	d.output.MoveCursor(position.Y, position.X)

	return nil
}

// ShowCursor implements TerminalBackend.
func (d *DefaultBackend) ShowCursor() error {
	d.output.ShowCursor()

	return nil
}

func (d *DefaultBackend) queue(commands ...ansi.Command) error {
	return queue(d.writer, commands...)
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
