package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	Agent struct {
		ID       uint64 `json:"agentID,string"`
		Handle   string `json:"handle"`
		Status   string `json:"status"`
		Revision int    `json:"revision"`

		Meta       AgentMeta       `json:"meta"`
		Behavior   AgentBehavior   `json:"behavior"`
		Execution  AgentExecution  `json:"execution"`
		Access     AgentAccess     `json:"access"`
		Invocation AgentInvocation `json:"invocation"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	AgentMeta struct {
		Short       string `json:"short"`
		Description string `json:"description,omitempty"`
	}

	AgentBehavior struct {
		SystemPrompt string   `json:"systemPrompt,omitempty"`
		Guardrails   []string `json:"guardrails,omitempty"`
	}

	AgentExecution struct {
		Model  AgentExecutionModel  `json:"model"`
		Limits AgentExecutionLimits `json:"limits"`
	}

	AgentExecutionModel struct {
		LLMProviderID uint64  `json:"llmProviderID,string,omitempty"`
		Model         string  `json:"model,omitempty"`
		Temperature   float64 `json:"temperature,omitempty"`
		MaxTokens     int     `json:"maxTokens,omitempty"`
	}

	AgentExecutionLimits struct {
		MaxIterations  int     `json:"maxIterations,omitempty"`
		MaxTokens      int     `json:"maxTokens,omitempty"`
		Timeout        string  `json:"timeout,omitempty"`
		SoftLimitRatio float64 `json:"softLimitRatio,omitempty"`
	}

	AgentAccess struct {
		Context AgentAccessContext `json:"context"`
		Tools   []AgentAccessTool  `json:"tools,omitempty"`
		Allow   []AgentAccessAllow `json:"allow,omitempty"`
	}

	AgentAccessContext struct {
		Namespace string         `json:"namespace,omitempty"`
		Module    string         `json:"module,omitempty"`
		Defaults  map[string]any `json:"defaults,omitempty"`
	}

	AgentAccessTool struct {
		Name    string                 `json:"name"`
		Hints   string                 `json:"hints,omitempty"`
		Context AgentAccessToolContext `json:"context,omitempty"`
	}

	AgentAccessToolContext struct {
		Defaults  map[string]any `json:"defaults,omitempty"`
		Overrides map[string]any `json:"overrides,omitempty"`
	}

	AgentAccessAllow struct {
		Resource   string                     `json:"resource"`
		Filter     string                     `json:"filter,omitempty"`
		Properties []AgentAccessAllowProperty `json:"properties,omitempty"`
	}

	AgentAccessAllowProperty struct {
		Name   string `json:"name"`
		Access string `json:"access"`
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
		ServiceAccount string          `json:"serviceAccount,omitempty"`
		InputSchema    json.RawMessage `json:"inputSchema,omitempty"`
		OutputFormat   string          `json:"outputFormat,omitempty"`
	}

	AgentFilter struct {
		AgentID []string `json:"agentID"`
		Handle  string   `json:"handle"`
		Status  string   `json:"status"`
		Query   string   `json:"query"`

		Deleted filter.State `json:"deleted"`

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
