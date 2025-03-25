package bento

import (
	"fmt"
	"io"
)

type PositionedCell struct {
	Cell

	Position Position
}

var _ io.Reader = (*Terminal)(nil)

type Terminal struct {
	backend      TerminalBackend
	viewport     Viewport
	viewportArea Rect

	buffers [2]Buffer
	current int

	lastKnownArea      Rect
	lastKnownCursorPos Position

	hiddenCursor bool

	frameCount int
}

func NewTerminal(backend TerminalBackend, viewport Viewport) (*Terminal, error) {
	var area Rect

	switch v := viewport.(type) {
	case ViewportFullscreen, ViewportInline:
		size, _, err := backend.GetSize()
		if err != nil {
			return nil, fmt.Errorf("get size: %w", err)
		}

		area = NewRect(size.Width, size.Height)
	case ViewportFixed:
		area = Rect(v)
	}

	var (
		viewportArea Rect
		cursorPos    Position
	)

	switch v := viewport.(type) {
	case ViewportFullscreen:
		viewportArea = area
		cursorPos = Position{}
	case ViewportInline:
		panic("unimplemented")
	case ViewportFixed:
		viewportArea = Rect(v)
		cursorPos = viewportArea.Position
	}

	return &Terminal{
		backend:      backend,
		viewport:     viewport,
		viewportArea: viewportArea,
		buffers: [2]Buffer{
			NewBufferEmpty(viewportArea),
			NewBufferEmpty(viewportArea),
		},
		current:            0,
		lastKnownArea:      area,
		lastKnownCursorPos: cursorPos,
		hiddenCursor:       false,
		frameCount:         0,
	}, nil
}

func (t *Terminal) Viewport() Viewport {
	return t.viewport
}

func (t *Terminal) Input() io.Reader {
	return t.backend.Input()
}

func (t *Terminal) Output() io.Writer {
	return t.backend.Output()
}

// Read implements io.Reader.
func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.backend.Read(p)
}

func (t *Terminal) DisableBracketedPaste() error {
	return t.backend.DisableBracketedPaste()
}

func (t *Terminal) EnableBracketedPaste() error {
	return t.backend.EnableBracketedPaste()
}

func (t *Terminal) EnableAlternateScreen() error {
	return t.backend.EnableAlternateScreen()
}

func (t *Terminal) LeaveAlternateScreen() error {
	return t.backend.LeaveAlternateScreen()
}

func (t *Terminal) EnableRawMode() error {
	return t.backend.EnableRawMode()
}

func (t *Terminal) DisableRawMode() error {
	return t.backend.DisableRawMode()
}

func (t *Terminal) Draw(widget Widget) (CompletedFrame, error) {
	frame := t.GetFrame()

	frame.RenderWidget(widget, frame.Area())

	if err := t.Flush(); err != nil {
		return CompletedFrame{}, fmt.Errorf("flush: %w", err)
	}

	if frame.cursorPosition == nil {
		if err := t.HideCursor(); err != nil {
			return CompletedFrame{}, fmt.Errorf("hide cursor: %w", err)
		}
	} else {
		if err := t.ShowCursor(); err != nil {
			return CompletedFrame{}, fmt.Errorf("show cursor: %w", err)
		}

		if err := t.SetCursorPosition(*frame.cursorPosition); err != nil {
			return CompletedFrame{}, fmt.Errorf("set cursor position: %w", err)
		}
	}

	t.SwapBuffers()

	if err := t.backend.Flush(); err != nil {
		return CompletedFrame{}, fmt.Errorf("backend flush: %w", err)
	}

	completedFrame := CompletedFrame{
		buffer: t.PreviousBuffer(),
		area:   t.lastKnownArea,
		count:  t.frameCount,
	}

	t.frameCount++

	return completedFrame, nil
}

func (t *Terminal) SetCursorPosition(position Position) error {
	if err := t.backend.SetCursorPosition(position); err != nil {
		return err
	}

	t.lastKnownCursorPos = position
	return nil
}

func (t *Terminal) HideCursor() error {
	if err := t.backend.HideCursor(); err != nil {
		return err
	}

	t.hiddenCursor = true
	return nil
}

func (t *Terminal) ShowCursor() error {
	if err := t.backend.ShowCursor(); err != nil {
		return err
	}

	t.hiddenCursor = false
	return nil
}

// Flush obtains a difference between the previous and the current buffer and passes it to the
// current backend for drawing.
func (t *Terminal) Flush() error {
	previous := t.PreviousBuffer()
	current := t.CurrentBuffer()

	updates := previous.Diff(current)

	if len(updates) > 0 {
		last := updates[len(updates)-1]

		t.lastKnownCursorPos = last.Position
	}

	if err := t.backend.Draw(updates); err != nil {
		return fmt.Errorf("draw: %w", err)
	}

	return nil
}

func (t *Terminal) GetFrame() Frame {
	return Frame{
		cursorPosition: nil,
		viewportArea:   t.viewportArea,
		buffer:         t.CurrentBuffer(),
		count:          t.frameCount,
	}
}

func (t *Terminal) CurrentBuffer() *Buffer {
	return &t.buffers[t.current]
}

func (t *Terminal) PreviousBuffer() *Buffer {
	return &t.buffers[1-t.current]
}

func (t *Terminal) Resize(area Rect) error {
	var nextArea Rect

	switch t.viewport.(type) {
	case ViewportInline:
		panic("unimplemented")
	case ViewportFullscreen, ViewportFixed:
		nextArea = area
	}

	t.setViewportArea(nextArea)

	if err := t.Clear(); err != nil {
		return fmt.Errorf("clear: %w", err)
	}

	t.lastKnownArea = area

	return nil
}

func (t *Terminal) Clear() error {
	switch t.viewport.(type) {
	case ViewportFullscreen:
		if err := t.backend.ClearAll(); err != nil {
			return fmt.Errorf("clear all: %w", err)
		}
	case ViewportInline:
		if err := t.backend.SetCursorPosition(t.viewportArea.Position); err != nil {
			return fmt.Errorf("set cursor position: %w", err)
		}

		if err := t.backend.ClearAfterCursor(); err != nil {
			return fmt.Errorf("clear after cursor: %w", err)
		}
	case ViewportFixed:
		area := t.viewportArea

		for y := area.Top(); y < area.Bottom(); y++ {
			if err := t.backend.SetCursorPosition(Position{X: 0, Y: y}); err != nil {
				return fmt.Errorf("set cursor position: %w", err)
			}

			if err := t.backend.ClearAfterCursor(); err != nil {
				return fmt.Errorf("clear after cursor: %w", err)
			}
		}
	}

	t.PreviousBuffer().Reset()

	return nil
}

func (t *Terminal) SwapBuffers() {
	t.PreviousBuffer().Reset()
	t.current = 1 - t.current
}

func (t *Terminal) Size() (Size, bool, error) {
	return t.backend.GetSize()
}

func (t *Terminal) setViewportArea(area Rect) {
	t.CurrentBuffer().Resize(area)
	t.PreviousBuffer().Resize(area)

	t.viewportArea = area
}
