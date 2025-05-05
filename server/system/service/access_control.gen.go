package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
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
		rbac.NewResource(types.ApplicationRbacResource(0)),
		rbac.NewResource(types.ApigwRouteRbacResource(0)),
		rbac.NewResource(types.AuthClientRbacResource(0)),
		rbac.NewResource(types.DataPrivacyRequestRbacResource(0)),
		rbac.NewResource(types.QueueRbacResource(0)),
		rbac.NewResource(types.ReportRbacResource(0)),
		rbac.NewResource(types.RoleRbacResource(0)),
		rbac.NewResource(types.TemplateRbacResource(0)),
		rbac.NewResource(types.UserRbacResource(0)),
		rbac.NewResource(types.DalConnectionRbacResource(0)),
		rbac.NewResource(types.ComponentRbacResource()),
	}
}

// List returns list of operations on all resources
//
// This function is auto-generated
func (svc accessControl) List() (out []map[string]string) {
	def := []map[string]string{
		{
			"type": types.ApplicationResourceType,
			"any":  types.ApplicationRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.ApplicationResourceType,
			"any":  types.ApplicationRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.ApplicationResourceType,
			"any":  types.ApplicationRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.ApigwRouteResourceType,
			"any":  types.ApigwRouteRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.ApigwRouteResourceType,
			"any":  types.ApigwRouteRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.ApigwRouteResourceType,
			"any":  types.ApigwRouteRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.AuthClientResourceType,
			"any":  types.AuthClientRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.AuthClientResourceType,
			"any":  types.AuthClientRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.AuthClientResourceType,
			"any":  types.AuthClientRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.AuthClientResourceType,
			"any":  types.AuthClientRbacResource(0),
			"op":   "authorize",
		},
		{
			"type": types.DataPrivacyRequestResourceType,
			"any":  types.DataPrivacyRequestRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.DataPrivacyRequestResourceType,
			"any":  types.DataPrivacyRequestRbacResource(0),
			"op":   "approve",
		},
		{
			"type": types.QueueResourceType,
			"any":  types.QueueRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.QueueResourceType,
			"any":  types.QueueRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.QueueResourceType,
			"any":  types.QueueRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.QueueResourceType,
			"any":  types.QueueRbacResource(0),
			"op":   "queue.read",
		},
		{
			"type": types.QueueResourceType,
			"any":  types.QueueRbacResource(0),
			"op":   "queue.write",
		},
		{
			"type": types.ReportResourceType,
			"any":  types.ReportRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.ReportResourceType,
			"any":  types.ReportRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.ReportResourceType,
			"any":  types.ReportRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.ReportResourceType,
			"any":  types.ReportRbacResource(0),
			"op":   "run",
		},
		{
			"type": types.RoleResourceType,
			"any":  types.RoleRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.RoleResourceType,
			"any":  types.RoleRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.RoleResourceType,
			"any":  types.RoleRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.RoleResourceType,
			"any":  types.RoleRbacResource(0),
			"op":   "members.manage",
		},
		{
			"type": types.TemplateResourceType,
			"any":  types.TemplateRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.TemplateResourceType,
			"any":  types.TemplateRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.TemplateResourceType,
			"any":  types.TemplateRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.TemplateResourceType,
			"any":  types.TemplateRbacResource(0),
			"op":   "render",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "suspend",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "unsuspend",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "email.unmask",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "name.unmask",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "impersonate",
		},
		{
			"type": types.UserResourceType,
			"any":  types.UserRbacResource(0),
			"op":   "credentials.manage",
		},
		{
			"type": types.DalConnectionResourceType,
			"any":  types.DalConnectionRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.DalConnectionResourceType,
			"any":  types.DalConnectionRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.DalConnectionResourceType,
			"any":  types.DalConnectionRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.DalConnectionResourceType,
			"any":  types.DalConnectionRbacResource(0),
			"op":   "dal-config.manage",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "grant",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "action-log.read",
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
			"op":   "auth-client.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "auth-clients.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "role.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "roles.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "user.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "users.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "dal-connection.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "dal-connections.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "dal-sensitivity-level.manage",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "application.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "applications.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "application.flag.self",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "application.flag.global",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "template.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "templates.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "report.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "reports.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "reminder.assign",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "queue.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "queues.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "apigw-route.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "apigw-routes.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "resource-translations.manage",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "dal-schema-alterations.manage",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "data-privacy-request.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "data-privacy-requests.search",
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

// CanReadApplication checks if current user can read application
//
// This function is auto-generated
func (svc accessControl) CanReadApplication(ctx context.Context, r *types.Application) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateApplication checks if current user can update application
//
// This function is auto-generated
func (svc accessControl) CanUpdateApplication(ctx context.Context, r *types.Application) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteApplication checks if current user can delete application
//
// This function is auto-generated
func (svc accessControl) CanDeleteApplication(ctx context.Context, r *types.Application) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadApigwRoute checks if current user can read api gateway route
//
// This function is auto-generated
func (svc accessControl) CanReadApigwRoute(ctx context.Context, r *types.ApigwRoute) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateApigwRoute checks if current user can update api gateway route
//
// This function is auto-generated
func (svc accessControl) CanUpdateApigwRoute(ctx context.Context, r *types.ApigwRoute) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteApigwRoute checks if current user can delete api gateway route
//
// This function is auto-generated
func (svc accessControl) CanDeleteApigwRoute(ctx context.Context, r *types.ApigwRoute) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadAuthClient checks if current user can read authorization client
//
// This function is auto-generated
func (svc accessControl) CanReadAuthClient(ctx context.Context, r *types.AuthClient) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateAuthClient checks if current user can update authorization client
//
// This function is auto-generated
func (svc accessControl) CanUpdateAuthClient(ctx context.Context, r *types.AuthClient) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteAuthClient checks if current user can delete authorization client
//
// This function is auto-generated
func (svc accessControl) CanDeleteAuthClient(ctx context.Context, r *types.AuthClient) bool {
	return svc.can(ctx, "delete", r)
}

// CanAuthorizeAuthClient checks if current user can authorize authorization client
//
// This function is auto-generated
func (svc accessControl) CanAuthorizeAuthClient(ctx context.Context, r *types.AuthClient) bool {
	return svc.can(ctx, "authorize", r)
}

// CanReadDataPrivacyRequest checks if current user can read data privacy request
//
// This function is auto-generated
func (svc accessControl) CanReadDataPrivacyRequest(ctx context.Context, r *types.DataPrivacyRequest) bool {
	return svc.can(ctx, "read", r)
}

// CanApproveDataPrivacyRequest checks if current user can approve/reject data privacy request
//
// This function is auto-generated
func (svc accessControl) CanApproveDataPrivacyRequest(ctx context.Context, r *types.DataPrivacyRequest) bool {
	return svc.can(ctx, "approve", r)
}

// CanReadQueue checks if current user can read queue
//
// This function is auto-generated
func (svc accessControl) CanReadQueue(ctx context.Context, r *types.Queue) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateQueue checks if current user can update queue
//
// This function is auto-generated
func (svc accessControl) CanUpdateQueue(ctx context.Context, r *types.Queue) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteQueue checks if current user can delete queue
//
// This function is auto-generated
func (svc accessControl) CanDeleteQueue(ctx context.Context, r *types.Queue) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadQueueOnQueue checks if current user can read from queue
//
// This function is auto-generated
func (svc accessControl) CanReadQueueOnQueue(ctx context.Context, r *types.Queue) bool {
	return svc.can(ctx, "queue.read", r)
}

// CanWriteQueueOnQueue checks if current user can write to queue
//
// This function is auto-generated
func (svc accessControl) CanWriteQueueOnQueue(ctx context.Context, r *types.Queue) bool {
	return svc.can(ctx, "queue.write", r)
}

// CanReadReport checks if current user can read report
//
// This function is auto-generated
func (svc accessControl) CanReadReport(ctx context.Context, r *types.Report) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateReport checks if current user can update report
//
// This function is auto-generated
func (svc accessControl) CanUpdateReport(ctx context.Context, r *types.Report) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteReport checks if current user can delete report
//
// This function is auto-generated
func (svc accessControl) CanDeleteReport(ctx context.Context, r *types.Report) bool {
	return svc.can(ctx, "delete", r)
}

// CanRunReport checks if current user can run report
//
// This function is auto-generated
func (svc accessControl) CanRunReport(ctx context.Context, r *types.Report) bool {
	return svc.can(ctx, "run", r)
}

// CanReadRole checks if current user can read role
//
// This function is auto-generated
func (svc accessControl) CanReadRole(ctx context.Context, r *types.Role) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateRole checks if current user can update role
//
// This function is auto-generated
func (svc accessControl) CanUpdateRole(ctx context.Context, r *types.Role) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteRole checks if current user can delete role
//
// This function is auto-generated
func (svc accessControl) CanDeleteRole(ctx context.Context, r *types.Role) bool {
	return svc.can(ctx, "delete", r)
}

// CanManageMembersOnRole checks if current user can manage members
//
// This function is auto-generated
func (svc accessControl) CanManageMembersOnRole(ctx context.Context, r *types.Role) bool {
	return svc.can(ctx, "members.manage", r)
}

// CanReadTemplate checks if current user can read template
//
// This function is auto-generated
func (svc accessControl) CanReadTemplate(ctx context.Context, r *types.Template) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateTemplate checks if current user can update template
//
// This function is auto-generated
func (svc accessControl) CanUpdateTemplate(ctx context.Context, r *types.Template) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteTemplate checks if current user can delete template
//
// This function is auto-generated
func (svc accessControl) CanDeleteTemplate(ctx context.Context, r *types.Template) bool {
	return svc.can(ctx, "delete", r)
}

// CanRenderTemplate checks if current user can render template
//
// This function is auto-generated
func (svc accessControl) CanRenderTemplate(ctx context.Context, r *types.Template) bool {
	return svc.can(ctx, "render", r)
}

// CanReadUser checks if current user can read user
//
// This function is auto-generated
func (svc accessControl) CanReadUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateUser checks if current user can update user
//
// This function is auto-generated
func (svc accessControl) CanUpdateUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteUser checks if current user can delete user
//
// This function is auto-generated
func (svc accessControl) CanDeleteUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "delete", r)
}

// CanSuspendUser checks if current user can suspend user
//
// This function is auto-generated
func (svc accessControl) CanSuspendUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "suspend", r)
}

