package agentic

import (
	"context"
	_ "embed"
	"regexp"
	"strconv"
	"strings"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/weburl"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCP Apps pages for compose tools, rendered by hosts that support the
// extension. Everything a page needs is inlined, so none declares a CSP origin.

const (
	uiRecordLookup = "ui://human/compose/record-lookup"
	uiRecordReport = "ui://human/compose/record-report"
)

var (
	//go:embed mcpui/record_lookup.gen.html
	uiRecordLookupHTML []byte

	//go:embed mcpui/record_report.gen.html
	uiRecordReportHTML []byte
)

func registerUIResources(reg toolRegistrar) {
	reg.RegisterUIResource(hmcp.UIResource{
		URI:         uiRecordLookup,
		Name:        "Record lookup",
		Description: "Renders the records a compose_record_lookup call returned.",
		HTML:        uiRecordLookupHTML,
	})
	reg.RegisterUIResource(hmcp.UIResource{
		URI:         uiRecordReport,
		Name:        "Record report",
		Description: "Charts the groups a compose_record_report call returned.",
		HTML:        uiRecordReportHTML,
	})
}

type (
	// recordView is what the record lookup page renders beside the records:
	// the columns and where each record opens in Human.
	recordView struct {
		Module recordViewModule  `json:"module"`
		Fields []recordViewField `json:"fields"`
		Links  map[string]string `json:"links,omitempty"`
	}

	// reportView is what the record report page needs to label a report: the
	// field the rows are grouped by and each metric, in the order they were
	// asked for.
	reportView struct {
		Module    recordViewModule `json:"module"`
		Dimension *recordViewField `json:"dimension,omitempty"`
		Metrics   []reportMetric   `json:"metrics"`
	}

	// reportMetric is a metric's key in the result rows and the field it
	// aggregates, which is what the result's units are keyed by.
	reportMetric struct {
		Key   string `json:"key"`
		Field string `json:"field,omitempty"`
	}

	recordViewModule struct {
		Name   string `json:"name"`
		Handle string `json:"handle,omitempty"`
	}

	recordViewField struct {
		Name    string `json:"name"`
		Label   string `json:"label,omitempty"`
		Kind    string `json:"kind"`
		Multi   bool   `json:"multi,omitempty"`
		Options any    `json:"options,omitempty"`
	}
)

// withRecordView attaches the record lookup page's data to a lookup result.
// Without the module there are no columns to describe, so the result goes out
// as it is.
func withRecordView(ctx context.Context, res *mcp.CallToolResult, mod *cmpTypes.Module, set cmpTypes.RecordSet) *mcp.CallToolResult {
	if res == nil || mod == nil {
		return res
	}

	slug := ""
	if len(set) > 0 {
		slug = nsURLPart(ctx, mod.NamespaceID)
	}

	return hmcp.WithViewData(res, recordViewOf(mod, slug, set))
}

// recordViewOf describes mod's fields as columns and links each record in set
// under the namespace's url part. An empty slug means no links.
func recordViewOf(mod *cmpTypes.Module, slug string, set cmpTypes.RecordSet) recordView {
	v := recordView{
		Module: recordViewModule{Name: mod.Name, Handle: mod.Handle},
		Fields: make([]recordViewField, 0, len(mod.Fields)),
	}

	for _, f := range mod.Fields {
		v.Fields = append(v.Fields, viewFieldOf(f))
	}

	if slug != "" && len(set) > 0 {
		v.Links = make(map[string]string, len(set))
		for _, rec := range set {
			if url := weburl.ComposeRecord(slug, mod.ID, rec.ID); url != "" {
				v.Links[strconv.FormatUint(rec.ID, 10)] = url
			}
		}
	}

	return v
}

// viewFieldOf describes one field as a column.
func viewFieldOf(f *cmpTypes.ModuleField) recordViewField {
	vf := recordViewField{Name: f.Name, Label: f.Label, Kind: f.Kind, Multi: f.Multi}
	if f.Kind == "Select" {
		vf.Options = f.Options["options"]
	}
	return vf
}

// reportViewOf labels a report over mod grouped by dimension, which may wrap
// its field in DATE(). A dimension naming no field of mod is left out.
func reportViewOf(mod *cmpTypes.Module, dimension, metrics string) reportView {
	v := reportView{
		Module:  recordViewModule{Name: mod.Name, Handle: mod.Handle},
		Metrics: reportMetrics(metrics),
	}

	name := strings.TrimSpace(dimension)
	if m := dateBucket.FindStringSubmatch(name); m != nil {
		name = m[1]
	}
	if f := mod.Fields.FindByName(name); f != nil {
		vf := viewFieldOf(f)
		v.Dimension = &vf
	}

	return v
}

var dateBucket = regexp.MustCompile(`(?i)^DATE\(\s*([A-Za-z_][A-Za-z0-9_]*)\s*\)$`)

// reportMetrics names each metric by its result key — its alias, or the
// expression itself when it has none — and the field inside its aggregate.
func reportMetrics(metrics string) []reportMetric {
	out := []reportMetric{}
	for _, m := range splitMetrics(metrics) {
		rm := reportMetric{Key: strings.TrimSpace(m)}
		if a := asAlias.FindString(m); a != "" {
			words := strings.Fields(a)
			rm.Key = words[len(words)-1]
		}
		if f := metricFieldRef.FindStringSubmatch(m); f != nil {
			if ids := metricIdent.FindAllString(f[1], 1); len(ids) > 0 {
				rm.Field = ids[0]
			}
		}
		out = append(out, rm)
	}
	return out
}
