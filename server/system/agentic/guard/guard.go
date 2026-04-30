package guard

import (
	"context"

	"github.com/crusttech/human/server/system/types"
)

type (
	// GuardService evaluates user inputs for prompt injection attempts.
	// Both the built-in guard and external providers implement this interface.
	GuardService interface {
		CheckInput(ctx context.Context, input string, history []types.AiConversationMessage) (*GuardResult, error)
	}

	// GuardResult contains the classification outcome from a guard check.
	GuardResult struct {
		Safe       bool               `json:"safe"`
		Categories map[string]float64 `json:"categories,omitempty"`
		Blocked    bool               `json:"blocked"`
		Reason     string             `json:"reason,omitempty"`
	}
)
