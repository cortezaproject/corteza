package provision

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/sass"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

// SystemUsers creates or updates system users
func SystemUsers(ctx context.Context, log *zap.Logger, s store.Users) (uu []*types.User, err error) {
	uu = types.UserSet{
		&types.User{
			Email:  "provision@human.local",
			Name:   "Human Provisioner",
			Handle: auth.ProvisionUserHandle,
			Kind:   types.SystemUser,
		},
		&types.User{
			Email:  "service@human.local",
			Name:   "Human Service",
			Handle: auth.ServiceUserHandle,
			Kind:   types.SystemUser,
		},
		&types.User{
			Email:  "federation@human.local",
			Name:   "Human Federation",
			Handle: auth.FederationUserHandle,
			Kind:   types.SystemUser,
		},
	}

	m, err := loadSystemUsers(ctx, s)
	if err != nil {
		return
	}

	for i := range uu {
		u := uu[i]
		if m[u.Handle] == nil {
			// this is a new user
			u.ID = id.Next()
			u.CreatedAt = *now()

			if err := store.UpsertUser(ctx, s, u); err != nil {
				return nil, fmt.Errorf("failed to provision system user %s: %w", u.Handle, err)
			}

			log.Info("creating system user", zap.String("handle", u.Handle), logger.Uint64("ID", u.ID))
		} else {

			u.ID = m[u.Handle].ID

			// There is no need to update system users if they are unchanged
			if m[u.Handle].UpdatedAt == nil &&
				m[u.Handle].SuspendedAt == nil &&
				m[u.Handle].DeletedAt == nil {
				continue
			}

			// Make sure all values are as they should be
			u.CreatedAt = m[u.Handle].CreatedAt
			u.Email = m[u.Handle].Email
			u.Name = m[u.Handle].Name
			u.UpdatedAt = nil
			u.SuspendedAt = nil
			u.DeletedAt = nil

			if err := store.UpsertUser(ctx, s, u); err != nil {
				return nil, fmt.Errorf("failed to provision system user %s: %w", u.Handle, err)
			}

			log.Info("updating system user", zap.String("handle", u.Handle), logger.Uint64("ID", u.ID))
		}
	}

	return
}

func loadSystemUsers(ctx context.Context, s store.Users) (m map[string]*types.User, err error) {
	var (
		f = types.UserFilter{
			Suspended: filter.StateInclusive,
			Deleted:   filter.StateInclusive,
			Kind:      types.SystemUser,
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

func setUsersTheme(ctx context.Context, log *zap.Logger, s store.Users) (err error) {
	log.Info("provision start")
	defer log.Info("provision end")

	var (
		f = types.UserFilter{
			Deleted: filter.StateInclusive,
		}

		users []*types.User
	)

	if set, _, err := store.SearchUsers(ctx, s, f); err == nil {
		for _, r := range set {
			if r.Meta.Theme == "" {
				r.Meta.Theme = sass.LightTheme

				users = append(users, r)
			}
		}
	}

	err = store.UpdateUser(ctx, s, users...)
	if err != nil {
		log.Error("error provisioning users theme", zap.Error(err))
		return err
	}

	return nil
}
