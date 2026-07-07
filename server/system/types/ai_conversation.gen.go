package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	AiConversation struct {
		ID         uint64                 `json:"aiConversationID,string"`
		TenantID   uint64                 `json:"tenantID,string,omitempty"`
		ProjectID  uint64                 `json:"projectID,string,omitempty"`
		AgentID    uint64                 `json:"agentID,string"`
		Messages   AiConversationMessages `json:"messages"`
		TokenCount int                    `json:"tokenCount"`
		CreatedAt  time.Time              `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time             `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time             `json:"deletedAt,omitempty"`
		CreatedBy  uint64                 `json:"createdBy,string"`
		UpdatedBy  uint64                 `json:"updatedBy,string,omitempty"`
		DeletedBy  uint64                 `json:"deletedBy,string,omitempty"`
	}

	AiConversationMessage struct {
		Role        string                     `json:"role"`
		Content     string                     `json:"content"`
		Operator    string                     `json:"operator,omitempty"`
		ToolCalls   []AiConversationToolCall   `json:"toolCalls,omitempty"`
		ToolResults []AiConversationToolResult `json:"toolResults,omitempty"`
	}

	AiConversationToolCall struct {
		CallID string `json:"callID"`
		Name   string `json:"name"`
		Data   string `json:"data"`
	}

	AiConversationToolResult struct {
		CallID string `json:"callID"`
		Data   string `json:"data"`
		Error  string `json:"error,omitempty"`
	}
)

func (r AiConversation) Clone() *AiConversation {
	dup := r
	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	return &dup
}

func (r AiConversation) Diff(cmp *AiConversation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AiConversation{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "aiConversationID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.AgentID != cmp.AgentID {
		out = append(out, &revisions.Change{Key: "agentID", Old: []any{cmp.AgentID}, New: []any{r.AgentID}})
	}

	if !reflect.DeepEqual(r.Messages, cmp.Messages) {
		out = append(out, &revisions.Change{Key: "messages", Old: []any{cmp.Messages}, New: []any{r.Messages}})
	}

	if r.TokenCount != cmp.TokenCount {
		out = append(out, &revisions.Change{Key: "tokenCount", Old: []any{cmp.TokenCount}, New: []any{r.TokenCount}})
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

	return out
}

func (r AiConversationMessage) Clone() *AiConversationMessage {
	dup := r
	if r.ToolCalls != nil {
		dup.ToolCalls = make([]AiConversationToolCall, len(r.ToolCalls))
		for i := range r.ToolCalls {
			dup.ToolCalls[i] = *r.ToolCalls[i].Clone()
		}
	}

	if r.ToolResults != nil {
		dup.ToolResults = make([]AiConversationToolResult, len(r.ToolResults))
		for i := range r.ToolResults {
			dup.ToolResults[i] = *r.ToolResults[i].Clone()
		}
	}

	return &dup
}

func (r AiConversationMessage) Diff(cmp *AiConversationMessage) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AiConversationMessage{}
	}
	if r.Role != cmp.Role {
		out = append(out, &revisions.Change{Key: "role", Old: []any{cmp.Role}, New: []any{r.Role}})
	}

	if r.Content != cmp.Content {
		out = append(out, &revisions.Change{Key: "content", Old: []any{cmp.Content}, New: []any{r.Content}})
	}

	if r.Operator != cmp.Operator {
		out = append(out, &revisions.Change{Key: "operator", Old: []any{cmp.Operator}, New: []any{r.Operator}})
	}

	if !reflect.DeepEqual(r.ToolCalls, cmp.ToolCalls) {
		out = append(out, &revisions.Change{Key: "toolCalls", Old: []any{cmp.ToolCalls}, New: []any{r.ToolCalls}})
	}

	if !reflect.DeepEqual(r.ToolResults, cmp.ToolResults) {
		out = append(out, &revisions.Change{Key: "toolResults", Old: []any{cmp.ToolResults}, New: []any{r.ToolResults}})
	}

	return out
}

func (r *AiConversationMessage) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AiConversationMessage) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AiConversationToolCall) Clone() *AiConversationToolCall {
	dup := r
	return &dup
}

func (r AiConversationToolCall) Diff(cmp *AiConversationToolCall) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AiConversationToolCall{}
	}
	if r.CallID != cmp.CallID {
		out = append(out, &revisions.Change{Key: "callID", Old: []any{cmp.CallID}, New: []any{r.CallID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Data != cmp.Data {
		out = append(out, &revisions.Change{Key: "data", Old: []any{cmp.Data}, New: []any{r.Data}})
	}

	return out
}

func (r *AiConversationToolCall) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AiConversationToolCall) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AiConversationToolResult) Clone() *AiConversationToolResult {
	dup := r
	return &dup
}

func (r AiConversationToolResult) Diff(cmp *AiConversationToolResult) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AiConversationToolResult{}
	}
	if r.CallID != cmp.CallID {
		out = append(out, &revisions.Change{Key: "callID", Old: []any{cmp.CallID}, New: []any{r.CallID}})
	}

	if r.Data != cmp.Data {
		out = append(out, &revisions.Change{Key: "data", Old: []any{cmp.Data}, New: []any{r.Data}})
	}

	if r.Error != cmp.Error {
		out = append(out, &revisions.Change{Key: "error", Old: []any{cmp.Error}, New: []any{r.Error}})
	}

	return out
}

func (r *AiConversationToolResult) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AiConversationToolResult) Value() (driver.Value, error) { return json.Marshal(r) }

func (m *AiConversationMessages) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AiConversationMessages) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAiConversationMessages(ss []string) (p AiConversationMessages, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
