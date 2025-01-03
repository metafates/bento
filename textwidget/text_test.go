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

		want := Lines([]Line{*NewLine("foo  ")}).NewBuffer()

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

		want := Lines([]Line{*NewLine("  foo")}).NewBuffer()

		require.Equal(t, want, buffer)
	})

	t.Run("center aligned odd", func(t *testing.T) {
		text := NewText("foo")
		text.Alignment = bento.AlignmentRight

		area := bento.Rect{
			Width:  5,
			Height: 1,
		}

		buffer := bento.NewBufferEmpty(area)

		want := Lines([]Line{*NewLine(" foo ")})

		require.Equal(t, want, buffer)
	})
}
