package types

import (
	"time"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ChatbotSession struct {
		ID           uint64     `json:"id,string"`
		ChatbotID    uint64     `json:"chatbotID,string"`
		Status       string     `json:"status"`
		CurrentStep  int        `json:"currentStep"`
		CreatedAt    time.Time  `json:"createdAt,omitempty"`
		CreatedBy    uint64     `json:"createdBy,string"`
		UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy    uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt    *time.Time `json:"deletedAt,omitempty"`
		DeletedBy    uint64     `json:"deletedBy,string,omitempty"`
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
