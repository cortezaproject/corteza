package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// TAQ is Trigger Action Query. The service and REST resource behind these tools
// is ngAutomation; the tool surface never says so — see mcp/RESOURCES.md.
//
// This file holds declarations only. Implementations are in taq_handler.go, in
// the same order.

// The three graph params below carry the resource's own JSON, verbatim — the
// same shape automation_taq_lookup returns and the REST endpoint accepts — so a
// lookup result can be edited and sent straight back. A flattened tool-specific
// shape would have to be translated in both directions and would break that
// round trip.
//
// Checked against server/automation/types/ng_automation.gen.go (the wire types),
// server/automation/service/ng_automation_converter.go (what the runtime does
// with them) and the running construct library.

const taqTriggersDoc = `JSON array of triggers — what starts the TAQ. One trigger:
{"triggerID":"1001","handle":"onAgent","enabled":true,"resourceType":"automation:trigger:agentic","eventType":"onAgentic","meta":{"short":"Invoked by an agent"},"inputSchema":[{"name":"subject","type":"String","required":true}]}
resourceType and eventType are a pair from the trigger catalogue: GET /automation/construct-library/triggers, 22 entries. "automation:trigger:agentic" with "onAgentic" is the pair automation_taq_exec runs, so include one of those if you want to be able to run the TAQ on demand; the compose:record and system:user pairs fire on real events instead.
triggerID is a string of digits. Omit it on create and one is minted for you and returned. Supply it when a path references the trigger, and echo the existing triggerID on update — a trigger sent without one is replaced by a new trigger with a new ID, orphaning any path that pointed at the old one.
handle names the trigger: it is what automation_taq_exec takes as entryPoint, and what a step reads from as a scope.
inputSchema declares the parameters automation_taq_exec accepts under "input"; it is the only place those names are written down.`

const taqStepsDoc = `JSON array of steps — what the TAQ does. One step:
{"stepID":"1","handle":"notify","kind":"function","ref":"notificationSend","meta":{"short":"Notify the owner"},"arguments":[{"argumentName":"recipient","value":"507566326668132353","type":"ID"},{"argumentName":"title","value":"Hello","type":"String"}]}
stepID is a string of digits you choose, unique across steps AND triggers; paths reference steps by it.
kind is one of: function, iterator, gatewayExclusive, gatewayInclusive, termination, error. Those are the only six, and they are NOT the workflow step kinds — a TAQ has no "expressions" step. A kind outside the six is rejected before anything is written.
ref, on a function or an iterator, names a construct-library function. List them with GET /automation/construct-library/functions — 17 entries, each with its parameters, their types and whether they are required. Do NOT use GET /automation/functions/: that is the workflow function registry, and 78 of its 93 entries do not exist here.
arguments bind by argumentName, which must equal one of that function's parameter argumentNames exactly; every required parameter must be present. "type" is compared literally against that parameter's "types" array and must be spelled as listed there (ID, String, Integer, Boolean, ComposeRecord, …) — omitting it is the same as sending "" and fails. Supply a literal with "value", copy an earlier step's output with "source" plus "scope" (the producing step's handle), or compute one with "expr".
A step's results cannot be set: they are derived from the function definition and anything you send is overwritten. An error step's message likewise cannot be set through the API.
A termination step is optional — every leaf step is wired to an auto-injected one.`

