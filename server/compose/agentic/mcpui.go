package agentic

import (
	"context"
	_ "embed"
	"strconv"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/weburl"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCP Apps pages for compose tools, rendered by hosts that support the
// extension. Everything a page needs is inlined, so none declares a CSP origin.

const uiRecordLookup = "ui://human/compose/record-lookup"

//go:embed mcpui/record_lookup.gen.html
var uiRecordLookupHTML []byte

func registerUIResources(reg toolRegistrar) {
	reg.RegisterUIResource(hmcp.UIResource{
		URI:         uiRecordLookup,
		Name:        "Record lookup",
		Description: "Renders the records a compose_record_lookup call returned.",
		HTML:        uiRecordLookupHTML,
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
		vf := recordViewField{Name: f.Name, Label: f.Label, Kind: f.Kind, Multi: f.Multi}
		if f.Kind == "Select" {
			vf.Options = f.Options["options"]
		}
		v.Fields = append(v.Fields, vf)
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
