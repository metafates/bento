//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || aix || zos

package bento

import (
	"fmt"
	"os"
)

func (a *appRunner) initInput() (err error) {
	if err := a.terminal.EnableRawMode(); err != nil {
		return fmt.Errorf("enter raw mode: %w", err)
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
