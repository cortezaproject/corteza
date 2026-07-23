---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue
touched-by:
  - client/web/unify/src/sections/compose/views
  - client/web/unify/src/sections/compose/components
tests:
  - client/web/unify/src/sections/compose/composables/usePageVisibility.test.ts
---

# Compose composables

## Intention

Section-scoped composition functions shared by compose views and components,
covering the two cross-cutting concerns that don't belong to any single
screen: page/block visibility evaluation and resource-translation settings.

## Map

- `usePageVisibility.ts` — layout selection and block visibility for public
  pages and the page builder. `buildExpressionVariables` (user, serialized
  record, screen size/breakpoint, `isView/isEdit/isCreate` on record pages),
  `determineLayout` (first layout whose visibility expression + role list
  pass; a requested layoutID that fails falls back to unconstrained pick),
  `evaluateBlocks` (returns the set of invisible block IDs; a Tabs block with
  no visible child tab hides itself). Expressions are batch-evaluated via
  `$SystemAPI.expressionEvaluate`; API failure treats them all as false
  (fail-closed). Also exports `fetchBlockID` (persisted blockID or
  `meta.tempID` for unsaved builder blocks) and `getBreakpoint`.
- `usePageVisibility.test.ts` — vitest spec for the above.
- `useResourceTranslations.ts` — gates the translator UI: languages come from
  `$Settings` `resourceTranslations.languages` (default `['en']`); the
  translator button shows only when >1 language is configured AND the user
  holds `compose/` `resource-translations.manage` (RBAC store).

## When changing this

- Visibility must stay fail-closed: an expression that cannot be evaluated
  hides the layout/block rather than showing it.
- `fetchBlockID` is the identity contract between builder (tempID) and saved
  pages — RecordView/View.vue and Builder.vue rely on identical behavior.
