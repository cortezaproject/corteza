package gig

import (
	"context"
	"fmt"
)

func (w *workerAttachment) noop(_ context.Context, _ preprocessorNoop) error {
	return nil
}

func (w *workerAttachment) attachmentRemove(ctx context.Context, t preprocessorAttachmentRemove) error {
	out := make([]Source, 0, len(w.sources))
	for _, src := range w.sources {
		if !compareMeme(t.mimeType, src.MimeType()) {
			out = append(out, src)
		}
	}

	w.sources = out
	return nil
}

func (w *workerAttachment) attachmentTransform(ctx context.Context, t preprocessorAttachmentTransform) error {
	return fmt.Errorf("preprocessor not implemented: %s", PreprocessorHandleAttachmentTransform)
}
