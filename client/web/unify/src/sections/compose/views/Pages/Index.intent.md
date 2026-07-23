---
kind: file
covers: Index.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/usePageStore.js
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests: []
---

# Pages Index (namespace landing)

## Intention

The empty-path landing inside a namespace. If the namespace has a home page,
silently forward there; otherwise onboard the user toward building their first
module and page.

## UX capabilities

- Auto-redirects (replace) to the first visible, root-level, non-record page by weight as soon as one exists in the store — reactive, so it also fires when pages load in later.
- Onboarding screen otherwise: "create module" → `admin.modules` and "create page" → `admin.pages`, each shown only with the matching `canCreateModule` / `canCreatePage` flag.
- While a home page exists nothing is rendered (prevents onboarding flash before redirect).

## Routes

`pages` at the empty child path under `namespace.view` (`/namespace/:slug`). This is the canonical post-create landing target — `Namespace/Edit.vue` navigates here after creating a namespace.

## When changing this

- Home-page selection (visible + `selfID === NoID` + `moduleID === NoID`, lowest weight) is duplicated in `Namespace/View.vue`'s redirect — keep the two in sync.
- Redirect must use `router.replace` so Back does not bounce.
