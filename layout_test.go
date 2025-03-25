package bento_test

import (
	"strings"
	"testing"

	"github.com/metafates/bento"
	"github.com/metafates/bento/widget/textwidget"
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
			Constraints: []bento.Constraint{bento.ConstraintLen(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(3)},
			Want:        "aa",
		},
		{
			Name:        "width 1 zero zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(0)},
			Want:        "b",
		},
		{
			Name:        "width 1 zero exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(1)},
			Want:        "b",
		},
		{
			Name:        "width 1 zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(2)},
			Want:        "b",
		},
		{
			Name:        "width 1 exact zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(2)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(0)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(1)},
			Want:        "a",
		},
		{
			Name:        "width 1 overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       1,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(2)},
			Want:        "a",
		},
		{
			Name:        "width 2 zero zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(0)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(1)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(2)},
			Want:        "bb",
		},
		{
			Name:        "width 2 zero overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(0), bento.ConstraintLen(3)},
			Want:        "bb",
		},
		{
			Name:        "width 2 underflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(0)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(1)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(2)},
			Want:        "ab",
		},
		{
			Name:        "width 2 underflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(1), bento.ConstraintLen(3)},
			Want:        "ab",
		},
		{
			Name:        "width 2 exact zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 exact overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(3)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow zero",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(3), bento.ConstraintLen(0)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow underflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(3), bento.ConstraintLen(1)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow exact",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(3), bento.ConstraintLen(2)},
			Want:        "aa",
		},
		{
			Name:        "width 2 overflow overflow",
			Flex:        bento.FlexLegacy,
			Width:       2,
			Constraints: []bento.Constraint{bento.ConstraintLen(3), bento.ConstraintLen(3)},
			Want:        "aa",
		},
		{
			Name:        "width 3 with stretch last",
			Flex:        bento.FlexLegacy,
			Width:       3,
			Constraints: []bento.Constraint{bento.ConstraintLen(2), bento.ConstraintLen(2)},
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

func TestPercentageFlexStart(t *testing.T) {
	testCases := []LayoutSplitTestCase{
		{
			Name:        "Flex Start with Percentage 0, 0",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(0)},
			Want:        "          ",
		},
		{
			Name:        "Flex Start with Percentage 0, 25",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(25)},
			Want:        "bbb       ",
		},
		{
			Name:        "Flex Start with Percentage 0, 50",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(50)},
			Want:        "bbbbb     ",
		},
		{
			Name:        "Flex Start with Percentage 0, 100",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(100)},
			Want:        "bbbbbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 0, 200",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(200)},
			Want:        "bbbbbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 10, 0",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(0)},
			Want:        "a         ",
		},
		{
			Name:        "Flex Start with Percentage 10, 25",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(25)},
			Want:        "abbb      ",
		},
		{
			Name:        "Flex Start with Percentage 10, 50",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(50)},
			Want:        "abbbbb    ",
		},
		{
			Name:        "Flex Start with Percentage 10, 100",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(100)},
			Want:        "abbbbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 10, 200",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(200)},
			Want:        "abbbbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 25, 0",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(0)},
			Want:        "aaa       ",
		},
		{
			Name:        "Flex Start with Percentage 25, 25",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(25)},
			Want:        "aaabb     ",
		},
		{
			Name:        "Flex Start with Percentage 25, 50",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(50)},
			Want:        "aaabbbbb  ",
		},
		{
			Name:        "Flex Start with Percentage 25, 100",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(100)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 25, 200",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(200)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 33, 0",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(0)},
			Want:        "aaa       ",
		},
		{
			Name:        "Flex Start with Percentage 33, 25",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(25)},
			Want:        "aaabbb    ",
		},
		{
			Name:        "Flex Start with Percentage 33, 50",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(50)},
			Want:        "aaabbbbb  ",
		},
		{
			Name:        "Flex Start with Percentage 33, 100",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(100)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 33, 200",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(200)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex Start with Percentage 50, 0",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(50), bento.ConstraintPercentage(0)},
			Want:        "aaaaa     ",
		},
		{
			Name:        "Flex Start with Percentage 50, 50",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(50), bento.ConstraintPercentage(50)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex Start with Percentage 50, 100",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(50), bento.ConstraintPercentage(100)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex Start with Percentage 100, 0",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(0)},
			Want:        "aaaaaaaaaa",
		},
		{
			Name:        "Flex Start with Percentage 100, 50",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(50)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex Start with Percentage 100, 100",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(100)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex Start with Percentage 100, 200",
			Flex:        bento.FlexStart,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(200)},
			Want:        "aaaaabbbbb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, tc.Test)
	}
}

