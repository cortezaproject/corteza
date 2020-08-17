package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/compose_attachments.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	composeAttachmentsStore interface {
		SearchComposeAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		LookupComposeAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error)
		CreateComposeAttachment(ctx context.Context, rr ...*types.Attachment) error
		UpdateComposeAttachment(ctx context.Context, rr ...*types.Attachment) error
		PartialUpdateComposeAttachment(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) error
		RemoveComposeAttachment(ctx context.Context, rr ...*types.Attachment) error
		RemoveComposeAttachmentByID(ctx context.Context, ID uint64) error

		TruncateComposeAttachments(ctx context.Context) error
	}
)
