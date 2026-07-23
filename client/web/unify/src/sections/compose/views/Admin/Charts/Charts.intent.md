---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useChartStore.js
  - client/web/unify/src/sections/compose/lib/charts.js
  - client/web/unify/src/sections/compose/components/Chart
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Charts admin

## Intention

Manage a namespace's charts: list them, create one of four chart categories
(generic, funnel, gauge, radar), and edit its report definition with an
immediate visual preview.

## Data touched

- `useChartStore` — findByID/create/update/delete keep the shared set in sync.
- `$ComposeAPI.chartListCancellable` (paginated list) and `recordReport` (the
  preview's data source).
- Chart classes from `@planetcrust/human-js` (`compose.Chart` and subclasses);
  `../../../lib/charts` `chartConstructor` re-types raw API charts.
- RBAC: namespace `canCreateChart`/`canExportCharts`, per-chart `can*` flags;
  permission resource `corteza::compose:chart/<namespaceID>/<id|*>`.

## Map

- `List.vue` — paginated chart list + type-selector create dialog.
- `Edit.vue` — settings/report editor with live preview.

## When changing this

- Chart category is fixed at creation (query param → subclass); the editor
  picks its report sub-editor by instance type, so constructing with the wrong
  class silently degrades to the generic editor.
- Export is client-side JSON (`{ type: 'chart', list: [...] }`) — keep the shape
  importable.
