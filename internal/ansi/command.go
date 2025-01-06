package ansi

import (
	"fmt"
	"io"
	"strconv"

	"github.com/muesli/termenv"
)

const (
	// Escape character
	ESC = '\x1b'
	// Bell
	BEL = '\a'
	// Control Sequence Introducer
	CSI = string(ESC) + "["
	// Operating System Command
	OSC = string(ESC) + "]"
	// String Terminator
	ST = string(ESC) + `\`
)

type Command interface {
	WriteANSI(w io.Writer) error
}

var _ Command = (*MoveTo)(nil)

type MoveTo struct{ Column, Row int }

func (m MoveTo) WriteANSI(w io.Writer) error {
	column := strconv.Itoa(m.Column + 1)
	row := strconv.Itoa(m.Row + 1)

	return write(w, CSI+row+";"+column+"H")
}

var _ Command = (*SetAttribute)(nil)

type SetAttribute Attr

func (a SetAttribute) WriteANSI(w io.Writer) error {
	sgr := Attr(a).SGR()

	return write(w, CSI+sgr+"m")
}

type Colors struct {
	Foreground, Background termenv.Color
}

var _ Command = (*SetColors)(nil)

type SetColors Colors

func (c SetColors) WriteANSI(w io.Writer) error {
	switch {
	case c.Foreground != nil && c.Background != nil:
		return write(w, termenv.CSI+c.Foreground.Sequence(false)+";"+c.Background.Sequence(true)+"m")
	case c.Foreground == nil && c.Background != nil:
		return write(w, termenv.CSI+c.Foreground.Sequence(false)+"m")
	case c.Foreground != nil && c.Background == nil:
		return write(w, termenv.CSI+c.Background.Sequence(true)+"m")
	default:
		return nil
	}
}

var _ Command = (*Print)(nil)

type Print string

func (p Print) WriteANSI(w io.Writer) error {
	return write(w, string(p))
}

var _ Command = (*ShowCursor)(nil)

type ShowCursor struct{}

func (ShowCursor) WriteANSI(w io.Writer) error {
	return write(w, CSI+"?25h")
}

var _ Command = (*HideCursor)(nil)

type HideCursor struct{}

func (HideCursor) WriteANSI(w io.Writer) error {
	return write(w, CSI+"?25l")
}

type ClearAll struct{}

func (ClearAll) WriteANSI(w io.Writer) error {
	return write(w, CSI+"2j")
}

type ClearAfterCursor struct{}

func (ClearAfterCursor) WriteANSI(w io.Writer) error {
	return write(w, CSI+"J")
}

type ClearBeforeCursor struct{}

func (ClearBeforeCursor) WriteANSI(w io.Writer) error {
	return write(w, CSI+"1J")
}

type ClearCurrentLine struct{}

func (ClearCurrentLine) WriteANSI(w io.Writer) error {
	return write(w, CSI+"2K")
}

type ClearUntilNewLine struct{}

func (ClearUntilNewLine) WriteANSI(w io.Writer) error {
	return write(w, CSI+"K")
}

type EnterAlternateScreen struct{}

func (EnterAlternateScreen) WriteANSI(w io.Writer) error {
	return write(w, CSI+"?1049h")
}

type LeaveAlternateScreen struct{}

func (LeaveAlternateScreen) WriteANSI(w io.Writer) error {
	return write(w, CSI+"?1049l")
}

func write(w io.Writer, a ...any) error {
	_, err := fmt.Fprint(w, a...)

	return err
}

func writef(w io.Writer, format string, a ...any) error {
	_, err := fmt.Fprintf(w, format, a...)

	return err
}
