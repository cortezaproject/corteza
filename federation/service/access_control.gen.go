package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
	"strings"
)

type (
	roleMemberSearcher interface {
		SearchRoleMembers(context.Context, systemTypes.RoleMemberFilter) (systemTypes.RoleMemberSet, systemTypes.RoleMemberFilter, error)
	}

	rbacService interface {
		Evaluate(rbac.Session, string, rbac.Resource) rbac.Evaluated
		Grant(context.Context, ...*rbac.Rule) error
		FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
		CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) error
	}

	accessControl struct {
		actionlog actionlog.Recorder

		store roleMemberSearcher
		rbac  rbacService
	}
)

func AccessControl(rms roleMemberSearcher) *accessControl {
	return &accessControl{
		store:     rms,
		rbac:      rbac.Global(),
		actionlog: DefaultActionlog,
	}
}

func (svc accessControl) can(ctx context.Context, op string, res rbac.Resource) bool {
	return svc.rbac.Evaluate(rbac.ContextToSession(ctx), op, res).Can
}

// Effective returns a list of effective permissions for all given resource
//
// This function is auto-generated
func (svc accessControl) Effective(ctx context.Context, rr ...rbac.Resource) (ee rbac.EffectiveSet) {
	for _, res := range rr {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee.Push(r, op, svc.can(ctx, op, res))
		}
	}

	return
}

// Evaluate returns a list of permissions evaluated for the given user/roles combo
//
// This function is auto-generated
func (svc accessControl) Evaluate(ctx context.Context, userID uint64, roles []uint64, rr ...string) (ee rbac.EvaluatedSet, err error) {
	// Reusing the grant permission since this is who the feature is for
	if !svc.CanGrant(ctx) {
		// @todo should be altered to check grant permissions PER resource
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	var (
		resources []rbac.Resource
		members   systemTypes.RoleMemberSet
	)
	if len(rr) > 0 {
		resources = make([]rbac.Resource, 0, len(rr))
		for _, r := range rr {
			resources = append(resources, rbac.NewResource(r))
		}
	} else {
		resources = svc.Resources()
	}

	// User ID specified, load its roles
	if userID != 0 {
		if len(roles) > 0 {
			// should be prevented on the client
			return nil, fmt.Errorf("userID and roles are mutually exclusive")
		}

		members, _, err = svc.store.SearchRoleMembers(ctx, systemTypes.RoleMemberFilter{UserID: userID})
		if err != nil {
			return nil, err
		}

		for _, m := range members {
			roles = append(roles, m.RoleID)
		}
	}

	if len(roles) == 0 {
		// should be prevented on the client
		return nil, fmt.Errorf("no roles specified")
	}

	session := rbac.ParamsToSession(ctx, userID, roles...)
	for _, res := range resources {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee = append(ee, svc.rbac.Evaluate(session, op, res))
		}
	}

	return
}

// Resources returns list of resources
//
// This function is auto-generated
func (svc accessControl) Resources() []rbac.Resource {
	return []rbac.Resource{
		rbac.NewResource(types.NodeRbacResource(0)),
		rbac.NewResource(types.ExposedModuleRbacResource(0, 0)),
		rbac.NewResource(types.SharedModuleRbacResource(0, 0)),
		rbac.NewResource(types.ComponentRbacResource()),
	}
}

// List returns list of operations on all resources
//
// This function is auto-generated
func (svc accessControl) List() (out []map[string]string) {
	def := []map[string]string{
		{
			"type": types.NodeResourceType,
			"any":  types.NodeRbacResource(0),
			"op":   "manage",
		},
		{
			"type": types.NodeResourceType,
			"any":  types.NodeRbacResource(0),
			"op":   "module.create",
		},
		{
			"type": types.ExposedModuleResourceType,
			"any":  types.ExposedModuleRbacResource(0, 0),
			"op":   "manage",
		},
		{
			"type": types.SharedModuleResourceType,
			"any":  types.SharedModuleRbacResource(0, 0),
			"op":   "map",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "grant",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "pair",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "settings.read",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "settings.manage",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "node.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "nodes.search",
		},
	}

	func(svc interface{}) {
		if svc, is := svc.(interface{}).(interface{ list() []map[string]string }); is {
			def = append(def, svc.list()...)
		}
	}(svc)

	return def
}

// Grant applies one or more RBAC rules
//
// This function is auto-generated
func (svc accessControl) Grant(ctx context.Context, rr ...*rbac.Rule) error {
	if !svc.CanGrant(ctx) {
		// @todo should be altered to check grant permissions PER resource
		return AccessControlErrNotAllowedToSetPermissions()
	}

	for _, r := range rr {
		err := rbacResourceValidator(r.Resource, r.Operation)
		if err != nil {
			return err
		}
	}

	if err := svc.rbac.Grant(ctx, rr...); err != nil {
		return AccessControlErrGeneric().Wrap(err)
	}

	svc.logGrants(ctx, rr)

	return nil
}

