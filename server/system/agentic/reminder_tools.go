package agentic

import (
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

// Reminder is the reference implementation for MCP tool families. See
// mcp/CONVENTIONS.md. It was chosen as the exemplar because it exercises the
// awkward cases most resources hit one at a time: soft delete with no undelete,
// domain operations beyond CRUD, no handle (IDs only), ownership-scoped reads,
// a free-form JSON payload, and time-range filters.
//
// This file holds declarations only. Implementations are in reminder_handler.go,
// in the same order.

func (h *reminderHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_lookup",
			mcp.WithDescription(
				"Look up a reminder by ID, or list and filter reminders. Provide 'reminderID' to fetch one, "+
					"which returns the full reminder including its payload; omit it to list, which returns a "+
					"compact form (ID, resource, remindAt, assignedTo, dismissedAt) — fetch by ID when you need "+
					"the payload. "+
					"You only ever see reminders assigned to you: the service filters by assignee regardless of "+
					"what you pass in 'assignedTo', so this tool cannot read another user's reminders. "+
					"Dismissed reminders are included unless you set 'excludeDismissed'.",
			),
			mcp.WithString("reminderID", mcp.Description("Reminder ID as a string, to prevent precision loss. Omit to list instead.")),
			mcp.WithString("resource", mcp.Description("Filter by the resource string the reminder is attached to, e.g. \"compose:record\".")),
			mcp.WithString("assignedTo", mcp.Description("Assignee user ID as a string. Reads are scoped to you regardless; this narrows, it cannot widen.")),
			mcp.WithString("scheduledFrom", mcp.Description("Only reminders due at or after this time. RFC3339, e.g. 2026-08-03T09:00:00Z.")),
			mcp.WithString("scheduledUntil", mcp.Description("Only reminders due at or before this time. RFC3339.")),
			mcp.WithString("excludeDismissed", mcp.Description("Set \"true\" to omit reminders that have been dismissed.")),
			mcp.WithString("includeDeleted", mcp.Description("Set \"true\" to include soft-deleted reminders. Deletes are soft and there is no undelete tool, so this is the only way to see them.")),
			mcp.WithString("scheduledOnly", mcp.Description("Set \"true\" to return only reminders that have a remindAt set.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup reminder",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_create",
			mcp.WithDescription(
				"Create a reminder. 'resource' identifies what the reminder is about and 'payload' carries "+
					"whatever the consumer needs to render it — both are free-form, so mirror an existing "+
					"reminder's shape rather than inventing one; call system_reminder_lookup first if unsure. "+
					"Assigning to a user other than yourself requires the assign-reminder permission and fails "+
					"without it.",
			),
			mcp.WithString("resource", mcp.Required(), mcp.Description("Resource the reminder refers to, e.g. \"compose:record/1/2/3\".")),
			mcp.WithString("payload", mcp.Description("JSON object with the reminder's contents.")),
			mcp.WithString("remindAt", mcp.Description("When to surface the reminder. RFC3339. Omit for a reminder with no schedule.")),
			mcp.WithString("assignedTo", mcp.Description("Assignee user ID as a string. Defaults to you; assigning to anyone else needs permission.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Create reminder",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_update",
			mcp.WithDescription(
				"Update a reminder. Omit a field to leave it unchanged; pass it empty to clear it. "+
					"'payload' replaces the stored object wholesale rather than merging into it. "+
					"To change only the schedule, prefer system_reminder_snooze — it also records that the "+
					"reminder was snoozed, which this tool does not.",
			),
			mcp.WithString("reminderID", mcp.Required(), mcp.Description("Reminder ID as a string.")),
			mcp.WithString("resource", mcp.Description("New resource string. Empty clears it.")),
			mcp.WithString("payload", mcp.Description("Replacement JSON object. Empty clears it. Not merged.")),
			mcp.WithString("remindAt", mcp.Description("New due time, RFC3339. Empty clears the schedule.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Update reminder",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_delete",
			mcp.WithDescription(
				"Delete a reminder. The delete is soft — the record is retained and can be seen again by "+
					"passing includeDeleted to system_reminder_lookup — but there is no undelete tool, because "+
					"the reminder service exposes no undelete operation. Treat it as final. "+
					"If you only want to stop a reminder surfacing, use system_reminder_dismiss instead: that "+
					"is reversible.",
			),
			mcp.WithString("reminderID", mcp.Required(), mcp.Description("Reminder ID as a string.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete reminder",
		h.delete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_dismiss",
			mcp.WithDescription(
				"Dismiss a reminder so it stops surfacing, keeping the record. Reversible with "+
					"system_reminder_undismiss. Prefer this over system_reminder_delete when the user is done "+
					"with a reminder rather than wanting it gone.",
			),
			mcp.WithString("reminderID", mcp.Required(), mcp.Description("Reminder ID as a string.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Dismiss reminder",
		h.dismiss,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_undismiss",
			mcp.WithDescription("Reverse a dismissal, so the reminder surfaces again."),
			mcp.WithString("reminderID", mcp.Required(), mcp.Description("Reminder ID as a string.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undismiss reminder",
		h.undismiss,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_reminder_snooze",
			mcp.WithDescription(
				"Push a reminder's due time out and record that it was snoozed, incrementing its snooze count. "+
					"Use this rather than system_reminder_update when the user is deferring a reminder — the "+
					"count is what tells you a reminder keeps being put off.",
			),
			mcp.WithString("reminderID", mcp.Required(), mcp.Description("Reminder ID as a string.")),
			mcp.WithString("remindAt", mcp.Required(), mcp.Description("New due time, RFC3339.")),
			hmcp.InGroup(hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Snooze reminder",
		h.snooze,
	)
}
