---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/RecordListBlock.vue
  - client/web/unify/src/sections/compose/components/PageBlocks/Configurators/RecordListConfigurator.vue
  - client/web/unify/src/sections/compose/components/PageBlocks/Configurators
  - client/web/unify/src/sections/compose/components/Chart/Report/ReportEdit.vue
tests: []
---

# Common compose components

## Intention

Cross-family shared pieces of the compose section that belong to no single
resource: the record-list filter builder and the interpolation-variable
footnote shown next to every configurator input whose value gets interpolated.

## Map

- RecordListFilter.vue — visual AND/OR filter-group builder over a module's fields; produces/consumes the record query filter used both at runtime (record list block filtering) and at design time (configurator presets)
- InterpolationFootnote.vue — props `isRecordPage` (record pages additionally get `${record.values.x}`/`${recordID}`/`${ownerID}`; both kinds get `${userID}`/`${user.name}`) and `dependsOnPlacement` (adds a note that the value is authored away from the page it renders on, e.g. a namespace-level chart placed on either page kind); used by every block Configurator with an interpolated input, keyed off `!!page.moduleID && page.moduleID !== '0'`, and by Chart/Report/ReportEdit.vue

## When changing this

The same component serves runtime and builder contexts — verify both consumers
when changing emitted filter shape. The filter output must remain a valid record
query expression the backend accepts; field-kind-specific operators must track
the module field kinds.
InterpolationFootnote's variable list is a static hint, not derived from code —
keep it in sync with what `evaluatePrefilter` (lib/record-filter.js) actually
supports.