// This function is auto-generated
func (svc accessControl) logGrants(ctx context.Context, rr []*rbac.Rule) {
	if svc.actionlog == nil {
		return
	}

	for _, r := range rr {
		g := AccessControlActionGrant(&accessControlActionProps{r})
		g.log = r.String()
		g.resource = r.Resource

		svc.actionlog.Record(ctx, g.ToAction())
	}
}

// FindRulesByRoleID find all rules for a specific role
//
// This function is auto-generated
func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.rbac.FindRulesByRoleID(roleID), nil
}

// CloneRulesByRoleID clone all rules of a Role S to a specific Role T
//
// This function is auto-generated
func (svc accessControl) CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) error {
	if !svc.CanGrant(ctx) {
		return AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.rbac.CloneRulesByRoleID(ctx, fromRoleID, toRoleID...)
}

// CanManageNode checks if current user can manage federation node
//
// This function is auto-generated
func (svc accessControl) CanManageNode(ctx context.Context, r *types.Node) bool {
	return svc.can(ctx, "manage", r)
}

// CanCreateModuleOnNode checks if current user can create shared module
//
// This function is auto-generated
func (svc accessControl) CanCreateModuleOnNode(ctx context.Context, r *types.Node) bool {
	return svc.can(ctx, "module.create", r)
}

// CanManageExposedModule checks if current user can manage exposed module module
//
// This function is auto-generated
func (svc accessControl) CanManageExposedModule(ctx context.Context, r *types.ExposedModule) bool {
	return svc.can(ctx, "manage", r)
}

// CanMapSharedModule checks if current user can map shared module
//
// This function is auto-generated
func (svc accessControl) CanMapSharedModule(ctx context.Context, r *types.SharedModule) bool {
	return svc.can(ctx, "map", r)
}

// CanGrant checks if current user can manage federation permissions
//
// This function is auto-generated
func (svc accessControl) CanGrant(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "grant", r)
}

// CanPair checks if current user can pair federation nodes
//
// This function is auto-generated
func (svc accessControl) CanPair(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "pair", r)
}

// CanReadSettings checks if current user can read settings
//
// This function is auto-generated
func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "settings.read", r)
}

// CanManageSettings checks if current user can manage settings
//
// This function is auto-generated
func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "settings.manage", r)
}

// CanCreateNode checks if current user can create new federation node
//
// This function is auto-generated
func (svc accessControl) CanCreateNode(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "node.create", r)
}

// CanSearchNodes checks if current user can list, search or filter federation nodes
//
// This function is auto-generated
func (svc accessControl) CanSearchNodes(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "nodes.search", r)
}

// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceType(r) {
	case types.NodeResourceType:
		return rbacNodeResourceValidator(r, oo...)
	case types.ExposedModuleResourceType:
		return rbacExposedModuleResourceValidator(r, oo...)
	case types.SharedModuleResourceType:
		return rbacSharedModuleResourceValidator(r, oo...)
	case types.ComponentResourceType:
		return rbacComponentResourceValidator(r, oo...)
	}

	return fmt.Errorf("unknown resource type '%q'", r)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceType(r) {
	case types.NodeResourceType:
		return map[string]bool{
			"manage":        true,
			"module.create": true,
		}
	case types.ExposedModuleResourceType:
		return map[string]bool{
			"manage": true,
		}
	case types.SharedModuleResourceType:
		return map[string]bool{
			"map": true,
		}
	case types.ComponentResourceType:
		return map[string]bool{
			"grant":           true,
			"pair":            true,
			"settings.read":   true,
			"settings.manage": true,
			"node.create":     true,
			"nodes.search":    true,
		}
	}

	return nil
}

// rbacNodeResourceValidator checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacNodeResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.NodeResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for node resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.NodeResourceType):], sep), sep)
		prc = []string{
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for node resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacExposedModuleResourceValidator checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacExposedModuleResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ExposedModuleResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for exposedModule resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ExposedModuleResourceType):], sep), sep)
		prc = []string{
			"NodeID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for exposedModule resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacSharedModuleResourceValidator checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacSharedModuleResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.SharedModuleResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for sharedModule resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.SharedModuleResourceType):], sep), sep)
		prc = []string{
			"NodeID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for sharedModule resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacComponentResourceValidator checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacComponentResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ComponentResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for federation component resource", o)
		}
	}

	return nil
}
