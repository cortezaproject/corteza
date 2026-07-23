---
kind: file
covers: languages.ts
backfilled: true
owner: fe
depends-on: []
touched-by:
  - client/web/unify/src/sections/compose/components/Translator/CTranslatorDialog.vue
tests: []
---

# languages store

## Intention

Lazy, load-once cache of the server's locale list (`{ tag, name,
localizedName }`), so the resource translator can offer language columns
without refetching per dialog open.

## State owned

`set` (Language[]), `loaded`, `loading`.

## API surface consumed

`$SystemAPI.localeList()` (injected).

## Consumers

Translator components (CTranslatorDialog); anything needing the configured
server locales inside compose.

## Invariants

- `load()` is idempotent and re-entrancy-safe: it no-ops while `loading` or
  once `loaded` — callers may fire it unconditionally.
- No invalidation path exists; a locale change server-side requires a page
  reload.
