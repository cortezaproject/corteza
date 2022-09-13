package store

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	automationType "github.com/cortezaproject/corteza-server/automation/types"
	composeType "github.com/cortezaproject/corteza-server/compose/types"
	federationType "github.com/cortezaproject/corteza-server/federation/types"
	actionlogType "github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	discoveryType "github.com/cortezaproject/corteza-server/pkg/discovery/types"
	flagType "github.com/cortezaproject/corteza-server/pkg/flag/types"
	labelsType "github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	rbacType "github.com/cortezaproject/corteza-server/pkg/rbac"
	systemType "github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type (
	// Storer interface combines interfaces of all supported store interfaces
	Storer interface {
		// SetLogger sets new logging facility
		//
		// Store facility should fallback to logger.Default when no logging facility is set
		//
		// Intentionally closely coupled with Zap logger since this is not some public lib
		// and it's highly unlikely we'll support different/multiple logging "backend"
		SetLogger(*zap.Logger)

		// Returns underlying store as DAL connection
		ToDalConn() dal.Connection

		// Tx is a transaction handler
		Tx(context.Context, func(context.Context, Storer) error) error

		// Upgrade store's schema to the latest version
		Upgrade(context.Context) error
		Actionlogs
		ApigwFilters
		ApigwRoutes
		Applications
		Attachments
		AuthClients
		AuthConfirmedClients
		AuthOa2tokens
		AuthSessions
		AutomationSessions
		AutomationTriggers
		AutomationWorkflows
		ComposeAttachments
		ComposeCharts
		ComposeModules
		ComposeModuleFields
		ComposeNamespaces
		ComposePages
		Credentials
		DalConnections
		DalSensitivityLevels
		DataPrivacyRequests
		DataPrivacyRequestComments
		FederationExposedModules
		FederationModuleMappings
		FederationNodes
		FederationNodeSyncs
		FederationSharedModules
		Flags
		Labels
		Queues
		QueueMessages
		RbacRules
		Reminders
		Reports
		ResourceActivitys
		ResourceTranslations
		Roles
		RoleMembers
		SettingValues
		Templates
		Users
	}

	Actionlogs interface {
		SearchActionlogs(ctx context.Context, f actionlogType.Filter) (actionlogType.ActionSet, actionlogType.Filter, error)
		CreateActionlog(ctx context.Context, rr ...*actionlogType.Action) error
		UpdateActionlog(ctx context.Context, rr ...*actionlogType.Action) error
		UpsertActionlog(ctx context.Context, rr ...*actionlogType.Action) error
		DeleteActionlog(ctx context.Context, rr ...*actionlogType.Action) error
		DeleteActionlogByID(ctx context.Context, id uint64) error
		TruncateActionlogs(ctx context.Context) error
		LookupActionlogByID(ctx context.Context, id uint64) (*actionlogType.Action, error)
	}

	ApigwFilters interface {
		SearchApigwFilters(ctx context.Context, f systemType.ApigwFilterFilter) (systemType.ApigwFilterSet, systemType.ApigwFilterFilter, error)
		CreateApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) error
		UpdateApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) error
		UpsertApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) error
		DeleteApigwFilter(ctx context.Context, rr ...*systemType.ApigwFilter) error
		DeleteApigwFilterByID(ctx context.Context, id uint64) error
		TruncateApigwFilters(ctx context.Context) error
		LookupApigwFilterByID(ctx context.Context, id uint64) (*systemType.ApigwFilter, error)
		LookupApigwFilterByRoute(ctx context.Context, route uint64) (*systemType.ApigwFilter, error)
	}

	ApigwRoutes interface {
		SearchApigwRoutes(ctx context.Context, f systemType.ApigwRouteFilter) (systemType.ApigwRouteSet, systemType.ApigwRouteFilter, error)
		CreateApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) error
		UpdateApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) error
		UpsertApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) error
		DeleteApigwRoute(ctx context.Context, rr ...*systemType.ApigwRoute) error
		DeleteApigwRouteByID(ctx context.Context, id uint64) error
		TruncateApigwRoutes(ctx context.Context) error
		LookupApigwRouteByID(ctx context.Context, id uint64) (*systemType.ApigwRoute, error)
		LookupApigwRouteByEndpoint(ctx context.Context, endpoint string) (*systemType.ApigwRoute, error)
	}

	Applications interface {
		SearchApplications(ctx context.Context, f systemType.ApplicationFilter) (systemType.ApplicationSet, systemType.ApplicationFilter, error)
		CreateApplication(ctx context.Context, rr ...*systemType.Application) error
		UpdateApplication(ctx context.Context, rr ...*systemType.Application) error
		UpsertApplication(ctx context.Context, rr ...*systemType.Application) error
		DeleteApplication(ctx context.Context, rr ...*systemType.Application) error
		DeleteApplicationByID(ctx context.Context, id uint64) error
		TruncateApplications(ctx context.Context) error
		LookupApplicationByID(ctx context.Context, id uint64) (*systemType.Application, error)
		ApplicationMetrics(ctx context.Context) (*systemType.ApplicationMetrics, error)
		ReorderApplications(ctx context.Context, order []uint64) error
	}

	Attachments interface {
		SearchAttachments(ctx context.Context, f systemType.AttachmentFilter) (systemType.AttachmentSet, systemType.AttachmentFilter, error)
		CreateAttachment(ctx context.Context, rr ...*systemType.Attachment) error
		UpdateAttachment(ctx context.Context, rr ...*systemType.Attachment) error
		UpsertAttachment(ctx context.Context, rr ...*systemType.Attachment) error
		DeleteAttachment(ctx context.Context, rr ...*systemType.Attachment) error
		DeleteAttachmentByID(ctx context.Context, id uint64) error
		TruncateAttachments(ctx context.Context) error
		LookupAttachmentByID(ctx context.Context, id uint64) (*systemType.Attachment, error)
	}

	AuthClients interface {
		SearchAuthClients(ctx context.Context, f systemType.AuthClientFilter) (systemType.AuthClientSet, systemType.AuthClientFilter, error)
		CreateAuthClient(ctx context.Context, rr ...*systemType.AuthClient) error
		UpdateAuthClient(ctx context.Context, rr ...*systemType.AuthClient) error
		UpsertAuthClient(ctx context.Context, rr ...*systemType.AuthClient) error
		DeleteAuthClient(ctx context.Context, rr ...*systemType.AuthClient) error
		DeleteAuthClientByID(ctx context.Context, id uint64) error
		TruncateAuthClients(ctx context.Context) error
		LookupAuthClientByID(ctx context.Context, id uint64) (*systemType.AuthClient, error)
		LookupAuthClientByHandle(ctx context.Context, handle string) (*systemType.AuthClient, error)
	}

	AuthConfirmedClients interface {
		SearchAuthConfirmedClients(ctx context.Context, f systemType.AuthConfirmedClientFilter) (systemType.AuthConfirmedClientSet, systemType.AuthConfirmedClientFilter, error)
		CreateAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) error
		UpdateAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) error
		UpsertAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) error
		DeleteAuthConfirmedClient(ctx context.Context, rr ...*systemType.AuthConfirmedClient) error
		DeleteAuthConfirmedClientByUserIDClientID(ctx context.Context, userID uint64, clientID uint64) error
		TruncateAuthConfirmedClients(ctx context.Context) error
		LookupAuthConfirmedClientByUserIDClientID(ctx context.Context, userID uint64, clientID uint64) (*systemType.AuthConfirmedClient, error)
	}

	AuthOa2tokens interface {
		SearchAuthOa2tokens(ctx context.Context, f systemType.AuthOa2tokenFilter) (systemType.AuthOa2tokenSet, systemType.AuthOa2tokenFilter, error)
		CreateAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) error
		UpdateAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) error
		UpsertAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) error
		DeleteAuthOa2token(ctx context.Context, rr ...*systemType.AuthOa2token) error
		DeleteAuthOa2tokenByID(ctx context.Context, id uint64) error
		TruncateAuthOa2tokens(ctx context.Context) error
		LookupAuthOa2tokenByID(ctx context.Context, id uint64) (*systemType.AuthOa2token, error)
		LookupAuthOa2tokenByCode(ctx context.Context, code string) (*systemType.AuthOa2token, error)
		LookupAuthOa2tokenByAccess(ctx context.Context, access string) (*systemType.AuthOa2token, error)
		LookupAuthOa2tokenByRefresh(ctx context.Context, refresh string) (*systemType.AuthOa2token, error)
		DeleteExpiredAuthOA2Tokens(ctx context.Context) error
		DeleteAuthOA2TokenByCode(ctx context.Context, code string) error
		DeleteAuthOA2TokenByAccess(ctx context.Context, access string) error
		DeleteAuthOA2TokenByRefresh(ctx context.Context, refresh string) error
		DeleteAuthOA2TokenByUserID(ctx context.Context, userID uint64) error
	}

	AuthSessions interface {
		SearchAuthSessions(ctx context.Context, f systemType.AuthSessionFilter) (systemType.AuthSessionSet, systemType.AuthSessionFilter, error)
		CreateAuthSession(ctx context.Context, rr ...*systemType.AuthSession) error
		UpdateAuthSession(ctx context.Context, rr ...*systemType.AuthSession) error
		UpsertAuthSession(ctx context.Context, rr ...*systemType.AuthSession) error
		DeleteAuthSession(ctx context.Context, rr ...*systemType.AuthSession) error
		DeleteAuthSessionByID(ctx context.Context, id string) error
		TruncateAuthSessions(ctx context.Context) error
		LookupAuthSessionByID(ctx context.Context, id string) (*systemType.AuthSession, error)
		DeleteExpiredAuthSessions(ctx context.Context) error
		DeleteAuthSessionsByUserID(ctx context.Context, userID uint64) error
	}

	AutomationSessions interface {
		SearchAutomationSessions(ctx context.Context, f automationType.SessionFilter) (automationType.SessionSet, automationType.SessionFilter, error)
		CreateAutomationSession(ctx context.Context, rr ...*automationType.Session) error
		UpdateAutomationSession(ctx context.Context, rr ...*automationType.Session) error
		UpsertAutomationSession(ctx context.Context, rr ...*automationType.Session) error
		DeleteAutomationSession(ctx context.Context, rr ...*automationType.Session) error
		DeleteAutomationSessionByID(ctx context.Context, id uint64) error
		TruncateAutomationSessions(ctx context.Context) error
		LookupAutomationSessionByID(ctx context.Context, id uint64) (*automationType.Session, error)
	}

	AutomationTriggers interface {
		SearchAutomationTriggers(ctx context.Context, f automationType.TriggerFilter) (automationType.TriggerSet, automationType.TriggerFilter, error)
		CreateAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) error
		UpdateAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) error
		UpsertAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) error
		DeleteAutomationTrigger(ctx context.Context, rr ...*automationType.Trigger) error
		DeleteAutomationTriggerByID(ctx context.Context, id uint64) error
		TruncateAutomationTriggers(ctx context.Context) error
		LookupAutomationTriggerByID(ctx context.Context, id uint64) (*automationType.Trigger, error)
	}

	AutomationWorkflows interface {
		SearchAutomationWorkflows(ctx context.Context, f automationType.WorkflowFilter) (automationType.WorkflowSet, automationType.WorkflowFilter, error)
		CreateAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) error
		UpdateAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) error
		UpsertAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) error
		DeleteAutomationWorkflow(ctx context.Context, rr ...*automationType.Workflow) error
		DeleteAutomationWorkflowByID(ctx context.Context, id uint64) error
		TruncateAutomationWorkflows(ctx context.Context) error
		LookupAutomationWorkflowByID(ctx context.Context, id uint64) (*automationType.Workflow, error)
		LookupAutomationWorkflowByHandle(ctx context.Context, handle string) (*automationType.Workflow, error)
	}

	ComposeAttachments interface {
		SearchComposeAttachments(ctx context.Context, f composeType.AttachmentFilter) (composeType.AttachmentSet, composeType.AttachmentFilter, error)
		CreateComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) error
		UpdateComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) error
		UpsertComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) error
		DeleteComposeAttachment(ctx context.Context, rr ...*composeType.Attachment) error
		DeleteComposeAttachmentByID(ctx context.Context, id uint64) error
		TruncateComposeAttachments(ctx context.Context) error
		LookupComposeAttachmentByID(ctx context.Context, id uint64) (*composeType.Attachment, error)
	}

	ComposeCharts interface {
		SearchComposeCharts(ctx context.Context, f composeType.ChartFilter) (composeType.ChartSet, composeType.ChartFilter, error)
		CreateComposeChart(ctx context.Context, rr ...*composeType.Chart) error
		UpdateComposeChart(ctx context.Context, rr ...*composeType.Chart) error
		UpsertComposeChart(ctx context.Context, rr ...*composeType.Chart) error
		DeleteComposeChart(ctx context.Context, rr ...*composeType.Chart) error
		DeleteComposeChartByID(ctx context.Context, id uint64) error
		TruncateComposeCharts(ctx context.Context) error
		LookupComposeChartByID(ctx context.Context, id uint64) (*composeType.Chart, error)
		LookupComposeChartByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (*composeType.Chart, error)
	}

	ComposeModules interface {
		SearchComposeModules(ctx context.Context, f composeType.ModuleFilter) (composeType.ModuleSet, composeType.ModuleFilter, error)
		CreateComposeModule(ctx context.Context, rr ...*composeType.Module) error
		UpdateComposeModule(ctx context.Context, rr ...*composeType.Module) error
		UpsertComposeModule(ctx context.Context, rr ...*composeType.Module) error
		DeleteComposeModule(ctx context.Context, rr ...*composeType.Module) error
		DeleteComposeModuleByID(ctx context.Context, id uint64) error
		TruncateComposeModules(ctx context.Context) error
		LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (*composeType.Module, error)
		LookupComposeModuleByNamespaceIDName(ctx context.Context, namespaceID uint64, name string) (*composeType.Module, error)
		LookupComposeModuleByID(ctx context.Context, id uint64) (*composeType.Module, error)
	}

	ComposeModuleFields interface {
		SearchComposeModuleFields(ctx context.Context, f composeType.ModuleFieldFilter) (composeType.ModuleFieldSet, composeType.ModuleFieldFilter, error)
		CreateComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) error
		UpdateComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) error
		UpsertComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) error
		DeleteComposeModuleField(ctx context.Context, rr ...*composeType.ModuleField) error
		DeleteComposeModuleFieldByID(ctx context.Context, id uint64) error
		TruncateComposeModuleFields(ctx context.Context) error
		LookupComposeModuleFieldByModuleIDName(ctx context.Context, moduleID uint64, name string) (*composeType.ModuleField, error)
		LookupComposeModuleFieldByID(ctx context.Context, id uint64) (*composeType.ModuleField, error)
	}

	ComposeNamespaces interface {
		SearchComposeNamespaces(ctx context.Context, f composeType.NamespaceFilter) (composeType.NamespaceSet, composeType.NamespaceFilter, error)
		CreateComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) error
		UpdateComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) error
		UpsertComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) error
		DeleteComposeNamespace(ctx context.Context, rr ...*composeType.Namespace) error
		DeleteComposeNamespaceByID(ctx context.Context, id uint64) error
		TruncateComposeNamespaces(ctx context.Context) error
		LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*composeType.Namespace, error)
		LookupComposeNamespaceByID(ctx context.Context, id uint64) (*composeType.Namespace, error)
	}

	ComposePages interface {
		SearchComposePages(ctx context.Context, f composeType.PageFilter) (composeType.PageSet, composeType.PageFilter, error)
		CreateComposePage(ctx context.Context, rr ...*composeType.Page) error
		UpdateComposePage(ctx context.Context, rr ...*composeType.Page) error
		UpsertComposePage(ctx context.Context, rr ...*composeType.Page) error
		DeleteComposePage(ctx context.Context, rr ...*composeType.Page) error
		DeleteComposePageByID(ctx context.Context, id uint64) error
		TruncateComposePages(ctx context.Context) error
		LookupComposePageByNamespaceIDHandle(ctx context.Context, namespaceID uint64, handle string) (*composeType.Page, error)
		LookupComposePageByNamespaceIDModuleID(ctx context.Context, namespaceID uint64, moduleID uint64) (*composeType.Page, error)
		LookupComposePageByID(ctx context.Context, id uint64) (*composeType.Page, error)
		ReorderComposePages(ctx context.Context, namespace_id uint64, parent_id uint64, page_ids []uint64) error
	}

	Credentials interface {
		SearchCredentials(ctx context.Context, f systemType.CredentialFilter) (systemType.CredentialSet, systemType.CredentialFilter, error)
		CreateCredential(ctx context.Context, rr ...*systemType.Credential) error
		UpdateCredential(ctx context.Context, rr ...*systemType.Credential) error
		UpsertCredential(ctx context.Context, rr ...*systemType.Credential) error
		DeleteCredential(ctx context.Context, rr ...*systemType.Credential) error
		DeleteCredentialByID(ctx context.Context, id uint64) error
		TruncateCredentials(ctx context.Context) error
		LookupCredentialByID(ctx context.Context, id uint64) (*systemType.Credential, error)
	}

	DalConnections interface {
		SearchDalConnections(ctx context.Context, f systemType.DalConnectionFilter) (systemType.DalConnectionSet, systemType.DalConnectionFilter, error)
		CreateDalConnection(ctx context.Context, rr ...*systemType.DalConnection) error
		UpdateDalConnection(ctx context.Context, rr ...*systemType.DalConnection) error
		UpsertDalConnection(ctx context.Context, rr ...*systemType.DalConnection) error
		DeleteDalConnection(ctx context.Context, rr ...*systemType.DalConnection) error
		DeleteDalConnectionByID(ctx context.Context, id uint64) error
		TruncateDalConnections(ctx context.Context) error
		LookupDalConnectionByID(ctx context.Context, id uint64) (*systemType.DalConnection, error)
		LookupDalConnectionByHandle(ctx context.Context, handle string) (*systemType.DalConnection, error)
	}

	DalSensitivityLevels interface {
		SearchDalSensitivityLevels(ctx context.Context, f systemType.DalSensitivityLevelFilter) (systemType.DalSensitivityLevelSet, systemType.DalSensitivityLevelFilter, error)
		CreateDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) error
		UpdateDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) error
		UpsertDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) error
		DeleteDalSensitivityLevel(ctx context.Context, rr ...*systemType.DalSensitivityLevel) error
		DeleteDalSensitivityLevelByID(ctx context.Context, id uint64) error
		TruncateDalSensitivityLevels(ctx context.Context) error
		LookupDalSensitivityLevelByID(ctx context.Context, id uint64) (*systemType.DalSensitivityLevel, error)
	}

	DataPrivacyRequests interface {
		SearchDataPrivacyRequests(ctx context.Context, f systemType.DataPrivacyRequestFilter) (systemType.DataPrivacyRequestSet, systemType.DataPrivacyRequestFilter, error)
		CreateDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) error
		UpdateDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) error
		UpsertDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) error
		DeleteDataPrivacyRequest(ctx context.Context, rr ...*systemType.DataPrivacyRequest) error
		DeleteDataPrivacyRequestByID(ctx context.Context, id uint64) error
		TruncateDataPrivacyRequests(ctx context.Context) error
		LookupDataPrivacyRequestByID(ctx context.Context, id uint64) (*systemType.DataPrivacyRequest, error)
	}

	DataPrivacyRequestComments interface {
		SearchDataPrivacyRequestComments(ctx context.Context, f systemType.DataPrivacyRequestCommentFilter) (systemType.DataPrivacyRequestCommentSet, systemType.DataPrivacyRequestCommentFilter, error)
		CreateDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) error
		UpdateDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) error
		UpsertDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) error
		DeleteDataPrivacyRequestComment(ctx context.Context, rr ...*systemType.DataPrivacyRequestComment) error
		DeleteDataPrivacyRequestCommentByID(ctx context.Context, id uint64) error
		TruncateDataPrivacyRequestComments(ctx context.Context) error
	}

	FederationExposedModules interface {
		SearchFederationExposedModules(ctx context.Context, f federationType.ExposedModuleFilter) (federationType.ExposedModuleSet, federationType.ExposedModuleFilter, error)
		CreateFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) error
		UpdateFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) error
		UpsertFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) error
		DeleteFederationExposedModule(ctx context.Context, rr ...*federationType.ExposedModule) error
		DeleteFederationExposedModuleByID(ctx context.Context, id uint64) error
		TruncateFederationExposedModules(ctx context.Context) error
		LookupFederationExposedModuleByID(ctx context.Context, id uint64) (*federationType.ExposedModule, error)
	}

	FederationModuleMappings interface {
		SearchFederationModuleMappings(ctx context.Context, f federationType.ModuleMappingFilter) (federationType.ModuleMappingSet, federationType.ModuleMappingFilter, error)
		CreateFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) error
		UpdateFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) error
		UpsertFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) error
		DeleteFederationModuleMapping(ctx context.Context, rr ...*federationType.ModuleMapping) error
		DeleteFederationModuleMappingByNodeID(ctx context.Context, nodeID uint64) error
		TruncateFederationModuleMappings(ctx context.Context) error
		LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID(ctx context.Context, federationModuleID uint64, composeModuleID uint64, composeNamespaceID uint64) (*federationType.ModuleMapping, error)
		LookupFederationModuleMappingByFederationModuleID(ctx context.Context, federationModuleID uint64) (*federationType.ModuleMapping, error)
	}

	FederationNodes interface {
		SearchFederationNodes(ctx context.Context, f federationType.NodeFilter) (federationType.NodeSet, federationType.NodeFilter, error)
		CreateFederationNode(ctx context.Context, rr ...*federationType.Node) error
		UpdateFederationNode(ctx context.Context, rr ...*federationType.Node) error
		UpsertFederationNode(ctx context.Context, rr ...*federationType.Node) error
		DeleteFederationNode(ctx context.Context, rr ...*federationType.Node) error
		DeleteFederationNodeByID(ctx context.Context, id uint64) error
		TruncateFederationNodes(ctx context.Context) error
		LookupFederationNodeByID(ctx context.Context, id uint64) (*federationType.Node, error)
		LookupFederationNodeByBaseURLSharedNodeID(ctx context.Context, baseURL string, sharedNodeID uint64) (*federationType.Node, error)
		LookupFederationNodeBySharedNodeID(ctx context.Context, sharedNodeID uint64) (*federationType.Node, error)
	}

	FederationNodeSyncs interface {
		SearchFederationNodeSyncs(ctx context.Context, f federationType.NodeSyncFilter) (federationType.NodeSyncSet, federationType.NodeSyncFilter, error)
		CreateFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) error
		UpdateFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) error
		UpsertFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) error
		DeleteFederationNodeSync(ctx context.Context, rr ...*federationType.NodeSync) error
		DeleteFederationNodeSyncByNodeID(ctx context.Context, nodeID uint64) error
		TruncateFederationNodeSyncs(ctx context.Context) error
		LookupFederationNodeSyncByNodeID(ctx context.Context, nodeID uint64) (*federationType.NodeSync, error)
		LookupFederationNodeSyncByNodeIDModuleIDSyncTypeSyncStatus(ctx context.Context, nodeID uint64, moduleID uint64, syncType string, syncStatus string) (*federationType.NodeSync, error)
	}

	FederationSharedModules interface {
		SearchFederationSharedModules(ctx context.Context, f federationType.SharedModuleFilter) (federationType.SharedModuleSet, federationType.SharedModuleFilter, error)
		CreateFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) error
		UpdateFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) error
		UpsertFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) error
		DeleteFederationSharedModule(ctx context.Context, rr ...*federationType.SharedModule) error
		DeleteFederationSharedModuleByID(ctx context.Context, id uint64) error
		TruncateFederationSharedModules(ctx context.Context) error
		LookupFederationSharedModuleByID(ctx context.Context, id uint64) (*federationType.SharedModule, error)
	}

	Flags interface {
		SearchFlags(ctx context.Context, f flagType.FlagFilter) (flagType.FlagSet, flagType.FlagFilter, error)
		CreateFlag(ctx context.Context, rr ...*flagType.Flag) error
		UpdateFlag(ctx context.Context, rr ...*flagType.Flag) error
		UpsertFlag(ctx context.Context, rr ...*flagType.Flag) error
		DeleteFlag(ctx context.Context, rr ...*flagType.Flag) error
		DeleteFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resourceID uint64, ownedBy uint64, name string) error
		TruncateFlags(ctx context.Context) error
		LookupFlagByKindResourceIDOwnedByName(ctx context.Context, kind string, resourceID uint64, ownedBy uint64, name string) (*flagType.Flag, error)
	}

	Labels interface {
		SearchLabels(ctx context.Context, f labelsType.LabelFilter) (labelsType.LabelSet, labelsType.LabelFilter, error)
		CreateLabel(ctx context.Context, rr ...*labelsType.Label) error
		UpdateLabel(ctx context.Context, rr ...*labelsType.Label) error
		UpsertLabel(ctx context.Context, rr ...*labelsType.Label) error
		DeleteLabel(ctx context.Context, rr ...*labelsType.Label) error
		DeleteLabelByKindResourceIDName(ctx context.Context, kind string, resourceID uint64, name string) error
		TruncateLabels(ctx context.Context) error
		LookupLabelByKindResourceIDName(ctx context.Context, kind string, resourceID uint64, name string) (*labelsType.Label, error)
		DeleteExtraLabels(ctx context.Context, kind string, resourceId uint64, name ...string) error
	}

	Queues interface {
		SearchQueues(ctx context.Context, f systemType.QueueFilter) (systemType.QueueSet, systemType.QueueFilter, error)
		CreateQueue(ctx context.Context, rr ...*systemType.Queue) error
		UpdateQueue(ctx context.Context, rr ...*systemType.Queue) error
		UpsertQueue(ctx context.Context, rr ...*systemType.Queue) error
		DeleteQueue(ctx context.Context, rr ...*systemType.Queue) error
		DeleteQueueByID(ctx context.Context, id uint64) error
		TruncateQueues(ctx context.Context) error
		LookupQueueByID(ctx context.Context, id uint64) (*systemType.Queue, error)
		LookupQueueByQueue(ctx context.Context, queue string) (*systemType.Queue, error)
	}

	QueueMessages interface {
		SearchQueueMessages(ctx context.Context, f systemType.QueueMessageFilter) (systemType.QueueMessageSet, systemType.QueueMessageFilter, error)
		CreateQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) error
		UpdateQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) error
		UpsertQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) error
		DeleteQueueMessage(ctx context.Context, rr ...*systemType.QueueMessage) error
		DeleteQueueMessageByID(ctx context.Context, id uint64) error
		TruncateQueueMessages(ctx context.Context) error
	}

	RbacRules interface {
		SearchRbacRules(ctx context.Context, f rbacType.RuleFilter) (rbacType.RuleSet, rbacType.RuleFilter, error)
		CreateRbacRule(ctx context.Context, rr ...*rbacType.Rule) error
		UpdateRbacRule(ctx context.Context, rr ...*rbacType.Rule) error
		UpsertRbacRule(ctx context.Context, rr ...*rbacType.Rule) error
		DeleteRbacRule(ctx context.Context, rr ...*rbacType.Rule) error
		DeleteRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error
		TruncateRbacRules(ctx context.Context) error
		TransferRbacRules(ctx context.Context, src uint64, dst uint64) error
	}

	Reminders interface {
		SearchReminders(ctx context.Context, f systemType.ReminderFilter) (systemType.ReminderSet, systemType.ReminderFilter, error)
		CreateReminder(ctx context.Context, rr ...*systemType.Reminder) error
		UpdateReminder(ctx context.Context, rr ...*systemType.Reminder) error
		UpsertReminder(ctx context.Context, rr ...*systemType.Reminder) error
		DeleteReminder(ctx context.Context, rr ...*systemType.Reminder) error
		DeleteReminderByID(ctx context.Context, id uint64) error
		TruncateReminders(ctx context.Context) error
		LookupReminderByID(ctx context.Context, id uint64) (*systemType.Reminder, error)
	}

	Reports interface {
		SearchReports(ctx context.Context, f systemType.ReportFilter) (systemType.ReportSet, systemType.ReportFilter, error)
		CreateReport(ctx context.Context, rr ...*systemType.Report) error
		UpdateReport(ctx context.Context, rr ...*systemType.Report) error
		UpsertReport(ctx context.Context, rr ...*systemType.Report) error
		DeleteReport(ctx context.Context, rr ...*systemType.Report) error
		DeleteReportByID(ctx context.Context, id uint64) error
		TruncateReports(ctx context.Context) error
		LookupReportByID(ctx context.Context, id uint64) (*systemType.Report, error)
		LookupReportByHandle(ctx context.Context, handle string) (*systemType.Report, error)
	}

	ResourceActivitys interface {
		SearchResourceActivitys(ctx context.Context, f discoveryType.ResourceActivityFilter) (discoveryType.ResourceActivitySet, discoveryType.ResourceActivityFilter, error)
		CreateResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) error
		UpdateResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) error
		UpsertResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) error
		DeleteResourceActivity(ctx context.Context, rr ...*discoveryType.ResourceActivity) error
		DeleteResourceActivityByID(ctx context.Context, id uint64) error
		TruncateResourceActivitys(ctx context.Context) error
	}

	ResourceTranslations interface {
		SearchResourceTranslations(ctx context.Context, f systemType.ResourceTranslationFilter) (systemType.ResourceTranslationSet, systemType.ResourceTranslationFilter, error)
		CreateResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) error
		UpdateResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) error
		UpsertResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) error
		DeleteResourceTranslation(ctx context.Context, rr ...*systemType.ResourceTranslation) error
		DeleteResourceTranslationByID(ctx context.Context, id uint64) error
		TruncateResourceTranslations(ctx context.Context) error
		LookupResourceTranslationByID(ctx context.Context, id uint64) (*systemType.ResourceTranslation, error)
		TransformResource(ctx context.Context, lang language.Tag) (map[string]map[string]*locale.ResourceTranslation, error)
	}

	Roles interface {
		SearchRoles(ctx context.Context, f systemType.RoleFilter) (systemType.RoleSet, systemType.RoleFilter, error)
		CreateRole(ctx context.Context, rr ...*systemType.Role) error
		UpdateRole(ctx context.Context, rr ...*systemType.Role) error
		UpsertRole(ctx context.Context, rr ...*systemType.Role) error
		DeleteRole(ctx context.Context, rr ...*systemType.Role) error
		DeleteRoleByID(ctx context.Context, id uint64) error
		TruncateRoles(ctx context.Context) error
		LookupRoleByID(ctx context.Context, id uint64) (*systemType.Role, error)
		LookupRoleByHandle(ctx context.Context, handle string) (*systemType.Role, error)
		LookupRoleByName(ctx context.Context, name string) (*systemType.Role, error)
		RoleMetrics(ctx context.Context) (*systemType.RoleMetrics, error)
	}

	RoleMembers interface {
		SearchRoleMembers(ctx context.Context, f systemType.RoleMemberFilter) (systemType.RoleMemberSet, systemType.RoleMemberFilter, error)
		CreateRoleMember(ctx context.Context, rr ...*systemType.RoleMember) error
		UpdateRoleMember(ctx context.Context, rr ...*systemType.RoleMember) error
		UpsertRoleMember(ctx context.Context, rr ...*systemType.RoleMember) error
		DeleteRoleMember(ctx context.Context, rr ...*systemType.RoleMember) error
		DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error
		TruncateRoleMembers(ctx context.Context) error
		TransferRoleMembers(ctx context.Context, src uint64, dst uint64) error
	}

	SettingValues interface {
		SearchSettingValues(ctx context.Context, f systemType.SettingsFilter) (systemType.SettingValueSet, systemType.SettingsFilter, error)
		CreateSettingValue(ctx context.Context, rr ...*systemType.SettingValue) error
		UpdateSettingValue(ctx context.Context, rr ...*systemType.SettingValue) error
		UpsertSettingValue(ctx context.Context, rr ...*systemType.SettingValue) error
		DeleteSettingValue(ctx context.Context, rr ...*systemType.SettingValue) error
		DeleteSettingValueByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) error
		TruncateSettingValues(ctx context.Context) error
		LookupSettingValueByNameOwnedBy(ctx context.Context, name string, ownedBy uint64) (*systemType.SettingValue, error)
	}

	Templates interface {
		SearchTemplates(ctx context.Context, f systemType.TemplateFilter) (systemType.TemplateSet, systemType.TemplateFilter, error)
		CreateTemplate(ctx context.Context, rr ...*systemType.Template) error
		UpdateTemplate(ctx context.Context, rr ...*systemType.Template) error
		UpsertTemplate(ctx context.Context, rr ...*systemType.Template) error
		DeleteTemplate(ctx context.Context, rr ...*systemType.Template) error
		DeleteTemplateByID(ctx context.Context, id uint64) error
		TruncateTemplates(ctx context.Context) error
		LookupTemplateByID(ctx context.Context, id uint64) (*systemType.Template, error)
		LookupTemplateByHandle(ctx context.Context, handle string) (*systemType.Template, error)
	}

	Users interface {
		SearchUsers(ctx context.Context, f systemType.UserFilter) (systemType.UserSet, systemType.UserFilter, error)
		CreateUser(ctx context.Context, rr ...*systemType.User) error
		UpdateUser(ctx context.Context, rr ...*systemType.User) error
		UpsertUser(ctx context.Context, rr ...*systemType.User) error
		DeleteUser(ctx context.Context, rr ...*systemType.User) error
		DeleteUserByID(ctx context.Context, id uint64) error
		TruncateUsers(ctx context.Context) error
		LookupUserByID(ctx context.Context, id uint64) (*systemType.User, error)
		LookupUserByEmail(ctx context.Context, email string) (*systemType.User, error)
		LookupUserByHandle(ctx context.Context, handle string) (*systemType.User, error)
		LookupUserByUsername(ctx context.Context, username string) (*systemType.User, error)
		CountUsers(ctx context.Context, u systemType.UserFilter) (uint, error)
		UserMetrics(ctx context.Context) (*systemType.UserMetrics, error)
	}
)

