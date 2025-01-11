package bento

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/muesli/cancelreader"
)

func (a *appRunner) initTerminal() error {
	if err := a.initInput(); err != nil {
		return err
	}

	a.terminal.HideCursor()
	return nil
}

// restoreTerminalState restores the terminal to the state prior to running the
// Bubble Tea program.
func (a *appRunner) restoreTerminalState() error {
	a.terminal.DisableBracketedPaste()
	a.terminal.ShowCursor()
	// a.disableMouse()

	// if p.renderer.reportFocus() {
	// p.renderer.disableReportFocus()
	// }

	// if a.renderer.altScreen() {
	// p.renderer.exitAltScreen()

	// give the terminal a moment to catch up
	// time.Sleep(time.Millisecond * 10) //nolint:gomnd
	// }

	return a.restoreInput()
}

// restoreInput restores the tty input to its original state.
func (a *appRunner) restoreInput() error {
	if a.ttyInput != nil && a.previousInputState != nil {
		if err := term.Restore(a.ttyInput.Fd(), a.previousInputState); err != nil {
			return fmt.Errorf("restore input: %w", err)
		}
	}

	if a.ttyOutput != nil && a.previousOutputState != nil {
		if err := term.Restore(a.ttyOutput.Fd(), a.previousOutputState); err != nil {
			return fmt.Errorf("restore output: %w", err)
		}
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
	if a.ttyOutput == nil {
		// can't query window size
		return
	}

	size, err := a.terminal.Size()
	if err != nil {
		select {
		case <-a.ctx.Done():
		case a.errs <- err:
		}

		return
	}

	a.Send(WindowSizeMsg(size))
}
