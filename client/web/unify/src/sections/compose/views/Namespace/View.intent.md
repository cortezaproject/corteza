---
kind: file
covers: View.vue
backfilled: true
owner: fe
depends-on:
  - lib/vue/src/stores/useNamespaceStore.js
  - lib/vue/src/stores/useModuleStore.js
  - lib/vue/src/stores/usePageStore.js
  - lib/vue/src/stores/useChartStore.js
  - lib/vue/src/stores/usePageLayoutStore.js
  - client/web/unify/src/sections/compose/components/Record/RecordModal.vue
touched-by:
  - client/web/unify/src/sections/compose/routes.js
tests:
  - client/web/unify/src/sections/compose/views/Namespace/View.redirect.test.js
---

# Namespace View (context shell)

## Intention

Layout route for everything inside a namespace: resolves the `:slug` to a
namespace, prepares all namespace-scoped data, then renders the matched child
view. Every nested compose screen depends on the context established here.

## UX capabilities

- Resolves slug-or-ID; unknown namespaces toast and bounce to `root`.
- Disabled namespaces: admins may still use `admin.*` routes; users with update rights are redirected to the namespace editor; everyone else to `root`.
- Every one of those bounces replaces the history entry rather than pushing one, so Back cannot return to a URL this view refuses and be bounced again.
- Clears then preloads module/page/chart/page-layout stores in parallel on every namespace switch — children may assume the stores are populated.
- Provides `$namespace` and `$pageStore` for deeply nested consumers (e.g. record field editors), passes `:namespace` as a prop to the child view, and mounts the global `RecordModal`.
- Landing on the bare namespace route redirects to the first visible root-level non-record page (by weight).
- Loading state holds for a minimum duration to avoid spinner flash.

## Routes

`namespace.view` at `/namespace/:slug` (props: `slug`); parent of `pages`, `page`, `page.record`, and all `admin.*` compose routes.

## When changing this

- Store preloading must complete before the first-page redirect — the redirect reads `pageStore.set`.
- Re-runs fully on slug change; anything cached per-namespace must be cleared in `prepareNamespace`.
