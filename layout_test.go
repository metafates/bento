package bento_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
	"github.com/stretchr/testify/require"
)

func TestLengthSimple(t *testing.T) {
	letters(t, bento.FlexLegacy, []bento.Constraint{bento.ConstraintLength(0)}, 1, "a")
}

func TestLength(t *testing.T) {
	testCases := []struct {
		Flex        bento.Flex
		Width       int
		Constraints []bento.Constraint
		Want        string
	}{
		{
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(0)},
			Want:        "          ",
		},
		{
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(25)},
			Want:        "bbb       ",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case #%d", i+1), func(t *testing.T) {
			letters(
				t,
				tc.Flex,
				tc.Constraints,
				tc.Width,
				tc.Want,
			)
		})
	}
}

func letters(t *testing.T, flex bento.Flex, constraints []bento.Constraint, width int, expected string) {
	area := bento.Rect{Width: width, Height: 1}

	layout := bento.Layout{
		Direction:   bento.DirectionHorizontal,
		Constraints: constraints,
		Flex:        flex,
		Spacing:     bento.SpacingSpace(0),
	}.Split(area)

	fmt.Printf("layout: %+v\n", layout)

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
