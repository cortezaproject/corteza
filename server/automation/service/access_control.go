// Hand-written companion to access_control.gen.go.
package service

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/rbac"
	sysTypes "github.com/crusttech/human/server/system/types"
)

// CanImpersonateUser checks if current user can impersonate the given user
func (svc accessControl) CanImpersonateUser(ctx context.Context, u *sysTypes.User) bool {
	return svc.can(ctx, "impersonate", u)
}

// EffectiveFor answers the caller's effective permissions on one resource,
// addressed by its RBAC resource string.
//
// The generated Effective() takes typed resources, and a REST caller has only
// the string. A context role (owner, creator) is evaluated against the record
// itself, so the resource has to be loaded before it can be answered at all —
// which is what separates this from asking about the component.
func (svc accessControl) EffectiveFor(ctx context.Context, resource string) (ee rbac.EffectiveSet, err error) {
	if err = rbacResourceValidator(resource); err != nil {
		return nil, fmt.Errorf("can not use resource %q: %w", resource, err)
	}

	res, err := svc.resourceLoader(ctx, resource)
	if err != nil {
		return nil, err
	}

	return svc.Effective(ctx, res), nil
}