// CanUnsuspendUser checks if current user can unsuspend user
//
// This function is auto-generated
func (svc accessControl) CanUnsuspendUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "unsuspend", r)
}

// CanUnmaskEmailOnUser checks if current user can unmask email
//
// This function is auto-generated
func (svc accessControl) CanUnmaskEmailOnUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "email.unmask", r)
}

// CanUnmaskNameOnUser checks if current user can unmask name
//
// This function is auto-generated
func (svc accessControl) CanUnmaskNameOnUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "name.unmask", r)
}

// CanImpersonateUser checks if current user can impersonate user
//
// This function is auto-generated
func (svc accessControl) CanImpersonateUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "impersonate", r)
}

// CanManageCredentialsOnUser checks if current user can manage user's credentials
//
// This function is auto-generated
func (svc accessControl) CanManageCredentialsOnUser(ctx context.Context, r *types.User) bool {
	return svc.can(ctx, "credentials.manage", r)
}

// CanReadDalConnection checks if current user can read connection
//
// This function is auto-generated
func (svc accessControl) CanReadDalConnection(ctx context.Context, r *types.DalConnection) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateDalConnection checks if current user can update connection
//
// This function is auto-generated
func (svc accessControl) CanUpdateDalConnection(ctx context.Context, r *types.DalConnection) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteDalConnection checks if current user can delete connection
//
// This function is auto-generated
func (svc accessControl) CanDeleteDalConnection(ctx context.Context, r *types.DalConnection) bool {
	return svc.can(ctx, "delete", r)
}

