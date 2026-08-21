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

	return baseMarkerMissing(rr, pp, log)
}

// baseMarker is a grant the base config makes to one of its own roles, chosen
// so that its absence dates the install.
//
// It is scoped to the role on purpose. "Does any rule of this shape exist" is
// satisfied by any ad-hoc role a deployment happens to have made — a probe
// role carrying one user-group rule was enough to convince this check that an
// install had permissions it had never been given — and then the re-import it
// is supposed to trigger never happens.
type baseMarker struct {
	what      string
	role      string
	resource  string
	operation string
}

// Deliberately a named few of the grants the base config makes rather than all
// of them: a re-import restates the whole file, so a deployment that has
// removed one of the base grants would get it back. Once per newly-introduced
// model is a defensible price for the feature being reachable at all; once per
// anything would not be.
//
// Adding a grant to the base config means adding a row here, or it reaches new
// installs only.
var baseMarkers = []baseMarker{
	// The whole project feature shipped without a single rule, so it was usable
	// only by super-admin, which bypasses RBAC entirely.
	{"project permissions", "admin", "corteza::system:project/*", "read"},
	// Applications predate webapp access control, so every install already holds
	// `read` — `access` is the operation that dates this one.
	{"webapp access permissions", "admin", "corteza::system:application/*", "access"},
	// User groups, connections and data sources have admin screens that were
	// never granted to the role those screens belong to.
	{"connection permissions", "admin", "corteza::system/", "connections.search"},
	// Labels and Corredor scripts were listed by endpoints that checked nothing
	// at all, so no rule for them existed anywhere.
	{"label permissions", "admin", "corteza::system/", "labels.search"},
}

// baseMarkerMissing reports whether the base config should be re-imported,
// which it should exactly when one of its markers is not already granted.
func baseMarkerMissing(rr types.RoleSet, pp rbac.RuleSet, log *zap.Logger) bool {
	held := make(map[baseMarker]bool, len(pp))

	for _, r := range pp {
		if r.Access != rbac.Allow {
			continue
		}

		role := rr.FindByID(r.RoleID)
		if role == nil {
			continue
		}

		held[baseMarker{role: role.Handle, resource: r.Resource, operation: r.Operation}] = true
	}

	for _, m := range baseMarkers {
		probe := baseMarker{role: m.role, resource: m.resource, operation: m.operation}
		if held[probe] {
			continue
		}

		log.Info("base config carries "+m.what+" this install has never seen; re-importing it",
			zap.String("role", m.role),
			zap.String("resource", m.resource),
			zap.String("operation", m.operation))

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
