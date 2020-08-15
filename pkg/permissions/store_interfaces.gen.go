package permissions

import (
	"context"
)

type (
	rbacRulesStore interface {
		SearchRbacRules(ctx context.Context, f RuleFilter) (RuleSet, RuleFilter, error)
		CreateRbacRule(ctx context.Context, rr ...*Rule) error
		UpdateRbacRule(ctx context.Context, rr ...*Rule) error
		PartialUpdateRbacRule(ctx context.Context, onlyColumns []string, rr ...*Rule) error
		RemoveRbacRule(ctx context.Context, rr ...*Rule) error
		RemoveRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error

		TruncateRbacRules(ctx context.Context) error
	}
)
