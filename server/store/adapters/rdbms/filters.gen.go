package rdbms

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	automationType "github.com/cortezaproject/corteza/server/automation/types"
	composeType "github.com/cortezaproject/corteza/server/compose/types"
	federationType "github.com/cortezaproject/corteza/server/federation/types"
	actionlogType "github.com/cortezaproject/corteza/server/pkg/actionlog"
	discoveryType "github.com/cortezaproject/corteza/server/pkg/discovery/types"
	flagType "github.com/cortezaproject/corteza/server/pkg/flag/types"
	labelsType "github.com/cortezaproject/corteza/server/pkg/label/types"
	rbacType "github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	systemType "github.com/cortezaproject/corteza/server/system/types"
	"github.com/doug-martin/goqu/v9"
	"strings"
)

type (
	// extendedFilters allows special per-resource
	// filters to be attached to store
	//
	// when optional filter is set, generated filter function is NOT called automatically
	// (but can be called from the optional filter)
	extendedFilters struct {
		// Filter extensions for search/query functions

		// optional actionlog filter function called after the generated function
		Actionlog func(*Store, actionlogType.Filter) ([]goqu.Expression, actionlogType.Filter, error)

		// optional apigwFilter filter function called after the generated function
		ApigwFilter func(*Store, systemType.ApigwFilterFilter) ([]goqu.Expression, systemType.ApigwFilterFilter, error)

		// optional apigwRoute filter function called after the generated function
		ApigwRoute func(*Store, systemType.ApigwRouteFilter) ([]goqu.Expression, systemType.ApigwRouteFilter, error)

		// optional application filter function called after the generated function
		Application func(*Store, systemType.ApplicationFilter) ([]goqu.Expression, systemType.ApplicationFilter, error)

		// optional attachment filter function called after the generated function
		Attachment func(*Store, systemType.AttachmentFilter) ([]goqu.Expression, systemType.AttachmentFilter, error)

		// optional authClient filter function called after the generated function
		AuthClient func(*Store, systemType.AuthClientFilter) ([]goqu.Expression, systemType.AuthClientFilter, error)

		// optional authConfirmedClient filter function called after the generated function
		AuthConfirmedClient func(*Store, systemType.AuthConfirmedClientFilter) ([]goqu.Expression, systemType.AuthConfirmedClientFilter, error)

		// optional authOa2token filter function called after the generated function
		AuthOa2token func(*Store, systemType.AuthOa2tokenFilter) ([]goqu.Expression, systemType.AuthOa2tokenFilter, error)

		// optional authSession filter function called after the generated function
		AuthSession func(*Store, systemType.AuthSessionFilter) ([]goqu.Expression, systemType.AuthSessionFilter, error)

		// optional automationSession filter function called after the generated function
		AutomationSession func(*Store, automationType.SessionFilter) ([]goqu.Expression, automationType.SessionFilter, error)

		// optional automationTrigger filter function called after the generated function
		AutomationTrigger func(*Store, automationType.TriggerFilter) ([]goqu.Expression, automationType.TriggerFilter, error)

		// optional automationWorkflow filter function called after the generated function
		AutomationWorkflow func(*Store, automationType.WorkflowFilter) ([]goqu.Expression, automationType.WorkflowFilter, error)

		// optional composeAttachment filter function called after the generated function
		ComposeAttachment func(*Store, composeType.AttachmentFilter) ([]goqu.Expression, composeType.AttachmentFilter, error)

		// optional composeChart filter function called after the generated function
		ComposeChart func(*Store, composeType.ChartFilter) ([]goqu.Expression, composeType.ChartFilter, error)

		// optional composeModule filter function called after the generated function
		ComposeModule func(*Store, composeType.ModuleFilter) ([]goqu.Expression, composeType.ModuleFilter, error)

		// optional composeModuleField filter function called after the generated function
		ComposeModuleField func(*Store, composeType.ModuleFieldFilter) ([]goqu.Expression, composeType.ModuleFieldFilter, error)

		// optional composeNamespace filter function called after the generated function
		ComposeNamespace func(*Store, composeType.NamespaceFilter) ([]goqu.Expression, composeType.NamespaceFilter, error)

		// optional composePage filter function called after the generated function
		ComposePage func(*Store, composeType.PageFilter) ([]goqu.Expression, composeType.PageFilter, error)

		// optional credential filter function called after the generated function
		Credential func(*Store, systemType.CredentialFilter) ([]goqu.Expression, systemType.CredentialFilter, error)

		// optional dalConnection filter function called after the generated function
		DalConnection func(*Store, systemType.DalConnectionFilter) ([]goqu.Expression, systemType.DalConnectionFilter, error)

		// optional dalSensitivityLevel filter function called after the generated function
		DalSensitivityLevel func(*Store, systemType.DalSensitivityLevelFilter) ([]goqu.Expression, systemType.DalSensitivityLevelFilter, error)

		// optional dataPrivacyRequest filter function called after the generated function
		DataPrivacyRequest func(*Store, systemType.DataPrivacyRequestFilter) ([]goqu.Expression, systemType.DataPrivacyRequestFilter, error)

		// optional dataPrivacyRequestComment filter function called after the generated function
		DataPrivacyRequestComment func(*Store, systemType.DataPrivacyRequestCommentFilter) ([]goqu.Expression, systemType.DataPrivacyRequestCommentFilter, error)

		// optional federationExposedModule filter function called after the generated function
		FederationExposedModule func(*Store, federationType.ExposedModuleFilter) ([]goqu.Expression, federationType.ExposedModuleFilter, error)

		// optional federationModuleMapping filter function called after the generated function
		FederationModuleMapping func(*Store, federationType.ModuleMappingFilter) ([]goqu.Expression, federationType.ModuleMappingFilter, error)

		// optional federationNode filter function called after the generated function
		FederationNode func(*Store, federationType.NodeFilter) ([]goqu.Expression, federationType.NodeFilter, error)

		// optional federationNodeSync filter function called after the generated function
		FederationNodeSync func(*Store, federationType.NodeSyncFilter) ([]goqu.Expression, federationType.NodeSyncFilter, error)

		// optional federationSharedModule filter function called after the generated function
		FederationSharedModule func(*Store, federationType.SharedModuleFilter) ([]goqu.Expression, federationType.SharedModuleFilter, error)

		// optional flag filter function called after the generated function
		Flag func(*Store, flagType.FlagFilter) ([]goqu.Expression, flagType.FlagFilter, error)

		// optional label filter function called after the generated function
		Label func(*Store, labelsType.LabelFilter) ([]goqu.Expression, labelsType.LabelFilter, error)

		// optional queue filter function called after the generated function
		Queue func(*Store, systemType.QueueFilter) ([]goqu.Expression, systemType.QueueFilter, error)

		// optional queueMessage filter function called after the generated function
		QueueMessage func(*Store, systemType.QueueMessageFilter) ([]goqu.Expression, systemType.QueueMessageFilter, error)

		// optional rbacRule filter function called after the generated function
		RbacRule func(*Store, rbacType.RuleFilter) ([]goqu.Expression, rbacType.RuleFilter, error)

		// optional reminder filter function called after the generated function
		Reminder func(*Store, systemType.ReminderFilter) ([]goqu.Expression, systemType.ReminderFilter, error)

		// optional report filter function called after the generated function
		Report func(*Store, systemType.ReportFilter) ([]goqu.Expression, systemType.ReportFilter, error)

		// optional resourceActivity filter function called after the generated function
		ResourceActivity func(*Store, discoveryType.ResourceActivityFilter) ([]goqu.Expression, discoveryType.ResourceActivityFilter, error)

		// optional resourceTranslation filter function called after the generated function
		ResourceTranslation func(*Store, systemType.ResourceTranslationFilter) ([]goqu.Expression, systemType.ResourceTranslationFilter, error)

		// optional role filter function called after the generated function
		Role func(*Store, systemType.RoleFilter) ([]goqu.Expression, systemType.RoleFilter, error)

		// optional roleMember filter function called after the generated function
		RoleMember func(*Store, systemType.RoleMemberFilter) ([]goqu.Expression, systemType.RoleMemberFilter, error)

		// optional settingValue filter function called after the generated function
		SettingValue func(*Store, systemType.SettingsFilter) ([]goqu.Expression, systemType.SettingsFilter, error)

		// optional template filter function called after the generated function
		Template func(*Store, systemType.TemplateFilter) ([]goqu.Expression, systemType.TemplateFilter, error)

		// optional user filter function called after the generated function
		User func(*Store, systemType.UserFilter) ([]goqu.Expression, systemType.UserFilter, error)
	}
)

