package widget

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

// Controller wires together the public chatbot widget HTTP API.
// All endpoints live under /api/widget/v1 and are unauthenticated except
// by the per-session JWT issued from POST /session.
type Controller struct {
	chatbotLookup *ChatbotByKeyLookup
	obsBus        *observability.Bus
	runtime       AgenticRunner
	convStore     AiConversationStore
	store         store.Storer
	serverSecret  string
	sessionSvc    ChatbotSessionService
}

// AgenticRunner is the minimal subset of the runtime the widget needs.
type AgenticRunner interface {
	Run(ctx context.Context, req *runtime.AgentRequest) (*runtime.AgentResponse, error)
}

// AiConversationStore is satisfied by service.DefaultAiConversation.
type AiConversationStore interface {
	Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
}

// ChatbotSessionService is satisfied by service.ChatbotSession().
type ChatbotSessionService interface {
	Create(ctx context.Context, new *types.ChatbotSession) (*types.ChatbotSession, error)
	FindByID(ctx context.Context, ID uint64) (*types.ChatbotSession, error)
	UpdateStatus(ctx context.Context, id uint64, status string) error
	UpdateCurrentStep(ctx context.Context, id uint64, idx int) error
	CreateStep(ctx context.Context, step *types.ChatbotSessionStep) (*types.ChatbotSessionStep, error)
	FindStepByID(ctx context.Context, id uint64) (*types.ChatbotSessionStep, error)
	FindStepBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionStep, error)
	FindStepBySessionAndIndex(ctx context.Context, sessionID uint64, scenarioIndex int) (*types.ChatbotSessionStep, error)
	CompleteStep(ctx context.Context, stepID uint64) error
	StartStep(ctx context.Context, sessionID uint64, scenario *types.ChatbotScenario, conversationID uint64, scenarioIndex int, vars *expr.Vars) (*types.ChatbotSessionStep, error)
	FinalizeStep(ctx context.Context, stepID uint64, scenario *types.ChatbotScenario, vars *expr.Vars) error
	RequestHandoff(ctx context.Context, sessionID, stepID uint64) (*types.ChatbotSessionHandoff, error)
	ActivateHandoff(ctx context.Context, handoffID uint64) error
	CloseHandoff(ctx context.Context, handoffID uint64) error
	FindHandoffByID(ctx context.Context, id uint64) (*types.ChatbotSessionHandoff, error)
	FindHandoffBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionHandoff, error)
}

func New(
	s store.Storer,
	obsBus *observability.Bus,
	rt AgenticRunner,
	conv AiConversationStore,
	serverSecret string,
	sessionSvc ChatbotSessionService,
) *Controller {
	return &Controller{
		chatbotLookup: NewChatbotByKeyLookup(s),
		obsBus:        obsBus,
		runtime:       rt,
		convStore:     conv,
		store:         s,
		serverSecret:  serverSecret,
		sessionSvc:    sessionSvc,
	}
}

// resolveAgent loads an Agent by ID under service identity.
func (c *Controller) resolveAgent(ctx context.Context, agentID uint64) (*types.Agent, error) {
	if agentID == 0 {
		return nil, fmt.Errorf("widget: chatbot has no agent")
	}
	svcCtx := auth.SetIdentityToContext(ctx, auth.ServiceUser())
	return store.LookupAgentByID(svcCtx, c.store, agentID)
}

// impersonateServiceAccount builds a context authenticated as the given user
// with their full role memberships loaded from the store. Widget runs without
// caller identity, so RBAC checks downstream rely entirely on this.
func (c *Controller) impersonateServiceAccount(ctx context.Context, userID uint64) context.Context {
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

// MountRoutes mounts the widget routes onto r at the /widget/v1 prefix.
func (c *Controller) MountRoutes(r chi.Router) {
	r.Route("/widget/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(c.keyOnly)
			r.Options("/asset/{attachmentID}", okOptions)
			r.Get("/asset/{attachmentID}", c.readAsset)
			r.Options("/asset/{attachmentID}/{name}", okOptions)
			r.Get("/asset/{attachmentID}/{name}", c.readAsset)
		})

		r.Group(func(r chi.Router) {
			r.Use(c.corsForKey)

			r.Options("/config", okOptions)
			r.Get("/config", c.readConfig)

			r.Options("/session", okOptions)
			r.Post("/session", c.createSession)

			r.Group(func(r chi.Router) {
				r.Use(c.requireSessionJWT)

				r.Options("/session/{id}/submit", okOptions)
				r.Post("/session/{id}/submit", c.submit)

				r.Options("/session/{id}/advance-step", okOptions)
				r.Post("/session/{id}/advance-step", c.advanceStep)

				r.Options("/session/{id}/close", okOptions)
				r.Post("/session/{id}/close", c.closeSession)

				r.Options("/session/{id}/handoff", okOptions)
				r.Post("/session/{id}/handoff", c.requestHandoff)

				r.Options("/session/{id}/operator-message", okOptions)
				r.Post("/session/{id}/operator-message", c.sendOperatorMessage)

				r.Options("/session/{id}/handoff-complete", okOptions)
				r.Post("/session/{id}/handoff-complete", c.closeHandoff)

				r.Get("/session/{id}/stream", c.stream)
			})
		})
	})
}

