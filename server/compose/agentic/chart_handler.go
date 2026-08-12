package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
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
		// moduleID decodes as a quoted uint64, so a module handle fails here
		// rather than in the checks below, and does it in encoding/json's words.
		if strings.Contains(err.Error(), "invalid use of ,string struct tag") {
			return cmpTypes.ChartConfig{}, fmt.Errorf(
				"config: moduleID must be the numeric module ID as a string, from compose_module_lookup, not a module handle (%w)", err,
			)
		}
		return cmpTypes.ChartConfig{}, fmt.Errorf("config must be a JSON object: %w", err)
	}

	if len(cfg.Reports) == 0 {
		return cmpTypes.ChartConfig{}, fmt.Errorf("config must contain at least one report")
	}

	if err := checkChartReports(cfg.Reports); err != nil {
		return cmpTypes.ChartConfig{}, err
	}

	return cfg, nil
}

// What the webapp's chart renderer accepts, mirrored from
// lib/js/src/compose/types/chart/util.ts (aggregateFunctions, dimensionFunctions
// and the ChartType enum). The renderer validates a chart only when it draws it,
// so a config it rejects is stored happily and fails in front of a person
// instead — which is why these are checked here.
var (
	// STD is deliberately absent: it was dropped from the webapp's aggregate
	// options, and only the MySQL store backend ever implemented it — the
	// postgres and sqlite QL handlers list it as unsupported.
	chartAggregates = []string{"SUM", "MAX", "MIN", "AVG"}
	chartTypes      = []string{"pie", "bar", "line", "doughnut", "funnel", "gauge", "radar", "scatter"}

	// The first one groups by the field's own values; the rest bucket a date
	// field. Its odd spelling is the literal value the renderer looks for.
	chartNoGrouping = "(no grouping / buckets)"
	chartModifiers  = []string{chartNoGrouping, "DATE", "WEEK", "MONTH", "QUARTER", "YEAR"}
)

// checkChartReports rejects the report configs that render as an error message
// where the chart should be, and fills in the one default that is unambiguous.
//
// Every rule here is the renderer's own (BaseChart.isValid / mtrCheck / dimCheck
// in lib/js). The names in the errors are the ones a person sees on the page —
// "Metrics aggregate not defined" is what a metric without an aggregate looks
// like once it is too late.
func checkChartReports(reports []*cmpTypes.ChartConfigReport) error {
	for i, r := range reports {
		if r == nil {
			return fmt.Errorf("report %d is empty", i+1)
		}

		// The renderer needs a module to query and only ever looks at the ID; a
		// handle here reads as "module not found" and the chart draws nothing.
		if r.ModuleID == 0 {
			return fmt.Errorf(
				"report %d: moduleID is required — pass the numeric ID from compose_module_lookup as a string, not a module handle",
				i+1,
			)
		}

		if len(r.Metrics) == 0 {
			return fmt.Errorf("report %d: at least one metric is required — what the chart plots", i+1)
		}

		gauge := false
		for j, m := range r.Metrics {
			isGauge, err := checkChartMetric(i+1, j+1, m)
			if err != nil {
				return err
			}
			gauge = gauge || isGauge
		}

		for j, d := range r.Dimensions {
			if err := checkChartDimension(i+1, j+1, d, gauge); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkChartMetric validates one metric in place and reports whether it is a
// gauge, which changes what its dimensions must carry.
func checkChartMetric(report, idx int, m map[string]any) (bool, error) {
	where := fmt.Sprintf("report %d, metric %d", report, idx)

	field, _ := m["field"].(string)
	if field == "" {
		return false, fmt.Errorf(
			`%s: "field" is required — "count" to count records, or the name of a numeric module field to aggregate`,
			where,
		)
	}

	// "count" is the one field the renderer computes itself; everything else is
	// a column it has to reduce, and it has no default for how.
	if !strings.EqualFold(field, "count") {
		aggregate, _ := m["aggregate"].(string)
		if aggregate == "" {
			return false, fmt.Errorf(
				`%s (field %q): "aggregate" is required for any field other than "count" — one of %s. `+
					`Without it the chart renders as "Metrics aggregate not defined"`,
				where, field, strings.Join(chartAggregates, ", "),
			)
		}

		normalized := strings.ToUpper(strings.TrimSpace(aggregate))
		if !slices.Contains(chartAggregates, normalized) {
			return false, fmt.Errorf(
				`%s (field %q): aggregate %q is not one of %s`,
				where, field, aggregate, strings.Join(chartAggregates, ", "),
			)
		}
		m["aggregate"] = normalized
	}

	kind, _ := m["type"].(string)
	if kind == "" {
		return false, fmt.Errorf(
			`%s: "type" is required — one of %s. It is also what picks the chart's shape, so a funnel or gauge is made by naming it here`,
			where, strings.Join(chartTypes, ", "),
		)
	}

	normalized := strings.ToLower(strings.TrimSpace(kind))
	if !slices.Contains(chartTypes, normalized) {
		return false, fmt.Errorf(`%s: type %q is not one of %s`, where, kind, strings.Join(chartTypes, ", "))
	}
	m["type"] = normalized

	return normalized == "gauge", nil
}

// checkChartDimension validates one dimension in place. A dimension is what the
// metric is grouped by, so it needs a field and a way to group it.
func checkChartDimension(report, idx int, d map[string]any, gauge bool) error {
	where := fmt.Sprintf("report %d, dimension %d", report, idx)

	field, _ := d["field"].(string)
	if field == "" {
		return fmt.Errorf(
			`%s: "field" is required — the module field whose values the metric is grouped by (a Select field groups well; a date field pairs with a modifier)`,
			where,
		)
	}

	// An absent modifier means the plainest grouping there is, by the field's own
	// values, so it is filled in rather than refused. A wrong one is refused: the
	// renderer would silently fall back to this same default and quietly draw a
	// chart nobody asked for.
	modifier, _ := d["modifier"].(string)
	if modifier == "" {
		d["modifier"] = chartNoGrouping
	} else if !slices.Contains(chartModifiers, modifier) {
		if normalized := strings.ToUpper(strings.TrimSpace(modifier)); slices.Contains(chartModifiers, normalized) {
			d["modifier"] = normalized
		} else {
			return fmt.Errorf(
				`%s: modifier %q is not one of %s`,
				where, modifier, strings.Join(quoted(chartModifiers), ", "),
			)
		}
	}

	// A gauge draws its bands from the dimension, and renders as "Dimensions
	// steps not defined" without them.
	if gauge {
		meta, _ := d["meta"].(map[string]any)
		steps, _ := meta["steps"].([]any)
		if len(steps) == 0 {
			return fmt.Errorf(
				`%s: a gauge metric needs this dimension's meta.steps — e.g. "meta":{"steps":[{"value":0},{"value":50},{"value":100}]}`,
				where,
			)
		}
	}

	return nil
}

// quoted wraps values for an error message, since one chart modifier contains
// spaces and slashes and is unreadable bare.
func quoted(vv []string) []string {
	out := make([]string, len(vv))
	for i, v := range vv {
		out[i] = strconv.Quote(v)
	}
	return out
}
