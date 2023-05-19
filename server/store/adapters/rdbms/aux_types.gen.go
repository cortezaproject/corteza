package rdbms

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	automationType "github.com/cortezaproject/corteza/server/automation/types"
	composeType "github.com/cortezaproject/corteza/server/compose/types"
	discoveryType "github.com/cortezaproject/corteza/server/discovery/types"
	federationType "github.com/cortezaproject/corteza/server/federation/types"
	actionlogType "github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	flagType "github.com/cortezaproject/corteza/server/pkg/flag/types"
	labelsType "github.com/cortezaproject/corteza/server/pkg/label/types"
	rbacType "github.com/cortezaproject/corteza/server/pkg/rbac"
	systemType "github.com/cortezaproject/corteza/server/system/types"
	"time"
)

type (

	// auxActionlog is an auxiliary structure used for transporting to/from RDBMS store
	auxActionlog struct {
		ID            uint64                 `db:"id"`
		Timestamp     time.Time              `db:"timestamp"`
		ActorIPAddr   string                 `db:"actor_ip_addr"`
		ActorID       uint64                 `db:"actor_id"`
		RequestOrigin string                 `db:"request_origin"`
		RequestID     string                 `db:"request_id"`
		Resource      string                 `db:"resource"`
		Action        string                 `db:"action"`
		Error         string                 `db:"error"`
		Severity      actionlogType.Severity `db:"severity"`
		Description   string                 `db:"description"`
		Meta          actionlogType.Meta     `db:"meta"`
	}

	// auxApigwFilter is an auxiliary structure used for transporting to/from RDBMS store
	auxApigwFilter struct {
		ID        uint64                       `db:"id"`
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

	// auxAutomationSession is an auxiliary structure used for transporting to/from RDBMS store
	auxAutomationSession struct {
		ID           uint64                       `db:"id"`
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

	// auxComposeAttachment is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeAttachment struct {
		ID          uint64                     `db:"id"`
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
		NamespaceID uint64                  `db:"namespace_id"`
		Name        string                  `db:"name"`
		Config      composeType.ChartConfig `db:"config"`
		CreatedAt   time.Time               `db:"created_at"`
		UpdatedAt   *time.Time              `db:"updated_at"`
		DeletedAt   *time.Time              `db:"deleted_at"`
	}

	// auxComposeModule is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeModule struct {
		ID          uint64                   `db:"id"`
		NamespaceID uint64                   `db:"namespace_id"`
		Handle      string                   `db:"handle"`
		Name        string                   `db:"name"`
		Meta        rawJson                  `db:"meta"`
		Config      composeType.ModuleConfig `db:"config"`
		CreatedAt   time.Time                `db:"created_at"`
		UpdatedAt   *time.Time               `db:"updated_at"`
		DeletedAt   *time.Time               `db:"deleted_at"`
	}

	// auxComposeModuleField is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeModuleField struct {
		ID           uint64                         `db:"id"`
		ModuleID     uint64                         `db:"module_id"`
		Place        int                            `db:"place"`
		Kind         string                         `db:"kind"`
		Options      composeType.ModuleFieldOptions `db:"options"`
		Name         string                         `db:"name"`
		Label        string                         `db:"label"`
		Config       composeType.ModuleFieldConfig  `db:"config"`
		Required     bool                           `db:"required"`
		Multi        bool                           `db:"multi"`
		DefaultValue composeType.RecordValueSet     `db:"default_value"`
		Expressions  composeType.ModuleFieldExpr    `db:"expressions"`
		CreatedAt    time.Time                      `db:"created_at"`
		UpdatedAt    *time.Time                     `db:"updated_at"`
		DeletedAt    *time.Time                     `db:"deleted_at"`
	}

	// auxComposeNamespace is an auxiliary structure used for transporting to/from RDBMS store
	auxComposeNamespace struct {
		ID        uint64                    `db:"id"`
		Slug      string                    `db:"slug"`
		Enabled   bool                      `db:"enabled"`
		Meta      composeType.NamespaceMeta `db:"meta"`
		Name      string                    `db:"name"`
		CreatedAt time.Time                 `db:"created_at"`
		UpdatedAt *time.Time                `db:"updated_at"`
		DeletedAt *time.Time                `db:"deleted_at"`
	}

	// auxComposePage is an auxiliary structure used for transporting to/from RDBMS store
	auxComposePage struct {
		ID          uint64                 `db:"id"`
		Title       string                 `db:"title"`
		Handle      string                 `db:"handle"`
		SelfID      uint64                 `db:"self_id"`
		ModuleID    uint64                 `db:"module_id"`
		NamespaceID uint64                 `db:"namespace_id"`
		Meta        composeType.PageMeta   `db:"meta"`
		Config      composeType.PageConfig `db:"config"`
		Blocks      composeType.PageBlocks `db:"blocks"`
		Visible     bool                   `db:"visible"`
		Weight      int                    `db:"weight"`
		Description string                 `db:"description"`
		CreatedAt   time.Time              `db:"created_at"`
		UpdatedAt   *time.Time             `db:"updated_at"`
		DeletedAt   *time.Time             `db:"deleted_at"`
	}

	// auxComposePageLayout is an auxiliary structure used for transporting to/from RDBMS store
	auxComposePageLayout struct {
		ID          uint64                       `db:"id"`
		Handle      string                       `db:"handle"`
		PageID      uint64                       `db:"page_id"`
		ParentID    uint64                       `db:"parent_id"`
		NamespaceID uint64                       `db:"namespace_id"`
		Weight      int                          `db:"weight"`
		Meta        composeType.PageLayoutMeta   `db:"meta"`
		Config      composeType.PageLayoutConfig `db:"config"`
		Blocks      composeType.PageLayoutBlocks `db:"blocks"`
		OwnedBy     uint64                       `db:"owned_by"`
		CreatedAt   time.Time                    `db:"created_at"`
		UpdatedAt   *time.Time                   `db:"updated_at"`
		DeletedAt   *time.Time                   `db:"deleted_at"`
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
		ID        uint64                      `db:"id"`
		Handle    string                      `db:"handle"`
		Type      string                      `db:"type"`
		Config    systemType.ConnectionConfig `db:"config"`
		Meta      systemType.ConnectionMeta   `db:"meta"`
		CreatedAt time.Time                   `db:"created_at"`
		UpdatedAt *time.Time                  `db:"updated_at"`
		DeletedAt *time.Time                  `db:"deleted_at"`
		CreatedBy uint64                      `db:"created_by"`
		UpdatedBy uint64                      `db:"updated_by"`
		DeletedBy uint64                      `db:"deleted_by"`
	}

	// auxDalSchemaAlteration is an auxiliary structure used for transporting to/from RDBMS store
	auxDalSchemaAlteration struct {
		ID          uint64                                `db:"id"`
		BatchID     uint64                                `db:"batchID"`
		DependsOn   uint64                                `db:"dependsOn"`
		Kind        string                                `db:"kind"`
		Params      *systemType.DalSchemaAlterationParams `db:"params"`
		CreatedAt   time.Time                             `db:"created_at"`
		UpdatedAt   *time.Time                            `db:"updated_at"`
		DeletedAt   *time.Time                            `db:"deleted_at"`
		CompletedAt *time.Time                            `db:"completed_at"`
		CreatedBy   uint64                                `db:"created_by"`
		UpdatedBy   uint64                                `db:"updated_by"`
		DeletedBy   uint64                                `db:"deleted_by"`
		CompletedBy uint64                                `db:"completed_by"`
	}

	// auxDalSensitivityLevel is an auxiliary structure used for transporting to/from RDBMS store
	auxDalSensitivityLevel struct {
		ID        uint64                             `db:"id"`
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

	// auxFederationExposedModule is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationExposedModule struct {
		ID                 uint64                        `db:"id"`
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
		NodeID             uint64                               `db:"node_id"`
		FederationModuleID uint64                               `db:"federation_module_id"`
		ComposeModuleID    uint64                               `db:"compose_module_id"`
		ComposeNamespaceID uint64                               `db:"compose_namespace_id"`
		FieldMapping       federationType.ModuleFieldMappingSet `db:"field_mapping"`
	}

	// auxFederationNode is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationNode struct {
		ID           uint64     `db:"id"`
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
		NodeID       uint64    `db:"node_id"`
		ModuleID     uint64    `db:"module_id"`
		SyncType     string    `db:"sync_type"`
		SyncStatus   string    `db:"sync_status"`
		TimeOfAction time.Time `db:"time_of_action"`
	}

	// auxFederationSharedModule is an auxiliary structure used for transporting to/from RDBMS store
	auxFederationSharedModule struct {
		ID                         uint64                        `db:"id"`
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

	// auxLabel is an auxiliary structure used for transporting to/from RDBMS store
	auxLabel struct {
		Kind       string `db:"kind"`
		ResourceID uint64 `db:"resource_id"`
		Name       string `db:"name"`
		Value      string `db:"value"`
	}

	// auxQueue is an auxiliary structure used for transporting to/from RDBMS store
	auxQueue struct {
		ID        uint64               `db:"id"`
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
		UserID uint64 `db:"user_id"`
		RoleID uint64 `db:"role_id"`
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

	// auxUser is an auxiliary structure used for transporting to/from RDBMS store
	auxUser struct {
		ID             uint64               `db:"id"`
		Email          string               `db:"email"`
		EmailConfirmed bool                 `db:"email_confirmed"`
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
	aux.Action = res.Action
	aux.Error = res.Error
	aux.Severity = res.Severity
	aux.Description = res.Description
	aux.Meta = res.Meta
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
	res.Action = aux.Action
	res.Error = aux.Error
	res.Severity = aux.Severity
	res.Description = aux.Description
	res.Meta = aux.Meta
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
		&aux.Action,
		&aux.Error,
		&aux.Severity,
		&aux.Description,
		&aux.Meta,
	)
}

// encodes ApigwFilter to auxApigwFilter
//
// This function is auto-generated
func (aux *auxApigwFilter) encode(res *systemType.ApigwFilter) (_ error) {
	aux.ID = res.ID
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

// encodes AutomationSession to auxAutomationSession
//
// This function is auto-generated
func (aux *auxAutomationSession) encode(res *automationType.Session) (_ error) {
	aux.ID = res.ID
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

// encodes ComposeAttachment to auxComposeAttachment
//
// This function is auto-generated
func (aux *auxComposeAttachment) encode(res *composeType.Attachment) (_ error) {
	aux.ID = res.ID
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
	aux.NamespaceID = res.NamespaceID
	aux.Handle = res.Handle
	aux.Name = res.Name
	aux.Meta = res.Meta
	aux.Config = res.Config
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes ComposeModule from auxComposeModule
//
// This function is auto-generated
func (aux auxComposeModule) decode() (res *composeType.Module, _ error) {
	res = new(composeType.Module)
	res.ID = aux.ID
	res.NamespaceID = aux.NamespaceID
	res.Handle = aux.Handle
	res.Name = aux.Name
	res.Meta = aux.Meta
	res.Config = aux.Config
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxComposeModule fields
//
// This function is auto-generated
func (aux *auxComposeModule) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.NamespaceID,
		&aux.Handle,
		&aux.Name,
		&aux.Meta,
		&aux.Config,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes ComposeModuleField to auxComposeModuleField
//
// This function is auto-generated
func (aux *auxComposeModuleField) encode(res *composeType.ModuleField) (_ error) {
	aux.ID = res.ID
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
	return
}

// decodes ComposeModuleField from auxComposeModuleField
//
// This function is auto-generated
func (aux auxComposeModuleField) decode() (res *composeType.ModuleField, _ error) {
	res = new(composeType.ModuleField)
	res.ID = aux.ID
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
	return
}

// scans row and fills auxComposeModuleField fields
//
// This function is auto-generated
func (aux *auxComposeModuleField) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
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
	)
}

// encodes ComposeNamespace to auxComposeNamespace
//
// This function is auto-generated
func (aux *auxComposeNamespace) encode(res *composeType.Namespace) (_ error) {
	aux.ID = res.ID
	aux.Slug = res.Slug
	aux.Enabled = res.Enabled
	aux.Meta = res.Meta
	aux.Name = res.Name
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	return
}

// decodes ComposeNamespace from auxComposeNamespace
//
// This function is auto-generated
func (aux auxComposeNamespace) decode() (res *composeType.Namespace, _ error) {
	res = new(composeType.Namespace)
	res.ID = aux.ID
	res.Slug = aux.Slug
	res.Enabled = aux.Enabled
	res.Meta = aux.Meta
	res.Name = aux.Name
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	return
}

// scans row and fills auxComposeNamespace fields
//
// This function is auto-generated
func (aux *auxComposeNamespace) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
		&aux.Slug,
		&aux.Enabled,
		&aux.Meta,
		&aux.Name,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
	)
}

// encodes ComposePage to auxComposePage
//
// This function is auto-generated
func (aux *auxComposePage) encode(res *composeType.Page) (_ error) {
	aux.ID = res.ID
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
	return
}

// decodes ComposePage from auxComposePage
//
// This function is auto-generated
func (aux auxComposePage) decode() (res *composeType.Page, _ error) {
	res = new(composeType.Page)
	res.ID = aux.ID
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
	return
}

// scans row and fills auxComposePage fields
//
// This function is auto-generated
func (aux *auxComposePage) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
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
	)
}

