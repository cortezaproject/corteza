package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - federation.exposed-module.yaml
// - federation.node.yaml
// - federation.shared-module.yaml
// - federation.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/spf13/cast"
	"strings"
)

type (
	accessControl struct {
		actionlog actionlog.Recorder

		rbac interface {
			Can([]uint64, string, rbac.Resource) bool
			Grant(context.Context, ...*rbac.Rule) error
			FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
		}
	}
)

func AccessControl() *accessControl {
	return &accessControl{
		rbac:      rbac.Global(),
		actionlog: DefaultActionlog,
	}
}

func (svc accessControl) can(ctx context.Context, op string, res rbac.Resource) bool {
	var (
		identity = internalAuth.GetIdentityFromContext(ctx)
	)

	if identity == nil {
		panic("expecting identity in context")
	}

	return svc.rbac.Can(identity.Roles(), op, res)
}

// Effective returns a list of effective permissions for all given resource
func (svc accessControl) Effective(ctx context.Context, rr ...rbac.Resource) (ee rbac.EffectiveSet) {
	for _, res := range rr {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee.Push(r, op, svc.can(ctx, op, res))
		}
	}

	return
}

func (svc accessControl) List() (out []map[string]string) {
	return []map[string]string{
		{"resource": "corteza+federation.exposed-module", "operation": "manage"},
		{"resource": "corteza+federation.node", "operation": "manage"},
		{"resource": "corteza+federation.node", "operation": "module.create"},
		{"resource": "corteza+federation.shared-module", "operation": "map"},
		{"resource": "corteza+federation", "operation": "grant"},
		{"resource": "corteza+federation", "operation": "pair"},
		{"resource": "corteza+federation", "operation": "node.create"},
		{"resource": "corteza+federation", "operation": "settings.read"},
		{"resource": "corteza+federation", "operation": "settings.manage"},
	}
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

// CanManageExposedModule checks if current user can manage shared module
//
// This function is auto-generated
func (svc accessControl) CanManageExposedModule(ctx context.Context, r *types.ExposedModule) bool {
	return svc.can(ctx, "manage", r)
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
	return svc.can(ctx, "grant", &types.Component{})
}

// CanPair checks if current user can pair federation nodes
//
// This function is auto-generated
func (svc accessControl) CanPair(ctx context.Context) bool {
	return svc.can(ctx, "pair", &types.Component{})
}

// CanCreateNode checks if current user can create new federation node
//
// This function is auto-generated
func (svc accessControl) CanCreateNode(ctx context.Context) bool {
	return svc.can(ctx, "node.create", &types.Component{})
}

// CanReadSettings checks if current user can read settings
//
// This function is auto-generated
func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, "settings.read", &types.Component{})
}

// CanManageSettings checks if current user can manage settings
//
// This function is auto-generated
func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, "settings.manage", &types.Component{})
}

// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceSchema(r) {
	case "corteza+federation.exposed-module":
		return rbacExposedModuleResourceValidator(r, oo...)
	case "corteza+federation.node":
		return rbacNodeResourceValidator(r, oo...)
	case "corteza+federation.shared-module":
		return rbacSharedModuleResourceValidator(r, oo...)
	case "corteza+federation":
		return rbacComponentResourceValidator(r, oo...)
	}

	return fmt.Errorf("unknown resource schema '%q'", r)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceSchema(r) {
	case "corteza+federation.exposed-module":
		return map[string]bool{
			"manage": true,
		}
	case "corteza+federation.node":
		return map[string]bool{
			"manage":        true,
			"module.create": true,
		}
	case "corteza+federation.shared-module":
		return map[string]bool{
			"map": true,
		}
	case "corteza+federation":
		return map[string]bool{
			"grant":           true,
			"pair":            true,
			"node.create":     true,
			"settings.read":   true,
			"settings.manage": true,
		}
	}

	return nil
}

// rbacExposedModuleResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacExposedModuleResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for federation ExposedModule resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ExposedModuleRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.ExposedModuleRbacResourceSchema)+2:], "/")
	if len(pp) != 2 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"NodeID",
			"ID",
		}
	)

	for i, p := range pp {
		if p == "*" {
			ppWildcard = true
			continue
		}

		if !ppWildcard {
			return fmt.Errorf("invalid resource path wildcard level")
		}

		if _, err := cast.ToUint64E(p); err != nil {
			return fmt.Errorf("invalid ID for %s: '%s'", pathElements[i], p)
		}
	}

	return nil
}

// rbacNodeResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacNodeResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for federation Node resource", o)
		}
	}

	if !strings.HasPrefix(r, types.NodeRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.NodeRbacResourceSchema)+2:], "/")
	if len(pp) != 1 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"ID",
		}
	)

	for i, p := range pp {
		if p == "*" {
			ppWildcard = true
			continue
		}

		if !ppWildcard {
			return fmt.Errorf("invalid resource path wildcard level")
		}

		if _, err := cast.ToUint64E(p); err != nil {
			return fmt.Errorf("invalid ID for %s: '%s'", pathElements[i], p)
		}
	}

	return nil
}

// rbacSharedModuleResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacSharedModuleResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for federation SharedModule resource", o)
		}
	}

	if !strings.HasPrefix(r, types.SharedModuleRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.SharedModuleRbacResourceSchema)+2:], "/")
	if len(pp) != 2 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"NodeID",
			"ID",
		}
	)

	for i, p := range pp {
		if p == "*" {
			ppWildcard = true
			continue
		}

		if !ppWildcard {
			return fmt.Errorf("invalid resource path wildcard level")
		}

		if _, err := cast.ToUint64E(p); err != nil {
			return fmt.Errorf("invalid ID for %s: '%s'", pathElements[i], p)
		}
	}

	return nil
}

// rbacComponentResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacComponentResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for federation resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ComponentRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	return nil
}
