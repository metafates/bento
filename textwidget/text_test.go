package textwidget

import (
	"testing"

	"github.com/metafates/bento"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/require"
)

func TestText_Render(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		text := NewText(NewLines("foo"))
		area := bento.Rect{
			X:      0,
			Y:      0,
			Width:  5,
			Height: 1,
		}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines("foo  ").NewBuffer()

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

		NewText(NewLines("Hello, World!")).Render(outOfBoundsArea, smallBuffer)

		require.Equal(t, bento.NewBufferEmpty(smallBuffer.Area), smallBuffer)
	})

	t.Run("right aligned", func(t *testing.T) {
		text := NewText(NewLines("foo"))
		text.Alignment = bento.AlignmentRight

		area := bento.Rect{
			Width:  5,
			Height: 1,
		}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines("  foo").NewBuffer()

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
				text := NewText(NewLines(tc.text))
				text.Alignment = bento.AlignmentCenter

				area := bento.Rect{
					Width:  tc.width,
					Height: 1,
				}

				buffer := bento.NewBufferEmpty(area)

				text.Render(area, buffer)

				want := NewLines(tc.want).NewBuffer()

				require.Equal(t, want, buffer)
			})
		}
	})

	t.Run("right aligned with truncation", func(t *testing.T) {
		text := NewText(NewLines("123456789"))
		text.Alignment = bento.AlignmentRight

		area := bento.Rect{Width: 5, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines("56789").NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("center aligned odd with truncation", func(t *testing.T) {
		text := NewText(NewLines("123456789"))
		text.Alignment = bento.AlignmentCenter

		area := bento.Rect{Width: 5, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines("34567").NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("center aligned even with truncation", func(t *testing.T) {
		text := NewText(NewLines("123456789"))
		text.Alignment = bento.AlignmentCenter

		area := bento.Rect{Width: 6, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines("234567").NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("one line right", func(t *testing.T) {
		line1 := NewLine("foo")
		line2 := NewLine("bar")
		line2.Alignment = bento.AlignmentCenter

		text := NewText([]Line{line1, line2})
		text.Alignment = bento.AlignmentRight

		area := bento.Rect{Width: 5, Height: 2}
		buffer := bento.NewBufferEmpty(area)

		text.Render(area, buffer)

		want := NewLines("  foo", " bar ").NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("only styles line area", func(t *testing.T) {
		area := bento.Rect{Width: 5, Height: 1}

		buffer := bento.NewBufferEmpty(area)

		line := NewLine("foo").WithStyle(bento.NewStyle().WithBackground(termenv.ANSIBlue))

		NewText([]Line{line}).Render(area, buffer)

		want := NewLines("foo  ").NewBuffer()
		want.SetStyle(bento.Rect{Width: 3, Height: 1}, bento.NewStyle().WithBackground(termenv.ANSIBlue))

		require.Equal(t, want, buffer)
	})

	t.Run("truncate", func(t *testing.T) {
		buffer := bento.NewBufferEmpty(bento.Rect{Width: 6, Height: 1})

		NewText([]Line{
			NewLine("foobar").WithStyle(bento.NewStyle().WithBackground(termenv.ANSIBlue)),
		}).Render(bento.Rect{Width: 3, Height: 1}, buffer)

		want := NewLines("foo   ").NewBuffer()
		want.SetStyle(bento.Rect{Width: 3, Height: 1}, bento.NewStyle().WithBackground(termenv.ANSIBlue))

		require.Equal(t, want, buffer)
	})
}
