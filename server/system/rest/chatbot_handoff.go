package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

// ChatbotHandoffController exposes the operator-facing surface for DB-backed
// chatbot sessions. Mirrors the relevant subset of /api/widget/v1 (handoff
// accept, operator-message, handoff-complete, stream) but authenticates with
// the standard Corteza admin token instead of a per-session JWT.
//
// Mounted under the protected system rest group; the upstream HTTP token
// validator has already run by the time these handlers fire.
type ChatbotHandoffController struct {
	store      store.Storer
	bus        *observability.Bus
	sessionSvc chatbotSession
}

// chatbotSession narrows service.DefaultChatbotSession to the surface used here
// and keeps the controller testable.
type chatbotSession interface {
	FindByID(ctx context.Context, ID uint64) (*types.ChatbotSession, error)
	FindStepByID(ctx context.Context, id uint64) (*types.ChatbotSessionStep, error)
	FindStepBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionStep, error)
	FindHandoffByID(ctx context.Context, id uint64) (*types.ChatbotSessionHandoff, error)
	FindLiveHandoffBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionHandoff, error)
	ActivateHandoff(ctx context.Context, handoffID uint64) error
	CloseHandoff(ctx context.Context, handoffID uint64) error
	AppendOperatorMessage(ctx context.Context, handoffID uint64, message string) (uint64, error)
}

// NewChatbotHandoffController wires the controller with shared services. All
// dependencies come from system/service defaults.
func NewChatbotHandoffController() *ChatbotHandoffController {
	return &ChatbotHandoffController{
		store:      service.DefaultStore,
		bus:        service.DefaultObsBus,
		sessionSvc: service.DefaultChatbotSession,
	}
}

// MountRoutes attaches the controller to the supplied router. Caller must
// invoke this inside the protected (admin-authed) group so identity is set.
func (c *ChatbotHandoffController) MountRoutes(r chi.Router) {
	r.Route("/chatbots/sessions/{sessionID}", func(r chi.Router) {
		r.Post("/handoff-accept", c.acceptHandoff)
		r.Post("/operator-message", c.sendOperatorMessage)
		r.Post("/handoff-complete", c.closeHandoff)
		r.Get("/stream", c.stream)
	})
}

// acceptHandoff activates the session's pending handoff. Operator must have
// passed the RBAC gate (see rbacGate). Idempotent on already-active handoff.
func (c *ChatbotHandoffController) acceptHandoff(w http.ResponseWriter, r *http.Request) {
	sess := c.sessionFromURL(w, r)
	if sess == nil {
		return
	}
	cb, ok := c.rbacGate(w, r, sess)
	if !ok {
		return
	}

	h, err := c.sessionSvc.FindLiveHandoffBySession(r.Context(), sess.ID)
	if err != nil {
		http.Error(w, "handoff: lookup failed", http.StatusInternalServerError)
		return
	}
	if h == nil {
		http.Error(w, "handoff: no pending handoff", http.StatusNotFound)
		return
	}
	if h.Status == "active" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if h.Status != "requested" {
		http.Error(w, "handoff: not in requested state", http.StatusConflict)
		return
	}

	if err := c.sessionSvc.ActivateHandoff(r.Context(), h.ID); err != nil {
		http.Error(w, "handoff: cannot activate", http.StatusInternalServerError)
		return
	}

	convID, _ := c.handoffConvID(r.Context(), h)
	c.emit(convID, "handoff_active", map[string]any{
		"handoffID": strconv.FormatUint(h.ID, 10),
		"operator":  operatorLabel(r.Context(), cb),
	})

	w.WriteHeader(http.StatusNoContent)
}

// sendOperatorMessage relays an operator message to the user widget via SSE and
// appends it to the conversation (as an assistant turn) so the LLM sees it on
// resume. Requires the session's handoff to be in "active" state.
func (c *ChatbotHandoffController) sendOperatorMessage(w http.ResponseWriter, r *http.Request) {
	sess := c.sessionFromURL(w, r)
	if sess == nil {
		return
	}
	cb, ok := c.rbacGate(w, r, sess)
	if !ok {
		return
	}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		http.Error(w, "handoff: bad body", http.StatusBadRequest)
		return
	}

	h, err := c.sessionSvc.FindLiveHandoffBySession(r.Context(), sess.ID)
	if err != nil {
		http.Error(w, "handoff: lookup failed", http.StatusInternalServerError)
		return
	}
	if h == nil || h.Status != "active" {
		http.Error(w, "handoff: not accepted", http.StatusConflict)
		return
	}

	convID, err := c.sessionSvc.AppendOperatorMessage(r.Context(), h.ID, body.Message)
	if err != nil {
		http.Error(w, "handoff: send failed", http.StatusInternalServerError)
		return
	}

	c.emit(convID, "operator_message", map[string]any{
		"content":  body.Message,
		"operator": operatorLabel(r.Context(), cb),
	})

	w.WriteHeader(http.StatusNoContent)
}

