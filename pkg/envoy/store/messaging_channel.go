package store

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	messagingChannelState struct {
		cfg *EncoderConfig

		res *resource.MessagingChannel
		ch  *types.Channel
	}
)

func NewMessagingChannelState(res *resource.MessagingChannel, cfg *EncoderConfig) resourceState {
	return &messagingChannelState{
		cfg: cfg,
		res: res,
	}
}

func (n *messagingChannelState) Prepare(ctx context.Context, s store.Storer, rs *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
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

	// This is not possible, but let's do it anyway
	if state.Conflicting {
		return nil
	}

	// Create fresh messagingChannel
	if !exists {
		return store.CreateMessagingChannel(ctx, s, res)
	}

	// Update existing messagingChannel
	switch n.cfg.OnExisting {
	case Skip:
		return nil

	case MergeLeft:
		res = mergeMessagingChannels(n.ch, res)

	case MergeRight:
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

	q := gf.handle
	if q == "" {
		q = gf.name
	}

	if q != "" {
		var aa types.ChannelSet
		aa, _, err = store.SearchMessagingChannels(ctx, s, types.ChannelFilter{Query: q})
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}
		if len(aa) > 0 {
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
