package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Event types are a reference catalogue, not a resource: they are defined by
// Human's own code, identical on every instance, and there is nothing to create,
// change or delete. Hence one read tool and no CRUD.
//
// It sits in the configuring group because it is read while configuring
// automation — it is the lookup table for automation_trigger_create, which needs
// a valid eventType/resourceType pair and has no other way to discover one.
//
// It declares no limit and no pageCursor. The catalogue is a fixed in-memory
// slice with no store behind it, so a cursor would advertise paging that cannot
// work; the whole list is well inside the tool-result size ceiling, and
// 'resourceType' narrows it for a caller who does not want all of it.
//
// This file holds declarations only. The implementation is in
// event_type_handler.go.

func (h *eventTypeHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("automation_event_type_lookup",
			mcp.WithDescription(
				"List the automation events Human can react to: every valid eventType and resourceType pair, "+
					"with the properties each event carries and the constraints a trigger may narrow it by. "+
					"This is a fixed reference list built into Human — read-only, the same on every instance, "+
					"and not something a caller can add to. "+
					"Read it before calling automation_trigger_create or automation_trigger_update: those need "+
					"an eventType and resourceType that exist together, this is the only place they are written "+
					"down, and a pair that is not in this list is stored without complaint and then never "+
					"fires. It is also how you find what a workflow will receive when the event fires — the "+
					"'properties' of an entry are the variables put into the workflow's scope — and which "+
					"constraint names the workflow editor will accept for that event. "+
					"The whole catalogue is around 120 entries and is returned in full when you pass nothing, "+
					"which is usually what you want. Pass 'resourceType' to narrow it: an exact value like "+
					"\"compose:record\" returns just that resource's events, while a bare prefix like "+
					"\"compose\" returns that resource and all of its sub-resources. There is no paging and no "+
					"per-entry fetch — one call returns everything there is to know.",
			),
			mcp.WithString("resourceType", mcp.Description("Narrow to one resource, e.g. \"compose:record\" for exactly that one or \"compose\" for it and every compose sub-resource. Case-insensitive. Omit for the whole catalogue.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup automation event types",
		h.lookup,
	)
}
