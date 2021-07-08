package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - automation.workflow.yaml
// - automation.yaml

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

		rbac interface {
			Can(rbac.Session, string, rbac.Resource) bool
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
	return svc.rbac.Can(rbac.ContextToSession(ctx), op, res)
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
	return svc.can(ctx, "grant", &types.Component{})
}

// CanCreateWorkflow checks if current user can create workflows
//
// This function is auto-generated
func (svc accessControl) CanCreateWorkflow(ctx context.Context) bool {
	return svc.can(ctx, "workflow.create", &types.Component{})
}

// CanSearchTriggers checks if current user can search triggers
//
// This function is auto-generated
func (svc accessControl) CanSearchTriggers(ctx context.Context) bool {
	return svc.can(ctx, "triggers.search", &types.Component{})
}

// CanSearchSessions checks if current user can search sessions
//
// This function is auto-generated
func (svc accessControl) CanSearchSessions(ctx context.Context) bool {
	return svc.can(ctx, "sessions.search", &types.Component{})
}

// CanSearchWorkflows checks if current user can search workflows
//
// This function is auto-generated
func (svc accessControl) CanSearchWorkflows(ctx context.Context) bool {
	return svc.can(ctx, "workflows.search", &types.Component{})
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
			"grant":            true,
			"workflow.create":  true,
			"triggers.search":  true,
			"sessions.search":  true,
			"workflows.search": true,
		}
	}

	return nil
}

// rbacWorkflowResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacWorkflowResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for automation Workflow resource", o)
		}
	}

	if !strings.HasPrefix(r, types.WorkflowResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	const sep = "/"
	var (
		specIdUsed = true

		pp  = strings.Split(strings.Trim(r[len(types.WorkflowResourceType):], sep), sep)
		prc = []string{
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i, p := range pp {
		if p == "*" {
			if !specIdUsed {
				return fmt.Errorf("invalid resource path wildcard level (%d) for Workflow", i)
			}

			specIdUsed = false
			continue
		}

		specIdUsed = true
		if _, err := cast.ToUint64E(p); err != nil {
			return fmt.Errorf("invalid reference for %s: '%s'", prc[i], p)
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
			return fmt.Errorf("invalid operation '%s' for automation resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ComponentResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	return nil
}
