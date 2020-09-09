package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_channels.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingChannels interface {
		SearchMessagingChannels(ctx context.Context, f types.ChannelFilter) (types.ChannelSet, types.ChannelFilter, error)
		LookupMessagingChannelByID(ctx context.Context, id uint64) (*types.Channel, error)

		CreateMessagingChannel(ctx context.Context, rr ...*types.Channel) error

		UpdateMessagingChannel(ctx context.Context, rr ...*types.Channel) error

		UpsertMessagingChannel(ctx context.Context, rr ...*types.Channel) error

		DeleteMessagingChannel(ctx context.Context, rr ...*types.Channel) error
		DeleteMessagingChannelByID(ctx context.Context, ID uint64) error

		TruncateMessagingChannels(ctx context.Context) error

		// Additional custom functions

		// LookupMessagingChannelByMemberSet (custom function)
		LookupMessagingChannelByMemberSet(ctx context.Context, _memberIDs ...uint64) (*types.Channel, error)
	}
)

var _ *types.Channel
var _ context.Context

// SearchMessagingChannels returns all matching MessagingChannels from store
func SearchMessagingChannels(ctx context.Context, s MessagingChannels, f types.ChannelFilter) (types.ChannelSet, types.ChannelFilter, error) {
	return s.SearchMessagingChannels(ctx, f)
}

// LookupMessagingChannelByID searches for attachment by its ID
//
// It returns attachment even if deleted
func LookupMessagingChannelByID(ctx context.Context, s MessagingChannels, id uint64) (*types.Channel, error) {
	return s.LookupMessagingChannelByID(ctx, id)
}

// CreateMessagingChannel creates one or more MessagingChannels in store
func CreateMessagingChannel(ctx context.Context, s MessagingChannels, rr ...*types.Channel) error {
	return s.CreateMessagingChannel(ctx, rr...)
}

// UpdateMessagingChannel updates one or more (existing) MessagingChannels in store
func UpdateMessagingChannel(ctx context.Context, s MessagingChannels, rr ...*types.Channel) error {
	return s.UpdateMessagingChannel(ctx, rr...)
}

// UpsertMessagingChannel creates new or updates existing one or more MessagingChannels in store
func UpsertMessagingChannel(ctx context.Context, s MessagingChannels, rr ...*types.Channel) error {
	return s.UpsertMessagingChannel(ctx, rr...)
}

// DeleteMessagingChannel Deletes one or more MessagingChannels from store
func DeleteMessagingChannel(ctx context.Context, s MessagingChannels, rr ...*types.Channel) error {
	return s.DeleteMessagingChannel(ctx, rr...)
}

// DeleteMessagingChannelByID Deletes MessagingChannel from store
func DeleteMessagingChannelByID(ctx context.Context, s MessagingChannels, ID uint64) error {
	return s.DeleteMessagingChannelByID(ctx, ID)
}

// TruncateMessagingChannels Deletes all MessagingChannels from store
func TruncateMessagingChannels(ctx context.Context, s MessagingChannels) error {
	return s.TruncateMessagingChannels(ctx)
}

func LookupMessagingChannelByMemberSet(ctx context.Context, s MessagingChannels, _memberIDs ...uint64) (*types.Channel, error) {
	return s.LookupMessagingChannelByMemberSet(ctx, _memberIDs...)
}
