package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	AgentMeta struct {
		Short        string   `json:"short"`
		Description  string   `json:"description,omitempty"`
		SidebarRoles []string `json:"sidebarRoles,omitempty"`
	}

	AgentBehavior struct {
		SystemPrompt        string              `json:"systemPrompt,omitempty"`
		Guardrails          []string            `json:"guardrails,omitempty"`
		InjectSystemContext bool                `json:"injectSystemContext"`
		KnowledgeBases      KnowledgeBaseIDList `json:"knowledgeBases,omitempty"`
		TreatyCLEnabled     *bool               `json:"treatyCLEnabled"`          // nil = enabled by default
		TreatyCLTemperature int                 `json:"tclTemperature,omitempty"` // 1–10, controls citation verbosity
		TreatyCLArticles    []string            `json:"tclArticles,omitempty"`    // selected article IDs
	}

	AgentExecution struct {
		Model  AgentExecutionModel  `json:"model"`
		Limits AgentExecutionLimits `json:"limits"`
	}

	AgentExecutionModel struct {
		LLMProviderID uint64   `json:"llmProviderID,string,omitempty"`
		Model         string   `json:"model,omitempty"`
		Temperature   *float64 `json:"temperature,omitempty"`
	}

	AgentExecutionLimits struct {
		MaxIterations  int     `json:"maxIterations,omitempty"`
		Timeout        string  `json:"timeout"`
		SoftLimitRatio float64 `json:"softLimitRatio,omitempty"`
		ContextWindow  int     `json:"contextWindow"`
		OutputTokens   int     `json:"outputTokens"`
	}

	AgentAccess struct {
		Context   AgentAccessContext    `json:"context"`
		Tools     []AgentAccessTool     `json:"tools,omitempty"`
		TAQs      []AgentAccessTAQ      `json:"taqs,omitempty"`
		Workflows []AgentAccessWorkflow `json:"workflows,omitempty"`
	}

	AgentAccessTAQ struct {
		ID          uint64            `json:"id,string"`
		Description string            `json:"description,omitempty"`
		Params      map[string]string `json:"params,omitempty"`
	}

	AgentAccessWorkflow struct {
		ID          uint64 `json:"id,string"`
		Description string `json:"description,omitempty"`
	}

	AgentAccessContext struct {
		Namespace string         `json:"namespace,omitempty"`
		Module    string         `json:"module,omitempty"`
		Defaults  map[string]any `json:"defaults,omitempty"`
	}

	AgentAccessTool struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Allow       []AgentAccessAllow     `json:"allow"`
		Context     AgentAccessToolContext `json:"context,omitempty"`
	}

	AgentAccessToolContext struct {
		Defaults  map[string]any `json:"defaults,omitempty"`
		Overrides map[string]any `json:"overrides,omitempty"`
	}

	AgentAccessAllow struct {
		NamespaceID uint64            `json:"namespaceID,string"`
		ModuleIDs   AgentAccessIDList `json:"moduleIDs"`
	}

	AgentInvocation struct {
		User   AgentInvocationUser   `json:"user"`
		System AgentInvocationSystem `json:"system"`
	}

	AgentInvocationUser struct {
		Enabled bool `json:"enabled"`
	}

	AgentInvocationSystem struct {
		Enabled        bool            `json:"enabled"`
		ServiceAccount uint64          `json:"serviceAccount,string,omitempty"`
		InputSchema    json.RawMessage `json:"inputSchema,omitempty"`
		OutputFormat   string          `json:"outputFormat,omitempty"`
	}

	AgentFilter struct {
		AgentID []string `json:"agentID"`
		Handle  string   `json:"handle"`
		Status  string   `json:"status"`
		Query   string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Agent) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
