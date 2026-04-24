package widget

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
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

// refererOrigin returns scheme://host from a Referer URL, or "" if unparseable.
// Browsers omit the Origin header on some same-origin GETs (notably Firefox)
// but still send Referer; we use it as a fallback source for the allowlist
// check.
func refererOrigin(referer string) string {
	if referer == "" {
		return ""
	}
	u, err := url.Parse(referer)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

// requestOrigin returns the request Origin, falling back to Referer-derived
// origin when Origin is absent.
func requestOrigin(r *http.Request) (origin string, fromReferer bool) {
	if o := r.Header.Get("Origin"); o != "" {
		return o, false
	}
	return refererOrigin(r.Header.Get("Referer")), true
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

// keyOnly resolves the chatbot from widgetKey but skips the Origin allowlist
// and the Enabled gate. Intended for public read-only assets (<img> loads):
// browsers don't reliably send Origin and won't honor CORS on image loads,
// and admin preview needs assets to load while the chatbot is still disabled.
// The widget key is sufficient authz for public branding assets.
func (c *Controller) keyOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		widgetKey := r.URL.Query().Get("widgetKey")
		if widgetKey == "" {
			widgetKey = r.Header.Get("X-Widget-Key")
		}

		cb, err := c.chatbotLookup.Find(r.Context(), widgetKey)
		if err != nil || cb == nil {
			http.Error(w, "widget: forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyChatbot, cb)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

		rawOrigin := r.Header.Get("Origin")
		checkedOrigin, fromReferer := requestOrigin(r)
		if !originAllowed(checkedOrigin, cb.AllowedOrigins) {
			reason := "mismatch"
			if checkedOrigin == "" {
				reason = "empty origin and no referer"
			} else if fromReferer {
				reason = "derived from referer; not in allowlist"
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":   "widget: origin not allowed",
				"origin":  checkedOrigin,
				"header":  rawOrigin,
				"referer": r.Header.Get("Referer"),
				"reason":  reason,
			})
			return
		}
		writeCORS(w, rawOrigin)

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
