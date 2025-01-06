package ansi

import (
	"golang.org/x/term"
)

func GetSize(fd int) (width, height int, err error) {
	return term.GetSize(fd)
}

func GetCursorPosition() (column, row int, err error) {
	panic("unimplemented")
}
