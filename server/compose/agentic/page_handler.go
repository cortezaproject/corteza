package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type (
	pageHandler struct {
		reg toolRegistrar
	}

	// pageItem is the slim list projection. A page's blocks are its heavy part —
	// and they are incidental when the caller is picking one page out of many —
	// so a listing drops them along with config, meta and description. parentID
	// and weight stay: they are what the hierarchy and the navigation order are
	// made of, and dropping them would make a flat listing unreadable as a tree.
	// A single-page lookup returns the full service type.
	pageItem struct {
		ID       uint64 `json:"pageID,string"`
		Title    string `json:"title"`
		Handle   string `json:"handle,omitempty"`
		ParentID uint64 `json:"parentID,string"`
		ModuleID uint64 `json:"moduleID,string,omitempty"`
		Visible  bool   `json:"visible"`
		Weight   int    `json:"weight"`
	}

	// pageNode is pageItem nested, for the tree listing.
	pageNode struct {
		pageItem
		Children []pageNode `json:"children,omitempty"`
	}
)

// Declarations for these handlers are in page_tools.go, in the same order.
func PageHandler(reg toolRegistrar) *pageHandler {
	h := &pageHandler{reg: reg}
	h.register()
	return h
}

// resolveNs resolves the namespace ref every page tool takes.
func (h *pageHandler) resolveNs(ctx context.Context, args map[string]any) (uint64, error) {
	nsRef, err := toolkit.ReqRef(args, "namespace")
	if err != nil {
		return 0, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		return 0, toolkit.Errf("namespace lookup", err)
	}

	return ns.ID, nil
}

// resolvePage resolves the namespace and the required page ref that the
// single-page tools all take. The lookup tool, which may omit the page ref,
// resolves the namespace on its own.
func (h *pageHandler) resolvePage(ctx context.Context, args map[string]any) (uint64, *cmpTypes.Page, error) {
	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return 0, nil, err
	}

	pageRef, err := toolkit.ReqRef(args, "page")
	if err != nil {
		return 0, nil, err
	}

	pg, err := findPageByAny(ctx, nsID, pageRef)
	if err != nil {
		return 0, nil, err
	}

	return nsID, pg, nil
}

