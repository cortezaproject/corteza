---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/project/stores/projects.js
  - client/web/unify/src/sections/project/config/fieldTypes.js
  - client/web/unify/src/sections/project/utils/fields.js
  - client/web/unify/src/sections/project/components/permissions
touched-by:
  - client/web/unify/src/sections/project/views/Wizard.vue
  - client/web/unify/src/sections/project/components/wizard/steps/DataModelStep.vue
  - client/web/unify/src/sections/project/components/permissions/ProjectPermissionMatrix.vue
tests: []
---

# Data model dialogs (modules & fields)

## Intention

The module resource-kind's implementation of the project resource dialog
standard: a minimal CreateDialog (name only — everything else is edited
afterwards) and an editable DetailDialog, both rendered once by the Wizard and
opened via its `createResource('module')` / `inspectResource('module', id)`.
Field editing is a third, stacked dialog so module meta stays staged while
fields persist immediately.

## Data touched

- `useProjectsStore`: `resourcesFor` (modules live in the generic resource
  list, `kind: 'module'`), `addResource`, `updateResource`, `addField`,
  `updateField`, `removeField`.
- No direct API calls — all persistence goes through the store.

## Map

- `ModuleCreateDialog.vue` — name-only create; validates slug-handle clashes
  (`fieldName` lowercased) against sibling modules; `createdId` guard makes a
  re-save after partial failure reuse the created module instead of duplicating.
- `ModuleDetailDialog.vue` — staged name/description (committed on Save); live
  field list with immediate, confirm-guarded removal; embeds
  `ResourcePermissionsSection` (`kind="module"`); footer deep-links to the full
  admin module editor (new tab) when the project has a namespace. Delegates
  field add/edit to the Wizard-provided `createField` / `editField` injects.
- `FieldDialog.vue` — the shared field editor the Wizard renders (stacked over
  the module dialog): two-step flow (visual type picker → config form), staged
  draft, persists on Save via `addField`/`updateField`; supports `readonly`.
  Select option values and record label-fields use `fieldName`-derived machine
  handles — clash validation mirrors what the store will persist.
- `FieldKindTag.vue` — neutral badge (icon chip + localized type name) for a
  field's kind; also used by DataModelStep and the permission matrix.

## When changing this

- Uphold the dialog standard: dialogs are opened only through the Wizard's
  provides; never mount them ad hoc in steps.
- Handle derivation (`utils/fields.js fieldName`) is the source of clash rules
  here AND in the store's persistence — keep both in sync.
- Module meta is staged (Save commits); fields persist immediately. Do not
  merge the two models without a ruling.
