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

All handles/slugs `agent-` prefixed. Base: `dev/agent/api.sh`.

1. Namespace: `POST /compose/namespace/ {"name", "slug", "enabled": true}`
2. Modules: `POST /compose/namespace/{ns}/module/` with
   `{"name", "handle", "fields": [...]}` — fields is an **array**:
   `{"name", "label", "kind", "options": {}}`.
   Field kinds: `String` (default), `Email`, `Url`, `Number`, `DateTime`,
   `Bool`, `Select` (`options: {"options": [{"value","text"}]}`),
   `Record` (`options: {"moduleID": "<id-as-string>"}`), `User`, `File`.
   Multi-value: `"multi": true` on the field.
   REST-created modules register DAL models live — records work immediately.
3. Verify: GET the modules back; POST one probe record per module, then
   delete it, before bulk-creating data.

## Phase 3 — Records

`POST /compose/namespace/{ns}/module/{mod}/record/` with
`{"values": [{"name": "<field>", "value": "<string>"}]}` — every value a
string; Record-ref values are the target recordID string. Create referenced
records first, keep an ident→recordID map for refs.

## Phase 4 — Charts + pages (use pagebuild, don't hand-roll)

Write a spec JSON and apply it with
`python3 dev/agent/pagebuild.py <namespace-slug> <spec.json>` (idempotent
upsert by handle; format documented in the script header). It resolves
`{"module": "<handle>"}` / `{"chart": "<handle>"}` refs to IDs.

Key shapes (ground truth: `client/web/unify/src/sections/compose/components/PageBlocks/Blocks/*.vue`):

- Blocks live on a **48-column grid** (`Grid.vue`: `COLS = 48`, cell height
  10px, defaults w=24 h=18): `xywh: [x, y, w, h]` — required. Full-width
  block = w 48; a tall list ≈ h 36.
- `RecordList`: `options: {"module": "<handle>"}`.
- `Chart`: `options: {"chart": "<chart-handle>"}` (chart resource; blocks
  with invented options like `chartKind` render nothing).
- `Record` (record pages only): `options: {}` shows all module fields, or
  `{"fields": ["a","b"]}` for a subset.
- `Metric`: `options: {"metrics": [{"label", "module": "<handle>",
"metricField": "count", "operation": "", "filter": ""}]}` — per-metric
  moduleID; for sums use `"metricField": "<numField>", "operation": "sum"`.
- Record page = page with `"module": "<handle>"` + a Record block,
  `"visible": false`. Dashboards: `"visible": true`, optional
  `"icon": "font-awesome://<name>"`, `"weight"` for nav order.
- Chart resource config:
  `{"reports": [{"module": "<handle>", "filter": "", "dimensions":
[{"field", "modifier": "(no grouping / buckets)", "conditions": {}}],
"metrics": [{"field": "count", "type": "doughnut|bar|line|pie"}]}],
"colorScheme": "tableau.Tableau10"}`.

**Never create pages via envoy YAML import** — block refs stay unresolved
(handles instead of IDs) and the pages are broken in the UI.

## Phase 5 — Automation / TAQs

Workflows and TAQs are API-created only (not envoy-importable here). Check
`server/automation/rest/` and TAQ handlers for payload shapes when needed —
verify with a GET of an existing resource before inventing shapes.

## Phase 6 — Verify + hand over

- Re-fetch pages/charts/records via API; confirm every block option holds a
  real ID (a handle string left in options = broken).
- Tell the human what to eyeball in the webapp (namespace name, pages).
- If the system should persist as a fixture: put modules+records in
  `dev/fixtures/<slug>/def.yaml` + CSVs, presentation in `ui.json`
  (see `/dev-seed`). Records export: `export compose-namespace <slug>` gives
  module YAML; pages must stay in ui.json.
