package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// discovery_search is one of the three grandfathered names in
// mcp/CONVENTIONS.md §6: two segments, app "discovery" in package
// system/agentic, and hardcoded in eight places in runtime/executor.go. It is
// kept rather than renamed.
//
// It declares no pageCursor, and that is deliberate rather than an oversight.
// The backing discovery service is an Elasticsearch-style API that offsets with
// `from`; it has no cursor to hand back, so a pageCursor param would advertise
// paging that cannot work. `limit` is declared because the service does cap
// result size. See §8.1 for why an advertised-but-broken cursor is worse than
// no cursor.
//
// This file holds declarations only. Implementations are in
// discovery_handler.go, in the same order.

func (h *discoveryHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("discovery_search",
			mcp.WithDescription(
				"Search or list records you have access to. Use this when the user asks to find, list, or "+
					"show existing records. Leave 'query' empty to list records, or give a specific field "+
					"value (a name, email, phone) to search within them. "+
					"Scope the search with 'namespace', and optionally 'module', to avoid searching "+
					"everything you can read. Omitting 'namespace' searches every namespace you have "+
					"access to, which is rarely what you want and can be slow. "+
					"When running as a configured agent, only the namespaces in your DISCOVERY ACCESS "+
					"section are permitted and anything else is denied. "+
					"Results are capped and cannot be paged — narrow the query rather than asking for more.",
			),
			mcp.WithString("namespace", mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss). Omit to search every namespace you can access.")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID (as string to prevent precision loss). Requires 'namespace'. Omit to search all modules in the namespace.")),
			mcp.WithString("query", mcp.Description("A specific value to search for inside record fields (e.g. a name, email, phone). Leave empty to list records.")),
			mcp.WithString("limit", mcp.Description("Maximum records to return, default 50, capped at 200. There is no cursor: this search cannot be paged.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Discover Records",
		h.search,
		hmcp.Available(h.isAvailable),
	)
}
