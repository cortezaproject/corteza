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
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/types"
)

// Controller wires together the public chatbot widget HTTP API.
// All endpoints live under /api/widget/v1 and are unauthenticated except
// by the per-session JWT issued from POST /session.
type Controller struct {
	chatbotLookup *ChatbotByKeyLookup
	obsBus        *observability.Bus
	runtime       AgenticRunner
	convStore     AiConversationCreator
	store         store.Storer
	serverSecret  string
}

// AgenticRunner is the minimal subset of the runtime the widget needs.
type AgenticRunner interface {
	Run(ctx context.Context, req *runtime.AgentRequest) (*runtime.AgentResponse, error)
}

// AiConversationCreator is satisfied by service.DefaultAiConversation.
type AiConversationCreator interface {
	Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
}

func New(
	s store.Storer,
	obsBus *observability.Bus,
	rt AgenticRunner,
	conv AiConversationCreator,
	serverSecret string,
) *Controller {
	return &Controller{
		chatbotLookup: NewChatbotByKeyLookup(s),
		obsBus:        obsBus,
		runtime:       rt,
		convStore:     conv,
		store:         s,
		serverSecret:  serverSecret,
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
	// resolve via service identity so RoleMember search itself isn't gated
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
// The caller is expected to mount this group under the API base URL and
// *outside* the admin token validator so visitors on third-party sites can
// reach it (canonical URL becomes /api/widget/v1/...).
func (c *Controller) MountRoutes(r chi.Router) {
	r.Route("/widget/v1", func(r chi.Router) {
		r.Use(c.corsForKey)

		r.Options("/config", okOptions)
		r.Get("/config", c.readConfig)

		r.Options("/session", okOptions)
		r.Post("/session", c.createSession)

		r.Group(func(r chi.Router) {
			r.Use(c.requireSessionJWT)

			r.Options("/session/{id}/message", okOptions)
			r.Post("/session/{id}/message", c.sendMessage)

			r.Get("/session/{id}/stream", c.stream)
		})
	})
}

// ---- handlers ---------------------------------------------------------------

// readConfig returns a redacted view of the chatbot/agent config safe to
// ship to a public client. Secrets/ops fields are stripped.
func (c *Controller) readConfig(w http.ResponseWriter, r *http.Request) {
	cb := chatbotFromCtx(r.Context())
	if cb == nil {
		http.Error(w, "widget: no chatbot", http.StatusForbidden)
		return
	}
	writeJSON(w, http.StatusOK, publicConfigChatbot(cb))
}

// createSession opens a new AiConversation and returns a signed JWT.
// No agent is bound at session open — each sendMessage names its scenario.
func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	cb := chatbotFromCtx(r.Context())
	if cb == nil {
		http.Error(w, "widget: no chatbot", http.StatusForbidden)
		return
	}

	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	conv, err := c.convStore.Create(svcCtx, &types.AiConversation{})
	if err != nil {
		http.Error(w, "widget: cannot open session", http.StatusInternalServerError)
		return
	}

	ttl := parseTTL(cb.SessionTTL, 2*time.Hour)
	now := time.Now()
	claims := sessionClaims{
		Sid: randomID(),
		Cid: conv.ID,
		Iat: now.Unix(),
		Exp: now.Add(ttl).Unix(),
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
	})
}

// sendMessage runs the scenario's agent against the session's conversation.
// The scenario is identified by the request body's scenarioID; it must be a
// conversation-type scenario with a non-zero agentID. Responds 202 immediately.
func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	cb := chatbotFromCtx(r.Context())
	if claims == nil || cb == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	var body struct {
		ScenarioID string `json:"scenarioID"`
		Input      string `json:"input"`
	}
	if err := readJSON(r.Body, &body); err != nil {
		http.Error(w, "widget: bad body", http.StatusBadRequest)
		return
	}

	var scenario *types.ChatbotScenario
	for i := range cb.Scenarios {
		if cb.Scenarios[i].ID == body.ScenarioID {
			scenario = &cb.Scenarios[i]
			break
		}
	}
	if scenario == nil || scenario.Type != "conversation" || scenario.AgentID == 0 {
		http.Error(w, "widget: scenario not runnable", http.StatusBadRequest)
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

	// First-turn stamp: pin the conversation to this agent if it isn't already,
	// so downstream listings/automation see a stable per-conversation agent.
	svcCtx := auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
	if conv, err := store.LookupAiConversationByID(svcCtx, c.store, claims.Cid); err == nil && conv != nil && conv.AgentID == 0 {
		conv.AgentID = agent.ID
		_ = store.UpdateAiConversation(svcCtx, c.store, conv)
	}

	saCtx := c.impersonateServiceAccount(context.Background(), agent.Invocation.System.ServiceAccount)
	go func() {
		_, _ = c.runtime.Run(saCtx, &runtime.AgentRequest{
			AgentID:        agent.ID,
			Input:          body.Input,
			ConversationID: claims.Cid,
		})
	}()

	w.WriteHeader(http.StatusAccepted)
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

	d := newSSEDispatcher(claims.Cid)
	c.obsBus.Register(d)
	defer d.close()

	pumpSSE(w, r, d)
}

// ---- helpers ----------------------------------------------------------------

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
	return map[string]any{
		"styling":   cb.Styling,
		"scenarios": cb.Scenarios,
		"handoff": map[string]any{
			"enabled":        cb.Handoff.Enabled,
			"notImplemented": true,
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
	// restore body for downstream reads
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

// ensure unused imports don't creep in
var _ = fmt.Sprintf