// SearchActionlogs returns all matching Actionlogs from store
//
// This function is auto-generated
func SearchActionlogs(ctx context.Context, s Actionlogs, f actionlogType.Filter) (actionlogType.ActionSet, actionlogType.Filter, error) {
	return s.SearchActionlogs(ctx, f)
}

// CreateActionlog creates one or more Actionlogs in store
//
// This function is auto-generated
func CreateActionlog(ctx context.Context, s Actionlogs, rr ...*actionlogType.Action) error {
	return s.CreateActionlog(ctx, rr...)
}

// UpdateActionlog updates one or more (existing) Actionlogs in store
//
// This function is auto-generated
func UpdateActionlog(ctx context.Context, s Actionlogs, rr ...*actionlogType.Action) error {
	return s.UpdateActionlog(ctx, rr...)
}

// UpsertActionlog creates new or updates existing one or more Actionlogs in store
//
// This function is auto-generated
func UpsertActionlog(ctx context.Context, s Actionlogs, rr ...*actionlogType.Action) error {
	return s.UpsertActionlog(ctx, rr...)
}

// DeleteActionlog deletes one or more Actionlogs from store
//
// This function is auto-generated
func DeleteActionlog(ctx context.Context, s Actionlogs, rr ...*actionlogType.Action) error {
	return s.DeleteActionlog(ctx, rr...)
}