// ---- handlers ---------------------------------------------------------------

func (c *Controller) readConfig(w http.ResponseWriter, r *http.Request) {
	cb := chatbotFromCtx(r.Context())
	if cb == nil {
		http.Error(w, "widget: no chatbot", http.StatusForbidden)
		return
	}
	writeJSON(w, http.StatusOK, publicConfigChatbot(cb))
}

func (c *Controller) readAsset(w http.ResponseWriter, r *http.Request) {
	cb := chatbotFromCtx(r.Context())
	if cb == nil {
		http.Error(w, "widget: no chatbot", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "attachmentID"), 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "widget: bad asset id", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	att, err := store.LookupAttachmentByID(svcCtx, c.store, id)
	if err != nil || att == nil || att.DeletedAt != nil {
		http.Error(w, "widget: asset not found", http.StatusNotFound)
		return
	}

	if att.Kind != types.AttachmentKindChatbot {
		http.Error(w, "widget: asset not found", http.StatusNotFound)
		return
	}
	if att.Meta.Labels["chatbotID"] != strconv.FormatUint(cb.ID, 10) {
		http.Error(w, "widget: asset not found", http.StatusNotFound)
		return
	}

	fh, err := service.DefaultAttachment.OpenOriginal(att)
	if err != nil || fh == nil {
		http.Error(w, "widget: asset unavailable", http.StatusNotFound)
		return
	}
	defer fh.Close()

	w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
	if mt := att.Meta.Original.Mimetype; mt != "" {
		w.Header().Set("Content-Type", mt)
	}
	http.ServeContent(w, r, att.Name, att.CreatedAt, fh)
}

// createSession opens a new ChatbotSession + AiConversation and starts step 0
// (runs its before-automation, emits step_start). Returns JWT.
func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	cb := chatbotFromCtx(r.Context())
	if cb == nil {
		http.Error(w, "widget: no chatbot", http.StatusForbidden)
		return
	}
	if len(cb.Scenarios) == 0 {
		http.Error(w, "widget: chatbot has no scenarios", http.StatusServiceUnavailable)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())

	conv, err := c.convStore.Create(svcCtx, &types.AiConversation{})
	if err != nil {
		http.Error(w, "widget: cannot open session", http.StatusInternalServerError)
		return
	}

	session, err := c.sessionSvc.Create(svcCtx, &types.ChatbotSession{
		ChatbotID: cb.ID,
		Status:    "active",
	})
	if err != nil {
		http.Error(w, "widget: cannot create session", http.StatusInternalServerError)
		return
	}

	ttl := parseTTL(cb.SessionTTL, 2*time.Hour)
	nowT := time.Now()
	claims := sessionClaims{
		Sid:   randomID(),
		Cid:   conv.ID,
		Dbsid: session.ID,
		Iat:   nowT.Unix(),
		Exp:   nowT.Add(ttl).Unix(),
	}
	tok, err := signSession(claims, deriveSessionSecret(cb.WidgetKey, c.serverSecret))
	if err != nil {
		http.Error(w, "widget: cannot sign token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"sessionID":      claims.Sid,
		"token":          tok,
		"conversationID": strconv.FormatUint(conv.ID, 10),
		"dbSessionID":    strconv.FormatUint(session.ID, 10),
	})

	// Kick off step 0 asynchronously: invoke before-automation and emit
	// step_start so the SSE client receives the first scenario payload as soon
	// as it connects. Errors are surfaced via SSE.
	go c.startStep(context.Background(), cb, session.ID, conv.ID, 0, nil)
}

