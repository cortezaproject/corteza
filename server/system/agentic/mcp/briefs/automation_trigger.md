# Brief: `automation_trigger`

Tool file: `server/automation/agentic/trigger_tools.go`
Handler:   `server/automation/agentic/trigger_handler.go`
Group: `configuring` · Package: `automation/agentic` (reuse the existing
`toolRegistrar` in `taq_handler.go`; do not redeclare it)

## 1. Service

`automation/service/trigger.go` + `trigger.gen.go` —
`autoService.DefaultTrigger`.

| Method | Signature | Becomes |
|---|---|---|
| `FindByID` | `(ctx, triggerID uint64)` | single-fetch in `lookup` |
| `Search` | `(ctx, types.TriggerFilter)` | list mode in `lookup` |
| `Create` | `(ctx, *types.Trigger)` | `create` |
| `Update` | `(ctx, *types.Trigger)` | **blocked — see §2** |
| `DeleteByID` | `(ctx, ID uint64)` | **blocked — see §2** |
| `UndeleteByID` | `(ctx, ID uint64)` | `undelete` |
| `SearchOnManual` | `(ctx, workflowID, stepID uint64)` | skip — an internal lookup for the manual-run path, not a configuration op |

Risks: `lookup` read; `create`, `undelete` write.

## 2. Authorization (§8.6) — PARTIAL PASS. Two tools are blocked.

Verified 2026-08-04, and the picture is inconsistent:

| Op | Check | Verdict |
|---|---|---|
| `FindByID` (`trigger.go:146`) | `CanSearchTriggers` | ✅ |
| `Search` | `CanSearchTriggers` | ✅ |
| `onCreate` (`trigger.go:192`) | `CanManageTriggersOnWorkflow` | ✅ |
| `onUndelete` (`trigger.go:315`) | `CanManageTriggersOnWorkflow` | ✅ |
| `onUpdate` (`trigger.go:227`) | **none** | ❌ |
| `onDelete` (`trigger.go:285`) | **none** | ❌ |

`onUpdate` and `onDelete` have no `ac.Can*` call and no identity scoping. The
generated wrapper calls `svc.guard` → `guardProjectWritable`, which asks whether
the *project* is writable — not whether this caller may manage triggers. So a
caller who can write anything in the project can retarget or delete any trigger
in it, without ever holding `CanManageTriggersOnWorkflow`.

**Do not write `automation_trigger_update` or `automation_trigger_delete`.**
File both as §8.6 gaps. That is the rule: a service method with neither an
`ac.Can*` call nor an ownership scope does not get a tool.

Note what this means for the resource as a whole: a trigger can be created and
restored through MCP but not changed or removed. Say so plainly in the `create`
description rather than leaving a caller to discover it — and do not offer
delete-then-recreate as a workaround, because delete is the blocked op.

## 3. Identifier strategy

**No `FindByAny`, no `FindByHandle`** — triggers have no handle. They are
identified by ID, or found by the workflow and event they belong to.

`lookup` takes an optional `triggerID` for single-fetch. There is nothing else
to resolve against, so unlike most resources there is no name-or-handle ref.
`undelete` takes a required `triggerID`, and its description must say where the
ID comes from — `TriggerFilter.Deleted` defaults to `StateExcluded`, so verify
whether a deleted trigger is reachable through `lookup` with `includeDeleted`
before claiming it is.

## 4. Filter fields → params

`TriggerFilter` (`automation/types/trigger.go:10`):

| Field | Param? | Note |
|---|---|---|
| `WorkflowID []string` | yes, as `workflow` | accept a workflow ref and resolve it via the workflow service — the most useful way to ask "what fires this workflow?" |
| `EventType` | yes | pair it with `automation_event_type_lookup` in the description; the valid values come from there |
| `ResourceType` | yes | same |
| `Deleted` | yes, `includeDeleted` boolean | |
| `Disabled` | yes, `includeDisabled` boolean | mirrors `automation_taq_lookup` |
| `TriggerID []string` | **no** | the ref param covers single fetch |
| `LabeledIDs`, `Labels` | **no** | internal |

**Verify each of these reaches SQL before exposing it.** Three filter fields in
this codebase have turned out to be declared and never translated —
`RoleFilter.MemberID`, `UserGroupFilter.MemberID` and `Name`,
`ApplicationFilter.ApplicationID`. Check `store/adapters/rdbms/filters.gen.go`
for the trigger filter and drop any param the store ignores; a filter that
silently returns everything is the `discovery_search` failure again.

Plus `limit` and `pageCursor` (§8.1).

## 5. List projection

`{triggerID, workflowID, eventType, resourceType, enabled, deletedAt}`. Leave
out `Constraints` and `Meta` — the constraint set is the heavy part and is
incidental when picking one trigger from a list.

## 6. Traps

**A trigger is not standalone.** It binds a workflow to an event, so creating
one changes when that workflow runs. The `create` description must say which
workflow it will fire and that the workflow starts responding to the event
immediately.

**Registration is live.** `registerTriggers` / `unregisterTriggers`
(`trigger.go:426`, `:533`) mutate an in-memory registry as triggers change.
Check whether `Create` and `UndeleteByID` re-register, and whether a trigger
created through the service takes effect without a restart. If it does not, the
description must say so — a tool that appears to work and silently does nothing
until a restart is worse than no tool.

**`onManual` is special.** `SearchOnManual` exists because manual triggers are
looked up by workflow and step rather than by event. Do not fold it into
`lookup`.

**Constraints are the expressive part.** A trigger's constraints decide whether
it actually fires for a given event. Exposing them on `create` means a caller
can write a constraint that silently never matches. Consider omitting them from
the first version and saying that constrained triggers are configured in the
webapp — flag this rather than deciding it alone.