func TestPercentageFlexSpaceBetween(t *testing.T) {
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
		{
			Name:        "Flex SpaceBetween with Percentage 0, 50",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(50)},
			Want:        "     bbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 0, 100",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(100)},
			Want:        "bbbbbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 0, 200",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(0), bento.ConstraintPercentage(200)},
			Want:        "bbbbbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 10, 0",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(0)},
			Want:        "a         ",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 10, 25",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(25)},
			Want:        "a       bb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 10, 50",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(50)},
			Want:        "a    bbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 10, 100",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(100)},
			Want:        "abbbbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 10, 200",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(10), bento.ConstraintPercentage(200)},
			Want:        "abbbbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 25, 0",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(0)},
			Want:        "aaa       ",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 25, 25",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(25)},
			Want:        "aaa     bb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 25, 50",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(50)},
			Want:        "aaa  bbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 25, 100",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(100)},
			Want:        "aaabbbbbbb",
		},
		{
			Name: "Flex SpaceBetween with Percentage 25, 200", Flex: bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(25), bento.ConstraintPercentage(200)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 33, 0",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(0)},
			Want:        "aaa       ",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 33, 25",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(25)},
			Want:        "aaa     bb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 33, 50",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(50)},
			Want:        "aaa  bbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 33, 100",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(100)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 33, 200",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(33), bento.ConstraintPercentage(200)},
			Want:        "aaabbbbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 50, 0",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(50), bento.ConstraintPercentage(0)},
			Want:        "aaaaa     ",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 50, 50",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(50), bento.ConstraintPercentage(50)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 50, 100",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(50), bento.ConstraintPercentage(100)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 100, 0",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(0)},
			Want:        "aaaaaaaaaa",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 100, 50",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(50)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 100, 100",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(100)},
			Want:        "aaaaabbbbb",
		},
		{
			Name:        "Flex SpaceBetween with Percentage 100, 200",
			Flex:        bento.FlexSpaceBetween,
			Width:       10,
			Constraints: []bento.Constraint{bento.ConstraintPercentage(100), bento.ConstraintPercentage(200)},
			Want:        "aaaaabbbbb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, tc.Test)
	}
}

func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		constraints []bento.Constraint
		direction   bento.Direction
		split       bento.Rect
		want        []bento.Rect
	}{
		{
			name: "50% 50% min(0) stretches into last",
			constraints: []bento.Constraint{
				bento.ConstraintPercentage(50),
				bento.ConstraintPercentage(50),
				bento.ConstraintMin(0),
			},
			direction: bento.DirectionVertical,
			split:     bento.Rect{Width: 1, Height: 1},
			want: []bento.Rect{
				{Position: bento.Position{X: 0, Y: 0}, Width: 1, Height: 1},
				{Position: bento.Position{X: 0, Y: 1}, Width: 1, Height: 0},
				{Position: bento.Position{X: 0, Y: 1}, Width: 1, Height: 0},
			},
		},
		{
			name: "max(1) 99% min(0) stretches into last",
			constraints: []bento.Constraint{
				bento.ConstraintMax(1),
				bento.ConstraintPercentage(99),
				bento.ConstraintMin(0),
			},
			direction: bento.DirectionVertical,
			split:     bento.Rect{Width: 1, Height: 1},
			want: []bento.Rect{
				{Position: bento.Position{X: 0, Y: 0}, Width: 1, Height: 0},
				{Position: bento.Position{X: 0, Y: 0}, Width: 1, Height: 1},
				{Position: bento.Position{X: 0, Y: 1}, Width: 1, Height: 0},
			},
		},
		{
			name: "min(1) length(0) min(1)",
			constraints: []bento.Constraint{
				bento.ConstraintMin(1),
				bento.ConstraintLen(0),
				bento.ConstraintMin(1),
			},
			direction: bento.DirectionHorizontal,
			split:     bento.Rect{Width: 1, Height: 1},
			want: []bento.Rect{
				{Position: bento.Position{X: 0, Y: 0}, Width: 1, Height: 1},
				{Position: bento.Position{X: 1, Y: 0}, Width: 0, Height: 1},
				{Position: bento.Position{X: 1, Y: 0}, Width: 0, Height: 1},
			},
		},
		{
			name: "stretches the 2nd last length instead of the last min based on ranking",
			constraints: []bento.Constraint{
				bento.ConstraintLen(3),
				bento.ConstraintMin(4),
				bento.ConstraintLen(1),
				bento.ConstraintMin(4),
			},
			direction: bento.DirectionHorizontal,
			split:     bento.Rect{Width: 7, Height: 1},
			want: []bento.Rect{
				{Position: bento.Position{X: 0, Y: 0}, Width: 0, Height: 1},
				{Position: bento.Position{X: 0, Y: 0}, Width: 4, Height: 1},
				{Position: bento.Position{X: 4, Y: 0}, Width: 0, Height: 1},
				{Position: bento.Position{X: 4, Y: 0}, Width: 3, Height: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			layout := bento.Layout{
				Constraints: tc.constraints,
				Direction:   tc.direction,
			}.Split(tc.split)

			require.Equal(t, tc.want, layout)
		})
	}
}

func TestFlexConstraint(t *testing.T) {
	testCases := []struct {
		name        string
		constraints []bento.Constraint
		want        [][]int
		flex        bento.Flex
	}{
		{
			name: "length center",
			constraints: []bento.Constraint{
				bento.ConstraintLen(50),
			},
			want: [][]int{{50, 100}},
			flex: bento.FlexEnd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rect := bento.Rect{
				Width:  100,
				Height: 1,
			}

			rects := bento.NewLayout(tc.constraints...).Horizontal().WithFlex(tc.flex).Split(rect)

			ranges := make([][]int, 0, len(rects))

			for _, r := range rects {
				ranges = append(ranges, []int{r.Left(), r.Right()})
			}

			require.Equal(t, tc.want, ranges)
		})
	}
}

func letters(t *testing.T, flex bento.Flex, constraints []bento.Constraint, width int, expected string) {
	t.Helper()

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

		textwidget.NewText(textwidget.NewLinesStr(s)...).Render(area, &buffer)
	}

	want := textwidget.NewLinesStr(expected).NewBuffer()

	require.Equal(t, want, buffer)
}
