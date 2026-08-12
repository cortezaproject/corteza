package mcpkit

import (
	"context"
	"net/http"
	"strings"
)

// Scope narrows what a single MCP session may see and do. It is resolved per
// request from the URL and carried in the context.
//
// The two dimensions are deliberately different in kind:
//
//   - Group is presentation. It decides which tools are listed, so a session
//     working on configuration is not shown 100 data tools. It is not a
//     security boundary — see CONVENTIONS.md §2.1 — and a caller can pick any
//     group, because RBAC is what actually governs.
//   - Risk is a ceiling and IS enforced on dispatch, so an out-of-ceiling tool
//     cannot be called even if a client already knew its name. It is a seatbelt
//     against accidents (pointing a session at production and having it delete
//     a namespace), not a lock against a hostile caller, who can simply not set
//     it.
type Scope struct {
	Group   Group
	MaxRisk Risk

	// FullDocs sends every tool's complete description and per-parameter
	// documentation in the listing, instead of the one-line summary.
	//
	// For tooling, not for agents: dev/agent/mcp-verify.py has to audit what the
	// tools actually declare, and a human debugging "why did the model call it
	// like that" needs to read what the model was shown. An agent using it pays
	// ~40k tokens per request instead of ~7k for a surface it can already call
	// in full, so it is deliberately not mentioned in any tool description —
	// human_tool_load is the supported way to get the same prose for the two or
	// three tools that actually need it.
	FullDocs bool
}

type scopeCtxKey struct{}

// WithScope returns ctx carrying s.
func WithScope(ctx context.Context, s Scope) context.Context {
	return context.WithValue(ctx, scopeCtxKey{}, s)
}

// ScopeFromContext returns the scope for this request. The zero value means
// unnarrowed: every group, no ceiling.
func ScopeFromContext(ctx context.Context) Scope {
	s, _ := ctx.Value(scopeCtxKey{}).(Scope)
	return s
}

// Permits reports whether a tool is visible and callable under this scope.
func (s Scope) Permits(groups []Group, risk Risk) bool {
	if s.Group != "" {
		var found bool
		for _, g := range groups {
			if g == s.Group {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if s.MaxRisk != "" && !risk.AtOrBelow(s.MaxRisk) {
		return false
	}

	return true
}

// scopeFromRequest reads the scope out of the URL.
//
// Group comes from the path segment after /mcp, so a client picks it once in
// its config rather than per call. Risk comes from ?maxRisk= on the same URL.
// Both are optional and an unrecognised value is ignored rather than rejected:
// a typo should not silently narrow the surface to nothing, and the structural
// guarantees do not depend on either being present.
func scopeFromRequest(base string, r *http.Request) Scope {
	var s Scope

	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, base), "/")
	if seg := strings.SplitN(rest, "/", 2)[0]; seg != "" {
		switch Group(seg) {
		case GroupDevelopment, GroupConfiguring, GroupUsage:
			s.Group = Group(seg)
		}
	}

	// `tools=all` is the historical spelling from when this opted out of
	// progressive disclosure. Every tool is listed either way now, so what it
	// selects is the documentation depth; both spellings are accepted because
	// the old one is in checked-in tooling and in people's notes.
	switch {
	case r.URL.Query().Get("docs") == "full", r.URL.Query().Get("tools") == "all":
		s.FullDocs = true
	}

	switch Risk(r.URL.Query().Get("maxRisk")) {
	case RiskRead:
		s.MaxRisk = RiskRead
	case RiskWrite:
		s.MaxRisk = RiskWrite
	case RiskDestructive:
		s.MaxRisk = RiskDestructive
	}

	return s
}
