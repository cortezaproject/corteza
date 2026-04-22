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
	lookup       *AgentByKeyLookup
	obsBus       *observability.Bus
	runtime      AgenticRunner
	convStore    AiConversationCreator
	store        store.Storer
	serverSecret string
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
		lookup:       NewAgentByKeyLookup(s),
		obsBus:       obsBus,
		runtime:      rt,
		convStore:    conv,
		store:        s,
		serverSecret: serverSecret,
	}
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

// readConfig returns a redacted view of AgentChatbot safe to ship to a public
// client. Secrets/ops fields are stripped.
func (c *Controller) readConfig(w http.ResponseWriter, r *http.Request) {
	a := agentFromCtx(r.Context())
	if a == nil {
		http.Error(w, "widget: no agent", http.StatusForbidden)
		return
	}
	out := publicConfig(a)
	writeJSON(w, http.StatusOK, out)
}

// createSession opens a new AiConversation under the agent's service account
// and returns a signed JWT plus initial scenario ID.
func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	a := agentFromCtx(r.Context())
	if a == nil {
		http.Error(w, "widget: no agent", http.StatusForbidden)
		return
	}
	if !a.Invocation.System.Enabled || a.Invocation.System.ServiceAccount == 0 {
		http.Error(w, "widget: agent not configured for widget runtime", http.StatusServiceUnavailable)
		return
	}

	saCtx := c.impersonateServiceAccount(r.Context(), a.Invocation.System.ServiceAccount)
	conv, err := c.convStore.Create(saCtx, &types.AiConversation{AgentID: a.ID})
	if err != nil {
		http.Error(w, "widget: cannot open session", http.StatusInternalServerError)
		return
	}

	ttl := parseTTL(a.Chatbot.SessionTTL, 2*time.Hour)
	now := time.Now()
	claims := sessionClaims{
		Sid: randomID(),
		Aid: a.ID,
		Cid: conv.ID,
		Sa:  a.Invocation.System.ServiceAccount,
		Iat: now.Unix(),
		Exp: now.Add(ttl).Unix(),
	}
	tok, err := signSession(claims, deriveSessionSecret(a.Chatbot.WidgetKey, c.serverSecret))
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

// sendMessage hands the user input to DefaultAgenticRuntime under the service
// account. Responds 202 immediately — any tokens/events flow over SSE.
func (c *Controller) sendMessage(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r.Context())
	if claims == nil {
		http.Error(w, "widget: no session", http.StatusUnauthorized)
		return
	}

	var body struct {
		Input string `json:"input"`
	}
	if err := readJSON(r.Body, &body); err != nil {
		http.Error(w, "widget: bad body", http.StatusBadRequest)
		return
	}

	saCtx := c.impersonateServiceAccount(context.Background(), claims.Sa)
	go func() {
		_, _ = c.runtime.Run(saCtx, &runtime.AgentRequest{
			AgentID:        claims.Aid,
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

func agentFromCtx(ctx context.Context) *types.Agent {
	if v, ok := ctx.Value(ctxKeyAgent).(*types.Agent); ok {
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

// publicConfig redacts operator-only fields (allowedOrigins, handoff roles).
func publicConfig(a *types.Agent) map[string]any {
	c := a.Chatbot
	return map[string]any{
		"styling":   c.Styling,
		"scenarios": c.Scenarios,
		"handoff": map[string]any{
			"enabled":        c.Handoff.Enabled,
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
