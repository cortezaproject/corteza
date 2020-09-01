package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_unread.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingUnreads interface {
		CreateMessagingUnread(ctx context.Context, rr ...*types.Unread) error

		UpdateMessagingUnread(ctx context.Context, rr ...*types.Unread) error
		PartialMessagingUnreadUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Unread) error

		UpsertMessagingUnread(ctx context.Context, rr ...*types.Unread) error

		DeleteMessagingUnread(ctx context.Context, rr ...*types.Unread) error
		DeleteMessagingUnreadByChannelIDReplyToUserID(ctx context.Context, channelID uint64, replyTo uint64, userID uint64) error

		TruncateMessagingUnreads(ctx context.Context) error
	}
)

var _ *types.Unread
var _ context.Context

// CreateMessagingUnread creates one or more MessagingUnreads in store
func CreateMessagingUnread(ctx context.Context, s MessagingUnreads, rr ...*types.Unread) error {
	return s.CreateMessagingUnread(ctx, rr...)
}

// UpdateMessagingUnread updates one or more (existing) MessagingUnreads in store
func UpdateMessagingUnread(ctx context.Context, s MessagingUnreads, rr ...*types.Unread) error {
	return s.UpdateMessagingUnread(ctx, rr...)
}

// PartialMessagingUnreadUpdate updates one or more existing MessagingUnreads in store
func PartialMessagingUnreadUpdate(ctx context.Context, s MessagingUnreads, onlyColumns []string, rr ...*types.Unread) error {
	return s.PartialMessagingUnreadUpdate(ctx, onlyColumns, rr...)
}

// UpsertMessagingUnread creates new or updates existing one or more MessagingUnreads in store
func UpsertMessagingUnread(ctx context.Context, s MessagingUnreads, rr ...*types.Unread) error {
	return s.UpsertMessagingUnread(ctx, rr...)
}

// DeleteMessagingUnread Deletes one or more MessagingUnreads from store
func DeleteMessagingUnread(ctx context.Context, s MessagingUnreads, rr ...*types.Unread) error {
	return s.DeleteMessagingUnread(ctx, rr...)
}

// DeleteMessagingUnreadByChannelIDReplyToUserID Deletes MessagingUnread from store
func DeleteMessagingUnreadByChannelIDReplyToUserID(ctx context.Context, s MessagingUnreads, channelID uint64, replyTo uint64, userID uint64) error {
	return s.DeleteMessagingUnreadByChannelIDReplyToUserID(ctx, channelID, replyTo, userID)
}

// TruncateMessagingUnreads Deletes all MessagingUnreads from store
func TruncateMessagingUnreads(ctx context.Context, s MessagingUnreads) error {
	return s.TruncateMessagingUnreads(ctx)
}
