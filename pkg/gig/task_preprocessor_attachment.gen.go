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

	PreprocessorHandleAttachmentRemove    = "attachmentRemove"
	PreprocessorHandleAttachmentTransform = "attachmentTransform"
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

func (t preprocessorAttachmentRemove) Ref() string {
	return PreprocessorHandleAttachmentRemove
}

func (t preprocessorAttachmentRemove) Worker() []string {
	return []string{AttachmentWorkerName}
}

func (t preprocessorAttachmentRemove) Params() map[string]interface{} {
	return map[string]interface{}{
		"mimeType": t.MimeType,
	}
}

// ...

func (t preprocessorAttachmentTransform) Ref() string {
	return PreprocessorHandleAttachmentTransform
}

func (t preprocessorAttachmentTransform) Worker() []string {
	return []string{AttachmentWorkerName}
}

func (t preprocessorAttachmentTransform) Params() map[string]interface{} {
	return map[string]interface{}{
		"width":  t.Width,
		"height": t.Height,
	}
}