// closeHandoff terminates the active or requested handoff and resumes the
// session at the same step.
func (c *ChatbotHandoffController) closeHandoff(w http.ResponseWriter, r *http.Request) {
	sess := c.sessionFromURL(w, r)
	if sess == nil {
		return
	}
	if _, ok := c.rbacGate(w, r, sess); !ok {
		return
	}

	h, err := c.sessionSvc.FindLiveHandoffBySession(r.Context(), sess.ID)
	if err != nil {
		http.Error(w, "handoff: lookup failed", http.StatusInternalServerError)
		return
	}
	if h == nil {
		http.Error(w, "handoff: no live handoff", http.StatusNotFound)
		return
	}

	if err := c.sessionSvc.CloseHandoff(r.Context(), h.ID); err != nil {
		http.Error(w, "handoff: cannot close", http.StatusInternalServerError)
		return
	}

	convID, _ := c.handoffConvID(r.Context(), h)
	c.emit(convID, "handoff_complete", map[string]any{
		"handoffID": strconv.FormatUint(h.ID, 10),
	})

	w.WriteHeader(http.StatusNoContent)
}

// stream opens an SSE connection scoped to the conversation backing the
// session's current step. The dispatcher's allowlist (see observability/sse.go)
// applies, so internal trace events stay filtered out.
func (c *ChatbotHandoffController) stream(w http.ResponseWriter, r *http.Request) {
	sess := c.sessionFromURL(w, r)
	if sess == nil {
		return
	}
	if _, ok := c.rbacGate(w, r, sess); !ok {
		return
	}
	if c.bus == nil {
		http.Error(w, "handoff: no bus", http.StatusServiceUnavailable)
		return
	}

	convID := c.sessionConvID(r.Context(), sess)
	if convID == 0 {
		http.Error(w, "handoff: no conversation", http.StatusNotFound)
		return
	}

	d := observability.NewSSEDispatcher(convID)
	c.bus.Register(d)
	defer d.Close()
	observability.PumpSSE(w, r, d)
}

// ---- helpers ----------------------------------------------------------------

// sessionFromURL parses the {sessionID} path param and loads the session.
// Writes 4xx and returns nil on failure.
func (c *ChatbotHandoffController) sessionFromURL(w http.ResponseWriter, r *http.Request) *types.ChatbotSession {
	idRaw := chi.URLParam(r, "sessionID")
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "handoff: bad session id", http.StatusBadRequest)
		return nil
	}
	sess, err := c.sessionSvc.FindByID(r.Context(), id)
	if err != nil || sess == nil {
		http.Error(w, "handoff: session not found", http.StatusNotFound)
		return nil
	}
	return sess
}

// rbacGate enforces that the caller is authenticated and (when Handoff.TargetRoles
// is non-empty on the parent chatbot) a member of at least one allowed role.
// Returns the resolved chatbot so callers can avoid a second lookup.
func (c *ChatbotHandoffController) rbacGate(
	w http.ResponseWriter,
	r *http.Request,
	sess *types.ChatbotSession,
) (*types.Chatbot, bool) {
	id := auth.GetIdentityFromContext(r.Context())
	if id == nil || !id.Valid() {
		http.Error(w, "handoff: unauthorized", http.StatusUnauthorized)
		return nil, false
	}

	cb, err := store.LookupChatbotByID(r.Context(), c.store, sess.ChatbotID)
	if err != nil || cb == nil {
		http.Error(w, "handoff: chatbot not found", http.StatusNotFound)
		return nil, false
	}

	allowed := cb.Handoff.TargetRoles
	if len(allowed) == 0 {
		return cb, true
	}

	roles := id.Roles()
	roleSet := make(map[uint64]struct{}, len(roles))
	for _, rid := range roles {
		roleSet[rid] = struct{}{}
	}
	for _, want := range allowed {
		if _, ok := roleSet[want]; ok {
			return cb, true
		}
	}

	http.Error(w, "handoff: not allowed", http.StatusForbidden)
	return nil, false
}

// sessionConvID returns the conversation ID for the session's current step, or
// 0 when no step is available.
func (c *ChatbotHandoffController) sessionConvID(ctx context.Context, sess *types.ChatbotSession) uint64 {
	step, err := c.sessionSvc.FindStepBySession(ctx, sess.ID)
	if err != nil || step == nil {
		return 0
	}
	return step.ConversationID
}

// handoffConvID resolves the conversation ID for a handoff's owning step.
func (c *ChatbotHandoffController) handoffConvID(ctx context.Context, h *types.ChatbotSessionHandoff) (uint64, error) {
	if h == nil {
		return 0, nil
	}
	step, err := c.sessionSvc.FindStepByID(ctx, h.StepID)
	if err != nil || step == nil {
		return 0, err
	}
	return step.ConversationID, nil
}

// emit broadcasts an SSE event for the given conversation. No-op when bus or
// convID is zero — emit failures must not break the HTTP response.
func (c *ChatbotHandoffController) emit(convID uint64, name string, details map[string]any) {
	if c.bus == nil || convID == 0 {
		return
	}
	c.bus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          name,
		Details:        details,
		Timestamp:      time.Now(),
	})
}

// operatorLabel resolves the operator display label from the request identity.
// Falls back to "operator" when the user lookup is unavailable so SSE consumers
// always receive a non-empty operator field.
func operatorLabel(ctx context.Context, _ *types.Chatbot) string {
	id := auth.GetIdentityFromContext(ctx)
	if id == nil {
		return "operator"
	}
	u, err := store.LookupUserByID(ctx, service.DefaultStore, id.Identity())
	if err != nil || u == nil {
		return "operator"
	}
	if u.Name != "" {
		return u.Name
	}
	if u.Handle != "" {
		return u.Handle
	}
	if u.Email != "" {
		return u.Email
	}
	return "operator"
}
