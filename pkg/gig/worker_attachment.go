package gig

import (
	"context"
	"fmt"
	"mime"
	"strings"
)

type (
	workerAttachmentState struct {
		Attachments []SourceWrap
	}

	workerAttachment struct {
		sources SourceSet
	}
)

func WorkerAttachment() Worker {
	return &workerAttachment{}
}

func (w *workerAttachment) Ref() string {
	return WorkerHandleAttachment
}

func (w *workerAttachment) Prepare(_ context.Context, sources ...Source) error {
	return fmt.Errorf("worker not implemented: %s", WorkerHandleAttachment)

	// No need to do anything special here here
	w.sources = sources
	return nil
}

func (w *workerAttachment) Exec(ctx context.Context) (output SourceSet, meta WorkMeta, err error) {
	// @todo processing and stuff
	output = w.sources
	meta = w.collectMeta()

	return
}

// @todo ...
func (w *workerAttachment) collectMeta() (meta WorkMeta) {
	meta = make(WorkMeta)
	return
}

func (w *workerAttachment) State(context.Context) (WorkerState, error) {
	out := workerAttachmentState{
		Attachments: make([]SourceWrap, len(w.sources)),
	}

	for i, src := range w.sources {
		out.Attachments[i] = SourceWrap{
			Name:     src.Name(),
			Size:     src.Size(),
			Mime:     src.MimeType(),
			Checksum: src.Checksum(),
		}
	}

	return out, nil
}

func (w *workerAttachment) Cleanup(context.Context) error {
	return nil
}

// ---

// @todo use something built in
func compareMeme(rawA, rawB string) bool {
	a, _, err := mime.ParseMediaType(rawA)
	if err != nil {
		return false
	}
	b, _, err := mime.ParseMediaType(rawB)
	if err != nil {
		return false
	}

	if a == b {
		return true
	}

	pa := strings.Split(a, "/")
	pb := strings.Split(b, "/")

	return pa[1] == "*" && pa[0] == pb[0]
}