// CanManageDalConfigOnDalConnection checks if current user can manage dal configuration
//
// This function is auto-generated
func (svc accessControl) CanManageDalConfigOnDalConnection(ctx context.Context, r *types.DalConnection) bool {
	return svc.can(ctx, "dal-config.manage", r)
}

// CanGrant checks if current user can manage system permissions
//
// This function is auto-generated
func (svc accessControl) CanGrant(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "grant", r)
}

// CanReadActionLog checks if current user can access to action log
//
// This function is auto-generated
func (svc accessControl) CanReadActionLog(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "action-log.read", r)
}

// CanReadSettings checks if current user can read system settings
//
// This function is auto-generated
func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "settings.read", r)
}

// CanManageSettings checks if current user can manage system settings
//
// This function is auto-generated
func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "settings.manage", r)
}

// CanCreateAuthClient checks if current user can create auth clients
//
// This function is auto-generated
func (svc accessControl) CanCreateAuthClient(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "auth-client.create", r)
}

// CanSearchAuthClients checks if current user can list, search or filter auth clients
//
// This function is auto-generated
func (svc accessControl) CanSearchAuthClients(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "auth-clients.search", r)
}

// CanCreateRole checks if current user can create roles
//
// This function is auto-generated
func (svc accessControl) CanCreateRole(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "role.create", r)
}

