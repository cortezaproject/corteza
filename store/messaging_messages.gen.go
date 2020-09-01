package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_messages.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingMessages interface {
		SearchMessagingMessages(ctx context.Context, f types.MessageFilter) (types.MessageSet, types.MessageFilter, error)
		LookupMessagingMessageByID(ctx context.Context, id uint64) (*types.Message, error)

		CreateMessagingMessage(ctx context.Context, rr ...*types.Message) error

		UpdateMessagingMessage(ctx context.Context, rr ...*types.Message) error
		PartialMessagingMessageUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Message) error

		UpsertMessagingMessage(ctx context.Context, rr ...*types.Message) error

		DeleteMessagingMessage(ctx context.Context, rr ...*types.Message) error
		DeleteMessagingMessageByID(ctx context.Context, ID uint64) error

		TruncateMessagingMessages(ctx context.Context) error
	}
)

var _ *types.Message
var _ context.Context

// SearchMessagingMessages returns all matching MessagingMessages from store
func SearchMessagingMessages(ctx context.Context, s MessagingMessages, f types.MessageFilter) (types.MessageSet, types.MessageFilter, error) {
	return s.SearchMessagingMessages(ctx, f)
}

// LookupMessagingMessageByID searches for attachment by its ID
//
// It returns attachment even if deleted
func LookupMessagingMessageByID(ctx context.Context, s MessagingMessages, id uint64) (*types.Message, error) {
	return s.LookupMessagingMessageByID(ctx, id)
}

// CreateMessagingMessage creates one or more MessagingMessages in store
func CreateMessagingMessage(ctx context.Context, s MessagingMessages, rr ...*types.Message) error {
	return s.CreateMessagingMessage(ctx, rr...)
}

// UpdateMessagingMessage updates one or more (existing) MessagingMessages in store
func UpdateMessagingMessage(ctx context.Context, s MessagingMessages, rr ...*types.Message) error {
	return s.UpdateMessagingMessage(ctx, rr...)
}

// PartialMessagingMessageUpdate updates one or more existing MessagingMessages in store
func PartialMessagingMessageUpdate(ctx context.Context, s MessagingMessages, onlyColumns []string, rr ...*types.Message) error {
	return s.PartialMessagingMessageUpdate(ctx, onlyColumns, rr...)
}

// UpsertMessagingMessage creates new or updates existing one or more MessagingMessages in store
func UpsertMessagingMessage(ctx context.Context, s MessagingMessages, rr ...*types.Message) error {
	return s.UpsertMessagingMessage(ctx, rr...)
}

// DeleteMessagingMessage Deletes one or more MessagingMessages from store
func DeleteMessagingMessage(ctx context.Context, s MessagingMessages, rr ...*types.Message) error {
	return s.DeleteMessagingMessage(ctx, rr...)
}

// DeleteMessagingMessageByID Deletes MessagingMessage from store
func DeleteMessagingMessageByID(ctx context.Context, s MessagingMessages, ID uint64) error {
	return s.DeleteMessagingMessageByID(ctx, ID)
}

// TruncateMessagingMessages Deletes all MessagingMessages from store
func TruncateMessagingMessages(ctx context.Context, s MessagingMessages) error {
	return s.TruncateMessagingMessages(ctx)
}
