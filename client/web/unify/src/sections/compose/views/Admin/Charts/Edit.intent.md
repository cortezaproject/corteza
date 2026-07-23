---
kind: file
covers: Edit.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useChartStore.js
  - client/web/unify/src/sections/compose/components/Chart/ChartRenderer.vue
  - client/web/unify/src/sections/compose/components/Chart/Report
  - client/web/unify/src/sections/compose/lib/charts
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Chart Edit view

## Intention

Define a chart: general settings (name, handle, color scheme, animation,
toolbox), the category-specific report (module, dimensions, metrics), and see
the result immediately in a live preview fed by real record data.

## UX capabilities

- Report sub-editor chosen by chart class (generic/funnel/gauge/radar); the draft report is provided to it via `reportDraft`.
- Live preview renders through `ChartRenderer` with `$ComposeAPI.recordReport` as reporter; manual refresh disabled until every report has a module.
- Name required, handle validated; save hidden when the user lacks `canUpdateChart`; save-as-copy clones with blank handle; delete returns to the list.
- Edit mode extras: translator in the topbar, export as JSON, per-chart permissions button; unsaved-changes guard throughout.

## Routes

`admin.charts.create` at `admin/charts/create` (`?category` selects the subclass) and `admin.charts.edit` at `admin/charts/:chartID/edit`. Create navigates to the edit route after saving; load failure bounces to `admin.charts`.

## When changing this

- Loaded charts must pass through `chartConstructor` — a plain object loses the subclass and the wrong report editor mounts.
- Report editors receive `modules` from the preloaded module store; they mutate the draft in place via the `editReport` computed.
- Manual preview refresh intentionally disables animation before re-keying the renderer.
