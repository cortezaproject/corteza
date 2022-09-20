package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	attachmentPayload struct {
		*types.Attachment
	}

	attachmentSetPayload struct {
		Filter types.AttachmentFilter `json:"filter"`
		Set    []*attachmentPayload   `json:"set"`
	}

	Attachment struct {
		attachment service.AttachmentService
	}
)

func (Attachment) New() *Attachment {
	return &Attachment{
		attachment: service.DefaultAttachment,
	}
}

func (ctrl Attachment) Read(ctx context.Context, r *request.AttachmentRead) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.New("Unauthorized")
	}

	a, err := ctrl.attachment.FindByID(ctx, r.AttachmentID)
	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl Attachment) Delete(ctx context.Context, r *request.AttachmentDelete) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.New("Unauthorized")
	}

	_, err := ctrl.attachment.FindByID(ctx, r.AttachmentID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.attachment.DeleteByID(ctx, r.AttachmentID)
}

func (ctrl Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	if err := ctrl.isAccessible(r.Kind, r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(ctx, r.AttachmentID, false, r.Download)
}

func (ctrl Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	if err := ctrl.isAccessible(r.Kind, r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(ctx, r.AttachmentID, true, false)
}

func (ctrl Attachment) isAccessible(kind string, attachmentID, userID uint64, signature string) error {
	if kind == types.AttachmentKindSettings {
		// Attachments on settings are public
		return nil
	}

	if signature == "" {
		return errors.New("Unauthorized")
	}

	if userID == 0 {
		return errors.New("missing or invalid user ID")
	}

	if attachmentID == 0 {
		return errors.New("missing or invalid attachment ID")
	}

	if !auth.DefaultSigner.Verify(signature, userID, attachmentID) {
		return errors.New("missing or invalid signature")
	}

	return nil
}

func (ctrl Attachment) serve(ctx context.Context, attachmentID uint64, preview, download bool) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		att, err := ctrl.attachment.FindByID(ctx, attachmentID)
		if err != nil {
			// Simplify error handling for now
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var fh io.ReadSeekCloser

		if preview {
			fh, err = ctrl.attachment.OpenPreview(att)
		} else {
			fh, err = ctrl.attachment.OpenOriginal(att)
		}

		defer fh.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := url.QueryEscape(att.Name)

		if download {
			w.Header().Add("Content-Disposition", "attachment; filename="+name)
		} else {
			w.Header().Add("Content-Disposition", "inline; filename="+name)
			w.Header().Add("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
		}

		http.ServeContent(w, req, name, att.CreatedAt, fh)
	}, nil
}

func (ctrl Attachment) makeFilterPayload(ctx context.Context, aa types.AttachmentSet, f types.AttachmentFilter, err error) (*attachmentSetPayload, error) {
	if err != nil {
		return nil, err
	}

	asp := &attachmentSetPayload{Filter: f, Set: make([]*attachmentPayload, len(aa))}

	for i := range aa {
		asp.Set[i], _ = makeAttachmentPayload(ctx, aa[i], nil)
	}

	return asp, nil
}

func makeAttachmentPayload(ctx context.Context, a *types.Attachment, err error) (*attachmentPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	var (
		userID     = auth.GetIdentityFromContext(ctx).Identity()
		signParams = ""

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

	switch a.Kind {
	case types.AttachmentKindSettings:
		// No URL signing, settings attachments are public.
	default:
		signParams = fmt.Sprintf("?sign=%s&userID=%d", auth.DefaultSigner.Sign(userID, a.ID), userID)
	}

	ap := &attachmentPayload{a}

	ap.Url = baseURL + fmt.Sprintf("original/%s", url.PathEscape(a.Name)) + signParams
	ap.PreviewUrl = preview + signParams

	return ap, nil
}
