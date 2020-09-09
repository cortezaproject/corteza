package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/attachments.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Attachments interface {
		SearchAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		LookupAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error)

		CreateAttachment(ctx context.Context, rr ...*types.Attachment) error

		UpdateAttachment(ctx context.Context, rr ...*types.Attachment) error

		UpsertAttachment(ctx context.Context, rr ...*types.Attachment) error

		DeleteAttachment(ctx context.Context, rr ...*types.Attachment) error
		DeleteAttachmentByID(ctx context.Context, ID uint64) error

		TruncateAttachments(ctx context.Context) error
	}
)

var _ *types.Attachment
var _ context.Context

// SearchAttachments returns all matching Attachments from store
func SearchAttachments(ctx context.Context, s Attachments, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	return s.SearchAttachments(ctx, f)
}

// LookupAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func LookupAttachmentByID(ctx context.Context, s Attachments, id uint64) (*types.Attachment, error) {
	return s.LookupAttachmentByID(ctx, id)
}

// CreateAttachment creates one or more Attachments in store
func CreateAttachment(ctx context.Context, s Attachments, rr ...*types.Attachment) error {
	return s.CreateAttachment(ctx, rr...)
}

// UpdateAttachment updates one or more (existing) Attachments in store
func UpdateAttachment(ctx context.Context, s Attachments, rr ...*types.Attachment) error {
	return s.UpdateAttachment(ctx, rr...)
}

// UpsertAttachment creates new or updates existing one or more Attachments in store
func UpsertAttachment(ctx context.Context, s Attachments, rr ...*types.Attachment) error {
	return s.UpsertAttachment(ctx, rr...)
}

// DeleteAttachment Deletes one or more Attachments from store
func DeleteAttachment(ctx context.Context, s Attachments, rr ...*types.Attachment) error {
	return s.DeleteAttachment(ctx, rr...)
}

// DeleteAttachmentByID Deletes Attachment from store
func DeleteAttachmentByID(ctx context.Context, s Attachments, ID uint64) error {
	return s.DeleteAttachmentByID(ctx, ID)
}

// TruncateAttachments Deletes all Attachments from store
func TruncateAttachments(ctx context.Context, s Attachments) error {
	return s.TruncateAttachments(ctx)
}
