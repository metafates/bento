package ansi

import (
	"fmt"

	"golang.org/x/term"
)

type State = *term.State

func EnableRawMode(fd int) (State, error) {
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("make raw: %w", err)
	}

	return state, nil
}

func Restore(fd int, state State) error {
	return term.Restore(fd, state)
}
