package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// TAQ is Trigger Action Query. The service and REST resource behind these tools
// is ngAutomation; the tool surface never says so — see mcp/RESOURCES.md.
//
// This file holds declarations only. Implementations are in taq_handler.go, in
// the same order.

func (h *taqHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_lookup",
			mcp.WithDescription(
				"List TAQs (Trigger Action Query automations), or fetch one whole. "+
					"Omit 'taq' to list: a list row is a compact form — automationID, handle, name, enabled — "+
					"enough to pick a TAQ, not enough to run one, because a TAQ carries its entire trigger, "+
					"step and path graph and one graph per row would swamp the result. "+
					"Pass 'taq' — an ID or a handle — to fetch that single TAQ in full, including its triggers "+
					"with their entry-point handles and input schemas. Do this before automation_taq_exec: the "+
					"trigger's input schema is the only place the expected input parameter names are written "+
					"down. "+
					"'query' matches the handle only — not the name — as a case-insensitive substring. A query "+
					"containing spaces that matches nothing is retried lowercased with spaces turned into "+
					"underscores, because handles are slugs. "+
					"Disabled TAQs are omitted from the list unless you set 'includeDisabled'; deleted ones are "+
					"always omitted.",
			),
			mcp.WithString("taq", mcp.Description("TAQ ID as a string (to prevent precision loss), or handle. Omit to list instead.")),
			mcp.WithString("query", mcp.Description("Case-insensitive substring of the handle. Ignored when 'taq' is given.")),
			mcp.WithBoolean("includeDisabled", mcp.Description("Also list disabled TAQs. Ignored when 'taq' is given.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup TAQ",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_delete",
			mcp.WithDescription(
				"Delete a TAQ (Trigger Action Query automation). The delete is soft — its triggers, steps "+
					"and paths are all retained — and automation_taq_undelete restores the whole definition. "+
					"Deleting stops the TAQ running: its triggers are unregistered, so nothing it was "+
					"listening for fires it any more. "+
					"Prefer disabling a TAQ over deleting it when you only want to pause it, because a "+
					"disabled TAQ still appears in automation_taq_lookup with includeDisabled, whereas a "+
					"deleted one cannot be found at all — not by handle, not by search — and restoring it "+
					"needs an ID you must have kept.",
			),
			mcp.WithString("taqID", mcp.Required(), mcp.Description("TAQ ID as a string (to prevent precision loss). Find it with automation_taq_lookup.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete TAQ",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_undelete",
			mcp.WithDescription(
				"Restore a deleted TAQ (Trigger Action Query automation). Deleting a TAQ is soft — its "+
					"triggers, steps and paths are all kept — so restoring puts the whole definition back, "+
					"re-registers its triggers, and makes it runnable again if it is enabled. "+
					"You cannot search for a deleted TAQ: automation_taq_lookup omits deleted TAQs from its "+
					"list and does not find them by handle either, so you must already hold the TAQ's numeric "+
					"ID — the 'automationID' automation_taq_lookup reported before the delete, or an ID from "+
					"an audit trail. That ID does still resolve: pass it as 'taq' to automation_taq_lookup to "+
					"see the TAQ, with 'deletedAt' set, before you restore it. "+
					"Restoring a TAQ that is not deleted succeeds and changes nothing.",
			),
			mcp.WithString("taqID", mcp.Required(), mcp.Description("TAQ ID as a string (to prevent precision loss). A handle does not work here — a deleted TAQ cannot be resolved by handle.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete TAQ",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_exec",
			mcp.WithDescription(
				"Run a TAQ and wait for it to finish. Returns the execution result — executionID, status, "+
					"error, timings — not the values the TAQ produced; pass the returned executionID to "+
					"automation_taq_execution_trace to see each step's input and output. "+
					"This runs the automation for real, doing whatever its steps do. There is no dry run, so "+
					"read the TAQ with automation_taq_lookup first if you are unsure what it will do. "+
					"'entryPoint' names the trigger to start from, by trigger handle; omitted, the TAQ's first "+
					"trigger is used. 'input' must satisfy that trigger's input schema, which "+
					"automation_taq_lookup returns: keys are matched to the schema case-insensitively, and a "+
					"missing required parameter fails the call before the TAQ starts.",
			),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID as a string (to prevent precision loss), or handle.")),
			mcp.WithString("entryPoint", mcp.Description("Handle of the trigger to start from. Defaults to the TAQ's first trigger.")),
			mcp.WithString("input", mcp.Description("JSON object of input parameters, shaped by the entry-point trigger's input schema. Omit when the trigger takes no input.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Execute TAQ",
		h.exec,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_executions",
			mcp.WithDescription(
				"List past runs of one TAQ: executionID, status, any error, start and end time and duration. "+
					"Use it to check whether a TAQ ran and whether it succeeded, then pass an executionID to "+
					"automation_taq_execution_trace for the step-by-step detail of one run. "+
					"Every retained execution of the TAQ is returned — the service takes no filter and no "+
					"paging — so a heavily used TAQ can exceed the result size limit.",
			),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID as a string (to prevent precision loss), or handle.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"List TAQ executions",
		h.executions,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_execution_trace",
			mcp.WithDescription(
				"Get the step-by-step trace of one TAQ run: a stack frame per executed step, with its handle, "+
					"kind, arguments, input, output, timings and error. This is the tool for diagnosing why a "+
					"run failed or produced what it did. "+
					"Get the executionID from automation_taq_executions, or from the result of "+
					"automation_taq_exec. A long run traces every step it took, so a trace can exceed the "+
					"result size limit.",
			),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID as a string (to prevent precision loss), or handle.")),
			mcp.WithString("executionID", mcp.Required(), mcp.Description("Execution ID as a string (to prevent precision loss), from automation_taq_executions or automation_taq_exec.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Get TAQ execution trace",
		h.executionTrace,
	)
}
