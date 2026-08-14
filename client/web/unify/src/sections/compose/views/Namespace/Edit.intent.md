---
kind: file
covers: Edit.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useNamespaceStore.js
  - client/web/unify/src/sections/compose/components/Namespaces/NamespaceTranslator.vue
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Namespace Edit view

## Intention

Create a new namespace or configure an existing one: identity (name, slug,
labels), presentation (logo, subtitle, description), and behavior (enabled,
hide sidebar).

## UX capabilities

- Name is required; slug validated as a handle (letter first, then alphanumerics/underscores).
- Logo: opt-in toggle plus drop-zone upload to the namespace attachment endpoint; preview falls back to the global `ui.mainLogo` setting; clearable.
- Edit mode adds topbar tools (visit — disabled while namespace is disabled, configure → `admin.modules`, translator) and export/permissions buttons.
- Save gated by `canUpdateNamespace` on edit; clone creates a copy with blank slug; delete returns to the list.
- Unsaved-changes guard (deep compare against loaded snapshot); after create the user lands on the namespace's `pages` route so onboarding shows.

## Routes

`namespace.create` at `/namespaces/create` and `namespace.edit` at `/namespaces/edit/:slug` (both carry the namespace nav sidebar, `meta.sidebar: 'namespaces'`). Loads by slug-or-ID via `namespaceStore.getByUrlPart`; unknown slug bounces to `namespace.list`.

## When changing this

- After create, navigate to the leaf route `pages` (not the parent `namespace.view`) so the empty-path child renders the onboarding screen.
- `meta.logo` stores a relative attachment URL; preview must prefix `$ComposeAPI.baseURL` unless it is already absolute.
- Call `markSaved()` before programmatic navigation or the unsaved guard blocks it.
