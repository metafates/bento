//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || aix || zos

package bento

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/term"
)

func (a *appRunner) initInput() (err error) {
	if err := a.terminal.EnableRawMode(); err != nil {
		return fmt.Errorf("enter raw mode: %w", err)
	}

	if f, ok := a.terminal.Output().(term.File); ok && term.IsTerminal(f.Fd()) {
		a.ttyOutput = f
	}

	return nil
}

func openInputTTY() (*os.File, error) {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return nil, fmt.Errorf("open tty: %w", err)
	}

	return f, nil
}
