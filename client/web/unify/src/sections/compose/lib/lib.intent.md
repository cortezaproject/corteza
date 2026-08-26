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
  - client/web/unify/src/sections/compose/lib/chart-color-schemes.test.js
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
  template evaluation), the valueless `IS EMPTY`/`IS NOT EMPTY` operators, and
  the page-scoped storage-key helpers, plus small operator/format utilities.
  Three wrappers sit on `evaluatePrefilter`, each carrying the failure contract
  its callers need: `usesRecordVariables` reports whether a template reads the
  page record at all; `evaluatePlacementFilter` returns `undefined` for a filter
  that needs a record there is none of, so the caller reports nothing rather
  than everything; `interpolateDisplayString` returns the template as authored
  for the same case and for anything that will not evaluate, because a title or
  a label must not degrade into an error message.
- `record-filter.test.ts` — vitest spec for the above (escaping, grouping,
  empty-operators and storage keys; keep it green when touching escaping).
- `charts.js` — `chartConstructor(c)`: inspects report metric types to pick
  the concrete chart class (Funnel/Gauge/Radar, else base `compose.Chart`).
- `chart-color-schemes.js` — the instance's own chart palettes, as against the
  tables lib/js compiles in. Owns the one setting they live in
  (`COLOR_SCHEMES_SETTING` = `ui.charts.colorSchemes`) and `readColorSchemes`,
  which answers with a list whatever the setting holds or whether settings
  loaded at all. `colorSchemeOptions` builds the picker's list — the instance's
  own first, then every built-in family named `Family: Label (n colors)` by
  `builtinColorSchemes` — and `upsertColorScheme`/`removeColorScheme` produce
  the whole new setting value without touching the list handed in.
  `isCustomScheme` names the substring lib/js selects a custom scheme by, and
  `newCustomScheme` mints an id carrying it.
- `resource-translations.ts` — applies fetched resource-translation sets
  in-place onto modules/fields/pages/layouts/namespaces (`apply*Translations`).
  Owns the resource-ID string formats (`compose:module/{ns}/{mod}`,
  `compose:module-field/...`, `compose:page/...`, `compose:page-layout/...`,
  `compose:namespace/{ns}`) and the translation key names per entity.
  `moduleFieldKeyLabel`, `pageKeyLabel`, `chartKeyLabel` and
  `namespaceKeyLabel` name one of those keys for display, each returning `''`
  for a key outside its entity so the translator falls back to its own
  formatting. One mapper per resource rather than one that reads the key alone:
  `name` is both a module's key and a namespace's, and `meta.description` is
  both a namespace's and a layout's. Key names must match what
  the server's locale keys emit — a layout's title is `meta.title`, a page's is
  `title`, and the page never carried a `recordToolbar.*` key.

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
- A custom colour scheme is recognised by `custom` appearing in its id, and the
  decision is lib/js's (`getColorschemeColors`), not this module's —
  `isCustomScheme` and `newCustomScheme` exist so the two cannot drift. An id
  minted without it is looked up in the built-in tables and found nowhere.
- The setting is written whole. Both mutators return a new array so a write that
  fails leaves the picker showing what the server still holds.
