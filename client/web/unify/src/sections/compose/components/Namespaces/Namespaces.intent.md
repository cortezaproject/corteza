---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/Translator
touched-by:
  - client/web/unify/src/sections/compose/views/Namespace/List.vue
  - client/web/unify/src/sections/compose/views/Namespace/Edit.vue
tests: []
---

# Namespace components

## Intention

Non-route helpers for namespace management: importing a namespace from an export
archive and translating namespace metadata.

## Data touched

NamespaceImporter drives the Compose API namespace import endpoints (upload +
import run). NamespaceTranslator reads/writes namespace translations through the
shared Translator machinery.

## Map

- NamespaceImporter.vue — upload namespace export, trigger server-side import, report progress/result to the list view
- NamespaceTranslator.vue — CTranslatorButton wrapper supplying the namespace resource ident, titles, fetcher and updater

## When changing this

Import consumes the archive produced by namespace export — keep formats in sync.
The translator wrapper follows the same per-resource pattern as the Admin/
translators; change the shared Translator contract there, not here.
