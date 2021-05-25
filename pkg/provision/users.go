package provision

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

// SystemUsers creates or updates system users
func SystemUsers(ctx context.Context, log *zap.Logger, s store.Users) (uu []*types.User, err error) {
	var (
		now = time.Now().Round(time.Second)
	)

	uu = types.UserSet{
		&types.User{
			Email:  "provision@corteza.local",
			Name:   "Corteza Provisioner",
			Handle: auth.ProvisionUserHandle,
			Kind:   types.SystemUser,
		},
		&types.User{
			Email:  "service@corteza.local",
			Name:   "Corteza Service",
			Handle: auth.ServiceUserHandle,
			Kind:   types.SystemUser,
		},
		&types.User{
			Email:  "federation@corteza.local",
			Name:   "Corteza Federation",
			Handle: auth.FederationUserHandle,
			Kind:   types.SystemUser,
		},
	}

	m, err := loadUsers(ctx, s)
	if err != nil {
		return
	}

	for i := range uu {
		u := uu[i]
		if m[u.Handle] == nil {
			log.Info("creating user", zap.String("handle", u.Handle))
			// this is a new user
			u.ID = id.Next()
			u.CreatedAt = now

			if err := store.UpsertUser(ctx, s, u); err != nil {
				return nil, fmt.Errorf("failed to provision user %s: %w", u.Handle, err)
			}
		} else {
			u.ID = m[u.Handle].ID
			u.Email = m[u.Handle].Email
			u.Name = m[u.Handle].Name
			u.SuspendedAt = nil
			u.DeletedAt = nil

			if err := store.UpsertUser(ctx, s, u); err != nil {
				return nil, fmt.Errorf("failed to provision user %s: %w", u.Handle, err)
			}

		}
	}

	return
}

func loadUsers(ctx context.Context, s store.Users) (m map[string]*types.User, err error) {
	var (
		f = types.UserFilter{
			Suspended: filter.StateInclusive,
			Deleted:   filter.StateInclusive,
		}
	)

	m = make(map[string]*types.User)
	if set, _, err := store.SearchUsers(ctx, s, f); err == nil {
		for _, r := range set {
			m[r.Handle] = r
		}
	}

	return
}
