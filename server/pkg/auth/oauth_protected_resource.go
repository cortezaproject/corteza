package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// WellKnownProtectedResource is where RFC 9728 puts the document. A client
// looking for the metadata of a resource at /api/mcp asks the resource's own
// origin for /.well-known/oauth-protected-resource/api/mcp — the resource path
// is appended to the well-known prefix rather than the other way around, which
// is why this has to be mounted on the root router and not under the API.
const WellKnownProtectedResource = "/.well-known/oauth-protected-resource"

// OAuthProtectedResource publishes the two things an OAuth-capable MCP client
// needs to find its way from a rejected request to a login screen:
//
//   - a 401 that says where the metadata is, via WWW-Authenticate
//   - the metadata itself, which names the authorization server
//
// Both are required, and the 401 pointer especially so here. A client that gets
// no pointer falls back to probing the origin's well-known paths, and this
// server answers every unrouted path with an empty 200 rather than a 404, so
// the probe finds a document-shaped nothing instead of failing cleanly.
//
// This is deliberately not in pkg/mcpkit. That package is vendored into
// dev/mcp, which speaks stdio and has no authorization server to discover, and
// its locked contract keeps it free of anything that knows about Human.
type OAuthProtectedResource struct {
	// BasePath is where the protected surface is mounted, as an absolute path
	// including any HTTP_BASE_URL prefix, e.g. "/api/mcp". Metadata is served
	// for this path and anything under it, because a client may narrow its
	// session with a longer path (/api/mcp/{group}) and the document has to
	// echo whichever one it asked about.
	BasePath string

	// Issuer is the authorization server's issuer URL (AUTH_BASE_URL), the
	// same value its OpenID configuration reports as "issuer".
	Issuer string

	// Scopes the client should ask the authorization server for. The access
	// token has to carry "api" or HttpTokenValidator rejects it.
	Scopes []string

	// SslTerminated reports that TLS ends at a proxy in front of us, so a
	// request arriving as plain HTTP was still https:// to the client. Only
	// consulted when the proxy sets no X-Forwarded-Proto.
	SslTerminated bool
}

// MetadataHandler serves the RFC 9728 document.
//
// The resource identifier is built from the request rather than from
// configuration: it has to match the URL the user typed into their client
// exactly, and the request is the only thing that knows what that was. Deriving
// it from DOMAIN via options.FullURL would publish "http://localhost/api/mcp"
// on any deployment that left DOMAIN unset — a document that looks valid and
// fails to connect for a reason nothing in the error message points at.
func (o OAuthProtectedResource) MetadataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resource := strings.TrimPrefix(r.URL.Path, WellKnownProtectedResource)
		resource = strings.TrimSuffix(resource, "/")

		// A bare probe carries no resource path. It can only be asking about
		// the surface this is mounted for.
		if resource == "" {
			resource = o.BasePath
		}

		// Only answer for the surface we actually protect. Without this the
		// handler would happily declare any path on the origin to be an OAuth
		// protected resource, including paths that are not protected at all.
		if !o.covers(resource) {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"resource":                 externalOrigin(r, o.SslTerminated) + resource,
			"authorization_servers":    []string{o.Issuer},
			"scopes_supported":         o.Scopes,
			"bearer_methods_supported": []string{"header"},
		})
	}
}

// Challenge adds the WWW-Authenticate pointer to any 401 the wrapped handler
// produces.
//
// It has to be a response wrapper rather than a header set upfront: the header
// belongs on a 401 and nothing else, and whether this request is going to be
// one is not known until the auth middlewares below have run.
func (o OAuthProtectedResource) Challenge() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				metadata = externalOrigin(r, o.SslTerminated) +
					WellKnownProtectedResource +
					strings.TrimSuffix(r.URL.Path, "/")

				challenge = fmt.Sprintf(
					"Bearer resource_metadata=%q, scope=%q",
					metadata,
					strings.Join(o.Scopes, " "),
				)
			)

			next.ServeHTTP(&challengeWriter{ResponseWriter: w, challenge: challenge}, r)
		})
	}
}

// covers reports whether path is the protected surface or sits under it.
func (o OAuthProtectedResource) covers(path string) bool {
	return path == o.BasePath || strings.HasPrefix(path, o.BasePath+"/")
}

// externalOrigin reconstructs the scheme and host the client used, which is not
// necessarily the one that reached us.
func externalOrigin(r *http.Request, sslTerminated bool) string {
	var (
		scheme = "http"
		host   = r.Host
	)

	if r.TLS != nil || sslTerminated {
		scheme = "https"
	}

	// A proxy that tells us what it terminated is more reliable than the
	// static option, which describes the usual case rather than this request.
	if v := firstForwarded(r.Header.Get("X-Forwarded-Proto")); v != "" {
		scheme = v
	}

	if v := firstForwarded(r.Header.Get("X-Forwarded-Host")); v != "" {
		host = v
	}

	return scheme + "://" + host
}

// firstForwarded takes the client-side value out of a comma-chained X-Forwarded
// header, which accumulates one entry per proxy hop.
func firstForwarded(h string) string {
	return strings.TrimSpace(strings.Split(h, ",")[0])
}

// challengeWriter sets the challenge header if, and only if, the status turns
// out to be 401.
type challengeWriter struct {
	http.ResponseWriter

	challenge string
	written   bool
}

func (w *challengeWriter) WriteHeader(code int) {
	if !w.written {
		w.written = true

		if code == http.StatusUnauthorized && w.Header().Get("WWW-Authenticate") == "" {
			w.Header().Set("WWW-Authenticate", w.challenge)
		}
	}

	w.ResponseWriter.WriteHeader(code)
}

// Flush and Unwrap keep the streamable HTTP transport working. Wrapping a
// ResponseWriter hides the interfaces the embedded one implements, and an MCP
// server that cannot flush cannot stream a server-sent event — the response
// would sit in the buffer until the handler returned.
func (w *challengeWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *challengeWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
