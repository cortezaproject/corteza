package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/mark3labs/mcp-go/mcp"
)

type pageHandler struct {
	reg toolRegistrar
}

func PageHandler(reg toolRegistrar) *pageHandler {
	h := &pageHandler{reg: reg}
	h.register()
	return h
}

func (h *pageHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_lookup",
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDescription("Look up pages in a namespace. Omit the page argument to get the full page hierarchy as a nested tree (parent/child relationships and navigation order). Provide a title, handle, or ID to get a single page with its full block configuration."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("page", mcp.Description("Page title, handle, or ID. Omit to get the namespace's page tree.")),
		),
		"Lookup page",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_create",
			mcp.WithDescription(`Create a new page in a namespace. There are two distinct page types:

1. Record list page — shows all records in a table. Do NOT set the module parameter at page level. Add a RecordList block with moduleID in its options.
2. Record detail page — the form for viewing or editing a single record. Set the module parameter at page level. Add a Record block with the fields to display. Only one record detail page can exist per module.

The layout grid is 48 columns wide (cell height 10px; default block size is w=24 h=18). A full-width block uses xywh [0,0,48,20]. Blocks CLIP their content when too short and fail silently as blank UI — give Metric blocks h>=20 and RecordList/Chart blocks h>=30. Call compose_page_block_schema with the block kind to get its options before creating blocks.`),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier")),
			mcp.WithString("description", mcp.Description("Page description")),
			mcp.WithString("parent", mcp.Description("Parent page title, handle, or ID. Omit for a root-level page.")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID. Set ONLY for record pages (the single-record form). Do not set for record list pages — put the module in the RecordList block options instead.")),
			mcp.WithBoolean("visible", mcp.Description("Show page in navigation (default: true)")),
			mcp.WithString("blocks", mcp.Description(`JSON array of page blocks. Grid is 48 columns wide — full-width block: [{"kind":"RecordList","title":"My Block","xywh":[0,0,48,20],"options":{...}}]. Call compose_page_block_schema first for kind-specific options.`)),
			mcp.WithString("icon", mcp.Description(`JSON object for nav icon: {"type":"library","src":"font-awesome://home"} or {"type":"link","src":"https://..."} or {"type":"svg","src":"<svg>..."}`)),
			mcp.WithString("config", mcp.Description(`JSON object for page configuration. Example: {"navItem":{"expanded":true}}`)),
			mcp.WithString("meta", mcp.Description(`JSON object for page meta. Example: {"allowPersonalLayouts":true}`)),
		),
		"Create page",
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_update",
			mcp.WithDescription("Update an existing page. Pass only the fields you want to change. Blocks are merged by blockID: blocks with a matching blockID overwrite existing ones, blocks without a blockID are appended, and existing blocks not in the payload are kept."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID")),
			mcp.WithString("title", mcp.Description("New title")),
			mcp.WithString("handle", mcp.Description("New handle")),
			mcp.WithString("description", mcp.Description("New description")),
			mcp.WithString("parent", mcp.Description("New parent page title, handle, or ID. Pass empty string to move to root.")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID for record detail pages. Pass empty string to clear.")),
			mcp.WithBoolean("visible", mcp.Description("Show page in navigation")),
			mcp.WithString("blocks", mcp.Description(`JSON array of page blocks. Merged by blockID — include blockID to update an existing block, omit blockID to add a new one. Grid is 48 columns wide; a full-width block uses xywh [0,0,48,20].`)),
			mcp.WithString("icon", mcp.Description(`JSON object for nav icon: {"type":"library","src":"font-awesome://home"} or {"type":"link","src":"https://..."} or {"type":"svg","src":"<svg>..."}`)),
			mcp.WithString("config", mcp.Description(`JSON object for page configuration. Replaces existing config. Example: {"navItem":{"expanded":true}}`)),
			mcp.WithString("meta", mcp.Description(`JSON object for page meta. Replaces existing meta. Example: {"allowPersonalLayouts":true}`)),
		),
		"Update page",
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_delete",
			mcp.WithDescription(`Delete a page. Choose how child pages are handled with the strategy parameter.`),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID")),
			mcp.WithString("strategy", mcp.Description(`How to handle child pages: "abort" (default — fail if children exist), "cascade" (delete children too), "rebase" (move children to grandparent), "force" (delete unconditionally)`)),
		),
		"Delete page",
		h.del,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_reorder",
			mcp.WithDescription("Reorder child pages under a parent. Provide page IDs in the desired display order."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("parent", mcp.Description("Parent page title, handle, or ID. Omit to reorder root-level pages.")),
			mcp.WithString("pageIDs", mcp.Required(), mcp.Description(`JSON array of page IDs in desired order, e.g. ["123","456","789"]`)),
		),
		"Reorder pages",
		h.reorder,
	)
	h.reg.RegisterHiddenTool(
		mcp.NewTool("compose_page_block_schema",
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDescription(`Get the options structure for a page block kind — field names and types as a zero-value skeleton (not examples from live pages). Call this before creating blocks of an unfamiliar kind. Semantics the skeleton cannot express: Metric items use metricField "count" with empty operation for record counts, or a numeric field with operation sum/avg/min/max; Chart blocks reference an existing chart resource by chartID; RecordList/Record moduleID/fields take IDs and field names from compose_module_lookup.`),
			mcp.WithString("kind", mcp.Required(), mcp.Description("Block kind: Record, RecordList, Chart, Automation, Content, Metric, Progress, Comment, Calendar, RecordOrganizer, SocialFeed")),
		),
		"Get page block schema",
		h.blockSchema,
	)
}