// ActionlogFilter returns logical expressions
//
// This function is called from Store.QueryActionlogs() and can be extended
// by setting Store.Filters.Actionlog. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ActionlogFilter(d drivers.Dialect, f actionlogType.Filter) (ee []goqu.Expression, _ actionlogType.Filter, err error) {

	if val := strings.TrimSpace(f.Action); len(val) > 0 {
		ee = append(ee, goqu.C("action").Eq(f.Action))
	}

	if val := strings.TrimSpace(f.Resource); len(val) > 0 {
		ee = append(ee, goqu.C("resource").Eq(f.Resource))
	}

	if val := strings.TrimSpace(f.Origin); len(val) > 0 {
		ee = append(ee, goqu.C("origin").Eq(f.Origin))
	}

	if len(f.ActorID) > 0 {
		ee = append(ee, goqu.C("actor_id").In(f.ActorID))
	}

	return ee, f, err
}

// ApigwFilterFilter returns logical expressions
//
// This function is called from Store.QueryApigwFilters() and can be extended
// by setting Store.Filters.ApigwFilter. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ApigwFilterFilter(d drivers.Dialect, f systemType.ApigwFilterFilter) (ee []goqu.Expression, _ systemType.ApigwFilterFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateFalseComparison(d, "enabled", f.Disabled); expr != nil {
		ee = append(ee, expr)
	}

	if f.RouteID > 0 {
		ee = append(ee, goqu.C("rel_route").Eq(f.RouteID))
	}

	return ee, f, err
}

