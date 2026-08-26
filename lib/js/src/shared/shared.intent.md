---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/cast.ts
touched-by:
  - lib/vue
  - client/web/unify
  - lib/js/src/corredor
tests: []
---

# Shared types

## Intention

Types used by all three subsystems (compose, system, automation) that belong
to none of them.

## Map

- `types/attachment.ts` — `Attachment` class: file attachments with URL/preview handling. Only the corredor helpers construct it today; the file field/page-block carry bare attachment IDs and their viewer rebuilds the same preview/download URLs by hand.
- `types/chart/colorschemes` — the bundled chart color scheme catalog; `types/chart/helper.ts` — `getColorschemeColors()` resolves a scheme name to a color array. A name containing `custom` is looked up in the caller's `customColorSchemes` list instead of the bundled tables; anything that resolves to nothing — no name, a custom scheme since deleted, a built-in name that never existed — comes back as the default 13-colour palette.

## When changing this

- Color scheme names are persisted in chart/page-block configs — removing or renaming a scheme breaks saved charts; only add.
- `getColorschemeColors` must never hand back an empty or undefined palette. echarts takes one as "no colours" and draws the chart's legend with none of its series, which reads on screen as an empty result set rather than as a missing scheme. `helper.test.ts` holds every branch that used to do it.
