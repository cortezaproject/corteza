package rbac

import (
	"context"
)

type (
	// All
	rbacRulesStore interface {
		SearchRbacRules(ctx context.Context, f RuleFilter) (RuleSet, RuleFilter, error)
		UpsertRbacRule(ctx context.Context, rr ...*Rule) error
		DeleteRbacRule(ctx context.Context, rr ...*Rule) error
		TruncateRbacRules(ctx context.Context) error
	}
)
