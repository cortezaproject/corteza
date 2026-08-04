package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations only. Implementations are in workflow_handler.go, in this order.

// The graph documentation below is shared by create and update, because the
// shape of a step or a path does not change between them. It carries what no
// JSON schema can: which fields each step kind expects, that a function's 'ref'
// comes from a registry only the running server knows, and that the order of
// the paths out of a branching step is semantics rather than presentation.
const (
	workflowStepsDoc = `JSON array of steps. Give every step a "stepID" of your own choosing, as a QUOTED string, unique within this workflow — paths reference it. ` +
		`Shape: {"stepID":"1","kind":"expressions","ref":"","arguments":[],"results":[],"meta":{"name":"","visual":{"id":"1","parent":"1","value":"Label","xywh":[240,120,180,64]}}}. ` +
		`An argument or a result is {"target":"...","expr":"...","value":...,"type":"..."}: "target" is the name (on a function step, the parameter the argument feeds; on an expressions step or on any result, the scope variable it writes), ` +
		`"expr" is an expression evaluated against the workflow scope, "value" is a literal used when there is no "expr", and "type" is the expression type — "Any" when omitted, except on a function or iterator argument, where it must match one of that parameter's declared types exactly. ` +
		`Step kinds and what each one expects: ` +
		`"expressions" — no ref, one or more arguments, at most one outbound path, e.g. {"target":"total","expr":"1 + 1","type":"Integer"}. ` +
		`"function" — "ref" names a construct from GET /automation/functions/ (93 of them; NOT /automation/construct-library/functions, which is the 17-entry registry Trigger-Action-Query automations use and which a workflow must not be built from), ` +
		`arguments match that construct's parameters by "target", results copy its outputs into scope as {"target":"<variable>","expr":"<result name>","type":"<type>"}, at most one outbound path. ` +
		`"iterator" — like function but the ref must be an iterating construct (composeRecordsEach and friends), and it needs exactly two outbound paths. ` +
		`"gateway" — "ref" is "fork" (run every branch), "join" (wait for branches to meet), "excl" (take the first matching branch) or "incl" (take every matching branch); no arguments and no results. ` +
		`"termination" — no ref, no arguments, no outbound path; ends the workflow. ` +
		`"error" — a single {"target":"message","type":"String","value":"..."} argument and no outbound path; fails the run. ` +
		`"error-handler" — no ref, one or two outbound paths. ` +
		`"delay" — exactly one argument, either {"target":"timestamp","type":"DateTime"} or {"target":"offset","type":"Duration"}. ` +
		`"prompt" — "ref" names a prompt construct such as "notification". ` +
		`"exec-workflow" — no ref, a required {"target":"workflow","type":"Handle"} (or type "ID") argument naming another workflow, and an optional {"target":"scope","type":"Vars"}. ` +
		`"break" and "continue" — inside an iterator body, nothing else set. "debug" — logs the whole scope. "visual" — a canvas note that never runs. ` +
		`"meta.visual.xywh" is [x, y, width, height] on the workflow editor's canvas; steps without it all land on (0,0) and cover each other, so set it if a person will ever open this workflow.`

	workflowPathsDoc = `JSON array of connections between steps: [{"parentID":"1","childID":"2","expr":"","meta":{"visual":{"id":"e1","parent":"1","points":[],"style":""}}}]. ` +
		`"parentID" and "childID" are the QUOTED stepIDs from 'steps'. Omit paths entirely for a single-step workflow — a step with no outbound path ends the run. ` +
		`ARRAY ORDER IS SEMANTICS, not presentation, wherever a step branches: the FIRST path out of an iterator is the loop body and the SECOND is the exit; the SECOND path out of an error-handler is the catch branch; ` +
		`an "excl" gateway tests the "expr" of each path leaving it in array order and takes the first that passes, so the fallback path — the one with an empty "expr" — must be last. ` +
		`"expr" is read only on paths leaving an "excl" or "incl" gateway; on any other path it is ignored.`

	// Both write tools repeat this. It is the single most important thing a
	// caller can be told about this resource: the write succeeding says nothing
	// about the workflow working.
	workflowIssuesDoc = `A BROKEN WORKFLOW IS STILL STORED. Only a missing name, an invalid or already-taken handle, a run-as user that cannot be loaded and a stale 'updatedAt' are refused. ` +
		`Everything else — an unknown function ref, an argument whose type does not match the parameter, a step kind given arguments it does not take, a path pointing at a step that does not exist, an unparseable expression — ` +
		`comes back in the "issues" array of this tool's result, and the workflow is written anyway. A workflow with issues never runs and its triggers are never registered, so a call that "succeeds" with a non-empty "issues" ` +
		`is a failure to fix, not a warning to note. Read "issues" before reporting success. `
)

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
			mcp.WithBoolean("includeDisabled", mcp.Description("Include disabled workflows in the listing. Off by default, matching the rest of the product. Naming one workflow by handle always finds it, enabled or not; this flag only affects the listing.")),
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
		mcp.NewTool("automation_workflow_create",
			mcp.WithDescription(
				"Create an automation workflow: a graph of steps joined by paths, which runs when a trigger "+
					"fires it or when automation_workflow_exec calls it. Workflows are the step-and-path "+
					"automations; Trigger-Action-Query automations are a different resource with their own "+
					"tools, and nothing here transfers to them. "+
					"The smallest workflow that runs is one step and no paths: steps "+
					"[{\"stepID\":\"1\",\"kind\":\"expressions\",\"arguments\":[{\"target\":\"result\",\"expr\":\"1 + 1\",\"type\":\"Integer\"}]}] "+
					"with no 'paths' at all. Execution starts at the one step nothing points at and follows the "+
					"paths out of it until it reaches a step with none. A workflow may have only ONE such "+
					"starting step, and this tool refuses a graph with more than one rather than storing "+
					"something that fails on its first run. "+
					workflowIssuesDoc+
					"Before you write a 'function' or 'iterator' step, fetch the construct catalogue with a REST "+
					"GET on /automation/functions/ — 93 entries, each with the 'ref' you put on the step, "+
					"its parameters, and the exact type names each parameter accepts (a parameter declares "+
					"'types' as a list; your argument's single 'type' has to be one of them, spelled the same). "+
					"There is no tool for this catalogue and a guessed ref stores a workflow that will not run. "+
					"NOTHING FIRES THIS WORKFLOW ON ITS OWN. Triggers are a separate resource: call "+
					"automation_trigger_create with this workflow's handle or ID afterwards, and give it 'stepID' "+
					"when you want it to start somewhere other than the workflow's only starting step. Until then "+
					"the workflow runs only when automation_workflow_exec is called. "+
					"Verify with automation_workflow_exec once created — storing and running are different "+
					"claims.",
			),
			mcp.WithString("handle", mcp.Required(), mcp.Description("Unique handle, e.g. \"order_sync\". Must start with a letter, be at least 2 characters, and use only letters, digits, underscores, dashes and dots. Handles are unique across the whole instance, not per project, and a taken one is refused.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Display name shown in the workflow list and editor. A workflow without one is refused.")),
			mcp.WithString("description", mcp.Description("What this workflow is for. Shown in the workflow editor.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the workflow may run. Defaults to true. A disabled workflow keeps its definition but never fires, and automation_workflow_exec reports it as not found. It is also left out of the automation_workflow_lookup listing unless includeDisabled is set — but it still resolves by handle, so it can be enabled again without its ID.")),
			mcp.WithString("steps", mcp.Description(workflowStepsDoc)),
			mcp.WithString("paths", mcp.Description(workflowPathsDoc)),
			mcp.WithString("scope", mcp.Description("JSON object of variables every run starts with, e.g. {\"retries\":3}. Input passed to automation_workflow_exec, and input carried by a trigger, is merged over this.")),
			mcp.WithString("runAs", mcp.Description("User ID as a string (to prevent precision loss) whose permissions the workflow's steps run with. Omit to run as whoever or whatever started it. Interval and timestamp triggers only register on a workflow that has one. A user that cannot be loaded fails the whole call.")),
			mcp.WithBoolean("trace", mcp.Description("Keep a full step-by-step trace of every run. Useful while building; it makes every run store more, so turn it off once the workflow works.")),
			mcp.WithBoolean("subWorkflow", mcp.Description("Mark this workflow as one that only other workflows call, through their 'exec-workflow' steps. Triggers on a sub-workflow are stored but never registered, so set this only when nothing should fire it directly.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create workflow",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_workflow_update",
			mcp.WithDescription(
				"Change an existing workflow: rename it, enable or disable it, or replace its step and path "+
					"graph. This is a read-modify-write — the workflow is loaded, only the arguments you pass "+
					"are applied, and its stored 'updatedAt' is echoed back so an edit someone else made in the "+
					"meantime is refused instead of being silently overwritten. "+
					"An argument you omit is left exactly as it is, so renaming or disabling a workflow is safe "+
					"to send on its own and needs no graph. But 'steps' and 'paths' REPLACE what is stored when "+
					"you do pass them — they are never merged, and passing an empty array [] is a request to "+
					"delete every step or every path, not a way to say \"unchanged\". Send the complete graph: "+
					"call automation_workflow_lookup with 'workflow' first and edit what it returns. Everything "+
					"you cannot set here — the owner, how many sessions are kept, labels — is preserved. "+
					"Changing the graph takes effect at once: the workflow is re-converted and re-registered on "+
					"the running instance, so the next event or the next automation_workflow_exec uses the new "+
					"steps. "+
					"A workflow may have only ONE step that nothing points at, and this tool refuses a graph "+
					"with more than one rather than storing something that fails on its first run. "+
					workflowIssuesDoc+
					"Fixing the issues and calling this tool again clears them. "+
					"Triggers are a separate resource and are not touched here: use automation_trigger_update to "+
					"change when the workflow fires, and check automation_trigger_lookup for this workflow after "+
					"you renumber steps, since a trigger pinned to a stepID that no longer exists fails at run "+
					"time.",
			),
			mcp.WithString("workflow", mcp.Required(), mcp.Description("Workflow ID as a string (to prevent precision loss) or handle. Find it with automation_workflow_lookup; a disabled workflow resolves by handle too, though it is left out of the listing unless includeDisabled is set.")),
			mcp.WithString("handle", mcp.Description("New handle. Must be unique across the instance. An empty string clears it, after which the workflow can only be referenced by ID.")),
			mcp.WithString("name", mcp.Description("New display name. Cannot be emptied — a workflow with no name is refused.")),
			mcp.WithString("description", mcp.Description("New description. An empty string clears it.")),
			mcp.WithBoolean("enabled", mcp.Description("Enable or disable the workflow. The definition is kept either way. A disabled workflow is left out of the automation_workflow_lookup listing unless includeDisabled is set, but it still resolves by handle, so this is reversible with the handle alone.")),
			mcp.WithString("steps", mcp.Description("Replaces the stored steps wholesale. "+workflowStepsDoc)),
			mcp.WithString("paths", mcp.Description("Replaces the stored paths wholesale. "+workflowPathsDoc)),
			mcp.WithString("scope", mcp.Description("JSON object of variables every run starts with. Replaces the stored set wholesale; {} clears it.")),
			mcp.WithString("runAs", mcp.Description("User ID as a string (to prevent precision loss) whose permissions the steps run with. An empty string clears it, and interval or timestamp triggers on this workflow then stop registering.")),
			mcp.WithBoolean("trace", mcp.Description("Keep a full step-by-step trace of every run.")),
			mcp.WithBoolean("subWorkflow", mcp.Description("Mark the workflow as one that only other workflows call. Turning this on unregisters its triggers.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update workflow",
		h.update,
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
					"still resolves by handle, a deleted one is not findable by handle at all and restoring it needs "+
					"an ID you must have kept.",
			),
			mcp.WithString("workflow", mcp.Required(), mcp.Description("Workflow ID as a string (to prevent precision loss) or handle. Find it with automation_workflow_lookup.")),
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