// ApigwRouteFilter returns logical expressions
//
// This function is called from Store.QueryApigwRoutes() and can be extended
// by setting Store.Filters.ApigwRoute. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ApigwRouteFilter(d drivers.Dialect, f systemType.ApigwRouteFilter) (ee []goqu.Expression, _ systemType.ApigwRouteFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateFalseComparison(d, "enabled", f.Disabled); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Route); len(val) > 0 {
		ee = append(ee, goqu.C("id").Eq(f.Route))
	}

	if val := strings.TrimSpace(f.Method); len(val) > 0 {
		ee = append(ee, goqu.C("method").Eq(f.Method))
	}

	return ee, f, err
}

// ApplicationFilter returns logical expressions
//
// This function is called from Store.QueryApplications() and can be extended
// by setting Store.Filters.Application. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ApplicationFilter(d drivers.Dialect, f systemType.ApplicationFilter) (ee []goqu.Expression, _ systemType.ApplicationFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Name); len(val) > 0 {
		ee = append(ee, goqu.C("name").Eq(f.Name))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if len(f.FlaggedIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.FlaggedIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("name").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// AttachmentFilter returns logical expressions
//
// This function is called from Store.QueryAttachments() and can be extended
// by setting Store.Filters.Attachment. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AttachmentFilter(d drivers.Dialect, f systemType.AttachmentFilter) (ee []goqu.Expression, _ systemType.AttachmentFilter, err error) {

	if val := strings.TrimSpace(f.Kind); len(val) > 0 {
		ee = append(ee, goqu.C("kind").Eq(f.Kind))
	}

	return ee, f, err
}

// AuthClientFilter returns logical expressions
//
// This function is called from Store.QueryAuthClients() and can be extended
// by setting Store.Filters.AuthClient. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AuthClientFilter(d drivers.Dialect, f systemType.AuthClientFilter) (ee []goqu.Expression, _ systemType.AuthClientFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	return ee, f, err
}

// AuthConfirmedClientFilter returns logical expressions
//
// This function is called from Store.QueryAuthConfirmedClients() and can be extended
// by setting Store.Filters.AuthConfirmedClient. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AuthConfirmedClientFilter(d drivers.Dialect, f systemType.AuthConfirmedClientFilter) (ee []goqu.Expression, _ systemType.AuthConfirmedClientFilter, err error) {

	if f.UserID > 0 {
		ee = append(ee, goqu.C("rel_user").Eq(f.UserID))
	}

	return ee, f, err
}

// AuthOa2tokenFilter returns logical expressions
//
// This function is called from Store.QueryAuthOa2tokens() and can be extended
// by setting Store.Filters.AuthOa2token. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AuthOa2tokenFilter(d drivers.Dialect, f systemType.AuthOa2tokenFilter) (ee []goqu.Expression, _ systemType.AuthOa2tokenFilter, err error) {

	if f.UserID > 0 {
		ee = append(ee, goqu.C("user_id").Eq(f.UserID))
	}

	return ee, f, err
}

