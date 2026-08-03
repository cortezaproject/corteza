package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in record_handler.go, in this order.

func (h *recordHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_lookup",
			mcp.WithDescription(
				"Look up a record by ID, or list and filter records in a module. Provide 'recordID' to fetch "+
					"one; omit it and use 'filter' to search by field values. "+
					"Do NOT call this before creating a record — only use it when the user explicitly asks to "+
					"search or check for existing records. "+
					"Records are returned in full, including their values, so use 'limit' and narrow with "+
					"'filter' rather than listing a whole module: a large module will exceed the result size "+
					"limit and the call will fail.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Description("Record ID (as string to prevent precision loss). Omit to list or filter instead.")),
			mcp.WithString("filter", mcp.Description("Filter expression when no recordID is given, e.g. \"name = 'John'\" or \"status = 'open'\".")),
			mcp.WithString("limit", mcp.Description("Maximum records to return, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup record",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_create",
			mcp.WithDescription("Create a new record. If you do not know the field names, call compose_module_lookup first to get them."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("values", mcp.Required(), mcp.Description("JSON object of field name-value pairs")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create record",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_update",
			mcp.WithDescription(
				"Update an existing record. Requires a record ID — use compose_record_lookup with a filter to "+
					"find it if unknown. Only the fields present in 'values' are written; fields you omit keep "+
					"their current value.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID (as string to prevent precision loss)")),
			mcp.WithString("values", mcp.Required(), mcp.Description("JSON object of field name-value pairs to update")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update record",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_record_delete",
			mcp.WithDescription(
				"Delete a record by ID. Requires a record ID — use compose_record_lookup with a filter to find "+
					"it if unknown. The delete is soft: the record stops appearing in lookups but is retained.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("module", mcp.Required(), mcp.Description("Module name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("recordID", mcp.Required(), mcp.Description("Record ID (as string to prevent precision loss)")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete record",
		h.del,
	)
}
