package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_attachments.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingAttachments interface {
		SearchMessagingAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		LookupMessagingAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error)

		CreateMessagingAttachment(ctx context.Context, rr ...*types.Attachment) error

		UpdateMessagingAttachment(ctx context.Context, rr ...*types.Attachment) error
		PartialMessagingAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) error

		UpsertMessagingAttachment(ctx context.Context, rr ...*types.Attachment) error

		DeleteMessagingAttachment(ctx context.Context, rr ...*types.Attachment) error
		DeleteMessagingAttachmentByID(ctx context.Context, ID uint64) error

		TruncateMessagingAttachments(ctx context.Context) error
	}
)

var _ *types.Attachment
var _ context.Context

// SearchMessagingAttachments returns all matching MessagingAttachments from store
func SearchMessagingAttachments(ctx context.Context, s MessagingAttachments, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	return s.SearchMessagingAttachments(ctx, f)
}

// LookupMessagingAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func LookupMessagingAttachmentByID(ctx context.Context, s MessagingAttachments, id uint64) (*types.Attachment, error) {
	return s.LookupMessagingAttachmentByID(ctx, id)
}

// CreateMessagingAttachment creates one or more MessagingAttachments in store
func CreateMessagingAttachment(ctx context.Context, s MessagingAttachments, rr ...*types.Attachment) error {
	return s.CreateMessagingAttachment(ctx, rr...)
}

// UpdateMessagingAttachment updates one or more (existing) MessagingAttachments in store
func UpdateMessagingAttachment(ctx context.Context, s MessagingAttachments, rr ...*types.Attachment) error {
	return s.UpdateMessagingAttachment(ctx, rr...)
}

// PartialMessagingAttachmentUpdate updates one or more existing MessagingAttachments in store
func PartialMessagingAttachmentUpdate(ctx context.Context, s MessagingAttachments, onlyColumns []string, rr ...*types.Attachment) error {
	return s.PartialMessagingAttachmentUpdate(ctx, onlyColumns, rr...)
}

// UpsertMessagingAttachment creates new or updates existing one or more MessagingAttachments in store
func UpsertMessagingAttachment(ctx context.Context, s MessagingAttachments, rr ...*types.Attachment) error {
	return s.UpsertMessagingAttachment(ctx, rr...)
}

// DeleteMessagingAttachment Deletes one or more MessagingAttachments from store
func DeleteMessagingAttachment(ctx context.Context, s MessagingAttachments, rr ...*types.Attachment) error {
	return s.DeleteMessagingAttachment(ctx, rr...)
}

// DeleteMessagingAttachmentByID Deletes MessagingAttachment from store
func DeleteMessagingAttachmentByID(ctx context.Context, s MessagingAttachments, ID uint64) error {
	return s.DeleteMessagingAttachmentByID(ctx, ID)
}

// TruncateMessagingAttachments Deletes all MessagingAttachments from store
func TruncateMessagingAttachments(ctx context.Context, s MessagingAttachments) error {
	return s.TruncateMessagingAttachments(ctx)
}
