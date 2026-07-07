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
	ChatbotSession struct {
		ID          uint64              `json:"id,string"`
		TenantID    uint64              `json:"tenantID,string,omitempty"`
		ProjectID   uint64              `json:"projectID,string,omitempty"`
		ChatbotID   uint64              `json:"chatbotID,string"`
		Status      string              `json:"status"`
		CurrentStep int                 `json:"currentStep"`
		State       ChatbotSessionState `json:"state,omitempty"`
		CreatedAt   time.Time           `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time          `json:"deletedAt,omitempty"`
		CreatedBy   uint64              `json:"createdBy,string"`
		UpdatedBy   uint64              `json:"updatedBy,string,omitempty"`
		DeletedBy   uint64              `json:"deletedBy,string,omitempty"`
	}

	ChatbotSessionStepState struct {
		ScenarioID   string                        `json:"scenarioID"`
		Type         string                        `json:"type"`
		Conversation *ChatbotConversationStepState `json:"conversation,omitempty"`
		Form         *ChatbotFormStepState         `json:"form,omitempty"`
		Consent      *ChatbotConsentStepState      `json:"consent,omitempty"`
	}

	ChatbotConsentStepState struct {
		Accepted bool      `json:"accepted"`
		At       time.Time `json:"at"`
	}

	ChatbotConversationStepState struct {
		ConversationID uint64                           `json:"conversationID,string,omitempty"`
		History        []AiConversationMessage          `json:"history,omitempty"`
		Handoff        *ChatbotConversationHandoffState `json:"handoff,omitempty"`
	}

	ChatbotConversationHandoffState struct {
		OperatorID   uint64 `json:"operatorID,string,omitempty"`
		OperatorName string `json:"operatorName,omitempty"`
	}

	ChatbotFormStepState struct {
		Fields    map[string]string `json:"fields,omitempty"`
		Submitted bool              `json:"submitted,omitempty"`
	}
)

func (r ChatbotSession) Clone() *ChatbotSession {
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

func (r ChatbotSession) Diff(cmp *ChatbotSession) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotSession{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "id", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.ChatbotID != cmp.ChatbotID {
		out = append(out, &revisions.Change{Key: "chatbotID", Old: []any{cmp.ChatbotID}, New: []any{r.ChatbotID}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.CurrentStep != cmp.CurrentStep {
		out = append(out, &revisions.Change{Key: "currentStep", Old: []any{cmp.CurrentStep}, New: []any{r.CurrentStep}})
	}

	if !reflect.DeepEqual(r.State, cmp.State) {
		out = append(out, &revisions.Change{Key: "state", Old: []any{cmp.State}, New: []any{r.State}})
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

func (r ChatbotSessionStepState) Clone() *ChatbotSessionStepState {
	dup := r
	if r.Conversation != nil {
		dup.Conversation = r.Conversation.Clone()
	}

	if r.Form != nil {
		dup.Form = r.Form.Clone()
	}

	if r.Consent != nil {
		dup.Consent = r.Consent.Clone()
	}

	return &dup
}

func (r ChatbotSessionStepState) Diff(cmp *ChatbotSessionStepState) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotSessionStepState{}
	}
	if r.ScenarioID != cmp.ScenarioID {
		out = append(out, &revisions.Change{Key: "scenarioID", Old: []any{cmp.ScenarioID}, New: []any{r.ScenarioID}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if (r.Conversation == nil) != (cmp.Conversation == nil) {
		out = append(out, &revisions.Change{Key: "conversation", Old: []any{cmp.Conversation}, New: []any{r.Conversation}})
	} else if r.Conversation != nil {
		for _, c := range r.Conversation.Diff(cmp.Conversation) {
			c.Key = "conversation." + c.Key
			out = append(out, c)
		}
	}

	if (r.Form == nil) != (cmp.Form == nil) {
		out = append(out, &revisions.Change{Key: "form", Old: []any{cmp.Form}, New: []any{r.Form}})
	} else if r.Form != nil {
		for _, c := range r.Form.Diff(cmp.Form) {
			c.Key = "form." + c.Key
			out = append(out, c)
		}
	}

	if (r.Consent == nil) != (cmp.Consent == nil) {
		out = append(out, &revisions.Change{Key: "consent", Old: []any{cmp.Consent}, New: []any{r.Consent}})
	} else if r.Consent != nil {
		for _, c := range r.Consent.Diff(cmp.Consent) {
			c.Key = "consent." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ChatbotSessionStepState) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotSessionStepState) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotConsentStepState) Clone() *ChatbotConsentStepState {
	dup := r
	return &dup
}

func (r ChatbotConsentStepState) Diff(cmp *ChatbotConsentStepState) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotConsentStepState{}
	}
	if r.Accepted != cmp.Accepted {
		out = append(out, &revisions.Change{Key: "accepted", Old: []any{cmp.Accepted}, New: []any{r.Accepted}})
	}

	if !reflect.DeepEqual(r.At, cmp.At) {
		out = append(out, &revisions.Change{Key: "at", Old: []any{cmp.At}, New: []any{r.At}})
	}

	return out
}

func (r *ChatbotConsentStepState) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotConsentStepState) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotConversationStepState) Clone() *ChatbotConversationStepState {
	dup := r
	if r.History != nil {
		dup.History = make([]AiConversationMessage, len(r.History))
		copy(dup.History, r.History)
	}

	if r.Handoff != nil {
		dup.Handoff = r.Handoff.Clone()
	}

	return &dup
}

func (r ChatbotConversationStepState) Diff(cmp *ChatbotConversationStepState) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotConversationStepState{}
	}
	if r.ConversationID != cmp.ConversationID {
		out = append(out, &revisions.Change{Key: "conversationID", Old: []any{cmp.ConversationID}, New: []any{r.ConversationID}})
	}

	if !reflect.DeepEqual(r.History, cmp.History) {
		out = append(out, &revisions.Change{Key: "history", Old: []any{cmp.History}, New: []any{r.History}})
	}

	if (r.Handoff == nil) != (cmp.Handoff == nil) {
		out = append(out, &revisions.Change{Key: "handoff", Old: []any{cmp.Handoff}, New: []any{r.Handoff}})
	} else if r.Handoff != nil {
		for _, c := range r.Handoff.Diff(cmp.Handoff) {
			c.Key = "handoff." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ChatbotConversationStepState) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotConversationStepState) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotConversationHandoffState) Clone() *ChatbotConversationHandoffState {
	dup := r
	return &dup
}

func (r ChatbotConversationHandoffState) Diff(cmp *ChatbotConversationHandoffState) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotConversationHandoffState{}
	}
	if r.OperatorID != cmp.OperatorID {
		out = append(out, &revisions.Change{Key: "operatorID", Old: []any{cmp.OperatorID}, New: []any{r.OperatorID}})
	}

	if r.OperatorName != cmp.OperatorName {
		out = append(out, &revisions.Change{Key: "operatorName", Old: []any{cmp.OperatorName}, New: []any{r.OperatorName}})
	}

	return out
}

func (r *ChatbotConversationHandoffState) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotConversationHandoffState) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotFormStepState) Clone() *ChatbotFormStepState {
	dup := r
	if r.Fields != nil {
		dup.Fields = make(map[string]string, len(r.Fields))
		for k, v := range r.Fields {
			dup.Fields[k] = v
		}
	}

	return &dup
}

func (r ChatbotFormStepState) Diff(cmp *ChatbotFormStepState) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotFormStepState{}
	}
	if !reflect.DeepEqual(r.Fields, cmp.Fields) {
		out = append(out, &revisions.Change{Key: "fields", Old: []any{cmp.Fields}, New: []any{r.Fields}})
	}

	if r.Submitted != cmp.Submitted {
		out = append(out, &revisions.Change{Key: "submitted", Old: []any{cmp.Submitted}, New: []any{r.Submitted}})
	}

	return out
}

func (r *ChatbotFormStepState) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotFormStepState) Value() (driver.Value, error) { return json.Marshal(r) }

func (m *ChatbotSessionState) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotSessionState) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotSessionState(ss []string) (p ChatbotSessionState, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
