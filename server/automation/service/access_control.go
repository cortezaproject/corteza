package service

import (
	"context"

	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

// CanImpersonateUser checks if current user can impersonate the given user
func (svc accessControl) CanImpersonateUser(ctx context.Context, u *sysTypes.User) bool {
	return svc.can(ctx, "impersonate", u)
}
