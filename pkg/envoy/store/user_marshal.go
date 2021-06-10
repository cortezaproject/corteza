package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func NewUserFromResource(res *resource.User, cfg *EncoderConfig) resourceState {
	return &user{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *user) Prepare(ctx context.Context, pl *payload) (err error) {
	if n.cfg.IgnoreStore {
		n.res.Res.ID = 0
		return nil
	}

	// Try to get the original user
	n.u, err = findUserStore(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.u != nil {
		n.res.Res.ID = n.u.ID
	}
	return nil
}

func (n *user) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.u != nil && n.u.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.u.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

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
		if ts.SuspendedAt != nil {
			res.SuspendedAt = ts.SuspendedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh user
	if !exists {
		return store.CreateUser(ctx, pl.s, res)
	}

	// Update existing user
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeUsers(n.u, res)

	case resource.MergeRight:
		res = mergeUsers(res, n.u)
	}

	err = store.UpdateUser(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
