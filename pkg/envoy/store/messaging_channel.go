package store

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/messaging/types"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	messagingChannelState struct {
		cfg *EncoderConfig

		res *resource.MessagingChannel
		ch  *types.Channel

		relUsr map[string]uint64
	}
)

func NewMessagingChannelState(res *resource.MessagingChannel, cfg *EncoderConfig) resourceState {
	return &messagingChannelState{
		cfg: mergeConfig(cfg, res.Config()),
		res: res,
	}
}

func (n *messagingChannelState) Prepare(ctx context.Context, s store.Storer, rs *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
	}

	// Sys users
	n.relUsr = make(map[string]uint64)
	if err = resolveUserRefs(ctx, s, rs.ParentResources, n.res.UserRefs(), n.relUsr); err != nil {
		return err
	}

	// Get the existing channel
	n.ch, err = findMessagingChannelS(ctx, s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.ch != nil {
		n.res.Res.ID = n.ch.ID
	}
	return nil
}

// Encode encodes the given messagingChannel
func (n *messagingChannelState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	res := n.res.Res
	exists := n.ch != nil && n.ch.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.ch.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// Sys users
	for idf, ID := range n.relUsr {
		if ID != 0 {
			continue
		}
		u := findUserR(ctx, state.ParentResources, resource.MakeIdentifiers(idf))
		n.relUsr[idf] = u.ID
	}

	// This is not possible, but let's do it anyway
	if state.Conflicting {
		return nil
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != "" {
			t := toTime(ts.CreatedAt)
			if t != nil {
				res.CreatedAt = *t
			}
		}
		if ts.UpdatedAt != "" {
			res.UpdatedAt = toTime(ts.UpdatedAt)
		}
		if ts.DeletedAt != "" {
			res.DeletedAt = toTime(ts.DeletedAt)
		}
		if ts.ArchivedAt != "" {
			res.ArchivedAt = toTime(ts.ArchivedAt)
		}
	}

	// Userstamps
	us := n.res.Userstamps()
	if us != nil {
		if us.CreatedBy != "" {
			res.CreatorID = n.relUsr[us.CreatedBy]
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create fresh messagingChannel
	if !exists {
		return store.CreateMessagingChannel(ctx, s, res)
	}

	// Update existing messagingChannel
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeMessagingChannels(n.ch, res)

	case resource.MergeRight:
		res = mergeMessagingChannels(res, n.ch)
	}

	err = store.UpdateMessagingChannel(ctx, s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

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
	ap = findMessagingChannelR(rr, ii)
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

// findMessagingChannelR looks for the ch in the resource set
func findMessagingChannelR(rr resource.InterfaceSet, ii resource.Identifiers) (ap *types.Channel) {
	var chRes *resource.MessagingChannel

	rr.Walk(func(r resource.Interface) error {
		ar, ok := r.(*resource.MessagingChannel)
		if !ok {
			return nil
		}

		if ar.Identifiers().HasAny(ii) {
			chRes = ar
		}

		return nil
	})

	// Found it
	if chRes != nil {
		return chRes.Res
	}

	return nil
}