// DeleteActionlogByID deletes one or more Actionlogs from store
//
// This function is auto-generated
func DeleteActionlogByID(ctx context.Context, s Actionlogs, id uint64) error {
	return s.DeleteActionlogByID(ctx, id)
}

// TruncateActionlogs Deletes all Actionlogs from store
//
// This function is auto-generated
func TruncateActionlogs(ctx context.Context, s Actionlogs) error {
	return s.TruncateActionlogs(ctx)
}

// LookupActionlogByID searches for action log by ID
//
// This function is auto-generated
func LookupActionlogByID(ctx context.Context, s Actionlogs, id uint64) (*actionlogType.Action, error) {
	return s.LookupActionlogByID(ctx, id)
}

// SearchApigwFilters returns all matching ApigwFilters from store
//
// This function is auto-generated
func SearchApigwFilters(ctx context.Context, s ApigwFilters, f systemType.ApigwFilterFilter) (systemType.ApigwFilterSet, systemType.ApigwFilterFilter, error) {
	return s.SearchApigwFilters(ctx, f)
}

// CreateApigwFilter creates one or more ApigwFilters in store
//
// This function is auto-generated
func CreateApigwFilter(ctx context.Context, s ApigwFilters, rr ...*systemType.ApigwFilter) error {
	return s.CreateApigwFilter(ctx, rr...)
}

// UpdateApigwFilter updates one or more (existing) ApigwFilters in store
//
// This function is auto-generated
func UpdateApigwFilter(ctx context.Context, s ApigwFilters, rr ...*systemType.ApigwFilter) error {
	return s.UpdateApigwFilter(ctx, rr...)
}

// UpsertApigwFilter creates new or updates existing one or more ApigwFilters in store
//
// This function is auto-generated
func UpsertApigwFilter(ctx context.Context, s ApigwFilters, rr ...*systemType.ApigwFilter) error {
	return s.UpsertApigwFilter(ctx, rr...)
}

// DeleteApigwFilter deletes one or more ApigwFilters from store
//
// This function is auto-generated
func DeleteApigwFilter(ctx context.Context, s ApigwFilters, rr ...*systemType.ApigwFilter) error {
	return s.DeleteApigwFilter(ctx, rr...)
}

// DeleteApigwFilterByID deletes one or more ApigwFilters from store
//
// This function is auto-generated
func DeleteApigwFilterByID(ctx context.Context, s ApigwFilters, id uint64) error {
	return s.DeleteApigwFilterByID(ctx, id)
}

// TruncateApigwFilters Deletes all ApigwFilters from store
//
// This function is auto-generated
func TruncateApigwFilters(ctx context.Context, s ApigwFilters) error {
	return s.TruncateApigwFilters(ctx)
}

// LookupApigwFilterByID searches for filter by ID
//
// This function is auto-generated
func LookupApigwFilterByID(ctx context.Context, s ApigwFilters, id uint64) (*systemType.ApigwFilter, error) {
	return s.LookupApigwFilterByID(ctx, id)
}

// LookupApigwFilterByRoute searches for filter by route
//
// This function is auto-generated
func LookupApigwFilterByRoute(ctx context.Context, s ApigwFilters, route uint64) (*systemType.ApigwFilter, error) {
	return s.LookupApigwFilterByRoute(ctx, route)
}

// SearchApigwRoutes returns all matching ApigwRoutes from store
//
// This function is auto-generated
func SearchApigwRoutes(ctx context.Context, s ApigwRoutes, f systemType.ApigwRouteFilter) (systemType.ApigwRouteSet, systemType.ApigwRouteFilter, error) {
	return s.SearchApigwRoutes(ctx, f)
}

// CreateApigwRoute creates one or more ApigwRoutes in store
//
// This function is auto-generated
func CreateApigwRoute(ctx context.Context, s ApigwRoutes, rr ...*systemType.ApigwRoute) error {
	return s.CreateApigwRoute(ctx, rr...)
}

// UpdateApigwRoute updates one or more (existing) ApigwRoutes in store
//
// This function is auto-generated
func UpdateApigwRoute(ctx context.Context, s ApigwRoutes, rr ...*systemType.ApigwRoute) error {
	return s.UpdateApigwRoute(ctx, rr...)
}

// UpsertApigwRoute creates new or updates existing one or more ApigwRoutes in store
//
// This function is auto-generated
func UpsertApigwRoute(ctx context.Context, s ApigwRoutes, rr ...*systemType.ApigwRoute) error {
	return s.UpsertApigwRoute(ctx, rr...)
}

// DeleteApigwRoute deletes one or more ApigwRoutes from store
//
// This function is auto-generated
func DeleteApigwRoute(ctx context.Context, s ApigwRoutes, rr ...*systemType.ApigwRoute) error {
	return s.DeleteApigwRoute(ctx, rr...)
}

// DeleteApigwRouteByID deletes one or more ApigwRoutes from store
//
// This function is auto-generated
func DeleteApigwRouteByID(ctx context.Context, s ApigwRoutes, id uint64) error {
	return s.DeleteApigwRouteByID(ctx, id)
}

// TruncateApigwRoutes Deletes all ApigwRoutes from store
//
// This function is auto-generated
func TruncateApigwRoutes(ctx context.Context, s ApigwRoutes) error {
	return s.TruncateApigwRoutes(ctx)
}

// LookupApigwRouteByID searches for route by ID
//
// It returns route even if deleted or suspended
//
// This function is auto-generated
func LookupApigwRouteByID(ctx context.Context, s ApigwRoutes, id uint64) (*systemType.ApigwRoute, error) {
	return s.LookupApigwRouteByID(ctx, id)
}

// LookupApigwRouteByEndpoint searches for route by endpoint
//
// It returns route even if deleted or suspended
//
// This function is auto-generated
func LookupApigwRouteByEndpoint(ctx context.Context, s ApigwRoutes, endpoint string) (*systemType.ApigwRoute, error) {
	return s.LookupApigwRouteByEndpoint(ctx, endpoint)
}

// SearchApplications returns all matching Applications from store
//
// This function is auto-generated
func SearchApplications(ctx context.Context, s Applications, f systemType.ApplicationFilter) (systemType.ApplicationSet, systemType.ApplicationFilter, error) {
	return s.SearchApplications(ctx, f)
}

// CreateApplication creates one or more Applications in store
//
// This function is auto-generated
func CreateApplication(ctx context.Context, s Applications, rr ...*systemType.Application) error {
	return s.CreateApplication(ctx, rr...)
}

// UpdateApplication updates one or more (existing) Applications in store
//
// This function is auto-generated
func UpdateApplication(ctx context.Context, s Applications, rr ...*systemType.Application) error {
	return s.UpdateApplication(ctx, rr...)
}

// UpsertApplication creates new or updates existing one or more Applications in store
//
// This function is auto-generated
func UpsertApplication(ctx context.Context, s Applications, rr ...*systemType.Application) error {
	return s.UpsertApplication(ctx, rr...)
}

// DeleteApplication deletes one or more Applications from store
//
// This function is auto-generated
func DeleteApplication(ctx context.Context, s Applications, rr ...*systemType.Application) error {
	return s.DeleteApplication(ctx, rr...)
}

// DeleteApplicationByID deletes one or more Applications from store
//
// This function is auto-generated
func DeleteApplicationByID(ctx context.Context, s Applications, id uint64) error {
	return s.DeleteApplicationByID(ctx, id)
}

// TruncateApplications Deletes all Applications from store
//
// This function is auto-generated
func TruncateApplications(ctx context.Context, s Applications) error {
	return s.TruncateApplications(ctx)
}

// LookupApplicationByID searches for role by ID
//
// It returns role even if deleted or suspended
//
// This function is auto-generated
func LookupApplicationByID(ctx context.Context, s Applications, id uint64) (*systemType.Application, error) {
	return s.LookupApplicationByID(ctx, id)
}

// ApplicationMetrics
//
// This function is auto-generated
func ApplicationMetrics(ctx context.Context, s Applications) (*systemType.ApplicationMetrics, error) {
	return s.ApplicationMetrics(ctx)
}

// ReorderApplications
//
// This function is auto-generated
func ReorderApplications(ctx context.Context, s Applications, order []uint64) error {
	return s.ReorderApplications(ctx, order)
}

// SearchAttachments returns all matching Attachments from store
//
// This function is auto-generated
func SearchAttachments(ctx context.Context, s Attachments, f systemType.AttachmentFilter) (systemType.AttachmentSet, systemType.AttachmentFilter, error) {
	return s.SearchAttachments(ctx, f)
}

// CreateAttachment creates one or more Attachments in store
//
// This function is auto-generated
func CreateAttachment(ctx context.Context, s Attachments, rr ...*systemType.Attachment) error {
	return s.CreateAttachment(ctx, rr...)
}

// UpdateAttachment updates one or more (existing) Attachments in store
//
// This function is auto-generated
func UpdateAttachment(ctx context.Context, s Attachments, rr ...*systemType.Attachment) error {
	return s.UpdateAttachment(ctx, rr...)
}

// UpsertAttachment creates new or updates existing one or more Attachments in store
//
// This function is auto-generated
func UpsertAttachment(ctx context.Context, s Attachments, rr ...*systemType.Attachment) error {
	return s.UpsertAttachment(ctx, rr...)
}

// DeleteAttachment deletes one or more Attachments from store
//
// This function is auto-generated
func DeleteAttachment(ctx context.Context, s Attachments, rr ...*systemType.Attachment) error {
	return s.DeleteAttachment(ctx, rr...)
}

// DeleteAttachmentByID deletes one or more Attachments from store
//
// This function is auto-generated
func DeleteAttachmentByID(ctx context.Context, s Attachments, id uint64) error {
	return s.DeleteAttachmentByID(ctx, id)
}

// TruncateAttachments Deletes all Attachments from store
//
// This function is auto-generated
func TruncateAttachments(ctx context.Context, s Attachments) error {
	return s.TruncateAttachments(ctx)
}

// LookupAttachmentByID
//
// This function is auto-generated
func LookupAttachmentByID(ctx context.Context, s Attachments, id uint64) (*systemType.Attachment, error) {
	return s.LookupAttachmentByID(ctx, id)
}

// SearchAuthClients returns all matching AuthClients from store
//
// This function is auto-generated
func SearchAuthClients(ctx context.Context, s AuthClients, f systemType.AuthClientFilter) (systemType.AuthClientSet, systemType.AuthClientFilter, error) {
	return s.SearchAuthClients(ctx, f)
}

// CreateAuthClient creates one or more AuthClients in store
//
// This function is auto-generated
func CreateAuthClient(ctx context.Context, s AuthClients, rr ...*systemType.AuthClient) error {
	return s.CreateAuthClient(ctx, rr...)
}

// UpdateAuthClient updates one or more (existing) AuthClients in store
//
// This function is auto-generated
func UpdateAuthClient(ctx context.Context, s AuthClients, rr ...*systemType.AuthClient) error {
	return s.UpdateAuthClient(ctx, rr...)
}

// UpsertAuthClient creates new or updates existing one or more AuthClients in store
//
// This function is auto-generated
func UpsertAuthClient(ctx context.Context, s AuthClients, rr ...*systemType.AuthClient) error {
	return s.UpsertAuthClient(ctx, rr...)
}

// DeleteAuthClient deletes one or more AuthClients from store
//
// This function is auto-generated
func DeleteAuthClient(ctx context.Context, s AuthClients, rr ...*systemType.AuthClient) error {
	return s.DeleteAuthClient(ctx, rr...)
}

// DeleteAuthClientByID deletes one or more AuthClients from store
//
// This function is auto-generated
func DeleteAuthClientByID(ctx context.Context, s AuthClients, id uint64) error {
	return s.DeleteAuthClientByID(ctx, id)
}

// TruncateAuthClients Deletes all AuthClients from store
//
// This function is auto-generated
func TruncateAuthClients(ctx context.Context, s AuthClients) error {
	return s.TruncateAuthClients(ctx)
}

// LookupAuthClientByID 	searches for auth client by ID
//
// 	It returns auth clint even if deleted
//
// This function is auto-generated
func LookupAuthClientByID(ctx context.Context, s AuthClients, id uint64) (*systemType.AuthClient, error) {
	return s.LookupAuthClientByID(ctx, id)
}

// LookupAuthClientByHandle searches for auth client by ID
//
// It returns auth clint even if deleted
//
// This function is auto-generated
func LookupAuthClientByHandle(ctx context.Context, s AuthClients, handle string) (*systemType.AuthClient, error) {
	return s.LookupAuthClientByHandle(ctx, handle)
}

// SearchAuthConfirmedClients returns all matching AuthConfirmedClients from store
//
// This function is auto-generated
func SearchAuthConfirmedClients(ctx context.Context, s AuthConfirmedClients, f systemType.AuthConfirmedClientFilter) (systemType.AuthConfirmedClientSet, systemType.AuthConfirmedClientFilter, error) {
	return s.SearchAuthConfirmedClients(ctx, f)
}

// CreateAuthConfirmedClient creates one or more AuthConfirmedClients in store
//
// This function is auto-generated
func CreateAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*systemType.AuthConfirmedClient) error {
	return s.CreateAuthConfirmedClient(ctx, rr...)
}

// UpdateAuthConfirmedClient updates one or more (existing) AuthConfirmedClients in store
//
// This function is auto-generated
func UpdateAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*systemType.AuthConfirmedClient) error {
	return s.UpdateAuthConfirmedClient(ctx, rr...)
}

// UpsertAuthConfirmedClient creates new or updates existing one or more AuthConfirmedClients in store
//
// This function is auto-generated
func UpsertAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*systemType.AuthConfirmedClient) error {
	return s.UpsertAuthConfirmedClient(ctx, rr...)
}

// DeleteAuthConfirmedClient deletes one or more AuthConfirmedClients from store
//
// This function is auto-generated
func DeleteAuthConfirmedClient(ctx context.Context, s AuthConfirmedClients, rr ...*systemType.AuthConfirmedClient) error {
	return s.DeleteAuthConfirmedClient(ctx, rr...)
}

// DeleteAuthConfirmedClientByID deletes one or more AuthConfirmedClients from store
//
// This function is auto-generated
func DeleteAuthConfirmedClientByUserIDClientID(ctx context.Context, s AuthConfirmedClients, userID uint64, clientID uint64) error {
	return s.DeleteAuthConfirmedClientByUserIDClientID(ctx, userID, clientID)
}

// TruncateAuthConfirmedClients Deletes all AuthConfirmedClients from store
//
// This function is auto-generated
func TruncateAuthConfirmedClients(ctx context.Context, s AuthConfirmedClients) error {
	return s.TruncateAuthConfirmedClients(ctx)
}

// LookupAuthConfirmedClientByUserIDClientID
//
// This function is auto-generated
func LookupAuthConfirmedClientByUserIDClientID(ctx context.Context, s AuthConfirmedClients, userID uint64, clientID uint64) (*systemType.AuthConfirmedClient, error) {
	return s.LookupAuthConfirmedClientByUserIDClientID(ctx, userID, clientID)
}

// SearchAuthOa2tokens returns all matching AuthOa2tokens from store
//
// This function is auto-generated
func SearchAuthOa2tokens(ctx context.Context, s AuthOa2tokens, f systemType.AuthOa2tokenFilter) (systemType.AuthOa2tokenSet, systemType.AuthOa2tokenFilter, error) {
	return s.SearchAuthOa2tokens(ctx, f)
}

