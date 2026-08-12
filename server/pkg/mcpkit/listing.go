package mcpkit

import (
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Slim listing: every tool is listed and callable from the first request, but
// summarised to its first sentence, with the per-parameter prose left out.
//
// The numbers are measured, not assumed. The full surface is 96 tools and
// ~40,000 tokens, of which ~28,000 is prose: 15,700 in tool descriptions and
// 12,200 in per-parameter descriptions. What a client needs to *register and
// call* a tool — names, types, required fields — is only ~12,400, and slimming
// lands the whole listing at ~7,400.
//
// This replaced progressive disclosure, which showed five tools until a session
// searched and cost ~1,500. That was cheaper and it did not work. Disclosure
// makes a tool callable by adding it to a later tools/list, so it needs a client
// that honours notifications/tools/list_changed and re-lists on one session. The
// claude.ai connector does not: its search returned the full definitions and its
// own note promised they were "now available", and every call then came back
// "tool not found" from the client, having never reached Human at all. Returning
// definitions inline was supposed to be the fallback for exactly that client,
// and it is not one — a client that will not call an unlisted tool is not
// helped by being handed its schema.
//
// Slimming has no such dependency. Nothing has to arrive later, so nothing can
// fail to arrive. human_tool_load survives as documentation rather than
// registration: the tool it describes was already callable, and loading only
// adds the prose the listing dropped. That works on every MCP client, because
// it is just text in a tool result.
//
// The cost is ~7,400 tokens per request against disclosure's ~1,500, on a
// surface heading for ~200 tools. Slimming scales with the tool count where
// disclosure was flat, which is the real thing being traded away — and the
// mitigation is that summarising is per-tool, so 200 tools cost ~15,000 rather
// than the ~84,000 a full listing would.

const (
	toolSearchName = "human_tool_search"
	toolLoadName   = "human_tool_load"

	// maxSearchResults bounds a vague query. Search now returns full
	// documentation for each hit, which is the expensive part, so the ceiling
	// matters more than it did when it merely bounded a listing.
	maxSearchResults = 12
)

// loadHint is appended to the summary of a tool marked NeedsFullDocs.
//
// Phrased as an instruction rather than a note because a model reading "see the
// full documentation" treats it as optional and calls the tool anyway.
const loadHint = " IMPORTANT: call " + toolLoadName + " for this tool before using it — its full " +
	"documentation carries rules that are not visible in the parameter list, and guessing them " +
	"produces a call that succeeds and is wrong."

// slimTool returns t summarised for the listing: first sentence only, and no
// per-parameter prose.
//
// The copying is not incidental. mcp.Tool is a value, but InputSchema.Properties
// is a map shared with the Registry, and so is every property inside it.
// Stripping in place would delete those descriptions permanently — for every
// later request, for human_tool_load, and for Human's own in-process agentic
// runtime, which reads the same tools. The first slim listing would quietly
// destroy the documentation it exists to defer.
func slimTool(t mcp.Tool) mcp.Tool {
	// Shallow copy. Name, Annotations, Required and the Meta pointer are shared
	// with the registry and none of them is written below.
	out := t

	out.Description = summarize(t.Description)
	if WantsFullDocs(t) {
		out.Description += loadHint
	}

	if len(t.InputSchema.Properties) == 0 {
		return out
	}

	props := make(map[string]any, len(t.InputSchema.Properties))
	for name, raw := range t.InputSchema.Properties {
		prop, ok := raw.(map[string]any)
		if !ok {
			props[name] = raw
			continue
		}

		cp := make(map[string]any, len(prop))
		for k, v := range prop {
			if k == "description" {
				continue
			}
			cp[k] = v
		}
		props[name] = cp
	}
	out.InputSchema.Properties = props

	return out
}

// abbreviations end in a period without ending a sentence. Without this, the
// summary of a description whose first sentence contains "e.g." stops there,
// and the listing says half of what the tool does.
var abbreviations = map[string]bool{
	"e.g": true, "i.e": true, "etc": true, "cf": true, "vs": true, "no": true,
}

// summarize cuts a description to its first sentence, or its first line when
// that comes sooner.
//
// A line break first is deliberate: several descriptions open with a one-line
// statement of what the tool is and then a paragraph of rules, and cutting at
// the break gives a better summary than hunting for a period that may be
// several sentences further on.
func summarize(desc string) string {
	desc = strings.TrimSpace(desc)

	for i := 0; i < len(desc)-1; i++ {
		if desc[i] == '\n' {
			return strings.TrimSpace(desc[:i])
		}

		if desc[i] != '.' || desc[i+1] != ' ' {
			continue
		}

		word := desc[:i]
		if j := strings.LastIndexAny(word, " ("); j >= 0 {
			word = word[j+1:]
		}
		if abbreviations[strings.ToLower(word)] {
			continue
		}

		return desc[:i+1]
	}

	return desc
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
