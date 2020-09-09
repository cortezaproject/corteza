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

		UpsertMessagingUnread(ctx context.Context, rr ...*types.Unread) error

		DeleteMessagingUnread(ctx context.Context, rr ...*types.Unread) error
		DeleteMessagingUnreadByChannelIDReplyToUserID(ctx context.Context, channelID uint64, replyTo uint64, userID uint64) error

		TruncateMessagingUnreads(ctx context.Context) error

		// Additional custom functions

		// CountMessagingUnreadThreads (custom function)
		CountMessagingUnreadThreads(ctx context.Context, _userID uint64, _channelID uint64) (types.UnreadSet, error)

		// CountMessagingUnread (custom function)
		CountMessagingUnread(ctx context.Context, _userID uint64, _channelID uint64, _threadIDs ...uint64) (types.UnreadSet, error)

		// ResetMessagingUnreadThreads (custom function)
		ResetMessagingUnreadThreads(ctx context.Context, _userID uint64, _channelID uint64) error

		// PresetMessagingUnread (custom function)
		PresetMessagingUnread(ctx context.Context, _channelID uint64, _threadIDs uint64, _userID ...uint64) error

		// IncMessagingUnreadCount (custom function)
		IncMessagingUnreadCount(ctx context.Context, _channelID uint64, _threadIDs uint64, _userID uint64) error

		// DecMessagingUnreadCount (custom function)
		DecMessagingUnreadCount(ctx context.Context, _channelID uint64, _threadIDs uint64, _userID uint64) error
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

func CountMessagingUnreadThreads(ctx context.Context, s MessagingUnreads, _userID uint64, _channelID uint64) (types.UnreadSet, error) {
	return s.CountMessagingUnreadThreads(ctx, _userID, _channelID)
}

func CountMessagingUnread(ctx context.Context, s MessagingUnreads, _userID uint64, _channelID uint64, _threadIDs ...uint64) (types.UnreadSet, error) {
	return s.CountMessagingUnread(ctx, _userID, _channelID, _threadIDs...)
}

func ResetMessagingUnreadThreads(ctx context.Context, s MessagingUnreads, _userID uint64, _channelID uint64) error {
	return s.ResetMessagingUnreadThreads(ctx, _userID, _channelID)
}

func PresetMessagingUnread(ctx context.Context, s MessagingUnreads, _channelID uint64, _threadIDs uint64, _userID ...uint64) error {
	return s.PresetMessagingUnread(ctx, _channelID, _threadIDs, _userID...)
}

func IncMessagingUnreadCount(ctx context.Context, s MessagingUnreads, _channelID uint64, _threadIDs uint64, _userID uint64) error {
	return s.IncMessagingUnreadCount(ctx, _channelID, _threadIDs, _userID)
}

func DecMessagingUnreadCount(ctx context.Context, s MessagingUnreads, _channelID uint64, _threadIDs uint64, _userID uint64) error {
	return s.DecMessagingUnreadCount(ctx, _channelID, _threadIDs, _userID)
}
