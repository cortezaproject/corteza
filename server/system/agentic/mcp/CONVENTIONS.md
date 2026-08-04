# MCP tool conventions

**Status: agreed, partially implemented.** See §14 for what is built and what
is not.

This document is the input to a fan-out that will write roughly 165 more tools
across the 42 resources in RESOURCES.md. Everything here is inherited ~170 times. Read it as a
spec, not as background.

> **Why this is not `mcp.intent.md`.** The intent system covers `server` in
> `.intent/config.mjs` `covered` but **not** in `enforced`, and
> `.intent/TODO.md` reserves Phase 4 server-backfill scoping for the human
> (package granularity, project-file exclusion, which top-level dirs are in
> scope). Writing `mcp.intent.md` would be the first server intent doc in the
> repo and would set precedent for ~197 packages. Convert this file once Phase
> 4 is scoped.

---

## 1. What this system is

Human exposes its capabilities to LLM agents through a single MCP tool
registry (`mcp_registry.go`). That registry has two consumers:

| Surface | Path | Consumer |
|---|---|---|
| In-process | `Registry.GetTools` / `ExecuteTool` | Human's own agentic runtime (agents, chatbots, magic buttons) |
| HTTP | `/api/mcp` (streamable) | External MCP clients, primarily Claude Code |

Both surfaces call the same handler functions. Where they legitimately differ
is documented in §2.3; where they differed by accident, that has been fixed.

36 tools are registered today across 9 resources. The target is roughly 200,
so the fan-out is roughly 165 tools.

### Layers this serves

Four usage layers, **not** four servers:

1. **Development** — changing Human's own code. Repo-side.
2. **Configuring** — shaping the system: modules, pages, TAQs, workflows, roles.
3. **Usage** — operating the configured system: records, running automations,
   chatbots, reports.
4. **Connections** — third-party integrations. Human Brain is the catalog;
   `connection` / `configuredConnection` already exist as Human resources.
   Deferred.

---

## 2. Locked decisions

| Decision | Ruling |
|---|---|
| Topology | One server, one registry. No per-layer servers. |
| Groups | `development` / `configuring` / `usage`. Connections deferred. |
| Group semantics | **Filtering only.** Decides what gets listed. Not a security boundary. |
| Risk semantics | `read` / `write` / `destructive`. Sole writer of the protocol's annotation hints, and an enforced per-session ceiling. |
| Scope | Per request, from the URL: `/api/mcp/{group}` narrows the listing, `?maxRisk=` caps and is refused on dispatch. See §2.5. |
| Security boundary | authclient + RBAC — with documented exceptions, see §2.2. |
| Coverage target | Hand-written tools, full resource coverage, ~200 tools. |
| Tenancy | Single instance. Multi-tenant MCP is out of scope. |

### 2.1 Why groups are not enforcement

Most handlers pass the request context straight to the service layer, so RBAC
governs what a caller can do. A caller in Claude Code is authenticated as
themselves and could perform the same operations via REST or the webapp. A
group check would block nothing real.

Groups exist to keep the advertised tool list small, and the cost is measured
rather than assumed: 35 tools is ~11,400 tokens per request, which projects to
**~65,000 tokens per request at 200 tools** — paid on every call before any
work happens. `dev/agent/mcp-verify.py cost` reprints the number.

**Consequence for tool authors:** never rely on a group tag for safety. If an
operation needs authorisation it comes from RBAC in the service layer — and
§8.6 says what to do when the service has none.

### 2.2 The RBAC premise is not an invariant

RBAC is the boundary *for most services*. Known exceptions:

- `compose/agentic/module_handler.go` and `namespace_handler.go` both call
  `a.SetIdentityToContext(ctx, a.ServiceUser())` — the super-admin bypass role
  (`pkg/auth/system.go`). The comment says the elevation is deliberate.
- `system/service/chatbot_session.go` `Search` goes straight to the store with
  no access check; authorization lives in the REST controller.

§8.6 exists because "most" is not "all".

### 2.3 Surface differences

**By design:** `ToolAliases` is applied on the in-process path only. Agent
definitions store tool names by value, so a rename would strand them; remote
clients list tools live and never hold a stale name. Registering aliases with
the mcp-go server would advertise every historical name to every client.

