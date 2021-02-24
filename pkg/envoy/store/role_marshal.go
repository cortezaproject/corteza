package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func NewRoleFromResource(res *resource.Role, cfg *EncoderConfig) resourceState {
	return &role{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

// Prepare prepares the role to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *role) Prepare(ctx context.Context, pl *payload) (err error) {
	n.rl, err = findRoleS(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.rl != nil {
		n.res.Res.ID = n.rl.ID
	}
	return nil
}

// Encode encodes the composeChart to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *role) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.rl != nil && n.rl.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.rl.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// This is not possible, but let's do it anyway
	if pl.state.Conflicting {
		return nil
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
		if ts.ArchivedAt != nil {
			res.ArchivedAt = ts.ArchivedAt.T
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

	// Create a fresh role
	if !exists {
		return store.CreateRole(ctx, pl.s, res)
	}

	// Update existing roles
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeRoles(n.rl, res)

	case resource.MergeRight:
		res = mergeRoles(res, n.rl)
	}

	err = store.UpdateRole(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
