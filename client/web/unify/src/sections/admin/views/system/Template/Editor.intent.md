---
kind: file
covers: Editor.vue
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/admin/routes.js
  - client/web/unify/src/sections/admin/components/Template/CCodeEditor.vue
  - client/web/unify/src/sections/admin/components/Template/CTemplateToolbox.vue
  - client/web/unify/src/sections/admin/components/Template/CTemplatePreview.vue
  - lib/js
  - lib/vue
touched-by: []
tests: []
---

# Template Editor view

## Intention

Author one template: metadata, content (in a code editor with a partials
toolbox), and composition — header/footer partials for full templates.

## UX capabilities

- Basic info: name (`meta.short`, required), handle, content type
  (html/plain/pdf — drives editor syntax mode), language, partial toggle.
- Full (non-partial) templates pick header/footer from existing partials and
  get a live preview panel (edit mode only); partials get neither.
- Per-template permissions button; delete gated by `canDeleteTemplate`;
  unsaved-changes guard on leave.

## Routes

- `system.templates.create` → `/system/templates/new` (default type
  `text/html`); `system.templates.edit` → `/system/templates/:templateID`.
- Create redirects to `.edit`; back action → `system.templates`.

## When changing this

- Header/footer references live in `meta.headerTemplateID/footerTemplateID`;
  selecting "none" (`NoID`) must delete the key, not store the sentinel.
- The partial picker excludes the template itself — keep that filter or
  users can self-reference.
