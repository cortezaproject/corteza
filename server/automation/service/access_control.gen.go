package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
	"strings"
)

type (
	rbacService interface {
		Can(rbac.Session, string, rbac.Resource) bool
		Trace(rbac.Session, string, rbac.Resource) (*rbac.Trace, error)
		Grant(context.Context, ...*rbac.Rule) error
		FindRulesByRoleID(ctx context.Context, roleID uint64) (rr rbac.RuleSet, err error)
	}

	accessControl struct {
		actionlog actionlog.Recorder

		store store.Storer
		rbac  rbacService
	}
)

func AccessControl(s store.Storer) *accessControl {
	return &accessControl{
		store:     s,
		rbac:      rbac.Global(),
		actionlog: DefaultActionlog,
	}
}

func (svc accessControl) can(ctx context.Context, op string, res rbac.Resource) bool {
	return svc.rbac.Can(rbac.ContextToSession(ctx), op, res)
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
func (svc accessControl) Trace(ctx context.Context, userID uint64, roles []uint64, rr ...string) (ee []*rbac.Trace, err error) {
	// Reusing the grant permission since this is who the feature is for
	if !svc.CanGrant(ctx) {
		// @todo should be altered to check grant permissions PER resource
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	var (
		resource  rbac.Resource
		resources []rbac.Resource
		members   systemTypes.RoleMemberSet
	)
	if len(rr) > 0 {
		resources = make([]rbac.Resource, 0, len(rr))
		for _, r := range rr {
			if err = rbacResourceValidator(r); err != nil {
				return nil, fmt.Errorf("can not use resource %q: %w", r, err)
			}

			resource, err = svc.resourceLoader(ctx, r)
			if err != nil {
				return
			}

			resources = append(resources, resource)
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

		members, _, err = store.SearchRoleMembers(ctx, svc.store, systemTypes.RoleMemberFilter{UserID: userID})
		if err != nil {
			return nil, err
		}

		for _, m := range members {
			roles = append(roles, m.RoleID)
		}

		for _, r := range internalAuth.AuthenticatedRoles() {
			roles = append(roles, r.ID)
		}
	}

	if len(roles) == 0 {
		// should be prevented on the client
		return nil, fmt.Errorf("no roles specified")
	}

	var auxTrace *rbac.Trace
	session := rbac.ParamsToSession(ctx, userID, roles...)
	for _, res := range resources {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			auxTrace, err = svc.rbac.Trace(session, op, res)
			if err != nil {
				return
			}
			ee = append(ee, auxTrace)
		}
	}

	return
}

// Resources returns list of resources
//
// This function is auto-generated
func (svc accessControl) Resources() []rbac.Resource {
	return []rbac.Resource{
		rbac.NewResource(types.WorkflowRbacResource(0)),
		rbac.NewResource(types.ComponentRbacResource()),
	}
}

// List returns list of operations on all resources
//
// This function is auto-generated
func (svc accessControl) List() (out []map[string]string) {
	def := []map[string]string{
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "undelete",
		},
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "execute",
		},
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "triggers.manage",
		},
		{
			"type": types.WorkflowResourceType,
			"any":  types.WorkflowRbacResource(0),
			"op":   "sessions.manage",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "grant",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "workflow.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "triggers.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "sessions.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "workflows.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "resource-translations.manage",
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

// FindRules find all rules based on filters
//
// This function is auto-generated
func (svc accessControl) FindRules(ctx context.Context, roleID uint64, rr ...string) (out rbac.RuleSet, err error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	out, err = svc.FindRulesByRoleID(ctx, roleID)
	if err != nil {
		return
	}

	var resources []rbac.Resource
	if len(rr) > 0 {
		resources = make([]rbac.Resource, 0, len(rr))
		for _, r := range rr {
			if err = rbacResourceValidator(r); err != nil {
				return nil, fmt.Errorf("can not use resource %q: %w", r, err)
			}

			resources = append(resources, rbac.NewResource(r))
		}
	} else {
		resources = svc.Resources()
	}

	return out.FilterResource(resources...), nil
}

// FindRulesByRoleID find all rules for a specific role
//
// This function is auto-generated
func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.rbac.FindRulesByRoleID(ctx, roleID)
}

// CanReadWorkflow checks if current user can read workflow
//
// This function is auto-generated
func (svc accessControl) CanReadWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateWorkflow checks if current user can update workflow
//
// This function is auto-generated
func (svc accessControl) CanUpdateWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteWorkflow checks if current user can delete workflow
//
// This function is auto-generated
func (svc accessControl) CanDeleteWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "delete", r)
}

