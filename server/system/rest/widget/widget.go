package widget

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// Controller wires together the public chatbot widget HTTP API.
// All endpoints live under /api/widget/v1 and are unauthenticated except
// by the per-session JWT issued from POST /session.
type Controller struct {
	chatbotLookup *ChatbotByKeyLookup
	obsBus        *observability.Bus
	convStore     AiConversationStore
	store         store.Storer
	serverSecret  string
	sessionSvc    ChatbotSessionService
}

// AiConversationStore is satisfied by service.DefaultAiConversation. Kept for
// the lookup middleware that needs to talk to the store directly.
type AiConversationStore interface {
	Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
}

// ChatbotSessionService is satisfied by service.ChatbotSession().
type ChatbotSessionService interface {
	Open(ctx context.Context, cb *types.Chatbot) (*types.ChatbotSession, *types.AiConversation, error)
	Start(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64, scenarioIndex int)
	SubmitMessage(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64, input string) error
	SubmitForm(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64, fields map[string]string) (map[string]string, error)
	AdvanceStep(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64) error
	CloseSession(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64) error
	RequestHandoffPublic(ctx context.Context, sessionID, convID uint64, reason string) (*types.ChatbotSessionHandoff, error)
	AcceptHandoffByID(ctx context.Context, cb *types.Chatbot, sessionID, convID, handoffID uint64, operatorID uint64, operatorName string) error
	SendOperatorMessage(ctx context.Context, cb *types.Chatbot, sessionID, convID, handoffID uint64, msg, operator string) error
	CloseHandoffPublic(ctx context.Context, convID, handoffID uint64) error
}

func New(
	s store.Storer,
	obsBus *observability.Bus,
	_ any, // legacy runtime; retained for caller compat, unused
	conv AiConversationStore,
	serverSecret string,
	sessionSvc ChatbotSessionService,
) *Controller {
	return &Controller{
		chatbotLookup: NewChatbotByKeyLookup(s),
		obsBus:        obsBus,
		convStore:     conv,
		store:         s,
		serverSecret:  serverSecret,
		sessionSvc:    sessionSvc,
	}
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

				r.Options("/session/{id}/handoff-accept", okOptions)
				r.Post("/session/{id}/handoff-accept", c.acceptHandoff)

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

// createSession opens a new ChatbotSession + AiConversation, mints a JWT and
// kicks off step 0 async.
func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	cb := chatbotFromCtx(r.Context())
	if cb == nil {
		http.Error(w, "widget: no chatbot", http.StatusForbidden)
		return
	}

	session, conv, err := c.sessionSvc.Open(r.Context(), cb)
	if err != nil {
		c.writeWidgetError(w, err)
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

	// Kick off step 0 asynchronously so SSE listeners attach in time.
	go c.sessionSvc.Start(context.Background(), cb, session.ID, conv.ID, 0)
}

// submit dispatches on payload.type — message or form.
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

	switch body.Type {
	case "message":
		var d struct {
			Input string `json:"input"`
		}
		if err := json.Unmarshal(body.Data, &d); err != nil {
			http.Error(w, "widget: bad message payload", http.StatusBadRequest)
			return
		}
		if err := c.sessionSvc.SubmitMessage(r.Context(), cb, claims.Dbsid, claims.Cid, d.Input); err != nil {
			c.writeWidgetError(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	case "form":
		var d struct {
			Fields map[string]string `json:"fields"`
		}
		if err := json.Unmarshal(body.Data, &d); err != nil {
			http.Error(w, "widget: bad form payload", http.StatusBadRequest)
			return
		}
		errs, err := c.sessionSvc.SubmitForm(r.Context(), cb, claims.Dbsid, claims.Cid, d.Fields)
		if err != nil {
			c.writeWidgetError(w, err)
			return
		}
		if len(errs) > 0 {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "widget: unknown submit type", http.StatusBadRequest)
	}
}

func (c *Controller) advanceStep(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	cb := chatbotFromCtx(r.Context())
	if claims == nil || cb == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}
	if err := c.sessionSvc.AdvanceStep(r.Context(), cb, claims.Dbsid, claims.Cid); err != nil {
		c.writeWidgetError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) closeSession(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	cb := chatbotFromCtx(r.Context())
	if claims == nil || cb == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}
	if err := c.sessionSvc.CloseSession(r.Context(), cb, claims.Dbsid, claims.Cid); err != nil {
		c.writeWidgetError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
	h, err := c.sessionSvc.RequestHandoffPublic(r.Context(), claims.Dbsid, claims.Cid, body.Reason)
	if err != nil {
		c.writeWidgetError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "handoff_requested",
		"handoffID": strconv.FormatUint(h.ID, 10),
	})
}

func (c *Controller) acceptHandoff(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	if claims == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}
	var body struct {
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
	if err := c.sessionSvc.AcceptHandoffByID(r.Context(), chatbotFromCtx(r.Context()), claims.Dbsid, claims.Cid, hid, 0, body.Operator); err != nil {
		c.writeWidgetError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
	if err := c.sessionSvc.SendOperatorMessage(r.Context(), chatbotFromCtx(r.Context()), claims.Dbsid, claims.Cid, hid, body.Message, body.Operator); err != nil {
		c.writeWidgetError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
	if err := c.sessionSvc.CloseHandoffPublic(r.Context(), claims.Cid, hid); err != nil {
		c.writeWidgetError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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

// ---- helpers ----------------------------------------------------------------

// writeWidgetError maps service-layer sentinel errors to HTTP status codes.
func (c *Controller) writeWidgetError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ChatbotSessionErrChatbotNoScenarios()):
		http.Error(w, "widget: "+err.Error(), http.StatusServiceUnavailable)
	case errors.Is(err, service.ChatbotSessionErrScenarioOutOfRange()),
		errors.Is(err, service.ChatbotSessionErrBadHandoffID()),
		errors.Is(err, service.ChatbotSessionErrScenarioNotMessage()):
		http.Error(w, "widget: "+err.Error(), http.StatusBadRequest)
	case errors.Is(err, service.ChatbotSessionErrSessionNotActive()),
		errors.Is(err, service.ChatbotSessionErrScenarioNotForm()),
		errors.Is(err, service.ChatbotSessionErrHandoffNotRequested()),
		errors.Is(err, service.ChatbotSessionErrHandoffNotActive()):
		http.Error(w, "widget: "+err.Error(), http.StatusConflict)
	case errors.Is(err, service.ChatbotSessionErrSessionNotFound()):
		http.Error(w, "widget: "+err.Error(), http.StatusUnauthorized)
	case errors.Is(err, service.ChatbotSessionErrNoActiveStep()),
		errors.Is(err, service.ChatbotSessionErrHandoffNotFound()):
		http.Error(w, "widget: "+err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ChatbotSessionErrAgentUnavailable()),
		errors.Is(err, service.ChatbotSessionErrAgentNotConfigured()):
		http.Error(w, "widget: "+err.Error(), http.StatusServiceUnavailable)
	default:
		http.Error(w, "widget: "+err.Error(), http.StatusInternalServerError)
	}
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
