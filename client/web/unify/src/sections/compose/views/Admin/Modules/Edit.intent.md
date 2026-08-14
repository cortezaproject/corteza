---
kind: file
covers: Edit.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
  - client/web/unify/src/sections/compose/components/ModuleFields/Configurator/index.vue
  - client/web/unify/src/sections/compose/components/Admin/Module
  - client/web/unify/src/sections/compose/lib/resource-translations.ts
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Module Edit view

## Intention

The module editor: define fields and their kinds/config, storage and integrity
settings, and wire the module to its record pages. Central admin surface for a
module's whole lifecycle.

## UX capabilities

- Tabs: Fields (name/handle + draggable field list), DAL, unique values, record revisions, and an Issues tab shown only when the server reports module issues (opening it triggers the schema-alterations dialog).
- Per-field row: name/label/kind, required/multi flags, kind-specific configurator modal (shares the draft via `moduleDraft` provide), field permissions, and translation actions (label, select options, bool labels) behind the translator gate; system fields shown read-only below.
- Validation: module name required (free text — no pattern), handle matches the handle pattern, field names valid identifiers and unique, field labels required — errors surface on the Fields tab and block save.
- Related-page actions (`canManageNamespace`): open-or-create the module's record page and record-list page (created pages get a seeded `primary` layout; record page parents under the list page), edit in builder.
- Discovery/federation settings modals gated by `$Settings` feature flags; export JSON; permissions menu covering module, all fields, and all records wildcards.
- Save/clone/delete with unsaved-changes guard; create redirects into edit mode.

## Routes

`admin.modules.create` at `admin/modules/create`, `admin.modules.edit` at `admin/modules/:moduleID/edit`. Topbar links to `admin.modules.record.list` and the module translator. Load failure bounces to `admin.modules`.

## When changing this

- Auto-created pages must use the `compose.PageLayout` type and seed `layout.blocks` from the page blocks — a plain object or empty layout persists disabled toolbar buttons / hides every block (see in-code comments).
- Record-list detection scans page blocks for a `RecordList` targeting this module; record-page detection is `page.moduleID` — creating the record page requires the list page decision first (`selfID`).
- Modules with issues auto-open the schema alterations dialog on load.
- The module name is a label, the handle is the identifier — never apply
  `isValidFieldName` (the field-name identifier regex) to it; the server validates
  `Handle` and field names, never `Name`. The name/handle inputs sit in a resolver
  `Form`, which paints them invalid and renders the message via `CFormGroup`: a
  hand-bound `:invalid` there is a second rule that `handleSubmit` never checks,
  so it can show red on something that saves fine.
