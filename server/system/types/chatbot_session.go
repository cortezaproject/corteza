package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/sql"
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
		CreatedBy   uint64              `json:"createdBy,string"`
		UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
		UpdatedBy   uint64              `json:"updatedBy,string,omitempty"`
		DeletedAt   *time.Time          `json:"deletedAt,omitempty"`
		DeletedBy   uint64              `json:"deletedBy,string,omitempty"`
	}

	// ChatbotSessionState is per-scenario state for a session. Discriminated
	// by Type — only the matching sub-field is populated. Stored as a JSON
	// blob on the chatbot_sessions row.
	ChatbotSessionState []ChatbotSessionStepState

	ChatbotSessionStepState struct {
		ScenarioID string `json:"scenarioID"`
		Type       string `json:"type"`

		Conversation *ChatbotConversationStepState `json:"conversation,omitempty"`
		Form         *ChatbotFormStepState         `json:"form,omitempty"`
		Consent      *ChatbotConsentStepState      `json:"consent,omitempty"`
	}

	// ChatbotConsentStepState captures the visitor's response to a consent
	// step (e.g. TOS acceptance) for audit purposes. Decision is final per
	// step; At is set to the server-side timestamp at submit time.
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

	ChatbotSessionFilter struct {
		ChatbotSessionID []uint64 `json:"chatbotSessionID"`
		ChatbotID        uint64   `json:"chatbotID"`
		Status           []string `json:"status"`
		Query            string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		Check func(*ChatbotSession) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	ChatbotSessionStep struct {
		ID             uint64     `json:"id,string"`
		SessionID      uint64     `json:"sessionID,string"`
		ScenarioIndex  int        `json:"scenarioIndex"`
		ConversationID uint64     `json:"conversationID,string"`
		Status         string     `json:"status"`
		CreatedAt      time.Time  `json:"createdAt,omitempty"`
		CreatedBy      uint64     `json:"createdBy,string"`
		UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy      uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt      *time.Time `json:"deletedAt,omitempty"`
		DeletedBy      uint64     `json:"deletedBy,string,omitempty"`
	}

	ChatbotSessionStepFilter struct {
		ChatbotSessionStepID []uint64 `json:"chatbotSessionStepID"`
		SessionID            uint64   `json:"sessionID"`
		ConversationID       uint64   `json:"conversationID"`
		Status               string   `json:"status"`
		Query                string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		Check func(*ChatbotSessionStep) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	ChatbotSessionHandoff struct {
		ID          uint64     `json:"id,string"`
		SessionID   uint64     `json:"sessionID,string"`
		StepID      uint64     `json:"stepID,string"`
		Status      string     `json:"status"`
		InitiatedAt time.Time  `json:"initiatedAt,omitempty"`
		ClosedAt    *time.Time `json:"closedAt,omitempty"`
		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		CreatedBy   uint64     `json:"createdBy,string"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy   uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty"`
		DeletedBy   uint64     `json:"deletedBy,string,omitempty"`
	}

	ChatbotSessionHandoffFilter struct {
		ChatbotSessionHandoffID []uint64 `json:"chatbotSessionHandoffID"`
		SessionID               uint64   `json:"sessionID"`
		StepID                  uint64   `json:"stepID"`
		Status                  string   `json:"status"`
		Query                   string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		Check func(*ChatbotSessionHandoff) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (s *ChatbotSessionState) Scan(src any) error          { return sql.ParseJSON(src, s) }
func (s ChatbotSessionState) Value() (driver.Value, error) { return json.Marshal(s) }

// ForScenario returns the existing per-scenario state entry, or nil.
func (s ChatbotSessionState) ForScenario(scenarioID string) *ChatbotSessionStepState {
	for i := range s {
		if s[i].ScenarioID == scenarioID {
			return &s[i]
		}
	}
	return nil
}

// Upsert returns the entry for scenarioID, creating it (with the given type)
// if missing. Returned pointer is valid for in-place mutation.
func (s *ChatbotSessionState) Upsert(scenarioID, typ string) *ChatbotSessionStepState {
	if e := s.ForScenario(scenarioID); e != nil {
		return e
	}
	*s = append(*s, ChatbotSessionStepState{ScenarioID: scenarioID, Type: typ})
	return &(*s)[len(*s)-1]
}
