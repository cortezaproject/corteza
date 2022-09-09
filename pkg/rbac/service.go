package rbac

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"go.uber.org/zap"
)

type (
	service struct {
		l      *sync.RWMutex
		logger *zap.Logger

		//  service will flush values on TRUE or just reload on FALSE
		f chan bool

		rules   RuleSet
		indexed OptRuleSet

		roles []*Role

		store rbacRulesStore
	}

	// RuleFilter is a dummy struct to satisfy store codegen
	RuleFilter struct {
		Limit uint
	}

	RoleSettings struct {
		Bypass        []uint64
		Authenticated []uint64
		Anonymous     []uint64
	}
)

var (
	// Global RBAC service
	gRBAC *service
)

const (
	watchInterval = time.Hour

	RuleResourceType = "corteza::generic:rbac-rule"
)

// Global returns global RBAC service
func Global() *service {
	return gRBAC
}

// SetGlobal re-sets global service
func SetGlobal(svc *service) {
	gRBAC = svc
}

// NewService initializes service{} struct
//
// service{} struct preloads, checks, grants and flushes privileges to and from store
// It acts as a caching layer
func NewService(logger *zap.Logger, s rbacRulesStore) (svc *service) {
	svc = &service{
		l: &sync.RWMutex{},
		f: make(chan bool),

		store: s,
	}

	if logger != nil {
		svc.logger = logger.Named("rbac")
	}

	return
}

// Can function performs permission check for roles in context
//
// First extracts roles from context, then
// use Check() to test against permission rules and
// iterate over all fallback functions
//
// System user is always allowed to do everything
func (svc *service) Can(ses Session, op string, res Resource) bool {
	return svc.Check(ses, op, res) == Allow
}

// Check verifies if role has access to perform an operation on a resource
//
// See RuleSet's Check() func for details
func (svc *service) Check(ses Session, op string, res Resource) (a Access) {
	var (
		fRoles = getContextRoles(ses, res, svc.roles)
	)

	if hasWildcards(res.RbacResource()) {
		// prevent use of wildcard resources for checking permissions
		return Inherit
	}

	a = check(svc.indexed, fRoles, op, res.RbacResource(), nil)

	svc.logger.Debug(a.String()+" "+op+" for "+res.RbacResource(),
		append(
			fRoles.LogFields(),
			zap.Uint64("identity", ses.Identity()),
			zap.Any("indexed", len(svc.indexed)),
			zap.Any("rules", len(svc.rules)),
		)...,
	)

	return
}

// Trace checks RBAC rules and returns all decision trace log
func (svc *service) Trace(ses Session, op string, res Resource) *Trace {
	var (
		t = new(Trace)
	)

	if hasWildcards(res.RbacResource()) {
		// a special case for when user has contextual roles
		// AND trace is done on a resource with wildcards
		ctxRolesDebug := partRoles{ContextRole: make(map[uint64]bool)}
		for _, memberOf := range ses.Roles() {
			for _, role := range svc.roles {
				if role.kind != ContextRole {
					continue
				}

				if role.id != memberOf {
					continue
				}

				// member of contextual role
				//
				// this is a tricky situation:
				// when doing regular check this is an unlikely scenario since
				// check can not be done on a resource with wildcards
				//
				// all contextual roles we're members off will be collected
				ctxRolesDebug[ContextRole][memberOf] = true

			}
		}

		if len(ctxRolesDebug[ContextRole]) > 0 {
			// session has at least one contextual role
			// and since we're checking this on a wildcard resource
			// there is no need to procede with RBAC check
			baseTraceInfo(t, res.RbacResource(), op, ctxRolesDebug)
			resolve(t, Inherit, unknownContext)
			return t
		}
	}

	var (
		fRoles = getContextRoles(ses, res, svc.roles)
	)

	_ = check(svc.indexed, fRoles, op, res.RbacResource(), t)

	return t
}

// Grant appends and/or overwrites internal rules slice
//
// All rules with Inherit are removed
func (svc *service) Grant(ctx context.Context, rules ...*Rule) (err error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	for _, r := range rules {
		svc.logger.Debug(r.Access.String() + " " + r.Operation + " on " + r.Resource + " to " + strconv.FormatUint(r.RoleID, 10))
	}

	svc.grant(rules...)
	return svc.flush(ctx)
}

func (svc *service) grant(rules ...*Rule) {
	svc.rules = merge(svc.rules, rules...)
	svc.indexed = indexRules(svc.rules)
}

// Watch reloads RBAC rules in intervals and on request
func (svc *service) Watch(ctx context.Context) {
	go func() {
		defer sentry.Recover()

		var ticker = time.NewTicker(watchInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.Reload(ctx)
			case <-svc.f:
				svc.Reload(ctx)
			}
		}
	}()

	svc.logger.Debug("watcher initialized")
}