// CreateAuthOa2token creates one or more AuthOa2tokens in store
//
// This function is auto-generated
func CreateAuthOa2token(ctx context.Context, s AuthOa2tokens, rr ...*systemType.AuthOa2token) error {
	return s.CreateAuthOa2token(ctx, rr...)
}

// UpdateAuthOa2token updates one or more (existing) AuthOa2tokens in store
//
// This function is auto-generated
func UpdateAuthOa2token(ctx context.Context, s AuthOa2tokens, rr ...*systemType.AuthOa2token) error {
	return s.UpdateAuthOa2token(ctx, rr...)
}

// UpsertAuthOa2token creates new or updates existing one or more AuthOa2tokens in store
//
// This function is auto-generated
func UpsertAuthOa2token(ctx context.Context, s AuthOa2tokens, rr ...*systemType.AuthOa2token) error {
	return s.UpsertAuthOa2token(ctx, rr...)
}

// DeleteAuthOa2token deletes one or more AuthOa2tokens from store
//
// This function is auto-generated
func DeleteAuthOa2token(ctx context.Context, s AuthOa2tokens, rr ...*systemType.AuthOa2token) error {
	return s.DeleteAuthOa2token(ctx, rr...)
}

// DeleteAuthOa2tokenByID deletes one or more AuthOa2tokens from store
//
// This function is auto-generated
func DeleteAuthOa2tokenByID(ctx context.Context, s AuthOa2tokens, id uint64) error {
	return s.DeleteAuthOa2tokenByID(ctx, id)
}

// TruncateAuthOa2tokens Deletes all AuthOa2tokens from store
//
// This function is auto-generated
func TruncateAuthOa2tokens(ctx context.Context, s AuthOa2tokens) error {
	return s.TruncateAuthOa2tokens(ctx)
}

// LookupAuthOa2tokenByID
//
// This function is auto-generated
func LookupAuthOa2tokenByID(ctx context.Context, s AuthOa2tokens, id uint64) (*systemType.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByID(ctx, id)
}

// LookupAuthOa2tokenByCode
//
// This function is auto-generated
func LookupAuthOa2tokenByCode(ctx context.Context, s AuthOa2tokens, code string) (*systemType.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByCode(ctx, code)
}

// LookupAuthOa2tokenByAccess
//
// This function is auto-generated
func LookupAuthOa2tokenByAccess(ctx context.Context, s AuthOa2tokens, access string) (*systemType.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByAccess(ctx, access)
}

// LookupAuthOa2tokenByRefresh
//
// This function is auto-generated
func LookupAuthOa2tokenByRefresh(ctx context.Context, s AuthOa2tokens, refresh string) (*systemType.AuthOa2token, error) {
	return s.LookupAuthOa2tokenByRefresh(ctx, refresh)
}

// DeleteExpiredAuthOA2Tokens
//
// This function is auto-generated
func DeleteExpiredAuthOA2Tokens(ctx context.Context, s AuthOa2tokens) error {
	return s.DeleteExpiredAuthOA2Tokens(ctx)
}

// DeleteAuthOA2TokenByCode
//
// This function is auto-generated
func DeleteAuthOA2TokenByCode(ctx context.Context, s AuthOa2tokens, code string) error {
	return s.DeleteAuthOA2TokenByCode(ctx, code)
}

// DeleteAuthOA2TokenByAccess
//
// This function is auto-generated
func DeleteAuthOA2TokenByAccess(ctx context.Context, s AuthOa2tokens, access string) error {
	return s.DeleteAuthOA2TokenByAccess(ctx, access)
}

// DeleteAuthOA2TokenByRefresh
//
// This function is auto-generated
func DeleteAuthOA2TokenByRefresh(ctx context.Context, s AuthOa2tokens, refresh string) error {
	return s.DeleteAuthOA2TokenByRefresh(ctx, refresh)
}

// DeleteAuthOA2TokenByUserID
//
// This function is auto-generated
func DeleteAuthOA2TokenByUserID(ctx context.Context, s AuthOa2tokens, userID uint64) error {
	return s.DeleteAuthOA2TokenByUserID(ctx, userID)
}

// SearchAuthSessions returns all matching AuthSessions from store
//
// This function is auto-generated
func SearchAuthSessions(ctx context.Context, s AuthSessions, f systemType.AuthSessionFilter) (systemType.AuthSessionSet, systemType.AuthSessionFilter, error) {
	return s.SearchAuthSessions(ctx, f)
}

// CreateAuthSession creates one or more AuthSessions in store
//
// This function is auto-generated
func CreateAuthSession(ctx context.Context, s AuthSessions, rr ...*systemType.AuthSession) error {
	return s.CreateAuthSession(ctx, rr...)
}

// UpdateAuthSession updates one or more (existing) AuthSessions in store
//
// This function is auto-generated
func UpdateAuthSession(ctx context.Context, s AuthSessions, rr ...*systemType.AuthSession) error {
	return s.UpdateAuthSession(ctx, rr...)
}

// UpsertAuthSession creates new or updates existing one or more AuthSessions in store
//
// This function is auto-generated
func UpsertAuthSession(ctx context.Context, s AuthSessions, rr ...*systemType.AuthSession) error {
	return s.UpsertAuthSession(ctx, rr...)
}

// DeleteAuthSession deletes one or more AuthSessions from store
//
// This function is auto-generated
func DeleteAuthSession(ctx context.Context, s AuthSessions, rr ...*systemType.AuthSession) error {
	return s.DeleteAuthSession(ctx, rr...)
}

// DeleteAuthSessionByID deletes one or more AuthSessions from store
//
// This function is auto-generated
func DeleteAuthSessionByID(ctx context.Context, s AuthSessions, id string) error {
	return s.DeleteAuthSessionByID(ctx, id)
}

// TruncateAuthSessions Deletes all AuthSessions from store
//
// This function is auto-generated
func TruncateAuthSessions(ctx context.Context, s AuthSessions) error {
	return s.TruncateAuthSessions(ctx)
}

// LookupAuthSessionByID
//
// This function is auto-generated
func LookupAuthSessionByID(ctx context.Context, s AuthSessions, id string) (*systemType.AuthSession, error) {
	return s.LookupAuthSessionByID(ctx, id)
}

// DeleteExpiredAuthSessions
//
// This function is auto-generated
func DeleteExpiredAuthSessions(ctx context.Context, s AuthSessions) error {
	return s.DeleteExpiredAuthSessions(ctx)
}

// DeleteAuthSessionsByUserID
//
// This function is auto-generated
func DeleteAuthSessionsByUserID(ctx context.Context, s AuthSessions, userID uint64) error {
	return s.DeleteAuthSessionsByUserID(ctx, userID)
}

// SearchAutomationSessions returns all matching AutomationSessions from store
//
// This function is auto-generated
func SearchAutomationSessions(ctx context.Context, s AutomationSessions, f automationType.SessionFilter) (automationType.SessionSet, automationType.SessionFilter, error) {
	return s.SearchAutomationSessions(ctx, f)
}

// CreateAutomationSession creates one or more AutomationSessions in store
//
// This function is auto-generated
func CreateAutomationSession(ctx context.Context, s AutomationSessions, rr ...*automationType.Session) error {
	return s.CreateAutomationSession(ctx, rr...)
}

// UpdateAutomationSession updates one or more (existing) AutomationSessions in store
//
// This function is auto-generated
func UpdateAutomationSession(ctx context.Context, s AutomationSessions, rr ...*automationType.Session) error {
	return s.UpdateAutomationSession(ctx, rr...)
}

// UpsertAutomationSession creates new or updates existing one or more AutomationSessions in store
//
// This function is auto-generated
func UpsertAutomationSession(ctx context.Context, s AutomationSessions, rr ...*automationType.Session) error {
	return s.UpsertAutomationSession(ctx, rr...)
}

// DeleteAutomationSession deletes one or more AutomationSessions from store
//
// This function is auto-generated
func DeleteAutomationSession(ctx context.Context, s AutomationSessions, rr ...*automationType.Session) error {
	return s.DeleteAutomationSession(ctx, rr...)
}

// DeleteAutomationSessionByID deletes one or more AutomationSessions from store
//
// This function is auto-generated
func DeleteAutomationSessionByID(ctx context.Context, s AutomationSessions, id uint64) error {
	return s.DeleteAutomationSessionByID(ctx, id)
}

// TruncateAutomationSessions Deletes all AutomationSessions from store
//
// This function is auto-generated
func TruncateAutomationSessions(ctx context.Context, s AutomationSessions) error {
	return s.TruncateAutomationSessions(ctx)
}

// LookupAutomationSessionByID searches for session by ID
//
// It returns session even if deleted
//
// This function is auto-generated
func LookupAutomationSessionByID(ctx context.Context, s AutomationSessions, id uint64) (*automationType.Session, error) {
	return s.LookupAutomationSessionByID(ctx, id)
}

// SearchAutomationTriggers returns all matching AutomationTriggers from store
//
// This function is auto-generated
func SearchAutomationTriggers(ctx context.Context, s AutomationTriggers, f automationType.TriggerFilter) (automationType.TriggerSet, automationType.TriggerFilter, error) {
	return s.SearchAutomationTriggers(ctx, f)
}

// CreateAutomationTrigger creates one or more AutomationTriggers in store
//
// This function is auto-generated
func CreateAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*automationType.Trigger) error {
	return s.CreateAutomationTrigger(ctx, rr...)
}

// UpdateAutomationTrigger updates one or more (existing) AutomationTriggers in store
//
// This function is auto-generated
func UpdateAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*automationType.Trigger) error {
	return s.UpdateAutomationTrigger(ctx, rr...)
}

// UpsertAutomationTrigger creates new or updates existing one or more AutomationTriggers in store
//
// This function is auto-generated
func UpsertAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*automationType.Trigger) error {
	return s.UpsertAutomationTrigger(ctx, rr...)
}

// DeleteAutomationTrigger deletes one or more AutomationTriggers from store
//
// This function is auto-generated
func DeleteAutomationTrigger(ctx context.Context, s AutomationTriggers, rr ...*automationType.Trigger) error {
	return s.DeleteAutomationTrigger(ctx, rr...)
}

// DeleteAutomationTriggerByID deletes one or more AutomationTriggers from store
//
// This function is auto-generated
func DeleteAutomationTriggerByID(ctx context.Context, s AutomationTriggers, id uint64) error {
	return s.DeleteAutomationTriggerByID(ctx, id)
}

// TruncateAutomationTriggers Deletes all AutomationTriggers from store
//
// This function is auto-generated
func TruncateAutomationTriggers(ctx context.Context, s AutomationTriggers) error {
	return s.TruncateAutomationTriggers(ctx)
}

// LookupAutomationTriggerByID searches for trigger by ID
//
// It returns trigger even if deleted
//
// This function is auto-generated
func LookupAutomationTriggerByID(ctx context.Context, s AutomationTriggers, id uint64) (*automationType.Trigger, error) {
	return s.LookupAutomationTriggerByID(ctx, id)
}

// SearchAutomationWorkflows returns all matching AutomationWorkflows from store
//
// This function is auto-generated
func SearchAutomationWorkflows(ctx context.Context, s AutomationWorkflows, f automationType.WorkflowFilter) (automationType.WorkflowSet, automationType.WorkflowFilter, error) {
	return s.SearchAutomationWorkflows(ctx, f)
}

// CreateAutomationWorkflow creates one or more AutomationWorkflows in store
//
// This function is auto-generated
func CreateAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*automationType.Workflow) error {
	return s.CreateAutomationWorkflow(ctx, rr...)
}

// UpdateAutomationWorkflow updates one or more (existing) AutomationWorkflows in store
//
// This function is auto-generated
func UpdateAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*automationType.Workflow) error {
	return s.UpdateAutomationWorkflow(ctx, rr...)
}

// UpsertAutomationWorkflow creates new or updates existing one or more AutomationWorkflows in store
//
// This function is auto-generated
func UpsertAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*automationType.Workflow) error {
	return s.UpsertAutomationWorkflow(ctx, rr...)
}

// DeleteAutomationWorkflow deletes one or more AutomationWorkflows from store
//
// This function is auto-generated
func DeleteAutomationWorkflow(ctx context.Context, s AutomationWorkflows, rr ...*automationType.Workflow) error {
	return s.DeleteAutomationWorkflow(ctx, rr...)
}

// DeleteAutomationWorkflowByID deletes one or more AutomationWorkflows from store
//
// This function is auto-generated
func DeleteAutomationWorkflowByID(ctx context.Context, s AutomationWorkflows, id uint64) error {
	return s.DeleteAutomationWorkflowByID(ctx, id)
}

// TruncateAutomationWorkflows Deletes all AutomationWorkflows from store
//
// This function is auto-generated
func TruncateAutomationWorkflows(ctx context.Context, s AutomationWorkflows) error {
	return s.TruncateAutomationWorkflows(ctx)
}

// LookupAutomationWorkflowByID searches for workflow by ID
//
// It returns workflow even if deleted
//
// This function is auto-generated
func LookupAutomationWorkflowByID(ctx context.Context, s AutomationWorkflows, id uint64) (*automationType.Workflow, error) {
	return s.LookupAutomationWorkflowByID(ctx, id)
}

// LookupAutomationWorkflowByHandle searches for workflow by their handle
//
// It returns only valid workflows
//
// This function is auto-generated
func LookupAutomationWorkflowByHandle(ctx context.Context, s AutomationWorkflows, handle string) (*automationType.Workflow, error) {
	return s.LookupAutomationWorkflowByHandle(ctx, handle)
}

// SearchComposeAttachments returns all matching ComposeAttachments from store
//
// This function is auto-generated
func SearchComposeAttachments(ctx context.Context, s ComposeAttachments, f composeType.AttachmentFilter) (composeType.AttachmentSet, composeType.AttachmentFilter, error) {
	return s.SearchComposeAttachments(ctx, f)
}

// CreateComposeAttachment creates one or more ComposeAttachments in store
//
// This function is auto-generated
func CreateComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*composeType.Attachment) error {
	return s.CreateComposeAttachment(ctx, rr...)
}

// UpdateComposeAttachment updates one or more (existing) ComposeAttachments in store
//
// This function is auto-generated
func UpdateComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*composeType.Attachment) error {
	return s.UpdateComposeAttachment(ctx, rr...)
}

// UpsertComposeAttachment creates new or updates existing one or more ComposeAttachments in store
//
// This function is auto-generated
func UpsertComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*composeType.Attachment) error {
	return s.UpsertComposeAttachment(ctx, rr...)
}

// DeleteComposeAttachment deletes one or more ComposeAttachments from store
//
// This function is auto-generated
func DeleteComposeAttachment(ctx context.Context, s ComposeAttachments, rr ...*composeType.Attachment) error {
	return s.DeleteComposeAttachment(ctx, rr...)
}

// DeleteComposeAttachmentByID deletes one or more ComposeAttachments from store
//
// This function is auto-generated
func DeleteComposeAttachmentByID(ctx context.Context, s ComposeAttachments, id uint64) error {
	return s.DeleteComposeAttachmentByID(ctx, id)
}

// TruncateComposeAttachments Deletes all ComposeAttachments from store
//
// This function is auto-generated
func TruncateComposeAttachments(ctx context.Context, s ComposeAttachments) error {
	return s.TruncateComposeAttachments(ctx)
}

// LookupComposeAttachmentByID
//
// This function is auto-generated
func LookupComposeAttachmentByID(ctx context.Context, s ComposeAttachments, id uint64) (*composeType.Attachment, error) {
	return s.LookupComposeAttachmentByID(ctx, id)
}

// SearchComposeCharts returns all matching ComposeCharts from store
//
// This function is auto-generated
func SearchComposeCharts(ctx context.Context, s ComposeCharts, f composeType.ChartFilter) (composeType.ChartSet, composeType.ChartFilter, error) {
	return s.SearchComposeCharts(ctx, f)
}

// CreateComposeChart creates one or more ComposeCharts in store
//
// This function is auto-generated
func CreateComposeChart(ctx context.Context, s ComposeCharts, rr ...*composeType.Chart) error {
	return s.CreateComposeChart(ctx, rr...)
}

// UpdateComposeChart updates one or more (existing) ComposeCharts in store
//
// This function is auto-generated
func UpdateComposeChart(ctx context.Context, s ComposeCharts, rr ...*composeType.Chart) error {
	return s.UpdateComposeChart(ctx, rr...)
}

// UpsertComposeChart creates new or updates existing one or more ComposeCharts in store
//
// This function is auto-generated
func UpsertComposeChart(ctx context.Context, s ComposeCharts, rr ...*composeType.Chart) error {
	return s.UpsertComposeChart(ctx, rr...)
}

// DeleteComposeChart deletes one or more ComposeCharts from store
//
// This function is auto-generated
func DeleteComposeChart(ctx context.Context, s ComposeCharts, rr ...*composeType.Chart) error {
	return s.DeleteComposeChart(ctx, rr...)
}

