package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/agentic/runtime"
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
	store    store.Storer
	preview  previewStore
	bus      *observability.Bus
	runtime  previewRuntime
	convCtor previewConvCreator
}

// previewStore narrows service.DefaultChatbotPreview to the surface used here.
type previewStore interface {
	Open(cb types.Chatbot, convID uint64) *service.ChatbotPreviewSession
	Find(id string) *service.ChatbotPreviewSession
	StartStep(ps *service.ChatbotPreviewSession, scenarioIndex int)
	FinalizeStep(ps *service.ChatbotPreviewSession)
	FailStep(ps *service.ChatbotPreviewSession, reason string)
	CurrentStep(ps *service.ChatbotPreviewSession) *service.ChatbotPreviewStep
	RequestHandoff(ps *service.ChatbotPreviewSession) *service.ChatbotPreviewHandoff
	ActivateHandoff(ps *service.ChatbotPreviewSession, operator string)
	CloseHandoff(ps *service.ChatbotPreviewSession)
	Close(ps *service.ChatbotPreviewSession, reason string)
	EmitUserMessage(ps *service.ChatbotPreviewSession, content string)
	EmitOperatorMessage(ps *service.ChatbotPreviewSession, content, operator string)
	EmitFormError(ps *service.ChatbotPreviewSession, scenarioID string, errs map[string]string)
}

type previewRuntime interface {
	Run(ctx context.Context, req *runtime.AgentRequest) (*runtime.AgentResponse, error)
}

type previewConvCreator interface {
	Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
}

// NewChatbotPreviewController wires the controller with shared services. All
// dependencies come from system/service defaults.
func NewChatbotPreviewController() *ChatbotPreviewController {
	return &ChatbotPreviewController{
		store:    service.DefaultStore,
		preview:  service.DefaultChatbotPreview,
		bus:      service.DefaultObsBus,
		runtime:  service.DefaultAgenticRuntime,
		convCtor: service.DefaultAiConversation,
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
		r.Post("/session/{id}/handoff-accept", c.acceptHandoff)
		r.Post("/session/{id}/operator-message", c.sendOperatorMessage)
		r.Post("/session/{id}/handoff-complete", c.closeHandoff)
		r.Get("/session/{id}/stream", c.stream)
	})
}

// openSession allocates an in-memory preview session bound to a fresh
// AiConversation (real, so the agent runtime can persist history). The
// inline chatbot snapshot is the source of truth for the rest of the flow.
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
	if len(body.Chatbot.Scenarios) == 0 {
		http.Error(w, "preview: chatbot has no scenarios", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	conv, err := c.convCtor.Create(svcCtx, &types.AiConversation{})
	if err != nil {
		http.Error(w, "preview: cannot open conversation", http.StatusInternalServerError)
		return
	}

	ps := c.preview.Open(body.Chatbot, conv.ID)

	previewJSON(w, http.StatusOK, map[string]any{
		"sessionID":      ps.ID,
		"token":          "",
		"conversationID": strconv.FormatUint(conv.ID, 10),
		"dbSessionID":    ps.ID,
	})
}

// startSession triggers the first step_start emission. Called by the client
// after the SSE stream has been opened so the event is not dropped by the
// bus (which has no per-subscriber buffering).
func (c *ChatbotPreviewController) startSession(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	// Only fire if no step has been started yet (idempotent on retry).
	if c.preview.CurrentStep(ps) == nil {
		c.preview.StartStep(ps, 0)
	}
	w.WriteHeader(http.StatusNoContent)
}

// submit dispatches on payload.type the same way the widget does.
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
		c.submitMessage(w, r, ps, body.Data)
	case "form":
		c.submitForm(w, r, ps, body.Data)
	default:
		http.Error(w, "preview: unknown submit type", http.StatusBadRequest)
	}
}