// submit is the single user-input endpoint. The payload is shape-typed:
//
//	{ "type": "message", "data": { "input": "..." } }   // conversation / handoff
//	{ "type": "form",    "data": { "fields": {..} } }   // form scenario
//
// The server reads the current step to resolve the scenario; the client does
// not pass a scenarioID. Type-vs-step mismatches return 400.
func (c *Controller) submit(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	cb := chatbotFromCtx(r.Context())
	if claims == nil || cb == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	var body struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := readJSON(r.Body, &body); err != nil {
		http.Error(w, "widget: bad body", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())

	session, err := c.sessionSvc.FindByID(svcCtx, claims.Dbsid)
	if err != nil || session == nil {
		http.Error(w, "widget: session not found", http.StatusUnauthorized)
		return
	}

	switch body.Type {
	case "message":
		c.submitMessage(w, r, svcCtx, claims, cb, session, body.Data)
	case "form":
		c.submitForm(w, r, svcCtx, claims, cb, session, body.Data)
	default:
		http.Error(w, "widget: unknown submit type", http.StatusBadRequest)
	}
}

// submitMessage handles { type: "message", data: { input } }. Active session →
// agent runtime; handoff_active/requested → bypass agent and relay to operator.
func (c *Controller) submitMessage(
	w http.ResponseWriter,
	r *http.Request,
	svcCtx context.Context,
	claims *sessionClaims,
	cb *types.Chatbot,
	session *types.ChatbotSession,
	raw json.RawMessage,
) {
	var data struct {
		Input string `json:"input"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		http.Error(w, "widget: bad message payload", http.StatusBadRequest)
		return
	}

	switch session.Status {
	case "handoff_requested", "handoff_active":
		c.appendConversationMessage(svcCtx, claims.Cid, types.AiConversationMessage{
			Role:    "user",
			Content: data.Input,
		})
		c.emit(claims.Cid, "user_message", map[string]any{"content": data.Input})
		w.WriteHeader(http.StatusAccepted)
		return
	case "active":
		// fall through
	default:
		http.Error(w, "widget: session not active", http.StatusConflict)
		return
	}

	step, err := c.sessionSvc.FindStepBySession(svcCtx, claims.Dbsid)
	if err != nil || step == nil {
		http.Error(w, "widget: step not found", http.StatusNotFound)
		return
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		http.Error(w, "widget: scenario out of range", http.StatusBadRequest)
		return
	}
	scenario := &cb.Scenarios[step.ScenarioIndex]
	if scenario.Type != "conversation" || scenario.AgentID == 0 {
		http.Error(w, "widget: current step does not accept messages", http.StatusBadRequest)
		return
	}

	agent, err := c.resolveAgent(r.Context(), scenario.AgentID)
	if err != nil || agent == nil {
		http.Error(w, "widget: agent unavailable", http.StatusServiceUnavailable)
		return
	}
	if !agent.Invocation.System.Enabled || agent.Invocation.System.ServiceAccount == 0 {
		http.Error(w, "widget: agent not configured for widget runtime", http.StatusServiceUnavailable)
		return
	}

	if conv, err := store.LookupAiConversationByID(svcCtx, c.store, claims.Cid); err == nil && conv != nil && conv.AgentID == 0 {
		conv.AgentID = agent.ID
		_ = store.UpdateAiConversation(svcCtx, c.store, conv)
	}

	saCtx := c.impersonateServiceAccount(context.Background(), agent.Invocation.System.ServiceAccount)
	cid := claims.Cid
	go func() {
		resp, err := c.runtime.Run(saCtx, &runtime.AgentRequest{
			AgentID:        agent.ID,
			Input:          data.Input,
			ConversationID: cid,
		})
		// Surface runtime errors to the user as a system bubble.
		if err != nil {
			c.emit(cid, "agent_error", map[string]any{"error": err.Error()})
		} else if resp != nil && resp.Output != "" {
			// Non-streaming providers don't emit `token` events. Fan out
			// the final assistant text so clients have a single shape.
			c.emit(cid, "token", map[string]any{"text": resp.Output})
		}
		c.emit(cid, "done", map[string]any{})
	}()

	w.WriteHeader(http.StatusAccepted)
}

// submitForm handles { type: "form", data: { fields } }. Validates against the
// current scenario's config; on success finalizes the step and starts the next.
func (c *Controller) submitForm(
	w http.ResponseWriter,
	_ *http.Request,
	svcCtx context.Context,
	claims *sessionClaims,
	cb *types.Chatbot,
	session *types.ChatbotSession,
	raw json.RawMessage,
) {
	if session.Status != "active" {
		http.Error(w, "widget: session not active", http.StatusConflict)
		return
	}

	var data struct {
		Fields map[string]string `json:"fields"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		http.Error(w, "widget: bad form payload", http.StatusBadRequest)
		return
	}

	step, err := c.sessionSvc.FindStepBySession(svcCtx, claims.Dbsid)
	if err != nil || step == nil {
		http.Error(w, "widget: step not found", http.StatusNotFound)
		return
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		http.Error(w, "widget: scenario out of range", http.StatusBadRequest)
		return
	}
	scenario := &cb.Scenarios[step.ScenarioIndex]
	if scenario.Type != "form" {
		http.Error(w, "widget: current step is not a form", http.StatusConflict)
		return
	}

	errs := service.ValidateChatbotFormFields(scenario.Config, data.Fields)
	if len(errs) > 0 {
		c.emit(claims.Cid, "form_error", map[string]any{
			"scenarioID": scenario.ID,
			"errors":     errs,
		})
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}

	c.appendConversationMessage(svcCtx, claims.Cid, types.AiConversationMessage{
		Role:    "user",
		Content: service.EncodeChatbotFormSubmission(data.Fields),
	})

	vars := service.ChatbotVarsFromMap(data.Fields)
	if err := c.sessionSvc.FinalizeStep(svcCtx, step.ID, scenario, vars); err != nil {
		http.Error(w, fmt.Sprintf("widget: finalize: %v", err), http.StatusInternalServerError)
		return
	}
	c.emit(claims.Cid, "step_complete", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": step.ScenarioIndex,
	})

	c.startStep(svcCtx, cb, claims.Dbsid, claims.Cid, step.ScenarioIndex+1, vars)

	w.WriteHeader(http.StatusNoContent)
}

