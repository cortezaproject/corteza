package provision

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
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
//
// The "already done" test is per-install, not per-file-version, so a base
// config that gains grants for a NEW resource type would never reach an
// existing database. That is how projects shipped without a single
// corteza::system:project rule anywhere: the whole feature was usable only by
// super-admin, which bypasses RBAC entirely. Hence the per-resource-type probe
// below, in the shape provisionPartialAuthClients already established — it
// fires once on an install that predates the resource type, then never again.
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

	// Deliberately keyed on the project resource type alone rather than on
	// every type the base config mentions: a re-import restates the whole
	// file, so a deployment that has removed one of the base grants would get
	// it back. Once per newly-introduced resource type is a defensible price
	// for the feature being reachable at all; once per anything would not be.
	for _, r := range pp {
		if rbac.ResourceType(r.Resource) == types.ProjectResourceType {
			return false
		}
	}

	log.Info("base config carries project permissions this install has never seen; re-importing it")
	return true
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