// CanUndeleteWorkflow checks if current user can undelete workflow
//
// This function is auto-generated
func (svc accessControl) CanUndeleteWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "undelete", r)
}

// CanExecuteWorkflow checks if current user can execute workflow
//
// This function is auto-generated
func (svc accessControl) CanExecuteWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "execute", r)
}

// CanManageTriggersOnWorkflow checks if current user can manage workflow triggers
//
// This function is auto-generated
func (svc accessControl) CanManageTriggersOnWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "triggers.manage", r)
}

// CanManageSessionsOnWorkflow checks if current user can manage workflow sessions
//
// This function is auto-generated
func (svc accessControl) CanManageSessionsOnWorkflow(ctx context.Context, r *types.Workflow) bool {
	return svc.can(ctx, "sessions.manage", r)
}

// CanGrant checks if current user can manage automation permissions
//
// This function is auto-generated
func (svc accessControl) CanGrant(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "grant", r)
}

// CanCreateWorkflow checks if current user can create workflows
//
// This function is auto-generated
func (svc accessControl) CanCreateWorkflow(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "workflow.create", r)
}

// CanSearchTriggers checks if current user can list, search or filter triggers
//
// This function is auto-generated
func (svc accessControl) CanSearchTriggers(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "triggers.search", r)
}

// CanSearchSessions checks if current user can list, search or filter sessions
//
// This function is auto-generated
func (svc accessControl) CanSearchSessions(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "sessions.search", r)
}

// CanSearchWorkflows checks if current user can list, search or filter workflows
//
// This function is auto-generated
func (svc accessControl) CanSearchWorkflows(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "workflows.search", r)
}

// CanManageResourceTranslations checks if current user can list, search, create, or update resource translations
//
// This function is auto-generated
func (svc accessControl) CanManageResourceTranslations(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "resource-translations.manage", r)
}

// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceType(r) {
	case types.WorkflowResourceType:
		return rbacWorkflowResourceValidator(r, oo...)
	case types.ComponentResourceType:
		return rbacComponentResourceValidator(r, oo...)
	}

	return fmt.Errorf("unknown resource type %q", r)
}

// resourceLoader loads resource from store
//
// # Notes
// Function assumes existence of loader functions for all resource types
//
// This function is auto-generated
func (svc accessControl) resourceLoader(ctx context.Context, resource string) (rbac.Resource, error) {
	var (
		hasWildcard       = false
		resourceType, ids = rbac.ParseResourceID(resource)
	)

	for _, id := range ids {
		if id == 0 {
			hasWildcard = true
			break
		}
	}

	switch rbac.ResourceType(resourceType) {
	case types.WorkflowResourceType:
		if hasWildcard {
			return rbac.NewResource(types.WorkflowRbacResource(ids[0])), nil
		}

		return loadWorkflow(ctx, svc.store, ids[0])
	case types.ComponentResourceType:
		return &types.Component{}, nil
	}

	_ = ids
	return nil, fmt.Errorf("unknown resource type %q", resourceType)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceType(r) {
	case types.WorkflowResourceType:
		return map[string]bool{
			"read":            true,
			"update":          true,
			"delete":          true,
			"undelete":        true,
			"execute":         true,
			"triggers.manage": true,
			"sessions.manage": true,
		}
	case types.ComponentResourceType:
		return map[string]bool{
			"grant":                        true,
			"workflow.create":              true,
			"triggers.search":              true,
			"sessions.search":              true,
			"workflows.search":             true,
			"resource-translations.manage": true,
		}
	}

	return nil
}

// rbacWorkflowResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacWorkflowResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.WorkflowResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for workflow resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.WorkflowResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for workflow resource", i)
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
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacComponentResourceValidator(r string, oo ...string) error {
	if r != types.ComponentResourceType+"/" {
		// expecting resource to always include path
		return fmt.Errorf("invalid component resource, expecting " + types.ComponentResourceType + "/")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for automation component resource", o)
		}
	}

	return nil
}