// advanceStep marks the current step complete (after-automation) and starts the
// next scenario (before-automation + step_start emission).
func (c *Controller) advanceStep(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	cb := chatbotFromCtx(r.Context())
	if claims == nil || cb == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())

	// If a handoff is active or requested, close it first so the session resumes.
	if h, _ := c.sessionSvc.FindHandoffBySession(svcCtx, claims.Dbsid); h != nil && h.ClosedAt == nil {
		_ = c.sessionSvc.CloseHandoff(svcCtx, h.ID)
		c.emit(claims.Cid, "handoff_complete", map[string]any{
			"handoffID": strconv.FormatUint(h.ID, 10),
		})
	}

	step, err := c.sessionSvc.FindStepBySession(svcCtx, claims.Dbsid)
	if err != nil || step == nil {
		http.Error(w, "widget: step not found", http.StatusNotFound)
		return
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		http.Error(w, "widget: scenario not found", http.StatusBadRequest)
		return
	}
	scenario := &cb.Scenarios[step.ScenarioIndex]

	if err := c.sessionSvc.FinalizeStep(svcCtx, step.ID, scenario, nil); err != nil {
		http.Error(w, fmt.Sprintf("widget: cannot advance step: %v", err), http.StatusInternalServerError)
		return
	}
	c.emit(claims.Cid, "step_complete", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": step.ScenarioIndex,
	})

	c.startStep(svcCtx, cb, claims.Dbsid, claims.Cid, step.ScenarioIndex+1, nil)

	w.WriteHeader(http.StatusNoContent)
}

// closeSession finalizes the current step and marks the session closed.
// Intended for the final outro scenario.
func (c *Controller) closeSession(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	cb := chatbotFromCtx(r.Context())
	if claims == nil || cb == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())

	if step, err := c.sessionSvc.FindStepBySession(svcCtx, claims.Dbsid); err == nil && step != nil {
		if step.ScenarioIndex >= 0 && step.ScenarioIndex < len(cb.Scenarios) {
			scenario := &cb.Scenarios[step.ScenarioIndex]
			_ = c.sessionSvc.FinalizeStep(svcCtx, step.ID, scenario, nil)
			c.emit(claims.Cid, "step_complete", map[string]any{
				"scenarioID":    scenario.ID,
				"scenarioIndex": step.ScenarioIndex,
			})
		}
	}

	if err := c.sessionSvc.UpdateStatus(svcCtx, claims.Dbsid, "closed"); err != nil {
		http.Error(w, "widget: cannot close session", http.StatusInternalServerError)
		return
	}
	c.emit(claims.Cid, "session_closed", map[string]any{})

	w.WriteHeader(http.StatusNoContent)
}

