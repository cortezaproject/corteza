package mcpkit

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

// searchTools matches a query against tool names, keywords and descriptions.
//
// Substring rather than anything cleverer: the names are already structured
// (`{app}_{resource}_{op}`), so "page" and "delete role" both hit reliably, and
// a miss is obvious rather than mysterious. An index or an embedding would add
// a dependency and a failure mode to the one tool that is always loaded.
//
// Three tiers, because two were not enough. A name hit beats a keyword hit
// beats a description hit: the name is what the tool *is*, a keyword is the
// outside word for it, and a description mention is usually incidental. With
// only name and description, a word Human does not use — "dashboard",
// "report", "permission" — could never outrank a tool that mentioned it in
// passing, so the tie fell through to alphabetical order and the answer was
// wrong in a way no amount of description editing could fix.
func (m *MCPServer) searchTools(query string, scope Scope) []mcp.Tool {
	terms := searchTerms(query)

	type scored struct {
		tool    mcp.Tool
		score   int
		matched int
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

		var (
			segments = strings.Split(strings.ToLower(t.Name), "_")
			desc     = strings.ToLower(t.Description)
			keywords = KeywordsOf(t)
		)

		score, matched := 0, 0
		for _, term := range terms {
			if s := segmentScore(segments, term); s > 0 {
				// A name hit is worth more than a description hit: someone
				// searching "role" wants the role tools, not every tool whose
				// description happens to mention roles.
				score += s
				matched++
				continue
			}

			if s := keywordScore(keywords, term); s > 0 {
				score += s
				matched++
				continue
			}

			if strings.Contains(desc, term) {
				score++
				matched++
			}
		}

		// One term is enough to be a candidate; matching more ranks higher.
		//
		// Requiring every term used to be the rule, and it made a natural
		// question the worst possible input: "workflow create tool" returned
		// nothing, because no single tool contains all three words, while
		// "workflow" returned the family. A model reading "no tool matches"
		// concludes the capability does not exist, so the stricter rule did not
		// merely narrow results — it hid the surface.
		if matched == 0 {
			continue
		}
		hits = append(hits, scored{t, score, matched})
	}

	sort.SliceStable(hits, func(i, j int) bool {
		// Breadth first: a tool matching three of the caller's words belongs
		// above one that matched a single word in its name, however strongly.
		if hits[i].matched != hits[j].matched {
			return hits[i].matched > hits[j].matched
		}
		if hits[i].score != hits[j].score {
			return hits[i].score > hits[j].score
		}
		// Shorter name before longer: fewer segments means the more general
		// tool of the family, which is the better guess when nothing else
		// separates them. Falling straight to alphabetical order instead meant
		// `automation_*` won every tie it was in, so "delete a customer" led
		// with automation_taq_delete.
		if len(hits[i].tool.Name) != len(hits[j].tool.Name) {
			return len(hits[i].tool.Name) < len(hits[j].tool.Name)
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

// minTermLen is the shortest word that carries meaning here.
//
// Matching is substring-based, so a one- or two-letter word is a substring of
// nearly every tool name and description: "a" matched everything, which made a
// naturally phrased question rank worse than a single word. "make a dashboard"
// returned twelve unrelated tools while "dashboard" honestly returned none —
// and a confident wrong answer is worse than an admitted miss, because the
// caller acts on it.
const minTermLen = 3

// searchTerms splits a query into the words worth matching on.
//
// A query made entirely of short words keeps them rather than matching nothing:
// the caller meant something, and the old behaviour is better than silence.
func searchTerms(query string) []string {
	all := strings.Fields(strings.ToLower(query))

	kept := make([]string, 0, len(all))
	for _, term := range all {
		if len(term) >= minTermLen {
			kept = append(kept, term)
		}
	}

	if len(kept) == 0 {
		return all
	}
	return kept
}

// segmentScore matches a term against the parts of a `{app}_{resource}_{op}`
// name, rather than against the name as one string.
//
// Whole-name substring matching cannot tell "role" landing on the resource of
// system_role_create from "ole" landing in the middle of a word. Splitting
// first makes an exact segment hit — the common case, since callers use
// Human's own nouns — outrank a partial one.
func segmentScore(segments []string, term string) int {
	best := 0
	for _, seg := range segments {
		switch {
		case seg == term:
			return 20
		case strings.Contains(seg, term) || strings.Contains(term, seg):
			// Covers a plural or a possessive the caller typed: "roles"
			// against the segment "role".
			best = 12
		}
	}
	return best
}

// keywordScore matches a term against a tool's declared synonyms. Scored below
// any name hit and well above a description hit — see WithKeywords for why the
// middle tier has to exist.
func keywordScore(keywords []string, term string) int {
	best := 0
	for _, kw := range keywords {
		switch {
		case kw == term:
			return 8
		case strings.Contains(kw, term) || strings.Contains(term, kw):
			best = 6
		}
	}
	return best
}