// DeleteComposeChartByID deletes one or more ComposeCharts from store
//
// This function is auto-generated
func DeleteComposeChartByID(ctx context.Context, s ComposeCharts, id uint64) error {
	return s.DeleteComposeChartByID(ctx, id)
}

// TruncateComposeCharts Deletes all ComposeCharts from store
//
// This function is auto-generated
func TruncateComposeCharts(ctx context.Context, s ComposeCharts) error {
	return s.TruncateComposeCharts(ctx)
}

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
//
// This function is auto-generated
func LookupComposeChartByID(ctx context.Context, s ComposeCharts, id uint64) (*composeType.Chart, error) {
	return s.LookupComposeChartByID(ctx, id)
}

// LookupComposeChartByNamespaceIDHandle searches for compose chart by handle (case-insensitive)
//
// This function is auto-generated
func LookupComposeChartByNamespaceIDHandle(ctx context.Context, s ComposeCharts, namespaceID uint64, handle string) (*composeType.Chart, error) {
	return s.LookupComposeChartByNamespaceIDHandle(ctx, namespaceID, handle)
}

// SearchComposeModules returns all matching ComposeModules from store
//
// This function is auto-generated
func SearchComposeModules(ctx context.Context, s ComposeModules, f composeType.ModuleFilter) (composeType.ModuleSet, composeType.ModuleFilter, error) {
	return s.SearchComposeModules(ctx, f)
}

// CreateComposeModule creates one or more ComposeModules in store
//
// This function is auto-generated
func CreateComposeModule(ctx context.Context, s ComposeModules, rr ...*composeType.Module) error {
	return s.CreateComposeModule(ctx, rr...)
}

// UpdateComposeModule updates one or more (existing) ComposeModules in store
//
// This function is auto-generated
func UpdateComposeModule(ctx context.Context, s ComposeModules, rr ...*composeType.Module) error {
	return s.UpdateComposeModule(ctx, rr...)
}

// UpsertComposeModule creates new or updates existing one or more ComposeModules in store
//
// This function is auto-generated
func UpsertComposeModule(ctx context.Context, s ComposeModules, rr ...*composeType.Module) error {
	return s.UpsertComposeModule(ctx, rr...)
}

// DeleteComposeModule deletes one or more ComposeModules from store
//
// This function is auto-generated
func DeleteComposeModule(ctx context.Context, s ComposeModules, rr ...*composeType.Module) error {
	return s.DeleteComposeModule(ctx, rr...)
}

// DeleteComposeModuleByID deletes one or more ComposeModules from store
//
// This function is auto-generated
func DeleteComposeModuleByID(ctx context.Context, s ComposeModules, id uint64) error {
	return s.DeleteComposeModuleByID(ctx, id)
}

// TruncateComposeModules Deletes all ComposeModules from store
//
// This function is auto-generated
func TruncateComposeModules(ctx context.Context, s ComposeModules) error {
	return s.TruncateComposeModules(ctx)
}

// LookupComposeModuleByNamespaceIDHandle searches for compose module by handle (case-insensitive)
//
// This function is auto-generated
func LookupComposeModuleByNamespaceIDHandle(ctx context.Context, s ComposeModules, namespaceID uint64, handle string) (*composeType.Module, error) {
	return s.LookupComposeModuleByNamespaceIDHandle(ctx, namespaceID, handle)
}

// LookupComposeModuleByNamespaceIDName searches for compose module by name (case-insensitive)
//
// This function is auto-generated
func LookupComposeModuleByNamespaceIDName(ctx context.Context, s ComposeModules, namespaceID uint64, name string) (*composeType.Module, error) {
	return s.LookupComposeModuleByNamespaceIDName(ctx, namespaceID, name)
}

// LookupComposeModuleByID searches for compose module by ID
//
// It returns compose module even if deleted
//
// This function is auto-generated
func LookupComposeModuleByID(ctx context.Context, s ComposeModules, id uint64) (*composeType.Module, error) {
	return s.LookupComposeModuleByID(ctx, id)
}

// SearchComposeModuleFields returns all matching ComposeModuleFields from store
//
// This function is auto-generated
func SearchComposeModuleFields(ctx context.Context, s ComposeModuleFields, f composeType.ModuleFieldFilter) (composeType.ModuleFieldSet, composeType.ModuleFieldFilter, error) {
	return s.SearchComposeModuleFields(ctx, f)
}

// CreateComposeModuleField creates one or more ComposeModuleFields in store
//
// This function is auto-generated
func CreateComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*composeType.ModuleField) error {
	return s.CreateComposeModuleField(ctx, rr...)
}

// UpdateComposeModuleField updates one or more (existing) ComposeModuleFields in store
//
// This function is auto-generated
func UpdateComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*composeType.ModuleField) error {
	return s.UpdateComposeModuleField(ctx, rr...)
}

// UpsertComposeModuleField creates new or updates existing one or more ComposeModuleFields in store
//
// This function is auto-generated
func UpsertComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*composeType.ModuleField) error {
	return s.UpsertComposeModuleField(ctx, rr...)
}

// DeleteComposeModuleField deletes one or more ComposeModuleFields from store
//
// This function is auto-generated
func DeleteComposeModuleField(ctx context.Context, s ComposeModuleFields, rr ...*composeType.ModuleField) error {
	return s.DeleteComposeModuleField(ctx, rr...)
}

// DeleteComposeModuleFieldByID deletes one or more ComposeModuleFields from store
//
// This function is auto-generated
func DeleteComposeModuleFieldByID(ctx context.Context, s ComposeModuleFields, id uint64) error {
	return s.DeleteComposeModuleFieldByID(ctx, id)
}

// TruncateComposeModuleFields Deletes all ComposeModuleFields from store
//
// This function is auto-generated
func TruncateComposeModuleFields(ctx context.Context, s ComposeModuleFields) error {
	return s.TruncateComposeModuleFields(ctx)
}

// LookupComposeModuleFieldByModuleIDName searches for compose module field by name (case-insensitive)
//
// This function is auto-generated
func LookupComposeModuleFieldByModuleIDName(ctx context.Context, s ComposeModuleFields, moduleID uint64, name string) (*composeType.ModuleField, error) {
	return s.LookupComposeModuleFieldByModuleIDName(ctx, moduleID, name)
}

// LookupComposeModuleFieldByID searches for compose module field by ID
//
// This function is auto-generated
func LookupComposeModuleFieldByID(ctx context.Context, s ComposeModuleFields, id uint64) (*composeType.ModuleField, error) {
	return s.LookupComposeModuleFieldByID(ctx, id)
}

// SearchComposeNamespaces returns all matching ComposeNamespaces from store
//
// This function is auto-generated
func SearchComposeNamespaces(ctx context.Context, s ComposeNamespaces, f composeType.NamespaceFilter) (composeType.NamespaceSet, composeType.NamespaceFilter, error) {
	return s.SearchComposeNamespaces(ctx, f)
}

// CreateComposeNamespace creates one or more ComposeNamespaces in store
//
// This function is auto-generated
func CreateComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*composeType.Namespace) error {
	return s.CreateComposeNamespace(ctx, rr...)
}

// UpdateComposeNamespace updates one or more (existing) ComposeNamespaces in store
//
// This function is auto-generated
func UpdateComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*composeType.Namespace) error {
	return s.UpdateComposeNamespace(ctx, rr...)
}

// UpsertComposeNamespace creates new or updates existing one or more ComposeNamespaces in store
//
// This function is auto-generated
func UpsertComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*composeType.Namespace) error {
	return s.UpsertComposeNamespace(ctx, rr...)
}

// DeleteComposeNamespace deletes one or more ComposeNamespaces from store
//
// This function is auto-generated
func DeleteComposeNamespace(ctx context.Context, s ComposeNamespaces, rr ...*composeType.Namespace) error {
	return s.DeleteComposeNamespace(ctx, rr...)
}

// DeleteComposeNamespaceByID deletes one or more ComposeNamespaces from store
//
// This function is auto-generated
func DeleteComposeNamespaceByID(ctx context.Context, s ComposeNamespaces, id uint64) error {
	return s.DeleteComposeNamespaceByID(ctx, id)
}

// TruncateComposeNamespaces Deletes all ComposeNamespaces from store
//
// This function is auto-generated
func TruncateComposeNamespaces(ctx context.Context, s ComposeNamespaces) error {
	return s.TruncateComposeNamespaces(ctx)
}

// LookupComposeNamespaceBySlug searches for namespace by slug (case-insensitive)
//
// This function is auto-generated
func LookupComposeNamespaceBySlug(ctx context.Context, s ComposeNamespaces, slug string) (*composeType.Namespace, error) {
	return s.LookupComposeNamespaceBySlug(ctx, slug)
}

// LookupComposeNamespaceByID searches for compose namespace by ID
//
// It returns compose namespace even if deleted
//
// This function is auto-generated
func LookupComposeNamespaceByID(ctx context.Context, s ComposeNamespaces, id uint64) (*composeType.Namespace, error) {
	return s.LookupComposeNamespaceByID(ctx, id)
}

// SearchComposePages returns all matching ComposePages from store
//
// This function is auto-generated
func SearchComposePages(ctx context.Context, s ComposePages, f composeType.PageFilter) (composeType.PageSet, composeType.PageFilter, error) {
	return s.SearchComposePages(ctx, f)
}

// CreateComposePage creates one or more ComposePages in store
//
// This function is auto-generated
func CreateComposePage(ctx context.Context, s ComposePages, rr ...*composeType.Page) error {
	return s.CreateComposePage(ctx, rr...)
}

// UpdateComposePage updates one or more (existing) ComposePages in store
//
// This function is auto-generated
func UpdateComposePage(ctx context.Context, s ComposePages, rr ...*composeType.Page) error {
	return s.UpdateComposePage(ctx, rr...)
}

// UpsertComposePage creates new or updates existing one or more ComposePages in store
//
// This function is auto-generated
func UpsertComposePage(ctx context.Context, s ComposePages, rr ...*composeType.Page) error {
	return s.UpsertComposePage(ctx, rr...)
}

// DeleteComposePage deletes one or more ComposePages from store
//
// This function is auto-generated
func DeleteComposePage(ctx context.Context, s ComposePages, rr ...*composeType.Page) error {
	return s.DeleteComposePage(ctx, rr...)
}

// DeleteComposePageByID deletes one or more ComposePages from store
//
// This function is auto-generated
func DeleteComposePageByID(ctx context.Context, s ComposePages, id uint64) error {
	return s.DeleteComposePageByID(ctx, id)
}

// TruncateComposePages Deletes all ComposePages from store
//
// This function is auto-generated
func TruncateComposePages(ctx context.Context, s ComposePages) error {
	return s.TruncateComposePages(ctx)
}

// LookupComposePageByNamespaceIDHandle searches for page by handle (case-insensitive)
//
// This function is auto-generated
func LookupComposePageByNamespaceIDHandle(ctx context.Context, s ComposePages, namespaceID uint64, handle string) (*composeType.Page, error) {
	return s.LookupComposePageByNamespaceIDHandle(ctx, namespaceID, handle)
}

