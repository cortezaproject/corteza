---
kind: file
covers: labels.js
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/workflow/views/Home.vue
  - client/web/unify/src/sections/workflow/components/WorkflowEditor.vue
tests: []
---

# labels store

## Intention

Resolve compose namespace and module IDs (extracted from workflows'
`ref_namespace`/`ref_module` labels) to display names, so list rows and the
editor overlay can render human-readable tags without repeated fetches.

## State owned

- Two reactive id→name maps (`namespaces`, `modules`). `null` marks an
  in-flight resolve (dedup); a failed or nameless lookup stores the id itself
  as the name, so entries always settle.

## API surface consumed

- Injected `$ComposeAPI`: `namespaceRead`, `moduleRead` (module resolution
  needs both `namespaceID` and `moduleID`).

## Consumers

- `views/Home.vue` (batch-resolves after each list page), `components/WorkflowEditor.vue` (resolves the open workflow's labels).

## Invariants

- `resolve*` never rejects; callers may fire-and-forget.
- Getters return `undefined`/`null` until resolved — callers must fall back to the raw id.
- Cache never invalidates within a session; a rename shows stale until reload.
- Pinia id is the app-global `'labels'` — a second `'labels'` store anywhere in the unified app would collide.
