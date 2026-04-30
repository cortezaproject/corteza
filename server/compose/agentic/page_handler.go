package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
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
			mcp.WithDescription("Look up pages in a namespace. Omit the page argument to list all pages. Provide a title, handle, or ID to get a single page with its full block configuration."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("page", mcp.Description("Page title, handle, or ID. Omit to list all pages.")),
		),
		"Lookup page",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_tree",
			mcp.WithDescription("Get the full page hierarchy for a namespace as a nested tree showing parent/child relationships and navigation order."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
		),
		"Get page tree",
		h.tree,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_create",
			mcp.WithDescription(`Create a new page in a namespace. There are two distinct page types:

1. Record list page — shows all records in a table. Do NOT set the module parameter at page level. Add a RecordList block with moduleID in its options.
2. Record detail page — the form for viewing or editing a single record. Set the module parameter at page level. Add a Record block with the fields to display. Only one record detail page can exist per module.

The layout grid is 12 columns wide. A full-width block uses xywh [0,0,12,20]. Call compose_page_block_schema with the block kind to get its options before creating blocks.`),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier")),
			mcp.WithString("description", mcp.Description("Page description")),
			mcp.WithString("parent", mcp.Description("Parent page title, handle, or ID. Omit for a root-level page.")),
			mcp.WithString("module", mcp.Description("Module name, handle, or ID. Set ONLY for record pages (the single-record form). Do not set for record list pages — put the module in the RecordList block options instead.")),
			mcp.WithBoolean("visible", mcp.Description("Show page in navigation (default: true)")),
			mcp.WithString("blocks", mcp.Description(`JSON array of page blocks. Grid is 12 columns wide — full-width block: [{"kind":"RecordList","title":"My Block","xywh":[0,0,12,20],"options":{...}}]. Call compose_page_block_schema first for kind-specific options.`)),
		),
		"Create page",
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_update",
			mcp.WithDescription("Update an existing page. When blocks are provided they replace all existing blocks — call compose_page_lookup first if you need to preserve existing blocks."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("page", mcp.Required(), mcp.Description("Page title, handle, or ID")),
			mcp.WithString("title", mcp.Description("New title")),
			mcp.WithString("handle", mcp.Description("New handle")),
			mcp.WithString("description", mcp.Description("New description")),
			mcp.WithString("parent", mcp.Description("New parent page title, handle, or ID. Pass empty string to move to root.")),
			mcp.WithBoolean("visible", mcp.Description("Show page in navigation")),
			mcp.WithString("blocks", mcp.Description(`JSON array of page blocks. Replaces all existing blocks. Grid is 12 columns wide — full-width block uses xywh [0,0,12,20].`)),
			mcp.WithString("icon", mcp.Description(`JSON object for nav icon: {"type":"library","src":"font-awesome://home"} or {"type":"link","src":"https://..."} or {"type":"inline-svg","src":"<svg>..."}`)),
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
	h.reg.RegisterTool(
		mcp.NewTool("compose_page_block_schema",
			mcp.WithDescription("Get real examples of a page block kind from existing pages. Returns actual block configurations from the system so you know what options to use. Call this before creating blocks of an unfamiliar kind."),
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
		set, _, err := cmpService.DefaultPage.Find(ctx, cmpTypes.PageFilter{NamespaceID: ns.ID})
		if err != nil {
			return nil, fmt.Errorf("page list failed: %w", err)
		}
		type pageItem struct {
			ID       uint64 `json:"pageID,string"`
			SelfID   uint64 `json:"selfID,string"`
			ModuleID uint64 `json:"moduleID,string,omitempty"`
			Title    string `json:"title"`
			Handle   string `json:"handle,omitempty"`
			Visible  bool   `json:"visible"`
			Weight   int    `json:"weight"`
		}
		items := make([]pageItem, len(set))
		for i, p := range set {
			items[i] = pageItem{
				ID: p.ID, SelfID: p.SelfID, ModuleID: p.ModuleID,
				Title: p.Title, Handle: p.Handle,
				Visible: p.Visible, Weight: p.Weight,
			}
		}
		out, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal pages: %w", err)
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

func (h *pageHandler) tree(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

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
		NamespaceID: ns.ID,
		Title:       title,
		Visible:     parseBoolArg(args["visible"], true),
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
		pg.Blocks = blocks
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

	if rawBlocks, ok := args["blocks"]; ok && rawBlocks != nil {
		blocks, err := parsePageBlocks(rawBlocks)
		if err != nil {
			return nil, fmt.Errorf("invalid blocks: %w", err)
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
func findPageByAny(ctx context.Context, namespaceID uint64, ref string) (*cmpTypes.Page, error) {
	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		return cmpService.DefaultPage.FindByID(ctx, namespaceID, id)
	}
	if p, err := cmpService.DefaultPage.FindByHandle(ctx, namespaceID, ref); err == nil {
		return p, nil
	}
	// Fall back to title search
	set, _, err := cmpService.DefaultPage.Find(ctx, cmpTypes.PageFilter{
		NamespaceID: namespaceID,
		Title:       ref,
	})
	if err != nil {
		return nil, fmt.Errorf("page lookup failed: %w", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("page %q not found", ref)
	}
	return set[0], nil
}

// parsePageBlocks unmarshals a JSON string or array into PageBlocks.
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
	return blocks, nil
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
