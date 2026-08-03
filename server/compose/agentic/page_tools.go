package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in page_handler.go, in this order.

// Layout facts shared by the create and update descriptions. Verified against the
// gridstack setup the webapp page builder uses
// (client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue: 48
// columns, 10px cell height, default block 24x18).
const pageGridDoc = `The layout grid is 48 columns wide (cell height 10px; default block size is w=24 h=18). A full-width block uses xywh [0,0,48,20]. Blocks CLIP their content when too short and fail silently as blank UI — give Metric blocks h>=20 and RecordList/Chart blocks h>=30.`

func (h *pageHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_lookup",
			mcp.WithDescription(
				"Look up pages in a namespace, or fetch one page in full. "+
					"Provide 'page' (title, handle, or ID) to get that single page with its complete block "+
					"layout — this is the only mode that returns blocks, and it is what you call before "+
					"changing them, because compose_page_update needs the existing blockIDs. "+
					"Omit 'page' to list the namespace's pages as "+
					"{pageID, title, handle, parentID, moduleID, visible, weight}, ordered by navigation "+
					"weight and paged with 'limit' and 'pageCursor'. parentID is the page's parent, 0 for a "+
					"root-level page, so the hierarchy is reconstructible from a listing. "+
					"Set 'tree' to true to get that same slim view nested as a hierarchy instead. A tree is "+
					"always returned whole and ignores 'limit' and 'pageCursor': a hierarchy cannot be cut "+
					"into pages without severing subtrees from their parents, so paging it would be "+
					"meaningless. Use the flat, paged listing for a namespace with many pages. "+
					"Neither listing carries blocks, config or meta — fetch the single page for those.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Description("Page title, handle, or ID (as string to prevent precision loss). Omit to list the namespace's pages instead.")),
			mcp.WithBoolean("tree", mcp.Description("List the pages nested as a hierarchy instead of a flat list. Ignored when 'page' is given, and not paged.")),
			mcp.WithString("limit", mcp.Description("Maximum pages to return when listing, default 50, capped at 200. Ignored in tree mode.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page of the listing. Ignored in tree mode.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup page",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_create",
			mcp.WithDescription(`Create a new page in a namespace. A page is a screen in the namespace's navigation; it holds blocks that render records, charts and content. There are two distinct page types:

1. Record list page — shows all records in a table. Do NOT set the module parameter at page level. Add a RecordList block with moduleID in its options.
2. Record detail page — the form for viewing or editing a single record. Set the module parameter at page level. Add a Record block with the fields to display. Only one record detail page can exist per module, and creating a second one for the same module is rejected.

`+pageGridDoc+`

Call compose_page_block_schema with the block kind to get its options before creating blocks. Blocks are numbered on create — the returned page carries the assigned blockIDs, which compose_page_update needs to change a block later.`),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier (lowercase letters, digits, and underscores; at least 2 characters)")),
			mcp.WithString("description", mcp.Description("Page description")),
			mcp.WithString("parent", mcp.Description("Parent page title, handle, or ID (as string to prevent precision loss). Omit for a root-level page.")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID (as string to prevent precision loss). Set ONLY for record detail pages (the single-record form). Do not set for record list pages — put the module in the RecordList block options instead.")),
			mcp.WithBoolean("visible", mcp.Description("Show page in navigation (default: true)")),
			mcp.WithString("blocks", mcp.Description(`JSON array of page blocks. Grid is 48 columns wide — full-width block: [{"kind":"RecordList","title":"My Block","xywh":[0,0,48,20],"options":{...}}]. Call compose_page_block_schema first for kind-specific options.`)),
			mcp.WithString("icon", mcp.Description(`JSON object for nav icon: {"type":"library","src":"font-awesome://home"} or {"type":"link","src":"https://..."} or {"type":"svg","src":"<svg>..."}`)),
			mcp.WithString("config", mcp.Description(`JSON object for page configuration. Example: {"navItem":{"expanded":true}}`)),
			mcp.WithString("meta", mcp.Description(`JSON object for page meta. Example: {"allowPersonalLayouts":true}`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create page",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_update",
			mcp.WithDescription(
				"Update an existing page. An argument you omit is left unchanged. Sending an empty string "+
					"clears the value: 'parent' moves the page to the root, 'module' unlinks the module, "+
					"'handle' and 'description' are emptied. 'title' is the exception — a page without one is "+
					"unusable, so an empty title is ignored rather than applied. "+
					"'blocks' is the documented exception to replacing wholesale: blocks are MERGED by "+
					"blockID. A block carrying a blockID overwrites the existing block with that ID, a block "+
					"without one is appended and is assigned a fresh blockID, and existing blocks you do not "+
					"mention are kept — so send only the blocks you are adding or changing. An unknown "+
					"blockID is rejected. Because the merge never removes, this tool cannot take a block away: "+
					"use compose_page_remove_blocks for that. "+
					"'config' and 'meta' each replace their whole object; 'icon' is applied after 'config', "+
					"so passing both leaves 'icon' in charge of the navigation icon. "+
					"Call compose_page_lookup with 'page' first when you need the current blocks or their "+
					"blockIDs.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("title", mcp.Description("New title. An empty string is ignored — a page cannot be left without a title.")),
			mcp.WithString("handle", mcp.Description("New handle. Pass an empty string to clear it.")),
			mcp.WithString("description", mcp.Description("New description. Pass an empty string to clear it.")),
			mcp.WithString("parent", mcp.Description("New parent page title, handle, or ID (as string to prevent precision loss). Pass an empty string to move the page to the root.")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID (as string to prevent precision loss) for record detail pages. Pass an empty string to unlink the module.")),
			mcp.WithBoolean("visible", mcp.Description("Show page in navigation")),
			mcp.WithString("blocks", mcp.Description(`JSON array of page blocks. Merged by blockID — include blockID to update an existing block, omit blockID to add a new one. Blocks are never removed here; omitting one keeps it, and compose_page_remove_blocks is what deletes one. Grid is 48 columns wide; a full-width block uses xywh [0,0,48,20].`)),
			mcp.WithString("icon", mcp.Description(`JSON object for nav icon: {"type":"library","src":"font-awesome://home"} or {"type":"link","src":"https://..."} or {"type":"svg","src":"<svg>..."}`)),
			mcp.WithString("config", mcp.Description(`JSON object for page configuration. Replaces existing config. Example: {"navItem":{"expanded":true}}`)),
			mcp.WithString("meta", mcp.Description(`JSON object for page meta. Replaces existing meta. Example: {"allowPersonalLayouts":true}`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update page",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_delete",
			mcp.WithDescription(
				"Delete a page by title, handle, or ID. The delete is soft: nothing is erased — the page "+
					"stops appearing in lookups and in navigation, but it is retained. Records are not "+
					"touched at all, because a page only displays data that lives in modules. "+
					"Child pages are handled by 'strategy', which defaults to refusing the delete when the "+
					"page has children, so decide explicitly what should happen to them. "+
					"compose_page_undelete reverses this one page at a time, but note the page ID first: a "+
					"deleted page can no longer be found by title or handle.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("strategy", mcp.Description(`How to handle child pages: "abort" (default — fail if the page has any undeleted children), "cascade" (delete the whole subtree), "rebase" (move the children up to this page's parent, then delete it), "force" (delete only this page and leave the children pointing at it — they then show up as root-level pages)`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete page",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted page, reversing compose_page_delete. A delete only marks the page — "+
					"its blocks, config, meta and navigation weight were all retained — so the page returns to "+
					"navigation exactly as it was and no block needs rebuilding. "+
					"Requires the numeric pageID: deleted pages are excluded from every lookup path, so "+
					"compose_page_lookup can no longer resolve one by title or handle, its tree mode does not "+
					"show it, and it exposes no includeDeleted-style filter. Use the ID compose_page_delete "+
					"reported, or one from a compose_page_lookup taken before the delete. "+
					"This restores one page and never its subtree: a page deleted with the 'cascade' strategy "+
					"needs one call per child, and a child restored while its parent is still deleted shows up "+
					"as a root-level page until the parent is restored too. "+
					"Calling this on a page that is not deleted is accepted and changes nothing.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("pageID", mcp.Required(), mcp.Description("ID of the deleted page (as string to prevent precision loss). A title or handle will not work — deleted pages are not resolvable by either.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete page",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_remove_blocks",
			mcp.WithDescription(
				"Remove one or more blocks from a page, by blockID. This is the only way to take a block off a "+
					"page: compose_page_update MERGES the blocks it is given by blockID, so it can add a block "+
					"or overwrite one but never drop one. Use that tool to change a block and this one to "+
					"delete it. "+
					"This edits a page's layout, it does not delete the page — that is compose_page_delete — "+
					"and it touches nothing a block pointed at: the module, records or chart a block rendered "+
					"are left exactly as they were, because a block is only a view onto them. "+
					"Call compose_page_lookup with 'page' first to read the current blocks and their blockIDs; "+
					"blockIDs are assigned per page as blocks are created, so they mean nothing on another "+
					"page. Every ID you pass must exist on this page — an unknown blockID is rejected and "+
					"nothing at all is removed, so a stale layout fails loudly instead of half applying. Blocks "+
					"you do not list keep their position, which means removing one leaves a gap in the grid "+
					"rather than reflowing the rest; reposition the survivors with compose_page_update if the "+
					"layout should close up. The updated page is returned, so the remaining blocks and their "+
					"blockIDs come back in the response.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID (as string to prevent precision loss)")),
			mcp.WithString("blockIDs", mcp.Required(), mcp.Description(`JSON array of blockIDs to remove, as strings to prevent precision loss, e.g. ["1","3"]. Every ID must exist on the page; read them from compose_page_lookup with 'page'.`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Remove page blocks",
		h.removeBlocks,
	)

	h.reg.RegisterTool(
		mcp.NewTool("compose_page_reorder",
			mcp.WithDescription(
				"Set the navigation order of the pages under one parent. The pages you list are weighted in "+
					"the order given, so pass the COMPLETE ordered list of that parent's children: any child "+
					"you leave out is pushed behind the ones you list, in an unspecified order. "+
					"Omit 'parent' to order the root-level pages — note that this mode re-weights every "+
					"other page in the namespace as well, so only use it with the full list of root pages. "+
					"Call compose_page_lookup first to get the pageIDs and the current order; this tool "+
					"changes ordering only, never the parent of a page — use compose_page_update's 'parent' "+
					"to move a page.",
			),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID (as string to prevent precision loss)")),
			mcp.WithString("parent", mcp.Description("Parent page title, handle, or ID (as string to prevent precision loss). Omit to reorder root-level pages.")),
			mcp.WithString("pageIDs", mcp.Required(), mcp.Description(`JSON array of page IDs in the desired order, as strings to prevent precision loss, e.g. ["123","456","789"]`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Reorder pages",
		h.reorder,
	)

	// Grandfathered four-segment name: the grammar is {app}_{resource}_{op}, but
	// this one shipped and agent definitions store tool names by value.
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_block_schema",
			mcp.WithDescription(`Get the options structure for a page block kind — field names and types as a zero-value skeleton (not examples from live pages). Call this before creating blocks of an unfamiliar kind, and pass the result's field names into the block's "options" object in compose_page_create or compose_page_update. Semantics the skeleton cannot express: Metric items use metricField "count" with empty operation for record counts, or a numeric field with operation sum/avg/min/max; Chart blocks reference an existing chart resource by chartID, created with compose_chart_create; RecordList/Record moduleID/fields take IDs and field names from compose_module_lookup. This reads a static schema — it touches no namespace and no data.`),
			mcp.WithString("kind", mcp.Required(), mcp.Description("Block kind: Record, RecordList, Chart, Automation, Content, Metric, Progress, Comment, Calendar, RecordOrganizer, SocialFeed, ChatbotInbox")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Get page block schema",
		h.blockSchema,
	)
}
