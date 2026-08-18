---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/lib/record-filter.js
  - lib/vue/src/components/expression/CExpressionHint.vue
touched-by:
  - client/web/unify/src/sections/compose/views/Pages/View.vue
  - client/web/unify/src/sections/compose/views/Pages/RecordView.vue
  - client/web/unify/src/sections/compose/views/Admin/Pages/Builder.vue
tests:
  - client/web/unify/src/sections/compose/components/PageBlocks/Configurators/TabsConfigurator.test.js
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/PageBlock.test.js
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/TabsBlock.test.js
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/MetricBlock.test.js
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
- Blocks/ — runtime components; contract: each accepts `block, namespace, page` plus `record` on record pages, and typically wraps content in Blocks/PageBlock.vue (Card chrome: title, refresh, magnify dialog), passing `record` on so the chrome can interpolate too. Two deliberate exceptions: TabsBlock also takes `blocks` (it renders siblings), and RecordBlock takes no `record` prop — it reads the active record from the injected `recordViewContext` and passes that on as the chrome's record. Calendar, Geometry, Metric, Navigation, Chart, Progress and Automation interpolate author-typed strings via `lib/record-filter.js#evaluatePrefilter` (`${record.values.x}`, `${recordID}`, `${ownerID}`, `${userID}`, `${user.name}`). There is no single guard contract on a non-record page for those: Calendar/Geometry/Metric skip the feed/metric with a `console.warn`, Chart resolves to an empty set silently, Navigation and Automation fall back to the string as authored, and Progress does not guard at all.
- Displayed strings — every block's title and description, and TabsBlock's tab labels — are interpolated too, through `lib/record-filter.js#interpolateDisplayString`, which has one guard contract for all of them: a template reading a record there is none of, or one that fails to evaluate, renders exactly as authored. Chrome must not degrade into an error message. The same helper feeds the magnify dialog header, the Chart and Metric drill-down modal headers, and the iframe's `title` attribute, so one authored string reads the same everywhere it appears.
- Blocks/MetricBlock.vue — keys fetched values by the metric's configured index (`rtr[mi]`), matching how `formatResponse()` looks them up, so a skipped metric doesn't shift later ones. Its drill-down passes the metric filter **raw** (not pre-interpolated) — the child RecordListBlock it opens interpolates that filter itself; evaluating it here would double-evaluate.
- Blocks/RecordBlock.vue — `provide('$recordContext', activeRecord)` so nested field editors (lib/vue) can interpolate a Record field's prefilter. Its `fieldConditions` (server-evaluated via `$SystemAPI.expressionEvaluate` — gval expressions, not `${}` templates) re-evaluate on every record-value change, debounced 300ms, with a sequence guard so a slower earlier response cannot overwrite a newer one or clear values based on a stale record.
- Blocks/GeometryBlock.vue + Configurators/GeometryConfigurator.vue — `options.bounds` means the locked area and nothing else: the configurator captures the viewport when the lock goes on, replaces it on demand through the update button beside the toggle (offered only while the view and the saved area actually differ, which is what zooming in makes possible), and clears it when the lock goes off. The block applies it as leaflet `maxBounds` while `lockBounds` holds and never re-fits the view to it; the starting view is always `center` + `zoomStarting`. Capturing also writes `zoomMin`, because an area already implies that floor and the zoom-range slider beside it would otherwise contradict what the map enforces. The configurator's own preview is bounded by the same box the block will apply, so authoring shows what viewing gets: at the zoom it was locked at the viewport fills the area and panning snaps straight back, and zooming in leaves room to move around inside it.
- Configurators/ — per-kind options editors; NOT in the registry — statically imported and dispatched by name in views/Admin/Pages/Builder.vue; they mutate `block.options`. An author-typed value the block later interpolates or queries with goes in a `CInputExpression` of the matching dialect (`interpolation` for `${}` templates, `ql` for record filters, `expr` for server-evaluated conditions) followed by a `CExpressionHint` bound to the same scope — the hint is what tells an author which variables the input accepts, and `useExpressionScope` is what varies it by page kind. A list whose rows are keyed on a value being edited remounts the row on every keystroke — tab rows are keyed on index and block, never on the title.
- Grid.vue — gridstack layout (48 cols, `block.xywh`); view mode (mobile reflow, single-block fills page) vs builder mode (drag/resize, `layout-updated` emit, `item-overlay` slot for the builder toolbox)
- Shared/ — AutomationButtons(+Editor): workflow/script trigger buttons used by the RecordList block and the Automation _configurator_. The Automation block itself does not use them — it carries its own near-duplicate button rendering.
- Blocks/CBulkRecordEditModal.vue, Blocks/Comment/, Blocks/Metric/ — private helpers of their blocks

## When changing this

A new block kind needs all three: Block component + registry entry + Configurator
(plus its import/case in Builder.vue) — miss one and the kind renders the
"no configuration" fallback or is unconfigurable. A kind that cannot render
without an option states so in its class's `validate()` (lib/js) and marks the
control `required`; the builder does the rest.

Three house rules keep the eighteen reading as one surface: a configurator picks
a module with `CInputModule` and a field with `CInputModuleField` rather than
building option lists (labelling, system fields and the name fallback live
there); a boolean is a `CInputToggleCard` when it carries a description and a
plain `Checkbox` when it is one of a compact group of sibling flags; and
anything downstream of a module choice appears once the module is chosen rather
than sitting there disabled. Block IDs may be temp
(`meta.tempID`) before save; Grid keys/rebuilds on `getBlockId`, so keep it stable.
`meta.hidden` blocks are filtered out before layout. Grid CSS comments encode
hard-won gridstack workarounds — read them before restyling.
