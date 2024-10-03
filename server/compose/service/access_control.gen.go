package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/types"
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
		Trace(rbac.Session, string, rbac.Resource) *rbac.Trace
		Grant(context.Context, ...*rbac.Rule) error
		FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
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

	session := rbac.ParamsToSession(ctx, userID, roles...)
	for _, res := range resources {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee = append(ee, svc.rbac.Trace(session, op, res))
		}
	}

	return
}

// Resources returns list of resources
//
// This function is auto-generated
func (svc accessControl) Resources() []rbac.Resource {
	return []rbac.Resource{
		rbac.NewResource(types.ChartRbacResource(0, 0)),
		rbac.NewResource(types.ModuleRbacResource(0, 0)),
		rbac.NewResource(types.ModuleFieldRbacResource(0, 0, 0)),
		rbac.NewResource(types.NamespaceRbacResource(0)),
		rbac.NewResource(types.PageRbacResource(0, 0)),
		rbac.NewResource(types.PageLayoutRbacResource(0, 0, 0)),
		rbac.NewResource(types.RecordRbacResource(0, 0, 0)),
		rbac.NewResource(types.ComponentRbacResource()),
	}
}

// List returns list of operations on all resources
//
// This function is auto-generated
func (svc accessControl) List() (out []map[string]string) {
	def := []map[string]string{
		{
			"type": types.ChartResourceType,
			"any":  types.ChartRbacResource(0, 0),
			"op":   "read",
		},
		{
			"type": types.ChartResourceType,
			"any":  types.ChartRbacResource(0, 0),
			"op":   "update",
		},
		{
			"type": types.ChartResourceType,
			"any":  types.ChartRbacResource(0, 0),
			"op":   "delete",
		},
		{
			"type": types.ModuleResourceType,
			"any":  types.ModuleRbacResource(0, 0),
			"op":   "read",
		},
		{
			"type": types.ModuleResourceType,
			"any":  types.ModuleRbacResource(0, 0),
			"op":   "update",
		},
		{
			"type": types.ModuleResourceType,
			"any":  types.ModuleRbacResource(0, 0),
			"op":   "delete",
		},
		{
			"type": types.ModuleResourceType,
			"any":  types.ModuleRbacResource(0, 0),
			"op":   "record.create",
		},
		{
			"type": types.ModuleResourceType,
			"any":  types.ModuleRbacResource(0, 0),
			"op":   "owned-record.create",
		},
		{
			"type": types.ModuleResourceType,
			"any":  types.ModuleRbacResource(0, 0),
			"op":   "records.search",
		},
		{
			"type": types.ModuleFieldResourceType,
			"any":  types.ModuleFieldRbacResource(0, 0, 0),
			"op":   "record.value.read",
		},
		{
			"type": types.ModuleFieldResourceType,
			"any":  types.ModuleFieldRbacResource(0, 0, 0),
			"op":   "record.value.update",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "read",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "update",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "delete",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "export",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "manage",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "module.create",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "modules.search",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "modules.export",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "chart.create",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "charts.search",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "charts.export",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "page.create",
		},
		{
			"type": types.NamespaceResourceType,
			"any":  types.NamespaceRbacResource(0),
			"op":   "pages.search",
		},
		{
			"type": types.PageResourceType,
			"any":  types.PageRbacResource(0, 0),
			"op":   "read",
		},
		{
			"type": types.PageResourceType,
			"any":  types.PageRbacResource(0, 0),
			"op":   "update",
		},
		{
			"type": types.PageResourceType,
			"any":  types.PageRbacResource(0, 0),
			"op":   "delete",
		},
		{
			"type": types.PageResourceType,
			"any":  types.PageRbacResource(0, 0),
			"op":   "page-layout.create",
		},
		{
			"type": types.PageResourceType,
			"any":  types.PageRbacResource(0, 0),
			"op":   "page-layouts.search",
		},
		{
			"type": types.PageLayoutResourceType,
			"any":  types.PageLayoutRbacResource(0, 0, 0),
			"op":   "read",
		},
		{
			"type": types.PageLayoutResourceType,
			"any":  types.PageLayoutRbacResource(0, 0, 0),
			"op":   "update",
		},
		{
			"type": types.PageLayoutResourceType,
			"any":  types.PageLayoutRbacResource(0, 0, 0),
			"op":   "delete",
		},
		{
			"type": types.RecordResourceType,
			"any":  types.RecordRbacResource(0, 0, 0),
			"op":   "read",
		},
		{
			"type": types.RecordResourceType,
			"any":  types.RecordRbacResource(0, 0, 0),
			"op":   "update",
		},
		{
			"type": types.RecordResourceType,
			"any":  types.RecordRbacResource(0, 0, 0),
			"op":   "delete",
		},
		{
			"type": types.RecordResourceType,
			"any":  types.RecordRbacResource(0, 0, 0),
			"op":   "undelete",
		},
		{
			"type": types.RecordResourceType,
			"any":  types.RecordRbacResource(0, 0, 0),
			"op":   "owner.manage",
		},
		{
			"type": types.RecordResourceType,
			"any":  types.RecordRbacResource(0, 0, 0),
			"op":   "revisions.search",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "grant",
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
			"op":   "namespace.create",
		},
		{
			"type": types.ComponentResourceType,
			"any":  types.ComponentRbacResource(),
			"op":   "namespaces.search",
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

	return svc.rbac.FindRulesByRoleID(roleID), nil
}

// CanReadChart checks if current user can read
//
// This function is auto-generated
func (svc accessControl) CanReadChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateChart checks if current user can update
//
// This function is auto-generated
func (svc accessControl) CanUpdateChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteChart checks if current user can delete
//
// This function is auto-generated
func (svc accessControl) CanDeleteChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadModule checks if current user can read
//
// This function is auto-generated
func (svc accessControl) CanReadModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateModule checks if current user can update
//
// This function is auto-generated
func (svc accessControl) CanUpdateModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteModule checks if current user can delete
//
// This function is auto-generated
func (svc accessControl) CanDeleteModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "delete", r)
}

