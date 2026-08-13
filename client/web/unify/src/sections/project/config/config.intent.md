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

- `pipeline.js` — the build pipeline: STEPS (key/type/tab/kind), `resolveSteps`
  /`stepsForTab`, and `PUBLISH_GOVERNANCE_STEP_KEY` — publish is the Publish
  tab's own flow, not a wizard step.
- `kinds.js` — per-resource-kind visual config (icon/colors), RESOURCE_KINDS
  (pipeline order) and OVERVIEW_KINDS (the graph's initial visible set).
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
  Board, Activity, five Categories, plus an insights group: Backlog, Reports).
- `manageNav.js` — the M&M left-rail sections; the tab renders from this,
  never from STEPS. Meant to hold the SAME section set as `dashboard.js`
  (ruled 2026-07-28: one surface, two scopes), differing only in navigation
  (child routes there, `?section=` here).
- `aiSystemClasses.js`, `friaDeterminationForm.js`, `friaScenario.js` — the
  Govern tab's AI-system and FRIA vocabularies and form schemas.
- `publishState.js` — `chainHasPublished`/`projectStatusTag`: the one place
  lifecycle status becomes a routing decision or a chip.
- `resourceRefs.js` — A PERSISTED INTERFACE (see its header): how a resource
  declares the refs the graph and the revision clone walk.

- `categories.js` — per-category dashboard config: columns, charts, KPIs,
  badges, trend grouping; reuses eventForm's schemas and visual identity.
- `eventForm.js` — New Event form schemas per category + EVENT_STATUS;
  consumed by the GovernanceForm renderer. All five categories share one
  four-value status set (review included) so the M&M board's columns hold
  every item type. `REVISION_FIELD` assigns an item to a revision; like the
  `user()` factory it carries a `source` the consuming dialog resolves at
  render time, keeping GovernanceForm option-agnostic. Shown on edit always;
  on create only where the surface isn't already implying a revision.
- `friaTaxonomies.js` — the FRIA regulatory taxonomies (Charter rights by
  chapter, vulnerable groups, AI harm vectors, trigger conditions, impacted
  parties). Entry `key`s are persisted the moment a scenario stores a
  selection — treat them as an interface, never rename to relabel.
- `eventKinds.js` — action-log resource type → project kind mapping so audit
  events wear their resource's icon/color; unmapped types fall back neutral.
- `chartColors.js` — validated chart palettes (status/severity/risk/…) +
  ordering helpers; one hex per label passes light AND dark checks. Status
  mirrors EventBadge's hue FAMILIES, never its hexes — pills are
  theme-aware, marks are one hex (unification reverted 2026-07-29).
- `trend.js` — day/week/month bucketing, range presets, `adaptiveWindow`.
- `summaryForm.js` — Project Summary governance step schema.
- `resourceManagementForm.js` — Resource Management step values shape +
  defaults merge (Art. 17(l)); custom-rendered, no GovernanceForm schema.

> **WIP:** Govern/FRIA content — `summaryForm.js`, `resourceManagementForm.js`
> and the Govern step entries in `pipeline.js` are governance CONTENT and due
> to change; only the step mechanism, tab shapes and kind catalogs are locked.

> **DRIFT:** `manageNav.js` is missing dashboard.js's whole insights group
> (Backlog, Reports), so the two section sets are not in step.

## When changing this

- Adding a resource step to STEPS must follow the locked dialog/permissions
  standards; kinds join graphs/metrics automatically via kinds.js.
- Re-run the dataviz palette validator when any chartColors hex changes.
