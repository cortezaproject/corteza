package store

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleState struct {
		cfg *EncoderConfig

		res *resource.Role
		rl  *types.Role
	}
)

func NewRole(res *resource.Role, cfg *EncoderConfig) resourceState {
	return &roleState{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *roleState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
	}

	n.rl, err = findRoleS(ctx, s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.rl != nil {
		n.res.Res.ID = n.rl.ID
	}
	return nil
}

func (n *roleState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	rl := n.res.Res
	exists := n.rl != nil && n.rl.ID > 0

	// Determine the ID
	if rl.ID <= 0 && exists {
		rl.ID = n.rl.ID
	}
	if rl.ID <= 0 {
		rl.ID = NextID()
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
				rl.CreatedAt = *t
			}
		}
		if ts.UpdatedAt != "" {
			rl.UpdatedAt = toTime(ts.UpdatedAt)
		}
		if ts.DeletedAt != "" {
			rl.DeletedAt = toTime(ts.DeletedAt)
		}
		if ts.ArchivedAt != "" {
			rl.ArchivedAt = toTime(ts.ArchivedAt)
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
		return store.CreateRole(ctx, s, rl)
	}

	// Update existing roles
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		rl = mergeRoles(n.rl, rl)

	case resource.MergeRight:
		rl = mergeRoles(rl, n.rl)
	}

	err = store.UpdateRole(ctx, s, rl)
	if err != nil {
		return err
	}

	n.res.Res = rl
	return nil
}

// mergeRoles merges b into a, prioritising a
func mergeRoles(a, b *types.Role) *types.Role {
	c := *a

	if c.Name == "" {
		c.Name = b.Name
	}
	if c.Handle == "" {
		c.Handle = b.Handle
	}

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

// findRoleRS looks for the role in the resources & the store
//
// Provided resources are prioritized.
func findRoleRS(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (rl *types.Role, err error) {
	rl = findRoleR(rr, ii)
	if rl != nil {
		return rl, nil
	}

	return findRoleS(ctx, s, makeGenericFilter(ii))
}

// findRoleS looks for the role in the store
func findRoleS(ctx context.Context, s store.Storer, gf genericFilter) (rl *types.Role, err error) {
	if gf.id > 0 {
		rl, err = store.LookupRoleByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if rl != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		rl, err = store.LookupRoleByHandle(ctx, s, i)
		if err == store.ErrNotFound {
			var rr types.RoleSet
			rr, _, err = store.SearchRoles(ctx, s, types.RoleFilter{
				Name: i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(rr) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(rr) == 1 {
				rl = rr[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if rl != nil {
			return
		}
	}

	return nil, nil
}

// findRoleR looks for the role in the resources
func findRoleR(rr resource.InterfaceSet, ii resource.Identifiers) (rl *types.Role) {
	var rlRes *resource.Role

	rr.Walk(func(r resource.Interface) error {
		rr, ok := r.(*resource.Role)
		if !ok {
			return nil
		}

		if rr.Identifiers().HasAny(ii) {
			rlRes = rr
		}
		return nil
	})

	// Found it
	if rlRes != nil {
		return rlRes.Res
	}

	return nil
}

func roleErrUnresolved(ii resource.Identifiers) error {
	return fmt.Errorf("role unresolved %v", ii.StringSlice())
}
