package provision

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

type (
	uConfigFn func(context.Context, store.Storer, *zap.Logger) bool
	uConfig   struct {
		dir string
		fn  uConfigFn
	}
)

// provisionPartialBase check for roles and permissions
//
// It checks if there are any roles and any RBAC rules. If there are, we assume
// the provision for the base dir was already done.
func provisionPartialBase(ctx context.Context, s store.Storer, log *zap.Logger) bool {
	rr, _, err := store.SearchRoles(ctx, s, types.RoleFilter{Deleted: filter.StateInclusive})
	if err != nil {
		log.Warn("could not make a partial import of base: roles", zap.Error(err))
		return false
	}
	if len(rr) == 0 {
		return true
	}

	pp, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	if err != nil {
		log.Warn("could not make a partial import of base: permissions", zap.Error(err))
		return false
	}
	if len(pp) == 0 {
		return true
	}

	return false
}

// provisionPartialAuthClients checks for a specific set of auth client rbac rules
func provisionPartialAuthClients(ctx context.Context, s store.Storer, log *zap.Logger) bool {
	set, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})

	if err != nil {
		log.Warn("could not make a partial import of templates", zap.Error(err))
		return false
	}

	for _, r := range set {
		if rbac.ResourceType(r.Resource) == types.AuthClientResourceType {
			return false
		}
	}

	return true
}

// provisionPartialTemplates checks if any templates are in the store at all
func provisionPartialTemplates(ctx context.Context, s store.Storer, log *zap.Logger) bool {
	set, _, err := store.SearchTemplates(ctx, s, types.TemplateFilter{Deleted: filter.StateInclusive})
	if err != nil {
		log.Warn("could not make a partial import of templates", zap.Error(err))
	}

	return err != nil || len(set) == 0
}
