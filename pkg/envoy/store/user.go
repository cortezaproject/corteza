package store

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	userState struct {
		cfg *EncoderConfig

		res *resource.User
		u   *types.User
	}
)

func NewUserState(res *resource.User, cfg *EncoderConfig) resourceState {
	return &userState{
		cfg: cfg,

		res: res,
	}
}

func (n *userState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
	}

	// Try to get the original user
	// @todo make filtering more flexible (email, username, ...)
	n.u, err = findUserS(ctx, s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.u != nil {
		n.res.Res.ID = n.u.ID
	}
	return nil
}

func (n *userState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	res := n.res.Res
	exists := n.u != nil && n.u.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.u.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
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
		if ts.SuspendedAt != "" {
			res.SuspendedAt = toTime(ts.SuspendedAt)
		}
	}

	// Create a fresh user
	if !exists {
		return store.CreateUser(ctx, s, res)
	}

	// Update existing user
	switch n.cfg.OnExisting {
	case Skip:
		return nil

	case MergeLeft:
		res = mergeUsers(n.u, res)

	case MergeRight:
		res = mergeUsers(res, n.u)
	}

	err = store.UpdateUser(ctx, s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

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
	if *c.Meta == (types.UserMeta{}) {
		c.Meta = b.Meta
	}

	return &c
}

// findUserRS looks for the user in the resources & the store
//
// Provided resources are prioritized.
func findUserRS(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (u *types.User, err error) {
	u = findUserR(ctx, rr, ii)
	if u != nil {
		return u, nil
	}

	return findUserS(ctx, s, makeGenericFilter(ii))
}

// findUserS looks for the user in the store
func findUserS(ctx context.Context, s store.Storer, gf genericFilter) (u *types.User, err error) {
	if gf.id > 0 {
		u, err = store.LookupUserByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if u != nil {
			return
		}
	}

	if gf.handle != "" {
		u, err = store.LookupUserByHandle(ctx, s, gf.handle)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if u != nil {
			return
		}
	}

	return nil, nil
}

// findUserR looks for the user in the resources
func findUserR(ctx context.Context, rr resource.InterfaceSet, ii resource.Identifiers) (u *types.User) {
	var uRes *resource.User

	rr.Walk(func(r resource.Interface) error {
		ur, ok := r.(*resource.User)
		if !ok {
			return nil
		}

		if ur.Identifiers().HasAny(ii) {
			uRes = ur
		}
		return nil
	})

	// Found it
	if uRes != nil {
		return uRes.Res
	}

	return nil
}

func userErrUnresolved(ii resource.Identifiers) error {
	return fmt.Errorf("user unresolved %v", ii.StringSlice())
}
