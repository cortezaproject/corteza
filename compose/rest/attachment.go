package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/crusttech/crust/compose/internal/service"
	"github.com/crusttech/crust/compose/rest/request"
	"github.com/crusttech/crust/compose/types"
	"github.com/crusttech/crust/internal/auth"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

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

	Attachment struct {
		attachment service.AttachmentService
	}
)

func (Attachment) New() *Attachment {
	return &Attachment{attachment: service.DefaultAttachment}
}

// Attachments returns list of all files attached to records
func (ctrl *Attachment) List(ctx context.Context, r *request.AttachmentList) (interface{}, error) {
	f := types.AttachmentFilter{
		Kind:      r.Kind,
		ModuleID:  r.ModuleID,
		RecordID:  r.RecordID,
		FieldName: r.FieldName,
		// Filter:    r.Filter,
		PerPage: r.PerPage,
		Page:    r.Page,
		// Sort:      r.Sort,
	}

	aa, meta, err := ctrl.attachment.Find(f)
	if err != nil {
		return nil, err
	}

	pp := make([]*attachmentPayload, len(aa))
	for i := range aa {
		pp[i] = makeAttachmentPayload(aa[i], auth.GetIdentityFromContext(ctx).Identity())
	}

	return map[string]interface{}{"meta": meta, "attachments": pp}, nil
}

func (ctrl Attachment) Details(ctx context.Context, r *request.AttachmentDetails) (interface{}, error) {
	if a, err := ctrl.attachment.FindByID(r.AttachmentID); err != nil {
		return nil, err
	} else {
		return makeAttachmentPayload(a, auth.GetIdentityFromContext(ctx).Identity()), nil
	}
}

func (ctrl Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	if err := ctrl.isAccessible(r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(ctx, r.AttachmentID, false, r.Download)
}

func (ctrl *Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	if err := ctrl.isAccessible(r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(ctx, r.AttachmentID, true, false)
}

func (ctrl Attachment) isAccessible(attachmentID, userID uint64, signature string) error {
	if userID == 0 {
		return errors.New("missing or invalid user ID")
	}

	if attachmentID == 0 {
		return errors.New("missing or invalid attachment ID")
	}

	if auth.DefaultSigner.Verify(signature, userID, attachmentID) {
		return errors.New("missing or invalid signature")
	}

	return nil
}

func (ctrl Attachment) serve(ctx context.Context, ID uint64, preview, download bool) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		att, err := ctrl.attachment.With(ctx).FindByID(ID)
		if err != nil {
			// Simplify error handling for now
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var fh io.ReadSeeker

		if preview {
			fh, err = ctrl.attachment.OpenPreview(att)
		} else {
			fh, err = ctrl.attachment.OpenOriginal(att)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := url.QueryEscape(att.Name)

		if download {
			w.Header().Add("Content-Disposition", "attachment; filename="+name)
		} else {
			w.Header().Add("Content-Disposition", "inline; filename="+name)
		}

		http.ServeContent(w, req, name, att.CreatedAt, fh)
	}, nil
}

func makeAttachmentPayload(a *types.Attachment, userID uint64) *attachmentPayload {
	if a == nil {
		return nil
	}

	var (
		signParams = fmt.Sprintf("?sign=%s&userID=%d", auth.DefaultSigner.Sign(userID, a.ID), userID)

		preview string
		baseURL = fmt.Sprintf("/attachment/%s/%d/", a.Kind, a.ID)
	)
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
		Url:        baseURL + fmt.Sprintf("original/%s", url.PathEscape(a.Name)) + signParams,
		PreviewUrl: preview + signParams,
		Meta:       a.Meta,
		Name:       a.Name,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

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
