package provision

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	uConfigFn func(context.Context, store.Storer, *zap.Logger) bool
	uConfig   struct {
		dir string
		fn  uConfigFn
	}
)

// provisionPartialAuthClients checks for a specific set of auth client rbac rules
func provisionPartialAuthClients(ctx context.Context, s store.Storer, log *zap.Logger) bool {
	set, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})

	if err != nil {
		log.Warn("could not make a partial import of templates", zap.Error(err))
		return false
	}

	set, _ = set.Filter(func(r *rbac.Rule) (bool, error) {
		// check only auth client rbac rules
		if r.Resource.String() != "system:auth-client:*" {
			return false, nil
		}
		return true, nil
	})

	return len(set) == 0
}

// provisionPartialTemplates checks if any templates are in the store at all
func provisionPartialTemplates(ctx context.Context, s store.Storer, log *zap.Logger) bool {
	set, _, err := store.SearchTemplates(ctx, s, types.TemplateFilter{Deleted: filter.StateInclusive})
	if err != nil {
		log.Warn("could not make a partial import of templates", zap.Error(err))
	}

	return err != nil || len(set) == 0
}
