package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/spf13/cast"
	"strings"
)

type (
	accessControl struct {
		actionlog actionlog.Recorder

		roleFinder func(ctx context.Context, id uint64) ([]uint64, error)
		rbac       interface {
			Evaluate(rbac.Session, string, rbac.Resource) rbac.Evaluated
			Grant(context.Context, ...*rbac.Rule) error
			FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
			CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) error
		}
	}
)

func AccessControl(rf func(ctx context.Context, id uint64) ([]uint64, error)) *accessControl {
	return &accessControl{
		roleFinder: rf,
		rbac:       rbac.Global(),
		actionlog:  DefaultActionlog,
	}
}

func (svc accessControl) can(ctx context.Context, op string, res rbac.Resource) bool {
	return svc.rbac.Evaluate(rbac.ContextToSession(ctx), op, res).Can
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

// Evaluate returns a list of permissions evaluated for the given user/roles combo
func (svc accessControl) Evaluate(ctx context.Context, user uint64, roles []uint64, rr ...rbac.Resource) (ee rbac.EvaluatedSet, err error) {
	// Reusing the grant permission since this is who the feature is for
	if !svc.CanGrant(ctx) {
		// @todo should be altered to check grant permissions PER resource
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	// Load roles for this user
	//
	// User's roles take priority over specified ones
	if user != 0 {
		rr, err := svc.roleFinder(ctx, user)
		if err != nil {
			return nil, err
		}

		roles = append(rr, roles...)
	}

	session := rbac.ParamsToSession(ctx, user, roles...)
	for _, res := range rr {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			eval := svc.rbac.Evaluate(session, op, res)

			ee = append(ee, eval)
		}
	}

	return
}

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

	return fmt.Errorf("unknown resource type '%q'", r)
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
			return fmt.Errorf("invalid operation '%s' for automation component resource", o)
		}
	}

	return nil
}
