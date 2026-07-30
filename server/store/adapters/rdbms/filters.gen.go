package rdbms

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	automationType "github.com/crusttech/human/server/automation/types"
	composeType "github.com/crusttech/human/server/compose/types"
	discoveryType "github.com/crusttech/human/server/discovery/types"
	federationType "github.com/crusttech/human/server/federation/types"
	actionlogType "github.com/crusttech/human/server/pkg/actionlog"
	flagType "github.com/crusttech/human/server/pkg/flag/types"
	labelsType "github.com/crusttech/human/server/pkg/label/types"
	rbacType "github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers"
	systemType "github.com/crusttech/human/server/system/types"
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

		// optional agent filter function called after the generated function
		Agent func(*Store, systemType.AgentFilter) ([]goqu.Expression, systemType.AgentFilter, error)

		// optional aiConversation filter function called after the generated function
		AiConversation func(*Store, systemType.AiConversationFilter) ([]goqu.Expression, systemType.AiConversationFilter, error)

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

		// optional automationNgAutomation filter function called after the generated function
		AutomationNgAutomation func(*Store, automationType.NgAutomationFilter) ([]goqu.Expression, automationType.NgAutomationFilter, error)

		// optional automationSession filter function called after the generated function
		AutomationSession func(*Store, automationType.SessionFilter) ([]goqu.Expression, automationType.SessionFilter, error)

		// optional automationTrigger filter function called after the generated function
		AutomationTrigger func(*Store, automationType.TriggerFilter) ([]goqu.Expression, automationType.TriggerFilter, error)

		// optional automationWorkflow filter function called after the generated function
		AutomationWorkflow func(*Store, automationType.WorkflowFilter) ([]goqu.Expression, automationType.WorkflowFilter, error)

		// optional chatbot filter function called after the generated function
		Chatbot func(*Store, systemType.ChatbotFilter) ([]goqu.Expression, systemType.ChatbotFilter, error)

		// optional chatbotSession filter function called after the generated function
		ChatbotSession func(*Store, systemType.ChatbotSessionFilter) ([]goqu.Expression, systemType.ChatbotSessionFilter, error)

		// optional chatbotSessionHandoff filter function called after the generated function
		ChatbotSessionHandoff func(*Store, systemType.ChatbotSessionHandoffFilter) ([]goqu.Expression, systemType.ChatbotSessionHandoffFilter, error)

		// optional chatbotSessionStep filter function called after the generated function
		ChatbotSessionStep func(*Store, systemType.ChatbotSessionStepFilter) ([]goqu.Expression, systemType.ChatbotSessionStepFilter, error)

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

		// optional composePageLayout filter function called after the generated function
		ComposePageLayout func(*Store, composeType.PageLayoutFilter) ([]goqu.Expression, composeType.PageLayoutFilter, error)

		// optional configuredConnection filter function called after the generated function
		ConfiguredConnection func(*Store, systemType.ConfiguredConnectionFilter) ([]goqu.Expression, systemType.ConfiguredConnectionFilter, error)

		// optional connection filter function called after the generated function
		Connection func(*Store, systemType.ConnectionFilter) ([]goqu.Expression, systemType.ConnectionFilter, error)

		// optional credential filter function called after the generated function
		Credential func(*Store, systemType.CredentialFilter) ([]goqu.Expression, systemType.CredentialFilter, error)

		// optional dalConnection filter function called after the generated function
		DalConnection func(*Store, systemType.DalConnectionFilter) ([]goqu.Expression, systemType.DalConnectionFilter, error)

		// optional dalSchemaAlteration filter function called after the generated function
		DalSchemaAlteration func(*Store, systemType.DalSchemaAlterationFilter) ([]goqu.Expression, systemType.DalSchemaAlterationFilter, error)

		// optional dalSensitivityLevel filter function called after the generated function
		DalSensitivityLevel func(*Store, systemType.DalSensitivityLevelFilter) ([]goqu.Expression, systemType.DalSensitivityLevelFilter, error)

		// optional dataPrivacyRequest filter function called after the generated function
		DataPrivacyRequest func(*Store, systemType.DataPrivacyRequestFilter) ([]goqu.Expression, systemType.DataPrivacyRequestFilter, error)

		// optional dataPrivacyRequestComment filter function called after the generated function
		DataPrivacyRequestComment func(*Store, systemType.DataPrivacyRequestCommentFilter) ([]goqu.Expression, systemType.DataPrivacyRequestCommentFilter, error)

		// optional dmlConnection filter function called after the generated function
		DmlConnection func(*Store, systemType.DmlConnectionFilter) ([]goqu.Expression, systemType.DmlConnectionFilter, error)

		// optional dmlImportRun filter function called after the generated function
		DmlImportRun func(*Store, systemType.DmlImportRunFilter) ([]goqu.Expression, systemType.DmlImportRunFilter, error)

		// optional dmlMapping filter function called after the generated function
		DmlMapping func(*Store, systemType.DmlMappingFilter) ([]goqu.Expression, systemType.DmlMappingFilter, error)

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

		// optional knowledgeBase filter function called after the generated function
		KnowledgeBase func(*Store, systemType.KnowledgeBaseFilter) ([]goqu.Expression, systemType.KnowledgeBaseFilter, error)

		// optional label filter function called after the generated function
		Label func(*Store, labelsType.LabelFilter) ([]goqu.Expression, labelsType.LabelFilter, error)

		// optional llmProvider filter function called after the generated function
		LlmProvider func(*Store, systemType.LlmProviderFilter) ([]goqu.Expression, systemType.LlmProviderFilter, error)

		// optional notification filter function called after the generated function
		Notification func(*Store, systemType.NotificationFilter) ([]goqu.Expression, systemType.NotificationFilter, error)

		// optional project filter function called after the generated function
		Project func(*Store, systemType.ProjectFilter) ([]goqu.Expression, systemType.ProjectFilter, error)

		// optional projectAiSystem filter function called after the generated function
		ProjectAiSystem func(*Store, systemType.ProjectAiSystemFilter) ([]goqu.Expression, systemType.ProjectAiSystemFilter, error)

		// optional projectAiSystemEntry filter function called after the generated function
		ProjectAiSystemEntry func(*Store, systemType.ProjectAiSystemEntryFilter) ([]goqu.Expression, systemType.ProjectAiSystemEntryFilter, error)

		// optional projectBacklogItem filter function called after the generated function
		ProjectBacklogItem func(*Store, systemType.ProjectBacklogItemFilter) ([]goqu.Expression, systemType.ProjectBacklogItemFilter, error)

		// optional projectFeature filter function called after the generated function
		ProjectFeature func(*Store, systemType.ProjectFeatureFilter) ([]goqu.Expression, systemType.ProjectFeatureFilter, error)

		// optional projectFriaScenario filter function called after the generated function
		ProjectFriaScenario func(*Store, systemType.ProjectFriaScenarioFilter) ([]goqu.Expression, systemType.ProjectFriaScenarioFilter, error)

		// optional projectIncident filter function called after the generated function
		ProjectIncident func(*Store, systemType.ProjectIncidentFilter) ([]goqu.Expression, systemType.ProjectIncidentFilter, error)

		// optional projectMember filter function called after the generated function
		ProjectMember func(*Store, systemType.ProjectMemberFilter) ([]goqu.Expression, systemType.ProjectMemberFilter, error)

		// optional projectPrivacy filter function called after the generated function
		ProjectPrivacy func(*Store, systemType.ProjectPrivacyFilter) ([]goqu.Expression, systemType.ProjectPrivacyFilter, error)

		// optional projectReview filter function called after the generated function
		ProjectReview func(*Store, systemType.ProjectReviewFilter) ([]goqu.Expression, systemType.ProjectReviewFilter, error)

		// optional projectTask filter function called after the generated function
		ProjectTask func(*Store, systemType.ProjectTaskFilter) ([]goqu.Expression, systemType.ProjectTaskFilter, error)

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

		// optional tenant filter function called after the generated function
		Tenant func(*Store, systemType.TenantFilter) ([]goqu.Expression, systemType.TenantFilter, error)

		// optional tenantMembership filter function called after the generated function
		TenantMembership func(*Store, systemType.TenantMembershipFilter) ([]goqu.Expression, systemType.TenantMembershipFilter, error)

		// optional user filter function called after the generated function
		User func(*Store, systemType.UserFilter) ([]goqu.Expression, systemType.UserFilter, error)

		// optional userGroup filter function called after the generated function
		UserGroup func(*Store, systemType.UserGroupFilter) ([]goqu.Expression, systemType.UserGroupFilter, error)
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
		ee = append(ee, goqu.C("request_origin").Eq(f.Origin))
	}

	if len(f.ActorID) > 0 {
		ee = append(ee, goqu.C("actor_id").In(f.ActorID))
	}

	if f.TenantID > 0 {
		ee = append(ee, goqu.C("rel_tenant").Eq(f.TenantID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RootProjectID > 0 {
		ee = append(ee, goqu.C("rel_root_project").Eq(f.RootProjectID))
	}

	if f.ResourceProjectID > 0 {
		ee = append(ee, goqu.C("rel_resource_project").Eq(f.ResourceProjectID))
	}

	if f.ResourceRevisionID > 0 {
		ee = append(ee, goqu.C("rel_resource_revision").Eq(f.ResourceRevisionID))
	}

	return ee, f, err
}

// AgentFilter returns logical expressions
//
// This function is called from Store.QueryAgents() and can be extended
// by setting Store.Filters.Agent. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AgentFilter(d drivers.Dialect, f systemType.AgentFilter) (ee []goqu.Expression, _ systemType.AgentFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.AgentID) > 0 {
		ee = append(ee, goqu.C("id").In(f.AgentID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// AiConversationFilter returns logical expressions
//
// This function is called from Store.QueryAiConversations() and can be extended
// by setting Store.Filters.AiConversation. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AiConversationFilter(d drivers.Dialect, f systemType.AiConversationFilter) (ee []goqu.Expression, _ systemType.AiConversationFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.AiConversationID) > 0 {
		ee = append(ee, goqu.C("id").In(f.AiConversationID))
	}

	if f.AgentID > 0 {
		ee = append(ee, goqu.C("rel_agent").Eq(f.AgentID))
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

	if len(f.ApigwFilterID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ApigwFilterID))
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

	if len(f.ApigwRouteID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ApigwRouteID))
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

// AutomationNgAutomationFilter returns logical expressions
//
// This function is called from Store.QueryAutomationNgAutomations() and can be extended
// by setting Store.Filters.AutomationNgAutomation. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func AutomationNgAutomationFilter(d drivers.Dialect, f automationType.NgAutomationFilter) (ee []goqu.Expression, _ automationType.NgAutomationFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateFalseComparison(d, "enabled", f.Disabled); expr != nil {
		ee = append(ee, expr)
	}

	if ss := trimStringSlice(f.AutomationID); len(ss) > 0 {
		ee = append(ee, goqu.C("id").In(ss))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
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
		))
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

	if ss := trimStringSlice(f.SessionID); len(ss) > 0 {
		ee = append(ee, goqu.C("id").In(ss))
	}

	if ss := trimStringSlice(f.WorkflowID); len(ss) > 0 {
		ee = append(ee, goqu.C("rel_workflow").In(ss))
	}

	if val := strings.TrimSpace(f.EventType); len(val) > 0 {
		ee = append(ee, goqu.C("event_type").Eq(f.EventType))
	}

	if val := strings.TrimSpace(f.ResourceType); len(val) > 0 {
		ee = append(ee, goqu.C("resource_type").Eq(f.ResourceType))
	}

	if ss := trimStringSlice(f.CreatedBy); len(ss) > 0 {
		ee = append(ee, goqu.C("created_by").In(ss))
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

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
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

// ChatbotFilter returns logical expressions
//
// This function is called from Store.QueryChatbots() and can be extended
// by setting Store.Filters.Chatbot. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ChatbotFilter(d drivers.Dialect, f systemType.ChatbotFilter) (ee []goqu.Expression, _ systemType.ChatbotFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ChatbotID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ChatbotID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if val := strings.TrimSpace(f.WidgetKey); len(val) > 0 {
		ee = append(ee, goqu.C("widget_key").Eq(f.WidgetKey))
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

// ChatbotSessionFilter returns logical expressions
//
// This function is called from Store.QueryChatbotSessions() and can be extended
// by setting Store.Filters.ChatbotSession. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ChatbotSessionFilter(d drivers.Dialect, f systemType.ChatbotSessionFilter) (ee []goqu.Expression, _ systemType.ChatbotSessionFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ChatbotSessionID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ChatbotSessionID))
	}

	if f.ChatbotID > 0 {
		ee = append(ee, goqu.C("rel_chatbot").Eq(f.ChatbotID))
	}

	if ss := trimStringSlice(f.Status); len(ss) > 0 {
		ee = append(ee, goqu.C("status").In(ss))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ChatbotSessionHandoffFilter returns logical expressions
//
// This function is called from Store.QueryChatbotSessionHandoffs() and can be extended
// by setting Store.Filters.ChatbotSessionHandoff. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ChatbotSessionHandoffFilter(d drivers.Dialect, f systemType.ChatbotSessionHandoffFilter) (ee []goqu.Expression, _ systemType.ChatbotSessionHandoffFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ChatbotSessionHandoffID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ChatbotSessionHandoffID))
	}

	if f.SessionID > 0 {
		ee = append(ee, goqu.C("rel_session").Eq(f.SessionID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ChatbotSessionStepFilter returns logical expressions
//
// This function is called from Store.QueryChatbotSessionSteps() and can be extended
// by setting Store.Filters.ChatbotSessionStep. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ChatbotSessionStepFilter(d drivers.Dialect, f systemType.ChatbotSessionStepFilter) (ee []goqu.Expression, _ systemType.ChatbotSessionStepFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ChatbotSessionStepID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ChatbotSessionStepID))
	}

	if f.SessionID > 0 {
		ee = append(ee, goqu.C("rel_session").Eq(f.SessionID))
	}

	if f.ConversationID > 0 {
		ee = append(ee, goqu.C("rel_conversation").Eq(f.ConversationID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("status").ILike("%"+f.Query+"%"),
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

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
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

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
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
		ee = append(ee, goqu.C("id").In(f.NamespaceID))
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

	if len(f.PageID) > 0 {
		ee = append(ee, goqu.C("id").In(f.PageID))
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

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
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

// ComposePageLayoutFilter returns logical expressions
//
// This function is called from Store.QueryComposePageLayouts() and can be extended
// by setting Store.Filters.ComposePageLayout. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ComposePageLayoutFilter(d drivers.Dialect, f composeType.PageLayoutFilter) (ee []goqu.Expression, _ composeType.PageLayoutFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if f.ParentID > 0 {
		ee = append(ee, goqu.C("parent_id").Eq(f.ParentID))
	}

	if f.NamespaceID > 0 {
		ee = append(ee, goqu.C("rel_namespace").Eq(f.NamespaceID))
	}

	if f.PageID > 0 {
		ee = append(ee, goqu.C("page_id").Eq(f.PageID))
	}

	if len(f.PageLayoutID) > 0 {
		ee = append(ee, goqu.C("id").In(f.PageLayoutID))
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

// ConfiguredConnectionFilter returns logical expressions
//
// This function is called from Store.QueryConfiguredConnections() and can be extended
// by setting Store.Filters.ConfiguredConnection. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ConfiguredConnectionFilter(d drivers.Dialect, f systemType.ConfiguredConnectionFilter) (ee []goqu.Expression, _ systemType.ConfiguredConnectionFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if f.ConnectionID > 0 {
		ee = append(ee, goqu.C("rel_connection").Eq(f.ConnectionID))
	}

	if ss := trimStringSlice(f.Status); len(ss) > 0 {
		ee = append(ee, goqu.C("status").In(ss))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
	}

	return ee, f, err
}

// ConnectionFilter returns logical expressions
//
// This function is called from Store.QueryConnections() and can be extended
// by setting Store.Filters.Connection. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ConnectionFilter(d drivers.Dialect, f systemType.ConnectionFilter) (ee []goqu.Expression, _ systemType.ConnectionFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if ss := trimStringSlice(f.Status); len(ss) > 0 {
		ee = append(ee, goqu.C("status").In(ss))
	}

	if val := strings.TrimSpace(f.Source); len(val) > 0 {
		ee = append(ee, goqu.C("source").Eq(f.Source))
	}

	if len(f.LabeledIDs) > 0 {
		ee = append(ee, goqu.I("id").In(f.LabeledIDs))
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

	if len(f.DalConnectionID) > 0 {
		ee = append(ee, goqu.C("id").In(f.DalConnectionID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if val := strings.TrimSpace(f.Type); len(val) > 0 {
		ee = append(ee, goqu.C("type").Eq(f.Type))
	}

	return ee, f, err
}

// DalSchemaAlterationFilter returns logical expressions
//
// This function is called from Store.QueryDalSchemaAlterations() and can be extended
// by setting Store.Filters.DalSchemaAlteration. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DalSchemaAlterationFilter(d drivers.Dialect, f systemType.DalSchemaAlterationFilter) (ee []goqu.Expression, _ systemType.DalSchemaAlterationFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "completed_at", f.Completed); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "dismissed_at", f.Dismissed); expr != nil {
		ee = append(ee, expr)
	}

	if val := strings.TrimSpace(f.Kind); len(val) > 0 {
		ee = append(ee, goqu.C("kind").Eq(f.Kind))
	}

	if ss := trimStringSlice(f.Resource); len(ss) > 0 {
		ee = append(ee, goqu.C("resource").In(ss))
	}

	if val := strings.TrimSpace(f.ResourceType); len(val) > 0 {
		ee = append(ee, goqu.C("resource_type").Eq(f.ResourceType))
	}

	if len(f.AlterationID) > 0 {
		ee = append(ee, goqu.C("id").In(f.AlterationID))
	}

	if len(f.BatchID) > 0 {
		ee = append(ee, goqu.C("batch_id").In(f.BatchID))
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

	if len(f.DalSensitivityLevelID) > 0 {
		ee = append(ee, goqu.C("id").In(f.DalSensitivityLevelID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
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

// DmlConnectionFilter returns logical expressions
//
// This function is called from Store.QueryDmlConnections() and can be extended
// by setting Store.Filters.DmlConnection. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DmlConnectionFilter(d drivers.Dialect, f systemType.DmlConnectionFilter) (ee []goqu.Expression, _ systemType.DmlConnectionFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ConnectionID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ConnectionID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	return ee, f, err
}

// DmlImportRunFilter returns logical expressions
//
// This function is called from Store.QueryDmlImportRuns() and can be extended
// by setting Store.Filters.DmlImportRun. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DmlImportRunFilter(d drivers.Dialect, f systemType.DmlImportRunFilter) (ee []goqu.Expression, _ systemType.DmlImportRunFilter, err error) {

	if len(f.ImportRunID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ImportRunID))
	}

	if f.MappingID > 0 {
		ee = append(ee, goqu.C("rel_mapping").Eq(f.MappingID))
	}

	if ss := trimStringSlice(f.Status); len(ss) > 0 {
		ee = append(ee, goqu.C("status").In(ss))
	}

	return ee, f, err
}

// DmlMappingFilter returns logical expressions
//
// This function is called from Store.QueryDmlMappings() and can be extended
// by setting Store.Filters.DmlMapping. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func DmlMappingFilter(d drivers.Dialect, f systemType.DmlMappingFilter) (ee []goqu.Expression, _ systemType.DmlMappingFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.MappingID) > 0 {
		ee = append(ee, goqu.C("id").In(f.MappingID))
	}

	if f.ConnectionID > 0 {
		ee = append(ee, goqu.C("rel_connection").Eq(f.ConnectionID))
	}

	if val := strings.TrimSpace(f.SourceIdent); len(val) > 0 {
		ee = append(ee, goqu.C("source_ident").Eq(f.SourceIdent))
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

// KnowledgeBaseFilter returns logical expressions
//
// This function is called from Store.QueryKnowledgeBases() and can be extended
// by setting Store.Filters.KnowledgeBase. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func KnowledgeBaseFilter(d drivers.Dialect, f systemType.KnowledgeBaseFilter) (ee []goqu.Expression, _ systemType.KnowledgeBaseFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.KnowledgeBaseID) > 0 {
		ee = append(ee, goqu.C("id").In(f.KnowledgeBaseID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
			goqu.C("title").ILike("%"+f.Query+"%"),
		))
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

// LlmProviderFilter returns logical expressions
//
// This function is called from Store.QueryLlmProviders() and can be extended
// by setting Store.Filters.LlmProvider. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func LlmProviderFilter(d drivers.Dialect, f systemType.LlmProviderFilter) (ee []goqu.Expression, _ systemType.LlmProviderFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.LlmProviderID) > 0 {
		ee = append(ee, goqu.C("id").In(f.LlmProviderID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if val := strings.TrimSpace(f.Provider); len(val) > 0 {
		ee = append(ee, goqu.C("provider").Eq(f.Provider))
	}

	return ee, f, err
}

// NotificationFilter returns logical expressions
//
// This function is called from Store.QueryNotifications() and can be extended
// by setting Store.Filters.Notification. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func NotificationFilter(d drivers.Dialect, f systemType.NotificationFilter) (ee []goqu.Expression, _ systemType.NotificationFilter, err error) {

	if expr := stateNilComparison(d, "read_at", f.Read); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.NotificationID) > 0 {
		ee = append(ee, goqu.C("id").In(f.NotificationID))
	}

	if f.Recipient > 0 {
		ee = append(ee, goqu.C("recipient").Eq(f.Recipient))
	}

	return ee, f, err
}

// ProjectFilter returns logical expressions
//
// This function is called from Store.QueryProjects() and can be extended
// by setting Store.Filters.Project. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectFilter(d drivers.Dialect, f systemType.ProjectFilter) (ee []goqu.Expression, _ systemType.ProjectFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "archived_at", f.Archived); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ProjectID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ProjectID))
	}

	if f.RootProjectID > 0 {
		ee = append(ee, goqu.C("root_project_id").Eq(f.RootProjectID))
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
		))
	}

	return ee, f, err
}