// CanCreateRecordOnModule checks if current user can create record
//
// This function is auto-generated
func (svc accessControl) CanCreateRecordOnModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "record.create", r)
}

// CanCreateOwnedRecordOnModule checks if current user can create record with custom owner
//
// This function is auto-generated
func (svc accessControl) CanCreateOwnedRecordOnModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "owned-record.create", r)
}

// CanSearchRecordsOnModule checks if current user can list, search or filter records
//
// This function is auto-generated
func (svc accessControl) CanSearchRecordsOnModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "records.search", r)
}

// CanReadRecordValueOnModuleField checks if current user can read field value on records
//
// This function is auto-generated
func (svc accessControl) CanReadRecordValueOnModuleField(ctx context.Context, r *types.ModuleField) bool {
	return svc.can(ctx, "record.value.read", r)
}

// CanUpdateRecordValueOnModuleField checks if current user can update field value on records
//
// This function is auto-generated
func (svc accessControl) CanUpdateRecordValueOnModuleField(ctx context.Context, r *types.ModuleField) bool {
	return svc.can(ctx, "record.value.update", r)
}

// CanReadNamespace checks if current user can read
//
// This function is auto-generated
func (svc accessControl) CanReadNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateNamespace checks if current user can update
//
// This function is auto-generated
func (svc accessControl) CanUpdateNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteNamespace checks if current user can delete
//
// This function is auto-generated
func (svc accessControl) CanDeleteNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "delete", r)
}

// CanExportNamespace checks if current user can access to export the entire namespace
//
// This function is auto-generated
func (svc accessControl) CanExportNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "export", r)
}

// CanManageNamespace checks if current user can access to namespace admin panel
//
// This function is auto-generated
func (svc accessControl) CanManageNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "manage", r)
}

// CanCreateModuleOnNamespace checks if current user can create module on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreateModuleOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "module.create", r)
}

// CanSearchModulesOnNamespace checks if current user can list, search or filter module on namespace
//
// This function is auto-generated
func (svc accessControl) CanSearchModulesOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "modules.search", r)
}

// CanExportModulesOnNamespace checks if current user can export modules on namespace
//
// This function is auto-generated
func (svc accessControl) CanExportModulesOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "modules.export", r)
}

// CanCreateChartOnNamespace checks if current user can create chart on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreateChartOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "chart.create", r)
}

// CanSearchChartsOnNamespace checks if current user can list, search or filter chart on namespace
//
// This function is auto-generated
func (svc accessControl) CanSearchChartsOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "charts.search", r)
}

