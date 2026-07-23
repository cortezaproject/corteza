---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/Translator
touched-by:
  - client/web/unify/src/sections/compose/views/Admin/Modules/Edit.vue
  - client/web/unify/src/sections/compose/views/Admin/Pages/Edit.vue
  - client/web/unify/src/sections/compose/views/Admin/Pages/Builder.vue
  - client/web/unify/src/sections/compose/views/Admin/Charts/Edit.vue
tests: []
---

# Compose admin panels

## Intention

Sub-panels embedded in the compose admin editor views (module, page, chart).
Each encapsulates one advanced concern so the editor views stay orchestration-only.

## Data touched

Module panels edit `module.config.*` (DAL, discovery, record revisions,
privacy/federation) in place on the module object the parent passes; the parent
saves. DAL panels also read `$SystemAPI` (connection list, schema alterations,
alteration apply/dismiss). Translators call the Compose API translation endpoints
through the shared Translator components.

## Map

- Module/DalSettings.vue, DalFieldStoreEncoding.vue, encoding-strategy.js — per-field store-encoding strategy (omit/plain/alias/json); encoding-strategy.js maps strategy → config payload and knows system-field defaults differ (null)
- Module/DalSchemaAlterations.vue — pending DAL schema alterations for the module
- Module/RecordRevisionsSettings.vue, DiscoverySettings.vue, FederationSettings.vue — feature toggles on module.config
- Module/ModuleIssues.vue, UniqueValues.vue — issue display and unique-value constraints
- Module/ModuleTranslator.vue, Page/PageTranslator.vue, Chart/ChartTranslator.vue — per-resource CTranslatorButton wrappers: each supplies the resource ident, titles, fetcher and updater for its resource (page variant also covers block-level keys)

## When changing this

These panels mutate the parent-owned resource object; keep them side-effect free
until the parent view saves. Translator wrappers are the pattern to copy for a new
translatable compose resource — the shared Translator family defines the contract.
