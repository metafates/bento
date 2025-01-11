package bento

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/muesli/cancelreader"
)

func (a *appRunner) initTerminal() error {
	if err := a.terminal.EnableRawMode(); err != nil {
		return fmt.Errorf("enable raw mode: %w", err)
	}

	if err := a.terminal.HideCursor(); err != nil {
		return fmt.Errorf("hide cursor: %w", err)
	}

	return nil
}

// restoreTerminal restores the terminal to the state prior to running the
// Bento app.
func (a *appRunner) restoreTerminal() error {
	if err := a.terminal.DisableBracketedPaste(); err != nil {
		return fmt.Errorf("disable bracketed paste: %w", err)
	}

	if err := a.terminal.ShowCursor(); err != nil {
		return fmt.Errorf("show cursor: %w", err)
	}
	// a.disableMouse()

	// if p.renderer.reportFocus() {
	// p.renderer.disableReportFocus()
	// }

	// if a.renderer.altScreen() {
	// p.renderer.exitAltScreen()

	// give the terminal a moment to catch up
	// time.Sleep(time.Millisecond * 10) //nolint:gomnd
	// }
	if err := a.terminal.DisableRawMode(); err != nil {
		return fmt.Errorf("disable raw mode: %w", err)
	}

	return nil
}

func (a *appRunner) initCancelReader() error {
	r, err := newInputReader(a.terminal)
	if err != nil {
		return fmt.Errorf("new reader: %w", err)
	}

	a.cancelReader = r
	a.readLoopDone = make(chan struct{})

	go a.readLoop()

	return nil
}

func (a *appRunner) readLoop() {
	defer close(a.readLoopDone)

	err := readInputs(a.ctx, a.msgs, a.cancelReader)
	if !errors.Is(err, io.EOF) && !errors.Is(err, cancelreader.ErrCanceled) {
		select {
		case <-a.ctx.Done():
		case a.errs <- err:
		}
	}
}

// waitForReadLoop waits for the cancelReader to finish its read loop.
func (a *appRunner) waitForReadLoop() {
	select {
	case <-a.readLoopDone:
	case <-time.After(500 * time.Millisecond):
		// The read loop hangs, which means the input
		// cancelReader's cancel function has returned true even
		// though it was not able to cancel the read.
	}
}

// checkResize detects the current size of the output and informs the program
// via a WindowSizeMsg.
func (a *appRunner) checkResize() {
	size, ok, err := a.terminal.Size()
	if !ok {
		return
	}

	if err != nil {
		select {
		case <-a.ctx.Done():
		case a.errs <- err:
		}

		return
	}

	a.Send(WindowSizeMsg(size))
}
