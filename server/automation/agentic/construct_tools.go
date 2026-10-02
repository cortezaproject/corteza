package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Two catalogues, two tools, and the split is the point.
//
// A TAQ step's ref comes from the construct library; a workflow step's ref comes
// from the workflow function registry. They are different registries with
// overlapping names, and naming one from the other stores a definition that
// passes every write check and then never runs. Both authoring tools' own
// descriptions warn about it harder than anything else they say, and until these
// tools existed the only place either catalogue was written down was a REST
// endpoint an MCP-only caller cannot reach.
//
// So the tool names carry the distinction rather than a parameter: a caller
// reaching for automation_taq_create finds automation_taq_construct_lookup
// beside it, and a caller reaching for automation_workflow_create finds
// automation_workflow_function_lookup. There is no single "list the functions"
// tool, deliberately — it would be the one call both callers make and half of
// them would be wrong.
//
// Like automation_event_type_lookup these declare neither limit nor pageCursor:
// there is no store to page, the whole catalogue is well inside the tool-result
// ceiling once the webapp's form spec is dropped, and the filters narrow it for
// a caller who does not want all of it.
//
// This file holds declarations only. The implementation is in
// construct_handler.go.

func (h *constructHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_construct_lookup",
			mcp.WithDescription(
				"List the construct library — the vocabulary a TAQ (Trigger Action Query) is built from. "+
					"Two catalogues in one call: 'functions', every ref a TAQ step of kind function or "+
					"iterator may name, each with its parameters, their exact type names and which are "+
					"required; and 'triggers', every legal resourceType/eventType pair a TAQ trigger may "+
					"use. There are around 19 functions and 22 triggers, so the default is to return both "+
					"whole. "+
					"Read this before automation_taq_create or automation_taq_update. A step whose 'ref' is "+
					"not in this catalogue is still stored and still returns success — it comes back with a "+
					"'function.unknown' issue, 'runnable' false, and automation_taq_exec then answers "+
					"'manager: executable not found'. The same goes for a trigger pair that is not listed: "+
					"it stores without complaint and never fires. "+
					"This is NOT the workflow function registry. A TAQ and a workflow are different engines "+
					"with different vocabularies, and most of the 93 workflow functions do not exist here. "+
					"For a workflow use automation_workflow_function_lookup instead — mixing the two is the "+
					"most common way to author an automation that stores cleanly and never runs. "+
					"Bind a step's arguments by 'argumentName', spelled exactly as listed here, and give "+
					"each one a 'type' spelled exactly as it appears in that parameter's 'types' array. "+
					"Note the workflow registry keys its parameters by 'name' instead, which is another "+
					"reason a ref copied across engines does not work. "+
					"The webapp's form layout for each entry ('segments') is omitted; nothing an author "+
					"needs is in it and it is four fifths of the payload.",
			),
			mcp.WithString("catalogue", mcp.Enum("functions", "triggers"), mcp.Description("Return one catalogue only: \"functions\" for step refs, \"triggers\" for resourceType/eventType pairs. Omit for both.")),
			mcp.WithString("ref", mcp.Description("Narrow the functions to one. An exact ref is matched first, then any ref containing it, case-insensitively. A ref that matches nothing answers with the full list of refs rather than an empty result. Ignored when catalogue is \"triggers\".")),
			mcp.WithString("resourceType", mcp.Description("Narrow the triggers to one resource, e.g. \"compose:record\" for exactly that or \"compose\" for it and every sub-resource. Case-insensitive. Ignored when catalogue is \"functions\".")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup TAQ construct library",
		h.taqConstructLookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_workflow_function_lookup",
			mcp.WithDescription(
				"List the workflow function registry — every ref a workflow step of kind 'function' or "+
					"'iterator' may name, with its parameters, their exact type names and which are "+
					"required, plus the results it puts into scope for later steps. There are around 93, "+
					"returned whole by default. "+
					"Read this before automation_workflow_create or automation_workflow_update. A step "+
					"naming a ref that is not here is stored and returns success, and the workflow then "+
					"fails to convert into anything runnable — a guessed ref is not a guess that gets "+
					"corrected, it is a workflow that silently does nothing. "+
					"This is NOT the TAQ construct library. Most of these functions are unavailable to a "+
					"TAQ, and a TAQ step naming one of them stores with a 'function.unknown' issue and "+
					"never runs. For a TAQ use automation_taq_construct_lookup instead. "+
					"A workflow function's parameters are keyed by 'name' — the construct library keys "+
					"its own by 'argumentName' — so an argument list copied from one engine to the other "+
					"is wrong even where the ref happens to exist in both. "+
					"'kind' says how the step is wired: a 'function' runs once, an 'iterator' runs its "+
					"body once per item. 'labels' records where the registry advertises the function.",
			),
			mcp.WithString("ref", mcp.Description("Narrow to one function. An exact ref is matched first, then any ref containing it, case-insensitively. A ref that matches nothing answers with the full list of refs rather than an empty result.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup workflow functions",
		h.workflowFunctionLookup,
	)
}
