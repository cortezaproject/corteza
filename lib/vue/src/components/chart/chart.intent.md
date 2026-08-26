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
- The option reaches echarts untouched, but the **update** is not left to
  vue-echarts to plan: `series`, `xAxis` and `yAxis` are replaced wholesale on
  every option change, everything else merges. A consumer that rebuilds its
  whole option — which is what the compose renderer does — gets the removals it
  expects rather than components merged over by index.
- The palette (`color`) is deliberately not in that list. vue-echarts' own diff
  treats every top-level array as a component list and puts any that lost an
  entry into `replaceMerge`; echarts then refuses `color` as a main type, and
  the update after it lands mid-process. Swapping a 13-colour scheme for a
  10-colour one was enough.
- Autoresizes to its container; the instance is disposed on unmount.

## When changing this

The `chart` prop is a pass-through contract: consumers hand-build ECharts
options, so never transform or "fix up" the option object here. `update-options`
is the exception the contract allows — it says how `setOption` is called, not
what is in the option — and widening it is a change every app renders through.
