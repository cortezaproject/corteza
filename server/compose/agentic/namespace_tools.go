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
			// Reading a definition is what a data agent does before it can read the
			// data: a "usage" grant that cannot see a module has no way to name one
			// or to interpret the values it gets back. Writing one stays configuring.
			hmcp.InGroup(hmcp.GroupConfiguring, hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup namespace",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_create",
			mcp.WithDescription(
				"Create a namespace: the top-level container an app is built in, holding its modules, "+
					"records, pages and charts. Call compose_namespace_lookup first — a slug must be unique, "+
					"and the app you are asked to build may already exist. The slug is what filters and "+
					"expressions refer to, so pick it as an identifier, not a label. A namespace starts "+
					"enabled unless told otherwise; create its modules next (compose_module_create), then "+
					"its pages.",
			),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name for the namespace")),
			mcp.WithString("slug", mcp.Required(), mcp.Description("URL-friendly identifier (lowercase letters, digits, and underscores only). Use snake_case: a hyphen is the subtraction operator wherever an identifier is parsed, so underscores keep the slug safe to reuse in filters and expressions.")),
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
					"is cleared. Renaming the slug changes every URL and filter that names it, so do it only "+
					"when asked. Disabling hides the namespace from users without deleting anything; "+
					"compose_namespace_delete is the one that removes it.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Description("New display name")),
			mcp.WithString("slug", mcp.Description("New URL-friendly identifier. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed.")),
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

	h.reg.RegisterTool(
		mcp.NewTool("compose_namespace_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted namespace, reversing compose_namespace_delete. The delete only set a "+
					"marker: the namespace, its modules, records, pages and charts were all retained, so it comes "+
					"back exactly as it was and appears in compose_namespace_lookup again. Agent allow-list "+
					"entries the delete dropped are not restored; re-grant them with system_agent_update. "+
					"Requires the numeric namespaceID, and nothing else will do: a deleted namespace is "+
					"excluded from every lookup path, so neither name nor slug resolves it. Take the ID from "+
					"what compose_namespace_delete reported, or from a compose_namespace_lookup made before the "+
					"delete. Calling this on a namespace that is not deleted is accepted and changes nothing.",
			),
			mcp.WithString("namespaceID", mcp.Required(), mcp.Description("ID of the deleted namespace (as string to prevent precision loss). A name or slug will not work — deleted namespaces are not resolvable by either.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete namespace",
		h.undelete,
	)
}
