---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js
touched-by:
  - client/web/unify/src/sections/compose/components
  - client/web/unify/src/sections/compose/views
tests:
  - client/web/unify/src/sections/compose/lib/record-filter.test.ts
---

# Compose lib helpers

## Intention

Pure, component-free helper modules for the compose section: record-list
filtering/QL generation, chart subtype construction, and applying resource
translations to compose entities. No Vue reactivity, no stores — safe to unit
test in isolation.

## Map

- `record-filter.js` — record-list filter model ↔ Human QL. Escaping
  (`escapeQlString` — QL knows only `\'` and `\\`; LIKE patterns need double
  backslash escaping plus `%`/`_`), filter→SQL (`getRecordListFilterSql`,
  `getFieldFilter`), query→filter groups (`queryToFilter` honoring
  `nonQueryableFieldNames`/`nonQueryableFieldKinds`), prefilter macro
  evaluation (`evaluatePrefilter` — signature/behaviour unchanged, delegates
  to lib/js's `compose.interpolateTemplate` so lib/vue can reuse the same
  template evaluation), plus small operator/format utilities.
- `record-filter.test.ts` — vitest spec for the above (escaping and grouping
  edge cases live here; keep it green when touching escaping).
- `charts.js` — `chartConstructor(c)`: inspects report metric types to pick
  the concrete chart class (Funnel/Gauge/Radar, else base `compose.Chart`).
- `resource-translations.ts` — applies fetched resource-translation sets
  in-place onto modules/fields/pages/layouts/namespaces (`apply*Translations`).
  Owns the resource-ID string formats (`compose:module/{ns}/{mod}`,
  `compose:module-field/...`, `compose:page/...`, `compose:page-layout/...`,
  `compose:namespace/{ns}`) and the translation key names per entity.

## When changing this

- Escaping in `record-filter.js` mirrors the server-side QL lexer + SQL LIKE
  semantics; changing it silently corrupts saved prefilters — extend the test
  first.
- The template-evaluation logic behind `evaluatePrefilter` (`${record.values.x}`,
  `${recordID}`, `${ownerID}`, `${userID}`, `${user.name}`) lives in
  `lib/js/src/compose/helpers/interpolate.ts`, not here — change semantics
  there, not by re-wrapping `eval`/`Function` locally.
- Translation resource-ID/key formats must match what the server's resource
  translation endpoints emit; they are shared with the Translator components.
