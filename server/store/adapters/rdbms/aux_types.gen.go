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
	"github.com/crusttech/human/server/pkg/expr"
	flagType "github.com/crusttech/human/server/pkg/flag/types"
	labelsType "github.com/crusttech/human/server/pkg/label/types"
	rbacType "github.com/crusttech/human/server/pkg/rbac"
	systemType "github.com/crusttech/human/server/system/types"
	"time"
)

type (

	// auxActionlog is an auxiliary structure used for transporting to/from RDBMS store
	auxActionlog struct {
		ID                 uint64                 `db:"id"`
		Timestamp          time.Time              `db:"timestamp"`
		ActorIPAddr        string                 `db:"actor_ip_addr"`
		ActorID            uint64                 `db:"actor_id"`
		RequestOrigin      string                 `db:"request_origin"`
		RequestID          string                 `db:"request_id"`
		Resource           string                 `db:"resource"`
		TenantID           uint64                 `db:"tenant_id"`
		ProjectID          uint64                 `db:"project_id"`
		RootProjectID      uint64                 `db:"root_project_id"`
		ResourceProjectID  uint64                 `db:"resource_project_id"`
		ResourceRevisionID uint64                 `db:"resource_revision_id"`
		Action             string                 `db:"action"`
		Error              string                 `db:"error"`
		Severity           actionlogType.Severity `db:"severity"`
		Description        string                 `db:"description"`
		Meta               actionlogType.Meta     `db:"meta"`
		Delta              actionlogType.Delta    `db:"delta"`
		OldState           actionlogType.OldState `db:"old_state"`
	}

	// auxAgent is an auxiliary structure used for transporting to/from RDBMS store
	auxAgent struct {
		ID         uint64                     `db:"id"`
		TenantID   uint64                     `db:"tenant_id"`
		ProjectID  uint64                     `db:"project_id"`
		Handle     string                     `db:"handle"`
		Status     string                     `db:"status"`
		Revision   int                        `db:"revision"`
		Meta       systemType.AgentMeta       `db:"meta"`
		Behavior   systemType.AgentBehavior   `db:"behavior"`
		Execution  systemType.AgentExecution  `db:"execution"`
		Access     systemType.AgentAccess     `db:"access"`
		Invocation systemType.AgentInvocation `db:"invocation"`
		CreatedAt  time.Time                  `db:"created_at"`
		UpdatedAt  *time.Time                 `db:"updated_at"`
		DeletedAt  *time.Time                 `db:"deleted_at"`
		CreatedBy  uint64                     `db:"created_by"`
		UpdatedBy  uint64                     `db:"updated_by"`
		DeletedBy  uint64                     `db:"deleted_by"`
	}

	// auxAiConversation is an auxiliary structure used for transporting to/from RDBMS store
	auxAiConversation struct {
		ID         uint64                            `db:"id"`
		TenantID   uint64                            `db:"tenant_id"`
		ProjectID  uint64                            `db:"project_id"`
		AgentID    uint64                            `db:"agentID"`
		Messages   systemType.AiConversationMessages `db:"messages"`
		TokenCount int                               `db:"tokenCount"`
		CreatedAt  time.Time                         `db:"created_at"`
		UpdatedAt  *time.Time                        `db:"updated_at"`
		DeletedAt  *time.Time                        `db:"deleted_at"`
		CreatedBy  uint64                            `db:"created_by"`
		UpdatedBy  uint64                            `db:"updated_by"`
		DeletedBy  uint64                            `db:"deleted_by"`
	}

	// auxApigwFilter is an auxiliary structure used for transporting to/from RDBMS store
	auxApigwFilter struct {
		ID        uint64                       `db:"id"`
		TenantID  uint64                       `db:"tenant_id"`
		ProjectID uint64                       `db:"project_id"`
		Route     uint64                       `db:"route"`
		Weight    uint64                       `db:"weight"`
		Kind      string                       `db:"kind"`
		Ref       string                       `db:"ref"`
		Enabled   bool                         `db:"enabled"`
		Params    systemType.ApigwFilterParams `db:"params"`
		CreatedAt time.Time                    `db:"created_at"`
		UpdatedAt *time.Time                   `db:"updated_at"`
		DeletedAt *time.Time                   `db:"deleted_at"`
		CreatedBy uint64                       `db:"created_by"`
		UpdatedBy uint64                       `db:"updated_by"`
		DeletedBy uint64                       `db:"deleted_by"`
	}

	// auxApigwRoute is an auxiliary structure used for transporting to/from RDBMS store
	auxApigwRoute struct {
		ID        uint64                    `db:"id"`
		TenantID  uint64                    `db:"tenant_id"`
		ProjectID uint64                    `db:"project_id"`
		Endpoint  string                    `db:"endpoint"`
		Method    string                    `db:"method"`
		Enabled   bool                      `db:"enabled"`
		Meta      systemType.ApigwRouteMeta `db:"meta"`
		Group     uint64                    `db:"group"`
		CreatedAt time.Time                 `db:"created_at"`
		UpdatedAt *time.Time                `db:"updated_at"`
		DeletedAt *time.Time                `db:"deleted_at"`
		CreatedBy uint64                    `db:"created_by"`
		UpdatedBy uint64                    `db:"updated_by"`
		DeletedBy uint64                    `db:"deleted_by"`
	}

	// auxApplication is an auxiliary structure used for transporting to/from RDBMS store
	auxApplication struct {
		ID        uint64                       `db:"id"`
		TenantID  uint64                       `db:"tenant_id"`
		ProjectID uint64                       `db:"project_id"`
		Name      string                       `db:"name"`
		Enabled   bool                         `db:"enabled"`
		Weight    int                          `db:"weight"`
		Unify     *systemType.ApplicationUnify `db:"unify"`
		OwnerID   uint64                       `db:"owner_id"`
		CreatedAt time.Time                    `db:"created_at"`
		UpdatedAt *time.Time                   `db:"updated_at"`
		DeletedAt *time.Time                   `db:"deleted_at"`
	}

	// auxAttachment is an auxiliary structure used for transporting to/from RDBMS store
	auxAttachment struct {
		ID         uint64                    `db:"id"`
		TenantID   uint64                    `db:"tenant_id"`
		ProjectID  uint64                    `db:"project_id"`
		OwnerID    uint64                    `db:"owner_id"`
		Kind       string                    `db:"kind"`
		Url        string                    `db:"url"`
		PreviewUrl string                    `db:"preview_url"`
		Name       string                    `db:"name"`
		Meta       systemType.AttachmentMeta `db:"meta"`
		CreatedAt  time.Time                 `db:"created_at"`
		UpdatedAt  *time.Time                `db:"updated_at"`
		DeletedAt  *time.Time                `db:"deleted_at"`
	}

	// auxAuthClient is an auxiliary structure used for transporting to/from RDBMS store
	auxAuthClient struct {
		ID          uint64                         `db:"id"`
		TenantID    uint64                         `db:"tenant_id"`
		ProjectID   uint64                         `db:"project_id"`
		Handle      string                         `db:"handle"`
		Meta        *systemType.AuthClientMeta     `db:"meta"`
		Secret      string                         `db:"secret"`
		Scope       string                         `db:"scope"`
		ValidGrant  string                         `db:"valid_grant"`
		RedirectURI string                         `db:"redirect_uri"`
		Enabled     bool                           `db:"enabled"`
		Trusted     bool                           `db:"trusted"`
		ValidFrom   *time.Time                     `db:"valid_from"`
		ExpiresAt   *time.Time                     `db:"expires_at"`
		Security    *systemType.AuthClientSecurity `db:"security"`
		OwnedBy     uint64                         `db:"owned_by"`
		CreatedAt   time.Time                      `db:"created_at"`
		UpdatedAt   *time.Time                     `db:"updated_at"`
		DeletedAt   *time.Time                     `db:"deleted_at"`
		CreatedBy   uint64                         `db:"created_by"`
		UpdatedBy   uint64                         `db:"updated_by"`
		DeletedBy   uint64                         `db:"deleted_by"`
	}

	// auxAuthConfirmedClient is an auxiliary structure used for transporting to/from RDBMS store
	auxAuthConfirmedClient struct {
		UserID      uint64    `db:"user_id"`
		ClientID    uint64    `db:"client_id"`
		ConfirmedAt time.Time `db:"confirmed_at"`
	}

	// auxAuthOa2token is an auxiliary structure used for transporting to/from RDBMS store
	auxAuthOa2token struct {
		ID         uint64    `db:"id"`
		Code       string    `db:"code"`
		Access     string    `db:"access"`
		Refresh    string    `db:"refresh"`
		Data       rawJson   `db:"data"`
		RemoteAddr string    `db:"remote_addr"`
		UserAgent  string    `db:"user_agent"`
		ClientID   uint64    `db:"client_id"`
		UserID     uint64    `db:"user_id"`
		CreatedAt  time.Time `db:"created_at"`
		ExpiresAt  time.Time `db:"expires_at"`
	}

	// auxAuthSession is an auxiliary structure used for transporting to/from RDBMS store
	auxAuthSession struct {
		ID         string    `db:"id"`
		Data       []byte    `db:"data"`
		UserID     uint64    `db:"user_id"`
		RemoteAddr string    `db:"remote_addr"`
		UserAgent  string    `db:"user_agent"`
		ExpiresAt  time.Time `db:"expires_at"`
		CreatedAt  time.Time `db:"created_at"`
	}

	// auxAutomationNgAutomation is an auxiliary structure used for transporting to/from RDBMS store
	auxAutomationNgAutomation struct {
		ID        uint64                                `db:"id"`
		TenantID  uint64                                `db:"tenant_id"`
		ProjectID uint64                                `db:"project_id"`
		Handle    string                                `db:"handle"`
		Meta      *automationType.NgAutomationMeta      `db:"meta"`
		Enabled   bool                                  `db:"enabled"`
		Scope     *expr.Vars                            `db:"scope"`
		Triggers  automationType.NgAutomationTriggerSet `db:"triggers"`
		Steps     automationType.NgAutomationStepSet    `db:"steps"`
		Paths     automationType.NgAutomationPathSet    `db:"paths"`
		Issues    automationType.NgAutomationIssueSet   `db:"issues"`
		RunAs     uint64                                `db:"run_as"`
		OwnedBy   uint64                                `db:"owned_by"`
		CreatedAt time.Time                             `db:"created_at"`
		UpdatedAt *time.Time                            `db:"updated_at"`
		DeletedAt *time.Time                            `db:"deleted_at"`
		CreatedBy uint64                                `db:"created_by"`
		UpdatedBy uint64                                `db:"updated_by"`
		DeletedBy uint64                                `db:"deleted_by"`
	}

	// auxAutomationSession is an auxiliary structure used for transporting to/from RDBMS store
	auxAutomationSession struct {
		ID           uint64                       `db:"id"`
		TenantID     uint64                       `db:"tenant_id"`
		ProjectID    uint64                       `db:"project_id"`
		WorkflowID   uint64                       `db:"workflow_id"`
		Status       automationType.SessionStatus `db:"status"`
		EventType    string                       `db:"event_type"`
		ResourceType string                       `db:"resource_type"`
		Input        *expr.Vars                   `db:"input"`
		Output       *expr.Vars                   `db:"output"`
		Stacktrace   automationType.Stacktrace    `db:"stacktrace"`
		CreatedBy    uint64                       `db:"created_by"`
		CreatedAt    time.Time                    `db:"created_at"`
		PurgeAt      *time.Time                   `db:"purge_at"`
		SuspendedAt  *time.Time                   `db:"suspended_at"`
		CompletedAt  *time.Time                   `db:"completed_at"`
		Error        string                       `db:"error"`
	}

	// auxAutomationTrigger is an auxiliary structure used for transporting to/from RDBMS store
	auxAutomationTrigger struct {
		ID           uint64                              `db:"id"`
		TenantID     uint64                              `db:"tenant_id"`
		ProjectID    uint64                              `db:"project_id"`
		WorkflowID   uint64                              `db:"workflow_id"`
		StepID       uint64                              `db:"step_id"`
		Enabled      bool                                `db:"enabled"`
		Meta         *automationType.TriggerMeta         `db:"meta"`
		ResourceType string                              `db:"resource_type"`
		EventType    string                              `db:"event_type"`
		Constraints  automationType.TriggerConstraintSet `db:"constraints"`
		Input        *expr.Vars                          `db:"input"`
		OwnedBy      uint64                              `db:"owned_by"`
		CreatedAt    time.Time                           `db:"created_at"`
		UpdatedAt    *time.Time                          `db:"updated_at"`
		DeletedAt    *time.Time                          `db:"deleted_at"`
		CreatedBy    uint64                              `db:"created_by"`
		UpdatedBy    uint64                              `db:"updated_by"`
		DeletedBy    uint64                              `db:"deleted_by"`
	}

	// auxAutomationWorkflow is an auxiliary structure used for transporting to/from RDBMS store
	auxAutomationWorkflow struct {
		ID           uint64                          `db:"id"`
		TenantID     uint64                          `db:"tenant_id"`
		ProjectID    uint64                          `db:"project_id"`
		Handle       string                          `db:"handle"`
		Meta         *automationType.WorkflowMeta    `db:"meta"`
		Enabled      bool                            `db:"enabled"`
		Trace        bool                            `db:"trace"`
		KeepSessions int                             `db:"keep_sessions"`
		Scope        *expr.Vars                      `db:"scope"`
		Steps        automationType.WorkflowStepSet  `db:"steps"`
		Paths        automationType.WorkflowPathSet  `db:"paths"`
		Issues       automationType.WorkflowIssueSet `db:"issues"`
		RunAs        uint64                          `db:"run_as"`
		OwnedBy      uint64                          `db:"owned_by"`
		CreatedAt    time.Time                       `db:"created_at"`
		UpdatedAt    *time.Time                      `db:"updated_at"`
		DeletedAt    *time.Time                      `db:"deleted_at"`
		CreatedBy    uint64                          `db:"created_by"`
		UpdatedBy    uint64                          `db:"updated_by"`
		DeletedBy    uint64                          `db:"deleted_by"`
	}

	// auxChatbot is an auxiliary structure used for transporting to/from RDBMS store
	auxChatbot struct {
		ID             uint64                           `db:"id"`
		TenantID       uint64                           `db:"tenant_id"`
		ProjectID      uint64                           `db:"project_id"`
		Handle         string                           `db:"handle"`
		Name           string                           `db:"name"`
		Enabled        bool                             `db:"enabled"`
		WidgetKey      string                           `db:"widget_key"`
		AllowedOrigins systemType.ChatbotAllowedOrigins `db:"allowed_origins"`
		SessionTTL     string                           `db:"session_ttl"`
		Handoff        systemType.ChatbotHandoff        `db:"handoff"`
		Styling        systemType.ChatbotStyling        `db:"styling"`
		Scenarios      systemType.ChatbotScenarios      `db:"scenarios"`
		CreatedAt      time.Time                        `db:"created_at"`
		UpdatedAt      *time.Time                       `db:"updated_at"`
		DeletedAt      *time.Time                       `db:"deleted_at"`
		CreatedBy      uint64                           `db:"created_by"`
		UpdatedBy      uint64                           `db:"updated_by"`
		DeletedBy      uint64                           `db:"deleted_by"`
	}

	// auxChatbotSession is an auxiliary structure used for transporting to/from RDBMS store
	auxChatbotSession struct {
		ID          uint64                         `db:"id"`
		TenantID    uint64                         `db:"tenant_id"`
		ProjectID   uint64                         `db:"project_id"`
		ChatbotID   uint64                         `db:"chatbot_id"`
		Status      string                         `db:"status"`
		CurrentStep int                            `db:"current_step"`
		State       systemType.ChatbotSessionState `db:"state"`
		CreatedAt   time.Time                      `db:"created_at"`
		UpdatedAt   *time.Time                     `db:"updated_at"`
		DeletedAt   *time.Time                     `db:"deleted_at"`
		CreatedBy   uint64                         `db:"created_by"`
		UpdatedBy   uint64                         `db:"updated_by"`
		DeletedBy   uint64                         `db:"deleted_by"`
	}

	// auxChatbotSessionHandoff is an auxiliary structure used for transporting to/from RDBMS store
	auxChatbotSessionHandoff struct {
		ID          uint64     `db:"id"`
		SessionID   uint64     `db:"session_id"`
		StepID      uint64     `db:"step_id"`
		Status      string     `db:"status"`
		InitiatedAt time.Time  `db:"initiated_at"`
		ClosedAt    *time.Time `db:"closed_at"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
		CreatedBy   uint64     `db:"created_by"`
		UpdatedBy   uint64     `db:"updated_by"`
		DeletedBy   uint64     `db:"deleted_by"`
	}

	// auxChatbotSessionStep is an auxiliary structure used for transporting to/from RDBMS store
	auxChatbotSessionStep struct {
		ID             uint64     `db:"id"`
		SessionID      uint64     `db:"session_id"`
		ScenarioIndex  int        `db:"scenario_index"`
		ConversationID uint64     `db:"conversation_id"`
		Status         string     `db:"status"`
		CreatedAt      time.Time  `db:"created_at"`
		UpdatedAt      *time.Time `db:"updated_at"`
		DeletedAt      *time.Time `db:"deleted_at"`
		CreatedBy      uint64     `db:"created_by"`
		UpdatedBy      uint64     `db:"updated_by"`
		DeletedBy      uint64     `db:"deleted_by"`
	}

	// auxComposeAttachment is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeAttachment struct {
		ID          uint64                     `db:"id"`
		TenantID    uint64                     `db:"tenant_id"`
		ProjectID   uint64                     `db:"project_id"`
		NamespaceID uint64                     `db:"namespace_id"`
		OwnerID     uint64                     `db:"owner_id"`
		Kind        string                     `db:"kind"`
		Url         string                     `db:"url"`
		PreviewUrl  string                     `db:"preview_url"`
		Name        string                     `db:"name"`
		Meta        composeType.AttachmentMeta `db:"meta"`
		CreatedAt   time.Time                  `db:"created_at"`
		UpdatedAt   *time.Time                 `db:"updated_at"`
		DeletedAt   *time.Time                 `db:"deleted_at"`
	}

	// auxComposeChart is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeChart struct {
		ID          uint64                  `db:"id"`
		Handle      string                  `db:"handle"`
		TenantID    uint64                  `db:"tenant_id"`
		ProjectID   uint64                  `db:"project_id"`
		NamespaceID uint64                  `db:"namespace_id"`
		Name        string                  `db:"name"`
		Config      composeType.ChartConfig `db:"config"`
		CreatedAt   time.Time               `db:"created_at"`
		UpdatedAt   *time.Time              `db:"updated_at"`
		DeletedAt   *time.Time              `db:"deleted_at"`
	}

	// auxComposeModule is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeModule struct {
		ID             uint64                   `db:"id"`
		TenantID       uint64                   `db:"tenant_id"`
		ProjectID      uint64                   `db:"project_id"`
		NamespaceID    uint64                   `db:"namespace_id"`
		Handle         string                   `db:"handle"`
		Name           string                   `db:"name"`
		Meta           rawJson                  `db:"meta"`
		Config         composeType.ModuleConfig `db:"config"`
		CreatedAt      time.Time                `db:"created_at"`
		UpdatedAt      *time.Time               `db:"updated_at"`
		DeletedAt      *time.Time               `db:"deleted_at"`
		CreatedByAgent uint64                   `db:"created_by_agent"`
	}

	// auxComposeModuleField is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeModuleField struct {
		ID             uint64                         `db:"id"`
		TenantID       uint64                         `db:"tenant_id"`
		ProjectID      uint64                         `db:"project_id"`
		ModuleID       uint64                         `db:"module_id"`
		Place          int                            `db:"place"`
		Kind           string                         `db:"kind"`
		Options        composeType.ModuleFieldOptions `db:"options"`
		Name           string                         `db:"name"`
		Label          string                         `db:"label"`
		Config         composeType.ModuleFieldConfig  `db:"config"`
		Required       bool                           `db:"required"`
		Multi          bool                           `db:"multi"`
		DefaultValue   composeType.RecordValueSet     `db:"default_value"`
		Expressions    composeType.ModuleFieldExpr    `db:"expressions"`
		CreatedAt      time.Time                      `db:"created_at"`
		UpdatedAt      *time.Time                     `db:"updated_at"`
		DeletedAt      *time.Time                     `db:"deleted_at"`
		CreatedByAgent uint64                         `db:"created_by_agent"`
	}

	// auxComposeNamespace is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeNamespace struct {
		ID             uint64                    `db:"id"`
		TenantID       uint64                    `db:"tenant_id"`
		ProjectID      uint64                    `db:"project_id"`
		Slug           string                    `db:"slug"`
		Enabled        bool                      `db:"enabled"`
		Meta           composeType.NamespaceMeta `db:"meta"`
		Name           string                    `db:"name"`
		CreatedAt      time.Time                 `db:"created_at"`
		UpdatedAt      *time.Time                `db:"updated_at"`
		DeletedAt      *time.Time                `db:"deleted_at"`
		CreatedByAgent uint64                    `db:"created_by_agent"`
	}

	// auxComposePage is an auxiliary structure used for transporting to/from RDBMS store
	auxComposePage struct {
		ID             uint64                 `db:"id"`
		TenantID       uint64                 `db:"tenant_id"`
		ProjectID      uint64                 `db:"project_id"`
		Title          string                 `db:"title"`
		Handle         string                 `db:"handle"`
		SelfID         uint64                 `db:"self_id"`
		ModuleID       uint64                 `db:"module_id"`
		NamespaceID    uint64                 `db:"namespace_id"`
		Meta           composeType.PageMeta   `db:"meta"`
		Config         composeType.PageConfig `db:"config"`
		Blocks         composeType.PageBlocks `db:"blocks"`
		Visible        bool                   `db:"visible"`
		Weight         int                    `db:"weight"`
		Description    string                 `db:"description"`
		CreatedAt      time.Time              `db:"created_at"`
		UpdatedAt      *time.Time             `db:"updated_at"`
		DeletedAt      *time.Time             `db:"deleted_at"`
		CreatedByAgent uint64                 `db:"created_by_agent"`
	}

	// auxComposePageLayout is an auxiliary structure used for transporting to/from RDBMS store
	auxComposePageLayout struct {
		ID             uint64                       `db:"id"`
		TenantID       uint64                       `db:"tenant_id"`
		ProjectID      uint64                       `db:"project_id"`
		Handle         string                       `db:"handle"`
		PageID         uint64                       `db:"page_id"`
		ParentID       uint64                       `db:"parent_id"`
		NamespaceID    uint64                       `db:"namespace_id"`
		Weight         int                          `db:"weight"`
		Meta           composeType.PageLayoutMeta   `db:"meta"`
		Config         composeType.PageLayoutConfig `db:"config"`
		Blocks         composeType.PageLayoutBlocks `db:"blocks"`
		OwnedBy        uint64                       `db:"owned_by"`
		CreatedAt      time.Time                    `db:"created_at"`
		UpdatedAt      *time.Time                   `db:"updated_at"`
		DeletedAt      *time.Time                   `db:"deleted_at"`
		CreatedByAgent uint64                       `db:"created_by_agent"`
	}

	// auxConfiguredConnection is an auxiliary structure used for transporting to/from RDBMS store
	auxConfiguredConnection struct {
		ID           uint64                                `db:"id"`
		TenantID     uint64                                `db:"tenant_id"`
		ProjectID    uint64                                `db:"project_id"`
		ConnectionID uint64                                `db:"connection_id"`
		Name         string                                `db:"name"`
		Status       string                                `db:"status"`
		Connection   systemType.Connection                 `db:"connection"`
		Config       systemType.ConfiguredConnectionConfig `db:"config"`
		CreatedAt    time.Time                             `db:"created_at"`
		UpdatedAt    *time.Time                            `db:"updated_at"`
		DeletedAt    *time.Time                            `db:"deleted_at"`
		CreatedBy    uint64                                `db:"created_by"`
		UpdatedBy    uint64                                `db:"updated_by"`
		DeletedBy    uint64                                `db:"deleted_by"`
	}

	// auxConnection is an auxiliary structure used for transporting to/from RDBMS store
	auxConnection struct {
		ID         uint64                          `db:"id"`
		Handle     string                          `db:"handle"`
		Revision   int                             `db:"revision"`
		Status     string                          `db:"status"`
		Source     string                          `db:"source"`
		Meta       systemType.ConnectionMeta       `db:"meta"`
		Service    systemType.ConnectionService    `db:"service"`
		Resources  systemType.ConnectionResources  `db:"resources"`
		Operations systemType.ConnectionOperations `db:"operations"`
		CreatedAt  time.Time                       `db:"created_at"`
		UpdatedAt  *time.Time                      `db:"updated_at"`
		DeletedAt  *time.Time                      `db:"deleted_at"`
		CreatedBy  uint64                          `db:"created_by"`
		UpdatedBy  uint64                          `db:"updated_by"`
		DeletedBy  uint64                          `db:"deleted_by"`
	}

	// auxCredential is an auxiliary structure used for transporting to/from RDBMS store
	auxCredential struct {
		ID          uint64     `db:"id"`
		OwnerID     uint64     `db:"owner_id"`
		Label       string     `db:"label"`
		Kind        string     `db:"kind"`
		Credentials string     `db:"credentials"`
		Meta        rawJson    `db:"meta"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
		LastUsedAt  *time.Time `db:"last_used_at"`
		ExpiresAt   *time.Time `db:"expires_at"`
	}

	// auxDalConnection is an auxiliary structure used for transporting to/from RDBMS store
	auxDalConnection struct {
		ID        uint64                         `db:"id"`
		Handle    string                         `db:"handle"`
		Type      string                         `db:"type"`
		Config    systemType.DalConnectionConfig `db:"config"`
		Meta      systemType.DalConnectionMeta   `db:"meta"`
		CreatedAt time.Time                      `db:"created_at"`
		UpdatedAt *time.Time                     `db:"updated_at"`
		DeletedAt *time.Time                     `db:"deleted_at"`
		CreatedBy uint64                         `db:"created_by"`
		UpdatedBy uint64                         `db:"updated_by"`
		DeletedBy uint64                         `db:"deleted_by"`
	}

	// auxDalSchemaAlteration is an auxiliary structure used for transporting to/from RDBMS store
	auxDalSchemaAlteration struct {
		ID           uint64                                `db:"id"`
		BatchID      uint64                                `db:"batchID"`
		DependsOn    uint64                                `db:"dependsOn"`
		Resource     string                                `db:"resource"`
		ResourceType string                                `db:"resourceType"`
		ConnectionID uint64                                `db:"connectionID"`
		Kind         string                                `db:"kind"`
		Params       *systemType.DalSchemaAlterationParams `db:"params"`
		Error        string                                `db:"error"`
		CreatedAt    time.Time                             `db:"created_at"`
		UpdatedAt    *time.Time                            `db:"updated_at"`
		DeletedAt    *time.Time                            `db:"deleted_at"`
		CompletedAt  *time.Time                            `db:"completed_at"`
		DismissedAt  *time.Time                            `db:"dismissed_at"`
		CreatedBy    uint64                                `db:"created_by"`
		UpdatedBy    uint64                                `db:"updated_by"`
		DeletedBy    uint64                                `db:"deleted_by"`
		CompletedBy  uint64                                `db:"completed_by"`
		DismissedBy  uint64                                `db:"dismissed_by"`
	}

	// auxDalSensitivityLevel is an auxiliary structure used for transporting to/from RDBMS store
	auxDalSensitivityLevel struct {
		ID        uint64                             `db:"id"`
		TenantID  uint64                             `db:"tenant_id"`
		ProjectID uint64                             `db:"project_id"`
		Handle    string                             `db:"handle"`
		Level     int                                `db:"level"`
		Meta      systemType.DalSensitivityLevelMeta `db:"meta"`
		CreatedAt time.Time                          `db:"created_at"`
		UpdatedAt *time.Time                         `db:"updated_at"`
		DeletedAt *time.Time                         `db:"deleted_at"`
		CreatedBy uint64                             `db:"created_by"`
		UpdatedBy uint64                             `db:"updated_by"`
		DeletedBy uint64                             `db:"deleted_by"`
	}

	// auxDataPrivacyRequest is an auxiliary structure used for transporting to/from RDBMS store
	auxDataPrivacyRequest struct {
		ID          uint64                                  `db:"id"`
		TenantID    uint64                                  `db:"tenant_id"`
		ProjectID   uint64                                  `db:"project_id"`
		Kind        systemType.RequestKind                  `db:"kind"`
		Status      systemType.RequestStatus                `db:"status"`
		Payload     systemType.DataPrivacyRequestPayloadSet `db:"payload"`
		RequestedAt time.Time                               `db:"requested_at"`
		RequestedBy uint64                                  `db:"requested_by"`
		CompletedAt *time.Time                              `db:"completed_at"`
		CompletedBy uint64                                  `db:"completed_by"`
		CreatedAt   time.Time                               `db:"created_at"`
		UpdatedAt   *time.Time                              `db:"updated_at"`
		DeletedAt   *time.Time                              `db:"deleted_at"`
		CreatedBy   uint64                                  `db:"created_by"`
		UpdatedBy   uint64                                  `db:"updated_by"`
		DeletedBy   uint64                                  `db:"deleted_by"`
	}

	// auxDataPrivacyRequestComment is an auxiliary structure used for transporting to/from RDBMS store
	auxDataPrivacyRequestComment struct {
		ID        uint64     `db:"id"`
		RequestID uint64     `db:"request_id"`
		Comment   string     `db:"comment"`
		CreatedAt time.Time  `db:"created_at"`
		UpdatedAt *time.Time `db:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at"`
		CreatedBy uint64     `db:"created_by"`
		UpdatedBy uint64     `db:"updated_by"`
		DeletedBy uint64     `db:"deleted_by"`
	}

	// auxDmlConnection is an auxiliary structure used for transporting to/from RDBMS store
	auxDmlConnection struct {
		ID        uint64                         `db:"id"`
		Handle    string                         `db:"handle"`
		Label     string                         `db:"label"`
		Params    systemType.DmlConnectionParams `db:"params"`
		CreatedAt time.Time                      `db:"created_at"`
		UpdatedAt *time.Time                     `db:"updated_at"`
		DeletedAt *time.Time                     `db:"deleted_at"`
	}

	// auxDmlImportRun is an auxiliary structure used for transporting to/from RDBMS store
	auxDmlImportRun struct {
		ID           uint64                     `db:"id"`
		ConnectionID uint64                     `db:"connection_id"`
		MappingID    uint64                     `db:"mapping_id"`
		Method       systemType.DmlImportMethod `db:"method"`
		Status       string                     `db:"status"`
		Processed    uint64                     `db:"processed"`
		Failed       uint64                     `db:"failed"`
		Error        string                     `db:"error"`
	}

	// auxDmlMapping is an auxiliary structure used for transporting to/from RDBMS store
	auxDmlMapping struct {
		ID              uint64                     `db:"id"`
		ConnectionID    uint64                     `db:"connection_id"`
		NamespaceHandle string                     `db:"namespace_handle"`
		SourceIdent     string                     `db:"source_ident"`
		ModuleHandle    string                     `db:"module_handle"`
		ModuleName      string                     `db:"module_name"`
		Skip            bool                       `db:"skip"`
		Identifier      string                     `db:"identifier"`
		Columns         systemType.DmlColumnMapSet `db:"columns"`
		CreatedAt       time.Time                  `db:"created_at"`
		UpdatedAt       *time.Time                 `db:"updated_at"`
		DeletedAt       *time.Time                 `db:"deleted_at"`
	}

	// auxFederationExposedModule is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationExposedModule struct {
		ID                 uint64                        `db:"id"`
		TenantID           uint64                        `db:"tenant_id"`
		ProjectID          uint64                        `db:"project_id"`
		Handle             string                        `db:"handle"`
		Name               string                        `db:"name"`
		NodeID             uint64                        `db:"node_id"`
		ComposeModuleID    uint64                        `db:"compose_module_id"`
		ComposeNamespaceID uint64                        `db:"compose_namespace_id"`
		Fields             federationType.ModuleFieldSet `db:"fields"`
		CreatedAt          time.Time                     `db:"created_at"`
		UpdatedAt          *time.Time                    `db:"updated_at"`
		DeletedAt          *time.Time                    `db:"deleted_at"`
		CreatedBy          uint64                        `db:"created_by"`
		UpdatedBy          uint64                        `db:"updated_by"`
		DeletedBy          uint64                        `db:"deleted_by"`
	}

	// auxFederationModuleMapping is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationModuleMapping struct {
		TenantID           uint64                               `db:"tenant_id"`
		ProjectID          uint64                               `db:"project_id"`
		NodeID             uint64                               `db:"node_id"`
		FederationModuleID uint64                               `db:"federation_module_id"`
		ComposeModuleID    uint64                               `db:"compose_module_id"`
		ComposeNamespaceID uint64                               `db:"compose_namespace_id"`
		FieldMapping       federationType.ModuleFieldMappingSet `db:"field_mapping"`
	}

	// auxFederationNode is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationNode struct {
		ID           uint64     `db:"id"`
		TenantID     uint64     `db:"tenant_id"`
		SharedNodeID uint64     `db:"shared_node_id"`
		Name         string     `db:"name"`
		BaseURL      string     `db:"base_url"`
		Status       string     `db:"status"`
		Contact      string     `db:"contact"`
		PairToken    string     `db:"pair_token"`
		AuthToken    string     `db:"auth_token"`
		CreatedAt    time.Time  `db:"created_at"`
		UpdatedAt    *time.Time `db:"updated_at"`
		DeletedAt    *time.Time `db:"deleted_at"`
		CreatedBy    uint64     `db:"created_by"`
		UpdatedBy    uint64     `db:"updated_by"`
		DeletedBy    uint64     `db:"deleted_by"`
	}

	// auxFederationNodeSync is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationNodeSync struct {
		NodeID       uint64    `db:"rel_node"`
		ModuleID     uint64    `db:"rel_module"`
		SyncType     string    `db:"sync_type"`
		SyncStatus   string    `db:"sync_status"`
		TimeOfAction time.Time `db:"time_of_action"`
	}

	// auxFederationSharedModule is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationSharedModule struct {
		ID                         uint64                        `db:"id"`
		TenantID                   uint64                        `db:"tenant_id"`
		ProjectID                  uint64                        `db:"project_id"`
		Handle                     string                        `db:"handle"`
		NodeID                     uint64                        `db:"node_id"`
		Name                       string                        `db:"name"`
		ExternalFederationModuleID uint64                        `db:"external_federation_module_id"`
		Fields                     federationType.ModuleFieldSet `db:"fields"`
		CreatedAt                  time.Time                     `db:"created_at"`
		UpdatedAt                  *time.Time                    `db:"updated_at"`
		DeletedAt                  *time.Time                    `db:"deleted_at"`
		CreatedBy                  uint64                        `db:"created_by"`
		UpdatedBy                  uint64                        `db:"updated_by"`
		DeletedBy                  uint64                        `db:"deleted_by"`
	}

	// auxFlag is an auxiliary structure used for transporting to/from RDBMS store
	auxFlag struct {
		Kind       string `db:"kind"`
		ResourceID uint64 `db:"resource_id"`
		OwnedBy    uint64 `db:"owned_by"`
		Name       string `db:"name"`
		Active     bool   `db:"active"`
	}

	// auxKnowledgeBase is an auxiliary structure used for transporting to/from RDBMS store
	auxKnowledgeBase struct {
		ID          uint64                           `db:"id"`
		TenantID    uint64                           `db:"tenant_id"`
		ProjectID   uint64                           `db:"project_id"`
		Handle      string                           `db:"handle"`
		Title       string                           `db:"title"`
		Description string                           `db:"description"`
		Context     *systemType.KnowledgeBaseContext `db:"context"`
		CreatedAt   time.Time                        `db:"created_at"`
		UpdatedAt   *time.Time                       `db:"updated_at"`
		DeletedAt   *time.Time                       `db:"deleted_at"`
		CreatedBy   uint64                           `db:"created_by"`
		UpdatedBy   uint64                           `db:"updated_by"`
		DeletedBy   uint64                           `db:"deleted_by"`
	}

	// auxLabel is an auxiliary structure used for transporting to/from RDBMS store
	auxLabel struct {
		Kind       string                `db:"kind"`
		ResourceID uint64                `db:"resource_id"`
		Name       string                `db:"name"`
		Value      labelsType.LabelValue `db:"value"`
	}

	// auxLlmProvider is an auxiliary structure used for transporting to/from RDBMS store
	auxLlmProvider struct {
		ID           uint64                       `db:"id"`
		TenantID     uint64                       `db:"tenant_id"`
		ProjectID    uint64                       `db:"project_id"`
		Handle       string                       `db:"handle"`
		Status       string                       `db:"status"`
		Provider     string                       `db:"provider"`
		CredentialID uint64                       `db:"credential_id"`
		Meta         systemType.LLMProviderMeta   `db:"meta"`
		Config       systemType.LLMProviderConfig `db:"config"`
		CreatedAt    time.Time                    `db:"created_at"`
		UpdatedAt    *time.Time                   `db:"updated_at"`
		DeletedAt    *time.Time                   `db:"deleted_at"`
		CreatedBy    uint64                       `db:"created_by"`
		UpdatedBy    uint64                       `db:"updated_by"`
		DeletedBy    uint64                       `db:"deleted_by"`
	}

	// auxNotification is an auxiliary structure used for transporting to/from RDBMS store
	auxNotification struct {
		ID        uint64                        `db:"id"`
		TenantID  uint64                        `db:"tenant_id"`
		ProjectID uint64                        `db:"project_id"`
		Kind      systemType.NotificationKind   `db:"kind"`
		Config    systemType.NotificationConfig `db:"config"`
		Recipient uint64                        `db:"recipient"`
		CreatedBy uint64                        `db:"created_by"`
		ReadAt    *time.Time                    `db:"read_at"`
		CreatedAt time.Time                     `db:"created_at"`
		UpdatedAt *time.Time                    `db:"updated_at"`
		DeletedAt *time.Time                    `db:"deleted_at"`
	}

	// auxProject is an auxiliary structure used for transporting to/from RDBMS store
	auxProject struct {
		ID                  uint64                           `db:"id"`
		TenantID            uint64                           `db:"tenant_id"`
		Handle              string                           `db:"handle"`
		Status              systemType.ProjectStatus         `db:"status"`
		Config              systemType.ProjectConfig         `db:"config"`
		Meta                systemType.ProjectMeta           `db:"meta"`
		ProjectID           uint64                           `db:"root_project_id"`
		ParentRevisionID    uint64                           `db:"parent_revision_id"`
		Revision            int                              `db:"revision"`
		ArchivedAt          *time.Time                       `db:"archived_at"`
		ApprovalStatus      systemType.ProjectApprovalStatus `db:"approval_status"`
		ApprovalPlan        string                           `db:"approval_plan"`
		ApprovalNote        string                           `db:"approval_note"`
		ApprovalSubmittedBy uint64                           `db:"approval_submitted_by"`
		ApprovalSubmittedAt *time.Time                       `db:"approval_submitted_at"`
		ApprovalDecidedBy   uint64                           `db:"approval_decided_by"`
		ApprovalDecidedAt   *time.Time                       `db:"approval_decided_at"`
		CreatedAt           time.Time                        `db:"created_at"`
		UpdatedAt           *time.Time                       `db:"updated_at"`
		DeletedAt           *time.Time                       `db:"deleted_at"`
		CreatedBy           uint64                           `db:"created_by"`
		UpdatedBy           uint64                           `db:"updated_by"`
		DeletedBy           uint64                           `db:"deleted_by"`
	}

	// auxProjectAiSystem is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectAiSystem struct {
		ID        uint64                         `db:"id"`
		TenantID  uint64                         `db:"tenant_id"`
		ProjectID uint64                         `db:"project_id"`
		Handle    string                         `db:"handle"`
		RiskClass string                         `db:"risk_class"`
		Meta      systemType.ProjectAiSystemMeta `db:"meta"`
		CreatedAt time.Time                      `db:"created_at"`
		UpdatedAt *time.Time                     `db:"updated_at"`
		DeletedAt *time.Time                     `db:"deleted_at"`
	}

	// auxProjectAiSystemEntry is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectAiSystemEntry struct {
		ProjectAiSystemID uint64    `db:"project_ai_system_id"`
		ResourceRef       string    `db:"resource_ref"`
		CreatedAt         time.Time `db:"created_at"`
	}

	// auxProjectBacklogItem is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectBacklogItem struct {
		ID          uint64     `db:"id"`
		TenantID    uint64     `db:"tenant_id"`
		ProjectID   uint64     `db:"project_id"`
		RevisionID  uint64     `db:"revision_id"`
		Title       string     `db:"title"`
		Description string     `db:"description"`
		Category    string     `db:"category"`
		EventID     uint64     `db:"event_id"`
		Assignee    uint64     `db:"assignee"`
		Priority    string     `db:"priority"`
		Status      string     `db:"status"`
		DateDue     string     `db:"date_due"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
		CreatedBy   uint64     `db:"created_by"`
		UpdatedBy   uint64     `db:"updated_by"`
		DeletedBy   uint64     `db:"deleted_by"`
	}

	// auxProjectFeature is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectFeature struct {
		ID               uint64     `db:"id"`
		TenantID         uint64     `db:"tenant_id"`
		ProjectID        uint64     `db:"project_id"`
		RevisionID       uint64     `db:"revision_id"`
		Title            string     `db:"title"`
		Description      string     `db:"description"`
		FeatureType      string     `db:"feature_type"`
		Status           string     `db:"status"`
		Severity         string     `db:"severity"`
		Risk             string     `db:"risk"`
		FeatureOwner     uint64     `db:"feature_owner"`
		ChangeOwner      uint64     `db:"change_owner"`
		ChangeApprovedBy uint64     `db:"change_approved_by"`
		RiskFeature      string     `db:"risk_feature"`
		ChangeRequired   string     `db:"change_required"`
		RiskChange       string     `db:"risk_change"`
		DateDue          string     `db:"date_due"`
		CreatedAt        time.Time  `db:"created_at"`
		UpdatedAt        *time.Time `db:"updated_at"`
		DeletedAt        *time.Time `db:"deleted_at"`
		CreatedBy        uint64     `db:"created_by"`
		UpdatedBy        uint64     `db:"updated_by"`
		DeletedBy        uint64     `db:"deleted_by"`
	}

	// auxProjectFriaScenario is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectFriaScenario struct {
		ID         uint64                             `db:"id"`
		TenantID   uint64                             `db:"tenant_id"`
		ProjectID  uint64                             `db:"project_id"`
		AiSystemID uint64                             `db:"ai_system_id"`
		Title      string                             `db:"title"`
		Severity   string                             `db:"severity"`
		Meta       systemType.ProjectFriaScenarioMeta `db:"meta"`
		CreatedAt  time.Time                          `db:"created_at"`
		UpdatedAt  *time.Time                         `db:"updated_at"`
		DeletedAt  *time.Time                         `db:"deleted_at"`
		CreatedBy  uint64                             `db:"created_by"`
		UpdatedBy  uint64                             `db:"updated_by"`
		DeletedBy  uint64                             `db:"deleted_by"`
	}

	// auxProjectIncident is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectIncident struct {
		ID               uint64     `db:"id"`
		TenantID         uint64     `db:"tenant_id"`
		ProjectID        uint64     `db:"project_id"`
		RevisionID       uint64     `db:"revision_id"`
		Title            string     `db:"title"`
		Description      string     `db:"description"`
		IncidentType     string     `db:"incident_type"`
		GroupSystem      string     `db:"group_system"`
		Status           string     `db:"status"`
		Severity         string     `db:"severity"`
		Risk             string     `db:"risk"`
		IssueOwner       uint64     `db:"issue_owner"`
		ChangeOwner      uint64     `db:"change_owner"`
		ChangeApprovedBy uint64     `db:"change_approved_by"`
		RiskIssue        string     `db:"risk_issue"`
		ChangeRequired   string     `db:"change_required"`
		RiskChange       string     `db:"risk_change"`
		DateDue          string     `db:"date_due"`
		CompletedDate    string     `db:"completed_date"`
		CreatedAt        time.Time  `db:"created_at"`
		UpdatedAt        *time.Time `db:"updated_at"`
		DeletedAt        *time.Time `db:"deleted_at"`
		CreatedBy        uint64     `db:"created_by"`
		UpdatedBy        uint64     `db:"updated_by"`
		DeletedBy        uint64     `db:"deleted_by"`
	}

	// auxProjectMember is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectMember struct {
		ID         uint64                       `db:"id"`
		TenantID   uint64                       `db:"tenant_id"`
		ProjectID  uint64                       `db:"project_id"`
		UserID     uint64                       `db:"user_id"`
		RolePreset systemType.ProjectMemberRole `db:"role_preset"`
		InvitedBy  uint64                       `db:"invited_by"`
		CreatedAt  time.Time                    `db:"created_at"`
		UpdatedAt  *time.Time                   `db:"updated_at"`
		DeletedAt  *time.Time                   `db:"deleted_at"`
	}

	// auxProjectPrivacy is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectPrivacy struct {
		ID               uint64     `db:"id"`
		TenantID         uint64     `db:"tenant_id"`
		ProjectID        uint64     `db:"project_id"`
		RevisionID       uint64     `db:"revision_id"`
		Title            string     `db:"title"`
		Description      string     `db:"description"`
		RequestType      string     `db:"request_type"`
		Status           string     `db:"status"`
		Severity         string     `db:"severity"`
		Risk             string     `db:"risk"`
		RequestOwner     uint64     `db:"request_owner"`
		ChangeOwner      uint64     `db:"change_owner"`
		ChangeApprovedBy uint64     `db:"change_approved_by"`
		RiskAssessment   string     `db:"risk_assessment"`
		ChangeRequired   string     `db:"change_required"`
		RiskChange       string     `db:"risk_change"`
		DateDue          string     `db:"date_due"`
		CreatedAt        time.Time  `db:"created_at"`
		UpdatedAt        *time.Time `db:"updated_at"`
		DeletedAt        *time.Time `db:"deleted_at"`
		CreatedBy        uint64     `db:"created_by"`
		UpdatedBy        uint64     `db:"updated_by"`
		DeletedBy        uint64     `db:"deleted_by"`
	}

	// auxProjectReview is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectReview struct {
		ID              uint64     `db:"id"`
		TenantID        uint64     `db:"tenant_id"`
		ProjectID       uint64     `db:"project_id"`
		RevisionID      uint64     `db:"revision_id"`
		Title           string     `db:"title"`
		Description     string     `db:"description"`
		ReviewType      string     `db:"review_type"`
		ReviewFrequency string     `db:"review_frequency"`
		Scope           string     `db:"scope"`
		Reviewer        uint64     `db:"reviewer"`
		ApprovedBy      uint64     `db:"approved_by"`
		Status          string     `db:"status"`
		DateDue         string     `db:"date_due"`
		CreatedAt       time.Time  `db:"created_at"`
		UpdatedAt       *time.Time `db:"updated_at"`
		DeletedAt       *time.Time `db:"deleted_at"`
		CreatedBy       uint64     `db:"created_by"`
		UpdatedBy       uint64     `db:"updated_by"`
		DeletedBy       uint64     `db:"deleted_by"`
	}

	// auxProjectTask is an auxiliary structure used for transporting to/from RDBMS store
	auxProjectTask struct {
		ID            uint64     `db:"id"`
		TenantID      uint64     `db:"tenant_id"`
		ProjectID     uint64     `db:"project_id"`
		RevisionID    uint64     `db:"revision_id"`
		Title         string     `db:"title"`
		Description   string     `db:"description"`
		TaskName      string     `db:"task_name"`
		TaskType      string     `db:"task_type"`
		Status        string     `db:"status"`
		Severity      string     `db:"severity"`
		Risk          string     `db:"risk"`
		Owner         uint64     `db:"owner"`
		ChangeOwner   uint64     `db:"change_owner"`
		DateDue       string     `db:"date_due"`
		CompletedDate string     `db:"completed_date"`
		CreatedAt     time.Time  `db:"created_at"`
		UpdatedAt     *time.Time `db:"updated_at"`
		DeletedAt     *time.Time `db:"deleted_at"`
		CreatedBy     uint64     `db:"created_by"`
		UpdatedBy     uint64     `db:"updated_by"`
		DeletedBy     uint64     `db:"deleted_by"`
	}

	// auxQueue is an auxiliary structure used for transporting to/from RDBMS store
	auxQueue struct {
		ID        uint64               `db:"id"`
		TenantID  uint64               `db:"tenant_id"`
		ProjectID uint64               `db:"project_id"`
		Consumer  string               `db:"consumer"`
		Queue     string               `db:"queue"`
		Meta      systemType.QueueMeta `db:"meta"`
		CreatedAt time.Time            `db:"created_at"`
		UpdatedAt *time.Time           `db:"updated_at"`
		DeletedAt *time.Time           `db:"deleted_at"`
		CreatedBy uint64               `db:"created_by"`
		UpdatedBy uint64               `db:"updated_by"`
		DeletedBy uint64               `db:"deleted_by"`
	}

	// auxQueueMessage is an auxiliary structure used for transporting to/from RDBMS store
	auxQueueMessage struct {
		ID        uint64     `db:"id"`
		Queue     string     `db:"queue"`
		Payload   []byte     `db:"payload"`
		Created   *time.Time `db:"created"`
		Processed *time.Time `db:"processed"`
	}

	// auxRbacRule is an auxiliary structure used for transporting to/from RDBMS store
	auxRbacRule struct {
		RoleID    uint64          `db:"role_id"`
		Resource  string          `db:"resource"`
		Operation string          `db:"operation"`
		Access    rbacType.Access `db:"access"`
	}

	// auxReminder is an auxiliary structure used for transporting to/from RDBMS store
	auxReminder struct {
		ID          uint64     `db:"id"`
		TenantID    uint64     `db:"tenant_id"`
		ProjectID   uint64     `db:"project_id"`
		Resource    string     `db:"resource"`
		Payload     rawJson    `db:"payload"`
		SnoozeCount uint       `db:"snooze_count"`
		AssignedTo  uint64     `db:"assigned_to"`
		AssignedBy  uint64     `db:"assigned_by"`
		AssignedAt  time.Time  `db:"assigned_at"`
		DismissedBy uint64     `db:"dismissed_by"`
		DismissedAt *time.Time `db:"dismissed_at"`
		RemindAt    *time.Time `db:"remind_at"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
	}

	// auxReport is an auxiliary structure used for transporting to/from RDBMS store
	auxReport struct {
		ID        uint64                         `db:"id"`
		TenantID  uint64                         `db:"tenant_id"`
		ProjectID uint64                         `db:"project_id"`
		Handle    string                         `db:"handle"`
		Meta      *systemType.ReportMeta         `db:"meta"`
		Scenarios systemType.ReportScenarioSet   `db:"scenarios"`
		Sources   systemType.ReportDataSourceSet `db:"sources"`
		Blocks    systemType.ReportBlockSet      `db:"blocks"`
		OwnedBy   uint64                         `db:"owned_by"`
		CreatedAt time.Time                      `db:"created_at"`
		UpdatedAt *time.Time                     `db:"updated_at"`
		DeletedAt *time.Time                     `db:"deleted_at"`
		CreatedBy uint64                         `db:"created_by"`
		UpdatedBy uint64                         `db:"updated_by"`
		DeletedBy uint64                         `db:"deleted_by"`
	}

	// auxResourceActivity is an auxiliary structure used for transporting to/from RDBMS store
	auxResourceActivity struct {
		ID             uint64    `db:"id"`
		Timestamp      time.Time `db:"timestamp"`
		ResourceType   string    `db:"resource_type"`
		ResourceAction string    `db:"resource_action"`
		ResourceID     uint64    `db:"resource_id"`
		Meta           rawJson   `db:"meta"`
	}

	// auxResourceTranslation is an auxiliary structure used for transporting to/from RDBMS store
	auxResourceTranslation struct {
		ID        uint64          `db:"id"`
		TenantID  uint64          `db:"tenant_id"`
		ProjectID uint64          `db:"project_id"`
		Lang      systemType.Lang `db:"lang"`
		Resource  string          `db:"resource"`
		K         string          `db:"k"`
		Message   string          `db:"message"`
		CreatedAt time.Time       `db:"created_at"`
		UpdatedAt *time.Time      `db:"updated_at"`
		DeletedAt *time.Time      `db:"deleted_at"`
		OwnedBy   uint64          `db:"owned_by"`
		CreatedBy uint64          `db:"created_by"`
		UpdatedBy uint64          `db:"updated_by"`
		DeletedBy uint64          `db:"deleted_by"`
	}

	// auxRole is an auxiliary structure used for transporting to/from RDBMS store
	auxRole struct {
		ID         uint64               `db:"id"`
		TenantID   uint64               `db:"tenant_id"`
		ProjectID  uint64               `db:"project_id"`
		Name       string               `db:"name"`
		Handle     string               `db:"handle"`
		Meta       *systemType.RoleMeta `db:"meta"`
		ArchivedAt *time.Time           `db:"archived_at"`
		CreatedAt  time.Time            `db:"created_at"`
		UpdatedAt  *time.Time           `db:"updated_at"`
		DeletedAt  *time.Time           `db:"deleted_at"`
	}

	// auxRoleMember is an auxiliary structure used for transporting to/from RDBMS store
	auxRoleMember struct {
		Resource string `db:"resource"`
		RoleID   uint64 `db:"role_id"`
	}

	// auxSettingValue is an auxiliary structure used for transporting to/from RDBMS store
	auxSettingValue struct {
		OwnedBy   uint64    `db:"owned_by"`
		Name      string    `db:"name"`
		Value     rawJson   `db:"value"`
		UpdatedBy uint64    `db:"updated_by"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	// auxTemplate is an auxiliary structure used for transporting to/from RDBMS store
	auxTemplate struct {
		ID         uint64                  `db:"id"`
		TenantID   uint64                  `db:"tenant_id"`
		ProjectID  uint64                  `db:"project_id"`
		OwnerID    uint64                  `db:"owner_id"`
		Handle     string                  `db:"handle"`
		Language   string                  `db:"language"`
		Type       systemType.DocumentType `db:"type"`
		Partial    bool                    `db:"partial"`
		Meta       systemType.TemplateMeta `db:"meta"`
		Template   string                  `db:"template"`
		CreatedAt  time.Time               `db:"created_at"`
		UpdatedAt  *time.Time              `db:"updated_at"`
		DeletedAt  *time.Time              `db:"deleted_at"`
		LastUsedAt *time.Time              `db:"last_used_at"`
	}

	// auxTenant is an auxiliary structure used for transporting to/from RDBMS store
	auxTenant struct {
		ID          uint64                  `db:"id"`
		Handle      string                  `db:"handle"`
		Status      systemType.TenantStatus `db:"status"`
		Config      systemType.TenantConfig `db:"config"`
		Meta        systemType.TenantMeta   `db:"meta"`
		CreatedAt   time.Time               `db:"created_at"`
		UpdatedAt   *time.Time              `db:"updated_at"`
		SuspendedAt *time.Time              `db:"suspended_at"`
		DeletedAt   *time.Time              `db:"deleted_at"`
		CreatedBy   uint64                  `db:"created_by"`
		UpdatedBy   uint64                  `db:"updated_by"`
		DeletedBy   uint64                  `db:"deleted_by"`
	}

	// auxTenantMembership is an auxiliary structure used for transporting to/from RDBMS store
	auxTenantMembership struct {
		ID        uint64                        `db:"id"`
		TenantID  uint64                        `db:"tenant_id"`
		ProjectID uint64                        `db:"project_id"`
		UserID    uint64                        `db:"user_id"`
		Role      systemType.TenantMemberRole   `db:"role"`
		Status    systemType.TenantMemberStatus `db:"status"`
		InvitedBy uint64                        `db:"invited_by"`
		CreatedAt time.Time                     `db:"created_at"`
		UpdatedAt *time.Time                    `db:"updated_at"`
	}

	// auxUser is an auxiliary structure used for transporting to/from RDBMS store
	auxUser struct {
		ID             uint64               `db:"id"`
		TenantID       uint64               `db:"tenant_id"`
		ProjectID      uint64               `db:"project_id"`
		Email          string               `db:"email"`
		EmailConfirmed bool                 `db:"email_confirmed"`
		UserGroupID    uint64               `db:"user_group_id"`
		Username       string               `db:"username"`
		Name           string               `db:"name"`
		Handle         string               `db:"handle"`
		Kind           systemType.UserKind  `db:"kind"`
		Meta           *systemType.UserMeta `db:"meta"`
		SuspendedAt    *time.Time           `db:"suspended_at"`
		CreatedAt      time.Time            `db:"created_at"`
		UpdatedAt      *time.Time           `db:"updated_at"`
		DeletedAt      *time.Time           `db:"deleted_at"`
	}

	// auxUserGroup is an auxiliary structure used for transporting to/from RDBMS store
	auxUserGroup struct {
		ID         uint64                      `db:"id"`
		TenantID   uint64                      `db:"tenant_id"`
		ProjectID  uint64                      `db:"project_id"`
		Handle     string                      `db:"handle"`
		Meta       *systemType.UserGroupMeta   `db:"meta"`
		Config     *systemType.UserGroupConfig `db:"config"`
		ArchivedAt *time.Time                  `db:"archived_at"`
		CreatedAt  time.Time                   `db:"created_at"`
		UpdatedAt  *time.Time                  `db:"updated_at"`
		DeletedAt  *time.Time                  `db:"deleted_at"`
	}
)

// encodes Actionlog to auxActionlog
//
// This function is auto-generated
func (aux *auxActionlog) encode(res *actionlogType.Action) (_ error) {
	aux.ID = res.ID
	aux.Timestamp = res.Timestamp
	aux.ActorIPAddr = res.ActorIPAddr
	aux.ActorID = res.ActorID
	aux.RequestOrigin = res.RequestOrigin
	aux.RequestID = res.RequestID
	aux.Resource = res.Resource
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RootProjectID = res.RootProjectID
	aux.ResourceProjectID = res.ResourceProjectID
	aux.ResourceRevisionID = res.ResourceRevisionID
	aux.Action = res.Action
	aux.Error = res.Error
	aux.Severity = res.Severity
	aux.Description = res.Description
	aux.Meta = res.Meta
	aux.Delta = res.Delta
	aux.OldState = res.OldState
	return
}

// decodes Actionlog from auxActionlog
//
// This function is auto-generated
func (aux auxActionlog) decode() (res *actionlogType.Action, _ error) {
	res = new(actionlogType.Action)
	res.ID = aux.ID
	res.Timestamp = aux.Timestamp
	res.ActorIPAddr = aux.ActorIPAddr
	res.ActorID = aux.ActorID
	res.RequestOrigin = aux.RequestOrigin
	res.RequestID = aux.RequestID
	res.Resource = aux.Resource
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RootProjectID = aux.RootProjectID
	res.ResourceProjectID = aux.ResourceProjectID
	res.ResourceRevisionID = aux.ResourceRevisionID
	res.Action = aux.Action
	res.Error = aux.Error
	res.Severity = aux.Severity
	res.Description = aux.Description
	res.Meta = aux.Meta
	res.Delta = aux.Delta
	res.OldState = aux.OldState
	return
}

// scans row and fills auxActionlog fields
//
// This function is auto-generated
func (aux *auxActionlog) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Timestamp,
		&aux.ActorIPAddr,
		&aux.ActorID,
		&aux.RequestOrigin,
		&aux.RequestID,
		&aux.Resource,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RootProjectID,
		&aux.ResourceProjectID,
		&aux.ResourceRevisionID,
		&aux.Action,
		&aux.Error,
		&aux.Severity,
		&aux.Description,
		&aux.Meta,
		&aux.Delta,
		&aux.OldState,
	)
}

// encodes Agent to auxAgent
//
// This function is auto-generated
func (aux *auxAgent) encode(res *systemType.Agent) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Status = res.Status
	aux.Revision = res.Revision
	aux.Meta = res.Meta
	aux.Behavior = res.Behavior
	aux.Execution = res.Execution
	aux.Access = res.Access
	aux.Invocation = res.Invocation
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Agent from auxAgent
//
// This function is auto-generated
func (aux auxAgent) decode() (res *systemType.Agent, _ error) {
	res = new(systemType.Agent)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Status = aux.Status
	res.Revision = aux.Revision
	res.Meta = aux.Meta
	res.Behavior = aux.Behavior
	res.Execution = aux.Execution
	res.Access = aux.Access
	res.Invocation = aux.Invocation
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxAgent fields
//
// This function is auto-generated
func (aux *auxAgent) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Status,
		&aux.Revision,
		&aux.Meta,
		&aux.Behavior,
		&aux.Execution,
		&aux.Access,
		&aux.Invocation,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes AiConversation to auxAiConversation
//
// This function is auto-generated
func (aux *auxAiConversation) encode(res *systemType.AiConversation) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.AgentID = res.AgentID
	aux.Messages = res.Messages
	aux.TokenCount = res.TokenCount
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes AiConversation from auxAiConversation
//
// This function is auto-generated
func (aux auxAiConversation) decode() (res *systemType.AiConversation, _ error) {
	res = new(systemType.AiConversation)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.AgentID = aux.AgentID
	res.Messages = aux.Messages
	res.TokenCount = aux.TokenCount
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxAiConversation fields
//
// This function is auto-generated
func (aux *auxAiConversation) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.AgentID,
		&aux.Messages,
		&aux.TokenCount,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ApigwFilter to auxApigwFilter
//
// This function is auto-generated
func (aux *auxApigwFilter) encode(res *systemType.ApigwFilter) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Route = res.Route
	aux.Weight = res.Weight
	aux.Kind = res.Kind
	aux.Ref = res.Ref
	aux.Enabled = res.Enabled
	aux.Params = res.Params
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ApigwFilter from auxApigwFilter
//
// This function is auto-generated
func (aux auxApigwFilter) decode() (res *systemType.ApigwFilter, _ error) {
	res = new(systemType.ApigwFilter)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Route = aux.Route
	res.Weight = aux.Weight
	res.Kind = aux.Kind
	res.Ref = aux.Ref
	res.Enabled = aux.Enabled
	res.Params = aux.Params
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxApigwFilter fields
//
// This function is auto-generated
func (aux *auxApigwFilter) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Route,
		&aux.Weight,
		&aux.Kind,
		&aux.Ref,
		&aux.Enabled,
		&aux.Params,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ApigwRoute to auxApigwRoute
//
// This function is auto-generated
func (aux *auxApigwRoute) encode(res *systemType.ApigwRoute) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Endpoint = res.Endpoint
	aux.Method = res.Method
	aux.Enabled = res.Enabled
	aux.Meta = res.Meta
	aux.Group = res.Group
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ApigwRoute from auxApigwRoute
//
// This function is auto-generated
func (aux auxApigwRoute) decode() (res *systemType.ApigwRoute, _ error) {
	res = new(systemType.ApigwRoute)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Endpoint = aux.Endpoint
	res.Method = aux.Method
	res.Enabled = aux.Enabled
	res.Meta = aux.Meta
	res.Group = aux.Group
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxApigwRoute fields
//
// This function is auto-generated
func (aux *auxApigwRoute) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Endpoint,
		&aux.Method,
		&aux.Enabled,
		&aux.Meta,
		&aux.Group,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Application to auxApplication
//
// This function is auto-generated
func (aux *auxApplication) encode(res *systemType.Application) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Name = res.Name
	aux.Enabled = res.Enabled
	aux.Weight = res.Weight
	aux.Unify = res.Unify
	aux.OwnerID = res.OwnerID
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes Application from auxApplication
//
// This function is auto-generated
func (aux auxApplication) decode() (res *systemType.Application, _ error) {
	res = new(systemType.Application)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Name = aux.Name
	res.Enabled = aux.Enabled
	res.Weight = aux.Weight
	res.Unify = aux.Unify
	res.OwnerID = aux.OwnerID
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxApplication fields
//
// This function is auto-generated
func (aux *auxApplication) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Name,
		&aux.Enabled,
		&aux.Weight,
		&aux.Unify,
		&aux.OwnerID,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes Attachment to auxAttachment
//
// This function is auto-generated
func (aux *auxAttachment) encode(res *systemType.Attachment) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.OwnerID = res.OwnerID
	aux.Kind = res.Kind
	aux.Url = res.Url
	aux.PreviewUrl = res.PreviewUrl
	aux.Name = res.Name
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes Attachment from auxAttachment
//
// This function is auto-generated
func (aux auxAttachment) decode() (res *systemType.Attachment, _ error) {
	res = new(systemType.Attachment)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.OwnerID = aux.OwnerID
	res.Kind = aux.Kind
	res.Url = aux.Url
	res.PreviewUrl = aux.PreviewUrl
	res.Name = aux.Name
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxAttachment fields
//
// This function is auto-generated
func (aux *auxAttachment) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.OwnerID,
		&aux.Kind,
		&aux.Url,
		&aux.PreviewUrl,
		&aux.Name,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes AuthClient to auxAuthClient
//
// This function is auto-generated
func (aux *auxAuthClient) encode(res *systemType.AuthClient) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Meta = res.Meta
	aux.Secret = res.Secret
	aux.Scope = res.Scope
	aux.ValidGrant = res.ValidGrant
	aux.RedirectURI = res.RedirectURI
	aux.Enabled = res.Enabled
	aux.Trusted = res.Trusted
	aux.ValidFrom = res.ValidFrom
	aux.ExpiresAt = res.ExpiresAt
	aux.Security = res.Security
	aux.OwnedBy = res.OwnedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes AuthClient from auxAuthClient
//
// This function is auto-generated
func (aux auxAuthClient) decode() (res *systemType.AuthClient, _ error) {
	res = new(systemType.AuthClient)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Meta = aux.Meta
	res.Secret = aux.Secret
	res.Scope = aux.Scope
	res.ValidGrant = aux.ValidGrant
	res.RedirectURI = aux.RedirectURI
	res.Enabled = aux.Enabled
	res.Trusted = aux.Trusted
	res.ValidFrom = aux.ValidFrom
	res.ExpiresAt = aux.ExpiresAt
	res.Security = aux.Security
	res.OwnedBy = aux.OwnedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxAuthClient fields
//
// This function is auto-generated
func (aux *auxAuthClient) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Meta,
		&aux.Secret,
		&aux.Scope,
		&aux.ValidGrant,
		&aux.RedirectURI,
		&aux.Enabled,
		&aux.Trusted,
		&aux.ValidFrom,
		&aux.ExpiresAt,
		&aux.Security,
		&aux.OwnedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes AuthConfirmedClient to auxAuthConfirmedClient
//
// This function is auto-generated
func (aux *auxAuthConfirmedClient) encode(res *systemType.AuthConfirmedClient) (_ error) {
	aux.UserID = res.UserID
	aux.ClientID = res.ClientID
	aux.ConfirmedAt = res.ConfirmedAt
	return
}

// decodes AuthConfirmedClient from auxAuthConfirmedClient
//
// This function is auto-generated
func (aux auxAuthConfirmedClient) decode() (res *systemType.AuthConfirmedClient, _ error) {
	res = new(systemType.AuthConfirmedClient)
	res.UserID = aux.UserID
	res.ClientID = aux.ClientID
	res.ConfirmedAt = aux.ConfirmedAt
	return
}

// scans row and fills auxAuthConfirmedClient fields
//
// This function is auto-generated
func (aux *auxAuthConfirmedClient) scan(row scanner) error {
	return row.Scan(
		&aux.UserID,
		&aux.ClientID,
		&aux.ConfirmedAt,
	)
}

// encodes AuthOa2token to auxAuthOa2token
//
// This function is auto-generated
func (aux *auxAuthOa2token) encode(res *systemType.AuthOa2token) (_ error) {
	aux.ID = res.ID
	aux.Code = res.Code
	aux.Access = res.Access
	aux.Refresh = res.Refresh
	aux.Data = res.Data
	aux.RemoteAddr = res.RemoteAddr
	aux.UserAgent = res.UserAgent
	aux.ClientID = res.ClientID
	aux.UserID = res.UserID
	aux.CreatedAt = res.CreatedAt
	aux.ExpiresAt = res.ExpiresAt
	return
}

// decodes AuthOa2token from auxAuthOa2token
//
// This function is auto-generated
func (aux auxAuthOa2token) decode() (res *systemType.AuthOa2token, _ error) {
	res = new(systemType.AuthOa2token)
	res.ID = aux.ID
	res.Code = aux.Code
	res.Access = aux.Access
	res.Refresh = aux.Refresh
	res.Data = aux.Data
	res.RemoteAddr = aux.RemoteAddr
	res.UserAgent = aux.UserAgent
	res.ClientID = aux.ClientID
	res.UserID = aux.UserID
	res.CreatedAt = aux.CreatedAt
	res.ExpiresAt = aux.ExpiresAt
	return
}

// scans row and fills auxAuthOa2token fields
//
// This function is auto-generated
func (aux *auxAuthOa2token) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Code,
		&aux.Access,
		&aux.Refresh,
		&aux.Data,
		&aux.RemoteAddr,
		&aux.UserAgent,
		&aux.ClientID,
		&aux.UserID,
		&aux.CreatedAt,
		&aux.ExpiresAt,
	)
}

// encodes AuthSession to auxAuthSession
//
// This function is auto-generated
func (aux *auxAuthSession) encode(res *systemType.AuthSession) (_ error) {
	aux.ID = res.ID
	aux.Data = res.Data
	aux.UserID = res.UserID
	aux.RemoteAddr = res.RemoteAddr
	aux.UserAgent = res.UserAgent
	aux.ExpiresAt = res.ExpiresAt
	aux.CreatedAt = res.CreatedAt
	return
}

// decodes AuthSession from auxAuthSession
//
// This function is auto-generated
func (aux auxAuthSession) decode() (res *systemType.AuthSession, _ error) {
	res = new(systemType.AuthSession)
	res.ID = aux.ID
	res.Data = aux.Data
	res.UserID = aux.UserID
	res.RemoteAddr = aux.RemoteAddr
	res.UserAgent = aux.UserAgent
	res.ExpiresAt = aux.ExpiresAt
	res.CreatedAt = aux.CreatedAt
	return
}

// scans row and fills auxAuthSession fields
//
// This function is auto-generated
func (aux *auxAuthSession) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Data,
		&aux.UserID,
		&aux.RemoteAddr,
		&aux.UserAgent,
		&aux.ExpiresAt,
		&aux.CreatedAt,
	)
}

// encodes AutomationNgAutomation to auxAutomationNgAutomation
//
// This function is auto-generated
func (aux *auxAutomationNgAutomation) encode(res *automationType.NgAutomation) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Meta = res.Meta
	aux.Enabled = res.Enabled
	aux.Scope = res.Scope
	aux.Triggers = res.Triggers
	aux.Steps = res.Steps
	aux.Paths = res.Paths
	aux.Issues = res.Issues
	aux.RunAs = res.RunAs
	aux.OwnedBy = res.OwnedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes AutomationNgAutomation from auxAutomationNgAutomation
//
// This function is auto-generated
func (aux auxAutomationNgAutomation) decode() (res *automationType.NgAutomation, _ error) {
	res = new(automationType.NgAutomation)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Meta = aux.Meta
	res.Enabled = aux.Enabled
	res.Scope = aux.Scope
	res.Triggers = aux.Triggers
	res.Steps = aux.Steps
	res.Paths = aux.Paths
	res.Issues = aux.Issues
	res.RunAs = aux.RunAs
	res.OwnedBy = aux.OwnedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxAutomationNgAutomation fields
//
// This function is auto-generated
func (aux *auxAutomationNgAutomation) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Meta,
		&aux.Enabled,
		&aux.Scope,
		&aux.Triggers,
		&aux.Steps,
		&aux.Paths,
		&aux.Issues,
		&aux.RunAs,
		&aux.OwnedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes AutomationSession to auxAutomationSession
//
// This function is auto-generated
func (aux *auxAutomationSession) encode(res *automationType.Session) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.WorkflowID = res.WorkflowID
	aux.Status = res.Status
	aux.EventType = res.EventType
	aux.ResourceType = res.ResourceType
	aux.Input = res.Input
	aux.Output = res.Output
	aux.Stacktrace = res.Stacktrace
	aux.CreatedBy = res.CreatedBy
	aux.CreatedAt = res.CreatedAt
	aux.PurgeAt = res.PurgeAt
	aux.SuspendedAt = res.SuspendedAt
	aux.CompletedAt = res.CompletedAt
	aux.Error = res.Error
	return
}

// decodes AutomationSession from auxAutomationSession
//
// This function is auto-generated
func (aux auxAutomationSession) decode() (res *automationType.Session, _ error) {
	res = new(automationType.Session)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.WorkflowID = aux.WorkflowID
	res.Status = aux.Status
	res.EventType = aux.EventType
	res.ResourceType = aux.ResourceType
	res.Input = aux.Input
	res.Output = aux.Output
	res.Stacktrace = aux.Stacktrace
	res.CreatedBy = aux.CreatedBy
	res.CreatedAt = aux.CreatedAt
	res.PurgeAt = aux.PurgeAt
	res.SuspendedAt = aux.SuspendedAt
	res.CompletedAt = aux.CompletedAt
	res.Error = aux.Error
	return
}

// scans row and fills auxAutomationSession fields
//
// This function is auto-generated
func (aux *auxAutomationSession) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.WorkflowID,
		&aux.Status,
		&aux.EventType,
		&aux.ResourceType,
		&aux.Input,
		&aux.Output,
		&aux.Stacktrace,
		&aux.CreatedBy,
		&aux.CreatedAt,
		&aux.PurgeAt,
		&aux.SuspendedAt,
		&aux.CompletedAt,
		&aux.Error,
	)
}

// encodes AutomationTrigger to auxAutomationTrigger
//
// This function is auto-generated
func (aux *auxAutomationTrigger) encode(res *automationType.Trigger) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.WorkflowID = res.WorkflowID
	aux.StepID = res.StepID
	aux.Enabled = res.Enabled
	aux.Meta = res.Meta
	aux.ResourceType = res.ResourceType
	aux.EventType = res.EventType
	aux.Constraints = res.Constraints
	aux.Input = res.Input
	aux.OwnedBy = res.OwnedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes AutomationTrigger from auxAutomationTrigger
//
// This function is auto-generated
func (aux auxAutomationTrigger) decode() (res *automationType.Trigger, _ error) {
	res = new(automationType.Trigger)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.WorkflowID = aux.WorkflowID
	res.StepID = aux.StepID
	res.Enabled = aux.Enabled
	res.Meta = aux.Meta
	res.ResourceType = aux.ResourceType
	res.EventType = aux.EventType
	res.Constraints = aux.Constraints
	res.Input = aux.Input
	res.OwnedBy = aux.OwnedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxAutomationTrigger fields
//
// This function is auto-generated
func (aux *auxAutomationTrigger) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.WorkflowID,
		&aux.StepID,
		&aux.Enabled,
		&aux.Meta,
		&aux.ResourceType,
		&aux.EventType,
		&aux.Constraints,
		&aux.Input,
		&aux.OwnedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes AutomationWorkflow to auxAutomationWorkflow
//
// This function is auto-generated
func (aux *auxAutomationWorkflow) encode(res *automationType.Workflow) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Meta = res.Meta
	aux.Enabled = res.Enabled
	aux.Trace = res.Trace
	aux.KeepSessions = res.KeepSessions
	aux.Scope = res.Scope
	aux.Steps = res.Steps
	aux.Paths = res.Paths
	aux.Issues = res.Issues
	aux.RunAs = res.RunAs
	aux.OwnedBy = res.OwnedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes AutomationWorkflow from auxAutomationWorkflow
//
// This function is auto-generated
func (aux auxAutomationWorkflow) decode() (res *automationType.Workflow, _ error) {
	res = new(automationType.Workflow)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Meta = aux.Meta
	res.Enabled = aux.Enabled
	res.Trace = aux.Trace
	res.KeepSessions = aux.KeepSessions
	res.Scope = aux.Scope
	res.Steps = aux.Steps
	res.Paths = aux.Paths
	res.Issues = aux.Issues
	res.RunAs = aux.RunAs
	res.OwnedBy = aux.OwnedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxAutomationWorkflow fields
//
// This function is auto-generated
func (aux *auxAutomationWorkflow) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Meta,
		&aux.Enabled,
		&aux.Trace,
		&aux.KeepSessions,
		&aux.Scope,
		&aux.Steps,
		&aux.Paths,
		&aux.Issues,
		&aux.RunAs,
		&aux.OwnedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Chatbot to auxChatbot
//
// This function is auto-generated
func (aux *auxChatbot) encode(res *systemType.Chatbot) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Name = res.Name
	aux.Enabled = res.Enabled
	aux.WidgetKey = res.WidgetKey
	aux.AllowedOrigins = res.AllowedOrigins
	aux.SessionTTL = res.SessionTTL
	aux.Handoff = res.Handoff
	aux.Styling = res.Styling
	aux.Scenarios = res.Scenarios
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Chatbot from auxChatbot
//
// This function is auto-generated
func (aux auxChatbot) decode() (res *systemType.Chatbot, _ error) {
	res = new(systemType.Chatbot)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Name = aux.Name
	res.Enabled = aux.Enabled
	res.WidgetKey = aux.WidgetKey
	res.AllowedOrigins = aux.AllowedOrigins
	res.SessionTTL = aux.SessionTTL
	res.Handoff = aux.Handoff
	res.Styling = aux.Styling
	res.Scenarios = aux.Scenarios
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxChatbot fields
//
// This function is auto-generated
func (aux *auxChatbot) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Name,
		&aux.Enabled,
		&aux.WidgetKey,
		&aux.AllowedOrigins,
		&aux.SessionTTL,
		&aux.Handoff,
		&aux.Styling,
		&aux.Scenarios,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ChatbotSession to auxChatbotSession
//
// This function is auto-generated
func (aux *auxChatbotSession) encode(res *systemType.ChatbotSession) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.ChatbotID = res.ChatbotID
	aux.Status = res.Status
	aux.CurrentStep = res.CurrentStep
	aux.State = res.State
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ChatbotSession from auxChatbotSession
//
// This function is auto-generated
func (aux auxChatbotSession) decode() (res *systemType.ChatbotSession, _ error) {
	res = new(systemType.ChatbotSession)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.ChatbotID = aux.ChatbotID
	res.Status = aux.Status
	res.CurrentStep = aux.CurrentStep
	res.State = aux.State
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxChatbotSession fields
//
// This function is auto-generated
func (aux *auxChatbotSession) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.ChatbotID,
		&aux.Status,
		&aux.CurrentStep,
		&aux.State,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ChatbotSessionHandoff to auxChatbotSessionHandoff
//
// This function is auto-generated
func (aux *auxChatbotSessionHandoff) encode(res *systemType.ChatbotSessionHandoff) (_ error) {
	aux.ID = res.ID
	aux.SessionID = res.SessionID
	aux.StepID = res.StepID
	aux.Status = res.Status
	aux.InitiatedAt = res.InitiatedAt
	aux.ClosedAt = res.ClosedAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ChatbotSessionHandoff from auxChatbotSessionHandoff
//
// This function is auto-generated
func (aux auxChatbotSessionHandoff) decode() (res *systemType.ChatbotSessionHandoff, _ error) {
	res = new(systemType.ChatbotSessionHandoff)
	res.ID = aux.ID
	res.SessionID = aux.SessionID
	res.StepID = aux.StepID
	res.Status = aux.Status
	res.InitiatedAt = aux.InitiatedAt
	res.ClosedAt = aux.ClosedAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxChatbotSessionHandoff fields
//
// This function is auto-generated
func (aux *auxChatbotSessionHandoff) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.SessionID,
		&aux.StepID,
		&aux.Status,
		&aux.InitiatedAt,
		&aux.ClosedAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ChatbotSessionStep to auxChatbotSessionStep
//
// This function is auto-generated
func (aux *auxChatbotSessionStep) encode(res *systemType.ChatbotSessionStep) (_ error) {
	aux.ID = res.ID
	aux.SessionID = res.SessionID
	aux.ScenarioIndex = res.ScenarioIndex
	aux.ConversationID = res.ConversationID
	aux.Status = res.Status
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ChatbotSessionStep from auxChatbotSessionStep
//
// This function is auto-generated
func (aux auxChatbotSessionStep) decode() (res *systemType.ChatbotSessionStep, _ error) {
	res = new(systemType.ChatbotSessionStep)
	res.ID = aux.ID
	res.SessionID = aux.SessionID
	res.ScenarioIndex = aux.ScenarioIndex
	res.ConversationID = aux.ConversationID
	res.Status = aux.Status
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxChatbotSessionStep fields
//
// This function is auto-generated
func (aux *auxChatbotSessionStep) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.SessionID,
		&aux.ScenarioIndex,
		&aux.ConversationID,
		&aux.Status,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ComposeAttachment to auxComposeAttachment
//
// This function is auto-generated
func (aux *auxComposeAttachment) encode(res *composeType.Attachment) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.NamespaceID = res.NamespaceID
	aux.OwnerID = res.OwnerID
	aux.Kind = res.Kind
	aux.Url = res.Url
	aux.PreviewUrl = res.PreviewUrl
	aux.Name = res.Name
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes ComposeAttachment from auxComposeAttachment
//
// This function is auto-generated
func (aux auxComposeAttachment) decode() (res *composeType.Attachment, _ error) {
	res = new(composeType.Attachment)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.NamespaceID = aux.NamespaceID
	res.OwnerID = aux.OwnerID
	res.Kind = aux.Kind
	res.Url = aux.Url
	res.PreviewUrl = aux.PreviewUrl
	res.Name = aux.Name
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxComposeAttachment fields
//
// This function is auto-generated
func (aux *auxComposeAttachment) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.NamespaceID,
		&aux.OwnerID,
		&aux.Kind,
		&aux.Url,
		&aux.PreviewUrl,
		&aux.Name,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes ComposeChart to auxComposeChart
//
// This function is auto-generated
func (aux *auxComposeChart) encode(res *composeType.Chart) (_ error) {
	aux.ID = res.ID
	aux.Handle = res.Handle
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.NamespaceID = res.NamespaceID
	aux.Name = res.Name
	aux.Config = res.Config
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes ComposeChart from auxComposeChart
//
// This function is auto-generated
func (aux auxComposeChart) decode() (res *composeType.Chart, _ error) {
	res = new(composeType.Chart)
	res.ID = aux.ID
	res.Handle = aux.Handle
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.NamespaceID = aux.NamespaceID
	res.Name = aux.Name
	res.Config = aux.Config
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxComposeChart fields
//
// This function is auto-generated
func (aux *auxComposeChart) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Handle,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.NamespaceID,
		&aux.Name,
		&aux.Config,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes ComposeModule to auxComposeModule
//
// This function is auto-generated
func (aux *auxComposeModule) encode(res *composeType.Module) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.NamespaceID = res.NamespaceID
	aux.Handle = res.Handle
	aux.Name = res.Name
	aux.Meta = res.Meta
	aux.Config = res.Config
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedByAgent = res.CreatedByAgent
	return
}

// decodes ComposeModule from auxComposeModule
//
// This function is auto-generated
func (aux auxComposeModule) decode() (res *composeType.Module, _ error) {
	res = new(composeType.Module)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.NamespaceID = aux.NamespaceID
	res.Handle = aux.Handle
	res.Name = aux.Name
	res.Meta = aux.Meta
	res.Config = aux.Config
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedByAgent = aux.CreatedByAgent
	return
}

// scans row and fills auxComposeModule fields
//
// This function is auto-generated
func (aux *auxComposeModule) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.NamespaceID,
		&aux.Handle,
		&aux.Name,
		&aux.Meta,
		&aux.Config,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedByAgent,
	)
}

// encodes ComposeModuleField to auxComposeModuleField
//
// This function is auto-generated
func (aux *auxComposeModuleField) encode(res *composeType.ModuleField) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.ModuleID = res.ModuleID
	aux.Place = res.Place
	aux.Kind = res.Kind
	aux.Options = res.Options
	aux.Name = res.Name
	aux.Label = res.Label
	aux.Config = res.Config
	aux.Required = res.Required
	aux.Multi = res.Multi
	aux.DefaultValue = res.DefaultValue
	aux.Expressions = res.Expressions
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedByAgent = res.CreatedByAgent
	return
}

// decodes ComposeModuleField from auxComposeModuleField
//
// This function is auto-generated
func (aux auxComposeModuleField) decode() (res *composeType.ModuleField, _ error) {
	res = new(composeType.ModuleField)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.ModuleID = aux.ModuleID
	res.Place = aux.Place
	res.Kind = aux.Kind
	res.Options = aux.Options
	res.Name = aux.Name
	res.Label = aux.Label
	res.Config = aux.Config
	res.Required = aux.Required
	res.Multi = aux.Multi
	res.DefaultValue = aux.DefaultValue
	res.Expressions = aux.Expressions
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedByAgent = aux.CreatedByAgent
	return
}

// scans row and fills auxComposeModuleField fields
//
// This function is auto-generated
func (aux *auxComposeModuleField) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.ModuleID,
		&aux.Place,
		&aux.Kind,
		&aux.Options,
		&aux.Name,
		&aux.Label,
		&aux.Config,
		&aux.Required,
		&aux.Multi,
		&aux.DefaultValue,
		&aux.Expressions,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedByAgent,
	)
}

// encodes ComposeNamespace to auxComposeNamespace
//
// This function is auto-generated
func (aux *auxComposeNamespace) encode(res *composeType.Namespace) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Slug = res.Slug
	aux.Enabled = res.Enabled
	aux.Meta = res.Meta
	aux.Name = res.Name
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedByAgent = res.CreatedByAgent
	return
}

// decodes ComposeNamespace from auxComposeNamespace
//
// This function is auto-generated
func (aux auxComposeNamespace) decode() (res *composeType.Namespace, _ error) {
	res = new(composeType.Namespace)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Slug = aux.Slug
	res.Enabled = aux.Enabled
	res.Meta = aux.Meta
	res.Name = aux.Name
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedByAgent = aux.CreatedByAgent
	return
}

// scans row and fills auxComposeNamespace fields
//
// This function is auto-generated
func (aux *auxComposeNamespace) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Slug,
		&aux.Enabled,
		&aux.Meta,
		&aux.Name,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedByAgent,
	)
}

// encodes ComposePage to auxComposePage
//
// This function is auto-generated
func (aux *auxComposePage) encode(res *composeType.Page) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Title = res.Title
	aux.Handle = res.Handle
	aux.SelfID = res.SelfID
	aux.ModuleID = res.ModuleID
	aux.NamespaceID = res.NamespaceID
	aux.Meta = res.Meta
	aux.Config = res.Config
	aux.Blocks = res.Blocks
	aux.Visible = res.Visible
	aux.Weight = res.Weight
	aux.Description = res.Description
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedByAgent = res.CreatedByAgent
	return
}

// decodes ComposePage from auxComposePage
//
// This function is auto-generated
func (aux auxComposePage) decode() (res *composeType.Page, _ error) {
	res = new(composeType.Page)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Title = aux.Title
	res.Handle = aux.Handle
	res.SelfID = aux.SelfID
	res.ModuleID = aux.ModuleID
	res.NamespaceID = aux.NamespaceID
	res.Meta = aux.Meta
	res.Config = aux.Config
	res.Blocks = aux.Blocks
	res.Visible = aux.Visible
	res.Weight = aux.Weight
	res.Description = aux.Description
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedByAgent = aux.CreatedByAgent
	return
}

// scans row and fills auxComposePage fields
//
// This function is auto-generated
func (aux *auxComposePage) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Title,
		&aux.Handle,
		&aux.SelfID,
		&aux.ModuleID,
		&aux.NamespaceID,
		&aux.Meta,
		&aux.Config,
		&aux.Blocks,
		&aux.Visible,
		&aux.Weight,
		&aux.Description,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedByAgent,
	)
}

// encodes ComposePageLayout to auxComposePageLayout
//
// This function is auto-generated
func (aux *auxComposePageLayout) encode(res *composeType.PageLayout) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.PageID = res.PageID
	aux.ParentID = res.ParentID
	aux.NamespaceID = res.NamespaceID
	aux.Weight = res.Weight
	aux.Meta = res.Meta
	aux.Config = res.Config
	aux.Blocks = res.Blocks
	aux.OwnedBy = res.OwnedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedByAgent = res.CreatedByAgent
	return
}

// decodes ComposePageLayout from auxComposePageLayout
//
// This function is auto-generated
func (aux auxComposePageLayout) decode() (res *composeType.PageLayout, _ error) {
	res = new(composeType.PageLayout)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.PageID = aux.PageID
	res.ParentID = aux.ParentID
	res.NamespaceID = aux.NamespaceID
	res.Weight = aux.Weight
	res.Meta = aux.Meta
	res.Config = aux.Config
	res.Blocks = aux.Blocks
	res.OwnedBy = aux.OwnedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedByAgent = aux.CreatedByAgent
	return
}

// scans row and fills auxComposePageLayout fields
//
// This function is auto-generated
func (aux *auxComposePageLayout) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.PageID,
		&aux.ParentID,
		&aux.NamespaceID,
		&aux.Weight,
		&aux.Meta,
		&aux.Config,
		&aux.Blocks,
		&aux.OwnedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedByAgent,
	)
}

// encodes ConfiguredConnection to auxConfiguredConnection
//
// This function is auto-generated
func (aux *auxConfiguredConnection) encode(res *systemType.ConfiguredConnection) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.ConnectionID = res.ConnectionID
	aux.Name = res.Name
	aux.Status = res.Status
	aux.Connection = res.Connection
	aux.Config = res.Config
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ConfiguredConnection from auxConfiguredConnection
//
// This function is auto-generated
func (aux auxConfiguredConnection) decode() (res *systemType.ConfiguredConnection, _ error) {
	res = new(systemType.ConfiguredConnection)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.ConnectionID = aux.ConnectionID
	res.Name = aux.Name
	res.Status = aux.Status
	res.Connection = aux.Connection
	res.Config = aux.Config
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxConfiguredConnection fields
//
// This function is auto-generated
func (aux *auxConfiguredConnection) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.ConnectionID,
		&aux.Name,
		&aux.Status,
		&aux.Connection,
		&aux.Config,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Connection to auxConnection
//
// This function is auto-generated
func (aux *auxConnection) encode(res *systemType.Connection) (_ error) {
	aux.ID = res.ID
	aux.Handle = res.Handle
	aux.Revision = res.Revision
	aux.Status = res.Status
	aux.Source = res.Source
	aux.Meta = res.Meta
	aux.Service = res.Service
	aux.Resources = res.Resources
	aux.Operations = res.Operations
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Connection from auxConnection
//
// This function is auto-generated
func (aux auxConnection) decode() (res *systemType.Connection, _ error) {
	res = new(systemType.Connection)
	res.ID = aux.ID
	res.Handle = aux.Handle
	res.Revision = aux.Revision
	res.Status = aux.Status
	res.Source = aux.Source
	res.Meta = aux.Meta
	res.Service = aux.Service
	res.Resources = aux.Resources
	res.Operations = aux.Operations
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxConnection fields
//
// This function is auto-generated
func (aux *auxConnection) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Handle,
		&aux.Revision,
		&aux.Status,
		&aux.Source,
		&aux.Meta,
		&aux.Service,
		&aux.Resources,
		&aux.Operations,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Credential to auxCredential
//
// This function is auto-generated
func (aux *auxCredential) encode(res *systemType.Credential) (_ error) {
	aux.ID = res.ID
	aux.OwnerID = res.OwnerID
	aux.Label = res.Label
	aux.Kind = res.Kind
	aux.Credentials = res.Credentials
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.LastUsedAt = res.LastUsedAt
	aux.ExpiresAt = res.ExpiresAt
	return
}

// decodes Credential from auxCredential
//
// This function is auto-generated
func (aux auxCredential) decode() (res *systemType.Credential, _ error) {
	res = new(systemType.Credential)
	res.ID = aux.ID
	res.OwnerID = aux.OwnerID
	res.Label = aux.Label
	res.Kind = aux.Kind
	res.Credentials = aux.Credentials
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.LastUsedAt = aux.LastUsedAt
	res.ExpiresAt = aux.ExpiresAt
	return
}

// scans row and fills auxCredential fields
//
// This function is auto-generated
func (aux *auxCredential) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.OwnerID,
		&aux.Label,
		&aux.Kind,
		&aux.Credentials,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.LastUsedAt,
		&aux.ExpiresAt,
	)
}

// encodes DalConnection to auxDalConnection
//
// This function is auto-generated
func (aux *auxDalConnection) encode(res *systemType.DalConnection) (_ error) {
	aux.ID = res.ID
	aux.Handle = res.Handle
	aux.Type = res.Type
	aux.Config = res.Config
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes DalConnection from auxDalConnection
//
// This function is auto-generated
func (aux auxDalConnection) decode() (res *systemType.DalConnection, _ error) {
	res = new(systemType.DalConnection)
	res.ID = aux.ID
	res.Handle = aux.Handle
	res.Type = aux.Type
	res.Config = aux.Config
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxDalConnection fields
//
// This function is auto-generated
func (aux *auxDalConnection) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Handle,
		&aux.Type,
		&aux.Config,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes DalSchemaAlteration to auxDalSchemaAlteration
//
// This function is auto-generated
func (aux *auxDalSchemaAlteration) encode(res *systemType.DalSchemaAlteration) (_ error) {
	aux.ID = res.ID
	aux.BatchID = res.BatchID
	aux.DependsOn = res.DependsOn
	aux.Resource = res.Resource
	aux.ResourceType = res.ResourceType
	aux.ConnectionID = res.ConnectionID
	aux.Kind = res.Kind
	aux.Params = res.Params
	aux.Error = res.Error
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CompletedAt = res.CompletedAt
	aux.DismissedAt = res.DismissedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	aux.CompletedBy = res.CompletedBy
	aux.DismissedBy = res.DismissedBy
	return
}

// decodes DalSchemaAlteration from auxDalSchemaAlteration
//
// This function is auto-generated
func (aux auxDalSchemaAlteration) decode() (res *systemType.DalSchemaAlteration, _ error) {
	res = new(systemType.DalSchemaAlteration)
	res.ID = aux.ID
	res.BatchID = aux.BatchID
	res.DependsOn = aux.DependsOn
	res.Resource = aux.Resource
	res.ResourceType = aux.ResourceType
	res.ConnectionID = aux.ConnectionID
	res.Kind = aux.Kind
	res.Params = aux.Params
	res.Error = aux.Error
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CompletedAt = aux.CompletedAt
	res.DismissedAt = aux.DismissedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	res.CompletedBy = aux.CompletedBy
	res.DismissedBy = aux.DismissedBy
	return
}

// scans row and fills auxDalSchemaAlteration fields
//
// This function is auto-generated
func (aux *auxDalSchemaAlteration) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.BatchID,
		&aux.DependsOn,
		&aux.Resource,
		&aux.ResourceType,
		&aux.ConnectionID,
		&aux.Kind,
		&aux.Params,
		&aux.Error,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CompletedAt,
		&aux.DismissedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
		&aux.CompletedBy,
		&aux.DismissedBy,
	)
}

// encodes DalSensitivityLevel to auxDalSensitivityLevel
//
// This function is auto-generated
func (aux *auxDalSensitivityLevel) encode(res *systemType.DalSensitivityLevel) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Level = res.Level
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes DalSensitivityLevel from auxDalSensitivityLevel
//
// This function is auto-generated
func (aux auxDalSensitivityLevel) decode() (res *systemType.DalSensitivityLevel, _ error) {
	res = new(systemType.DalSensitivityLevel)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Level = aux.Level
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxDalSensitivityLevel fields
//
// This function is auto-generated
func (aux *auxDalSensitivityLevel) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Level,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes DataPrivacyRequest to auxDataPrivacyRequest
//
// This function is auto-generated
func (aux *auxDataPrivacyRequest) encode(res *systemType.DataPrivacyRequest) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Kind = res.Kind
	aux.Status = res.Status
	aux.Payload = res.Payload
	aux.RequestedAt = res.RequestedAt
	aux.RequestedBy = res.RequestedBy
	aux.CompletedAt = res.CompletedAt
	aux.CompletedBy = res.CompletedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes DataPrivacyRequest from auxDataPrivacyRequest
//
// This function is auto-generated
func (aux auxDataPrivacyRequest) decode() (res *systemType.DataPrivacyRequest, _ error) {
	res = new(systemType.DataPrivacyRequest)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Kind = aux.Kind
	res.Status = aux.Status
	res.Payload = aux.Payload
	res.RequestedAt = aux.RequestedAt
	res.RequestedBy = aux.RequestedBy
	res.CompletedAt = aux.CompletedAt
	res.CompletedBy = aux.CompletedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxDataPrivacyRequest fields
//
// This function is auto-generated
func (aux *auxDataPrivacyRequest) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Kind,
		&aux.Status,
		&aux.Payload,
		&aux.RequestedAt,
		&aux.RequestedBy,
		&aux.CompletedAt,
		&aux.CompletedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes DataPrivacyRequestComment to auxDataPrivacyRequestComment
//
// This function is auto-generated
func (aux *auxDataPrivacyRequestComment) encode(res *systemType.DataPrivacyRequestComment) (_ error) {
	aux.ID = res.ID
	aux.RequestID = res.RequestID
	aux.Comment = res.Comment
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes DataPrivacyRequestComment from auxDataPrivacyRequestComment
//
// This function is auto-generated
func (aux auxDataPrivacyRequestComment) decode() (res *systemType.DataPrivacyRequestComment, _ error) {
	res = new(systemType.DataPrivacyRequestComment)
	res.ID = aux.ID
	res.RequestID = aux.RequestID
	res.Comment = aux.Comment
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxDataPrivacyRequestComment fields
//
// This function is auto-generated
func (aux *auxDataPrivacyRequestComment) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.RequestID,
		&aux.Comment,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes DmlConnection to auxDmlConnection
//
// This function is auto-generated
func (aux *auxDmlConnection) encode(res *systemType.DmlConnection) (_ error) {
	aux.ID = res.ID
	aux.Handle = res.Handle
	aux.Label = res.Label
	aux.Params = res.Params
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes DmlConnection from auxDmlConnection
//
// This function is auto-generated
func (aux auxDmlConnection) decode() (res *systemType.DmlConnection, _ error) {
	res = new(systemType.DmlConnection)
	res.ID = aux.ID
	res.Handle = aux.Handle
	res.Label = aux.Label
	res.Params = aux.Params
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxDmlConnection fields
//
// This function is auto-generated
func (aux *auxDmlConnection) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Handle,
		&aux.Label,
		&aux.Params,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes DmlImportRun to auxDmlImportRun
//
// This function is auto-generated
func (aux *auxDmlImportRun) encode(res *systemType.DmlImportRun) (_ error) {
	aux.ID = res.ID
	aux.ConnectionID = res.ConnectionID
	aux.MappingID = res.MappingID
	aux.Method = res.Method
	aux.Status = res.Status
	aux.Processed = res.Processed
	aux.Failed = res.Failed
	aux.Error = res.Error
	return
}

// decodes DmlImportRun from auxDmlImportRun
//
// This function is auto-generated
func (aux auxDmlImportRun) decode() (res *systemType.DmlImportRun, _ error) {
	res = new(systemType.DmlImportRun)
	res.ID = aux.ID
	res.ConnectionID = aux.ConnectionID
	res.MappingID = aux.MappingID
	res.Method = aux.Method
	res.Status = aux.Status
	res.Processed = aux.Processed
	res.Failed = aux.Failed
	res.Error = aux.Error
	return
}

// scans row and fills auxDmlImportRun fields
//
// This function is auto-generated
func (aux *auxDmlImportRun) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.ConnectionID,
		&aux.MappingID,
		&aux.Method,
		&aux.Status,
		&aux.Processed,
		&aux.Failed,
		&aux.Error,
	)
}

// encodes DmlMapping to auxDmlMapping
//
// This function is auto-generated
func (aux *auxDmlMapping) encode(res *systemType.DmlMapping) (_ error) {
	aux.ID = res.ID
	aux.ConnectionID = res.ConnectionID
	aux.NamespaceHandle = res.NamespaceHandle
	aux.SourceIdent = res.SourceIdent
	aux.ModuleHandle = res.ModuleHandle
	aux.ModuleName = res.ModuleName
	aux.Skip = res.Skip
	aux.Identifier = res.Identifier
	aux.Columns = res.Columns
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes DmlMapping from auxDmlMapping
//
// This function is auto-generated
func (aux auxDmlMapping) decode() (res *systemType.DmlMapping, _ error) {
	res = new(systemType.DmlMapping)
	res.ID = aux.ID
	res.ConnectionID = aux.ConnectionID
	res.NamespaceHandle = aux.NamespaceHandle
	res.SourceIdent = aux.SourceIdent
	res.ModuleHandle = aux.ModuleHandle
	res.ModuleName = aux.ModuleName
	res.Skip = aux.Skip
	res.Identifier = aux.Identifier
	res.Columns = aux.Columns
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxDmlMapping fields
//
// This function is auto-generated
func (aux *auxDmlMapping) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.ConnectionID,
		&aux.NamespaceHandle,
		&aux.SourceIdent,
		&aux.ModuleHandle,
		&aux.ModuleName,
		&aux.Skip,
		&aux.Identifier,
		&aux.Columns,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes FederationExposedModule to auxFederationExposedModule
//
// This function is auto-generated
func (aux *auxFederationExposedModule) encode(res *federationType.ExposedModule) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Name = res.Name
	aux.NodeID = res.NodeID
	aux.ComposeModuleID = res.ComposeModuleID
	aux.ComposeNamespaceID = res.ComposeNamespaceID
	aux.Fields = res.Fields
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes FederationExposedModule from auxFederationExposedModule
//
// This function is auto-generated
func (aux auxFederationExposedModule) decode() (res *federationType.ExposedModule, _ error) {
	res = new(federationType.ExposedModule)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Name = aux.Name
	res.NodeID = aux.NodeID
	res.ComposeModuleID = aux.ComposeModuleID
	res.ComposeNamespaceID = aux.ComposeNamespaceID
	res.Fields = aux.Fields
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxFederationExposedModule fields
//
// This function is auto-generated
func (aux *auxFederationExposedModule) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Name,
		&aux.NodeID,
		&aux.ComposeModuleID,
		&aux.ComposeNamespaceID,
		&aux.Fields,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes FederationModuleMapping to auxFederationModuleMapping
//
// This function is auto-generated
func (aux *auxFederationModuleMapping) encode(res *federationType.ModuleMapping) (_ error) {
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.NodeID = res.NodeID
	aux.FederationModuleID = res.FederationModuleID
	aux.ComposeModuleID = res.ComposeModuleID
	aux.ComposeNamespaceID = res.ComposeNamespaceID
	aux.FieldMapping = res.FieldMapping
	return
}

// decodes FederationModuleMapping from auxFederationModuleMapping
//
// This function is auto-generated
func (aux auxFederationModuleMapping) decode() (res *federationType.ModuleMapping, _ error) {
	res = new(federationType.ModuleMapping)
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.NodeID = aux.NodeID
	res.FederationModuleID = aux.FederationModuleID
	res.ComposeModuleID = aux.ComposeModuleID
	res.ComposeNamespaceID = aux.ComposeNamespaceID
	res.FieldMapping = aux.FieldMapping
	return
}

// scans row and fills auxFederationModuleMapping fields
//
// This function is auto-generated
func (aux *auxFederationModuleMapping) scan(row scanner) error {
	return row.Scan(
		&aux.TenantID,
		&aux.ProjectID,
		&aux.NodeID,
		&aux.FederationModuleID,
		&aux.ComposeModuleID,
		&aux.ComposeNamespaceID,
		&aux.FieldMapping,
	)
}

// encodes FederationNode to auxFederationNode
//
// This function is auto-generated
func (aux *auxFederationNode) encode(res *federationType.Node) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.SharedNodeID = res.SharedNodeID
	aux.Name = res.Name
	aux.BaseURL = res.BaseURL
	aux.Status = res.Status
	aux.Contact = res.Contact
	aux.PairToken = res.PairToken
	aux.AuthToken = res.AuthToken
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes FederationNode from auxFederationNode
//
// This function is auto-generated
func (aux auxFederationNode) decode() (res *federationType.Node, _ error) {
	res = new(federationType.Node)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.SharedNodeID = aux.SharedNodeID
	res.Name = aux.Name
	res.BaseURL = aux.BaseURL
	res.Status = aux.Status
	res.Contact = aux.Contact
	res.PairToken = aux.PairToken
	res.AuthToken = aux.AuthToken
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxFederationNode fields
//
// This function is auto-generated
func (aux *auxFederationNode) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.SharedNodeID,
		&aux.Name,
		&aux.BaseURL,
		&aux.Status,
		&aux.Contact,
		&aux.PairToken,
		&aux.AuthToken,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes FederationNodeSync to auxFederationNodeSync
//
// This function is auto-generated
func (aux *auxFederationNodeSync) encode(res *federationType.NodeSync) (_ error) {
	aux.NodeID = res.NodeID
	aux.ModuleID = res.ModuleID
	aux.SyncType = res.SyncType
	aux.SyncStatus = res.SyncStatus
	aux.TimeOfAction = res.TimeOfAction
	return
}

// decodes FederationNodeSync from auxFederationNodeSync
//
// This function is auto-generated
func (aux auxFederationNodeSync) decode() (res *federationType.NodeSync, _ error) {
	res = new(federationType.NodeSync)
	res.NodeID = aux.NodeID
	res.ModuleID = aux.ModuleID
	res.SyncType = aux.SyncType
	res.SyncStatus = aux.SyncStatus
	res.TimeOfAction = aux.TimeOfAction
	return
}

// scans row and fills auxFederationNodeSync fields
//
// This function is auto-generated
func (aux *auxFederationNodeSync) scan(row scanner) error {
	return row.Scan(
		&aux.NodeID,
		&aux.ModuleID,
		&aux.SyncType,
		&aux.SyncStatus,
		&aux.TimeOfAction,
	)
}

// encodes FederationSharedModule to auxFederationSharedModule
//
// This function is auto-generated
func (aux *auxFederationSharedModule) encode(res *federationType.SharedModule) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.NodeID = res.NodeID
	aux.Name = res.Name
	aux.ExternalFederationModuleID = res.ExternalFederationModuleID
	aux.Fields = res.Fields
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes FederationSharedModule from auxFederationSharedModule
//
// This function is auto-generated
func (aux auxFederationSharedModule) decode() (res *federationType.SharedModule, _ error) {
	res = new(federationType.SharedModule)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.NodeID = aux.NodeID
	res.Name = aux.Name
	res.ExternalFederationModuleID = aux.ExternalFederationModuleID
	res.Fields = aux.Fields
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxFederationSharedModule fields
//
// This function is auto-generated
func (aux *auxFederationSharedModule) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.NodeID,
		&aux.Name,
		&aux.ExternalFederationModuleID,
		&aux.Fields,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Flag to auxFlag
//
// This function is auto-generated
func (aux *auxFlag) encode(res *flagType.Flag) (_ error) {
	aux.Kind = res.Kind
	aux.ResourceID = res.ResourceID
	aux.OwnedBy = res.OwnedBy
	aux.Name = res.Name
	aux.Active = res.Active
	return
}

// decodes Flag from auxFlag
//
// This function is auto-generated
func (aux auxFlag) decode() (res *flagType.Flag, _ error) {
	res = new(flagType.Flag)
	res.Kind = aux.Kind
	res.ResourceID = aux.ResourceID
	res.OwnedBy = aux.OwnedBy
	res.Name = aux.Name
	res.Active = aux.Active
	return
}

// scans row and fills auxFlag fields
//
// This function is auto-generated
func (aux *auxFlag) scan(row scanner) error {
	return row.Scan(
		&aux.Kind,
		&aux.ResourceID,
		&aux.OwnedBy,
		&aux.Name,
		&aux.Active,
	)
}

// encodes KnowledgeBase to auxKnowledgeBase
//
// This function is auto-generated
func (aux *auxKnowledgeBase) encode(res *systemType.KnowledgeBase) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Title = res.Title
	aux.Description = res.Description
	aux.Context = res.Context
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes KnowledgeBase from auxKnowledgeBase
//
// This function is auto-generated
func (aux auxKnowledgeBase) decode() (res *systemType.KnowledgeBase, _ error) {
	res = new(systemType.KnowledgeBase)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Title = aux.Title
	res.Description = aux.Description
	res.Context = aux.Context
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxKnowledgeBase fields
//
// This function is auto-generated
func (aux *auxKnowledgeBase) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Title,
		&aux.Description,
		&aux.Context,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Label to auxLabel
//
// This function is auto-generated
func (aux *auxLabel) encode(res *labelsType.Label) (_ error) {
	aux.Kind = res.Kind
	aux.ResourceID = res.ResourceID
	aux.Name = res.Name
	aux.Value = res.Value
	return
}

// decodes Label from auxLabel
//
// This function is auto-generated
func (aux auxLabel) decode() (res *labelsType.Label, _ error) {
	res = new(labelsType.Label)
	res.Kind = aux.Kind
	res.ResourceID = aux.ResourceID
	res.Name = aux.Name
	res.Value = aux.Value
	return
}

// scans row and fills auxLabel fields
//
// This function is auto-generated
func (aux *auxLabel) scan(row scanner) error {
	return row.Scan(
		&aux.Kind,
		&aux.ResourceID,
		&aux.Name,
		&aux.Value,
	)
}

// encodes LlmProvider to auxLlmProvider
//
// This function is auto-generated
func (aux *auxLlmProvider) encode(res *systemType.LlmProvider) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Status = res.Status
	aux.Provider = res.Provider
	aux.CredentialID = res.CredentialID
	aux.Meta = res.Meta
	aux.Config = res.Config
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes LlmProvider from auxLlmProvider
//
// This function is auto-generated
func (aux auxLlmProvider) decode() (res *systemType.LlmProvider, _ error) {
	res = new(systemType.LlmProvider)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Status = aux.Status
	res.Provider = aux.Provider
	res.CredentialID = aux.CredentialID
	res.Meta = aux.Meta
	res.Config = aux.Config
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxLlmProvider fields
//
// This function is auto-generated
func (aux *auxLlmProvider) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Status,
		&aux.Provider,
		&aux.CredentialID,
		&aux.Meta,
		&aux.Config,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Notification to auxNotification
//
// This function is auto-generated
func (aux *auxNotification) encode(res *systemType.Notification) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Kind = res.Kind
	aux.Config = res.Config
	aux.Recipient = res.Recipient
	aux.CreatedBy = res.CreatedBy
	aux.ReadAt = res.ReadAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes Notification from auxNotification
//
// This function is auto-generated
func (aux auxNotification) decode() (res *systemType.Notification, _ error) {
	res = new(systemType.Notification)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Kind = aux.Kind
	res.Config = aux.Config
	res.Recipient = aux.Recipient
	res.CreatedBy = aux.CreatedBy
	res.ReadAt = aux.ReadAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxNotification fields
//
// This function is auto-generated
func (aux *auxNotification) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Kind,
		&aux.Config,
		&aux.Recipient,
		&aux.CreatedBy,
		&aux.ReadAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes Project to auxProject
//
// This function is auto-generated
func (aux *auxProject) encode(res *systemType.Project) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.Handle = res.Handle
	aux.Status = res.Status
	aux.Config = res.Config
	aux.Meta = res.Meta
	aux.ProjectID = res.ProjectID
	aux.ParentRevisionID = res.ParentRevisionID
	aux.Revision = res.Revision
	aux.ArchivedAt = res.ArchivedAt
	aux.ApprovalStatus = res.ApprovalStatus
	aux.ApprovalPlan = res.ApprovalPlan
	aux.ApprovalNote = res.ApprovalNote
	aux.ApprovalSubmittedBy = res.ApprovalSubmittedBy
	aux.ApprovalSubmittedAt = res.ApprovalSubmittedAt
	aux.ApprovalDecidedBy = res.ApprovalDecidedBy
	aux.ApprovalDecidedAt = res.ApprovalDecidedAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Project from auxProject
//
// This function is auto-generated
func (aux auxProject) decode() (res *systemType.Project, _ error) {
	res = new(systemType.Project)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.Handle = aux.Handle
	res.Status = aux.Status
	res.Config = aux.Config
	res.Meta = aux.Meta
	res.ProjectID = aux.ProjectID
	res.ParentRevisionID = aux.ParentRevisionID
	res.Revision = aux.Revision
	res.ArchivedAt = aux.ArchivedAt
	res.ApprovalStatus = aux.ApprovalStatus
	res.ApprovalPlan = aux.ApprovalPlan
	res.ApprovalNote = aux.ApprovalNote
	res.ApprovalSubmittedBy = aux.ApprovalSubmittedBy
	res.ApprovalSubmittedAt = aux.ApprovalSubmittedAt
	res.ApprovalDecidedBy = aux.ApprovalDecidedBy
	res.ApprovalDecidedAt = aux.ApprovalDecidedAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProject fields
//
// This function is auto-generated
func (aux *auxProject) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.Handle,
		&aux.Status,
		&aux.Config,
		&aux.Meta,
		&aux.ProjectID,
		&aux.ParentRevisionID,
		&aux.Revision,
		&aux.ArchivedAt,
		&aux.ApprovalStatus,
		&aux.ApprovalPlan,
		&aux.ApprovalNote,
		&aux.ApprovalSubmittedBy,
		&aux.ApprovalSubmittedAt,
		&aux.ApprovalDecidedBy,
		&aux.ApprovalDecidedAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectAiSystem to auxProjectAiSystem
//
// This function is auto-generated
func (aux *auxProjectAiSystem) encode(res *systemType.ProjectAiSystem) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.RiskClass = res.RiskClass
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes ProjectAiSystem from auxProjectAiSystem
//
// This function is auto-generated
func (aux auxProjectAiSystem) decode() (res *systemType.ProjectAiSystem, _ error) {
	res = new(systemType.ProjectAiSystem)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.RiskClass = aux.RiskClass
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxProjectAiSystem fields
//
// This function is auto-generated
func (aux *auxProjectAiSystem) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.RiskClass,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes ProjectAiSystemEntry to auxProjectAiSystemEntry
//
// This function is auto-generated
func (aux *auxProjectAiSystemEntry) encode(res *systemType.ProjectAiSystemEntry) (_ error) {
	aux.ProjectAiSystemID = res.ProjectAiSystemID
	aux.ResourceRef = res.ResourceRef
	aux.CreatedAt = res.CreatedAt
	return
}

// decodes ProjectAiSystemEntry from auxProjectAiSystemEntry
//
// This function is auto-generated
func (aux auxProjectAiSystemEntry) decode() (res *systemType.ProjectAiSystemEntry, _ error) {
	res = new(systemType.ProjectAiSystemEntry)
	res.ProjectAiSystemID = aux.ProjectAiSystemID
	res.ResourceRef = aux.ResourceRef
	res.CreatedAt = aux.CreatedAt
	return
}

// scans row and fills auxProjectAiSystemEntry fields
//
// This function is auto-generated
func (aux *auxProjectAiSystemEntry) scan(row scanner) error {
	return row.Scan(
		&aux.ProjectAiSystemID,
		&aux.ResourceRef,
		&aux.CreatedAt,
	)
}

// encodes ProjectBacklogItem to auxProjectBacklogItem
//
// This function is auto-generated
func (aux *auxProjectBacklogItem) encode(res *systemType.ProjectBacklogItem) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RevisionID = res.RevisionID
	aux.Title = res.Title
	aux.Description = res.Description
	aux.Category = res.Category
	aux.EventID = res.EventID
	aux.Assignee = res.Assignee
	aux.Priority = res.Priority
	aux.Status = res.Status
	aux.DateDue = res.DateDue
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectBacklogItem from auxProjectBacklogItem
//
// This function is auto-generated
func (aux auxProjectBacklogItem) decode() (res *systemType.ProjectBacklogItem, _ error) {
	res = new(systemType.ProjectBacklogItem)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RevisionID = aux.RevisionID
	res.Title = aux.Title
	res.Description = aux.Description
	res.Category = aux.Category
	res.EventID = aux.EventID
	res.Assignee = aux.Assignee
	res.Priority = aux.Priority
	res.Status = aux.Status
	res.DateDue = aux.DateDue
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectBacklogItem fields
//
// This function is auto-generated
func (aux *auxProjectBacklogItem) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RevisionID,
		&aux.Title,
		&aux.Description,
		&aux.Category,
		&aux.EventID,
		&aux.Assignee,
		&aux.Priority,
		&aux.Status,
		&aux.DateDue,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectFeature to auxProjectFeature
//
// This function is auto-generated
func (aux *auxProjectFeature) encode(res *systemType.ProjectFeature) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RevisionID = res.RevisionID
	aux.Title = res.Title
	aux.Description = res.Description
	aux.FeatureType = res.FeatureType
	aux.Status = res.Status
	aux.Severity = res.Severity
	aux.Risk = res.Risk
	aux.FeatureOwner = res.FeatureOwner
	aux.ChangeOwner = res.ChangeOwner
	aux.ChangeApprovedBy = res.ChangeApprovedBy
	aux.RiskFeature = res.RiskFeature
	aux.ChangeRequired = res.ChangeRequired
	aux.RiskChange = res.RiskChange
	aux.DateDue = res.DateDue
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectFeature from auxProjectFeature
//
// This function is auto-generated
func (aux auxProjectFeature) decode() (res *systemType.ProjectFeature, _ error) {
	res = new(systemType.ProjectFeature)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RevisionID = aux.RevisionID
	res.Title = aux.Title
	res.Description = aux.Description
	res.FeatureType = aux.FeatureType
	res.Status = aux.Status
	res.Severity = aux.Severity
	res.Risk = aux.Risk
	res.FeatureOwner = aux.FeatureOwner
	res.ChangeOwner = aux.ChangeOwner
	res.ChangeApprovedBy = aux.ChangeApprovedBy
	res.RiskFeature = aux.RiskFeature
	res.ChangeRequired = aux.ChangeRequired
	res.RiskChange = aux.RiskChange
	res.DateDue = aux.DateDue
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectFeature fields
//
// This function is auto-generated
func (aux *auxProjectFeature) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RevisionID,
		&aux.Title,
		&aux.Description,
		&aux.FeatureType,
		&aux.Status,
		&aux.Severity,
		&aux.Risk,
		&aux.FeatureOwner,
		&aux.ChangeOwner,
		&aux.ChangeApprovedBy,
		&aux.RiskFeature,
		&aux.ChangeRequired,
		&aux.RiskChange,
		&aux.DateDue,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectFriaScenario to auxProjectFriaScenario
//
// This function is auto-generated
func (aux *auxProjectFriaScenario) encode(res *systemType.ProjectFriaScenario) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.AiSystemID = res.AiSystemID
	aux.Title = res.Title
	aux.Severity = res.Severity
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectFriaScenario from auxProjectFriaScenario
//
// This function is auto-generated
func (aux auxProjectFriaScenario) decode() (res *systemType.ProjectFriaScenario, _ error) {
	res = new(systemType.ProjectFriaScenario)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.AiSystemID = aux.AiSystemID
	res.Title = aux.Title
	res.Severity = aux.Severity
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectFriaScenario fields
//
// This function is auto-generated
func (aux *auxProjectFriaScenario) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.AiSystemID,
		&aux.Title,
		&aux.Severity,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectIncident to auxProjectIncident
//
// This function is auto-generated
func (aux *auxProjectIncident) encode(res *systemType.ProjectIncident) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RevisionID = res.RevisionID
	aux.Title = res.Title
	aux.Description = res.Description
	aux.IncidentType = res.IncidentType
	aux.GroupSystem = res.GroupSystem
	aux.Status = res.Status
	aux.Severity = res.Severity
	aux.Risk = res.Risk
	aux.IssueOwner = res.IssueOwner
	aux.ChangeOwner = res.ChangeOwner
	aux.ChangeApprovedBy = res.ChangeApprovedBy
	aux.RiskIssue = res.RiskIssue
	aux.ChangeRequired = res.ChangeRequired
	aux.RiskChange = res.RiskChange
	aux.DateDue = res.DateDue
	aux.CompletedDate = res.CompletedDate
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectIncident from auxProjectIncident
//
// This function is auto-generated
func (aux auxProjectIncident) decode() (res *systemType.ProjectIncident, _ error) {
	res = new(systemType.ProjectIncident)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RevisionID = aux.RevisionID
	res.Title = aux.Title
	res.Description = aux.Description
	res.IncidentType = aux.IncidentType
	res.GroupSystem = aux.GroupSystem
	res.Status = aux.Status
	res.Severity = aux.Severity
	res.Risk = aux.Risk
	res.IssueOwner = aux.IssueOwner
	res.ChangeOwner = aux.ChangeOwner
	res.ChangeApprovedBy = aux.ChangeApprovedBy
	res.RiskIssue = aux.RiskIssue
	res.ChangeRequired = aux.ChangeRequired
	res.RiskChange = aux.RiskChange
	res.DateDue = aux.DateDue
	res.CompletedDate = aux.CompletedDate
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectIncident fields
//
// This function is auto-generated
func (aux *auxProjectIncident) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RevisionID,
		&aux.Title,
		&aux.Description,
		&aux.IncidentType,
		&aux.GroupSystem,
		&aux.Status,
		&aux.Severity,
		&aux.Risk,
		&aux.IssueOwner,
		&aux.ChangeOwner,
		&aux.ChangeApprovedBy,
		&aux.RiskIssue,
		&aux.ChangeRequired,
		&aux.RiskChange,
		&aux.DateDue,
		&aux.CompletedDate,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectMember to auxProjectMember
//
// This function is auto-generated
func (aux *auxProjectMember) encode(res *systemType.ProjectMember) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.UserID = res.UserID
	aux.RolePreset = res.RolePreset
	aux.InvitedBy = res.InvitedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes ProjectMember from auxProjectMember
//
// This function is auto-generated
func (aux auxProjectMember) decode() (res *systemType.ProjectMember, _ error) {
	res = new(systemType.ProjectMember)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.UserID = aux.UserID
	res.RolePreset = aux.RolePreset
	res.InvitedBy = aux.InvitedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxProjectMember fields
//
// This function is auto-generated
func (aux *auxProjectMember) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.UserID,
		&aux.RolePreset,
		&aux.InvitedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes ProjectPrivacy to auxProjectPrivacy
//
// This function is auto-generated
func (aux *auxProjectPrivacy) encode(res *systemType.ProjectPrivacy) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RevisionID = res.RevisionID
	aux.Title = res.Title
	aux.Description = res.Description
	aux.RequestType = res.RequestType
	aux.Status = res.Status
	aux.Severity = res.Severity
	aux.Risk = res.Risk
	aux.RequestOwner = res.RequestOwner
	aux.ChangeOwner = res.ChangeOwner
	aux.ChangeApprovedBy = res.ChangeApprovedBy
	aux.RiskAssessment = res.RiskAssessment
	aux.ChangeRequired = res.ChangeRequired
	aux.RiskChange = res.RiskChange
	aux.DateDue = res.DateDue
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectPrivacy from auxProjectPrivacy
//
// This function is auto-generated
func (aux auxProjectPrivacy) decode() (res *systemType.ProjectPrivacy, _ error) {
	res = new(systemType.ProjectPrivacy)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RevisionID = aux.RevisionID
	res.Title = aux.Title
	res.Description = aux.Description
	res.RequestType = aux.RequestType
	res.Status = aux.Status
	res.Severity = aux.Severity
	res.Risk = aux.Risk
	res.RequestOwner = aux.RequestOwner
	res.ChangeOwner = aux.ChangeOwner
	res.ChangeApprovedBy = aux.ChangeApprovedBy
	res.RiskAssessment = aux.RiskAssessment
	res.ChangeRequired = aux.ChangeRequired
	res.RiskChange = aux.RiskChange
	res.DateDue = aux.DateDue
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectPrivacy fields
//
// This function is auto-generated
func (aux *auxProjectPrivacy) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RevisionID,
		&aux.Title,
		&aux.Description,
		&aux.RequestType,
		&aux.Status,
		&aux.Severity,
		&aux.Risk,
		&aux.RequestOwner,
		&aux.ChangeOwner,
		&aux.ChangeApprovedBy,
		&aux.RiskAssessment,
		&aux.ChangeRequired,
		&aux.RiskChange,
		&aux.DateDue,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectReview to auxProjectReview
//
// This function is auto-generated
func (aux *auxProjectReview) encode(res *systemType.ProjectReview) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RevisionID = res.RevisionID
	aux.Title = res.Title
	aux.Description = res.Description
	aux.ReviewType = res.ReviewType
	aux.ReviewFrequency = res.ReviewFrequency
	aux.Scope = res.Scope
	aux.Reviewer = res.Reviewer
	aux.ApprovedBy = res.ApprovedBy
	aux.Status = res.Status
	aux.DateDue = res.DateDue
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectReview from auxProjectReview
//
// This function is auto-generated
func (aux auxProjectReview) decode() (res *systemType.ProjectReview, _ error) {
	res = new(systemType.ProjectReview)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RevisionID = aux.RevisionID
	res.Title = aux.Title
	res.Description = aux.Description
	res.ReviewType = aux.ReviewType
	res.ReviewFrequency = aux.ReviewFrequency
	res.Scope = aux.Scope
	res.Reviewer = aux.Reviewer
	res.ApprovedBy = aux.ApprovedBy
	res.Status = aux.Status
	res.DateDue = aux.DateDue
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectReview fields
//
// This function is auto-generated
func (aux *auxProjectReview) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RevisionID,
		&aux.Title,
		&aux.Description,
		&aux.ReviewType,
		&aux.ReviewFrequency,
		&aux.Scope,
		&aux.Reviewer,
		&aux.ApprovedBy,
		&aux.Status,
		&aux.DateDue,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ProjectTask to auxProjectTask
//
// This function is auto-generated
func (aux *auxProjectTask) encode(res *systemType.ProjectTask) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.RevisionID = res.RevisionID
	aux.Title = res.Title
	aux.Description = res.Description
	aux.TaskName = res.TaskName
	aux.TaskType = res.TaskType
	aux.Status = res.Status
	aux.Severity = res.Severity
	aux.Risk = res.Risk
	aux.Owner = res.Owner
	aux.ChangeOwner = res.ChangeOwner
	aux.DateDue = res.DateDue
	aux.CompletedDate = res.CompletedDate
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ProjectTask from auxProjectTask
//
// This function is auto-generated
func (aux auxProjectTask) decode() (res *systemType.ProjectTask, _ error) {
	res = new(systemType.ProjectTask)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.RevisionID = aux.RevisionID
	res.Title = aux.Title
	res.Description = aux.Description
	res.TaskName = aux.TaskName
	res.TaskType = aux.TaskType
	res.Status = aux.Status
	res.Severity = aux.Severity
	res.Risk = aux.Risk
	res.Owner = aux.Owner
	res.ChangeOwner = aux.ChangeOwner
	res.DateDue = aux.DateDue
	res.CompletedDate = aux.CompletedDate
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxProjectTask fields
//
// This function is auto-generated
func (aux *auxProjectTask) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.RevisionID,
		&aux.Title,
		&aux.Description,
		&aux.TaskName,
		&aux.TaskType,
		&aux.Status,
		&aux.Severity,
		&aux.Risk,
		&aux.Owner,
		&aux.ChangeOwner,
		&aux.DateDue,
		&aux.CompletedDate,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Queue to auxQueue
//
// This function is auto-generated
func (aux *auxQueue) encode(res *systemType.Queue) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Consumer = res.Consumer
	aux.Queue = res.Queue
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Queue from auxQueue
//
// This function is auto-generated
func (aux auxQueue) decode() (res *systemType.Queue, _ error) {
	res = new(systemType.Queue)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Consumer = aux.Consumer
	res.Queue = aux.Queue
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxQueue fields
//
// This function is auto-generated
func (aux *auxQueue) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Consumer,
		&aux.Queue,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes QueueMessage to auxQueueMessage
//
// This function is auto-generated
func (aux *auxQueueMessage) encode(res *systemType.QueueMessage) (_ error) {
	aux.ID = res.ID
	aux.Queue = res.Queue
	aux.Payload = res.Payload
	aux.Created = res.Created
	aux.Processed = res.Processed
	return
}

// decodes QueueMessage from auxQueueMessage
//
// This function is auto-generated
func (aux auxQueueMessage) decode() (res *systemType.QueueMessage, _ error) {
	res = new(systemType.QueueMessage)
	res.ID = aux.ID
	res.Queue = aux.Queue
	res.Payload = aux.Payload
	res.Created = aux.Created
	res.Processed = aux.Processed
	return
}

// scans row and fills auxQueueMessage fields
//
// This function is auto-generated
func (aux *auxQueueMessage) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Queue,
		&aux.Payload,
		&aux.Created,
		&aux.Processed,
	)
}

// encodes RbacRule to auxRbacRule
//
// This function is auto-generated
func (aux *auxRbacRule) encode(res *rbacType.Rule) (_ error) {
	aux.RoleID = res.RoleID
	aux.Resource = res.Resource
	aux.Operation = res.Operation
	aux.Access = res.Access
	return
}

// decodes RbacRule from auxRbacRule
//
// This function is auto-generated
func (aux auxRbacRule) decode() (res *rbacType.Rule, _ error) {
	res = new(rbacType.Rule)
	res.RoleID = aux.RoleID
	res.Resource = aux.Resource
	res.Operation = aux.Operation
	res.Access = aux.Access
	return
}

// scans row and fills auxRbacRule fields
//
// This function is auto-generated
func (aux *auxRbacRule) scan(row scanner) error {
	return row.Scan(
		&aux.RoleID,
		&aux.Resource,
		&aux.Operation,
		&aux.Access,
	)
}

// encodes Reminder to auxReminder
//
// This function is auto-generated
func (aux *auxReminder) encode(res *systemType.Reminder) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Resource = res.Resource
	aux.Payload = res.Payload
	aux.SnoozeCount = res.SnoozeCount
	aux.AssignedTo = res.AssignedTo
	aux.AssignedBy = res.AssignedBy
	aux.AssignedAt = res.AssignedAt
	aux.DismissedBy = res.DismissedBy
	aux.DismissedAt = res.DismissedAt
	aux.RemindAt = res.RemindAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes Reminder from auxReminder
//
// This function is auto-generated
func (aux auxReminder) decode() (res *systemType.Reminder, _ error) {
	res = new(systemType.Reminder)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Resource = aux.Resource
	res.Payload = aux.Payload
	res.SnoozeCount = aux.SnoozeCount
	res.AssignedTo = aux.AssignedTo
	res.AssignedBy = aux.AssignedBy
	res.AssignedAt = aux.AssignedAt
	res.DismissedBy = aux.DismissedBy
	res.DismissedAt = aux.DismissedAt
	res.RemindAt = aux.RemindAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxReminder fields
//
// This function is auto-generated
func (aux *auxReminder) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Resource,
		&aux.Payload,
		&aux.SnoozeCount,
		&aux.AssignedTo,
		&aux.AssignedBy,
		&aux.AssignedAt,
		&aux.DismissedBy,
		&aux.DismissedAt,
		&aux.RemindAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes Report to auxReport
//
// This function is auto-generated
func (aux *auxReport) encode(res *systemType.Report) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Meta = res.Meta
	aux.Scenarios = res.Scenarios
	aux.Sources = res.Sources
	aux.Blocks = res.Blocks
	aux.OwnedBy = res.OwnedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Report from auxReport
//
// This function is auto-generated
func (aux auxReport) decode() (res *systemType.Report, _ error) {
	res = new(systemType.Report)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Meta = aux.Meta
	res.Scenarios = aux.Scenarios
	res.Sources = aux.Sources
	res.Blocks = aux.Blocks
	res.OwnedBy = aux.OwnedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxReport fields
//
// This function is auto-generated
func (aux *auxReport) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Meta,
		&aux.Scenarios,
		&aux.Sources,
		&aux.Blocks,
		&aux.OwnedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes ResourceActivity to auxResourceActivity
//
// This function is auto-generated
func (aux *auxResourceActivity) encode(res *discoveryType.ResourceActivity) (_ error) {
	aux.ID = res.ID
	aux.Timestamp = res.Timestamp
	aux.ResourceType = res.ResourceType
	aux.ResourceAction = res.ResourceAction
	aux.ResourceID = res.ResourceID
	aux.Meta = res.Meta
	return
}

// decodes ResourceActivity from auxResourceActivity
//
// This function is auto-generated
func (aux auxResourceActivity) decode() (res *discoveryType.ResourceActivity, _ error) {
	res = new(discoveryType.ResourceActivity)
	res.ID = aux.ID
	res.Timestamp = aux.Timestamp
	res.ResourceType = aux.ResourceType
	res.ResourceAction = aux.ResourceAction
	res.ResourceID = aux.ResourceID
	res.Meta = aux.Meta
	return
}

// scans row and fills auxResourceActivity fields
//
// This function is auto-generated
func (aux *auxResourceActivity) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Timestamp,
		&aux.ResourceType,
		&aux.ResourceAction,
		&aux.ResourceID,
		&aux.Meta,
	)
}

// encodes ResourceTranslation to auxResourceTranslation
//
// This function is auto-generated
func (aux *auxResourceTranslation) encode(res *systemType.ResourceTranslation) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Lang = res.Lang
	aux.Resource = res.Resource
	aux.K = res.K
	aux.Message = res.Message
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.OwnedBy = res.OwnedBy
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes ResourceTranslation from auxResourceTranslation
//
// This function is auto-generated
func (aux auxResourceTranslation) decode() (res *systemType.ResourceTranslation, _ error) {
	res = new(systemType.ResourceTranslation)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Lang = aux.Lang
	res.Resource = aux.Resource
	res.K = aux.K
	res.Message = aux.Message
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.OwnedBy = aux.OwnedBy
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxResourceTranslation fields
//
// This function is auto-generated
func (aux *auxResourceTranslation) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Lang,
		&aux.Resource,
		&aux.K,
		&aux.Message,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.OwnedBy,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes Role to auxRole
//
// This function is auto-generated
func (aux *auxRole) encode(res *systemType.Role) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Name = res.Name
	aux.Handle = res.Handle
	aux.Meta = res.Meta
	aux.ArchivedAt = res.ArchivedAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes Role from auxRole
//
// This function is auto-generated
func (aux auxRole) decode() (res *systemType.Role, _ error) {
	res = new(systemType.Role)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Name = aux.Name
	res.Handle = aux.Handle
	res.Meta = aux.Meta
	res.ArchivedAt = aux.ArchivedAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxRole fields
//
// This function is auto-generated
func (aux *auxRole) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Name,
		&aux.Handle,
		&aux.Meta,
		&aux.ArchivedAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes RoleMember to auxRoleMember
//
// This function is auto-generated
func (aux *auxRoleMember) encode(res *systemType.RoleMember) (_ error) {
	aux.Resource = res.Resource
	aux.RoleID = res.RoleID
	return
}

// decodes RoleMember from auxRoleMember
//
// This function is auto-generated
func (aux auxRoleMember) decode() (res *systemType.RoleMember, _ error) {
	res = new(systemType.RoleMember)
	res.Resource = aux.Resource
	res.RoleID = aux.RoleID
	return
}

// scans row and fills auxRoleMember fields
//
// This function is auto-generated
func (aux *auxRoleMember) scan(row scanner) error {
	return row.Scan(
		&aux.Resource,
		&aux.RoleID,
	)
}

// encodes SettingValue to auxSettingValue
//
// This function is auto-generated
func (aux *auxSettingValue) encode(res *systemType.SettingValue) (_ error) {
	aux.OwnedBy = res.OwnedBy
	aux.Name = res.Name
	aux.Value = res.Value
	aux.UpdatedBy = res.UpdatedBy
	aux.UpdatedAt = res.UpdatedAt
	return
}

// decodes SettingValue from auxSettingValue
//
// This function is auto-generated
func (aux auxSettingValue) decode() (res *systemType.SettingValue, _ error) {
	res = new(systemType.SettingValue)
	res.OwnedBy = aux.OwnedBy
	res.Name = aux.Name
	res.Value = aux.Value
	res.UpdatedBy = aux.UpdatedBy
	res.UpdatedAt = aux.UpdatedAt
	return
}

// scans row and fills auxSettingValue fields
//
// This function is auto-generated
func (aux *auxSettingValue) scan(row scanner) error {
	return row.Scan(
		&aux.OwnedBy,
		&aux.Name,
		&aux.Value,
		&aux.UpdatedBy,
		&aux.UpdatedAt,
	)
}

// encodes Template to auxTemplate
//
// This function is auto-generated
func (aux *auxTemplate) encode(res *systemType.Template) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.OwnerID = res.OwnerID
	aux.Handle = res.Handle
	aux.Language = res.Language
	aux.Type = res.Type
	aux.Partial = res.Partial
	aux.Meta = res.Meta
	aux.Template = res.Template
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.LastUsedAt = res.LastUsedAt
	return
}

// decodes Template from auxTemplate
//
// This function is auto-generated
func (aux auxTemplate) decode() (res *systemType.Template, _ error) {
	res = new(systemType.Template)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.OwnerID = aux.OwnerID
	res.Handle = aux.Handle
	res.Language = aux.Language
	res.Type = aux.Type
	res.Partial = aux.Partial
	res.Meta = aux.Meta
	res.Template = aux.Template
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.LastUsedAt = aux.LastUsedAt
	return
}

// scans row and fills auxTemplate fields
//
// This function is auto-generated
func (aux *auxTemplate) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.OwnerID,
		&aux.Handle,
		&aux.Language,
		&aux.Type,
		&aux.Partial,
		&aux.Meta,
		&aux.Template,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.LastUsedAt,
	)
}

// encodes Tenant to auxTenant
//
// This function is auto-generated
func (aux *auxTenant) encode(res *systemType.Tenant) (_ error) {
	aux.ID = res.ID
	aux.Handle = res.Handle
	aux.Status = res.Status
	aux.Config = res.Config
	aux.Meta = res.Meta
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.SuspendedAt = res.SuspendedAt
	aux.DeletedAt = res.DeletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	return
}

// decodes Tenant from auxTenant
//
// This function is auto-generated
func (aux auxTenant) decode() (res *systemType.Tenant, _ error) {
	res = new(systemType.Tenant)
	res.ID = aux.ID
	res.Handle = aux.Handle
	res.Status = aux.Status
	res.Config = aux.Config
	res.Meta = aux.Meta
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.SuspendedAt = aux.SuspendedAt
	res.DeletedAt = aux.DeletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	return
}

// scans row and fills auxTenant fields
//
// This function is auto-generated
func (aux *auxTenant) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Handle,
		&aux.Status,
		&aux.Config,
		&aux.Meta,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.SuspendedAt,
		&aux.DeletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
	)
}

// encodes TenantMembership to auxTenantMembership
//
// This function is auto-generated
func (aux *auxTenantMembership) encode(res *systemType.TenantMembership) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.UserID = res.UserID
	aux.Role = res.Role
	aux.Status = res.Status
	aux.InvitedBy = res.InvitedBy
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	return
}

// decodes TenantMembership from auxTenantMembership
//
// This function is auto-generated
func (aux auxTenantMembership) decode() (res *systemType.TenantMembership, _ error) {
	res = new(systemType.TenantMembership)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.UserID = aux.UserID
	res.Role = aux.Role
	res.Status = aux.Status
	res.InvitedBy = aux.InvitedBy
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	return
}

// scans row and fills auxTenantMembership fields
//
// This function is auto-generated
func (aux *auxTenantMembership) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.UserID,
		&aux.Role,
		&aux.Status,
		&aux.InvitedBy,
		&aux.CreatedAt,
		&aux.UpdatedAt,
	)
}

// encodes User to auxUser
//
// This function is auto-generated
func (aux *auxUser) encode(res *systemType.User) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Email = res.Email
	aux.EmailConfirmed = res.EmailConfirmed
	aux.UserGroupID = res.UserGroupID
	aux.Username = res.Username
	aux.Name = res.Name
	aux.Handle = res.Handle
	aux.Kind = res.Kind
	aux.Meta = res.Meta
	aux.SuspendedAt = res.SuspendedAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes User from auxUser
//
// This function is auto-generated
func (aux auxUser) decode() (res *systemType.User, _ error) {
	res = new(systemType.User)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Email = aux.Email
	res.EmailConfirmed = aux.EmailConfirmed
	res.UserGroupID = aux.UserGroupID
	res.Username = aux.Username
	res.Name = aux.Name
	res.Handle = aux.Handle
	res.Kind = aux.Kind
	res.Meta = aux.Meta
	res.SuspendedAt = aux.SuspendedAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxUser fields
//
// This function is auto-generated
func (aux *auxUser) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Email,
		&aux.EmailConfirmed,
		&aux.UserGroupID,
		&aux.Username,
		&aux.Name,
		&aux.Handle,
		&aux.Kind,
		&aux.Meta,
		&aux.SuspendedAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes UserGroup to auxUserGroup
//
// This function is auto-generated
func (aux *auxUserGroup) encode(res *systemType.UserGroup) (_ error) {
	aux.ID = res.ID
	aux.TenantID = res.TenantID
	aux.ProjectID = res.ProjectID
	aux.Handle = res.Handle
	aux.Meta = res.Meta
	aux.Config = res.Config
	aux.ArchivedAt = res.ArchivedAt
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes UserGroup from auxUserGroup
//
// This function is auto-generated
func (aux auxUserGroup) decode() (res *systemType.UserGroup, _ error) {
	res = new(systemType.UserGroup)
	res.ID = aux.ID
	res.TenantID = aux.TenantID
	res.ProjectID = aux.ProjectID
	res.Handle = aux.Handle
	res.Meta = aux.Meta
	res.Config = aux.Config
	res.ArchivedAt = aux.ArchivedAt
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxUserGroup fields
//
// This function is auto-generated
func (aux *auxUserGroup) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.TenantID,
		&aux.ProjectID,
		&aux.Handle,
		&aux.Meta,
		&aux.Config,
		&aux.ArchivedAt,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}
