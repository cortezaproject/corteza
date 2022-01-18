package gig

import (
	"context"
)

func (w *workerImport) noop(_ context.Context, _ preprocessorNoop) error {
	return nil
}
