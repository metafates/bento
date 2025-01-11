package ansi

import (
	"fmt"

	"github.com/charmbracelet/x/term"
)

type State = *term.State

func EnableRawMode(fd uintptr) (State, error) {
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("make raw: %w", err)
	}

	return state, nil
}

func Restore(fd uintptr, state State) error {
	return term.Restore(fd, state)
}