func (h *pageHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	pageRef, _ := args["page"].(string)
	if pageRef == "" {
		tree, err := cmpService.DefaultPage.Tree(ctx, ns.ID)
		if err != nil {
			return nil, fmt.Errorf("page tree failed: %w", err)
		}
		out, err := json.Marshal(tree)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tree: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	pg, err := findPageByAny(ctx, ns.ID, pageRef)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(pg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal page: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *pageHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	title, _ := args["title"].(string)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	pg := &cmpTypes.Page{
		NamespaceID:    ns.ID,
		Title:          title,
		Visible:        parseBoolArg(args["visible"], true),
		CreatedByAgent: a.GetAgentIDFromContext(ctx),
	}

	if v, ok := args["handle"].(string); ok && v != "" {
		pg.Handle = v
	}
	if v, ok := args["description"].(string); ok {
		pg.Description = v
	}

	if parentRef, ok := args["parent"].(string); ok && parentRef != "" {
		parent, err := findPageByAny(ctx, ns.ID, parentRef)
		if err != nil {
			return nil, fmt.Errorf("parent page lookup failed: %w", err)
		}
		pg.SelfID = parent.ID
	}

	if modRef, ok := args["module"].(string); ok && modRef != "" {
		mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, modRef)
		if err != nil {
			return nil, fmt.Errorf("module lookup failed: %w", err)
		}
		pg.ModuleID = mod.ID
	}

	if rawBlocks, ok := args["blocks"]; ok && rawBlocks != nil {
		blocks, err := parsePageBlocks(rawBlocks)
		if err != nil {
			return nil, fmt.Errorf("invalid blocks: %w", err)
		}
		if err = resolveBlockRefs(ctx, ns.ID, blocks); err != nil {
			return nil, err
		}
		pg.Blocks = blocks
	}

	if rawIcon, ok := args["icon"]; ok && rawIcon != nil {
		icon, err := parsePageIcon(rawIcon)
		if err != nil {
			return nil, fmt.Errorf("invalid icon: %w", err)
		}
		pg.Config.NavItem.Icon = icon
	}

	if rawConfig, ok := args["config"]; ok && rawConfig != nil {
		cfg, err := parsePageConfig(rawConfig)
		if err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
		// preserve icon if it was set above and config didn't carry one
		if pg.Config.NavItem.Icon != nil && cfg.NavItem.Icon == nil {
			cfg.NavItem.Icon = pg.Config.NavItem.Icon
		}
		pg.Config = cfg
	}

	if rawMeta, ok := args["meta"]; ok && rawMeta != nil {
		meta, err := parsePageMeta(rawMeta)
		if err != nil {
			return nil, fmt.Errorf("invalid meta: %w", err)
		}
		pg.Meta = meta
	}

	pg, err = cmpService.DefaultPage.Create(ctx, pg)
	if err != nil {
		return nil, fmt.Errorf("page creation failed: %w", err)
	}

	out, err := json.Marshal(pg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal page: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *pageHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	pageRef, _ := args["page"].(string)
	pg, err := findPageByAny(ctx, ns.ID, pageRef)
	if err != nil {
		return nil, err
	}

	if v, ok := args["title"].(string); ok && v != "" {
		pg.Title = v
	}
	if v, ok := args["handle"].(string); ok && v != "" {
		pg.Handle = v
	}
	if v, ok := args["description"].(string); ok {
		pg.Description = v
	}
	if v, ok := args["visible"]; ok {
		pg.Visible = parseBoolArg(v, pg.Visible)
	}

	if parentRef, ok := args["parent"].(string); ok {
		if parentRef == "" {
			pg.SelfID = 0
		} else {
			parent, err := findPageByAny(ctx, ns.ID, parentRef)
			if err != nil {
				return nil, fmt.Errorf("parent page lookup failed: %w", err)
			}
			pg.SelfID = parent.ID
		}
	}

	if modRef, ok := args["module"].(string); ok {
		if modRef == "" {
			pg.ModuleID = 0
		} else {
			mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, modRef)
			if err != nil {
				return nil, fmt.Errorf("module lookup failed: %w", err)
			}
			pg.ModuleID = mod.ID
		}
	}

	if rawBlocks, ok := args["blocks"]; ok && rawBlocks != nil {
		blocks, err := parsePageBlocks(rawBlocks)
		if err != nil {
			return nil, fmt.Errorf("invalid blocks: %w", err)
		}
		if err = resolveBlockRefs(ctx, ns.ID, blocks); err != nil {
			return nil, err
		}
		merged, err := mergePageBlocks(pg.Blocks, blocks)
		if err != nil {
			return nil, err
		}
		pg.Blocks = merged
	}

	if rawConfig, ok := args["config"]; ok && rawConfig != nil {
		cfg, err := parsePageConfig(rawConfig)
		if err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
		pg.Config = cfg
	}

	if rawIcon, ok := args["icon"]; ok && rawIcon != nil {
		icon, err := parsePageIcon(rawIcon)
		if err != nil {
			return nil, fmt.Errorf("invalid icon: %w", err)
		}
		pg.Config.NavItem.Icon = icon
	}

	if rawMeta, ok := args["meta"]; ok && rawMeta != nil {
		meta, err := parsePageMeta(rawMeta)
		if err != nil {
			return nil, fmt.Errorf("invalid meta: %w", err)
		}
		pg.Meta = meta
	}

	pg, err = cmpService.DefaultPage.Update(ctx, pg)
	if err != nil {
		return nil, fmt.Errorf("page update failed: %w", err)
	}

	out, err := json.Marshal(pg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal page: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *pageHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	pageRef, _ := args["page"].(string)
	pg, err := findPageByAny(ctx, ns.ID, pageRef)
	if err != nil {
		return nil, err
	}

	strategyStr, _ := args["strategy"].(string)
	strategy := cmpTypes.PageChildrenDeleteStrategy(strategyStr)
	if strategy == "" {
		strategy = cmpTypes.PageChildrenOnDeleteAbort
	}

	if err = cmpService.DefaultPage.DeleteByID(ctx, ns.ID, pg.ID, strategy); err != nil {
		return nil, fmt.Errorf("page delete failed: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("page %d deleted", pg.ID)), nil
}

func (h *pageHandler) reorder(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	var parentID uint64
	if parentRef, ok := args["parent"].(string); ok && parentRef != "" {
		parent, err := findPageByAny(ctx, ns.ID, parentRef)
		if err != nil {
			return nil, fmt.Errorf("parent page lookup failed: %w", err)
		}
		parentID = parent.ID
	}

	rawIDs, ok := args["pageIDs"]
	if !ok || rawIDs == nil {
		return nil, fmt.Errorf("pageIDs is required")
	}
	idStrs, err := parseStringArray(rawIDs)
	if err != nil {
		return nil, fmt.Errorf("invalid pageIDs: %w", err)
	}
	pageIDs := make([]uint64, len(idStrs))
	for i, s := range idStrs {
		pageIDs[i], err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid page ID %q: %w", s, err)
		}
	}

	if err = cmpService.DefaultPage.Reorder(ctx, ns.ID, parentID, pageIDs); err != nil {
		return nil, fmt.Errorf("page reorder failed: %w", err)
	}

	return mcp.NewToolResultText("pages reordered"), nil
}

func (h *pageHandler) blockSchema(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	kind, _ := args["kind"].(string)
	if kind == "" {
		return nil, fmt.Errorf("kind is required")
	}

	schema, exists := cmpTypes.PageBlockOptionSchemas[kind]
	if !exists {
		return nil, fmt.Errorf("unknown block kind %q; supported: %s", kind, supportedBlockKinds())
	}

	out, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func supportedBlockKinds() string {
	kinds := make([]string, 0, len(cmpTypes.PageBlockOptionSchemas))
	for k := range cmpTypes.PageBlockOptionSchemas {
		kinds = append(kinds, k)
	}
	return fmt.Sprintf("%v", kinds)
}

// findPageByAny resolves a page reference (numeric ID, handle, or title) within a namespace.
// Title path uses Query because the auto-generated filter ignores the Title field.
func findPageByAny(ctx context.Context, namespaceID uint64, ref string) (*cmpTypes.Page, error) {
	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		return cmpService.DefaultPage.FindByID(ctx, namespaceID, id)
	}
	if p, err := cmpService.DefaultPage.FindByHandle(ctx, namespaceID, ref); err == nil {
		return p, nil
	}
	set, _, err := cmpService.DefaultPage.Search(ctx, cmpTypes.PageFilter{
		NamespaceID: namespaceID,
		Query:       ref,
	})
	if err != nil {
		return nil, fmt.Errorf("page lookup failed: %w", err)
	}
	for _, p := range set {
		if strings.EqualFold(p.Title, ref) || strings.EqualFold(p.Handle, ref) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("page %q not found", ref)
}

func mergePageBlocks(existing, incoming cmpTypes.PageBlocks) (cmpTypes.PageBlocks, error) {
	if len(incoming) == 0 {
		return existing, nil
	}
	indexByID := make(map[uint64]int, len(existing))
	for i, b := range existing {
		if b.BlockID != 0 {
			indexByID[b.BlockID] = i
		}
	}
	out := make(cmpTypes.PageBlocks, len(existing))
	copy(out, existing)
	for _, b := range incoming {
		if b.BlockID == 0 {
			out = append(out, b)
			continue
		}
		idx, ok := indexByID[b.BlockID]
		if !ok {
			return nil, fmt.Errorf("unknown block ID %d", b.BlockID)
		}
		out[idx] = b
	}
	return out, nil
}

// resolveBlockRefs translates human-readable module/chart references in block
// options to their uint64 ID strings, which the frontend SDK requires.
// Handles: RecordList.module, Chart.chart, and any moduleID field that looks
// like a name rather than an already-numeric ID.
func resolveBlockRefs(ctx context.Context, nsID uint64, blocks cmpTypes.PageBlocks) error {
	for i := range blocks {
		b := &blocks[i]
		if b.Options == nil {
			continue
		}

		// module → module ID (RecordList)
		if ref, _ := b.Options["module"].(string); ref != "" {
			if _, err := strconv.ParseUint(ref, 10, 64); err != nil {
				mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, ref)
				if err != nil {
					return fmt.Errorf("block %q: module %q not found: %w", b.Title, ref, err)
				}
				b.Options["module"] = strconv.FormatUint(mod.ID, 10)
			}
		}

		// moduleID → module ID (Metric, Progress, Comment, RecordOrganizer, etc.)
		if ref, _ := b.Options["moduleID"].(string); ref != "" {
			if _, err := strconv.ParseUint(ref, 10, 64); err != nil {
				mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, ref)
				if err != nil {
					return fmt.Errorf("block %q: moduleID %q not found: %w", b.Title, ref, err)
				}
				b.Options["moduleID"] = strconv.FormatUint(mod.ID, 10)
			}
		}

		// chart → chart ID (Chart block)
		if ref, _ := b.Options["chart"].(string); ref != "" {
			if _, err := strconv.ParseUint(ref, 10, 64); err != nil {
				ch, err := findChartByAny(ctx, nsID, ref)
				if err != nil {
					return fmt.Errorf("block %q: chart %q not found: %w", b.Title, ref, err)
				}
				b.Options["chart"] = strconv.FormatUint(ch.ID, 10)
			}
		}

		// Calendar feeds[].moduleID and Metric metrics[].moduleID
		for _, arrayKey := range []string{"feeds", "metrics"} {
			items, _ := b.Options[arrayKey].([]interface{})
			for _, item := range items {
				m, _ := item.(map[string]interface{})
				if m == nil {
					continue
				}
				if ref, _ := m["moduleID"].(string); ref != "" {
					if _, err := strconv.ParseUint(ref, 10, 64); err != nil {
						mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, ref)
						if err != nil {
							return fmt.Errorf("block %q %s item: moduleID %q not found: %w", b.Title, arrayKey, ref, err)
						}
						m["moduleID"] = strconv.FormatUint(mod.ID, 10)
					}
				}
			}
		}
	}
	return nil
}

func findChartByAny(ctx context.Context, nsID uint64, ref string) (*cmpTypes.Chart, error) {
	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		return cmpService.DefaultChart.FindByID(ctx, nsID, id)
	}
	if c, err := cmpService.DefaultChart.FindByHandle(ctx, nsID, ref); err == nil {
		return c, nil
	}
	set, _, err := cmpService.DefaultChart.Find(ctx, cmpTypes.ChartFilter{
		NamespaceID: nsID,
		Query:       ref,
	})
	if err != nil {
		return nil, fmt.Errorf("chart lookup failed: %w", err)
	}
	for _, c := range set {
		if strings.EqualFold(c.Name, ref) || strings.EqualFold(c.Handle, ref) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("chart %q not found", ref)
}

// parsePageBlocks unmarshals a JSON string or array into PageBlocks and
// normalizes layout: width defaults to 12 (full grid) and blocks are stacked
// vertically in the order given, so agents only need to supply kind + options.
func parsePageBlocks(raw interface{}) (cmpTypes.PageBlocks, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return nil, fmt.Errorf("cannot encode blocks: %w", err)
		}
	}
	var blocks cmpTypes.PageBlocks
	if err := json.Unmarshal(data, &blocks); err != nil {
		return nil, fmt.Errorf("blocks must be a JSON array: %w", err)
	}
	return autoLayoutBlocks(blocks), nil
}

