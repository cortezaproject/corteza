package observability

import "time"

type (
	SpanStatus string

	AgentSpan struct {
		ID             string
		ParentID       string
		TraceID        string
		Name           string
		AgentID        string
		UserID         string
		ConversationID string
		StartedAt      time.Time
		EndedAt        time.Time
		Attributes     map[string]any
		Status         SpanStatus
		Error          error
	}

	AgentEvent struct {
		ID             string
		TraceID        string
		SpanID         string
		Timestamp      time.Time
		Event          string
		AgentID        string
		UserID         string
		ConversationID string
		Details        map[string]any
	}
)

const (
	StatusOK    SpanStatus = "ok"
	StatusError SpanStatus = "error"
)
