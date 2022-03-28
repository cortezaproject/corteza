package provision

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

func migratePre202109Roles(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	const (
		obsoleteEveryoneID uint64 = 1
		obsoleteAdminsID   uint64 = 2
	)

	log.Info("migrating pre-2021.9 roles")
	m, err := loadRoles(ctx, s)
	if err != nil {
		return
	}

	// let's see if everyone role is still here:
	if m["everyone"] != nil && m["everyone"].ID == obsoleteEveryoneID {
		log.Info("migrating 'everyone' role to new ID")

		// everyone role still present and it is using "hardcoded" ID
		// we can remove it
		if err = store.DeleteRoleByID(ctx, s, obsoleteEveryoneID); err != nil {
			return
		}

		// transfer all rbac rules
		if err = s.TransferRbacRules(ctx, obsoleteEveryoneID, m["authenticated"].ID); err != nil {
			return
		}
	}

	// let's see if admin role is still here:
	if m["admins"] != nil && m["admins"].ID == obsoleteAdminsID {
		log.Info("migrating 'admins' role to new ID and renaming it to 'admin'")

		// everyone role still present and it is using "hardcoded" ID
		// we can remove it
		m["admins"].ID = id.Next()
		m["admins"].UpdatedAt = now()
		m["admins"].Handle = "admin"
		m["admins"].Name = "Administrator"

		if err = store.DeleteRoleByID(ctx, s, obsoleteAdminsID); err != nil {
			return
		}

		if err = store.CreateRole(ctx, s, m["admins"]); err != nil {
			return
		}

		// transfer all rbac rules
		if err = store.TransferRoleMembers(ctx, s, obsoleteAdminsID, m["admins"].ID); err != nil {
			return
		}

		// transfer all rbac rules
		if err = store.TransferRbacRules(ctx, s, obsoleteAdminsID, m["admins"].ID); err != nil {
			return
		}
	}

	return
}
