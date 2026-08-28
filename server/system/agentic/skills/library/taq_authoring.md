---
name: taq_authoring
description: Writing a TAQ — where the function and trigger vocabularies come from, what a stored-but-not-running TAQ looks like, and how expression scope and map literals go wrong.
triggers:
  - automation_taq_create
  - automation_taq_update
  - automation_taq_lookup
  - automation_taq_delete
  - automation_taq_undelete
  - automation_event_type_lookup
---

# Authoring a TAQ

TAQ is Trigger Action Query. It is a graph: **triggers** say what starts it,
**steps** say what it does, **paths** connect them. `automation_taq_create`
takes all three in one call.

Executing one is a different skill — see `automation`.

## Look the vocabulary up first, from the right place

Two endpoints, and they are the authority:

- `GET /automation/construct-library/functions` — every `ref` a function or
  iterator step may name, with its parameters, their exact type names, and
  which are required. There are 19.
- `GET /automation/construct-library/triggers` — the 22 legal
  `resourceType`/`eventType` pairs.

**Do not use `GET /automation/functions/`.** That is the workflow function
registry; most of its entries do not exist in a TAQ, and a step naming one
stores fine and never runs.

`automation:trigger:agentic` + `onAgentic` is the pair `automation_taq_exec`
runs, so give the TAQ one if you want to be able to test it. The
`compose:record` and `system:user` pairs fire on real events.

## The smallest thing that runs

One trigger, one function step, `"paths": []`. A lone unconnected trigger and a
lone unconnected step are wired together for you and termination is added.
More than one of either without paths gives a `graph.ambiguousEntry` issue.

## Stored is not running, and running is not working

- A TAQ whose graph fails validation is **still stored and still returns
  success**. What is wrong comes back under `issues`. Any issue at all, of any
  severity, leaves `runnable` false, and `automation_taq_exec` then answers
  "manager: executable not found" — which means the TAQ has issues, not that
  the ID is wrong. **No `issues` key at all is the clean result.**
- `enabled` defaults to true, which arms the triggers immediately. `false`
  authors without arming, but a disabled TAQ also refuses `automation_taq_exec`.
- After running, read `automation_taq_execution_trace`. `exec` reports status
  only, and a `completed` status with no working frames — nothing, or only
  `trigger` and `termination` frames — means no step ran.

A frame carries what it resolved under `args`. That is where you confirm an
expression evaluated to what you meant, and it is the fastest way to see a
filter that came out wrong.

## Arguments

Bind by `argumentName`, spelled exactly as the construct library lists it, and
supply every required parameter. `type` is compared literally against that
parameter's `types` array — omitting it is the same as sending `""` and fails.

Three ways to give a value:

- `"value"` — a literal.
- `"source"` + `"scope"` — an earlier step's output, `scope` being the
  producing step's handle.
- `"expr"` + `"scope"` — computed.

**An `expr` is evaluated inside its own `scope`, and without one that is the
trigger's.** So inside an iterator body the record variable is still the
trigger's record; reaching the loop item means giving that argument the
iterator's handle as `scope`. One expression sees one scope, so a formula
needing both the trigger's record and the loop item cannot be one expression —
split it across arguments, or compute it in a module field value expression
instead (see `field_expressions`).

The trigger scope also carries `invoker` and `runner`, so `invoker.email` is an
expression, not a parameter name.

## Map literals: quote the keys

A `values` argument for `composeRecordsCreate` is a KV, and the obvious literal
is wrong:

```
{"full_name": full_name, "stage": "applied"}   ✅ {"full_name":"Grace","stage":"applied"}
{full_name: full_name}                          ❌ {"Grace":"Grace"}
```

An unquoted key is evaluated as an expression like any other term, so the key
becomes the value. The run then fails with `no such field Grace` — the error
names the _value_ where you expect a field name, which is the tell.

A plain JSON object under `"value"` works too, when the contents are static.

## Record references in a filter

A Record-ref field holds the target's ID as a string, not a record, so a query
over them is string concatenation: `"card = " + record.values.card`. Writing
`record.values.card.recordID` yields nothing, the query matches nothing, and an
iterator over it runs **zero times in silence** — no error, no issue, and the
execution still reports completed. Only the absence of body frames after the
iterator's own frame shows it.
