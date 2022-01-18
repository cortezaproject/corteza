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
	WorkerHandleNoop       = "noop"
	WorkerHandleAttachment = "attachment"
	WorkerHandleImport     = "import"
	WorkerHandleExport     = "export"
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

func (w *workerImport) Preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
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

func (w *workerExport) Preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
	for _, t := range tasks {
		switch tc := t.(type) {
		case preprocessorExperimentalExport:
			err = w.experimentalExport(ctx, tc)
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

// ------------------------------------------------------------------------
// Worker registry

func workerDefinitions() WorkerDefSet {
	return WorkerDefSet{
		{
			Ref:         WorkerHandleNoop,
			Description: "Noop worker has no predefined operations -- it proxies decoder results into postprocessor input.",
		},
		{
			Ref:         WorkerHandleAttachment,
			Description: "@todo not implemented.",
		},
		{
			Ref:         WorkerHandleImport,
			Description: "Import worker is used to import external data into Corteza.",
		},
		{
			Ref:         WorkerHandleExport,
			Description: "Export worker is used to export internal data into a predefined format.",
		},
	}
}
