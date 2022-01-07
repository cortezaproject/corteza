package gig

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
)

const (
	WorkerNoopHandle       = "noop"
	WorkerAttachmentHandle = "attachment"
	WorkerEnvoyHandle      = "envoy"
)

func (w *workerNoop) Preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
	for _, t := range tasks {
		switch tc := t.(type) {
		case preprocessorNoop:
			err = w.noop(ctx, tc)
		default:
			err = fmt.Errorf("unknown preprocessor: %s", w.Ref())
		}

		if err != nil {
			return
		}
	}

	return nil
}

func (w *workerAttachment) Preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
	for _, t := range tasks {
		switch tc := t.(type) {
		case preprocessorAttachmentRemove:
			err = w.attachmentRemove(ctx, tc)
		case preprocessorAttachmentTransform:
			err = w.attachmentTransform(ctx, tc)
		case preprocessorNoop:
			err = w.noop(ctx, tc)
		default:
			err = fmt.Errorf("unknown preprocessor: %s", w.Ref())
		}

		if err != nil {
			return
		}
	}

	return nil
}

func (w *workerEnvoy) Preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
	for _, t := range tasks {
		switch tc := t.(type) {
		case preprocessorResourceRemove:
			err = w.resourceRemove(ctx, tc)
		case preprocessorResourceLoad:
			err = w.resourceLoad(ctx, tc)
		case preprocessorNamespaceLoad:
			err = w.namespaceLoad(ctx, tc)
		case preprocessorNoop:
			err = w.noop(ctx, tc)
		default:
			err = fmt.Errorf("unknown preprocessor: %s", w.Ref())
		}

		if err != nil {
			return
		}
	}

	return nil
}
