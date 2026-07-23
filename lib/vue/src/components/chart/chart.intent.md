---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by: []
tests: []
---

# chart/

## Intention

The single shared ECharts wrapper so every app renders charts with the same
renderer, registered chart types, and theme handling.

## Map

- `CChart.vue` — vue-echarts `<v-chart>` wrapper; `chart` prop is a raw
  ECharts option object passed through untouched.

## Contracts consumers rely on

- Tree-shaken ECharts: only registered types render (bar, line, pie, scatter,
  radar, funnel, gauge + title/tooltip/legend/grid/dataset/toolbox/dataZoom).
  A new chart kind means registering its module here first.
- Theme is chosen from `chart.darkMode` on the option object — callers embed
  the flag rather than passing a separate prop.
- Autoresizes to its container; the instance is disposed on unmount.

## When changing this

The `chart` prop is a pass-through contract: consumers hand-build ECharts
options, so never transform or "fix up" the option object here.
