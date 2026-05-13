package widget

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/crusttech/human/server/system/agentic/observability"
)

// sseDispatcher is a transient observability.Dispatcher scoped to one SSE stream.
// It is filtered to a single ConversationID so cross-tenant leaks are impossible.
type sseDispatcher struct {
	id        string
	convID    string
	events    chan observability.AgentEvent
	closed    chan struct{}
	closeOnce sync.Once
}

func newSSEDispatcher(convID uint64) *sseDispatcher {
	return &sseDispatcher{
		id:     fmt.Sprintf("sse-%d-%d", convID, time.Now().UnixNano()),
		convID: strconv.FormatUint(convID, 10),
		events: make(chan observability.AgentEvent, 64),
		closed: make(chan struct{}),
	}
}

func (d *sseDispatcher) ID() string { return d.id }

func (d *sseDispatcher) OnSpan(_ observability.AgentSpan) error   { return nil }
func (d *sseDispatcher) Flush() error                              { return nil }
func (d *sseDispatcher) Shutdown() error                           { d.close(); return nil }

func (d *sseDispatcher) OnEvent(e observability.AgentEvent) error {
	if e.ConversationID != d.convID {
		return nil
	}
	select {
	case <-d.closed:
		return nil
	case d.events <- e:
	default:
		// drop on backpressure — widget UX can re-sync from conversation store
	}
	return nil
}

func (d *sseDispatcher) close() {
	d.closeOnce.Do(func() { close(d.closed) })
}

// pumpSSE writes SSE frames until the client disconnects or the context
// deadline hits. Keepalive comment frames are sent every 15s so proxies don't
// sever the connection.
func pumpSSE(w http.ResponseWriter, r *http.Request, d *sseDispatcher) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "widget: streaming unsupported", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache, no-transform")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ka := time.NewTicker(15 * time.Second)
	defer ka.Stop()

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		case <-d.closed:
			return
		case <-ka.C:
			_, _ = w.Write([]byte(":\n\n"))
			flusher.Flush()
		case ev := <-d.events:
			writeSSEEvent(w, ev)
			flusher.Flush()
		}
	}
}

func writeSSEEvent(w http.ResponseWriter, ev observability.AgentEvent) {
	name := ev.Event
	if name == "" {
		name = "message"
	}
	data, _ := json.Marshal(ev.Details)
	fmt.Fprintf(w, "event: %s\n", name)
	fmt.Fprintf(w, "data: %s\n\n", string(data))
}

// writeSSEMessage sends a custom SSE event (non-agent).
// Used for scenario_step_complete, errors, etc.
func writeSSEMessage(w http.ResponseWriter, event string, payload any) {
	data, _ := json.Marshal(payload)
	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", string(data))
}
