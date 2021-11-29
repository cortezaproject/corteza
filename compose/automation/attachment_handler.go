package automation

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	attachmentService interface {
		FindByID(ctx context.Context, namespaceID, attachmentID uint64) (*types.Attachment, error)
		CreateRecordAttachment(ctx context.Context, namespaceID uint64, name string, size int64, fh io.ReadSeeker, moduleID, recordID uint64) (att *types.Attachment, err error)
		DeleteByID(ctx context.Context, namespaceID uint64, attachmentID uint64) error
		OpenOriginal(att *types.Attachment) (io.ReadSeeker, error)
		OpenPreview(att *types.Attachment) (io.ReadSeeker, error)
	}

	attachmentHandler struct {
		reg modulesHandlerRegistry
		svc attachmentService
	}

	attachmentLookup interface {
		GetAttachment() (bool, uint64, *types.Attachment)
	}
)

func AttachmentHandler(reg modulesHandlerRegistry, svc attachmentService) *attachmentHandler {
	h := &attachmentHandler{
		reg: reg,
		svc: svc,
	}

	h.register()
	return h
}

func (h attachmentHandler) lookup(ctx context.Context, args *attachmentLookupArgs) (results *attachmentLookupResults, err error) {
	results = &attachmentLookupResults{}
	results.Attachment, err = h.svc.FindByID(ctx, 0, args.Attachment)
	return
}

func (h attachmentHandler) delete(ctx context.Context, args *attachmentDeleteArgs) error {
	return h.svc.DeleteByID(ctx, 0, args.Attachment)
}

func (h attachmentHandler) openOriginal(ctx context.Context, args *attachmentOpenOriginalArgs) (*attachmentOpenOriginalResults, error) {
	att, err := lookupAttachment(ctx, h.svc, args)
	if err != nil {
		return nil, err
	}

	r := &attachmentOpenOriginalResults{}
	r.Content, err = h.svc.OpenOriginal(att)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (h attachmentHandler) openPreview(ctx context.Context, args *attachmentOpenPreviewArgs) (*attachmentOpenPreviewResults, error) {
	att, err := lookupAttachment(ctx, h.svc, args)
	if err != nil {
		return nil, err
	}

	r := &attachmentOpenPreviewResults{}
	r.Content, err = h.svc.OpenOriginal(att)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (h attachmentHandler) create(ctx context.Context, args *attachmentCreateArgs) (*attachmentCreateResults, error) {
	var (
		err error
		att *types.Attachment

		fh   io.ReadSeeker
		size int64
	)

	switch {
	case len(args.contentBytes) > 0:
		size = int64(len(args.contentString))
		fh = bytes.NewReader(args.contentBytes)

	case args.contentStream != nil:
		if rs, is := args.contentStream.(io.ReadSeeker); is {
			fh = rs
		}

	default:
		fh = strings.NewReader(args.contentString)
		size = int64(len(args.contentString))
	}

	switch {
	case args.Resource != nil:
		att, err = h.svc.CreateRecordAttachment(
			ctx,
			args.Resource.NamespaceID,
			args.Name,
			size,
			fh,
			args.Resource.ModuleID,
			args.Resource.ID,
		)
	default:
		return nil, fmt.Errorf("unknown resource")
	}

	if err != nil {
		return nil, err
	}

	return &attachmentCreateResults{att}, nil
}

func lookupAttachment(ctx context.Context, svc attachmentService, args attachmentLookup) (*types.Attachment, error) {
	_, ID, attachment := args.GetAttachment()

	switch {
	case attachment != nil:
		return attachment, nil
	case ID > 0:
		return svc.FindByID(ctx, 0, ID)
	}

	return nil, fmt.Errorf("empty attachment lookup params")
}