// CanSearchRoles checks if current user can list, search or filter roles
//
// This function is auto-generated
func (svc accessControl) CanSearchRoles(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "roles.search", r)
}

// CanCreateUser checks if current user can create users
//
// This function is auto-generated
func (svc accessControl) CanCreateUser(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "user.create", r)
}

// CanSearchUsers checks if current user can list, search or filter users
//
// This function is auto-generated
func (svc accessControl) CanSearchUsers(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "users.search", r)
}

// CanCreateDalConnection checks if current user can create dal connections
//
// This function is auto-generated
func (svc accessControl) CanCreateDalConnection(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "dal-connection.create", r)
}

// CanSearchDalConnections checks if current user can list, search or filter dal connections
//
// This function is auto-generated
func (svc accessControl) CanSearchDalConnections(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "dal-connections.search", r)
}

// CanManageDalSensitivityLevel checks if current user can can manage dal sensitivity levels
//
// This function is auto-generated
func (svc accessControl) CanManageDalSensitivityLevel(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "dal-sensitivity-level.manage", r)
}

// CanCreateApplication checks if current user can create applications
//
// This function is auto-generated
func (svc accessControl) CanCreateApplication(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "application.create", r)
}

// CanSearchApplications checks if current user can list, search or filter auth clients
//
// This function is auto-generated
func (svc accessControl) CanSearchApplications(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "applications.search", r)
}

// CanSelfApplicationFlag checks if current user can manage private flags for applications
//
// This function is auto-generated
func (svc accessControl) CanSelfApplicationFlag(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "application.flag.self", r)
}

// CanGlobalApplicationFlag checks if current user can manage global flags for applications
//
// This function is auto-generated
func (svc accessControl) CanGlobalApplicationFlag(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "application.flag.global", r)
}

// CanCreateTemplate checks if current user can create template
//
// This function is auto-generated
func (svc accessControl) CanCreateTemplate(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "template.create", r)
}

// CanSearchTemplates checks if current user can list, search or filter templates
//
// This function is auto-generated
func (svc accessControl) CanSearchTemplates(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "templates.search", r)
}

// CanCreateReport checks if current user can create report
//
// This function is auto-generated
func (svc accessControl) CanCreateReport(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "report.create", r)
}

// CanSearchReports checks if current user can list, search or filter reports
//
// This function is auto-generated
func (svc accessControl) CanSearchReports(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "reports.search", r)
}

// CanAssignReminder checks if current user can  assign reminders
//
// This function is auto-generated
func (svc accessControl) CanAssignReminder(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "reminder.assign", r)
}

// CanCreateQueue checks if current user can create messagebus queues
//
// This function is auto-generated
func (svc accessControl) CanCreateQueue(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "queue.create", r)
}

