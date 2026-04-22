package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"time"

	"github.com/crusttech/human/server/pkg/sql"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	Agent struct {
		ID       uint64 `json:"agentID,string"`
		Handle   string `json:"handle"`
		Status   string `json:"status"`
		Revision int    `json:"revision"`

		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Meta       AgentMeta       `json:"meta"`
		Behavior   AgentBehavior   `json:"behavior"`
		Execution  AgentExecution  `json:"execution"`
		Access     AgentAccess     `json:"access"`
		Invocation AgentInvocation `json:"invocation"`
		Chatbot    AgentChatbot    `json:"chatbot"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

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
		LLMProviderID uint64  `json:"llmProviderID,string,omitempty"`
		Model         string  `json:"model,omitempty"`
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
		Name    string                 `json:"name"`
		Description string             `json:"description"`
		Allow   []AgentAccessAllow     `json:"allow"`
		Context AgentAccessToolContext `json:"context,omitempty"`
	}

	AgentAccessToolContext struct {
		Defaults  map[string]any `json:"defaults,omitempty"`
		Overrides map[string]any `json:"overrides,omitempty"`
	}

	AgentAccessAllow struct {
		NamespaceID uint64            `json:"namespaceID,string"`
		ModuleIDs   AgentAccessIDList `json:"moduleIDs"`
	}

	// AgentAccessIDList is a []uint64 that serializes each element as a JSON string.
	AgentAccessIDList []uint64

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

func (m *AgentMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *AgentBehavior) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentBehavior) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *AgentExecution) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentExecution) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *AgentAccess) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentAccess) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *AgentInvocation) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentInvocation) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentMeta(ss []string) (p AgentMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseAgentBehavior(ss []string) (p AgentBehavior, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseAgentExecution(ss []string) (p AgentExecution, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseAgentAccess(ss []string) (p AgentAccess, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseAgentInvocation(ss []string) (p AgentInvocation, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (ll AgentAccessIDList) MarshalJSON() ([]byte, error) {
	strs := make([]string, len(ll))
	for i, id := range ll {
		strs[i] = strconv.FormatUint(id, 10)
	}
	return json.Marshal(strs)
}

func (ll *AgentAccessIDList) UnmarshalJSON(data []byte) error {
	var strs []string
	if err := json.Unmarshal(data, &strs); err != nil {
		return err
	}
	*ll = make(AgentAccessIDList, len(strs))
	for i, s := range strs {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		(*ll)[i] = id
	}
	return nil
}
