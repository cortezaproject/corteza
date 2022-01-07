package gig

import (
	"context"
)

func (w *workerNoop) noop(_ context.Context, _ preprocessorNoop) error {
	return nil
}
