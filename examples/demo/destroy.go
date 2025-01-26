package main

import (
	"math"
	"math/rand/v2"

	"github.com/metafates/bento"
)

func destroy(frame int, buffer *bento.Buffer) {
	drip(frame, buffer.Area(), buffer)
}

var randGen = rand.New(rand.NewChaCha8([32]byte{
	24, 220, 78, 229, 154, 40, 116, 47, 166, 40, 45, 197, 179, 3, 228, 157, 21, 234, 125, 156, 72, 0, 186, 102, 202, 213, 231, 116, 33, 190, 108, 10,
}))

func drip(frame int, area bento.Rect, buffer *bento.Buffer) {
	const (
		rampFrames = 450
		dripSpeed  = 500
	)

	fractionalSpeed := float64(frame) / float64(rampFrames)
	variableSpeed := float64(dripSpeed) * fractionalSpeed * fractionalSpeed * fractionalSpeed
	pixelCount := int(math.Floor(float64(frame) * variableSpeed))

	for i := 0; i < pixelCount; i++ {
		srcX := randRange(0, area.Width)
		srcY := randRange(1, area.Height-2)

		src := *buffer.CellAt(bento.NewPosition(srcX, srcY))

		if randRatio(1, 100) {
			destX := randRange(max(0, srcX-5), srcX+5)
			destX = max(area.Left(), min(area.Right()-1, destX))

			destY := area.Top() + 1

			dest := buffer.CellAt(bento.NewPosition(destX, destY))

			if randRatio(1, 10) {
				*dest = src
			} else {
				dest.Reset()
			}
		} else {
			destX := srcX
			destY := min(area.Bottom()-2, srcY+1)

			*buffer.CellAt(bento.NewPosition(destX, destY)) = src
		}
	}
}

func randRange(from, to int) int {
	return randGen.IntN(to-from) + from
}

func randRatio(num, den int) bool {
	if num == den {
		return true
	}

	return randGen.Float64() <= float64(num)/float64(den)
}
