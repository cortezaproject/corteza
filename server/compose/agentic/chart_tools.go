package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in chart_handler.go, in this order.

// canonical config contract validated against the unify webapp renderer
// (client/web/unify/src/sections/compose/components/Chart/, lib/js chart types)
const chartConfigDoc = `JSON chart configuration. Canonical shape:
{"colorScheme":"tableau.Tableau10","reports":[{"moduleID":"<id from compose_module_lookup>","filter":"","dimensions":[{"field":"<grouping field>","modifier":"(no grouping / buckets)","conditions":{}}],"metrics":[{"field":"count","type":"doughnut"}]}]}
Metric "field":"count" counts records; a numeric module field aggregates instead. "type" per metric: doughnut, pie, bar, line. The dimension field is what values are grouped by (Select fields work well).`

// Callers ask for a report or a dashboard; Human calls the thing a chart.
var chartKeywords = hmcp.WithKeywords("report", "dashboard", "graph", "analytics", "visualisation")

func (h *chartHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_lookup",
			mcp.WithDescription(
				"List charts in a namespace, or look one up by name, handle, or ID. Provide 'chart' to fetch "+
					"a single chart in full, including its config; omit it to list every chart in the "+
					"namespace. A listing is trimmed to chartID, name and handle — it never carries the "+
					"config, so fetch the chart itself before changing its configuration. "+
					"Chart blocks on pages reference charts by chartID.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("chart", mcp.Description("Chart name, handle, or ID (as string to prevent precision loss). Omit to list all charts in the namespace.")),
			mcp.WithString("limit", mcp.Description("Maximum charts to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			chartKeywords,
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup chart",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_create",
			mcp.WithDescription("Create a chart in a namespace. After creating, place it on a page with a Chart block: {\"kind\":\"Chart\",\"options\":{\"chartID\":\"<created ID>\"}}."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Chart name")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed.")),
			mcp.WithString("config", mcp.Required(), mcp.Description(chartConfigDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			chartKeywords,
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Create chart",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_update",
			mcp.WithDescription(
				"Update an existing chart. A field you omit is left unchanged; passing an empty string for "+
					"'name' or 'handle' clears it. 'config' replaces the whole configuration rather than "+
					"merging into it, so send the complete config, not only the part you are changing — "+
					"call compose_chart_lookup with 'chart' first if you need the current one.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("chart", mcp.Required(), mcp.Description("Chart name, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("name", mcp.Description("New name")),
			mcp.WithString("handle", mcp.Description("New handle")),
			mcp.WithString("config", mcp.Description(chartConfigDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			chartKeywords,
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update chart",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_delete",
			mcp.WithDescription("Delete a chart by name, handle, or ID. Remove Chart blocks referencing it from pages first — they render empty otherwise. The delete is soft and compose_chart_undelete reverses it, but note the chart ID first: a deleted chart can no longer be found by name or handle."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("chart", mcp.Required(), mcp.Description("Chart name, handle, or ID (as string to prevent precision loss)")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			chartKeywords,
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete chart",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted chart, reversing compose_chart_delete. A delete only marks the chart — "+
					"its name, handle and full report configuration were retained — so it comes back unchanged "+
					"and Chart blocks that reference its chartID render again, without the page needing an edit. "+
					"Requires the numeric chartID: deleted charts are excluded from every lookup path, so "+
					"compose_chart_lookup can no longer resolve one by name or handle, and it exposes no "+
					"includeDeleted-style filter. Use the ID compose_chart_delete reported, or one from a "+
					"compose_chart_lookup taken before the delete. "+
					"Calling this on a chart that is not deleted is accepted and changes nothing.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("chartID", mcp.Required(), mcp.Description("ID of the deleted chart (as string to prevent precision loss). A name or handle will not work — deleted charts are not resolvable by either.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			chartKeywords,
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete chart",
		h.undelete,
	)
}
