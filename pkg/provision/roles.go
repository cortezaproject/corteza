package provision

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func roles(ctx context.Context, s store.Storer) error {
	if set, _, err := store.SearchRoles(ctx, s, types.RoleFilter{}); err != nil {
		return err
	} else if len(set) > 0 {
		return nil
	}

	now := time.Now().Round(time.Second)

	rr := types.RoleSet{
		&types.Role{ID: rbac.AdminsRoleID, Name: "Administrators", Handle: "admins"},
		&types.Role{ID: rbac.EveryoneRoleID, Name: "Everyone", Handle: "everyone"},
	}

	err := rr.Walk(func(r *types.Role) error {
		r.CreatedAt = now
		return store.CreateRole(ctx, s, r)
	})

	return err
}
