package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/attachments.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	attachmentsStore interface {
		SearchAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		LookupAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error)
		CreateAttachment(ctx context.Context, rr ...*types.Attachment) error
		UpdateAttachment(ctx context.Context, rr ...*types.Attachment) error
		PartialUpdateAttachment(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) error
		RemoveAttachment(ctx context.Context, rr ...*types.Attachment) error
		RemoveAttachmentByID(ctx context.Context, ID uint64) error

		TruncateAttachments(ctx context.Context) error
	}
)