func (h *pageHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	pageRef, err := toolkit.Ref(args, "page")
	if err != nil {
		return nil, err
	}

	// Single-page mode returns the raw service type, blocks included: once the
	// caller has picked a page, the blocks are what they came for.
	if pageRef != "" {
		pg, err := findPageByAny(ctx, nsID, pageRef)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(pg)
	}

	// Tree mode is deliberately unpaged. A cursor slices an ordered row set, and
	// a hierarchy has no such ordering that survives being cut: any slice would
	// hand back child nodes whose parents are on another page, or drop whole
	// subtrees. So the tree is returned whole and limit/pageCursor are ignored
	// here; the flat listing below is the paged path. The slim projection keeps
	// the response small either way.
	if toolkit.Bool(args, "tree") {
		tree, err := cmpService.DefaultPage.Tree(ctx, nsID)
		if err != nil {
			return nil, toolkit.Errf("page tree", err)
		}
		return toolkit.JSONResult(map[string]any{"pages": pageNodes(tree)})
	}

	f := cmpTypes.PageFilter{NamespaceID: nsID}

	// Weight is the navigation order, and it is a sortable column, so the listing
	// comes back in the order the pages are displayed in. The store appends the
	// primary key to the sort itself, which keeps the cursor stable.
	if err = f.Sort.Set("weight ASC"); err != nil {
		return nil, toolkit.Errf("page list", err)
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := cmpService.DefaultPage.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("page list", err)
	}

	items := make([]pageItem, len(set))
	for i, p := range set {
		items[i] = newPageItem(p)
	}

	// The cursor goes in as the value, not as a string: marshalling a
	// PagingCursor emits the encoded form that pageCursor accepts back, while its
	// String method is a debug rendering that does not round-trip. A nil cursor
	// marshals to null, so the last page needs no special case.
	return toolkit.JSONResult(map[string]any{
		"pages":          items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *pageHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	title, err := toolkit.ReqStr(args, "title")
	if err != nil {
		return nil, err
	}

	pg := &cmpTypes.Page{
		NamespaceID:    nsID,
		Title:          title,
		Handle:         toolkit.Str(args, "handle"),
		Description:    toolkit.Str(args, "description"),
		Visible:        true,
		CreatedByAgent: a.GetAgentIDFromContext(ctx),
	}

	// Pages are visible in navigation unless the caller says otherwise, so the
	// flag defaults to true and only an argument that is actually present can
	// turn it off.
	if v, ok := args["visible"]; ok && v != nil {
		pg.Visible = toolkit.Bool(args, "visible")
	}

	parentRef, err := toolkit.Ref(args, "parent")
	if err != nil {
		return nil, err
	}
	if parentRef != "" {
		parent, err := findPageByAny(ctx, nsID, parentRef)
		if err != nil {
			return nil, toolkit.Errf("parent page lookup", err)
		}
		pg.SelfID = parent.ID
	}

	modRef, err := toolkit.Ref(args, "module")
	if err != nil {
		return nil, err
	}
	if modRef != "" {
		mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, modRef)
		if err != nil {
			return nil, toolkit.Errf("module lookup", err)
		}
		pg.ModuleID = mod.ID
	}

	if rawBlocks, ok := args["blocks"]; ok && rawBlocks != nil {
		if pg.Blocks, err = parsePageBlocks(rawBlocks); err != nil {
			return nil, fmt.Errorf("invalid blocks: %w", err)
		}

		if err = resolveBlockRefs(ctx, nsID, pg.Blocks); err != nil {
			return nil, err
		}

		pg.Blocks = autoLayoutBlocks(pg.Blocks)
	}

	if rawConfig, ok := args["config"]; ok && rawConfig != nil {
		if pg.Config, err = parsePageConfig(rawConfig); err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
	}

	if rawMeta, ok := args["meta"]; ok && rawMeta != nil {
		if pg.Meta, err = parsePageMeta(rawMeta); err != nil {
			return nil, fmt.Errorf("invalid meta: %w", err)
		}
	}

	pg, err = cmpService.DefaultPage.Create(ctx, pg)
	if err != nil {
		return nil, toolkit.Errf("page creation", err)
	}

	return toolkit.JSONResult(pg)
}

func (h *pageHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	// Absent leaves the value alone. Title is not clearable — a page without one
	// is unusable — so an empty string is ignored rather than treated as a clear.
	if v := toolkit.Str(args, "title"); v != "" {
		pg.Title = v
	}

	// Handle and description are optional on a page, so here present-and-empty
	// does clear.
	if v, ok := args["handle"].(string); ok {
		pg.Handle = v
	}
	if v, ok := args["description"].(string); ok {
		pg.Description = v
	}

	if v, ok := args["visible"]; ok && v != nil {
		pg.Visible = toolkit.Bool(args, "visible")
	}

	if raw, ok := args["parent"]; ok && raw != nil {
		parentRef, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("parent must be a string to avoid precision loss, got %T", raw)
		}

		if parentRef == "" {
			pg.SelfID = 0
		} else {
			parent, err := findPageByAny(ctx, nsID, parentRef)
			if err != nil {
				return nil, toolkit.Errf("parent page lookup", err)
			}
			pg.SelfID = parent.ID
		}
	}

	if raw, ok := args["module"]; ok && raw != nil {
		modRef, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("module must be a string to avoid precision loss, got %T", raw)
		}

		if modRef == "" {
			pg.ModuleID = 0
		} else {
			mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, modRef)
			if err != nil {
				return nil, toolkit.Errf("module lookup", err)
			}
			pg.ModuleID = mod.ID
		}
	}

	// Blocks are the documented merge exception: incoming blocks are matched by
	// blockID against the existing set instead of replacing it, so a caller can
	// add or amend one block without resending the page's whole layout.
	if rawBlocks, ok := args["blocks"]; ok && rawBlocks != nil {
		blocks, err := parsePageBlocks(rawBlocks)
		if err != nil {
			return nil, fmt.Errorf("invalid blocks: %w", err)
		}
		if err = resolveBlockRefs(ctx, nsID, blocks); err != nil {
			return nil, err
		}

		// A block sent to change its options only, with no xywh, keeps where it
		// already sits — moving it would be a side effect of an edit that never
		// mentioned layout. Blocks genuinely new to the page have nothing to
		// inherit and are placed by autoLayoutBlocks below the existing ones.
		inheritBlockLayout(pg.Blocks, blocks)

		if pg.Blocks, err = mergePageBlocks(pg.Blocks, blocks); err != nil {
			return nil, err
		}

		pg.Blocks = autoLayoutBlocks(pg.Blocks)
	}

	// Config and meta replace wholesale: present-and-empty clears them.
	if rawConfig, ok := args["config"]; ok && rawConfig != nil {
		if pg.Config, err = parsePageConfig(rawConfig); err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
	}

	if rawMeta, ok := args["meta"]; ok && rawMeta != nil {
		if pg.Meta, err = parsePageMeta(rawMeta); err != nil {
			return nil, fmt.Errorf("invalid meta: %w", err)
		}
	}

	pg, err = cmpService.DefaultPage.Update(ctx, pg)
	if err != nil {
		return nil, toolkit.Errf("page update", err)
	}

	return toolkit.JSONResult(pg)
}

