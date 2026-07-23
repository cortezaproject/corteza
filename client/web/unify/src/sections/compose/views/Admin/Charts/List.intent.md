---
kind: file
covers: List.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/composables/useResourceList.ts
  - lib/vue/src/stores/useChartStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Chart List view

## Intention

Browse the namespace's charts and reach their editors; start a new chart by
choosing its category; manage export, permissions, and deletion from the list.

## UX capabilities

- Paginated, searchable, sortable table (name/handle/changedAt) scoped to the namespace; rows open the editor only when the chart is updatable or deletable.
- Create opens a category picker dialog (generic, funnel, gauge, radar) and forwards the choice as `?category` to the create route.
- Export as JSON: per-chart from the row menu or all listed charts, gated by namespace `canExportCharts`; wildcard permissions button gated by `compose/` `grant`.
- Row menu: permissions (per-chart `canGrant`), export, delete (confirm + store delete + list refetch).

## Routes

`admin.charts` at `admin/charts` under `namespace.view`; navigates to `admin.charts.create` (with `?category`) and `admin.charts.edit` (`:chartID`).

## When changing this

- The category values must match the subclass switch in `Edit.vue` (`funnel`/`gauge`/`radar`, empty = generic).
- Permission resource includes the namespace: `corteza::compose:chart/<namespaceID>/<chartID>`.
