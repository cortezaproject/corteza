package provision

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

// SystemRoles creates system roles
func SystemRoles(ctx context.Context, log *zap.Logger, s store.Storer) (rr []*types.Role, err error) {
	rr = types.RoleSet{
		&types.Role{
			Name:   "Super administrator",
			Handle: "super-admin",
			Meta: &types.RoleMeta{
				Description: "Super admin is a 'bypass' role that auto-allows all operations to it's members",
				Context:     nil,
			},
		},

		&types.Role{
			Name:   "Authenticated",
			Handle: "authenticated",
			Meta: &types.RoleMeta{
				Description: "Authenticated role is auto-assigned to all authenticated sessions",
				Context:     nil,
			},
		},

		&types.Role{
			Name:   "Anonymous",
			Handle: "anonymous",
			Meta: &types.RoleMeta{
				Description: "Authenticated role is auto-assigned to all non-authenticated sessions",
				Context:     nil,
			}},
	}

	m, err := loadRoles(ctx, s)
	if err != nil {
		return
	}

	for i := range rr {
		r := rr[i]
		if m[r.Handle] == nil {
			log.Info("creating system role", zap.String("handle", r.Handle))
			// this is a new role
			r.ID = id.Next()
			r.CreatedAt = *now()

			m[r.Handle] = r
		} else {
			log.Info("updating system role", zap.String("handle", r.Handle))
			// use existing role
			rr[i] = m[r.Handle]

			// make sure it's not deleted or archived
			// and leave other props as they are
			r.DeletedAt = nil
			r.ArchivedAt = nil

		}
	}

	if err := store.UpsertRole(ctx, s, rr...); err != nil {
		return nil, fmt.Errorf("failed to provision system roles: %w", err)
	}

	return
}

func loadRoles(ctx context.Context, s store.Roles) (m map[string]*types.Role, err error) {
	var (
		f = types.RoleFilter{
			Archived: filter.StateInclusive,
			Deleted:  filter.StateInclusive,
		}
	)

	m = make(map[string]*types.Role)

	if set, _, err := store.SearchRoles(ctx, s, f); err == nil {
		for _, r := range set {
			m[r.Handle] = r
		}
	}

	return
}
