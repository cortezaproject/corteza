package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// discovery_search is one of the three grandfathered names in
// mcp/CONVENTIONS.md §6: two segments, app "discovery" in package
// system/agentic, and hardcoded in eight places in runtime/executor.go. It is
// kept rather than renamed.
//
// Known divergence, recorded in CONVENTIONS §2.3 and deliberately not fixed
// here: the 'namespace' and 'module' params declared below are never read by
// the handler. It reads 'namespaceIDs'/'moduleIDs', which only the in-process
// executor injects, so over HTTP the declared args are silently ignored.
//
// This file holds declarations only. Implementations are in
// discovery_handler.go, in the same order.

func (h *discoveryHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("discovery_search",
			mcp.WithDescription("Search or list records you have access to. Use this when the user asks to find, list, or show existing records. Only namespaces listed in your DISCOVERY ACCESS section are permitted — the executor will deny any other namespace. Leave query empty to list all accessible records, or provide a specific field value (name, email, phone) to search within them."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("The exact namespace name from your DISCOVERY ACCESS section. Do not use module names here.")),
			mcp.WithString("module", mcp.Description("The exact module name from your DISCOVERY ACCESS section. Leave empty to search across all accessible modules in the namespace.")),
			mcp.WithString("query", mcp.Description("A specific value to search for inside record fields (e.g. a name, email, phone). Leave empty to list all records.")),
			mcp.WithString("size", mcp.Description("Number of results to return (default: 10)")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Discover Records",
		h.search,
		hmcp.Available(h.isAvailable),
	)
}
