---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/lib/record-filter.js
  - client/web/unify/src/sections/compose/lib/record-sort.js
  - client/web/unify/src/sections/compose/lib/script-events.js
  - lib/vue/src/corredor
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
  - client/web/unify/src/sections/compose/components/PageBlocks/Shared/AutomationButtons.test.js
  - client/web/unify/src/sections/compose/components/PageBlocks/Shared/AutomationButtonsEditor.test.js
  - client/web/unify/e2e/sections/compose/corredor-record-scripts.spec.ts
  - client/web/unify/e2e/sections/compose/custom-block.spec.ts
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
- Blocks/RecordListBlock.vue — column sorting is multi-column (`sort-mode="multiple"`, `removable-sort`): a plain header click sorts by that column alone, ctrl/cmd-click adds or reverses one, and a click past descending drops it. The ctrl/cmd-click hint is a delayed tooltip on the sort icon (drawn through the `sorticon` slot with PrimeIcons). Picked columns replace the block's `presort`; until the user sorts, the presort applies and the headers show it exactly as if picked, so a ctrl/cmd-click appends to it. Clearing every column by clicking leaves the list unsorted, not back on the presort — only the Reset sort button returns to it. That button shows whenever the sort differs from the presort; it sits in the toolbar, just left of the deleted-records toggle, filter and search group. The headers are the only statement of the sort, so a presort key on a column the list does not show appears nowhere. The order badge is the block's own and numbers only shown columns, and only when two or more are sorted — PrimeVue's counts hidden keys too. From the presort, a plain click on a column it sorts descending sorts by that column ascending, where PrimeVue's asc → desc → off cycle would clear it on the first click. `lib/record-sort.js#sortExpression` builds the one sort string both the page fetch and the prev/next ID list send. A multi-value field sorts by its first value (server-side).
- Blocks/GeometryBlock.vue + Configurators/GeometryConfigurator.vue — `options.bounds` means the locked area and nothing else: the configurator captures the viewport when the lock goes on, replaces it on demand through the update button beside the toggle (offered only while the view and the saved area actually differ, which is what zooming in makes possible), and clears it when the lock goes off. The block applies it as leaflet `maxBounds` while `lockBounds` holds and never re-fits the view to it; the starting view is always `center` + `zoomStarting`. Capturing also writes `zoomMin`, because an area already implies that floor and the zoom-range slider beside it would otherwise contradict what the map enforces. The configurator's own preview is bounded by the same box the block will apply, so authoring shows what viewing gets: at the zoom it was locked at the viewport fills the area and panning snaps straight back, and zooming in leaves room to move around inside it.
- Configurators/ — per-kind options editors; NOT in the registry — statically imported and dispatched by name in views/Admin/Pages/Builder.vue; they mutate `block.options`. An author-typed value the block later interpolates or queries with goes in a `CInputExpression` of the matching dialect (`interpolation` for `${}` templates, `ql` for record filters, `expr` for server-evaluated conditions) followed by a `CExpressionHint` bound to the same scope — the hint is what tells an author which variables the input accepts — as does Ctrl-Space in the input itself, which offers the same set as `${…}` snippets — and `useExpressionScope` is what varies it by page kind. Where the value repeats per row (a metric, a tab), the inputs are held in a ref array and one hint under the list inserts into the row last focused; where there is no cursor to insert into (Content's rich text body) the hint is `:insertable="false"` and its chips read as labels. A list whose rows are keyed on a value being edited remounts the row on every keystroke — tab rows are keyed on index and block, never on the title.
- Grid.vue — gridstack layout (48 cols, `block.xywh`); view mode (mobile reflow, single-block fills page) vs builder mode (drag/resize, `layout-updated` emit, `item-overlay` slot for the builder toolbox). A lone visible, non-editable block skips gridstack and fills its parent; from two blocks on, each is exactly `h × 10px` tall (`cellHeight: 10`, `margin: 6`), so adding a second block turns a screen that filled and scrolled inside into fixed heights from `xywh`. `rebuildLayout` watches the joined block IDs only; once gridstack owns geometry, a height that varies with state is never applied.
- Shared/ — AutomationButtons(+Editor): the one path from a page button to a TAQ, a workflow or a script. Both callers go through it: the RecordList block's selection buttons and the Automation block; the Automation _configurator_ uses the Editor half to choose what a button runs. Every value it puts in the run's scope is a `{"@type","@value"}` envelope — `input` decodes into expr.Vars server-side and one bare value rejects the whole request, as HTTP 200 with an error body. It sends what its caller hands it: namespace, page, record and module; selected records and the list's current filter only from a record list. No block is given a module prop, so the Automation block reads one off its record or page. A button naming a `script` takes the script bus instead: the event is shaped by the button's `resourceType` (record, module, namespace, page or compose event, carrying namespace, module, page, the selected records and the filter) and goes through `$ScriptBus.Dispatch(ev, script)`, which runs a client script in the browser or forwards a server script to the API; the record a server script returns is applied to the one on the page, not saved. Scripts constrain on namespace slug and module handle through `lib/script-events.js`. The Editor's Scripts tab offers the `onManual` scripts the compose automation list carries for this app (uiProp `app` compose, unify or unset), filtered by fit: given `page`/`namespace`/`module` props it leaves out scripts whose trigger resource the page cannot hand them (`compose:record` needs a record page and a block that supplies a record — `canSupplyRecord`, false for a record list's selection buttons) and lists, greyed and not addable, those whose namespace/module constraints the page contradicts; each row carries its resource and constraint chips. At runtime a configured script button whose script is not registered (`$UIHooks.FindByScript`) or whose resource cannot be built on this page renders danger-outlined with a tooltip and still explains on click.
- Blocks/CustomBlock.vue + Configurators/CustomConfigurator.vue — custom HTML in the custom app sandbox (`sections/app/app.intent.md` holds the sandbox, the bridge and what a page declares). Two sources, one per block: a custom application's page by `applicationID` — shown only to a viewer who may open that application, otherwise the block says it is not available — or a page of the block's own in `options.source`, which reads the namespace the page is on. The configurator offers the second first and preselects it for a new block; choosing one clears the other. A new block written into the page, or one switched to that source, starts out declaring the record page's module, and only an empty declaration is filled. The server holds an inline page to what an application's page is held to when the page is stored (`server/compose/service/page_custom_block.go`): the sandbox check, writes within modules, every module in the page's namespace, origins reduced to bare origins, and `moduleIDs` set beside the handles; YAML export drops `moduleIDs`. Its `params` are edited as name/value rows and its origins as a one-column table (free text is a table, pickers stay pickers); a parameter value is read as JSON where it parses and kept as text otherwise. The block hands the frame where it is shown (`context`), its own height (`resizable` false: the grid owns geometry), and turns the frame's `changed` into `refetch-records`.
- Blocks/CBulkRecordEditModal.vue, Blocks/Comment/, Blocks/Metric/ — private helpers of their blocks

## When changing this

A new block kind needs all three: Block component + registry entry + Configurator
(plus its import/case in Builder.vue) — miss one and the kind renders the
"no configuration" fallback or is unconfigurable. A kind that cannot render
without an option states so in its class's `validate()` (lib/js) and marks the
control `required`; the builder does the rest.

Three house rules keep the nineteen reading as one surface: a configurator picks
a module with `CInputModule` and a field with `CInputModuleField` rather than
building option lists (labelling, system fields and the name fallback live
there); a boolean is a `CInputToggleCard` when it carries a description and a
plain `Checkbox` when it is one of a compact group of sibling flags; and
anything downstream of a module choice appears once the module is chosen rather
than sitting there disabled. Block IDs may be temp
(`meta.tempID`) before save; Grid keys/rebuilds on `getBlockId`, so keep it stable.
`meta.hidden` blocks are filtered out before layout. Grid CSS comments encode
hard-won gridstack workarounds — read them before restyling.
