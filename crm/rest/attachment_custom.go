package rest

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/crusttech/crust/crm/internal/service"
	"github.com/crusttech/crust/crm/rest/handlers"
	"github.com/crusttech/crust/crm/types"
)

type (
	attachmentPayload struct {
		ID         uint64      `json:"attachmentID,string"`
		OwnerID    uint64      `json:"ownerID,string"`
		Url        string      `json:"url"`
		PreviewUrl string      `json:"previewUrl,omitempty"`
		Meta       interface{} `json:"meta"`
		Name       string      `json:"name"`
		CreatedAt  time.Time   `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time  `json:"updatedAt,omitempty"`
	}

	file struct {
		*types.Attachment
		content  io.ReadSeeker
		download bool
	}
)

func (f *file) Download() bool {
	return f.download
}

func (f *file) Name() string {
	return f.Attachment.Name
}

func (f *file) ModTime() time.Time {
	return f.Attachment.CreatedAt
}

func (f *file) Content() io.ReadSeeker {
	return f.content
}

func (f *file) Valid() bool {
	return f.content != nil
}

func loadAttachedFile(svc service.AttachmentService, ID uint64, preview, download bool) (handlers.Downloadable, error) {
	rval := &file{download: download}

	if att, err := svc.FindByID(ID); err != nil {
		return nil, err
	} else {
		rval.Attachment = att
		if preview {
			rval.content, err = svc.OpenPreview(att)
		} else {
			rval.content, err = svc.OpenOriginal(att)
		}

		if err != nil {
			return nil, err
		}
	}

	return rval, nil
}

func makeAttachmentPayload(a *types.Attachment) *attachmentPayload {
	if a == nil {
		return nil
	}

	var preview string
	var baseURL = fmt.Sprintf("/attachment/%s/%d/", a.Kind, a.ID)

	if a.Meta.Preview != nil {
		var ext = a.Meta.Preview.Extension
		if ext == "" {
			ext = "jpg"
		}

		preview = baseURL + fmt.Sprintf("preview.%s", ext)
	}

	return &attachmentPayload{
		ID:         a.ID,
		OwnerID:    a.OwnerID,
		Url:        baseURL + fmt.Sprintf("original/%s", url.PathEscape(a.Name)),
		PreviewUrl: preview,
		Meta:       a.Meta,
		Name:       a.Name,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

func makeRecordAttachmentSetPayload(aa types.AttachmentSet, meta types.AttachmentFilter, err error) (map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}

	pp := make([]*attachmentPayload, len(aa))
	for i := range aa {
		pp[i] = makeAttachmentPayload(aa[i])
	}

	rval := map[string]interface{}{"meta": meta, "attachments": pp}

	return rval, err
}
