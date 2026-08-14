package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in page_layout_handler.go, in this order.

const layoutModelDoc = "A page holds what its blocks ARE — kind, options, title. A LAYOUT holds where " +
	"they go, as {blockID, xywh}, and it is the layout the page renders and the builder draws. " +
	"Every page is created with a primary layout seeded from its blocks, so these tools are for " +
	"rearranging that layout, changing its record toolbar, or adding further layouts — not for " +
	"giving a page its first one."

func (h *pageLayoutHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_layout_lookup",
			mcp.WithDescription(
				"List a page's layouts, or fetch one whole. "+layoutModelDoc+" "+
					"Provide 'layout' to fetch that one with its blocks and config; omit it to list "+
					"the page's layouts as {pageLayoutID, handle, weight, title}. Call this before "+
					"compose_page_layout_update: the update needs the layout's ID, and a block is "+
					"placed by the blockID compose_page_lookup reports.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("layout", mcp.Description("Layout handle or ID (as string). Omit to list the page's layouts.")),
			mcp.WithString("limit", mcp.Description("Maximum layouts to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup page layout",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_layout_create",
			mcp.WithDescription(
				"Add a layout to a page. "+layoutModelDoc+" "+
					"A second layout is an alternative arrangement of the same page — a different "+
					"subset of blocks, or a different record toolbar — shown when its visibility "+
					"expression passes. Layouts are evaluated in weight order and the first whose "+
					"expression passes is the one rendered, so an unconditional layout should sort "+
					"last or it will always win.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier, unique among the page's layouts (lowercase letters, digits, underscores)")),
			mcp.WithString("title", mcp.Description("Layout title, shown in the builder's layout picker")),
			mcp.WithString("blocks", mcp.Description(layoutBlocksDoc)),
			mcp.WithString("buttons", mcp.Description(layoutButtonsDoc)),
			mcp.WithString("visibility", mcp.Description(layoutVisibilityDoc)),
			mcp.WithNumber("weight", mcp.Description("Sort order among the page's layouts; lower is evaluated first.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create page layout",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_layout_update",
			mcp.WithDescription(
				"Change a layout: move or resize its blocks, retitle it, change its record toolbar "+
					"or its visibility expression. "+
					"Every parameter is optional and only what you send is changed — omitting "+
					"'blocks' leaves the arrangement alone, and omitting 'buttons' leaves the "+
					"toolbar alone. Sending 'blocks' REPLACES the arrangement wholesale, so send "+
					"every block that should be placed, not only the ones that moved.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("layout", mcp.Required(), mcp.Description("Layout handle or ID (as string), from compose_page_layout_lookup")),
			mcp.WithString("handle", mcp.Description("New URL-friendly identifier")),
			mcp.WithString("title", mcp.Description("New layout title")),
			mcp.WithString("blocks", mcp.Description(layoutBlocksDoc)),
			mcp.WithString("buttons", mcp.Description(layoutButtonsDoc)),
			mcp.WithString("visibility", mcp.Description(layoutVisibilityDoc)),
			mcp.WithNumber("weight", mcp.Description("Sort order among the page's layouts; lower is evaluated first.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update page layout",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_layout_delete",
			mcp.WithDescription(
				"Delete a layout from a page. The page and its blocks are untouched — only the "+
					"arrangement goes. Deleting a page's only layout leaves it with nothing to "+
					"render and is rejected; delete the page instead, with compose_page_delete.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("layout", mcp.Required(), mcp.Description("Layout handle or ID (as string), from compose_page_layout_lookup")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete page layout",
		h.delete,
	)
}

const layoutBlocksDoc = `JSON array placing the page's blocks: [{"blockID":"1","xywh":[0,0,24,20]}]. ` +
	`blockID comes from compose_page_lookup; a block the layout does not name is not drawn in it. ` +
	`The grid is 48 columns wide and a cell is 10px tall — full width is w=48, half 24, a quarter 12 ` +
	`— and blocks sit side by side by stepping x and keeping y. Blocks CLIP when too short: give ` +
	`Metric h>=20 and RecordList/Chart h>=30.`

const layoutButtonsDoc = `JSON object turning record-toolbar buttons on or off, each with an optional ` +
	`label override: {"new":{"enabled":true},"edit":{"enabled":true,"label":"Amend"},` +
	`"submit":{"enabled":true},"delete":{"enabled":false},"clone":{"enabled":true},` +
	`"back":{"enabled":true}}. Only meaningful on a record page. A button you omit keeps its ` +
	`current setting.`

const layoutVisibilityDoc = `JSON object deciding when this layout is the one rendered: ` +
	`{"expression":"record.values.status == 'closed'","roles":["<roleID>"]}. The expression is ` +
	`evaluated against the record being viewed, so it only makes sense on a record page. Layouts ` +
	`are tried in weight order and the first match wins.`
