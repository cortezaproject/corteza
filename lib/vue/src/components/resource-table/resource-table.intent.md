---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/admin
  - client/web/unify/src/sections/compose/components/PageBlocks/Configurators/TabsConfigurator.vue
tests: []
---

# Resource table

## Intention

Lightweight embedded table for small, already-loaded collections inside editor
screens (triggers, settings rows, code snippets, tab configs) — the non-lazy,
non-paginated sibling of CResourceList. Reach for it when the data is a plain
array in memory; reach for resource-list when the API pages.

## Contracts

- Renders all `items` directly (no pagination, no search, no lazy loading);
  `primaryKey` defaults to `_dataKey`.
- `fields` drives columns ({ key, header, sortable, style, frozen, pt, `hint`
  for an info-tooltip header, ... }); cell content via `body-<key>` slots with
  DataTable slotProps. Extra hand-written `<Column>`s pass through the default
  slot. `inheritAttrs` is off — unknown attrs land on the inner DataTable.
- `actionItems(row, index) => MenuItem[]` appends the same hover-revealed
  ellipsis column + centralized TieredMenu as CResourceList (items support
  `route`); a `fields` entry keyed `actions` is then dropped.
- `reorderableRows` adds a drag-handle column and emits `row-reorder` — the
  parent owns persisting the new order. Also emits `sort`, `row-click`.
- Layout knobs: `fillHeight` (+`scrollHeight`), `resizable`, `rowClass`, `pt`.
- Exposes `hideActionsMenu()` and the raw `dataTableRef`.

## When changing this

Keep the actionItems/TieredMenu behavior in step with CResourceList — the two
tables should feel identical to users. Reorder emits index-based events; verify
against a screen that persists order (e.g. workflow triggers) after changes.
