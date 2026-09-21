package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"io"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/auth"
)

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
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot list attachments")
	}

	f := types.AttachmentFilter{
		NamespaceID: r.NamespaceID,
		Kind:        r.Kind,
		ModuleID:    r.ModuleID,
		RecordID:    r.RecordID,
		FieldName:   r.FieldName,
	}

	set, filter, err := ctrl.attachment.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Attachment) Read(ctx context.Context, r *request.AttachmentRead) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot read attachment")
	}

	a, err := ctrl.attachment.FindByID(ctx, r.NamespaceID, r.AttachmentID)
	if err == nil && a.Kind != r.Kind {
		return nil, service.AttachmentErrNotFound()
	}

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl Attachment) Delete(ctx context.Context, r *request.AttachmentDelete) (interface{}, error) {
	if !auth.GetIdentityFromContext(ctx).Valid() {
		return nil, errors.Unauthorized("cannot delete attachment")
	}

	a, err := ctrl.attachment.FindByID(ctx, r.NamespaceID, r.AttachmentID)
	if err != nil {
		return nil, err
	} else if a.Kind != r.Kind {
		return nil, service.AttachmentErrNotFound()
	}

	return api.OK(), ctrl.attachment.DeleteByID(ctx, r.NamespaceID, r.AttachmentID)
}

func (ctrl Attachment) Original(ctx context.Context, r *request.AttachmentOriginal) (interface{}, error) {
	att, err := ctrl.attachment.FindForServing(ctx, r.NamespaceID, r.Kind, r.AttachmentID)
	if err != nil {
		return ctrl.notFound(), nil
	}

	if err = ctrl.isAccessible(att, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(att, false, r.Download)
}

func (ctrl Attachment) Preview(ctx context.Context, r *request.AttachmentPreview) (interface{}, error) {
	att, err := ctrl.attachment.FindForServing(ctx, r.NamespaceID, r.Kind, r.AttachmentID)
	if err != nil {
		return ctrl.notFound(), nil
	}

	if err = ctrl.isAccessible(att, r.UserID, r.Sign); err != nil {
		return nil, err
	}

	return ctrl.serve(att, true, false)
}

// isAccessible verifies the signature for non-public kinds
//
// Kind of the stored attachment decides, not the one from the request
func (ctrl Attachment) isAccessible(att *types.Attachment, userID uint64, signature string) error {
	if att.Kind == types.PageAttachment || att.Kind == types.IconAttachment || att.Kind == types.NamespaceAttachment {
		// Public Attachments
		return nil
	}

	if signature == "" {
		return errors.Unauthorized("missing signature")
	}

	if userID == 0 {
		return errors.InvalidData("missing or invalid user ID")
	}

	if !auth.DefaultSigner.Verify(signature, userID, att.NamespaceID, att.ID) {
		return errors.InvalidData("missing or invalid signature")
	}

	return nil
}

func (ctrl Attachment) notFound() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ctrl Attachment) serve(att *types.Attachment, preview, download bool) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		var (
			fh  io.ReadSeekCloser
			err error
		)

		if preview {
			fh, err = ctrl.attachment.OpenPreview(att)
		} else {
			fh, err = ctrl.attachment.OpenOriginal(att)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer fh.Close()

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
		signParams = fmt.Sprintf("?sign=%s&userID=%d", auth.DefaultSigner.Sign(userID, a.NamespaceID, a.ID), userID)

		preview string
		baseURL = fmt.Sprintf("/namespace/%d/attachment/%s/%d/", a.NamespaceID, a.Kind, a.ID)
	)

	if a.Meta.Preview != nil {
		var ext = a.Meta.Preview.Extension
		if ext == "" {
			ext = "jpg"
		}

		preview = baseURL + fmt.Sprintf("preview.%s", ext) + signParams
	}

	ap := &attachmentPayload{a}

	ap.Url = baseURL + fmt.Sprintf("original/%s", url.PathEscape(a.Name)) + signParams
	ap.PreviewUrl = preview

	return ap, nil
}