// LookupComposePageByNamespaceIDModuleID searches for page by moduleID
//
// This function is auto-generated
func LookupComposePageByNamespaceIDModuleID(ctx context.Context, s ComposePages, namespaceID uint64, moduleID uint64) (*composeType.Page, error) {
	return s.LookupComposePageByNamespaceIDModuleID(ctx, namespaceID, moduleID)
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
//
// This function is auto-generated
func LookupComposePageByID(ctx context.Context, s ComposePages, id uint64) (*composeType.Page, error) {
	return s.LookupComposePageByID(ctx, id)
}

// ReorderComposePages
//
// This function is auto-generated
func ReorderComposePages(ctx context.Context, s ComposePages, namespace_id uint64, parent_id uint64, page_ids []uint64) error {
	return s.ReorderComposePages(ctx, namespace_id, parent_id, page_ids)
}

// SearchCredentials returns all matching Credentials from store
//
// This function is auto-generated
func SearchCredentials(ctx context.Context, s Credentials, f systemType.CredentialFilter) (systemType.CredentialSet, systemType.CredentialFilter, error) {
	return s.SearchCredentials(ctx, f)
}

// CreateCredential creates one or more Credentials in store
//
// This function is auto-generated
func CreateCredential(ctx context.Context, s Credentials, rr ...*systemType.Credential) error {
	return s.CreateCredential(ctx, rr...)
}

// UpdateCredential updates one or more (existing) Credentials in store
//
// This function is auto-generated
func UpdateCredential(ctx context.Context, s Credentials, rr ...*systemType.Credential) error {
	return s.UpdateCredential(ctx, rr...)
}

// UpsertCredential creates new or updates existing one or more Credentials in store
//
// This function is auto-generated
func UpsertCredential(ctx context.Context, s Credentials, rr ...*systemType.Credential) error {
	return s.UpsertCredential(ctx, rr...)
}

// DeleteCredential deletes one or more Credentials from store
//
// This function is auto-generated
func DeleteCredential(ctx context.Context, s Credentials, rr ...*systemType.Credential) error {
	return s.DeleteCredential(ctx, rr...)
}

// DeleteCredentialByID deletes one or more Credentials from store
//
// This function is auto-generated
func DeleteCredentialByID(ctx context.Context, s Credentials, id uint64) error {
	return s.DeleteCredentialByID(ctx, id)
}

// TruncateCredentials Deletes all Credentials from store
//
// This function is auto-generated
func TruncateCredentials(ctx context.Context, s Credentials) error {
	return s.TruncateCredentials(ctx)
}

// LookupCredentialByID searches for credentials by ID
//
// It returns credentials even if deleted
//
// This function is auto-generated
func LookupCredentialByID(ctx context.Context, s Credentials, id uint64) (*systemType.Credential, error) {
	return s.LookupCredentialByID(ctx, id)
}

// SearchDalConnections returns all matching DalConnections from store
//
// This function is auto-generated
func SearchDalConnections(ctx context.Context, s DalConnections, f systemType.DalConnectionFilter) (systemType.DalConnectionSet, systemType.DalConnectionFilter, error) {
	return s.SearchDalConnections(ctx, f)
}

// CreateDalConnection creates one or more DalConnections in store
//
// This function is auto-generated
func CreateDalConnection(ctx context.Context, s DalConnections, rr ...*systemType.DalConnection) error {
	return s.CreateDalConnection(ctx, rr...)
}

// UpdateDalConnection updates one or more (existing) DalConnections in store
//
// This function is auto-generated
func UpdateDalConnection(ctx context.Context, s DalConnections, rr ...*systemType.DalConnection) error {
	return s.UpdateDalConnection(ctx, rr...)
}

// UpsertDalConnection creates new or updates existing one or more DalConnections in store
//
// This function is auto-generated
func UpsertDalConnection(ctx context.Context, s DalConnections, rr ...*systemType.DalConnection) error {
	return s.UpsertDalConnection(ctx, rr...)
}

// DeleteDalConnection deletes one or more DalConnections from store
//
// This function is auto-generated
func DeleteDalConnection(ctx context.Context, s DalConnections, rr ...*systemType.DalConnection) error {
	return s.DeleteDalConnection(ctx, rr...)
}

// DeleteDalConnectionByID deletes one or more DalConnections from store
//
// This function is auto-generated
func DeleteDalConnectionByID(ctx context.Context, s DalConnections, id uint64) error {
	return s.DeleteDalConnectionByID(ctx, id)
}

// TruncateDalConnections Deletes all DalConnections from store
//
// This function is auto-generated
func TruncateDalConnections(ctx context.Context, s DalConnections) error {
	return s.TruncateDalConnections(ctx)
}

// LookupDalConnectionByID searches for connection by ID
//
// It returns connection even if deleted or suspended
//
// This function is auto-generated
func LookupDalConnectionByID(ctx context.Context, s DalConnections, id uint64) (*systemType.DalConnection, error) {
	return s.LookupDalConnectionByID(ctx, id)
}

// LookupDalConnectionByHandle searches for connection by handle
//
// It returns only valid connection (not deleted)
//
// This function is auto-generated
func LookupDalConnectionByHandle(ctx context.Context, s DalConnections, handle string) (*systemType.DalConnection, error) {
	return s.LookupDalConnectionByHandle(ctx, handle)
}

// SearchDalSensitivityLevels returns all matching DalSensitivityLevels from store
//
// This function is auto-generated
func SearchDalSensitivityLevels(ctx context.Context, s DalSensitivityLevels, f systemType.DalSensitivityLevelFilter) (systemType.DalSensitivityLevelSet, systemType.DalSensitivityLevelFilter, error) {
	return s.SearchDalSensitivityLevels(ctx, f)
}

// CreateDalSensitivityLevel creates one or more DalSensitivityLevels in store
//
// This function is auto-generated
func CreateDalSensitivityLevel(ctx context.Context, s DalSensitivityLevels, rr ...*systemType.DalSensitivityLevel) error {
	return s.CreateDalSensitivityLevel(ctx, rr...)
}

// UpdateDalSensitivityLevel updates one or more (existing) DalSensitivityLevels in store
//
// This function is auto-generated
func UpdateDalSensitivityLevel(ctx context.Context, s DalSensitivityLevels, rr ...*systemType.DalSensitivityLevel) error {
	return s.UpdateDalSensitivityLevel(ctx, rr...)
}

// UpsertDalSensitivityLevel creates new or updates existing one or more DalSensitivityLevels in store
//
// This function is auto-generated
func UpsertDalSensitivityLevel(ctx context.Context, s DalSensitivityLevels, rr ...*systemType.DalSensitivityLevel) error {
	return s.UpsertDalSensitivityLevel(ctx, rr...)
}

// DeleteDalSensitivityLevel deletes one or more DalSensitivityLevels from store
//
// This function is auto-generated
func DeleteDalSensitivityLevel(ctx context.Context, s DalSensitivityLevels, rr ...*systemType.DalSensitivityLevel) error {
	return s.DeleteDalSensitivityLevel(ctx, rr...)
}

// DeleteDalSensitivityLevelByID deletes one or more DalSensitivityLevels from store
//
// This function is auto-generated
func DeleteDalSensitivityLevelByID(ctx context.Context, s DalSensitivityLevels, id uint64) error {
	return s.DeleteDalSensitivityLevelByID(ctx, id)
}

// TruncateDalSensitivityLevels Deletes all DalSensitivityLevels from store
//
// This function is auto-generated
func TruncateDalSensitivityLevels(ctx context.Context, s DalSensitivityLevels) error {
	return s.TruncateDalSensitivityLevels(ctx)
}

// LookupDalSensitivityLevelByID searches for user by ID
//
// It returns user even if deleted or suspended
//
// This function is auto-generated
func LookupDalSensitivityLevelByID(ctx context.Context, s DalSensitivityLevels, id uint64) (*systemType.DalSensitivityLevel, error) {
	return s.LookupDalSensitivityLevelByID(ctx, id)
}

// SearchDataPrivacyRequests returns all matching DataPrivacyRequests from store
//
// This function is auto-generated
func SearchDataPrivacyRequests(ctx context.Context, s DataPrivacyRequests, f systemType.DataPrivacyRequestFilter) (systemType.DataPrivacyRequestSet, systemType.DataPrivacyRequestFilter, error) {
	return s.SearchDataPrivacyRequests(ctx, f)
}

// CreateDataPrivacyRequest creates one or more DataPrivacyRequests in store
//
// This function is auto-generated
func CreateDataPrivacyRequest(ctx context.Context, s DataPrivacyRequests, rr ...*systemType.DataPrivacyRequest) error {
	return s.CreateDataPrivacyRequest(ctx, rr...)
}

// UpdateDataPrivacyRequest updates one or more (existing) DataPrivacyRequests in store
//
// This function is auto-generated
func UpdateDataPrivacyRequest(ctx context.Context, s DataPrivacyRequests, rr ...*systemType.DataPrivacyRequest) error {
	return s.UpdateDataPrivacyRequest(ctx, rr...)
}

// UpsertDataPrivacyRequest creates new or updates existing one or more DataPrivacyRequests in store
//
// This function is auto-generated
func UpsertDataPrivacyRequest(ctx context.Context, s DataPrivacyRequests, rr ...*systemType.DataPrivacyRequest) error {
	return s.UpsertDataPrivacyRequest(ctx, rr...)
}

// DeleteDataPrivacyRequest deletes one or more DataPrivacyRequests from store
//
// This function is auto-generated
func DeleteDataPrivacyRequest(ctx context.Context, s DataPrivacyRequests, rr ...*systemType.DataPrivacyRequest) error {
	return s.DeleteDataPrivacyRequest(ctx, rr...)
}

// DeleteDataPrivacyRequestByID deletes one or more DataPrivacyRequests from store
//
// This function is auto-generated
func DeleteDataPrivacyRequestByID(ctx context.Context, s DataPrivacyRequests, id uint64) error {
	return s.DeleteDataPrivacyRequestByID(ctx, id)
}

// TruncateDataPrivacyRequests Deletes all DataPrivacyRequests from store
//
// This function is auto-generated
func TruncateDataPrivacyRequests(ctx context.Context, s DataPrivacyRequests) error {
	return s.TruncateDataPrivacyRequests(ctx)
}

// LookupDataPrivacyRequestByID searches for data privacy request by ID
//
// It returns data privacy request even if deleted
//
// This function is auto-generated
func LookupDataPrivacyRequestByID(ctx context.Context, s DataPrivacyRequests, id uint64) (*systemType.DataPrivacyRequest, error) {
	return s.LookupDataPrivacyRequestByID(ctx, id)
}

// SearchDataPrivacyRequestComments returns all matching DataPrivacyRequestComments from store
//
// This function is auto-generated
func SearchDataPrivacyRequestComments(ctx context.Context, s DataPrivacyRequestComments, f systemType.DataPrivacyRequestCommentFilter) (systemType.DataPrivacyRequestCommentSet, systemType.DataPrivacyRequestCommentFilter, error) {
	return s.SearchDataPrivacyRequestComments(ctx, f)
}

// CreateDataPrivacyRequestComment creates one or more DataPrivacyRequestComments in store
//
// This function is auto-generated
func CreateDataPrivacyRequestComment(ctx context.Context, s DataPrivacyRequestComments, rr ...*systemType.DataPrivacyRequestComment) error {
	return s.CreateDataPrivacyRequestComment(ctx, rr...)
}

// UpdateDataPrivacyRequestComment updates one or more (existing) DataPrivacyRequestComments in store
//
// This function is auto-generated
func UpdateDataPrivacyRequestComment(ctx context.Context, s DataPrivacyRequestComments, rr ...*systemType.DataPrivacyRequestComment) error {
	return s.UpdateDataPrivacyRequestComment(ctx, rr...)
}

// UpsertDataPrivacyRequestComment creates new or updates existing one or more DataPrivacyRequestComments in store
//
// This function is auto-generated
func UpsertDataPrivacyRequestComment(ctx context.Context, s DataPrivacyRequestComments, rr ...*systemType.DataPrivacyRequestComment) error {
	return s.UpsertDataPrivacyRequestComment(ctx, rr...)
}

// DeleteDataPrivacyRequestComment deletes one or more DataPrivacyRequestComments from store
//
// This function is auto-generated
func DeleteDataPrivacyRequestComment(ctx context.Context, s DataPrivacyRequestComments, rr ...*systemType.DataPrivacyRequestComment) error {
	return s.DeleteDataPrivacyRequestComment(ctx, rr...)
}

// DeleteDataPrivacyRequestCommentByID deletes one or more DataPrivacyRequestComments from store
//
// This function is auto-generated
func DeleteDataPrivacyRequestCommentByID(ctx context.Context, s DataPrivacyRequestComments, id uint64) error {
	return s.DeleteDataPrivacyRequestCommentByID(ctx, id)
}

// TruncateDataPrivacyRequestComments Deletes all DataPrivacyRequestComments from store
//
// This function is auto-generated
func TruncateDataPrivacyRequestComments(ctx context.Context, s DataPrivacyRequestComments) error {
	return s.TruncateDataPrivacyRequestComments(ctx)
}

// SearchFederationExposedModules returns all matching FederationExposedModules from store
//
// This function is auto-generated
func SearchFederationExposedModules(ctx context.Context, s FederationExposedModules, f federationType.ExposedModuleFilter) (federationType.ExposedModuleSet, federationType.ExposedModuleFilter, error) {
	return s.SearchFederationExposedModules(ctx, f)
}

// CreateFederationExposedModule creates one or more FederationExposedModules in store
//
// This function is auto-generated
func CreateFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*federationType.ExposedModule) error {
	return s.CreateFederationExposedModule(ctx, rr...)
}

// UpdateFederationExposedModule updates one or more (existing) FederationExposedModules in store
//
// This function is auto-generated
func UpdateFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*federationType.ExposedModule) error {
	return s.UpdateFederationExposedModule(ctx, rr...)
}

// UpsertFederationExposedModule creates new or updates existing one or more FederationExposedModules in store
//
// This function is auto-generated
func UpsertFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*federationType.ExposedModule) error {
	return s.UpsertFederationExposedModule(ctx, rr...)
}

// DeleteFederationExposedModule deletes one or more FederationExposedModules from store
//
// This function is auto-generated
func DeleteFederationExposedModule(ctx context.Context, s FederationExposedModules, rr ...*federationType.ExposedModule) error {
	return s.DeleteFederationExposedModule(ctx, rr...)
}

// DeleteFederationExposedModuleByID deletes one or more FederationExposedModules from store
//
// This function is auto-generated
func DeleteFederationExposedModuleByID(ctx context.Context, s FederationExposedModules, id uint64) error {
	return s.DeleteFederationExposedModuleByID(ctx, id)
}

// TruncateFederationExposedModules Deletes all FederationExposedModules from store
//
// This function is auto-generated
func TruncateFederationExposedModules(ctx context.Context, s FederationExposedModules) error {
	return s.TruncateFederationExposedModules(ctx)
}

// LookupFederationExposedModuleByID searches for federation module by ID
//
// It returns federation module
//
// This function is auto-generated
func LookupFederationExposedModuleByID(ctx context.Context, s FederationExposedModules, id uint64) (*federationType.ExposedModule, error) {
	return s.LookupFederationExposedModuleByID(ctx, id)
}

// SearchFederationModuleMappings returns all matching FederationModuleMappings from store
//
// This function is auto-generated
func SearchFederationModuleMappings(ctx context.Context, s FederationModuleMappings, f federationType.ModuleMappingFilter) (federationType.ModuleMappingSet, federationType.ModuleMappingFilter, error) {
	return s.SearchFederationModuleMappings(ctx, f)
}

// CreateFederationModuleMapping creates one or more FederationModuleMappings in store
//
// This function is auto-generated
func CreateFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*federationType.ModuleMapping) error {
	return s.CreateFederationModuleMapping(ctx, rr...)
}

// UpdateFederationModuleMapping updates one or more (existing) FederationModuleMappings in store
//
// This function is auto-generated
func UpdateFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*federationType.ModuleMapping) error {
	return s.UpdateFederationModuleMapping(ctx, rr...)
}

// UpsertFederationModuleMapping creates new or updates existing one or more FederationModuleMappings in store
//
// This function is auto-generated
func UpsertFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*federationType.ModuleMapping) error {
	return s.UpsertFederationModuleMapping(ctx, rr...)
}

// DeleteFederationModuleMapping deletes one or more FederationModuleMappings from store
//
// This function is auto-generated
func DeleteFederationModuleMapping(ctx context.Context, s FederationModuleMappings, rr ...*federationType.ModuleMapping) error {
	return s.DeleteFederationModuleMapping(ctx, rr...)
}

// DeleteFederationModuleMappingByID deletes one or more FederationModuleMappings from store
//
// This function is auto-generated
func DeleteFederationModuleMappingByNodeID(ctx context.Context, s FederationModuleMappings, nodeID uint64) error {
	return s.DeleteFederationModuleMappingByNodeID(ctx, nodeID)
}

// TruncateFederationModuleMappings Deletes all FederationModuleMappings from store
//
// This function is auto-generated
func TruncateFederationModuleMappings(ctx context.Context, s FederationModuleMappings) error {
	return s.TruncateFederationModuleMappings(ctx)
}

// LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID searches for module mapping by federation module id and compose module id
//
// It returns module mapping
//
// This function is auto-generated
func LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID(ctx context.Context, s FederationModuleMappings, federationModuleID uint64, composeModuleID uint64, composeNamespaceID uint64) (*federationType.ModuleMapping, error) {
	return s.LookupFederationModuleMappingByFederationModuleIDComposeModuleIDComposeNamespaceID(ctx, federationModuleID, composeModuleID, composeNamespaceID)
}

// LookupFederationModuleMappingByFederationModuleID searches for module mapping by federation module id
//
// It returns module mapping
//
// This function is auto-generated
func LookupFederationModuleMappingByFederationModuleID(ctx context.Context, s FederationModuleMappings, federationModuleID uint64) (*federationType.ModuleMapping, error) {
	return s.LookupFederationModuleMappingByFederationModuleID(ctx, federationModuleID)
}

// SearchFederationNodes returns all matching FederationNodes from store
//
// This function is auto-generated
func SearchFederationNodes(ctx context.Context, s FederationNodes, f federationType.NodeFilter) (federationType.NodeSet, federationType.NodeFilter, error) {
	return s.SearchFederationNodes(ctx, f)
}

// CreateFederationNode creates one or more FederationNodes in store
//
// This function is auto-generated
func CreateFederationNode(ctx context.Context, s FederationNodes, rr ...*federationType.Node) error {
	return s.CreateFederationNode(ctx, rr...)
}

// UpdateFederationNode updates one or more (existing) FederationNodes in store
//
// This function is auto-generated
func UpdateFederationNode(ctx context.Context, s FederationNodes, rr ...*federationType.Node) error {
	return s.UpdateFederationNode(ctx, rr...)
}

// UpsertFederationNode creates new or updates existing one or more FederationNodes in store
//
// This function is auto-generated
func UpsertFederationNode(ctx context.Context, s FederationNodes, rr ...*federationType.Node) error {
	return s.UpsertFederationNode(ctx, rr...)
}

// DeleteFederationNode deletes one or more FederationNodes from store
//
// This function is auto-generated
func DeleteFederationNode(ctx context.Context, s FederationNodes, rr ...*federationType.Node) error {
	return s.DeleteFederationNode(ctx, rr...)
}

// DeleteFederationNodeByID deletes one or more FederationNodes from store
//
// This function is auto-generated
func DeleteFederationNodeByID(ctx context.Context, s FederationNodes, id uint64) error {
	return s.DeleteFederationNodeByID(ctx, id)
}

// TruncateFederationNodes Deletes all FederationNodes from store
//
// This function is auto-generated
func TruncateFederationNodes(ctx context.Context, s FederationNodes) error {
	return s.TruncateFederationNodes(ctx)
}

// LookupFederationNodeByID searches for federation node by ID
//
// It returns federation node
//
// This function is auto-generated
func LookupFederationNodeByID(ctx context.Context, s FederationNodes, id uint64) (*federationType.Node, error) {
	return s.LookupFederationNodeByID(ctx, id)
}

// LookupFederationNodeByBaseURLSharedNodeID searches for node by shared-node-id and base-url
//
// This function is auto-generated
func LookupFederationNodeByBaseURLSharedNodeID(ctx context.Context, s FederationNodes, baseURL string, sharedNodeID uint64) (*federationType.Node, error) {
	return s.LookupFederationNodeByBaseURLSharedNodeID(ctx, baseURL, sharedNodeID)
}

// LookupFederationNodeBySharedNodeID searches for node by shared-node-id
//
// This function is auto-generated
func LookupFederationNodeBySharedNodeID(ctx context.Context, s FederationNodes, sharedNodeID uint64) (*federationType.Node, error) {
	return s.LookupFederationNodeBySharedNodeID(ctx, sharedNodeID)
}

// SearchFederationNodeSyncs returns all matching FederationNodeSyncs from store
//
// This function is auto-generated
func SearchFederationNodeSyncs(ctx context.Context, s FederationNodeSyncs, f federationType.NodeSyncFilter) (federationType.NodeSyncSet, federationType.NodeSyncFilter, error) {
	return s.SearchFederationNodeSyncs(ctx, f)
}

// CreateFederationNodeSync creates one or more FederationNodeSyncs in store
//
// This function is auto-generated
func CreateFederationNodeSync(ctx context.Context, s FederationNodeSyncs, rr ...*federationType.NodeSync) error {
	return s.CreateFederationNodeSync(ctx, rr...)
}

// UpdateFederationNodeSync updates one or more (existing) FederationNodeSyncs in store
//
// This function is auto-generated
func UpdateFederationNodeSync(ctx context.Context, s FederationNodeSyncs, rr ...*federationType.NodeSync) error {
	return s.UpdateFederationNodeSync(ctx, rr...)
}

// UpsertFederationNodeSync creates new or updates existing one or more FederationNodeSyncs in store
//
// This function is auto-generated
func UpsertFederationNodeSync(ctx context.Context, s FederationNodeSyncs, rr ...*federationType.NodeSync) error {
	return s.UpsertFederationNodeSync(ctx, rr...)
}

// DeleteFederationNodeSync deletes one or more FederationNodeSyncs from store
//
// This function is auto-generated
func DeleteFederationNodeSync(ctx context.Context, s FederationNodeSyncs, rr ...*federationType.NodeSync) error {
	return s.DeleteFederationNodeSync(ctx, rr...)
}

// DeleteFederationNodeSyncByID deletes one or more FederationNodeSyncs from store
//
// This function is auto-generated
func DeleteFederationNodeSyncByNodeID(ctx context.Context, s FederationNodeSyncs, nodeID uint64) error {
	return s.DeleteFederationNodeSyncByNodeID(ctx, nodeID)
}

// TruncateFederationNodeSyncs Deletes all FederationNodeSyncs from store
//
// This function is auto-generated
func TruncateFederationNodeSyncs(ctx context.Context, s FederationNodeSyncs) error {
	return s.TruncateFederationNodeSyncs(ctx)
}