const taqPathsDoc = `JSON array of edges: [{"parentID":"<trigger or step ID>","childID":"<step ID>"}]. IDs are strings of digits.
Send [] for the common case. With exactly one unconnected trigger and exactly one unconnected step, the two are wired together for you; more than one of either and you get a graph.ambiguousEntry issue instead.
Order carries meaning, not presentation. A non-gateway step's second outbound path is its error handler rather than a second successor; an iterator's paths are its body first, then its exit.
A gatewayExclusive or gatewayInclusive step needs at least two outbound paths: each carries a "condition", except exactly one which carries none and is the else branch, always tested last. A condition is an AST object, not an expression string: {"ref":"eq","args":[{"symbol":"foo"},{"value":{"@type":"String","@value":"bar"}}]}. Operators: and, or, not, isNull, isNotNull, eq, ne, lt, gt, lte, gte. Nothing checks a condition when it is written — a wrong operator, symbol or scope surfaces only when the TAQ runs.`

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
		mcp.NewTool("automation_taq_create",
			mcp.WithDescription(
				"Create a TAQ (Trigger Action Query automation) — its triggers, steps and paths in one "+
					"call. A TAQ is a graph: triggers say what starts it, steps say what it does, paths "+
					"connect them. "+
					"Find out what a step can do before writing one. GET /automation/construct-library/functions "+
					"lists the 17 functions a step's 'ref' may name, each with its parameters, their exact "+
					"type names and whether they are required; GET /automation/construct-library/triggers lists "+
					"the 22 resourceType/eventType pairs a trigger may use. Do not use GET /automation/functions/ "+
					"— that is the workflow function registry and 78 of its 93 entries are unavailable to a TAQ. "+
					"The smallest TAQ that runs is one trigger, one function step, and no paths at all: a lone "+
					"unconnected trigger and a lone unconnected step are wired together for you, and termination "+
					"is added automatically. "+
					"Read the result, do not assume it worked. A TAQ whose graph fails validation is still "+
					"stored and still succeeds, and says what is wrong under 'issues' — a code, a severity, and "+
					"the exact step, parameter and legal types. Any issue at all, of any severity, stops the TAQ "+
					"being registered with the runtime: 'runnable' comes back false and automation_taq_exec "+
					"answers 'manager: executable not found', which means the TAQ has issues and not that the ID "+
					"is wrong. No 'issues' key at all is the clean result. Fix issues with automation_taq_update. "+
					"'enabled' defaults to true, which arms the TAQ immediately — its triggers begin listening "+
					"for real events. Pass false to author without arming, but a disabled TAQ also refuses "+
					"automation_taq_exec, so you cannot test one until you enable it. "+
					"Storing is not proof of working. Run it with automation_taq_exec and then read "+
					"automation_taq_execution_trace, because exec reports status only: a 'completed' status with "+
					"no step frames in the trace — nothing, or only the trigger frame — is a failure, it means "+
					"no step ran. A frame with a populated 'args' proves the arguments bound and nothing more — "+
					"where the step has a visible effect, check for that effect separately.",
			),
			mcp.WithString("handle", mcp.Required(), mcp.Description("URL-friendly identifier, unique among TAQs in the same project. Use snake_case (lowercase letters, digits, underscores); a hyphen is the subtraction operator wherever an identifier is parsed. This is what automation_taq_lookup searches and what automation_taq_exec resolves, so a TAQ without one can only ever be reached by its numeric ID.")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Human-readable name, shown in listings.")),
			mcp.WithString("description", mcp.Description("What this TAQ is for.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the TAQ is live. Defaults to true. False leaves its triggers unregistered AND makes automation_taq_exec refuse it, so it cannot be tested while disabled.")),
			mcp.WithString("triggers", mcp.Description(taqTriggersDoc)),
			mcp.WithString("steps", mcp.Description(taqStepsDoc)),
			mcp.WithString("paths", mcp.Description(taqPathsDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create TAQ",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_update",
			mcp.WithDescription(
				"Update a TAQ (Trigger Action Query automation): rename it, enable or disable it, or replace "+
					"its trigger, step and path graph. This is also the tool for fixing the 'issues' "+
					"automation_taq_create reported. "+
					"A field you omit is left alone, so a rename does not disturb the graph and a graph change "+
					"does not disturb the name. 'triggers', 'steps' and 'paths' each replace their whole "+
					"collection rather than merging into it: send the complete array, not only the part you "+
					"changed, and send [] only when you mean to remove every one. A TAQ with no steps stores "+
					"clean and then executes as 'completed' having run nothing. "+
					"Read the TAQ first with automation_taq_lookup — it returns exactly the shape these params "+
					"accept, so the safe edit is lookup, change, send back. Echo each trigger's existing "+
					"triggerID when you resend triggers, or the trigger is replaced by a new one with a new ID "+
					"and any path pointing at the old one is orphaned. "+
					"Anything these params do not carry is preserved from the stored TAQ: run-as user, owner, "+
					"labels and scope are never cleared by omitting them. "+
					"The write succeeds even when the graph is invalid. 'issues' says what is wrong, and any "+
					"issue at all unregisters the TAQ from the runtime — 'runnable' comes back false and "+
					"automation_taq_exec answers 'manager: executable not found' until the graph is clean. "+
					"Disabling stops the triggers firing and also makes automation_taq_exec refuse the TAQ. To "+
					"take one out of service reversibly, prefer enabled:false over automation_taq_delete — a "+
					"deleted TAQ can no longer be found by handle at all.",
			),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID as a string (to prevent precision loss), or handle. Find it with automation_taq_lookup.")),
			mcp.WithString("handle", mcp.Description("New handle. Omit to leave it alone; an empty string clears it, after which the TAQ is reachable only by its numeric ID.")),
			mcp.WithString("name", mcp.Description("New name. Omit to leave it alone; it cannot be cleared, because a TAQ must have a name.")),
			mcp.WithString("description", mcp.Description("New description. Omit to leave it alone; an empty string clears it.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the TAQ is live. Omit to leave it alone. False unregisters its triggers AND makes automation_taq_exec refuse it.")),
			mcp.WithString("triggers", mcp.Description("Replaces every trigger. Omit to leave the triggers alone. "+taqTriggersDoc)),
			mcp.WithString("steps", mcp.Description("Replaces every step. Omit to leave the steps alone. "+taqStepsDoc)),
			mcp.WithString("paths", mcp.Description("Replaces every path. Omit to leave the paths alone. "+taqPathsDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update TAQ",
		h.update,
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
