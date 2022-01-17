package gig

import (
	"context"
)

type (
	WorkerNoopState struct {
		Sources SourceWrapSet
	}

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
			ID:       src.ID(),
			Name:     src.FileName(),
			Mime:     src.MimeType(),
			Size:     src.Size(),
			Checksum: src.Checksum(),
		})
	}

	out := WorkerNoopState{
		Sources: sources,
	}

	return out, nil
}

func (w *workerNoop) Cleanup(context.Context) error {
	return nil
}
