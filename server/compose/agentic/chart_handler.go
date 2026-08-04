package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type chartHandler struct {
	reg toolRegistrar
}

// Declarations for these handlers are in chart_tools.go, in the same order.
func ChartHandler(reg toolRegistrar) *chartHandler {
	h := &chartHandler{reg: reg}
	h.register()
	return h
}

func (h *chartHandler) resolveNs(ctx context.Context, args map[string]any) (uint64, error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return 0, toolkit.Errf("namespace lookup", err)
	}

	return ns.ID, nil
}

func (h *chartHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	// Single-item mode returns the raw service type, config included — that is
	// what the caller is after once they have picked a chart.
	if ref := toolkit.Str(args, "chart"); ref != "" {
		c, err := findChartByAny(ctx, nsID, ref)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(c)
	}

	f := cmpTypes.ChartFilter{NamespaceID: nsID}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := cmpService.DefaultChart.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("chart list", err)
	}

	// List mode is a slim projection: a chart's config is incidental when the
	// caller is picking one out of a list, and a namespace's worth of configs
	// would swamp the result ceiling.
	type chartItem struct {
		ID     uint64 `json:"chartID,string"`
		Name   string `json:"name"`
		Handle string `json:"handle"`
	}

	items := make([]chartItem, len(set))
	for i, c := range set {
		items[i] = chartItem{ID: c.ID, Name: c.Name, Handle: c.Handle}
	}

	return toolkit.JSONResult(map[string]any{
		"charts": items,
		// The cursor marshals itself to the base64 form parseCursor accepts;
		// its String() is a debug rendering and does not round-trip.
		"nextPageCursor": out.NextPage,
	})
}

func (h *chartHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	cfg, err := parseChartConfig(args["config"])
	if err != nil {
		return nil, err
	}

	c := &cmpTypes.Chart{
		NamespaceID: nsID,
		Name:        name,
		Handle:      toolkit.Str(args, "handle"),
		Config:      cfg,
	}

	c, err = cmpService.DefaultChart.Create(ctx, c)
	if err != nil {
		return nil, toolkit.Errf("chart creation", err)
	}

	return toolkit.JSONResult(c)
}

func (h *chartHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqStr(args, "chart")
	if err != nil {
		return nil, err
	}

	c, err := findChartByAny(ctx, nsID, ref)
	if err != nil {
		return nil, err
	}

	// Absent leaves the field unchanged; present-and-empty clears it.
	if raw, ok := args["name"]; ok && raw != nil {
		v, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("name must be a string")
		}
		c.Name = v
	}

	if raw, ok := args["handle"]; ok && raw != nil {
		v, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("handle must be a string")
		}
		c.Handle = v
	}

	// Config deliberately replaces wholesale rather than merging: a report set
	// is a collection, and a partial merge would leave the caller unable to
	// remove a report or a metric.
	if raw, ok := args["config"]; ok && raw != nil {
		cfg, err := parseChartConfig(raw)
		if err != nil {
			return nil, err
		}
		c.Config = cfg
	}

	c, err = cmpService.DefaultChart.Update(ctx, c)
	if err != nil {
		return nil, toolkit.Errf("chart update", err)
	}

	return toolkit.JSONResult(c)
}

func (h *chartHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqStr(args, "chart")
	if err != nil {
		return nil, err
	}

	c, err := findChartByAny(ctx, nsID, ref)
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultChart.DeleteByID(ctx, nsID, c.ID); err != nil {
		return nil, toolkit.Errf("chart delete", err)
	}

	return toolkit.TextResult("chart %d deleted", c.ID), nil
}

func (h *chartHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, err := h.resolveNs(ctx, args)
	if err != nil {
		return nil, err
	}

	// An ID rather than the name-or-handle ref the other chart tools take:
	// findChartByAny goes through lookup and search, and both refuse a chart with
	// a deleted marker, so there is nothing left to resolve a name against.
	chartID, err := toolkit.ReqID(args, "chartID")
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultChart.UndeleteByID(ctx, nsID, chartID); err != nil {
		return nil, toolkit.Errf("chart undelete", err)
	}

	return toolkit.TextResult("chart %d restored", chartID), nil
}

// findChartByAny resolves a chart reference (numeric ID, handle, or name)
// within a namespace.
func findChartByAny(ctx context.Context, namespaceID uint64, ref string) (*cmpTypes.Chart, error) {
	if ref == "" {
		return nil, fmt.Errorf("chart reference is required")
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		return cmpService.DefaultChart.FindByID(ctx, namespaceID, id)
	}

	if c, err := cmpService.DefaultChart.FindByHandle(ctx, namespaceID, ref); err == nil {
		return c, nil
	}

	// Name is not a filterable column, so the fallback is a scan.
	set, _, err := cmpService.DefaultChart.Search(ctx, cmpTypes.ChartFilter{NamespaceID: namespaceID})
	if err != nil {
		return nil, toolkit.Errf("chart lookup", err)
	}

	for _, c := range set {
		if strings.EqualFold(c.Name, ref) || strings.EqualFold(c.Handle, ref) {
			return c, nil
		}
	}

	return nil, fmt.Errorf("chart %q not found", ref)
}

// parseChartConfig unmarshals a JSON string or object into ChartConfig, because
// models produce both.
func parseChartConfig(raw any) (cmpTypes.ChartConfig, error) {
	var data []byte

	switch v := raw.(type) {
	case nil:
		return cmpTypes.ChartConfig{}, fmt.Errorf("config is required")
	case string:
		data = []byte(v)
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return cmpTypes.ChartConfig{}, fmt.Errorf("cannot encode config: %w", err)
		}
	}

	var cfg cmpTypes.ChartConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cmpTypes.ChartConfig{}, fmt.Errorf("config must be a JSON object: %w", err)
	}

	if len(cfg.Reports) == 0 {
		return cmpTypes.ChartConfig{}, fmt.Errorf("config must contain at least one report")
	}

	return cfg, nil
}
