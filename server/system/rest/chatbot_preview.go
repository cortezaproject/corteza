package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

// ChatbotPreviewController mounts the admin-authed preview API that mirrors
// the public /api/widget/v1 endpoints but holds all session state in memory.
// The chatbot config is supplied inline at session-open so editor drafts can
// be exercised without persisting to the DB.
//
// Routes are scoped under /chatbot/preview/* in the protected system API group.
type ChatbotPreviewController struct {
	preview previewService
	bus     *observability.Bus
}

// previewService narrows service.DefaultChatbotPreview to the surface used
// by REST. Lifecycle/state methods stay on the service; this is just the
// orchestration surface invoked from handlers.
type previewService interface {
	OpenWithConversation(ctx context.Context, cb types.Chatbot) (*service.ChatbotPreviewSession, *types.AiConversation, error)
	Find(id string) *service.ChatbotPreviewSession
	Start(ps *service.ChatbotPreviewSession)
	SubmitMessage(ctx context.Context, ps *service.ChatbotPreviewSession, input string) error
	SubmitForm(ctx context.Context, ps *service.ChatbotPreviewSession, fields map[string]string) (map[string]string, error)
	SubmitConsent(ps *service.ChatbotPreviewSession, accepted bool) error
	AdvanceStep(ps *service.ChatbotPreviewSession) error
	CloseSession(ps *service.ChatbotPreviewSession) error
	RequestHandoff(ps *service.ChatbotPreviewSession) *service.ChatbotPreviewHandoff
}

// NewChatbotPreviewController wires the controller with shared services.
func NewChatbotPreviewController() *ChatbotPreviewController {
	return &ChatbotPreviewController{
		preview: service.DefaultChatbotPreview,
		bus:     service.DefaultObsBus,
	}
}

// MountRoutes attaches the preview routes to the supplied router. The caller
// is expected to invoke this inside the protected (admin-authed) group of
// system/rest.router so that auth.HttpTokenValidator has already run.
func (c *ChatbotPreviewController) MountRoutes(r chi.Router) {
	r.Route("/chatbot/preview", func(r chi.Router) {
		r.Post("/session", c.openSession)
		r.Post("/session/{id}/start", c.startSession)
		r.Post("/session/{id}/submit", c.submit)
		r.Post("/session/{id}/advance-step", c.advanceStep)
		r.Post("/session/{id}/close", c.closeSession)
		r.Post("/session/{id}/handoff", c.requestHandoff)
		r.Get("/session/{id}/stream", c.stream)
	})
}

func (c *ChatbotPreviewController) openSession(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	var body struct {
		Chatbot types.Chatbot `json:"chatbot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "preview: bad body", http.StatusBadRequest)
		return
	}

	ps, conv, err := c.preview.OpenWithConversation(r.Context(), body.Chatbot)
	if err != nil {
		c.writePreviewError(w, err)
		return
	}

	previewJSON(w, http.StatusOK, map[string]any{
		"sessionID":      ps.ID,
		"token":          "",
		"conversationID": strconv.FormatUint(conv.ID, 10),
		"dbSessionID":    ps.ID,
	})
}

func (c *ChatbotPreviewController) startSession(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	c.preview.Start(ps)
	w.WriteHeader(http.StatusNoContent)
}

func (c *ChatbotPreviewController) submit(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}

	var body struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "preview: bad body", http.StatusBadRequest)
		return
	}

	switch body.Type {
	case "message":
		var d struct {
			Input string `json:"input"`
		}
		if err := json.Unmarshal(body.Data, &d); err != nil {
			http.Error(w, "preview: bad message payload", http.StatusBadRequest)
			return
		}
		if err := c.preview.SubmitMessage(r.Context(), ps, d.Input); err != nil {
			c.writePreviewError(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	case "form":
		var d struct {
			Fields map[string]string `json:"fields"`
		}
		if err := json.Unmarshal(body.Data, &d); err != nil {
			http.Error(w, "preview: bad form payload", http.StatusBadRequest)
			return
		}
		errs, err := c.preview.SubmitForm(r.Context(), ps, d.Fields)
		if err != nil {
			c.writePreviewError(w, err)
			return
		}
		if len(errs) > 0 {
			previewJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "consent":
		var d struct {
			Accepted bool `json:"accepted"`
		}
		if err := json.Unmarshal(body.Data, &d); err != nil {
			http.Error(w, "preview: bad consent payload", http.StatusBadRequest)
			return
		}
		if err := c.preview.SubmitConsent(ps, d.Accepted); err != nil {
			c.writePreviewError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "preview: unknown submit type", http.StatusBadRequest)
	}
}

func (c *ChatbotPreviewController) advanceStep(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	if err := c.preview.AdvanceStep(ps); err != nil {
		c.writePreviewError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *ChatbotPreviewController) closeSession(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	_ = c.preview.CloseSession(ps)
	w.WriteHeader(http.StatusNoContent)
}

func (c *ChatbotPreviewController) requestHandoff(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	h := c.preview.RequestHandoff(ps)
	previewJSON(w, http.StatusOK, map[string]any{
		"status":    "handoff_requested",
		"handoffID": h.ID,
	})
}

func (c *ChatbotPreviewController) stream(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	if c.bus == nil {
		http.Error(w, "preview: no bus", http.StatusServiceUnavailable)
		return
	}
	d := observability.NewSSEDispatcher(ps.ConversationID)
	c.bus.Register(d)
	defer d.Close()
	observability.PumpSSE(w, r, d)
}

// ---- helpers ----------------------------------------------------------------

func (c *ChatbotPreviewController) sessionFromURL(w http.ResponseWriter, r *http.Request) *service.ChatbotPreviewSession {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "preview: missing session id", http.StatusBadRequest)
		return nil
	}
	ps := c.preview.Find(id)
	if ps == nil {
		http.Error(w, "preview: session not found", http.StatusNotFound)
		return nil
	}
	return ps
}

func (c *ChatbotPreviewController) adminGate(w http.ResponseWriter, r *http.Request) bool {
	id := auth.GetIdentityFromContext(r.Context())
	if id == nil || !id.Valid() {
		http.Error(w, "preview: unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// writePreviewError maps service-layer sentinel errors to HTTP status codes.
func (c *ChatbotPreviewController) writePreviewError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ChatbotSessionErrChatbotNoScenarios()),
		errors.Is(err, service.ChatbotSessionErrScenarioOutOfRange()),
		errors.Is(err, service.ChatbotSessionErrBadHandoffID()),
		errors.Is(err, service.ChatbotSessionErrScenarioNotMessage()):
		http.Error(w, "preview: "+err.Error(), http.StatusBadRequest)
	case errors.Is(err, service.ChatbotSessionErrSessionNotActive()),
		errors.Is(err, service.ChatbotSessionErrScenarioNotForm()),
		errors.Is(err, service.ChatbotSessionErrScenarioNotConsent()),
		errors.Is(err, service.ChatbotSessionErrHandoffNotRequested()),
		errors.Is(err, service.ChatbotSessionErrHandoffNotActive()):
		http.Error(w, "preview: "+err.Error(), http.StatusConflict)
	case errors.Is(err, service.ChatbotSessionErrNoActiveStep()),
		errors.Is(err, service.ChatbotSessionErrHandoffNotFound()):
		http.Error(w, "preview: "+err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ChatbotSessionErrAgentUnavailable()),
		errors.Is(err, service.ChatbotSessionErrAgentNotConfigured()):
		http.Error(w, "preview: "+err.Error(), http.StatusServiceUnavailable)
	default:
		http.Error(w, "preview: "+err.Error(), http.StatusInternalServerError)
	}
}

func previewJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

