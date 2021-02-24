package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func newMessagingChannelFromResource(res *resource.MessagingChannel, cfg *EncoderConfig) resourceState {
	return &messagingChannel{
		cfg: mergeConfig(cfg, res.Config()),
		res: res,
	}
}

// Prepare prepares the messagingChannel to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *messagingChannel) Prepare(ctx context.Context, pl *payload) (err error) {
	// Sys users
	n.relUsr = make(map[string]uint64)

	// Get the existing channel
	n.ch, err = findMessagingChannelS(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.ch != nil {
		n.res.Res.ID = n.ch.ID
	}
	return nil
}

// Encode encodes the messagingChannel to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *messagingChannel) Encode(ctx context.Context, pl *payload) (err error) {
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
	us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, n.res.Userstamps())
	if err != nil {
		return err
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
		if ts.ArchivedAt != nil {
			res.ArchivedAt = ts.ArchivedAt.T
		}
	}

	// Userstamps
	if us != nil {
		if us.CreatedBy != nil {
			res.CreatorID = us.CreatedBy.UserID
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
		return store.CreateMessagingChannel(ctx, pl.s, res)
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

	err = store.UpdateMessagingChannel(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