// requestHandoff user requests handoff from AI to human
func (c *Controller) requestHandoff(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	if claims == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := readJSON(r.Body, &body); err != nil {
		http.Error(w, "widget: bad body", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())

	step, err := c.sessionSvc.FindStepBySession(svcCtx, claims.Dbsid)
	if err != nil || step == nil {
		http.Error(w, "widget: step not found", http.StatusNotFound)
		return
	}

	h, err := c.sessionSvc.RequestHandoff(svcCtx, claims.Dbsid, step.ID)
	if err != nil {
		http.Error(w, "widget: cannot request handoff", http.StatusInternalServerError)
		return
	}

	c.emit(claims.Cid, "handoff_requested", map[string]any{
		"handoffID": strconv.FormatUint(h.ID, 10),
		"reason":    body.Reason,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "handoff_requested",
		"handoffID": strconv.FormatUint(h.ID, 10),
	})
}

// sendOperatorMessage: operator sends message to widget user. Persists into the
// AiConversation as an assistant turn (prefixed) so the LLM sees it on resume,
// and broadcasts operator_message SSE to the widget. Auto-activates the
// handoff on first operator turn.
func (c *Controller) sendOperatorMessage(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	if claims == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	var body struct {
		Message   string `json:"message"`
		HandoffID string `json:"handoffID"`
		Operator  string `json:"operator,omitempty"`
	}
	if err := readJSON(r.Body, &body); err != nil {
		http.Error(w, "widget: bad body", http.StatusBadRequest)
		return
	}

	hid, err := strconv.ParseUint(body.HandoffID, 10, 64)
	if err != nil || hid == 0 {
		http.Error(w, "widget: bad handoff id", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())

	h, err := c.sessionSvc.FindHandoffByID(svcCtx, hid)
	if err != nil || h == nil {
		http.Error(w, "widget: handoff not found", http.StatusNotFound)
		return
	}
	if h.Status == "requested" {
		if err := c.sessionSvc.ActivateHandoff(svcCtx, hid); err != nil {
			http.Error(w, "widget: cannot activate handoff", http.StatusInternalServerError)
			return
		}
		c.emit(claims.Cid, "handoff_active", map[string]any{
			"handoffID": body.HandoffID,
			"operator":  body.Operator,
		})
	}

	c.appendConversationMessage(svcCtx, claims.Cid, types.AiConversationMessage{
		Role:    "assistant",
		Content: "[operator] " + body.Message,
	})

	c.emit(claims.Cid, "operator_message", map[string]any{
		"content":  body.Message,
		"operator": body.Operator,
	})

	w.WriteHeader(http.StatusNoContent)
}

// closeHandoff terminates the handoff. Flow resumes at the same step with
// session.status="active"; the AI loop accepts further user messages.
func (c *Controller) closeHandoff(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	if claims == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	var body struct {
		HandoffID string `json:"handoffID"`
	}
	if err := readJSON(r.Body, &body); err != nil {
		http.Error(w, "widget: bad body", http.StatusBadRequest)
		return
	}

	hid, err := strconv.ParseUint(body.HandoffID, 10, 64)
	if err != nil {
		http.Error(w, "widget: bad handoff id", http.StatusBadRequest)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	if err := c.sessionSvc.CloseHandoff(svcCtx, hid); err != nil {
		http.Error(w, "widget: cannot close handoff", http.StatusInternalServerError)
		return
	}

	c.emit(claims.Cid, "handoff_complete", map[string]any{
		"handoffID": body.HandoffID,
	})

	w.WriteHeader(http.StatusNoContent)
}

// stream opens an SSE connection filtered to the JWT's conversation ID.
func (c *Controller) stream(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	if claims == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}
	if c.obsBus == nil {
		http.Error(w, "widget: no bus", http.StatusServiceUnavailable)
		return
	}

	d := observability.NewSSEDispatcher(claims.Cid)
	c.obsBus.Register(d)
	defer d.Close()

	observability.PumpSSE(w, r, d)
}

// ---- step orchestration -----------------------------------------------------

// startStep runs the scenario's before-automation, persists the new step row,
// and broadcasts step_start with the scenario payload. If the index is beyond
// the last scenario, emits session_closed instead.
func (c *Controller) startStep(
	ctx context.Context,
	cb *types.Chatbot,
	sessionID, convID uint64,
	scenarioIndex int,
	vars *expr.Vars,
) {
	svcCtx := auth.SetIdentityToContext(ctx, auth.ServiceUser())

	if scenarioIndex >= len(cb.Scenarios) {
		_ = c.sessionSvc.UpdateStatus(svcCtx, sessionID, "closed")
		c.emit(convID, "session_closed", map[string]any{})
		return
	}

	scenario := &cb.Scenarios[scenarioIndex]
	step, err := c.sessionSvc.StartStep(svcCtx, sessionID, scenario, convID, scenarioIndex, vars)
	if err != nil {
		c.emit(convID, "step_error", map[string]any{
			"scenarioID":    scenario.ID,
			"scenarioIndex": scenarioIndex,
			"phase":         "before",
			"error":         err.Error(),
		})
		return
	}

	c.emit(convID, "step_start", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": scenarioIndex,
		"type":          scenario.Type,
		"config":        service.ChatbotScenarioConfig(scenario),
		"stepID":        strconv.FormatUint(step.ID, 10),
	})
}

// ---- helpers ----------------------------------------------------------------

// emit broadcasts an SSE event for the conversation. No-op when obsBus is nil.
func (c *Controller) emit(convID uint64, name string, details map[string]any) {
	if c.obsBus == nil {
		return
	}
	c.obsBus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          name,
		Details:        details,
		Timestamp:      time.Now(),
	})
}

