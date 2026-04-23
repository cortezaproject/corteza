package widget

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey int

const (
	ctxKeyClaims ctxKey = iota
	ctxKeyChatbot
)

// originAllowed implements the simple allowlist:
//   - if list contains "*": allow all (including same-origin requests that
//     omit the Origin header)
//   - otherwise exact match (scheme+host+port); empty Origin is rejected
func originAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" {
			return true
		}
	}
	if origin == "" {
		return false
	}
	for _, a := range allowed {
		if strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}

// writeCORS sets the per-origin CORS headers. Credentials are NOT enabled —
// widget runs on a third-party origin with its own bearer token.
// Same-origin requests (no Origin header) skip CORS entirely.
func writeCORS(w http.ResponseWriter, origin string) {
	if origin == "" {
		return
	}
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", origin)
	h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	h.Set("Vary", "Origin")
}

// corsForKey builds a middleware that:
//   1. reads ?widgetKey= (or body field on POST) and resolves the agent,
//   2. validates Origin against the agent's AllowedOrigins,
//   3. answers OPTIONS preflights,
//   4. stashes the agent in request context so handlers skip a second lookup.
func (c *Controller) corsForKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		widgetKey := r.URL.Query().Get("widgetKey")
		if widgetKey == "" {
			widgetKey = r.Header.Get("X-Widget-Key")
		}

		// POST /session also carries the key in the body; parse it there.
		if widgetKey == "" && r.Method == http.MethodPost {
			if k := peekJSONString(r, "widgetKey"); k != "" {
				widgetKey = k
			}
		}

		ctx := r.Context()

		cb, err := c.chatbotLookup.Find(ctx, widgetKey)
		if err != nil || cb == nil || !cb.Enabled {
			http.Error(w, "widget: forbidden", http.StatusForbidden)
			return
		}

		origin := r.Header.Get("Origin")
		if !originAllowed(origin, cb.AllowedOrigins) {
			http.Error(w, "widget: origin not allowed", http.StatusForbidden)
			return
		}
		writeCORS(w, origin)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ctx = context.WithValue(ctx, ctxKeyChatbot, cb)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireSessionJWT loads the bearer (or ?token=) JWT, verifies it using a
// secret derived from the current widget key, and stashes claims in context.
// The middleware pairs with corsForKey — agent is already on the ctx.
func (c *Controller) requireSessionJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := bearerToken(r)
		if raw == "" {
			raw = r.URL.Query().Get("token")
		}
		if raw == "" {
			http.Error(w, "widget: missing token", http.StatusUnauthorized)
			return
		}

		cb := chatbotFromCtx(r.Context())
		if cb == nil {
			http.Error(w, "widget: no session", http.StatusUnauthorized)
			return
		}

		claims, err := verifySession(raw, deriveSessionSecret(cb.WidgetKey, c.serverSecret))
		if err != nil {
			http.Error(w, "widget: invalid session", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyClaims, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}
