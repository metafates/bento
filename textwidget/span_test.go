package textwidget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnicodeTruncateStart(t *testing.T) {
	testCases := []struct {
		Name   string
		String string
		Width  int

		WantString string
		WantLen    int
	}{
		{
			Name:   "empty",
			String: "",
			Width:  4,

			WantString: "",
			WantLen:    0,
		},
		{
			Name:   "zero width simple",
			String: "ab",
			Width:  0,

			WantString: "",
			WantLen:    0,
		},
		{
			Name:   "zero width",
			String: "你好",
			Width:  0,

			WantString: "",
			WantLen:    0,
		},
		{
			Name:   "less than limit simple",
			String: "abc",
			Width:  4,

			WantString: "abc",
			WantLen:    3,
		},
		{
			Name:   "less than limit",
			String: "你",
			Width:  4,

			WantString: "你",
			WantLen:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s, l := unicodeTruncateStart(tc.String, tc.Width)

			require.Equal(t, tc.WantString, s)
			require.Equal(t, tc.WantLen, l)
		})
	}
}
