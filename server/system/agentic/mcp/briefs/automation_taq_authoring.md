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

**The construct library is `GET /automation/construct-library/functions`** — 17
entries, served by `service.ConstructLibrary()`, the same singleton
`stepConvFunction` resolves `step.Ref` against. Its sibling
`GET /automation/construct-library/triggers` returns the 22
`resourceType`/`eventType` pairs, which a caller otherwise has no way to learn.
Neither takes a trailing slash; the bare directory path 404s.

**`GET /automation/functions/` is a different registry and a trap.** It is
`service.Registry()` — the *workflow* function registry, 93 entries, 78 of which
a TAQ step cannot use. The two overlap on `notificationSend`, which is why an
empirical probe that only ever sent `notificationSend` could not tell them apart
and reported the wrong endpoint. `ref:"logInfo"` is the cheap discriminator: it
is in the workflow registry and comes back `function.unknown` from a TAQ. The
tool's description must name the construct-library endpoint and warn off the
other one. (Their parameters are keyed differently too — `argumentName` in the
construct library, `name` in the workflow registry — and `paramArgKey` falls back
to `Name`, which hides the difference for the refs they share.)

**`types` is plural on the parameter, singular on the argument.** The parameter
declares `types: ["ID","Handle","String"]`; the argument carries one `type` that
must be a member, spelled exactly. `"Number"` normalises to `"Integer"` and is
then rejected against that set.

## 6. Minimum viable TAQ

One trigger, one `function` step, **and deliberately zero paths** — the
single-orphan inference (`ng_automation_converter.go:394-397`) wires the entry.
Termination is auto-injected.

Probe's working payload — `recipient` and `title` are required, `description` is
not:

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

A trigger→step path can be authored in one call — verified: a supplied
`triggerID` survives create and both `trigger→step` and `step→step` edges wire
up, running in order.

## 7. Traps that remain

**`enabled` is not what it sounds like** — it gated only trigger firing until
this branch. It now also blocks exec. Worth stating plainly in the description.

**Entry points are all parentless steps** (`step_scheduler.go:581-598`).
`EntryPoint` names the global scope alias and the trace frame; it does not
select a start structurally.

**`exec` returns no results.** Proof of success requires the trace endpoint. A
`completed` status with a `null` trace and sub-0.1ms duration is a TAQ that ran
nothing — the tool's description must say so, since it is the signature of the
class of bug just fixed. `steps: []` stores clean and produces a trace holding
the trigger frame and no step frames — "no step frames", not "a null trace", is
the signature to check.

**A populated `args` frame proves binding, not effect.** `args: {}` or `null`
means binding failed; a full `args` means the values reached the handler and
nothing more. The cold probe caught a real case only by going outside the
described protocol. Where a construct has an observable effect, the description
must tell the caller to verify it independently — for `notificationSend`, `GET
/system/notification/?limit=5&sort=id%20DESC`. Do not let a tool report delivery
on the strength of the trace alone.

**Runtime failures are not `issues`.** `issues` covers authoring defects; a bad
expression or a missing referenced entity surfaces at exec as `status:"failed"`
with a typed frame error. And `issues` is *absent* on a clean write, not empty —
treat a missing key as success and any `severity:"error"` entry as fatal. Such a
TAQ stores at 200 but never registers, and `exec` answers `manager: executable
not found`, which means "your TAQ has issues", not "wrong ID".

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

## 9. Assessment — shippable

A cold re-probe, run after the five fixes with only the knowledge a tool
description would carry, built a working TAQ in **one payload submission and one
execution** — down from six and three. It then authored fifteen deliberately
broken ones: six were rejected at write time with a code, a severity, the exact
parameter and the exact legal types; two more failed loudly at exec with a typed
error; `enabled:false` refused. A broken TAQ cannot reach the runtime at all.

The residual silent class is no longer structural but referential — a
well-formed argument naming an entity that does not exist. The probe found one
(`notificationSend` trusted a recipient ID without a lookup, so a valid-looking
ID wrote a notification addressed to nobody, undeletable because the
notification endpoints scope to the *recipient*, so only a notification
addressed to a user that does not exist is unreachable — one addressed to you is
deletable normally). That is fixed: all three recipient
forms resolve, and a bad ID now fails with the same `user not found` a bad
handle always did. Other constructs have not been audited for the same
asymmetry, which is the honest remaining caveat.
