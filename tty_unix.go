//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || aix || zos

package bento

import (
	"fmt"
	"os"
)

func openTTY() (*os.File, error) {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return nil, fmt.Errorf("open tty: %w", err)
	}

	return f, nil
}
