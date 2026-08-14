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
tests:
  - client/web/unify/src/sections/compose/components/Admin/Module/DalSettings.draft.test.js
---

# Compose admin panels

## Intention

Sub-panels embedded in the compose admin editor views (module, page, chart).
Each encapsulates one advanced concern so the editor views stay orchestration-only.

## Data touched

Most module panels edit `module.config.*` (DAL, discovery, record revisions,
privacy) in place on the module object the parent passes; the parent saves. DAL
panels also read `$SystemAPI` (connection list, schema alterations, alteration
apply/dismiss). Translators call the Compose API translation endpoints through
the shared Translator components.

## Map

- Module/DalSettings.vue, DalFieldStoreEncoding.vue, encoding-strategy.js — connection choice plus per-field store-encoding strategy (omit/plain/alias/json); encoding-strategy.js maps strategy → config payload and knows system-field defaults differ (null). An unset `connectionID` ('0') means "whichever connection is primary" and stays unset: the primary is resolved for display only
- Module/DalSchemaAlterations.vue — pending DAL schema alterations for the module
- Module/RecordRevisionsSettings.vue, DiscoverySettings.vue — feature toggles on module.config
- Module/FederationSettings.vue — standalone field-mapping modal, shown only when the global `federation.enabled` setting is on; persists through `$FederationAPI` itself and never touches `module.config`
- Module/ModuleIssues.vue, UniqueValues.vue — issue display and unique-value constraints
- Module/ModuleTranslator.vue, Page/PageTranslator.vue, Chart/ChartTranslator.vue — per-resource CTranslatorButton wrappers: each supplies the resource ident, titles, fetcher and updater for its resource (the page variant also covers block-level keys and, via a `layouts` prop + `update:layouts` emit, page-layout keys)

## When changing this

These panels mutate the parent-owned resource object; keep them side-effect free
until the parent view saves. Loading counts: a panel that writes a resolved
default into the draft on mount makes an untouched editor report unsaved
changes, and turns the next save into a change the user never asked for — show
the default, write only what the user picks. The two exceptions own their own writes and say so
by having their own action buttons: FederationSettings (its Save) and
DalSchemaAlterations (apply/dismiss). Translator wrappers are the pattern to copy for a new
translatable compose resource — the shared Translator family defines the contract.
