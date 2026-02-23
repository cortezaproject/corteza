package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type langfuseDispatcher struct {
	host      string
	publicKey string
	secretKey string
	mu        sync.Mutex
	spans     map[string][]AgentSpan  // spans waiting to be flushed
	events    map[string][]AgentEvent // events waiting to be flushed
}

func NewLangfuseDispatcher(host, publicKey, secretKey string) *langfuseDispatcher {
	return &langfuseDispatcher{
		host:      host,
		publicKey: publicKey,
		secretKey: secretKey,
		spans:     make(map[string][]AgentSpan),
		events:    make(map[string][]AgentEvent),
	}
}

func (d *langfuseDispatcher) ID() string { return "langfuse" }

func (d *langfuseDispatcher) OnSpan(span AgentSpan) error {
	d.mu.Lock()
	d.spans[span.TraceID] = append(d.spans[span.TraceID], span)
	isRoot := span.ParentID == ""
	d.mu.Unlock()

	if isRoot {
		return d.flush(span.TraceID)
	}
	return nil
}

func (d *langfuseDispatcher) OnEvent(event AgentEvent) error {
	d.mu.Lock()
	d.events[event.TraceID] = append(d.events[event.TraceID], event)
	d.mu.Unlock()
	return nil
}

func (d *langfuseDispatcher) flush(traceID string) error {
	d.mu.Lock()
	spans := d.spans[traceID]
	events := d.events[traceID]
	delete(d.spans, traceID)
	delete(d.events, traceID)
	d.mu.Unlock()

	var root AgentSpan
	for _, s := range spans {
		if s.ParentID == "" {
			root = s
			break
		}
	}

	var batch []map[string]any

	batch = append(batch, map[string]any{
		"id":        root.ID,
		"type":      "trace-create",
		"timestamp": root.StartedAt.Format(time.RFC3339Nano),
		"body": map[string]any{
			"id":        root.TraceID,
			"name":      "agent.run",
			"userId":    root.UserID,
			"sessionId": root.ConversationID,
			"metadata":  map[string]any{"agentID": root.AgentID},
		},
	})

	for _, s := range spans {
		if s.ParentID == "" {
			continue
		}

		body := map[string]any{
			"id":        s.ID,
			"traceId":   s.TraceID,
			"name":      s.Name,
			"startTime": s.StartedAt.Format(time.RFC3339Nano),
			"endTime":   s.EndedAt.Format(time.RFC3339Nano),
			"level":     langfuseLevel(s.Status),
		}

		if errStr(s.Error) != "" {
			body["statusMessage"] = errStr(s.Error)
		}

		// children of root have no parent observation in Langfuse
		if s.ParentID != root.ID {
			body["parentObservationId"] = s.ParentID
		}

		itemType := "span-create"
		if s.Name == "llm.chat" {
			itemType = "generation-create"
			usage := map[string]any{}
			if s.Attributes != nil {
				if v, ok := s.Attributes["inputTokens"]; ok {
					usage["input"] = v
				}
				if v, ok := s.Attributes["outputTokens"]; ok {
					usage["output"] = v
				}
			}
			body["usage"] = usage
		} else if s.Attributes != nil {
			body["metadata"] = s.Attributes
		}

		batch = append(batch, map[string]any{
			"id":        s.ID + "-item",
			"type":      itemType,
			"timestamp": s.StartedAt.Format(time.RFC3339Nano),
			"body":      body,
		})
	}

	for _, ev := range events {
		body := map[string]any{
			"id":        ev.ID,
			"traceId":   ev.TraceID,
			"name":      ev.Event,
			"startTime": ev.Timestamp.Format(time.RFC3339Nano),
			"metadata":  ev.Details,
		}

		if ev.SpanID != root.ID {
			body["parentObservationId"] = ev.SpanID
		}

		batch = append(batch, map[string]any{
			"id":        ev.ID + "-item",
			"type":      "event-create",
			"timestamp": ev.Timestamp.Format(time.RFC3339Nano),
			"body":      body,
		})
	}

	return d.post(batch)
}

func (d *langfuseDispatcher) post(batch []map[string]any) error {
	body, err := json.Marshal(map[string]any{"batch": batch})
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, d.host+"/api/public/ingestion", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(d.publicKey, d.secretKey)

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("langfuse request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMultiStatus {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("langfuse returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (d *langfuseDispatcher) Flush() error    { return nil }
func (d *langfuseDispatcher) Shutdown() error { return nil }

func langfuseLevel(status SpanStatus) string {
	if status == StatusError {
		return "ERROR"
	}
	return "DEFAULT"
}
