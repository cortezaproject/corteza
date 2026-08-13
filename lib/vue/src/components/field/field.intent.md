---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/components/field/registry.ts
  - lib/vue/src/components/input
  - lib/vue/src/stores/useRecordStore.js
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/useUserStore.js
  - lib/js/src/compose
touched-by:
  - client/web/unify/src/sections/compose/components/PageBlocks
  - client/web/unify/src/sections/compose/components/Common/RecordListFilter.vue
  - client/web/unify/src/sections/taq/components/builder
tests:
  - lib/vue/src/components/field/CFieldEditor.test.ts
  - lib/vue/src/components/field/CFieldEditor.formfield.test.ts
  - lib/vue/src/components/field/url.test.ts
  - lib/vue/src/components/field/viewers/CFieldBoolViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldDateTimeViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldNumberViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldRecordViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldSelectViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldStringViewer.test.ts
  - lib/vue/src/components/field/viewers/CFieldUserViewer.test.ts
---

# Field system

## Intention

One shared edit/view surface for compose module field values. Any screen that
shows or edits a record value goes through the two dispatchers here instead of
switching on `field.kind` itself, so a field kind behaves identically in record
pages, page blocks, filters and the TAQ builder.

## Map

- `registry.ts` — kind → { editor, viewer } lookup with String fallback (own sidecar).
- `CFieldEditor.vue` — editor dispatcher; owns single/multi-value orchestration.
- `CFieldViewer.vue` — viewer dispatcher; pure pass-through to the kind's viewer.
- `editors/` — one `CField<Kind>Editor.vue` per kind; heavier ones delegate to `../input` (CInputRecord, CInputUser, CInputFile, CInputLocation).
- `viewers/` — one `CField<Kind>Viewer.vue` per kind; read-only presentation.
- `url.ts` / `index.ts` — Url-viewer helpers; public exports (dispatchers, registry, url helpers).

## Contracts

- Sub-editors are rendered with `field`, `namespace`, `model-value`, `disabled`
  and must emit `update:modelValue`. Nothing else is guaranteed to be passed.
- Viewers are rendered with `field`, `record`, `namespace`, `valueOnly`,
  `extraOptions`, `disableClick`. Viewers read their own value from the record:
  `record[field.name]` for system fields, else `record.values[field.name]`.
- Multi-value logic lives ONLY in CFieldEditor: one sub-editor per entry with
  add/remove and stable entry ids (`allowEmpty` decides whether an emptied list
  keeps one blank entry). Exception ("multi-absorbing"): File always, and
  Select/User/Record with `options.selectType === 'multiple'`, receive the whole
  array as modelValue and manage it themselves.
- CFieldEditor severs the `@primevue/forms` injections (`$pcFormField`,
  `$pcForm`) for its subtree. A field editor placed inside a `<FormField>`
  would otherwise have all its multi-value inputs write one shared value,
  overwriting each other.
- Record value shape: viewers expect the normalized compose.Record shape
  (`values` keyed by field name). Raw API records (`values` = array of
  {name, value}) exist only in the record store's label cache when the module
  is not loaded; CFieldRecordViewer guards this by rendering a plain record ID
  unless the target module (and its label field) is available.
- CFieldRecordViewer resolves labels via `recordStore.resolveRecordLabels` and
  navigates on click (injected `$recordRoutes` wins; else page-store record page,
  honoring `extraOptions.recordSelectorDisplayOption`: sameTab/newTab/modal).
- CFieldRecordEditor resolves its Record field's `prefilter` through
  `compose.interpolateTemplate` before splicing it into the CQL filter, using
  `inject('$recordContext', null)` (the record being edited, provided by
  compose's RecordBlock) and `inject('$Auth', null)` (app-wide auth plugin) as
  the source of the `record`/`recordID`/`ownerID`/`user`/`userID` template
  variables. Falls back to the raw prefilter string on any interpolation
  failure, so existing non-templated configurations can't start throwing.

## When changing this

- Keep dispatcher prop sets stable; every consumer and sub-component relies on them.
- Unknown kinds must keep degrading to the String pair, never throw.
- Do not add per-kind branching to consumers — extend the registry instead.
