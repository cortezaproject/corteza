package observability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// SSEDispatcher is a transient Dispatcher scoped to one SSE stream. It is
// filtered to a single ConversationID so cross-tenant leaks are impossible,
// and to a fixed allowlist of event names so internal trace events (which
// carry raw user input + prompt internals) never reach end-user clients.
//
// Lives in observability so both the public widget and the admin preview
// can register dispatchers against the same bus instance.
type SSEDispatcher struct {
	id        string
	convID    string
	events    chan AgentEvent
	closed    chan struct{}
	closeOnce sync.Once
}

// clientEvents enumerates the only event names a public SSE client should
// receive. Anything emitted by the agentic runtime for tracing (agent.invoked,
// agent.decision, prompt.build, tool.*, etc.) is intentionally excluded; those
// belong on Langfuse / OTel dispatchers, not the widget wire.
var clientEvents = map[string]struct{}{
	"token":              {},
	"done":               {},
	"agent_error":        {},
	"step_start":         {},
	"step_complete":      {},
	"step_error":         {},
	"step_hook_skipped":  {},
	"form_error":         {},
	"handoff_requested":  {},
	"handoff_active":     {},
	"handoff_complete":   {},
	"operator_message":   {},
	"user_message":       {},
	"session_closed":     {},
}

// NewSSEDispatcher allocates a buffered SSE dispatcher for the given conversation.
func NewSSEDispatcher(convID uint64) *SSEDispatcher {
	return &SSEDispatcher{
		id:     fmt.Sprintf("sse-%d-%d", convID, time.Now().UnixNano()),
		convID: strconv.FormatUint(convID, 10),
		events: make(chan AgentEvent, 64),
		closed: make(chan struct{}),
	}
}

func (d *SSEDispatcher) ID() string                  { return d.id }
func (d *SSEDispatcher) OnSpan(_ AgentSpan) error    { return nil }
func (d *SSEDispatcher) Flush() error                { return nil }
func (d *SSEDispatcher) Shutdown() error             { d.Close(); return nil }

func (d *SSEDispatcher) OnEvent(e AgentEvent) error {
	if e.ConversationID != d.convID {
		return nil
	}
	if _, ok := clientEvents[e.Event]; !ok {
		return nil
	}
	select {
	case <-d.closed:
		return nil
	case d.events <- e:
	default:
		// drop on backpressure — UX can resync from conversation store
	}
	return nil
}

// Close releases the dispatcher and unblocks any pumpers. Idempotent.
func (d *SSEDispatcher) Close() {
	d.closeOnce.Do(func() { close(d.closed) })
}

// PumpSSE writes SSE frames until the client disconnects or the context
// deadline hits. Keepalive comment frames are sent every 15s so proxies don't
// sever the connection.
func PumpSSE(w http.ResponseWriter, r *http.Request, d *SSEDispatcher) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "stream: unsupported", http.StatusInternalServerError)
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

func writeSSEEvent(w http.ResponseWriter, ev AgentEvent) {
	name := ev.Event
	if name == "" {
		name = "message"
	}
	data, _ := json.Marshal(ev.Details)
	fmt.Fprintf(w, "event: %s\n", name)
	fmt.Fprintf(w, "data: %s\n\n", string(data))
}
