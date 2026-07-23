---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useAutomationStore.js
touched-by:
  - client/web/unify/src/sections/taq/index.js
tests: []
---

# TAQ sidebar

## Intention

Section drawer navigation: a single "Automations" root (routes to the list)
with every TAQ as a child item, alphabetical by `meta.short`, each routing to
`taq.builder-edit` — so users hop between automations without leaving the
builder.

## Map

- `TaqSidebar.vue` — `CSidebarNav` tree built from `useAutomationStore.list`; lazily calls `fetchList()` on mount when the list is empty.

## When changing this

- The sidebar is a lazy PrimeVue Drawer and may never mount — nothing else may
  rely on it loading the automation list; the builder's catalog is loaded by
  Builder.vue, not here.
- List freshness depends on views mutating through the automation store
  (create/delete/save update the same list this renders).