// AuthSessionFilter returns logical expressions
//
// This function is called from Store.QueryAuthSessions() and can be extended
// by setting Store.Filters.AuthSession. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AuthSessionFilter(d drivers.Dialect, f systemType.AuthSessionFilter) (ee []goqu.Expression, _ systemType.AuthSessionFilter, err error) {

	if f.UserID > 0 {
		ee = append(ee, goqu.C("rel_user").Eq(f.UserID))
	}

	return ee, f, err
}

// AutomationSessionFilter returns logical expressions
//
// This function is called from Store.QueryAutomationSessions() and can be extended
// by setting Store.Filters.AutomationSession. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AutomationSessionFilter(d drivers.Dialect, f automationType.SessionFilter) (ee []goqu.Expression, _ automationType.SessionFilter, err error) {

	if expr := stateNilComparison(d, "completed_at", f.Completed); expr != nil {
		ee = append(ee, expr)
	}

	// @todo codegen warning: filtering by Status ([]uint) not supported,
	//       see rdbms.go.tpl and add an exception

	if len(f.WorkflowID) > 0 {
		ee = append(ee, goqu.C("rel_workflow").In(f.WorkflowID))
	}

	if val := strings.TrimSpace(f.EventType); len(val) > 0 {
		ee = append(ee, goqu.C("event_type").Eq(f.EventType))
	}

	if val := strings.TrimSpace(f.ResourceType); len(val) > 0 {
		ee = append(ee, goqu.C("resource_type").Eq(f.ResourceType))
	}

	if len(f.CreatedBy) > 0 {
		ee = append(ee, goqu.C("created_by").In(f.CreatedBy))
	}

	return ee, f, err
}

// AutomationTriggerFilter returns logical expressions
//
// This function is called from Store.QueryAutomationTriggers() and can be extended
// by setting Store.Filters.AutomationTrigger. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AutomationTriggerFilter(d drivers.Dialect, f automationType.TriggerFilter) (ee []goqu.Expression, _ automationType.TriggerFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateFalseComparison(d, "enabled", f.Disabled); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.TriggerID) > 0 {
		ee = append(ee, goqu.C("id").In(f.TriggerID))
	}

	if len(f.WorkflowID) > 0 {
		ee = append(ee, goqu.C("rel_workflow").In(f.WorkflowID))
	}

	if val := strings.TrimSpace(f.EventType); len(val) > 0 {
		ee = append(ee, goqu.C("event_type").Eq(f.EventType))
	}

	if val := strings.TrimSpace(f.ResourceType); len(val) > 0 {
		ee = append(ee, goqu.C("resource_type").Eq(f.ResourceType))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	return ee, f, err
}

// AutomationWorkflowFilter returns logical expressions
//
// This function is called from Store.QueryAutomationWorkflows() and can be extended
// by setting Store.Filters.AutomationWorkflow. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AutomationWorkflowFilter(d drivers.Dialect, f automationType.WorkflowFilter) (ee []goqu.Expression, _ automationType.WorkflowFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateFalseComparison(d, "enabled", f.Disabled); expr != nil {
		ee = append(ee, expr)
	}

	if ss := trimStringSlice(f.WorkflowID); len(ss) > 0 {
		ee = append(ee, goqu.C("id").In(ss))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ComposeAttachmentFilter returns logical expressions
//
// This function is called from Store.QueryComposeAttachments() and can be extended
// by setting Store.Filters.ComposeAttachment. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposeAttachmentFilter(d drivers.Dialect, f composeType.AttachmentFilter) (ee []goqu.Expression, _ composeType.AttachmentFilter, err error) {

	if val := strings.TrimSpace(f.Kind); len(val) > 0 {
		ee = append(ee, goqu.C("kind").Eq(f.Kind))
	}

	if f.NamespaceID > 0 {
		ee = append(ee, goqu.C("namespace_id").Eq(f.NamespaceID))
	}

	return ee, f, err
}

