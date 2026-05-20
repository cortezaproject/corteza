package service

import (
	"context"
	"fmt"

	pkgAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// ImpersonateServiceAccount builds a context authenticated as the given user
// with their full role memberships loaded from the store. Used by chatbot
// surfaces (widget, preview) where the caller has no identity so RBAC checks
// downstream rely entirely on this synthesized auth context.
func ImpersonateServiceAccount(ctx context.Context, s store.Storer, userID uint64) context.Context {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	mm, _, err := store.SearchRoleMembers(svcCtx, s, types.RoleMemberFilter{
		Resource: fmt.Sprintf("corteza::system:user/%d", userID),
	})
	var roles []uint64
	if err == nil {
		for _, m := range mm {
			roles = append(roles, m.RoleID)
		}
	}
	return pkgAuth.SetIdentityToContext(ctx, pkgAuth.Authenticated(userID, roles...))
}
