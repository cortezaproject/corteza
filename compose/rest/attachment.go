package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"

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

// Attachments returns list of all files attached to records
func (ctrl Attachment) List(ctx context.Context, r *request.AttachmentList) (interface{}, error) {
	f := types.AttachmentFilter{
		NamespaceID: r.NamespaceID,
		Kind:        r.Kind,
		ModuleID:    r.ModuleID,
		RecordID:    r.RecordID,
		FieldName:   r.FieldName,
		// Filter:    r.Filter,
		PerPage: r.PerPage,
		Page:    r.Page,
		// Sort:      r.Sort,
	}

	set, filter, err := ctrl.attachment.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Attachment) Read(ctx context.Context, r *request.AttachmentRead) (interface{}, error) {
	a, err := ctrl.attachment.With(ctx).FindByID(r.NamespaceID, r.AttachmentID)
	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl Attachment) Delete(ctx context.Context, r *request.AttachmentDelete) (interface{}, error) {
	_, err := ctrl.attachment.With(ctx).FindByID(r.NamespaceID, r.AttachmentID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.attachment.With(ctx).DeleteByID(r.NamespaceID, r.AttachmentID)
}

func (ctrl Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	if err := ctrl.isAccessible(r.NamespaceID, r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(ctx, r.NamespaceID, r.AttachmentID, false, r.Download)
}

func (ctrl Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	if err := ctrl.isAccessible(r.NamespaceID, r.AttachmentID, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(ctx, r.NamespaceID, r.AttachmentID, true, false)
}

func (ctrl Attachment) isAccessible(namespaceID, attachmentID, userID uint64, signature string) error {
	if userID == 0 {
		return errors.New("missing or invalid user ID")
	}

	if attachmentID == 0 {
		return errors.New("missing or invalid attachment ID")
	}

	if auth.DefaultSigner.Verify(signature, userID, namespaceID, attachmentID) {
		return errors.New("missing or invalid signature")
	}

	return nil
}

func (ctrl Attachment) serve(ctx context.Context, namespaceID, attachmentID uint64, preview, download bool) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		att, err := ctrl.attachment.With(ctx).FindByID(namespaceID, attachmentID)
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
		signParams = fmt.Sprintf("?sign=%s&userID=%d", auth.DefaultSigner.Sign(userID, a.NamespaceID, a.ID), userID)

		preview string
		baseURL = fmt.Sprintf("/namespace/%d/attachment/%s/%d/", a.NamespaceID, a.Kind, a.ID)
	)

	if a.Meta.Preview != nil {
		var ext = a.Meta.Preview.Extension
		if ext == "" {
			ext = "jpg"
		}

		preview = baseURL + fmt.Sprintf("preview.%s", ext)
	}

	ap := &attachmentPayload{a}

	ap.Url = baseURL + fmt.Sprintf("original/%s", url.PathEscape(a.Name)) + signParams
	ap.PreviewUrl = preview + signParams

	return ap, nil
}
