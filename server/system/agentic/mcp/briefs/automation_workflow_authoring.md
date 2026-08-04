# Brief: `automation_workflow` create and update

Extends `server/automation/agentic/workflow_tools.go` / `workflow_handler.go`,
which already carry `lookup`, `exec`, `delete` and `undelete`.

Group `configuring`. `create` and `update` are both risk `write`.

Distilled from three independent investigations (webapp payload, server
validation, empirical probe). **Workflow and TAQ are not variants of each
other** — different endpoints, types, path-condition representation, gateway
vocabulary and trigger model. Do not generalise between them; see
`automation_taq_authoring.md` for the other one.

## 1. It works, and it worked first time

The probe created and ran one in **two calls, zero failures**:

```json
POST /automation/workflows/
{"handle":"agent-probe-wf","meta":{"name":"Agent Probe WF"},"enabled":true,
 "steps":[{"stepID":"1","kind":"expressions",
   "arguments":[{"target":"result","expr":"1 + 1","type":"Integer"}]}],
 "paths":[]}
```

then `POST /automation/workflows/{id}/exec` with `{"stepID":"1","wait":true}`
returned `results.result = {"@value":2,"@type":"Integer"}`.

That is the minimum viable workflow. One step, no paths, termination implicit.

## 2. Service

`automation/service/workflow.go`, `autoService.DefaultWorkflow`. `Create` and
`Update` both take `*types.Workflow`.

**Update is a read-modify-write.** Per-field nil guards at `workflow.go:319-367`
mean omitting `steps` preserves them and sending `[]` wipes them — so a sparse
update is safe, but an explicitly empty array is destructive. `isStale` gives
optimistic locking on `updatedAt`, so load first and echo it back.

## 3. Authorization (§8.6) — passes

`CanCreateWorkflow` on create; `CanUpdateWorkflow` on update. Verified by
reading the methods, not by grepping for `ac.Can` — that shortcut has produced a
wrong answer twice in this project.

## 4. What is fatal, and what merely stores

Fatal: access control, `meta.name` empty, invalid or duplicate handle, stale
`updatedAt`, run-as user load failure.

**Everything else is an "issue", and issues do not block the write.**
`workflow.go:196-204` stores the workflow regardless and skips trigger
registration. So a graph that will never run returns 200.

The tool must therefore **return `issues` verbatim in its result**. A create
that reports success while the response carried three issues is the failure mode
this whole brief exists to prevent.

## 5. Step kinds

14, a typed enum at `types/workflow.gen.go:486-501`, all 14 handled by the
runtime — no type/runtime gap, unlike TAQ. `verifyStep`
(`workflow_converter.go:624-886`) is a per-kind arity, ref and argument table;
generate the schema's guidance from it rather than restating it by hand.

Start with `expressions` and `function` in the description's examples: the probe
proved `expressions` works with nothing but a `target`, an `expr` and a `type`.

## 6. Traps

**`ref` must exist in a function registry assembled at boot** from ~30
`*_handler.gen.go` files, and arguments are verified by name and type. A schema
cannot carry that list. The description must tell the caller to discover
functions first, the way `automation_taq_exec` already tells callers to look a
TAQ up before running it.

**Exactly one parentless step**, checked only at exec time
(`session.Start:204-219`). A multi-entry workflow writes clean and fails on
first run. The handler can check this before writing and refuse — cheap, and it
converts a confusing runtime error into a clear one.

**Path array order is semantics, not presentation.** Iterator `out[0]` is the
body and `out[1]` the exit; an error-handler's `out[1]` is the catch; an
exclusive gateway is tested in order with else last. No JSON Schema expresses
this, so it belongs in the description of whatever param carries paths.

**Triggers are a separate resource.** A workflow has no `Triggers` field —
`automation_trigger_create` attaches one, with `workflowStepID`. The create
description must say so, or a caller will build a workflow that nothing fires.

**`runAs` and `ownedBy` are `json:",string"`** — they must be quoted.

**Client-only rules the server does not enforce**: out-degree caps, connection
legality (no edge into a trigger, none out of a termination, no self-loops),
positional edge labels, parallel-gateway `ref` auto-flipped between fork and
join from edge balance. A tool that emits graphs the editor could not produce is
not necessarily broken, but it will look wrong when someone opens it.

**`meta.visual` is where node positions live.** Omitting it collapses every node
to (0,0) in the editor. Cosmetic, but a caller building a workflow someone will
later edit should know.

## 7. What the tool should do beyond the schema

1. Return `issues` verbatim.
2. Refuse a graph with more than one parentless step before writing it.
3. Read-modify-write on update, echoing `updatedAt`.
4. Tell the caller triggers are a separate call.

## 8. Assessment

**Shippable.** Small shape, guessable names, loud failures, and the probe hit it
first time. The realistic gap is `ref` discovery, which a description can route
around by pointing at the function registry.