// ProjectAiSystemFilter returns logical expressions
//
// This function is called from Store.QueryProjectAiSystems() and can be extended
// by setting Store.Filters.ProjectAiSystem. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectAiSystemFilter(d drivers.Dialect, f systemType.ProjectAiSystemFilter) (ee []goqu.Expression, _ systemType.ProjectAiSystemFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ProjectAiSystemID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ProjectAiSystemID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if val := strings.TrimSpace(f.Handle); len(val) > 0 {
		ee = append(ee, goqu.C("handle").Eq(f.Handle))
	}

	if val := strings.TrimSpace(f.RiskClass); len(val) > 0 {
		ee = append(ee, goqu.C("risk_class").Eq(f.RiskClass))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("handle").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectAiSystemEntryFilter returns logical expressions
//
// This function is called from Store.QueryProjectAiSystemEntrys() and can be extended
// by setting Store.Filters.ProjectAiSystemEntry. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectAiSystemEntryFilter(d drivers.Dialect, f systemType.ProjectAiSystemEntryFilter) (ee []goqu.Expression, _ systemType.ProjectAiSystemEntryFilter, err error) {

	if f.ProjectAiSystemID > 0 {
		ee = append(ee, goqu.C("rel_project_ai_system").Eq(f.ProjectAiSystemID))
	}

	if val := strings.TrimSpace(f.ResourceRef); len(val) > 0 {
		ee = append(ee, goqu.C("resource_ref").Eq(f.ResourceRef))
	}

	return ee, f, err
}

