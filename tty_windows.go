//go:build windows

package bento

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/windows"
)

func (a *appRunner) initInput() error {
	// Save stdin state and enable VT input
	// We also need to enable VT
	// input here.
	if f, ok := a.terminal.Input().(term.File); ok && term.IsTerminal(f.Fd()) {
		if err := a.terminal.EnableRawMode(); err != nil {
			return fmt.Errorf("enable raw mode: %w", err)
		}

		state, err := term.MakeRaw(f.Fd())
		if err != nil {
			return err
		}

		a.ttyInput = f
		a.previousInputState = state

		// Enable VT input
		var mode uint32
		if err := windows.GetConsoleMode(windows.Handle(f.Fd()), &mode); err != nil {
			return fmt.Errorf("get console mode: %w", err)
		}

		if err := windows.SetConsoleMode(windows.Handle(f.Fd()), mode|windows.ENABLE_VIRTUAL_TERMINAL_INPUT); err != nil {
			return fmt.Errorf("set console mode: %w", err)
		}
	}

	// Save output screen buffer state and enable VT processing.
	if f, ok := a.terminal.Output().(term.File); ok && term.IsTerminal(f.Fd()) {
		state, err := term.GetState(f.Fd())
		if err != nil {
			return err
		}

		a.ttyOutput = f
		a.previousOutputState = state

		var mode uint32
		if err := windows.GetConsoleMode(windows.Handle(f.Fd()), &mode); err != nil {
			return fmt.Errorf("error getting console mode: %w", err)
		}

		if err := windows.SetConsoleMode(windows.Handle(f.Fd()), mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
			return fmt.Errorf("error setting console mode: %w", err)
		}
	}

	return nil
}

// Open the Windows equivalent of a TTY.
func openInputTTY() (*os.File, error) {
	f, err := os.OpenFile("CONIN$", os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return f, nil
}
