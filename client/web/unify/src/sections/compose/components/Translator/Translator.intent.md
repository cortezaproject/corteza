---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/stores/translator.ts
  - client/web/unify/src/sections/compose/composables/useResourceTranslations.ts
touched-by:
  - client/web/unify/src/sections/compose/ComposeHost.vue
  - client/web/unify/src/sections/compose/components/Admin
  - client/web/unify/src/sections/compose/components/Namespaces/NamespaceTranslator.vue
tests: []
---

# Resource translator

## Intention

Generic multi-language translation editing for any compose resource (namespace,
module, page, chart, fields, blocks). One dialog instance serves the whole
section; per-resource wrappers only describe _what_ to translate.

## Data touched

The trio itself is API-agnostic: callers supply `fetcher` (load translation
rows) and `updater` (persist changes) closures, so all endpoint knowledge stays
in the per-resource wrappers. Visibility is gated by `useResourceTranslations`
(feature enabled + permissions). Active config lives in `stores/translator`.

## Map

- CTranslatorButton.vue — the entry point wrappers render; on click passes `{resource, titles, fetcher, updater, highlightKey?, keyPrettifier?}` to `translatorStore.open()`
- CTranslatorDialog.vue — singleton dialog mounted in ComposeHost; opens whenever the store holds a config
- CTranslatorForm.vue — language-column manager + key × language editing table; saves via the config's `updater`

## When changing this

The store-config shape (`TranslatorConfig`) is the contract every wrapper
(Admin/\*Translator, NamespaceTranslator) depends on — extend it
backward-compatibly. Never mount a second dialog. The default language column is
not removable; `highlightKey` scrolls/marks one row for "translate this field"
entry points.
