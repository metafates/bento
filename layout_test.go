package bento_test

import (
	"strings"
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
	"github.com/stretchr/testify/require"
)

func TestLength(t *testing.T) {
	letters(
		t,
		bento.FlexStart,
		[]bento.Constraint{
			bento.ConstraintLength(0),
		},
		1,
		"a",
	)
}

func letters(t *testing.T, flex bento.Flex, constraints []bento.Constraint, width int, expected string) {
	area := bento.Rect{Width: width, Height: 1}

	layout := bento.Layout{
		Direction:   bento.DirectionHorizontal,
		Constraints: constraints,
		Flex:        flex,
		Spacing:     bento.SpacingSpace(0),
	}.Split(area)

	buffer := bento.NewBufferEmpty(area)

	latin := []rune("abcdefghijklmnopqrstuvwxyz")

	for i := 0; i < min(len(constraints), len(layout)); i++ {
		c := latin[i]
		area := layout[i]

		s := strings.Repeat(string(c), area.Width)

		textwidget.NewText(textwidget.NewLines(s)...).Render(area, buffer)
	}

	want := textwidget.NewLines(expected).NewBuffer()

	require.Equal(t, want, buffer)
}