// LookupFederationNodeSyncByNodeID searches for sync activity by node ID
//
// It returns sync activity
//
// This function is auto-generated
func LookupFederationNodeSyncByNodeID(ctx context.Context, s FederationNodeSyncs, nodeID uint64) (*federationType.NodeSync, error) {
	return s.LookupFederationNodeSyncByNodeID(ctx, nodeID)
}

// LookupFederationNodeSyncByNodeIDModuleIDSyncTypeSyncStatus searches for activity by node, type and status
//
// It returns sync activity
//
// This function is auto-generated
func LookupFederationNodeSyncByNodeIDModuleIDSyncTypeSyncStatus(ctx context.Context, s FederationNodeSyncs, nodeID uint64, moduleID uint64, syncType string, syncStatus string) (*federationType.NodeSync, error) {
	return s.LookupFederationNodeSyncByNodeIDModuleIDSyncTypeSyncStatus(ctx, nodeID, moduleID, syncType, syncStatus)
}

// SearchFederationSharedModules returns all matching FederationSharedModules from store
//
// This function is auto-generated
func SearchFederationSharedModules(ctx context.Context, s FederationSharedModules, f federationType.SharedModuleFilter) (federationType.SharedModuleSet, federationType.SharedModuleFilter, error) {
	return s.SearchFederationSharedModules(ctx, f)
}

// CreateFederationSharedModule creates one or more FederationSharedModules in store
//
// This function is auto-generated
func CreateFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*federationType.SharedModule) error {
	return s.CreateFederationSharedModule(ctx, rr...)
}

// UpdateFederationSharedModule updates one or more (existing) FederationSharedModules in store
//
// This function is auto-generated
func UpdateFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*federationType.SharedModule) error {
	return s.UpdateFederationSharedModule(ctx, rr...)
}

// UpsertFederationSharedModule creates new or updates existing one or more FederationSharedModules in store
//
// This function is auto-generated
func UpsertFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*federationType.SharedModule) error {
	return s.UpsertFederationSharedModule(ctx, rr...)
}

// DeleteFederationSharedModule deletes one or more FederationSharedModules from store
//
// This function is auto-generated
func DeleteFederationSharedModule(ctx context.Context, s FederationSharedModules, rr ...*federationType.SharedModule) error {
	return s.DeleteFederationSharedModule(ctx, rr...)
}

// DeleteFederationSharedModuleByID deletes one or more FederationSharedModules from store
//
// This function is auto-generated
func DeleteFederationSharedModuleByID(ctx context.Context, s FederationSharedModules, id uint64) error {
	return s.DeleteFederationSharedModuleByID(ctx, id)
}

// TruncateFederationSharedModules Deletes all FederationSharedModules from store
//
// This function is auto-generated
func TruncateFederationSharedModules(ctx context.Context, s FederationSharedModules) error {
	return s.TruncateFederationSharedModules(ctx)
}

// LookupFederationSharedModuleByID searches for shared federation module by ID
//
// It returns shared federation module
//
// This function is auto-generated
func LookupFederationSharedModuleByID(ctx context.Context, s FederationSharedModules, id uint64) (*federationType.SharedModule, error) {
	return s.LookupFederationSharedModuleByID(ctx, id)
}

// SearchFlags returns all matching Flags from store
//
// This function is auto-generated
func SearchFlags(ctx context.Context, s Flags, f flagType.FlagFilter) (flagType.FlagSet, flagType.FlagFilter, error) {
	return s.SearchFlags(ctx, f)
}

// CreateFlag creates one or more Flags in store
//
// This function is auto-generated
func CreateFlag(ctx context.Context, s Flags, rr ...*flagType.Flag) error {
	return s.CreateFlag(ctx, rr...)
}

// UpdateFlag updates one or more (existing) Flags in store
//
// This function is auto-generated
func UpdateFlag(ctx context.Context, s Flags, rr ...*flagType.Flag) error {
	return s.UpdateFlag(ctx, rr...)
}

// UpsertFlag creates new or updates existing one or more Flags in store
//
// This function is auto-generated
func UpsertFlag(ctx context.Context, s Flags, rr ...*flagType.Flag) error {
	return s.UpsertFlag(ctx, rr...)
}

// DeleteFlag deletes one or more Flags from store
//
// This function is auto-generated
func DeleteFlag(ctx context.Context, s Flags, rr ...*flagType.Flag) error {
	return s.DeleteFlag(ctx, rr...)
}

// DeleteFlagByID deletes one or more Flags from store
//
// This function is auto-generated
func DeleteFlagByKindResourceIDOwnedByName(ctx context.Context, s Flags, kind string, resourceID uint64, ownedBy uint64, name string) error {
	return s.DeleteFlagByKindResourceIDOwnedByName(ctx, kind, resourceID, ownedBy, name)
}

// TruncateFlags Deletes all Flags from store
//
// This function is auto-generated
func TruncateFlags(ctx context.Context, s Flags) error {
	return s.TruncateFlags(ctx)
}

// LookupFlagByKindResourceIDOwnedByName searches for flag by kind, resource ID, owner and name
//
// This function is auto-generated
func LookupFlagByKindResourceIDOwnedByName(ctx context.Context, s Flags, kind string, resourceID uint64, ownedBy uint64, name string) (*flagType.Flag, error) {
	return s.LookupFlagByKindResourceIDOwnedByName(ctx, kind, resourceID, ownedBy, name)
}

// SearchLabels returns all matching Labels from store
//
// This function is auto-generated
func SearchLabels(ctx context.Context, s Labels, f labelsType.LabelFilter) (labelsType.LabelSet, labelsType.LabelFilter, error) {
	return s.SearchLabels(ctx, f)
}

// CreateLabel creates one or more Labels in store
//
// This function is auto-generated
func CreateLabel(ctx context.Context, s Labels, rr ...*labelsType.Label) error {
	return s.CreateLabel(ctx, rr...)
}

// UpdateLabel updates one or more (existing) Labels in store
//
// This function is auto-generated
func UpdateLabel(ctx context.Context, s Labels, rr ...*labelsType.Label) error {
	return s.UpdateLabel(ctx, rr...)
}

// UpsertLabel creates new or updates existing one or more Labels in store
//
// This function is auto-generated
func UpsertLabel(ctx context.Context, s Labels, rr ...*labelsType.Label) error {
	return s.UpsertLabel(ctx, rr...)
}

// DeleteLabel deletes one or more Labels from store
//
// This function is auto-generated
func DeleteLabel(ctx context.Context, s Labels, rr ...*labelsType.Label) error {
	return s.DeleteLabel(ctx, rr...)
}

// DeleteLabelByID deletes one or more Labels from store
//
// This function is auto-generated
func DeleteLabelByKindResourceIDName(ctx context.Context, s Labels, kind string, resourceID uint64, name string) error {
	return s.DeleteLabelByKindResourceIDName(ctx, kind, resourceID, name)
}

// TruncateLabels Deletes all Labels from store
//
// This function is auto-generated
func TruncateLabels(ctx context.Context, s Labels) error {
	return s.TruncateLabels(ctx)
}

// LookupLabelByKindResourceIDName searches for label by kind, resource ID and name
//
// This function is auto-generated
func LookupLabelByKindResourceIDName(ctx context.Context, s Labels, kind string, resourceID uint64, name string) (*labelsType.Label, error) {
	return s.LookupLabelByKindResourceIDName(ctx, kind, resourceID, name)
}

// DeleteExtraLabels
//
// This function is auto-generated
func DeleteExtraLabels(ctx context.Context, s Labels, kind string, resourceId uint64, name ...string) error {
	return s.DeleteExtraLabels(ctx, kind, resourceId, name...)
}

// SearchQueues returns all matching Queues from store
//
// This function is auto-generated
func SearchQueues(ctx context.Context, s Queues, f systemType.QueueFilter) (systemType.QueueSet, systemType.QueueFilter, error) {
	return s.SearchQueues(ctx, f)
}

// CreateQueue creates one or more Queues in store
//
// This function is auto-generated
func CreateQueue(ctx context.Context, s Queues, rr ...*systemType.Queue) error {
	return s.CreateQueue(ctx, rr...)
}

// UpdateQueue updates one or more (existing) Queues in store
//
// This function is auto-generated
func UpdateQueue(ctx context.Context, s Queues, rr ...*systemType.Queue) error {
	return s.UpdateQueue(ctx, rr...)
}

// UpsertQueue creates new or updates existing one or more Queues in store
//
// This function is auto-generated
func UpsertQueue(ctx context.Context, s Queues, rr ...*systemType.Queue) error {
	return s.UpsertQueue(ctx, rr...)
}

// DeleteQueue deletes one or more Queues from store
//
// This function is auto-generated
func DeleteQueue(ctx context.Context, s Queues, rr ...*systemType.Queue) error {
	return s.DeleteQueue(ctx, rr...)
}

// DeleteQueueByID deletes one or more Queues from store
//
// This function is auto-generated
func DeleteQueueByID(ctx context.Context, s Queues, id uint64) error {
	return s.DeleteQueueByID(ctx, id)
}

// TruncateQueues Deletes all Queues from store
//
// This function is auto-generated
func TruncateQueues(ctx context.Context, s Queues) error {
	return s.TruncateQueues(ctx)
}

// LookupQueueByID searches for queue by ID
//
// This function is auto-generated
func LookupQueueByID(ctx context.Context, s Queues, id uint64) (*systemType.Queue, error) {
	return s.LookupQueueByID(ctx, id)
}

// LookupQueueByQueue searches for queue by queue name
//
// This function is auto-generated
func LookupQueueByQueue(ctx context.Context, s Queues, queue string) (*systemType.Queue, error) {
	return s.LookupQueueByQueue(ctx, queue)
}

// SearchQueueMessages returns all matching QueueMessages from store
//
// This function is auto-generated
func SearchQueueMessages(ctx context.Context, s QueueMessages, f systemType.QueueMessageFilter) (systemType.QueueMessageSet, systemType.QueueMessageFilter, error) {
	return s.SearchQueueMessages(ctx, f)
}

// CreateQueueMessage creates one or more QueueMessages in store
//
// This function is auto-generated
func CreateQueueMessage(ctx context.Context, s QueueMessages, rr ...*systemType.QueueMessage) error {
	return s.CreateQueueMessage(ctx, rr...)
}

// UpdateQueueMessage updates one or more (existing) QueueMessages in store
//
// This function is auto-generated
func UpdateQueueMessage(ctx context.Context, s QueueMessages, rr ...*systemType.QueueMessage) error {
	return s.UpdateQueueMessage(ctx, rr...)
}

// UpsertQueueMessage creates new or updates existing one or more QueueMessages in store
//
// This function is auto-generated
func UpsertQueueMessage(ctx context.Context, s QueueMessages, rr ...*systemType.QueueMessage) error {
	return s.UpsertQueueMessage(ctx, rr...)
}

// DeleteQueueMessage deletes one or more QueueMessages from store
//
// This function is auto-generated
func DeleteQueueMessage(ctx context.Context, s QueueMessages, rr ...*systemType.QueueMessage) error {
	return s.DeleteQueueMessage(ctx, rr...)
}

// DeleteQueueMessageByID deletes one or more QueueMessages from store
//
// This function is auto-generated
func DeleteQueueMessageByID(ctx context.Context, s QueueMessages, id uint64) error {
	return s.DeleteQueueMessageByID(ctx, id)
}

// TruncateQueueMessages Deletes all QueueMessages from store
//
// This function is auto-generated
func TruncateQueueMessages(ctx context.Context, s QueueMessages) error {
	return s.TruncateQueueMessages(ctx)
}

// SearchRbacRules returns all matching RbacRules from store
//
// This function is auto-generated
func SearchRbacRules(ctx context.Context, s RbacRules, f rbacType.RuleFilter) (rbacType.RuleSet, rbacType.RuleFilter, error) {
	return s.SearchRbacRules(ctx, f)
}

// CreateRbacRule creates one or more RbacRules in store
//
// This function is auto-generated
func CreateRbacRule(ctx context.Context, s RbacRules, rr ...*rbacType.Rule) error {
	return s.CreateRbacRule(ctx, rr...)
}

// UpdateRbacRule updates one or more (existing) RbacRules in store
//
// This function is auto-generated
func UpdateRbacRule(ctx context.Context, s RbacRules, rr ...*rbacType.Rule) error {
	return s.UpdateRbacRule(ctx, rr...)
}

// UpsertRbacRule creates new or updates existing one or more RbacRules in store
//
// This function is auto-generated
func UpsertRbacRule(ctx context.Context, s RbacRules, rr ...*rbacType.Rule) error {
	return s.UpsertRbacRule(ctx, rr...)
}

// DeleteRbacRule deletes one or more RbacRules from store
//
// This function is auto-generated
func DeleteRbacRule(ctx context.Context, s RbacRules, rr ...*rbacType.Rule) error {
	return s.DeleteRbacRule(ctx, rr...)
}

// DeleteRbacRuleByID deletes one or more RbacRules from store
//
// This function is auto-generated
func DeleteRbacRuleByRoleIDResourceOperation(ctx context.Context, s RbacRules, roleID uint64, resource string, operation string) error {
	return s.DeleteRbacRuleByRoleIDResourceOperation(ctx, roleID, resource, operation)
}

// TruncateRbacRules Deletes all RbacRules from store
//
// This function is auto-generated
func TruncateRbacRules(ctx context.Context, s RbacRules) error {
	return s.TruncateRbacRules(ctx)
}

// TransferRbacRules
//
// This function is auto-generated
func TransferRbacRules(ctx context.Context, s RbacRules, src uint64, dst uint64) error {
	return s.TransferRbacRules(ctx, src, dst)
}

// SearchReminders returns all matching Reminders from store
//
// This function is auto-generated
func SearchReminders(ctx context.Context, s Reminders, f systemType.ReminderFilter) (systemType.ReminderSet, systemType.ReminderFilter, error) {
	return s.SearchReminders(ctx, f)
}

// CreateReminder creates one or more Reminders in store
//
// This function is auto-generated
func CreateReminder(ctx context.Context, s Reminders, rr ...*systemType.Reminder) error {
	return s.CreateReminder(ctx, rr...)
}

// UpdateReminder updates one or more (existing) Reminders in store
//
// This function is auto-generated
func UpdateReminder(ctx context.Context, s Reminders, rr ...*systemType.Reminder) error {
	return s.UpdateReminder(ctx, rr...)
}

// UpsertReminder creates new or updates existing one or more Reminders in store
//
// This function is auto-generated
func UpsertReminder(ctx context.Context, s Reminders, rr ...*systemType.Reminder) error {
	return s.UpsertReminder(ctx, rr...)
}

// DeleteReminder deletes one or more Reminders from store
//
// This function is auto-generated
func DeleteReminder(ctx context.Context, s Reminders, rr ...*systemType.Reminder) error {
	return s.DeleteReminder(ctx, rr...)
}

// DeleteReminderByID deletes one or more Reminders from store
//
// This function is auto-generated
func DeleteReminderByID(ctx context.Context, s Reminders, id uint64) error {
	return s.DeleteReminderByID(ctx, id)
}

// TruncateReminders Deletes all Reminders from store
//
// This function is auto-generated
func TruncateReminders(ctx context.Context, s Reminders) error {
	return s.TruncateReminders(ctx)
}

// LookupReminderByID
//
// This function is auto-generated
func LookupReminderByID(ctx context.Context, s Reminders, id uint64) (*systemType.Reminder, error) {
	return s.LookupReminderByID(ctx, id)
}

// SearchReports returns all matching Reports from store
//
// This function is auto-generated
func SearchReports(ctx context.Context, s Reports, f systemType.ReportFilter) (systemType.ReportSet, systemType.ReportFilter, error) {
	return s.SearchReports(ctx, f)
}

// CreateReport creates one or more Reports in store
//
// This function is auto-generated
func CreateReport(ctx context.Context, s Reports, rr ...*systemType.Report) error {
	return s.CreateReport(ctx, rr...)
}

// UpdateReport updates one or more (existing) Reports in store
//
// This function is auto-generated
func UpdateReport(ctx context.Context, s Reports, rr ...*systemType.Report) error {
	return s.UpdateReport(ctx, rr...)
}

// UpsertReport creates new or updates existing one or more Reports in store
//
// This function is auto-generated
func UpsertReport(ctx context.Context, s Reports, rr ...*systemType.Report) error {
	return s.UpsertReport(ctx, rr...)
}

// DeleteReport deletes one or more Reports from store
//
// This function is auto-generated
func DeleteReport(ctx context.Context, s Reports, rr ...*systemType.Report) error {
	return s.DeleteReport(ctx, rr...)
}

// DeleteReportByID deletes one or more Reports from store
//
// This function is auto-generated
func DeleteReportByID(ctx context.Context, s Reports, id uint64) error {
	return s.DeleteReportByID(ctx, id)
}

// TruncateReports Deletes all Reports from store
//
// This function is auto-generated
func TruncateReports(ctx context.Context, s Reports) error {
	return s.TruncateReports(ctx)
}

