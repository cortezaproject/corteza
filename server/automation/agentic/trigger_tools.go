package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// A trigger binds one workflow to one event: when an event of the trigger's
// eventType fires on its resourceType, the workflow runs. A trigger is therefore
// never standalone — creating, changing or removing one changes when its
// workflow runs.
//
// The eventType/resourceType pairs a trigger may use are a fixed catalogue in
// Human's own code; automation_event_type_lookup is the only way to discover
// them, which is why every description here points at it.
//
// This file holds declarations only. Implementations are in trigger_handler.go,
// in the same order.

func (h *triggerHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("automation_trigger_lookup",
			mcp.WithDescription(
				"List automation triggers, or fetch one whole. A trigger binds a workflow to an event: when "+
					"an event of its eventType fires on its resourceType, the workflow runs. "+
					"Omit 'triggerID' to list: a list row is a compact form — triggerID, workflowID, eventType, "+
					"resourceType, enabled, deletedAt — enough to pick a trigger out of a list. Pass 'triggerID' "+
					"to fetch that one trigger in full, which additionally returns its constraints (the "+
					"conditions that decide whether it actually fires), its fixed input and its metadata; those "+
					"are left out of list rows because they are bulky and rarely what tells two rows apart. "+
					"Narrow the list with 'workflow' to answer \"what makes this workflow run?\", or with "+
					"'eventType' and 'resourceType' to answer \"what runs when this happens?\". Both match "+
					"exactly, not as substrings, and the values that exist are listed by "+
					"automation_event_type_lookup — guessing them returns an empty list rather than an error. "+
					"Disabled triggers are omitted unless you set 'includeDisabled', deleted ones unless you set "+
					"'includeDeleted'. Fetching by 'triggerID' ignores both flags and returns the trigger either "+
					"way, so this is also how you inspect a deleted trigger before restoring it with "+
					"automation_trigger_undelete.",
			),
			mcp.WithString("triggerID", mcp.Description("Trigger ID as a string (to prevent precision loss). Omit to list instead. Triggers have no handle; they are identified by ID only.")),
			mcp.WithString("workflow", mcp.Description("Workflow ID as a string (to prevent precision loss), or workflow handle. Lists only that workflow's triggers. Ignored when 'triggerID' is given.")),
			mcp.WithString("eventType", mcp.Description("Exact event type to match, e.g. \"onManual\" or \"afterUpdate\". Valid values come from automation_event_type_lookup.")),
			mcp.WithString("resourceType", mcp.Description("Exact resource type to match, e.g. \"compose:record\". Valid values come from automation_event_type_lookup.")),
			mcp.WithBoolean("includeDisabled", mcp.Description("Also list disabled triggers, which stay stored but never fire.")),
			mcp.WithBoolean("includeDeleted", mcp.Description("Also list deleted triggers. This is how you find the ID for automation_trigger_undelete.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup automation trigger",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_trigger_create",
			mcp.WithDescription(
				"Create a trigger, making a workflow run whenever a given event fires. This changes when the "+
					"named workflow runs, so only create one when you intend that workflow to start responding "+
					"to the event. It takes effect immediately: the trigger is registered as this call returns "+
					"and the workflow can fire on the very next matching event, with no restart and no separate "+
					"activation step. "+
					"'eventType' and 'resourceType' must be a pair that actually exists — call "+
					"automation_event_type_lookup for the catalogue and copy a pair from it. A pair that does "+
					"not exist is stored without complaint and simply never fires. "+
					"A trigger created here carries NO constraints, so it fires for every event of that type on "+
					"that resource type: an afterUpdate trigger on compose:record runs the workflow for every "+
					"record update in every module of every namespace. Constraints, which narrow a trigger to a "+
					"particular namespace, module or field value, cannot be set through this tool — create "+
					"constrained triggers in the workflow editor in the Human webapp instead. "+
					"Two cases are stored but silently never registered, so nothing fires and no error is "+
					"returned: an onInterval or onTimestamp trigger (and system:sink and system:queue events) on "+
					"a workflow with no run-as user, and any enabled trigger on a workflow marked as a "+
					"sub-workflow. Check the workflow with automation_workflow_lookup first if either might "+
					"apply. A disabled workflow, or one with unresolved issues, likewise registers nothing.",
			),
			mcp.WithString("workflow", mcp.Required(), mcp.Description("Workflow to run, as an ID string (to prevent precision loss) or a handle. This is the workflow that will start responding to the event.")),
			mcp.WithString("eventType", mcp.Required(), mcp.Description("Event that fires the workflow, e.g. \"onManual\" or \"afterCreate\". Must be a value automation_event_type_lookup reports for this resourceType.")),
			mcp.WithString("resourceType", mcp.Required(), mcp.Description("Resource the event happens on, e.g. \"compose:record\" or \"system:user\". Must be paired with 'eventType' as automation_event_type_lookup reports it.")),
			mcp.WithString("stepID", mcp.Description("ID as a string of the workflow step to start at. Omit only when the workflow has exactly one starting step; otherwise starting it fails at run time.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the trigger fires. Defaults to true. A disabled trigger is stored but never registered.")),
			mcp.WithString("description", mcp.Description("Human-readable note on what this trigger is for. Shown in the workflow editor.")),
			mcp.WithString("input", mcp.Description("JSON object of fixed input variables merged into the workflow's scope on every run. Omit when the event's own properties are all the workflow needs.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Create automation trigger",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_trigger_update",
			mcp.WithDescription(
				"Change an existing trigger: re-point it at another workflow or event, enable or disable it, or "+
					"edit its note and fixed input. Omit a field to leave it unchanged; pass 'description' or "+
					"'input' empty to clear it. "+
					"'workflow', 'eventType' and 'resourceType' can be changed but not cleared — an empty value "+
					"is refused rather than stored, because a trigger with no workflow can never be loaded again "+
					"and so becomes impossible to edit, delete or restore. "+
					"Constraints and labels are preserved untouched; neither can be set here. New "+
					"'eventType'/'resourceType' values must be a pair automation_event_type_lookup reports, and "+
					"an invalid pair is stored without complaint and simply never fires. "+
					"The change takes effect immediately: the trigger is re-registered on the running "+
					"instance, so a trigger you disable here stops firing at once and one you re-point starts "+
					"firing on its new event.",
			),
			mcp.WithString("triggerID", mcp.Required(), mcp.Description("Trigger ID as a string (to prevent precision loss). Find it with automation_trigger_lookup.")),
			mcp.WithString("workflow", mcp.Description("Move the trigger to this workflow: an ID string or a handle. Cannot be cleared.")),
			mcp.WithString("eventType", mcp.Description("New event type, from automation_event_type_lookup. Cannot be cleared.")),
			mcp.WithString("resourceType", mcp.Description("New resource type, from automation_event_type_lookup. Cannot be cleared.")),
			mcp.WithString("stepID", mcp.Description("ID as a string of the workflow step to start at. Empty resets it, which requires the workflow to have exactly one starting step.")),
			mcp.WithBoolean("enabled", mcp.Description("Enable or disable the trigger. See the note above: disabling does not stop it firing until the workflow is saved again or the server restarts.")),
			mcp.WithString("description", mcp.Description("New note. Empty clears it.")),
			mcp.WithString("input", mcp.Description("Replacement JSON object of fixed input variables. Replaces the stored set wholesale rather than merging. Empty clears it.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Update automation trigger",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_trigger_delete",
			mcp.WithDescription(
				"Delete a trigger, so its event stops starting its workflow. The workflow itself is untouched "+
					"and still exists: deleting every trigger only means nothing starts it automatically, and it "+
					"can still be run on demand with automation_workflow_exec. "+
					"The delete is soft and reversible — automation_trigger_undelete restores it, and "+
					"automation_trigger_lookup with 'includeDeleted' still lists it. Prefer disabling the "+
					"trigger with automation_trigger_update when you only want to pause it, since that keeps it "+
					"visible in the workflow editor. "+
					"The trigger stops firing immediately — it is unregistered from the running instance as "+
					"well as marked deleted — and automation_trigger_undelete brings it back.",
			),
			mcp.WithString("triggerID", mcp.Required(), mcp.Description("Trigger ID as a string (to prevent precision loss). Find it with automation_trigger_lookup.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete automation trigger",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("automation_trigger_undelete",
			mcp.WithDescription(
				"Restore a deleted trigger, putting its workflow back on the event it used to respond to. The "+
					"delete was soft, so the event type, resource type, constraints and input all come back as "+
					"they were. "+
					"Find the ID with automation_trigger_lookup: unlike most resources here, deleted triggers "+
					"are reachable — pass 'includeDeleted' to list them, and passing a deleted trigger's ID as "+
					"'triggerID' returns it in full with 'deletedAt' set, so you can check what you are about to "+
					"restore. Restoring a trigger that is not deleted succeeds and changes nothing. "+
					"The trigger starts firing again immediately: it is re-registered on the running instance "+
					"as part of the restore.",
			),
			mcp.WithString("triggerID", mcp.Required(), mcp.Description("Trigger ID as a string (to prevent precision loss), from automation_trigger_lookup with includeDeleted set.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete automation trigger",
		h.undelete,
	)
}
