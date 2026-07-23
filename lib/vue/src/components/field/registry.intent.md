---
kind: file
covers: registry.ts
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/components/field/editors
  - lib/vue/src/components/field/viewers
touched-by:
  - lib/vue/src/components/field/CFieldEditor.vue
  - lib/vue/src/components/field/CFieldViewer.vue
tests: []
---

# registry.ts

## Intention

Single source of truth for which component edits and which component views each
compose field kind. Consumers never import kind components directly — they call
`resolveFieldEditor(kind)` / `resolveFieldViewer(kind)` (or read `FIELD_REGISTRY`).

## Contract

- `FIELD_REGISTRY: Record<kind, { editor, viewer }>` — kinds: String, Number,
  Bool, DateTime, Select, Email, Url, User, Record, File, Geometry.
- Both resolvers fall back to the String pair for unknown/missing kinds; they
  never return undefined and never throw.
- All entries are `defineAsyncComponent` — kind components stay lazy-loaded;
  keep new entries async too.

## Adding a field kind

Create `editors/CField<Kind>Editor.vue` (props: field, modelValue, disabled,
namespace; emits update:modelValue) and `viewers/CField<Kind>Viewer.vue`
(props: field, record, namespace, valueOnly, extraOptions, disableClick),
then register both here. No dispatcher or consumer changes needed.

## When changing this

Removing or renaming a kind silently downgrades existing module fields to the
String fallback — migrate stored field kinds before dropping an entry.
