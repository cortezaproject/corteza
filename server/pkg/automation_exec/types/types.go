package types

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	Executable struct {
		ID       id.ID
		Revision int
		Handle   string

		RunAs auth.Identifiable

		Steps []Step

		// @todo
		// Limits        ExecutableLimits
	}

	// Expr represents an expression mapping for arguments/results
	// Cloned from automation/types.Expr to avoid circular dependencies
	Expr struct {
		// Context defines the evaluation context of this expression
		// Leave empty for global context
		Context string `json:"context,omitempty"`

		// Variable name to set results of the expression to
		Target string `json:"target"`

		// Source of the value / name of the variable from scope
		// Takes precedence over Value
		Source string `json:"source,omitempty"`

		// Expression to evaluate over the input variables
		Expression string `json:"expr,omitempty"`

		// Raw value to be set to target
		// If expression is set and fails, evaluation defaults to value
		Value interface{} `json:"value,omitempty"`

		// Expected type of the input value
		Type string `json:"type,omitempty"`
	}

	StepArg struct {
		*Expr
	}

	StepRst struct {
		Name string
	}

	// Step represents a single step in an executable workflow
	Step struct {
		ID     id.ID
		Handle string
		Kind   string

		Children []Step
		Parents  []Step

		Arguments []StepArg
		Results   []StepRst

		Handler StepHandler
	}

	ExecRequest struct {
		// Current scope
		Scope map[string]*expr.Vars
	}

	ExecResponse any

	ExecutableLimits struct{}

	// Executable represents a complete workflow definition

	Execution struct {
		ID           id.ID
		ExecutableID id.ID
		Revision     int
		Status       Status
		CreatedAt    time.Time
		UpdatedAt    time.Time
		EndedAt      *time.Time

		Error error

		Events []StepEvent
	}

	StepEvent struct {
		Type      EventType
		StepID    id.ID
		Timestamp time.Time
		Payload   any
		Error     error
	}

	StepHandler interface {
		ExecN(context.Context, *ExecRequest) (ExecResponse, error)
	}

	Budget struct {
		MaxOps int
		Window time.Duration // 0 = hard cap
	}

	RateLimit struct {
		MaxOps int
		Window time.Duration
	}

	EventType        string
	Status           string
	ExecutableStatus int

	ExecutionResult struct {
		ExecutionID  id.ID      `json:"executionID"`
		ExecutableID id.ID      `json:"executableID"`
		Revision     int        `json:"revision"`
		Status       Status     `json:"status"`
		Error        error      `json:"error,omitempty"`
		StartedAt    time.Time  `json:"startedAt"`
		EndedAt      *time.Time `json:"endedAt,omitempty"`
		Duration     string     `json:"duration"`
	}
)

const (
	StatusActive ExecutableStatus = iota
	StatusDeprecated
)

const (
	StatusCreated   Status = "created"
	StatusRunning   Status = "running"
	StatusPaused    Status = "paused"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

const (
	EventStepStarted   EventType = "step_started"
	EventStepCompleted EventType = "step_completed"
	EventStepFailed    EventType = "step_failed"
)

func (s ExecutableStatus) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusDeprecated:
		return "deprecated"
	default:
		return "unknown"
	}
}