// ComposeChartFilter returns logical expressions
//
// This function is called from Store.QueryComposeCharts() and can be extended
// by setting Store.Filters.ComposeChart. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposeChartFilter(d drivers.Dialect, f composeType.ChartFilter) (ee []goqu.Expression, _ composeType.ChartFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if len(f.ChartID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ChartID))
	}

	if f.NamespaceID > 0 {
		ee = append(ee, goqu.C("rel_namespace").Eq(f.NamespaceID))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("name").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ComposeModuleFilter returns logical expressions
//
// This function is called from Store.QueryComposeModules() and can be extended
// by setting Store.Filters.ComposeModule. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposeModuleFilter(d drivers.Dialect, f composeType.ModuleFilter) (ee []goqu.Expression, _ composeType.ModuleFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if len(f.ModuleID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ModuleID))
	}

	if f.NamespaceID > 0 {
		ee = append(ee, goqu.C("rel_namespace").Eq(f.NamespaceID))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("name").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ComposeModuleFieldFilter returns logical expressions
//
// This function is called from Store.QueryComposeModuleFields() and can be extended
// by setting Store.Filters.ComposeModuleField. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposeModuleFieldFilter(d drivers.Dialect, f composeType.ModuleFieldFilter) (ee []goqu.Expression, _ composeType.ModuleFieldFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ModuleID) > 0 {
		ee = append(ee, goqu.C("rel_module").In(f.ModuleID))
	}

	return ee, f, err
}

// ComposeNamespaceFilter returns logical expressions
//
// This function is called from Store.QueryComposeNamespaces() and can be extended
// by setting Store.Filters.ComposeNamespace. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposeNamespaceFilter(d drivers.Dialect, f composeType.NamespaceFilter) (ee []goqu.Expression, _ composeType.NamespaceFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.NamespaceID) > 0 {
		ee = append(ee, goqu.C("namespace_id").In(f.NamespaceID))
	}

	if val := strings.TrimSpace(f.Name); len(val) > 0 {
		ee = append(ee, goqu.C("name").Eq(f.Name))
	}

	if val := strings.TrimSpace(f.Slug); len(val) > 0 {
		ee = append(ee, goqu.C("slug").Eq(f.Slug))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("name").ILike("%"+f.Query+"%"),
			goqu.C("slug").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ComposePageFilter returns logical expressions
//
// This function is called from Store.QueryComposePages() and can be extended
// by setting Store.Filters.ComposePage. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposePageFilter(d drivers.Dialect, f composeType.PageFilter) (ee []goqu.Expression, _ composeType.PageFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if f.NamespaceID > 0 {
		ee = append(ee, goqu.C("rel_namespace").Eq(f.NamespaceID))
	}

	if f.ModuleID > 0 {
		ee = append(ee, goqu.C("rel_module").Eq(f.ModuleID))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("title").ILike("%"+f.Query+"%"),
			goqu.C("description").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// CredentialFilter returns logical expressions
//
// This function is called from Store.QueryCredentials() and can be extended
// by setting Store.Filters.Credential. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func CredentialFilter(d drivers.Dialect, f systemType.CredentialFilter) (ee []goqu.Expression, _ systemType.CredentialFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if f.OwnerID > 0 {
		ee = append(ee, goqu.C("rel_owner").Eq(f.OwnerID))
	}

	if val := strings.TrimSpace(f.Kind); len(val) > 0 {
		ee = append(ee, goqu.C("kind").Eq(f.Kind))
	}

	if val := strings.TrimSpace(f.Credentials); len(val) > 0 {
		ee = append(ee, goqu.C("credentials").Eq(f.Credentials))
	}

	return ee, f, err
}

// DalConnectionFilter returns logical expressions
//
// This function is called from Store.QueryDalConnections() and can be extended
// by setting Store.Filters.DalConnection. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DalConnectionFilter(d drivers.Dialect, f systemType.DalConnectionFilter) (ee []goqu.Expression, _ systemType.DalConnectionFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ConnectionID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ConnectionID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if val := strings.TrimSpace(f.Type); len(val) > 0 {
		ee = append(ee, goqu.C("type").Eq(f.Type))
	}

	return ee, f, err
}

// DalSensitivityLevelFilter returns logical expressions
//
// This function is called from Store.QueryDalSensitivityLevels() and can be extended
// by setting Store.Filters.DalSensitivityLevel. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DalSensitivityLevelFilter(d drivers.Dialect, f systemType.DalSensitivityLevelFilter) (ee []goqu.Expression, _ systemType.DalSensitivityLevelFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.SensitivityLevelID) > 0 {
		ee = append(ee, goqu.C("id").In(f.SensitivityLevelID))
	}

	return ee, f, err
}

