package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in workflow_handler.go, in this order.

func (h *workflowHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("automation_workflow_lookup",
			mcp.WithDescription(
				"Look up one automation workflow, or list the workflows on this instance. Provide 'workflow' "+
					"(ID or handle) to fetch a single workflow in full, including its step and path graph; omit "+
					"it to list workflows as summaries — workflow ID, handle, name and enabled flag — and then "+
					"fetch the one you want. List mode never returns the graph, so a listing cannot tell you "+
					"what a workflow does. "+
					"Narrow the list with 'query': it matches the handle as a case-insensitive substring and "+
					"does not search the display name; a query containing spaces that finds nothing is retried "+
					"once as a slug (\"My Workflow\" is retried as \"my_workflow\"). 'query' is ignored when "+
					"'workflow' is given. "+
					"Workflows are the step-and-path automations; Trigger-Action-Query automations are a "+
					"separate resource — use automation_taq_lookup for those.",
			),
			mcp.WithString("workflow", mcp.Description("Workflow ID as string (to prevent precision loss) or handle. Omit to list all.")),
			mcp.WithString("query", mcp.Description("Case-insensitive substring match on the workflow handle. Applies to list mode only.")),
			mcp.WithString("limit", mcp.Description("Maximum workflows to return in list mode, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous list response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup workflow",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_workflow_delete",
			mcp.WithDescription(
				"Delete an automation workflow. The delete is soft — its step and path graph is kept — and "+
					"automation_workflow_undelete restores it intact. Deleting unregisters its triggers, so "+
					"nothing fires the workflow any more and automation_workflow_exec will not find it. "+
					"Any TAQ or trigger pointing at this workflow stops working; check "+
					"automation_trigger_lookup with this workflow first if you are not sure what depends on "+
					"it. Prefer disabling the workflow when you only want to pause it — a disabled workflow "+
					"is still listed, a deleted one is not findable by handle at all and restoring it needs "+
					"an ID you must have kept.",
			),
			mcp.WithString("workflowID", mcp.Required(), mcp.Description("Workflow ID as a string (to prevent precision loss). Find it with automation_workflow_lookup.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete workflow",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_workflow_undelete",
			mcp.WithDescription(
				"Restore a deleted automation workflow. Deleting a workflow is soft — its step and path "+
					"graph is kept — so restoring brings the workflow back intact, re-registers its triggers "+
					"and makes it executable again if it is enabled. "+
					"You cannot search for a deleted workflow: automation_workflow_lookup leaves deleted "+
					"workflows out of its list and does not find them by handle, so you need the workflow's "+
					"numeric ID from before the delete, or from an audit trail. That ID still resolves: pass "+
					"it as 'workflow' to automation_workflow_lookup to see what you are about to restore — a "+
					"deleted workflow comes back with 'deletedAt' set. "+
					"Restoring a workflow that is not deleted succeeds and changes nothing. To restore a "+
					"Trigger-Action-Query automation use automation_taq_undelete instead.",
			),
			mcp.WithString("workflowID", mcp.Required(), mcp.Description("Workflow ID as a string (to prevent precision loss). A handle does not work here — a deleted workflow cannot be resolved by handle.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete workflow",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_workflow_exec",
			mcp.WithDescription(
				"Run a workflow by ID or handle, wait for it to finish, and return the variables it produced. "+
					"This performs the workflow's real side effects — it is not a way to inspect what a "+
					"workflow does; use automation_workflow_lookup for that. "+
					"The workflow must be enabled: executing a disabled workflow fails rather than queueing. "+
					"'input' supplies the workflow's input variables and overrides any input the workflow's own "+
					"trigger defines. "+
					"For Trigger-Action-Query automations use automation_taq_exec instead.",
			),
			mcp.WithString("workflow", mcp.Required(), mcp.Description("Workflow ID as string (to prevent precision loss) or handle")),
			mcp.WithString("input", mcp.Description("JSON object of input variables")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Execute workflow",
		h.exec,
	)
}
