package gig

import (
	"context"
)

type (
	WorkerNoopState map[string]interface{}

	workerNoop struct {
		sources SourceSet
	}
)

var (
	WorkerHandleNoop = "noop"
)

func WorkerNoop() Worker {
	return &workerNoop{}
}

func (w *workerNoop) Ref() string {
	return WorkerHandleNoop
}

func (w *workerNoop) Prepare(_ context.Context, sources ...Source) error {
	w.sources = SourceSet(sources)
	return nil
}

func (w *workerNoop) Exec(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	output = w.sources
	meta = w.collectMeta()

	return
}

func (w *workerNoop) collectMeta() (meta WorkMeta) {
	meta = make(WorkMeta)

	return
}

func (w *workerNoop) State(context.Context) (WorkerState, error) {
	sources := make([]SourceWrap, 0, len(w.sources))
	for _, src := range w.sources {
		sources = append(sources, SourceWrap{
			Name:     src.FileName(),
			Size:     src.Size(),
			Mime:     src.MimeType(),
			Checksum: src.Checksum(),
		})
	}

	out := WorkerNoopState{
		"sources": sources,
	}

	return out, nil
}

func (w *workerNoop) Cleanup(context.Context) error {
	return nil
}