// encodes ComposePageLayout to auxComposePageLayout
//
// This function is auto-generated
func (aux *auxComposePageLayout) encode(res *composeType.PageLayout) (_ error) {
	aux.ID = res.ID
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
	return
}

// decodes ComposePageLayout from auxComposePageLayout
//
// This function is auto-generated
func (aux auxComposePageLayout) decode() (res *composeType.PageLayout, _ error) {
	res = new(composeType.PageLayout)
	res.ID = aux.ID
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
	return
}

// scans row and fills auxComposePageLayout fields
//
// This function is auto-generated
func (aux *auxComposePageLayout) scan(row scanner) error {
	return row.Scan(
		&aux.ID,
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
	aux.Kind = res.Kind
	aux.Params = res.Params
	aux.CreatedAt = res.CreatedAt
	aux.UpdatedAt = res.UpdatedAt
	aux.DeletedAt = res.DeletedAt
	aux.CompletedAt = res.CompletedAt
	aux.CreatedBy = res.CreatedBy
	aux.UpdatedBy = res.UpdatedBy
	aux.DeletedBy = res.DeletedBy
	aux.CompletedBy = res.CompletedBy
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
	res.Kind = aux.Kind
	res.Params = aux.Params
	res.CreatedAt = aux.CreatedAt
	res.UpdatedAt = aux.UpdatedAt
	res.DeletedAt = aux.DeletedAt
	res.CompletedAt = aux.CompletedAt
	res.CreatedBy = aux.CreatedBy
	res.UpdatedBy = aux.UpdatedBy
	res.DeletedBy = aux.DeletedBy
	res.CompletedBy = aux.CompletedBy
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
		&aux.Kind,
		&aux.Params,
		&aux.CreatedAt,
		&aux.UpdatedAt,
		&aux.DeletedAt,
		&aux.CompletedAt,
		&aux.CreatedBy,
		&aux.UpdatedBy,
		&aux.DeletedBy,
		&aux.CompletedBy,
	)
}

