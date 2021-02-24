package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	messagingChannel struct {
		cfg *EncoderConfig

		res *resource.MessagingChannel
		ch  *types.Channel

		relUsr map[string]uint64
	}
)

// mergeMessagingChannels merges b into a, prioritising a
func mergeMessagingChannels(a, b *types.Channel) *types.Channel {
	c := *a

	if c.Name == "" {
		c.Name = b.Name
	}
	if c.Topic == "" {
		c.Topic = b.Topic
	}
	if c.Type == "" {
		c.Type = b.Type
	}
	if c.Meta == nil {
		c.Meta = b.Meta
	}
	c.MembershipPolicy = b.MembershipPolicy
	c.CreatorID = b.CreatorID

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}
	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}
	if c.ArchivedAt == nil {
		c.ArchivedAt = b.ArchivedAt
	}
	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return &c
}

// findMessagingChannelRS looks for the ch in the resources & the store
//
// Provided resources are prioritized.
func findMessagingChannelRS(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (ap *types.Channel, err error) {
	ap = resource.FindMessagingChannel(rr, ii)
	if ap != nil {
		return ap, nil
	}

	return findMessagingChannelS(ctx, s, makeGenericFilter(ii))
}

// findMessagingChannelS looks for the ch in the store
func findMessagingChannelS(ctx context.Context, s store.Storer, gf genericFilter) (ap *types.Channel, err error) {
	if gf.id > 0 {
		ap, err = store.LookupMessagingChannelByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if ap != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		var aa types.ChannelSet
		aa, _, err = store.SearchMessagingChannels(ctx, s, types.ChannelFilter{Query: i})
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}
		if len(aa) > 1 {
			return nil, resourceErrIdentifierNotUnique(i)
		}
		if len(aa) == 1 {
			ap = aa[0]
			return
		}
	}

	return nil, nil
}
