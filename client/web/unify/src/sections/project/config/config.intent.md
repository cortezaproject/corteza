---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/project/stores
  - client/web/unify/src/sections/project/views
  - client/web/unify/src/sections/project/components
tests: []
---

# Project config (locked definitions catalog)

## Intention

Declarative, single-source-of-truth definitions the whole section renders
from. Per the section's locked contracts these definitions are locked: the
step mechanism and step set, the resource kind catalog, and the dashboard
view set. Display strings are i18n keys resolved with `$t`; enum values that
persist (statuses, AI Act options) stay literal.

## Map

- `pipeline.js` — the build pipeline: STEPS (key/type/tab/kind), the Build vs
  Govern tab split (Manage & Monitor has no steps), `kindsThroughStep`
  (graph gating), and `PUBLISH_GOVERNANCE_STEP_KEY` — publish is a topbar
  action, not a wizard step.
- `kinds.js` — per-resource-kind visual config (icon/colors), RESOURCE_KINDS
  (pipeline order), OVERVIEW_KINDS, graph LAYERS (data/logic/experience),
  NODE_LAYER_KINDS, ACCESS_KINDS (role/user ride the access overlay).
- `roles.js` — fixed member role presets with capability flags
  (read/write/requestApproval/grantApproval); flags gate review actions,
  never tab/step visibility (AI Act Art. 17(m) accountability framework).
- `sensitivity.js` — the standard classification scheme the store seeds as
  real DAL sensitivity-level resources; `id` is the resource handle.
- `fieldTypes.js` — module field types, 1:1 with compose field kinds; labels
  reuse compose's shared `general.fieldKinds.*` keys.
- `connectors.js` — the connection catalog (catalogID + brand + tags),
  mirroring the Admin connection catalog.
- `dashboard.js` — the dashboard left-rail nav (locked view set: Overview,
  Events, five Categories, Backlog, Reports stub).
- `categories.js` — per-category dashboard config: columns, charts, KPIs,
  badges, trend grouping; reuses eventForm's schemas and visual identity.
- `eventForm.js` — New Event form schemas per category + EVENT_STATUS;
  consumed by the GovernanceForm renderer.
- `eventKinds.js` — action-log resource type → project kind mapping so audit
  events wear their resource's icon/color; unmapped types fall back neutral.
- `chartColors.js` — validated chart palettes (status/severity/risk/…) with
  ordering helpers; single hex per label must pass light AND dark checks.
- `trend.js` — day/week/month bucketing, range presets, `adaptiveWindow`.
- `summaryForm.js` — Project Summary governance step schema.
- `resourceManagementForm.js` — Resource Management step values shape +
  defaults merge (Art. 17(l)); custom-rendered, no GovernanceForm schema.

> **WIP:** Govern/FRIA content — `summaryForm.js`, `resourceManagementForm.js`
> and the Govern step entries in `pipeline.js` are governance CONTENT and due
> to change; only the step mechanism, tab shapes and kind catalogs are locked.

## When changing this

- Adding a resource step to STEPS must follow the locked dialog/permissions
  standards; kinds join graphs/metrics automatically via kinds.js.
- Re-run the dataviz palette validator when any chartColors hex changes.
