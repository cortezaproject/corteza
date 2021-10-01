package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - system.apigw-route.yaml
// - system.application.yaml
// - system.auth-client.yaml
// - system.report.yaml
// - system.role.yaml
// - system.template.yaml
// - system.user.yaml
// - system.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
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

// CanSuspendUser checks if current user can suspemd user
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

// CanGrant checks if current user can manage system permissions
//
// This function is auto-generated
func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, "grant", &types.Component{})
}

// CanReadActionLog checks if current user can access to action log
//
// This function is auto-generated
func (svc accessControl) CanReadActionLog(ctx context.Context) bool {
	return svc.can(ctx, "action-log.read", &types.Component{})
}

// CanReadSettings checks if current user can read system settings
//
// This function is auto-generated
func (svc accessControl) CanReadSettings(ctx context.Context) bool {
	return svc.can(ctx, "settings.read", &types.Component{})
}

// CanManageSettings checks if current user can manage system settings
//
// This function is auto-generated
func (svc accessControl) CanManageSettings(ctx context.Context) bool {
	return svc.can(ctx, "settings.manage", &types.Component{})
}

// CanCreateAuthClient checks if current user can create auth clients
//
// This function is auto-generated
func (svc accessControl) CanCreateAuthClient(ctx context.Context) bool {
	return svc.can(ctx, "auth-client.create", &types.Component{})
}

// CanSearchAuthClients checks if current user can list, search or filter auth clients
//
// This function is auto-generated
func (svc accessControl) CanSearchAuthClients(ctx context.Context) bool {
	return svc.can(ctx, "auth-clients.search", &types.Component{})
}

// CanCreateRole checks if current user can create roles
//
// This function is auto-generated
func (svc accessControl) CanCreateRole(ctx context.Context) bool {
	return svc.can(ctx, "role.create", &types.Component{})
}

// CanSearchRoles checks if current user can list, search or filter roles
//
// This function is auto-generated
func (svc accessControl) CanSearchRoles(ctx context.Context) bool {
	return svc.can(ctx, "roles.search", &types.Component{})
}

// CanCreateUser checks if current user can create users
//
// This function is auto-generated
func (svc accessControl) CanCreateUser(ctx context.Context) bool {
	return svc.can(ctx, "user.create", &types.Component{})
}

// CanSearchUsers checks if current user can list, search or filter users
//
// This function is auto-generated
func (svc accessControl) CanSearchUsers(ctx context.Context) bool {
	return svc.can(ctx, "users.search", &types.Component{})
}

// CanCreateApplication checks if current user can create applications
//
// This function is auto-generated
func (svc accessControl) CanCreateApplication(ctx context.Context) bool {
	return svc.can(ctx, "application.create", &types.Component{})
}

// CanSearchApplications checks if current user can list, search or filter auth clients
//
// This function is auto-generated
func (svc accessControl) CanSearchApplications(ctx context.Context) bool {
	return svc.can(ctx, "applications.search", &types.Component{})
}

// CanSelfApplicationFlag checks if current user can manage private flags for applications
//
// This function is auto-generated
func (svc accessControl) CanSelfApplicationFlag(ctx context.Context) bool {
	return svc.can(ctx, "application.flag.self", &types.Component{})
}

// CanGlobalApplicationFlag checks if current user can manage global flags for applications
//
// This function is auto-generated
func (svc accessControl) CanGlobalApplicationFlag(ctx context.Context) bool {
	return svc.can(ctx, "application.flag.global", &types.Component{})
}

// CanCreateTemplate checks if current user can create template
//
// This function is auto-generated
func (svc accessControl) CanCreateTemplate(ctx context.Context) bool {
	return svc.can(ctx, "template.create", &types.Component{})
}

// CanSearchTemplates checks if current user can list, search or filter templates
//
// This function is auto-generated
func (svc accessControl) CanSearchTemplates(ctx context.Context) bool {
	return svc.can(ctx, "templates.search", &types.Component{})
}

// CanCreateReport checks if current user can create report
//
// This function is auto-generated
func (svc accessControl) CanCreateReport(ctx context.Context) bool {
	return svc.can(ctx, "report.create", &types.Component{})
}

// CanSearchReports checks if current user can list, search or filter reports
//
// This function is auto-generated
func (svc accessControl) CanSearchReports(ctx context.Context) bool {
	return svc.can(ctx, "reports.search", &types.Component{})
}