// FindRulesByRoleID returns all RBAC rules that belong to a role
func (svc *service) FindRulesByRoleID(roleID uint64) (rr RuleSet) {
	svc.l.RLock()
	defer svc.l.RUnlock()

	return ruleByRole(svc.rules, roleID)
}

// Rules return all roles
func (svc *service) Rules() (rr RuleSet) {
	svc.l.RLock()
	defer svc.l.RUnlock()
	return svc.rules
}

// Reload store rules
func (svc *service) Reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()
	svc.reloadRules(ctx)
}

// Clear removes all access control rules
func (svc *service) Clear() {
	svc.l.Lock()
	defer svc.l.Unlock()
	svc.rules = RuleSet{}
	svc.indexed = OptRuleSet{}
}

func (svc *service) reloadRules(ctx context.Context) {
	rr, _, err := svc.store.SearchRbacRules(ctx, RuleFilter{})
	svc.logger.Debug(
		"reloading rules",
		zap.Error(err),
		zap.Int("before", len(svc.rules)),
		zap.Int("after", len(rr)),
	)

	if err == nil {
		svc.rules = rr
		svc.indexed = indexRules(rr)
	}
}

// UpdateRoles updates RBAC roles
//
// Warning: this REPLACES all existing roles that are recognized by RBAC subsystem
func (svc *service) UpdateRoles(rr ...*Role) {
	svc.l.Lock()
	defer svc.l.Unlock()

	stats := statRoles(rr...)
	svc.logger.Debug(
		"updating roles",
		zap.Int("before", len(svc.roles)),
		zap.Int("after", len(rr)),
		zap.Int("bypass", stats[BypassRole]),
		zap.Int("context", stats[ContextRole]),
		zap.Int("common", stats[CommonRole]),
		zap.Int("authenticated", stats[AuthenticatedRole]),
		zap.Int("anonymous", stats[AnonymousRole]),
	)
	svc.roles = rr
}

// flush pushes all changed rules to the store (if service is configured with one)
func (svc *service) flush(ctx context.Context) (err error) {
	if svc.store == nil {
		svc.logger.Debug("rule flushing disabled (no store)")
		return
	}

	deletable, updatable, final := flushable(svc.rules)

	err = svc.store.DeleteRbacRule(ctx, deletable...)
	if err != nil {
		return
	}

	err = svc.store.UpsertRbacRule(ctx, updatable...)
	if err != nil {
		return
	}

	clear(final)
	svc.rules = final
	svc.logger.Debug(
		"flushed rules",
		zap.Int("deleted", len(deletable)),
		zap.Int("updated", len(updatable)),
		zap.Int("final", len(final)),
	)

	return
}

// SignificantRoles returns two list of significant roles.
//
// See sigRoles on rules for more details
func (svc *service) SignificantRoles(res Resource, op string) (aRR, dRR []uint64) {
	svc.l.Lock()
	defer svc.l.Unlock()

	return svc.rules.sigRoles(res.RbacResource(), op)
}

func (svc *service) String() (out string) {
	tpl := "%-5v %-20s to %-20s %-30s\n"
	out += strings.Repeat("-", 120) + "\n"

	role := func(id uint64) string {
		for _, r := range svc.roles {
			if r.id == id {
				if r.handle != "" {
					return fmt.Sprintf("%s [%d]", r.handle, r.kind)
				}
				return fmt.Sprintf("%d [%d]", id, r.kind)
			}
		}

		return fmt.Sprintf("%d [?]", id)
	}

	for _, byRole := range svc.indexed {
		for _, rr := range byRole {
			for _, r := range rr {
				out += fmt.Sprintf(tpl, r.Access, r.Operation, role(r.RoleID), r.Resource)
			}
		}
	}

	out += strings.Repeat("-", 120) + "\n"

	return
}

// CloneRulesByRoleID clone all rules of a Role S to a specific Role T by removing its existing rules
func (svc *service) CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) (err error) {
	var (
		updatedRules RuleSet
	)

	// Make sure rules of fromRoleID stays intact
	rr := svc.FindRulesByRoleID(fromRoleID)

	for _, roleID := range toRoleID {
		// Remove existing rules
		existingRules := svc.FindRulesByRoleID(roleID)
		for _, rule := range existingRules {
			// Make sure to remove existing rules
			rule.Access = Inherit
		}
		updatedRules = append(updatedRules, existingRules...)

		// Clone rules from role S to role T
		for _, rule := range rr {
			// Make sure everything is properly set
			r := *rule
			r.RoleID = roleID
			updatedRules = append(updatedRules, &r)
		}
	}

	return svc.Grant(ctx, updatedRules...)
}
