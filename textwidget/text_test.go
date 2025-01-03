package textwidget

import (
	"testing"

	"github.com/metafates/bento"
	"github.com/stretchr/testify/require"
)

func TestText_Render(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		text := NewText("foo")
		area := bento.Rect{
			X:      0,
			Y:      0,
			Width:  5,
			Height: 1,
		}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines(*NewLine("foo  ")).NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("out of bounds", func(t *testing.T) {
		smallBuffer := bento.NewBufferEmpty(bento.Rect{
			Width:  10,
			Height: 1,
		})

		outOfBoundsArea := bento.Rect{
			X:      20,
			Y:      20,
			Width:  10,
			Height: 1,
		}

		NewText("Hello, World!").Render(outOfBoundsArea, smallBuffer)

		require.Equal(t, bento.NewBufferEmpty(smallBuffer.Area), smallBuffer)
	})

	t.Run("right aligned", func(t *testing.T) {
		text := NewText("foo")
		text.Alignment = bento.AlignmentRight

		area := bento.Rect{
			Width:  5,
			Height: 1,
		}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines(*NewLine("  foo")).NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("center aligned", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			width int
			text  string
			want  string
		}{
			{
				name:  "odd",
				width: 5,
				text:  "foo",
				want:  " foo ",
			},
			{
				name:  "even",
				width: 6,
				text:  "foo",
				want:  " foo  ",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				text := NewText(tc.text)
				text.Alignment = bento.AlignmentCenter

				area := bento.Rect{
					Width:  tc.width,
					Height: 1,
				}

				buffer := bento.NewBufferEmpty(area)

				text.Render(area, buffer)

				want := NewLines(*NewLine(tc.want)).NewBuffer()

				require.Equal(t, want, buffer)
			})
		}
	})

	t.Run("right aligned with truncation", func(t *testing.T) {
		text := NewText("123456789")
		text.Alignment = bento.AlignmentRight

		area := bento.Rect{Width: 5, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines(*NewLine("56789")).NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("center aligned odd with truncation", func(t *testing.T) {
		text := NewText("123456789")
		text.Alignment = bento.AlignmentCenter

		area := bento.Rect{Width: 5, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines(*NewLine("34567")).NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("center aligned even with truncation", func(t *testing.T) {
		text := NewText("123456789")
		text.Alignment = bento.AlignmentCenter

		area := bento.Rect{Width: 6, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines(*NewLine("234567")).NewBuffer()

		require.Equal(t, want, buffer)
	})
}
