package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - compose.chart.yaml
// - compose.module-field.yaml
// - compose.module.yaml
// - compose.namespace.yaml
// - compose.page.yaml
// - compose.record.yaml
// - compose.yaml

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
	return []map[string]string{
		{"resource": "corteza+compose.chart", "operation": "read"},
		{"resource": "corteza+compose.chart", "operation": "update"},
		{"resource": "corteza+compose.chart", "operation": "delete"},
		{"resource": "corteza+compose.module-field", "operation": "record.value.read"},
		{"resource": "corteza+compose.module-field", "operation": "record.value.update"},
		{"resource": "corteza+compose.module", "operation": "read"},
		{"resource": "corteza+compose.module", "operation": "update"},
		{"resource": "corteza+compose.module", "operation": "delete"},
		{"resource": "corteza+compose.module", "operation": "record.create"},
		{"resource": "corteza+compose.namespace", "operation": "read"},
		{"resource": "corteza+compose.namespace", "operation": "update"},
		{"resource": "corteza+compose.namespace", "operation": "delete"},
		{"resource": "corteza+compose.namespace", "operation": "module.create"},
		{"resource": "corteza+compose.namespace", "operation": "chart.create"},
		{"resource": "corteza+compose.namespace", "operation": "page.create"},
		{"resource": "corteza+compose.page", "operation": "read"},
		{"resource": "corteza+compose.page", "operation": "create"},
		{"resource": "corteza+compose.page", "operation": "update"},
		{"resource": "corteza+compose.page", "operation": "delete"},
		{"resource": "corteza+compose.record", "operation": "read"},
		{"resource": "corteza+compose.record", "operation": "update"},
		{"resource": "corteza+compose.record", "operation": "delete"},
		{"resource": "corteza+compose", "operation": "grant"},
		{"resource": "corteza+compose", "operation": "namespace.create"},
		{"resource": "corteza+compose", "operation": "settings.read"},
		{"resource": "corteza+compose", "operation": "settings.manage"},
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

// CanReadChart checks if current user can read chart
//
// This function is auto-generated
func (svc accessControl) CanReadChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateChart checks if current user can update chart
//
// This function is auto-generated
func (svc accessControl) CanUpdateChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteChart checks if current user can delete chart
//
// This function is auto-generated
func (svc accessControl) CanDeleteChart(ctx context.Context, r *types.Chart) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadRecordValue checks if current user can read field value on records
//
// This function is auto-generated
func (svc accessControl) CanReadRecordValue(ctx context.Context, r *types.ModuleField) bool {
	return svc.can(ctx, "record.value.read", r)
}

// CanUpdateRecordValue checks if current user can update field value on records
//
// This function is auto-generated
func (svc accessControl) CanUpdateRecordValue(ctx context.Context, r *types.ModuleField) bool {
	return svc.can(ctx, "record.value.update", r)
}

// CanReadModule checks if current user can read module
//
// This function is auto-generated
func (svc accessControl) CanReadModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateModule checks if current user can update module
//
// This function is auto-generated
func (svc accessControl) CanUpdateModule(ctx context.Context, r *types.Module) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteModule checks if current user can delete module
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

// CanReadNamespace checks if current user can read namespace
//
// This function is auto-generated
func (svc accessControl) CanReadNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateNamespace checks if current user can update namespace
//
// This function is auto-generated
func (svc accessControl) CanUpdateNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteNamespace checks if current user can delete namespace
//
// This function is auto-generated
func (svc accessControl) CanDeleteNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "delete", r)
}

// CanCreateModuleOnNamespace checks if current user can create module on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreateModuleOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "module.create", r)
}

// CanCreateChartOnNamespace checks if current user can create chart on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreateChartOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "chart.create", r)
}

// CanCreatePageOnNamespace checks if current user can create page on namespace
//
// This function is auto-generated
func (svc accessControl) CanCreatePageOnNamespace(ctx context.Context, r *types.Namespace) bool {
	return svc.can(ctx, "page.create", r)
}

// CanReadPage checks if current user can read page
//
// This function is auto-generated
func (svc accessControl) CanReadPage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "read", r)
}

// CanCreatePage checks if current user can create page
//
// This function is auto-generated
func (svc accessControl) CanCreatePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "create", r)
}

// CanUpdatePage checks if current user can update page
//
// This function is auto-generated
func (svc accessControl) CanUpdatePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "update", r)
}

