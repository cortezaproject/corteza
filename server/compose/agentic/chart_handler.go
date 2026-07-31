package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type chartHandler struct {
	reg toolRegistrar
}

func ChartHandler(reg toolRegistrar) *chartHandler {
	h := &chartHandler{reg: reg}
	h.register()
	return h
}

// canonical config contract validated against the unify webapp renderer
// (client/web/unify/src/sections/compose/components/Chart/, lib/js chart types)
const chartConfigDoc = `JSON chart configuration. Canonical shape:
{"colorScheme":"tableau.Tableau10","reports":[{"moduleID":"<id from compose_module_lookup>","filter":"","dimensions":[{"field":"<grouping field>","modifier":"(no grouping / buckets)","conditions":{}}],"metrics":[{"field":"count","type":"doughnut"}]}]}
Metric "field":"count" counts records; a numeric module field aggregates instead. "type" per metric: doughnut, pie, bar, line. The dimension field is what values are grouped by (Select fields work well).`

func (h *chartHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_lookup",
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDescription("List charts in a namespace, or look one up by name, handle, or ID. Omit 'chart' to list all. Chart blocks on pages reference charts by chartID."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("chart", mcp.Description("Chart name, handle, or ID. Omit to list all charts in the namespace.")),
		),
		"Lookup chart",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_create",
			mcp.WithDescription("Create a chart in a namespace. After creating, place it on a page with a Chart block: {\"kind\":\"Chart\",\"options\":{\"chartID\":\"<created ID>\"}}."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Chart name")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier")),
			mcp.WithString("config", mcp.Required(), mcp.Description(chartConfigDoc)),
		),
		"Create chart",
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_update",
			mcp.WithDescription("Update an existing chart. Pass only the fields you want to change; config replaces the whole configuration."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("chart", mcp.Required(), mcp.Description("Chart name, handle, or ID")),
			mcp.WithString("name", mcp.Description("New name")),
			mcp.WithString("handle", mcp.Description("New handle")),
			mcp.WithString("config", mcp.Description(chartConfigDoc)),
		),
		"Update chart",
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_chart_delete",
			mcp.WithDescription("Delete a chart by name, handle, or ID. Remove Chart blocks referencing it from pages first — they render empty otherwise."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace name, handle, slug, or ID")),
			mcp.WithString("chart", mcp.Required(), mcp.Description("Chart name, handle, or ID")),
		),
		"Delete chart",
		h.del,
	)
}

func (h *chartHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	chartRef, _ := args["chart"].(string)
	if chartRef == "" {
		set, _, err := cmpService.DefaultChart.Search(ctx, cmpTypes.ChartFilter{NamespaceID: ns.ID})
		if err != nil {
			return nil, fmt.Errorf("chart list failed: %w", err)
		}
		type chartItem struct {
			ID     uint64 `json:"chartID,string"`
			Name   string `json:"name"`
			Handle string `json:"handle"`
		}
		items := make([]chartItem, len(set))
		for i, c := range set {
			items[i] = chartItem{ID: c.ID, Name: c.Name, Handle: c.Handle}
		}
		out, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal charts: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	c, err := findChartByAny(ctx, ns.ID, chartRef)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chart: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *chartHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	name, _ := args["name"].(string)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	c := &cmpTypes.Chart{
		NamespaceID: ns.ID,
		Name:        name,
	}
	if v, ok := args["handle"].(string); ok && v != "" {
		c.Handle = v
	}

	cfg, err := parseChartConfig(args["config"])
	if err != nil {
		return nil, err
	}
	c.Config = cfg

	c, err = cmpService.DefaultChart.Create(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("chart creation failed: %w", err)
	}

	out, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chart: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *chartHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	chartRef, _ := args["chart"].(string)
	c, err := findChartByAny(ctx, ns.ID, chartRef)
	if err != nil {
		return nil, err
	}

	if v, ok := args["name"].(string); ok && v != "" {
		c.Name = v
	}
	if v, ok := args["handle"].(string); ok && v != "" {
		c.Handle = v
	}
	if raw, ok := args["config"]; ok && raw != nil {
		cfg, err := parseChartConfig(raw)
		if err != nil {
			return nil, err
		}
		c.Config = cfg
	}

	c, err = cmpService.DefaultChart.Update(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("chart update failed: %w", err)
	}

	out, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chart: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *chartHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	chartRef, _ := args["chart"].(string)
	c, err := findChartByAny(ctx, ns.ID, chartRef)
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultChart.DeleteByID(ctx, ns.ID, c.ID); err != nil {
		return nil, fmt.Errorf("chart deletion failed: %w", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("chart %d deleted", c.ID)), nil
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
	set, _, err := cmpService.DefaultChart.Search(ctx, cmpTypes.ChartFilter{NamespaceID: namespaceID})
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

// parseChartConfig unmarshals a JSON string or object into ChartConfig.
func parseChartConfig(raw interface{}) (cmpTypes.ChartConfig, error) {
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