// CanExportChartsOnNamespace checks if current user can export charts on namespace
//
// This function is auto-generated
func (svc accessControl) CanExportChartsOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "charts.export", r)
}

// CanCreatePageOnNamespace checks if current user can create page on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreatePageOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "page.create", r)
}

// CanSearchPagesOnNamespace checks if current user can list, search or filter pages on namespace
//
// This function is auto-generated
func (svc accessControl) CanSearchPagesOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "pages.search", r)
}

// CanReadPage checks if current user can read
//
// This function is auto-generated
func (svc accessControl) CanReadPage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdatePage checks if current user can update
//
// This function is auto-generated
func (svc accessControl) CanUpdatePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "update", r)
}

// CanDeletePage checks if current user can delete
//
// This function is auto-generated
func (svc accessControl) CanDeletePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "delete", r)
}

// CanCreatePageLayoutOnPage checks if current user can create page layout on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreatePageLayoutOnPage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "page-layout.create", r)
}

// CanSearchPageLayoutsOnPage checks if current user can list, search or filter page layouts on namespace
//
// This function is auto-generated
func (svc accessControl) CanSearchPageLayoutsOnPage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "page-layouts.search", r)
}

// CanReadPageLayout checks if current user can read
//
// This function is auto-generated
func (svc accessControl) CanReadPageLayout(ctx context.Context, r *types.PageLayout) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdatePageLayout checks if current user can update
//
// This function is auto-generated
func (svc accessControl) CanUpdatePageLayout(ctx context.Context, r *types.PageLayout) bool {
	return svc.can(ctx, "update", r)
}

// CanDeletePageLayout checks if current user can delete
//
// This function is auto-generated
func (svc accessControl) CanDeletePageLayout(ctx context.Context, r *types.PageLayout) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadRecord checks if current user can read
//
// This function is auto-generated
func (svc accessControl) CanReadRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateRecord checks if current user can update
//
// This function is auto-generated
func (svc accessControl) CanUpdateRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteRecord checks if current user can delete
//
// This function is auto-generated
func (svc accessControl) CanDeleteRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "delete", r)
}

// CanUndeleteRecord checks if current user can undelete
//
// This function is auto-generated
func (svc accessControl) CanUndeleteRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "undelete", r)
}

// CanManageOwnerOnRecord checks if current user can owner.manage
//
// This function is auto-generated
func (svc accessControl) CanManageOwnerOnRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "owner.manage", r)
}

// CanSearchRevisionsOnRecord checks if current user can revisions.search
//
// This function is auto-generated
func (svc accessControl) CanSearchRevisionsOnRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "revisions.search", r)
}

// CanGrant checks if current user can manage compose permissions
//
// This function is auto-generated
func (svc accessControl) CanGrant(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "grant", r)
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

// CanCreateNamespace checks if current user can create namespace
//
// This function is auto-generated
func (svc accessControl) CanCreateNamespace(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "namespace.create", r)
}

