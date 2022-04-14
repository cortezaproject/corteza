package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
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
			CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) error
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
	case types.AttachmentResourceType:
		return rbacAttachmentResourceValidator(r, oo...)
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
	case types.RecordResourceType:
		return rbacRecordResourceValidator(r, oo...)
	case types.RecordValueResourceType:
		return rbacRecordValueResourceValidator(r, oo...)
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
	case types.AttachmentResourceType:
		return map[string]bool{}
	case types.ChartResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.ModuleResourceType:
		return map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"record.create":  true,
			"records.search": true,
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
			"manage":         true,
			"module.create":  true,
			"modules.search": true,
			"chart.create":   true,
			"charts.search":  true,
			"page.create":    true,
			"pages.search":   true,
		}
	case types.PageResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.RecordResourceType:
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case types.RecordValueResourceType:
		return map[string]bool{}
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

// rbacAttachmentResourceValidator checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacAttachmentResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.AttachmentResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for attachment resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.AttachmentResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for attachment resource", i)
			}

			if _, err := cast.ToUint64E(pp[i]); err != nil {
				return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
			}
		}
	}
	return nil
}

// rbacChartResourceValidator checks validity of RBAC resource and operations
//
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

// rbacRecordResourceValidator checks validity of RBAC resource and operations
//
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

// rbacRecordValueResourceValidator checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacRecordValueResourceValidator(r string, oo ...string) error {
	if !strings.HasPrefix(r, types.RecordValueResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for recordValue resource", o)
		}
	}

	const sep = "/"
	var (
		pp  = strings.Split(strings.Trim(r[len(types.RecordValueResourceType):], sep), sep)
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
				return fmt.Errorf("invalid path wildcard level (%d) for recordValue resource", i)
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
			return fmt.Errorf("invalid operation '%s' for compose component resource", o)
		}
	}

	return nil
}
