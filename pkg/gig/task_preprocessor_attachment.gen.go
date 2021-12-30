package gig

// @todo generate the file; this is a placeholder until qlang support is added

import "context"

type (
	preprocessorAttachmentRemove struct {
		MimeType string `json:"mimeType"`
	}

	preprocessorAttachmentTransform struct {
		Width  uint `json:"width"`
		Height uint `json:"height"`
	}
)

var (
	AttachmentWorkerName = "attachment"

	PreprocessorHandleAttachmentRemove    preprocessor = "attachmentRemove"
	PreprocessorHandleAttachmentTransform preprocessor = "attachmentTransform"
)

// Utilities

func (w *workerAttachment) preprocess(ctx context.Context, tasks ...Preprocessor) (err error) {
	for _, t := range tasks {
		switch tc := t.(type) {
		case preprocessorAttachmentRemove:
			err = w.attachmentRemove(ctx, tc)
		case preprocessorAttachmentTransform:
			err = w.attachmentTransform(ctx, tc)
		}

		if err != nil {
			return
		}
	}

	return nil
}

// Preprocessors

func (t preprocessorAttachmentRemove) Ref() preprocessor {
	return PreprocessorHandleAttachmentRemove
}

func (t preprocessorAttachmentRemove) Worker() []string {
	return []string{AttachmentWorkerName}
}

func (t preprocessorAttachmentRemove) Params() interface{} {
	return t
}

// ...

func (t preprocessorAttachmentTransform) Ref() preprocessor {
	return PreprocessorHandleAttachmentTransform
}

func (t preprocessorAttachmentTransform) Worker() []string {
	return []string{AttachmentWorkerName}
}

func (t preprocessorAttachmentTransform) Params() interface{} {
	return t
}