func (h *pageHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	// Rejected here rather than at the service, so an unusable strategy names the
	// alternatives instead of returning an opaque error.
	strategy := cmpTypes.PageChildrenDeleteStrategy(toolkit.Str(args, "strategy"))
	switch strategy {
	case "":
		strategy = cmpTypes.PageChildrenOnDeleteAbort
	case cmpTypes.PageChildrenOnDeleteAbort,
		cmpTypes.PageChildrenOnDeleteCascade,
		cmpTypes.PageChildrenOnDeleteRebase,
		cmpTypes.PageChildrenOnDeleteForce:
	default:
		return nil, fmt.Errorf("unknown strategy %q: use abort, cascade, rebase or force", strategy)
	}

	if err = cmpService.DefaultPage.DeleteByID(ctx, nsID, pg.ID, strategy); err != nil {
		return nil, toolkit.Errf("page delete", err)
	}

	return toolkit.TextResult("page %d deleted", pg.ID), nil
}

func (h *pageHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	// An ID rather than the title-or-handle ref the other page tools take:
	// findPageByAny resolves through lookup and search, and both treat a page
	// carrying a deleted marker as not found, so a deleted page has no name left
	// to be resolved by. The service loads it by ID, outside those paths.
	pageID, err := toolkit.ReqID(args, "pageID")
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultPage.UndeleteByID(ctx, nsID, pageID); err != nil {
		return nil, toolkit.Errf("page undelete", err)
	}

	return toolkit.TextResult("page %d restored", pageID), nil
}

func (h *pageHandler) removeBlocks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	_, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	rawIDs, ok := args["blockIDs"]
	if !ok || rawIDs == nil {
		return nil, fmt.Errorf("blockIDs is required")
	}

	idStrs, err := parseStringArray(rawIDs)
	if err != nil {
		return nil, fmt.Errorf("invalid blockIDs: must be a JSON array of ID strings: %w", err)
	}
	if len(idStrs) == 0 {
		return nil, fmt.Errorf("blockIDs must list at least one block")
	}

	present := make(map[uint64]bool, len(pg.Blocks))
	for _, b := range pg.Blocks {
		present[b.BlockID] = true
	}

	// An unknown blockID fails the whole call rather than being skipped, matching
	// the merge in update: a caller naming a block the page does not have is
	// working from a stale layout, and quietly removing the rest would report
	// success for something that did not happen.
	remove := make(map[uint64]bool, len(idStrs))
	for _, s := range idStrs {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid block ID %q: %w", s, err)
		}
		if !present[id] {
			return nil, fmt.Errorf("unknown block ID %d", id)
		}
		remove[id] = true
	}

	kept := make(cmpTypes.PageBlocks, 0, len(pg.Blocks))
	for _, b := range pg.Blocks {
		if !remove[b.BlockID] {
			kept = append(kept, b)
		}
	}
	pg.Blocks = kept

	// The service replaces the stored block set with the one it is handed, so
	// removal is expressible here even though the update tool's merge never
	// shrinks the set. Surviving blocks keep their xywh: closing the gap would be
	// a layout decision this tool has no basis to make.
	pg, err = cmpService.DefaultPage.Update(ctx, pg)
	if err != nil {
		return nil, toolkit.Errf("page block removal", err)
	}

	return toolkit.JSONResult(pg)
}

