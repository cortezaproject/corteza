package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/messaging_flags.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingFlags interface {
		SearchMessagingFlags(ctx context.Context, f types.MessageFlagFilter) (types.MessageFlagSet, types.MessageFlagFilter, error)
		LookupMessagingFlagByID(ctx context.Context, id uint64) (*types.MessageFlag, error)

		CreateMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) error

		UpdateMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) error

		UpsertMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) error

		DeleteMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) error
		DeleteMessagingFlagByID(ctx context.Context, ID uint64) error

		TruncateMessagingFlags(ctx context.Context) error
	}
)

var _ *types.MessageFlag
var _ context.Context

// SearchMessagingFlags returns all matching MessagingFlags from store
func SearchMessagingFlags(ctx context.Context, s MessagingFlags, f types.MessageFlagFilter) (types.MessageFlagSet, types.MessageFlagFilter, error) {
	return s.SearchMessagingFlags(ctx, f)
}

// LookupMessagingFlagByID searches for flags by ID
func LookupMessagingFlagByID(ctx context.Context, s MessagingFlags, id uint64) (*types.MessageFlag, error) {
	return s.LookupMessagingFlagByID(ctx, id)
}

// CreateMessagingFlag creates one or more MessagingFlags in store
func CreateMessagingFlag(ctx context.Context, s MessagingFlags, rr ...*types.MessageFlag) error {
	return s.CreateMessagingFlag(ctx, rr...)
}

// UpdateMessagingFlag updates one or more (existing) MessagingFlags in store
func UpdateMessagingFlag(ctx context.Context, s MessagingFlags, rr ...*types.MessageFlag) error {
	return s.UpdateMessagingFlag(ctx, rr...)
}

// UpsertMessagingFlag creates new or updates existing one or more MessagingFlags in store
func UpsertMessagingFlag(ctx context.Context, s MessagingFlags, rr ...*types.MessageFlag) error {
	return s.UpsertMessagingFlag(ctx, rr...)
}

// DeleteMessagingFlag Deletes one or more MessagingFlags from store
func DeleteMessagingFlag(ctx context.Context, s MessagingFlags, rr ...*types.MessageFlag) error {
	return s.DeleteMessagingFlag(ctx, rr...)
}

// DeleteMessagingFlagByID Deletes MessagingFlag from store
func DeleteMessagingFlagByID(ctx context.Context, s MessagingFlags, ID uint64) error {
	return s.DeleteMessagingFlagByID(ctx, ID)
}

// TruncateMessagingFlags Deletes all MessagingFlags from store
func TruncateMessagingFlags(ctx context.Context, s MessagingFlags) error {
	return s.TruncateMessagingFlags(ctx)
}