**Fixed:** `Hidden` tools were being advertised over HTTP because
`mcp_server.go` copied the registry without checking the flag. `ExecuteTool`
did not resolve aliases, so an aliased tool listed fine and failed on dispatch.
Duplicate tool names silently overwrote each other.

**Fixed, but verified differently:** `discovery_search` used to declare
`namespace` and `module` params its handler never read — it read
`namespaceIDs`/`moduleIDs`, injected only by the in-process executor, so over
HTTP the declared scope was ignored and the search ran across everything the
caller's discovery token allowed. The handler now resolves the declared args
itself when nothing was injected; injected IDs still win, because in-process
they are the authorization narrowing already checked against the agent's
allow-list.

This is the one tool that is **not live-verified**. Its backing discovery
service does not run in the local dev environment (`/api/discovery` is 404), so
`isAvailable` filters it out of `tools/list` — correct behaviour, and also why
the fix could only be verified by build and by reading. Everything else in this
document was exercised against a running server.

`discovery_search` also declares no `pageCursor`: the discovery service offsets
with `from` and returns no cursor, so advertising one would promise paging that
cannot work (§8.1). It declares `limit` only, and says so.

### 2.6 Progressive disclosure

A session sees five tools until it asks for more: `human_tool_search`,
`human_tool_load`, and the compose read path (`namespace`, `module`, `record`
lookups). Everything else is loaded by searching for it.

Measured on the running server: **~1,442 tokens instead of ~27,368, a 19x
reduction**, and it stops growing — 200 tools cost a session the same as 20.

Why not the alternatives. Trimming descriptions fights §8.5, which made them
rich precisely because a caller without the repo has nothing else; disclosure
lets a description be as long as it needs to be because few are loaded. Group
endpoints help but not enough — `configuring` is 64 of 80, because most tools
genuinely are configuration. `tools/list` pagination does nothing, because
clients drain every page.

`human_tool_search` returns each match's **full definition inline** as well as
registering it for `notifications/tools/list_changed`. That is deliberate
redundancy: a client that honours the notification re-lists and sees the tools
properly, and one that ignores it can still call straight from the search
result. Without it, disclosure would break silently on any client that does not
re-list, and the failure would look like the tool not existing.

Neither disclosure nor group is enforced at dispatch. A caller naming a tool it
was never shown has done nothing RBAC would not already allow, and refusing
would break exactly the clients the inline schemas exist for. The risk ceiling
*is* enforced — see §2.5.

`?tools=all` opts out and lists everything the group and risk allow. It is for
tooling: `dev/agent/mcp-verify.py` has to audit the whole surface, and a human
debugging "why can the model not see X" needs the unfiltered list. It is
deliberately absent from every tool description, because an agent using it
would pay the cost this exists to avoid.

The two meta-tools live on the mcp-go server, not in the `Registry`: they are a
property of this transport, not of Human. The in-process runtime scopes an
agent with `allowedTools` and has no listing to shrink. This also keeps them
out of the coverage matrix, where they would read as resources they are not —
and out of the structural test, which asserts things about resource tools.

### 2.4 The in-process policy layer

`policy/policy.go` is the authorization gate for the in-process surface; the
HTTP path never calls it.

Per-TAQ tools are minted at runtime as `automation_<numeric id>`. Policy used
to treat *any* `automation_*` name outside a hardcoded three-item list as a TAQ
ID, which denied `automation_taq_exec`, `automation_taq_executions` and
`automation_taq_execution_trace` outright. Now matched on a numeric suffix,
which is how the names are actually minted.

A tool with no `buildResource` case gets no resource-level narrowing. That is
**not** denied at runtime — the agent's `Access.Tools` allow-list is already
deny-by-default, and denying here would break every newly added tool until
someone edited `policy.go`, a cross-package coupling a tool author has no
reason to discover. `policy.IsClassified` backs a CI-time assertion instead.
The problem was that the default was *silent*, not that it was permissive.

`policy.go` holds a second alias map that must stay in sync with
`mcp_registry.go`.

### 2.5 Scope: group filtering and the risk ceiling

A request carries a `Scope` resolved from its URL (`scope.go`):

| URL | Effect |
|---|---|
| `/api/mcp` | Everything. No narrowing, no ceiling. |
| `/api/mcp/configuring` | Lists only configuring tools. |
| `/api/mcp/usage?maxRisk=read` | Lists only usage reads, and refuses anything above read on dispatch. |

