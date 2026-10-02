# Compose gotchas

Facts about compose namespaces, modules, records, pages, blocks and charts that the code does not make obvious.

## Block `options.fields` has two shapes

A page block's `options.fields` is either `["title", "qty"]` (what the unify field pickers write: `CFieldPicker` emits names and `RecordConfigurator`/`RecordListConfigurator` `onFieldPickerUpdate` pass them straight on) or `[{ "name": "title" }]` (pages from `dev/agent/pagebuild.py`, envoy fixtures, the vue2-era compose). Touching a block's picker rewrites it to the name shape; the server never normalises.

**Why:** a reader assuming objects passes review and fails only on blocks whose picker was used (issue #41: the page builder's required-field guard refused every save).

**How to apply:** read `options.fields`, `editFields`, `inlineEditFields`, `searchableFields` with `f.name ?? f`; don't change the writers.

## Conditional visibility is one all-or-nothing batch

Layout `config.visibility`, block `meta.visibility` (both `usePageVisibility.ts`) and RecordBlock `options.fieldConditions` all POST `/system/expressions/evaluate`. `Evaluate` in `server/system/service/expression.go` returns `nil, err` if any expression fails, so one typo fails every expression in the request. Layouts and blocks then fail closed: every expression reads `false`, and a page with layouts but no match shows a warning and leaves. Field conditions keep their previous answer or show the field.

Evaluation: an unknown one-level key is `false`, a two-level miss (`record.values.x.y`) errors; `record.values.x == ""` is false for a never-set field but true for a cleared one, so `isEmpty()` is the only correct emptiness test. Results are not coerced: `"false"` is truthy. `screen.*` is captured once per page load. Clear-on-hide (per-condition `clearOnHide` OR block `clearConditionalFieldsOnHide`) runs only in edit/create mode and writes to the shared cached record, so it empties the field in every block showing it.

## Organize places the moved record at the requested position

`record.Organize` (`server/compose/service/record.go`) writes the moved record's position field to `position`, then renumbers records matching `filter AND posField >= position AND recordID != <id>` from `position + 1`. The client's `positionFor` in `RecordOrganizerBlock.vue` sends one past the card it was dropped behind (Corteza's `calcNewPosition`).

**How to apply:** a stored position is a sort key, not a list index; read siblings back from the server after a move rather than trusting an optimistic list.

## Page layouts: server-made, never adopting, fully replaced

`createPrimaryLayout` (`server/compose/service/page.go`) gives every new page a layout seeded from its blocks, in the page-create transaction. A layout never adopts a page's existing blocks: a block added later to `page.blocks` renders nowhere until its `{blockID, xywh}` is appended to the layout (`pagebuild.py` `sync_layout` does this; a plain REST page POST does not). A one-block page hides the gap because Grid takes its single-block path. List layouts with `GET /compose/namespace/{ns}/page-layout?pageID={id}` (no trailing slash).

Layout and page POSTs are full replaces: a partial body blanks handle, `meta.title`, `config` and `blocks`. GET, mutate, POST the whole object. One record page per module: a second fails `uniqueCheck` with `recordPageNotUnique`. `pagebuild.py` does not send a block's `description`; set it with a follow-up page POST.

## `pageReorder` errors for nested parents

The generated `ComposePageFilter` (`server/store/adapters/rdbms/filters.gen.go`) appends `parent_id = ParentID` when `ParentID > 0`, but `compose_page` has no `parent_id` column (it is `self_id`). The custom `f.ComposePage` in `filter.go` calls the generated filter first and only then adds `self_id`, so `store.ReorderComposePages` errors for any nested parent. Root reorder works (it sets `Root`).

**How to apply:** reorder nested pages via page update with a changed `weight` (and `selfID` to reparent); the update path honours both.

## Record save errors arrive on HTTP 200

Compose record endpoints put failures in the body of a 200: validation as `{error:{message:"2 issue(s) found", details:[{kind,message,meta:{field,id}}]}}`, everything else as `{error:{message, meta:{resource,type}, stack}}`. The summary is only a count; `details[].meta.field` names the field. `{{value}}` in `duplicateValue` is never substituted server-side; the fill is in `detail.meta.value` and `detailMessage()` in `client/web/unify/src/sections/compose/lib/record-errors.js` interpolates it.

Non-strict duplicate warnings ride a successful save in `valueErrors` (`compose.Record` copies them). A value-expression field has no editor, so its complaint belongs in the alert, not inline. Attachment uploads bypass `stdResolve`; read the body with `apiError()` (same file) — some upload failures are plain text. The server sniffs content (`extractMimetype`), so a text file named `.png` passes the client and is refused.

## A record value has two empty states

Never set means no value row (`Field IS NULL`); set then cleared means a row holding `''`, which counts as NOT NULL. `(Field != '')` alone is the correct "has a value" test. Only text-column kinds (String, Email, Url, Select) tolerate `= ''`; Number, Bool, DateTime, User, Record and File are typed columns where `= ''` is a hard postgres error. `getFieldFilter` in `sections/compose/lib/record-filter.js` keeps that list as an allowlist. Geometry cannot be queried at all.

The organizer disagrees with itself: its ungrouped column lists with `IS NULL`, while the server's `organize` writes `''` for an empty group, so a card dropped into "Ungrouped" is in no column on reload.

**How to apply:** test empty kind-aware; fix the ungrouped case on the server by clearing the value.

## The record report returns null for the no-value group

`GET .../module/{mod}/record/report` puts records with no dimension value in their own row with `dimension_0: null` and null aggregates (`recordReportCorrectTypes` in `server/compose/service/record.go`), distinct from records holding `0`.

**Why:** chart `skipMissing` and the dimension's `default` label (`pickLabel` in `lib/js/src/compose/types/chart/base.ts`) depend on that distinction.

**How to apply:** handle `null` in dimension and metric slots; never pass metrics through the dimension default; guard numeric folds with `Number.isFinite(parseFloat(v))` (`isNaN(null)` is false). Check `ChartBlock`, `MetricBlock`, `page-block/progress.ts`.

## Value expressions run through the full record pipeline

`RecordPreparer` (`server/compose/service/record.go`) runs sanitizer → `values.Expression` → sanitizer again → validator → `GetClean` → formatter, so computed values are sanitized (refs resolved), validated, required-checked and unique-checked like any other. The scope (`ValueExprScope` in `values/expr.go`) holds each field by name plus `new` and `old`; on create `old` is an empty record of the module. `options.isUnique` compares the typed value across the whole module, including records the saver cannot read, and skips multi-value fields (those use `isUniqueMultiValue`). The module editor does not expose `isUnique`; only API, envoy and agentic callers set it.

## An unknown chart colour scheme is silent

`getColorschemeColors` (`lib/js/src/shared/types/chart/helper.ts`) falls back to the default palette when a name resolves to nothing, so a misspelt scheme renders without error in the wrong colours. The MCP `compose_chart_create`/`_update` refuse an unknown name (`server/compose/agentic/chart_color_scheme.go`); REST, `pagebuild.py` and envoy do not. Valid keys are those in `lib/js/src/shared/types/chart/colorschemes/tableau.ts`; near-misses bite (`ClassicOrangeBlue13` exists, `ClassicOrangeBlue7` does not, `OrangeBlue7` does).

**How to apply:** grep the key before using it outside the MCP.

## Aggregate at read time, not in a stored rollup

A rollup stored on a parent record has nothing to keep it true: the TAQ construct library has no aggregate function and a value expression sees only its own record. A Chart with `{"aggregate":"SUM"}` over the child module is always live, and a Record-ref dimension renders the parent's label, so no denormalised name field is needed. A Metric block does the same for one number. Value expressions suit within-record derivations (`quantity > owned ? quantity - owned : 0`).

## Query literals escape with a backslash

`TokenConsumerString.Consume` (`server/pkg/ql/token_consumers.go`) treats `\` as the escape: `name = 'Urza\'s Saga'`. The SQL habit `'Urza''s Saga'` closes the literal and matches nothing, with no error. Double quotes are not string delimiters (`illegal token`). There is no `IN`: batch by ID with the `recordID` REST parameter (`RecordFilter.RecordID`) or `recordIDs` on `compose_record_lookup`. No subqueries or joins; a dotted path (`card.name = 'Bolt'`) is resolved by the agentic layer for `=` and `LIKE` only (`server/compose/agentic/record_filter_refs.go`).

## Corteza is the regression oracle

Human is a port of Corteza; the vue2 original is the Corteza checkout (path per machine). Block sources there are `client/web/compose/src/components/PageBlocks/<Kind>Base.vue`, shared types `lib/js/src/compose/types/`. Stored block options are shared, so an option Corteza wrote must keep working in Human.

The Bootstrap tell: Corteza was BootstrapVue, Human is Tailwind + PrimeVue. Where an option keeps Bootstrap vocabulary (`tabs`/`pills`/`justified`, a `variant`, `b-*`), grep the class across `client/web/unify/src` and `lib/vue/src`; written but never defined means inert.

**How to apply:** when a compose block behaves oddly, diff Human's `Blocks/<Kind>Block.vue` against Corteza's `<Kind>Base.vue` before theorising.

## Canvas fonts in echarts

A canvas rejects a font string whose family is `inherit` or empty whole, size included, and draws `10px sans-serif`; `ChartRenderer` passes `getComputedStyle(document.body).fontFamily` (PrimeVue defines no `--p-font-family`). echarts 6 hard-codes `fontSize: 12` on axis labels, legends and gauge text, so size each explicitly (`chartFontSize` in `lib/js/src/compose/types/chart/util.ts`). To measure what was drawn, hook `CanvasRenderingContext2D.prototype.fillText` in a playwright init script; `getOption()` shows only what was asked. vue-echarts 8's exposed `setOption` is a no-op without `manual-update`; use `getInstanceByDom(el)` from `echarts/core`. A radar's default radius is `'50%'`.

## vue-echarts puts the colour array in `replaceMerge`

vue-echarts 8 treats every top-level array that shrank as a component list, so a palette switch to fewer colours put `color` into `replaceMerge`; echarts throws `"color" is not valid component main type`, then `setOption should not be called during main process`. `CChart.vue` (`lib/vue/src/components/chart/`) pins `replaceMerge` to `series`, `xAxis`, `yAxis`.

**How to apply:** keep that explicit list when touching `CChart`; add a new component array to it rather than dropping it.

## Multi-value flips, retyping and select display

Values live as JSON arrays even for single fields (`{"color":["Blue"]}`), so setting `isMulti: true` reads old rows without migration. On a multi-value field, `field = 'X'` and sorting use the first value only, and report grouping expands values, so a doughnut double-counts rows with several. A stored field's name and kind are locked by field ID while the module has records (`renamedOrRetypedField` in `server/compose/service/module.go`); to retype, back up values, save the module without the field, save again with a new field of the same name (ID 0), rewrite every record.

Select `selectType: multiple` is one MultiSelect (a set); `each` is one dropdown per value (repeats, ordered). `displayType: badge` with option `style` colours `#RRGGBBAA`; `multiDelimiter` `'\n'` stacks badges. Don't combine `isUniqueMultiValue` with `multiple`: selected options vanish from the list and cannot be unticked.
