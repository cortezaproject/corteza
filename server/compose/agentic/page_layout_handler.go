package agentic

import (
	"context"
	"fmt"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cast"
)

type pageLayoutHandler struct {
	reg toolRegistrar
}

// Declarations for these handlers are in page_layout_tools.go, in the same order.
func PageLayoutHandler(reg toolRegistrar) *pageLayoutHandler {
	h := &pageLayoutHandler{reg: reg}
	h.register()
	return h
}

// resolvePage resolves the namespace and page every layout tool takes.
func (h *pageLayoutHandler) resolvePage(ctx context.Context, args map[string]any) (uint64, *cmpTypes.Page, error) {
	nsRef, err := toolkit.ReqRef(args, "namespace")
	if err != nil {
		return 0, nil, err
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		return 0, nil, toolkit.Errf("namespace lookup", err)
	}

	pageRef, err := toolkit.ReqRef(args, "page")
	if err != nil {
		return 0, nil, err
	}

	pg, err := findPageByAny(ctx, ns.ID, pageRef)
	if err != nil {
		return 0, nil, err
	}

	return ns.ID, pg, nil
}

// pageLayouts returns the page's layouts in weight order, which is the order
// they are evaluated in.
func pageLayouts(ctx context.Context, nsID, pageID uint64) (cmpTypes.PageLayoutSet, error) {
	set, _, err := cmpService.DefaultPageLayout.Find(ctx, cmpTypes.PageLayoutFilter{
		NamespaceID: nsID,
		PageID:      pageID,
		Paging:      filter.Paging{Limit: 200},
		Sorting:     filter.Sorting{Sort: filter.SortExprSet{&filter.SortExpr{Column: "weight"}}},
	})
	if err != nil {
		return nil, toolkit.Errf("page layout lookup", err)
	}
	return set, nil
}

// findLayout resolves a layout by ID or handle, within one page — a handle is
// only unique per page, so the page is always part of the lookup.
func findLayout(ctx context.Context, nsID, pageID uint64, ref string) (*cmpTypes.PageLayout, error) {
	set, err := pageLayouts(ctx, nsID, pageID)
	if err != nil {
		return nil, err
	}

	for _, l := range set {
		if ref == fmt.Sprintf("%d", l.ID) || (l.Handle != "" && ref == l.Handle) {
			return l, nil
		}
	}

	return nil, fmt.Errorf("page layout %q not found on this page", ref)
}

func (h *pageLayoutHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	if ref := toolkit.Str(args, "layout"); ref != "" {
		layout, err := findLayout(ctx, nsID, pg.ID, ref)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(layout)
	}

	set, err := pageLayouts(ctx, nsID, pg.ID)
	if err != nil {
		return nil, err
	}

	type brief struct {
		PageLayoutID string `json:"pageLayoutID"`
		Handle       string `json:"handle"`
		Title        string `json:"title"`
		Weight       int    `json:"weight"`
		Blocks       int    `json:"placedBlocks"`
	}

	out := make([]brief, 0, len(set))
	for _, l := range set {
		out = append(out, brief{
			PageLayoutID: fmt.Sprintf("%d", l.ID),
			Handle:       l.Handle,
			Title:        l.Meta.Title,
			Weight:       l.Weight,
			Blocks:       len(l.Blocks),
		})
	}

	return toolkit.JSONResult(map[string]any{"layouts": out})
}

func (h *pageLayoutHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	layout := &cmpTypes.PageLayout{
		NamespaceID: nsID,
		PageID:      pg.ID,
		Handle:      toolkit.Str(args, "handle"),
		Meta:        cmpTypes.PageLayoutMeta{Title: toolkit.Str(args, "title")},
	}

	if err = applyLayoutArgs(layout, args); err != nil {
		return nil, err
	}

	layout, err = cmpService.DefaultPageLayout.Create(ctx, layout)
	if err != nil {
		return nil, toolkit.Errf("page layout creation", err)
	}

	return toolkit.JSONResult(layout)
}

func (h *pageLayoutHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "layout")
	if err != nil {
		return nil, err
	}

	layout, err := findLayout(ctx, nsID, pg.ID, ref)
	if err != nil {
		return nil, err
	}

	// Only what the caller sent is changed. A partial update that quietly
	// cleared the rest is the failure this tool family was written after.
	if v := toolkit.Str(args, "handle"); v != "" {
		layout.Handle = v
	}
	if v := toolkit.Str(args, "title"); v != "" {
		layout.Meta.Title = v
	}

	if err = applyLayoutArgs(layout, args); err != nil {
		return nil, err
	}

	layout, err = cmpService.DefaultPageLayout.Update(ctx, layout)
	if err != nil {
		return nil, toolkit.Errf("page layout update", err)
	}

	return toolkit.JSONResult(layout)
}

func (h *pageLayoutHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, pg, err := h.resolvePage(ctx, args)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqRef(args, "layout")
	if err != nil {
		return nil, err
	}

	set, err := pageLayouts(ctx, nsID, pg.ID)
	if err != nil {
		return nil, err
	}

	// A page with no layout renders nothing and cannot be opened in the
	// builder, so the last one is not a thing to delete by accident.
	if len(set) <= 1 {
		return nil, fmt.Errorf(
			"this is the page's only layout; a page without one renders nothing — " +
				"delete the page with compose_page_delete, or add another layout first")
	}

	layout, err := findLayout(ctx, nsID, pg.ID, ref)
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultPageLayout.DeleteByID(ctx, nsID, pg.ID, layout.ID); err != nil {
		return nil, toolkit.Errf("page layout deletion", err)
	}

	return toolkit.TextResult("page layout %d deleted", layout.ID), nil
}

// applyLayoutArgs reads the parameters shared by create and update. Each is
// optional and absent means "leave alone", which is what lets update change one
// thing about a layout without restating the others.
func applyLayoutArgs(layout *cmpTypes.PageLayout, args map[string]any) error {
	var blocks cmpTypes.PageLayoutBlocks
	if present, err := toolkit.JSONArg(args, "blocks", "blocks", &blocks); err != nil {
		return err
	} else if present {
		layout.Blocks = blocks
	}

	// Unmarshalled over what is there, so a button the caller omits keeps the
	// setting it had.
	if _, err := toolkit.JSONArg(args, "buttons", "buttons", &layout.Config.Buttons); err != nil {
		return err
	}

	if _, err := toolkit.JSONArg(args, "visibility", "visibility", &layout.Config.Visibility); err != nil {
		return err
	}

	if w, ok := args["weight"]; ok && w != nil {
		weight, err := cast.ToIntE(w)
		if err != nil {
			return fmt.Errorf("weight must be a number: %w", err)
		}
		layout.Weight = weight
	}

	return nil
}
