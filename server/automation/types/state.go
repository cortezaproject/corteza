package types

import (
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"time"
)

type (
	// WorkflowState tracks suspended sessions
	// Session can have more than one state
	State struct {
		ID        uint64 `json:"stateID,string"`
		SessionID uint64 `json:"sessionID,string"`

		ResumeAt        *time.Time `json:"resumeAt"`
		WaitingForInput bool       `json:"waitingForInput"`

		CreatedAt time.Time `json:"createdAt,omitempty"`
		CreatedBy uint64    `json:"createdBy,string"`

		CallerID uint64     `json:"callerID,string"`
		StepID   uint64     `json:"stepID,string"`
		Scope    *expr.Vars `json:"scope"`
	}
)
