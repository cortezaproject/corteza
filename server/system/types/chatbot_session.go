package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	// ChatbotSessionState is per-scenario state for a session. Discriminated
	// by Type — only the matching sub-field is populated. Stored as a JSON
	// blob on the chatbot_sessions row.
	ChatbotSessionState []ChatbotSessionStepState

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
