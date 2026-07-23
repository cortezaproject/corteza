---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
  - lib/vue
touched-by:
  - lib/vue/src/components/field/viewers/CFieldUserViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldRecordViewer.test.ts
  - lib/vue/src/composables/useUserResolver.test.ts
  - lib/vue/src/composables/useRBAC.test.ts
  - client/web/unify/src/sections/compose/stores/record.test.ts
  - client/web/unify/src/sections/compose/stores/namespace.test.ts
  - client/web/unify/src/sections/compose/stores/module.test.ts
tests: []
---

# Test utils (`@planetcrust/human-test-utils`)

## Intention

Shared Vitest helpers so component and store specs across `lib/vue` and
`client/web/unify` do not each reinvent API mocks, fixture factories, and the
provide/inject scaffolding the apps set up in `App.vue`. Consumed straight from
source (`main` points at `src/index.ts`, no build step); vitest, Vue tooling and
PrimeVue are peer dependencies supplied by the consuming package.

## Map

- `src/api/` — `createMockComposeAPI` / `createMockSystemAPI` / `createMockAutomationAPI`: full `vi.fn()` mock surfaces resolving to empty list results (`{ set: [], filter: { total: 0 } }`) or `{}`; every factory accepts per-method overrides and an index signature admits extra endpoints.
- `src/fixtures/` — `makeUser`, `makeNamespace`, `makeModule`, `makeRecord` (real `compose.*` class instances from `lib/js`), plus `makeRawRecord` for the raw API shape (`values` as `[{ name, value }]`). Auto-incrementing IDs per fixture family (users 9000+, namespaces 10001+, modules 20001+, records 30001+).
- `src/stores/` — `createTestPinia(provides)`: fresh Pinia mounted on a throwaway app hosting the given provides, so `inject()` inside `defineStore()` setups resolves; also calls `setActivePinia`.
- `src/mount/` — `mountWithContext(component, ctx, mountOptions)`: `@vue/test-utils` mount pre-wired with Pinia (active one if set), memory router, empty vue-i18n, Teleport stub, and the app-level provides (`$ComposeAPI`, `$SystemAPI`, `$AutomationAPI`, `$Auth`/`$auth`, `$Settings`, `$userStore`, `$recordStore`, `$pageStore`, `$recordRoutes`, `$namespace`, `$toast`, `$eventBus`).

## When changing this

- The provide-key list in `mountWithContext` mirrors what the real apps
  provide (unify + compose `App.vue`); adding a provide to an app usually means
  adding it here too, or components under test silently get `null`.
- Distinguish explicit `null` from omitted in `TestContext`: passing `null`
  keeps the provide null, omitting falls back to the default (`$Auth`,
  `$Settings` have non-null defaults).
- `makeRecord` builds a matching module from the value keys when no module is
  given — specs asserting field kinds must pass their own module/overrides.
- Call `createTestPinia` in `beforeEach`; a stale active Pinia leaks store
  state (and its provides) into the next test.
