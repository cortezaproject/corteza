package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/compose_attachments.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeAttachments interface {
		SearchComposeAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		LookupComposeAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error)

		CreateComposeAttachment(ctx context.Context, rr ...*types.Attachment) error

		UpdateComposeAttachment(ctx context.Context, rr ...*types.Attachment) error

		UpsertComposeAttachment(ctx context.Context, rr ...*types.Attachment) error

		DeleteComposeAttachment(ctx context.Context, rr ...*types.Attachment) error
		DeleteComposeAttachmentByID(ctx context.Context, ID uint64) error

		TruncateComposeAttachments(ctx context.Context) error
	}
)

var _ *types.Attachment
var _ context.Context

// SearchComposeAttachments returns all matching ComposeAttachments from store
func SearchComposeAttachments(ctx context.Context, s ComposeAttachments, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	return s.SearchComposeAttachments(ctx, f)
}

// LookupComposeAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func LookupComposeAttachmentByID(ctx context.Context, s ComposeAttachments, id uint64) (*types.Attachment, error) {
	return s.LookupComposeAttachmentByID(ctx, id)
}

// CreateComposeAttachment creates one or more ComposeAttachments in store
func CreateComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*types.Attachment) error {
	return s.CreateComposeAttachment(ctx, rr...)
}

// UpdateComposeAttachment updates one or more (existing) ComposeAttachments in store
func UpdateComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*types.Attachment) error {
	return s.UpdateComposeAttachment(ctx, rr...)
}

// UpsertComposeAttachment creates new or updates existing one or more ComposeAttachments in store
func UpsertComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*types.Attachment) error {
	return s.UpsertComposeAttachment(ctx, rr...)
}

// DeleteComposeAttachment Deletes one or more ComposeAttachments from store
func DeleteComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*types.Attachment) error {
	return s.DeleteComposeAttachment(ctx, rr...)
}

// DeleteComposeAttachmentByID Deletes ComposeAttachment from store
func DeleteComposeAttachmentByID(ctx context.Context, s ComposeAttachments, ID uint64) error {
	return s.DeleteComposeAttachmentByID(ctx, ID)
}

// TruncateComposeAttachments Deletes all ComposeAttachments from store
func TruncateComposeAttachments(ctx context.Context, s ComposeAttachments) error {
	return s.TruncateComposeAttachments(ctx)
}
