package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
)

type (
	Attachment struct {
		att service.AttachmentService
	}
)

func (Attachment) New() *Attachment {
	ctrl := &Attachment{}
	ctrl.att = service.DefaultAttachment
	return ctrl
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
	if signature == "" {
		return fmt.Errorf("unauthorized")
	}

	if userID == 0 {
		return fmt.Errorf("missing or invalid user ID")
	}

	if attachmentID == 0 {
		return fmt.Errorf("missing or invalid attachment ID")
	}

	if !auth.DefaultSigner.Verify(signature, userID, attachmentID) {
		return fmt.Errorf("missing or invalid signature")
	}

	return nil
}

func (ctrl Attachment) serve(ctx context.Context, ID uint64, preview, download bool) (interface{}, error) {
	return func(w http.ResponseWriter, req *http.Request) {
		att, err := ctrl.att.FindByID(ctx, ID)
		if err != nil {
			if errors.IsNotFound(err) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		var fh io.ReadSeeker

		if preview {
			fh, err = ctrl.att.OpenPreview(ctx, att)
		} else {
			fh, err = ctrl.att.OpenOriginal(ctx, att)
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
