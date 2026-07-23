---
kind: folder
covers: '.'
backfilled: true
owner: fe
depends-on:
  - lib/vue
  - lib/js
  - lib/test-utils
touched-by:
  - client/web/unify/src/sections/compose/index.js
  - client/web/unify/src/sections/compose/ComposeHost.vue
  - client/web/unify/src/sections/compose/components
tests:
  - client/web/unify/src/sections/compose/stores/module.test.ts
  - client/web/unify/src/sections/compose/stores/namespace.test.ts
  - client/web/unify/src/sections/compose/stores/record.test.ts
---

# Compose section stores

## Intention

Compose-specific Pinia stores. The shared module/namespace/record/user/page
stores live in lib/vue (`@planetcrust/human-vue`) and are provided by the
shell — this folder holds only the extras compose alone needs: reminders and
the resource-translator UI state.

## Map

- `reminder.js` — reminders sidebar/toast state + System reminder API (own
  sidecar doc).
- `languages.ts` — server locale list for the translator (own sidecar doc).
- `translator.ts` — translator dialog visibility + per-invocation config (own
  sidecar doc).
- `module.test.ts`, `namespace.test.ts`, `record.test.ts` — vitest specs for
  the **shared** lib/vue stores (`useModuleStore`, `useNamespaceStore`,
  `useRecordStore`), kept here because compose is their main consumer. They
  are covered by this folder doc; the stores under test are not in this
  folder.

## When changing this

- New stores here should be compose-only; anything another section could use
  belongs in lib/vue.
- The lib/vue store specs must be updated in lockstep with lib/vue store
  changes even though they live in this section.
