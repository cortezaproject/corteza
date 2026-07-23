---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks/Blocks/RecordListBlock.vue
  - client/web/unify/src/sections/compose/components/PageBlocks/Configurators/RecordListConfigurator.vue
tests: []
---

# Common compose components

## Intention

Cross-family shared pieces of the compose section that belong to no single
resource. Currently only the record-list filter builder.

## Map

- RecordListFilter.vue — visual AND/OR filter-group builder over a module's fields; produces/consumes the record query filter used both at runtime (record list block filtering) and at design time (configurator presets)

## When changing this

The same component serves runtime and builder contexts — verify both consumers
when changing emitted filter shape. The filter output must remain a valid record
query expression the backend accepts; field-kind-specific operators must track
the module field kinds.