// CanSearchNamespaces checks if current user can list, search or filter namespaces
//
// This function is auto-generated
func (svc accessControl) CanSearchNamespaces(ctx context.Context) bool {
	r := &types.Component{}
	return svc.can(ctx, "namespaces.search", r)
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
	case types.ChartResourceType:
		return rbacChartResourceValidator(r, oo...)
	case types.ModuleResourceType:
		return rbacModuleResourceValidator(r, oo...)
	case types.ModuleFieldResourceType:
		return rbacModuleFieldResourceValidator(r, oo...)
	case types.NamespaceResourceType:
		return rbacNamespaceResourceValidator(r, oo...)
	case types.PageResourceType:
		return rbacPageResourceValidator(r, oo...)
	case types.PageLayoutResourceType:
		return rbacPageLayoutResourceValidator(r, oo...)
	case types.RecordResourceType:
		return rbacRecordResourceValidator(r, oo...)
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
	case types.ChartResourceType:
		if hasWildcard {
			return rbac.NewResource(types.ChartRbacResource(ids[0], ids[1])), nil
		}

		return loadChart(ctx, svc.store, ids[0], ids[1])
	case types.ModuleResourceType:
		if hasWildcard {
			return rbac.NewResource(types.ModuleRbacResource(ids[0], ids[1])), nil
		}

		return loadModule(ctx, svc.store, ids[0], ids[1])
	case types.ModuleFieldResourceType:
		if hasWildcard {
			return rbac.NewResource(types.ModuleFieldRbacResource(ids[0], ids[1], ids[2])), nil
		}

		return loadModuleField(ctx, svc.store, ids[0], ids[1], ids[2])
	case types.NamespaceResourceType:
		if hasWildcard {
			return rbac.NewResource(types.NamespaceRbacResource(ids[0])), nil
		}

		return loadNamespace(ctx, svc.store, ids[0])
	case types.PageResourceType:
		if hasWildcard {
			return rbac.NewResource(types.PageRbacResource(ids[0], ids[1])), nil
		}

		return loadPage(ctx, svc.store, ids[0], ids[1])
	case types.PageLayoutResourceType:
		if hasWildcard {
			return rbac.NewResource(types.PageLayoutRbacResource(ids[0], ids[1], ids[2])), nil
		}

		return loadPageLayout(ctx, svc.store, ids[0], ids[1], ids[2])
	case types.RecordResourceType:
		if hasWildcard {
			return rbac.NewResource(types.RecordRbacResource(ids[0], ids[1], ids[2])), nil
		}

		return loadRecord(ctx, svc.store, ids[0], ids[1], ids[2])
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
	case types.ChartResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.ModuleResourceType:
		return map[string]bool{
			"read":                true,
			"update":              true,
			"delete":              true,
			"record.create":       true,
			"owned-record.create": true,
			"records.search":      true,
		}
	case types.ModuleFieldResourceType:
		return map[string]bool{
			"record.value.read":   true,
			"record.value.update": true,
		}
	case types.NamespaceResourceType:
		return map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"export":         true,
			"manage":         true,
			"module.create":  true,
			"modules.search": true,
			"modules.export": true,
			"chart.create":   true,
			"charts.search":  true,
			"charts.export":  true,
			"page.create":    true,
			"pages.search":   true,
		}
	case types.PageResourceType:
		return map[string]bool{
			"read":                true,
			"update":              true,
			"delete":              true,
			"page-layout.create":  true,
			"page-layouts.search": true,
		}
	case types.PageLayoutResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.RecordResourceType:
		return map[string]bool{
			"read":             true,
			"update":           true,
			"delete":           true,
			"undelete":         true,
			"owner.manage":     true,
			"revisions.search": true,
		}
	case types.ComponentResourceType:
		return map[string]bool{
			"grant":                        true,
			"settings.read":                true,
			"settings.manage":              true,
			"namespace.create":             true,
			"namespaces.search":            true,
			"resource-translations.manage": true,
		}
	}

	return nil
}

// rbacChartResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacChartResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ChartResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for chart resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ChartResourceType):], sep), sep)
		prc = []string{
			"NamespaceID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for chart resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacModuleResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacModuleResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ModuleResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for module resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ModuleResourceType):], sep), sep)
		prc = []string{
			"NamespaceID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for module resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacModuleFieldResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacModuleFieldResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.ModuleFieldResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for moduleField resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.ModuleFieldResourceType):], sep), sep)
		prc = []string{
			"NamespaceID",
			"ModuleID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for moduleField resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacNamespaceResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacNamespaceResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.NamespaceResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for namespace resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.NamespaceResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for namespace resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacPageResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacPageResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.PageResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for page resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.PageResourceType):], sep), sep)
		prc = []string{
			"NamespaceID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for page resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacPageLayoutResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacPageLayoutResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.PageLayoutResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for pageLayout resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.PageLayoutResourceType):], sep), sep)
		prc = []string{
			"NamespaceID",
			"PageID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for pageLayout resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacRecordResourceValidator checks validity of RBAC resource and operations
//
// # Notes
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacRecordResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.RecordResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for record resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.RecordResourceType):], sep), sep)
		prc = []string{
			"NamespaceID",
			"ModuleID",
			"ID",
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i := 0; i < len(pp); i++ {
		if pp[i] != "*" {
			if i > 0 && pp[i-1] == "*" {
				return fmt.Errorf("invalid path wildcard level (%d) for record resource", i)
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
			return fmt.Errorf("invalid operation '%s' for compose component resource", o)
		}
	}

	return nil
}
