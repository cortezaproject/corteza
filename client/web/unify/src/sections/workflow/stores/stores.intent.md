---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/workflow/views/Home.vue
  - client/web/unify/src/sections/workflow/components/WorkflowEditor.vue
tests: []
---

# Workflow section stores

## Intention

Section-local Pinia stores. Only display-name caching lives here; the workflow
list itself is owned by the shared `workflowStore` in `lib/vue` (also consumed
by the shell and sidebar), so this folder must never grow a competing list
store — list mutations go through the shared store's
`updateInList`/`removeFromList`.

## Data touched

- `$ComposeAPI` (injected inside the store) — see `labels.intent.md`.

## Map

- `labels.js` — namespace/module display-name cache for `ref_namespace`/`ref_module` workflow labels (own sidecar).

## When changing this

- New stores need unique app-global Pinia ids — the unified app shares one Pinia across all sections.
- Stores here inject APIs, so they must be first instantiated from a component within the app's provide tree.
