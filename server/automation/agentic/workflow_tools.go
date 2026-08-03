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
