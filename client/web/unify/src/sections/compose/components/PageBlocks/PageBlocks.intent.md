---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/lib/record-filter.js
  - client/web/unify/src/sections/compose/components/Common/InterpolationFootnote.vue
touched-by:
  - client/web/unify/src/sections/compose/views/Pages/View.vue
  - client/web/unify/src/sections/compose/views/Pages/RecordView.vue
  - client/web/unify/src/sections/compose/views/Admin/Pages/Builder.vue
tests: []
---

# Page blocks

## Intention

The building blocks of compose pages: one runtime component per block kind, one
configurator per kind for the page builder, and the gridstack-based layout that
hosts them. This is the extension point for adding new page content types.

## Data touched

Blocks read module/record/page data via `@planetcrust/human-vue` stores and the
injected Compose API; they render inside pages, never fetch page structure
themselves. Auth is injected under the key `$Auth` (not `$auth`).

## Map

- registry.ts — kind → async Block component; `resolveBlock(kind)` is the only dispatch path for rendering
- Blocks/ — runtime components; contract: each accepts props `block, blocks, namespace, page, record` (record only on record pages) and typically wraps content in Blocks/PageBlock.vue (Card chrome: title, refresh, magnify dialog). Calendar (per-feed prefilter), Geometry (per-feed prefilter), Metric (metric filter + comparison custom filter) and Navigation (item/dropdown URLs) interpolate author-typed strings via `lib/record-filter.js#evaluatePrefilter` (`${record.values.x}`, `${recordID}`, `${ownerID}`, `${userID}`, `${user.name}`). Guard contract on a non-record page: a string referencing `${record` or `${ownerID}` makes Calendar/Geometry/Metric skip that feed/metric (`console.warn`), while Navigation falls back to the URL as authored.
- Blocks/MetricBlock.vue — keys fetched values by the metric's configured index (`rtr[mi]`), matching how `formatResponse()` looks them up, so a skipped metric doesn't shift later ones. Its drill-down passes the metric filter **raw** (not pre-interpolated) — the child RecordListBlock it opens interpolates that filter itself; evaluating it here would double-evaluate.
- Blocks/RecordBlock.vue — `provide('$recordContext', activeRecord)` so nested field editors (lib/vue) can interpolate a Record field's prefilter. Its `fieldConditions` (server-evaluated via `$SystemAPI.expressionEvaluate` — gval expressions, not `${}` templates) re-evaluate on every record-value change, debounced 300ms, with a sequence guard so a slower earlier response cannot overwrite a newer one or clear values based on a stale record.
- Configurators/ — per-kind options editors; NOT in the registry — statically imported and dispatched by name in views/Admin/Pages/Builder.vue; they mutate `block.options`. Every input whose value is interpolated renders Common's `InterpolationFootnote`, switched on whether the page is a record page (`!!page.moduleID && page.moduleID !== '0'`).
- Grid.vue — gridstack layout (48 cols, `block.xywh`); view mode (mobile reflow, single-block fills page) vs builder mode (drag/resize, `layout-updated` emit, `item-overlay` slot for the builder toolbox)
- Shared/ — AutomationButtons(+Editor): workflow/script trigger buttons reused by Automation and RecordList blocks
- Blocks/CBulkRecordEditModal.vue, Blocks/Comment/, Blocks/Metric/ — private helpers of their blocks

## When changing this

A new block kind needs all three: Block component + registry entry + Configurator
(plus its import/case in Builder.vue) — miss one and the kind renders the
"no configuration" fallback or is unconfigurable. Block IDs may be temp
(`meta.tempID`) before save; Grid keys/rebuilds on `getBlockId`, so keep it stable.
`meta.hidden` blocks are filtered out before layout. Grid CSS comments encode
hard-won gridstack workarounds — read them before restyling.