func (c *ChatbotPreviewController) submitMessage(w http.ResponseWriter, r *http.Request, ps *service.ChatbotPreviewSession, raw json.RawMessage) {
	var data struct {
		Input string `json:"input"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		http.Error(w, "preview: bad message payload", http.StatusBadRequest)
		return
	}

	switch ps.Status {
	case "handoff_requested", "handoff_active":
		c.appendUserMessage(ps, data.Input)
		c.preview.EmitUserMessage(ps, data.Input)
		w.WriteHeader(http.StatusAccepted)
		return
	case "active":
		// fall through
	default:
		http.Error(w, "preview: session not active", http.StatusConflict)
		return
	}

	step := c.preview.CurrentStep(ps)
	if step == nil {
		http.Error(w, "preview: no active step", http.StatusNotFound)
		return
	}
	scenario := c.scenarioAt(ps, step.ScenarioIndex)
	if scenario == nil || scenario.Type != "conversation" || scenario.AgentID == 0 {
		http.Error(w, "preview: current step does not accept messages", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	agent, err := store.LookupAgentByID(svcCtx, c.store, scenario.AgentID)
	if err != nil || agent == nil {
		http.Error(w, "preview: agent unavailable", http.StatusServiceUnavailable)
		return
	}
	if !agent.Invocation.System.Enabled || agent.Invocation.System.ServiceAccount == 0 {
		http.Error(w, "preview: agent not configured for system invocation", http.StatusServiceUnavailable)
		return
	}

	if conv, err := store.LookupAiConversationByID(svcCtx, c.store, ps.ConversationID); err == nil && conv != nil && conv.AgentID == 0 {
		conv.AgentID = agent.ID
		_ = store.UpdateAiConversation(svcCtx, c.store, conv)
	}

	saCtx := c.impersonate(context.Background(), agent.Invocation.System.ServiceAccount)
	cid := ps.ConversationID
	go func() {
		resp, err := c.runtime.Run(saCtx, &runtime.AgentRequest{
			AgentID:        agent.ID,
			Input:          data.Input,
			ConversationID: cid,
		})
		if err != nil {
			c.emitAgentError(cid, err.Error())
		} else if resp != nil && resp.Output != "" {
			c.emitToken(cid, resp.Output)
		}
		c.emitDone(cid)
	}()

	w.WriteHeader(http.StatusAccepted)
}

func (c *ChatbotPreviewController) submitForm(w http.ResponseWriter, r *http.Request, ps *service.ChatbotPreviewSession, raw json.RawMessage) {
	if ps.Status != "active" {
		http.Error(w, "preview: session not active", http.StatusConflict)
		return
	}
	var data struct {
		Fields map[string]string `json:"fields"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		http.Error(w, "preview: bad form payload", http.StatusBadRequest)
		return
	}

	step := c.preview.CurrentStep(ps)
	if step == nil {
		http.Error(w, "preview: no active step", http.StatusNotFound)
		return
	}
	scenario := c.scenarioAt(ps, step.ScenarioIndex)
	if scenario == nil || scenario.Type != "form" {
		http.Error(w, "preview: current step is not a form", http.StatusConflict)
		return
	}

	errs := service.ValidateChatbotFormFields(scenario.Config, data.Fields)
	if len(errs) > 0 {
		c.preview.EmitFormError(ps, scenario.ID, errs)
		previewJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	c.appendConvMessage(svcCtx, ps.ConversationID, types.AiConversationMessage{
		Role:    "user",
		Content: service.EncodeChatbotFormSubmission(data.Fields),
	})

	c.preview.FinalizeStep(ps)
	c.preview.StartStep(ps, step.ScenarioIndex+1)

	w.WriteHeader(http.StatusNoContent)
}

func (c *ChatbotPreviewController) advanceStep(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}

	if ps.Handoff != nil && ps.Handoff.Status != "closed" {
		c.preview.CloseHandoff(ps)
	}

	step := c.preview.CurrentStep(ps)
	if step == nil {
		http.Error(w, "preview: no active step", http.StatusNotFound)
		return
	}
	c.preview.FinalizeStep(ps)
	c.preview.StartStep(ps, step.ScenarioIndex+1)

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
	if step := c.preview.CurrentStep(ps); step != nil {
		c.preview.FinalizeStep(ps)
	}
	c.preview.Close(ps, "")
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

func (c *ChatbotPreviewController) acceptHandoff(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}

	var body struct {
		Operator string `json:"operator,omitempty"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	if ps.Handoff == nil {
		http.Error(w, "preview: no handoff", http.StatusNotFound)
		return
	}
	if ps.Handoff.Status == "active" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if ps.Handoff.Status != "requested" {
		http.Error(w, "preview: handoff not in requested state", http.StatusConflict)
		return
	}

	c.preview.ActivateHandoff(ps, body.Operator)
	w.WriteHeader(http.StatusNoContent)
}

func (c *ChatbotPreviewController) sendOperatorMessage(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}

	var body struct {
		Message  string `json:"message"`
		Operator string `json:"operator,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "preview: bad body", http.StatusBadRequest)
		return
	}

	if ps.Handoff == nil || ps.Handoff.Status != "active" {
		http.Error(w, "preview: handoff not accepted", http.StatusConflict)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	c.appendConvMessage(svcCtx, ps.ConversationID, types.AiConversationMessage{
		Role:    "assistant",
		Content: "[operator] " + body.Message,
	})

	c.preview.EmitOperatorMessage(ps, body.Message, body.Operator)
	w.WriteHeader(http.StatusNoContent)
}

func (c *ChatbotPreviewController) closeHandoff(w http.ResponseWriter, r *http.Request) {
	if !c.adminGate(w, r) {
		return
	}
	ps := c.sessionFromURL(w, r)
	if ps == nil {
		return
	}
	c.preview.CloseHandoff(ps)
	w.WriteHeader(http.StatusNoContent)
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

func (c *ChatbotPreviewController) scenarioAt(ps *service.ChatbotPreviewSession, idx int) *types.ChatbotScenario {
	if idx < 0 || idx >= len(ps.Chatbot.Scenarios) {
		return nil
	}
	return &ps.Chatbot.Scenarios[idx]
}

// adminGate enforces that the caller has a valid Corteza identity. The
// system rest private group already runs the HTTP token validator; this is a
// defensive secondary check so anyone reusing the controller from elsewhere
// can't bypass auth by accident.
func (c *ChatbotPreviewController) adminGate(w http.ResponseWriter, r *http.Request) bool {
	id := auth.GetIdentityFromContext(r.Context())
	if id == nil || !id.Valid() {
		http.Error(w, "preview: unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

// impersonate builds a context authenticated as the given user with role
// memberships loaded, mirroring widget controller behavior for agent runtime
// invocation.
func (c *ChatbotPreviewController) impersonate(ctx context.Context, userID uint64) context.Context {
	svcCtx := auth.SetIdentityToContext(ctx, auth.ServiceUser())
	mm, _, err := store.SearchRoleMembers(svcCtx, c.store, types.RoleMemberFilter{
		Resource: fmt.Sprintf("corteza::system:user/%d", userID),
	})
	var roles []uint64
	if err == nil {
		for _, m := range mm {
			roles = append(roles, m.RoleID)
		}
	}
	return auth.SetIdentityToContext(ctx, auth.Authenticated(userID, roles...))
}

func (c *ChatbotPreviewController) appendConvMessage(ctx context.Context, convID uint64, msg types.AiConversationMessage) {
	conv, err := store.LookupAiConversationByID(ctx, c.store, convID)
	if err != nil || conv == nil {
		return
	}
	conv.Messages = append(conv.Messages, msg)
	_ = store.UpdateAiConversation(ctx, c.store, conv)
}

func (c *ChatbotPreviewController) appendUserMessage(ps *service.ChatbotPreviewSession, content string) {
	svcCtx := auth.SetIdentityToContext(context.Background(), auth.ServiceUser())
	c.appendConvMessage(svcCtx, ps.ConversationID, types.AiConversationMessage{
		Role:    "user",
		Content: content,
	})
}

func (c *ChatbotPreviewController) emitToken(convID uint64, text string) {
	if c.bus == nil {
		return
	}
	c.bus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          "token",
		Details:        map[string]any{"text": text},
	})
}

func (c *ChatbotPreviewController) emitAgentError(convID uint64, msg string) {
	if c.bus == nil {
		return
	}
	c.bus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          "agent_error",
		Details:        map[string]any{"error": msg},
	})
}

func (c *ChatbotPreviewController) emitDone(convID uint64) {
	if c.bus == nil {
		return
	}
	c.bus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          "done",
		Details:        map[string]any{},
	})
}

func previewJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
