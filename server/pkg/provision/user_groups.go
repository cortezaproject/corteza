package provision

import (
	"context"
	"os"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

func setDefaultUserGroupRefs(ctx context.Context, log *zap.Logger, s store.Storer, authOpt options.AuthOpt) (err error) {
	h, ok := os.LookupEnv("AUTH_DEFAULT_USER_GROUP")
	if !ok {
		h = "default-root"
	}
	ug, err := store.LookupUserGroupByHandle(ctx, s, h)
	if err != nil {
		return
	}

	// ---
	// Users
	// ---

	uu, _, err := store.SearchUsers(ctx, s, types.UserFilter{
		Kind: types.NormalUser,
	})
	if err != nil {
		return
	}

	for _, u := range uu {
		if u.UserGroupID != 0 {
			continue
		}

		u.UserGroupID = ug.ID
		err = store.UpdateUser(ctx, s, u)
		if err != nil {
			return
		}
	}

	// ---
	// auth clients
	// ---

	cc, _, err := store.SearchAuthClients(ctx, s, types.AuthClientFilter{})
	if err != nil {
		return
	}

	for _, c := range cc {
		if c.Security.UserGroup != 0 {
			continue
		}

		c.Security.UserGroup = ug.ID
		err = store.UpdateAuthClient(ctx, s, c)
		if err != nil {
			return
		}
	}

	return
}