// DataPrivacyRequestFilter returns logical expressions
//
// This function is called from Store.QueryDataPrivacyRequests() and can be extended
// by setting Store.Filters.DataPrivacyRequest. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DataPrivacyRequestFilter(d drivers.Dialect, f systemType.DataPrivacyRequestFilter) (ee []goqu.Expression, _ systemType.DataPrivacyRequestFilter, err error) {

	// @todo codegen warning: filtering by Kind ([]types.RequestKind) not supported,
	//       see rdbms.go.tpl and add an exception

	// @todo codegen warning: filtering by Status ([]types.RequestStatus) not supported,
	//       see rdbms.go.tpl and add an exception

	if len(f.RequestedBy) > 0 {
		ee = append(ee, goqu.C("requested_by").In(f.RequestedBy))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("kind").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// DataPrivacyRequestCommentFilter returns logical expressions
//
// This function is called from Store.QueryDataPrivacyRequestComments() and can be extended
// by setting Store.Filters.DataPrivacyRequestComment. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DataPrivacyRequestCommentFilter(d drivers.Dialect, f systemType.DataPrivacyRequestCommentFilter) (ee []goqu.Expression, _ systemType.DataPrivacyRequestCommentFilter, err error) {

	if len(f.RequestID) > 0 {
		ee = append(ee, goqu.C("rel_request").In(f.RequestID))
	}

	return ee, f, err
}

// FederationExposedModuleFilter returns logical expressions
//
// This function is called from Store.QueryFederationExposedModules() and can be extended
// by setting Store.Filters.FederationExposedModule. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func FederationExposedModuleFilter(d drivers.Dialect, f federationType.ExposedModuleFilter) (ee []goqu.Expression, _ federationType.ExposedModuleFilter, err error) {

	if f.ComposeModuleID > 0 {
		ee = append(ee, goqu.C("rel_compose_module").Eq(f.ComposeModuleID))
	}

	if f.ComposeNamespaceID > 0 {
		ee = append(ee, goqu.C("rel_compose_namespace").Eq(f.ComposeNamespaceID))
	}

	if f.NodeID > 0 {
		ee = append(ee, goqu.C("rel_node").Eq(f.NodeID))
	}

	return ee, f, err
}

// FederationModuleMappingFilter returns logical expressions
//
// This function is called from Store.QueryFederationModuleMappings() and can be extended
// by setting Store.Filters.FederationModuleMapping. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func FederationModuleMappingFilter(d drivers.Dialect, f federationType.ModuleMappingFilter) (ee []goqu.Expression, _ federationType.ModuleMappingFilter, err error) {

	if f.ComposeModuleID > 0 {
		ee = append(ee, goqu.C("rel_compose_module").Eq(f.ComposeModuleID))
	}

	if f.ComposeNamespaceID > 0 {
		ee = append(ee, goqu.C("rel_compose_namespace").Eq(f.ComposeNamespaceID))
	}

	if f.FederationModuleID > 0 {
		ee = append(ee, goqu.C("rel_federation_module").Eq(f.FederationModuleID))
	}

	return ee, f, err
}

// FederationNodeFilter returns logical expressions
//
// This function is called from Store.QueryFederationNodes() and can be extended
// by setting Store.Filters.FederationNode. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func FederationNodeFilter(d drivers.Dialect, f federationType.NodeFilter) (ee []goqu.Expression, _ federationType.NodeFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("name").ILike("%"+f.Query+"%"),
			goqu.C("base_url").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// FederationNodeSyncFilter returns logical expressions
//
// This function is called from Store.QueryFederationNodeSyncs() and can be extended
// by setting Store.Filters.FederationNodeSync. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func FederationNodeSyncFilter(d drivers.Dialect, f federationType.NodeSyncFilter) (ee []goqu.Expression, _ federationType.NodeSyncFilter, err error) {

	if f.NodeID > 0 {
		ee = append(ee, goqu.C("rel_node").Eq(f.NodeID))
	}

	if f.ModuleID > 0 {
		ee = append(ee, goqu.C("rel_module").Eq(f.ModuleID))
	}

	if val := strings.TrimSpace(f.SyncStatus); len(val) > 0 {
		ee = append(ee, goqu.C("sync_status").Eq(f.SyncStatus))
	}

	if val := strings.TrimSpace(f.SyncType); len(val) > 0 {
		ee = append(ee, goqu.C("sync_type").Eq(f.SyncType))
	}

	return ee, f, err
}