func (h *pageHandler) reorder(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	var parentID uint64
	parentRef, err := toolkit.Ref(args, "parent")
	if err != nil {
		return nil, err
	}
	if parentRef != "" {
		parent, err := findPageByAny(ctx, nsID, parentRef)
		if err != nil {
			return nil, toolkit.Errf("parent page lookup", err)
		}
		parentID = parent.ID
	}

	rawIDs, ok := args["pageIDs"]
	if !ok || rawIDs == nil {
		return nil, fmt.Errorf("pageIDs is required")
	}

	// parseStringArray refuses a JSON array of numbers, which is the point: an ID
	// that arrived as a number has already lost precision.
	idStrs, err := parseStringArray(rawIDs)
	if err != nil {
		return nil, fmt.Errorf("invalid pageIDs: must be a JSON array of ID strings: %w", err)
	}
	if len(idStrs) == 0 {
		return nil, fmt.Errorf("pageIDs must list at least one page")
	}

	pageIDs := make([]uint64, len(idStrs))
	for i, s := range idStrs {
		if pageIDs[i], err = strconv.ParseUint(s, 10, 64); err != nil {
			return nil, fmt.Errorf("invalid page ID %q: %w", s, err)
		}
	}

	if err = cmpService.DefaultPage.Reorder(ctx, nsID, parentID, pageIDs); err != nil {
		return nil, toolkit.Errf("page reorder", err)
	}

	return toolkit.TextResult("%d pages reordered", len(pageIDs)), nil
}

func (h *pageHandler) blockSchema(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	kind, err := toolkit.ReqStr(args, "kind")
	if err != nil {
		return nil, err
	}

	schema, exists := cmpTypes.PageBlockOptionSchemas[kind]
	if !exists {
		return nil, fmt.Errorf("unknown block kind %q: supported kinds are %s", kind, supportedBlockKinds())
	}

	return toolkit.JSONResult(schema)
}

// newPageItem projects a page onto the slim list shape.
func newPageItem(p *cmpTypes.Page) pageItem {
	return pageItem{
		ID:       p.ID,
		Title:    p.Title,
		Handle:   p.Handle,
		ParentID: p.SelfID,
		ModuleID: p.ModuleID,
		Visible:  p.Visible,
		Weight:   p.Weight,
	}
}

// pageNodes projects a nested page set onto the slim node shape, dropping blocks
// at every level of the tree.
func pageNodes(set cmpTypes.PageSet) []pageNode {
	out := make([]pageNode, 0, len(set))
	for _, p := range set {
		out = append(out, pageNode{
			pageItem: newPageItem(p),
			Children: pageNodes(p.Children),
		})
	}
	return out
}

// supportedBlockKinds lists the known block kinds, sorted so the error message is
// stable across calls.
func supportedBlockKinds() string {
	kinds := make([]string, 0, len(cmpTypes.PageBlockOptionSchemas))
	for k := range cmpTypes.PageBlockOptionSchemas {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	return strings.Join(kinds, ", ")
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
		return nil, toolkit.Errf("page lookup", err)
	}

	for _, p := range set {
		if strings.EqualFold(p.Title, ref) || strings.EqualFold(p.Handle, ref) {
			return p, nil
		}
	}

	return nil, fmt.Errorf("page %q not found", ref)
}

