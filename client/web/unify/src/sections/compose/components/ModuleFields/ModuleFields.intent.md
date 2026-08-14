---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/js/src/compose
touched-by:
  - client/web/unify/src/sections/compose/views/Admin/Modules/Edit.vue
tests: []
---

# Module field configurator

## Intention

The dialog for configuring a single module field in the module admin editor:
general settings, kind-specific options, multi-value delimiter, and validation.
This is the _authoring_ side only — field value editors/viewers live in
`lib/vue/src/components/field/registry.ts` and are not touched here.

## Data touched

Works on an in-memory draft: `index.vue` deep-clones the incoming field and
rebuilds a class instance via `compose.ModuleFieldMaker`
(`@planetcrust/human-js`), so class-level defaults/capabilities apply. No API
calls except the User kind's panel, which reads `$SystemAPI.roleList` on mount
to offer role-restricted user pickers.

## Map

- index.vue — Dialog + tabs; owns the draft, `provide('fieldDraft', ref)`; emits `save` with the whole draft (caller persists), `update:visible`
- CConfiguratorBasic.vue / CConfiguratorMultiDelimiter.vue / CConfiguratorValidation.vue — fixed tabs
- kinds/<Kind>.vue — kind-specific options panel, dynamically imported by field kind. `kinds/Geometry.vue` is the field-side twin of the page-block Geometry configurator and is expected to stay in step with it: same starting-view preview, same lock-bounds behaviour (capture the viewport on lock, update it to the current view on demand, clear `options.bounds` on unlock, preview bounded by the saved area). Only the options each side actually stores differ.

## When changing this

Contract for kind panels: file name must exactly equal the field `kind` (dynamic
import `./kinds/${kind}.vue`); a missing file is legal and simply hides the kind
tab. Panels take no props — they `inject('fieldDraft')` and mutate the draft
directly. Cancel must discard everything, which works only because the draft is a
clone; never hand the original field to the dialog. The multi tab shows only when
`field.cap.multi` and enables only when `isMulti`.
