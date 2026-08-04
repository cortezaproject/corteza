package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in namespace_handler.go, in this order.

func (h *namespaceHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_lookup",
			mcp.WithDescription(
				"Look up namespaces. Call this whenever the user asks what they have, what exists, what's "+
					"set up, or anything about the current state of their data. Also call this to resolve a "+
					"namespace before any other operation. "+
					"Provide 'namespace' to fetch that one namespace in full; omit it to list namespaces as "+
					"{namespaceID, name, slug}. If the namespace you name is not found, the list is returned "+
					"instead so you can pick the right one without another call.",
			),
			mcp.WithString("namespace", mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss). Omit to list namespaces instead.")),
			mcp.WithString("limit", mcp.Description("Maximum namespaces to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup namespace",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_create",
			mcp.WithDescription("Create a new namespace. A namespace is a top-level container for modules and records in Human Compose."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name for the namespace")),
			mcp.WithString("slug", mcp.Required(), mcp.Description("URL-friendly identifier (lowercase letters, digits, and hyphens only)")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the namespace is enabled (default: true)")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create namespace",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_update",
			mcp.WithDescription(
				"Update an existing namespace's name, slug, or enabled state. Only the fields you send are "+
					"written; fields you omit keep their current value, and a field sent as an empty string "+
					"is cleared.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Description("New display name")),
			mcp.WithString("slug", mcp.Description("New URL-friendly identifier")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the namespace should be enabled")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update namespace",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_delete",
			mcp.WithDescription(
				"Delete a namespace by name, handle, slug, or ID. The delete is soft: the namespace stops "+
					"appearing in lookups but is retained, and the modules, records, pages and charts inside "+
					"it are left in place rather than removed. Any agent tool allow-list entry pointing at "+
					"the namespace is dropped.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete namespace",
		h.del,
	)
}
