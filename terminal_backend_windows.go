//go:build windows

package bento

import (
	"fmt"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/windows"
)

func (d *DefaultBackend) EnableRawMode() error {
	// Save stdin state and enable VT input
	// We also need to enable VT
	// input here.
	if f, ok := d.input.(term.File); ok && term.IsTerminal(f.Fd()) {
		state, err := term.MakeRaw(f.Fd())
		if err != nil {
			return err
		}

		d.prevInputState = state

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
	if f, ok := d.output.(term.File); ok && term.IsTerminal(f.Fd()) {
		state, err := term.GetState(f.Fd())
		if err != nil {
			return err
		}

		d.prevOutputState = state

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