// FederationSharedModuleFilter returns logical expressions
//
// This function is called from Store.QueryFederationSharedModules() and can be extended
// by setting Store.Filters.FederationSharedModule. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func FederationSharedModuleFilter(d drivers.Dialect, f federationType.SharedModuleFilter) (ee []goqu.Expression, _ federationType.SharedModuleFilter, err error) {

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if f.NodeID > 0 {
		ee = append(ee, goqu.C("rel_node").Eq(f.NodeID))
	}

	if val := strings.TrimSpace(f.Name); len(val) > 0 {
		ee = append(ee, goqu.C("name").Eq(f.Name))
	}

	if f.ExternalFederationModuleID > 0 {
		ee = append(ee, goqu.C("xref_module").Eq(f.ExternalFederationModuleID))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("name").ILike("%"+f.Query+"%"),
			goqu.C("handle").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// FlagFilter returns logical expressions
//
// This function is called from Store.QueryFlags() and can be extended
// by setting Store.Filters.Flag. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func FlagFilter(d drivers.Dialect, f flagType.FlagFilter) (ee []goqu.Expression, _ flagType.FlagFilter, err error) {

	if val := strings.TrimSpace(f.Kind); len(val) > 0 {
		ee = append(ee, goqu.C("kind").Eq(f.Kind))
	}

	if len(f.ResourceID) > 0 {
		ee = append(ee, goqu.C("rel_resource").In(f.ResourceID))
	}

	if len(f.OwnedBy) > 0 {
		ee = append(ee, goqu.C("owned_by").In(f.OwnedBy))
	}

	if ss := trimStringSlice(f.Name); len(ss) > 0 {
		ee = append(ee, goqu.C("name").In(ss))
	}

	return ee, f, err
}

// LabelFilter returns logical expressions
//
// This function is called from Store.QueryLabels() and can be extended
// by setting Store.Filters.Label. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func LabelFilter(d drivers.Dialect, f labelsType.LabelFilter) (ee []goqu.Expression, _ labelsType.LabelFilter, err error) {

	if val := strings.TrimSpace(f.Kind); len(val) > 0 {
		ee = append(ee, goqu.C("kind").Eq(f.Kind))
	}

	if len(f.ResourceID) > 0 {
		ee = append(ee, goqu.C("rel_resource").In(f.ResourceID))
	}

	return ee, f, err
}

// QueueFilter returns logical expressions
//
// This function is called from Store.QueryQueues() and can be extended
// by setting Store.Filters.Queue. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func QueueFilter(d drivers.Dialect, f systemType.QueueFilter) (ee []goqu.Expression, _ systemType.QueueFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("queue").ILike("%"+f.Query+"%"),
			goqu.C("consumer").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// QueueMessageFilter returns logical expressions
//
// This function is called from Store.QueryQueueMessages() and can be extended
// by setting Store.Filters.QueueMessage. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func QueueMessageFilter(d drivers.Dialect, f systemType.QueueMessageFilter) (ee []goqu.Expression, _ systemType.QueueMessageFilter, err error) {

	if expr := stateNilComparison(d, "processed", f.Processed); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Queue); len(val) > 0 {
		ee = append(ee, goqu.C("queue").Eq(f.Queue))
	}

	return ee, f, err
}

// RbacRuleFilter returns logical expressions
//
// This function is called from Store.QueryRbacRules() and can be extended
// by setting Store.Filters.RbacRule. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func RbacRuleFilter(d drivers.Dialect, f rbacType.RuleFilter) (ee []goqu.Expression, _ rbacType.RuleFilter, err error) {

	return ee, f, err
}

// ReminderFilter returns logical expressions
//
// This function is called from Store.QueryReminders() and can be extended
// by setting Store.Filters.Reminder. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ReminderFilter(d drivers.Dialect, f systemType.ReminderFilter) (ee []goqu.Expression, _ systemType.ReminderFilter, err error) {

	if len(f.ReminderID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ReminderID))
	}

	if f.AssignedTo > 0 {
		ee = append(ee, goqu.C("assigned_to").Eq(f.AssignedTo))
	}

	return ee, f, err
}

// ReportFilter returns logical expressions
//
// This function is called from Store.QueryReports() and can be extended
// by setting Store.Filters.Report. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ReportFilter(d drivers.Dialect, f systemType.ReportFilter) (ee []goqu.Expression, _ systemType.ReportFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if len(f.ReportID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ReportID))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ResourceActivityFilter returns logical expressions
//
// This function is called from Store.QueryResourceActivitys() and can be extended
// by setting Store.Filters.ResourceActivity. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ResourceActivityFilter(d drivers.Dialect, f discoveryType.ResourceActivityFilter) (ee []goqu.Expression, _ discoveryType.ResourceActivityFilter, err error) {

	return ee, f, err
}

