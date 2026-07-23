---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/plugins/primevue-components.ts
touched-by:
  - lib/vue/vitest.config.ts
tests: []
---

# Test setup

## Intention

Vitest global setup for this package's component/composable tests. Every test
mounts inside the same plugin environment lib components assume at runtime,
so individual specs never wire i18n or component registration themselves.

## Map

- setup.ts — registers on `@vue/test-utils` global config: vue-i18n (Composition mode, empty `en` messages), unstyled PrimeVue, and the shared `PrimeVueComponentsPlugin` registry; spies away `console.error`/`console.warn` before each test.

## When changing this

- Tests rely on globally registered PrimeVue/C\* components — removing the
  registry breaks mounted SFC tests package-wide.
- i18n messages are intentionally empty: specs assert on translation KEYS,
  not translated strings; seeding real messages would silently change what
  existing assertions mean.
- Console error/warn are mocked per test — assertions about warnings must go
  through the vi spies, never through captured output.