// encodes DalSensitivityLevel to auxDalSensitivityLevel
//
// This function is auto-generated
func (aux *auxDalSensitivityLevel) encode(res *systemType.DalSensitivityLevel) (_ error) {
	aux.ID = res.ID
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

// encodes FederationExposedModule to auxFederationExposedModule
//
// This function is auto-generated
func (aux *auxFederationExposedModule) encode(res *federationType.ExposedModule) (_ error) {
	aux.ID = res.ID
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

// encodes Queue to auxQueue
//
// This function is auto-generated
func (aux *auxQueue) encode(res *systemType.Queue) (_ error) {
	aux.ID = res.ID
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
	aux.UserID = res.UserID
	aux.RoleID = res.RoleID
	return
}

// decodes RoleMember from auxRoleMember
//
// This function is auto-generated
func (aux auxRoleMember) decode() (res *systemType.RoleMember, _ error) {
	res = new(systemType.RoleMember)
	res.UserID = aux.UserID
	res.RoleID = aux.RoleID
	return
}

// scans row and fills auxRoleMember fields
//
// This function is auto-generated
func (aux *auxRoleMember) scan(row scanner) error {
	return row.Scan(
		&aux.UserID,
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

// encodes User to auxUser
//
// This function is auto-generated
func (aux *auxUser) encode(res *systemType.User) (_ error) {
	aux.ID = res.ID
	aux.Email = res.Email
	aux.EmailConfirmed = res.EmailConfirmed
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
	res.Email = aux.Email
	res.EmailConfirmed = aux.EmailConfirmed
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
		&aux.Email,
		&aux.EmailConfirmed,
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
