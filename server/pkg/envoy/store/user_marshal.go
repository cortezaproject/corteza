package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

func NewUserFromResource(res *resource.User, cfg *EncoderConfig) resourceState {
	return &user{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *user) Prepare(ctx context.Context, pl *payload) (err error) {
	// Reset old identifiers
	n.res.Res.ID = 0

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
		if err = store.CreateUser(ctx, pl.s, res); err != nil {
			return
		}
		return n.membership(ctx, pl.s, res, pl.state.ParentResources, exists)
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

	if err = store.UpdateUser(ctx, pl.s, res); err != nil {
		return err
	}
	n.res.Res = res

	return n.membership(ctx, pl.s, res, pl.state.ParentResources, exists)
}

func (n *user) membership(ctx context.Context, s store.Storer, res *types.User, pp resource.InterfaceSet, exists bool) (err error) {
	// find all roles
	roles := make([]uint64, 0, 10)
	var r *types.Role
	for _, m := range n.res.RoleMembership {
		r, err = findRole(ctx, s, pp, m)
		if err != nil {
			return
		}

		roles = append(roles, r.ID)
	}

	// update
	// @todo some smarter diff calculations; should be fine for now but could be improved.
	var mm types.RoleMemberSet
	if exists {
		mm, _, err = s.SearchRoleMembers(ctx, types.RoleMemberFilter{UserID: n.u.ID})
		if err != nil {
			return
		}
	}

	for _, m := range mm {
		if err = store.DeleteRoleMember(ctx, s, m); err != nil {
			return
		}
	}
	for _, r := range roles {
		if err = store.CreateRoleMember(ctx, s, &types.RoleMember{UserID: res.ID, RoleID: r}); err != nil {
			return
		}
	}

	return
}