// ProjectBacklogItemFilter returns logical expressions
//
// This function is called from Store.QueryProjectBacklogItems() and can be extended
// by setting Store.Filters.ProjectBacklogItem. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectBacklogItemFilter(d drivers.Dialect, f systemType.ProjectBacklogItemFilter) (ee []goqu.Expression, _ systemType.ProjectBacklogItemFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RevisionID > 0 {
		ee = append(ee, goqu.C("rel_revision").Eq(f.RevisionID))
	}

	if f.EventID > 0 {
		ee = append(ee, goqu.C("event_id").Eq(f.EventID))
	}

	if val := strings.TrimSpace(f.Category); len(val) > 0 {
		ee = append(ee, goqu.C("category").Eq(f.Category))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectFeatureFilter returns logical expressions
//
// This function is called from Store.QueryProjectFeatures() and can be extended
// by setting Store.Filters.ProjectFeature. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectFeatureFilter(d drivers.Dialect, f systemType.ProjectFeatureFilter) (ee []goqu.Expression, _ systemType.ProjectFeatureFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.FeatureID) > 0 {
		ee = append(ee, goqu.C("id").In(f.FeatureID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RevisionID > 0 {
		ee = append(ee, goqu.C("rel_revision").Eq(f.RevisionID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectFriaScenarioFilter returns logical expressions
//
// This function is called from Store.QueryProjectFriaScenarios() and can be extended
// by setting Store.Filters.ProjectFriaScenario. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectFriaScenarioFilter(d drivers.Dialect, f systemType.ProjectFriaScenarioFilter) (ee []goqu.Expression, _ systemType.ProjectFriaScenarioFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ProjectFriaScenarioID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ProjectFriaScenarioID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.AiSystemID > 0 {
		ee = append(ee, goqu.C("rel_ai_system").Eq(f.AiSystemID))
	}

	if val := strings.TrimSpace(f.Severity); len(val) > 0 {
		ee = append(ee, goqu.C("severity").Eq(f.Severity))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectIncidentFilter returns logical expressions
//
// This function is called from Store.QueryProjectIncidents() and can be extended
// by setting Store.Filters.ProjectIncident. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectIncidentFilter(d drivers.Dialect, f systemType.ProjectIncidentFilter) (ee []goqu.Expression, _ systemType.ProjectIncidentFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.IncidentID) > 0 {
		ee = append(ee, goqu.C("id").In(f.IncidentID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RevisionID > 0 {
		ee = append(ee, goqu.C("rel_revision").Eq(f.RevisionID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectMemberFilter returns logical expressions
//
// This function is called from Store.QueryProjectMembers() and can be extended
// by setting Store.Filters.ProjectMember. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectMemberFilter(d drivers.Dialect, f systemType.ProjectMemberFilter) (ee []goqu.Expression, _ systemType.ProjectMemberFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ProjectMemberID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ProjectMemberID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.UserID > 0 {
		ee = append(ee, goqu.C("rel_user").Eq(f.UserID))
	}

	return ee, f, err
}

// ProjectPrivacyFilter returns logical expressions
//
// This function is called from Store.QueryProjectPrivacys() and can be extended
// by setting Store.Filters.ProjectPrivacy. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectPrivacyFilter(d drivers.Dialect, f systemType.ProjectPrivacyFilter) (ee []goqu.Expression, _ systemType.ProjectPrivacyFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.PrivacyID) > 0 {
		ee = append(ee, goqu.C("id").In(f.PrivacyID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RevisionID > 0 {
		ee = append(ee, goqu.C("rel_revision").Eq(f.RevisionID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectReviewFilter returns logical expressions
//
// This function is called from Store.QueryProjectReviews() and can be extended
// by setting Store.Filters.ProjectReview. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectReviewFilter(d drivers.Dialect, f systemType.ProjectReviewFilter) (ee []goqu.Expression, _ systemType.ProjectReviewFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.ReviewID) > 0 {
		ee = append(ee, goqu.C("id").In(f.ReviewID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RevisionID > 0 {
		ee = append(ee, goqu.C("rel_revision").Eq(f.RevisionID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
	}

	return ee, f, err
}

// ProjectTaskFilter returns logical expressions
//
// This function is called from Store.QueryProjectTasks() and can be extended
// by setting Store.Filters.ProjectTask. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func ProjectTaskFilter(d drivers.Dialect, f systemType.ProjectTaskFilter) (ee []goqu.Expression, _ systemType.ProjectTaskFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.TaskID) > 0 {
		ee = append(ee, goqu.C("id").In(f.TaskID))
	}

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
	}

	if f.RevisionID > 0 {
		ee = append(ee, goqu.C("rel_revision").Eq(f.RevisionID))
	}

	if val := strings.TrimSpace(f.Status); len(val) > 0 {
		ee = append(ee, goqu.C("status").Eq(f.Status))
	}

	if f.Query != "" {
		ee = append(ee, goqu.Or(
			goqu.C("title").ILike("%"+f.Query+"%"),
			goqu.C("status").ILike("%"+f.Query+"%"),
		))
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

	if len(f.QueueID) > 0 {
		ee = append(ee, goqu.C("id").In(f.QueueID))
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

	if f.ProjectID > 0 {
		ee = append(ee, goqu.C("rel_project").Eq(f.ProjectID))
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

	if val := strings.TrimSpace(f.Resource); len(val) > 0 {
		ee = append(ee, goqu.C("rel_resource").Eq(f.Resource))
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

// TenantFilter returns logical expressions
//
// This function is called from Store.QueryTenants() and can be extended
// by setting Store.Filters.Tenant. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func TenantFilter(d drivers.Dialect, f systemType.TenantFilter) (ee []goqu.Expression, _ systemType.TenantFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.TenantID) > 0 {
		ee = append(ee, goqu.C("id").In(f.TenantID))
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
		))
	}

	return ee, f, err
}

// TenantMembershipFilter returns logical expressions
//
// This function is called from Store.QueryTenantMemberships() and can be extended
// by setting Store.Filters.TenantMembership. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func TenantMembershipFilter(d drivers.Dialect, f systemType.TenantMembershipFilter) (ee []goqu.Expression, _ systemType.TenantMembershipFilter, err error) {

	if len(f.TenantMembershipID) > 0 {
		ee = append(ee, goqu.C("id").In(f.TenantMembershipID))
	}

	if f.TenantID > 0 {
		ee = append(ee, goqu.C("rel_tenant").Eq(f.TenantID))
	}

	if f.UserID > 0 {
		ee = append(ee, goqu.C("rel_user").Eq(f.UserID))
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

// UserGroupFilter returns logical expressions
//
// This function is called from Store.QueryUserGroups() and can be extended
// by setting Store.Filters.UserGroup. Extension is called after all expressions
// are generated and can choose to ignore or alter them.
//
// This function is auto-generated
func UserGroupFilter(d drivers.Dialect, f systemType.UserGroupFilter) (ee []goqu.Expression, _ systemType.UserGroupFilter, err error) {

	if expr := stateNilComparison(d, "deleted_at", f.Deleted); expr != nil {
		ee = append(ee, expr)
	}

	if expr := stateNilComparison(d, "archived_at", f.Archived); expr != nil {
		ee = append(ee, expr)
	}

	if len(f.UserGroupID) > 0 {
		ee = append(ee, goqu.C("id").In(f.UserGroupID))
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
