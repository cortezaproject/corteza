package types

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	// ExecutableID uniquely identifies an executable (workflow/automation)
	ExecutableID = id.ID
	Status       string

	ExecutableStatus int

	ExecutableLimits struct {
		// Max operations allowed per single request/step
		MaxOpsPerRequest int

		// Hard cap for total operations during the entire execution
		TotalOps int

		// Rate limit for operations (ops per window)
		RateOps    int
		RateWindow time.Duration
	}

	// Step represents a single step in an executable workflow
	Step struct {
		ID       id.ID
		Kind     string
		Config   map[string]any
		Children []id.ID
		Handler  stepHandler
	}

	stepHandler interface {
		Execute(ctx context.Context, scope map[string]any) (outScope map[string]any, err error)
	}

	// Relationship defines a connection between steps in a workflow
	Relationship struct {
		From      id.ID // Empty string indicates entry point
		To        id.ID
		Condition string // Optional condition for conditional branching
	}

	// Executable represents a complete workflow definition
	Executable struct {
		ID            ExecutableID
		Revision      int
		Label         string
		Description   string
		Limits        ExecutableLimits
		Steps         []Step
		Relationships []Relationship
		Metadata      map[string]any
	}

	Execution struct {
		ID           id.ID
		ExecutableID id.ID
		Revision     uint32
		Status       Status
		CreatedAt    time.Time
		UpdatedAt    time.Time
		EndedAt      *time.Time
		Events       []StepEvent
		Variables    map[string]any
	}

	StepEvent struct {
		Type      EventType
		StepID    id.ID
		Timestamp time.Time
		Payload   any
		Error error
	}

	EventType string
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
