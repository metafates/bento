package ansi

import "github.com/charmbracelet/x/term"

func GetSize(fd uintptr) (width, height int, err error) {
	return term.GetSize(fd)
}

func GetCursorPosition() (column, row int, err error) {
	panic("unimplemented")
}