// ResourceTranslationFilter returns logical expressions
//
// This function is called from Store.QueryResourceTranslations() and can be extended
// by setting Store.Filters.ResourceTranslation. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ResourceTranslationFilter(d drivers.Dialect, f systemType.ResourceTranslationFilter) (ee []goqu.Expression, _ systemType.ResourceTranslationFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Resource); len(val) > 0 {
		ee = append(ee, goqu.C("resource").Eq(f.Resource))
	}

	if val := strings.TrimSpace(f.Lang); len(val) > 0 {
		ee = append(ee, goqu.C("lang").Eq(f.Lang))
	}

	if len(f.TranslationID) > 0 {
		ee = append(ee, goqu.C("translation_id").In(f.TranslationID))
	}

	return ee, f, err
}

// RoleFilter returns logical expressions
//
// This function is called from Store.QueryRoles() and can be extended
// by setting Store.Filters.Role. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func RoleFilter(d drivers.Dialect, f systemType.RoleFilter) (ee []goqu.Expression, _ systemType.RoleFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "archived_at", f.Archived); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.RoleID) > 0 {
		ee = append(ee, goqu.C("id").In(f.RoleID))
	}

	if val := strings.TrimSpace(f.Name); len(val) > 0 {
		ee = append(ee, goqu.C("name").Eq(f.Name))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("name").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// RoleMemberFilter returns logical expressions
//
// This function is called from Store.QueryRoleMembers() and can be extended
// by setting Store.Filters.RoleMember. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func RoleMemberFilter(d drivers.Dialect, f systemType.RoleMemberFilter) (ee []goqu.Expression, _ systemType.RoleMemberFilter, err error) {

	if f.UserID > 0 {
		ee = append(ee, goqu.C("rel_user").Eq(f.UserID))
	}

	if f.RoleID > 0 {
		ee = append(ee, goqu.C("rel_role").Eq(f.RoleID))
	}

	return ee, f, err
}

// SettingValueFilter returns logical expressions
//
// This function is called from Store.QuerySettingValues() and can be extended
// by setting Store.Filters.SettingValue. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func SettingValueFilter(d drivers.Dialect, f systemType.SettingsFilter) (ee []goqu.Expression, _ systemType.SettingsFilter, err error) {

	if f.OwnedBy > 0 {
		ee = append(ee, goqu.C("rel_owner").Eq(f.OwnedBy))
	}

	return ee, f, err
}

// TemplateFilter returns logical expressions
//
// This function is called from Store.QueryTemplates() and can be extended
// by setting Store.Filters.Template. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func TemplateFilter(d drivers.Dialect, f systemType.TemplateFilter) (ee []goqu.Expression, _ systemType.TemplateFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.TemplateID) > 0 {
		ee = append(ee, goqu.C("id").In(f.TemplateID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if f.Partial {
		ee = append(ee, goqu.C("partial").IsTrue())
	}

	if val := strings.TrimSpace(f.Type); len(val) > 0 {
		ee = append(ee, goqu.C("type").Eq(f.Type))
	}

	if f.OwnerID > 0 {
		ee = append(ee, goqu.C("rel_owner").Eq(f.OwnerID))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("type").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// UserFilter returns logical expressions
//
// This function is called from Store.QueryUsers() and can be extended
// by setting Store.Filters.User. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func UserFilter(d drivers.Dialect, f systemType.UserFilter) (ee []goqu.Expression, _ systemType.UserFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "suspended_at", f.Suspended); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.UserID) > 0 {
		ee = append(ee, goqu.C("id").In(f.UserID))
	}

	if val := strings.TrimSpace(f.Email); len(val) > 0 {
		ee = append(ee, goqu.C("email").Eq(f.Email))
	}

	if val := strings.TrimSpace(f.Username); len(val) > 0 {
		ee = append(ee, goqu.C("username").Eq(f.Username))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("email").ILike("%"+f.Query+"%"),
			goqu.C("username").ILike("%"+f.Query+"%"),
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("name").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// trimStringSlice is a utility to trim all of the string slice elements and omit empty ones
func trimStringSlice(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if t := strings.TrimSpace(s); len(t) > 0 {
			out = append(out, t)
		}
	}
	return out
}
