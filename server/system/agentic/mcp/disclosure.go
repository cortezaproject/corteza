package mcp

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Progressive disclosure: a session starts with a handful of tools and asks for
// the rest.
//
// The problem it solves is measured, not assumed. 80 tools cost ~27,000 tokens
// in every request before any work happens, and the surface is heading for
// ~200. Group endpoints help but not enough — `configuring` alone is 64 of the
// 80 — because most tools genuinely are configuration.
//
// Trimming descriptions is the obvious alternative and the wrong one: §8.5
// deliberately made them rich because a caller working through Claude Code
// against a hosted instance has no source to read. Disclosure resolves that
// tension rather than trading against it — a description can be as long as it
// needs to be if only six are loaded by default. It also scales the right way:
// 200 tools cost a session the same as 20.
//
// The mechanism is the one Claude Code itself uses on its own tools: a small
// always-on set plus a search that pulls in the rest by name.

// alwaysOn is what a session sees before it searches for anything.
//
// The two meta-tools, plus the compose read path — because "what namespaces
// exist, what shape is this module, what records are in it" is how almost every
// task opens, and making that cost a search round-trip would be a tax on the
// common case. Everything that writes is behind a search.
var alwaysOn = map[string]bool{
	toolSearchName:             true,
	toolLoadName:               true,
	"compose_namespace_lookup": true,
	"compose_module_lookup":    true,
	"compose_record_lookup":    true,
}

const (
	toolSearchName = "human_tool_search"
	toolLoadName   = "human_tool_load"

	// maxSearchResults bounds a vague query. A search that returned sixty tools
	// would reintroduce exactly the payload this exists to avoid.
	maxSearchResults = 12
)

// disclosure tracks which tools each session has pulled in.
//
// Keyed by MCP session ID and cleaned up on session teardown, so a long-lived
// server does not accumulate state for clients that have gone away.
type disclosure struct {
	mu     sync.RWMutex
	loaded map[string]map[string]bool
}

func newDisclosure() *disclosure {
	return &disclosure{loaded: make(map[string]map[string]bool)}
}

func (d *disclosure) load(sessionID string, names ...string) {
	if sessionID == "" {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.loaded[sessionID] == nil {
		d.loaded[sessionID] = make(map[string]bool, len(names))
	}
	for _, n := range names {
		d.loaded[sessionID][n] = true
	}
}

func (d *disclosure) isLoaded(sessionID, name string) bool {
	if sessionID == "" {
		return false
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.loaded[sessionID][name]
}

func (d *disclosure) forget(sessionID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.loaded, sessionID)
}

// sessionID returns the MCP session for this request, or "" when there is none.
//
// A request with no session cannot accumulate loaded tools, so it sees the
// always-on set and whatever a search returns inline. That is the degradation
// path for a client that does not maintain a session, and it still works.
func sessionID(ctx context.Context) string {
	if s := server.ClientSessionFromContext(ctx); s != nil {
		return s.SessionID()
	}
	return ""
}

// searchTools matches a query against tool names, titles and descriptions.
//
// Substring rather than anything cleverer: the names are already structured
// (`{app}_{resource}_{op}`), so "page" and "delete role" both hit reliably, and
// a miss is obvious rather than mysterious. An index or an embedding would add
// a dependency and a failure mode to the one tool that is always loaded.
func (m *MCPServer) searchTools(query string, scope Scope) []mcp.Tool {
	terms := strings.Fields(strings.ToLower(query))

	type scored struct {
		tool  mcp.Tool
		score int
	}

	var hits []scored
	for _, t := range m.reg.Tools() {
		if m.reg.IsHidden(t.Name) {
			continue
		}
		if !scope.Permits(GroupsOf(t), RiskOf(t)) {
			continue
		}
		if alwaysOn[t.Name] {
			continue
		}

		name := strings.ToLower(t.Name)
		desc := strings.ToLower(t.Description)

		score, matchedAll := 0, true
		for _, term := range terms {
			switch {
			case strings.Contains(name, term):
				// A name hit is worth more than a description hit: someone
				// searching "role" wants the role tools, not every tool whose
				// description happens to mention roles.
				score += 10
			case strings.Contains(desc, term):
				score++
			default:
				matchedAll = false
			}
		}

		if !matchedAll || len(terms) == 0 {
			continue
		}
		hits = append(hits, scored{t, score})
	}

	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].score != hits[j].score {
			return hits[i].score > hits[j].score
		}
		return hits[i].tool.Name < hits[j].tool.Name
	})

	out := make([]mcp.Tool, 0, maxSearchResults)
	for i, h := range hits {
		if i >= maxSearchResults {
			break
		}
		out = append(out, h.tool)
	}
	return out
}
