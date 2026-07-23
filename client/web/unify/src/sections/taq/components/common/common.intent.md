---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAutomationStore.js
  - lib/js/src/automation/types/icon.ts
touched-by:
  - client/web/unify/src/sections/taq/views/List.vue
  - client/web/unify/src/sections/taq/views/Builder.vue
tests: []
---

# TAQ shared components

## Intention

Small components used by both TAQ views.

## Map

- `TaqConfigModal.vue` — create/edit dialog for TAQ metadata: name (required, `meta.short`), description, labels, run-as user. Create mode calls `ngAutomationCreate` (empty graph), registers the result in `useAutomationStore` and navigates into the builder; edit mode only emits `saved` with the form values — the Builder applies and persists them.
- `TaqIcon.vue` — renders an automation `IconDef` (`name` → PrimeIcon class, `url`/`attachment` → inline img) with an optional fallback icon.

## When changing this

- Run-as '0'/0/empty must normalize to null in the form and back to '0' on
  save — the backend uses '0' as "nobody".
- Create must not send `triggers`/`steps`/`paths`; the builder owns the graph.
