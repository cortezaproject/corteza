package rbac

import (
	"context"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	// All
	rbacRulesStore interface {
		SearchRbacRules(ctx context.Context, f RuleFilter) (RuleSet, RuleFilter, error)
		UpsertRbacRule(ctx context.Context, rr ...*Rule) error
		DeleteRbacRule(ctx context.Context, rr ...*Rule) error
		TruncateRbacRules(ctx context.Context) error

		// @todo this isn't ok since we're referencing sys types
		SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error)
	}

	rbacRoleStore interface {
		SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error)
	}
)
