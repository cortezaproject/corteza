---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on:
  - client/web/unify/src/sections/compose/components/CSidebarNamespaceSwitcher.vue
  - client/web/unify/src/sections/compose/components/CSidebarNavigation.vue
  - client/web/unify/src/sections/compose/components/CSidebarNamespaceNav.vue
touched-by:
  - client/web/unify/src/sections/compose/index.js
tests:
  - client/web/unify/src/sections/compose/sidebar/ComposeSidebar.scroll.test.js
---

# Compose sidebar

## Intention

The compose section's sidebar body, registered in the section contract
(`index.js` `sidebar:`). In the standalone app the namespace switcher lived in
the CSidebar `#header` slot and the nav tree in `#body`; here both share the
one sidebar body the shell provides.

## Map

- `ComposeSidebar.vue` — thin column composition in two shapes, chosen by
  `route.meta.sidebar`: inside a namespace, CSidebarNamespaceSwitcher on top
  and CSidebarNavigation (search + page tree + admin panel) filling the rest; on
  the namespace editor, CSidebarNamespaceNav alone — the page tree needs a
  namespace, which the editor has not entered.

## When changing this

- Keep it a pure composition — behavior belongs in the components under
  `components/` (owned by their own doc); this file only picks between them.
- The sidebar scrolls in exactly one place, and the namespace switcher is not
  in it. The shell (`CSidebar`) already scrolls whatever this body hands it, so
  this column has to fit that scroller exactly — a child claiming a full column
  on top of the switcher is taller than it, and the sidebar ends up with two
  scrollbars. Children take the height their siblings leave; the nav owns the
  one scroller and decides what goes inside it.
- The choice reads route meta, never route names, so a rename cannot silently
  swap sidebars. Routes with `meta.hideSidebar` (the namespace list) suppress
  the sidebar entirely; this component need not handle that case.