// inheritBlockLayout copies the stored xywh onto every incoming block that
// names an existing block but carries no layout of its own, so an options-only
// update leaves the page's layout untouched. An incoming block that does carry a
// width overrides the stored one, which is how a block gets moved or resized.
func inheritBlockLayout(existing, incoming cmpTypes.PageBlocks) {
	if len(existing) == 0 || len(incoming) == 0 {
		return
	}

	stored := make(map[uint64][4]int, len(existing))
	for _, b := range existing {
		if b.BlockID != 0 {
			stored[b.BlockID] = b.XYWH
		}
	}

	for i := range incoming {
		b := &incoming[i]
		if b.BlockID == 0 || b.XYWH[2] > 0 {
			continue
		}
		if xywh, ok := stored[b.BlockID]; ok {
			b.XYWH = xywh
		}
	}
}

// mergePageBlocks merges incoming blocks into the existing ones by blockID.
// Existing blocks not mentioned are kept; a block without a blockID is appended
// and numbered by the service.
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

// foldBlockOptionAliases renames the "module" and "chart" options an agent
// writes to the "moduleID" and "chartID" the block schemas actually name
// (RecordListBlockOptions, ChartBlockOptions).
//
// The webapp reads only the canonical key, so resolving an alias to a real ID
// and leaving it under its own name stored a value nothing renders: the block
// came up "No module selected" with a perfectly good module ID beside it.
// pagebuild.py has always renamed these; this is the same rule for the MCP path.
// An explicit canonical key wins — it is the one the caller meant.
func foldBlockOptionAliases(options map[string]any) {
	for alias, canonical := range map[string]string{"module": "moduleID", "chart": "chartID"} {
		ref, _ := options[alias].(string)
		if ref == "" {
			continue
		}

		if existing, _ := options[canonical].(string); existing == "" {
			options[canonical] = ref
		}
		delete(options, alias)
	}
}

