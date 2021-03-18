package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	role struct {
		cfg *EncoderConfig

		res *resource.Role
		rl  *types.Role
	}
)

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

// findRole looks for the role in the resources & the store
//
// Provided resources are prioritized.
func findRole(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (rl *types.Role, err error) {
	rl = resource.FindRole(rr, ii)
	if rl != nil {
		return rl, nil
	}

	return findRoleStore(ctx, s, makeGenericFilter(ii))
}

// findRoleStore looks for the role in the store
func findRoleStore(ctx context.Context, s store.Storer, gf genericFilter) (rl *types.Role, err error) {
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