The two dimensions are enforced differently, on purpose. **Group is filtered
only** — it decides what `tools/list` returns, because its job is to keep the
list small, and a caller naming a tool outside its group has done nothing RBAC
would not already allow. **Risk is filtered *and* refused at dispatch**, because
a ceiling that only hid tools would mean nothing to a client that already knew a
name.

An unrecognised group or risk is ignored rather than rejected: a typo must not
silently narrow the surface to nothing and leave the caller thinking the server
is broken.

The ceiling is self-selected, so it is a seatbelt against accidents — pointing a
session at production and having it delete a namespace — not a lock against a
hostile caller, who simply would not set it. Moving it somewhere a caller cannot
choose means putting it in the token; see §14.

---

## 3. Groups

| Group | Contains | Test |
|---|---|---|
| `configuring` | Schema and definition level | Changes what the system *is* |
| `usage` | Data and execution level | Changes what the system *holds*, or runs it |
| `development` | Repo-level tooling | Operates on source, not on an instance |

The line already exists in code: `compose_module_*` is configuring,
`compose_record_*` is usage. Definitions are configuring, runs are usage —
hence `automation_taq_lookup` is configuring while `automation_taq_exec` is
usage.

A tool may carry more than one group. Prefer one: a tool that feels like both
is usually two tools.

`development` is currently empty, defined so L1 tooling has a home without a
later retrofit.

---

## 4. Risk levels

| Level | Meaning | Ops |
|---|---|---|
| `read` | No state change | `lookup`, `executions`, `execution_trace` |
| `write` | Creates or modifies | `create`, `update`, `undelete`, `exec` |
| `destructive` | Removes from view or loses data | `delete`, revoke, purge |

`exec` is `write`, not `read`, even when the executed thing happens to be
read-only — the registry cannot know what a TAQ or workflow does.

**`read` is about state change, not disclosure.** A read-shaped op can expose
credentials (`system/service/auth_client.go` `ExposeSecret`). Disclosure
sensitivity is handled by §8.6, not by the risk ladder.

### 4.1 Risk is the sole writer of annotation hints

`mcp.NewTool` defaults `DestructiveHint=true` and `OpenWorldHint=true` and
always serializes annotations, so a tool that does not override them advertises
itself to every client as destructive and open-world. All 29 tools did exactly
that, read-only lookups included.

`hmcp.WithRisk` therefore sets all four hints and **authors never write them by
hand**:

| Risk | readOnly | destructive | idempotent | openWorld |
|---|---|---|---|---|
| `read` | true | false | true | false |
| `write` | false | false | false | false |
| `destructive` | false | true | true | false |

`openWorld` is always false: every Human tool acts on this instance's own data.

### 4.2 Deletes are soft

`compose/service/namespace.go` `handleDelete` sets only `DeletedAt`; there is
no cascade, and `handleUndelete` restores intact. 29 services expose
`UndeleteByID` against 45 with `DeleteByID`.

Delete stays `destructive` — a soft-deleted namespace is invisible to every
consumer, so the blast radius is real — but `undelete` is in the op vocabulary
so the model has a way back.

