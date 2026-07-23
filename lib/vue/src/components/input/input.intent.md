---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useNamespaceStore.js
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/useRecordStore.js
  - lib/vue/src/stores/useChartStore.js
touched-by:
  - client/web/unify/src/sections/taq/components/builder/form/inputs/registry.ts
  - client/web/unify/src/sections/admin
  - client/web/unify/src/sections/compose
  - lib/vue/src/components/field/editors
tests: []
---

# Shared input components

## Intention

App-agnostic form controls shared by every section: entity selectors that own
their own fetching, and generic primitives. All are v-model components
(`modelValue` in, `update:modelValue` out) unless noted. `index.ts` is the
public surface; the CForm\*/CEditorActions scaffolding is imported by path.

## Map

- Selector contract: v-model holds the entity ID (string); props `placeholder`, `disabled`; options are self-fetched with search.
- `CInputUser`, `CInputRole`, `CInputUserGroup`, `CInputAgent`, `CInputLLM`, `CInputLabel`, `CInputKnowledgeBase` — SystemAPI-backed selectors (User/Role support exclusion lists; Label/KnowledgeBase can create inline).
- `CInputModel` — model selector scoped to an LLM provider (SystemAPI).
- `CInputTAQ`, `CInputWorkflow` — AutomationAPI-backed selectors.
- `CInputNamespace`, `CInputModule`, `CInputChart` — store-backed compose selectors (module/chart scoped by `namespaceID` prop).
- `CInputRecord` — record selector; needs `namespaceID` + `moduleID`, label via `labelField`/`recordLabelField`, narrows with `prefilter`/`queryFields` (ComposeAPI + record/module stores).
- `CInputSelect` — thin PrimeVue Select wrapper: `options`, `optionLabel`, `loading`, `showClear`, `hideSearch`.
- `CInputSearch` — search text box used by CResourceList; plain v-model.
- `CInputSwitch` / `CInputToggleCard` — boolean toggles (inline label vs. card with description/warning).
- `CInputDateTime` — date/time picker (`showTime`, `timeOnly`, `onlyDate`, min/max).
- `CInputCron` — cron/interval expression editor.
- `CInputColorPicker`, `CInputLocation` (map dialog geometry), `CInputFile` + `CFileDropZone` (attachment upload), `CRichTextInput` (tiptap editor + emoji, `submitOnEnter`).
- `CInputDelete` — confirm-guarded delete button (emits confirmation, no v-model).
- `CFieldPicker` — dual-list picker of module fields (`allFields` ⇄ v-model selection).
- `CFormGroup`, `CFormItemContent`, `CFormItemList`, `CFormList`, `CEditorActions` — form/list layout scaffolding for editor screens.

## Cross-cutting

- Selectors resolve their APIs via inject (`$SystemAPI`, `$ComposeAPI`, `$AutomationAPI`) — host app must provide them.
- The TAQ builder input registry maps automation input types onto these components; renaming an export or changing a v-model shape breaks TAQ step config forms.

## When changing this

Keep the selector contract (ID-valued v-model, self-fetching) uniform — the TAQ
registry and field editors instantiate these interchangeably by type.