// appendConversationMessage appends a message to the AiConversation. Best-effort.
func (c *Controller) appendConversationMessage(ctx context.Context, convID uint64, msg types.AiConversationMessage) {
	conv, err := store.LookupAiConversationByID(ctx, c.store, convID)
	if err != nil || conv == nil {
		return
	}
	conv.Messages = append(conv.Messages, msg)
	_ = store.UpdateAiConversation(ctx, c.store, conv)
}


func okOptions(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }

func chatbotFromCtx(ctx context.Context) *types.Chatbot {
	if v, ok := ctx.Value(ctxKeyChatbot).(*types.Chatbot); ok {
		return v
	}
	return nil
}

func claimsFromCtx(ctx context.Context) *sessionClaims {
	if v, ok := ctx.Value(ctxKeyClaims).(*sessionClaims); ok {
		return v
	}
	return nil
}

func parseTTL(s string, def time.Duration) time.Duration {
	if s == "" {
		return def
	}
	if d, err := time.ParseDuration(s); err == nil && d > 0 {
		return d
	}
	return def
}

// publicConfigChatbot redacts operator-only fields (allowedOrigins, handoff roles).
func publicConfigChatbot(cb *types.Chatbot) map[string]any {
	styling := cb.Styling
	if styling.LogoAttachmentID != 0 {
		styling.LogoURL = fmt.Sprintf("/api/widget/v1/asset/%d?widgetKey=%s", styling.LogoAttachmentID, cb.WidgetKey)
	}
	if styling.Launcher.IconAttachmentID != 0 {
		styling.Launcher.IconURL = fmt.Sprintf("/api/widget/v1/asset/%d?widgetKey=%s", styling.Launcher.IconAttachmentID, cb.WidgetKey)
	}
	return map[string]any{
		"styling":   styling,
		"scenarios": cb.Scenarios,
		"handoff": map[string]any{
			"enabled":        cb.Handoff.Enabled,
			"notImplemented": false,
		},
	}
}

func randomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON(r io.Reader, dst any) error {
	return json.NewDecoder(r).Decode(dst)
}

// peekJSONString reads up to 4KB of the body looking for a string field
// without breaking the ability of a downstream handler to re-parse it.
func peekJSONString(r *http.Request, field string) string {
	if r.Body == nil || r.ContentLength == 0 {
		return ""
	}
	buf, err := io.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		return ""
	}
	r.Body = io.NopCloser(readSeekerFromBytes(buf))

	var m map[string]any
	if err := json.Unmarshal(buf, &m); err != nil {
		return ""
	}
	s, _ := m[field].(string)
	return s
}

func readSeekerFromBytes(b []byte) io.Reader {
	return &bytesReader{b: b}
}

type bytesReader struct {
	b []byte
	i int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