**Rule: every resource whose service exposes `UndeleteByID` gets an `undelete`
tool.** Not a judgement call — if the service can reverse a delete, the model
gets the same way back a human has through the UI. It is roughly fifteen lines,
risk `write`, and its description must say how to find a deleted item (the
resource's `includeDeleted` filter, where one exists).

Where the service has no `UndeleteByID`, the `delete` description must say so
plainly, because a model has no other way to learn that the call is one-way.
`system_reminder_delete` is the worked example.

Note `UndeleteByID` is generated into `*.gen.go`, not the hand-written service
file — grepping only the latter will wrongly conclude a resource has no
undelete.

`compose_module_delete`'s description still claims it "permanently removes the
module and all its records", which is false and should be corrected.

---

## 5. File structure

```
server/{app}/agentic/
  {resource}_tools.go          declarations only
  {resource}_handler.go        implementations, identical order
  {resource}_handler_test.go   tests
```

Both directions must be derivable without searching:

- `compose_page_create` → app `compose`, resource `page`, op `create`
  → `compose/agentic/page_tools.go` (declaration)
  → `compose/agentic/page_handler.go` (implementation)

### Why declarations are split out

`page_handler.go` is 569 lines for 6 tools. An agent asking "what page tools
exist" reads 569 lines to learn about 60 lines' worth. Across 58 resources
that is ~17k lines of handler code wrapping perhaps 3k lines of actual spec.

After the split, reading every `*_tools.go` gives a complete picture of the
tool surface cheaply. This matters because Human is open source — an agent
working on the repo has the source available, and the source should be the best
available specification.

`{resource}_tools.go`: the `register()` method, tool declarations,
descriptions, schemas, group and risk tags. Nothing else.

`{resource}_handler.go`: handler struct, constructor, and handler methods in
declaration order.

### Wiring

Every handler is hand-wired in `server/app/boot_levels.go`. That is a single
serial edit point and a merge-conflict hotspot for parallel fan-out: **agents
do not edit it**, the orchestrator wires each batch centrally after it lands.

### Intra-file order

Fixed: `lookup`, `create`, `update`, `delete`, `undelete`, then domain ops.
Handlers mirror declaration order exactly.

---

## 6. Naming

```
{app}_{resource}_{op}
```

- `app` — `compose`, `automation`, `system`
- `resource` — from the published resource-name table, not invented
- `op` — `lookup` | `create` | `update` | `delete` | `undelete` | `exec`, or a
  named domain op

`lookup` covers both fetch-one and list. One tool, not two — see §8.1 for the
contract that makes this unambiguous.

**Renames** go in `ToolAliases` *and* in `policy.go`'s alias map.

### Legacy names

Three shipped names do not match the grammar and are grandfathered by an
exception list rather than renamed: `discovery_search` (two segments, app
`discovery` in package `system/agentic`, hardcoded in eight places in
`executor.go`), `compose_page_block_schema` and
`automation_taq_execution_trace` (four segments).

---

## 7. Shared toolkit

`server/system/agentic/toolkit` removes the boilerplate every handler repeats
and, more importantly, creates single choke points.

| Helper | Purpose |
|---|---|
| `Args(req)` | Unwrap arguments, uniform error |
| `Str` / `ReqStr` | Optional / required string |
| `ID` / `ReqID` | Parse a string ID to `uint64`; **rejects** a numeric argument rather than coercing it |
| `Page(args)` | `limit` (default 50, capped 200) and `pageCursor` |
| `JSONResult(v)` | Marshal, enforce the 256 KiB ceiling, wrap |
| `TextResult(...)` | Non-JSON results — delete handlers acknowledge in plain text |
| `Errf(subject, err)` | Uniform error wrapping |

`JSONResult` is the only sanctioned path from a Go value to a tool result. The
size ceiling lives there, and if tool output ever needs to mark third-party
content as data rather than instruction, that is one function to change instead
of ~200 call sites.

---

## 8. Tool authoring rules

### 8.1 The lookup contract

**Reference param is never `Required`.** `compose_namespace_lookup` declares
`namespace` as required while its handler branches on an empty ref to list
everything — the list mode is unreachable to a schema-obeying caller.

**List mode returns a slim projection**, single-item mode returns the raw
service type. This is what the code already does: module, chart and namespace
lookups all build trimmed items for lists. Projection shape:
`{<res>ID, name/title, handle}` plus whatever disambiguates.

The exception is a resource whose heavy field *is* the answer. A module's
fields, a chart's config and a reminder's payload are incidental to picking one
out of a list, so they are projected away; a record's values are the thing the
caller asked for, so `compose_record_lookup` returns records in full. Where you
take this exception, paging and the `toolkit.JSONResult` ceiling do the work
that projection does elsewhere, and the description must tell the caller to
narrow rather than list.

**Every lookup declares `limit` and `pageCursor`.** `compose_record_lookup`
builds a `RecordFilter` with zero-value paging and `dalutils` loops
`for f.Limit == 0 || ...` — one call drains an entire module into the model's
context. REST caps this at 500/1000; MCP does not.

`toolkit.JSONResult` enforces a hard ceiling regardless. Schema paging is what
the model drives; the ceiling is the backstop.

**Return the cursor value, never its `String()`.** Put the
`*filter.PagingCursor` straight into the result and let `json.Marshal` call its
`MarshalJSON`, which emits the base64 form `filter.parseCursor` accepts. A nil
pointer marshals to `null`, so the last page needs no special case.
`PagingCursor.String()` is a human-readable debug rendering
(`<id: 123, [FWD]>`) and cannot be fed back — a lookup that returns it
advertises a cursor that always fails to decode, making paging one-way. This
rule exists because the first two tools written to this spec both got it wrong.

### 8.2 Partial updates

**Absent = unchanged. Present-and-empty = clear.** Collections always replace.
Merge only where a tool documents it and ships a `removeX` companion.

Today `page_handler.go` gives an empty string three different meanings inside
one function, and `chart_handler.go` says config replaces wholesale while
`module_handler.go` says fields are merged.

### 8.3 IDs are strings

Human IDs are `uint64` and exceed JavaScript's safe integer range. Every ID in
a tool schema is `WithString` and the description says so. Parse with
`toolkit.ID`, which rejects a numeric argument — by the time a JSON number
reaches the handler it has already lost precision, so coercing would launder a
corrupted value into the service layer.

### 8.4 Accept names, handles, and IDs — where you can

Where a service offers `FindByAny`, accept any of them. **Only 6 services have
it** (`role`, `namespace`, `template`, `user_group`, `user`, `module`) against
~58 target resources. For the rest the per-resource brief states the identifier
strategy; do not invent one.

Note `ns_mod_resolver.go` is **not** used by any MCP handler — it is wired into
the agentic runtime only. The handlers all resolve inline.

### 8.5 Descriptions

Descriptions are the entire interface for a caller without the repository. A
configurator using Claude Code against a hosted instance has no source to read.

Every description states: what the tool does; when to use it and, where there
is a common mistake, when *not* to; how it relates to adjacent tools; the shape
of anything non-obvious.

Worth copying: `compose_record_lookup` ("Do NOT call this before creating a
record" — prevents a specific observed failure); `compose_chart_create` (tells
the caller the next step including the exact block JSON);
`compose_chart_delete` (warns that pages referencing the chart render empty).

Review descriptions against the no-repo case. A description that only makes
sense to someone who can read the handler has failed.

Annotations are never hand-written — `WithRisk` owns them.

### 8.6 Confirm the service authorizes, before writing the tool

A handler passes the caller's context to a service and returns what it gets.
That is only safe if the service authorizes. Most do; some do not; and the
check is not always RBAC.

**Mechanical rule, no judgement required.** Read the service method. It must do
at least one of:

1. call `ac.Can*`, or
2. constrain the result by the caller's identity — an ownership or assignee
   scope applied inside the service.

If it does neither, do not write the tool. File it as a gap in the coverage
matrix.

Both forms are real. `system/service/reminder.go` has exactly one `ac.Can*`
call (`CanAssignReminder`, and only on assignment) yet is correctly scoped:
`onLookup` returns `ReminderErrNotAllowedToRead` unless the reminder is
assigned to the caller, and `onSearch` sets `filter.Check` to the same
predicate. An `ac.Can*`-only rule would have rejected it, and would have taught
every fan-out agent to reject every ownership-scoped service.

The failing shape is neither: `system/service/chatbot_session.go` `Search` is a
bare `store.SearchChatbotSessions` with no check of any kind — the
authorization lives in the REST controller, which a tool does not go through.

Seed deny-list: `chatbotSession`, `chatbotPreview`, `role.Membership`.

### 8.6b Disclosure is a separate reason not to write a tool

A method can pass §8.6 and still not deserve a tool.
`authClient.ExposeSecret` is the case: it *is* authorized — `lookupByID` calls
`CanReadAuthClient` — but it returns a working credential, and read permission
is a lower bar than a credential deserves. `LookupByID` deliberately blanks
`Secret` before returning; `ExposeSecret` does not.

The disclosure is also irreversible in a way a write is not: once a secret is
in a model's context it is in the transcript and the logs. Same reasoning
excludes `user.SetPassword`, which is authorized correctly and which no agent
workflow should be performing without a human in the loop.

**Rule:** a tool that returns a credential, or sets one, is not written. File
it in the coverage matrix with the reason, so the gap reads as deliberate
rather than missed.

### 8.7 Errors

`fmt.Errorf("<subject> <verb> failed: %w", err)` via `toolkit.Errf` — lowercase,
wrapped, no punctuation. Errors are read by a model deciding what to do next.

---

## 9. Registration and tagging

Registration-time concerns go through `RegisterOption`; tool-definition
concerns ride on the `mcp.Tool`:

```go
h.reg.RegisterTool(
    mcp.NewTool("compose_record_lookup",
        mcp.WithDescription("..."),
        mcp.WithString("namespace", mcp.Description("...")),
        hmcp.InGroup(hmcp.GroupUsage),
        hmcp.WithRisk(hmcp.RiskRead),
    ),
    "Lookup record",
    h.lookup,
    hmcp.Hidden(), // optional
)
```

Human's package is imported as `hmcp`, because `mcp` is already bound to
mark3labs in every handler file and is needed there for
`NewTool`/`WithString`/`CallToolRequest`.

Group and risk travel in `Tool.Meta` rather than the registrar so that
registrar interfaces and call sites keep their arity — adding tags to ~200
tools is a per-declaration edit, never an interface migration — and so
`server.ToolFilterFunc`, which receives `[]mcp.Tool`, can read the tag directly
when per-request group/risk filtering lands.

`RegisterTool` **panics on a duplicate name**. Registration is boot-time, and
two fan-out agents emitting the same name previously produced a green build
with one tool silently gone.

`Hidden()` excludes a tool from the mcp-go server entirely, so it is neither
listed nor callable over HTTP — a filter would not do, since mcp-go tool
filters apply to `list_tools` and not to dispatch. `Available(fn)` is a runtime
predicate (`discovery_search` health-checks the discovery service), so it runs
through a tool filter that mcp-go re-runs per request.

---

## 10. Coverage matrix and structural test

Generated from the registry into `TOOLS.md`: tool name, file, group, risk, plus
every known resource with no tool yet and every resource on the §8.6 deny-list
with its reason.

The structural test asserts, per tool: name matches the grammar or is on the
legacy exception list; the app segment matches the package; declaration lives
in `*_tools.go`; group and risk are set; description is non-empty and above a
length floor; ID-typed params are declared as strings; lookup tools declare
`limit` and `pageCursor` and do not mark the ref param required; no duplicate
names; `policy.IsClassified` is true for every registered tool; the two alias
maps agree.

It cannot live in the `mcp` package — that would invert the dependency, since
every handler package imports `mcp`. It belongs in a package that already
imports all of them.

---

## 11. Verification is required after every batch

`dev/agent/mcp-verify.py` asks the running server what it actually does. Run it
after every batch of tools, not at the end.

```sh
dev/agent/mcp-verify.py            # auth, contracts, scope, exercise, cost
dev/agent/mcp-verify.py scope      # one section
```

This is not belt-and-braces. The `/api/mcp` authentication fix passed its unit
tests and did nothing, because `HttpTokenValidator` passes a request with no
token through — visible only by asking the server. The paging cursor was
likewise declared correctly and unusable in practice. Unit tests assert what the
code declares; this asserts what the server answers.

Exercise mode writes to the dev server under an `agent-` prefix and removes what
it creates, per CLAUDE.md.

## 12. Per-resource briefs and the resource-name table

The draft alone is not a sufficient brief. Writing tools for `system/reminder`
from prose required inventing roughly ten things: which of `ReminderFilter`'s 7
fields become params, what to do absent `FindByAny`, how `*time.Time` and
`types.JSONText` params are shaped, whether `assignedTo` defaults to the
caller, which of `Dismiss`/`Undismiss`/`Snooze` become domain ops.

Three deliverables before fan-out:

**Resource-name table.** REST resource → tool resource segment, published and
non-negotiable. "Matching the REST resource name" has no consistent answer for
`system_project*` (9 REST resources under one prefix, ~36 tools), for
`connection` / `configuredConnection` / `dalConnection`, or for multi-word
names — `system_userGroup_lookup` mixes camelCase into a snake_case grammar.
Without the table, ~170 agents invent ~170 answers.

**Exemplar.** `system/reminder` implemented end-to-end by hand: soft delete,
domain ops, no handle, caller scoping, free-form JSON payload — more edge cases
than any resource in the first fan-out batch.

**Per-resource brief**, generated in one pass: service methods to expose,
filter fields to expose, identifier strategy, scoping defaults, RBAC-check
presence per §8.6, domain ops.

---

## 13. Known hazards

**Repo / instance version skew.** An agent may hold the repository at one
commit while the endpoint serves another build. Ruling: expose version, build
commit and registry hash — as a **tool**, not a resource.
`Registry.RegisterResource` has zero callers repo-wide so the resources map is
always empty, and the in-process consumer interface is `GetTools`/`ExecuteTool`
only, so a resource is invisible to half the consumers in §1. Note
`Opt.Agentic.McpServerVersion` defaults to the literal `"v1"`, not a build
version.

**Bypass.** An agent with repository access can call REST directly, modify the
database, or write a mutating test. Documented, not fought: the development
surface has no ceiling by construction. Configurator sessions should not have
the repository.

**Description atrophy.** Weak descriptions are invisible to anyone with the
source. Reviewed against the no-repo case explicitly; the structural test
enforces presence and length only.

**Compound injection.** `usage` tools pull arbitrary record content — and, once
connections land, third-party content — into a session that may also hold
`configuring` tools. Recorded now, designed at L4. `toolkit.JSONResult` exists
so a mitigation is a single change rather than a ~200-file retrofit. This is
the hazard most likely to force revisiting "groups are filtering only".

---

## 14. Implementation status

**Built:**

- `/api/mcp` requires **both** `auth.HttpTokenValidator("api")` and
  `auth.HttpAuthenticatedOnly()`, and they do different jobs. The validator
  rejects a token minted for another scope; on its own it does *not* require a
  token at all, because it deliberately passes `jwtauth.ErrNoTokenFound`
  through. That is fine for REST, whose handlers reach RBAC and RBAC denies the
  anonymous role — but MCP tool discovery never reaches RBAC, it lists every
  registered tool and description straight from the registry. Verified on a
  running server: validator alone still answered an unauthenticated
  `tools/list` with 200.

  Known and pre-existing: a **malformed** token yields 500 rather than 401.
  This is global to the API, not specific to `/mcp` (`/system/settings/current`
  behaves the same), and lives in the shared token verifier.
- `Hidden` enforced by exclusion; `Available` enforced per-request via
  `server.WithToolFilter`.
- `ExecuteTool` resolves aliases. Duplicate registration panics. Schema
  building factored out of three copies.
- `RegisterHiddenTool` / `RegisterToolWithAvailability` collapsed into
  `RegisterTool` plus `Hidden()` / `Available(fn)`.
- `policy.go` numeric-suffix TAQ matching; `IsClassified` for CI-time coverage.
- `toolkit` package with tests.
- `InGroup` / `WithRisk` options; every tool tagged; hand-written
  `WithReadOnlyHintAnnotation` calls removed in favour of `WithRisk`.
- Declaration/implementation file split across all nine resources (§5).
- `RESOURCES.md` name table; `system/reminder` exemplar (§12).
- Structural test in `server/tests/mcp` and generated `TOOLS.md` (§10).
- Scope: `/api/mcp/{group}` filtering and an enforced `?maxRisk=` ceiling
  (§2.5). Measured: 35 tools cost ~11,400 tokens per request.
- `dev/agent/mcp-verify.py`, the live check required after every batch (§11).

**Not built:**

- Per-resource briefs for the identity batch (§12).
- Version handshake tool (§13).
- Identity read-only slice — the first real fan-out batch.
- A token-carried risk ceiling. The URL ceiling is self-selected, so it stops
  accidents but not a caller who omits it; moving it into the token needs a
  dedicated `mcp` auth scope and an issuing flow.
- `undelete` tools. `page` and 28 other services expose `UndeleteByID` with no
  tool, so there is no way back from a delete through MCP (§4.2).
- Block removal. Page's block merge ships without the `removeX` companion §8.2
  requires, so a caller can add or overwrite a block but never remove one.

**Ruled:**

- Fixing the `Hidden` leak would have removed 7 tools from remote clients: the
  four `automation_taq_*`, both `automation_workflow_*`, and
  `compose_page_block_schema`. They had been reachable over HTTP only because
  of the bug. Ruled: all 7 stay visible, now deliberately rather than by
  accident. `Hidden()` consequently has no callers today — it is kept because
  the distinction is real and the next in-process-only tool will need it.
- `automation_taq_exec` now reaches the `buildResource` case that was dead
  behind the policy trap, so it is resource-scoped and requires an allow entry.
  Agents relying on it without one will start being denied. This is the case
  behaving as designed for the first time, not a new restriction.
