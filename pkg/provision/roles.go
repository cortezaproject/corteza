package provision

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

func roles(ctx context.Context, log *zap.Logger, s store.Storer) error {
	const (
		obsoleteEveryoneID uint64 = 1
		obsoleteAdminsID   uint64 = 2
	)

	var (
		f = types.RoleFilter{
			Archived: filter.StateInclusive,
			Deleted:  filter.StateInclusive,
		}

		now = time.Now().Round(time.Second)

		rr = types.RoleSet{
			&types.Role{
				Name:   "Super administrator",
				Handle: "super-admin",
				Meta: &types.RoleMeta{
					Description: "Super admin is a 'bypass' role that allows all actions it's members",
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
	)

	m := make(map[string]*types.Role)
	if set, _, err := store.SearchRoles(ctx, s, f); err != nil {
		return err
	} else {
		for _, r := range set {
			m[r.Handle] = r
		}
	}

	for i := range rr {
		r := rr[i]
		if m[r.Handle] == nil {
			log.Info("creating role", zap.String("handle", r.Handle))
			// this is a new role
			r.ID = id.Next()
			r.CreatedAt = now

			m[r.Handle] = r
		} else {
			log.Info("updating role", zap.String("handle", r.Handle))
			// use existing role
			rr[i] = m[r.Handle]

			// make sure it's not deleted or archived
			// and left other props as they are
			r.DeletedAt = nil
			r.ArchivedAt = nil

		}
	}

	if err := store.UpsertRole(ctx, s, rr...); err != nil {
		return err
	}

	// let's see if everyone role is still here:
	if m["everyone"] != nil && m["everyone"].ID == obsoleteEveryoneID {
		log.Info("migrating 'everyone' role")

		// everyone role still present and it is using "hardcoded" ID
		// we can remove it
		if err := store.DeleteRoleByID(ctx, s, obsoleteEveryoneID); err != nil {
			return err
		}

		// transfer all rbac rules
		if err := s.TransferRbacRules(ctx, obsoleteEveryoneID, m["authenticated"].ID); err != nil {
			return nil
		}
	}

	// let's see if everyone role is still here:
	if m["admins"] != nil && m["admins"].ID == obsoleteAdminsID {
		log.Info("migrating 'admins' role")

		// everyone role still present and it is using "hardcoded" ID
		// we can remove it
		m["admins"].ID = id.Next()
		m["admins"].UpdatedAt = &now

		if err := store.DeleteRoleByID(ctx, s, obsoleteAdminsID); err != nil {
			return err
		}

		if err := store.CreateRole(ctx, s, m["admins"]); err != nil {
			return err
		}

		// transfer all rbac rules
		if err := s.TransferRoleMembers(ctx, obsoleteAdminsID, m["admins"].ID); err != nil {
			return nil
		}

		// transfer all rbac rules
		if err := s.TransferRbacRules(ctx, obsoleteAdminsID, m["admins"].ID); err != nil {
			return nil
		}
	}

	return nil
}
