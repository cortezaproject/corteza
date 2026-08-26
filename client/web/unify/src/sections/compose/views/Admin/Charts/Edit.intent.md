---
kind: file
covers: Edit.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useChartStore.js
  - client/web/unify/src/sections/compose/components/Chart/ChartRenderer.vue
  - client/web/unify/src/sections/compose/components/Chart/Report
  - client/web/unify/src/sections/compose/lib/charts.js
  - client/web/unify/src/sections/compose/lib/chart-color-schemes.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests:
  - client/web/unify/src/sections/compose/views/Admin/Charts/Edit.colorScheme.test.js
---

# Chart Edit view

## Intention

Define a chart: general settings (name, handle, description, color scheme,
animation, toolbox), the category-specific report (module, dimensions,
metrics), and see the result immediately in a live preview fed by real record
data.

## UX capabilities

- Report sub-editor chosen by chart class (generic/funnel/gauge/radar); the draft report is provided to it via `reportDraft`.
- Live preview renders through `ChartRenderer` with `$ComposeAPI.recordReport` as reporter; manual refresh disabled until every report has a module.
- Description is free text on `meta.description`, kept off `config` so the
  chart's blurb and its rendering definition stay separable.
- The colour scheme picker offers the instance's own schemes first, then the
  bundled tables. With `system/` `settings.manage` it also manages them: add
  from the dropdown's header, edit the selected one from the pencil beside it,
  delete from inside that modal. Each writes the whole
  `ui.charts.colorSchemes` setting and refetches `$Settings`.
- Saving a scheme selects it on the draft and stops there — the chart is the
  author's to save, and committing it here would take every other in-flight
  edit with it. Deleting the one in use clears the draft's selection the same
  way; charts elsewhere keep the dead id and fall back to the default palette.
- Editing a scheme in place changes nothing about the chart, so the preview is
  refreshed through `ChartRenderer`'s exposed `updateChart()`. Where the
  selection does change, its own watcher covers it and a second call races the
  first.
- Name required, handle validated; save hidden when the user lacks `canUpdateChart`; save-as-copy clones with blank handle; delete returns to the list.
- Edit mode extras: translator in the topbar, export as JSON, per-chart permissions button; unsaved-changes guard throughout.

## Routes

`admin.charts.create` at `admin/charts/create` (`?category` selects the subclass) and `admin.charts.edit` at `admin/charts/:chartID/edit`. Create navigates to the edit route after saving; load failure bounces to `admin.charts`.

## When changing this

- Loaded charts must pass through `chartConstructor` — a plain object loses the subclass and the wrong report editor mounts.
- Report editors receive `modules` from the preloaded module store; they mutate the draft in place via the `editReport` computed.
- Manual preview refresh intentionally disables animation before re-keying the renderer.
- Nothing here validates a scheme's colours against the bundled tables: an id
  the picker cannot resolve is lib/js's to answer, and it answers with the
  default palette.
