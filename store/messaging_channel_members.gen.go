package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_channel_members.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingChannelMembers interface {
		SearchMessagingChannelMembers(ctx context.Context, f types.ChannelMemberFilter) (types.ChannelMemberSet, types.ChannelMemberFilter, error)

		CreateMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) error

		UpdateMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) error

		UpsertMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) error

		DeleteMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) error
		DeleteMessagingChannelMemberByChannelIDUserID(ctx context.Context, channelID uint64, userID uint64) error

		TruncateMessagingChannelMembers(ctx context.Context) error
	}
)

var _ *types.ChannelMember
var _ context.Context

// SearchMessagingChannelMembers returns all matching MessagingChannelMembers from store
func SearchMessagingChannelMembers(ctx context.Context, s MessagingChannelMembers, f types.ChannelMemberFilter) (types.ChannelMemberSet, types.ChannelMemberFilter, error) {
	return s.SearchMessagingChannelMembers(ctx, f)
}

// CreateMessagingChannelMember creates one or more MessagingChannelMembers in store
func CreateMessagingChannelMember(ctx context.Context, s MessagingChannelMembers, rr ...*types.ChannelMember) error {
	return s.CreateMessagingChannelMember(ctx, rr...)
}

// UpdateMessagingChannelMember updates one or more (existing) MessagingChannelMembers in store
func UpdateMessagingChannelMember(ctx context.Context, s MessagingChannelMembers, rr ...*types.ChannelMember) error {
	return s.UpdateMessagingChannelMember(ctx, rr...)
}

// UpsertMessagingChannelMember creates new or updates existing one or more MessagingChannelMembers in store
func UpsertMessagingChannelMember(ctx context.Context, s MessagingChannelMembers, rr ...*types.ChannelMember) error {
	return s.UpsertMessagingChannelMember(ctx, rr...)
}

// DeleteMessagingChannelMember Deletes one or more MessagingChannelMembers from store
func DeleteMessagingChannelMember(ctx context.Context, s MessagingChannelMembers, rr ...*types.ChannelMember) error {
	return s.DeleteMessagingChannelMember(ctx, rr...)
}

// DeleteMessagingChannelMemberByChannelIDUserID Deletes MessagingChannelMember from store
func DeleteMessagingChannelMemberByChannelIDUserID(ctx context.Context, s MessagingChannelMembers, channelID uint64, userID uint64) error {
	return s.DeleteMessagingChannelMemberByChannelIDUserID(ctx, channelID, userID)
}

// TruncateMessagingChannelMembers Deletes all MessagingChannelMembers from store
func TruncateMessagingChannelMembers(ctx context.Context, s MessagingChannelMembers) error {
	return s.TruncateMessagingChannelMembers(ctx)
}
