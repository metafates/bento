//go:build !windows

package bento

import (
	"context"
	"io"
)

func readInputs(ctx context.Context, msgs chan<- Msg, input io.Reader) error {
	return readAnsiInputs(ctx, msgs, input)
}
