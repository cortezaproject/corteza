---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/CSidebarNamespaceSwitcher.vue
  - client/web/unify/src/sections/compose/components/CSidebarNavigation.vue
touched-by:
  - client/web/unify/src/sections/compose/index.js
tests: []
---

# Compose sidebar

## Intention

The compose section's sidebar body, registered in the section contract
(`index.js` `sidebar:`). In the standalone app the namespace switcher lived in
the CSidebar `#header` slot and the nav tree in `#body`; here both share the
one sidebar body the shell provides.

## Map

- `ComposeSidebar.vue` — thin column composition: CSidebarNamespaceSwitcher on
  top, CSidebarNavigation (the namespace page tree) filling the rest.

## When changing this

- Keep it a pure composition — behavior belongs in the two components under
  `components/` (owned by their own doc).
- Routes with `meta.hideSidebar` (namespace chooser screens) suppress the
  sidebar entirely; this component need not handle that case.
