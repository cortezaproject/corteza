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
  - lib/js/src/compose
touched-by:
  - client/web/unify/src/sections/taq/components/builder/form/inputs/registry.ts
  - client/web/unify/src/sections/admin
  - client/web/unify/src/sections/compose
  - lib/vue/src/components/field/editors
tests:
  - lib/vue/src/components/input/CInputModuleField.test.ts
  - lib/vue/src/components/input/CInputPickers.descriptions.test.ts
  - lib/vue/src/components/input/CInputCorredorScript.test.ts
---

# Shared input components

## Intention

App-agnostic form controls shared by every section: entity selectors that own
their own fetching, and generic primitives. All are v-model components
(`modelValue` in, `update:modelValue` out) unless noted. `index.ts` is the
public surface; the CForm\*/CEditorActions scaffolding is imported by path.

## Map

- Selector contract: v-model holds the entity ID (string), or an array of them where the selector offers `multiple`; props `placeholder`, `disabled`; options are self-fetched with search.
- A dropdown option is two lines wherever the resource carries a description: the name, then the description muted beneath it, guarded so an option without one stays a single line. Agent, Chart, KnowledgeBase, LLM, Module, Namespace, Role, UserGroup and Workflow read it from `meta.description`; CorredorScript from a top-level `description`; TAQ (NgAutomation) from a top-level `description`. `CInputPickers.descriptions.test.ts` is the sweep that holds the set together.
- `CInputUser`, `CInputRole`, `CInputUserGroup`, `CInputAgent`, `CInputLLM`, `CInputLabel`, `CInputKnowledgeBase` — SystemAPI-backed selectors (User/Role/UserGroup support exclusion lists; User/Role support a `multiple` mode rendering picked entities as chips; User hands an `option` slot to its dropdown so a caller can add a line under a user without replacing the selector; Label/KnowledgeBase can create inline).
- `CInputModel` — model selector scoped to an LLM provider (SystemAPI).
- `CCodeEditor` — CodeMirror over whatever a view keeps as code: `language` is `html`, `json` or plain text, and `read-only` both refuses the change and stops the element taking the keystroke, so what is shown is what the document holds. Language and read-only are fixed when the editor is built, so changing either rebuilds it. Used by the template editor and by the application editor's custom app page.
- `CInputTAQ`, `CInputWorkflow` — AutomationAPI-backed selectors.
- `CInputCorredorScript` — manual `system` Corredor server scripts the caller may run (SystemAPI automation list); the value is the script name, and the Select is editable so a name that is not listed (a script deployed later, or one the caller cannot list) can still be typed. It declares no `options` prop: step forms pass `options=[]` to every input, which would otherwise replace the fetched list.
- `CInputNamespace`, `CInputModule`, `CInputChart` — store-backed compose selectors (module/chart scoped by `namespaceID` prop).
- `CInputRecord` — record selector; needs `namespaceID` + `moduleID`, label via `labelField`/`recordLabelField`, narrows with `prefilter`/`queryFields` (ComposeAPI + record/module stores).
- `CInputSelect` — thin PrimeVue Select wrapper: `options`, `optionLabel`, `loading`, `showClear`, `hideSearch`.
- `CInputSearch` — search text box used by CResourceList; plain v-model.
- `CInputSwitch` / `CInputToggleCard` — boolean toggles (inline label vs. card with description/warning).
- `CInputDateTime` — date/time picker (`showTime`, `timeOnly`, `onlyDate`, min/max).
- `CInputCron` — cron/interval expression editor.
- `CInputColorPicker`, `CInputLocation` (map dialog geometry), `CInputFile` + `CFileDropZone` (attachment upload; `previewStyle` sets the surface the preview sits on, e.g. the theme's sidebar colour), `CRichTextInput` (tiptap editor + emoji, `submitOnEnter`).
- `CInputDelete` — confirm-guarded delete button (emits confirmation, no v-model).
- `CFieldPicker` — dual-list picker of module fields (`allFields` ⇄ v-model selection).
- `CInputModuleField` — single/`multiple` select of one module's fields; takes `module` or a store-resolved `moduleID`, narrows with `kinds`/`excludeMulti`/`queryableOnly` plus a `filter` predicate for anything those cannot express, appends the record's system fields under `includeSystem`, and offers `extraOptions` (an aggregate like `count`, which is picked here but is not a field) above them. v-model holds `valueKey` (`name`, or `fieldID` for the record block). Labels fall back to the field name, system labels come from `field.system.<name>`, `showName` appends the technical name, and system fields sort last — a caller that builds its own option list drifts from all of it.
- `CFormGroup`, `CFormItemContent`, `CFormItemList`, `CFormList`, `CEditorActions` — form/list layout scaffolding for editor screens. `CEditorActions` carries `data-testid` `editor-actions` and `editor-back`: its Back button is the same control on every editor, so it is addressable as one rather than by each screen's label.

## Cross-cutting

- Selectors resolve their APIs via inject (`$SystemAPI`, `$ComposeAPI`, `$AutomationAPI`) — host app must provide them.
- A fetch replaces a selector's options wholesale, so `CInputUser`/`CInputRole` put an already-picked entity back in afterwards; without it a value past the page or filtered out renders as a bare ID. Pinning runs from the model watcher and from every fetch, so the append de-duplicates after its await.
- `CInputRecord` interpolates its `prefilter` prop through `compose.interpolateTemplate`
  before querying, using `inject('$recordContext', null)` (the record currently being
  edited; only present under compose's RecordBlock → CFieldEditor → CFieldRecordEditor
  chain, so record variables fall back to `'0'` for other hosts, e.g. workflow prompt
  inputs) and `inject('$Auth', null)` (app-wide auth plugin) for the `user`/`userID`
  variables. Falls back to the raw prefilter string on any interpolation failure.
- The TAQ builder input registry maps automation input types onto these components; renaming an export or changing a v-model shape breaks TAQ step config forms.
- `CInputLLM` / `CInputModel` mark their inner Select `formControl: { novalidate: true }`: hosts may wrap them in a named `FormField`/`CFormGroup` for resolver error display without PrimeVue form state hijacking the select value — v-model stays the only value channel.

## When changing this

- `CFormList` and `CFormItemList` carry their own controls — a remove button and
  a drag handle — so an editor cannot make them read-only by disabling its own
  fields. Both take `disabled`: rows still render, reordering and remove go.
- `CFormList` sizes its grid to its content and scrolls sideways under it. That
  is wrong for a row of inputs, whose content width is the browser's default
  `<input>` size rather than anything a caller chose: such a list is wider than
  its container from the first row and hides its right-hand columns. `fitWidth`
  makes the tracks share the width the list is given instead. It is opt-in
  because a flexible column can only share a width if it declares a `minmax()`
  floor: a bare `1fr` collapses to nothing once the fixed columns are paid for,
  and content sizing is the only thing holding such a list up. Opt in a list
  whose columns carry floors and whose grid is measurably wider than its box;
  leave the rest sized to content.

Keep the selector contract (ID-valued v-model, self-fetching) uniform — the TAQ
registry and field editors instantiate these interchangeably by type.
