---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/lib/charts.js
  - client/web/unify/src/sections/compose/components/Common/InterpolationFootnote.vue
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/ChartBlock.vue
  - client/web/unify/src/sections/compose/views/Admin/Charts/Edit.vue
tests: []
---

# Chart components

## Intention

Everything chart-shaped in compose: the runtime renderer shared by the chart
page block and the chart builder preview, plus the builder's report/config
editors.

## Data touched

ChartRenderer never fetches report data itself — the consumer passes a
`reporter` function (so blocks and the builder can scope/filter differently).
Chart model classes and echarts option assembly come from
`lib/charts` (`chartConstructor`) and `CChart` from `@planetcrust/human-vue`.

## Map

- ChartRenderer.vue — props `chart` (model) + `reporter` (+ optional `record` for prefilter vars); renders via CChart, emits `updated` and `drill-down` on datapoint click
- Report/ReportEdit.vue — builder editor for one report: source module, preset/custom filters (renders Common's `InterpolationFootnote` with `dependsOnPlacement`, since a chart is namespace-level and may end up placed on a record or non-record page), and the per-chart-type option editor
- Report/GenericChart.vue, FunnelChart.vue, GaugeChart.vue, RadarChart.vue — per-chart-type metric/dimension config editors, exported through Report/index.js and dispatched by chart type

## When changing this

The renderer/`reporter` split is the contract: keep data fetching outside this
folder. Drill-down emits must keep working from both the block and magnified
views. A new chart type touches lib/charts (model), a Report/ editor, and the
Report/index.js export.
