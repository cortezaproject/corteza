# Brief: `automation_event_type`

Tool file: `server/automation/agentic/event_type_tools.go`
Handler:   `server/automation/agentic/event_type_handler.go`
Group: `configuring` · Package: `automation/agentic`

**One tool: `automation_event_type_lookup`. Risk `read`. There is nothing else
to write.**

## 1. This is not a service

Unlike every other resource briefed so far, `eventTypes` has **no service, no
store and no CRUD**. It is a static catalogue generated into
`automation/rest/eventTypes.gen.go` as `getEventTypeDefinitions()`, exposed by
one REST handler:

```go
func (ctrl EventTypes) List(_ context.Context, _ *request.EventTypesList) (interface{}, error)
```

Both parameters are ignored. There is no filter, no paging, no identity.

Consequences, all of which the tool must respect rather than paper over:

- No `create`, `update`, `delete` or `undelete`. Event types are defined by
  Human's own code; a caller cannot add one. Do not write these tools and do not
  hint in the description that new event types can be registered.
- No §8.6 question to answer. The data is a compile-time constant, identical
  for every caller, and reveals nothing about the instance. It is
  authorization-neutral, which is why a read tool is safe despite the service
  having no checks — the *reason* matters, so record it rather than just
  asserting "passes".
- No `limit` or `pageCursor`. §8.1's paging rule assumes a store that can page;
  this is a fixed in-memory slice. Declaring a cursor would advertise paging
  that cannot work, which §8.1 forbids. Bound the result with a `resourceType`
  filter instead, and say in the description that the list is complete and
  small.

Confirm the size of `getEventTypeDefinitions()` before deciding whether even a
filter is needed — if the whole catalogue fits comfortably under
`toolkit.JSONResult`'s ceiling, returning it whole is simpler and more useful
than making a caller guess a filter value.

## 2. Shape

`eventTypeDef` (`automation/rest/eventTypes.go:20`):

```go
ResourceType string                   // "compose", "system", …
EventType    string                   // "onManual", "onInterval", "onTimestamp", …
Properties   []eventTypePropertyDef   // what the event carries
Constraints  []eventTypeConstraintDef // what a trigger may constrain on
```

Return all four. `Properties` and `Constraints` look like the heavy fields §8.1
would project away, but here they are the entire point: a caller reads this tool
precisely to learn what a trigger can constrain on and what an automation will
receive. There is no single-item mode to fetch them from later.

## 3. Where the data lives

`getEventTypeDefinitions()` is in the `rest` package, not `service`. Importing
`automation/rest` from `automation/agentic` may not be appropriate — check the
direction before doing it. If it is wrong, the honest options are to move the
generator somewhere both can reach or to call the REST controller; **do not
duplicate the list into the tool**, which would rot silently the first time the
generator output changed.

Report which you found and what you chose.

## 4. Why this tool matters

It is the lookup table for `automation_trigger`. A caller creating a trigger
needs a valid `eventType` and `resourceType` pair, and there is no other way to
discover them — they are not in any other tool's schema and not enumerable from
a filter. The two descriptions should reference each other explicitly.

## 5. Traps

**Do not describe it as configuration.** It reads as a configuring-group tool
because it is used while configuring, but nothing about it is editable. The
description should make clear it is a reference list.

**Name is singular per RESOURCES.md**: `automation_event_type`, from a REST
resource called `eventTypes`. That is the singularisation rule, not an error.