// autoLayoutBlocks stacks blocks top-to-bottom at full grid width (12 columns).
// This corrects the common agent mistake of passing wrong or missing xywh values.
// Explicit heights are preserved; missing heights get a per-kind default.
func autoLayoutBlocks(blocks cmpTypes.PageBlocks) cmpTypes.PageBlocks {
	y := 0
	for i := range blocks {
		b := &blocks[i]
		b.XYWH[0] = 0  // x: left edge
		b.XYWH[1] = y  // y: stacked below previous block
		b.XYWH[2] = 12 // w: full grid width
		if b.XYWH[3] == 0 {
			b.XYWH[3] = defaultBlockHeight(b.Kind)
		}
		y += b.XYWH[3]
	}
	return blocks
}

func defaultBlockHeight(kind string) int {
	switch kind {
	case "RecordList", "Record", "Calendar", "RecordOrganizer", "ChatbotInbox":
		return 20
	case "Chart":
		return 12
	case "Metric", "Progress":
		return 6
	case "Content", "SocialFeed", "Comment":
		return 10
	case "Automation":
		return 8
	default:
		return 15
	}
}

// parsePageIcon unmarshals a JSON string or object into PageConfigIcon.
func parsePageIcon(raw interface{}) (*cmpTypes.PageConfigIcon, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return nil, fmt.Errorf("cannot encode icon: %w", err)
		}
	}
	var icon cmpTypes.PageConfigIcon
	if err := json.Unmarshal(data, &icon); err != nil {
		return nil, fmt.Errorf("icon must be a JSON object: %w", err)
	}
	return &icon, nil
}

// parsePageConfig unmarshals a JSON string or object into PageConfig.
func parsePageConfig(raw interface{}) (cmpTypes.PageConfig, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return cmpTypes.PageConfig{}, fmt.Errorf("cannot encode config: %w", err)
		}
	}
	var cfg cmpTypes.PageConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cmpTypes.PageConfig{}, fmt.Errorf("config must be a JSON object: %w", err)
	}
	return cfg, nil
}

// parsePageMeta unmarshals a JSON string or object into PageMeta.
func parsePageMeta(raw interface{}) (cmpTypes.PageMeta, error) {
	var data []byte
	switch v := raw.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return cmpTypes.PageMeta{}, fmt.Errorf("cannot encode meta: %w", err)
		}
	}
	var meta cmpTypes.PageMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return cmpTypes.PageMeta{}, fmt.Errorf("meta must be a JSON object: %w", err)
	}
	return meta, nil
}
