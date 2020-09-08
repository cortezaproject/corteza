package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_message_attachments.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingMessageAttachments interface {
		LookupMessagingMessageAttachmentByMessageID(ctx context.Context, message_id uint64) (*types.MessageAttachment, error)

		CreateMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) error

		UpdateMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) error
		PartialMessagingMessageAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.MessageAttachment) error

		UpsertMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) error

		DeleteMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) error
		DeleteMessagingMessageAttachmentByMessageID(ctx context.Context, messageID uint64) error

		TruncateMessagingMessageAttachments(ctx context.Context) error
	}
)

var _ *types.MessageAttachment
var _ context.Context

// LookupMessagingMessageAttachmentByMessageID searches for message attachment by message ID
func LookupMessagingMessageAttachmentByMessageID(ctx context.Context, s MessagingMessageAttachments, message_id uint64) (*types.MessageAttachment, error) {
	return s.LookupMessagingMessageAttachmentByMessageID(ctx, message_id)
}

// CreateMessagingMessageAttachment creates one or more MessagingMessageAttachments in store
func CreateMessagingMessageAttachment(ctx context.Context, s MessagingMessageAttachments, rr ...*types.MessageAttachment) error {
	return s.CreateMessagingMessageAttachment(ctx, rr...)
}

// UpdateMessagingMessageAttachment updates one or more (existing) MessagingMessageAttachments in store
func UpdateMessagingMessageAttachment(ctx context.Context, s MessagingMessageAttachments, rr ...*types.MessageAttachment) error {
	return s.UpdateMessagingMessageAttachment(ctx, rr...)
}

// PartialMessagingMessageAttachmentUpdate updates one or more existing MessagingMessageAttachments in store
func PartialMessagingMessageAttachmentUpdate(ctx context.Context, s MessagingMessageAttachments, onlyColumns []string, rr ...*types.MessageAttachment) error {
	return s.PartialMessagingMessageAttachmentUpdate(ctx, onlyColumns, rr...)
}

// UpsertMessagingMessageAttachment creates new or updates existing one or more MessagingMessageAttachments in store
func UpsertMessagingMessageAttachment(ctx context.Context, s MessagingMessageAttachments, rr ...*types.MessageAttachment) error {
	return s.UpsertMessagingMessageAttachment(ctx, rr...)
}

// DeleteMessagingMessageAttachment Deletes one or more MessagingMessageAttachments from store
func DeleteMessagingMessageAttachment(ctx context.Context, s MessagingMessageAttachments, rr ...*types.MessageAttachment) error {
	return s.DeleteMessagingMessageAttachment(ctx, rr...)
}

// DeleteMessagingMessageAttachmentByMessageID Deletes MessagingMessageAttachment from store
func DeleteMessagingMessageAttachmentByMessageID(ctx context.Context, s MessagingMessageAttachments, messageID uint64) error {
	return s.DeleteMessagingMessageAttachmentByMessageID(ctx, messageID)
}

// TruncateMessagingMessageAttachments Deletes all MessagingMessageAttachments from store
func TruncateMessagingMessageAttachments(ctx context.Context, s MessagingMessageAttachments) error {
	return s.TruncateMessagingMessageAttachments(ctx)
}
