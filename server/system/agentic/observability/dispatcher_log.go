package observability

import (
	"encoding/json"
	"io"
	"os"
)

type logDispatcher struct {
	out io.Writer
}

func NewLogDispatcher() *logDispatcher {
	return &logDispatcher{out: os.Stdout}
}

func (d *logDispatcher) ID() string { return "log" }

func (d *logDispatcher) OnSpan(span AgentSpan) error {
	return json.NewEncoder(d.out).Encode(map[string]any{
		"type":           "span",
		"id":             span.ID,
		"parentID":       span.ParentID,
		"traceID":        span.TraceID,
		"name":           span.Name,
		"agentID":        span.AgentID,
		"userID":         span.UserID,
		"conversationID": span.ConversationID,
		"startedAt":      span.StartedAt,
		"endedAt":        span.EndedAt,
		"durationMs":     span.EndedAt.Sub(span.StartedAt).Milliseconds(),
		"attributes":     span.Attributes,
		"status":         span.Status,
		"error":          errStr(span.Error),
	})
}

func (d *logDispatcher) OnEvent(event AgentEvent) error {
	return json.NewEncoder(d.out).Encode(map[string]any{
		"type":           "event",
		"id":             event.ID,
		"traceID":        event.TraceID,
		"spanID":         event.SpanID,
		"timestamp":      event.Timestamp,
		"event":          event.Event,
		"agentID":        event.AgentID,
		"userID":         event.UserID,
		"conversationID": event.ConversationID,
		"details":        event.Details,
	})
}

func (d *logDispatcher) Flush() error    { return nil }
func (d *logDispatcher) Shutdown() error { return nil }

func errStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
