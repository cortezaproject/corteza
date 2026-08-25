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
   | `Number`                           | `{"precision": 0}`                                                                                             |
   | `DateTime`                         | `{"onlyDate": true}` (or `onlyTime`)                                                                           |
   | `Bool`                             | —                                                                                                              |
   | `Select`                           | `{"selectType": "default", "options": [{"value","text"}]}`                                                     |
   | `Record`                           | `{"moduleID": "<id-as-string>", "labelField": "<field>", "queryFields": ["<field>"], "selectType": "default"}` |
   | `User`                             | `{"selectType": "default"}`                                                                                    |
   | `File`                             | —                                                                                                              |

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
```

A Metric block's tile is one entry in `options.metrics`, needing at least
`{"moduleID", "metricField", "operation", "filter"}`. `operation` is `sum` \|
`max` \| `min` \| `avg`; counting records is `"metricField": "count"` with
`operation` left `""`. The block's own `title` and the metric's `label` both
render, so setting both prints the tile's name twice.

What the schemas cannot express (layout semantics):

- **48-column grid**, cell height 10px; `xywh` required in pagebuild specs.
  Blocks CLIP silently when too short — Metric h≥20, RecordList/Chart h≥30.
- **Lay pages out as rows**: blocks sit side by side by stepping x and
  keeping y — four tiles are `[0,0,12,20]`, `[12,0,12,20]`, `[24,0,12,20]`,
  `[36,0,12,20]`. Full width (48) is for lists and forms; half (24) suits
  charts, a quarter (12) suits Metric/Progress tiles. Stepping y for every
  block leaves the page empty down its right-hand side.
- Record page = page with `"module": "<handle>"` + a Record block (`{}` =
  all fields), `"visible": false`. Dashboards: `"visible": true`,
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

Workflows and TAQs are API-created only (not envoy-importable here). Check
`server/automation/rest/` and TAQ handlers for payload shapes when needed —
verify with a GET of an existing resource before inventing shapes.

**Default to a TAQ. Reach for a workflow only when a TAQ provably cannot do
it.** A TAQ reports `issues` and `runnable` on every write, so a malformed
automation says so at authoring time; a workflow stores clean and fails
silently at run time, visible only in `logs.sh` under `workflow.session.exec`.
A TAQ also holds triggers + steps in one resource, supports per-step
`maxRetries`/`recoverable`, and is what the MCP tooling targets. Evaluate
before building — a TAQ can only do what the construct library covers:

```sh
dev/agent/api.sh GET /automation/construct-library/functions   # 17 step refs
dev/agent/api.sh GET /automation/construct-library/triggers    # 22 rt/et pairs
dev/agent/mcp.py schema automation_taq_create                  # full contract
```

If the function or trigger the task needs is absent there, say so and ask the
human before falling back — a workflow reaches the full 93-entry registry
(`GET /automation/functions/`), but 78 of those are unavailable to a TAQ, so
the fallback is a real trade, not a formality.

TAQ gotchas that cost real time:

- **Trigger constraint names must be bare, and `@type` picks the field.**
  `prepConstraintBits` (`server/automation/service/ng_automation.go`) appends a
  suffix from the value type: `String`→`.name`, `Handle`→`.handle`, `ID`→`.id`.
  So `{"name":"namespace","@type":"Handle"}` matches `namespace.handle`;
  writing `{"name":"namespace.handle","@type":"String"}` silently builds
  `namespace.handle.name`, matches nothing, and only logs at Debug. The TAQ
  still reports `runnable: true` — a clean write is not a working trigger.
- Step kinds are exactly `function`, `iterator`, `gatewayExclusive`,
  `gatewayInclusive`, `termination`, `error`. There is **no `expressions`
  step** — that is a workflow kind.
- Arguments bind by `argumentName` (workflows use `target`), and `type` must
  be spelled exactly as the parameter's `types` array lists it.
- **Always author `paths` explicitly; never send `[]`.** The create contract
  says a lone trigger and lone step are wired together for you, and they are —
  but only in the runtime registration, never in what is stored. `paths` stays
  empty, so the TAQ runs correctly while the builder canvas draws two
  disconnected chains, each ending in its own `End`, which reads as broken to
  anyone who opens it. Send `[{"parentID":"<triggerID>","childID":"<stepID>"}]`.
- Automation updates are **PUT** `/automation/{workflows,triggers}/{id}` — the
  reverse of compose's POST — and live under `/api/automation/`.
- A `compose:record` automation has **no `invoker` in scope**. `EncodeVars`
  (`server/compose/service/event/events.gen.go`) provides only `record`,
  `oldRecord`, `module`, `namespace`, `recordValueErrors`; "who did this" is
  `record.createdBy`.

Always fire a real probe record and confirm the effect (the notification row,
the updated field), then check a near-miss case does _not_ fire. Delete the
probe records afterwards.

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
  `note:` lines are informational — a list scrolling inside its block is
  normal paging, not a defect.
- Tell the human what to eyeball in the webapp (namespace name, pages).
- If the system should persist as a fixture: put modules+records in
  `dev/fixtures/<slug>/def.yaml` + CSVs, presentation in `ui.json`
  (see `/dev-seed`). Records export: `export compose-namespace <slug>` gives
  module YAML; pages must stay in ui.json.