// resolveBlockRefs translates human-readable module/chart references in block
// options to their uint64 ID strings, which the frontend SDK requires.
// Handles moduleID and chartID — including the "module"/"chart" aliases agents
// write, which are folded into those keys first — and the moduleID inside
// Calendar feeds and Metric metrics.
func resolveBlockRefs(ctx context.Context, nsID uint64, blocks cmpTypes.PageBlocks) error {
	for i := range blocks {
		b := &blocks[i]
		if b.Options == nil {
			continue
		}

		foldBlockOptionAliases(b.Options)

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

		// chartID → chart ID (Chart block). A handle left here is what the
		// webapp reports as "invalid ID".
		if ref, _ := b.Options["chartID"].(string); ref != "" {
			if _, err := strconv.ParseUint(ref, 10, 64); err != nil {
				ch, err := findChartByAny(ctx, nsID, ref)
				if err != nil {
					return fmt.Errorf("block %q: chart %q not found: %w", b.Title, ref, err)
				}
				b.Options["chartID"] = strconv.FormatUint(ch.ID, 10)
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

// parsePageBlocks unmarshals a JSON string or array into PageBlocks. Layout is
// left exactly as sent here; autoLayoutBlocks is what places the blocks that
// carry none, and it runs after the update path has merged.
func parsePageBlocks(raw any) (cmpTypes.PageBlocks, error) {
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

// Grid facts, mirrored from the gridstack setup the webapp page builder uses
// (client/web/unify/src/sections/compose/components/PageBlocks/Grid.vue).
// A quarter here is 12 columns, which is why a block sized as if the grid had
// 12 columns renders as a narrow strip down the left edge.
const (
	gridColumns = 48
	gridHalf    = gridColumns / 2
	gridQuarter = gridColumns / 4
)

// autoLayoutBlocks places the blocks that carry no width and leaves every block
// that carries one exactly where the caller put it.
//
// A width is what marks a block as deliberately positioned: it is the one
// coordinate that cannot be a meaningful zero, since a block 0 columns wide
// renders nothing. Auto-placed blocks flow left to right at a per-kind width and
// wrap to a new row when the current one fills, so a handful of Metric tiles
// land side by side instead of one per row. They start below the lowest block
// the caller positioned, which keeps a page that mixes both from stacking an
// auto-placed block on top of an explicit one.
func autoLayoutBlocks(blocks cmpTypes.PageBlocks) cmpTypes.PageBlocks {
	rowY := 0
	for i := range blocks {
		b := &blocks[i]
		if b.XYWH[2] <= 0 {
			continue
		}

		// Normalize only what the grid cannot render: an off-grid or negative
		// coordinate. Anything that fits is the caller's layout and stands.
		if b.XYWH[0] < 0 {
			b.XYWH[0] = 0
		}
		if b.XYWH[1] < 0 {
			b.XYWH[1] = 0
		}
		if b.XYWH[2] > gridColumns {
			b.XYWH[2] = gridColumns
		}
		if b.XYWH[3] <= 0 {
			b.XYWH[3] = defaultBlockHeight(b.Kind)
		}

		if bottom := b.XYWH[1] + b.XYWH[3]; bottom > rowY {
			rowY = bottom
		}
	}

	x, rowHeight := 0, 0
	for i := range blocks {
		b := &blocks[i]
		if b.XYWH[2] > 0 {
			continue
		}

		w := defaultBlockWidth(b.Kind)
		if x+w > gridColumns {
			rowY += rowHeight
			x, rowHeight = 0, 0
		}

		if b.XYWH[3] <= 0 {
			b.XYWH[3] = defaultBlockHeight(b.Kind)
		}
		b.XYWH[0], b.XYWH[1], b.XYWH[2] = x, rowY, w

		x += w
		if b.XYWH[3] > rowHeight {
			rowHeight = b.XYWH[3]
		}
	}

	return blocks
}

// defaultBlockWidth is how wide a block goes when the caller supplied no width:
// as wide as the block kind needs to be readable, and no wider, so the page uses
// its horizontal space instead of running down the left edge.
func defaultBlockWidth(kind string) int {
	switch kind {
	case "Metric", "Progress":
		// Tiles: a single number each, four to a row.
		return gridQuarter
	case "Chart", "Calendar", "Comment", "Content", "SocialFeed", "Automation":
		// Readable in half a page, so two sit side by side.
		return gridHalf
	default:
		// RecordList, Record, RecordOrganizer, ChatbotInbox and any kind added
		// later: tabular or form content that needs the full width.
		return gridColumns
	}
}

// defaultBlockHeight is in grid cells of 10px. Blocks CLIP silently rather than
// scroll, so these are the heights at which a kind shows something rather than
// the least it can be given.
func defaultBlockHeight(kind string) int {
	switch kind {
	case "Metric", "Progress":
		return 20
	case "Content", "Automation":
		return 20
	default:
		return 30
	}
}

// parsePageConfig unmarshals a JSON string or object into PageConfig.
func parsePageConfig(raw any) (cmpTypes.PageConfig, error) {
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

	if err := checkNavIcon(cfg.NavItem.Icon); err != nil {
		return cmpTypes.PageConfig{}, err
	}

	return cfg, nil
}

// checkNavIcon rejects the icon types the webapp cannot draw.
//
// The navigation renders a page icon as an <img>: a "link" icon is a URL in its
// own right and an "attachment" is a path on this instance, so both resolve. A
// "library" icon ("font-awesome://home") and an inline SVG have no URL to point
// at — the webapp still builds an <img> src out of them and the page ends up
// with a broken image next to its name. The failure is silent server-side, so it
// is refused here rather than stored. Nothing else can set an icon through these
// tools; an icon is a human's choice made in the page editor, which uploads the
// image it then references.
func checkNavIcon(icon *cmpTypes.PageConfigIcon) error {
	if icon == nil {
		return nil
	}

	switch icon.Type {
	case cmpTypes.IconTypeLibrary, cmpTypes.IconTypeInlineSvg, "svg":
		return fmt.Errorf(
			"config.navItem.icon: type %q is not rendered by the webapp and shows as a broken image; "+
				"only %q (an absolute image URL) and %q (an image uploaded in the page editor) resolve, "+
				"and a page needs no icon at all",
			icon.Type, cmpTypes.IconTypeLink, cmpTypes.IconTypeAttachment,
		)
	}

	return nil
}

// parsePageMeta unmarshals a JSON string or object into PageMeta.
func parsePageMeta(raw any) (cmpTypes.PageMeta, error) {
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