// CanSearchQueues checks if current user can list, search or filter messagebus queues
//
// This function is auto-generated
func (svc accessControl) CanSearchQueues(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "queues.search", r)
}

// CanCreateApigwRoute checks if current user can create api gateway route
//
// This function is auto-generated
func (svc accessControl) CanCreateApigwRoute(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "apigw-route.create", r)
}

// CanSearchApigwRoutes checks if current user can list search or filter api gateway routes
//
// This function is auto-generated
func (svc accessControl) CanSearchApigwRoutes(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "apigw-routes.search", r)
}

// CanManageResourceTranslations checks if current user can list, search, create, or update resource translations
//
// This function is auto-generated
func (svc accessControl) CanManageResourceTranslations(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "resource-translations.manage", r)
}

// CanManageDalSchemaAlterations checks if current user can list, search, apply, or dismiss dal alterations
//
// This function is auto-generated
func (svc accessControl) CanManageDalSchemaAlterations(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "dal-schema-alterations.manage", r)
}

// CanCreateDataPrivacyRequest checks if current user can create data privacy requests
//
// This function is auto-generated
func (svc accessControl) CanCreateDataPrivacyRequest(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "data-privacy-request.create", r)
}

// CanSearchDataPrivacyRequests checks if current user can list, search or filter data privacy requests
//
// This function is auto-generated
func (svc accessControl) CanSearchDataPrivacyRequests(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "data-privacy-requests.search", r)
}

// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceType(r) {
	case types.ApplicationResourceType:
		return rbacApplicationResourceValidator(r, oo...)
	case types.ApigwRouteResourceType:
		return rbacApigwRouteResourceValidator(r, oo...)
	case types.AuthClientResourceType:
		return rbacAuthClientResourceValidator(r, oo...)
	case types.DataPrivacyRequestResourceType:
		return rbacDataPrivacyRequestResourceValidator(r, oo...)
	case types.QueueResourceType:
		return rbacQueueResourceValidator(r, oo...)
	case types.ReportResourceType:
		return rbacReportResourceValidator(r, oo...)
	case types.RoleResourceType:
		return rbacRoleResourceValidator(r, oo...)
	case types.TemplateResourceType:
		return rbacTemplateResourceValidator(r, oo...)
	case types.UserResourceType:
		return rbacUserResourceValidator(r, oo...)
	case types.DalConnectionResourceType:
		return rbacDalConnectionResourceValidator(r, oo...)
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
	case types.ApplicationResourceType:
		if hasWildcard {
			return rbac.NewResource(types.ApplicationRbacResource(ids[0])), nil
		}

		return loadApplication(ctx, svc.store, ids[0])
	case types.ApigwRouteResourceType:
		if hasWildcard {
			return rbac.NewResource(types.ApigwRouteRbacResource(ids[0])), nil
		}

		return loadApigwRoute(ctx, svc.store, ids[0])
	case types.AuthClientResourceType:
		if hasWildcard {
			return rbac.NewResource(types.AuthClientRbacResource(ids[0])), nil
		}

		return loadAuthClient(ctx, svc.store, ids[0])
	case types.DataPrivacyRequestResourceType:
		if hasWildcard {
			return rbac.NewResource(types.DataPrivacyRequestRbacResource(ids[0])), nil
		}

		return loadDataPrivacyRequest(ctx, svc.store, ids[0])
	case types.QueueResourceType:
		if hasWildcard {
			return rbac.NewResource(types.QueueRbacResource(ids[0])), nil
		}

		return loadQueue(ctx, svc.store, ids[0])
	case types.ReportResourceType:
		if hasWildcard {
			return rbac.NewResource(types.ReportRbacResource(ids[0])), nil
		}

		return loadReport(ctx, svc.store, ids[0])
	case types.RoleResourceType:
		if hasWildcard {
			return rbac.NewResource(types.RoleRbacResource(ids[0])), nil
		}

		return loadRole(ctx, svc.store, ids[0])
	case types.TemplateResourceType:
		if hasWildcard {
			return rbac.NewResource(types.TemplateRbacResource(ids[0])), nil
		}

		return loadTemplate(ctx, svc.store, ids[0])
	case types.UserResourceType:
		if hasWildcard {
			return rbac.NewResource(types.UserRbacResource(ids[0])), nil
		}

		return loadUser(ctx, svc.store, ids[0])
	case types.DalConnectionResourceType:
		if hasWildcard {
			return rbac.NewResource(types.DalConnectionRbacResource(ids[0])), nil
		}

		return loadDalConnection(ctx, svc.store, ids[0])
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
	case types.ApplicationResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.ApigwRouteResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.AuthClientResourceType:
		return map[string]bool{
			"read":      true,
			"update":    true,
			"delete":    true,
			"authorize": true,
		}
	case types.DataPrivacyRequestResourceType:
		return map[string]bool{
			"read":    true,
			"approve": true,
		}
	case types.QueueResourceType:
		return map[string]bool{
			"read":        true,
			"update":      true,
			"delete":      true,
			"queue.read":  true,
			"queue.write": true,
		}
	case types.ReportResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
			"run":    true,
		}
	case types.RoleResourceType:
		return map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"members.manage": true,
		}
	case types.TemplateResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
			"render": true,
		}
	case types.UserResourceType:
		return map[string]bool{
			"read":               true,
			"update":             true,
			"delete":             true,
			"suspend":            true,
			"unsuspend":          true,
			"email.unmask":       true,
			"name.unmask":        true,
			"impersonate":        true,
			"credentials.manage": true,
		}
	case types.DalConnectionResourceType:
		return map[string]bool{
			"read":              true,
			"update":            true,
			"delete":            true,
			"dal-config.manage": true,
		}
	case types.ComponentResourceType:
		return map[string]bool{
			"grant":                         true,
			"action-log.read":               true,
			"settings.read":                 true,
			"settings.manage":               true,
			"auth-client.create":            true,
			"auth-clients.search":           true,
			"role.create":                   true,
			"roles.search":                  true,
			"user.create":                   true,
			"users.search":                  true,
			"dal-connection.create":         true,
			"dal-connections.search":        true,
			"dal-sensitivity-level.manage":  true,
			"application.create":            true,
			"applications.search":           true,
			"application.flag.self":         true,
			"application.flag.global":       true,
			"template.create":               true,
			"templates.search":              true,
			"report.create":                 true,
			"reports.search":                true,
			"reminder.assign":               true,
			"queue.create":                  true,
			"queues.search":                 true,
			"apigw-route.create":            true,
			"apigw-routes.search":           true,
			"resource-translations.manage":  true,
			"dal-schema-alterations.manage": true,
			"data-privacy-request.create":   true,
			"data-privacy-requests.search":  true,
		}
	}

	return nil
}

// rbacApplicationResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacApplicationResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ApplicationResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for application resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ApplicationResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for application resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacApigwRouteResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacApigwRouteResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ApigwRouteResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for apigwRoute resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ApigwRouteResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for apigwRoute resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacAuthClientResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacAuthClientResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.AuthClientResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for authClient resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.AuthClientResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for authClient resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacDataPrivacyRequestResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacDataPrivacyRequestResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.DataPrivacyRequestResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for dataPrivacyRequest resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.DataPrivacyRequestResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for dataPrivacyRequest resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacQueueResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacQueueResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.QueueResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for queue resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.QueueResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for queue resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacReportResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacReportResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ReportResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for report resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ReportResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for report resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacRoleResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacRoleResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.RoleResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for role resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.RoleResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for role resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacTemplateResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacTemplateResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.TemplateResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for template resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.TemplateResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for template resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacUserResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacUserResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.UserResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for user resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.UserResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for user resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacDalConnectionResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacDalConnectionResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.DalConnectionResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for dalConnection resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.DalConnectionResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for dalConnection resource", i)
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
			return fmt.Errorf("invalid operation '%s' for system component resource", o)
		}
	}

	return nil
}
