package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_mentions.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingMentions interface {
		SearchMessagingMentions(ctx context.Context, f types.MentionFilter) (types.MentionSet, types.MentionFilter, error)
		LookupMessagingMentionByID(ctx context.Context, id uint64) (*types.Mention, error)

		CreateMessagingMention(ctx context.Context, rr ...*types.Mention) error

		UpdateMessagingMention(ctx context.Context, rr ...*types.Mention) error
		PartialMessagingMentionUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Mention) error

		UpsertMessagingMention(ctx context.Context, rr ...*types.Mention) error

		DeleteMessagingMention(ctx context.Context, rr ...*types.Mention) error
		DeleteMessagingMentionByID(ctx context.Context, ID uint64) error

		TruncateMessagingMentions(ctx context.Context) error
	}
)

var _ *types.Mention
var _ context.Context

// SearchMessagingMentions returns all matching MessagingMentions from store
func SearchMessagingMentions(ctx context.Context, s MessagingMentions, f types.MentionFilter) (types.MentionSet, types.MentionFilter, error) {
	return s.SearchMessagingMentions(ctx, f)
}

// LookupMessagingMentionByID searches for attachment by its ID
//
// It returns attachment even if deleted
func LookupMessagingMentionByID(ctx context.Context, s MessagingMentions, id uint64) (*types.Mention, error) {
	return s.LookupMessagingMentionByID(ctx, id)
}

// CreateMessagingMention creates one or more MessagingMentions in store
func CreateMessagingMention(ctx context.Context, s MessagingMentions, rr ...*types.Mention) error {
	return s.CreateMessagingMention(ctx, rr...)
}

// UpdateMessagingMention updates one or more (existing) MessagingMentions in store
func UpdateMessagingMention(ctx context.Context, s MessagingMentions, rr ...*types.Mention) error {
	return s.UpdateMessagingMention(ctx, rr...)
}

// PartialMessagingMentionUpdate updates one or more existing MessagingMentions in store
func PartialMessagingMentionUpdate(ctx context.Context, s MessagingMentions, onlyColumns []string, rr ...*types.Mention) error {
	return s.PartialMessagingMentionUpdate(ctx, onlyColumns, rr...)
}

// UpsertMessagingMention creates new or updates existing one or more MessagingMentions in store
func UpsertMessagingMention(ctx context.Context, s MessagingMentions, rr ...*types.Mention) error {
	return s.UpsertMessagingMention(ctx, rr...)
}

// DeleteMessagingMention Deletes one or more MessagingMentions from store
func DeleteMessagingMention(ctx context.Context, s MessagingMentions, rr ...*types.Mention) error {
	return s.DeleteMessagingMention(ctx, rr...)
}

// DeleteMessagingMentionByID Deletes MessagingMention from store
func DeleteMessagingMentionByID(ctx context.Context, s MessagingMentions, ID uint64) error {
	return s.DeleteMessagingMentionByID(ctx, ID)
}

// TruncateMessagingMentions Deletes all MessagingMentions from store
func TruncateMessagingMentions(ctx context.Context, s MessagingMentions) error {
	return s.TruncateMessagingMentions(ctx)
}
