package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"strconv"
	"time"
)

type (
	Agent struct {
		ID         uint64                           `json:"agentID,string"`
		TenantID   uint64                           `json:"tenantID,string,omitempty"`
		ProjectID  uint64                           `json:"projectID,string,omitempty"`
		Handle     string                           `json:"handle"`
		Status     string                           `json:"status"`
		Revision   int                              `json:"revision"`
		Meta       AgentMeta                        `json:"meta"`
		Behavior   AgentBehavior                    `json:"behavior"`
		Execution  AgentExecution                   `json:"execution"`
		Access     AgentAccess                      `json:"access"`
		Invocation AgentInvocation                  `json:"invocation"`
		CreatedAt  time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy  uint64                           `json:"createdBy,string"`
		UpdatedBy  uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy  uint64                           `json:"deletedBy,string,omitempty"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
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
		TreatyCLEnabled     *bool               `json:"treatyCLEnabled"`
		TreatyCLTemperature int                 `json:"tclTemperature,omitempty"`
		TreatyCLArticles    []string            `json:"tclArticles,omitempty"`
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

	AgentAccessContext struct {
		Namespace string         `json:"namespace,omitempty"`
		Module    string         `json:"module,omitempty"`
		Defaults  map[string]any `json:"defaults,omitempty"`
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

	// AgentAccessIDList is a []uint64 that serializes each element as a JSON string.
	AgentAccessIDList []uint64
)

func (r Agent) Clone() *Agent {
	dup := r
	dup.Meta = *r.Meta.Clone()

	dup.Behavior = *r.Behavior.Clone()

	dup.Execution = *r.Execution.Clone()

	dup.Access = *r.Access.Clone()

	dup.Invocation = *r.Invocation.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
	return &dup
}

func (r Agent) Diff(cmp *Agent) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Agent{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "agentID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.Revision != cmp.Revision {
		out = append(out, &revisions.Change{Key: "revision", Old: []any{cmp.Revision}, New: []any{r.Revision}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Behavior.Diff(&cmp.Behavior) {
		c.Key = "behavior." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Execution.Diff(&cmp.Execution) {
		c.Key = "execution." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Access.Diff(&cmp.Access) {
		c.Key = "access." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Invocation.Diff(&cmp.Invocation) {
		c.Key = "invocation." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r AgentMeta) Clone() *AgentMeta {
	dup := r
	if r.SidebarRoles != nil {
		dup.SidebarRoles = make([]string, len(r.SidebarRoles))
		copy(dup.SidebarRoles, r.SidebarRoles)
	}

	return &dup
}

func (r AgentMeta) Diff(cmp *AgentMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if !reflect.DeepEqual(r.SidebarRoles, cmp.SidebarRoles) {
		out = append(out, &revisions.Change{Key: "sidebarRoles", Old: []any{cmp.SidebarRoles}, New: []any{r.SidebarRoles}})
	}

	return out
}

func (r *AgentMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentBehavior) Clone() *AgentBehavior {
	dup := r
	if r.Guardrails != nil {
		dup.Guardrails = make([]string, len(r.Guardrails))
		copy(dup.Guardrails, r.Guardrails)
	}

	if r.TreatyCLEnabled != nil {
		v := *r.TreatyCLEnabled
		dup.TreatyCLEnabled = &v
	}

	if r.TreatyCLArticles != nil {
		dup.TreatyCLArticles = make([]string, len(r.TreatyCLArticles))
		copy(dup.TreatyCLArticles, r.TreatyCLArticles)
	}

	return &dup
}

func (r AgentBehavior) Diff(cmp *AgentBehavior) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentBehavior{}
	}
	if r.SystemPrompt != cmp.SystemPrompt {
		out = append(out, &revisions.Change{Key: "systemPrompt", Old: []any{cmp.SystemPrompt}, New: []any{r.SystemPrompt}})
	}

	if !reflect.DeepEqual(r.Guardrails, cmp.Guardrails) {
		out = append(out, &revisions.Change{Key: "guardrails", Old: []any{cmp.Guardrails}, New: []any{r.Guardrails}})
	}

	if r.InjectSystemContext != cmp.InjectSystemContext {
		out = append(out, &revisions.Change{Key: "injectSystemContext", Old: []any{cmp.InjectSystemContext}, New: []any{r.InjectSystemContext}})
	}

	if !reflect.DeepEqual(r.KnowledgeBases, cmp.KnowledgeBases) {
		out = append(out, &revisions.Change{Key: "knowledgeBases", Old: []any{cmp.KnowledgeBases}, New: []any{r.KnowledgeBases}})
	}

	if !reflect.DeepEqual(r.TreatyCLEnabled, cmp.TreatyCLEnabled) {
		out = append(out, &revisions.Change{Key: "treatyCLEnabled", Old: []any{cmp.TreatyCLEnabled}, New: []any{r.TreatyCLEnabled}})
	}

	if r.TreatyCLTemperature != cmp.TreatyCLTemperature {
		out = append(out, &revisions.Change{Key: "tclTemperature", Old: []any{cmp.TreatyCLTemperature}, New: []any{r.TreatyCLTemperature}})
	}

	if !reflect.DeepEqual(r.TreatyCLArticles, cmp.TreatyCLArticles) {
		out = append(out, &revisions.Change{Key: "tclArticles", Old: []any{cmp.TreatyCLArticles}, New: []any{r.TreatyCLArticles}})
	}

	return out
}

func (r *AgentBehavior) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentBehavior) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentExecution) Clone() *AgentExecution {
	dup := r
	dup.Model = *r.Model.Clone()

	dup.Limits = *r.Limits.Clone()

	return &dup
}

func (r AgentExecution) Diff(cmp *AgentExecution) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentExecution{}
	}
	for _, c := range r.Model.Diff(&cmp.Model) {
		c.Key = "model." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Limits.Diff(&cmp.Limits) {
		c.Key = "limits." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *AgentExecution) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentExecution) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentExecutionModel) Clone() *AgentExecutionModel {
	dup := r
	if r.Temperature != nil {
		v := *r.Temperature
		dup.Temperature = &v
	}

	return &dup
}

func (r AgentExecutionModel) Diff(cmp *AgentExecutionModel) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentExecutionModel{}
	}
	if r.LLMProviderID != cmp.LLMProviderID {
		out = append(out, &revisions.Change{Key: "llmProviderID", Old: []any{cmp.LLMProviderID}, New: []any{r.LLMProviderID}})
	}

	if r.Model != cmp.Model {
		out = append(out, &revisions.Change{Key: "model", Old: []any{cmp.Model}, New: []any{r.Model}})
	}

	if !reflect.DeepEqual(r.Temperature, cmp.Temperature) {
		out = append(out, &revisions.Change{Key: "temperature", Old: []any{cmp.Temperature}, New: []any{r.Temperature}})
	}

	return out
}

func (r *AgentExecutionModel) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentExecutionModel) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentExecutionLimits) Clone() *AgentExecutionLimits {
	dup := r
	return &dup
}

func (r AgentExecutionLimits) Diff(cmp *AgentExecutionLimits) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentExecutionLimits{}
	}
	if r.MaxIterations != cmp.MaxIterations {
		out = append(out, &revisions.Change{Key: "maxIterations", Old: []any{cmp.MaxIterations}, New: []any{r.MaxIterations}})
	}

	if r.Timeout != cmp.Timeout {
		out = append(out, &revisions.Change{Key: "timeout", Old: []any{cmp.Timeout}, New: []any{r.Timeout}})
	}

	if r.SoftLimitRatio != cmp.SoftLimitRatio {
		out = append(out, &revisions.Change{Key: "softLimitRatio", Old: []any{cmp.SoftLimitRatio}, New: []any{r.SoftLimitRatio}})
	}

	if r.ContextWindow != cmp.ContextWindow {
		out = append(out, &revisions.Change{Key: "contextWindow", Old: []any{cmp.ContextWindow}, New: []any{r.ContextWindow}})
	}

	if r.OutputTokens != cmp.OutputTokens {
		out = append(out, &revisions.Change{Key: "outputTokens", Old: []any{cmp.OutputTokens}, New: []any{r.OutputTokens}})
	}

	return out
}

func (r *AgentExecutionLimits) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentExecutionLimits) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentAccess) Clone() *AgentAccess {
	dup := r
	dup.Context = *r.Context.Clone()

	if r.Tools != nil {
		dup.Tools = make([]AgentAccessTool, len(r.Tools))
		copy(dup.Tools, r.Tools)
	}

	if r.TAQs != nil {
		dup.TAQs = make([]AgentAccessTAQ, len(r.TAQs))
		copy(dup.TAQs, r.TAQs)
	}

	if r.Workflows != nil {
		dup.Workflows = make([]AgentAccessWorkflow, len(r.Workflows))
		copy(dup.Workflows, r.Workflows)
	}

	return &dup
}

func (r AgentAccess) Diff(cmp *AgentAccess) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentAccess{}
	}
	for _, c := range r.Context.Diff(&cmp.Context) {
		c.Key = "context." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Tools, cmp.Tools) {
		out = append(out, &revisions.Change{Key: "tools", Old: []any{cmp.Tools}, New: []any{r.Tools}})
	}

	if !reflect.DeepEqual(r.TAQs, cmp.TAQs) {
		out = append(out, &revisions.Change{Key: "taqs", Old: []any{cmp.TAQs}, New: []any{r.TAQs}})
	}

	if !reflect.DeepEqual(r.Workflows, cmp.Workflows) {
		out = append(out, &revisions.Change{Key: "workflows", Old: []any{cmp.Workflows}, New: []any{r.Workflows}})
	}

	return out
}

func (r *AgentAccess) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentAccess) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentAccessContext) Clone() *AgentAccessContext {
	dup := r
	if r.Defaults != nil {
		dup.Defaults = make(map[string]any, len(r.Defaults))
		for k, v := range r.Defaults {
			dup.Defaults[k] = v
		}
	}

	return &dup
}

func (r AgentAccessContext) Diff(cmp *AgentAccessContext) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentAccessContext{}
	}
	if r.Namespace != cmp.Namespace {
		out = append(out, &revisions.Change{Key: "namespace", Old: []any{cmp.Namespace}, New: []any{r.Namespace}})
	}

	if r.Module != cmp.Module {
		out = append(out, &revisions.Change{Key: "module", Old: []any{cmp.Module}, New: []any{r.Module}})
	}

	if !reflect.DeepEqual(r.Defaults, cmp.Defaults) {
		out = append(out, &revisions.Change{Key: "defaults", Old: []any{cmp.Defaults}, New: []any{r.Defaults}})
	}

	return out
}

func (r *AgentAccessContext) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentAccessContext) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentInvocation) Clone() *AgentInvocation {
	dup := r
	dup.User = *r.User.Clone()

	dup.System = *r.System.Clone()

	return &dup
}

func (r AgentInvocation) Diff(cmp *AgentInvocation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentInvocation{}
	}
	for _, c := range r.User.Diff(&cmp.User) {
		c.Key = "user." + c.Key
		out = append(out, c)
	}

	for _, c := range r.System.Diff(&cmp.System) {
		c.Key = "system." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *AgentInvocation) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentInvocation) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentInvocationUser) Clone() *AgentInvocationUser {
	dup := r
	return &dup
}

func (r AgentInvocationUser) Diff(cmp *AgentInvocationUser) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentInvocationUser{}
	}
	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	return out
}

func (r *AgentInvocationUser) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentInvocationUser) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AgentInvocationSystem) Clone() *AgentInvocationSystem {
	dup := r
	return &dup
}

func (r AgentInvocationSystem) Diff(cmp *AgentInvocationSystem) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AgentInvocationSystem{}
	}
	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.ServiceAccount != cmp.ServiceAccount {
		out = append(out, &revisions.Change{Key: "serviceAccount", Old: []any{cmp.ServiceAccount}, New: []any{r.ServiceAccount}})
	}

	if !reflect.DeepEqual(r.InputSchema, cmp.InputSchema) {
		out = append(out, &revisions.Change{Key: "inputSchema", Old: []any{cmp.InputSchema}, New: []any{r.InputSchema}})
	}

	if r.OutputFormat != cmp.OutputFormat {
		out = append(out, &revisions.Change{Key: "outputFormat", Old: []any{cmp.OutputFormat}, New: []any{r.OutputFormat}})
	}

	return out
}

func (r *AgentInvocationSystem) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AgentInvocationSystem) Value() (driver.Value, error) { return json.Marshal(r) }

func (ll AgentAccessIDList) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(ll))
	for i, id := range ll {
		ss[i] = strconv.FormatUint(id, 10)
	}
	return json.Marshal(ss)
}

func (ll *AgentAccessIDList) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ll = make(AgentAccessIDList, 0, len(raw))
	for _, r := range raw {
		s := string(r)
		if len(s) >= 2 && s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*ll = append(*ll, id)
	}
	return nil
}

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