// CanDeletePage checks if current user can delete page
//
// This function is auto-generated
func (svc accessControl) CanDeletePage(ctx context.Context, r *types.Page) bool {
	return svc.can(ctx, "delete", r)
}

// CanReadRecord checks if current user can read record
//
// This function is auto-generated
func (svc accessControl) CanReadRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "read", r)
}

// CanUpdateRecord checks if current user can update record
//
// This function is auto-generated
func (svc accessControl) CanUpdateRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "update", r)
}

// CanDeleteRecord checks if current user can delete record
//
// This function is auto-generated
func (svc accessControl) CanDeleteRecord(ctx context.Context, r *types.Record) bool {
	return svc.can(ctx, "delete", r)
}

// CanGrant checks if current user can manage compose permissions
//
// This function is auto-generated
func (svc accessControl) CanGrant(ctx context.Context) bool {
	return svc.can(ctx, "grant", &types.Component{})
}

// CanCreateNamespace checks if current user can create namespace
//
// This function is auto-generated
func (svc accessControl) CanCreateNamespace(ctx context.Context) bool {
	return svc.can(ctx, "namespace.create", &types.Component{})
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
	case "corteza+compose.chart":
		return rbacChartResourceValidator(r, oo...)
	case "corteza+compose.module-field":
		return rbacModuleFieldResourceValidator(r, oo...)
	case "corteza+compose.module":
		return rbacModuleResourceValidator(r, oo...)
	case "corteza+compose.namespace":
		return rbacNamespaceResourceValidator(r, oo...)
	case "corteza+compose.page":
		return rbacPageResourceValidator(r, oo...)
	case "corteza+compose.record":
		return rbacRecordResourceValidator(r, oo...)
	case "corteza+compose":
		return rbacComponentResourceValidator(r, oo...)
	}

	return fmt.Errorf("unknown resource schema '%q'", r)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceSchema(r) {
	case "corteza+compose.chart":
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case "corteza+compose.module-field":
		return map[string]bool{
			"record.value.read":   true,
			"record.value.update": true,
		}
	case "corteza+compose.module":
		return map[string]bool{
			"read":          true,
			"update":        true,
			"delete":        true,
			"record.create": true,
		}
	case "corteza+compose.namespace":
		return map[string]bool{
			"read":          true,
			"update":        true,
			"delete":        true,
			"module.create": true,
			"chart.create":  true,
			"page.create":   true,
		}
	case "corteza+compose.page":
		return map[string]bool{
			"read":   true,
			"create": true,
			"update": true,
			"delete": true,
		}
	case "corteza+compose.record":
		return map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		}
	case "corteza+compose":
		return map[string]bool{
			"grant":            true,
			"namespace.create": true,
			"settings.read":    true,
			"settings.manage":  true,
		}
	}

	return nil
}

// rbacChartResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacChartResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for compose Chart resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ChartRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.ChartRbacResourceSchema)+2:], "/")
	if len(pp) != 2 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"namespaceID",
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

// rbacModuleFieldResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacModuleFieldResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for compose ModuleField resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ModuleFieldRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.ModuleFieldRbacResourceSchema)+2:], "/")
	if len(pp) != 3 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"namespaceID",
			"moduleID",
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

// rbacModuleResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacModuleResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for compose Module resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ModuleRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.ModuleRbacResourceSchema)+2:], "/")
	if len(pp) != 2 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"namespaceID",
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

// rbacNamespaceResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacNamespaceResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for compose Namespace resource", o)
		}
	}

	if !strings.HasPrefix(r, types.NamespaceRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.NamespaceRbacResourceSchema)+2:], "/")
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

// rbacPageResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacPageResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for compose Page resource", o)
		}
	}

	if !strings.HasPrefix(r, types.PageRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.PageRbacResourceSchema)+2:], "/")
	if len(pp) != 2 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"namespaceID",
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

// rbacRecordResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbacRecordResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for compose Record resource", o)
		}
	}

	if !strings.HasPrefix(r, types.RecordRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	pp := strings.Split(r[len(types.RecordRbacResourceSchema)+2:], "/")
	if len(pp) != 3 {
		return fmt.Errorf("invalid resource path")
	}

	var (
		ppWildcard   bool
		pathElements = []string{
			"namespaceID",
			"moduleID",
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
			return fmt.Errorf("invalid operation '%s' for compose resource", o)
		}
	}

	if !strings.HasPrefix(r, types.ComponentRbacResourceSchema+":/") {
		return fmt.Errorf("invalid schema")
	}

	return nil
}