// CanAssignReminder checks if current user can assign reminders
//
// This function is auto-generated
func (svc accessControl) CanAssignReminder(ctx context.Context) bool {
	return svc.can(ctx, "reminder.assign", &types.Component{})
}

// CanCreateQueue checks if current user can create messagebus queues
//
// This function is auto-generated
func (svc accessControl) CanCreateQueue(ctx context.Context) bool {
	return svc.can(ctx, "queue.create", &types.Component{})
}

// CanSearchQueues checks if current user can list, search or filter messagebus queues
//
// This function is auto-generated
func (svc accessControl) CanSearchQueues(ctx context.Context) bool {
	return svc.can(ctx, "queues.search", &types.Component{})
}

// CanCreateApigwRoute checks if current user can create api gateway route
//
// This function is auto-generated
func (svc accessControl) CanCreateApigwRoute(ctx context.Context) bool {
	return svc.can(ctx, "apigw-route.create", &types.Component{})
}

// CanSearchApigwRoutes checks if current user can list search or filter api gateway routes
//
// This function is auto-generated
func (svc accessControl) CanSearchApigwRoutes(ctx context.Context) bool {
	return svc.can(ctx, "apigw-routes.search", &types.Component{})
}

// CanManageResourceTranslations checks if current user can list, search, create, or update resource translations
//
// This function is auto-generated
func (svc accessControl) CanManageResourceTranslations(ctx context.Context) bool {
	return svc.can(ctx, "resource-translations.manage", &types.Component{})
}

// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceType(r) {
	case types.ApigwRouteResourceType:
		return rbacApigwRouteResourceValidator(r, oo...)
	case types.ApplicationResourceType:
		return rbacApplicationResourceValidator(r, oo...)
	case types.AuthClientResourceType:
		return rbacAuthClientResourceValidator(r, oo...)
	case types.ReportResourceType:
		return rbacReportResourceValidator(r, oo...)
	case types.RoleResourceType:
		return rbacRoleResourceValidator(r, oo...)
	case types.TemplateResourceType:
		return rbacTemplateResourceValidator(r, oo...)
	case types.UserResourceType:
		return rbacUserResourceValidator(r, oo...)
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
	case types.ApigwRouteResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.ApplicationResourceType:
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
			"read":         true,
			"update":       true,
			"delete":       true,
			"suspend":      true,
			"unsuspend":    true,
			"email.unmask": true,
			"name.unmask":  true,
			"impersonate":  true,
		}
	case types.ComponentResourceType:
		return map[string]bool{
			"grant":                        true,
			"action-log.read":              true,
			"settings.read":                true,
			"settings.manage":              true,
			"auth-client.create":           true,
			"auth-clients.search":          true,
			"role.create":                  true,
			"roles.search":                 true,
			"user.create":                  true,
			"users.search":                 true,
			"application.create":           true,
			"applications.search":          true,
			"application.flag.self":        true,
			"application.flag.global":      true,
			"template.create":              true,
			"templates.search":             true,
			"report.create":                true,
			"reports.search":               true,
			"reminder.assign":              true,
			"queue.create":                 true,
			"queues.search":                true,
			"apigw-route.create":           true,
			"apigw-routes.search":          true,
			"resource-translations.manage": true,
		}
	}

	return nil
}

// rbacApigwRouteResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacApigwRouteResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system ApigwRoute resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ApigwRouteResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for ApigwRoute", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacApplicationResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacApplicationResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system Application resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ApplicationResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for Application", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacAuthClientResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacAuthClientResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system AuthClient resource", o)
		}
	}

	if !strings.HasPrefix(r, types.AuthClientResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for AuthClient", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacReportResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacReportResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system Report resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ReportResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for Report", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacRoleResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacRoleResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system Role resource", o)
		}
	}

	if !strings.HasPrefix(r, types.RoleResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for Role", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacTemplateResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacTemplateResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system Template resource", o)
		}
	}

	if !strings.HasPrefix(r, types.TemplateResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for Template", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacUserResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacUserResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for system User resource", o)
		}
	}

	if !strings.HasPrefix(r, types.UserResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
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
				return fmt.Errorf("invalid resource path wildcard level (%d) for User", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
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
			return fmt.Errorf("invalid operation '%s' for system resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ComponentResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	return nil
}
