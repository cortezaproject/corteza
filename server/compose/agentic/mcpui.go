package agentic

import (
	_ "embed"

	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
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
