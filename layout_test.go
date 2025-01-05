package bento_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/textwidget"
	"github.com/stretchr/testify/require"
)

type LayoutSplitTestCase struct {
	Name        string
	Flex        bento.Flex
	Width       int
	Constraints []bento.Constraint
	Want        string
}

func (tc LayoutSplitTestCase) Test(t *testing.T) {
	letters(t, tc.Flex, tc.Constraints, tc.Width, tc.Want)
}

func TestLength(t *testing.T) {
	testCases := []LayoutSplitTestCase{
		{
			Name:        "width 1 zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(3)},
			Want:        "aa",
		},
		{
			Name:        "width 1 zero zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(0)},
			Want:        "b",
		},
		{
			Name:        "width 1 zero exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(1)},
			Want:        "b",
		},
		{
			Name:        "width 1 zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(2)},
			Want:        "b",
		},
		{
			Name:        "width 1 exact zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(2)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 zero zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(0)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(1)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(2)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(0), bento.ConstraintLength(3)},
			Want:        "bb",
		},
		{
			Name:        "width 2 underflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(0)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(1)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(2)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(1), bento.ConstraintLength(3)},
			Want:        "ab",
		},
		{
			Name:        "width 2 exact zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(3)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(3), bento.ConstraintLength(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(3), bento.ConstraintLength(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(3), bento.ConstraintLength(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLength(3), bento.ConstraintLength(3)},
			Want:        "aa",
		},
		{
			Name:        "width 3 with stretch last",
			Flex:        bento.FlexLegacy,
			Width:       3,
			Constraints: []bento.Constraint{bento.ConstraintLength(2), bento.ConstraintLength(2)},
			Want:        "aab",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, tc.Test)
	}
}

func TestMax(t *testing.T) {
	testCases := []LayoutSplitTestCase{
		{
			Name:        "width 1 zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(3)},
			Want:        "aa",
		},
		{
			Name:        "width 1 zero zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(0)},
			Want:        "b",
		},
		{
			Name:        "width 1 zero exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(1)},
			Want:        "b",
		},
		{
			Name:        "width 1 zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(2)},
			Want:        "b",
		},
		{
			Name:        "width 1 exact zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(2)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 zero zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(0)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(1)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(2)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(0), bento.ConstraintMax(3)},
			Want:        "bb",
		},
		{
			Name:        "width 2 underflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(0)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(1)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(2)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(1), bento.ConstraintMax(3)},
			Want:        "ab",
		},
		{
			Name:        "width 2 exact zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(3)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(3), bento.ConstraintMax(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(3), bento.ConstraintMax(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(3), bento.ConstraintMax(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMax(3), bento.ConstraintMax(3)},
			Want:        "aa",
		},
		{
			Name:        "width 3 with stretch last",
			Flex:        bento.FlexLegacy,
			Width:       3,
			Constraints: []bento.Constraint{bento.ConstraintMax(2), bento.ConstraintMax(2)},
			Want:        "aab",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, tc.Test)
	}
}

func TestMin(t *testing.T) {
	testCases := []LayoutSplitTestCase{
		{
			Name:        "width 1 min zero zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(0)},
			Want:        "b",
		},
		{
			Name:        "width 1 min zero exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(1)},
			Want:        "b",
		},
		{
			Name:        "width 1 min zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(2)},
			Want:        "b",
		},
		{
			Name:        "width 1 min exact zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 min exact exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 min exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(2)},
			Want:        "a",
		},
		{
			Name:        "width 1 min overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 min overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 min overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 min zero zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(0)},
			Want:        "bb",
		},
		{
			Name:        "width 2 min zero underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(1)},
			Want:        "bb",
		},
		{
			Name:        "width 2 min zero exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(2)},
			Want:        "bb",
		},
		{
			Name:        "width 2 min zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(0), bento.ConstraintMin(3)},
			Want:        "bb",
		},
		{
			Name:        "width 2 min underflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(0)},
			Want:        "ab",
		},
		{
			Name:        "width 2 min underflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(1)},
			Want:        "ab",
		},
		{
			Name:        "width 2 min underflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(2)},
			Want:        "ab",
		},
		{
			Name:        "width 2 min underflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(1), bento.ConstraintMin(3)},
			Want:        "ab",
		},
		{
			Name:        "width 2 min exact zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min exact underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min exact exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(3)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(3), bento.ConstraintMin(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min overflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(3), bento.ConstraintMin(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(3), bento.ConstraintMin(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 min overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintMin(3), bento.ConstraintMin(3)},
			Want:        "aa",
		},
		{
			Name:        "width 3 min with stretch last",
			Flex:        bento.FlexLegacy,
			Width:       3,
			Constraints: []bento.Constraint{bento.ConstraintMin(2), bento.ConstraintMin(2)},
			Want:        "aab",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, tc.Test)
	}
}

func TestPercentage(t *testing.T) {
	testCases := []LayoutSplitTestCase{
		{
			Name:        "Flex SpaceBetween with Percentage 0, 0",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(0)},
			Want:        "          ",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 0, 25",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(25)},
			Want:        "        bb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, tc.Test)
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
