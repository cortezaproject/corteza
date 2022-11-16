package store

import (
	"context"
	"net/mail"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	user struct {
		cfg *EncoderConfig

		res *resource.User
		u   *types.User
	}
)

// mergeUsers merges b into a, prioritising a
func mergeUsers(a, b *types.User) *types.User {
	c := *a

	if c.Username == "" {
		c.Username = b.Username
	}
	if c.Email == "" {
		c.Email = b.Email
	}
	if c.Name == "" {
		c.Name = b.Name
	}
	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Kind == "" {
		c.Kind = b.Kind
	}
	if c.Meta == nil {
		c.Meta = b.Meta
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}
	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}
	if c.SuspendedAt == nil {
		c.SuspendedAt = b.SuspendedAt
	}
	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return &c
}

// findUser looks for the user in the resources & the store
//
// Provided resources are prioritized.
func findUser(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (u *types.User, err error) {
	u = resource.FindUser(rr, ii)
	if u != nil {
		return u, nil
	}

	return findUserStore(ctx, s, makeGenericFilter(ii))
}

// findUserStore looks for the user in the store
func findUserStore(ctx context.Context, s store.Storer, gf genericFilter) (u *types.User, err error) {
	if gf.id > 0 {
		u, err = store.LookupUserByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if u != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		// email
		if _, err = mail.ParseAddress(i); err == nil {
			u, err = store.LookupUserByEmail(ctx, s, i)
			if err == store.ErrNotFound {
				return nil, nil
			} else if err != nil {
				return nil, err
			}
		}

		// Handle & username
		u, err = store.LookupUserByHandle(ctx, s, i)
		if err == store.ErrNotFound {
			u, err = store.LookupUserByUsername(ctx, s, i)
		}
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if u != nil {
			return
		}
	}

	return nil, nil
}
