---
name: sys-design
description: Design and build a full system inside Human on the local dev server — namespace, data model, records, charts, pages — from a spec or requirements. Use when asked to model/build an application, demo system, or prototype in Human (not for changing Human's own source code).
---

# /sys-design <what to build>

Build order matters; each step is verified via API before the next. Requires
the dev/agent toolkit (`/dev-api` skill) — run `dev/agent/smoke.sh` first.

## Phase 1 — Design on paper

From the requirements, write down: modules with fields (+ kinds and Record
refs between them), what records exemplify the data, which charts/metrics
matter, and what pages the user needs (dashboards + one record page per
module). Confirm the design with the human if scope is ambiguous.

## Phase 2 — Data model (REST, never envoy YAML for live builds)

No name prefix: `api.sh`/`mcp.py` record what they create in
`.state/created.jsonl` and `cleanup.sh` deletes only that. Base:
`dev/agent/api.sh`.

**Everything you name is snake_case** — namespace slugs, module and page and
chart and TAQ handles, and every module field name. Underscores only, never
hyphens or dots. A hyphen is the subtraction operator everywhere an identifier
is parsed: a field named `close-date` lexes as `close` minus `date`, so any
prefilter, presort, chart dimension or `record.values.` expression naming it
breaks. Neither the API nor the webapp validator stops you creating one
(`FieldNameValidator` in `lib/js/src/compose/types/module-field/base.ts` admits
`-`), so the discipline has to come from here.

1. Namespace: `POST /compose/namespace/ {"name", "slug", "enabled": true}`
2. Modules: `POST /compose/namespace/{ns}/module/` with
   `{"name", "handle", "fields": [...]}` — fields is an **array**:
   `{"name", "label", "kind", "options": {}}`.
   Field kinds and the options that matter:

   | kind                               | options                                                                                                        |
   | ---------------------------------- | -------------------------------------------------------------------------------------------------------------- |
   | `String` (default), `Email`, `Url` | —                                                                                                              |
   | `Number`                           | `{"precision": 0}` storage rounding — **display is `format`**: `{"format": "0,0.00", "suffix": " EUR"}`        |
   | `DateTime`                         | `{"onlyDate": true}` (or `onlyTime`)                                                                           |
   | `Bool`                             | —                                                                                                              |
   | `Select`                           | `{"selectType": "default", "options": [{"value","text"}]}`                                                     |
   | `Record`                           | `{"moduleID": "<id-as-string>", "labelField": "<field>", "queryFields": ["<field>"], "selectType": "default"}` |
   | `User`                             | `{"selectType": "default"}`                                                                                    |
   | `File`                             | —                                                                                                              |

   `precision` rounds what is **stored** and has no effect on what is shown:
   `formatValue` reads `format`/`presetFormat`, `prefix` and `suffix` and never
   looks at it (`lib/js/src/compose/types/module-field/number.ts:99`). A money
   field declared with `precision: 2` alone stores `245.00`, echoes `"245"` and
   renders `245` — cents gone, no error at any layer. Give it a numeral.js
   `format` as well.
   `selectType` is one of `default` \| `multiple` \| `each` for Select, Record
   and User. On a Record field, **set `labelField`** — without it the picker
   and every viewer fall back to showing the raw record ID.
   Flags are `isRequired` and `isMulti` — spelled that way, not `required` /
   `multi` (`server/compose/types/module_field.go:40`). A misspelt flag is
   accepted and silently ignored.
   REST-created modules register DAL models live — records work immediately.

3. Verify: GET the modules back; POST one probe record per module, then
   delete it, before bulk-creating data.

## Phase 3 — Records

`POST /compose/namespace/{ns}/module/{mod}/record/` with
`{"values": [{"name": "<field>", "value": "<string>"}]}` — every value a
string; Record-ref values are the target recordID string. Create referenced
records first, keep an ident→recordID map for refs.

**Reading values back: a falsy value has no `value` key at all.** A `Bool`
set to `"0"` returns `{"name": "done"}` — the key is absent, not empty. So
`{v["name"]: v["value"] for v in record["values"]}` raises `KeyError` on the
first false checkbox and takes the whole run down mid-way. Use `v.get("value")`.
The **`name` key is still there** — only `value` is gone — so a write-back check
that looks for missing _names_ passes on a value that never landed. Compare on
`v.get("value")`, never on key presence.

**Filtering records: the list endpoint takes `query`, not `filter`.** Only
`/record/report` takes `filter` (`server/compose/rest/request/record.go:590`
vs `:691`). Passing `filter` to the list endpoint is not an error — unknown
params are dropped, so the call returns **the whole module** and reads as a
filter that matched everything. Confirm any filter with a negative control:
a query that should match nothing must come back empty.

## Alternative surface: Human's own MCP tools

`dev/agent/mcp.py tools|schema|call` reaches the server's first-party MCP
tools (compose namespace/module/record/page/chart CRUD, TAQ/workflow exec).
They resolve names/handles server-side and are the product's official agent
surface — prefer them when they cover the operation; REST/api.sh covers the
rest. Known caveat: namespaces created via CLI import may not resolve by
slug through MCP lookup or for some users (server bug, logged).

## Phase 4 — Charts + pages (use pagebuild, don't hand-roll)

Write a spec JSON and apply it with
`python3 dev/agent/pagebuild.py <namespace-slug> <spec.json>` (idempotent
upsert by handle; format documented in the script header). It resolves
`{"module": "<handle>"}` / `{"chart": "<handle>"}` refs to IDs.

**Block option shapes: query the server, don't trust prose.** The canonical
contract is served by the MCP schema tool (CI-enforced against the actual
webapp classes via `lib/js/.../page-block/schema-contract.test.ts`):

```sh
dev/agent/mcp.py call compose_page_block_schema '{"kind":"Metric"}'
dev/agent/mcp.py schema compose_chart_create   # chart config contract
dev/agent/mcp.py schema compose_page_create    # grid + page-type guidance
dev/agent/mcp.py call system_skill_lookup '{"tool":"compose_page_update"}'
```

**Read `system_skill_lookup` before using an unfamiliar write tool.** It holds
the semantics a parameter list cannot — the page/layout split among them — and
a tool description can be wrong where the skill is right: `compose_page_update`
documents `xywh` on its `blocks` param as moving a block, and it does not (the
same tool's own summary says "existing blocks are never re-placed"). Moving a
block needs `compose_page_layout_update`.

A Metric block's tile is one entry in `options.metrics`, needing at least
`{"moduleID", "metricField", "operation", "filter"}`. `operation` is `sum` \|
`max` \| `min` \| `avg`; counting records is `"metricField": "count"` with
`operation` left `""`. The block's own `title` and the metric's `label` both
render, so setting both prints the tile's name twice.

What the schemas cannot express (layout semantics):

- **48-column grid**, cell height 10px; `xywh` required in pagebuild specs.
  Those minima are the floor at which a block is not _broken_ (Metric h≥20,
  RecordList/Chart h≥30) — they are not a usable size. Budget instead:
  - a **RecordList** spends ~19 cells on toolbar, header and pager and ~4.6
    per row, so `h ≈ 19 + 4.6 × rows`. h=30 shows **two** rows; 4 rows needs
    h≈38, 8 rows h≈56.
  - **width** holds far fewer columns than 48 suggests: a text column renders
    140–190px against ~24px per grid column, so `w=48` fits about 5–6 fields
    and `w=24` about 3. Content decides it, so treat this as a starting guess
    and let `verify-ui.mjs` measure the truth — it reports the exact pixels
    cut off.
    Getting either wrong is silent: the block renders, just short or clipped.
- **Lay pages out as rows**: blocks sit side by side by stepping x and
  keeping y — four tiles are `[0,0,12,20]`, `[12,0,12,20]`, `[24,0,12,20]`,
  `[36,0,12,20]`. Full width (48) is for lists and forms; half (24) suits
  charts, a quarter (12) suits Metric/Progress tiles. Stepping y for every
  block leaves the page empty down its right-hand side.
- Record page = page with `"module": "<handle>"` + a Record block (`{}` =
  all fields), `"visible": false`. Its URL takes a record:
  `/compose/namespace/<slug>/pages/<pageID>/records/<recordID>` — **`records`
  plural**, unlike the API's `/record/`. Verifying one without a recordID
  silently renders the empty-state page. Dashboards: `"visible": true`,
  `"weight"` for nav order. **No nav icons** — the webapp draws one as an
  image, so a `font-awesome://` icon renders as a broken image; icons are
  for a human to upload in the page editor.
- In pagebuild specs, `{"module"/"chart": "<handle>"}` are resolved to
  `moduleID`/`chartID`; via MCP/REST you pass real IDs yourself.
- **The page's primary layout owns the geometry, not `page.blocks`.** The
  server derives a layout on first create, so a new page agrees with its spec;
  a later page write does not touch the layout, so an edited `xywh` applies to
  nothing. `pagebuild.py` syncs both and says `N block(s) re-placed`. Building
  pages through MCP/REST instead means writing the layout yourself:
  `POST /compose/namespace/{ns}/page/{pageID}/layout/{pageLayoutID}`.
- A page with **one** block ignores `xywh` entirely — it renders through a
  flex wrapper that fills the view, so height there is neither honoured nor
  worth tuning.
- Related lists on a record page filter with `prefilter`, and the record
  variable is `${recordID}` — `${record.values.<field>}` and `${ownerID}`
  also interpolate. `${record.recordID}` is **not** a thing and silently
  yields nothing.

**Never create pages via envoy YAML import** — block refs stay unresolved
(handles instead of IDs) and the pages are broken in the UI.

## Phase 5 — Automation / TAQs

Workflows and TAQs are API-created only (not envoy-importable here). TAQs live
at **`/automation/ng-automation/`** — POST to create, **PUT
`/automation/ng-automation/{automationID}`** to update, and
`GET .../{id}/executions` + `GET .../{id}/execution/{eid}/trace` to inspect a
run. (The `/automation/{workflows,triggers}/` paths are the _workflow_ API and
do not reach a TAQ.)

**Default to a TAQ. Reach for a workflow only when a TAQ provably cannot do
it.** A TAQ create/update returns `issues`, so many malformed automations say
so at authoring time; a workflow stores clean and fails silently at run time,
visible only in `logs.sh` under `workflow.session.exec`. But **`issues` ride a
200 and `api.sh` exits 0** — a TAQ the runtime has refused looks exactly like a
working one unless you parse `issues[].severity == "error"` yourself. Do that
on every write. (`runnable` is an MCP-layer field; REST responses never carry
it, and read/list only echo the `issues` persisted at the last write.)
A TAQ also holds triggers + steps in one resource, supports per-step
`maxRetries`/`recoverable`, and is what the MCP tooling targets. Evaluate
before building — a TAQ can only do what the construct library covers:

```sh
dev/agent/api.sh GET /automation/construct-library/functions   # 19 step refs
dev/agent/api.sh GET /automation/construct-library/triggers    # 22 rt/et pairs
dev/agent/mcp.py schema automation_taq_create                  # full contract
```

If the function or trigger the task needs is absent there, say so and ask the
human before falling back — a workflow reaches the full 93-entry registry
(`GET /automation/functions/`), but 74 of those are unavailable to a TAQ, so
the fallback is a real trade, not a formality.

TAQ gotchas that cost real time:

- **A trigger constraint is `{"name", "op", "values":[{"@type","@value"}]}`.**
  The type lives on each value, not beside the name:

  ```json
  { "name": "module", "op": "eq", "values": [{ "@type": "Handle", "@value": "tickets" }] }
  ```

  Get the shape wrong and it fails silently, always **closed** — the trigger
  registers and never fires, with `issues: null`:
  - Names must be **bare**. `prepConstraintBits`
    (`server/automation/service/ng_automation.go:905`) appends the field from
    the value type — `String`→`.name`, `Handle`→`.handle`, `ID`→`.id` — so
    `{"name":"module"}` with a `Handle` matches `module.handle`, while
    `{"name":"module.handle"}` with a `String` builds `module.handle.name`.
  - An **empty `values` array** (what you get by copying a `{"name","@type"}`
    shape) makes the constraint match nothing at all.

- **Only `String`, `Handle` and `ID` work as constraint types — and using any
  other fails OPEN.** `GET /automation/construct-library/triggers` advertises
  `ComposeNamespace`, `ComposeModule` and `ComposeRecord`; `prepConstraintBits`
  handles none of them, its caller logs the failure at **Debug** and
  `continue`s (`:868`), and the constraint is dropped _entirely_. The trigger
  then registers **broader than you wrote it** — a `compose:record` trigger
  meant for one module fires on every record create in every namespace, and a
  write step reaches data you never intended to touch. `issues` stays null.
  Use `Handle` or `ID`, and confirm any near-miss case does not fire.
- Step kinds are exactly `function`, `iterator`, `gatewayExclusive`,
  `gatewayInclusive`, `termination`, `error`. There is **no `expressions`
  step** — that is a workflow kind.
- Arguments bind by `argumentName` (workflows use `target`), and `type` must
  be spelled exactly as the parameter's `types` array lists it.
- **An aggregate parameter uses BOTH `argumentName` and `target`.** For a
  parameter with `Aggregate: true` — `composeRecordsUpdate.values`, whose
  construct-library `types` read `KV|KVV|Any` as though a map were wanted —
  the record field goes in `target`, one argument per field:

  ```json
  { "argumentName": "values", "target": "note", "type": "String", "expr": "..." }
  ```

  `argumentName` must be the bare parameter name: `VerifyArguments`
  (`server/automation/types/param.go:100`) rejects anything else, so
  `"values.note"` errors with `unknown parameter values.note is used`.

- **Always author `paths` explicitly; never send `[]`.** The create contract
  says a lone trigger and lone step are wired together for you, and they are —
  but only in the runtime registration, never in what is stored. `paths` stays
  empty, so the TAQ runs correctly while the builder canvas draws two
  disconnected chains, each ending in its own `End`, which reads as broken to
  anyone who opens it. Send `[{"parentID":"<triggerID>","childID":"<stepID>"}]`.
- A TAQ **does** have `invoker` and `runner` in scope, despite `EncodeVars`
  (`server/compose/service/event/events.gen.go`) listing only `record`,
  `oldRecord`, `module`, `namespace`, `recordValueErrors`: `injectIdentities`
  (`server/automation/service/ng_automation.go:322`) adds them to the input, so
  `{"expr": "invoker.email"}` resolves. That injection is TAQ-only — in a
  workflow, "who did this" is still `record.createdBy`.
- **A trace is in-memory and expires within minutes.** Once it has, the trace
  endpoint returns `execution not found`, which reads as a wrong ID while
  `GET .../executions` still lists the run — so read a trace in the same breath
  as the probe that produced it. That endpoint also answers **`text/plain` on
  error**, unlike every other one, so use `api.sh --json` for it. And
  `GET .../executions` returns a **bare array** in `response`, not the
  `{"set": [...]}` every compose list endpoint returns.

Always fire a real probe record and confirm the effect (the notification row,
the updated field), then check a near-miss case does _not_ fire. Delete the
probe records afterwards.

**Check `GET .../{id}/executions` as part of that, not just the effect.** The
record write cannot tell you the automation worked: `compose/service/record.go`
discards the eventbus error (`_ = svc.eventbus.WaitFor(...)`), so a step that
blew up mid-`afterCreate` still returns HTTP 200 and a clean record. A probe
that silently no-opped is indistinguishable from one that never triggered
unless you look at the execution list.

**Render-verify the graph too** — a correct-at-runtime automation can still be
drawn wrong, and the API cannot see it:
`node dev/agent/verify-ui.mjs '/taq/builder/<automationID>'`, then Read the
screenshot and confirm one connected chain from trigger to `End`.

## Phase 6 — Verify + hand over

- Re-fetch pages/charts/records via API; confirm every block option holds a
  real ID (a handle string left in options = broken).
- Prove filters bind rather than assuming: every prefilter and metric filter
  needs a case that must **not** match, checked to come back empty.
- **Render-verify in a real browser** (standard final gate — API checks
  cannot see clipped blocks, broken charts, or raw IDs in the UI):
  `node dev/agent/verify-ui.mjs '/compose/namespace/<slug>/pages/<pageID>'`.
  It prints a **text report per path** before naming the screenshot: block
  geometry, blocks whose content is clipped (with the `h` to add), tables with
  columns cut off, empty blocks, raw IDs where a label belongs, and
  uninterpolated `${...}`. Iterate on that report — it is the cheap loop, and
  it catches the faults that actually recur. `OK` with no findings means the
  page is sound; **read the screenshot when it reports a finding, when the
  check is about visual design, or once at the end** to confirm the thing
  looks right. Passing several paths in one call is one browser launch.
  `note:` lines flag a list scrolling inside its block. That is normal only
  when the list genuinely holds more rows than `perPage`; on a block sized
  below its own row count it means rows are hidden behind an inner scrollbar,
  so check the count against the budget above before dismissing one.
- Tell the human what to eyeball in the webapp (namespace name, pages).
- If the system should persist as a fixture: put modules+records in
  `dev/fixtures/<slug>/def.yaml` + CSVs, presentation in `ui.json`
  (see `/dev-seed`). Records export: `export compose-namespace <slug>` gives
  module YAML; pages must stay in ui.json.
