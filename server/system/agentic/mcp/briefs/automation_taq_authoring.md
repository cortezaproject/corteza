# Brief: `automation_taq` create and update

Extends `server/automation/agentic/taq_tools.go` / `taq_handler.go`, which
already carry `lookup`, `exec`, `executions`, `execution_trace`, `delete` and
`undelete`.

Group `configuring`. Both risk `write`.

TAQ = Trigger Action Query. Never "Task Queue". The service resource is
`ngAutomation`; the product says TAQ.

**Read `automation_workflow_authoring.md` too, then assume none of it transfers.**
Three investigations independently concluded the two are not variants: different
endpoints, Go types, path-condition representation (`condition` AST object vs
`expr` gval string), gateway vocabulary (`gatewayExclusive` vs
`kind:"gateway"`+`ref:"excl"`), and trigger model (embedded vs separate
resource).

## 1. Five defects were fixed before this brief was written

The empirical probe took **six calls with two hard failures and one silent
failure** to get a TAQ running. Most of that was the server, not the graph.
Fixed on this branch, so the ground has moved:

| Was | Now |
|---|---|
| `spew.Dump(err); continue` dropped failed steps with no issue | reported as `step.invalid` with the converter's own message |
| omitting `meta` nil-panicked, HTTP 500 | `Error: missing name` |
| create overwrote supplied trigger IDs unconditionally | only mints missing ones |
| a disabled TAQ executed anyway | refuses, `NgAutomationErrDisabled` |
| `target`-named arguments validated clean and bound to nothing | bind by the same name validation resolved |

The first is the important one. **Six classes of broken TAQ used to store clean,
report `issues: null`, and exec as `"completed"` having run nothing**:
unsupported kind, unknown function ref, missing required argument, misnamed
argument, wrong argument type, unparseable expression. All six now surface.

The fifth is the same failure wearing the other three's clothes: `Expr.ArgKey`
resolves an argument's name as ArgumentName-falling-back-to-Target, and
`VerifyArguments` had always honoured that fallback, but `functionStep.ExecN`
grouped by `ArgumentName` alone. A `target`-named argument therefore passed
validation — whose error message even named the parameter it had matched — and
then bound to the empty key, so the step ran with `args: {}` and reported
`completed`. Fixed in `stepConvFunction` by resolving names onto copies before
the handler is built; workflow is untouched.

The third means a **one-call create is possible for the first time** — the
webapp still does create-empty-then-update because it predates the fix.

## 2. Service

`automation/service/ng_automation.go`, `autoService.DefaultNgAutomation`.
Name lives in `Meta.Short`, not `Meta.Name`.

Handle uniqueness is scoped **within project**, unlike workflow's global check.

`svc.engine(ctx)` can hard-fail with `automation runtime missing for
tenant/project` — so a *clean* automation can fail to create where a *broken*
one succeeds, because registration only runs when there are no issues.

## 3. Authorization (§8.6) — passes

`CanCreateNgAutomation` / `CanUpdateNgAutomation`. Read the methods; do not grep.

## 4. Step kinds — six, and they are now named

`types.NgAutomationStepKinds` and `IsValidNgAutomationStepKind` were added on
this branch: `function`, `iterator`, `gatewayExclusive`, `gatewayInclusive`,
`termination`, `error`.

**Use them.** `Kind` is a plain string on the generated type where Workflow has
a typed enum, which is precisely why `expressions` and other workflow kinds
leaked in silently. The tool should validate against
`IsValidNgAutomationStepKind` before writing and name the six in its description
— a model that has just written a workflow will reach for `expressions`.

## 5. Arguments and the construct library

A `function` or `iterator` step needs `ref` naming a construct-library function,
and `arguments` matching its parameters by `argumentName`.

Name an argument with `argumentName`. `target` now works as a synonym (see §1),
but it is the wrong word — `target` names where a *result* is written — and a
tool's schema should ask for `argumentName` and nothing else.

`type` is load-bearing: `ParamSet.VerifyArguments` compares `p.HasType(e.Type)`
literally, so omitting it fails. That failure is now visible; before this branch
it was one of the six silent classes.

The description must route the caller to the construct library the same way
`automation_taq_exec` routes them to `lookup`. A 17-entry list is not something
a schema can carry.

## 6. Minimum viable TAQ

One trigger, one `function` step, **and deliberately zero paths** — the
single-orphan inference (`ng_automation_converter.go:394-397`) wires the entry.
Termination is auto-injected.

Probe's working payload, all three arguments required:

```json
{"handle":"agent-x","meta":{"short":"X"},"enabled":true,
 "triggers":[{"handle":"entry","meta":{"short":"e"},"enabled":true,
   "resourceType":"automation:trigger:agentic","eventType":"onAgentic"}],
 "steps":[{"stepID":"1","handle":"notify","kind":"function","ref":"notificationSend",
   "meta":{"short":"Notify"},
   "arguments":[{"argumentName":"recipient","value":"<userID>","type":"ID"},
                {"argumentName":"title","value":"t","type":"String"},
                {"argumentName":"description","value":"d","type":"String"}]}],
 "paths":[]}
```

`recipient` needs a numeric ID with `"type":"ID"` — a handle string fails at run
time with `user not found`.

Now that trigger IDs survive create, a trigger→step path *can* be authored in
one call. Verify that before relying on it: the fix is new and the single-orphan
inference may make the explicit path unnecessary anyway.

## 7. Traps that remain

**`enabled` is not what it sounds like** — it gated only trigger firing until
this branch. It now also blocks exec. Worth stating plainly in the description.

**Entry points are all parentless steps** (`step_scheduler.go:581-598`).
`EntryPoint` names the global scope alias and the trace frame; it does not
select a start structurally.

**`exec` returns no results.** Proof of success requires the trace endpoint. A
`completed` status with a `null` trace and sub-0.1ms duration is a TAQ that ran
nothing — the tool's description must say so, since it is the signature of the
class of bug just fixed.

**Path `Condition` is an AST**, not a string, and **nothing validates it at
write time** — not operator, symbol or scope. Operators: `and, or, not, isNull,
isNotNull, eq, ne, lt, gt, lte, gte`.

**`buildExecSteps` overwrites `step.Results`** from the function definition, on
the object that then gets persisted. A caller's submitted `Results` are silently
replaced. Do not expose `results`.

**An `error` step's message can never be set through the API** —
`stepConvError` type-asserts `*ast.ASTNode` on a `map[string]any` that came from
JSON, so it always fails. Known, unfixed, narrow. Do not expose the message
field; it would do nothing.

**Scope refs** pass write-time validation for any declared handle but resolve at
run time only for `""`, `"global"`, the active entry point, and completed steps
on reachable branches.

## 8. What the tool should do beyond the schema

1. **Return `issues` verbatim.** Now that they are populated, this is the whole
   feedback loop.
2. **Diff submitted steps against returned steps.** A step dropped from the exec
   graph used to be invisible; it is now an issue, but the diff is a cheap
   second net.
3. Validate `kind` against `IsValidNgAutomationStepKind` before writing.
4. Treat a `completed` execution with an empty trace as a failure, not a success.

## 9. Assessment — provisional

Before the four fixes: **not shippable**, because the failure mode was silence.
After them the loud path works, which is the precondition rather than the proof.

A cold re-probe is running to measure how many attempts a model needs now. Do
not start writing tools until its result is in — if the answer is still "many,
and it believed it had succeeded", the honest outcome is to ship workflow
authoring alone and record TAQ as a deliberate gap with this brief as the
reason.
