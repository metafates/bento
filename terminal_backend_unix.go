//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || aix || zos

package bento

import (
	"github.com/charmbracelet/x/term"
)

func (d *DefaultBackend) EnableRawMode() error {
	if f, ok := d.input.(term.File); ok && term.IsTerminal(f.Fd()) {
		state, err := term.MakeRaw(f.Fd())
		if err != nil {
			return err
		}

		d.prevInputState = state
	}

	return nil
}