// LookupReportByID searches for report by ID
//
// It returns report even if deleted
//
// This function is auto-generated
func LookupReportByID(ctx context.Context, s Reports, id uint64) (*systemType.Report, error) {
	return s.LookupReportByID(ctx, id)
}

// LookupReportByHandle searches for report by handle
//
// It returns report if deleted
//
// This function is auto-generated
func LookupReportByHandle(ctx context.Context, s Reports, handle string) (*systemType.Report, error) {
	return s.LookupReportByHandle(ctx, handle)
}

// SearchResourceActivitys returns all matching ResourceActivitys from store
//
// This function is auto-generated
func SearchResourceActivitys(ctx context.Context, s ResourceActivitys, f discoveryType.ResourceActivityFilter) (discoveryType.ResourceActivitySet, discoveryType.ResourceActivityFilter, error) {
	return s.SearchResourceActivitys(ctx, f)
}

// CreateResourceActivity creates one or more ResourceActivitys in store
//
// This function is auto-generated
func CreateResourceActivity(ctx context.Context, s ResourceActivitys, rr ...*discoveryType.ResourceActivity) error {
	return s.CreateResourceActivity(ctx, rr...)
}

// UpdateResourceActivity updates one or more (existing) ResourceActivitys in store
//
// This function is auto-generated
func UpdateResourceActivity(ctx context.Context, s ResourceActivitys, rr ...*discoveryType.ResourceActivity) error {
	return s.UpdateResourceActivity(ctx, rr...)
}

// UpsertResourceActivity creates new or updates existing one or more ResourceActivitys in store
//
// This function is auto-generated
func UpsertResourceActivity(ctx context.Context, s ResourceActivitys, rr ...*discoveryType.ResourceActivity) error {
	return s.UpsertResourceActivity(ctx, rr...)
}

// DeleteResourceActivity deletes one or more ResourceActivitys from store
//
// This function is auto-generated
func DeleteResourceActivity(ctx context.Context, s ResourceActivitys, rr ...*discoveryType.ResourceActivity) error {
	return s.DeleteResourceActivity(ctx, rr...)
}

// DeleteResourceActivityByID deletes one or more ResourceActivitys from store
//
// This function is auto-generated
func DeleteResourceActivityByID(ctx context.Context, s ResourceActivitys, id uint64) error {
	return s.DeleteResourceActivityByID(ctx, id)
}

// TruncateResourceActivitys Deletes all ResourceActivitys from store
//
// This function is auto-generated
func TruncateResourceActivitys(ctx context.Context, s ResourceActivitys) error {
	return s.TruncateResourceActivitys(ctx)
}

// SearchResourceTranslations returns all matching ResourceTranslations from store
//
// This function is auto-generated
func SearchResourceTranslations(ctx context.Context, s ResourceTranslations, f systemType.ResourceTranslationFilter) (systemType.ResourceTranslationSet, systemType.ResourceTranslationFilter, error) {
	return s.SearchResourceTranslations(ctx, f)
}

// CreateResourceTranslation creates one or more ResourceTranslations in store
//
// This function is auto-generated
func CreateResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*systemType.ResourceTranslation) error {
	return s.CreateResourceTranslation(ctx, rr...)
}

// UpdateResourceTranslation updates one or more (existing) ResourceTranslations in store
//
// This function is auto-generated
func UpdateResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*systemType.ResourceTranslation) error {
	return s.UpdateResourceTranslation(ctx, rr...)
}

// UpsertResourceTranslation creates new or updates existing one or more ResourceTranslations in store
//
// This function is auto-generated
func UpsertResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*systemType.ResourceTranslation) error {
	return s.UpsertResourceTranslation(ctx, rr...)
}

// DeleteResourceTranslation deletes one or more ResourceTranslations from store
//
// This function is auto-generated
func DeleteResourceTranslation(ctx context.Context, s ResourceTranslations, rr ...*systemType.ResourceTranslation) error {
	return s.DeleteResourceTranslation(ctx, rr...)
}

// DeleteResourceTranslationByID deletes one or more ResourceTranslations from store
//
// This function is auto-generated
func DeleteResourceTranslationByID(ctx context.Context, s ResourceTranslations, id uint64) error {
	return s.DeleteResourceTranslationByID(ctx, id)
}

// TruncateResourceTranslations Deletes all ResourceTranslations from store
//
// This function is auto-generated
func TruncateResourceTranslations(ctx context.Context, s ResourceTranslations) error {
	return s.TruncateResourceTranslations(ctx)
}

// LookupResourceTranslationByID searches for resource translation by ID
// It also returns deleted resource translations.
//
// This function is auto-generated
func LookupResourceTranslationByID(ctx context.Context, s ResourceTranslations, id uint64) (*systemType.ResourceTranslation, error) {
	return s.LookupResourceTranslationByID(ctx, id)
}

// TransformResource
//
// This function is auto-generated
func TransformResource(ctx context.Context, s ResourceTranslations, lang language.Tag) (map[string]map[string]*locale.ResourceTranslation, error) {
	return s.TransformResource(ctx, lang)
}

// SearchRoles returns all matching Roles from store
//
// This function is auto-generated
func SearchRoles(ctx context.Context, s Roles, f systemType.RoleFilter) (systemType.RoleSet, systemType.RoleFilter, error) {
	return s.SearchRoles(ctx, f)
}

// CreateRole creates one or more Roles in store
//
// This function is auto-generated
func CreateRole(ctx context.Context, s Roles, rr ...*systemType.Role) error {
	return s.CreateRole(ctx, rr...)
}

// UpdateRole updates one or more (existing) Roles in store
//
// This function is auto-generated
func UpdateRole(ctx context.Context, s Roles, rr ...*systemType.Role) error {
	return s.UpdateRole(ctx, rr...)
}

// UpsertRole creates new or updates existing one or more Roles in store
//
// This function is auto-generated
func UpsertRole(ctx context.Context, s Roles, rr ...*systemType.Role) error {
	return s.UpsertRole(ctx, rr...)
}

// DeleteRole deletes one or more Roles from store
//
// This function is auto-generated
func DeleteRole(ctx context.Context, s Roles, rr ...*systemType.Role) error {
	return s.DeleteRole(ctx, rr...)
}

// DeleteRoleByID deletes one or more Roles from store
//
// This function is auto-generated
func DeleteRoleByID(ctx context.Context, s Roles, id uint64) error {
	return s.DeleteRoleByID(ctx, id)
}

// TruncateRoles Deletes all Roles from store
//
// This function is auto-generated
func TruncateRoles(ctx context.Context, s Roles) error {
	return s.TruncateRoles(ctx)
}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
//
// This function is auto-generated
func LookupRoleByID(ctx context.Context, s Roles, id uint64) (*systemType.Role, error) {
	return s.LookupRoleByID(ctx, id)
}

// LookupRoleByHandle searches for role by handle
//
// It returns only valid role (not deleted, not suspended)
//
// This function is auto-generated
func LookupRoleByHandle(ctx context.Context, s Roles, handle string) (*systemType.Role, error) {
	return s.LookupRoleByHandle(ctx, handle)
}

// LookupRoleByName searches for role by name
//
// It returns only valid role (not deleted, not suspended)
//
// This function is auto-generated
func LookupRoleByName(ctx context.Context, s Roles, name string) (*systemType.Role, error) {
	return s.LookupRoleByName(ctx, name)
}

// RoleMetrics
//
// This function is auto-generated
func RoleMetrics(ctx context.Context, s Roles) (*systemType.RoleMetrics, error) {
	return s.RoleMetrics(ctx)
}

// SearchRoleMembers returns all matching RoleMembers from store
//
// This function is auto-generated
func SearchRoleMembers(ctx context.Context, s RoleMembers, f systemType.RoleMemberFilter) (systemType.RoleMemberSet, systemType.RoleMemberFilter, error) {
	return s.SearchRoleMembers(ctx, f)
}

// CreateRoleMember creates one or more RoleMembers in store
//
// This function is auto-generated
func CreateRoleMember(ctx context.Context, s RoleMembers, rr ...*systemType.RoleMember) error {
	return s.CreateRoleMember(ctx, rr...)
}

// UpdateRoleMember updates one or more (existing) RoleMembers in store
//
// This function is auto-generated
func UpdateRoleMember(ctx context.Context, s RoleMembers, rr ...*systemType.RoleMember) error {
	return s.UpdateRoleMember(ctx, rr...)
}

// UpsertRoleMember creates new or updates existing one or more RoleMembers in store
//
// This function is auto-generated
func UpsertRoleMember(ctx context.Context, s RoleMembers, rr ...*systemType.RoleMember) error {
	return s.UpsertRoleMember(ctx, rr...)
}

// DeleteRoleMember deletes one or more RoleMembers from store
//
// This function is auto-generated
func DeleteRoleMember(ctx context.Context, s RoleMembers, rr ...*systemType.RoleMember) error {
	return s.DeleteRoleMember(ctx, rr...)
}

// DeleteRoleMemberByID deletes one or more RoleMembers from store
//
// This function is auto-generated
func DeleteRoleMemberByUserIDRoleID(ctx context.Context, s RoleMembers, userID uint64, roleID uint64) error {
	return s.DeleteRoleMemberByUserIDRoleID(ctx, userID, roleID)
}

// TruncateRoleMembers Deletes all RoleMembers from store
//
// This function is auto-generated
func TruncateRoleMembers(ctx context.Context, s RoleMembers) error {
	return s.TruncateRoleMembers(ctx)
}

// TransferRoleMembers
//
// This function is auto-generated
func TransferRoleMembers(ctx context.Context, s RoleMembers, src uint64, dst uint64) error {
	return s.TransferRoleMembers(ctx, src, dst)
}

// SearchSettingValues returns all matching SettingValues from store
//
// This function is auto-generated
func SearchSettingValues(ctx context.Context, s SettingValues, f systemType.SettingsFilter) (systemType.SettingValueSet, systemType.SettingsFilter, error) {
	return s.SearchSettingValues(ctx, f)
}

// CreateSettingValue creates one or more SettingValues in store
//
// This function is auto-generated
func CreateSettingValue(ctx context.Context, s SettingValues, rr ...*systemType.SettingValue) error {
	return s.CreateSettingValue(ctx, rr...)
}

// UpdateSettingValue updates one or more (existing) SettingValues in store
//
// This function is auto-generated
func UpdateSettingValue(ctx context.Context, s SettingValues, rr ...*systemType.SettingValue) error {
	return s.UpdateSettingValue(ctx, rr...)
}

// UpsertSettingValue creates new or updates existing one or more SettingValues in store
//
// This function is auto-generated
func UpsertSettingValue(ctx context.Context, s SettingValues, rr ...*systemType.SettingValue) error {
	return s.UpsertSettingValue(ctx, rr...)
}

// DeleteSettingValue deletes one or more SettingValues from store
//
// This function is auto-generated
func DeleteSettingValue(ctx context.Context, s SettingValues, rr ...*systemType.SettingValue) error {
	return s.DeleteSettingValue(ctx, rr...)
}

// DeleteSettingValueByID deletes one or more SettingValues from store
//
// This function is auto-generated
func DeleteSettingValueByNameOwnedBy(ctx context.Context, s SettingValues, name string, ownedBy uint64) error {
	return s.DeleteSettingValueByNameOwnedBy(ctx, name, ownedBy)
}

// TruncateSettingValues Deletes all SettingValues from store
//
// This function is auto-generated
func TruncateSettingValues(ctx context.Context, s SettingValues) error {
	return s.TruncateSettingValues(ctx)
}

// LookupSettingValueByNameOwnedBy searches for settings by name and owner
//
// This function is auto-generated
func LookupSettingValueByNameOwnedBy(ctx context.Context, s SettingValues, name string, ownedBy uint64) (*systemType.SettingValue, error) {
	return s.LookupSettingValueByNameOwnedBy(ctx, name, ownedBy)
}

// SearchTemplates returns all matching Templates from store
//
// This function is auto-generated
func SearchTemplates(ctx context.Context, s Templates, f systemType.TemplateFilter) (systemType.TemplateSet, systemType.TemplateFilter, error) {
	return s.SearchTemplates(ctx, f)
}

// CreateTemplate creates one or more Templates in store
//
// This function is auto-generated
func CreateTemplate(ctx context.Context, s Templates, rr ...*systemType.Template) error {
	return s.CreateTemplate(ctx, rr...)
}

// UpdateTemplate updates one or more (existing) Templates in store
//
// This function is auto-generated
func UpdateTemplate(ctx context.Context, s Templates, rr ...*systemType.Template) error {
	return s.UpdateTemplate(ctx, rr...)
}

// UpsertTemplate creates new or updates existing one or more Templates in store
//
// This function is auto-generated
func UpsertTemplate(ctx context.Context, s Templates, rr ...*systemType.Template) error {
	return s.UpsertTemplate(ctx, rr...)
}

// DeleteTemplate deletes one or more Templates from store
//
// This function is auto-generated
func DeleteTemplate(ctx context.Context, s Templates, rr ...*systemType.Template) error {
	return s.DeleteTemplate(ctx, rr...)
}

// DeleteTemplateByID deletes one or more Templates from store
//
// This function is auto-generated
func DeleteTemplateByID(ctx context.Context, s Templates, id uint64) error {
	return s.DeleteTemplateByID(ctx, id)
}

// TruncateTemplates Deletes all Templates from store
//
// This function is auto-generated
func TruncateTemplates(ctx context.Context, s Templates) error {
	return s.TruncateTemplates(ctx)
}

// LookupTemplateByID searches for template by ID
//
// It also returns deleted templates.
//
// This function is auto-generated
func LookupTemplateByID(ctx context.Context, s Templates, id uint64) (*systemType.Template, error) {
	return s.LookupTemplateByID(ctx, id)
}

// LookupTemplateByHandle searches for template by handle
//
// It returns only valid templates (not deleted)
//
// This function is auto-generated
func LookupTemplateByHandle(ctx context.Context, s Templates, handle string) (*systemType.Template, error) {
	return s.LookupTemplateByHandle(ctx, handle)
}

// SearchUsers returns all matching Users from store
//
// This function is auto-generated
func SearchUsers(ctx context.Context, s Users, f systemType.UserFilter) (systemType.UserSet, systemType.UserFilter, error) {
	return s.SearchUsers(ctx, f)
}

// CreateUser creates one or more Users in store
//
// This function is auto-generated
func CreateUser(ctx context.Context, s Users, rr ...*systemType.User) error {
	return s.CreateUser(ctx, rr...)
}

// UpdateUser updates one or more (existing) Users in store
//
// This function is auto-generated
func UpdateUser(ctx context.Context, s Users, rr ...*systemType.User) error {
	return s.UpdateUser(ctx, rr...)
}

// UpsertUser creates new or updates existing one or more Users in store
//
// This function is auto-generated
func UpsertUser(ctx context.Context, s Users, rr ...*systemType.User) error {
	return s.UpsertUser(ctx, rr...)
}

// DeleteUser deletes one or more Users from store
//
// This function is auto-generated
func DeleteUser(ctx context.Context, s Users, rr ...*systemType.User) error {
	return s.DeleteUser(ctx, rr...)
}

// DeleteUserByID deletes one or more Users from store
//
// This function is auto-generated
func DeleteUserByID(ctx context.Context, s Users, id uint64) error {
	return s.DeleteUserByID(ctx, id)
}

// TruncateUsers Deletes all Users from store
//
// This function is auto-generated
func TruncateUsers(ctx context.Context, s Users) error {
	return s.TruncateUsers(ctx)
}

// LookupUserByID searches for user by ID
//
// It returns user even if deleted or suspended
//
// This function is auto-generated
func LookupUserByID(ctx context.Context, s Users, id uint64) (*systemType.User, error) {
	return s.LookupUserByID(ctx, id)
}

// LookupUserByEmail searches for user by email
//
// It returns only valid user (not deleted, not suspended)
//
// This function is auto-generated
func LookupUserByEmail(ctx context.Context, s Users, email string) (*systemType.User, error) {
	return s.LookupUserByEmail(ctx, email)
}

// LookupUserByHandle searches for user by handle
//
// It returns only valid user (not deleted, not suspended)
//
// This function is auto-generated
func LookupUserByHandle(ctx context.Context, s Users, handle string) (*systemType.User, error) {
	return s.LookupUserByHandle(ctx, handle)
}

// LookupUserByUsername searches for user by username
//
// It returns only valid user (not deleted, not suspended)
//
// This function is auto-generated
func LookupUserByUsername(ctx context.Context, s Users, username string) (*systemType.User, error) {
	return s.LookupUserByUsername(ctx, username)
}

// CountUsers
//
// This function is auto-generated
func CountUsers(ctx context.Context, s Users, u systemType.UserFilter) (uint, error) {
	return s.CountUsers(ctx, u)
}

// UserMetrics
//
// This function is auto-generated
func UserMetrics(ctx context.Context, s Users) (*systemType.UserMetrics, error) {
	return s.UserMetrics(ctx)
}
